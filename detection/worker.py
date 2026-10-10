"""Persistent, CPU-only YOLO26n ONNX worker. One JSON request/response per line.

FFmpeg samples the entire video; its exit code and coverage must both pass.
Only confirmed intervals qualify. Ambiguous observations prevent pruning.
"""

import argparse
import ast
import hashlib
import json
import math
import os
import subprocess
import sys
import tempfile
from datetime import datetime, timezone

VERSION = 1
# Keep weak candidates for the retention decision even when confirmation is
# stricter. Raising confidence must not silently turn them into empty footage.
UNCERTAINTY_CONFIDENCE = .4
LABELS = {
    0: "person", 1: "bicycle", 2: "car", 3: "motorcycle", 5: "bus", 7: "truck",
    14: "bird", 15: "cat", 16: "dog", 17: "horse", 18: "sheep", 19: "cow",
    20: "elephant", 21: "bear", 22: "zebra", 23: "giraffe",
}


def kind(class_id):
    if class_id == 0:
        return "person"
    if class_id in (1, 3):
        return "bike"
    if class_id in (2, 5, 7):
        return "moving_car"
    return "animal"


def center(box):
    return ((box[0] + box[2]) / 2, (box[1] + box[3]) / 2)


def distance(a, b):
    return math.dist(center(a), center(b))


def iou(a, b):
    intersection = max(0, min(a[2], b[2]) - max(a[0], b[0])) * max(0, min(a[3], b[3]) - max(a[1], b[1]))
    area = (a[2] - a[0]) * (a[3] - a[1]) + (b[2] - b[0]) * (b[3] - b[1]) - intersection
    return intersection / area if area > 0 else 0


def coherent_motion(rows, threshold):
    """Require translation across multiple steps, rather than a box-size jump."""
    first, last = rows[0][2], rows[-1][2]

    def shared_shift(a, b):
        return math.copysign(min(abs(a), abs(b)), a) if a * b > 0 else 0

    # Both edges must shift in the same direction. Growing or shrinking a
    # parked-car box can move its center without translating the vehicle.
    shift = (shared_shift(last[0] - first[0], last[2] - first[2]),
             shared_shift(last[1] - first[1], last[3] - first[3]))
    movement = math.hypot(*shift)
    if movement < threshold:
        return False
    deltas = [(center(b[2])[0] - center(a[2])[0], center(b[2])[1] - center(a[2])[1])
              for a, b in zip(rows, rows[1:])]
    path = sum(math.hypot(*delta) for delta in deltas)
    net = distance(first, last)
    if path == 0 or net / path < .8:
        return False
    # A single relocation or spike is still insufficient evidence of motion.
    return sum((dx * shift[0] + dy * shift[1]) / movement >= threshold * .1
               for dx, dy in deltas) >= 2


class Timeline:
    """Associate same-class boxes, then confirm presence and local car motion."""

    def __init__(self, fps, min_frames, merge_gap, car_movement, confidence=.65):
        self.step = 1 / fps
        self.min_frames = min_frames
        self.merge_gap = merge_gap
        self.confidence = confidence
        self.car_movement = car_movement
        self.tracks = []

    def observe(self, timestamp, detections):
        # Match each box and track at most once. A bounded displacement avoids
        # stitching unrelated cars across a scene into a fictitious moving car.
        candidates = []
        for track_index, track in enumerate(self.tracks):
            last = track[-1]
            if timestamp - last[0] > self.step * 1.1:
                continue
            for box_index, (class_id, box, confidence) in enumerate(detections):
                if class_id != last[1]:
                    continue
                overlap = iou(box, last[2])
                displacement = distance(box, last[2])
                ratios = [(box[axis + 2] - box[axis]) / (last[2][axis + 2] - last[2][axis])
                          for axis in (0, 1)]
                limit = max(.03, (last[2][2] - last[2][0]) * .75, (last[2][3] - last[2][1]) * .75)
                if all(.5 <= ratio <= 2 for ratio in ratios) and (overlap >= .1 or displacement <= limit):
                    candidates.append((1 - overlap + 3 * displacement, track_index, box_index))
        used_tracks, used_boxes = set(), set()
        for _, track_index, box_index in sorted(candidates):
            if track_index in used_tracks or box_index in used_boxes:
                continue
            class_id, box, confidence = detections[box_index]
            self.tracks[track_index].append((timestamp, class_id, box, confidence))
            used_tracks.add(track_index)
            used_boxes.add(box_index)
        for box_index, (class_id, box, confidence) in enumerate(detections):
            if box_index not in used_boxes:
                self.tracks.append([(timestamp, class_id, box, confidence)])

    def finish(self, duration):
        intervals = []
        uncertain = False
        confirmed_tracks = []
        for track in self.tracks:
            run = []
            for row in track:
                if row[3] >= self.confidence:
                    run.append(row)
                else:
                    uncertain = True
                    if run:
                        confirmed_tracks.append(run)
                        run = []
            if run:
                confirmed_tracks.append(run)
        for track in confirmed_tracks:
            class_id = track[0][1]
            required = max(3, self.min_frames) if kind(class_id) == "moving_car" else self.min_frames
            if len(track) < required:
                uncertain = True
                continue
            ranges = []
            if kind(class_id) != "moving_car":
                ranges.append((track[0][0], track[-1][0], max(row[3] for row in track)))
            else:
                # Confirm a coherent trajectory within two seconds. A single
                # jitter spike or changing box size cannot qualify a parked car.
                suspicious_motion = False
                for index, row in enumerate(track):
                    for start_index in range(max(0, index - math.ceil(2 / self.step)), index):
                        anchor = track[start_index]
                        if row[0] - anchor[0] > 2 or distance(row[2], anchor[2]) < self.car_movement:
                            continue
                        suspicious_motion = True
                        window = track[start_index:index + 1]
                        if len(window) >= 3 and coherent_motion(window, self.car_movement):
                            ranges.append((anchor[0], row[0], max(p[3] for p in window)))
                if not ranges and suspicious_motion:
                    uncertain = True
            for start, end, confidence in ranges:
                intervals.append({
                    "kind": kind(class_id), "labels": [LABELS[class_id]],
                    "start": round(max(0, start - self.step / 2), 3),
                    "end": round(min(duration, end + self.step), 3),
                    "confidence": round(confidence, 4),
                })
        # Merge only within a category; intervals for different categories can
        # overlap and are displayed in separate lanes in the player.
        merged = []
        for category in ("person", "animal", "bike", "moving_car"):
            rows = sorted((r for r in intervals if r["kind"] == category), key=lambda r: r["start"])
            for row in rows:
                if merged and merged[-1]["kind"] == category and row["start"] <= merged[-1]["end"] + self.merge_gap:
                    merged[-1]["end"] = max(merged[-1]["end"], row["end"])
                    merged[-1]["confidence"] = max(merged[-1]["confidence"], row["confidence"])
                    merged[-1]["labels"] = sorted(set(merged[-1]["labels"] + row["labels"]))
                else:
                    merged.append(row.copy())
        return sorted(merged, key=lambda r: (r["start"], r["kind"])), uncertain


class Detector:
    def __init__(self, model, threads):
        import numpy as np
        import onnxruntime as ort
        self.np = np
        options = ort.SessionOptions()
        options.intra_op_num_threads = threads
        options.inter_op_num_threads = 1
        options.execution_mode = ort.ExecutionMode.ORT_SEQUENTIAL
        # Let idle CPU cores sleep between frames instead of spinning.
        options.add_session_config_entry("session.intra_op.allow_spinning", "0")
        self.session = ort.InferenceSession(model, sess_options=options, providers=["CPUExecutionProvider"])
        metadata = self.session.get_modelmeta().custom_metadata_map
        names = ast.literal_eval(metadata.get("names", "{}"))
        if any(names.get(class_id) != label for class_id, label in LABELS.items()):
            raise ValueError("model must use the standard COCO object classes")
        if "YOLO26n" not in metadata.get("description", ""):
            raise ValueError("model must be an official YOLO26n detection export")
        model_input = self.session.get_inputs()[0]
        shape = model_input.shape
        if len(shape) != 4 or shape[:2] != [1, 3] or not all(isinstance(n, int) for n in shape):
            raise ValueError("export a fixed-size batch=1 YOLO26n ONNX model")
        self.height, self.width = shape[2:]
        if min(self.height, self.width) < 64 or max(self.height, self.width) > 1280:
            raise ValueError("ONNX input dimensions must be between 64 and 1280")
        self.input_name = model_input.name
        # Validate the end-to-end output before any recording can be pruned.
        self.predict(bytes(self.height * self.width * 3), .4)
        with open(model, "rb") as source:
            self.model_sha256 = hashlib.file_digest(source, "sha256").hexdigest()

    def predict(self, raw, confidence):
        np = self.np
        image = np.frombuffer(raw, dtype=np.uint8).reshape(self.height, self.width, 3)
        tensor = np.ascontiguousarray(image[:, :, ::-1].transpose(2, 0, 1)[None], dtype=np.float32) / 255
        rows = self.session.run(None, {self.input_name: tensor})[0]
        if rows.ndim != 3 or rows.shape[0] != 1 or rows.shape[2] != 6:
            raise ValueError("expected YOLO26n end-to-end output [1, N, 6]; export with nms=False")
        if not np.isfinite(rows).all():
            raise ValueError("non-finite model output")
        detections = []
        for x1, y1, x2, y2, score, class_id in rows[0]:
            if score < confidence or int(class_id) not in LABELS:
                continue
            box = [float(x1) / self.width, float(y1) / self.height, float(x2) / self.width, float(y2) / self.height]
            if box[2] > box[0] and box[3] > box[1]:
                detections.append((int(class_id), box, float(score)))
        return detections

    def analyze(self, request):
        path, fps = request["path"], request["fps"]
        probe = subprocess.run([
            request["ffprobe"], "-v", "error", "-select_streams", "v:0",
            "-show_entries", "stream=duration:format=duration", "-of", "json", path,
        ], capture_output=True, check=True, timeout=30)
        info = json.loads(probe.stdout)
        if not info.get("streams"):
            raise ValueError("no video stream")
        duration = float(info["streams"][0].get("duration", info.get("format", {}).get("duration", "nan")))
        if not math.isfinite(duration) or duration <= 0:
            raise ValueError("unknown video duration")
        timeline = Timeline(fps, request["min_frames"], request["merge_gap"], request["car_movement"], request["confidence"])
        observation_confidence = min(UNCERTAINTY_CONFIDENCE, request["confidence"])
        filters = (f"fps=fps={fps}:start_time=0:round=up:eof_action=pass,"
                   f"scale={self.width}:{self.height}:force_original_aspect_ratio=decrease,"
                   f"pad={self.width}:{self.height}:(ow-iw)/2:(oh-ih)/2:color=0x727272")
        frames, frame_size = 0, self.width * self.height * 3
        with tempfile.TemporaryFile() as errors:
            process = subprocess.Popen([
                request["ffmpeg"], "-nostdin", "-v", "error", "-xerror", "-err_detect", "explode",
                "-threads", "1", "-i", path, "-map", "0:v:0", "-an", "-sn", "-dn",
                "-vf", filters, "-filter_threads", "1", "-threads", "1",
                "-f", "rawvideo", "-pix_fmt", "bgr24", "pipe:1",
            ], stdout=subprocess.PIPE, stderr=errors)
            try:
                while True:
                    raw = process.stdout.read(frame_size)
                    if not raw:
                        break
                    if len(raw) != frame_size:
                        raise ValueError("incomplete sampled frame")
                    timestamp = frames / fps
                    timeline.observe(timestamp, self.predict(raw, observation_confidence))
                    frames += 1
                if process.wait(timeout=30) != 0:
                    errors.seek(0)
                    raise ValueError("video decode failed: " + errors.read(2000).decode(errors="replace"))
            finally:
                process.stdout.close()
                if process.poll() is None:
                    process.kill()
                process.wait()
        if frames == 0 or abs(frames / fps - duration) > 1.5 / fps:
            raise ValueError("incomplete video coverage; retaining recording")
        intervals, uncertain = timeline.finish(duration)
        return {
            "version": VERSION, "model": "yolo26n", "model_sha256": self.model_sha256,
            "policy": request["policy"], "complete": True, "duration": duration,
            "sample_fps": fps, "frames": frames, "uncertain": uncertain,
            "intervals": intervals, "analyzed_at": datetime.now(timezone.utc).isoformat(),
        }


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--model", required=True)
    parser.add_argument("--threads", type=int, default=2)
    args = parser.parse_args()
    model = os.path.expanduser(args.model)
    if not os.path.isabs(model) and os.environ.get("NEONVR_MODEL_DIR"):
        model = os.path.join(os.environ["NEONVR_MODEL_DIR"], model)
    detector = Detector(model, args.threads)
    print(json.dumps({"ready": True, "model_sha256": detector.model_sha256}), flush=True)
    for line in sys.stdin:
        try:
            result = {"result": detector.analyze(json.loads(line))}
        except Exception as error:
            result = {"error": str(error)}
        print(json.dumps(result, separators=(",", ":"), allow_nan=False), flush=True)


if __name__ == "__main__":
    main()
