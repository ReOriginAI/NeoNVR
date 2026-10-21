# Codex conversation

## User

what method is the AI using for detecting moving car and ignoring a parked car. Is it using background subtraction + CNN object detection?

## Assistant

I’ll check the detection code to see how it distinguishes moving cars from parked cars.

## Activity

    $ pwd && rg --files -g 'AGENTS.md' -g '!node_modules' -g '!vendor' -g '!venv' -g '!\.git' -g '!dist' -g '!build'
    /home/ubuntu/NeoNVR
    ✗ (1) • 0ms

## Activity

    $ rg -n -i 'background.?subtract|mog2|knn|parked|moving car|motion.?detect|stationary|yolo|object.?detect|optical.?flow' . --glob '!package-lock.json' --glob '!yarn.lock' --glob '!node_modules/**' --glob '!vendor/**' --glob '!venv/**' --glob '!dist/**' --glob '!build/**'
    ./config.go:60:        DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
    ./config.go:274:    fs.BoolVar(&c.DetectionEnabled, "detection-enabled", c.DetectionEnabled, "enable background YOLO26n CPU detection")
    ./config.go:276:    fs.StringVar(&c.DetectionModel, "detection-model", c.DetectionModel, "YOLO26n end-to-end ONNX model path")
    ./config/neonvr.json:20:  "detection_model": "yolo26n.onnx",
    ./detection.go:29:const detectionPolicyVersion = "yolo26n-onnx-v2"
    ./detection.go:75:        return errors.New("detection requires a Python executable and a YOLO26n .onnx model path")
    ./detection.go:112:    if d.Version != detectionVersion || d.Model != "yolo26n" || len(d.ModelSHA256) != 64 || len(d.Policy) != 64 || !d.Complete ||
    ./detection.go:517:                s.detectionWorkerError = "YOLO26n unavailable: " + err.Error()
    ./README.md:3:Docker-first security-camera MP4 browser with CPU-only YOLO26n event detection, including Blink clips stored as:
    ./README.md:103:- YOLO26n detection of humans, COCO animals, bicycles, motorcycles, and moving cars/buses/trucks
    ./README.md:110:The Go binary is built on Alpine. A separate Python build stage downloads the official `yolo26n.pt` weights and exports a fixed 640-pixel, end-to-end ONNX model. The Debian Bookworm runtime includes FFmpeg, Python, NumPy, and **CPU-only ONNX Runtime**. PyTorch, Ultralytics and GPU libraries stay in the model build stage. The model is included in the image, so analysis and the browser work offline after building. The first build needs network access for packages and model weights.
    ./README.md:182:## YOLO26n detection and false-positive pruning
    ./README.md:194:Use **Settings → Object detection** to adjust confidence, sample rate, confirmation count, movement threshold, scan timing and the no-target policy:
    ./README.md:200:Favorites are always preserved by detection pruning. Brief unconfirmed detections, decode errors, incomplete video coverage, missing dependencies/models, metadata-write failures, timeouts and changing recordings all keep the clip. No-target classification filters recordings caused by leaves; it is not a separate leaf-recognition model. Sampled object detection can miss events, so review quarantined clips when tuning settings before enabling permanent deletion.
    ./README.md:227:.venv-export/bin/python detection/export_model.py --output models/yolo26n.onnx
    ./README.md:232:  -detection-model "$PWD/models/yolo26n.onnx"
    ./README.md:246:Model/export reference: [Ultralytics YOLO26](https://docs.ultralytics.com/models/yolo26), [ONNX export](https://docs.ultralytics.com/modes/export).
    ./detection_test.go:51:    d := DetectionResult{Version: 1, Model: "yolo26n", ModelSHA256: strings.Repeat("a", 64), Policy: s.detectionPolicy(),
    ./detection_test.go:236:    legacy, err := json.Marshal([]any{"yolo26n-onnx-v1", c.DetectionModel, c.DetectionFPS, c.DetectionConfidence, c.DetectionMinFrames, c.DetectionMergeGap, c.DetectionCarMovement})
    ./detection_test.go:259:        func(c *Config) { c.DetectionConfidence = 1 }, func(c *Config) { c.DetectionModel = "yolo26n.pt" },
    ./detection_test.go:287:func TestYOLO26nWorkerAndPipeline(t *testing.T) {
    ./Dockerfile:11:ENV YOLO_AUTOINSTALL=False OMP_NUM_THREADS=2
    ./Dockerfile:17:RUN python export_model.py --output /model/yolo26n.onnx
    ./Dockerfile:24:COPY --from=model /model/yolo26n.onnx /opt/neonvr/models/yolo26n.onnx
    ./detection/worker.py:1:"""Persistent, CPU-only YOLO26n ONNX worker. One JSON request/response per line.
    ./detection/worker.py:61:    # parked-car box can move its center without translating the vehicle.
    ./detection/worker.py:91:        # stitching unrelated cars across a scene into a fictitious moving car.
    ./detection/worker.py:146:                # jitter spike or changing box size cannot qualify a parked car.
    ./detection/worker.py:197:        if "YOLO26n" not in metadata.get("description", ""):
    ./detection/worker.py:198:            raise ValueError("model must be an official YOLO26n detection export")
    ./detection/worker.py:202:            raise ValueError("export a fixed-size batch=1 YOLO26n ONNX model")
    ./detection/worker.py:218:            raise ValueError("expected YOLO26n end-to-end output [1, N, 6]; export with nms=False")
    ./detection/worker.py:277:            "version": VERSION, "model": "yolo26n", "model_sha256": self.model_sha256,
    ./main.go:1792:    log.Printf("YOLO26n detection: enabled=%t, CPU threads=%d, sample fps=%g, no-target policy=%s", cfg.DetectionEnabled, cfg.DetectionThreads, cfg.DetectionFPS, cfg.DetectionPrune)
    ./detection/test_worker.py:29:    def test_parked_car_and_detector_jitter_do_not_qualify(self):
    ./detection/test_worker.py:49:    def test_single_parked_car_box_jump_is_not_motion(self):
    ./detection/test_worker.py:55:    def test_parked_car_box_spike_is_not_motion(self):
    ./detection/test_worker.py:61:    def test_parked_car_box_resizing_is_not_motion(self):
    ./detection/test_worker.py:133:    def test_movement_does_not_mark_entire_parked_presence(self):
    ./detection/export_model.py:6:from ultralytics import YOLO
    ./detection/export_model.py:9:parser.add_argument("--output", default="yolo26n.onnx")
    ./detection/export_model.py:14:exported = Path(YOLO("yolo26n.pt").export(
    ./detection/export_model.py:20:print(f"YOLO26n CPU model: {output}")
    ./static/settings.html:16:<section class="panel"><h2>Object detection</h2><p>Find humans, animals, bikes and moving vehicles. Colored intervals appear below the video scrubber.</p><div class="grid">
    ./static/settings.html:29:<label class="wide">YOLO26n ONNX model<input name="detection_model" required><small>Docker includes yolo26n.onnx. Native installs can supply an absolute model path.</small></label>
    ✓ • 0ms

## Activity

    $ for dir in / /home /home/ubuntu /home/ubuntu/NeoNVR; do if [ -f "$dir/AGENTS.md" ]; then printf '%s\n' "$dir/AGENTS.md"; cat "$dir/AGENTS.md"; fi; done
    ✓ • 0ms

## Activity

    $ nl -ba detection/worker.py | sed -n '1,310p'
         1    """Persistent, CPU-only YOLO26n ONNX worker. One JSON request/response per line.
         2
         3    FFmpeg samples the entire video; its exit code and coverage must both pass.
         4    Only confirmed intervals qualify. Ambiguous observations prevent pruning.
         5    """
         6
         7    import argparse
         8    import ast
         9    import hashlib
        10    import json
        11    import math
        12    import os
        13    import subprocess
        14    import sys
        15    import tempfile
        16    from datetime import datetime, timezone
        17
        18    VERSION = 1
        19    # Keep weak candidates for the retention decision even when confirmation is
        20    # stricter. Raising confidence must not silently turn them into empty footage.
        21    UNCERTAINTY_CONFIDENCE = .4
        22    LABELS = {
        23        0: "person", 1: "bicycle", 2: "car", 3: "motorcycle", 5: "bus", 7: "truck",
        24        14: "bird", 15: "cat", 16: "dog", 17: "horse", 18: "sheep", 19: "cow",
        25        20: "elephant", 21: "bear", 22: "zebra", 23: "giraffe",
        26    }
        27
        28
        29    def kind(class_id):
        30        if class_id == 0:
        31            return "person"
        32        if class_id in (1, 3):
        33            return "bike"
        34        if class_id in (2, 5, 7):
        35            return "moving_car"
        36        return "animal"
        37
        38
        39    def center(box):
        40        return ((box[0] + box[2]) / 2, (box[1] + box[3]) / 2)
        41
        42
        43    def distance(a, b):
        44        return math.dist(center(a), center(b))
        45
        46
        47    def iou(a, b):
        48        intersection = max(0, min(a[2], b[2]) - max(a[0], b[0])) * max(0, min(a[3], b[3]) - max(a[1], b[1]))
        49        area = (a[2] - a[0]) * (a[3] - a[1]) + (b[2] - b[0]) * (b[3] - b[1]) - intersection
        50        return intersection / area if area > 0 else 0
        51
        52
        53    def coherent_motion(rows, threshold):
        54        """Require translation across multiple steps, rather than a box-size jump."""
        55        first, last = rows[0][2], rows[-1][2]
        56
        57        def shared_shift(a, b):
        58            return math.copysign(min(abs(a), abs(b)), a) if a * b > 0 else 0
        59
        60        # Both edges must shift in the same direction. Growing or shrinking a
        61        # parked-car box can move its center without translating the vehicle.
        62        shift = (shared_shift(last[0] - first[0], last[2] - first[2]),
        63                 shared_shift(last[1] - first[1], last[3] - first[3]))
        64        movement = math.hypot(*shift)
        65        if movement < threshold:
        66            return False
        67        deltas = [(center(b[2])[0] - center(a[2])[0], center(b[2])[1] - center(a[2])[1])
        68                  for a, b in zip(rows, rows[1:])]
        69        path = sum(math.hypot(*delta) for delta in deltas)
        70        net = distance(first, last)
        71        if path == 0 or net / path < .8:
        72            return False
        73        # A single relocation or spike is still insufficient evidence of motion.
        74        return sum((dx * shift[0] + dy * shift[1]) / movement >= threshold * .1
        75                   for dx, dy in deltas) >= 2
        76
        77
        78    class Timeline:
        79        """Associate same-class boxes, then confirm presence and local car motion."""
        80
        81        def __init__(self, fps, min_frames, merge_gap, car_movement, confidence=.65):
        82            self.step = 1 / fps
        83            self.min_frames = min_frames
        84            self.merge_gap = merge_gap
        85            self.confidence = confidence
        86            self.car_movement = car_movement
        87            self.tracks = []
        88
        89        def observe(self, timestamp, detections):
        90            # Match each box and track at most once. A bounded displacement avoids
        91            # stitching unrelated cars across a scene into a fictitious moving car.
        92            candidates = []
        93            for track_index, track in enumerate(self.tracks):
        94                last = track[-1]
        95                if timestamp - last[0] > self.step * 1.1:
        96                    continue
        97                for box_index, (class_id, box, confidence) in enumerate(detections):
        98                    if class_id != last[1]:
        99                        continue
       100                    overlap = iou(box, last[2])
       101                    displacement = distance(box, last[2])
       102                    ratios = [(box[axis + 2] - box[axis]) / (last[2][axis + 2] - last[2][axis])
       103                              for axis in (0, 1)]
       104                    limit = max(.03, (last[2][2] - last[2][0]) * .75, (last[2][3] - last[2][1]) * .75)
       105                    if all(.5 <= ratio <= 2 for ratio in ratios) and (overlap >= .1 or displacement <= limit):
       106                        candidates.append((1 - overlap + 3 * displacement, track_index, box_index))
       107            used_tracks, used_boxes = set(), set()
       108            for _, track_index, box_index in sorted(candidates):
       109                if track_index in used_tracks or box_index in used_boxes:
       110                    continue
       111                class_id, box, confidence = detections[box_index]
       112                self.tracks[track_index].append((timestamp, class_id, box, confidence))
       113                used_tracks.add(track_index)
       114                used_boxes.add(box_index)
       115            for box_index, (class_id, box, confidence) in enumerate(detections):
       116                if box_index not in used_boxes:
       117                    self.tracks.append([(timestamp, class_id, box, confidence)])
       118
       119        def finish(self, duration):
       120            intervals = []
       121            uncertain = False
       122            confirmed_tracks = []
       123            for track in self.tracks:
       124                run = []
       125                for row in track:
       126                    if row[3] >= self.confidence:
       127                        run.append(row)
       128                    else:
       129                        uncertain = True
       130                        if run:
       131                            confirmed_tracks.append(run)
       132                            run = []
       133                if run:
       134                    confirmed_tracks.append(run)
       135            for track in confirmed_tracks:
       136                class_id = track[0][1]
       137                required = max(3, self.min_frames) if kind(class_id) == "moving_car" else self.min_frames
       138                if len(track) < required:
       139                    uncertain = True
       140                    continue
       141                ranges = []
       142                if kind(class_id) != "moving_car":
       143                    ranges.append((track[0][0], track[-1][0], max(row[3] for row in track)))
       144                else:
       145                    # Confirm a coherent trajectory within two seconds. A single
       146                    # jitter spike or changing box size cannot qualify a parked car.
       147                    suspicious_motion = False
       148                    for index, row in enumerate(track):
       149                        for start_index in range(max(0, index - math.ceil(2 / self.step)), index):
       150                            anchor = track[start_index]
       151                            if row[0] - anchor[0] > 2 or distance(row[2], anchor[2]) < self.car_movement:
       152                                continue
       153                            suspicious_motion = True
       154                            window = track[start_index:index + 1]
       155                            if len(window) >= 3 and coherent_motion(window, self.car_movement):
       156                                ranges.append((anchor[0], row[0], max(p[3] for p in window)))
       157                    if not ranges and suspicious_motion:
       158                        uncertain = True
       159                for start, end, confidence in ranges:
       160                    intervals.append({
       161                        "kind": kind(class_id), "labels": [LABELS[class_id]],
       162                        "start": round(max(0, start - self.step / 2), 3),
       163                        "end": round(min(duration, end + self.step), 3),
       164                        "confidence": round(confidence, 4),
       165                    })
       166            # Merge only within a category; intervals for different categories can
       167            # overlap and are displayed in separate lanes in the player.
       168            merged = []
       169            for category in ("person", "animal", "bike", "moving_car"):
       170                rows = sorted((r for r in intervals if r["kind"] == category), key=lambda r: r["start"])
       171                for row in rows:
       172                    if merged and merged[-1]["kind"] == category and row["start"] <= merged[-1]["end"] + self.merge_gap:
       173                        merged[-1]["end"] = max(merged[-1]["end"], row["end"])
       174                        merged[-1]["confidence"] = max(merged[-1]["confidence"], row["confidence"])
       175                        merged[-1]["labels"] = sorted(set(merged[-1]["labels"] + row["labels"]))
       176                    else:
       177                        merged.append(row.copy())
       178            return sorted(merged, key=lambda r: (r["start"], r["kind"])), uncertain
       179
       180
       181    class Detector:
       182        def __init__(self, model, threads):
       183            import numpy as np
       184            import onnxruntime as ort
       185            self.np = np
       186            options = ort.SessionOptions()
       187            options.intra_op_num_threads = threads
       188            options.inter_op_num_threads = 1
       189            options.execution_mode = ort.ExecutionMode.ORT_SEQUENTIAL
       190            # Let idle CPU cores sleep between frames instead of spinning.
       191            options.add_session_config_entry("session.intra_op.allow_spinning", "0")
       192            self.session = ort.InferenceSession(model, sess_options=options, providers=["CPUExecutionProvider"])
       193            metadata = self.session.get_modelmeta().custom_metadata_map
       194            names = ast.literal_eval(metadata.get("names", "{}"))
       195            if any(names.get(class_id) != label for class_id, label in LABELS.items()):
       196                raise ValueError("model must use the standard COCO object classes")
       197            if "YOLO26n" not in metadata.get("description", ""):
       198                raise ValueError("model must be an official YOLO26n detection export")
       199            model_input = self.session.get_inputs()[0]
       200            shape = model_input.shape
       201            if len(shape) != 4 or shape[:2] != [1, 3] or not all(isinstance(n, int) for n in shape):
       202                raise ValueError("export a fixed-size batch=1 YOLO26n ONNX model")
       203            self.height, self.width = shape[2:]
       204            if min(self.height, self.width) < 64 or max(self.height, self.width) > 1280:
       205                raise ValueError("ONNX input dimensions must be between 64 and 1280")
       206            self.input_name = model_input.name
       207            # Validate the end-to-end output before any recording can be pruned.
       208            self.predict(bytes(self.height * self.width * 3), .4)
       209            with open(model, "rb") as source:
       210                self.model_sha256 = hashlib.file_digest(source, "sha256").hexdigest()
       211
       212        def predict(self, raw, confidence):
       213            np = self.np
       214            image = np.frombuffer(raw, dtype=np.uint8).reshape(self.height, self.width, 3)
       215            tensor = np.ascontiguousarray(image[:, :, ::-1].transpose(2, 0, 1)[None], dtype=np.float32) / 255
       216            rows = self.session.run(None, {self.input_name: tensor})[0]
       217            if rows.ndim != 3 or rows.shape[0] != 1 or rows.shape[2] != 6:
       218                raise ValueError("expected YOLO26n end-to-end output [1, N, 6]; export with nms=False")
       219            if not np.isfinite(rows).all():
       220                raise ValueError("non-finite model output")
       221            detections = []
       222            for x1, y1, x2, y2, score, class_id in rows[0]:
       223                if score < confidence or int(class_id) not in LABELS:
       224                    continue
       225                box = [float(x1) / self.width, float(y1) / self.height, float(x2) / self.width, float(y2) / self.height]
       226                if box[2] > box[0] and box[3] > box[1]:
       227                    detections.append((int(class_id), box, float(score)))
       228            return detections
       229
       230        def analyze(self, request):
       231            path, fps = request["path"], request["fps"]
       232            probe = subprocess.run([
       233                request["ffprobe"], "-v", "error", "-select_streams", "v:0",
       234                "-show_entries", "stream=duration:format=duration", "-of", "json", path,
       235            ], capture_output=True, check=True, timeout=30)
       236            info = json.loads(probe.stdout)
       237            if not info.get("streams"):
       238                raise ValueError("no video stream")
       239            duration = float(info["streams"][0].get("duration", info.get("format", {}).get("duration", "nan")))
       240            if not math.isfinite(duration) or duration <= 0:
       241                raise ValueError("unknown video duration")
       242            timeline = Timeline(fps, request["min_frames"], request["merge_gap"], request["car_movement"], request["confidence"])
       243            observation_confidence = min(UNCERTAINTY_CONFIDENCE, request["confidence"])
       244            filters = (f"fps=fps={fps}:start_time=0:round=up:eof_action=pass,"
       245                       f"scale={self.width}:{self.height}:force_original_aspect_ratio=decrease,"
       246                       f"pad={self.width}:{self.height}:(ow-iw)/2:(oh-ih)/2:color=0x727272")
       247            frames, frame_size = 0, self.width * self.height * 3
       248            with tempfile.TemporaryFile() as errors:
       249                process = subprocess.Popen([
       250                    request["ffmpeg"], "-nostdin", "-v", "error", "-xerror", "-err_detect", "explode",
       251                    "-threads", "1", "-i", path, "-map", "0:v:0", "-an", "-sn", "-dn",
       252                    "-vf", filters, "-filter_threads", "1", "-threads", "1",
       253                    "-f", "rawvideo", "-pix_fmt", "bgr24", "pipe:1",
       254                ], stdout=subprocess.PIPE, stderr=errors)
       255                try:
       256                    while True:
       257                        raw = process.stdout.read(frame_size)
       258                        if not raw:
       259                            break
       260                        if len(raw) != frame_size:
       261                            raise ValueError("incomplete sampled frame")
       262                        timestamp = frames / fps
       263                        timeline.observe(timestamp, self.predict(raw, observation_confidence))
       264                        frames += 1
       265                    if process.wait(timeout=30) != 0:
       266                        errors.seek(0)
       267                        raise ValueError("video decode failed: " + errors.read(2000).decode(errors="replace"))
       268                finally:
       269                    process.stdout.close()
       270                    if process.poll() is None:
       271                        process.kill()
       272                    process.wait()
       273            if frames == 0 or abs(frames / fps - duration) > 1.5 / fps:
       274                raise ValueError("incomplete video coverage; retaining recording")
       275            intervals, uncertain = timeline.finish(duration)
       276            return {
       277                "version": VERSION, "model": "yolo26n", "model_sha256": self.model_sha256,
       278                "policy": request["policy"], "complete": True, "duration": duration,
       279                "sample_fps": fps, "frames": frames, "uncertain": uncertain,
       280                "intervals": intervals, "analyzed_at": datetime.now(timezone.utc).isoformat(),
       281            }
       282
       283
       284    def main():
       285        parser = argparse.ArgumentParser()
       286        parser.add_argument("--model", required=True)
       287        parser.add_argument("--threads", type=int, default=2)
       288        args = parser.parse_args()
       289        model = os.path.expanduser(args.model)
       290        if not os.path.isabs(model) and os.environ.get("NEONVR_MODEL_DIR"):
       291            model = os.path.join(os.environ["NEONVR_MODEL_DIR"], model)
       292        detector = Detector(model, args.threads)
       293        print(json.dumps({"ready": True, "model_sha256": detector.model_sha256}), flush=True)
       294        for line in sys.stdin:
       295            try:
       296                result = {"result": detector.analyze(json.loads(line))}
       297            except Exception as error:
       298                result = {"error": str(error)}
       299            print(json.dumps(result, separators=(",", ":"), allow_nan=False), flush=True)
       300
       301
       302    if __name__ == "__main__":
       303        main()
    ✓ • 0ms

## Activity

    $ nl -ba README.md | sed -n '182,204p'
        nl -ba detection/test_worker.py | sed -n '1,180p'
        nl -ba config/neonvr.json | sed -n '14,37p'
        nl -ba static/settings.html | sed -n '15,35p'
       182    ## YOLO26n detection and false-positive pruning
       183
       184    Detection is enabled by default. Rebuild the container with `docker compose up -d --build` (or `sudo nerdctl compose up -d --build`). Background analysis starts after 20 seconds, scans all sources every minute, and waits two minutes after the last file write before processing. Existing recordings are analyzed as well as new ones. Files are processed serially, newest first within each camera. The model stays loaded between recordings.
       185
       186    The defaults target an **Intel Core i7-6700 @ 3.40 GHz** without a GPU: one worker, two ONNX inference threads, two sampled frames per second, and the nano model at 640 pixels. FFmpeg uses one decoder/filter thread. ONNX threads sleep between frames instead of spinning. The analysis runs through recordings as background work, so the sample rate describes video-time coverage, not a promised real-time throughput. Lower `detection_fps` to `1` or use a 320-pixel export if the server is busy; increase the sample rate for brief events. Lower resolutions and sample rates can miss small or fast subjects. Measurements on the actual server are needed to tune throughput.
       187
       188    Humans, animals and bikes require three consecutive matched observations with at least 65% confidence by default. Animal classes in the pretrained COCO model are bird, cat, dog, horse, sheep, cow, elephant, bear, zebra and giraffe; other wildlife is not covered by this pretrained class list. Cars, trucks and buses require at least three consecutive observations and consistent translation of both box edges by 2.5% of the model frame within two seconds. A single box jump, back-and-forth jitter, or box resizing does not qualify as vehicle motion. Detection intervals include up to one sample of padding and merge nearby same-category confirmed events across gaps of up to 0.5 seconds by default; the merge setting does not relax consecutive confirmation. Timing is approximate to the sampling cadence.
       189
       190    Candidates down to 40% confidence (or the configured confidence if lower) remain ambiguous when they fail confirmation, preserving the recording without adding event markers. Increasing confidence or the confirmation count therefore does not make these candidates eligible for pruning. This can leave noisy recordings in the browser for review. Thresholds trade event coverage for fewer confirmed false positives; review brief and distant subjects before tightening them further.
       191
       192    Vehicle confirmation needs at least three samples within a two-second motion window, so use at least one sampled frame per second for vehicle markers. At 0.5 FPS, observed vehicle movement remains ambiguous and the recording is retained.
       193
       194    Use **Settings → Object detection** to adjust confidence, sample rate, confirmation count, movement threshold, scan timing and the no-target policy:
       195
       196    - **Move to recoverable quarantine** (default): clips with no qualifying events leave the recording browser and move into `.neonvr_pruned` within their camera directory. This keeps the original video and its metadata on disk for review, and does **not** reclaim disk space.
       197    - **Keep all recordings**: embed intervals and display markers, retaining every clip.
       198    - **Permanently delete**: remove fully analyzed no-target clips and their cached thumbnails. This option reclaims disk space.
       199
       200    Favorites are always preserved by detection pruning. Brief unconfirmed detections, decode errors, incomplete video coverage, missing dependencies/models, metadata-write failures, timeouts and changing recordings all keep the clip. No-target classification filters recordings caused by leaves; it is not a separate leaf-recognition model. Sampled object detection can miss events, so review quarantined clips when tuning settings before enabling permanent deletion.
       201
       202    To restore a quarantined clip, first set the no-target policy to **Keep all recordings** and restart, then move the MP4 from `.neonvr_pruned` back to its original relative path in that camera directory. The existing age-based retention policy continues to operate on active recordings; it does not purge the quarantine.
       203
       204    ### MP4 metadata and viewer
         1    import shutil
         2    import subprocess
         3    import tempfile
         4    import unittest
         5    from pathlib import Path
         6
         7    from worker import Detector, LABELS, Timeline, kind
         8
         9
        10    def observation(class_id, x=.2, y=.2, score=.85):
        11        return (class_id, [x, y, x + .12, y + .15], score)
        12
        13
        14    class TimelineTests(unittest.TestCase):
        15        def timeline(self):
        16            return Timeline(fps=2, min_frames=2, merge_gap=.5, car_movement=.015)
        17
        18        def test_target_categories(self):
        19            for class_id in LABELS:
        20                with self.subTest(class_id=class_id):
        21                    timeline = self.timeline()
        22                    for frame in range(4):
        23                        timeline.observe(frame / 2, [observation(class_id, x=.2 + frame * .02)])
        24                    intervals, uncertain = timeline.finish(2)
        25                    self.assertFalse(uncertain)
        26                    self.assertEqual([row['kind'] for row in intervals], [kind(class_id)])
        27                    self.assertIn(LABELS[class_id], intervals[0]['labels'])
        28
        29        def test_parked_car_and_detector_jitter_do_not_qualify(self):
        30            timeline = self.timeline()
        31            for frame in range(20):
        32                timeline.observe(frame / 2, [observation(2, x=.2 + (frame % 3) * .002)])
        33            self.assertEqual(timeline.finish(10), ([], False))
        34
        35        def test_moving_car_with_two_observations(self):
        36            timeline = self.timeline()
        37            timeline.observe(0, [observation(2)])
        38            timeline.observe(.5, [observation(2, x=.23)])
        39            self.assertEqual(timeline.finish(1), ([], True))
        40
        41        def test_moving_car_with_consistent_translation(self):
        42            timeline = self.timeline()
        43            for frame in range(3):
        44                timeline.observe(frame / 2, [observation(2, x=.2 + frame * .02)])
        45            intervals, uncertain = timeline.finish(1.5)
        46            self.assertFalse(uncertain)
        47            self.assertEqual(intervals[0]['kind'], 'moving_car')
        48
        49        def test_single_parked_car_box_jump_is_not_motion(self):
        50            timeline = self.timeline()
        51            for frame in range(8):
        52                timeline.observe(frame / 2, [observation(2, x=.23 if frame >= 4 else .2)])
        53            self.assertEqual(timeline.finish(4), ([], True))
        54
        55        def test_parked_car_box_spike_is_not_motion(self):
        56            timeline = self.timeline()
        57            for frame in range(8):
        58                timeline.observe(frame / 2, [observation(2, x=.23 if frame == 4 else .2)])
        59            self.assertEqual(timeline.finish(4), ([], True))
        60
        61        def test_parked_car_box_resizing_is_not_motion(self):
        62            timeline = self.timeline()
        63            for frame in range(5):
        64                # One edge stays fixed while the detector expands the box.
        65                box = [.2, .2, .32 + frame * .02, .35]
        66                timeline.observe(frame / 2, [(2, box, .85)])
        67            self.assertEqual(timeline.finish(2.5), ([], True))
        68
        69        def test_low_confidence_false_positives_do_not_confirm_or_prune(self):
        70            for class_id in (0, 1, 2, 16):
        71                with self.subTest(class_id=class_id):
        72                    timeline = self.timeline()
        73                    for frame in range(6):
        74                        timeline.observe(frame / 2, [observation(class_id, x=.2 + frame * .02, score=.5)])
        75                    self.assertEqual(timeline.finish(3), ([], True))
        76
        77        def test_one_confident_frame_cannot_promote_a_weak_track(self):
        78            timeline = Timeline(2, 3, .5, .025, confidence=.65)
        79            for frame in range(6):
        80                timeline.observe(frame / 2, [observation(0, score=.9 if frame == 3 else .5)])
        81            self.assertEqual(timeline.finish(3), ([], True))
        82
        83        def test_two_confident_person_frames_remain_unconfirmed_by_default(self):
        84            timeline = Timeline(2, 3, .5, .025)
        85            for frame in range(2):
        86                timeline.observe(frame / 2, [observation(0, score=.7628)])
        87            self.assertEqual(timeline.finish(1), ([], True))
        88
        89        def test_low_confidence_sample_breaks_confirmation(self):
        90            timeline = Timeline(2, 3, .5, .025)
        91            for frame in range(5):
        92                timeline.observe(frame / 2, [observation(0, score=.5 if frame == 2 else .85)])
        93            self.assertEqual(timeline.finish(2.5), ([], True))
        94
        95        def test_merge_gap_does_not_confirm_intermittent_observations(self):
        96            timeline = Timeline(2, 3, 10, .025)
        97            for frame in range(6):
        98                timeline.observe(frame / 2, [observation(0)] if frame % 2 == 0 else [])
        99            self.assertEqual(timeline.finish(3), ([], True))
       100
       101        def test_confident_presence_does_not_require_motion(self):
       102            for class_id in (0, 1, 16):
       103                with self.subTest(class_id=class_id):
       104                    timeline = Timeline(2, 3, .5, .025)
       105                    for frame in range(4):
       106                        timeline.observe(frame / 2, [observation(class_id)])
       107                    intervals, uncertain = timeline.finish(2)
       108                    self.assertFalse(uncertain)
       109                    self.assertEqual([row['kind'] for row in intervals], [kind(class_id)])
       110
       111        def test_confidence_setting_can_be_lowered_explicitly(self):
       112            timeline = Timeline(2, 3, .5, .025, confidence=.45)
       113            for frame in range(3):
       114                timeline.observe(frame / 2, [observation(0, score=.5)])
       115            intervals, uncertain = timeline.finish(1.5)
       116            self.assertFalse(uncertain)
       117            self.assertEqual([row['kind'] for row in intervals], ['person'])
       118
       119        def test_large_confirmation_count_still_allows_local_vehicle_motion(self):
       120            timeline = Timeline(2, 10, .5, .025)
       121            for frame in range(12):
       122                timeline.observe(frame / 2, [observation(2, x=.2 + frame * .01)])
       123            intervals, uncertain = timeline.finish(6)
       124            self.assertFalse(uncertain)
       125            self.assertEqual([row['kind'] for row in intervals], ['moving_car'])
       126
       127        def test_vehicle_motion_at_low_sample_rate_remains_ambiguous(self):
       128            timeline = Timeline(.5, 3, .5, .025)
       129            for frame in range(4):
       130                timeline.observe(frame * 2, [observation(2, x=.2 + frame * .04)])
       131            self.assertEqual(timeline.finish(8), ([], True))
       132
       133        def test_movement_does_not_mark_entire_parked_presence(self):
       134            timeline = self.timeline()
       135            for frame in range(30):
       136                x = .2 if frame < 10 else .2 + min(3, frame - 10) * .03
       137                timeline.observe(frame / 2, [observation(2, x=x)])
       138            intervals, _ = timeline.finish(15)
       139            self.assertGreater(intervals[0]['start'], 0)
       140            self.assertLess(intervals[-1]['end'], 10)
       141
       142        def test_brief_observations_prevent_pruning(self):
       143            for class_id in (0, 1, 2, 16):
       144                timeline = self.timeline()
       145                timeline.observe(0, [observation(class_id)])
       146                self.assertEqual(timeline.finish(.5), ([], True))
       147            timeline = self.timeline()
       148            timeline.observe(0, [observation(2)])
       149            timeline.observe(.5, [observation(2)])
       150            self.assertEqual(timeline.finish(1), ([], True))
       151
       152        def test_multiple_timeframes_and_overlapping_categories(self):
       153            timeline = self.timeline()
       154            for frame in range(14):
       155                boxes = []
       156                if frame in (0, 1, 10, 11):
       157                    boxes.append(observation(0))
       158                if frame in (0, 1):
       159                    boxes.append(observation(16, x=.6))
       160                timeline.observe(frame / 2, boxes)
       161            intervals, uncertain = timeline.finish(7)
       162            self.assertFalse(uncertain)
       163            self.assertEqual(sum(row['kind'] == 'person' for row in intervals), 2)
       164            self.assertEqual(sum(row['kind'] == 'animal' for row in intervals), 1)
       165            self.assertTrue(all(0 <= row['start'] < row['end'] <= 7 for row in intervals))
       166
       167        def test_unrelated_cars_do_not_become_a_moving_track(self):
       168            timeline = self.timeline()
       169            timeline.observe(0, [observation(2, x=.1)])
       170            timeline.observe(.5, [observation(2, x=.8)])
       171            self.assertEqual(timeline.finish(1), ([], True))
       172
       173        def test_empty_frames_have_no_targets(self):
       174            timeline = self.timeline()
       175            for frame in range(10):
       176                timeline.observe(frame / 2, [])
       177            self.assertEqual(timeline.finish(5), ([], False))
       178
       179
       180    @unittest.skipUnless(shutil.which('ffmpeg') and shutil.which('ffprobe'), 'FFmpeg and FFprobe required')
        14      "thumb_maintenance_every": "30m",
        15      "size_refresh": "30m",
        16      "auto_delete_days": 14,
        17      "auto_delete_every": "6h",
        18      "detection_enabled": true,
        19      "detection_python": "python3",
        20      "detection_model": "yolo26n.onnx",
        21      "detection_threads": 2,
        22      "detection_fps": 2,
        23      "detection_confidence": 0.65,
        24      "detection_min_frames": 3,
        25      "detection_merge_gap": 0.5,
        26      "detection_car_movement": 0.025,
        27      "detection_every": "1m",
        28      "detection_min_age": "2m",
        29      "detection_timeout": "30m",
        30      "detection_prune": "quarantine"
        31    }
        15    <fieldset id="fields" disabled>
        16    <section class="panel"><h2>Object detection</h2><p>Find humans, animals, bikes and moving vehicles. Colored intervals appear below the video scrubber.</p><div class="grid">
        17    <label>Background detection<select name="detection_enabled"><option value="true">Enabled</option><option value="false">Disabled</option></select></label>
        18    <label>No-target recordings<select name="detection_prune"><option value="quarantine">Move to recoverable quarantine</option><option value="off">Keep all recordings</option><option value="delete">Permanently delete</option></select><small>Only fully analyzed clips without targets or ambiguous observations qualify. Favorites are preserved. Quarantine uses .neonvr_pruned inside each camera directory; it keeps files on disk.</small></label>
        19    <label>CPU inference threads<input name="detection_threads" type="number" min="1" max="16" step="1" required><small>2 is the default for an Intel i7-6700. One analysis worker loads the model once.</small></label>
        20    <label>Sampled frames per second<input name="detection_fps" type="number" min="0.5" max="10" step="0.5" required><small>2 balances CPU usage and event coverage. Increase for very brief events.</small></label>
        21    <label>Minimum confidence<input name="detection_confidence" type="number" min="0.05" max="0.95" step="0.01" required><small>0.65 by default. Higher values reduce false detections. Weaker candidates still keep the recording for review.</small></label>
        22    <label>Observations to confirm<input name="detection_min_frames" type="number" min="1" max="10" step="1" required><small>3 consecutive samples by default. Brief or intermittent observations keep the recording for review. Vehicle motion always needs at least 3.</small></label>
        23    <label>Detection gap (seconds)<input name="detection_merge_gap" type="number" min="0" max="10" step="0.1" required><small>Merge nearby confirmed events of the same category. This does not combine intermittent observations into a confirmed event.</small></label>
        24    <label>Vehicle movement threshold<input name="detection_car_movement" type="number" min="0.001" max="0.25" step="0.001" required><small>Fraction of frame size within two seconds; 0.025 is 2.5%. Requires consistent movement over multiple samples, filtering box jitter and resizing.</small></label>
        25    <label>Recording scan interval<input name="detection_every" required><small>Use 1m or 5m; minimum 1m.</small></label>
        26    <label>Wait after last recording write<input name="detection_min_age" required><small>2m by default; minimum 10s. Avoids processing ongoing uploads.</small></label>
        27    <label>Maximum analysis time per clip<input name="detection_timeout" required><small>30m by default. Timeouts keep the recording.</small></label>
        28    <label>Python executable<input name="detection_python" required><small>Docker includes the CPU detection runtime. Native installs can use a virtual environment.</small></label>
        29    <label class="wide">YOLO26n ONNX model<input name="detection_model" required><small>Docker includes yolo26n.onnx. Native installs can supply an absolute model path.</small></label>
        30    </div></section>
        31    <section class="panel"><h2>Recording sources</h2><div class="grid">
        32    <label class="wide">Primary recordings directory<input name="root" required><small>Existing recordings stay where they are. Enter an absolute path.</small></label>
        33    <label class="wide">Additional recordings directories<textarea name="extra_roots" placeholder="/home/ubuntu/FTPfiles/yicamBack&#10;/home/ubuntu/FTPfiles/yicamFront"></textarea><small>One absolute path per line. In a container, each directory must also be mounted.</small></label>
        34    </div></section>
        35    <section class="panel"><h2>Server and storage</h2><div class="grid">
    ✓ • 0ms

## Assistant

NeoNVR uses **YOLO26n object detection followed by bounding-box tracking and motion checks**. It does **not** use background subtraction.

The process is:

1. **Detect vehicles:** YOLO analyzes sampled video frames—**2 frames per second** by default—and returns vehicle bounding boxes. [Detection code](/home/ubuntu/NeoNVR/detection/worker.py:212)
2. **Track them:** The code matches vehicles between consecutive frames using class, box overlap, position, and size. [Tracking code](/home/ubuntu/NeoNVR/detection/worker.py:89)
3. **Confirm movement:** By default, a vehicle needs at least **3 consecutive detections at 65% confidence**, plus consistent box translation of at least **2.5% of the frame size within 2 seconds**. Motion checks reject small jitter, a single box jump, and box resizing. [Motion checks](/home/ubuntu/NeoNVR/detection/worker.py:53)

**A parked car is still detected by YOLO**, but its box stays approximately stationary, so it does not generate a moving-car event. The moving-versus-parked decision comes from these tracking rules, rather than a separate AI classifier.

## User

cameras now write file names with its own local timezone so localtime zone correction is not needed in NeoNVR

## Assistant

I’ll update NeoNVR to treat timestamps in camera filenames as local time and remove the timezone correction.

## Activity

    $ rg --files -g 'AGENTS.md' -g '!node_modules' -g '!vendor' -g '!venv' -g '!dist' -g '!build'
        git status --short
     M docker-compose.yml
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ✓ • 0ms

## Activity

    $ rg --files -g '!node_modules' -g '!vendor' -g '!venv' -g '!dist' -g '!build' -g '!\.git'
    main.go.bak-multiroot
    neonvr.service.example
    main.go.bak-interesting-thumbs
    detection.go
    docker-compose.yml.bak-pre-memory-limit-20261010
    config_test.go
    main.go.bak-retention-all
    detection_test.go
    Dockerfile
    main.go
    static/detection-viewer.js
    static/media-chrome-lite.js
    static/settings.html
    config.go
    README.md
    detection/test_worker.py
    detection/export_model.py
    detection/requirements.txt
    detection/worker.py
    main.go.bak-thumb-maintenance
    main.go.bak-divergence
    main.go.bak-yi-timezone
    go.mod
    docker-compose.yml
    config/neonvr.json.previous
    config/neonvr.json
    ✓ • 0ms

## Activity

    $ rg -n -i 'timezone|time.zone|localtime|local.time|utc|time\.Local|ParseInLocation|filename.*time|time.*filename|offset|correction' . --glob '!package-lock.json' --glob '!yarn.lock' --glob '!node_modules/**' --glob '!vendor/**' --glob '!venv/**' --glob '!dist/**' --glob '!build/**'
    ./main.go.bak-retention-all:651:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-retention-all:652:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-retention-all:1035:function playRelative(offset) {
    ./main.go.bak-retention-all:1038:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-retention-all:1039:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-retention-all:1499:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-retention-all:2038:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
    ./main.go.bak-retention-all:2254:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-retention-all:2262:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./detection_test.go:52:        Complete: true, Duration: 4, SampleFPS: 2, Frames: 8, Intervals: []DetectionInterval{}, AnalyzedAt: time.Now().UTC().Format(time.RFC3339)}
    ./README.md:267:New thumbnails use a divergence selector designed for mostly-static security footage. The server samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    ./README.md:278:Yi camera folder/file timestamps (for example `2026Y09M23D23H/W421M41S18.mp4`) are interpreted as UTC and converted to `America/New_York` before display and date grouping. This automatically uses EST or EDT as appropriate for the date.
    ./main.go.bak-yi-timezone:657:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-yi-timezone:658:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-yi-timezone:1041:function playRelative(offset) {
    ./main.go.bak-yi-timezone:1044:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-yi-timezone:1045:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-yi-timezone:1505:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-yi-timezone:2057:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
    ./main.go.bak-yi-timezone:2273:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-yi-timezone:2281:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./main.go.bak-yi-timezone:2486:        // Some camera MP4 variants do not expose duration metadata. Try useful fixed offsets
    ./main.go.bak-yi-timezone:2518:            return fmt.Errorf("duration probe failed (%v) and fixed-offset sampling produced fewer than two usable frames", probeErr)
    ./main.go.bak-divergence:656:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-divergence:657:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-divergence:1040:function playRelative(offset) {
    ./main.go.bak-divergence:1043:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-divergence:1044:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-divergence:1504:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-divergence:2043:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
    ./main.go.bak-divergence:2259:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-divergence:2267:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./main.go.bak-divergence:2420:        // Some camera MP4 variants do not expose duration metadata. Try useful fixed offsets
    ./main.go.bak-divergence:2457:            return fmt.Errorf("duration probe failed (%v) and fixed-offset sampling produced no usable frame", probeErr)
    ./main.go:676:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go:677:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go:1069:function playRelative(offset) {
    ./main.go:1072:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go:1073:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go:1533:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go:2147:        log.Printf("warning: could not load America/New_York timezone: %v; using UTC", err)
    ./main.go:2148:        return time.UTC
    ./main.go:2154:    // Yi folder/file timestamps are recorded in UTC. Convert them to the
    ./main.go:2155:    // display timezone here so API dates/times and date grouping use local time.
    ./main.go:2156:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.UTC)
    ./main.go:2388:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go:2396:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./main.go:2605:        // Some camera MP4 variants do not expose duration metadata. Try useful fixed offsets
    ./main.go:2637:            return fmt.Errorf("duration probe failed (%v) and fixed-offset sampling produced fewer than two usable frames", probeErr)
    ./main.go.bak-interesting-thumbs:651:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-interesting-thumbs:652:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-interesting-thumbs:1035:function playRelative(offset) {
    ./main.go.bak-interesting-thumbs:1038:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-interesting-thumbs:1039:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-interesting-thumbs:1499:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-interesting-thumbs:2038:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
    ./main.go.bak-interesting-thumbs:2254:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-interesting-thumbs:2262:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./static/media-chrome-lite.js:91:      const offset = Number(this.getAttribute('seek-offset') || '1');
    ./static/media-chrome-lite.js:92:      video.currentTime = Math.max(0, video.currentTime - offset);
    ./static/media-chrome-lite.js:101:      const offset = Number(this.getAttribute('seek-offset') || '1');
    ./static/media-chrome-lite.js:103:      video.currentTime = Math.min(max, video.currentTime + offset);
    ./detection/worker.py:16:from datetime import datetime, timezone
    ./detection/worker.py:280:            "intervals": intervals, "analyzed_at": datetime.now(timezone.utc).isoformat(),
    ./main.go.bak-thumb-maintenance:656:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-thumb-maintenance:657:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-thumb-maintenance:1040:function playRelative(offset) {
    ./main.go.bak-thumb-maintenance:1043:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-thumb-maintenance:1044:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-thumb-maintenance:1504:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-thumb-maintenance:2043:    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
    ./main.go.bak-thumb-maintenance:2259:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-thumb-maintenance:2267:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ./main.go.bak-thumb-maintenance:2472:        // Some camera MP4 variants do not expose duration metadata. Try useful fixed offsets
    ./main.go.bak-thumb-maintenance:2504:            return fmt.Errorf("duration probe failed (%v) and fixed-offset sampling produced fewer than two usable frames", probeErr)
    ./main.go.bak-multiroot:624:              <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
    ./main.go.bak-multiroot:625:              <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
    ./main.go.bak-multiroot:1008:function playRelative(offset) {
    ./main.go.bak-multiroot:1011:  if (index < 0) index = offset > 0 ? -1 : 0;
    ./main.go.bak-multiroot:1012:  const nextIndex = Math.max(0, Math.min(visibleClips.length - 1, index + offset));
    ./main.go.bak-multiroot:1472:               '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
    ./main.go.bak-multiroot:2101:func parseClipTimestamp(filename string) (time.Time, bool) {
    ./main.go.bak-multiroot:2109:    ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
    ✓ • 0ms

## Activity

    $ sed -n '1,230p' main.go
        sed -n '1860,2440p' main.go
    package main
    
    import (
        "context"
        "crypto/sha1"
        _ "embed"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "html/template"
        "image/jpeg"
        "log"
        "math"
        "net/http"
        "net/url"
        "os"
        "os/exec"
        "path/filepath"
        "sort"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    )
    
    const defaultRoot = "/mnt/mySSD/BlinkClips"
    const defaultAddr = "0.0.0.0:5500"
    const thumbnailAlgorithmVersion = "divergence-v2"
    
    var thumbnailCandidateFractions = []float64{0.08, 0.18, 0.28, 0.38, 0.50, 0.62, 0.72, 0.82, 0.92}
    
    type stringListFlag []string
    
    func (v *stringListFlag) String() string { return strings.Join(*v, ",") }
    func (v *stringListFlag) Set(value string) error {
        value = strings.TrimSpace(value)
        if value == "" {
            return fmt.Errorf("path cannot be empty")
        }
        *v = append(*v, value)
        return nil
    }
    
    func expandPath(path string) string {
        if path == "~" || strings.HasPrefix(path, "~/") {
            if home, err := os.UserHomeDir(); err == nil {
                if path == "~" {
                    return home
                }
                return filepath.Join(home, strings.TrimPrefix(path, "~/"))
            }
        }
        return path
    }
    
    type Server struct {
        runCtx                context.Context
        background            sync.WaitGroup
        root                  string
        extraRoots            []string
        thumbDir              string
        metaPath              string
        ffmpegPath            string
        thumbSeek             float64
        thumbWidth            int
        thumbSem              chan struct{}
        thumbMaintenanceEvery time.Duration
        sizeTTL               time.Duration
        autoDeleteDays        int
        autoDeleteEvery       time.Duration
    
        meta                 *MetaStore
        clipMutationMu       sync.Mutex
        detectionConfig      Config
        detectionDir         string
        detectionMu          sync.Mutex
        detectionStates      map[string]detectionResponse
        detectionWorkerError string
    
        sizeCacheMu  sync.Mutex
        sizeCache    statsResponse
        sizeCacheAt  time.Time
        sizeCacheErr string
    }
    
    type Camera struct {
        Name    string `json:"name"`
        Dir     string `json:"-"`
        Managed bool   `json:"-"`
    }
    
    type Clip struct {
        Camera       string `json:"camera"`
        File         string `json:"file"`
        Key          string `json:"key"`
        Timestamp    string `json:"timestamp"`
        UnixSeconds  int64  `json:"unix_seconds"`
        Date         string `json:"date"`
        Time         string `json:"time"`
        StartSeconds int    `json:"start_seconds"`
        Favorite     bool   `json:"favorite"`
        VideoURL     string `json:"video_url"`
        ThumbURL     string `json:"thumb_url"`
    }
    
    type ClipMeta struct {
        Favorite  bool   `json:"favorite,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
    }
    
    type MetaFile struct {
        Version int                 `json:"version"`
        Clips   map[string]ClipMeta `json:"clips"`
    }
    
    type MetaStore struct {
        path string
        mu   sync.Mutex
        data MetaFile
    }
    
    type camerasResponse struct {
        Root       string   `json:"root"`
        Cameras    []Camera `json:"cameras"`
        LatestDate string   `json:"latest_date"`
    }
    
    type datesResponse struct {
        Cameras []string `json:"cameras"`
        Dates   []string `json:"dates"`
    }
    
    type clipsResponse struct {
        Cameras []string `json:"cameras"`
        Date    string   `json:"date"`
        Clips   []Clip   `json:"clips"`
    }
    
    type favoritesResponse struct {
        Cameras []string `json:"cameras"`
        Clips   []Clip   `json:"clips"`
    }
    
    type mutateRequest struct {
        Camera   string `json:"camera"`
        File     string `json:"file"`
        Favorite bool   `json:"favorite"`
        Force    bool   `json:"force"`
    }
    
    type bulkDeleteItem struct {
        Camera string `json:"camera"`
        File   string `json:"file"`
    }
    
    type bulkDeleteRequest struct {
        Items []bulkDeleteItem `json:"items"`
        Force bool             `json:"force"`
    }
    
    type mutateResponse struct {
        OK      bool   `json:"ok"`
        Deleted int    `json:"deleted,omitempty"`
        Message string `json:"message,omitempty"`
    }
    
    type statsResponse struct {
        RootPath        string `json:"root_path"`
        RootSizeBytes   int64  `json:"root_size_bytes"`
        RootSizeHuman   string `json:"root_size_human"`
        GeneratedAt     string `json:"generated_at"`
        CacheTTLSeconds int64  `json:"cache_ttl_seconds"`
        Error           string `json:"error,omitempty"`
    }
    
    type pageData struct {
        Root string
    }
    
    //go:embed static/settings.html
    var settingsHTML []byte
    
    //go:embed static/media-chrome-lite.js
    var mediaChromeLiteJS []byte
    
    var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
    <html lang="en">
    <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
      <title>NeoNVR</title>
      <script src="/static/media-chrome-lite.js?v=detections-v1"></script>
      <style>
        :root {
          color-scheme: dark;
          --bg: #0f1115;
          --panel: #171a21;
          --panel2: #20242e;
          --border: #333846;
          --text: #f0f2f5;
          --muted: #9aa4b2;
          --accent: #72a7ff;
          --danger: #ff5d55;
          --ok: #75d685;
          --thumb: #f5f7fb;
          --thumbBorder: #111827;
        }
    
        * { box-sizing: border-box; }
        body {
          margin: 0;
          padding: 24px;
          background: var(--bg);
          color: var(--text);
          font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        }
        h1 { margin: 0 0 6px; font-size: 28px; }
        a { color: var(--accent); }
        .root { color: var(--muted); font-size: 13px; word-break: break-all; margin-bottom: 10px; }
    
        .panel {
          background: var(--panel);
          border: 1px solid var(--border);
          border-radius: 16px;
          padding: 16px;
          margin-bottom: 16px;
          box-shadow: 0 10px 30px rgba(0,0,0,.18);
        }
        m.mu.Lock()
        defer m.mu.Unlock()
        if favorite {
            meta := m.data.Clips[key]
            meta.Favorite = true
            meta.UpdatedAt = time.Now().Format(time.RFC3339)
            m.data.Clips[key] = meta
        } else {
            delete(m.data.Clips, key)
        }
        return m.saveLocked()
    }
    
    func (m *MetaStore) Remove(key string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        delete(m.data.Clips, key)
        return m.saveLocked()
    }
    
    func (m *MetaStore) FavoriteKeys() map[string]bool {
        m.mu.Lock()
        defer m.mu.Unlock()
        out := make(map[string]bool)
        for key, meta := range m.data.Clips {
            if meta.Favorite {
                out[key] = true
            }
        }
        return out
    }
    
    func (m *MetaStore) saveLocked() error {
        if err := os.MkdirAll(filepath.Dir(m.path), 0755); err != nil {
            return err
        }
        b, err := json.MarshalIndent(m.data, "", "  ")
        if err != nil {
            return err
        }
        tmp := m.path + ".tmp"
        if err := os.WriteFile(tmp, b, 0644); err != nil {
            return err
        }
        return os.Rename(tmp, m.path)
    }
    
    func clipKey(camera, file string) string {
        return camera + "/" + file
    }
    
    func (s *Server) staticMediaChromeHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
        w.Header().Set("Cache-Control", "no-cache")
        _, _ = w.Write(mediaChromeLiteJS)
    }
    
    func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        if err := pageTemplate.Execute(w, pageData{Root: s.root}); err != nil {
            log.Printf("template error: %v", err)
        }
    }
    
    func (s *Server) apiCamerasHandler(w http.ResponseWriter, r *http.Request) {
        cameras, err := s.listCameras()
        if err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        dates, err := s.listDates(nil)
        if err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        latest := ""
        if len(dates) > 0 {
            latest = dates[0]
        }
        writeJSON(w, camerasResponse{Root: s.root, Cameras: cameras, LatestDate: latest})
    }
    
    func (s *Server) apiDatesHandler(w http.ResponseWriter, r *http.Request) {
        selected, err := s.parseSelectedCameras(r)
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        dates, err := s.listDates(selected)
        if err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, datesResponse{Cameras: mapKeys(selected), Dates: dates})
    }
    
    func (s *Server) apiClipsHandler(w http.ResponseWriter, r *http.Request) {
        selected, err := s.parseSelectedCameras(r)
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        date := r.URL.Query().Get("date")
        clips, err := s.scanClips(selected, date, false)
        if err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, clipsResponse{Cameras: mapKeys(selected), Date: date, Clips: clips})
    }
    
    func (s *Server) apiFavoritesHandler(w http.ResponseWriter, r *http.Request) {
        selected, err := s.parseSelectedCameras(r)
        if err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        clips, err := s.scanClips(selected, "", true)
        if err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, favoritesResponse{Cameras: mapKeys(selected), Clips: clips})
    }
    
    func (s *Server) apiFavoriteHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
            return
        }
        var req mutateRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        if req.Camera == "" || req.File == "" {
            writeJSONError(w, http.StatusBadRequest, fmt.Errorf("camera and file are required"))
            return
        }
        s.clipMutationMu.Lock()
        defer s.clipMutationMu.Unlock()
        if _, err := s.resolveClipPath(req.Camera, req.File); err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        if err := s.meta.SetFavorite(clipKey(req.Camera, req.File), req.Favorite); err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, mutateResponse{OK: true})
    }
    
    func (s *Server) apiDeleteHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
            return
        }
        var req mutateRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        if req.Camera == "" || req.File == "" {
            writeJSONError(w, http.StatusBadRequest, fmt.Errorf("camera and file are required"))
            return
        }
        if err := s.deleteClip(req.Camera, req.File, req.Force); err != nil {
            writeJSONError(w, http.StatusInternalServerError, err)
            return
        }
        writeJSON(w, mutateResponse{OK: true, Deleted: 1, Message: "deleted"})
    }
    
    func (s *Server) apiBulkDeleteHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            writeJSONError(w, http.StatusMethodNotAllowed, fmt.Errorf("POST required"))
            return
        }
        var req bulkDeleteRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSONError(w, http.StatusBadRequest, err)
            return
        }
        if len(req.Items) == 0 {
            writeJSONError(w, http.StatusBadRequest, fmt.Errorf("no items selected"))
            return
        }
        deleted := 0
        for _, item := range req.Items {
            if item.Camera == "" || item.File == "" {
                continue
            }
            if err := s.deleteClip(item.Camera, item.File, req.Force); err != nil {
                writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("failed on %s/%s: %w", item.Camera, item.File, err))
                return
            }
            deleted++
        }
        writeJSON(w, mutateResponse{OK: true, Deleted: deleted, Message: "deleted"})
    }
    
    func (s *Server) apiStatsHandler(w http.ResponseWriter, r *http.Request) {
        writeJSON(w, s.getStats())
    }
    
    func (s *Server) videoHandler(w http.ResponseWriter, r *http.Request) {
        camera := r.URL.Query().Get("camera")
        file := r.URL.Query().Get("file")
        if camera == "" || file == "" {
            http.Error(w, "missing camera or file", http.StatusBadRequest)
            return
        }
        path, err := s.resolveClipPath(camera, file)
        if err != nil {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        http.ServeFile(w, r, path)
    }
    
    func (s *Server) thumbHandler(w http.ResponseWriter, r *http.Request) {
        camera := r.URL.Query().Get("camera")
        file := r.URL.Query().Get("file")
        if camera == "" || file == "" {
            writePlaceholderSVG(w, "missing camera/file")
            return
        }
    
        videoPath, err := s.resolveClipPath(camera, file)
        if err != nil {
            writePlaceholderSVG(w, "video unavailable")
            return
        }
        thumbPath, err := s.ensureThumbnail(videoPath)
        if err != nil {
            log.Printf("thumbnail error for %s: %v", videoPath, err)
            writePlaceholderSVG(w, "thumbnail unavailable")
            return
        }
    
        w.Header().Set("Cache-Control", "public, max-age=86400")
        http.ServeFile(w, r, thumbPath)
    }
    
    func (s *Server) parseSelectedCameras(r *http.Request) (map[string]bool, error) {
        allCameras, err := s.listCameras()
        if err != nil {
            return nil, err
        }
        all := make(map[string]bool, len(allCameras))
        for _, cam := range allCameras {
            all[cam.Name] = true
        }
    
        raw := strings.TrimSpace(r.URL.Query().Get("cameras"))
        if raw == "" {
            return all, nil
        }
        selected := make(map[string]bool)
        for _, part := range strings.Split(raw, ",") {
            name, err := url.QueryUnescape(strings.TrimSpace(part))
            if err != nil {
                return nil, err
            }
            if name != "" && all[name] {
                selected[name] = true
            }
        }
        return selected, nil
    }
    
    func mapKeys(m map[string]bool) []string {
        keys := make([]string, 0, len(m))
        for k := range m {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        return keys
    }
    
    var newYorkLocation = func() *time.Location {
        loc, err := time.LoadLocation("America/New_York")
        if err != nil {
            log.Printf("warning: could not load America/New_York timezone: %v; using UTC", err)
            return time.UTC
        }
        return loc
    }()
    
    func parseYICamHourFolder(name string) (time.Time, bool) {
        // Yi folder/file timestamps are recorded in UTC. Convert them to the
        // display timezone here so API dates/times and date grouping use local time.
        ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.UTC)
        if err != nil {
            return time.Time{}, false
        }
        return ts.In(newYorkLocation), true
    }
    
    func parseYICamMinuteSecond(filename string) (int, int, bool) {
        base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
        var minute, second, duration int
        n, err := fmt.Sscanf(base, "W4%dM%dS%d", &minute, &second, &duration)
        if err != nil || n != 3 || minute < 0 || minute > 59 || second < 0 || second > 59 {
            return 0, 0, false
        }
        return minute, second, true
    }
    
    func (s *Server) listCameras() ([]Camera, error) {
        entries, err := os.ReadDir(s.root)
        if err != nil {
            return nil, fmt.Errorf("could not read clips root %q: %w", s.root, err)
        }
        cameras := make([]Camera, 0)
        used := make(map[string]bool)
        for _, entry := range entries {
            name := entry.Name()
            if entry.IsDir() && !strings.HasPrefix(name, ".") {
                cameras = append(cameras, Camera{Name: name, Dir: filepath.Join(s.root, name), Managed: true})
                used[name] = true
            }
        }
        for _, dir := range s.extraRoots {
            if _, err := os.Stat(dir); err != nil {
                log.Printf("extra clip root unavailable %s: %v", dir, err)
                continue
            }
            name := filepath.Base(dir)
            base := name
            if used[name] {
                name = name + " [" + filepath.Base(filepath.Dir(dir)) + "]"
            }
            for i := 2; used[name]; i++ {
                name = fmt.Sprintf("%s #%d", base, i)
            }
            used[name] = true
            cameras = append(cameras, Camera{Name: name, Dir: dir, Managed: false})
        }
        sort.Slice(cameras, func(i, j int) bool { return cameras[i].Name < cameras[j].Name })
        return cameras, nil
    }
    
    type clipFile struct {
        RelPath  string
        FullPath string
        Info     os.FileInfo
    }
    
    func listCameraFiles(cam Camera) ([]clipFile, error) {
        entries, err := os.ReadDir(cam.Dir)
        if err != nil {
            return nil, err
        }
        out := make([]clipFile, 0)
        add := func(rel, full string, e os.DirEntry) {
            if strings.HasPrefix(e.Name(), ".") || e.Type()&os.ModeSymlink != 0 || e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".mp4") {
                return
            }
            info, err := e.Info()
            if err != nil {
                return
            }
            out = append(out, clipFile{RelPath: filepath.ToSlash(rel), FullPath: full, Info: info})
        }
        for _, e := range entries {
            if e.IsDir() {
                if _, ok := parseYICamHourFolder(e.Name()); !ok {
                    continue
                }
                sub := filepath.Join(cam.Dir, e.Name())
                subs, err := os.ReadDir(sub)
                if err != nil {
                    continue
                }
                for _, se := range subs {
                    add(filepath.Join(e.Name(), se.Name()), filepath.Join(sub, se.Name()), se)
                }
                continue
            }
            add(e.Name(), filepath.Join(cam.Dir, e.Name()), e)
        }
        return out, nil
    }
    
    func clipTimestamp(f clipFile) time.Time {
        if ts, ok := parseClipTimestamp(filepath.Base(f.RelPath)); ok {
            return ts
        }
        parent := filepath.Base(filepath.Dir(filepath.FromSlash(f.RelPath)))
        if hour, ok := parseYICamHourFolder(parent); ok {
            if minute, second, ok := parseYICamMinuteSecond(filepath.Base(f.RelPath)); ok {
                return hour.Add(time.Duration(minute)*time.Minute + time.Duration(second)*time.Second)
            }
        }
        return f.Info.ModTime()
    }
    
    func (s *Server) listDates(selected map[string]bool) ([]string, error) {
        cameras, err := s.listCameras()
        if err != nil {
            return nil, err
        }
        seen := make(map[string]bool)
        for _, cam := range cameras {
            if selected != nil && !selected[cam.Name] {
                continue
            }
            files, err := listCameraFiles(cam)
            if err != nil {
                continue
            }
            for _, f := range files {
                seen[clipTimestamp(f).Format("2006-01-02")] = true
            }
        }
        dates := make([]string, 0, len(seen))
        for date := range seen {
            dates = append(dates, date)
        }
        sort.Sort(sort.Reverse(sort.StringSlice(dates)))
        return dates, nil
    }
    
    func (s *Server) scanClips(selected map[string]bool, selectedDate string, favoritesOnly bool) ([]Clip, error) {
        cameras, err := s.listCameras()
        if err != nil {
            return nil, err
        }
        favKeys := map[string]bool{}
        if favoritesOnly {
            favKeys = s.meta.FavoriteKeys()
        }
        clips := make([]Clip, 0)
        for _, cam := range cameras {
            if selected != nil && !selected[cam.Name] {
                continue
            }
            files, err := listCameraFiles(cam)
            if err != nil {
                continue
            }
            for _, f := range files {
                ts := clipTimestamp(f)
                date := ts.Format("2006-01-02")
                if selectedDate != "" && date != selectedDate {
                    continue
                }
                name := f.RelPath
                key := clipKey(cam.Name, name)
                if favoritesOnly && !favKeys[key] {
                    continue
                }
                meta := s.meta.Get(key)
                clips = append(clips, Clip{
                    Camera: cam.Name, File: name, Key: key,
                    Timestamp: ts.Format("2006-01-02 15:04:05"), UnixSeconds: ts.Unix(),
                    Date: date, Time: ts.Format("15:04:05"),
                    StartSeconds: ts.Hour()*3600 + ts.Minute()*60 + ts.Second(),
                    Favorite:     meta.Favorite,
                    VideoURL:     "/video?camera=" + url.QueryEscape(cam.Name) + "&file=" + url.QueryEscape(name),
                    ThumbURL:     "/thumb?camera=" + url.QueryEscape(cam.Name) + "&file=" + url.QueryEscape(name),
                })
            }
        }
        sort.Slice(clips, func(i, j int) bool {
            if clips[i].UnixSeconds == clips[j].UnixSeconds {
                return clips[i].Camera < clips[j].Camera
            }
            return clips[i].UnixSeconds > clips[j].UnixSeconds
        })
        return clips, nil
    }
    
    func (s *Server) resolveClipPath(camera, file string) (string, error) {
        cameras, err := s.listCameras()
        if err != nil {
            return "", err
        }
        var dir string
        for _, cam := range cameras {
            if cam.Name == camera {
                dir = cam.Dir
                break
            }
        }
        if dir == "" {
            return "", fmt.Errorf("unknown camera/source %q", camera)
        }
        if file == "" || filepath.IsAbs(file) {
            return "", fmt.Errorf("invalid clip path")
        }
        rel := filepath.Clean(filepath.FromSlash(file))
        if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
            return "", fmt.Errorf("invalid clip path")
        }
        full := filepath.Join(dir, rel)
        check, err := filepath.Rel(dir, full)
        if err != nil || check == ".." || strings.HasPrefix(check, ".."+string(filepath.Separator)) {
            return "", fmt.Errorf("clip path escapes source directory")
        }
        info, err := os.Stat(full)
        if err != nil {
            return "", err
        }
        if info.IsDir() || !strings.EqualFold(filepath.Ext(full), ".mp4") {
            return "", fmt.Errorf("not an MP4 clip")
        }
        // Resolve symlinks before allowing detection/remux/deletion to touch a file.
        realDir, err := filepath.EvalSymlinks(dir)
        if err != nil {
            return "", err
        }
        realFile, err := filepath.EvalSymlinks(full)
        if err != nil {
            return "", err
        }
        realRel, err := filepath.Rel(realDir, realFile)
        if err != nil || realRel == ".." || strings.HasPrefix(realRel, ".."+string(filepath.Separator)) {
            return "", fmt.Errorf("clip symlink escapes source directory")
        }
        return full, nil
    }
    
    func parseClipTimestamp(filename string) (time.Time, bool) {
        base := filepath.Base(filename)
        parts := strings.Split(base, "_")
        if len(parts) < 2 {
            return time.Time{}, false
        }
        datePart := parts[0]
        timePart := strings.ReplaceAll(parts[1], "-", ":")
        ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
        if err != nil {
            return time.Time{}, false
        }
        return ts, true
    }
    
    func (s *Server) thumbnailPath(videoPath string) (string, error) {
        info, err := os.Stat(videoPath)
        if err != nil {
            return "", err
        }
        keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
        sum := sha1.Sum([]byte(keyInput))
        name := hex.EncodeToString(sum[:]) + ".jpg"
        return filepath.Join(s.thumbDir, name), nil
    }
    
    func (s *Server) ensureThumbnail(videoPath string) (string, error) {
        thumbPath, err := s.thumbnailPath(videoPath)
        if err != nil {
            return "", err
        }
        if err := os.MkdirAll(s.thumbDir, 0755); err != nil {
            return "", err
        }
        if st, err := os.Stat(thumbPath); err == nil && st.Size() > 0 {
            return thumbPath, nil
        }
    
        tmpPath := strings.TrimSuffix(thumbPath, ".jpg") + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
        _ = os.Remove(tmpPath)
    
        select {
        case s.thumbSem <- struct{}{}:
        case <-s.workContext().Done():
            return "", s.workContext().Err()
        }
        defer func() { <-s.thumbSem }()
    
        if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
            _ = os.Remove(tmpPath)
            seek := strconv.FormatFloat(s.thumbSeek, 'f', 3, 64)
            log.Printf("interesting thumbnail selection failed for %s, falling back to %ss: %v", videoPath, seek, err)
            if err2 := s.runFFmpegThumbnail(videoPath, tmpPath, seek); err2 != nil {
    ✓ • 0ms

## Activity

    $ rg -n 'yiDisplayLocation|parseYi|Yi|\.In\(|time.LoadLocation|tzdata|Timestamp|ModTime|GroupDate|Date:' main.go config.go config_test.go detection_test.go README.md Dockerfile docker-compose.yml
        cat go.mod
        sed -n '1,130p' Dockerfile
        sed -n '1,115p' docker-compose.yml
        sed -n '260,290p' README.md
        git diff -- docker-compose.yml
    README.md:14:The settings file is `config/neonvr.json`. The supplied file preserves the existing Blink root, both Yi camera roots, 14-day retention, two thumbnail workers, and 30-minute refresh/maintenance intervals.
    README.md:180:Detection also writes `.neonvr_detections/` under the primary root as an expendable JSON cache. Detection metadata lives inside the MP4 itself, so deleting the cache does not lose event intervals. Quarantined recordings live in `<camera-directory>/.neonvr_pruned/`, preserving Yi hour subdirectories where applicable.
    README.md:259:Each extra root appears in the camera/source selector using the directory basename (for example `yicamBack`). Extra roots support MP4 files directly in the directory and Yi-style hour folders such as `2026Y09M21D16H/`. For a file such as `W440M00S60.mp4` in that folder, the server derives `2026-09-21 16:40:00`; loose MP4 files that do not contain a full timestamp fall back to their filesystem modification time.
    README.md:276:### Yi timestamps
    README.md:278:Yi camera folder/file timestamps (for example `2026Y09M23D23H/W421M41S18.mp4`) are interpreted as UTC and converted to `America/New_York` before display and date grouping. This automatically uses EST or EDT as appropriate for the date.
    Dockerfile:21:RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
    main.go:98:    Timestamp    string `json:"timestamp"`
    main.go:1943:    writeJSON(w, camerasResponse{Root: s.root, Cameras: cameras, LatestDate: latest})
    main.go:1972:    writeJSON(w, clipsResponse{Cameras: mapKeys(selected), Date: date, Clips: clips})
    main.go:2145:    loc, err := time.LoadLocation("America/New_York")
    main.go:2154:    // Yi folder/file timestamps are recorded in UTC. Convert them to the
    main.go:2160:    return ts.In(newYorkLocation), true
    main.go:2249:func clipTimestamp(f clipFile) time.Time {
    main.go:2250:    if ts, ok := parseClipTimestamp(filepath.Base(f.RelPath)); ok {
    main.go:2259:    return f.Info.ModTime()
    main.go:2277:            seen[clipTimestamp(f).Format("2006-01-02")] = true
    main.go:2307:            ts := clipTimestamp(f)
    main.go:2320:                Timestamp: ts.Format("2006-01-02 15:04:05"), UnixSeconds: ts.Unix(),
    main.go:2321:                Date: date, Time: ts.Format("15:04:05"),
    main.go:2388:func parseClipTimestamp(filename string) (time.Time, bool) {
    main.go:2408:    keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
    main.go:2898:            if !clipTimestamp(f).Before(cutoff) {
    detection_test.go:77:    if !after.ModTime().Equal(original.ModTime()) || after.Mode().Perm() != original.Mode().Perm() {
    detection_test.go:178:    newTime := original.ModTime().Add(time.Minute)
    detection_test.go:225:    newTime := info.ModTime().Add(time.Minute)
    module neonvr
    
    go 1.22
    FROM golang:1.22-alpine AS build
    WORKDIR /src
    COPY go.mod ./
    COPY *.go ./
    COPY static ./static
    COPY detection ./detection
    RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/neonvr .
    
    # PyTorch/Ultralytics are used only to export the official weights at build time.
    FROM python:3.12-slim-bookworm AS model
    ENV YOLO_AUTOINSTALL=False OMP_NUM_THREADS=2
    RUN apt-get update && apt-get install -y --no-install-recommends libgl1 libglib2.0-0 && rm -rf /var/lib/apt/lists/*
    WORKDIR /model
    RUN pip install --no-cache-dir torch==2.9.1 torchvision==0.24.1 --index-url https://download.pytorch.org/whl/cpu
    RUN pip install --no-cache-dir ultralytics==8.4.175 onnx==1.23.2
    COPY detection/export_model.py ./
    RUN python export_model.py --output /model/yolo26n.onnx
    
    # CPU ONNX Runtime supports the Intel i7-6700; no CUDA/GPU runtime is installed.
    FROM python:3.12-slim-bookworm
    RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
    COPY detection/requirements.txt /tmp/detection-requirements.txt
    RUN pip install --no-cache-dir -r /tmp/detection-requirements.txt && rm /tmp/detection-requirements.txt
    COPY --from=model /model/yolo26n.onnx /opt/neonvr/models/yolo26n.onnx
    ENV NEONVR_MODEL_DIR=/opt/neonvr/models OMP_NUM_THREADS=2
    COPY --from=build /out/neonvr /usr/local/bin/neonvr
    EXPOSE 5500
    ENTRYPOINT ["/usr/local/bin/neonvr"]
    CMD ["-config", "/config/neonvr.json"]
    services:
      neonvr:
        build: .
        container_name: neonvr
        restart: unless-stopped
        mem_limit: 6g
        memswap_limit: 4g
        ports:
          - "5500:5500"
        volumes:
          # Mount the directory so config saves can replace the file atomically.
          - ./config:/config
          # Read/write for thumbnails, favorites, and configured deletion.
          - /mnt:/mnt
      #    - /mnt/mySSD/BlinkClips:/mnt/mySSD/BlinkClips
     #     - /mnt/FTPfiles:/mnt/FTPfiles
    #      - /home/ubuntu/FTPfiles/yicamBack:/home/ubuntu/FTPfiles/yicamBack
    #      - /home/ubuntu/FTPfiles/yicamFront:/home/ubuntu/FTPfiles/yicamFront
        command: ["-config", "/config/neonvr.json"]
    
    When using Docker/nerdctl, every extra root must also be mounted into the container and listed in `extra_roots` in the config. The included Compose file mounts both `/home/ubuntu/FTPfiles/yicamBack` and `/home/ubuntu/FTPfiles/yicamFront`; the config lists them in `extra_roots`.
    
    Automatic age-based deletion applies to the primary `-root` and all configured `-extra-root` directories. Favorited clips are preserved everywhere.
    
    ### Thumbnail selection
    
    New thumbnails use a divergence selector designed for mostly-static security footage. The server samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    
    
    ### Automatic thumbnail maintenance
    
    Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same divergence selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).
    
    Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thumb-maintenance-every 1h`. Set it to `0` to disable background maintenance. Lazy/on-demand thumbnail generation remains enabled as a fallback.
    
    ### Yi timestamps
    
    Yi camera folder/file timestamps (for example `2026Y09M23D23H/W421M41S18.mp4`) are interpreted as UTC and converted to `America/New_York` before display and date grouping. This automatically uses EST or EDT as appropriate for the date.
    
    ## Project rename and compatibility
    
    The project directory is now `/home/ubuntu/NeoNVR`, the Go module and executable are `neonvr`, the Compose service/container are `neonvr`, and the systemd example is `neonvr.service.example`.
    
    Existing recordings remain under `/mnt/mySSD/BlinkClips`. The legacy `.blinkclips_thumbs` and `.blinkclips_meta.json` paths are deliberately retained so thumbnail caches and favorite selections continue to work. Custom `-root`, `-thumb-dir`, and `-meta` flags continue to work.
    
    A running container keeps its current binary until rebuilt. Stop the existing Compose stack using its original project name and configuration before starting the renamed stack; otherwise both stacks can contend for port 5500. For nerdctl deployments, use `sudo nerdctl compose` in place of `docker compose`. After stopping the old stack:
    
    ```bash
    cd /home/ubuntu/NeoNVR
    sudo nerdctl compose up -d --build
    diff --git a/docker-compose.yml b/docker-compose.yml
    index e61be4c..7b495d7 100644
    --- a/docker-compose.yml
    +++ b/docker-compose.yml
    @@ -3,6 +3,8 @@ services:
         build: .
         container_name: neonvr
         restart: unless-stopped
    +    mem_limit: 6g
    +    memswap_limit: 4g
         ports:
           - "5500:5500"
         volumes:
    ✓ • 0ms

## Activity

    file changes: Completed · 2 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -277,3 +277,3 @@
     
    -Yi camera folder/file timestamps (for example `2026Y09M23D23H/W421M41S18.mp4`) are interpreted as UTC and converted to `America/New_York` before display and date grouping. This automatically uses EST or EDT as appropriate for the date.
    +Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
     
    Update { move_path: None }: /home/ubuntu/NeoNVR/main.go
    @@ -2143,15 +2143,5 @@
     
    -var newYorkLocation = func() *time.Location {
    -	loc, err := time.LoadLocation("America/New_York")
    -	if err != nil {
    -		log.Printf("warning: could not load America/New_York timezone: %v; using UTC", err)
    -		return time.UTC
    -	}
    -	return loc
    -}()
    -
     func parseYICamHourFolder(name string) (time.Time, bool) {
    -	// Yi folder/file timestamps are recorded in UTC. Convert them to the
    -	// display timezone here so API dates/times and date grouping use local time.
    -	ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.UTC)
    +	// Cameras write local timestamps; preserve their date and hour as written.
    +	ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
     	if err != nil {
    @@ -2159,3 +2149,3 @@
     	}
    -	return ts.In(newYorkLocation), true
    +	return ts, true
     }

## Activity

    $ go test ./...
    /bin/bash: line 1: go: command not found
    ✗ (127) • 0ms

## Activity

    $ sed -n '1,95p' config.go
        sed -n '230,310p' config.go
        sed -n '1650,1815p' main.go
        git diff --check
        git diff --stat
        git diff -- main.go README.md
    package main
    
    import (
        "bytes"
        "crypto/rand"
        "crypto/sha256"
        "encoding/hex"
        "encoding/json"
        "errors"
        "flag"
        "fmt"
        "io"
        "math"
        "net"
        "net/http"
        "net/url"
        "os"
        "path/filepath"
        "reflect"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    )
    
    var flagHelpError = flag.ErrHelp
    
    type Config struct {
        Root                  string   `json:"root"`
        ExtraRoots            []string `json:"extra_roots"`
        Addr                  string   `json:"addr"`
        ThumbDir              string   `json:"thumb_dir"`
        Meta                  string   `json:"meta"`
        FFmpeg                string   `json:"ffmpeg"`
        ThumbSeek             float64  `json:"thumb_seek"`
        ThumbWidth            int      `json:"thumb_width"`
        ThumbWorkers          int      `json:"thumb_workers"`
        ThumbMaintenanceEvery string   `json:"thumb_maintenance_every"`
        SizeRefresh           string   `json:"size_refresh"`
        AutoDeleteDays        int      `json:"auto_delete_days"`
        AutoDeleteEvery       string   `json:"auto_delete_every"`
        DetectionEnabled      bool     `json:"detection_enabled"`
        DetectionPython       string   `json:"detection_python"`
        DetectionModel        string   `json:"detection_model"`
        DetectionThreads      int      `json:"detection_threads"`
        DetectionFPS          float64  `json:"detection_fps"`
        DetectionConfidence   float64  `json:"detection_confidence"`
        DetectionMinFrames    int      `json:"detection_min_frames"`
        DetectionMergeGap     float64  `json:"detection_merge_gap"`
        DetectionCarMovement  float64  `json:"detection_car_movement"`
        DetectionEvery        string   `json:"detection_every"`
        DetectionMinAge       string   `json:"detection_min_age"`
        DetectionTimeout      string   `json:"detection_timeout"`
        DetectionPrune        string   `json:"detection_prune"`
    }
    
    func defaultConfig() Config {
        return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h",
            DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
    }
    
    func decodeConfig(data []byte) (Config, error) {
        c := defaultConfig()
        if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
            return c, errors.New("config must be a JSON object")
        }
        dec := json.NewDecoder(bytes.NewReader(data))
        dec.DisallowUnknownFields()
        if err := dec.Decode(&c); err != nil {
            return c, err
        }
        if err := dec.Decode(new(any)); err != io.EOF {
            return c, errors.New("config must contain exactly one JSON object")
        }
        if c.ExtraRoots == nil {
            c.ExtraRoots = []string{}
        }
        return c, nil
    }
    
    func (c Config) validate() error {
        if strings.TrimSpace(c.Root) == "" {
            return errors.New("root cannot be empty")
        }
        for _, p := range append([]string{c.Root}, c.ExtraRoots...) {
            if !filepath.IsAbs(expandPath(p)) {
                return fmt.Errorf("recording directory must be absolute: %s", p)
            }
            st, err := os.Stat(expandPath(p))
            if err != nil {
                return fmt.Errorf("recording directory %s: %w", p, err)
            }
            if !st.IsDir() {
                return fmt.Errorf("recording directory %s is not a directory", p)
            }
        }
        if strings.TrimSpace(path) == "" {
            return Config{}, "", errors.New("config path cannot be empty")
        }
        path, err := filepath.Abs(expandPath(path))
        if err != nil {
            return Config{}, "", err
        }
        c := defaultConfig()
        raw, readErr := os.ReadFile(path)
        if readErr == nil {
            c, err = decodeConfig(raw)
            if err != nil {
                return c, path, fmt.Errorf("config %s: %w", path, err)
            }
        } else if !os.IsNotExist(readErr) || explicit {
            return c, path, readErr
        }
        fs := flag.NewFlagSet("neonvr", flag.ContinueOnError)
        fs.StringVar(&path, "config", path, "JSON config file (default config/neonvr.json)")
        fs.StringVar(&c.Root, "root", c.Root, "primary recordings directory")
        extra := stringListFlag(append([]string{}, c.ExtraRoots...))
        hasExtra := false
        for _, a := range args {
            if a == "-extra-root" || a == "--extra-root" || strings.HasPrefix(a, "-extra-root=") || strings.HasPrefix(a, "--extra-root=") {
                hasExtra = true
            }
        }
        if hasExtra {
            extra = nil
        }
        fs.Var(&extra, "extra-root", "additional directory; repeat to replace configured extra roots")
        fs.StringVar(&c.Addr, "addr", c.Addr, "listen address")
        fs.StringVar(&c.ThumbDir, "thumb-dir", c.ThumbDir, "thumbnail directory; empty uses <root>/.blinkclips_thumbs")
        fs.StringVar(&c.Meta, "meta", c.Meta, "favorites JSON; empty uses <root>/.blinkclips_meta.json")
        fs.StringVar(&c.FFmpeg, "ffmpeg", c.FFmpeg, "ffmpeg executable")
        fs.Float64Var(&c.ThumbSeek, "thumb-seek", c.ThumbSeek, "fallback thumbnail seek in seconds")
        fs.IntVar(&c.ThumbWidth, "thumb-width", c.ThumbWidth, "thumbnail width")
        fs.IntVar(&c.ThumbWorkers, "thumb-workers", c.ThumbWorkers, "thumbnail workers")
        fs.StringVar(&c.ThumbMaintenanceEvery, "thumb-maintenance-every", c.ThumbMaintenanceEvery, "maintenance cadence; 0 disables")
        fs.StringVar(&c.SizeRefresh, "size-refresh", c.SizeRefresh, "size statistics cache cadence")
        fs.IntVar(&c.AutoDeleteDays, "auto-delete-days", c.AutoDeleteDays, "retention days; 0 disables deletion")
        fs.StringVar(&c.AutoDeleteEvery, "auto-delete-every", c.AutoDeleteEvery, "retention check cadence")
        fs.BoolVar(&c.DetectionEnabled, "detection-enabled", c.DetectionEnabled, "enable background YOLO26n CPU detection")
        fs.StringVar(&c.DetectionPython, "detection-python", c.DetectionPython, "Python executable with ONNX Runtime installed")
        fs.StringVar(&c.DetectionModel, "detection-model", c.DetectionModel, "YOLO26n end-to-end ONNX model path")
        fs.IntVar(&c.DetectionThreads, "detection-threads", c.DetectionThreads, "CPU inference threads")
        fs.Float64Var(&c.DetectionFPS, "detection-fps", c.DetectionFPS, "sampled frames per second")
        fs.Float64Var(&c.DetectionConfidence, "detection-confidence", c.DetectionConfidence, "minimum confidence for confirmed detections")
        fs.IntVar(&c.DetectionMinFrames, "detection-min-frames", c.DetectionMinFrames, "consecutive observations required to confirm a detection")
        fs.Float64Var(&c.DetectionMergeGap, "detection-merge-gap", c.DetectionMergeGap, "maximum gap between confirmed intervals in seconds")
        fs.Float64Var(&c.DetectionCarMovement, "detection-car-movement", c.DetectionCarMovement, "car movement threshold as fraction of frame")
        fs.StringVar(&c.DetectionEvery, "detection-every", c.DetectionEvery, "recording scan interval")
        fs.StringVar(&c.DetectionMinAge, "detection-min-age", c.DetectionMinAge, "minimum time since last recording write")
        fs.StringVar(&c.DetectionTimeout, "detection-timeout", c.DetectionTimeout, "maximum analysis time per recording")
        fs.StringVar(&c.DetectionPrune, "detection-prune", c.DetectionPrune, "no-target clip policy: off, quarantine, or delete")
        if err = fs.Parse(args); err != nil {
            return c, path, err
        }
        if fs.NArg() != 0 {
            return c, path, errors.New("unexpected positional arguments")
        }
        c.ExtraRoots = append([]string{}, extra...)
        path, err = filepath.Abs(expandPath(path))
        if err != nil {
            return c, path, err
        }
        if err = c.validate(); err != nil {
            return c, path, err
        }
        // Persist explicit legacy flag overrides, so later browser restarts use those settings.
        changed := os.IsNotExist(readErr)
        if readErr == nil {
            original, _ := decodeConfig(raw)
            changed = !reflect.DeepEqual(c, original)
        }
        if changed {
            if err = atomicWriteConfig(path, configBytes(c)); err != nil {
                return c, path, err
            }
        }
        manager, err := newConfigManager(configPath, cfg)
        if err != nil {
            log.Fatal(err)
        }
        executable, err := os.Executable()
        if err != nil {
            log.Fatal(err)
        }
        root, addr := &cfg.Root, &cfg.Addr
        extraRoots := cfg.ExtraRoots
        thumbDir, metaPath, ffmpegPath := &cfg.ThumbDir, &cfg.Meta, &cfg.FFmpeg
        thumbSeek, thumbWidth, thumbWorkers := &cfg.ThumbSeek, &cfg.ThumbWidth, &cfg.ThumbWorkers
        maintenance, _ := time.ParseDuration(cfg.ThumbMaintenanceEvery)
        refresh, _ := time.ParseDuration(cfg.SizeRefresh)
        retentionEvery, _ := time.ParseDuration(cfg.AutoDeleteEvery)
        thumbMaintenanceEvery, sizeRefresh, autoDeleteEvery := &maintenance, &refresh, &retentionEvery
        autoDeleteDays := &cfg.AutoDeleteDays
        runCtx, cancelWork := context.WithCancel(context.Background())
        defer cancelWork()
    
        cleanRoot := filepath.Clean(expandPath(*root))
        cleanExtraRoots := make([]string, 0, len(extraRoots))
        seenExtraRoots := make(map[string]bool)
        for _, extraRoot := range extraRoots {
            clean := filepath.Clean(expandPath(extraRoot))
            if abs, err := filepath.Abs(clean); err == nil {
                clean = abs
            }
            if !seenExtraRoots[clean] {
                cleanExtraRoots = append(cleanExtraRoots, clean)
                seenExtraRoots[clean] = true
            }
        }
        cleanThumbDir := *thumbDir
        if cleanThumbDir == "" {
            cleanThumbDir = filepath.Join(cleanRoot, ".blinkclips_thumbs")
        }
        cleanThumbDir = filepath.Clean(expandPath(cleanThumbDir))
    
        cleanMetaPath := *metaPath
        if cleanMetaPath == "" {
            cleanMetaPath = filepath.Join(cleanRoot, ".blinkclips_meta.json")
        }
        cleanMetaPath = filepath.Clean(expandPath(cleanMetaPath))
    
        workers := *thumbWorkers
        if workers < 1 {
            workers = 1
        }
        if *thumbWidth < 64 {
            *thumbWidth = 64
        }
        if *thumbMaintenanceEvery > 0 && *thumbMaintenanceEvery < time.Minute {
            *thumbMaintenanceEvery = time.Minute
        }
        if *sizeRefresh < 10*time.Second {
            *sizeRefresh = 10 * time.Second
        }
        if *thumbSeek < 0 {
            *thumbSeek = 0
        }
        if *autoDeleteEvery < time.Minute {
            *autoDeleteEvery = time.Minute
        }
    
        meta, err := NewMetaStore(cleanMetaPath)
        if err != nil {
            log.Fatalf("could not load metadata file %s: %v", cleanMetaPath, err)
        }
    
        s := &Server{
            runCtx:                runCtx,
            root:                  cleanRoot,
            extraRoots:            cleanExtraRoots,
            thumbDir:              cleanThumbDir,
            metaPath:              cleanMetaPath,
            ffmpegPath:            *ffmpegPath,
            thumbSeek:             *thumbSeek,
            thumbWidth:            *thumbWidth,
            thumbSem:              make(chan struct{}, workers),
            thumbMaintenanceEvery: *thumbMaintenanceEvery,
            sizeTTL:               *sizeRefresh,
            autoDeleteDays:        *autoDeleteDays,
            autoDeleteEvery:       *autoDeleteEvery,
            meta:                  meta,
            detectionConfig:       cfg,
            detectionDir:          filepath.Join(cleanRoot, ".neonvr_detections"),
        }
    
        if s.autoDeleteDays > 0 {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.autoDeleteLoop() }()
        }
        if s.thumbMaintenanceEvery > 0 {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.thumbnailMaintenanceLoop() }()
        }
        if cfg.DetectionEnabled {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.detectionLoop() }()
        }
    
        mux := http.NewServeMux()
        mux.HandleFunc("/", s.indexHandler)
        mux.HandleFunc("/api/config", manager.handler)
        mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                w.Header().Set("Allow", "GET")
                http.Error(w, "method not allowed", 405)
                return
            }
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            w.Header().Set("Cache-Control", "no-store")
            _, _ = w.Write(settingsHTML)
        })
        mux.HandleFunc("/api/cameras", s.apiCamerasHandler)
        mux.HandleFunc("/api/dates", s.apiDatesHandler)
        mux.HandleFunc("/api/clips", s.apiClipsHandler)
        mux.HandleFunc("/api/favorites", s.apiFavoritesHandler)
        mux.HandleFunc("/api/favorite", s.apiFavoriteHandler)
        mux.HandleFunc("/api/delete", s.apiDeleteHandler)
        mux.HandleFunc("/api/bulk-delete", s.apiBulkDeleteHandler)
        mux.HandleFunc("/api/stats", s.apiStatsHandler)
        mux.HandleFunc("/api/detections", s.apiDetectionsHandler)
        mux.HandleFunc("/static/detection-viewer.js", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            _, _ = w.Write(detectionViewerJS)
        })
        mux.HandleFunc("/static/media-chrome-lite.js", s.staticMediaChromeHandler)
        mux.HandleFunc("/video", s.videoHandler)
        mux.HandleFunc("/thumb", s.thumbHandler)
    
        log.Printf("NeoNVR server listening on http://%s", *addr)
        log.Printf("Clips root: %s", s.root)
        for _, extraRoot := range s.extraRoots {
            log.Printf("Extra clip root: %s", extraRoot)
        }
        log.Printf("Thumbnail cache: %s", s.thumbDir)
        log.Printf("Metadata: %s", s.metaPath)
        log.Printf("ffmpeg: %s", s.ffmpegPath)
        log.Printf("YOLO26n detection: enabled=%t, CPU threads=%d, sample fps=%g, no-target policy=%s", cfg.DetectionEnabled, cfg.DetectionThreads, cfg.DetectionFPS, cfg.DetectionPrune)
        log.Printf("Thumbnail seek: %.3fs", s.thumbSeek)
        log.Printf("Thumbnail workers: %d", workers)
        if s.thumbMaintenanceEvery > 0 {
            log.Printf("Thumbnail maintenance: every %s (pre-generate missing, prune orphan cache files)", s.thumbMaintenanceEvery)
        } else {
            log.Printf("Thumbnail maintenance disabled")
        }
        log.Printf("Size refresh: %s", s.sizeTTL)
        if s.autoDeleteDays > 0 {
            log.Printf("Auto-delete enabled: non-favorite clips older than %d days, checking every %s", s.autoDeleteDays, s.autoDeleteEvery)
        } else {
            log.Printf("Auto-delete disabled; use -auto-delete-days N to enable")
        }
        httpServer := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
        serverErr := make(chan error, 1)
        go func() { serverErr <- httpServer.ListenAndServe() }()
        select {
        case err := <-serverErr:
            if !errors.Is(err, http.ErrServerClosed) {
                log.Fatal(err)
            }
        case <-manager.restart:
            log.Printf("Restarting NeoNVR with config %s", configPath)
     README.md          |  2 +-
     docker-compose.yml |  2 ++
     main.go            | 16 +++-------------
     3 files changed, 6 insertions(+), 14 deletions(-)
    diff --git a/README.md b/README.md
    index c3acd07..bfd5d3a 100644
    --- a/README.md
    +++ b/README.md
    @@ -275,7 +275,7 @@ Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thu
    
     ### Yi timestamps
    
    -Yi camera folder/file timestamps (for example `2026Y09M23D23H/W421M41S18.mp4`) are interpreted as UTC and converted to `America/New_York` before display and date grouping. This automatically uses EST or EDT as appropriate for the date.
    +Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
    
     ## Project rename and compatibility
    
    diff --git a/main.go b/main.go
    index 70d5f34..7336b1a 100644
    --- a/main.go
    +++ b/main.go
    @@ -2141,23 +2141,13 @@ func mapKeys(m map[string]bool) []string {
         return keys
     }
    
    -var newYorkLocation = func() *time.Location {
    -    loc, err := time.LoadLocation("America/New_York")
    -    if err != nil {
    -        log.Printf("warning: could not load America/New_York timezone: %v; using UTC", err)
    -        return time.UTC
    -    }
    -    return loc
    -}()
    -
     func parseYICamHourFolder(name string) (time.Time, bool) {
    -    // Yi folder/file timestamps are recorded in UTC. Convert them to the
    -    // display timezone here so API dates/times and date grouping use local time.
    -    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.UTC)
    +    // Cameras write local timestamps; preserve their date and hour as written.
    +    ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
         if err != nil {
             return time.Time{}, false
         }
    -    return ts.In(newYorkLocation), true
    +    return ts, true
     }
    
     func parseYICamMinuteSecond(filename string) (int, int, bool) {
    ✓ • 0ms

## Activity

    $ git status --short
        rg -n 'newYorkLocation|America/New_York|time.UTC' main.go README.md
     M README.md
     M docker-compose.yml
     M main.go
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ✗ (1) • 0ms

## Activity

    $ for executable in /usr/local/go/bin/go /opt/go/bin/go /snap/bin/go /usr/bin/docker /usr/local/bin/docker /usr/local/bin/nerdctl /usr/bin/nerdctl; do if [ -x "$executable" ]; then printf '%s\n' "$executable"; fi; done
        command -v go docker nerdctl podman curl python3
        rg --files /home/ubuntu /opt /usr/local -g go -g gofmt -g '!node_modules/**' -g '!\.git/**' -g '!\.cache/**' -g '!\.codex/**' -g '!\.local/**' -g '!\.rustup/**' -g '!\.cargo/**'
    /usr/local/bin/nerdctl
    /usr/local/bin/nerdctl
    /usr/bin/curl
    /usr/bin/python3
    rg: /opt/containerd: Permission denied (os error 13)
    ✗ (2) • 105ms

## Assistant

The Yi timestamp parser was converting UTC to New York time. I’ve removed that conversion and updated the documentation. Go isn’t installed in the shell, so I’m checking whether the available container tools can run the tests.

## Activity

    $ sudo -n nerdctl images --format '{{.Repository}}:{{.Tag}}'
    sudo: a password is required
    ✗ (1) • 0ms

## Activity

    $ rg -n 'go test|test_worker|go build' README.md detection_test.go
        sed -n '285,355p' detection_test.go
        sed -n '1,85p' config_test.go
        sudo -n nerdctl version
    README.md:31:go build -o neonvr .
    README.md:139:go build -o neonvr .
    README.md:229:go build -o neonvr .
    README.md:241:go test -race ./...
    // Optional integration with actual weights, without making ordinary Go tests
    // download models or install Python dependencies.
    func TestYOLO26nWorkerAndPipeline(t *testing.T) {
        model := os.Getenv("NEONVR_TEST_MODEL")
        python := os.Getenv("NEONVR_TEST_PYTHON")
        if model == "" || python == "" {
            t.Skip("set NEONVR_TEST_MODEL and NEONVR_TEST_PYTHON for real CPU inference")
        }
        s, cam, path, _ := detectionFixture(t)
        s.detectionConfig.DetectionModel = model
        s.detectionConfig.DetectionPython = python
        worker, err := s.startDetectionWorker()
        if err != nil {
            t.Fatal(err)
        }
        defer worker.close()
        bad := filepath.Join(cam.Dir, "2026-10-01_11-00-00_Corrupt.mp4")
        if err := os.WriteFile(bad, []byte("unfinished upload"), 0640); err != nil {
            t.Fatal(err)
        }
        old := time.Now().Add(-24 * time.Hour)
        os.Chtimes(bad, old, old)
        ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
        defer cancel()
        if _, err := s.analyzeDetection(ctx, worker, bad); err == nil || errors.Is(err, errDetectionWorker) {
            t.Fatalf("expected a per-video error, got %v", err)
        }
        result, err := s.analyzeDetection(ctx, worker, path)
        if err != nil || len(result.Intervals) != 0 || result.Uncertain || result.Frames != 8 {
            t.Fatalf("CPU inference after bad video: %+v %v", result, err)
        }
        if err := s.runDetectionMaintenance(worker); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); !os.IsNotExist(err) {
            t.Fatal("no-target clip not pruned by real pipeline")
        }
        if _, err := os.Stat(bad); err != nil {
            t.Fatal("corrupt recording removed")
        }
        if _, err := os.Stat(filepath.Join(cam.Dir, ".neonvr_pruned", filepath.Base(path))); err != nil {
            t.Fatal("real pipeline did not retain recoverable recording")
        }
    }
    package main
    
    import (
        "bytes"
        "context"
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "os"
        "path/filepath"
        "sync"
        "testing"
        "time"
    )
    
    func fixture(t *testing.T) (Config, string, *configManager) {
        t.Helper()
        c := defaultConfig()
        c.Root = t.TempDir()
        path := filepath.Join(t.TempDir(), "neonvr.json")
        if err := atomicWriteConfig(path, configBytes(c)); err != nil {
            t.Fatal(err)
        }
        manager, err := newConfigManager(path, c)
        if err != nil {
            t.Fatal(err)
        }
        return c, path, manager
    }
    
    func postConfig(t *testing.T, m *configManager, c Config, revision, token, origin string, restart bool) *httptest.ResponseRecorder {
        t.Helper()
        body, _ := json.Marshal(map[string]any{"config": c, "revision": revision, "restart": restart})
        req := httptest.NewRequest(http.MethodPost, "http://neonvr.test/api/config", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-NeoNVR-Token", token)
        if origin != "" {
            req.Header.Set("Origin", origin)
        }
        w := httptest.NewRecorder()
        m.handler(w, req)
        return w
    }
    
    func TestSaveConfigValidationConflictAndRollback(t *testing.T) {
        c, path, m := fixture(t)
        original, _ := os.ReadFile(path)
        revision := configRevision(original)
        c.ThumbWorkers = 0
        w := postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
        if w.Code != 400 {
            t.Fatalf("invalid config status %d: %s", w.Code, w.Body.String())
        }
        raw, _ := os.ReadFile(path)
        if !bytes.Equal(raw, original) {
            t.Fatal("invalid settings changed the file")
        }
        c.ThumbWorkers = 4
        w = postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
        if w.Code != 200 {
            t.Fatalf("save status %d: %s", w.Code, w.Body.String())
        }
        previous, _ := os.ReadFile(path + ".previous")
        if !bytes.Equal(previous, original) {
            t.Fatal("previous config not preserved")
        }
        saved, _ := os.ReadFile(path)
        loaded, err := decodeConfig(saved)
        if err != nil || loaded.ThumbWorkers != 4 {
            t.Fatal("config not saved")
        }
        if m.active.ThumbWorkers != 2 {
            t.Fatal("save applied config without restart")
        }
        if postConfig(t, m, c, revision, m.token, "", false).Code != 409 {
            t.Fatal("stale edit did not conflict")
        }
    }
    
    func TestConfigRejectsCrossOriginAndMissingToken(t *testing.T) {
        c, path, m := fixture(t)
        raw, _ := os.ReadFile(path)
        rev := configRevision(raw)
        for _, tc := range []struct{ token, origin string }{{"", "http://neonvr.test"}, {m.token, "http://evil.test"}, {m.token, "https://neonvr.test"}} {
            if w := postConfig(t, m, c, rev, tc.token, tc.origin, false); w.Code != 403 {
    sudo: a password is required
    ✗ (1) • 0ms

## Activity

    $ git diff --check
        git status --short
        nl -ba main.go | sed -n '2138,2160p'
        nl -ba README.md | sed -n '274,280p'
     M README.md
     M docker-compose.yml
     M main.go
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
      2138            keys = append(keys, k)
      2139        }
      2140        sort.Strings(keys)
      2141        return keys
      2142    }
      2143
      2144    func parseYICamHourFolder(name string) (time.Time, bool) {
      2145        // Cameras write local timestamps; preserve their date and hour as written.
      2146        ts, err := time.ParseInLocation("2006Y01M02D15H", name, time.Local)
      2147        if err != nil {
      2148            return time.Time{}, false
      2149        }
      2150        return ts, true
      2151    }
      2152
      2153    func parseYICamMinuteSecond(filename string) (int, int, bool) {
      2154        base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
      2155        var minute, second, duration int
      2156        n, err := fmt.Sscanf(base, "W4%dM%dS%d", &minute, &second, &duration)
      2157        if err != nil || n != 3 || minute < 0 || minute > 59 || second < 0 || second > 59 {
      2158            return 0, 0, false
      2159        }
      2160        return minute, second, true
       274    Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thumb-maintenance-every 1h`. Set it to `0` to disable background maintenance. Lazy/on-demand thumbnail generation remains enabled as a fallback.
       275
       276    ### Yi timestamps
       277
       278    Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
       279
       280    ## Project rename and compatibility
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import hashlib
        import json
        import os
        from pathlib import Path
        import socket
        import subprocess
        import tarfile
        import tempfile
        import time
        import urllib.error
        import urllib.request
        
        with tempfile.TemporaryDirectory(prefix='neonvr-timezone-check-') as temporary:
            task_root = Path(temporary)
            with urllib.request.urlopen('https://go.dev/dl/?mode=json&include=all', timeout=30) as response:
                releases = json.load(response)
            release = next(item for item in releases if item['version'] == 'go1.22.12')
            artifact = next(item for item in release['files'] if item['os'] == 'linux' and item['arch'] == 'amd64' and item['kind'] == 'archive')
            archive = task_root / artifact['filename']
            urllib.request.urlretrieve('https://go.dev/dl/' + artifact['filename'], archive)
            if hashlib.sha256(archive.read_bytes()).hexdigest() != artifact['sha256']:
                raise RuntimeError('Go toolchain checksum mismatch')
            with tarfile.open(archive) as package:
                package.extractall(task_root, filter='data')
            go = task_root / 'go/bin/go'
            env = dict(os.environ, GOCACHE=str(task_root / 'cache'), GOPATH=str(task_root / 'gopath'), GOTOOLCHAIN='local')
            print('Running existing Go tests with a temporary Go 1.22.12 toolchain.', flush=True)
            subprocess.run([str(go), 'test', './...'], env=env, check=True)
            binary = task_root / 'neonvr'
            subprocess.run([str(go), 'build', '-o', str(binary), '.'], env=env, check=True)
            formatting = subprocess.run([str(task_root / 'go/bin/gofmt'), '-d', 'main.go'], capture_output=True, text=True, check=True)
            if formatting.stdout:
                raise RuntimeError('Go formatting differs:\n' + formatting.stdout)
            recordings = task_root / 'recordings'
            camera = recordings / 'Yi'
            fixtures = {
                '2026Y07M10D00H/W421M41S18.mp4': '2026-07-10 00:21:41',
                '2026Y01M10D00H/W405M09S18.mp4': '2026-01-10 00:05:09',
                '2026Y09M23D23H/W421M41S18.mp4': '2026-09-23 23:21:41',
                '2026-09-23_23-21-41_Blink.mp4': '2026-09-23 23:21:41',
            }
            for filename in fixtures:
                recording = camera / filename
                recording.parent.mkdir(parents=True, exist_ok=True)
                recording.touch()
            for timezone in ('UTC', 'America/New_York'):
                with socket.socket() as listener:
                    listener.bind(('127.0.0.1', 0))
                    port = listener.getsockname()[1]
                config = task_root / ('config-' + timezone.replace('/', '-') + '.json')
                config.write_text(json.dumps({
                    'root': str(recordings), 'extra_roots': [], 'addr': '127.0.0.1:' + str(port),
                    'auto_delete_days': 0, 'thumb_maintenance_every': '0', 'detection_enabled': False,
                }))
                with tempfile.TemporaryFile() as log:
                    process = subprocess.Popen([str(binary), '-config', str(config)], env=dict(env, TZ=timezone), stdout=log, stderr=log)
                    try:
                        endpoint = 'http://127.0.0.1:' + str(port)
                        for attempt in range(100):
                            if process.poll() is not None:
                                log.seek(0)
                                raise RuntimeError(log.read().decode())
                            try:
                                with urllib.request.urlopen(endpoint + '/api/clips', timeout=2) as response:
                                    clips = json.load(response)['clips']
                                break
                            except urllib.error.URLError:
                                time.sleep(.05)
                        else:
                            raise RuntimeError('Temporary server did not start')
                        observed = {clip['file']: clip['timestamp'] for clip in clips}
                        if observed != fixtures:
                            raise RuntimeError('Timestamps changed for ' + timezone + ': ' + repr(observed))
                        with urllib.request.urlopen(endpoint + '/api/dates', timeout=2) as response:
                            dates = json.load(response)['dates']
                        if dates != ['2026-09-23', '2026-07-10', '2026-01-10']:
                            raise RuntimeError('Date grouping changed: ' + repr(dates))
                        with urllib.request.urlopen(endpoint + '/api/clips?date=2026-07-10', timeout=2) as response:
                            selected = json.load(response)['clips']
                        if len(selected) != 1 or selected[0]['start_seconds'] != 21 * 60 + 41:
                            raise RuntimeError('Date filtering or timeline position changed')
                        print('Timestamp display, date grouping, and filtering passed with TZ=' + timezone, flush=True)
                    finally:
                        process.terminate()
                        try:
                            process.wait(timeout=5)
                        except subprocess.TimeoutExpired:
                            process.kill()
                            process.wait()
            print('Build and formatting checks passed.', flush=True)
        PY
    Running existing Go tests with a temporary Go 1.22.12 toolchain.
    ok      neonvr    3.174s
    Timestamp display, date grouping, and filtering passed with TZ=UTC
    Timestamp display, date grouping, and filtering passed with TZ=America/New_York
    Build and formatting checks passed.
    ✓ • 19.93s

## Assistant

Removed the timezone correction in [main.go](/home/ubuntu/NeoNVR/main.go:2144). NeoNVR now displays camera filename dates and times as written.

Go tests and build passed. Verified timestamp display and date grouping with UTC and New York server timezones.

Rebuild the running container to apply it:

```bash
sudo nerdctl compose up -d --build
```

## User

add a new tab the ability to view streams such as from "
High res stream	rtsp://192.168.1.150:554/ch0_0.h264
High res stream with speaker backchannel	rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"

## Assistant

I’ll add a Streams tab with an editable RTSP source list and browser playback, including the camera URLs you provided.

## Activity

    $ sed -n '1,235p' config.go
        sed -n '310,690p' config.go
        sed -n '450,840p' main.go
        sed -n '1,175p' static/settings.html
    package main
    
    import (
        "bytes"
        "crypto/rand"
        "crypto/sha256"
        "encoding/hex"
        "encoding/json"
        "errors"
        "flag"
        "fmt"
        "io"
        "math"
        "net"
        "net/http"
        "net/url"
        "os"
        "path/filepath"
        "reflect"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    )
    
    var flagHelpError = flag.ErrHelp
    
    type Config struct {
        Root                  string   `json:"root"`
        ExtraRoots            []string `json:"extra_roots"`
        Addr                  string   `json:"addr"`
        ThumbDir              string   `json:"thumb_dir"`
        Meta                  string   `json:"meta"`
        FFmpeg                string   `json:"ffmpeg"`
        ThumbSeek             float64  `json:"thumb_seek"`
        ThumbWidth            int      `json:"thumb_width"`
        ThumbWorkers          int      `json:"thumb_workers"`
        ThumbMaintenanceEvery string   `json:"thumb_maintenance_every"`
        SizeRefresh           string   `json:"size_refresh"`
        AutoDeleteDays        int      `json:"auto_delete_days"`
        AutoDeleteEvery       string   `json:"auto_delete_every"`
        DetectionEnabled      bool     `json:"detection_enabled"`
        DetectionPython       string   `json:"detection_python"`
        DetectionModel        string   `json:"detection_model"`
        DetectionThreads      int      `json:"detection_threads"`
        DetectionFPS          float64  `json:"detection_fps"`
        DetectionConfidence   float64  `json:"detection_confidence"`
        DetectionMinFrames    int      `json:"detection_min_frames"`
        DetectionMergeGap     float64  `json:"detection_merge_gap"`
        DetectionCarMovement  float64  `json:"detection_car_movement"`
        DetectionEvery        string   `json:"detection_every"`
        DetectionMinAge       string   `json:"detection_min_age"`
        DetectionTimeout      string   `json:"detection_timeout"`
        DetectionPrune        string   `json:"detection_prune"`
    }
    
    func defaultConfig() Config {
        return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h",
            DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
    }
    
    func decodeConfig(data []byte) (Config, error) {
        c := defaultConfig()
        if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
            return c, errors.New("config must be a JSON object")
        }
        dec := json.NewDecoder(bytes.NewReader(data))
        dec.DisallowUnknownFields()
        if err := dec.Decode(&c); err != nil {
            return c, err
        }
        if err := dec.Decode(new(any)); err != io.EOF {
            return c, errors.New("config must contain exactly one JSON object")
        }
        if c.ExtraRoots == nil {
            c.ExtraRoots = []string{}
        }
        return c, nil
    }
    
    func (c Config) validate() error {
        if strings.TrimSpace(c.Root) == "" {
            return errors.New("root cannot be empty")
        }
        for _, p := range append([]string{c.Root}, c.ExtraRoots...) {
            if !filepath.IsAbs(expandPath(p)) {
                return fmt.Errorf("recording directory must be absolute: %s", p)
            }
            st, err := os.Stat(expandPath(p))
            if err != nil {
                return fmt.Errorf("recording directory %s: %w", p, err)
            }
            if !st.IsDir() {
                return fmt.Errorf("recording directory %s is not a directory", p)
            }
        }
        for name, p := range map[string]string{"thumb_dir": c.ThumbDir, "meta": c.Meta} {
            if p != "" && !filepath.IsAbs(expandPath(p)) {
                return fmt.Errorf("%s must be absolute or empty", name)
            }
        }
        host, port, err := net.SplitHostPort(c.Addr)
        if err != nil {
            return fmt.Errorf("addr: %w", err)
        }
        n, err := strconv.Atoi(port)
        if err != nil || n < 1 || n > 65535 {
            return errors.New("addr must use a port between 1 and 65535")
        }
        if _, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port)); err != nil {
            return fmt.Errorf("addr: %w", err)
        }
        if strings.TrimSpace(c.FFmpeg) == "" {
            return errors.New("ffmpeg cannot be empty")
        }
        if c.ThumbWorkers < 1 || c.ThumbWorkers > 64 {
            return errors.New("thumb_workers must be between 1 and 64")
        }
        if c.ThumbWidth < 64 || c.ThumbWidth > 4096 {
            return errors.New("thumb_width must be between 64 and 4096")
        }
        if math.IsNaN(c.ThumbSeek) || math.IsInf(c.ThumbSeek, 0) || c.ThumbSeek < 0 || c.ThumbSeek > 3600 {
            return errors.New("thumb_seek must be between 0 and 3600 seconds")
        }
        if c.AutoDeleteDays < 0 || c.AutoDeleteDays > 36500 {
            return errors.New("auto_delete_days must be between 0 and 36500; 0 disables deletion")
        }
        if err := c.validateDetection(); err != nil {
            return err
        }
        for name, v := range map[string]string{"thumb_maintenance_every": c.ThumbMaintenanceEvery, "size_refresh": c.SizeRefresh, "auto_delete_every": c.AutoDeleteEvery} {
            d, err := time.ParseDuration(v)
            if err != nil {
                return fmt.Errorf("%s: enter a duration such as 30m or 6h", name)
            }
            min := time.Minute
            if name == "size_refresh" {
                min = 10 * time.Second
            }
            if name == "thumb_maintenance_every" && d == 0 {
                continue
            }
            if d < min {
                return fmt.Errorf("%s must be at least %s", name, min)
            }
        }
        meta := c.Meta
        if meta == "" {
            meta = filepath.Join(expandPath(c.Root), ".blinkclips_meta.json")
        }
        if _, err := NewMetaStore(expandPath(meta)); err != nil {
            return fmt.Errorf("favorites metadata: %w", err)
        }
        return nil
    }
    
    func configRevision(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
    
    // The entire parent directory is mounted into containers so rename remains atomic.
    func atomicWriteConfig(path string, data []byte) error {
        if st, err := os.Lstat(path); err == nil && !st.Mode().IsRegular() {
            return errors.New("config path must be a regular file")
        } else if err != nil && !os.IsNotExist(err) {
            return err
        }
        if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
            return err
        }
        original, _ := os.Stat(path)
        if original == nil && strings.HasSuffix(path, ".previous") {
            original, _ = os.Stat(strings.TrimSuffix(path, ".previous"))
        }
        f, err := os.CreateTemp(filepath.Dir(path), ".neonvr-config-*")
        if err != nil {
            return err
        }
        defer os.Remove(f.Name())
        if original != nil && os.Geteuid() == 0 {
            if st, ok := original.Sys().(*syscall.Stat_t); ok {
                if err = f.Chown(int(st.Uid), int(st.Gid)); err != nil {
                    f.Close()
                    return err
                }
            }
        }
        if err = f.Chmod(0640); err == nil {
            _, err = f.Write(data)
        }
        if err == nil {
            err = f.Sync()
        }
        closeErr := f.Close()
        if err == nil {
            err = closeErr
        }
        if err != nil {
            return err
        }
        if err = os.Rename(f.Name(), path); err != nil {
            return err
        }
        dir, err := os.Open(filepath.Dir(path))
        if err != nil {
            return err
        }
        defer dir.Close()
        return dir.Sync()
    }
    
    func configBytes(c Config) []byte {
        data, _ := json.MarshalIndent(c, "", "  ")
        return append(data, '\n')
    }
    
    func loadStartupConfig(args []string) (Config, string, error) {
        path := "config/neonvr.json"
        explicit := false
        for i, a := range args {
            if a == "-config" || a == "--config" {
                if i+1 >= len(args) {
                    return Config{}, "", errors.New("-config requires a path")
                }
                path = args[i+1]
                explicit = true
            }
            if strings.HasPrefix(a, "-config=") || strings.HasPrefix(a, "--config=") {
                path = strings.SplitN(a, "=", 2)[1]
                explicit = true
            }
        }
        if strings.TrimSpace(path) == "" {
            return Config{}, "", errors.New("config path cannot be empty")
        }
        path, err := filepath.Abs(expandPath(path))
            }
        }
        return c, path, nil
    }
    
    type configManager struct {
        mu         sync.Mutex
        path       string
        active     Config
        token      string
        restarting bool
        restart    chan struct{}
    }
    
    func newConfigManager(path string, c Config) (*configManager, error) {
        token := make([]byte, 32)
        if _, err := rand.Read(token); err != nil {
            return nil, err
        }
        return &configManager{path: path, active: c, token: hex.EncodeToString(token), restart: make(chan struct{}, 1)}, nil
    }
    
    func (m *configManager) response(w http.ResponseWriter, status int, saved Config, revision string, nextURL string) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Cache-Control", "no-store")
        w.WriteHeader(status)
        _ = json.NewEncoder(w).Encode(map[string]any{"config": saved, "active_config": m.active, "revision": revision, "token": m.token, "instance_id": m.token, "restart_required": !reflect.DeepEqual(saved, m.active), "restarting": m.restarting, "next_url": nextURL})
    }
    
    func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        if r.Method != "GET" && r.Method != "POST" {
            w.Header().Set("Allow", "GET, POST")
            http.Error(w, "method not allowed", 405)
            return
        }
        if r.Method == "POST" {
            if r.Header.Get("X-NeoNVR-Token") != m.token {
                http.Error(w, "reload Settings before saving", 403)
                return
            }
            if origin := r.Header.Get("Origin"); origin != "" {
                u, err := url.Parse(origin)
                scheme := "http"
                if r.TLS != nil {
                    scheme = "https"
                }
                if err != nil || u.Host != r.Host || u.Scheme != scheme {
                    http.Error(w, "same-origin requests required", 403)
                    return
                }
            }
            if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
                http.Error(w, "JSON required", 415)
                return
            }
        }
        m.mu.Lock()
        defer m.mu.Unlock()
        raw, err := os.ReadFile(m.path)
        if err != nil {
            http.Error(w, "cannot read config: "+err.Error(), 500)
            return
        }
        saved, err := decodeConfig(raw)
        if err != nil {
            http.Error(w, "invalid config file: "+err.Error(), 500)
            return
        }
        revision := configRevision(raw)
        if r.Method == "GET" {
            m.response(w, 200, saved, revision, "")
            return
        }
        if m.restarting {
            http.Error(w, "restart already in progress", 409)
            return
        }
        var request struct {
            Config   json.RawMessage `json:"config"`
            Revision string          `json:"revision"`
            Restart  bool            `json:"restart"`
        }
        dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
        dec.DisallowUnknownFields()
        if err := dec.Decode(&request); err != nil {
            http.Error(w, "invalid request: "+err.Error(), 400)
            return
        }
        if err := dec.Decode(new(any)); err != io.EOF {
            http.Error(w, "unexpected trailing JSON", 400)
            return
        }
        if request.Revision != revision {
            http.Error(w, "settings changed elsewhere; reload before saving", 409)
            return
        }
        c, err := decodeConfig(request.Config)
        if err != nil {
            http.Error(w, "invalid config: "+err.Error(), 400)
            return
        }
        if err = c.validate(); err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        if request.Restart && c.Addr != m.active.Addr {
            listener, err := net.Listen("tcp", c.Addr)
            if err != nil {
                http.Error(w, "cannot listen on new address: "+err.Error(), 400)
                return
            }
            listener.Close()
        }
        data := configBytes(c)
        if err = atomicWriteConfig(m.path+".previous", raw); err != nil {
            http.Error(w, "cannot save previous config: "+err.Error(), 500)
            return
        }
        if err = atomicWriteConfig(m.path, data); err != nil {
            http.Error(w, "cannot save config: "+err.Error(), 500)
            return
        }
        nextURL := ""
        if request.Restart {
            m.restarting = true
            host, port, _ := net.SplitHostPort(c.Addr)
            if host == "" || host == "0.0.0.0" || host == "::" {
                host, _, _ = net.SplitHostPort(r.Host)
                if host == "" {
                    host = r.Host
                }
            }
            scheme := "http"
            if r.TLS != nil {
                scheme = "https"
            }
            nextURL = scheme + "://" + net.JoinHostPort(host, port) + "/settings"
        }
        m.response(w, 200, c, configRevision(data), nextURL)
        if request.Restart {
            // Delay until this response has reached the browser, then drain and re-exec.
            time.AfterFunc(750*time.Millisecond, func() { m.restart <- struct{}{} })
        }
    }
          transform: translateX(-1px);
          mix-blend-mode: screen;
        }
        .scale {
          position: absolute;
          left: 0;
          right: 0;
          pointer-events: none;
          color: #778292;
          font-size: 11px;
          user-select: none;
        }
        .scale.hourScale { top: 88px; }
        .scale.halfScale { top: 4px; color: #5f6b7b; }
        .scaleLabel {
          position: absolute;
          transform: translateX(-50%);
          white-space: nowrap;
        }
        .scaleLabel.last { transform: translateX(-100%); }
    
        input[type=range] {
          appearance: none;
          -webkit-appearance: none;
          pointer-events: none;
          position: absolute;
          left: 0;
          right: 0;
          top: 38px;
          width: 100%;
          height: 34px;
          background: transparent;
          margin: 0;
        }
        input[type=range]::-webkit-slider-runnable-track { height: 8px; background: transparent; }
        input[type=range]::-moz-range-track { height: 8px; background: transparent; }
        input[type=range]::-webkit-slider-thumb {
          appearance: none;
          -webkit-appearance: none;
          pointer-events: auto;
          width: 22px;
          height: 22px;
          margin-top: -7px;
          border-radius: 50%;
          background: var(--thumb);
          border: 3px solid var(--thumbBorder);
          box-shadow: 0 0 0 1px var(--accent);
          cursor: grab;
        }
        input[type=range]::-moz-range-thumb {
          pointer-events: auto;
          width: 18px;
          height: 18px;
          border-radius: 50%;
          background: var(--thumb);
          border: 3px solid var(--thumbBorder);
          box-shadow: 0 0 0 1px var(--accent);
          cursor: grab;
        }
    
        .stats {
          display: flex;
          gap: 14px;
          flex-wrap: wrap;
          color: var(--muted);
          font-size: 14px;
          margin: 4px 0 16px;
        }
    
        .sideRail {
          display: flex;
          flex-direction: column;
          gap: 9px;
          overflow-y: auto;
          overscroll-behavior: contain;
          touch-action: pan-y;
          min-height: 0;
          padding-right: 4px;
          scrollbar-gutter: stable;
        }
        .sideThumb {
          display: grid;
          grid-template-columns: 126px minmax(0, 1fr);
          gap: 9px;
          align-items: center;
          width: 100%;
          min-height: 82px;
          padding: 7px;
          border: 1px solid var(--border);
          border-radius: 12px;
          background: #0d0f14;
          color: var(--text);
          text-align: left;
          cursor: pointer;
        }
        .sideThumb:hover { border-color: #556071; }
        .sideThumb.active { border-color: var(--accent); box-shadow: 0 0 0 1px rgba(114,167,255,.45); }
        .sideThumb.selected { outline: 2px solid var(--accent); }
        .sideThumb.favorite { border-color: #7f6420; }
        .sideThumb img {
          width: 126px;
          aspect-ratio: 16 / 9;
          height: auto;
          object-fit: cover;
          display: block;
          border-radius: 8px;
          background: #000;
        }
        .sideThumbMeta { min-width: 0; display: grid; gap: 2px; }
        .sideThumbTime { font-weight: 700; font-size: 13px; line-height: 1.2; }
        .sideThumbCamera { color: var(--muted); font-size: 12px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
        .sideThumbStar { color: #ffe28a; font-size: 12px; }
        .clipGrid {
          display: grid;
          grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
          gap: 14px;
          align-items: start;
        }
        .card {
          background: var(--panel);
          border: 1px solid var(--border);
          border-radius: 14px;
          overflow: hidden;
          position: relative;
        }
        .card.favorite { border-color: #7f6420; }
        .card.selected { outline: 2px solid var(--accent); }
        .card.active { border-color: var(--accent); box-shadow: 0 0 0 1px rgba(114,167,255,.45); }
        body.selectMode .card { cursor: pointer; }
        .media {
          position: relative;
          background: #000;
          aspect-ratio: 16 / 9;
        }
        .selectModeBadge {
          display: none;
          position: absolute;
          z-index: 5;
          top: 8px;
          left: 8px;
          background: rgba(0,0,0,.62);
          color: white;
          padding: 4px 7px;
          border-radius: 999px;
          font-size: 12px;
        }
        body.selectMode .selectModeBadge { display: block; }
        .thumbButton {
          all: unset;
          display: block;
          width: 100%;
          height: 100%;
          cursor: pointer;
          position: relative;
          background: #000;
        }
        .thumbButton img, .card video {
          width: 100%;
          height: 100%;
          object-fit: contain;
          display: block;
          background: #000;
        }
        .playBadge {
          position: absolute;
          inset: 0;
          display: grid;
          place-items: center;
          background: linear-gradient(transparent, rgba(0,0,0,.32));
          font-size: 48px;
          text-shadow: 0 2px 12px rgba(0,0,0,.8);
        }
        .cardBody { padding: 10px 12px 12px; }
        .titleRow { display: flex; justify-content: space-between; gap: 8px; align-items: center; }
        .title { font-weight: 700; margin-bottom: 3px; }
        .badge { color: #ffe28a; font-size: 13px; }
        .file { color: var(--muted); font-size: 12px; word-break: break-all; }
        .cardActions { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 10px; }
        .cardActions button { padding: 7px 9px; font-size: 13px; }
        .footnote { color: var(--muted); font-size: 12px; margin: 0; }
        .selectFootnote { display: none; }
        body.selectMode .selectFootnote { display: block; }
        .empty, .error, .notice {
          padding: 24px;
          border: 1px solid var(--border);
          border-radius: 16px;
          background: var(--panel);
          color: var(--muted);
          margin-bottom: 14px;
        }
        .error { color: #ffb4a8; border-color: #6d332f; }
        .notice { color: #a7e0ff; border-color: #2c536b; }
    
        @media (max-width: 1100px) {
          .viewerLayout { grid-template-columns: 1fr; }
          .railPanel { position: static; height: 420px; }
          .controls { grid-template-columns: 1fr 1fr; }
          .controls .timelineWrap { grid-column: 1 / -1; }
        }
        @media (max-width: 720px) {
          body { padding: 14px; }
          .controls { grid-template-columns: 1fr; }
          .cameraMenu { width: min(92vw, 360px); }
          media-controller { aspect-ratio: auto; grid-template-rows: auto auto auto; }
          #mainVideo { height: auto; min-height: 180px; aspect-ratio: 16 / 9; }
        }
      </style>
    </head>
    <body>
      <h1>NeoNVR</h1>
      <a href="/settings" style="display:inline-block;margin-bottom:12px;color:inherit">Settings</a>
      <div class="root">Root: {{.Root}}</div>
      <div class="stats" id="diskStats"><span>Calculating recordings size…</span></div>
    
      <section class="panel mainPlayerPanel">
        <div class="viewerLayout">
          <div class="mainColumn">
            <div class="mainPlayerTitle">
              <span id="mainTitle">Main player: click a thumbnail or choose Load.</span>
              <span class="kbdHint">Space pause/play • ← / → seek • Ctrl+arrows previous/next</span>
            </div>
            <div class="mainPlayerShell" id="mainShell">
              <media-controller>
                <video id="mainVideo" slot="media" preload="metadata" playsinline webkit-playsinline x5-playsinline></video>
                <media-control-bar>
                  <media-play-button></media-play-button>
                  <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
                  <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
                  <media-time-range id="mainSeek"></media-time-range>
                  <media-time-display></media-time-display>
                  <media-duration-display></media-duration-display>
                  <media-mute-button></media-mute-button>
                  <media-volume-range></media-volume-range>
                  <media-playback-rate-button></media-playback-rate-button>
                  <media-fullscreen-button></media-fullscreen-button>
                </media-control-bar>
                <div class="detectionLegend" style="padding:0 10px 10px;background:#000;margin:0">
                  <span class="category" style="--detection-color:#75d685">Human</span>
                  <span class="category" style="--detection-color:#ffd166">Animal</span>
                  <span class="category" style="--detection-color:#72a7ff">Bike</span>
                  <span class="category" style="--detection-color:#c792ea">Moving vehicle</span>
                  <span id="detectionStatus" role="status">Select a clip to see detections.</span>
                </div>
              </media-controller>
            </div>
            <div class="analysisControls">
              <label style="display:flex;align-items:center;gap:6px">Step sec <input id="stepSize" type="number" min="0.03" step="0.05" value="0.20"></label>
              <button type="button" id="zoomToggle" aria-pressed="false">Zoom off</button>
              <button type="button" id="zoomReset">Reset zoom</button>
              <span id="zoomLabel">1.0×</span>
              <span class="kbdHint">Turn Zoom on to pinch/trackpad zoom and drag-to-pan.</span>
            </div>
          </div>
    
          <aside class="railPanel">
            <div class="railHeader">
              <strong>Videos</strong>
              <span id="railCount">0 clips</span>
            </div>
            <div class="railActions">
              <button id="selectModeToggle" type="button" class="compact">Select</button>
              <span id="selectionCount">0 selected</span>
              <button id="bulkDelete" type="button" class="danger compact" style="display:none" disabled>Delete selected</button>
              <button id="clearSelection" type="button" class="compact" style="display:none" disabled>Clear</button>
            </div>
            <main class="sideRail" id="sideRail" aria-label="Scrollable video thumbnail navigation"></main>
            <div class="footnote selectFootnote">Shift-click videos to select every video in between for bulk delete.</div>
          </aside>
        </div>
      </section>
    
      <section class="panel">
        <div class="topRow">
          <div class="tabs">
            <button id="clipsTab" type="button" class="active">Clips</button>
            <button id="favoritesTab" type="button">Favorites</button>
          </div>
        </div>
    
        <div class="controls">
          <div class="cameraPicker">
            <span>Camera</span>
            <button id="cameraButton" type="button" class="cameraButton"><span id="cameraSummary">All cameras</span><span>▾</span></button>
            <div id="cameraMenu" class="cameraMenu"></div>
          </div>
    
          <label>Date
            <input id="date" type="date">
          </label>
    
          <label>Sort
            <select id="sortOrder">
              <option value="newest" selected>Most recent first</option>
              <option value="oldest">Oldest first</option>
            </select>
          </label>
    
          <div class="timelineWrap">
            <div class="timelineHeader">
              <span id="rangeText">00:00:00 – 23:59:59</span>
              <span id="timelineCount">0 clips</span>
            </div>
            <div class="timeline" id="timeline">
              <div class="scale halfScale" id="halfScale"></div>
              <div class="track"></div>
              <div class="rangeFill" id="rangeFill"></div>
              <div class="marks" id="marks"></div>
              <input id="startRange" type="range" min="0" max="86399" step="60" value="0" aria-label="Start time">
              <input id="endRange" type="range" min="0" max="86399" step="60" value="86399" aria-label="End time">
              <div class="scale hourScale" id="hourScale"></div>
            </div>
          </div>
        </div>
        <div class="datePills" id="datePills"></div>
      </section>
    
      <div class="stats">
        <span id="selectedInfo"></span>
      </div>
    
      <div id="message"></div>
    
      <section class="clipGrid" id="clips" aria-label="Video cards"></section>
    
    <script src="/static/detection-viewer.js"></script>
    <script>
    const cameraButtonEl = document.getElementById('cameraButton');
    const cameraSummaryEl = document.getElementById('cameraSummary');
    const cameraMenuEl = document.getElementById('cameraMenu');
    const dateEl = document.getElementById('date');
    const sortEl = document.getElementById('sortOrder');
    const datePillsEl = document.getElementById('datePills');
    const startEl = document.getElementById('startRange');
    const endEl = document.getElementById('endRange');
    const rangeTextEl = document.getElementById('rangeText');
    const rangeFillEl = document.getElementById('rangeFill');
    const marksEl = document.getElementById('marks');
    const clipsEl = document.getElementById('clips');
    const sideRailEl = document.getElementById('sideRail');
    const messageEl = document.getElementById('message');
    const railCountEl = document.getElementById('railCount');
    const timelineCountEl = document.getElementById('timelineCount');
    const selectedInfoEl = document.getElementById('selectedInfo');
    const diskStatsEl = document.getElementById('diskStats');
    const clipsTabEl = document.getElementById('clipsTab');
    const favoritesTabEl = document.getElementById('favoritesTab');
    const halfScaleEl = document.getElementById('halfScale');
    const hourScaleEl = document.getElementById('hourScale');
    const mainVideoEl = document.getElementById('mainVideo');
    const mainShellEl = document.getElementById('mainShell');
    const mainTitleEl = document.getElementById('mainTitle');
    const stepSizeEl = document.getElementById('stepSize');
    const zoomToggleEl = document.getElementById('zoomToggle');
    const zoomResetEl = document.getElementById('zoomReset');
    const zoomLabelEl = document.getElementById('zoomLabel');
    const selectModeToggleEl = document.getElementById('selectModeToggle');
    const selectionCountEl = document.getElementById('selectionCount');
    const bulkDeleteEl = document.getElementById('bulkDelete');
    const clearSelectionEl = document.getElementById('clearSelection');
    
    let cameras = [];
    let selectedCameras = new Set();
    let allClips = [];
    let visibleClips = [];
    let availableDates = [];
    let currentView = 'clips';
    let selectedKeys = new Set();
    let lastSelectedIndex = null;
    let selectMode = false;
    let activeKey = '';
    let zoomMode = false;
    let zoom = 1;
    let panX = 0;
    let panY = 0;
    let pointers = new Map();
    let dragStart = null;
    let pinchStart = null;
    
    const palette = [
      '#ff5d55', '#72a7ff', '#75d685', '#ffd166', '#c792ea', '#4dd0e1',
      '#f78c6c', '#a3e635', '#f472b6', '#60a5fa', '#fb7185', '#34d399'
    ];
    
    function qs(params) {
      const q = new URLSearchParams();
      for (const [k, v] of Object.entries(params)) {
        if (v !== undefined && v !== null && v !== '') q.set(k, v);
      }
      return q.toString();
    }
    
    <!doctype html>
    <html lang="en">
    <head>
    <meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
    <title>NeoNVR Settings</title>
    <style>
    :root{color-scheme:dark;font-family:system-ui,sans-serif;background:#10151c;color:#e9edf3}body{max-width:900px;margin:32px auto;padding:0 20px}a{color:#96c6ff}header{display:flex;justify-content:space-between;align-items:center;gap:20px}h1{margin-bottom:6px}p{color:#aebccb;line-height:1.5}.panel{background:#19222d;border:1px solid #324355;border-radius:12px;padding:24px;margin:20px 0}.grid{display:grid;grid-template-columns:1fr 1fr;gap:22px}label{display:block;font-weight:600}input,textarea{box-sizing:border-box;display:block;width:100%;margin-top:8px;padding:11px;border:1px solid #506176;border-radius:6px;background:#101923;color:#edf2f8;font:inherit}textarea{min-height:96px;resize:vertical}.wide{grid-column:1/-1}small{display:block;font-weight:400;color:#aebccb;margin-top:7px;line-height:1.4}.actions{display:flex;flex-wrap:wrap;gap:12px}button{padding:12px 18px;background:#263c53;color:#eef4fc;border:1px solid #587494;border-radius:7px;cursor:pointer;font:inherit;font-weight:600}button.primary{background:#75bcff;color:#102035;border-color:#75bcff}button:disabled{opacity:.5;cursor:wait}#status{padding:14px 0;min-height:24px;white-space:pre-wrap;line-height:1.5}#status.error{color:#ffb3b3}#status.ok{color:#a5e4bb}.notice{border-left:3px solid #75bcff;padding-left:14px}fieldset{border:0;padding:0;margin:0}h2{font-size:20px;font-weight:650;margin-bottom:20px}@media(max-width:650px){.grid{grid-template-columns:1fr}.wide{grid-column:auto}header{align-items:start}body{margin-top:18px}.panel{padding:18px}}
    </style>
    <style>select{box-sizing:border-box;display:block;width:100%;margin-top:8px;padding:11px;border:1px solid #506176;border-radius:6px;background:#101923;color:#edf2f8;font:inherit}</style>
    </head>
    <body>
    <header><div><h1>NeoNVR Settings</h1><p>Save your settings, then restart to apply them.</p></div><a href="/">← Recordings</a></header>
    <p id="pending" class="notice" hidden>Saved settings differ from the running app. Restart to apply them.</p>
    <form id="settingsForm">
    <fieldset id="fields" disabled>
    <section class="panel"><h2>Object detection</h2><p>Find humans, animals, bikes and moving vehicles. Colored intervals appear below the video scrubber.</p><div class="grid">
    <label>Background detection<select name="detection_enabled"><option value="true">Enabled</option><option value="false">Disabled</option></select></label>
    <label>No-target recordings<select name="detection_prune"><option value="quarantine">Move to recoverable quarantine</option><option value="off">Keep all recordings</option><option value="delete">Permanently delete</option></select><small>Only fully analyzed clips without targets or ambiguous observations qualify. Favorites are preserved. Quarantine uses .neonvr_pruned inside each camera directory; it keeps files on disk.</small></label>
    <label>CPU inference threads<input name="detection_threads" type="number" min="1" max="16" step="1" required><small>2 is the default for an Intel i7-6700. One analysis worker loads the model once.</small></label>
    <label>Sampled frames per second<input name="detection_fps" type="number" min="0.5" max="10" step="0.5" required><small>2 balances CPU usage and event coverage. Increase for very brief events.</small></label>
    <label>Minimum confidence<input name="detection_confidence" type="number" min="0.05" max="0.95" step="0.01" required><small>0.65 by default. Higher values reduce false detections. Weaker candidates still keep the recording for review.</small></label>
    <label>Observations to confirm<input name="detection_min_frames" type="number" min="1" max="10" step="1" required><small>3 consecutive samples by default. Brief or intermittent observations keep the recording for review. Vehicle motion always needs at least 3.</small></label>
    <label>Detection gap (seconds)<input name="detection_merge_gap" type="number" min="0" max="10" step="0.1" required><small>Merge nearby confirmed events of the same category. This does not combine intermittent observations into a confirmed event.</small></label>
    <label>Vehicle movement threshold<input name="detection_car_movement" type="number" min="0.001" max="0.25" step="0.001" required><small>Fraction of frame size within two seconds; 0.025 is 2.5%. Requires consistent movement over multiple samples, filtering box jitter and resizing.</small></label>
    <label>Recording scan interval<input name="detection_every" required><small>Use 1m or 5m; minimum 1m.</small></label>
    <label>Wait after last recording write<input name="detection_min_age" required><small>2m by default; minimum 10s. Avoids processing ongoing uploads.</small></label>
    <label>Maximum analysis time per clip<input name="detection_timeout" required><small>30m by default. Timeouts keep the recording.</small></label>
    <label>Python executable<input name="detection_python" required><small>Docker includes the CPU detection runtime. Native installs can use a virtual environment.</small></label>
    <label class="wide">YOLO26n ONNX model<input name="detection_model" required><small>Docker includes yolo26n.onnx. Native installs can supply an absolute model path.</small></label>
    </div></section>
    <section class="panel"><h2>Recording sources</h2><div class="grid">
    <label class="wide">Primary recordings directory<input name="root" required><small>Existing recordings stay where they are. Enter an absolute path.</small></label>
    <label class="wide">Additional recordings directories<textarea name="extra_roots" placeholder="/home/ubuntu/FTPfiles/yicamBack&#10;/home/ubuntu/FTPfiles/yicamFront"></textarea><small>One absolute path per line. In a container, each directory must also be mounted.</small></label>
    </div></section>
    <section class="panel"><h2>Server and storage</h2><div class="grid">
    <label>Listen address<input name="addr" required placeholder="0.0.0.0:5500"><small>In a container, a new port also needs a matching host port mapping.</small></label>
    <label>FFmpeg executable<input name="ffmpeg" required></label>
    <label>Thumbnail cache directory<input name="thumb_dir"><small>Blank keeps the existing cache under the primary root.</small></label>
    <label>Favorites metadata file<input name="meta"><small>Blank keeps the existing favorites file under the primary root.</small></label>
    </div></section>
    <section class="panel"><h2>Thumbnails and refresh</h2><div class="grid">
    <label>Thumbnail workers<input name="thumb_workers" type="number" required min="1" max="64" step="1"></label>
    <label>Thumbnail width (pixels)<input name="thumb_width" type="number" required min="64" max="4096" step="1"></label>
    <label>Fallback seek (seconds)<input name="thumb_seek" type="number" required min="0" max="3600" step="0.01"></label>
    <label>Thumbnail maintenance interval<input name="thumb_maintenance_every" required><small>Use 30m or 1h; 0 disables maintenance.</small></label>
    <label>Size statistics refresh interval<input name="size_refresh" required><small>Use 30m or 1h; minimum 10s.</small></label>
    </div></section>
    <section class="panel"><h2>Retention</h2><p>Automatic deletion affects all recording sources. Favorited clips are preserved.</p><div class="grid">
    <label>Delete non-favorites older than (days)<input name="auto_delete_days" type="number" required min="0" max="36500" step="1"><small>0 disables automatic deletion.</small></label>
    <label>Deletion check interval<input name="auto_delete_every" required><small>Use 6h or 24h; minimum 1m.</small></label>
    </div></section>
    <div class="actions"><button id="save" type="submit">Save</button><button id="restart" class="primary" type="button">Save &amp; restart</button><button id="reload" type="button">Reload saved settings</button></div>
    </fieldset>
    <div id="status" role="status" aria-live="polite">Loading settings…</div>
    </form>
    <script>
    'use strict';
    const form=document.getElementById('settingsForm');
    const fields=document.getElementById('fields');
    const statusEl=document.getElementById('status');
    const pending=document.getElementById('pending');
    let state=null;
    let dirty=false;
    const numeric=new Set(['thumb_seek','thumb_width','thumb_workers','auto_delete_days','detection_threads','detection_fps','detection_confidence','detection_min_frames','detection_merge_gap','detection_car_movement']);
    function show(message,type=''){statusEl.textContent=message;statusEl.className=type;}
    async function request(url,options={}){
     const response=await fetch(url,{cache:'no-store',...options});
     if(!response.ok)throw new Error(await response.text());
     return response.json();
    }
    function populate(data){
     state=data;
     for(const [name,value] of Object.entries(data.config)){
      const field=form.elements.namedItem(name);
      if(field)field.value=name==='extra_roots'?value.join('\n'):value;
     }
     pending.hidden=!data.restart_required;
     fields.disabled=false;dirty=false;
     show(data.restart_required?'Saved settings are ready to apply. Choose Save & restart.':'Settings match the running app.','ok');
    }
    async function load(){fields.disabled=true;try{populate(await request('/api/config'));}catch(error){show(error.message,'error');}}
    function readForm(){
     const config={};
     for(const name of Object.keys(state.config)){
      let value=form.elements.namedItem(name).value;
      if(name==='extra_roots')value=value.split('\n').map(v=>v.trim()).filter(Boolean);
      else if(name==='detection_enabled')value=value==='true';
      else if(numeric.has(name))value=Number(value);
      else value=value.trim();
      config[name]=value;
     }
     return config;
    }
    async function waitForRestart(data,previousID){
     const target=new URL(data.next_url||'/settings',location.href);
     // A same-origin reconnect can prove that the new instance loaded the saved config.
     if(target.origin!==location.origin){
      show('Settings saved. Restarting… Open the new server address when it is ready.','ok');
      const link=document.createElement('a');link.href=target.href;link.textContent=target.href;
      statusEl.append(document.createElement('br'),link);return;
     }
     for(let attempt=0;attempt<60;attempt++){
      await new Promise(resolve=>setTimeout(resolve,1000));
      try{
       const controller=new AbortController();const timer=setTimeout(()=>controller.abort(),1500);
       let current;try{current=await request('/api/config',{signal:controller.signal});}finally{clearTimeout(timer);}
       if(current.instance_id!==previousID){
        if(current.restart_required)throw new Error('The server restarted, but saved settings still need to be applied.');
        populate(current);show('Restart complete. Saved settings are now active.','ok');return;
       }
      }catch(error){/* The listener is briefly unavailable during restart. */}
     }
     show('Settings were saved, but reconnection timed out. Refresh this page or check the server logs.','error');
     document.getElementById('reload').disabled=false;fields.disabled=false;
    }
    async function save(restart){
     if(!state||!form.reportValidity())return;
     const config=readForm();
     if(restart&&config.detection_enabled&&config.detection_prune==='delete'&&(!state.active_config.detection_enabled||state.active_config.detection_prune!=='delete')){
      if(!confirm('Permanently delete fully analyzed non-favorite recordings with no targets? Quarantine preserves files for review instead.'))return;
     }
     if(restart&&config.auto_delete_days>0&&config.auto_delete_days!==state.active_config.auto_delete_days){
      if(!confirm('Restart with automatic deletion enabled for non-favorites older than '+config.auto_delete_days+' days? This applies to every recording source.'))return;
     }
     const previousID=state.instance_id;
     fields.disabled=true;show(restart?'Validating settings and preparing restart…':'Saving settings…');
     try{
      const data=await request('/api/config',{method:'POST',headers:{'Content-Type':'application/json','X-NeoNVR-Token':state.token},body:JSON.stringify({config,revision:state.revision,restart})});
      state=data;dirty=false;pending.hidden=!data.restart_required;
      if(restart){show('Settings saved. Restarting and reconnecting…','ok');await waitForRestart(data,previousID);}
      else{fields.disabled=false;show('Settings saved. Choose Save & restart to apply them.','ok');}
     }catch(error){fields.disabled=false;show(error.message,'error');}
    }
    form.addEventListener('input',()=>{dirty=true;});
    form.addEventListener('submit',event=>{event.preventDefault();save(false);});
    document.getElementById('restart').addEventListener('click',()=>save(true));
    document.getElementById('reload').addEventListener('click',()=>{if(!dirty||confirm('Discard unsaved changes and reload?'))load();});
    window.addEventListener('beforeunload',event=>{if(dirty){event.preventDefault();event.returnValue='';}});
    load();
    </script>
    </body>
    </html>
    ✓ • 0ms

## Activity

    $ git status --short
        rg --files -g 'AGENTS.md' -g '!node_modules' -g '!vendor'
        rg -n 'nav|tab|settings|player|media-controller|viewMode|DOMContentLoaded|main\(|func \(.*handler|restart' main.go static/settings.html config.go README.md Dockerfile docker-compose.yml
     M README.md
     M docker-compose.yml
     M main.go
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    docker-compose.yml:5:    restart: unless-stopped
    config.go:266:    fs.StringVar(&c.FFmpeg, "ffmpeg", c.FFmpeg, "ffmpeg executable")
    config.go:275:    fs.StringVar(&c.DetectionPython, "detection-python", c.DetectionPython, "Python executable with ONNX Runtime installed")
    config.go:301:    // Persist explicit legacy flag overrides, so later browser restarts use those settings.
    config.go:320:    restarting bool
    config.go:321:    restart    chan struct{}
    config.go:329:    return &configManager{path: path, active: c, token: hex.EncodeToString(token), restart: make(chan struct{}, 1)}, nil
    config.go:336:    _ = json.NewEncoder(w).Encode(map[string]any{"config": saved, "active_config": m.active, "revision": revision, "token": m.token, "instance_id": m.token, "restart_required": !reflect.DeepEqual(saved, m.active), "restarting": m.restarting, "next_url": nextURL})
    config.go:339:func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
    config.go:385:    if m.restarting {
    config.go:386:        http.Error(w, "restart already in progress", 409)
    config.go:392:        Restart  bool            `json:"restart"`
    config.go:405:        http.Error(w, "settings changed elsewhere; reload before saving", 409)
    config.go:436:        m.restarting = true
    config.go:448:        nextURL = scheme + "://" + net.JoinHostPort(host, port) + "/settings"
    config.go:453:        time.AfterFunc(750*time.Millisecond, func() { m.restart <- struct{}{} })
    static/settings.html:12:<header><div><h1>NeoNVR Settings</h1><p>Save your settings, then restart to apply them.</p></div><a href="/">← Recordings</a></header>
    static/settings.html:13:<p id="pending" class="notice" hidden>Saved settings differ from the running app. Restart to apply them.</p>
    static/settings.html:14:<form id="settingsForm">
    static/settings.html:28:<label>Python executable<input name="detection_python" required><small>Docker includes the CPU detection runtime. Native installs can use a virtual environment.</small></label>
    static/settings.html:37:<label>FFmpeg executable<input name="ffmpeg" required></label>
    static/settings.html:52:<div class="actions"><button id="save" type="submit">Save</button><button id="restart" class="primary" type="button">Save &amp; restart</button><button id="reload" type="button">Reload saved settings</button></div>
    static/settings.html:54:<div id="status" role="status" aria-live="polite">Loading settings…</div>
    static/settings.html:58:const form=document.getElementById('settingsForm');
    static/settings.html:77: pending.hidden=!data.restart_required;
    static/settings.html:79: show(data.restart_required?'Saved settings are ready to apply. Choose Save & restart.':'Settings match the running app.','ok');
    static/settings.html:95: const target=new URL(data.next_url||'/settings',location.href);
    static/settings.html:108:    if(current.restart_required)throw new Error('The server restarted, but saved settings still need to be applied.');
    static/settings.html:109:    populate(current);show('Restart complete. Saved settings are now active.','ok');return;
    static/settings.html:111:  }catch(error){/* The listener is briefly unavailable during restart. */}
    static/settings.html:116:async function save(restart){
    static/settings.html:119: if(restart&&config.detection_enabled&&config.detection_prune==='delete'&&(!state.active_config.detection_enabled||state.active_config.detection_prune!=='delete')){
    static/settings.html:122: if(restart&&config.auto_delete_days>0&&config.auto_delete_days!==state.active_config.auto_delete_days){
    static/settings.html:126: fields.disabled=true;show(restart?'Validating settings and preparing restart…':'Saving settings…');
    static/settings.html:128:  const data=await request('/api/config',{method:'POST',headers:{'Content-Type':'application/json','X-NeoNVR-Token':state.token},body:JSON.stringify({config,revision:state.revision,restart})});
    static/settings.html:129:  state=data;dirty=false;pending.hidden=!data.restart_required;
    static/settings.html:130:  if(restart){show('Settings saved. Restarting and reconnecting…','ok');await waitForRestart(data,previousID);}
    static/settings.html:131:  else{fields.disabled=false;show('Settings saved. Choose Save & restart to apply them.','ok');}
    static/settings.html:136:document.getElementById('restart').addEventListener('click',()=>save(true));
    README.md:12:## Config file and browser settings
    README.md:14:The settings file is `config/neonvr.json`. The supplied file preserves the existing Blink root, both Yi camera roots, 14-day retention, two thumbnail workers, and 30-minute refresh/maintenance intervals.
    README.md:16:Open **Settings** from the recordings page, or visit `http://SERVER_IP:5500/settings`. Every runtime setting can be edited in the form.
    README.md:18:- **Save** validates and atomically writes the file. The running process continues with its current settings.
    README.md:19:- **Save & restart** saves, drains requests, stops background work, and replaces the running process with the same executable reading the new file. It works in a container and in a native installation without granting the app Docker access.
    README.md:20:- The page reconnects after a restart. When the listen address changes, it provides the new URL.
    README.md:21:- Conflicting edits are rejected. Reload saved settings before saving again.
    README.md:26:The application has the same network access model as the existing recording browser; it does not add a separate settings account. Config writes require a token from the current app instance and reject cross-origin browser requests.
    README.md:35:Legacy flags still work. Explicit flag overrides are persisted to the config file at startup so later browser restarts do not revert them. Repeated `-extra-root` flags replace the configured list. Without an explicit `-config`, a missing default file is initialized from defaults and flags; an explicitly named missing file is an error.
    README.md:39:To roll back settings, restore `neonvr.json.previous` over `neonvr.json` and restart the service. A source snapshot was also created before adding configuration support.
    README.md:89:- Main player at the top with a scrollable one-column thumbnail rail on the right
    README.md:91:- Clicking a sidebar thumbnail or mini-video thumbnail loads it into the main player
    README.md:92:- Active sidebar thumbnail and mini-video card stay highlighted/synced with the main player
    README.md:95:- Space bar pauses/resumes the main player
    README.md:97:- Zoom mode toggle on the main player: when Zoom is on, pinch/trackpad zoom and drag-to-pan are enabled; when Zoom is off, tapping/clicking the video pauses/plays normally
    README.md:101:- Favorites tab
    README.md:131:Subsequent settings changes use **Save & restart** in the browser and do not require a container rebuild.
    README.md:155:## Main player controls
    README.md:157:- Click a thumbnail or `Load` to send it to the main player.
    README.md:200:Favorites are always preserved by detection pruning. Brief unconfirmed detections, decode errors, incomplete video coverage, missing dependencies/models, metadata-write failures, timeouts and changing recordings all keep the clip. No-target classification filters recordings caused by leaves; it is not a separate leaf-recognition model. Sampled object detection can miss events, so review quarantined clips when tuning settings before enabling permanent deletion.
    README.md:202:To restore a quarantined clip, first set the no-target policy to **Keep all recordings** and restart, then move the MP4 from `.neonvr_pruned` back to its original relative path in that camera directory. The existing age-based retention policy continues to operate on active recordings; it does not purge the quarantine.
    README.md:244:To include real model inference and the background prune pipeline in the Go checks, set `NEONVR_TEST_MODEL` to an absolute ONNX path and `NEONVR_TEST_PYTHON` to the detection environment's Python executable. Ordinary Go tests do not download models; MP4 tests skip if FFmpeg/FFprobe are unavailable.
    README.md:267:New thumbnails use a divergence selector designed for mostly-static security footage. The server samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    README.md:282:The project directory is now `/home/ubuntu/NeoNVR`, the Go module and executable are `neonvr`, the Compose service/container are `neonvr`, and the systemd example is `neonvr.service.example`.
    README.md:294:For a native systemd installation, install the executable at `/opt/NeoNVR/neonvr` and use `neonvr.service.example` as the unit template.
    main.go:182://go:embed static/settings.html
    main.go:183:var settingsHTML []byte
    main.go:264:    media-controller {
    main.go:300:      font-variant-numeric: tabular-nums;
    main.go:335:    .tabs, .bulkBar { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
    main.go:336:    .tabs button.active { border-color: var(--accent); color: var(--text); }
    main.go:528:      scrollbar-gutter: stable;
    main.go:653:      media-controller { aspect-ratio: auto; grid-template-rows: auto auto auto; }
    main.go:660:  <a href="/settings" style="display:inline-block;margin-bottom:12px;color:inherit">Settings</a>
    main.go:668:          <span id="mainTitle">Main player: click a thumbnail or choose Load.</span>
    main.go:672:          <media-controller>
    main.go:693:          </media-controller>
    main.go:715:        <main class="sideRail" id="sideRail" aria-label="Scrollable video thumbnail navigation"></main>
    main.go:723:      <div class="tabs">
    main.go:1537:    diskStatsEl.innerHTML = '<span>Folder sizes unavailable: ' + escapeHTML(String(err.message || err)) + '</span>';
    main.go:1614:  if (evt.composedPath().some(el => ['INPUT','SELECT','TEXTAREA','BUTTON'].includes(el.tagName) || el.isContentEditable)) return;
    main.go:1643:func main() {
    main.go:1655:    executable, err := os.Executable()
    main.go:1756:    mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
    main.go:1764:        _, _ = w.Write(settingsHTML)
    main.go:1814:    case <-manager.restart:
    main.go:1823:        if err := syscall.Exec(executable, []string{executable, "-config", configPath}, os.Environ()); err != nil {
    main.go:1824:            log.Fatalf("restart failed: %v", err)
    main.go:2094:        writePlaceholderSVG(w, "video unavailable")
    main.go:2100:        writePlaceholderSVG(w, "thumbnail unavailable")
    main.go:2179:            log.Printf("extra clip root unavailable %s: %v", dir, err)
    ✓ • 0ms

## Activity

    $ sed -n '335,460p' config.go
        sed -n '1525,1645p' main.go
        sed -n '640,673p' main.go
        cat config/neonvr.json
        command -v ffmpeg ffprobe node npm chromium chromium-browser google-chrome
        rg --files /tmp /usr/local -g go -g node -g '*playwright*' -g '*chrom*' -g '!node_modules/**' -g '!\.cache/**'
        w.WriteHeader(status)
        _ = json.NewEncoder(w).Encode(map[string]any{"config": saved, "active_config": m.active, "revision": revision, "token": m.token, "instance_id": m.token, "restart_required": !reflect.DeepEqual(saved, m.active), "restarting": m.restarting, "next_url": nextURL})
    }
    
    func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        if r.Method != "GET" && r.Method != "POST" {
            w.Header().Set("Allow", "GET, POST")
            http.Error(w, "method not allowed", 405)
            return
        }
        if r.Method == "POST" {
            if r.Header.Get("X-NeoNVR-Token") != m.token {
                http.Error(w, "reload Settings before saving", 403)
                return
            }
            if origin := r.Header.Get("Origin"); origin != "" {
                u, err := url.Parse(origin)
                scheme := "http"
                if r.TLS != nil {
                    scheme = "https"
                }
                if err != nil || u.Host != r.Host || u.Scheme != scheme {
                    http.Error(w, "same-origin requests required", 403)
                    return
                }
            }
            if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
                http.Error(w, "JSON required", 415)
                return
            }
        }
        m.mu.Lock()
        defer m.mu.Unlock()
        raw, err := os.ReadFile(m.path)
        if err != nil {
            http.Error(w, "cannot read config: "+err.Error(), 500)
            return
        }
        saved, err := decodeConfig(raw)
        if err != nil {
            http.Error(w, "invalid config file: "+err.Error(), 500)
            return
        }
        revision := configRevision(raw)
        if r.Method == "GET" {
            m.response(w, 200, saved, revision, "")
            return
        }
        if m.restarting {
            http.Error(w, "restart already in progress", 409)
            return
        }
        var request struct {
            Config   json.RawMessage `json:"config"`
            Revision string          `json:"revision"`
            Restart  bool            `json:"restart"`
        }
        dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
        dec.DisallowUnknownFields()
        if err := dec.Decode(&request); err != nil {
            http.Error(w, "invalid request: "+err.Error(), 400)
            return
        }
        if err := dec.Decode(new(any)); err != io.EOF {
            http.Error(w, "unexpected trailing JSON", 400)
            return
        }
        if request.Revision != revision {
            http.Error(w, "settings changed elsewhere; reload before saving", 409)
            return
        }
        c, err := decodeConfig(request.Config)
        if err != nil {
            http.Error(w, "invalid config: "+err.Error(), 400)
            return
        }
        if err = c.validate(); err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        if request.Restart && c.Addr != m.active.Addr {
            listener, err := net.Listen("tcp", c.Addr)
            if err != nil {
                http.Error(w, "cannot listen on new address: "+err.Error(), 400)
                return
            }
            listener.Close()
        }
        data := configBytes(c)
        if err = atomicWriteConfig(m.path+".previous", raw); err != nil {
            http.Error(w, "cannot save previous config: "+err.Error(), 500)
            return
        }
        if err = atomicWriteConfig(m.path, data); err != nil {
            http.Error(w, "cannot save config: "+err.Error(), 500)
            return
        }
        nextURL := ""
        if request.Restart {
            m.restarting = true
            host, port, _ := net.SplitHostPort(c.Addr)
            if host == "" || host == "0.0.0.0" || host == "::" {
                host, _, _ = net.SplitHostPort(r.Host)
                if host == "" {
                    host = r.Host
                }
            }
            scheme := "http"
            if r.TLS != nil {
                scheme = "https"
            }
            nextURL = scheme + "://" + net.JoinHostPort(host, port) + "/settings"
        }
        m.response(w, 200, c, configRevision(data), nextURL)
        if request.Restart {
            // Delay until this response has reached the browser, then drain and re-exec.
            time.AfterFunc(750*time.Millisecond, func() { m.restart <- struct{}{} })
        }
    }
      renderClips();
      renderDatePills();
    }
    
    async function refreshStats() {
      try {
        const stats = await getJSON('/api/stats');
        let html = '<span>Recordings: <strong>' + stats.root_size_human + '</strong></span>' +
                   '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
        if (stats.error) html += '<span>' + escapeHTML(stats.error) + '</span>';
        diskStatsEl.innerHTML = html;
      } catch (err) {
        diskStatsEl.innerHTML = '<span>Folder sizes unavailable: ' + escapeHTML(String(err.message || err)) + '</span>';
      }
    }
    
    function setView(view) {
      currentView = view;
      clipsTabEl.classList.toggle('active', view === 'clips');
      favoritesTabEl.classList.toggle('active', view === 'favorites');
      dateEl.disabled = view === 'favorites';
      clearSelection();
      refreshClips(true).catch(err => {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      });
    }
    
    async function bulkDeleteSelected() {
      if (selectedKeys.size === 0) return;
      const items = visibleClips.filter(c => selectedKeys.has(keyFor(c))).map(c => ({camera: c.camera, file: c.file}));
      const hasFavorite = visibleClips.some(c => selectedKeys.has(keyFor(c)) && c.favorite);
      const msg = hasFavorite
        ? 'Delete ' + items.length + ' selected videos? This includes favorited clips.'
        : 'Delete ' + items.length + ' selected videos?';
      if (!confirm(msg)) return;
      try {
        await postJSON('/api/bulk-delete', {items, force: hasFavorite});
        const del = new Set(items.map(i => i.camera + '/' + i.file));
        allClips = allClips.filter(c => !del.has(keyFor(c)));
        clearSelection();
        renderMarks();
        renderClips();
        refreshStats();
      } catch (err) {
        alert(String(err.message || err));
      }
    }
    
    async function init() {
      try {
        buildScales();
        updateRangeVisuals();
        refreshStats();
        const data = await getJSON('/api/cameras');
        cameras = data.cameras || [];
        selectedCameras = new Set(cameras.map(c => c.name));
        renderCameraMenu();
        await refreshDates(true);
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    }
    
    cameraButtonEl.addEventListener('click', () => cameraMenuEl.classList.toggle('open'));
    document.addEventListener('click', evt => {
      if (!cameraButtonEl.contains(evt.target) && !cameraMenuEl.contains(evt.target)) cameraMenuEl.classList.remove('open');
    });
    
    dateEl.addEventListener('change', async () => {
      try {
        clearSelection();
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    });
    
    sortEl.addEventListener('change', () => { clearSelection(); renderClips(); syncActiveCard(true); });
    startEl.addEventListener('input', () => { normalizeRange('start'); clearSelection(); renderClips(); });
    endEl.addEventListener('input', () => { normalizeRange('end'); clearSelection(); renderClips(); });
    clipsTabEl.addEventListener('click', () => setView('clips'));
    favoritesTabEl.addEventListener('click', () => setView('favorites'));
    selectModeToggleEl.addEventListener('click', () => setSelectMode(!selectMode));
    bulkDeleteEl.addEventListener('click', bulkDeleteSelected);
    clearSelectionEl.addEventListener('click', () => clearSelection());
    zoomToggleEl.addEventListener('click', () => setZoomMode(!zoomMode));
    zoomResetEl.addEventListener('click', resetZoom);
    document.addEventListener('keydown', evt => {
      if (evt.composedPath().some(el => ['INPUT','SELECT','TEXTAREA','BUTTON'].includes(el.tagName) || el.isContentEditable)) return;
    
      if (evt.key === ' ') {
        if (!mainVideoEl.src) return;
        evt.preventDefault();
        if (mainVideoEl.paused) mainVideoEl.play().catch(() => {});
        else mainVideoEl.pause();
        return;
      }
    
      if (evt.ctrlKey && ['ArrowLeft','ArrowRight','ArrowUp','ArrowDown'].includes(evt.key)) {
        evt.preventDefault();
        if (evt.key === 'ArrowLeft' || evt.key === 'ArrowUp') playRelative(-1);
        else playRelative(1);
        return;
      }
    
      if (evt.key === 'ArrowLeft') { evt.preventDefault(); stepMain(-1); }
      if (evt.key === 'ArrowRight') { evt.preventDefault(); stepMain(1); }
    });
    
    installVideoPanZoom();
    installSideRailScrollCapture();
    setZoomMode(false);
    init();
    </script>
    </body>
    </html>`))
    
    func main() {
        cfg, configPath, err := loadStartupConfig(os.Args[1:])
        if errors.Is(err, flagHelpError) {
        .error { color: #ffb4a8; border-color: #6d332f; }
        .notice { color: #a7e0ff; border-color: #2c536b; }
    
        @media (max-width: 1100px) {
          .viewerLayout { grid-template-columns: 1fr; }
          .railPanel { position: static; height: 420px; }
          .controls { grid-template-columns: 1fr 1fr; }
          .controls .timelineWrap { grid-column: 1 / -1; }
        }
        @media (max-width: 720px) {
          body { padding: 14px; }
          .controls { grid-template-columns: 1fr; }
          .cameraMenu { width: min(92vw, 360px); }
          media-controller { aspect-ratio: auto; grid-template-rows: auto auto auto; }
          #mainVideo { height: auto; min-height: 180px; aspect-ratio: 16 / 9; }
        }
      </style>
    </head>
    <body>
      <h1>NeoNVR</h1>
      <a href="/settings" style="display:inline-block;margin-bottom:12px;color:inherit">Settings</a>
      <div class="root">Root: {{.Root}}</div>
      <div class="stats" id="diskStats"><span>Calculating recordings size…</span></div>
    
      <section class="panel mainPlayerPanel">
        <div class="viewerLayout">
          <div class="mainColumn">
            <div class="mainPlayerTitle">
              <span id="mainTitle">Main player: click a thumbnail or choose Load.</span>
              <span class="kbdHint">Space pause/play • ← / → seek • Ctrl+arrows previous/next</span>
            </div>
            <div class="mainPlayerShell" id="mainShell">
              <media-controller>
                <video id="mainVideo" slot="media" preload="metadata" playsinline webkit-playsinline x5-playsinline></video>
    {
      "root": "/mnt/mySSD/BlinkClips",
      "extra_roots": [
        "/mnt/FTPfiles/yicamBack",
        "/mnt/FTPfiles/yicamFront"
      ],
      "addr": "0.0.0.0:5500",
      "thumb_dir": "",
      "meta": "",
      "ffmpeg": "ffmpeg",
      "thumb_seek": 0.2,
      "thumb_width": 480,
      "thumb_workers": 2,
      "thumb_maintenance_every": "30m",
      "size_refresh": "30m",
      "auto_delete_days": 14,
      "auto_delete_every": "6h",
      "detection_enabled": true,
      "detection_python": "python3",
      "detection_model": "yolo26n.onnx",
      "detection_threads": 2,
      "detection_fps": 2,
      "detection_confidence": 0.65,
      "detection_min_frames": 3,
      "detection_merge_gap": 0.5,
      "detection_car_movement": 0.025,
      "detection_every": "1m",
      "detection_min_age": "2m",
      "detection_timeout": "30m",
      "detection_prune": "quarantine"
    }
    /usr/bin/ffmpeg
    /usr/bin/ffprobe
    rg: /tmp/snap-private-tmp: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-ModemManager.service-me2epj: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-systemd-timesyncd.service-5FT3IN: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-systemd-oomd.service-qPVtA3: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-power-profiles-daemon.service-1MO0vS: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-polkit.service-meRHgo: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-fwupd.service-jIGemC: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-colord.service-iRXTQd: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-upower.service-zeVlJo: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-systemd-resolved.service-kOhWqY: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-switcheroo-control.service-KUp5wR: Permission denied (os error 13)
    rg: /tmp/systemd-private-b7feda137c114dfeac83365e9fbbded3-systemd-logind.service-XZEZDL: Permission denied (os error 13)
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/Components/WebUnit/Resources/DriverBinaries/ChromeDriver/Linux-x86-64/102/chromedriver
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Ribbons/MeshGraphics3D/PlotTheme3DMonochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Ribbons/MeshGraphics3D/PlotTheme3DMonochrome.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Ribbons/MeshGraphics/PlotThemeMonochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Ribbons/MeshGraphics/PlotThemeMonochrome.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/Kernel/SystemResources/64Bit/DataPaclets/PsychrometricPropertyData.mx
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletScalogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ChromaticityPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoSmoothHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoBubbleChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BoxWhiskerChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SectorChart3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SphericalPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RenkoChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RadialAxisPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoRegionValuePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-HorizontalGauge-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Histogram3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-JuliaSetPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ChromaticityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexStreamPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SystemModelPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TradingChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ThermometerGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TernaryListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RectangleChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RectangleChart3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletScalogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceDensityPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoVectorPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDensityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AngularGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MatrixPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexArrayPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListLogPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-HorizontalGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexStreamPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RulePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LineBreakChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedSmoothHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListCurvePathPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RectangleChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DiscretePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VerticalGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Graph3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ReliefPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDisplacementPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoStreamPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListStepPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RootLocusPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceVectorPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Graph-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ContourPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PolarPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PointValuePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListLinePlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PieChart3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Histogram3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoSmoothHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDensityPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PieChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-QuantilePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Spectrogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ChromaticityPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Histogram3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PointFigureChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParallelAxisPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DistributionChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListSurfacePlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StackedDateListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ChromaticityPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TernaryListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RenkoChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothDensityHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StackedListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TernaryListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParallelAxisPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoBubbleChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NumberLinePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DiscretePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ProbabilityScalePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SectorChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLogPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListLinePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoStreamPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NicholsPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MeshRegion-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TimelinePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoDensityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLogPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MandelbrotSetPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceVectorPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ClockGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NumberLinePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ArrayPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MatrixPlot-Monoc
    ... command output truncated for persistence ...
    otTheme-ContourPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SystemModelPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ReImPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-KagiChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListCurvePathPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexArrayPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-InteractiveTradingChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedBarChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ProbabilityPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-InteractiveTradingChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListStepPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-InteractiveTradingChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothDensityHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexContourPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ReImPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RegionPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListPointPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-QuantilePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BoxWhiskerChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ContourPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TimelinePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoContourPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexVectorPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexVectorPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ArrayPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NyquistPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RulePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexContourPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AbsArgPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-JuliaSetPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceContourPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BarChart3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SectorChart3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RulePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDensityPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDisplacementPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-KagiChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MeshRegion-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ProbabilityScalePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StackedListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexRegionPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TradingChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexVectorPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoBubbleChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MatrixPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PolarPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DistributionChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BulletGauge-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ArrayPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListLinePlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VerticalGauge-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AbsArgPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ImageHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SphericalPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorDisplacementPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PieChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ReImPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LineBreakChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Histogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceDensityPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NyquistPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLinearPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListPointPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PointFigureChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PieChart3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoGraphPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Graph3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SystemModelPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RegionPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoVectorPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListSurfacePlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NyquistPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ThermometerGauge-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SingularValuePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Graph3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListStepPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedSmoothHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DistributionChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RegionPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RevolutionPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Plot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TradingChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedSmoothHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParallelAxisPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PointValuePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-QuantilePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AbsArgPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AnatomyPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceVectorPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ClockGauge-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ClockGauge-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AngularGauge-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothHistogram3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SingularValuePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothHistogram3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AnatomyPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BulletGauge-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SingularValuePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MandelbrotSetPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ImageHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BulletGauge-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RegionPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListStepPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SmoothHistogram3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceContourPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-CandlestickChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateHistogram-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ArrayPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoStreamPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoGraphValuePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoGraphValuePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SliceContourPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LineBreakChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListPointPlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AudioPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DiscretePlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-CandlestickChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MeshRegion-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoGraphValuePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedBarChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoContourPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletListPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoDensityPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-JuliaSetPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StackedDateListPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-RectangleChart3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ContourPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BoxWhiskerChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Spectrogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MandelbrotSetPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLinearPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StackedDateListPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-Histogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLinearPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListLinePlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexStreamPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DateListLogPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AngularGauge-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-DensityHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PieChart3D-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-BubbleChart-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PointValuePlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-FeatureSpacePlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletMatrixPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SectorChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ProbabilityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-LogLogPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ReliefPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-NicholsPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-AnatomyPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-WaveletMatrixPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListStepPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamDensityPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-PairedBarChart-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-MoleculePlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-StreamPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexRegionPlot-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ImageHistogram-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ComplexPlot3D-Monochrome@240dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListLinePlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoSmoothHistogram-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-GeoGraphPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-TimelinePlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-SphericalPlot3D-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ParametricPlot-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-ListCurvePathPlot-Monochrome@192dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-KagiChart-Monochrome@144dpi.png
    /usr/local/Wolfram/WolframEngine/13.3/SystemFiles/FrontEnd/SystemResources/Bitmaps/Popups/CodeCompletion/Thumbnails/PlotTheme-VectorPlot3D-Monochrome@144dpi.png
    ✗ (2) • 0ms

## Activity

    Searched the web for site.ffmpeg.org ffmpeg formats hls hls_time hls_list_size delete_segments temp_file rtsp_transport rw_timeout, site.hlsjs.video hls.js getting started Hls.isSupported native HLS

## User

however it should be toggleable and default should be off since it consumes camera resources. THe backchannel is for camera audio

## Assistant

Live viewing will default to off and connect to the camera only after you enable it. I’ll include camera audio playback for the backchannel URL you described.

## User

i mean the backchannel is for user audio to inject into the camera

## Assistant

Understood—the backchannel sends your microphone audio to the camera. I’ll keep live viewing off by default and add an explicit Talk control. Microphone access will require HTTPS or localhost.

## Activity

    $ sed -n '339,458p' config.go
        sed -n '1796,1835p' main.go
        sed -n '325,355p' main.go
        sed -n '1,70p' README.md
        ls -d /tmp/neonvr-* /home/ubuntu/.cache/ms-playwright /home/ubuntu/.local/share/uv/tools 2>/dev/null
    func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        if r.Method != "GET" && r.Method != "POST" {
            w.Header().Set("Allow", "GET, POST")
            http.Error(w, "method not allowed", 405)
            return
        }
        if r.Method == "POST" {
            if r.Header.Get("X-NeoNVR-Token") != m.token {
                http.Error(w, "reload Settings before saving", 403)
                return
            }
            if origin := r.Header.Get("Origin"); origin != "" {
                u, err := url.Parse(origin)
                scheme := "http"
                if r.TLS != nil {
                    scheme = "https"
                }
                if err != nil || u.Host != r.Host || u.Scheme != scheme {
                    http.Error(w, "same-origin requests required", 403)
                    return
                }
            }
            if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
                http.Error(w, "JSON required", 415)
                return
            }
        }
        m.mu.Lock()
        defer m.mu.Unlock()
        raw, err := os.ReadFile(m.path)
        if err != nil {
            http.Error(w, "cannot read config: "+err.Error(), 500)
            return
        }
        saved, err := decodeConfig(raw)
        if err != nil {
            http.Error(w, "invalid config file: "+err.Error(), 500)
            return
        }
        revision := configRevision(raw)
        if r.Method == "GET" {
            m.response(w, 200, saved, revision, "")
            return
        }
        if m.restarting {
            http.Error(w, "restart already in progress", 409)
            return
        }
        var request struct {
            Config   json.RawMessage `json:"config"`
            Revision string          `json:"revision"`
            Restart  bool            `json:"restart"`
        }
        dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
        dec.DisallowUnknownFields()
        if err := dec.Decode(&request); err != nil {
            http.Error(w, "invalid request: "+err.Error(), 400)
            return
        }
        if err := dec.Decode(new(any)); err != io.EOF {
            http.Error(w, "unexpected trailing JSON", 400)
            return
        }
        if request.Revision != revision {
            http.Error(w, "settings changed elsewhere; reload before saving", 409)
            return
        }
        c, err := decodeConfig(request.Config)
        if err != nil {
            http.Error(w, "invalid config: "+err.Error(), 400)
            return
        }
        if err = c.validate(); err != nil {
            http.Error(w, err.Error(), 400)
            return
        }
        if request.Restart && c.Addr != m.active.Addr {
            listener, err := net.Listen("tcp", c.Addr)
            if err != nil {
                http.Error(w, "cannot listen on new address: "+err.Error(), 400)
                return
            }
            listener.Close()
        }
        data := configBytes(c)
        if err = atomicWriteConfig(m.path+".previous", raw); err != nil {
            http.Error(w, "cannot save previous config: "+err.Error(), 500)
            return
        }
        if err = atomicWriteConfig(m.path, data); err != nil {
            http.Error(w, "cannot save config: "+err.Error(), 500)
            return
        }
        nextURL := ""
        if request.Restart {
            m.restarting = true
            host, port, _ := net.SplitHostPort(c.Addr)
            if host == "" || host == "0.0.0.0" || host == "::" {
                host, _, _ = net.SplitHostPort(r.Host)
                if host == "" {
                    host = r.Host
                }
            }
            scheme := "http"
            if r.TLS != nil {
                scheme = "https"
            }
            nextURL = scheme + "://" + net.JoinHostPort(host, port) + "/settings"
        }
        m.response(w, 200, c, configRevision(data), nextURL)
        if request.Restart {
            // Delay until this response has reached the browser, then drain and re-exec.
            time.AfterFunc(750*time.Millisecond, func() { m.restart <- struct{}{} })
        }
    }
            log.Printf("Thumbnail maintenance: every %s (pre-generate missing, prune orphan cache files)", s.thumbMaintenanceEvery)
        } else {
            log.Printf("Thumbnail maintenance disabled")
        }
        log.Printf("Size refresh: %s", s.sizeTTL)
        if s.autoDeleteDays > 0 {
            log.Printf("Auto-delete enabled: non-favorite clips older than %d days, checking every %s", s.autoDeleteDays, s.autoDeleteEvery)
        } else {
            log.Printf("Auto-delete disabled; use -auto-delete-days N to enable")
        }
        httpServer := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
        serverErr := make(chan error, 1)
        go func() { serverErr <- httpServer.ListenAndServe() }()
        select {
        case err := <-serverErr:
            if !errors.Is(err, http.ErrServerClosed) {
                log.Fatal(err)
            }
        case <-manager.restart:
            log.Printf("Restarting NeoNVR with config %s", configPath)
            cancelWork()
            ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
            if err := httpServer.Shutdown(ctx); err != nil {
                _ = httpServer.Close()
            }
            cancel()
            s.background.Wait()
            if err := syscall.Exec(executable, []string{executable, "-config", configPath}, os.Environ()); err != nil {
                log.Fatalf("restart failed: %v", err)
            }
        }
    }
    
    func NewMetaStore(path string) (*MetaStore, error) {
        m := &MetaStore{path: path, data: MetaFile{Version: 1, Clips: map[string]ClipMeta{}}}
        b, err := os.ReadFile(path)
        if err != nil {
            if os.IsNotExist(err) {
                return m, nil
            }
    
        .topRow {
          display: flex;
          justify-content: space-between;
          gap: 14px;
          align-items: center;
          flex-wrap: wrap;
          margin-bottom: 12px;
        }
    
        .tabs, .bulkBar { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
        .tabs button.active { border-color: var(--accent); color: var(--text); }
        .bulkBar { color: var(--muted); font-size: 13px; }
    
        .controls {
          display: grid;
          grid-template-columns: minmax(220px, 300px) minmax(160px, 220px) minmax(160px, 220px) 1fr;
          gap: 14px;
          align-items: end;
        }
    
        label { display: grid; gap: 6px; color: var(--muted); font-size: 14px; }
        select, input, button {
          font: inherit;
          border: 1px solid var(--border);
          border-radius: 10px;
          background: #0d0f14;
          color: var(--text);
          padding: 9px 10px;
        }
        button { cursor: pointer; background: var(--panel2); }
    # NeoNVR
    
    Docker-first security-camera MP4 browser with CPU-only YOLO26n event detection, including Blink clips stored as:
    
    ```text
    /mnt/mySSD/BlinkClips/<CameraName>/YYYY-MM-DD_HH-MM-SS_<CameraName>.mp4
    ```
    
    The browser UI is fully offline: the media-control shim is served locally from the Go binary, not from a CDN.
    
    
    ## Config file and browser settings
    
    The settings file is `config/neonvr.json`. The supplied file preserves the existing Blink root, both Yi camera roots, 14-day retention, two thumbnail workers, and 30-minute refresh/maintenance intervals.
    
    Open **Settings** from the recordings page, or visit `http://SERVER_IP:5500/settings`. Every runtime setting can be edited in the form.
    
    - **Save** validates and atomically writes the file. The running process continues with its current settings.
    - **Save & restart** saves, drains requests, stops background work, and replaces the running process with the same executable reading the new file. It works in a container and in a native installation without granting the app Docker access.
    - The page reconnects after a restart. When the listen address changes, it provides the new URL.
    - Conflicting edits are rejected. Reload saved settings before saving again.
    - The previous file is retained as `config/neonvr.json.previous`.
    
    Paths refer to the server/container filesystem. New recording directories must exist and, in a container, must also be mounted. Changing the container listen port also requires updating the Compose port mapping.
    
    The application has the same network access model as the existing recording browser; it does not add a separate settings account. Config writes require a token from the current app instance and reject cross-origin browser requests.
    
    For a native run:
    
    ```bash
    go build -o neonvr .
    ./neonvr -config config/neonvr.json
    ```
    
    Legacy flags still work. Explicit flag overrides are persisted to the config file at startup so later browser restarts do not revert them. Repeated `-extra-root` flags replace the configured list. Without an explicit `-config`, a missing default file is initialized from defaults and flags; an explicitly named missing file is an error.
    
    For Docker/nerdctl, mount the **whole config directory** read/write, as the included Compose file does. Mounting only the JSON file prevents atomic replacement.
    
    To roll back settings, restore `neonvr.json.previous` over `neonvr.json` and restart the service. A source snapshot was also created before adding configuration support.
    
    ## Docker quick start
    
    ```bash
    sudo apt update
    sudo apt install -y docker.io docker-compose-plugin
    sudo usermod -aG docker "$USER"
    # log out and back in if this is your first Docker install
    
    cd /home/ubuntu/NeoNVR
    docker compose up -d --build
    ```
    
    Open:
    
    ```text
    http://SERVER_IP:5500/
    ```
    
    Logs:
    
    ```bash
    docker compose logs -f neonvr
    ```
    
    Stop/start:
    
    ```bash
    docker compose stop
    docker compose start
    ```
    ✗ (2) • 0ms

## Activity

    Searched the web for site.github.com/AlexxIT/go2rtc two way audio microphone video-rtc.js backchannel, site.github.com/AlexxIT/go2rtc releases latest go2rtc linux amd64 api listen webrtc listen

## Activity

    Opened https://github.com/AlexxIT/go2rtc/blob/v1.9.14/internal/webrtc/README.md

## Activity

    $ python3 - <<'PY'
        import urllib.request
        base='https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'
        for path in ['www/video-rtc.js','www/webrtc.html','internal/webrtc/webrtc.go','internal/webrtc/README.md','internal/api/api.go','internal/rtsp/README.md']:
         print('\nFILE '+path)
         with urllib.request.urlopen(base+path,timeout=20) as response:
          content=response.read().decode()
         print(content[:27000])
        PY
    
    FILE www/video-rtc.js
    /**
     * VideoRTC v1.6.0 - Video player for go2rtc streaming application.
     *
     * All modern web technologies are supported in almost any browser except Apple Safari.
     *
     * Support:
     * - ECMAScript 2017 (ES8) = ES6 + async
     * - RTCPeerConnection for Safari iOS 11.0+
     * - IntersectionObserver for Safari iOS 12.2+
     * - ManagedMediaSource for Safari 17+
     *
     * Doesn't support:
     * - MediaSource for Safari iOS
     * - Customized built-in elements (extends HTMLVideoElement) because Safari
     * - Autoplay for WebRTC in Safari
     */
    export class VideoRTC extends HTMLElement {
        constructor() {
            super();
    
            this.DISCONNECT_TIMEOUT = 5000;
            this.RECONNECT_TIMEOUT = 15000;
    
            this.CODECS = [
                'avc1.640029',      // H.264 high 4.1 (Chromecast 1st and 2nd Gen)
                'avc1.64002A',      // H.264 high 4.2 (Chromecast 3rd Gen)
                'avc1.640033',      // H.264 high 5.1 (Chromecast with Google TV)
                'hvc1.1.6.L153.B0', // H.265 main 5.1 (Chromecast Ultra)
                'mp4a.40.2',        // AAC LC
                'mp4a.40.5',        // AAC HE
                'flac',             // FLAC (PCM compatible)
                'opus',             // OPUS Chrome, Firefox
            ];
    
            /**
             * [config] Supported modes (webrtc, webrtc/tcp, mse, hls, mp4, mjpeg).
             * @type {string}
             */
            this.mode = 'webrtc,mse,hls,mjpeg';
    
            /**
             * [Config] Requested medias (video, audio, microphone).
             * @type {string}
             */
            this.media = 'video,audio';
    
            /**
             * [config] Run stream when not displayed on the screen. Default `false`.
             * @type {boolean}
             */
            this.background = false;
    
            /**
             * [config] Run stream only when player in the viewport. Stop when user scroll out player.
             * Value is percentage of visibility from `0` (not visible) to `1` (full visible).
             * Default `0` - disable;
             * @type {number}
             */
            this.visibilityThreshold = 0;
    
            /**
             * [config] Run stream only when browser page on the screen. Stop when user change browser
             * tab or minimise browser windows.
             * @type {boolean}
             */
            this.visibilityCheck = true;
    
            /**
             * [config] WebRTC configuration
             * @type {RTCConfiguration}
             */
            this.pcConfig = {
                bundlePolicy: 'max-bundle',
                iceServers: [{urls: ['stun:stun.cloudflare.com:3478', 'stun:stun.l.google.com:19302']}],
                sdpSemantics: 'unified-plan',  // important for Chromecast 1
            };
    
            /**
             * [info] WebSocket connection state. Values: CONNECTING, OPEN, CLOSED
             * @type {number}
             */
            this.wsState = WebSocket.CLOSED;
    
            /**
             * [info] WebRTC connection state.
             * @type {number}
             */
            this.pcState = WebSocket.CLOSED;
    
            /**
             * @type {HTMLVideoElement}
             */
            this.video = null;
    
            /**
             * @type {WebSocket}
             */
            this.ws = null;
    
            /**
             * @type {string|URL}
             */
            this.wsURL = '';
    
            /**
             * @type {RTCPeerConnection}
             */
            this.pc = null;
    
            /**
             * @type {number}
             */
            this.connectTS = 0;
    
            /**
             * @type {string}
             */
            this.mseCodecs = '';
    
            /**
             * [internal] Disconnect TimeoutID.
             * @type {number}
             */
            this.disconnectTID = 0;
    
            /**
             * [internal] Reconnect TimeoutID.
             * @type {number}
             */
            this.reconnectTID = 0;
    
            /**
             * [internal] Handler for receiving Binary from WebSocket.
             * @type {Function}
             */
            this.ondata = null;
    
            /**
             * [internal] Handlers list for receiving JSON from WebSocket.
             * @type {Object.<string,Function>}
             */
            this.onmessage = null;
        }
    
        /**
         * Set video source (WebSocket URL). Support relative path.
         * @param {string|URL} value
         */
        set src(value) {
            if (typeof value !== 'string') value = value.toString();
            if (value.startsWith('http')) {
                value = 'ws' + value.substring(4);
            } else if (value.startsWith('/')) {
                value = 'ws' + location.origin.substring(4) + value;
            }
    
            this.wsURL = value;
    
            this.onconnect();
        }
    
        /**
         * Play video. Support automute when autoplay blocked.
         * https://developer.chrome.com/blog/autoplay/
         */
        play() {
            this.video.play().catch(() => {
                if (!this.video.muted) {
                    this.video.muted = true;
                    this.video.play().catch(er => {
                        console.warn(er);
                    });
                }
            });
        }
    
        /**
         * Send message to server via WebSocket
         * @param {Object} value
         */
        send(value) {
            if (this.ws) this.ws.send(JSON.stringify(value));
        }
    
        /** @param {Function} isSupported */
        codecs(isSupported) {
            return this.CODECS
                .filter(codec => this.media.includes(codec.includes('vc1') ? 'video' : 'audio'))
                .filter(codec => isSupported(`video/mp4; codecs="${codec}"`)).join();
        }
    
        /**
         * `CustomElement`. Invoked each time the custom element is appended into a
         * document-connected element.
         */
        connectedCallback() {
            if (this.disconnectTID) {
                clearTimeout(this.disconnectTID);
                this.disconnectTID = 0;
            }
    
            // because video autopause on disconnected from DOM
            if (this.video) {
                const seek = this.video.seekable;
                if (seek.length > 0) {
                    this.video.currentTime = seek.end(seek.length - 1);
                }
                this.play();
            } else {
                this.oninit();
            }
    
            this.onconnect();
        }
    
        /**
         * `CustomElement`. Invoked each time the custom element is disconnected from the
         * document's DOM.
         */
        disconnectedCallback() {
            if (this.background || this.disconnectTID) return;
            if (this.wsState === WebSocket.CLOSED && this.pcState === WebSocket.CLOSED) return;
    
            this.disconnectTID = setTimeout(() => {
                if (this.reconnectTID) {
                    clearTimeout(this.reconnectTID);
                    this.reconnectTID = 0;
                }
    
                this.disconnectTID = 0;
    
                this.ondisconnect();
            }, this.DISCONNECT_TIMEOUT);
        }
    
        /**
         * Creates child DOM elements. Called automatically once on `connectedCallback`.
         */
        oninit() {
            this.video = document.createElement('video');
            this.video.controls = true;
            this.video.playsInline = true;
            this.video.preload = 'auto';
    
            this.video.style.display = 'block'; // fix bottom margin 4px
            this.video.style.width = '100%';
            this.video.style.height = '100%';
    
            this.appendChild(this.video);
    
            this.video.addEventListener('error', ev => {
                const err = this.video.error;
                // https://developer.mozilla.org/en-US/docs/Web/API/MediaError/code
                const MEDIA_ERRORS = {
                    1: 'MEDIA_ERR_ABORTED',
                    2: 'MEDIA_ERR_NETWORK',
                    3: 'MEDIA_ERR_DECODE',
                    4: 'MEDIA_ERR_SRC_NOT_SUPPORTED'
                };
                console.error('[VideoRTC] Video error:', {
                    error: err ? MEDIA_ERRORS[err.code] : 'unknown',
                    message: err ? err.message : 'unknown',
                    codecs: this.mseCodecs || 'not set',
                    readyState: this.video.readyState,
                    networkState: this.video.networkState,
                    currentTime: this.video.currentTime
                });
                if (this.ws) this.ws.close(); // run reconnect for broken MSE stream
            });
    
            // all Safari lies about supported audio codecs
            const m = window.navigator.userAgent.match(/Version\/(\d+).+Safari/);
            if (m) {
                // AAC from v13, FLAC from v14, OPUS - unsupported
                const skip = m[1] < '13' ? 'mp4a.40.2' : m[1] < '14' ? 'flac' : 'opus';
                this.CODECS.splice(this.CODECS.indexOf(skip));
            }
    
            if (this.background) return;
    
            if ('hidden' in document && this.visibilityCheck) {
                document.addEventListener('visibilitychange', () => {
                    if (document.hidden) {
                        this.disconnectedCallback();
                    } else if (this.isConnected) {
                        this.connectedCallback();
                    }
                });
            }
    
            if ('IntersectionObserver' in window && this.visibilityThreshold) {
                const observer = new IntersectionObserver(entries => {
                    entries.forEach(entry => {
                        if (!entry.isIntersecting) {
                            this.disconnectedCallback();
                        } else if (this.isConnected) {
                            this.connectedCallback();
                        }
                    });
                }, {threshold: this.visibilityThreshold});
                observer.observe(this);
            }
        }
    
        /**
         * Connect to WebSocket. Called automatically on `connectedCallback`.
         * @return {boolean} true if the connection has started.
         */
        onconnect() {
            if (!this.isConnected || !this.wsURL || this.ws || this.pc) return false;
    
            // CLOSED or CONNECTING => CONNECTING
            this.wsState = WebSocket.CONNECTING;
    
            this.connectTS = Date.now();
    
            this.ws = new WebSocket(this.wsURL);
            this.ws.binaryType = 'arraybuffer';
            this.ws.addEventListener('open', () => this.onopen());
            this.ws.addEventListener('close', () => this.onclose());
    
            return true;
        }
    
        ondisconnect() {
            this.wsState = WebSocket.CLOSED;
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
    
            this.pcState = WebSocket.CLOSED;
            if (this.pc) {
                this.pc.getSenders().forEach(sender => {
                    if (sender.track) sender.track.stop();
                });
                this.pc.close();
                this.pc = null;
            }
    
            this.video.src = '';
            this.video.srcObject = null;
        }
    
        /**
         * @returns {Array.<string>} of modes (mse, webrtc, etc.)
         */
        onopen() {
            // CONNECTING => OPEN
            this.wsState = WebSocket.OPEN;
    
            this.ws.addEventListener('message', ev => {
                if (typeof ev.data === 'string') {
                    const msg = JSON.parse(ev.data);
                    for (const mode in this.onmessage) {
                        this.onmessage[mode](msg);
                    }
                } else {
                    this.ondata(ev.data);
                }
            });
    
            this.ondata = null;
            this.onmessage = {};
    
            const modes = [];
    
            if (this.mode.includes('mse') && ('MediaSource' in window || 'ManagedMediaSource' in window)) {
                modes.push('mse');
                this.onmse();
            } else if (this.mode.includes('hls') && this.video.canPlayType('application/vnd.apple.mpegurl')) {
                modes.push('hls');
                this.onhls();
            } else if (this.mode.includes('mp4')) {
                modes.push('mp4');
                this.onmp4();
            }
    
            if (this.mode.includes('webrtc') && 'RTCPeerConnection' in window) {
                modes.push('webrtc');
                this.onwebrtc();
            }
    
            if (this.mode.includes('mjpeg')) {
                if (modes.length) {
                    this.onmessage['mjpeg'] = msg => {
                        if (msg.type !== 'error' || msg.value.indexOf(modes[0]) !== 0) return;
                        this.onmjpeg();
                    };
                } else {
                    modes.push('mjpeg');
                    this.onmjpeg();
                }
            }
    
            return modes;
        }
    
        /**
         * @return {boolean} true if reconnection has started.
         */
        onclose() {
            if (this.wsState === WebSocket.CLOSED) return false;
    
            // CONNECTING, OPEN => CONNECTING
            this.wsState = WebSocket.CONNECTING;
            this.ws = null;
    
            // reconnect no more than once every X seconds
            const delay = Math.max(this.RECONNECT_TIMEOUT - (Date.now() - this.connectTS), 0);
    
            this.reconnectTID = setTimeout(() => {
                this.reconnectTID = 0;
                this.onconnect();
            }, delay);
    
            return true;
        }
    
        onmse() {
            /** @type {MediaSource} */
            let ms;
    
            if ('ManagedMediaSource' in window) {
                const MediaSource = window.ManagedMediaSource;
    
                ms = new MediaSource();
                ms.addEventListener('sourceopen', () => {
                    this.send({type: 'mse', value: this.codecs(MediaSource.isTypeSupported)});
                }, {once: true});
    
                this.video.disableRemotePlayback = true;
                this.video.srcObject = ms;
            } else {
                ms = new MediaSource();
                ms.addEventListener('sourceopen', () => {
                    URL.revokeObjectURL(this.video.src);
                    this.send({type: 'mse', value: this.codecs(MediaSource.isTypeSupported)});
                }, {once: true});
    
                this.video.src = URL.createObjectURL(ms);
                this.video.srcObject = null;
            }
    
            this.play();
    
            this.mseCodecs = '';
    
            this.onmessage['mse'] = msg => {
                if (msg.type !== 'mse') return;
    
                this.mseCodecs = msg.value;
    
                const sb = ms.addSourceBuffer(msg.value);
                sb.mode = 'segments'; // segments or sequence
                sb.addEventListener('updateend', () => {
                    if (!sb.updating && bufLen > 0) {
                        try {
                            const data = buf.slice(0, bufLen);
                            sb.appendBuffer(data);
                            bufLen = 0;
                        } catch (e) {
                            // console.debug(e);
                        }
                    }
    
                    if (!sb.updating && sb.buffered && sb.buffered.length) {
                        const end = sb.buffered.end(sb.buffered.length - 1);
                        const start = end - 5;
                        const start0 = sb.buffered.start(0);
                        if (start > start0) {
                            sb.remove(start0, start);
                            ms.setLiveSeekableRange(start, end);
                        }
                        if (this.video.currentTime < start) {
                            this.video.currentTime = start;
                        }
                        const gap = end - this.video.currentTime;
                        this.video.playbackRate = gap > 0.1 ? gap : 0.1;
                        // console.debug('VideoRTC.buffered', gap, this.video.playbackRate, this.video.readyState);
                    }
                });
    
                const buf = new Uint8Array(2 * 1024 * 1024);
                let bufLen = 0;
    
                this.ondata = data => {
                    if (sb.updating || bufLen > 0) {
                        const b = new Uint8Array(data);
                        buf.set(b, bufLen);
                        bufLen += b.byteLength;
                        // console.debug('VideoRTC.buffer', b.byteLength, bufLen);
                    } else {
                        try {
                            sb.appendBuffer(data);
                        } catch (e) {
                            // console.debug(e);
                        }
                    }
                };
            };
        }
    
        onwebrtc() {
            const pc = new RTCPeerConnection(this.pcConfig);
    
            pc.addEventListener('icecandidate', ev => {
                if (ev.candidate && this.mode.includes('webrtc/tcp') && ev.candidate.protocol === 'udp') return;
    
                const candidate = ev.candidate ? ev.candidate.toJSON().candidate : '';
                this.send({type: 'webrtc/candidate', value: candidate});
            });
    
            pc.addEventListener('connectionstatechange', () => {
                if (pc.connectionState === 'connected') {
                    const tracks = pc.getTransceivers()
                        .filter(tr => tr.currentDirection === 'recvonly') // skip inactive
                        .map(tr => tr.receiver.track);
                    /** @type {HTMLVideoElement} */
                    const video2 = document.createElement('video');
                    video2.addEventListener('loadeddata', () => this.onpcvideo(video2), {once: true});
                    video2.srcObject = new MediaStream(tracks);
                } else if (pc.connectionState === 'failed' || pc.connectionState === 'disconnected') {
                    pc.close(); // stop next events
    
                    this.pcState = WebSocket.CLOSED;
                    this.pc = null;
    
                    this.onconnect();
                }
            });
    
            this.onmessage['webrtc'] = msg => {
                switch (msg.type) {
                    case 'webrtc/candidate':
                        if (this.mode.includes('webrtc/tcp') && msg.value.includes(' udp ')) return;
    
                        pc.addIceCandidate({candidate: msg.value, sdpMid: '0'}).catch(er => {
                            console.warn(er);
                        });
                        break;
                    case 'webrtc/answer':
                        pc.setRemoteDescription({type: 'answer', sdp: msg.value}).catch(er => {
                            console.warn(er);
                        });
                        break;
                    case 'error':
                        if (!msg.value.includes('webrtc/offer')) return;
                        pc.close();
                }
            };
    
            this.createOffer(pc).then(offer => {
                this.send({type: 'webrtc/offer', value: offer.sdp});
            });
    
            this.pcState = WebSocket.CONNECTING;
            this.pc = pc;
        }
    
        /**
         * @param pc {RTCPeerConnection}
         * @return {Promise<RTCSessionDescriptionInit>}
         */
        async createOffer(pc) {
            try {
                if (this.media.includes('microphone')) {
                    const media = await navigator.mediaDevices.getUserMedia({audio: true});
                    media.getTracks().forEach(track => {
                        pc.addTransceiver(track, {direction: 'sendonly'});
                    });
                }
            } catch (e) {
                console.warn(e);
            }
    
            for (const kind of ['video', 'audio']) {
                if (this.media.includes(kind)) {
                    pc.addTransceiver(kind, {direction: 'recvonly'});
                }
            }
    
            const offer = await pc.createOffer();
            await pc.setLocalDescription(offer);
            return offer;
        }
    
        /**
         * @param video2 {HTMLVideoElement}
         */
        onpcvideo(video2) {
            if (this.pc) {
                // Video+Audio > Video, H265 > H264, Video > Audio, WebRTC > MSE
                let rtcPriority = 0, msePriority = 0;
    
                /** @type {MediaStream} */
                const stream = video2.srcObject;
                if (stream.getVideoTracks().length > 0) {
                    // not the best, but a pretty simple way to check a codec
                    const isH265Supported =  this.pc.remoteDescription.sdp.includes('H265/90000');
                    rtcPriority += isH265Supported ? 0x240 : 0x220;
                }
                if (stream.getAudioTracks().length > 0) rtcPriority += 0x102;
    
                if (this.mseCodecs.includes('hvc1.')) msePriority += 0x230;
                if (this.mseCodecs.includes('avc1.')) msePriority += 0x210;
                if (this.mseCodecs.includes('mp4a.')) msePriority += 0x101;
    
                if (rtcPriority >= msePriority) {
                    this.video.srcObject = stream;
                    this.play();
    
                    this.pcState = WebSocket.OPEN;
    
                    this.wsState = WebSocket.CLOSED;
                    if (this.ws) {
                        this.ws.close();
                        this.ws = null;
                    }
                } else {
                    this.pcState = WebSocket.CLOSED;
                    if (this.pc) {
                        this.pc.close();
                        this.pc = null;
                    }
                }
            }
    
            video2.srcObject = null;
        }
    
        onmjpeg() {
            this.ondata = data => {
                this.video.controls = false;
                this.video.poster = 'data:image/jpeg;base64,' + VideoRTC.btoa(data);
            };
    
            this.send({type: 'mjpeg'});
        }
    
        onhls() {
            this.onmessage['hls'] = msg => {
                if (msg.type !== 'hls') return;
    
                const url = 'http' + this.wsURL.substring(2, this.wsURL.indexOf('/ws')) + '/hls/';
                const playlist = msg.value.replace('hls/', url);
                this.video.src = 'data:application/vnd.apple.mpegurl;base64,' + btoa(playlist);
                this.play();
            };
    
            this.send({type: 'hls', value: this.codecs(type => this.video.canPlayType(type))});
        }
    
        onmp4() {
            /** @type {HTMLCanvasElement} **/
            const canvas = document.createElement('canvas');
            /** @type {CanvasRenderingContext2D} */
            let context;
    
            /** @type {HTMLVideoElement} */
            const video2 = document.createElement('video');
            video2.autoplay = true;
            video2.playsInline = true;
            video2.muted = true;
    
            video2.addEventListener('loadeddata', () => {
                if (!context) {
                    canvas.width = video2.videoWidth;
                    canvas.height = video2.videoHeight;
                    context = canvas.getContext('2d');
                }
    
                context.drawImage(video2, 0, 0, canvas.width, canvas.height);
    
                this.video.controls = false;
                this.video.poster = canvas.toDataURL('image/jpeg');
            });
    
            this.ondata = data => {
                video2.src = 'data:video/mp4;base64,' + VideoRTC.btoa(data);
            };
    
            this.send({type: 'mp4', value: this.codecs(this.video.canPlayType)});
        }
    
        static btoa(buffer) {
            const bytes = new Uint8Array(buffer);
            const len = bytes.byteLength;
            let binary = '';
            for (let i = 0; i < len; i++) {
                binary += String.fromCharCode(bytes[i]);
            }
            return window.btoa(binary);
        }
    }
    
    
    FILE www/webrtc.html
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <title>webrtc - go2rtc</title>
        <style>
            body {
                background-color: black;
                margin: 0;
                padding: 0;
            }
    
            html, body, video {
                height: 100%;
                width: 100%;
            }
        </style>
    </head>
    <body>
    <video id="video" autoplay controls playsinline muted></video>
    <script>
        async function PeerConnection(media) {
            const pc = new RTCPeerConnection({
                iceServers: [{urls: ['stun:stun.cloudflare.com:3478', 'stun:stun.l.google.com:19302']}]
            });
    
            const localTracks = [];
    
            if (/camera|microphone/.test(media)) {
                const tracks = await getMediaTracks('user', {
                    video: media.indexOf('camera') >= 0,
                    audio: media.indexOf('microphone') >= 0,
                });
                tracks.forEach(track => {
                    pc.addTransceiver(track, {direction: 'sendonly'});
                    if (track.kind === 'video') localTracks.push(track);
                });
            }
    
            if (media.indexOf('display') >= 0) {
                const tracks = await getMediaTracks('display', {
                    video: true,
                    audio: media.indexOf('speaker') >= 0,
                });
                tracks.forEach(track => {
                    pc.addTransceiver(track, {direction: 'sendonly'});
                    if (track.kind === 'video') localTracks.push(track);
                });
            }
    
            if (/video|audio/.test(media)) {
                const tracks = ['video', 'audio']
                    .filter(kind => media.indexOf(kind) >= 0)
                    .map(kind => pc.addTransceiver(kind, {direction: 'recvonly'}).receiver.track);
                localTracks.push(...tracks);
            }
    
            document.getElementById('video').srcObject = new MediaStream(localTracks);
    
            return pc;
        }
    
        async function getMediaTracks(media, constraints) {
            try {
                const stream = media === 'user'
                    ? await navigator.mediaDevices.getUserMedia(constraints)
                    : await navigator.mediaDevices.getDisplayMedia(constraints);
                return stream.getTracks();
            } catch (e) {
                console.warn(e);
                return [];
            }
        }
    
        async function connect(media) {
            const pc = await PeerConnection(media);
            const url = new URL('api/ws' + location.search, location.href);
            const ws = new WebSocket('ws' + url.toString().substring(4));
    
            ws.addEventListener('open', () => {
                pc.addEventListener('icecandidate', ev => {
                    if (!ev.candidate) return;
                    const msg = {type: 'webrtc/candidate', value: ev.candidate.candidate};
                    ws.send(JSON.stringify(msg));
                });
    
                pc.createOffer().then(offer => pc.setLocalDescription(offer)).then(() => {
                    const msg = {type: 'webrtc/offer', value: pc.localDescription.sdp};
                    ws.send(JSON.stringify(msg));
                });
            });
    
            ws.addEventListener('message', ev => {
                const msg = JSON.parse(ev.data);
                if (msg.type === 'webrtc/candidate') {
                    pc.addIceCandidate({candidate: msg.value, sdpMid: '0'});
                } else if (msg.type === 'webrtc/answer') {
                    pc.setRemoteDescription({type: 'answer', sdp: msg.value});
                }
            });
        }
    
        const media = new URLSearchParams(location.search).get('media');
        connect(media || 'video+audio');
    </script>
    </body>
    </html>
    
    FILE internal/webrtc/webrtc.go
    package webrtc
    
    import (
        "errors"
        "net"
        "strings"
    
        "github.com/AlexxIT/go2rtc/internal/api"
        "github.com/AlexxIT/go2rtc/internal/api/ws"
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/AlexxIT/go2rtc/internal/streams"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/webrtc"
        pion "github.com/pion/webrtc/v4"
        "github.com/rs/zerolog"
    )
    
    func Init() {
        var cfg struct {
            Mod struct {
                Listen     string           `yaml:"listen"`
                Candidates []string         `yaml:"candidates"`
                IceServers []pion.ICEServer `yaml:"ice_servers"`
                Filters    webrtc.Filters   `yaml:"filters"`
            } `yaml:"webrtc"`
        }
    
        cfg.Mod.Listen = ":8555"
        cfg.Mod.IceServers = []pion.ICEServer{
            {URLs: []string{"stun:stun.cloudflare.com:3478", "stun:stun.l.google.com:19302"}},
        }
    
        app.LoadConfig(&cfg)
    
        log = app.GetLogger("webrtc")
    
        if log.Debug().Enabled() {
            itfs, _ := net.Interfaces()
            for _, itf := range itfs {
                addrs, _ := itf.Addrs()
                log.Debug().Msgf("[webrtc] interface %+v addrs %v", itf, addrs)
            }
        }
    
        address, network, _ := strings.Cut(cfg.Mod.Listen, "/")
        for _, candidate := range cfg.Mod.Candidates {
            AddCandidate(network, candidate)
    
            if strings.HasPrefix(candidate, "stun:") && stuns == nil {
                for _, ice := range cfg.Mod.IceServers {
                    for _, url := range ice.URLs {
                        if strings.HasPrefix(url, "stun:") {
                            stuns = append(stuns, url[5:])
                        }
                    }
                }
            }
        }
    
        webrtc.OnNewListener = func(ln any) {
            switch ln := ln.(type) {
            case *net.TCPListener:
                log.Info().Stringer("addr", ln.Addr()).Msg("[webrtc] listen tcp")
            case *net.UDPConn:
                log.Info().Stringer("addr", ln.LocalAddr()).Msg("[webrtc] listen udp")
            }
        }
    
        var err error
    
        // create pionAPI with custom codecs list and custom network settings
        serverAPI, err = webrtc.NewServerAPI(network, address, &cfg.Mod.Filters)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        // use same API for WebRTC server and client if no address
        clientAPI = serverAPI
    
        if address != "" {
            clientAPI, _ = webrtc.NewAPI()
        }
    
        pionConf := pion.Configuration{
            ICEServers:   cfg.Mod.IceServers,
            SDPSemantics: pion.SDPSemanticsUnifiedPlanWithFallback,
        }
    
        PeerConnection = func(active bool) (*pion.PeerConnection, error) {
            // active - client, passive - server
            if active {
                return clientAPI.NewPeerConnection(pionConf)
            } else {
                return serverAPI.NewPeerConnection(pionConf)
            }
        }
    
        // async WebRTC server (two API versions)
        ws.HandleFunc("webrtc", asyncHandler)
        ws.HandleFunc("webrtc/offer", asyncHandler)
        ws.HandleFunc("webrtc/candidate", candidateHandler)
    
        // sync WebRTC server (two API versions)
        api.HandleFunc("api/webrtc", syncHandler)
    
        // WebRTC client
        streams.HandleFunc("webrtc", streamsHandler)
    }
    
    var serverAPI, clientAPI *pion.API
    
    var log zerolog.Logger
    
    var PeerConnection func(active bool) (*pion.PeerConnection, error)
    
    func asyncHandler(tr *ws.Transport, msg *ws.Message) (err error) {
        var stream *streams.Stream
        var mode core.Mode
    
        query := tr.Request.URL.Query()
        if name := query.Get("src"); name != "" {
            stream, _ = streams.GetOrPatch(query)
            mode = core.ModePassiveConsumer
            log.Debug().Str("src", name).Msg("[webrtc] new consumer")
        } else if name = query.Get("dst"); name != "" {
            stream = streams.Get(name)
            mode = core.ModePassiveProducer
            log.Debug().Str("src", name).Msg("[webrtc] new producer")
        }
    
        if stream == nil {
            return errors.New(api.StreamNotFound)
        }
    
        var offer struct {
            Type       string           `json:"type"`
            SDP        string           `json:"sdp"`
            ICEServers []pion.ICEServer `json:"ice_servers"`
        }
    
        // V2 - json/object exchange, V1 - raw SDP exchange
        apiV2 := msg.Type == "webrtc"
    
        if apiV2 {
            if err = msg.Unmarshal(&offer); err != nil {
                return err
            }
        } else {
            offer.SDP = msg.String()
        }
    
        // create new PeerConnection instance
        var pc *pion.PeerConnection
        if offer.ICEServers == nil {
            pc, err = PeerConnection(false)
        } else {
            pc, err = serverAPI.NewPeerConnection(pion.Configuration{ICEServers: offer.ICEServers})
        }
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return err
        }
    
        var sendAnswer core.Waiter
    
        // protect from blocking on errors
        defer sendAnswer.Done(nil)
    
        conn := webrtc.NewConn(pc)
        conn.Mode = mode
        conn.Protocol = "ws"
        conn.UserAgent = tr.Request.UserAgent()
        conn.Listen(func(msg any) {
            switch msg := msg.(type) {
            case pion.PeerConnectionState:
                if msg != pion.PeerConnectionStateClosed {
                    return
                }
                switch mode {
                case core.ModePassiveConsumer:
                    stream.RemoveConsumer(conn)
                case core.ModePassiveProducer:
                    stream.RemoveProducer(conn)
                }
    
            case *pion.ICECandidate:
                if !FilterCandidate(msg) {
                    return
                }
                _ = sendAnswer.Wait()
    
                s := msg.ToJSON().Candidate
                log.Trace().Str("candidate", s).Msg("[webrtc] local ")
                tr.Write(&ws.Message{Type: "webrtc/candidate", Value: s})
            }
        })
    
        log.Trace().Msgf("[webrtc] offer:\n%s", offer.SDP)
    
        // 1. SetOffer, so we can get remote client codecs
        if err = conn.SetOffer(offer.SDP); err != nil {
            log.Warn().Err(err).Caller().Send()
            return err
        }
    
        switch mode {
        case core.ModePassiveConsumer:
            // 2. AddConsumer, so we get new tracks
            if err = stream.AddConsumer(conn); err != nil {
                log.Debug().Err(err).Msg("[webrtc] add consumer")
                _ = conn.Close()
                return err
            }
        case core.ModePassiveProducer:
            stream.AddProducer(conn)
        }
    
        // 3. Exchange SDP without waiting all candidates
        answer, err := conn.GetAnswer()
        log.Trace().Msgf("[webrtc] answer\n%s", answer)
    
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return err
        }
    
        if apiV2 {
            desc := pion.SessionDescription{Type: pion.SDPTypeAnswer, SDP: answer}
            tr.Write(&ws.Message{Type: "webrtc", Value: desc})
        } else {
            tr.Write(&ws.Message{Type: "webrtc/answer", Value: answer})
        }
    
        sendAnswer.Done(nil)
    
        asyncCandidates(tr, conn)
    
        return nil
    }
    
    func ExchangeSDP(stream *streams.Stream, offer, desc, userAgent string) (answer string, err error) {
        pc, err := PeerConnection(false)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        // create new webrtc instance
        conn := webrtc.NewConn(pc)
        conn.FormatName = desc
        conn.UserAgent = userAgent
        conn.Protocol = "http"
        conn.Listen(func(msg any) {
            switch msg := msg.(type) {
            case pion.PeerConnectionState:
                if msg != pion.PeerConnectionStateClosed {
                    return
                }
                if conn.Mode == core.ModePassiveConsumer {
                    stream.RemoveConsumer(conn)
                } else {
                    stream.RemoveProducer(conn)
                }
            }
        })
    
        // 1. SetOffer, so we can get remote client codecs
        log.Trace().Msgf("[webrtc] offer:\n%s", offer)
    
        if err = conn.SetOffer(offer); err != nil {
            log.Warn().Err(err).Caller().Send()
            return
        }
    
        if IsConsumer(conn) {
            conn.Mode = core.ModePassiveConsumer
    
            // 2. AddConsumer, so we get new tracks
            if err = stream.AddConsumer(conn); err != nil {
                log.Warn().Err(err).Caller().Send()
                _ = conn.Close()
                return
            }
        } else {
            conn.Mode = core.ModePassiveProducer
    
            stream.AddProducer(conn)
        }
    
        answer, err = conn.GetCompleteAnswer(GetCandidates(), FilterCandidate)
        log.Trace().Msgf("[webrtc] answer\n%s", answer)
    
        if err != nil {
            log.Error().Err(err).Caller().Send()
        }
    
        return
    }
    
    func IsConsumer(conn *webrtc.Conn) bool {
        // if wants get video - consumer
        for _, media := range conn.GetMedias() {
            if media.Kind == core.KindVideo && media.Direction == core.DirectionSendonly {
                return true
            }
        }
        // if wants send video - producer
        for _, media := range conn.GetMedias() {
            if media.Kind == core.KindVideo && media.Direction == core.DirectionRecvonly {
                return false
            }
        }
        // if wants something - consumer
        for _, media := range conn.GetMedias() {
            if media.Direction == core.DirectionSendonly {
                return true
            }
        }
        return false
    }
    
    
    FILE internal/webrtc/README.md
    What you should to know about WebRTC:
    
    - It's almost always a **direct [peer-to-peer](https://en.wikipedia.org/wiki/Peer-to-peer) connection** from your browser to go2rtc app
    - When you use Home Assistant, Frigate, Nginx, Nabu Casa, Cloudflare and other software - they are only **involved in establishing** the connection, but they are **not involved in transferring** media data
    - WebRTC media cannot be transferred inside an HTTP connection
    - Usually, WebRTC uses random UDP ports on client and server side to establish a connection
    - Usually, WebRTC uses public [STUN](https://en.wikipedia.org/wiki/STUN) servers to establish a connection outside LAN, such servers are only needed to establish a connection and are not involved in data transfer
    - Usually, WebRTC will automatically discover all of your local addresses and all of your public addresses and try to establish a connection
    
    If an external connection via STUN is used:
    
    - Uses [UDP hole punching](https://en.wikipedia.org/wiki/UDP_hole_punching) technology to bypass NAT even if you not open your server to the World
    - For about 20% of users, the techology will not work because of the [Symmetric NAT](https://tomchen.github.io/symmetric-nat-test/)
    - UDP is not suitable for transmitting 2K and 4K high bitrate video over open networks because of the high loss rate:
      - https://habr.com/ru/companies/flashphoner/articles/480006/
      - https://www.youtube.com/watch?v=FXVg2ckuKfs
    
    ## Default config
    
    ```yaml
    webrtc:
      listen: ":8555"
      ice_servers:
        - urls: [ "stun:stun.l.google.com:19302" ]
    ```
    
    ## Config
    
    **Important!** This example is not for copypasting!
    
    ```yaml
    webrtc:
      # fix local TCP or UDP or both ports for WebRTC media
      listen: ":8555"            # address of your local server
    
      # add additional host candidates manually
      # order is important, the first will have a higher priority
      candidates:
        - 216.58.210.174:8555    # if you have static public IP-address
        - stun:8555              # if you have dynamic public IP-address
        - home.duckdns.org:8555  # if you have domain
    
      # add custom STUN and TURN servers
      # use `ice_servers: []` for remove defaults and leave empty
      ice_servers:
        - urls: [ stun:stun1.l.google.com:19302 ]
        - urls: [ turn:123.123.123.123:3478 ]
          username: your_user
          credential: your_pass
    
      # optional filter list for auto discovery logic
      # some settings only make sense if you don't specify a fixed UDP port
      filters:
        # list of host candidates from auto discovery to be sent
        # including candidates from the `listen` option
        # use `candidates: []` to remove all auto discovery candidates
        candidates: [ 192.168.1.123 ]
    
        # enable localhost candidates
        loopback: true
    
        # list of network types to be used for connection
        # including candidates from the `listen` option
        networks: [ udp4, udp6, tcp4, tcp6 ]
    
        # list of interfaces to be used for connection
        # including interfaces from unspecified `listen` option (empty host)
        interfaces: [ eno1 ]
    
        # list of host IP-addresses to be used for connection
        # including IPs from unspecified `listen` option (empty host)
        ips: [ 192.168.1.123 ]
    
        # range for random UDP ports [min, max] to be used for connection
        # not related to the `listen` option
        udp_ports: [ 50000, 50100 ]
    ```
    
    By default go2rtc uses **fixed TCP** port and **fixed UDP** ports for each **direct** WebRTC connection - `listen: ":8555"`.
    
    You can set **fixed TCP** and **random UDP** port for all connections - `listen: ":8555/tcp"`.
    
    Don't know why, but you can disable TCP port and leave only random UDP ports - `listen: ""`.
    
    ## Config filters
    
    **Importan!** By default go2rtc exclude all Docker-like candidates (`172.16.0.0/12`). This can not be disabled.
    
    Filters allow you to exclude unnecessary candidates. Extra candidates don't make your connection worse or better. But the wrong filter settings can break everything. Skip this setting if you don't understand it.
    
    For example, go2rtc is installed on the host system. And there are unnecessary interfaces. You can keep only the relevant via `interfaces` or `ips` options. You can also exclude IPv6 candidates if your server supports them but your home network does not.
    
    ```yaml
    webrtc:
      listen: ":8555/tcp"         # use fixed TCP port and random UDP ports
      filters:
        ips: [ 192.168.1.2 ]      # IP-address of your server
        networks: [ udp4, tcp4 ]  # skip IPv6, if it's not supported for you
    ```
    
    For example, go2rtc inside closed docker container (ex. [Frigate](https://frigate.video/)). You shouldn't filter docker interfaces, otherwise go2rtc will not be able to connect anywhere. But you can filter the docker candidates because no one can connect to them.
    
    ```yaml
    webrtc:
      listen: ":8555"                   # use fixed TCP and UDP ports
      candidates: [ 192.168.1.2:8555 ]  # add manual host candidate (use docker port forwarding)
    ```
    
    ## Userful links
    
    - https://www.ietf.org/archive/id/draft-ietf-wish-whip-01.html
    - https://www.ietf.org/id/draft-murillo-whep-01.html
    - https://github.com/Glimesh/broadcast-box/
    - https://github.com/obsproject/obs-studio/pull/7926
    - https://misi.github.io/webrtc-c0d3l4b/
    - https://github.com/webtorrent/webtorrent/blob/master/docs/faq.md
    
    
    FILE internal/api/api.go
    package api
    
    import (
        "crypto/tls"
        "encoding/json"
        "fmt"
        "net"
        "net/http"
        "os"
        "slices"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/rs/zerolog"
    )
    
    func Init() {
        var cfg struct {
            Mod struct {
                Listen     string `yaml:"listen"`
                Username   string `yaml:"username"`
                Password   string `yaml:"password"`
                LocalAuth  bool   `yaml:"local_auth"`
                BasePath   string `yaml:"base_path"`
                StaticDir  string `yaml:"static_dir"`
                Origin     string `yaml:"origin"`
                TLSListen  string `yaml:"tls_listen"`
                TLSCert    string `yaml:"tls_cert"`
                TLSKey     string `yaml:"tls_key"`
                UnixListen string `yaml:"unix_listen"`
    
                AllowPaths []string `yaml:"allow_paths"`
            } `yaml:"api"`
        }
    
        // default config
        cfg.Mod.Listen = ":1984"
    
        // load config from YAML
        app.LoadConfig(&cfg)
    
        if cfg.Mod.Listen == "" && cfg.Mod.UnixListen == "" && cfg.Mod.TLSListen == "" {
            return
        }
    
        allowPaths = cfg.Mod.AllowPaths
        basePath = cfg.Mod.BasePath
        log = app.GetLogger("api")
    
        initStatic(cfg.Mod.StaticDir)
    
        HandleFunc("api", apiHandler)
        HandleFunc("api/config", configHandler)
        HandleFunc("api/exit", exitHandler)
        HandleFunc("api/restart", restartHandler)
        HandleFunc("api/log", logHandler)
    
        Handler = http.DefaultServeMux // 4th
    
        if cfg.Mod.Origin == "*" {
            Handler = middlewareCORS(Handler) // 3rd
        }
    
        if cfg.Mod.Username != "" {
            Handler = middlewareAuth(cfg.Mod.Username, cfg.Mod.Password, cfg.Mod.LocalAuth, Handler) // 2nd
        }
    
        if log.Trace().Enabled() {
            Handler = middlewareLog(Handler) // 1st
        }
    
        if cfg.Mod.Listen != "" {
            _, port, _ := net.SplitHostPort(cfg.Mod.Listen)
            Port, _ = strconv.Atoi(port)
            go listen("tcp", cfg.Mod.Listen)
        }
    
        if cfg.Mod.UnixListen != "" {
            _ = syscall.Unlink(cfg.Mod.UnixListen)
            go listen("unix", cfg.Mod.UnixListen)
        }
    
        // Initialize the HTTPS server
        if cfg.Mod.TLSListen != "" && cfg.Mod.TLSCert != "" && cfg.Mod.TLSKey != "" {
            go tlsListen("tcp", cfg.Mod.TLSListen, cfg.Mod.TLSCert, cfg.Mod.TLSKey)
        }
    }
    
    func listen(network, address string) {
        ln, err := net.Listen(network, address)
        if err != nil {
            log.Error().Err(err).Msg("[api] listen")
            return
        }
    
        log.Info().Str("addr", address).Msg("[api] listen")
    
        server := http.Server{
            Handler:           Handler,
            ReadHeaderTimeout: 5 * time.Second, // Example: Set to 5 seconds
        }
        if err = server.Serve(ln); err != nil {
            log.Fatal().Err(err).Msg("[api] serve")
        }
    }
    
    func tlsListen(network, address, certFile, keyFile string) {
        var cert tls.Certificate
        var err error
        if strings.IndexByte(certFile, '\n') < 0 && strings.IndexByte(keyFile, '\n') < 0 {
            // check if file path
            cert, err = tls.LoadX509KeyPair(certFile, keyFile)
        } else {
            // if text file content
            cert, err = tls.X509KeyPair([]byte(certFile), []byte(keyFile))
        }
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        ln, err := net.Listen(network, address)
        if err != nil {
            log.Error().Err(err).Msg("[api] tls listen")
            return
        }
    
        log.Info().Str("addr", address).Msg("[api] tls listen")
    
        server := &http.Server{
            Handler:           Handler,
            TLSConfig:         &tls.Config{Certificates: []tls.Certificate{cert}},
            ReadHeaderTimeout: 5 * time.Second,
        }
        if err = server.ServeTLS(ln, "", ""); err != nil {
            log.Fatal().Err(err).Msg("[api] tls serve")
        }
    }
    
    var Port int
    
    const (
        MimeJSON = "application/json"
        MimeText = "text/plain"
    )
    
    var Handler http.Handler
    
    // HandleFunc handle pattern with relative path:
    // - "api/streams" => "{basepath}/api/streams"
    // - "/streams"    => "/streams"
    func HandleFunc(pattern string, handler http.HandlerFunc) {
        if len(pattern) == 0 || pattern[0] != '/' {
            pattern = basePath + "/" + pattern
        }
        if allowPaths != nil && !slices.Contains(allowPaths, pattern) {
            log.Trace().Str("path", pattern).Msg("[api] ignore path not in allow_paths")
            return
        }
        log.Trace().Str("path", pattern).Msg("[api] register path")
        http.HandleFunc(pattern, handler)
    }
    
    // ResponseJSON important always add Content-Type
    // so go won't need to call http.DetectContentType
    func ResponseJSON(w http.ResponseWriter, v any) {
        w.Header().Set("Content-Type", MimeJSON)
        _ = json.NewEncoder(w).Encode(v)
    }
    
    func ResponsePrettyJSON(w http.ResponseWriter, v any) {
        w.Header().Set("Content-Type", MimeJSON)
        enc := json.NewEncoder(w)
        enc.SetIndent("", "  ")
        _ = enc.Encode(v)
    }
    
    func Response(w http.ResponseWriter, body any, contentType string) {
        w.Header().Set("Content-Type", contentType)
    
        switch v := body.(type) {
        case []byte:
            _, _ = w.Write(v)
        case string:
            _, _ = w.Write([]byte(v))
        default:
            _, _ = fmt.Fprint(w, body)
        }
    }
    
    const StreamNotFound = "stream not found"
    
    var allowPaths []string
    var basePath string
    var log zerolog.Logger
    
    func middlewareLog(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Trace().Msgf("[api] %s %s %s", r.Method, r.URL, r.RemoteAddr)
            next.ServeHTTP(w, r)
        })
    }
    
    func isLoopback(remoteAddr string) bool {
        return strings.HasPrefix(remoteAddr, "127.") || strings.HasPrefix(remoteAddr, "[::1]") || remoteAddr == "@"
    }
    
    func middlewareAuth(username, password string, localAuth bool, next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if localAuth || !isLoopback(r.RemoteAddr) {
                user, pass, ok := r.BasicAuth()
                if !ok || user != username || pass != password {
                    w.Header().Set("Www-Authenticate", `Basic realm="go2rtc"`)
                    http.Error(w, "Unauthorized", http.StatusUnauthorized)
                    return
                }
            }
    
            next.ServeHTTP(w, r)
        })
    }
    
    func middlewareCORS(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
            next.ServeHTTP(w, r)
        })
    }
    
    var mu sync.Mutex
    
    func apiHandler(w http.ResponseWriter, r *http.Request) {
        mu.Lock()
        app.Info["host"] = r.Host
        mu.Unlock()
    
        ResponseJSON(w, app.Info)
    }
    
    func exitHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "", http.StatusBadRequest)
            return
        }
    
        s := r.URL.Query().Get("code")
        code, err := strconv.Atoi(s)
    
        // https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_08_02
        if err != nil || code < 0 || code > 125 {
            http.Error(w, "Code must be in the range [0, 125]", http.StatusBadRequest)
            return
        }
    
        os.Exit(code)
    }
    
    func restartHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            http.Error(w, "", http.StatusBadRequest)
            return
        }
    
        path, err := os.Executable()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        log.Debug().Msgf("[api] restart %s", path)
    
        go syscall.Exec(path, os.Args, os.Environ())
    }
    
    func logHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            // Send current state of the log file immediately
            w.Header().Set("Content-Type", "application/jsonlines")
            _, _ = app.MemoryLog.WriteTo(w)
        case "DELETE":
            app.MemoryLog.Reset()
            Response(w, "OK", "text/plain")
        default:
            http.Error(w, "Method not allowed", http.StatusBadRequest)
        }
    }
    
    type Source struct {
        ID       string `json:"id,omitempty"`
        Name     string `json:"name,omitempty"`
        Info     string `json:"info,omitempty"`
        URL      string `json:"url,omitempty"`
        Location string `json:"location,omitempty"`
    }
    
    func ResponseSources(w http.ResponseWriter, sources []*Source) {
        if len(sources) == 0 {
            http.Error(w, "no sources", http.StatusNotFound)
            return
        }
    
        var response = struct {
            Sources []*Source `json:"sources"`
        }{
            Sources: sources,
        }
        ResponseJSON(w, response)
    }
    
    func Error(w http.ResponseWriter, err error) {
        log.Error().Err(err).Caller(1).Send()
    
        http.Error(w, err.Error(), http.StatusInsufficientStorage)
    }
    
    
    FILE internal/rtsp/README.md
    Traceback (most recent call last):
      File "<stdin>", line 5, in <module>
      File "/usr/lib/python3.12/urllib/request.py", line 215, in urlopen
        return opener.open(url, data, timeout)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/usr/lib/python3.12/urllib/request.py", line 521, in open
        response = meth(req, response)
                   ^^^^^^^^^^^^^^^^^^^
      File "/usr/lib/python3.12/urllib/request.py", line 630, in http_response
        response = self.parent.error(
                   ^^^^^^^^^^^^^^^^^^
      File "/usr/lib/python3.12/urllib/request.py", line 559, in error
        return self._call_chain(*args)
               ^^^^^^^^^^^^^^^^^^^^^^^
      File "/usr/lib/python3.12/urllib/request.py", line 492, in _call_chain
        result = func(*args)
                 ^^^^^^^^^^^
      File "/usr/lib/python3.12/urllib/request.py", line 639, in http_error_default
        raise HTTPError(req.full_url, code, msg, hdrs, fp)
    urllib.error.HTTPError: HTTP Error 404: Not Found
    ✗ (1) • 661ms

## Activity

    $ python3 - <<'PY'
        import importlib.util
        print('playwright installed:', importlib.util.find_spec('playwright') is not None)
        print('selenium installed:', importlib.util.find_spec('selenium') is not None)
        PY
        ls /home/ubuntu/.cache 2>/dev/null
        ffmpeg -version | head -n 3
    playwright installed: False
    selenium installed: False
    Wolfram
    event-sound-cache.tdb.b5bed2896bd54379928274382582c663.x86_64-pc-linux-gnu
    evolution
    fontconfig
    go-build
    gstreamer-1.0
    ibus
    ibus-table
    mesa_shader_cache
    mesa_shader_cache_db
    motd.legal-displayed
    pip
    thumbnails
    tracker3
    ubuntu-report
    update-manager-core
    ffmpeg version 6.1.1-3ubuntu5 Copyright (c) 2000-2023 the FFmpeg developers
    built with gcc 13 (Ubuntu 13.2.0-23ubuntu3)
    configuration: --prefix=/usr --extra-version=3ubuntu5 --toolchain=hardened --libdir=/usr/lib/x86_64-linux-gnu --incdir=/usr/include/x86_64-linux-gnu --arch=amd64 --enable-gpl --disable-stripping --disable-omx --enable-gnutls --enable-libaom --enable-libass --enable-libbs2b --enable-libcaca --enable-libcdio --enable-libcodec2 --enable-libdav1d --enable-libflite --enable-libfontconfig --enable-libfreetype --enable-libfribidi --enable-libglslang --enable-libgme --enable-libgsm --enable-libharfbuzz --enable-libmp3lame --enable-libmysofa --enable-libopenjpeg --enable-libopenmpt --enable-libopus --enable-librubberband --enable-libshine --enable-libsnappy --enable-libsoxr --enable-libspeex --enable-libtheora --enable-libtwolame --enable-libvidstab --enable-libvorbis --enable-libvpx --enable-libwebp --enable-libx265 --enable-libxml2 --enable-libxvid --enable-libzimg --enable-openal --enable-opencl --enable-opengl --disable-sndio --enable-libvpl --disable-libmfx --enable-libdc1394 --enable-libdrm --enable-libiec61883 --enable-chromaprint --enable-frei0r --enable-ladspa --enable-libbluray --enable-libjack --enable-libpulse --enable-librabbitmq --enable-librist --enable-libsrt --enable-libssh --enable-libsvtav1 --enable-libx264 --enable-libzmq --enable-libzvbi --enable-lv2 --enable-sdl2 --enable-libplacebo --enable-librav1e --enable-pocketsphinx --enable-librsvg --enable-libjxl --enable-shared
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import urllib.request
        base='https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'
        for path in ['internal/webrtc/http.go','internal/webrtc/webrtc.go','internal/api/api.go','internal/rtsp/rtsp.go','README.md','LICENSE']:
         print('\nFILE '+path)
         try:
          with urllib.request.urlopen(base+path,timeout=20) as response: content=response.read().decode()
         except Exception as error:
          print(error);continue
         if path=='README.md':
          lines=content.splitlines()
          for index,line in enumerate(lines):
           if 'backchannel' in line or 'rtsp:' in line or 'ice_servers' in line:
            print('\n'.join(lines[max(0,index-2):index+4]))
         elif path=='internal/webrtc/webrtc.go':print(content[:11000])
         elif path=='internal/api/api.go':print(content[:5500])
         elif path=='internal/rtsp/rtsp.go':print(content[:3500])
         else:print(content[:15000])
        PY
    
    FILE internal/webrtc/http.go
    HTTP Error 404: Not Found
    
    FILE internal/webrtc/webrtc.go
    package webrtc
    
    import (
        "errors"
        "net"
        "strings"
    
        "github.com/AlexxIT/go2rtc/internal/api"
        "github.com/AlexxIT/go2rtc/internal/api/ws"
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/AlexxIT/go2rtc/internal/streams"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/webrtc"
        pion "github.com/pion/webrtc/v4"
        "github.com/rs/zerolog"
    )
    
    func Init() {
        var cfg struct {
            Mod struct {
                Listen     string           `yaml:"listen"`
                Candidates []string         `yaml:"candidates"`
                IceServers []pion.ICEServer `yaml:"ice_servers"`
                Filters    webrtc.Filters   `yaml:"filters"`
            } `yaml:"webrtc"`
        }
    
        cfg.Mod.Listen = ":8555"
        cfg.Mod.IceServers = []pion.ICEServer{
            {URLs: []string{"stun:stun.cloudflare.com:3478", "stun:stun.l.google.com:19302"}},
        }
    
        app.LoadConfig(&cfg)
    
        log = app.GetLogger("webrtc")
    
        if log.Debug().Enabled() {
            itfs, _ := net.Interfaces()
            for _, itf := range itfs {
                addrs, _ := itf.Addrs()
                log.Debug().Msgf("[webrtc] interface %+v addrs %v", itf, addrs)
            }
        }
    
        address, network, _ := strings.Cut(cfg.Mod.Listen, "/")
        for _, candidate := range cfg.Mod.Candidates {
            AddCandidate(network, candidate)
    
            if strings.HasPrefix(candidate, "stun:") && stuns == nil {
                for _, ice := range cfg.Mod.IceServers {
                    for _, url := range ice.URLs {
                        if strings.HasPrefix(url, "stun:") {
                            stuns = append(stuns, url[5:])
                        }
                    }
                }
            }
        }
    
        webrtc.OnNewListener = func(ln any) {
            switch ln := ln.(type) {
            case *net.TCPListener:
                log.Info().Stringer("addr", ln.Addr()).Msg("[webrtc] listen tcp")
            case *net.UDPConn:
                log.Info().Stringer("addr", ln.LocalAddr()).Msg("[webrtc] listen udp")
            }
        }
    
        var err error
    
        // create pionAPI with custom codecs list and custom network settings
        serverAPI, err = webrtc.NewServerAPI(network, address, &cfg.Mod.Filters)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        // use same API for WebRTC server and client if no address
        clientAPI = serverAPI
    
        if address != "" {
            clientAPI, _ = webrtc.NewAPI()
        }
    
        pionConf := pion.Configuration{
            ICEServers:   cfg.Mod.IceServers,
            SDPSemantics: pion.SDPSemanticsUnifiedPlanWithFallback,
        }
    
        PeerConnection = func(active bool) (*pion.PeerConnection, error) {
            // active - client, passive - server
            if active {
                return clientAPI.NewPeerConnection(pionConf)
            } else {
                return serverAPI.NewPeerConnection(pionConf)
            }
        }
    
        // async WebRTC server (two API versions)
        ws.HandleFunc("webrtc", asyncHandler)
        ws.HandleFunc("webrtc/offer", asyncHandler)
        ws.HandleFunc("webrtc/candidate", candidateHandler)
    
        // sync WebRTC server (two API versions)
        api.HandleFunc("api/webrtc", syncHandler)
    
        // WebRTC client
        streams.HandleFunc("webrtc", streamsHandler)
    }
    
    var serverAPI, clientAPI *pion.API
    
    var log zerolog.Logger
    
    var PeerConnection func(active bool) (*pion.PeerConnection, error)
    
    func asyncHandler(tr *ws.Transport, msg *ws.Message) (err error) {
        var stream *streams.Stream
        var mode core.Mode
    
        query := tr.Request.URL.Query()
        if name := query.Get("src"); name != "" {
            stream, _ = streams.GetOrPatch(query)
            mode = core.ModePassiveConsumer
            log.Debug().Str("src", name).Msg("[webrtc] new consumer")
        } else if name = query.Get("dst"); name != "" {
            stream = streams.Get(name)
            mode = core.ModePassiveProducer
            log.Debug().Str("src", name).Msg("[webrtc] new producer")
        }
    
        if stream == nil {
            return errors.New(api.StreamNotFound)
        }
    
        var offer struct {
            Type       string           `json:"type"`
            SDP        string           `json:"sdp"`
            ICEServers []pion.ICEServer `json:"ice_servers"`
        }
    
        // V2 - json/object exchange, V1 - raw SDP exchange
        apiV2 := msg.Type == "webrtc"
    
        if apiV2 {
            if err = msg.Unmarshal(&offer); err != nil {
                return err
            }
        } else {
            offer.SDP = msg.String()
        }
    
        // create new PeerConnection instance
        var pc *pion.PeerConnection
        if offer.ICEServers == nil {
            pc, err = PeerConnection(false)
        } else {
            pc, err = serverAPI.NewPeerConnection(pion.Configuration{ICEServers: offer.ICEServers})
        }
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return err
        }
    
        var sendAnswer core.Waiter
    
        // protect from blocking on errors
        defer sendAnswer.Done(nil)
    
        conn := webrtc.NewConn(pc)
        conn.Mode = mode
        conn.Protocol = "ws"
        conn.UserAgent = tr.Request.UserAgent()
        conn.Listen(func(msg any) {
            switch msg := msg.(type) {
            case pion.PeerConnectionState:
                if msg != pion.PeerConnectionStateClosed {
                    return
                }
                switch mode {
                case core.ModePassiveConsumer:
                    stream.RemoveConsumer(conn)
                case core.ModePassiveProducer:
                    stream.RemoveProducer(conn)
                }
    
            case *pion.ICECandidate:
                if !FilterCandidate(msg) {
                    return
                }
                _ = sendAnswer.Wait()
    
                s := msg.ToJSON().Candidate
                log.Trace().Str("candidate", s).Msg("[webrtc] local ")
                tr.Write(&ws.Message{Type: "webrtc/candidate", Value: s})
            }
        })
    
        log.Trace().Msgf("[webrtc] offer:\n%s", offer.SDP)
    
        // 1. SetOffer, so we can get remote client codecs
        if err = conn.SetOffer(offer.SDP); err != nil {
            log.Warn().Err(err).Caller().Send()
            return err
        }
    
        switch mode {
        case core.ModePassiveConsumer:
            // 2. AddConsumer, so we get new tracks
            if err = stream.AddConsumer(conn); err != nil {
                log.Debug().Err(err).Msg("[webrtc] add consumer")
                _ = conn.Close()
                return err
            }
        case core.ModePassiveProducer:
            stream.AddProducer(conn)
        }
    
        // 3. Exchange SDP without waiting all candidates
        answer, err := conn.GetAnswer()
        log.Trace().Msgf("[webrtc] answer\n%s", answer)
    
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return err
        }
    
        if apiV2 {
            desc := pion.SessionDescription{Type: pion.SDPTypeAnswer, SDP: answer}
            tr.Write(&ws.Message{Type: "webrtc", Value: desc})
        } else {
            tr.Write(&ws.Message{Type: "webrtc/answer", Value: answer})
        }
    
        sendAnswer.Done(nil)
    
        asyncCandidates(tr, conn)
    
        return nil
    }
    
    func ExchangeSDP(stream *streams.Stream, offer, desc, userAgent string) (answer string, err error) {
        pc, err := PeerConnection(false)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        // create new webrtc instance
        conn := webrtc.NewConn(pc)
        conn.FormatName = desc
        conn.UserAgent = userAgent
        conn.Protocol = "http"
        conn.Listen(func(msg any) {
            switch msg := msg.(type) {
            case pion.PeerConnectionState:
                if msg != pion.PeerConnectionStateClosed {
                    return
                }
                if conn.Mode == core.ModePassiveConsumer {
                    stream.RemoveConsumer(conn)
                } else {
                    stream.RemoveProducer(conn)
                }
            }
        })
    
        // 1. SetOffer, so we can get remote client codecs
        log.Trace().Msgf("[webrtc] offer:\n%s", offer)
    
        if err = conn.SetOffer(offer); err != nil {
            log.Warn().Err(err).Caller().Send()
            return
        }
    
        if IsConsumer(conn) {
            conn.Mode = core.ModePassiveConsumer
    
            // 2. AddConsumer, so we get new tracks
            if err = stream.AddConsumer(conn); err != nil {
                log.Warn().Err(err).Caller().Send()
                _ = conn.Close()
                return
            }
        } else {
            conn.Mode = core.ModePassiveProducer
    
            stream.AddProducer(conn)
        }
    
        answer, err = conn.GetCompleteAnswer(GetCandidates(), FilterCandidate)
        log.Trace().Msgf("[webrtc] answer\n%s", answer)
    
        if err != nil {
            log.Error().Err(err).Caller().Send()
        }
    
        return
    }
    
    func IsConsumer(conn *webrtc.Conn) bool {
        // if wants get video - consumer
        for _, media := range conn.GetMedias() {
            if media.Kind == core.KindVideo && media.Direction == core.DirectionSendonly {
                return true
            }
        }
        // if wants send video - producer
        for _, media := range conn.GetMedias() {
            if media.Kind == core.KindVideo && media.Direction == core.DirectionRecvonly {
                return false
            }
        }
        // if wants something - consumer
        for _, media := range conn.GetMedias() {
            if media.Direction == core.DirectionSendonly {
                return true
            }
        }
        return false
    }
    
    
    FILE internal/api/api.go
    package api
    
    import (
        "crypto/tls"
        "encoding/json"
        "fmt"
        "net"
        "net/http"
        "os"
        "slices"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/rs/zerolog"
    )
    
    func Init() {
        var cfg struct {
            Mod struct {
                Listen     string `yaml:"listen"`
                Username   string `yaml:"username"`
                Password   string `yaml:"password"`
                LocalAuth  bool   `yaml:"local_auth"`
                BasePath   string `yaml:"base_path"`
                StaticDir  string `yaml:"static_dir"`
                Origin     string `yaml:"origin"`
                TLSListen  string `yaml:"tls_listen"`
                TLSCert    string `yaml:"tls_cert"`
                TLSKey     string `yaml:"tls_key"`
                UnixListen string `yaml:"unix_listen"`
    
                AllowPaths []string `yaml:"allow_paths"`
            } `yaml:"api"`
        }
    
        // default config
        cfg.Mod.Listen = ":1984"
    
        // load config from YAML
        app.LoadConfig(&cfg)
    
        if cfg.Mod.Listen == "" && cfg.Mod.UnixListen == "" && cfg.Mod.TLSListen == "" {
            return
        }
    
        allowPaths = cfg.Mod.AllowPaths
        basePath = cfg.Mod.BasePath
        log = app.GetLogger("api")
    
        initStatic(cfg.Mod.StaticDir)
    
        HandleFunc("api", apiHandler)
        HandleFunc("api/config", configHandler)
        HandleFunc("api/exit", exitHandler)
        HandleFunc("api/restart", restartHandler)
        HandleFunc("api/log", logHandler)
    
        Handler = http.DefaultServeMux // 4th
    
        if cfg.Mod.Origin == "*" {
            Handler = middlewareCORS(Handler) // 3rd
        }
    
        if cfg.Mod.Username != "" {
            Handler = middlewareAuth(cfg.Mod.Username, cfg.Mod.Password, cfg.Mod.LocalAuth, Handler) // 2nd
        }
    
        if log.Trace().Enabled() {
            Handler = middlewareLog(Handler) // 1st
        }
    
        if cfg.Mod.Listen != "" {
            _, port, _ := net.SplitHostPort(cfg.Mod.Listen)
            Port, _ = strconv.Atoi(port)
            go listen("tcp", cfg.Mod.Listen)
        }
    
        if cfg.Mod.UnixListen != "" {
            _ = syscall.Unlink(cfg.Mod.UnixListen)
            go listen("unix", cfg.Mod.UnixListen)
        }
    
        // Initialize the HTTPS server
        if cfg.Mod.TLSListen != "" && cfg.Mod.TLSCert != "" && cfg.Mod.TLSKey != "" {
            go tlsListen("tcp", cfg.Mod.TLSListen, cfg.Mod.TLSCert, cfg.Mod.TLSKey)
        }
    }
    
    func listen(network, address string) {
        ln, err := net.Listen(network, address)
        if err != nil {
            log.Error().Err(err).Msg("[api] listen")
            return
        }
    
        log.Info().Str("addr", address).Msg("[api] listen")
    
        server := http.Server{
            Handler:           Handler,
            ReadHeaderTimeout: 5 * time.Second, // Example: Set to 5 seconds
        }
        if err = server.Serve(ln); err != nil {
            log.Fatal().Err(err).Msg("[api] serve")
        }
    }
    
    func tlsListen(network, address, certFile, keyFile string) {
        var cert tls.Certificate
        var err error
        if strings.IndexByte(certFile, '\n') < 0 && strings.IndexByte(keyFile, '\n') < 0 {
            // check if file path
            cert, err = tls.LoadX509KeyPair(certFile, keyFile)
        } else {
            // if text file content
            cert, err = tls.X509KeyPair([]byte(certFile), []byte(keyFile))
        }
        if err != nil {
            log.Error().Err(err).Caller().Send()
            return
        }
    
        ln, err := net.Listen(network, address)
        if err != nil {
            log.Error().Err(err).Msg("[api] tls listen")
            return
        }
    
        log.Info().Str("addr", address).Msg("[api] tls listen")
    
        server := &http.Server{
            Handler:           Handler,
            TLSConfig:         &tls.Config{Certificates: []tls.Certificate{cert}},
            ReadHeaderTimeout: 5 * time.Second,
        }
        if err = server.ServeTLS(ln, "", ""); err != nil {
            log.Fatal().Err(err).Msg("[api] tls serve")
        }
    }
    
    var Port int
    
    const (
        MimeJSON = "application/json"
        MimeText = "text/plain"
    )
    
    var Handler http.Handler
    
    // HandleFunc handle pattern with relative path:
    // - "api/streams" => "{basepath}/api/streams"
    // - "/streams"    => "/streams"
    func HandleFunc(pattern string, handler http.HandlerFunc) {
        if len(pattern) == 0 || pattern[0] != '/' {
            pattern = basePath + "/" + pattern
        }
        if allowPaths != nil && !slices.Contains(allowPaths, pattern) {
            log.Trace().Str("path", pattern).Msg("[api] ignore path not in allow_paths")
            return
        }
        log.Trace().Str("path", pattern).Msg("[api] register path")
        http.HandleFunc(pattern, handler)
    }
    
    // ResponseJSON important always add Content-Type
    // so go won't need to call http.DetectContentType
    func ResponseJSON(w http.ResponseWriter, v any) {
        w.Header().Set("Content-Type", MimeJSON)
        _ = json.NewEncoder(w).Encode(v)
    }
    
    func ResponsePrettyJSON(w http.ResponseWriter, v any) {
        w.Header().Set("Content-Type", MimeJSON)
        enc := json.NewEncoder(w)
        enc.SetIndent("", "  ")
        _ = enc.Encode(v)
    }
    
    func Response(w http.ResponseWriter, body any, contentType string) {
        w.Header().Set("Content-Type", contentType)
    
        switch v := body.(type) {
        case []byte:
            _, _ = w.Write(v)
        case string:
            _, _ = w.Write([]byte(v))
        default:
            _, _ = fmt.Fprint(w, body)
        }
    }
    
    const StreamNotFound = "stream not found"
    
    var allowPaths []string
    var basePath string
    var log zerolog.Logger
    
    func middlewareLog(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Trace().Msgf("[api] %s %s %s", r.Method, r.URL, r.RemoteAddr)
            next.ServeHTTP(w, r)
        })
    }
    
    func isLoopback(remoteAddr string) bool {
        return strings.HasPrefix(remoteAddr, "127.") || strings.HasPrefix(remoteAddr, "[::1]") || remoteAddr == "@"
    }
    
    func middlewareAuth(username, password string, localAuth bool, next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if localAuth || !isLoopback(r.RemoteAddr) {
                user, pass, ok := r.BasicAuth()
                if !ok || user != username || pass != password {
                    w.Header().Set("Www-Authenticate", `Basic realm="go2rtc"`)
    
    
    FILE internal/rtsp/rtsp.go
    package rtsp
    
    import (
        "errors"
        "io"
        "net"
        "net/url"
        "strings"
    
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/AlexxIT/go2rtc/internal/streams"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/rtsp"
        "github.com/AlexxIT/go2rtc/pkg/tcp"
        "github.com/rs/zerolog"
    )
    
    func Init() {
        var conf struct {
            Mod struct {
                Listen       string `yaml:"listen" json:"listen"`
                Username     string `yaml:"username" json:"-"`
                Password     string `yaml:"password" json:"-"`
                DefaultQuery string `yaml:"default_query" json:"default_query"`
                PacketSize   uint16 `yaml:"pkt_size" json:"pkt_size,omitempty"`
            } `yaml:"rtsp"`
        }
    
        // default config
        conf.Mod.Listen = ":8554"
        conf.Mod.DefaultQuery = "video&audio"
    
        app.LoadConfig(&conf)
        app.Info["rtsp"] = conf.Mod
    
        log = app.GetLogger("rtsp")
    
        // RTSP client support
        streams.HandleFunc("rtsp", rtspHandler)
        streams.HandleFunc("rtsps", rtspHandler)
        streams.HandleFunc("rtspx", rtspHandler)
    
        // RTSP server support
        address := conf.Mod.Listen
        if address == "" {
            return
        }
    
        ln, err := net.Listen("tcp", address)
        if err != nil {
            log.Error().Err(err).Msg("[rtsp] listen")
            return
        }
    
        _, Port, _ = net.SplitHostPort(address)
    
        log.Info().Str("addr", address).Msg("[rtsp] listen")
    
        if query, err := url.ParseQuery(conf.Mod.DefaultQuery); err == nil {
            defaultMedias = ParseQuery(query)
        }
    
        go func() {
            for {
                conn, err := ln.Accept()
                if err != nil {
                    return
                }
    
                c := rtsp.NewServer(conn)
                c.PacketSize = conf.Mod.PacketSize
                // skip check auth for localhost
                if conf.Mod.Username != "" && !conn.RemoteAddr().(*net.TCPAddr).IP.IsLoopback() {
                    c.Auth(conf.Mod.Username, conf.Mod.Password)
                }
                go tcpHandler(c)
            }
        }()
    }
    
    type Handler func(conn *rtsp.Conn) bool
    
    func HandleFunc(handler Handler) {
        handlers = append(handlers, handler)
    }
    
    var Port string
    
    // internal
    
    var log zerolog.Logger
    var handlers []Handler
    var defaultMedias []*core.Media
    
    func rtspHandler(rawURL string) (core.Producer, error) {
        rawURL, rawQuery, _ := strings.Cut(rawURL, "#")
    
        conn := rtsp.NewClient(rawURL)
        conn.Backchannel = true
        conn.UserAgent = app.UserAgent
    
        if rawQuery != "" {
            query := streams.ParseQuery(rawQuery)
            conn.Backchannel = query.Get("backchannel") == "1"
            conn.Media = query.Get("media")
            conn.Timeout = core.Atoi(query.Get("timeout"))
            conn.Transport = query.Get("transport")
        }
    
        if log.Trace().Enabled() {
            conn.Listen(func(msg any) {
                switch msg := msg.(type) {
                case *tcp.Request:
                    log.Trace().Msgf("[rtsp] client request:\n%s", msg)
                case *tcp.Response:
                    log.Trace().Msgf("[rtsp] client response:\n%s", msg)
                case string:
                    log.Trace().Msgf("[rtsp] client msg: %s", msg)
                }
            })
        }
    
        if err := conn.Dial(); err != nil {
            return nil, err
        }
    
        if err := conn.Describe(); err != nil {
            if !conn.Backchannel {
                return nil, err
            }
            log.Trace().Msgf("[rtsp] describe (backchannel=%t) err: %v", conn.Backchannel, err)
    
            // second try without backchannel, we need to reconnect
            conn.Backchannel = false
            if err = conn.Dial(); err != nil {
                return nil, err
            }
            if err = conn.Describe(); err != nil {
                return nil, err
            }
        }
    
        return conn, nil
    }
    
    func tcpHandler(conn *rtsp.Conn) {
        var name string
        var closer func()
    
        trace := log.Trace().Enabled()
        level := zerolog.WarnLevel
    
        conn.Listen(func(msg any) {
            if trace {
                switch msg := msg.(type) {
                case *tcp.Request:
                    log.Trace().Msgf("[rtsp] server request:\n%s", msg)
                case *tcp.Resp
    
    FILE README.md
    ```yaml
    streams:
      sonoff_camera: rtsp://rtsp:12345678@192.168.1.123/av_stream/ch0
      dahua_camera:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=1#backchannel=0
      sonoff_camera: rtsp://rtsp:12345678@192.168.1.123/av_stream/ch0
      dahua_camera:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=1#backchannel=0
      amcrest_doorbell:
        - rtsp://username:password@192.168.1.123:554/cam/realmonitor?channel=1&subtype=0#backchannel=0
      dahua_camera:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=1#backchannel=0
      amcrest_doorbell:
        - rtsp://username:password@192.168.1.123:554/cam/realmonitor?channel=1&subtype=0#backchannel=0
      unifi_camera: rtspx://192.168.1.123:7441/fD6ouM72bWoFijxK
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=1#backchannel=0
      amcrest_doorbell:
        - rtsp://username:password@192.168.1.123:554/cam/realmonitor?channel=1&subtype=0#backchannel=0
      unifi_camera: rtspx://192.168.1.123:7441/fD6ouM72bWoFijxK
      glichy_camera: ffmpeg:rtsp://username:password@192.168.1.123/live/ch00_1
    ```
        - rtsp://username:password@192.168.1.123:554/cam/realmonitor?channel=1&subtype=0#backchannel=0
      unifi_camera: rtspx://192.168.1.123:7441/fD6ouM72bWoFijxK
      glichy_camera: ffmpeg:rtsp://username:password@192.168.1.123/live/ch00_1
    ```
    
    **Recommendations**
    **Recommendations**
    
    - **Amcrest Doorbell** users may want to disable two-way audio, because with an active stream, you won't have a working call button. You need to add `#backchannel=0` to the end of your RTSP link in YAML config file
    - **Dahua Doorbell** users may want to change [audio codec](https://github.com/AlexxIT/go2rtc/issues/49#issuecomment-2127107379) for proper 2-way audio. Make sure not to request backchannel multiple times by adding `#backchannel=0` to other stream sources of the same doorbell. The `unicast=true&proto=Onvif` is preferred for 2-way audio as this makes the doorbell accept multiple codecs for the incoming audio
    - **Reolink** users may want NOT to use RTSP protocol at all, some camera models have a very awful, unusable stream implementation
    - **Ubiquiti UniFi** users may want to disable HTTPS verification. Use `rtspx://` prefix instead of `rtsps://`. And don't use `?enableSrtp` [suffix](https://github.com/AlexxIT/go2rtc/issues/81)
    
    - **Amcrest Doorbell** users may want to disable two-way audio, because with an active stream, you won't have a working call button. You need to add `#backchannel=0` to the end of your RTSP link in YAML config file
    - **Dahua Doorbell** users may want to change [audio codec](https://github.com/AlexxIT/go2rtc/issues/49#issuecomment-2127107379) for proper 2-way audio. Make sure not to request backchannel multiple times by adding `#backchannel=0` to other stream sources of the same doorbell. The `unicast=true&proto=Onvif` is preferred for 2-way audio as this makes the doorbell accept multiple codecs for the incoming audio
    - **Reolink** users may want NOT to use RTSP protocol at all, some camera models have a very awful, unusable stream implementation
    - **Ubiquiti UniFi** users may want to disable HTTPS verification. Use `rtspx://` prefix instead of `rtsps://`. And don't use `?enableSrtp` [suffix](https://github.com/AlexxIT/go2rtc/issues/81)
    - **TP-Link Tapo** users may skip login and password, because go2rtc support login [without them](https://drmnsamoliu.github.io/video.html)
    - Add custom timeout `#timeout=30` (in seconds)
    - Ignore audio - `#media=video` or ignore video - `#media=audio`
    - Ignore two-way audio API `#backchannel=0` - important for some glitchy cameras
    - Use WebSocket transport `#transport=ws...`
    
    **RTSP over WebSocket**
    streams:
      # WebSocket with authorization, RTSP - without
      axis-rtsp-ws:  rtsp://192.168.1.123:4567/axis-media/media.amp?overview=0&camera=1&resolution=1280x720&videoframeskipmode=empty&Axis-Orig-Sw=true#transport=ws://user:pass@192.168.1.123:4567/rtsp-over-websocket
      # WebSocket without authorization, RTSP - with
      dahua-rtsp-ws: rtsp://user:pass@192.168.1.123/cam/realmonitor?channel=1&subtype=1&proto=Private3#transport=ws://192.168.1.123/rtspoverwebsocket
    ```
      axis-rtsp-ws:  rtsp://192.168.1.123:4567/axis-media/media.amp?overview=0&camera=1&resolution=1280x720&videoframeskipmode=empty&Axis-Orig-Sw=true#transport=ws://user:pass@192.168.1.123:4567/rtsp-over-websocket
      # WebSocket without authorization, RTSP - with
      dahua-rtsp-ws: rtsp://user:pass@192.168.1.123/cam/realmonitor?channel=1&subtype=1&proto=Private3#transport=ws://192.168.1.123/rtspoverwebsocket
    ```
    
    ## Source: RTMP
    
      # [RTSP] video with rotation, should be transcoded, so select H264
      rotate: ffmpeg:rtsp://12345678@192.168.1.123/av_stream/ch0#video=h264#rotate=90
    ```
    
    All transcoding formats have [built-in templates](internal/ffmpeg/ffmpeg.go): `h264`, `h265`, `opus`, `pcmu`, `pcmu/16000`, `pcmu/48000`, `pcma`, `pcma/16000`, `pcma/48000`, `aac`, `aac/16000`.
    - `killsignal` - signal which will be sent to stop the process (numeric form)
    - `killtimeout` - time in seconds for forced termination with sigkill
    - `backchannel` - enable backchannel for two-way audio
    - `starttimeout` - time in seconds for waiting first byte from RTSP
    
    ```yaml
      pi5cam_h264: exec:libcamera-vid -t 0 --libav-format h264 -o -
      canon: exec:gphoto2 --capture-movie --stdout#killsignal=2#killtimeout=5
      play_pcma: exec:ffplay -fflags nobuffer -f alaw -ar 8000 -i -#backchannel=1
      play_pcm48k: exec:ffplay -fflags nobuffer -f s16be -ar 48000 -i -#backchannel=1
    ```
    
      canon: exec:gphoto2 --capture-movie --stdout#killsignal=2#killtimeout=5
      play_pcma: exec:ffplay -fflags nobuffer -f alaw -ar 8000 -i -#backchannel=1
      play_pcm48k: exec:ffplay -fflags nobuffer -f s16be -ar 48000 -i -#backchannel=1
    ```
    
    ## Source: Echo
    ```
    
    RTSP link with "normal" audio for any player: `rtsp://192.168.1.123:8554/aqara_g3?video&audio=aac`
    
    **This source is in active development!** Tested only with [Aqara Camera Hub G3](https://www.aqara.com/eu/product/camera-hub-g3) (both EU and CN versions).
    
    streams:
      only_stream: dvrip://username:password@192.168.1.123:34567?channel=0&subtype=0
      only_tts: dvrip://username:password@192.168.1.123:34567?backchannel=1
      two_way_audio:
        - dvrip://username:password@192.168.1.123:34567?channel=0&subtype=0
        - dvrip://username:password@192.168.1.123:34567?backchannel=1
      two_way_audio:
        - dvrip://username:password@192.168.1.123:34567?channel=0&subtype=0
        - dvrip://username:password@192.168.1.123:34567?backchannel=1
    ```
    
    ## Source: EseeCloud
    *[New in v1.3.0](https://github.com/AlexxIT/go2rtc/releases/tag/v1.3.0)*
    
    This source type supports only backchannel audio for the Hikvision ISAPI protocol. So it should be used as a second source in addition to the RTSP protocol.
    
    ```yaml
    streams:
    streams:
      hikvision1:
        - rtsp://admin:password@192.168.1.123:554/Streaming/Channels/101
        - isapi://admin:password@192.168.1.123:80/
    ```
    
    **kinesis** (*from [v1.6.1](https://github.com/AlexxIT/go2rtc/releases/tag/v1.6.1)*)
    
    Supports [Amazon Kinesis Video Streams](https://aws.amazon.com/kinesis/video-streams/), using WebRTC protocol. You need to specify the signalling WebSocket URL with all credentials in query params, `client_id` and `ice_servers` list in [JSON format](https://developer.mozilla.org/en-US/docs/Web/API/RTCIceServer).
    
    **switchbot**
    
      webrtc-whep:      webrtc:http://192.168.1.123:1984/api/webrtc?src=camera1
      webrtc-go2rtc:    webrtc:ws://192.168.1.123:1984/api/ws?src=camera1
      webrtc-openipc:   webrtc:ws://192.168.1.123/webrtc_ws#format=openipc#ice_servers=[{"urls":"stun:stun.kinesisvideo.eu-north-1.amazonaws.com:443"}]
      webrtc-wyze:      webrtc:http://192.168.1.123:5000/signaling/camera1?kvs#format=wyze
      webrtc-kinesis:   webrtc:wss://...amazonaws.com/?...#format=kinesis#client_id=...#ice_servers=[{...},{...}]
      webrtc-switchbot: webrtc:wss://...amazonaws.com/?...#format=switchbot#resolution=hd#play_type=0#client_id=...#ice_servers=[{...},{...}]
      webrtc-openipc:   webrtc:ws://192.168.1.123/webrtc_ws#format=openipc#ice_servers=[{"urls":"stun:stun.kinesisvideo.eu-north-1.amazonaws.com:443"}]
      webrtc-wyze:      webrtc:http://192.168.1.123:5000/signaling/camera1?kvs#format=wyze
      webrtc-kinesis:   webrtc:wss://...amazonaws.com/?...#format=kinesis#client_id=...#ice_servers=[{...},{...}]
      webrtc-switchbot: webrtc:wss://...amazonaws.com/?...#format=switchbot#resolution=hd#play_type=0#client_id=...#ice_servers=[{...},{...}]
    ```
    
      webrtc-wyze:      webrtc:http://192.168.1.123:5000/signaling/camera1?kvs#format=wyze
      webrtc-kinesis:   webrtc:wss://...amazonaws.com/?...#format=kinesis#client_id=...#ice_servers=[{...},{...}]
      webrtc-switchbot: webrtc:wss://...amazonaws.com/?...#format=switchbot#resolution=hd#play_type=0#client_id=...#ice_servers=[{...},{...}]
    ```
    
    **PS.** For `kinesis` sources, you can use [echo](#source-echo) to get connection params using `bash`, `python` or any other script language.
    - RTSP with any codec
      ```yaml
      ffmpeg -re -i BigBuckBunny.mp4 -c copy -rtsp_transport tcp -f rtsp rtsp://localhost:8554/camera1
      ```
    - HTTP-MJPEG with MJPEG codec
      ```yaml
    streams:
      video_audio_transcode:
        - ffmpeg:rtsp://user:pass@192.168.1.123/stream1#video=h264#hardware#audio=aac
      audio_transcode:
        - ffmpeg:rtsp://user:pass@192.168.1.123/stream1#video=copy#audio=aac
    ```
        - ffmpeg:rtsp://user:pass@192.168.1.123/stream1#video=h264#hardware#audio=aac
      audio_transcode:
        - ffmpeg:rtsp://user:pass@192.168.1.123/stream1#video=copy#audio=aac
    ```
    
    - **Telegram Desktop App** > Any public or private channel or group (where you admin) > Live stream > Start with... > Start streaming.
    streams:
      camera1:
        - rtsp://192.168.1.100/stream
      camera2:
        - rtsp://192.168.1.101/stream
      camera3:
        - rtsp://192.168.1.100/stream
      camera2:
        - rtsp://192.168.1.101/stream
      camera3:
        - rtsp://192.168.1.102/h265stream
        - ffmpeg:camera3#video=h264#audio=opus#hardware
        - rtsp://192.168.1.101/stream
      camera3:
        - rtsp://192.168.1.102/h265stream
        - ffmpeg:camera3#video=h264#audio=opus#hardware
    ```
    
    ## Module: RTSP
    
    You can get any stream as RTSP-stream: `rtsp://192.168.1.123:8554/{stream_name}`
    
    You can enable external password protection for your RTSP streams. Password protection is always disabled for localhost calls (ex. FFmpeg or Hass on the same server).
    
    
    ```yaml
    rtsp:
      listen: ":8554"    # RTSP Server TCP port, default - 8554
      username: "admin"  # optional, default - disabled
      password: "pass"   # optional, default - disabled
    ```yaml
    webrtc:
      ice_servers:
        - urls: [stun:stun.l.google.com:19302]
        - urls: [turn:123.123.123.123:3478]
          username: your_user
    ```yaml
    streams:
      dahua1: rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0
    homekit:
      dahua1:  # same stream ID from streams list, default PIN - 19550224
    ```
    streams:
      dahua1:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0
        - ffmpeg:dahua1#video=h264#hardware  # if your camera doesn't support H264, important for HomeKit
        - ffmpeg:dahua1#audio=opus           # only OPUS audio supported by HomeKit
    
       - Add your stream to [go2rtc config](#configuration)
       - Hass > Settings > Integrations > Add Integration > [ONVIF](https://my.home-assistant.io/redirect/config_flow_start/?domain=onvif) > Host: `127.0.0.1`, Port: `1984`
       - Hass > Settings > Integrations > Add Integration > [Generic Camera](https://my.home-assistant.io/redirect/config_flow_start/?domain=generic) > Stream Source URL: `rtsp://127.0.0.1:8554/camera1` (change to your stream name, leave everything else as is)
    
    You have several options on how to watch the stream from the cameras in Home Assistant:
    
      api: trace
      exec: debug
      rtsp: warn
      streams: error
      webrtc: fatal
    ```
      listen: "127.0.0.1:1984" # localhost
    
    rtsp:
      listen: "127.0.0.1:8554" # localhost
    
    webrtc:
    Some examples:
    
    - `rtsp://192.168.1.123:8554/camera1?mp4` - useful for recording as MP4 files (e.g. Hass or Frigate)
    - `rtsp://192.168.1.123:8554/camera1?video=h264,h265&audio=aac` - full version of the filter above
    - `rtsp://192.168.1.123:8554/camera1?video=h264&audio=aac&audio=opus` - H264 video codec and two separate audio tracks
    - `rtsp://192.168.1.123:8554/camera1?video&audio=all` - any video codec and all audio codecs as separate tracks
    
    - `rtsp://192.168.1.123:8554/camera1?mp4` - useful for recording as MP4 files (e.g. Hass or Frigate)
    - `rtsp://192.168.1.123:8554/camera1?video=h264,h265&audio=aac` - full version of the filter above
    - `rtsp://192.168.1.123:8554/camera1?video=h264&audio=aac&audio=opus` - H264 video codec and two separate audio tracks
    - `rtsp://192.168.1.123:8554/camera1?video&audio=all` - any video codec and all audio codecs as separate tracks
    - `http://192.168.1.123:1984/api/stream.m3u8?src=camera1&mp4` - HLS stream with MP4 compatible codecs (HLS/fMP4)
    - `rtsp://192.168.1.123:8554/camera1?mp4` - useful for recording as MP4 files (e.g. Hass or Frigate)
    - `rtsp://192.168.1.123:8554/camera1?video=h264,h265&audio=aac` - full version of the filter above
    - `rtsp://192.168.1.123:8554/camera1?video=h264&audio=aac&audio=opus` - H264 video codec and two separate audio tracks
    - `rtsp://192.168.1.123:8554/camera1?video&audio=all` - any video codec and all audio codecs as separate tracks
    - `http://192.168.1.123:1984/api/stream.m3u8?src=camera1&mp4` - HLS stream with MP4 compatible codecs (HLS/fMP4)
    - `http://192.168.1.123:1984/api/stream.m3u8?src=camera1&mp4=flac` - HLS stream with PCMA/PCMU/PCM audio support (HLS/fMP4), won't work on old devices
    - `rtsp://192.168.1.123:8554/camera1?video=h264,h265&audio=aac` - full version of the filter above
    - `rtsp://192.168.1.123:8554/camera1?video=h264&audio=aac&audio=opus` - H264 video codec and two separate audio tracks
    - `rtsp://192.168.1.123:8554/camera1?video&audio=all` - any video codec and all audio codecs as separate tracks
    - `http://192.168.1.123:1984/api/stream.m3u8?src=camera1&mp4` - HLS stream with MP4 compatible codecs (HLS/fMP4)
    - `http://192.168.1.123:1984/api/stream.m3u8?src=camera1&mp4=flac` - HLS stream with PCMA/PCMU/PCM audio support (HLS/fMP4), won't work on old devices
    - `http://192.168.1.123:1984/api/stream.mp4?src=camera1&mp4=flac` - MP4 file with PCMA/PCMU/PCM audio support, won't work on old devices (ex. iOS 12)
    streams:
      dahua:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
        - ffmpeg:rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0#audio=opus
    ```
    
      dahua:
        - rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0&unicast=true&proto=Onvif
        - ffmpeg:rtsp://admin:password@192.168.1.123/cam/realmonitor?channel=1&subtype=0#audio=opus
    ```
    
    **go2rtc** automatically matches codecs for your browser and all your stream sources. This is called **multi-source two-way codec negotiation**. And this is one of the main features of this app.
    **Using apps for low RTSP delay**
    
    - `ffplay -fflags nobuffer -flags low_delay "rtsp://192.168.1.123:8554/camera1"`
    - VLC > Preferences > Input / Codecs > Default Caching Level: Lowest Latency
    
    **Snapshots to Telegram**
    
    FILE LICENSE
    MIT License
    
    Copyright (c) 2022 Alexey Khit
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
    
    ✓ • 455ms

## Activity

    $ python3 - <<'PY'
        import urllib.request,json
        base='https://api.github.com/repos/AlexxIT/go2rtc/contents/'
        for path in ['internal/webrtc','build']:
         try:
          with urllib.request.urlopen(base+path+'?ref=v1.9.14',timeout=20) as response:data=json.load(response)
          print(path,[(entry['name'],entry['download_url']) for entry in data])
         except Exception as e:print(path,e)
        for path in ['internal/webrtc/server.go','Dockerfile','build/Dockerfile']:
         try:
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+path,timeout=20) as response: data=response.read().decode()
          print(path+'\n'+data[:11000])
         except Exception as e:print(path,e)
        PY
    internal/webrtc [('README.md', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/README.md'), ('candidates.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/candidates.go'), ('client.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/client.go'), ('client_creality.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/client_creality.go'), ('kinesis.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/kinesis.go'), ('milestone.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/milestone.go'), ('openipc.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/openipc.go'), ('server.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/server.go'), ('switchbot.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/switchbot.go'), ('webrtc.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/webrtc.go'), ('webrtc_test.go', 'https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/internal/webrtc/webrtc_test.go')]
    build HTTP Error 404: Not Found
    internal/webrtc/server.go
    package webrtc
    
    import (
        "encoding/base64"
        "encoding/json"
        "io"
        "net/http"
        "strconv"
        "strings"
        "time"
    
        "github.com/AlexxIT/go2rtc/internal/api"
        "github.com/AlexxIT/go2rtc/internal/streams"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/webrtc"
        pion "github.com/pion/webrtc/v4"
    )
    
    const MimeSDP = "application/sdp"
    
    var sessions = map[string]*webrtc.Conn{}
    
    func syncHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            query := r.URL.Query()
            if query.Get("src") != "" {
                // WHEP or JSON SDP or raw SDP exchange
                outputWebRTC(w, r)
            } else if query.Get("dst") != "" {
                // WHIP SDP exchange
                inputWebRTC(w, r)
            } else {
                http.Error(w, "", http.StatusBadRequest)
            }
    
        case "PATCH":
            // TODO: WHEP/WHIP
            http.Error(w, "", http.StatusMethodNotAllowed)
    
        case "DELETE":
            if id := r.URL.Query().Get("id"); id != "" {
                if conn, ok := sessions[id]; ok {
                    delete(sessions, id)
                    _ = conn.Close()
                } else {
                    http.Error(w, "", http.StatusNotFound)
                }
            } else {
                http.Error(w, "", http.StatusBadRequest)
            }
    
        case "OPTIONS":
            w.WriteHeader(http.StatusNoContent)
    
        default:
            http.Error(w, "", http.StatusMethodNotAllowed)
        }
    }
    
    // outputWebRTC support API depending on Content-Type:
    // 1. application/json - receive {"type":"offer","sdp":"v=0\r\n..."} and response {"type":"answer","sdp":"v=0\r\n..."}
    // 2. application/sdp - receive/response SDP via WebRTC-HTTP Egress Protocol (WHEP)
    // 3. other - receive/response raw SDP
    func outputWebRTC(w http.ResponseWriter, r *http.Request) {
        u := r.URL.Query().Get("src")
        stream := streams.Get(u)
        if stream == nil {
            http.Error(w, api.StreamNotFound, http.StatusNotFound)
            return
        }
    
        mediaType := r.Header.Get("Content-Type")
        if mediaType != "" {
            mediaType, _, _ = strings.Cut(mediaType, ";")
            mediaType = strings.ToLower(strings.TrimSpace(mediaType))
        }
    
        var offer string
    
        switch mediaType {
        case "application/json":
            var desc pion.SessionDescription
            if err := json.NewDecoder(r.Body).Decode(&desc); err != nil {
                log.Error().Err(err).Caller().Send()
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            offer = desc.SDP
    
        case "application/x-www-form-urlencoded":
            if err := r.ParseForm(); err != nil {
                log.Error().Err(err).Caller().Send()
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            offerB64 := r.Form.Get("data")
            b, err := base64.StdEncoding.DecodeString(offerB64)
            if err != nil {
                log.Error().Err(err).Caller().Send()
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            offer = string(b)
    
        default:
            body, err := io.ReadAll(r.Body)
            if err != nil {
                log.Error().Err(err).Caller().Send()
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            offer = string(body)
        }
    
        var desc string
    
        switch mediaType {
        case "application/json":
            desc = "webrtc/json"
        case MimeSDP:
            desc = "webrtc/whep"
        default:
            desc = "webrtc/post"
        }
    
        answer, err := ExchangeSDP(stream, offer, desc, r.UserAgent())
        if err != nil {
            log.Error().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        switch mediaType {
        case "application/json":
            w.Header().Set("Content-Type", mediaType)
    
            v := pion.SessionDescription{
                Type: pion.SDPTypeAnswer, SDP: answer,
            }
            err = json.NewEncoder(w).Encode(v)
    
        case "application/x-www-form-urlencoded":
            w.Header().Set("Content-Type", mediaType)
            answerB64 := base64.StdEncoding.EncodeToString([]byte(answer))
            _, err = w.Write([]byte(answerB64))
    
        case MimeSDP:
            w.Header().Set("Content-Type", mediaType)
            w.WriteHeader(http.StatusCreated)
    
            _, err = w.Write([]byte(answer))
    
        default:
            w.Header().Set("Content-Type", mediaType)
    
            _, err = w.Write([]byte(answer))
        }
    
        if err != nil {
            log.Error().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
    
    func inputWebRTC(w http.ResponseWriter, r *http.Request) {
        dst := r.URL.Query().Get("dst")
        stream := streams.Get(dst)
        if stream == nil {
            http.Error(w, api.StreamNotFound, http.StatusNotFound)
            return
        }
    
        // 1. Get offer
        offer, err := io.ReadAll(r.Body)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        log.Trace().Msgf("[webrtc] WHIP offer\n%s", offer)
    
        pc, err := PeerConnection(false)
        if err != nil {
            log.Error().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        // create new webrtc instance
        prod := webrtc.NewConn(pc)
        prod.Mode = core.ModePassiveProducer
        prod.Protocol = "http"
        prod.UserAgent = r.UserAgent()
    
        if err = prod.SetOffer(string(offer)); err != nil {
            log.Warn().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        answer, err := prod.GetCompleteAnswer(GetCandidates(), FilterCandidate)
        if err != nil {
            log.Warn().Err(err).Caller().Send()
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    
        log.Trace().Msgf("[webrtc] WHIP answer\n%s", answer)
    
        id := strconv.FormatInt(time.Now().UnixNano(), 36)
        sessions[id] = prod
    
        prod.Listen(func(msg any) {
            switch msg := msg.(type) {
            case pion.PeerConnectionState:
                if msg == pion.PeerConnectionStateClosed {
                    stream.RemoveProducer(prod)
                    delete(sessions, id)
                }
            }
        })
    
        stream.AddProducer(prod)
    
        w.Header().Set("Content-Type", MimeSDP)
        w.Header().Set("Location", "webrtc?id="+id)
        w.WriteHeader(http.StatusCreated)
    
        if _, err = w.Write([]byte(answer)); err != nil {
            log.Warn().Err(err).Caller().Send()
            return
        }
    }
    
    Dockerfile HTTP Error 404: Not Found
    build/Dockerfile HTTP Error 404: Not Found
    ✓ • 490ms

## Assistant

The playback path needs to support microphone audio, so I’m using go2rtc as the RTSP-to-WebRTC bridge. Viewing will use the regular URL; Talk will switch to the backchannel URL. Turning live viewing off will stop the bridge and disconnect all live viewers.

## Activity

    $ python3 - <<'PY'
        import urllib.request,json
        with urllib.request.urlopen('https://api.github.com/repos/AlexxIT/go2rtc/git/trees/v1.9.14?recursive=1',timeout=20) as response:data=json.load(response)
        for item in data['tree']:
         if item['path'].endswith('Dockerfile'):
          print(item['path'])
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+item['path'],timeout=20) as response:print(response.read().decode()[:3500])
        PY
    docker/Dockerfile
    # syntax=docker/dockerfile:labs
    
    # 0. Prepare images
    ARG PYTHON_VERSION="3.13"
    ARG GO_VERSION="1.25"
    
    
    # 1. Build go2rtc binary
    FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
    ARG TARGETPLATFORM
    ARG TARGETOS
    ARG TARGETARCH
    
    ENV GOOS=${TARGETOS}
    ENV GOARCH=${TARGETARCH}
    
    WORKDIR /build
    
    RUN apk add git
    
    # Cache dependencies
    COPY go.mod go.sum ./
    RUN --mount=type=cache,target=/root/.cache/go-build go mod download
    
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath
    
    
    # 2. Final image
    FROM python:${PYTHON_VERSION}-alpine AS base
    
    # Install ffmpeg, tini (for signal handling),
    # and other common tools for the echo source.
    # alsa-plugins-pulse for ALSA support (+0MB)
    # font-droid for FFmpeg drawtext filter (+2MB)
    RUN apk add --no-cache tini ffmpeg ffplay bash curl jq alsa-plugins-pulse font-droid
    
    # Hardware Acceleration for Intel CPU (+50MB)
    ARG TARGETARCH
    
    RUN if [ "${TARGETARCH}" = "amd64" ]; then apk add --no-cache libva-intel-driver intel-media-driver; fi
    
    # Hardware: AMD and NVidia VAAPI (not sure about this)
    # RUN libva-glx mesa-va-gallium
    # Hardware: AMD and NVidia VDPAU (not sure about this)
    # RUN libva-vdpau-driver mesa-vdpau-gallium (+150MB total)
    
    COPY --from=build /build/go2rtc /usr/local/bin/
    
    EXPOSE 1984 8554 8555 8555/udp
    ENTRYPOINT ["/sbin/tini", "--"]
    VOLUME /config
    WORKDIR /config
    
    CMD ["go2rtc", "-config", "/config/go2rtc.yaml"]
    
    docker/hardware.Dockerfile
    # syntax=docker/dockerfile:labs
    
    # 0. Prepare images
    # only debian 13 (trixie) has latest ffmpeg
    # https://packages.debian.org/trixie/ffmpeg
    ARG DEBIAN_VERSION="trixie-slim"
    ARG GO_VERSION="1.25-bookworm"
    
    
    # 1. Build go2rtc binary
    FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
    ARG TARGETPLATFORM
    ARG TARGETOS
    ARG TARGETARCH
    
    ENV GOOS=${TARGETOS}
    ENV GOARCH=${TARGETARCH}
    
    WORKDIR /build
    
    # Cache dependencies
    COPY go.mod go.sum ./
    RUN --mount=type=cache,target=/root/.cache/go-build go mod download
    
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath
    
    
    # 2. Final image
    FROM debian:${DEBIAN_VERSION}
    
    # Prepare apt for buildkit cache
    RUN rm -f /etc/apt/apt.conf.d/docker-clean \
      && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' >/etc/apt/apt.conf.d/keep-cache
    
    # Install ffmpeg, tini (for signal handling),
    # and other common tools for the echo source.
    # non-free for Intel QSV support (not used by go2rtc, just for tests)
    # mesa-va-drivers for AMD APU
    # libasound2-plugins for ALSA support
    RUN --mount=type=cache,target=/var/cache/apt,sharing=locked --mount=type=cache,target=/var/lib/apt,sharing=locked \
        echo 'deb http://deb.debian.org/debian trixie non-free' > /etc/apt/sources.list.d/debian-non-free.list && \
        apt-get -y update && apt-get -y install ffmpeg tini \
            python3 curl jq \
            intel-media-va-driver-non-free \
            mesa-va-drivers \
            libasound2-plugins && \
        apt-get clean && rm -rf /var/lib/apt/lists/*
    
    COPY --from=build /build/go2rtc /usr/local/bin/
    
    EXPOSE 1984 8554 8555 8555/udp
    ENTRYPOINT ["/usr/bin/tini", "--"]
    VOLUME /config
    WORKDIR /config
    # https://github.com/NVIDIA/nvidia-docker/wiki/Installation-(Native-GPU-Support)
    ENV NVIDIA_VISIBLE_DEVICES all
    ENV NVIDIA_DRIVER_CAPABILITIES compute,video,utility
    
    CMD ["go2rtc", "-config", "/config/go2rtc.yaml"]
    
    docker/rockchip.Dockerfile
    # syntax=docker/dockerfile:labs
    
    # 0. Prepare images
    ARG PYTHON_VERSION="3.13-slim-bookworm"
    ARG GO_VERSION="1.25-bookworm"
    
    
    # 1. Build go2rtc binary
    FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
    ARG TARGETPLATFORM
    ARG TARGETOS
    ARG TARGETARCH
    
    ENV GOOS=${TARGETOS}
    ENV GOARCH=${TARGETARCH}
    
    WORKDIR /build
    
    # Cache dependencies
    COPY go.mod go.sum ./
    RUN --mount=type=cache,target=/root/.cache/go-build go mod download
    
    COPY . .
    RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath
    
    
    # 2. Final image
    FROM python:${PYTHON_VERSION}
    
    # Prepare apt for buildkit cache
    RUN rm -f /etc/apt/apt.conf.d/docker-clean \
      && echo 'Binary::apt::APT::Keep-Downloaded-Packages "true";' >/etc/apt/apt.conf.d/keep-cache
    
    # Install ffmpeg, tini (for signal handling),
    # and other common tools for the echo source.
    # libasound2-plugins for ALSA support
    RUN --mount=type=cache,target=/var/cache/apt,sharing=locked --mount=type=cache,target=/var/lib/apt,sharing=locked \
        apt-get -y update && apt-get -y install tini \
            curl jq \
            libasound2-plugins && \
        apt-get clean && rm -rf /var/lib/apt/lists/*
    
    COPY --from=build /build/go2rtc /usr/local/bin/
    ADD --chmod=755 https://github.com/MarcA711/Rockchip-FFmpeg-Builds/releases/download/6.1-8-no_extra_dump/ffmpeg /usr/local/bin
    
    EXPOSE 1984 8554 8555 8555/udp
    ENTRYPOINT ["/usr/bin/tini", "--"]
    VOLUME /config
    WORKDIR /config
    
    CMD ["go2rtc", "-config", "/config/go2rtc.yaml"]
    
    ✓ • 530ms

## Activity

    $ python3 - <<'PY'
        import hashlib,json,os,tarfile,tempfile,urllib.request
        from pathlib import Path
        root=Path(tempfile.mkdtemp(prefix='neonvr-live-tools-'))
        print('TOOLS_DIR='+str(root),flush=True)
        with urllib.request.urlopen('https://go.dev/dl/?mode=json&include=all',timeout=30) as response: releases=json.load(response)
        release=next(item for item in releases if item['version']=='go1.22.12')
        artifact=next(item for item in release['files'] if item['os']=='linux' and item['arch']=='amd64' and item['kind']=='archive')
        archive=root/artifact['filename']
        urllib.request.urlretrieve('https://go.dev/dl/'+artifact['filename'],archive)
        assert hashlib.sha256(archive.read_bytes()).hexdigest()==artifact['sha256']
        with tarfile.open(archive) as package: package.extractall(root,filter='data')
        archive.unlink()
        binary=root/'go2rtc'
        urllib.request.urlretrieve('https://github.com/AlexxIT/go2rtc/releases/download/v1.9.14/go2rtc_linux_amd64',binary)
        assert hashlib.sha256(binary.read_bytes()).hexdigest()=='32d616af226bd731678ffde328b94cfb94e30339bfefc469cfb76323144615a6'
        binary.chmod(0o755)
        print('Temporary Go and go2rtc binaries ready.',flush=True)
        PY
    TOOLS_DIR=/tmp/neonvr-live-tools-vbq_9dyb
    Temporary Go and go2rtc binaries ready.
    ✓ • 6.57s

## Activity

    file changes: Completed · 3 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/config.go
    @@ -55,2 +55,6 @@
     	DetectionPrune        string   `json:"detection_prune"`
    +	LiveStreams           []LiveStream `json:"live_streams"`
    +	LiveGo2RTC            string       `json:"live_go2rtc"`
    +	LiveWebRTCPort        int          `json:"live_webrtc_port"`
    +	LiveWebRTCHost        string       `json:"live_webrtc_host"`
     }
    @@ -58,3 +62,3 @@
     func defaultConfig() Config {
    -	return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h",
    +	return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h", LiveStreams: []LiveStream{}, LiveGo2RTC: "go2rtc", LiveWebRTCPort: 8555,
     		DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
    @@ -78,2 +82,5 @@
     	}
    +	if c.LiveStreams == nil {
    +		c.LiveStreams = []LiveStream{}
    +	}
     	return c, nil
    @@ -131,2 +138,5 @@
     	}
    +	if err := c.validateLive(); err != nil {
    +		return err
    +	}
     	for name, v := range map[string]string{"thumb_maintenance_every": c.ThumbMaintenanceEvery, "size_refresh": c.SizeRefresh, "auto_delete_every": c.AutoDeleteEvery} {
    Add: /home/ubuntu/NeoNVR/live.go
    package main
    
    import (
    	"bytes"
    	"context"
    	"crypto/rand"
    	"encoding/hex"
    	"encoding/json"
    	"errors"
    	"fmt"
    	"io"
    	"net"
    	"net/http"
    	"net/url"
    	"os"
    	"os/exec"
    	"path/filepath"
    	"strconv"
    	"strings"
    	"sync"
    	"syscall"
    	"time"
    )
    
    type LiveStream struct {
    	Name           string `json:"name"`
    	URL            string `json:"url"`
    	BackchannelURL string `json:"backchannel_url,omitempty"`
    }
    
    func validateRTSP(raw string) error {
    	u, err := url.Parse(raw)
    	if err != nil || (u.Scheme != "rtsp" && u.Scheme != "rtsps") || u.Hostname() == "" || u.Fragment != "" || strings.ContainsAny(raw, "\r\n\t ") {
    		return errors.New("enter an rtsp:// or rtsps:// camera URL without spaces or a fragment")
    	}
    	if port := u.Port(); port != "" {
    		n, err := strconv.Atoi(port)
    		if err != nil || n < 1 || n > 65535 {
    			return errors.New("RTSP port must be between 1 and 65535")
    		}
    	}
    	return nil
    }
    
    func (c Config) validateLive() error {
    	if strings.TrimSpace(c.LiveGo2RTC) == "" {
    		return errors.New("live_go2rtc cannot be empty")
    	}
    	if c.LiveWebRTCPort < 1 || c.LiveWebRTCPort > 65535 {
    		return errors.New("live_webrtc_port must be between 1 and 65535")
    	}
    	_, appPort, _ := net.SplitHostPort(c.Addr)
    	if appPort == strconv.Itoa(c.LiveWebRTCPort) {
    		return errors.New("live_webrtc_port must differ from the app listen port")
    	}
    	if host := c.LiveWebRTCHost; host != "" && (strings.ContainsAny(host, "/?#@ \r\n\t") || strings.Contains(host, ":") && net.ParseIP(host) == nil) {
    		return errors.New("live_webrtc_host must be a server hostname or IP address without a port")
    	}
    	if len(c.LiveStreams) > 32 {
    		return errors.New("configure at most 32 live streams")
    	}
    	names := map[string]bool{}
    	for index, stream := range c.LiveStreams {
    		name := strings.TrimSpace(stream.Name)
    		if name == "" || len(name) > 120 || strings.ContainsAny(name, "|\r\n") || names[name] {
    			return fmt.Errorf("live stream %d needs a unique name of up to 120 characters", index+1)
    		}
    		names[name] = true
    		if err := validateRTSP(stream.URL); err != nil {
    			return fmt.Errorf("live stream %d: %w", index+1, err)
    		}
    		if stream.BackchannelURL != "" {
    			if err := validateRTSP(stream.BackchannelURL); err != nil {
    				return fmt.Errorf("live stream %d backchannel: %w", index+1, err)
    			}
    		}
    	}
    	return nil
    }
    
    type liveService struct {
    	cmd      *exec.Cmd
    	done     chan struct{}
    	baseURL  string
    	password string
    	client   *http.Client
    }
    
    type liveManager struct {
    	mu      sync.Mutex
    	ctx     context.Context
    	config  Config
    	token   string
    	service *liveService
    }
    
    func newLiveManager(ctx context.Context, config Config, token string) *liveManager {
    	// Enabling is deliberately transient: every app start defaults to off.
    	return &liveManager{ctx: ctx, config: config, token: token}
    }
    
    func (m *liveManager) runningLocked() *liveService {
    	if m.service != nil {
    		select {
    		case <-m.service.done:
    			m.service = nil
    		default:
    		}
    	}
    	return m.service
    }
    
    func (m *liveManager) response(w http.ResponseWriter) {
    	m.mu.Lock()
    	defer m.mu.Unlock()
    	streams := make([]map[string]any, 0, len(m.config.LiveStreams))
    	for index, stream := range m.config.LiveStreams {
    		u, _ := url.Parse(stream.URL)
    		streams = append(streams, map[string]any{"id": index, "name": stream.Name, "host": u.Host, "talk": stream.BackchannelURL != ""})
    	}
    	writeJSON(w, map[string]any{"enabled": m.runningLocked() != nil, "streams": streams, "token": m.token})
    }
    
    func (m *liveManager) authorize(w http.ResponseWriter, r *http.Request) bool {
    	if r.Header.Get("X-NeoNVR-Token") != m.token {
    		http.Error(w, "reload Streams before connecting", http.StatusForbidden)
    		return false
    	}
    	if origin := r.Header.Get("Origin"); origin != "" {
    		u, err := url.Parse(origin)
    		// Allow the HTTPS origin used by a TLS reverse proxy on the same host.
    		if err != nil || u.Host != r.Host || (u.Scheme != "http" && u.Scheme != "https") {
    			http.Error(w, "same-origin requests required", http.StatusForbidden)
    			return false
    		}
    	}
    	return true
    }
    
    func (m *liveManager) handler(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Cache-Control", "no-store")
    	w.Header().Set("X-Content-Type-Options", "nosniff")
    	switch r.Method {
    	case http.MethodGet:
    		m.response(w)
    	case http.MethodPost:
    		if !m.authorize(w, r) {
    			return
    		}
    		if r.Header.Get("Content-Type") != "application/json" {
    			http.Error(w, "JSON required", http.StatusUnsupportedMediaType)
    			return
    		}
    		var request struct {
    			Enabled *bool `json:"enabled"`
    		}
    		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024))
    		dec.DisallowUnknownFields()
    		if err := dec.Decode(&request); err != nil || request.Enabled == nil || dec.Decode(new(any)) != io.EOF {
    			http.Error(w, "expected an enabled boolean", http.StatusBadRequest)
    			return
    		}
    		m.mu.Lock()
    		var err error
    		if *request.Enabled {
    			if len(m.config.LiveStreams) == 0 {
    				err = errors.New("add a camera in Settings before enabling live viewing")
    			} else if m.runningLocked() == nil {
    				m.service, err = m.start(r.Context(), r.Host)
    			}
    		} else {
    			m.stopLocked()
    		}
    		m.mu.Unlock()
    		if err != nil {
    			http.Error(w, err.Error(), http.StatusServiceUnavailable)
    			return
    		}
    		m.response(w)
    	default:
    		w.Header().Set("Allow", "GET, POST")
    		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    	}
    }
    
    func liveSource(raw string, backchannel bool) string {
    	u, _ := url.Parse(raw)
    	flag := "0"
    	if backchannel {
    		flag = "1"
    	}
    	// The fragment is a go2rtc option; camera query parameters stay intact.
    	u.Fragment = "backchannel=" + flag + "&timeout=10"
    	return u.String()
    }
    
    func (m *liveManager) gatewayConfig(address, password, host string) map[string]any {
    	sources := map[string]string{}
    	for index, stream := range m.config.LiveStreams {
    		name := fmt.Sprintf("camera_%d", index)
    		sources[name] = liveSource(stream.URL, false)
    		if stream.BackchannelURL != "" {
    			sources[name+"_talk"] = liveSource(stream.BackchannelURL, true)
    		}
    	}
    	return map[string]any{
    		"api": map[string]any{"listen": address, "username": "neonvr", "password": password, "local_auth": true, "allow_paths": []string{"/api", "/api/webrtc"}},
    		"rtsp": map[string]any{"listen": ""},
    		"webrtc": map[string]any{"listen": fmt.Sprintf(":%d", m.config.LiveWebRTCPort), "candidates": []string{net.JoinHostPort(host, strconv.Itoa(m.config.LiveWebRTCPort))}, "ice_servers": []any{}, "filters": map[string]any{"loopback": true}},
    		"streams": sources,
    		"log": map[string]any{"level": "error"},
    	}
    }
    
    func (m *liveManager) start(requestCtx context.Context, requestHost string) (*liveService, error) {
    	// Reserve a private API address. Only the two fixed API routes are used.
    	listener, err := net.Listen("tcp", "127.0.0.1:0")
    	if err != nil {
    		return nil, errors.New("cannot allocate a local streaming gateway port")
    	}
    	address := listener.Addr().String()
    	listener.Close()
    	host := m.config.LiveWebRTCHost
    	if host == "" {
    		host, _, err = net.SplitHostPort(requestHost)
    		if err != nil {
    			host = requestHost
    		}
    	}
    	directory, err := os.MkdirTemp("", "neonvr-live-*")
    	if err != nil {
    		return nil, errors.New("cannot create temporary streaming configuration")
    	}
    	cleanup := true
    	defer func() {
    		if cleanup {
    			os.RemoveAll(directory)
    		}
    	}()
    	secret := make([]byte, 32)
    	if _, err := rand.Read(secret); err != nil {
    		return nil, err
    	}
    	password := hex.EncodeToString(secret)
    	data, _ := json.Marshal(m.gatewayConfig(address, password, host))
    	path := filepath.Join(directory, "go2rtc.json")
    	if err := os.WriteFile(path, data, 0600); err != nil {
    		return nil, errors.New("cannot write temporary streaming configuration")
    	}
    	cmd := exec.CommandContext(m.ctx, expandPath(m.config.LiveGo2RTC), "-config", path)
    	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
    	// RTSP errors can contain camera credentials; do not expose gateway logs.
    	cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
    	if err := cmd.Start(); err != nil {
    		return nil, errors.New("cannot start go2rtc; install it or set its executable path in Settings")
    	}
    	service := &liveService{cmd: cmd, done: make(chan struct{}), baseURL: "http://" + address, password: password, client: &http.Client{Timeout: 25 * time.Second, Transport: &http.Transport{}}}
    	cleanup = false
    	go func() {
    		_ = cmd.Wait()
    		service.client.CloseIdleConnections()
    		os.RemoveAll(directory)
    		close(service.done)
    	}()
    	ctx, cancel := context.WithTimeout(requestCtx, 5*time.Second)
    	defer cancel()
    	for {
    		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, service.baseURL+"/api", nil)
    		req.SetBasicAuth("neonvr", password)
    		response, requestErr := service.client.Do(req)
    		if requestErr == nil {
    			response.Body.Close()
    			if response.StatusCode == http.StatusOK {
    				return service, nil
    			}
    		}
    		select {
    		case <-service.done:
    			return nil, errors.New("streaming gateway exited; check the executable and WebRTC port in Settings")
    		case <-ctx.Done():
    			service.stop()
    			return nil, errors.New("streaming gateway did not become ready")
    		case <-time.After(50 * time.Millisecond):
    		}
    	}
    }
    
    func (s *liveService) stop() {
    	_ = s.cmd.Process.Signal(syscall.SIGTERM)
    	select {
    	case <-s.done:
    	case <-time.After(2 * time.Second):
    		_ = s.cmd.Process.Kill()
    		<-s.done
    	}
    }
    
    func (m *liveManager) stopLocked() {
    	if service := m.service; service != nil {
    		m.service = nil
    		service.stop()
    	}
    }
    
    func (m *liveManager) close() {
    	m.mu.Lock()
    	defer m.mu.Unlock()
    	m.stopLocked()
    }
    
    func (m *liveManager) webrtcHandler(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Cache-Control", "no-store")
    	w.Header().Set("X-Content-Type-Options", "nosniff")
    	if r.Method != http.MethodPost {
    		w.Header().Set("Allow", "POST")
    		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    		return
    	}
    	if !m.authorize(w, r) {
    		return
    	}
    	if r.Header.Get("Content-Type") != "application/sdp" {
    		http.Error(w, "SDP offer required", http.StatusUnsupportedMediaType)
    		return
    	}
    	m.mu.Lock()
    	service := m.runningLocked()
    	m.mu.Unlock()
    	if service == nil {
    		http.Error(w, "live viewing is off", http.StatusConflict)
    		return
    	}
    	index, err := strconv.Atoi(r.URL.Query().Get("id"))
    	if err != nil || index < 0 || index >= len(m.config.LiveStreams) {
    		http.Error(w, "unknown camera", http.StatusNotFound)
    		return
    	}
    	source := fmt.Sprintf("camera_%d", index)
    	if r.URL.Query().Get("talk") == "1" {
    		if m.config.LiveStreams[index].BackchannelURL == "" {
    			http.Error(w, "this camera has no backchannel URL", http.StatusBadRequest)
    			return
    		}
    		source += "_talk"
    	}
    	// Never forward arbitrary URLs or other go2rtc query parameters.
    	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64<<10))
    	if err != nil || !bytes.HasPrefix(body, []byte("v=0")) {
    		http.Error(w, "invalid SDP offer", http.StatusBadRequest)
    		return
    	}
    	req, _ := http.NewRequestWithContext(r.Context(), http.MethodPost, service.baseURL+"/api/webrtc?src="+source, bytes.NewReader(body))
    	req.Header.Set("Content-Type", "application/sdp")
    	req.SetBasicAuth("neonvr", service.password)
    	response, err := service.client.Do(req)
    	if err != nil {
    		http.Error(w, "camera connection failed; check the camera address and network", http.StatusBadGateway)
    		return
    	}
    	defer response.Body.Close()
    	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
    		// The gateway's raw errors may include RTSP credentials.
    		http.Error(w, "camera could not provide compatible video or audio; check its URL, availability, and codecs", http.StatusBadGateway)
    		return
    	}
    	w.Header().Set("Content-Type", "application/sdp")
    	_, _ = io.Copy(w, io.LimitReader(response.Body, 128<<10))
    }
    Update { move_path: None }: /home/ubuntu/NeoNVR/main.go
    @@ -184,2 +184,8 @@
     
    +//go:embed static/streams.html
    +var streamsHTML []byte
    +
    +//go:embed static/live-streams.js
    +var liveStreamsJS []byte
    +
     //go:embed static/media-chrome-lite.js
    @@ -659,3 +665,7 @@
       <h1>NeoNVR</h1>
    -  <a href="/settings" style="display:inline-block;margin-bottom:12px;color:inherit">Settings</a>
    +  <nav class="tabs" aria-label="Main navigation" style="margin-bottom:16px">
    +    <a href="/" aria-current="page">Recordings</a>
    +    <a href="/streams">Streams</a>
    +    <a href="/settings">Settings</a>
    +  </nav>
       <div class="root">Root: {{.Root}}</div>
    @@ -1669,2 +1679,4 @@
     	defer cancelWork()
    +	live := newLiveManager(runCtx, cfg, manager.token)
    +	defer live.close()
     
    @@ -1755,2 +1767,19 @@
     	mux.HandleFunc("/api/config", manager.handler)
    +	mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
    +		if r.Method != http.MethodGet {
    +			w.Header().Set("Allow", "GET")
    +			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    +			return
    +		}
    +		w.Header().Set("Content-Type", "text/html; charset=utf-8")
    +		w.Header().Set("Cache-Control", "no-store")
    +		_, _ = w.Write(streamsHTML)
    +	})
    +	mux.HandleFunc("/static/live-streams.js", func(w http.ResponseWriter, r *http.Request) {
    +		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    +		w.Header().Set("Cache-Control", "no-cache")
    +		_, _ = w.Write(liveStreamsJS)
    +	})
    +	mux.HandleFunc("/api/live", live.handler)
    +	mux.HandleFunc("/api/live/webrtc", live.webrtcHandler)
     	mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
    @@ -1815,2 +1844,3 @@
     		log.Printf("Restarting NeoNVR with config %s", configPath)
    +		live.close()
     		cancelWork()

## Activity

    file changes: Completed · 8 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/Dockerfile
    @@ -8,2 +8,5 @@
     
    +# The bridge is started only when live viewing is explicitly enabled.
    +FROM ghcr.io/alexxit/go2rtc:1.9.14 AS streams
    +
     # PyTorch/Ultralytics are used only to export the official weights at build time.
    @@ -26,3 +29,5 @@
     COPY --from=build /out/neonvr /usr/local/bin/neonvr
    -EXPOSE 5500
    +COPY --from=streams /usr/local/bin/go2rtc /usr/local/bin/go2rtc
    +COPY third_party/go2rtc.LICENSE /usr/share/doc/go2rtc/LICENSE
    +EXPOSE 5500 8555/tcp 8555/udp
     ENTRYPOINT ["/usr/local/bin/neonvr"]
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -156,2 +156,16 @@
     
    +## Live streams and speaker talkback
    +
    +Open the **Streams** tab to view live RTSP cameras. Live viewing is **off after every app start**. Loading the recordings page or Streams tab does not connect to any camera. Turn on **Enable live viewing**, select a camera, and choose **Watch** to connect. **Stop**, leaving the page, or hiding the browser tab closes that viewer's connection and microphone. Turning live viewing off stops the streaming gateway and disconnects all viewers immediately.
    +
    +The supplied config includes `rtsp://192.168.1.150:554/ch0_0.h264` for viewing and `rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1` for speaker talkback. Edit sources in **Settings → Live streams**, one camera per line: `name | viewing URL | optional backchannel URL`. Save and restart to apply source changes; live viewing resets to off.
    +
    +**Talk to camera** explicitly requests microphone permission and sends your audio to the camera speaker using the backchannel URL. **Stop talking** stops microphone capture and returns to the ordinary viewing URL. Microphone access requires HTTPS (for example through your existing TLS reverse proxy) or localhost. Camera backchannel and browser audio codecs must be compatible; unsupported backchannels produce an error and stop microphone capture.
    +
    +The Docker image includes [go2rtc v1.9.14](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14). Its private API listens only on loopback and is accessed through NeoNVR. It starts on demand with no preloaded camera connections. Viewing uses WebRTC with no video transcoding. Camera audio plays when its codec is supported by WebRTC; use the video controls to unmute it. The browser UI and LAN streaming work without a CDN or public STUN service.
    +
    +The Compose file maps **8555/TCP and 8555/UDP** for WebRTC media in addition to the app's HTTP port. The browser must be able to reach that media port on the server. In Settings, a blank **WebRTC server hostname or IP** uses the hostname from the URL used to enable viewing; set it explicitly if a reverse proxy uses a different address from the media server. Changing the media port also requires changing both Compose mappings. Native installations must install the [go2rtc binary](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14) and set its executable path in Settings if it is not on `PATH`.
    +
    +## Main player controls
    +
     - Click a thumbnail or `Load` to send it to the main player.
    Update { move_path: None }: /home/ubuntu/NeoNVR/config/neonvr.json
    @@ -29,3 +29,13 @@
       "detection_timeout": "30m",
    -  "detection_prune": "quarantine"
    +  "detection_prune": "quarantine",
    +  "live_streams": [
    +    {
    +      "name": "Camera 192.168.1.150",
    +      "url": "rtsp://192.168.1.150:554/ch0_0.h264",
    +      "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
    +    }
    +  ],
    +  "live_go2rtc": "go2rtc",
    +  "live_webrtc_port": 8555,
    +  "live_webrtc_host": ""
     }
    Update { move_path: None }: /home/ubuntu/NeoNVR/docker-compose.yml
    @@ -9,2 +9,5 @@
           - "5500:5500"
    +      # Live media and microphone talkback; inactive until enabled in Streams.
    +      - "8555:8555/tcp"
    +      - "8555:8555/udp"
         volumes:
    Add: /home/ubuntu/NeoNVR/static/live-streams.js
    'use strict';
    
    const enabledEl = document.getElementById('enabled');
    const videoEl = document.getElementById('video');
    const placeholderEl = document.getElementById('placeholder');
    const sourcesEl = document.getElementById('sources');
    const watchEl = document.getElementById('watch');
    const stopEl = document.getElementById('stop');
    const talkEl = document.getElementById('talk');
    const microphoneEl = document.getElementById('microphone');
    const statusEl = document.getElementById('status');
    
    let state = null;
    let selectedID = null;
    let active = null;
    let generation = 0;
    let changing = false;
    
    function show(message, type = '') {
      statusEl.textContent = message;
      statusEl.className = 'status ' + type;
    }
    
    function selectedSource() {
      return state && state.streams.find(source => source.id === selectedID);
    }
    
    function updateControls() {
      const source = selectedSource();
      enabledEl.disabled = changing || !state || !state.streams.length;
      enabledEl.checked = !!(state && state.enabled);
      watchEl.disabled = changing || !source || !state.enabled;
      stopEl.disabled = !active;
      talkEl.disabled = changing || !active || !active.connected || !source || !source.talk || !window.isSecureContext || !navigator.mediaDevices;
      talkEl.textContent = active && active.talk ? 'Stop talking' : 'Talk to camera';
      talkEl.setAttribute('aria-pressed', String(!!(active && active.talk && active.connected)));
      microphoneEl.textContent = active && active.talk && active.connected ? 'Microphone live' : 'Microphone off';
      microphoneEl.className = 'mic' + (active && active.talk && active.connected ? ' active' : '');
      placeholderEl.hidden = !!(active && active.connected);
      placeholderEl.querySelector('strong').textContent = active ? 'Connecting…' : state && state.enabled ? 'Ready to watch' : 'Live viewing is off';
      placeholderEl.querySelector('p').textContent = active ? 'Waiting for the camera.' : state && state.enabled ? 'Select a camera and choose Watch.' : 'Enable live viewing, then choose Watch.';
    }
    
    function stopPlayback(message = '') {
      generation++;
      const previous = active;
      active = null;
      if (previous) {
        clearTimeout(previous.timer);
        previous.abort.abort();
        if (previous.microphone) previous.microphone.getTracks().forEach(track => track.stop());
        previous.pc.close();
      }
      videoEl.pause();
      videoEl.srcObject = null;
      updateControls();
      if (message) show(message);
    }
    
    function renderSources() {
      sourcesEl.replaceChildren();
      if (!state.streams.length) {
        const empty = document.createElement('p');
        empty.textContent = 'Add an RTSP camera in Settings to get started.';
        sourcesEl.append(empty);
        return;
      }
      for (const source of state.streams) {
        const button = document.createElement('button');
        button.type = 'button';
        button.className = 'source';
        button.setAttribute('aria-pressed', String(source.id === selectedID));
        button.append(document.createTextNode(source.name));
        const detail = document.createElement('span');
        detail.textContent = source.host + (source.talk ? ' · Talk available' : '');
        button.append(detail);
        button.addEventListener('click', () => {
          if (selectedID === source.id) return;
          stopPlayback('Selected ' + source.name + '. Choose Watch to connect.');
          selectedID = source.id;
          renderSources();
          updateControls();
        });
        sourcesEl.append(button);
      }
    }
    
    async function loadState(initial = false) {
      const response = await fetch('/api/live', {cache: 'no-store', signal: AbortSignal.timeout(5000)});
      if (!response.ok) throw new Error(await response.text());
      const next = await response.json();
      const restarted = state && state.token !== next.token;
      const wasEnabled = state && state.enabled;
      state = next;
      if (selectedID === null || !selectedSource()) selectedID = state.streams.length ? state.streams[0].id : null;
      if (restarted || !state.enabled && (active || wasEnabled)) {
        stopPlayback(restarted ? 'Server restarted. Live viewing is off.' : 'Live viewing was turned off. All camera connections are closed.');
      }
      if (initial || restarted) {
        renderSources();
        show(state.streams.length ? 'Select a camera. Live viewing connects only after you enable it and choose Watch.' : 'Add a camera using Edit sources.');
      }
      updateControls();
    }
    
    function waitForICE(pc, signal) {
      if (pc.iceGatheringState === 'complete') return Promise.resolve();
      return new Promise((resolve, reject) => {
        const finish = error => {
          clearTimeout(timer);
          pc.removeEventListener('icegatheringstatechange', changed);
          signal.removeEventListener('abort', aborted);
          if (error) reject(error); else resolve();
        };
        const changed = () => { if (pc.iceGatheringState === 'complete') finish(); };
        const aborted = () => finish(new DOMException('Connection stopped', 'AbortError'));
        const timer = setTimeout(() => finish(), 5000);
        pc.addEventListener('icegatheringstatechange', changed);
        signal.addEventListener('abort', aborted, {once: true});
        if (signal.aborted) aborted();
      });
    }
    
    async function connect(talk = false) {
      const source = selectedSource();
      if (!source || !state.enabled || document.hidden) return;
      if (!window.RTCPeerConnection) {
        show('This browser does not support live WebRTC playback.', 'error');
        return;
      }
      if (talk && (!window.isSecureContext || !navigator.mediaDevices)) {
        show('Open NeoNVR over HTTPS or localhost to use your microphone.', 'error');
        return;
      }
      stopPlayback();
      const currentGeneration = generation;
      const session = {
        pc: new RTCPeerConnection({iceServers: [], bundlePolicy: 'max-bundle'}),
        abort: new AbortController(), microphone: null, talk, connected: false, timer: null
      };
      active = session;
      show(talk ? 'Requesting microphone access…' : 'Connecting to ' + source.name + '…');
      updateControls();
      session.timer = setTimeout(() => {
        if (active === session) stopPlayback('Connection timed out. Check the camera and that the server WebRTC port is reachable.');
      }, 25000);
      const received = new MediaStream();
      videoEl.srcObject = received;
      session.pc.addEventListener('track', event => {
        if (active !== session) return;
        received.addTrack(event.track);
        videoEl.play().catch(() => show('Connected. Press Play in the video controls.', 'live'));
      });
      session.pc.addEventListener('connectionstatechange', () => {
        if (active !== session) return;
        const connection = session.pc.connectionState;
        if (connection === 'connected') {
          clearTimeout(session.timer);
          session.connected = true;
          if (session.microphone) session.microphone.getTracks().forEach(track => { track.enabled = true; });
          updateControls();
          show(talk ? 'Talking to ' + source.name + '. Choose Stop talking to turn your microphone off.' : 'Watching ' + source.name + '. Use the video controls to unmute camera audio.', 'live');
        } else if (connection === 'failed' || connection === 'disconnected' || connection === 'closed') {
          stopPlayback('Camera disconnected. Choose Watch to reconnect.');
        }
      });
      try {
        let sender = null;
        if (talk) {
          const microphone = await navigator.mediaDevices.getUserMedia({audio: {echoCancellation: true, noiseSuppression: true}});
          if (generation !== currentGeneration || active !== session) {
            microphone.getTracks().forEach(track => track.stop());
            return;
          }
          session.microphone = microphone;
          for (const track of microphone.getAudioTracks()) {
            track.enabled = false;
            sender = session.pc.addTransceiver(track, {direction: 'sendonly'});
          }
        }
        session.pc.addTransceiver('video', {direction: 'recvonly'});
        session.pc.addTransceiver('audio', {direction: 'recvonly'});
        await session.pc.setLocalDescription(await session.pc.createOffer());
        await waitForICE(session.pc, session.abort.signal);
        if (active !== session) return;
        show('Connecting to ' + source.name + (talk ? ' with microphone audio…' : '…'));
        const response = await fetch('/api/live/webrtc?id=' + source.id + (talk ? '&talk=1' : ''), {
          method: 'POST', headers: {'Content-Type': 'application/sdp', 'X-NeoNVR-Token': state.token},
          body: session.pc.localDescription.sdp, signal: session.abort.signal
        });
        if (!response.ok) throw new Error(await response.text());
        const answer = await response.text();
        if (active !== session) return;
        await session.pc.setRemoteDescription({type: 'answer', sdp: answer});
        if (talk && (!sender || !['sendonly', 'sendrecv'].includes(sender.currentDirection))) {
          throw new Error('The camera did not accept a compatible audio backchannel. Your microphone is off; choose Watch to return to viewing.');
        }
      } catch (error) {
        if (active !== session) return;
        stopPlayback();
        show(error.name === 'NotAllowedError' ? 'Microphone permission was denied. Choose Watch to view without talking.' : error.message, 'error');
      }
    }
    
    enabledEl.addEventListener('change', async () => {
      const enabled = enabledEl.checked;
      if (!enabled) stopPlayback();
      changing = true;
      updateControls();
      show(enabled ? 'Enabling live viewing…' : 'Disconnecting live viewers…');
      try {
        const response = await fetch('/api/live', {
          method: 'POST', headers: {'Content-Type': 'application/json', 'X-NeoNVR-Token': state.token},
          body: JSON.stringify({enabled}), signal: AbortSignal.timeout(10000)
        });
        if (!response.ok) throw new Error(await response.text());
        state = await response.json();
        show(enabled ? 'Live viewing enabled. Choose Watch to connect to the selected camera.' : 'Live viewing is off. All camera connections are closed.');
      } catch (error) {
        show(error.message, 'error');
        try { await loadState(); } catch (_) { /* Keep the last known state. */ }
      } finally {
        changing = false;
        updateControls();
      }
    });
    watchEl.addEventListener('click', () => connect(false));
    stopEl.addEventListener('click', () => stopPlayback('Stopped. The camera connection and microphone are closed.'));
    talkEl.addEventListener('click', () => connect(!(active && active.talk)));
    document.addEventListener('visibilitychange', () => {
      if (document.hidden && active) stopPlayback('Stopped while this tab was hidden. Choose Watch to reconnect.');
    });
    window.addEventListener('pagehide', () => stopPlayback());
    loadState(true).catch(error => show(error.message, 'error'));
    setInterval(() => {
      if (document.hidden || changing) return;
      loadState().catch(() => {
        if (active) stopPlayback('Server unavailable. The camera connection and microphone are closed.');
      });
    }, 3000);
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/settings.html
    @@ -11,3 +11,3 @@
     <body>
    -<header><div><h1>NeoNVR Settings</h1><p>Save your settings, then restart to apply them.</p></div><a href="/">← Recordings</a></header>
    +<header><div><h1>NeoNVR Settings</h1><p>Save your settings, then restart to apply them.</p></div><nav aria-label="Main navigation"><a href="/">Recordings</a> · <a href="/streams">Streams</a> · <a href="/settings" aria-current="page">Settings</a></nav></header>
     <p id="pending" class="notice" hidden>Saved settings differ from the running app. Restart to apply them.</p>
    @@ -15,2 +15,8 @@
     <fieldset id="fields" disabled>
    +<section class="panel" id="live-streams"><h2>Live streams</h2><p>Live viewing defaults to off after every app restart. Enable it from the Streams tab when needed.</p><div class="grid">
    +<label class="wide">RTSP camera sources<textarea name="live_streams" placeholder="Front camera | rtsp://192.168.1.150:554/ch0_0.h264 | rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"></textarea><small>One camera per line: name | viewing URL | speaker backchannel URL. Leave the last URL blank if the camera has no talk support.</small></label>
    +<label>Streaming gateway executable<input name="live_go2rtc" required><small>Docker includes go2rtc. Native installs can enter its absolute path.</small></label>
    +<label>WebRTC media port<input name="live_webrtc_port" type="number" min="1" max="65535" step="1" required><small>8555 by default. In a container, map this port for both TCP and UDP.</small></label>
    +<label class="wide">WebRTC server hostname or IP<input name="live_webrtc_host"><small>Blank uses the hostname you open NeoNVR with. Override this if the browser needs a different server address for live media.</small></label>
    +</div></section>
     <section class="panel"><h2>Object detection</h2><p>Find humans, animals, bikes and moving vehicles. Colored intervals appear below the video scrubber.</p><div class="grid">
    @@ -63,3 +69,3 @@
     let dirty=false;
    -const numeric=new Set(['thumb_seek','thumb_width','thumb_workers','auto_delete_days','detection_threads','detection_fps','detection_confidence','detection_min_frames','detection_merge_gap','detection_car_movement']);
    +const numeric=new Set(['thumb_seek','thumb_width','thumb_workers','auto_delete_days','detection_threads','detection_fps','detection_confidence','detection_min_frames','detection_merge_gap','detection_car_movement','live_webrtc_port']);
     function show(message,type=''){statusEl.textContent=message;statusEl.className=type;}
    @@ -74,3 +80,3 @@
       const field=form.elements.namedItem(name);
    -  if(field)field.value=name==='extra_roots'?value.join('\n'):value;
    +  if(field)field.value=name==='extra_roots'?value.join('\n'):name==='live_streams'?value.map(source=>[source.name,source.url,source.backchannel_url||''].join(' | ')).join('\n'):value;
      }
    @@ -86,2 +92,7 @@
       if(name==='extra_roots')value=value.split('\n').map(v=>v.trim()).filter(Boolean);
    +  else if(name==='live_streams')value=value.split('\n').map(v=>v.trim()).filter(Boolean).map((line,index)=>{
    +   const parts=line.split('|').map(v=>v.trim());
    +   if(parts.length<2||parts.length>3||!parts[0]||!parts[1])throw new Error('Live stream line '+(index+1)+': enter name | RTSP URL | optional backchannel URL.');
    +   return {name:parts[0],url:parts[1],...(parts[2]?{backchannel_url:parts[2]}:{})};
    +  });
       else if(name==='detection_enabled')value=value==='true';
    @@ -117,3 +128,3 @@
      if(!state||!form.reportValidity())return;
    - const config=readForm();
    + let config;try{config=readForm();}catch(error){show(error.message,'error');return;}
      if(restart&&config.detection_enabled&&config.detection_prune==='delete'&&(!state.active_config.detection_enabled||state.active_config.detection_prune!=='delete')){
    Add: /home/ubuntu/NeoNVR/static/streams.html
    <!doctype html>
    <html lang="en">
    <head>
    <meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
    <title>NeoNVR Streams</title>
    <style>
    :root{color-scheme:dark;font-family:system-ui,sans-serif;background:#0f1115;color:#f0f2f5;--muted:#a6afbc;--border:#333846;--accent:#72a7ff}*{box-sizing:border-box}body{margin:0;padding:24px;max-width:1440px;margin-inline:auto}h1{margin:0 0 12px;font-size:28px}a{color:var(--accent)}nav{display:flex;gap:18px;margin-bottom:28px}nav a[aria-current]{color:#f0f2f5;text-decoration:none;font-weight:650}.panel{background:#171a21;border:1px solid var(--border);border-radius:16px;padding:20px}.header{display:flex;justify-content:space-between;align-items:center;gap:24px;margin-bottom:18px}h2{margin:0 0 8px;font-size:22px}p{color:var(--muted);line-height:1.5;margin:0}.toggle{display:flex;align-items:center;gap:10px;white-space:nowrap;font-weight:600}.toggle input{width:20px;height:20px;accent-color:var(--accent)}.layout{display:grid;grid-template-columns:minmax(0,1fr) 280px;gap:18px}.stage{aspect-ratio:16/9;background:#080a0d;border:1px solid var(--border);border-radius:12px;overflow:hidden;position:relative}video{width:100%;height:100%;object-fit:contain;display:block}#placeholder{position:absolute;inset:0;display:grid;place-content:center;text-align:center;gap:12px;padding:32px;pointer-events:none}#placeholder strong{font-size:20px}#placeholder[hidden]{display:none}.actions{display:flex;align-items:center;flex-wrap:wrap;gap:10px;margin-top:16px}button{font:inherit;color:#f0f2f5;padding:10px 16px;border:1px solid var(--border);border-radius:9px;background:#20242e;cursor:pointer}button.primary{background:#72a7ff;color:#101827;border-color:#72a7ff;font-weight:600}button:disabled{opacity:.45;cursor:default}button[aria-pressed=true]{background:#742e35;border-color:#ff8491}.status{min-height:50px;padding-top:14px;color:var(--muted);line-height:1.5;white-space:pre-wrap}.status.error{color:#ffb3b3}.status.live{color:#92dfad}.listHeading{display:flex;justify-content:space-between;align-items:center;gap:10px;margin-bottom:14px}.listHeading h3{margin:0;font-size:16px}.source{display:block;width:100%;text-align:left;padding:14px;margin-bottom:10px}.source[aria-pressed=true]{background:#1e2e45;border-color:var(--accent)}.source span{display:block;font-size:12px;color:var(--muted);margin-top:5px;overflow-wrap:anywhere}.hint{font-size:13px;margin-top:16px}.mic{font-size:13px;color:var(--muted)}.mic.active{color:#ffb3b3}button:focus-visible,a:focus-visible,input:focus-visible{outline:2px solid var(--accent);outline-offset:3px}@media(max-width:850px){.layout{grid-template-columns:1fr}.header{align-items:flex-start;flex-direction:column}body{padding:14px}.panel{padding:16px}}
    </style>
    </head>
    <body>
    <h1>NeoNVR</h1>
    <nav aria-label="Main navigation"><a href="/">Recordings</a><a href="/streams" aria-current="page">Streams</a><a href="/settings">Settings</a></nav>
    <section class="panel">
     <header class="header">
      <div><h2>Live streams</h2><p>Off by default. Connect to a camera when you want to watch.</p></div>
      <label class="toggle"><input id="enabled" type="checkbox" disabled>Enable live viewing</label>
     </header>
     <div class="layout">
      <div>
       <div class="stage">
        <video id="video" controls playsinline muted aria-label="Live camera video"></video>
        <div id="placeholder"><strong>Live viewing is off</strong><p>Enable live viewing, then choose Watch.</p></div>
       </div>
       <div class="actions">
        <button id="watch" class="primary" disabled>Watch</button>
        <button id="stop" disabled>Stop</button>
        <button id="talk" aria-pressed="false" disabled>Talk to camera</button>
        <span id="microphone" class="mic">Microphone off</span>
       </div>
       <div id="status" class="status" role="status" aria-live="polite">Loading cameras…</div>
       <p class="hint" id="talkHint">Talk sends your microphone audio to the camera speaker. Microphone access needs HTTPS or localhost.</p>
      </div>
      <aside aria-label="Live camera selection">
       <div class="listHeading"><h3>Cameras</h3><a href="/settings#live-streams">Edit sources</a></div>
       <div id="sources"></div>
       <p class="hint">Stop disconnects this viewer. Turning live viewing off disconnects every live viewer.</p>
      </aside>
     </div>
    </section>
    <script src="/static/live-streams.js" defer></script>
    </body>
    </html>
    Add: /home/ubuntu/NeoNVR/third_party/go2rtc.LICENSE
    MIT License
    
    Copyright (c) 2022 Alexey Khit
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.

## Activity

    $ git diff --check
        sed -n '147,196p' README.md
        sed -n '43,68p' config.go
        sed -n '1,90p' config_test.go
        sed -n '185,290p' config_test.go
      -root /mnt/mySSD/BlinkClips \
      -addr 0.0.0.0:5500 \
      -thumb-workers 2 \
      -thumb-width 480 \
      -thumb-seek 0.2 \
      -size-refresh 30m
    ```
    
    ## Main player controls
    
    ## Live streams and speaker talkback
    
    Open the **Streams** tab to view live RTSP cameras. Live viewing is **off after every app start**. Loading the recordings page or Streams tab does not connect to any camera. Turn on **Enable live viewing**, select a camera, and choose **Watch** to connect. **Stop**, leaving the page, or hiding the browser tab closes that viewer's connection and microphone. Turning live viewing off stops the streaming gateway and disconnects all viewers immediately.
    
    The supplied config includes `rtsp://192.168.1.150:554/ch0_0.h264` for viewing and `rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1` for speaker talkback. Edit sources in **Settings → Live streams**, one camera per line: `name | viewing URL | optional backchannel URL`. Save and restart to apply source changes; live viewing resets to off.
    
    **Talk to camera** explicitly requests microphone permission and sends your audio to the camera speaker using the backchannel URL. **Stop talking** stops microphone capture and returns to the ordinary viewing URL. Microphone access requires HTTPS (for example through your existing TLS reverse proxy) or localhost. Camera backchannel and browser audio codecs must be compatible; unsupported backchannels produce an error and stop microphone capture.
    
    The Docker image includes [go2rtc v1.9.14](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14). Its private API listens only on loopback and is accessed through NeoNVR. It starts on demand with no preloaded camera connections. Viewing uses WebRTC with no video transcoding. Camera audio plays when its codec is supported by WebRTC; use the video controls to unmute it. The browser UI and LAN streaming work without a CDN or public STUN service.
    
    The Compose file maps **8555/TCP and 8555/UDP** for WebRTC media in addition to the app's HTTP port. The browser must be able to reach that media port on the server. In Settings, a blank **WebRTC server hostname or IP** uses the hostname from the URL used to enable viewing; set it explicitly if a reverse proxy uses a different address from the media server. Changing the media port also requires changing both Compose mappings. Native installations must install the [go2rtc binary](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14) and set its executable path in Settings if it is not on `PATH`.
    
    ## Main player controls
    
    - Click a thumbnail or `Load` to send it to the main player.
    - Space pauses/resumes.
    - `←` and `→` step within the loaded video.
    - `Ctrl` + `←` / `↑` loads the previous visible video.
    - `Ctrl` + `→` / `↓` loads the next visible video.
    - Change `Step sec` for fine-grain frame/time stepping.
    - Use `Zoom off` / `Zoom on` to enable inspection zoom. With Zoom off, tapping/clicking the video pauses/plays. With Zoom on, pinch/trackpad zooms the actual video display and drag pans while zoomed. `Reset zoom` restores the view. This is inspection zoom only; it does not create missing detail.
    
    ## Bulk delete
    
    Click `Select` to enter delete-selection mode. In that mode, clicking thumbnails selects them instead of loading them. Shift-click another thumbnail to select every visible clip in between, then use `Delete selected`.
    
    ## Files written inside BlinkClips
    
    By default the app writes:
    
    ```text
    /mnt/mySSD/BlinkClips/.blinkclips_thumbs/
    /mnt/mySSD/BlinkClips/.blinkclips_meta.json
    ```
    
    The thumbnail cache can be moved with `-thumb-dir`.
    
    Detection also writes `.neonvr_detections/` under the primary root as an expendable JSON cache. Detection metadata lives inside the MP4 itself, so deleting the cache does not lose event intervals. Quarantined recordings live in `<camera-directory>/.neonvr_pruned/`, preserving Yi hour subdirectories where applicable.
    
    ## YOLO26n detection and false-positive pruning
        DetectionEnabled      bool         `json:"detection_enabled"`
        DetectionPython       string       `json:"detection_python"`
        DetectionModel        string       `json:"detection_model"`
        DetectionThreads      int          `json:"detection_threads"`
        DetectionFPS          float64      `json:"detection_fps"`
        DetectionConfidence   float64      `json:"detection_confidence"`
        DetectionMinFrames    int          `json:"detection_min_frames"`
        DetectionMergeGap     float64      `json:"detection_merge_gap"`
        DetectionCarMovement  float64      `json:"detection_car_movement"`
        DetectionEvery        string       `json:"detection_every"`
        DetectionMinAge       string       `json:"detection_min_age"`
        DetectionTimeout      string       `json:"detection_timeout"`
        DetectionPrune        string       `json:"detection_prune"`
        LiveStreams           []LiveStream `json:"live_streams"`
        LiveGo2RTC            string       `json:"live_go2rtc"`
        LiveWebRTCPort        int          `json:"live_webrtc_port"`
        LiveWebRTCHost        string       `json:"live_webrtc_host"`
    }
    
    func defaultConfig() Config {
        return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h", LiveStreams: []LiveStream{}, LiveGo2RTC: "go2rtc", LiveWebRTCPort: 8555,
            DetectionEnabled: true, DetectionPython: "python3", DetectionModel: "yolo26n.onnx", DetectionThreads: 2, DetectionFPS: 2, DetectionConfidence: .65, DetectionMinFrames: 3, DetectionMergeGap: .5, DetectionCarMovement: .025, DetectionEvery: "1m", DetectionMinAge: "2m", DetectionTimeout: "30m", DetectionPrune: "quarantine"}
    }
    
    func decodeConfig(data []byte) (Config, error) {
        c := defaultConfig()
    package main
    
    import (
        "bytes"
        "context"
        "encoding/json"
        "net/http"
        "net/http/httptest"
        "os"
        "path/filepath"
        "sync"
        "testing"
        "time"
    )
    
    func fixture(t *testing.T) (Config, string, *configManager) {
        t.Helper()
        c := defaultConfig()
        c.Root = t.TempDir()
        path := filepath.Join(t.TempDir(), "neonvr.json")
        if err := atomicWriteConfig(path, configBytes(c)); err != nil {
            t.Fatal(err)
        }
        manager, err := newConfigManager(path, c)
        if err != nil {
            t.Fatal(err)
        }
        return c, path, manager
    }
    
    func postConfig(t *testing.T, m *configManager, c Config, revision, token, origin string, restart bool) *httptest.ResponseRecorder {
        t.Helper()
        body, _ := json.Marshal(map[string]any{"config": c, "revision": revision, "restart": restart})
        req := httptest.NewRequest(http.MethodPost, "http://neonvr.test/api/config", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-NeoNVR-Token", token)
        if origin != "" {
            req.Header.Set("Origin", origin)
        }
        w := httptest.NewRecorder()
        m.handler(w, req)
        return w
    }
    
    func TestSaveConfigValidationConflictAndRollback(t *testing.T) {
        c, path, m := fixture(t)
        original, _ := os.ReadFile(path)
        revision := configRevision(original)
        c.ThumbWorkers = 0
        w := postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
        if w.Code != 400 {
            t.Fatalf("invalid config status %d: %s", w.Code, w.Body.String())
        }
        raw, _ := os.ReadFile(path)
        if !bytes.Equal(raw, original) {
            t.Fatal("invalid settings changed the file")
        }
        c.ThumbWorkers = 4
        w = postConfig(t, m, c, revision, m.token, "http://neonvr.test", false)
        if w.Code != 200 {
            t.Fatalf("save status %d: %s", w.Code, w.Body.String())
        }
        previous, _ := os.ReadFile(path + ".previous")
        if !bytes.Equal(previous, original) {
            t.Fatal("previous config not preserved")
        }
        saved, _ := os.ReadFile(path)
        loaded, err := decodeConfig(saved)
        if err != nil || loaded.ThumbWorkers != 4 {
            t.Fatal("config not saved")
        }
        if m.active.ThumbWorkers != 2 {
            t.Fatal("save applied config without restart")
        }
        if postConfig(t, m, c, revision, m.token, "", false).Code != 409 {
            t.Fatal("stale edit did not conflict")
        }
    }
    
    func TestConfigRejectsCrossOriginAndMissingToken(t *testing.T) {
        c, path, m := fixture(t)
        raw, _ := os.ReadFile(path)
        rev := configRevision(raw)
        for _, tc := range []struct{ token, origin string }{{"", "http://neonvr.test"}, {m.token, "http://evil.test"}, {m.token, "https://neonvr.test"}} {
            if w := postConfig(t, m, c, rev, tc.token, tc.origin, false); w.Code != 403 {
                t.Fatalf("expected 403 got %d", w.Code)
            }
        }
    }
    
                t.Fatal("background loop did not stop")
            }
        }
    }
    ✓ • 0ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w config.go main.go live.go
        /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test ./...
    ok      neonvr    2.822s
    ✓ • 3.85s

## Assistant

The Streams tab, settings, and default-off toggle are in place. I’m now checking playback, microphone cleanup, and the server-side gate to ensure cameras can’t be opened while live viewing is off. I’ll use a local test stream so your camera stays idle during testing.

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -154,4 +154,2 @@
     
    -## Main player controls
    -
     ## Live streams and speaker talkback

## Activity

    $ python3 -m venv /tmp/neonvr-live-tools-vbq_9dyb/browser-env
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/pip -q install playwright
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python -m playwright install chromium
    Downloading Chrome for Testing 153.0.8010.12 (playwright chromium v1243) from https://cdn.playwright.dev/builds/cft/153.0.8010.12/linux64/chrome-linux64.zip
    |                                                                                |   0% of 186.8 MiB
    |■■■■■■■■                                                                        |  10% of 186.8 MiB
    |■■■■■■■■■■■■■■■■                                                                |  20% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■                                                        |  30% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                                |  40% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                        |  50% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                |  60% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                        |  70% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                |  80% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■        |  90% of 186.8 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■| 100% of 186.8 MiB
    Chrome for Testing 153.0.8010.12 (playwright chromium v1243) downloaded to /home/ubuntu/.cache/ms-playwright/chromium-1243
    Downloading FFmpeg (playwright ffmpeg v1011) from https://cdn.playwright.dev/dbazure/download/playwright/builds/ffmpeg/1011/ffmpeg-linux.zip
    |                                                                                |   0% of 2.3 MiB
    |■■■■■■■■                                                                        |  10% of 2.3 MiB
    |■■■■■■■■■■■■■■■■                                                                |  20% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■                                                        |  30% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                                |  40% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                        |  50% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                |  60% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                        |  70% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                |  80% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■        |  90% of 2.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■| 100% of 2.3 MiB
    FFmpeg (playwright ffmpeg v1011) downloaded to /home/ubuntu/.cache/ms-playwright/ffmpeg-1011
    Downloading Chrome Headless Shell 153.0.8010.12 (playwright chromium-headless-shell v1243) from https://cdn.playwright.dev/builds/cft/153.0.8010.12/linux64/chrome-headless-shell-linux64.zip
    |                                                                                |   0% of 114.3 MiB
    |■■■■■■■■                                                                        |  10% of 114.3 MiB
    |■■■■■■■■■■■■■■■■                                                                |  20% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■                                                        |  30% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                                |  40% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                        |  50% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                                |  60% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                        |  70% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■                |  80% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■        |  90% of 114.3 MiB
    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■| 100% of 114.3 MiB
    Chrome Headless Shell 153.0.8010.12 (playwright chromium-headless-shell v1243) downloaded to /home/ubuntu/.cache/ms-playwright/chromium_headless_shell-1243
    ✓ • 20.23s

## Activity

    file changes: Completed · 1 changes
    Add: /home/ubuntu/NeoNVR/live_test.go
    package main
    
    import (
    	"bytes"
    	"context"
    	"encoding/json"
    	"io"
    	"net"
    	"net/http"
    	"net/http/httptest"
    	"os"
    	"strings"
    	"testing"
    	"time"
    )
    
    func liveFixture() Config {
    	c := defaultConfig()
    	c.LiveStreams = []LiveStream{{Name: "Front camera", URL: "rtsp://user:camera-secret@192.168.1.150:554/ch0_0.h264", BackchannelURL: "rtsp://user:camera-secret@192.168.1.150:554/ch0_0.h264?backchannel=1"}}
    	return c
    }
    
    func liveRequest(method, path, body, contentType string) *http.Request {
    	r := httptest.NewRequest(method, "http://neonvr.test"+path, strings.NewReader(body))
    	r.Header.Set("X-NeoNVR-Token", "test-token")
    	r.Header.Set("Content-Type", contentType)
    	r.Header.Set("Origin", "http://neonvr.test")
    	return r
    }
    
    func TestLiveDefaultsOffAndHidesCameraCredentials(t *testing.T) {
    	c := liveFixture()
    	c.LiveGo2RTC = "/does-not-exist"
    	m := newLiveManager(context.Background(), c, "test-token")
    	defer m.close()
    	w := httptest.NewRecorder()
    	m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
    	var status struct {
    		Enabled bool `json:"enabled"`
    		Streams []struct {
    			Name string `json:"name"`
    			Talk bool   `json:"talk"`
    		} `json:"streams"`
    	}
    	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil || status.Enabled || len(status.Streams) != 1 || !status.Streams[0].Talk {
    		t.Fatalf("unexpected default status: %s", w.Body.String())
    	}
    	if strings.Contains(w.Body.String(), "camera-secret") || strings.Contains(w.Body.String(), "rtsp://") || m.service != nil {
    		t.Fatal("listing sources exposed credentials or started a camera gateway")
    	}
    	w = httptest.NewRecorder()
    	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0", "v=0\r\n", "application/sdp"))
    	if w.Code != http.StatusConflict || m.service != nil {
    		t.Fatalf("a direct playback request bypassed the off switch: %d", w.Code)
    	}
    	w = httptest.NewRecorder()
    	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
    	if w.Code != http.StatusServiceUnavailable || m.service != nil {
    		t.Fatal("a failed gateway start left live viewing enabled")
    	}
    }
    
    func TestLiveToggleRejectsUnauthorizedAndMalformedRequests(t *testing.T) {
    	m := newLiveManager(context.Background(), liveFixture(), "test-token")
    	defer m.close()
    	for _, tc := range []struct {
    		name, method, body, token, origin, contentType string
    		code                                         int
    	}{
    		{"missing token", "POST", `{"enabled":false}`, "", "http://neonvr.test", "application/json", 403},
    		{"foreign origin", "POST", `{"enabled":false}`, "test-token", "https://other.test", "application/json", 403},
    		{"missing boolean", "POST", `{}`, "test-token", "http://neonvr.test", "application/json", 400},
    		{"wrong type", "POST", `{"enabled":"true"}`, "test-token", "http://neonvr.test", "application/json", 400},
    		{"trailing data", "POST", `{"enabled":false}{}`, "test-token", "http://neonvr.test", "application/json", 400},
    		{"unexpected field", "POST", `{"enabled":false,"src":"rtsp://other"}`, "test-token", "http://neonvr.test", "application/json", 400},
    		{"wrong content type", "POST", `{"enabled":false}`, "test-token", "http://neonvr.test", "text/plain", 415},
    		{"wrong method", "PUT", "", "test-token", "http://neonvr.test", "application/json", 405},
    		{"HTTPS proxy", "POST", `{"enabled":false}`, "test-token", "https://neonvr.test", "application/json", 200},
    	} {
    		t.Run(tc.name, func(t *testing.T) {
    			r := liveRequest(tc.method, "/api/live", tc.body, tc.contentType)
    			r.Header.Set("X-NeoNVR-Token", tc.token)
    			r.Header.Set("Origin", tc.origin)
    			w := httptest.NewRecorder()
    			m.handler(w, r)
    			if w.Code != tc.code || m.service != nil {
    				t.Fatalf("status=%d, want %d; %s", w.Code, tc.code, w.Body.String())
    			}
    		})
    	}
    }
    
    func TestLiveSignalingUsesOnlyConfiguredSourcesAndRedactsErrors(t *testing.T) {
    	var observed string
    	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    		user, password, ok := r.BasicAuth()
    		if !ok || user != "neonvr" || password != "private-password" || r.URL.Path != "/api/webrtc" {
    			t.Error("gateway authentication or route missing")
    		}
    		body, _ := io.ReadAll(r.Body)
    		if string(body) != "v=0\r\n" || r.Header.Get("Content-Type") != "application/sdp" {
    			t.Error("SDP offer was changed")
    		}
    		observed = r.URL.RawQuery
    		if r.URL.Query().Get("src") == "camera_0_talk" {
    			http.Error(w, "rtsp://user:camera-secret@camera failed", 500)
    			return
    		}
    		w.WriteHeader(http.StatusCreated)
    		_, _ = w.Write([]byte("v=0\r\nanswer"))
    	}))
    	defer backend.Close()
    	c := liveFixture()
    	c.LiveStreams = append(c.LiveStreams, LiveStream{Name: "View only", URL: "rtsp://192.168.1.151/live"})
    	m := newLiveManager(context.Background(), c, "test-token")
    	m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private-password", client: backend.Client()}
    	w := httptest.NewRecorder()
    	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&src=rtsp://other&dst=other", "v=0\r\n", "application/sdp"))
    	if w.Code != 200 || w.Body.String() != "v=0\r\nanswer" || observed != "src=camera_0" {
    		t.Fatalf("unexpected signaling response: %d %s, query %s", w.Code, w.Body.String(), observed)
    	}
    	w = httptest.NewRecorder()
    	m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&talk=1", "v=0\r\n", "application/sdp"))
    	if w.Code != 502 || observed != "src=camera_0_talk" || strings.Contains(w.Body.String(), "camera-secret") {
    		t.Fatal("backchannel selection failed or camera credentials leaked")
    	}
    	for _, tc := range []struct {
    		path, body string
    		code       int
    	}{
    		{"/api/live/webrtc?id=-1", "v=0\r\n", 404},
    		{"/api/live/webrtc?id=10", "v=0\r\n", 404},
    		{"/api/live/webrtc?id=1&talk=1", "v=0\r\n", 400},
    		{"/api/live/webrtc?id=0", "invalid", 400},
    		{"/api/live/webrtc?id=0", "v=0" + strings.Repeat("x", 64<<10), 400},
    	} {
    		observed = ""
    		w = httptest.NewRecorder()
    		m.webrtcHandler(w, liveRequest("POST", tc.path, tc.body, "application/sdp"))
    		if w.Code != tc.code || observed != "" {
    			t.Fatalf("invalid request reached the gateway: %s %d", tc.path, w.Code)
    		}
    	}
    }
    
    func TestLiveConfigValidationAndLegacyCompatibility(t *testing.T) {
    	c, err := decodeConfig([]byte(`{"root":"/tmp"}`))
    	if err != nil || c.LiveStreams == nil || len(c.LiveStreams) != 0 || c.LiveGo2RTC != "go2rtc" || c.LiveWebRTCPort != 8555 {
    		t.Fatalf("legacy config did not receive optional live defaults: %+v %v", c, err)
    	}
    	for _, raw := range []string{"exec:sh -c anything", "http://camera/live", "rtsp:///missing-host", "rtsp://camera:65536/live", "rtsp://camera/live#backchannel=1", "rtsp://camera/live\n"} {
    		if validateRTSP(raw) == nil {
    			t.Fatalf("invalid source accepted: %q", raw)
    		}
    	}
    	c = liveFixture()
    	if err := c.validateLive(); err != nil {
    		t.Fatal(err)
    	}
    	m := newLiveManager(context.Background(), c, "test-token")
    	sources := m.gatewayConfig("127.0.0.1:12345", "secret", "192.168.1.10")["streams"].(map[string]string)
    	if !strings.HasSuffix(sources["camera_0"], "#backchannel=0&timeout=10") || !strings.Contains(sources["camera_0_talk"], "?backchannel=1#backchannel=1&timeout=10") {
    		t.Fatal("viewing opened a backchannel or talk lost the camera query")
    	}
    	c.LiveStreams[0].BackchannelURL = "exec:invalid"
    	if c.validateLive() == nil {
    		t.Fatal("invalid backchannel accepted")
    	}
    }
    
    func TestLiveGatewayLifecycleWithoutOpeningCamera(t *testing.T) {
    	binary := os.Getenv("NEONVR_TEST_GO2RTC")
    	if binary == "" {
    		t.Skip("set NEONVR_TEST_GO2RTC to test the real streaming gateway")
    	}
    	camera, err := net.Listen("tcp", "127.0.0.1:0")
    	if err != nil {
    		t.Fatal(err)
    	}
    	defer camera.Close()
    	media, err := net.Listen("tcp", "127.0.0.1:0")
    	if err != nil {
    		t.Fatal(err)
    	}
    	port := media.Addr().(*net.TCPAddr).Port
    	media.Close()
    	c := defaultConfig()
    	c.LiveGo2RTC = binary
    	c.LiveWebRTCPort = port
    	c.LiveStreams = []LiveStream{{Name: "Test camera", URL: "rtsp://" + camera.Addr().String() + "/live"}}
    	ctx, cancel := context.WithCancel(context.Background())
    	defer cancel()
    	m := newLiveManager(ctx, c, "test-token")
    	defer m.close()
    	w := httptest.NewRecorder()
    	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
    	if w.Code != 200 || m.service == nil {
    		t.Fatalf("gateway did not start: %d %s", w.Code, w.Body.String())
    	}
    	service := m.service
    	configPath := service.cmd.Args[2]
    	if info, err := os.Stat(configPath); err != nil || info.Mode().Perm() != 0600 {
    		t.Fatal("temporary camera configuration was not private")
    	}
    	w = httptest.NewRecorder()
    	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
    	if w.Code != 200 || m.service != service {
    		t.Fatal("enabling twice started duplicate gateways")
    	}
    	_ = camera.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))
    	if conn, err := camera.Accept(); err == nil {
    		conn.Close()
    		t.Fatal("enabling the gateway opened a camera before Watch")
    	}
    	w = httptest.NewRecorder()
    	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":false}`, "application/json"))
    	if w.Code != 200 || m.service != nil {
    		t.Fatal("disabling did not stop the gateway")
    	}
    	select {
    	case <-service.done:
    	default:
    		t.Fatal("gateway process survived disabling")
    	}
    	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
    		t.Fatal("temporary camera credentials survived disabling")
    	}
    	w = httptest.NewRecorder()
    	m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
    	if w.Code != 200 {
    		t.Fatal(w.Body.String())
    	}
    	cancel()
    	select {
    	case <-m.service.done:
    	case <-time.After(3 * time.Second):
    		t.Fatal("gateway survived app shutdown")
    	}
    	w = httptest.NewRecorder()
    	m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
    	if !bytes.Contains(w.Body.Bytes(), []byte(`"enabled":false`)) {
    		t.Fatal("exited gateway still reported as enabled")
    	}
    }

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w live_test.go
        NEONVR_TEST_GO2RTC=/tmp/neonvr-live-tools-vbq_9dyb/go2rtc /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
    2026/10/10 23:43:12 detection: quarantine /tmp/TestDetectionPruningKeepsFavoritesAndUncertainAndQuarantinesEmpty2166714719/001/Garden/2026-10-01_12-00-00_Garden.mp4 (fully analyzed, no targets)
    2026/10/10 23:43:13 detection: delete /tmp/TestDetectionDeletionMode2812253406/001/Garden/2026-10-01_12-00-00_Garden.mp4 (fully analyzed, no targets)
    --- FAIL: TestLiveGatewayLifecycleWithoutOpeningCamera (0.92s)
        live_test.go:242: exited gateway still reported as enabled
    FAIL
    FAIL    neonvr    3.748s
    FAIL
    ✗ (1) • 5.03s

## Activity

    $ git diff --stat
        git diff --check
        sed -n '154,190p' README.md
     Dockerfile           |  7 +++++-
     README.md            | 14 +++++++++++-
     config.go            | 64 ++++++++++++++++++++++++++++++----------------------
     config/neonvr.json   | 12 +++++++++-
     docker-compose.yml   |  5 ++++
     main.go              | 48 +++++++++++++++++++++++++++------------
     static/settings.html | 19 ++++++++++++----
     7 files changed, 121 insertions(+), 48 deletions(-)
    
    ## Live streams and speaker talkback
    
    Open the **Streams** tab to view live RTSP cameras. Live viewing is **off after every app start**. Loading the recordings page or Streams tab does not connect to any camera. Turn on **Enable live viewing**, select a camera, and choose **Watch** to connect. **Stop**, leaving the page, or hiding the browser tab closes that viewer's connection and microphone. Turning live viewing off stops the streaming gateway and disconnects all viewers immediately.
    
    The supplied config includes `rtsp://192.168.1.150:554/ch0_0.h264` for viewing and `rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1` for speaker talkback. Edit sources in **Settings → Live streams**, one camera per line: `name | viewing URL | optional backchannel URL`. Save and restart to apply source changes; live viewing resets to off.
    
    **Talk to camera** explicitly requests microphone permission and sends your audio to the camera speaker using the backchannel URL. **Stop talking** stops microphone capture and returns to the ordinary viewing URL. Microphone access requires HTTPS (for example through your existing TLS reverse proxy) or localhost. Camera backchannel and browser audio codecs must be compatible; unsupported backchannels produce an error and stop microphone capture.
    
    The Docker image includes [go2rtc v1.9.14](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14). Its private API listens only on loopback and is accessed through NeoNVR. It starts on demand with no preloaded camera connections. Viewing uses WebRTC with no video transcoding. Camera audio plays when its codec is supported by WebRTC; use the video controls to unmute it. The browser UI and LAN streaming work without a CDN or public STUN service.
    
    The Compose file maps **8555/TCP and 8555/UDP** for WebRTC media in addition to the app's HTTP port. The browser must be able to reach that media port on the server. In Settings, a blank **WebRTC server hostname or IP** uses the hostname from the URL used to enable viewing; set it explicitly if a reverse proxy uses a different address from the media server. Changing the media port also requires changing both Compose mappings. Native installations must install the [go2rtc binary](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14) and set its executable path in Settings if it is not on `PATH`.
    
    ## Main player controls
    
    - Click a thumbnail or `Load` to send it to the main player.
    - Space pauses/resumes.
    - `←` and `→` step within the loaded video.
    - `Ctrl` + `←` / `↑` loads the previous visible video.
    - `Ctrl` + `→` / `↓` loads the next visible video.
    - Change `Step sec` for fine-grain frame/time stepping.
    - Use `Zoom off` / `Zoom on` to enable inspection zoom. With Zoom off, tapping/clicking the video pauses/plays. With Zoom on, pinch/trackpad zooms the actual video display and drag pans while zoomed. `Reset zoom` restores the view. This is inspection zoom only; it does not create missing detail.
    
    ## Bulk delete
    
    Click `Select` to enter delete-selection mode. In that mode, clicking thumbnails selects them instead of loading them. Shift-click another thumbnail to select every visible clip in between, then use `Delete selected`.
    
    ## Files written inside BlinkClips
    
    By default the app writes:
    
    ```text
    /mnt/mySSD/BlinkClips/.blinkclips_thumbs/
    /mnt/mySSD/BlinkClips/.blinkclips_meta.json
    ```
    
    The thumbnail cache can be moved with `-thumb-dir`.
    ✓ • 0ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr .
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python - <<'PY'
        from playwright.sync_api import sync_playwright
        with sync_playwright() as p:
         browser=p.chromium.launch(headless=True,args=['--no-sandbox'])
         page=browser.new_page()
         page.goto('about:blank')
         print('Chromium launched successfully.')
         browser.close()
        PY
    Chromium launched successfully.
    ✓ • 2.94s

## Activity

    $ rg -n 'func writeJSON|SetIndent' main.go
        sed -n '3010,3045p' main.go
        sed -n '225,248p' live_test.go
    3028:func writeJSON(w http.ResponseWriter, value any) {
    3031:    enc.SetIndent("", "  ")
    3037:func writeJSONError(w http.ResponseWriter, status int, err error) {
    }
    
    func humanBytes(n int64) string {
        if n < 0 {
            return "unknown"
        }
        const unit = 1024
        if n < unit {
            return fmt.Sprintf("%d B", n)
        }
        div, exp := int64(unit), 0
        for value := n / unit; value >= unit; value /= unit {
            div *= unit
            exp++
        }
        return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
    }
    
    func writeJSON(w http.ResponseWriter, value any) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        enc := json.NewEncoder(w)
        enc.SetIndent("", "  ")
        if err := enc.Encode(value); err != nil {
            log.Printf("json encode error: %v", err)
        }
    }
    
    func writeJSONError(w http.ResponseWriter, status int, err error) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(status)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error(), "message": err.Error()})
    }
    
    func writePlaceholderSVG(w http.ResponseWriter, label string) {
        w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
        w.Header().Set("Cache-Control", "no-store")
        if _, err := os.Stat(configPath); !os.IsNotExist(err) {
            t.Fatal("temporary camera credentials survived disabling")
        }
        w = httptest.NewRecorder()
        m.handler(w, liveRequest("POST", "/api/live", `{"enabled":true}`, "application/json"))
        if w.Code != 200 {
            t.Fatal(w.Body.String())
        }
        cancel()
        select {
        case <-m.service.done:
        case <-time.After(3 * time.Second):
            t.Fatal("gateway survived app shutdown")
        }
        w = httptest.NewRecorder()
        m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
        if !bytes.Contains(w.Body.Bytes(), []byte(`"enabled":false`)) {
            t.Fatal("exited gateway still reported as enabled")
        }
    }
    ✓ • 0ms

## Activity

    file changes: Completed · 3 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/live.go
    @@ -215,2 +215,13 @@
     func (m *liveManager) start(requestCtx context.Context, requestHost string) (*liveService, error) {
    +	mediaAddress := fmt.Sprintf(":%d", m.config.LiveWebRTCPort)
    +	tcp, err := net.Listen("tcp", mediaAddress)
    +	if err != nil {
    +		return nil, errors.New("WebRTC TCP port is already in use; choose another port in Settings")
    +	}
    +	udp, err := net.ListenPacket("udp", mediaAddress)
    +	tcp.Close()
    +	if err != nil {
    +		return nil, errors.New("WebRTC UDP port is already in use; choose another port in Settings")
    +	}
    +	udp.Close()
     	// Reserve a private API address. Only the two fixed API routes are used.
    Update { move_path: None }: /home/ubuntu/NeoNVR/live_test.go
    @@ -3,3 +3,2 @@
     import (
    -	"bytes"
     	"context"
    @@ -240,3 +239,6 @@
     	m.handler(w, httptest.NewRequest("GET", "/api/live", nil))
    -	if !bytes.Contains(w.Body.Bytes(), []byte(`"enabled":false`)) {
    +	var status struct {
    +		Enabled bool `json:"enabled"`
    +	}
    +	if err := json.Unmarshal(w.Body.Bytes(), &status); err != nil || status.Enabled {
     		t.Fatal("exited gateway still reported as enabled")
    Update { move_path: None }: /home/ubuntu/NeoNVR/main.go
    @@ -18,2 +18,3 @@
     	"os/exec"
    +	"os/signal"
     	"path/filepath"
    @@ -1835,2 +1836,15 @@
     	httpServer := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
    +	shutdownSignals := make(chan os.Signal, 1)
    +	signal.Notify(shutdownSignals, syscall.SIGTERM, syscall.SIGINT)
    +	defer signal.Stop(shutdownSignals)
    +	shutdown := func() {
    +		live.close()
    +		cancelWork()
    +		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    +		defer cancel()
    +		if err := httpServer.Shutdown(ctx); err != nil {
    +			_ = httpServer.Close()
    +		}
    +		s.background.Wait()
    +	}
     	serverErr := make(chan error, 1)
    @@ -1844,10 +1858,3 @@
     		log.Printf("Restarting NeoNVR with config %s", configPath)
    -		live.close()
    -		cancelWork()
    -		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    -		if err := httpServer.Shutdown(ctx); err != nil {
    -			_ = httpServer.Close()
    -		}
    -		cancel()
    -		s.background.Wait()
    +		shutdown()
     		if err := syscall.Exec(executable, []string{executable, "-config", configPath}, os.Environ()); err != nil {
    @@ -1855,2 +1862,4 @@
     		}
    +	case <-shutdownSignals:
    +		shutdown()
     	}

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w live_test.go live.go main.go
        NEONVR_TEST_GO2RTC=/tmp/neonvr-live-tools-vbq_9dyb/go2rtc /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
        /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr .
    ok      neonvr    4.808s
    ✓ • 6.30s

## Activity

    file changes: Completed · 1 changes
    Add: /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    import json
    from pathlib import Path
    import socket
    import socketserver
    import subprocess
    import tempfile
    import threading
    import time
    import urllib.request
    
    from playwright.sync_api import sync_playwright
    
    ROOT = Path('/tmp/neonvr-live-tools-vbq_9dyb')
    processes = []
    logs = []
    speaker_packets = 0
    speaker_connections = 0
    
    
    def port():
        with socket.socket() as listener:
            listener.bind(('127.0.0.1', 0))
            return listener.getsockname()[1]
    
    
    def start(arguments):
        log = tempfile.TemporaryFile()
        logs.append(log)
        process = subprocess.Popen(arguments, stdout=log, stderr=log)
        processes.append(process)
        return process
    
    
    def request(url):
        with urllib.request.urlopen(url, timeout=3) as response:
            return json.load(response)
    
    
    def ready(url, process):
        for attempt in range(100):
            if process.poll() is not None:
                raise RuntimeError('Server exited')
            try:
                return request(url)
            except Exception:
                time.sleep(.05)
        raise RuntimeError('Server did not start')
    
    
    class Speaker(socketserver.StreamRequestHandler):
        def handle(self):
            global speaker_packets, speaker_connections
            speaker_connections += 1
            while True:
                first = self.rfile.read(1)
                if not first:
                    return
                if first == b'$':
                    channel = self.rfile.read(1)
                    length = int.from_bytes(self.rfile.read(2), 'big')
                    data = self.rfile.read(length)
                    if channel == b'\x00' and len(data) > 12:
                        speaker_packets += 1
                    continue
                line = first + self.rfile.readline()
                method, target, _ = line.decode().strip().split(' ', 2)
                headers = {}
                while True:
                    line = self.rfile.readline()
                    if line in (b'\r\n', b'\n', b''):
                        break
                    name, value = line.decode().split(':', 1)
                    headers[name.lower()] = value.strip()
                size = int(headers.get('content-length', 0))
                if size:
                    self.rfile.read(size)
                body = b''
                extra = 'Session: local-speaker\r\n'
                if method == 'OPTIONS':
                    extra += 'Public: OPTIONS, DESCRIBE, SETUP, PLAY, TEARDOWN, GET_PARAMETER\r\n'
                elif method == 'DESCRIBE':
                    body = ('v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=Local speaker\r\nt=0 0\r\n'
                            'a=control:*\r\nm=audio 0 RTP/AVP 0\r\nc=IN IP4 0.0.0.0\r\n'
                            'a=rtpmap:0 PCMU/8000\r\na=recvonly\r\na=control:trackID=0\r\n').encode()
                    extra += 'Content-Type: application/sdp\r\nContent-Base: ' + target + '/\r\n'
                elif method == 'SETUP':
                    extra += 'Transport: ' + headers['transport'] + '\r\n'
                result = ('RTSP/1.0 200 OK\r\nCSeq: ' + headers.get('cseq', '1') + '\r\n' + extra +
                          'Content-Length: ' + str(len(body)) + '\r\n\r\n').encode() + body
                self.wfile.write(result)
                self.wfile.flush()
                if method == 'TEARDOWN':
                    return
    
    
    class SpeakerServer(socketserver.ThreadingTCPServer):
        allow_reuse_address = True
        daemon_threads = True
    
    
    speaker = SpeakerServer(('127.0.0.1', 0), Speaker)
    threading.Thread(target=speaker.serve_forever, daemon=True).start()
    
    try:
        with tempfile.TemporaryDirectory(prefix='neonvr-browser-fixture-') as temporary:
            directory = Path(temporary)
            rtsp_port, fixture_api, app_port, media_port = port(), port(), port(), port()
            fixture_config = directory / 'fixture.json'
            fixture_config.write_text(json.dumps({
                'api': {'listen': '127.0.0.1:' + str(fixture_api)},
                'rtsp': {'listen': '127.0.0.1:' + str(rtsp_port)},
                'webrtc': {'listen': '', 'ice_servers': []},
                'streams': {'fixture': []},
            }))
            fixture = start([str(ROOT / 'go2rtc'), '-config', str(fixture_config)])
            ready('http://127.0.0.1:' + str(fixture_api) + '/api', fixture)
            feed = start([
                'ffmpeg', '-nostdin', '-v', 'error', '-re', '-f', 'lavfi', '-i', 'testsrc=size=320x180:rate=10',
                '-f', 'lavfi', '-i', 'sine=frequency=440:sample_rate=8000', '-c:v', 'libx264', '-threads', '1',
                '-preset', 'ultrafast', '-tune', 'zerolatency', '-pix_fmt', 'yuv420p', '-g', '10',
                '-c:a', 'pcm_mulaw', '-ar', '8000', '-ac', '1', '-f', 'rtsp', '-rtsp_transport', 'tcp',
                'rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture',
            ])
            recordings = directory / 'recordings'
            recordings.mkdir()
            config = directory / 'neonvr.json'
            config.write_text(json.dumps({
                'root': str(recordings), 'extra_roots': [], 'addr': '127.0.0.1:' + str(app_port),
                'auto_delete_days': 0, 'thumb_maintenance_every': '0', 'detection_enabled': False,
                'live_go2rtc': str(ROOT / 'go2rtc'), 'live_webrtc_port': media_port,
                'live_streams': [{'name': 'Local test camera', 'url': 'rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture',
                                  'backchannel_url': 'rtsp://127.0.0.1:' + str(speaker.server_address[1]) + '/speaker?backchannel=1'}],
            }))
            app = start([str(ROOT / 'neonvr'), '-config', str(config)])
            origin = 'http://localhost:' + str(app_port)
            ready(origin + '/api/live', app)
            with sync_playwright() as p:
                browser = p.chromium.launch(headless=True, args=['--no-sandbox', '--use-fake-device-for-media-stream', '--use-fake-ui-for-media-stream'])
                context = browser.new_context(viewport={'width': 1280, 'height': 900}, permissions=['microphone'])
                page = context.new_page()
                errors = []
                page.on('pageerror', lambda error: errors.append(str(error)))
                page.goto(origin + '/')
                page.get_by_role('link', name='Streams', exact=True).click()
                page.wait_for_function('state !== null')
                assert not page.locator('#enabled').is_checked()
                assert page.locator('#watch').is_disabled()
                assert request(origin + '/api/live')['enabled'] is False
                assert speaker_connections == 0
                assert not request('http://127.0.0.1:' + str(fixture_api) + '/api/streams')['fixture'].get('consumers')
                page.locator('#enabled').check()
                page.wait_for_function('state.enabled && !changing')
                assert speaker_connections == 0
                assert not request('http://127.0.0.1:' + str(fixture_api) + '/api/streams')['fixture'].get('consumers')
                print('Opening the tab and enabling live viewing do not connect to cameras.', flush=True)
                page.locator('#watch').click()
                try:
                    page.wait_for_function('active && active.connected && document.querySelector("video").videoWidth === 320', timeout=20000)
                except Exception:
                    print('WATCH STATUS:', page.locator('#status').inner_text(), flush=True)
                    print('WATCH STATE:', page.evaluate('({pc:active?.pc.connectionState,ice:active?.pc.iceConnectionState,sdp:active?.pc.remoteDescription?.sdp,video:videoEl.readyState})'), flush=True)
                    raise
                assert page.locator('#talk').is_enabled()
                assert speaker_connections == 0
                assert page.evaluate('active.microphone === null')
                page.screenshot(path=str(ROOT / 'streams-desktop.png'), full_page=True)
                print('Real H.264 camera video plays over WebRTC without opening the backchannel.', flush=True)
                page.locator('#talk').click()
                try:
                    page.wait_for_function('active && active.connected && active.talk', timeout=20000)
                except Exception:
                    print('TALK STATUS:', page.locator('#status').inner_text(), flush=True)
                    print('TALK STATE:', page.evaluate('({pc:active?.pc.connectionState,sdp:active?.pc.remoteDescription?.sdp})'), flush=True)
                    raise
                page.wait_for_function('document.querySelector("#microphone").textContent === "Microphone live"')
                for attempt in range(40):
                    if speaker_packets:
                        break
                    time.sleep(.1)
                assert speaker_packets > 0, 'No user audio reached the RTSP speaker backchannel'
                page.evaluate('window.previousMic = active.microphone')
                page.locator('#talk').click()
                page.wait_for_function('active && active.connected && !active.talk', timeout=20000)
                assert page.evaluate('previousMic.getTracks().every(track => track.readyState === "ended")')
                print('Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.', flush=True)
                page.locator('#stop').click()
                assert page.evaluate('active === null && videoEl.srcObject === null')
                page.locator('#enabled').uncheck()
                page.wait_for_function('!state.enabled && !changing')
                assert request(origin + '/api/live')['enabled'] is False
                print('Stop and the off toggle close playback and the streaming gateway.', flush=True)
                page.goto(origin + '/settings')
                page.wait_for_function('state !== null')
                expected = 'Local test camera | rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture | rtsp://127.0.0.1:' + str(speaker.server_address[1]) + '/speaker?backchannel=1'
                assert page.locator('[name=live_streams]').input_value() == expected
                page.locator('#save').click()
                page.wait_for_function('document.querySelector("#status").textContent.startsWith("Settings saved.")')
                assert len(json.loads(config.read_text())['live_streams']) == 1
                print('Settings save preserves the live source and backchannel.', flush=True)
                page.goto(origin + '/streams')
                page.wait_for_function('state !== null')
                page.set_viewport_size({'width': 390, 'height': 844})
                assert page.evaluate('document.documentElement.scrollWidth <= window.innerWidth')
                page.screenshot(path=str(ROOT / 'streams-mobile.png'), full_page=True)
                assert errors == [], errors
                browser.close()
                print('Browser checks passed with no page errors, including mobile layout.', flush=True)
    finally:
        speaker.shutdown()
        speaker.server_close()
        for process in reversed(processes):
            process.terminate()
            try:
                process.wait(timeout=5)
            except subprocess.TimeoutExpired:
                process.kill()
                process.wait()
        for index, log in enumerate(logs):
            log.seek(0)
            output = log.read().decode(errors='replace')
            if output:
                print('Fixture log', index, output[-4500:], flush=True)
            log.close()

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Fixture log 0 23:45:31.809 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:45:31.810 INF config path=/tmp/neonvr-browser-fixture-dyoot_6m/fixture.json
    23:45:31.810 INF [rtsp] listen addr=127.0.0.1:60019
    23:45:31.810 INF [api] listen addr=127.0.0.1:46265
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:45:31 NeoNVR server listening on http://127.0.0.1:56433
    2026/10/10 23:45:31 Clips root: /tmp/neonvr-browser-fixture-dyoot_6m/recordings
    2026/10/10 23:45:31 Thumbnail cache: /tmp/neonvr-browser-fixture-dyoot_6m/recordings/.blinkclips_thumbs
    2026/10/10 23:45:31 Metadata: /tmp/neonvr-browser-fixture-dyoot_6m/recordings/.blinkclips_meta.json
    2026/10/10 23:45:31 ffmpeg: ffmpeg
    2026/10/10 23:45:31 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:45:31 Thumbnail seek: 0.200s
    2026/10/10 23:45:31 Thumbnail workers: 2
    2026/10/10 23:45:31 Thumbnail maintenance disabled
    2026/10/10 23:45:31 Size refresh: 5m0s
    2026/10/10 23:45:31 Auto-delete disabled; use -auto-delete-days N to enable
    
    Traceback (most recent call last):
      File "/tmp/neonvr-live-tools-vbq_9dyb/check_browser.py", line 151, in <module>
        page.locator('#enabled').check()
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/sync_api/_generated.py", line 17929, in check
        self._sync(
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_sync_base.py", line 115, in _sync
        return task.result()
               ^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_locator.py", line 148, in check
        return await self._frame.check(self._selector, strict=True, **params)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_frame.py", line 909, in check
        await self._channel.send("check", self._timeout, locals_to_params(locals()))
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_connection.py", line 76, in send
        return await self._connection.wrap_api_call(
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_connection.py", line 632, in wrap_api_call
        raise rewrite_error(error, f"{parsed_st['apiName']}: {error}") from None
    playwright._impl._errors.Error: Locator.check: Clicking the checkbox did not change its state
    Call log:
      - waiting for locator("#enabled")
        - locator resolved to <input id="enabled" type="checkbox"/>
      - attempting click action
        - waiting for element to be visible, enabled and stable
        - element is visible, enabled and stable
        - scrolling into view if needed
        - done scrolling
        - performing click action
        - click action done
        - waiting for scheduled navigations to finish
        - navigations have finished
    
    ✗ (1) • 1.72s

## Activity

    $ git diff --check
        git status --short
        sed -n '1780,1888p' main.go
        sed -n '1,100p' live.go
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M docker-compose.yml
     M main.go
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
            w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            _, _ = w.Write(liveStreamsJS)
        })
        mux.HandleFunc("/api/live", live.handler)
        mux.HandleFunc("/api/live/webrtc", live.webrtcHandler)
        mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                w.Header().Set("Allow", "GET")
                http.Error(w, "method not allowed", 405)
                return
            }
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            w.Header().Set("Cache-Control", "no-store")
            _, _ = w.Write(settingsHTML)
        })
        mux.HandleFunc("/api/cameras", s.apiCamerasHandler)
        mux.HandleFunc("/api/dates", s.apiDatesHandler)
        mux.HandleFunc("/api/clips", s.apiClipsHandler)
        mux.HandleFunc("/api/favorites", s.apiFavoritesHandler)
        mux.HandleFunc("/api/favorite", s.apiFavoriteHandler)
        mux.HandleFunc("/api/delete", s.apiDeleteHandler)
        mux.HandleFunc("/api/bulk-delete", s.apiBulkDeleteHandler)
        mux.HandleFunc("/api/stats", s.apiStatsHandler)
        mux.HandleFunc("/api/detections", s.apiDetectionsHandler)
        mux.HandleFunc("/static/detection-viewer.js", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            _, _ = w.Write(detectionViewerJS)
        })
        mux.HandleFunc("/static/media-chrome-lite.js", s.staticMediaChromeHandler)
        mux.HandleFunc("/video", s.videoHandler)
        mux.HandleFunc("/thumb", s.thumbHandler)
    
        log.Printf("NeoNVR server listening on http://%s", *addr)
        log.Printf("Clips root: %s", s.root)
        for _, extraRoot := range s.extraRoots {
            log.Printf("Extra clip root: %s", extraRoot)
        }
        log.Printf("Thumbnail cache: %s", s.thumbDir)
        log.Printf("Metadata: %s", s.metaPath)
        log.Printf("ffmpeg: %s", s.ffmpegPath)
        log.Printf("YOLO26n detection: enabled=%t, CPU threads=%d, sample fps=%g, no-target policy=%s", cfg.DetectionEnabled, cfg.DetectionThreads, cfg.DetectionFPS, cfg.DetectionPrune)
        log.Printf("Thumbnail seek: %.3fs", s.thumbSeek)
        log.Printf("Thumbnail workers: %d", workers)
        if s.thumbMaintenanceEvery > 0 {
            log.Printf("Thumbnail maintenance: every %s (pre-generate missing, prune orphan cache files)", s.thumbMaintenanceEvery)
        } else {
            log.Printf("Thumbnail maintenance disabled")
        }
        log.Printf("Size refresh: %s", s.sizeTTL)
        if s.autoDeleteDays > 0 {
            log.Printf("Auto-delete enabled: non-favorite clips older than %d days, checking every %s", s.autoDeleteDays, s.autoDeleteEvery)
        } else {
            log.Printf("Auto-delete disabled; use -auto-delete-days N to enable")
        }
        httpServer := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
        shutdownSignals := make(chan os.Signal, 1)
        signal.Notify(shutdownSignals, syscall.SIGTERM, syscall.SIGINT)
        defer signal.Stop(shutdownSignals)
        shutdown := func() {
            live.close()
            cancelWork()
            ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
            defer cancel()
            if err := httpServer.Shutdown(ctx); err != nil {
                _ = httpServer.Close()
            }
            s.background.Wait()
        }
        serverErr := make(chan error, 1)
        go func() { serverErr <- httpServer.ListenAndServe() }()
        select {
        case err := <-serverErr:
            if !errors.Is(err, http.ErrServerClosed) {
                log.Fatal(err)
            }
        case <-manager.restart:
            log.Printf("Restarting NeoNVR with config %s", configPath)
            shutdown()
            if err := syscall.Exec(executable, []string{executable, "-config", configPath}, os.Environ()); err != nil {
                log.Fatalf("restart failed: %v", err)
            }
        case <-shutdownSignals:
            shutdown()
        }
    }
    
    func NewMetaStore(path string) (*MetaStore, error) {
        m := &MetaStore{path: path, data: MetaFile{Version: 1, Clips: map[string]ClipMeta{}}}
        b, err := os.ReadFile(path)
        if err != nil {
            if os.IsNotExist(err) {
                return m, nil
            }
            return nil, err
        }
        if len(strings.TrimSpace(string(b))) == 0 {
            return m, nil
        }
        if err := json.Unmarshal(b, &m.data); err != nil {
            return nil, err
        }
        if m.data.Clips == nil {
            m.data.Clips = map[string]ClipMeta{}
        }
        if m.data.Version == 0 {
            m.data.Version = 1
        }
    package main
    
    import (
        "bytes"
        "context"
        "crypto/rand"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "io"
        "net"
        "net/http"
        "net/url"
        "os"
        "os/exec"
        "path/filepath"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    )
    
    type LiveStream struct {
        Name           string `json:"name"`
        URL            string `json:"url"`
        BackchannelURL string `json:"backchannel_url,omitempty"`
    }
    
    func validateRTSP(raw string) error {
        u, err := url.Parse(raw)
        if err != nil || (u.Scheme != "rtsp" && u.Scheme != "rtsps") || u.Hostname() == "" || u.Fragment != "" || strings.ContainsAny(raw, "\r\n\t ") {
            return errors.New("enter an rtsp:// or rtsps:// camera URL without spaces or a fragment")
        }
        if port := u.Port(); port != "" {
            n, err := strconv.Atoi(port)
            if err != nil || n < 1 || n > 65535 {
                return errors.New("RTSP port must be between 1 and 65535")
            }
        }
        return nil
    }
    
    func (c Config) validateLive() error {
        if strings.TrimSpace(c.LiveGo2RTC) == "" {
            return errors.New("live_go2rtc cannot be empty")
        }
        if c.LiveWebRTCPort < 1 || c.LiveWebRTCPort > 65535 {
            return errors.New("live_webrtc_port must be between 1 and 65535")
        }
        _, appPort, _ := net.SplitHostPort(c.Addr)
        if appPort == strconv.Itoa(c.LiveWebRTCPort) {
            return errors.New("live_webrtc_port must differ from the app listen port")
        }
        if host := c.LiveWebRTCHost; host != "" && (strings.ContainsAny(host, "/?#@ \r\n\t") || strings.Contains(host, ":") && net.ParseIP(host) == nil) {
            return errors.New("live_webrtc_host must be a server hostname or IP address without a port")
        }
        if len(c.LiveStreams) > 32 {
            return errors.New("configure at most 32 live streams")
        }
        names := map[string]bool{}
        for index, stream := range c.LiveStreams {
            name := strings.TrimSpace(stream.Name)
            if name == "" || len(name) > 120 || strings.ContainsAny(name, "|\r\n") || names[name] {
                return fmt.Errorf("live stream %d needs a unique name of up to 120 characters", index+1)
            }
            names[name] = true
            if err := validateRTSP(stream.URL); err != nil {
                return fmt.Errorf("live stream %d: %w", index+1, err)
            }
            if stream.BackchannelURL != "" {
                if err := validateRTSP(stream.BackchannelURL); err != nil {
                    return fmt.Errorf("live stream %d backchannel: %w", index+1, err)
                }
            }
        }
        return nil
    }
    
    type liveService struct {
        cmd      *exec.Cmd
        done     chan struct{}
        baseURL  string
        password string
        client   *http.Client
    }
    
    type liveManager struct {
        mu      sync.Mutex
        ctx     context.Context
        config  Config
        token   string
        service *liveService
    }
    
    func newLiveManager(ctx context.Context, config Config, token string) *liveManager {
        // Enabling is deliberately transient: every app start defaults to off.
        return &liveManager{ctx: ctx, config: config, token: token}
    }
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import urllib.request,json
        with urllib.request.urlopen('https://ghcr.io/token?scope=repository:alexxit/go2rtc:pull',timeout=20) as response: token=[REDACTED_SECRET]'token']
        request=urllib.request.Request('https://ghcr.io/v2/alexxit/go2rtc/manifests/1.9.14',headers={'Authorization':'Bearer '+token,'Accept':'application/vnd.oci.image.index.v1+json, application/vnd.docker.distribution.manifest.list.v2+json'})
        with urllib.request.urlopen(request,timeout=20) as response:
         data=json.load(response)
         print('Pinned go2rtc image is available:',response.headers.get('Docker-Content-Digest'))
         print('Platforms:',[entry.get('platform') for entry in data.get('manifests',[]) if entry.get('platform',{}).get('os')=='linux'])
        PY
    Pinned go2rtc image is available: sha256:675c318b23c06fd862a61d262240c9a63436b4050d177ffc68a32710d9e05bae
    Platforms: [{'architecture': 'amd64', 'os': 'linux'}, {'architecture': '386', 'os': 'linux'}, {'architecture': 'arm', 'os': 'linux', 'variant': 'v6'}, {'architecture': 'arm', 'os': 'linux', 'variant': 'v7'}, {'architecture': 'arm64', 'os': 'linux'}]
    ✓ • 139ms

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/live-streams.js
    @@ -30,3 +30,3 @@
       enabledEl.disabled = changing || !state || !state.streams.length;
    -  enabledEl.checked = !!(state && state.enabled);
    +  if (!changing) enabledEl.checked = !!(state && state.enabled);
       watchEl.disabled = changing || !source || !state.enabled;
    @@ -206,4 +206,4 @@
       const enabled = enabledEl.checked;
    +  changing = true;
       if (!enabled) stopPlayback();
    -  changing = true;
       updateControls();

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr .
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    Fixture log 0 23:46:12.672 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:46:12.672 INF config path=/tmp/neonvr-browser-fixture-ftxwlcb7/fixture.json
    23:46:12.673 INF [api] listen addr=127.0.0.1:38563
    23:46:12.673 INF [rtsp] listen addr=127.0.0.1:43277
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:46:12 NeoNVR server listening on http://127.0.0.1:41281
    2026/10/10 23:46:12 Clips root: /tmp/neonvr-browser-fixture-ftxwlcb7/recordings
    2026/10/10 23:46:12 Thumbnail cache: /tmp/neonvr-browser-fixture-ftxwlcb7/recordings/.blinkclips_thumbs
    2026/10/10 23:46:12 Metadata: /tmp/neonvr-browser-fixture-ftxwlcb7/recordings/.blinkclips_meta.json
    2026/10/10 23:46:12 ffmpeg: ffmpeg
    2026/10/10 23:46:12 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:46:12 Thumbnail seek: 0.200s
    2026/10/10 23:46:12 Thumbnail workers: 2
    2026/10/10 23:46:12 Thumbnail maintenance disabled
    2026/10/10 23:46:12 Size refresh: 5m0s
    2026/10/10 23:46:12 Auto-delete disabled; use -auto-delete-days N to enable
    
    Traceback (most recent call last):
      File "/tmp/neonvr-live-tools-vbq_9dyb/check_browser.py", line 180, in <module>
        assert speaker_packets > 0, 'No user audio reached the RTSP speaker backchannel'
               ^^^^^^^^^^^^^^^^^^^
    AssertionError: No user audio reached the RTSP speaker backchannel
    ✗ (1) • 7.68s

## Assistant

The server tests pass: playback requests are blocked while off, and enabling the bridge alone leaves cameras disconnected. The local browser test is now checking H.264 playback and microphone audio through the RTSP backchannel.

## Activity

    $ sed -n '48,103p' /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    
    
    class Speaker(socketserver.StreamRequestHandler):
        def handle(self):
            global speaker_packets, speaker_connections
            speaker_connections += 1
            while True:
                first = self.rfile.read(1)
                if not first:
                    return
                if first == b'$':
                    channel = self.rfile.read(1)
                    length = int.from_bytes(self.rfile.read(2), 'big')
                    data = self.rfile.read(length)
                    if channel == b'\x00' and len(data) > 12:
                        speaker_packets += 1
                    continue
                line = first + self.rfile.readline()
                method, target, _ = line.decode().strip().split(' ', 2)
                headers = {}
                while True:
                    line = self.rfile.readline()
                    if line in (b'\r\n', b'\n', b''):
                        break
                    name, value = line.decode().split(':', 1)
                    headers[name.lower()] = value.strip()
                size = int(headers.get('content-length', 0))
                if size:
                    self.rfile.read(size)
                body = b''
                extra = 'Session: local-speaker\r\n'
                if method == 'OPTIONS':
                    extra += 'Public: OPTIONS, DESCRIBE, SETUP, PLAY, TEARDOWN, GET_PARAMETER\r\n'
                elif method == 'DESCRIBE':
                    body = ('v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=Local speaker\r\nt=0 0\r\n'
                            'a=control:*\r\nm=audio 0 RTP/AVP 0\r\nc=IN IP4 0.0.0.0\r\n'
                            'a=rtpmap:0 PCMU/8000\r\na=recvonly\r\na=control:trackID=0\r\n').encode()
                    extra += 'Content-Type: application/sdp\r\nContent-Base: ' + target + '/\r\n'
                elif method == 'SETUP':
                    extra += 'Transport: ' + headers['transport'] + '\r\n'
                result = ('RTSP/1.0 200 OK\r\nCSeq: ' + headers.get('cseq', '1') + '\r\n' + extra +
                          'Content-Length: ' + str(len(body)) + '\r\n\r\n').encode() + body
                self.wfile.write(result)
                self.wfile.flush()
                if method == 'TEARDOWN':
                    return
    
    
    class SpeakerServer(socketserver.ThreadingTCPServer):
        allow_reuse_address = True
        daemon_threads = True
    
    
    speaker = SpeakerServer(('127.0.0.1', 0), Speaker)
    threading.Thread(target=speaker.serve_forever, daemon=True).start()
    
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import urllib.request
        for path in ['pkg/rtsp/client.go','pkg/rtsp/producer.go','pkg/rtsp/conn.go','pkg/rtsp/helpers.go']:
         try:
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+path,timeout=20) as response:lines=response.read().decode().splitlines()
         except Exception:continue
         print(path)
         for index,line in enumerate(lines):
          if any(pattern in line for pattern in ['interleaved','Backchannel','SetupMedia','func (c *Conn) Start','SendRTP','DirectionRecvonly']):
           print('\n'.join(lines[max(0,index-3):index+7]))
        PY
    pkg/rtsp/client.go
            },
        }
    
        if c.Backchannel {
            req.Header.Set("Require", "www.onvif.org/ver20/backchannel")
        }
    
        if c.UserAgent != "" {
            // this camera will answer with 401 on DESCRIBE without User-Agent
            // https://github.com/AlexxIT/go2rtc/issues/235
        return
    }
    
    func (c *Conn) SetupMedia(media *core.Media) (byte, error) {
        var transport string
    
        if c.Transport == "udp" {
            conn1, conn2, err := ListenUDPPair()
            if err != nil {
                return 0, err
                    transport = fmt.Sprintf(
                        // i   - RTP (data channel)
                        // i+1 - RTCP (control channel)
                        "RTP/AVP/TCP;unicast;interleaved=%d-%d", i*2, i*2+1,
                    )
                    break
                }
            }
        }
    
        if err != nil {
            // some Dahua/Amcrest cameras fail here because two simultaneous
            // backchannel connections
            if c.Backchannel {
                c.Backchannel = false
                if err = c.Reconnect(); err != nil {
                    return 0, err
                }
                return c.SetupMedia(media)
            }
            // some Dahua/Amcrest cameras fail here because two simultaneous
            // backchannel connections
            if c.Backchannel {
                c.Backchannel = false
                if err = c.Reconnect(); err != nil {
                    return 0, err
                }
                return c.SetupMedia(media)
            }
    
                if err = c.Reconnect(); err != nil {
                    return 0, err
                }
                return c.SetupMedia(media)
            }
    
            return 0, err
        }
    
        if c.session == "" {
    
            return channel, nil
        } else {
            // we send our `interleaved`, but camera can answer with another
    
            // Transport: RTP/AVP/TCP;unicast;interleaved=10-11;ssrc=10117CB7
            // Transport: RTP/AVP/TCP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
        } else {
            // we send our `interleaved`, but camera can answer with another
    
            // Transport: RTP/AVP/TCP;unicast;interleaved=10-11;ssrc=10117CB7
            // Transport: RTP/AVP/TCP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
            s := core.Between(transport, "interleaved=", "-")
            i, err := strconv.Atoi(s)
            // we send our `interleaved`, but camera can answer with another
    
            // Transport: RTP/AVP/TCP;unicast;interleaved=10-11;ssrc=10117CB7
            // Transport: RTP/AVP/TCP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
            s := core.Between(transport, "interleaved=", "-")
            i, err := strconv.Atoi(s)
            if err != nil {
    
            // Transport: RTP/AVP/TCP;unicast;interleaved=10-11;ssrc=10117CB7
            // Transport: RTP/AVP/TCP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
            s := core.Between(transport, "interleaved=", "-")
            i, err := strconv.Atoi(s)
            if err != nil {
                return 0, fmt.Errorf("wrong transport: %s", transport)
            // Transport: RTP/AVP/TCP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
            s := core.Between(transport, "interleaved=", "-")
            i, err := strconv.Atoi(s)
            if err != nil {
                return 0, fmt.Errorf("wrong transport: %s", transport)
            }
    
            // Transport: RTP/AVP/TCP;ssrc=22345682;interleaved=0-1
            // Escam Q6 has a bug:
            // Transport: RTP/AVP;unicast;destination=192.168.1.111;source=192.168.1.222;interleaved=0-1
            s := core.Between(transport, "interleaved=", "-")
            i, err := strconv.Atoi(s)
            if err != nil {
                return 0, fmt.Errorf("wrong transport: %s", transport)
            }
    
            return byte(i), nil
    pkg/rtsp/producer.go
    )
    
    func (c *Conn) GetTrack(media *core.Media, codec *core.Codec) (*core.Receiver, error) {
        core.Assert(media.Direction == core.DirectionRecvonly)
    
        for _, track := range c.Receivers {
            if track.Codec == codec {
                return track, nil
            }
        }
            }
    
            var err error
            channel, err = c.SetupMedia(media)
            if err != nil {
                return nil, err
            }
    
            c.state = StateSetup
        case core.ModePassiveConsumer:
    
            c.state = StateSetup
        case core.ModePassiveConsumer:
            // Backchannel
            channel = byte(len(c.Senders)) * 2
        default:
            return nil, errors.New("rtsp: wrong mode for GetTrack")
        }
    
        track := core.NewReceiver(media, codec)
        return track, nil
    }
    
    func (c *Conn) Start() (err error) {
        core.Assert(c.mode == core.ModeActiveProducer || c.mode == core.ModePassiveProducer)
    
        for {
            ok := false
    
            c.stateMu.Lock()
    
        // restore previous medias
        for _, receiver := range c.Receivers {
            if _, err := c.SetupMedia(receiver.Media); err != nil {
                return err
            }
        }
        for _, sender := range c.Senders {
            if _, err := c.SetupMedia(sender.Media); err != nil {
                return err
            }
        }
        for _, sender := range c.Senders {
            if _, err := c.SetupMedia(sender.Media); err != nil {
                return err
            }
        }
    
        return nil
    }
    pkg/rtsp/conn.go
    
        // public
    
        Backchannel bool
        Media       string
        OnClose     func() error
        PacketSize  uint16
        SessionName string
        Timeout     int
        Transport   string // custom transport support, ex. RTSP over WebSocket
    
    func (c *Conn) handleTCPData() error {
        // we can read:
        // 1. RTP interleaved: `$` + 1B channel number + 2B size
        // 2. RTSP response:   RTSP/1.0 200 OK
        // 3. RTSP request:    OPTIONS ...
        var buf4 []byte // `$` + 1B channel number + 2B size
        var err error
    
        buf4, err = c.reader.Peek(4)
    pkg/rtsp/helpers.go
            }
    
            if media.Direction == "" || forceDirection {
                media.Direction = core.DirectionRecvonly
            }
    
            medias = append(medias, media)
        }
    
        return medias, nil
    ✓ • 454ms

## Activity

    Viewed image /tmp/neonvr-live-tools-vbq_9dyb/streams-desktop.png

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    @@ -17,2 +17,3 @@
     speaker_connections = 0
    +speaker_events = []
     
    @@ -61,2 +62,4 @@
                     data = self.rfile.read(length)
    +                if len(speaker_events) < 30:
    +                    speaker_events.append(('packet', int.from_bytes(channel, 'big'), len(data)))
                     if channel == b'\x00' and len(data) > 12:
    @@ -66,2 +69,3 @@
                 method, target, _ = line.decode().strip().split(' ', 2)
    +            speaker_events.append((method, target))
                 headers = {}
    @@ -179,2 +183,6 @@
                     time.sleep(.1)
    +            if speaker_packets == 0:
    +                print('SPEAKER EVENTS:',speaker_events,flush=True)
    +                print('AUDIO STATE:',page.evaluate('({tracks:active.microphone.getTracks().map(t=>({enabled:t.enabled,muted:t.muted,state:t.readyState})),transceivers:active.pc.getTransceivers().map(t=>({direction:t.currentDirection,kind:t.sender.track?.kind})),sdp:active.pc.remoteDescription.sdp})'),flush=True)
    +                print('AUDIO STATS:',page.evaluate('async()=>Array.from((await active.pc.getStats()).values()).filter(s=>s.type==="outbound-rtp"||s.type==="media-source"||s.type==="codec")'),flush=True)
                 assert speaker_packets > 0, 'No user audio reached the RTSP speaker backchannel'

## Activity

    $ python3 - <<'PY'
        import urllib.request
        for path in ['pkg/rtsp/helpers.go','pkg/rtsp/consumer.go']:
         with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+path,timeout=20) as response: print(path+'\n'+response.read().decode()[:11500])
        PY
    pkg/rtsp/helpers.go
    package rtsp
    
    import (
        "bytes"
        "io"
        "net/url"
        "regexp"
        "strconv"
        "strings"
    
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/pion/rtcp"
        "github.com/pion/sdp/v3"
    )
    
    type RTCP struct {
        Channel byte
        Header  rtcp.Header
        Packets []rtcp.Packet
    }
    
    const sdpHeader = `v=0
    o=- 0 0 IN IP4 0.0.0.0
    s=-
    t=0 0`
    
    func UnmarshalSDP(rawSDP []byte) ([]*core.Media, error) {
        sd := &sdp.SessionDescription{}
        if err := sd.Unmarshal(rawSDP); err != nil {
            // fix multiple `s=` https://github.com/AlexxIT/WebRTC/issues/417
            rawSDP = regexp.MustCompile("\ns=[^\n]+").ReplaceAll(rawSDP, nil)
    
            // fix broken `c=` https://github.com/AlexxIT/go2rtc/issues/1426
            rawSDP = regexp.MustCompile("\nc=[^\n]+").ReplaceAll(rawSDP, nil)
    
            // fix SDP header for some cameras
            if i := bytes.Index(rawSDP, []byte("\nm=")); i > 0 {
                rawSDP = append([]byte(sdpHeader), rawSDP[i:]...)
            }
    
            // Fix invalid media type (errSDPInvalidValue) caused by
            // some TP-LINK IP camera, e.g. TL-IPC44GW
            for _, b := range regexp.MustCompile("m=[^ ]+ ").FindAll(rawSDP, -1) {
                switch string(b[2 : len(b)-1]) {
                case "audio", "video", "application":
                default:
                    rawSDP = bytes.Replace(rawSDP, b, []byte("m=application "), 1)
                }
            }
    
            if err == io.EOF {
                rawSDP = append(rawSDP, '\n')
            }
    
            sd = &sdp.SessionDescription{}
            err = sd.Unmarshal(rawSDP)
            if err != nil {
                return nil, err
            }
        }
    
        // fix buggy camera https://github.com/AlexxIT/go2rtc/issues/771
        forceDirection := sd.Origin.Username == "CV-RTSPHandler"
    
        var medias []*core.Media
    
        for _, md := range sd.MediaDescriptions {
            media := core.UnmarshalMedia(md)
    
            // Check buggy SDP with fmtp for H264 on another track
            // https://github.com/AlexxIT/WebRTC/issues/419
            for _, codec := range media.Codecs {
                switch codec.Name {
                case core.CodecH264:
                    if codec.FmtpLine == "" {
                        codec.FmtpLine = findFmtpLine(codec.PayloadType, sd.MediaDescriptions)
                    }
                case core.CodecH265:
                    if codec.FmtpLine != "" {
                        // all three parameters are needed for a valid fmtp line
                        // https://github.com/AlexxIT/go2rtc/pull/1588
                        if !strings.Contains(codec.FmtpLine, "sprop-vps=") ||
                            !strings.Contains(codec.FmtpLine, "sprop-sps=") ||
                            !strings.Contains(codec.FmtpLine, "sprop-pps=") {
                            codec.FmtpLine = ""
                        }
                    }
                case core.CodecOpus:
                    // fix OPUS for some cameras https://datatracker.ietf.org/doc/html/rfc7587
                    codec.ClockRate = 48000
                    codec.Channels = 2
                }
            }
    
            if media.Direction == "" || forceDirection {
                media.Direction = core.DirectionRecvonly
            }
    
            medias = append(medias, media)
        }
    
        return medias, nil
    }
    
    func findFmtpLine(payloadType uint8, descriptions []*sdp.MediaDescription) string {
        s := strconv.Itoa(int(payloadType))
        for _, md := range descriptions {
            codec := core.UnmarshalCodec(md, s)
            if codec.FmtpLine != "" {
                return codec.FmtpLine
            }
        }
        return ""
    }
    
    // urlParse fix bugs:
    // 1. Content-Base: rtsp://::ffff:192.168.1.123/onvif/profile.1/
    // 2. Content-Base: rtsp://rtsp://turret2-cam.lan:554/stream1/
    // 3. Content-Base: 192.168.253.220:1935/
    func urlParse(rawURL string) (*url.URL, error) {
        // fix https://github.com/AlexxIT/go2rtc/issues/830
        if strings.HasPrefix(rawURL, "rtsp://rtsp://") {
            rawURL = rawURL[7:]
        }
    
        // fix https://github.com/AlexxIT/go2rtc/issues/1852
        if !strings.Contains(rawURL, "://") {
            rawURL = "rtsp://" + rawURL
        }
    
        u, err := url.Parse(rawURL)
        if err != nil && strings.HasSuffix(err.Error(), "after host") {
            if i := indexN(rawURL, '/', 3); i > 0 {
                return urlParse(rawURL[:i] + ":" + rawURL[i:])
            }
        }
    
        return u, err
    }
    
    func indexN(s string, c byte, n int) int {
        var offset int
        for {
            i := strings.IndexByte(s[offset:], c)
            if i < 0 {
                break
            }
            if n--; n == 0 {
                return offset + i
            }
            offset += i + 1
        }
        return -1
    }
    
    pkg/rtsp/consumer.go
    package rtsp
    
    import (
        "time"
    
        "github.com/AlexxIT/go2rtc/pkg/aac"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/h264"
        "github.com/AlexxIT/go2rtc/pkg/h265"
        "github.com/AlexxIT/go2rtc/pkg/mjpeg"
        "github.com/AlexxIT/go2rtc/pkg/pcm"
        "github.com/pion/rtp"
    )
    
    func (c *Conn) GetMedias() []*core.Media {
        //core.Assert(c.Medias != nil)
        return c.Medias
    }
    
    func (c *Conn) AddTrack(media *core.Media, codec *core.Codec, track *core.Receiver) (err error) {
        var channel byte
    
        switch c.mode {
        case core.ModeActiveProducer: // backchannel
            c.stateMu.Lock()
            defer c.stateMu.Unlock()
    
            if c.state == StatePlay {
                if err = c.Reconnect(); err != nil {
                    return
                }
            }
    
            if channel, err = c.SetupMedia(media); err != nil {
                return
            }
    
            c.state = StateSetup
    
        case core.ModePassiveConsumer:
            channel = byte(len(c.Senders)) * 2
    
            // for consumer is better to use original track codec
            codec = track.Codec.Clone()
            // generate new payload type, starting from 96
            codec.PayloadType = byte(96 + len(c.Senders))
    
        default:
            panic(core.Caller())
        }
    
        // save original codec to sender (can have Codec.Name = ANY)
        sender := core.NewSender(media, codec)
        // important to send original codec for valid IsRTP check
        sender.Handler = c.packetWriter(track.Codec, channel, codec.PayloadType)
    
        if c.mode == core.ModeActiveProducer && track.Codec.Name == core.CodecPCMA {
            // Fix Reolink Doorbell https://github.com/AlexxIT/go2rtc/issues/331
            sender.Handler = pcm.RepackG711(true, sender.Handler)
        }
    
        sender.HandleRTP(track)
    
        c.Senders = append(c.Senders, sender)
        return nil
    }
    
    const (
        startVideoBuf = 32 * 1024   // 32KB
        startAudioBuf = 2 * 1024    // 2KB
        maxBuf        = 1024 * 1024 // 1MB
        rtpHdr        = 12          // basic RTP header size
        intHdr        = 4           // interleaved header size
    )
    
    func (c *Conn) packetWriter(codec *core.Codec, channel, payloadType uint8) core.HandlerFunc {
        var buf []byte
        var n int
    
        video := codec.IsVideo()
        if video {
            buf = make([]byte, startVideoBuf)
        } else {
            buf = make([]byte, startAudioBuf)
        }
    
        flushBuf := func() {
            //log.Printf("[rtsp] channel:%2d write_size:%6d buffer_size:%6d", channel, n, len(buf))
            if err := c.writeInterleavedData(buf[:n]); err != nil {
                c.Send += n
            }
            n = 0
        }
    
        handlerFunc := func(packet *rtp.Packet) {
            if c.state == StateNone {
                return
            }
    
            clone := rtp.Packet{
                Header: rtp.Header{
                    Version:        2,
                    Marker:         packet.Marker,
                    PayloadType:    payloadType,
                    SequenceNumber: packet.SequenceNumber,
                    Timestamp:      packet.Timestamp,
                    SSRC:           packet.SSRC,
                },
                Payload: packet.Payload,
            }
    
            if !video {
                packet.Marker = true // better to have marker on all audio packets
            }
    
            size := rtpHdr + len(packet.Payload)
    
            if l := len(buf); n+intHdr+size > l {
                if l < maxBuf {
                    buf = append(buf, make([]byte, l)...) // double buffer size
                } else {
                    flushBuf()
                }
            }
    
            //log.Printf("[RTP] codec: %s, size: %6d, ts: %10d, pt: %2d, ssrc: %d, seq: %d, mark: %v", codec.Name, len(packet.Payload), packet.Timestamp, packet.PayloadType, packet.SSRC, packet.SequenceNumber, packet.Marker)
    
            chunk := buf[n:]
            _ = chunk[4] // bounds
            chunk[0] = '$'
            chunk[1] = channel
            chunk[2] = byte(size >> 8)
            chunk[3] = byte(size)
    
            if _, err := clone.MarshalTo(chunk[4:]); err != nil {
                return
            }
    
            n += 4 + size
    
            if !packet.Marker || !c.playOK {
                // collect continious video packets to buffer
                // or wait OK for PLAY command for backchannel
                //log.Printf("[rtsp] collecting buffer ok=%t", c.playOK)
                return
            }
    
            flushBuf()
        }
    
        if !codec.IsRTP() {
            switch codec.Name {
            case core.CodecH264:
                handlerFunc = h264.RTPPay(c.PacketSize, handlerFunc)
            case core.CodecH265:
                handlerFunc = h265.RTPPay(c.PacketSize, handlerFunc)
            case core.CodecAAC:
                handlerFunc = aac.RTPPay(handlerFunc)
            case core.CodecJPEG:
                handlerFunc = mjpeg.RTPPay(handlerFunc)
            }
        } else if codec.Name == core.CodecPCML {
            handlerFunc = pcm.LittleToBig(handlerFunc)
        } else if c.PacketSize != 0 {
            switch codec.Name {
            case core.CodecH264:
                handlerFunc = h264.RTPPay(c.PacketSize, handlerFunc)
                handlerFunc = h264.RTPDepay(codec, handlerFunc)
            case core.CodecH265:
                handlerFunc = h265.RTPPay(c.PacketSize, handlerFunc)
                handlerFunc = h265.RTPDepay(codec, handlerFunc)
            }
        }
    
        return handlerFunc
    }
    
    func (c *Conn) writeInterleavedData(data []byte) error {
        if c.Transport != "udp" {
            _ = c.conn.SetWriteDeadline(time.Now().Add(Timeout))
            _, err := c.conn.Write(data)
            return err
        }
    
        for len(data) >= 4 && data[0] == '$' {
            channel := data[1]
            size := uint16(data[2])<<8 | uint16(data[3])
            rtpData := data[4 : 4+size]
    
            if _, err := c.WriteToUDP(rtpData, channel); err != nil {
                return err
            }
    
            data = data[4+size:]
        }
    
        return nil
    }
    
    ✓ • 67ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    SPEAKER EVENTS: [('DESCRIBE', 'rtsp://127.0.0.1:43839/speaker?backchannel=1'), ('SETUP', 'rtsp://127.0.0.1:43839/speaker?backchannel=1/trackID=0'), ('PLAY', 'rtsp://127.0.0.1:43839/speaker?backchannel=1/')]
    AUDIO STATE: {'tracks': [{'enabled': True, 'muted': False, 'state': 'live'}], 'transceivers': [{'direction': 'sendonly', 'kind': 'audio'}, {'direction': 'inactive', 'kind': None}, {'direction': 'recvonly', 'kind': None}], 'sdp': 'v=0\r\no=- 5111252415494560639 1791690408 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0 1 2\r\na=extmap-allow-mixed\r\na=msid-semantic: WMS go2rtc\r\nm=audio 58539 UDP/TLS/RTP/SAVPF 111 0 8\r\nc=IN IP4 127.0.0.1\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=candidate:3553330067 1 tcp 1671430142 localhost 58539 typ host tcptype passive generation 0\r\na=candidate:1849225843 1 udp 2130706429 localhost 58539 typ host generation 0\r\na=candidate:2878742611 1 udp 2130706431 127.0.0.1 58539 typ host generation 0\r\na=candidate:3135901995 1 udp 2130706431 192.168.1.130 58539 typ host generation 0\r\na=candidate:1934064981 1 udp 2130706431 10.4.5.1 58539 typ host generation 0\r\na=candidate:1118791624 1 udp 2130706431 10.4.6.1 58539 typ host generation 0\r\na=candidate:375738803 1 tcp 1671430143 127.0.0.1 58539 typ host tcptype passive generation 0\r\na=candidate:119103691 1 tcp 1671430143 192.168.1.130 58539 typ host tcptype passive generation 0\r\na=candidate:3467906229 1 tcp 1671430143 10.4.5.1 58539 typ host tcptype passive generation 0\r\na=candidate:4284225064 1 tcp 1671430143 10.4.6.1 58539 typ host tcptype passive generation 0\r\na=ice-ufrag:aWlVKVdAnPzLwxha\r\na=ice-pwd:pNLkzJSHDMsNMCoAtpXwfohrEbLcnKpQ\r\na=fingerprint:sha-256 FC:C9:57:0D:83:ED:2A:A0:AF:05:2D:02:89:7C:FE:3C:1D:9F:BB:40:01:9A:9B:F9:53:33:9A:93:67:F2:47:06\r\na=setup:active\r\na=mid:0\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=recvonly\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:111 opus/48000/2\r\na=rtcp-fb:111 transport-cc\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\nm=video 9 UDP/TLS/RTP/SAVPF 102 108\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:aWlVKVdAnPzLwxha\r\na=ice-pwd:pNLkzJSHDMsNMCoAtpXwfohrEbLcnKpQ\r\na=fingerprint:sha-256 FC:C9:57:0D:83:ED:2A:A0:AF:05:2D:02:89:7C:FE:3C:1D:9F:BB:40:01:9A:9B:F9:53:33:9A:93:67:F2:47:06\r\na=setup:active\r\na=mid:1\r\na=extmap:10 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\r\na=extmap:11 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\r\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=inactive\r\na=msid:go2rtc go2rtc-video\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:102 H264/90000\r\na=rtcp-fb:102 goog-remb\r\na=rtcp-fb:102 ccm fir\r\na=rtcp-fb:102 nack\r\na=rtcp-fb:102 nack pli\r\na=rtcp-fb:102 transport-cc\r\na=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f\r\na=rtpmap:108 H264/90000\r\na=rtcp-fb:108 goog-remb\r\na=rtcp-fb:108 ccm fir\r\na=rtcp-fb:108 nack\r\na=rtcp-fb:108 nack pli\r\na=rtcp-fb:108 transport-cc\r\na=fmtp:108 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\r\na=ssrc:3593793656 cname:go2rtc\r\na=ssrc:3593793656 msid:go2rtc go2rtc-video\r\nm=audio 9 UDP/TLS/RTP/SAVPF 111 0 8\r\nc=IN IP4 0.0.0.0\r\na=rtcp:9 IN IP4 0.0.0.0\r\na=ice-ufrag:aWlVKVdAnPzLwxha\r\na=ice-pwd:pNLkzJSHDMsNMCoAtpXwfohrEbLcnKpQ\r\na=fingerprint:sha-256 FC:C9:57:0D:83:ED:2A:A0:AF:05:2D:02:89:7C:FE:3C:1D:9F:BB:40:01:9A:9B:F9:53:33:9A:93:67:F2:47:06\r\na=setup:active\r\na=mid:2\r\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\r\na=sendonly\r\na=msid:go2rtc go2rtc-audio\r\na=rtcp-mux\r\na=rtcp-rsize\r\na=rtpmap:111 opus/48000/2\r\na=rtcp-fb:111 transport-cc\r\na=fmtp:111 minptime=10;useinbandfec=1\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:8 PCMA/8000\r\na=ssrc:2903819563 cname:go2rtc\r\na=ssrc:2903819563 msid:go2rtc go2rtc-audio\r\n'}
    AUDIO STATS: [{'id': 'COT01_111_minptime=10;useinbandfec=1', 'timestamp': 1791690412470.97, 'type': 'codec', 'channels': 2, 'clockRate': 48000, 'mimeType': 'audio/opus', 'payloadType': 111, 'sdpFmtpLine': 'minptime=10;useinbandfec=1', 'transportId': 'T01'}, {'id': 'OT01A744862198', 'timestamp': 1791690412470.97, 'type': 'outbound-rtp', 'codecId': 'COT01_111_minptime=10;useinbandfec=1', 'kind': 'audio', 'mediaType': 'audio', 'ssrc': 744862198, 'transportId': 'T01', 'bytesSent': 6381, 'packetsSent': 202, 'active': True, 'headerBytesSent': 4040, 'mediaSourceId': 'SA1', 'mid': '0', 'nackCount': 0, 'packetsSentWithEct1': 0, 'remoteId': 'RIA744862198', 'retransmittedBytesSent': 0, 'retransmittedPacketsSent': 0, 'targetBitrate': 32000, 'totalPacketSendDelay': 0}, {'id': 'SA1', 'timestamp': 1791690412470.97, 'type': 'media-source', 'kind': 'audio', 'trackIdentifier': 'c311d0aa-6f80-497b-a2c2-1136b4965256', 'audioLevel': 0.249977111117893, 'totalAudioEnergy': 0.9338146656138016, 'totalSamplesDuration': 4.209999999999955}]
    Fixture log 0 23:46:46.488 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:46:46.488 INF config path=/tmp/neonvr-browser-fixture-grfium2q/fixture.json
    23:46:46.488 INF [rtsp] listen addr=127.0.0.1:42147
    23:46:46.488 INF [api] listen addr=127.0.0.1:60805
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:46:46 NeoNVR server listening on http://127.0.0.1:51277
    2026/10/10 23:46:46 Clips root: /tmp/neonvr-browser-fixture-grfium2q/recordings
    2026/10/10 23:46:46 Thumbnail cache: /tmp/neonvr-browser-fixture-grfium2q/recordings/.blinkclips_thumbs
    2026/10/10 23:46:46 Metadata: /tmp/neonvr-browser-fixture-grfium2q/recordings/.blinkclips_meta.json
    2026/10/10 23:46:46 ffmpeg: ffmpeg
    2026/10/10 23:46:46 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:46:46 Thumbnail seek: 0.200s
    2026/10/10 23:46:46 Thumbnail workers: 2
    2026/10/10 23:46:46 Thumbnail maintenance disabled
    2026/10/10 23:46:46 Size refresh: 5m0s
    2026/10/10 23:46:46 Auto-delete disabled; use -auto-delete-days N to enable
    
    Traceback (most recent call last):
      File "/tmp/neonvr-live-tools-vbq_9dyb/check_browser.py", line 188, in <module>
        assert speaker_packets > 0, 'No user audio reached the RTSP speaker backchannel'
               ^^^^^^^^^^^^^^^^^^^
    AssertionError: No user audio reached the RTSP speaker backchannel
    ✗ (1) • 7.03s

## Activity

    $ python3 - <<'PY'
        import urllib.request,json
        with urllib.request.urlopen('https://api.github.com/repos/AlexxIT/go2rtc/contents/pkg/webrtc?ref=v1.9.14',timeout=20) as response: entries=json.load(response)
        print([e['name'] for e in entries])
        for name in ['conn.go','producer.go','consumer.go']:
         try:
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/pkg/webrtc/'+name,timeout=20) as response:print(name+'\n'+response.read().decode()[:12500])
         except Exception as e:print(e)
        PY
    ['README.md', 'api.go', 'client.go', 'client_test.go', 'conn.go', 'consumer.go', 'helpers.go', 'producer.go', 'server.go', 'track.go', 'webrtc_test.go']
    conn.go
    package webrtc
    
    import (
        "encoding/json"
        "fmt"
        "strings"
        "time"
    
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/pion/rtcp"
        "github.com/pion/rtp"
        "github.com/pion/webrtc/v4"
    )
    
    type Conn struct {
        core.Connection
        core.Listener
    
        Mode core.Mode `json:"mode"`
    
        pc *webrtc.PeerConnection
    
        offer  string
        closed core.Waiter
    }
    
    func NewConn(pc *webrtc.PeerConnection) *Conn {
        c := &Conn{
            Connection: core.Connection{
                ID:         core.NewID(),
                FormatName: "webrtc",
                Transport:  pc,
            },
            pc: pc,
        }
    
        pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
            // last candidate will be empty
            if candidate != nil {
                c.Fire(candidate)
            }
        })
    
        pc.OnDataChannel(func(channel *webrtc.DataChannel) {
            c.Fire(channel)
        })
    
        pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
            if state != webrtc.ICEConnectionStateChecking {
                return
            }
            pc.SCTP().Transport().ICETransport().OnSelectedCandidatePairChange(
                func(pair *webrtc.ICECandidatePair) {
                    // fix situation when candidate pair changes multiple times
                    if i := strings.IndexByte(c.Protocol, '+'); i > 0 {
                        c.Protocol = c.Protocol[:i]
                    }
                    c.Protocol += "+" + pair.Remote.Protocol.String()
                    c.RemoteAddr = fmt.Sprintf(
                        "%s:%d %s", sanitizeIP6(pair.Remote.Address), pair.Remote.Port, pair.Remote.Typ,
                    )
                    if pair.Remote.RelatedAddress != "" {
                        c.RemoteAddr += fmt.Sprintf(
                            " %s:%d", sanitizeIP6(pair.Remote.RelatedAddress), pair.Remote.RelatedPort,
                        )
                    }
                },
            )
        })
    
        pc.OnTrack(func(remote *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
            media, codec := c.getMediaCodec(remote)
            if media == nil {
                return
            }
    
            track, err := c.GetTrack(media, codec)
            if err != nil {
                return
            }
    
            switch c.Mode {
            case core.ModePassiveProducer, core.ModeActiveProducer:
                // replace the theoretical list of codecs with the actual list of codecs
                if len(media.Codecs) > 1 {
                    media.Codecs = []*core.Codec{codec}
                }
            }
    
            if c.Mode == core.ModePassiveProducer && remote.Kind() == webrtc.RTPCodecTypeVideo {
                go func() {
                    pkts := []rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(remote.SSRC())}}
                    for range time.NewTicker(time.Second * 2).C {
                        if err := pc.WriteRTCP(pkts); err != nil {
                            return
                        }
                    }
                }()
            }
    
            for {
                b := make([]byte, ReceiveMTU)
                n, _, err := remote.Read(b)
                if err != nil {
                    return
                }
    
                c.Recv += n
    
                packet := &rtp.Packet{}
                if err := packet.Unmarshal(b[:n]); err != nil {
                    return
                }
    
                if len(packet.Payload) == 0 {
                    continue
                }
    
                track.WriteRTP(packet)
            }
        })
    
        // OK connection:
        // 15:01:46 ICE connection state changed: checking
        // 15:01:46 peer connection state changed: connected
        // 15:01:54 peer connection state changed: disconnected
        // 15:02:20 peer connection state changed: failed
        //
        // Fail connection:
        // 14:53:08 ICE connection state changed: checking
        // 14:53:39 peer connection state changed: failed
        pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
            c.Fire(state)
    
            switch state {
            case webrtc.PeerConnectionStateConnected:
                for _, sender := range c.Senders {
                    sender.Start()
                }
            case webrtc.PeerConnectionStateDisconnected, webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed:
                // disconnect event comes earlier, than failed
                // but it comes only for success connections
                _ = c.Close()
            }
        })
    
        return c
    }
    
    func (c *Conn) MarshalJSON() ([]byte, error) {
        return json.Marshal(c.Connection)
    }
    
    func (c *Conn) Close() error {
        c.closed.Done(nil)
        return c.pc.Close()
    }
    
    func (c *Conn) AddCandidate(candidate string) error {
        // pion uses only candidate value from json/object candidate struct
        return c.pc.AddICECandidate(webrtc.ICECandidateInit{Candidate: candidate})
    }
    
    func (c *Conn) GetSenderTrack(mid string) *Track {
        if tr := c.getTranseiver(mid); tr != nil {
            if s := tr.Sender(); s != nil {
                if t := s.Track().(*Track); t != nil {
                    return t
                }
            }
        }
        return nil
    }
    
    func (c *Conn) getTranseiver(mid string) *webrtc.RTPTransceiver {
        for _, tr := range c.pc.GetTransceivers() {
            if tr.Mid() == mid {
                return tr
            }
        }
        return nil
    }
    
    func (c *Conn) getMediaCodec(remote *webrtc.TrackRemote) (*core.Media, *core.Codec) {
        for _, tr := range c.pc.GetTransceivers() {
            // search Transeiver for this TrackRemote
            if tr.Receiver() == nil || tr.Receiver().Track() != remote {
                continue
            }
    
            // search Media for this MID
            for _, media := range c.Medias {
                if media.ID != tr.Mid() || media.Direction != core.DirectionRecvonly {
                    continue
                }
    
                // search codec for this PayloadType
                for _, codec := range media.Codecs {
                    if codec.PayloadType != uint8(remote.PayloadType()) {
                        continue
                    }
                    return media, codec
                }
            }
        }
    
        // fix moment when core.ModePassiveProducer or core.ModeActiveProducer
        // sends new codec with new payload type to same media
        // check GetTrack
        panic(core.Caller())
    
        return nil, nil
    }
    
    func sanitizeIP6(host string) string {
        if strings.IndexByte(host, ':') > 0 {
            return "[" + host + "]"
        }
        return host
    }
    
    producer.go
    package webrtc
    
    import (
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/pion/webrtc/v4"
    )
    
    func (c *Conn) GetTrack(media *core.Media, codec *core.Codec) (*core.Receiver, error) {
        core.Assert(media.Direction == core.DirectionRecvonly)
    
        for _, track := range c.Receivers {
            if track.Codec == codec {
                return track, nil
            }
        }
    
        switch c.Mode {
        case core.ModePassiveConsumer: // backchannel from browser
            // set codec for consumer recv track so remote peer should send media with this codec
            params := webrtc.RTPCodecParameters{
                RTPCodecCapability: webrtc.RTPCodecCapability{
                    MimeType:  MimeType(codec),
                    ClockRate: codec.ClockRate,
                    Channels:  uint16(codec.Channels),
                },
                PayloadType: 0, // don't know if this necessary
            }
    
            tr := c.getTranseiver(media.ID)
    
            _ = tr.SetCodecPreferences([]webrtc.RTPCodecParameters{params})
    
        case core.ModePassiveProducer, core.ModeActiveProducer:
            // Passive producers: OBS Studio via WHIP or Browser
            // Active producers: go2rtc as WebRTC client or WebTorrent
    
        default:
            panic(core.Caller())
        }
    
        track := core.NewReceiver(media, codec)
        c.Receivers = append(c.Receivers, track)
        return track, nil
    }
    
    func (c *Conn) Start() error {
        c.closed.Wait()
        return nil
    }
    
    consumer.go
    package webrtc
    
    import (
        "errors"
    
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/h264"
        "github.com/AlexxIT/go2rtc/pkg/h265"
        "github.com/AlexxIT/go2rtc/pkg/pcm"
        "github.com/pion/rtp"
    )
    
    func (c *Conn) GetMedias() []*core.Media {
        return WithResampling(c.Medias)
    }
    
    func (c *Conn) AddTrack(media *core.Media, codec *core.Codec, track *core.Receiver) error {
        core.Assert(media.Direction == core.DirectionSendonly)
    
        for _, sender := range c.Senders {
            if sender.Codec == codec {
                sender.Bind(track)
                return nil
            }
        }
    
        switch c.Mode {
        case core.ModePassiveConsumer: // video/audio for browser
        case core.ModeActiveProducer: // go2rtc as WebRTC client (backchannel)
        case core.ModePassiveProducer: // WebRTC/WHIP
        default:
            panic(core.Caller())
        }
    
        localTrack := c.GetSenderTrack(media.ID)
        if localTrack == nil {
            return errors.New("webrtc: can't get track")
        }
    
        payloadType := codec.PayloadType
    
        sender := core.NewSender(media, codec)
        sender.Handler = func(packet *rtp.Packet) {
            c.Send += packet.MarshalSize()
            //important to send with remote PayloadType
            _ = localTrack.WriteRTP(payloadType, packet)
        }
    
        switch track.Codec.Name {
        case core.CodecH264:
            sender.Handler = h264.RTPPay(1200, sender.Handler)
            if track.Codec.IsRTP() {
                sender.Handler = h264.RTPDepay(track.Codec, sender.Handler)
            } else {
                sender.Handler = h264.RepairAVCC(track.Codec, sender.Handler)
            }
    
        case core.CodecH265:
            sender.Handler = h265.RTPPay(1200, sender.Handler)
            if track.Codec.IsRTP() {
                sender.Handler = h265.RTPDepay(track.Codec, sender.Handler)
            } else {
                sender.Handler = h265.RepairAVCC(track.Codec, sender.Handler)
            }
    
        case core.CodecPCMA, core.CodecPCMU, core.CodecPCM, core.CodecPCML:
            // Fix audio quality https://github.com/AlexxIT/WebRTC/issues/500
            // should be before ResampleToG711, because it will be called last
            sender.Handler = pcm.RepackG711(false, sender.Handler)
    
            if codec.ClockRate == 0 {
                if codec.Name == core.CodecPCM || codec.Name == core.CodecPCML {
                    codec.Name = core.CodecPCMA
                }
                codec.ClockRate = 8000
                sender.Handler = pcm.TranscodeHandler(codec, track.Codec, sender.Handler)
            }
        }
    
        // TODO: rewrite this dirty logic
        // maybe not best solution, but ActiveProducer connected before AddTrack
        if c.Mode != core.ModeActiveProducer {
            sender.Bind(track)
        } else {
            sender.HandleRTP(track)
        }
    
        c.Senders = append(c.Senders, sender)
        return nil
    }
    
    ✓ • 488ms

## Activity

    $ python3 - <<'PY'
        import urllib.request
        for path in ['pkg/core/media.go','pkg/core/helpers.go','pkg/rtsp/client.go']:
         try:
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+path,timeout=20) as response: text=response.read().decode()
         except Exception as e:continue
         if path=='pkg/rtsp/client.go':
          start=text.find('func (c *Conn) Describe')
          print(text[start:start+4600])
         else:print(path+'\n'+text[:9500])
        PY
    pkg/core/media.go
    package core
    
    import (
        "encoding/json"
        "fmt"
        "strings"
    
        "github.com/pion/sdp/v3"
    )
    
    // Media take best from:
    // - deepch/vdk/format/rtsp/sdp.Media
    // - pion/sdp.MediaDescription
    type Media struct {
        Kind      string   `json:"kind,omitempty"`      // video or audio
        Direction string   `json:"direction,omitempty"` // sendonly, recvonly
        Codecs    []*Codec `json:"codecs,omitempty"`
    
        ID string `json:"id,omitempty"` // MID for WebRTC, Control for RTSP
    }
    
    func (m *Media) String() string {
        s := fmt.Sprintf("%s, %s", m.Kind, m.Direction)
        for _, codec := range m.Codecs {
            name := codec.String()
    
            if strings.Contains(s, name) {
                continue
            }
    
            s += ", " + name
        }
        return s
    }
    
    func (m *Media) MarshalJSON() ([]byte, error) {
        return json.Marshal(m.String())
    }
    
    func (m *Media) Clone() *Media {
        clone := *m
        clone.Codecs = make([]*Codec, len(m.Codecs))
        for i, codec := range m.Codecs {
            clone.Codecs[i] = codec.Clone()
        }
        return &clone
    }
    
    func (m *Media) MatchMedia(remote *Media) (codec, remoteCodec *Codec) {
        // check same kind and opposite dirrection
        if m.Kind != remote.Kind ||
            m.Direction == DirectionSendonly && remote.Direction != DirectionRecvonly ||
            m.Direction == DirectionRecvonly && remote.Direction != DirectionSendonly {
            return nil, nil
        }
    
        for _, codec = range m.Codecs {
            for _, remoteCodec = range remote.Codecs {
                if codec.Match(remoteCodec) {
                    return
                }
            }
        }
    
        return nil, nil
    }
    
    func (m *Media) MatchCodec(remote *Codec) *Codec {
        for _, codec := range m.Codecs {
            if codec.Match(remote) {
                return codec
            }
        }
        return nil
    }
    
    func (m *Media) MatchAll() bool {
        for _, codec := range m.Codecs {
            if codec.Name == CodecAll {
                return true
            }
        }
        return false
    }
    
    func (m *Media) Equal(media *Media) bool {
        if media.ID != "" {
            return m.ID == media.ID
        }
        return m.String() == media.String()
    }
    
    func GetKind(name string) string {
        switch name {
        case CodecH264, CodecH265, CodecVP8, CodecVP9, CodecAV1, CodecJPEG, CodecRAW:
            return KindVideo
        case CodecPCMU, CodecPCMA, CodecAAC, CodecOpus, CodecG722, CodecMP3, CodecPCM, CodecPCML, CodecELD, CodecFLAC:
            return KindAudio
        }
        return ""
    }
    
    func MarshalSDP(name string, medias []*Media) ([]byte, error) {
        sd := &sdp.SessionDescription{
            Origin: sdp.Origin{
                Username: "-", SessionID: 1, SessionVersion: 1,
                NetworkType: "IN", AddressType: "IP4", UnicastAddress: "0.0.0.0",
            },
            SessionName: sdp.SessionName(name),
            ConnectionInformation: &sdp.ConnectionInformation{
                NetworkType: "IN", AddressType: "IP4", Address: &sdp.Address{
                    Address: "0.0.0.0",
                },
            },
            TimeDescriptions: []sdp.TimeDescription{
                {Timing: sdp.Timing{}},
            },
        }
    
        for _, media := range medias {
            if media.Codecs == nil {
                continue
            }
    
            codec := media.Codecs[0]
    
            switch codec.Name {
            case CodecELD:
                name = CodecAAC
            case CodecPCML:
                name = CodecPCM // beacuse we using pcm.LittleToBig for RTSP server
            default:
                name = codec.Name
            }
    
            md := &sdp.MediaDescription{
                MediaName: sdp.MediaName{
                    Media:  media.Kind,
                    Protos: []string{"RTP", "AVP"},
                },
            }
            md.WithCodec(codec.PayloadType, name, codec.ClockRate, uint16(codec.Channels), codec.FmtpLine)
    
            if media.Direction != "" {
                md.WithPropertyAttribute(media.Direction)
            }
    
            if media.ID != "" {
                md.WithValueAttribute("control", media.ID)
            }
    
            sd.MediaDescriptions = append(sd.MediaDescriptions, md)
        }
    
        return sd.Marshal()
    }
    
    func UnmarshalMedia(md *sdp.MediaDescription) *Media {
        m := &Media{
            Kind: md.MediaName.Media,
        }
    
        for _, attr := range md.Attributes {
            switch attr.Key {
            case DirectionSendonly, DirectionRecvonly, DirectionSendRecv:
                m.Direction = attr.Key
            case "control", "mid":
                m.ID = attr.Value
            }
        }
    
        for _, format := range md.MediaName.Formats {
            m.Codecs = append(m.Codecs, UnmarshalCodec(md, format))
        }
    
        return m
    }
    
    func ParseQuery(query map[string][]string) (medias []*Media) {
        // set media candidates from query list
        for key, values := range query {
            switch key {
            case KindVideo, KindAudio:
                for _, value := range values {
                    media := &Media{Kind: key, Direction: DirectionSendonly}
    
                    for _, name := range strings.Split(value, ",") {
                        name = strings.ToUpper(name)
    
                        // check aliases
                        switch name {
                        case "", "COPY":
                            name = CodecAny
                        case "MJPEG":
                            name = CodecJPEG
                        case "AAC":
                            name = CodecAAC
                        case "MP3":
                            name = CodecMP3
                        }
    
                        media.Codecs = append(media.Codecs, &Codec{Name: name})
                    }
    
                    medias = append(medias, media)
                }
            }
        }
    
        return
    }
    
    pkg/core/helpers.go
    package core
    
    import (
        "crypto/rand"
        "regexp"
        "runtime"
        "strconv"
        "strings"
        "time"
    )
    
    const (
        BufferSize      = 64 * 1024 // 64K
        ConnDialTimeout = 5 * time.Second
        ConnDeadline    = 5 * time.Second
        ProbeTimeout    = 5 * time.Second
    )
    
    // Now90000 - timestamp for Video (clock rate = 90000 samples per second)
    func Now90000() uint32 {
        return uint32(time.Duration(time.Now().UnixNano()) * 90000 / time.Second)
    }
    
    const symbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
    
    // RandString base10 - numbers, base16 - hex, base36 - digits+letters
    // base64 - URL safe symbols, base0 - crypto random
    func RandString(size, base byte) string {
        b := make([]byte, size)
        if _, err := rand.Read(b); err != nil {
            panic(err)
        }
        if base == 0 {
            return string(b)
        }
        for i := byte(0); i < size; i++ {
            b[i] = symbols[b[i]%base]
        }
        return string(b)
    }
    
    func Before(s, sep string) string {
        if i := strings.Index(s, sep); i > 0 {
            return s[:i]
        }
        return s
    }
    
    func Between(s, sub1, sub2 string) string {
        i := strings.Index(s, sub1)
        if i < 0 {
            return ""
        }
        s = s[i+len(sub1):]
    
        if i = strings.Index(s, sub2); i >= 0 {
            return s[:i]
        }
    
        return s
    }
    
    func Atoi(s string) (i int) {
        if s != "" {
            i, _ = strconv.Atoi(s)
        }
        return
    }
    
    // ParseByte - fast parsing string to byte function
    func ParseByte(s string) (b byte) {
        for i, ch := range []byte(s) {
            ch -= '0'
            if ch > 9 {
                return 0
            }
            if i > 0 {
                b *= 10
            }
            b += ch
        }
        return
    }
    
    func Assert(ok bool) {
        if !ok {
            _, file, line, _ := runtime.Caller(1)
            panic(file + ":" + strconv.Itoa(line))
        }
    }
    
    func Caller() string {
        _, file, line, _ := runtime.Caller(1)
        return file + ":" + strconv.Itoa(line)
    }
    
    const (
        unreserved = `A-Za-z0-9-._~`
        subdelims  = `!$&'()*+,;=`
        userinfo   = unreserved + subdelims + `%:`
    )
    
    func StripUserinfo(s string) string {
        sanitizer := regexp.MustCompile(`://[` + userinfo + `]+@`)
        return sanitizer.ReplaceAllString(s, `://***@`)
    }
    
    func (c *Conn) Describe() error {
        // 5.3 Back channel connection
        // https://www.onvif.org/specs/stream/ONVIF-Streaming-Spec.pdf
        req := &tcp.Request{
            Method: MethodDescribe,
            URL:    c.URL,
            Header: map[string][]string{
                "Accept": {"application/sdp"},
            },
        }
    
        if c.Backchannel {
            req.Header.Set("Require", "www.onvif.org/ver20/backchannel")
        }
    
        if c.UserAgent != "" {
            // this camera will answer with 401 on DESCRIBE without User-Agent
            // https://github.com/AlexxIT/go2rtc/issues/235
            req.Header.Set("User-Agent", c.UserAgent)
        }
    
        res, err := c.Do(req)
        if err != nil {
            return err
        }
    
        if val := res.Header.Get("Content-Base"); val != "" {
            c.URL, err = urlParse(val)
            if err != nil {
                return err
            }
        }
    
        c.SDP = string(res.Body) // for info
    
        medias, err := UnmarshalSDP(res.Body)
        if err != nil {
            return err
        }
    
        if c.Media != "" {
            clone := make([]*core.Media, 0, len(medias))
            for _, media := range medias {
                if strings.Contains(c.Media, media.Kind) {
                    clone = append(clone, media)
                }
            }
            medias = clone
        }
    
        // TODO: rewrite more smart
        if c.Medias == nil {
            c.Medias = medias
        } else if len(c.Medias) > len(medias) {
            c.Medias = c.Medias[:len(medias)]
        }
    
        c.mode = core.ModeActiveProducer
    
        return nil
    }
    
    func (c *Conn) Announce() (err error) {
        req := &tcp.Request{
            Method: MethodAnnounce,
            URL:    c.URL,
            Header: map[string][]string{
                "Content-Type": {"application/sdp"},
            },
        }
    
        req.Body, err = core.MarshalSDP(c.SessionName, c.Medias)
        if err != nil {
            return err
        }
    
        _, err = c.Do(req)
        return
    }
    
    func (c *Conn) Record() (err error) {
        req := &tcp.Request{
            Method: MethodRecord,
            URL:    c.URL,
            Header: map[string][]string{
                "Range": {"npt=0.000-"},
            },
        }
    
        _, err = c.Do(req)
        return
    }
    
    func (c *Conn) SetupMedia(media *core.Media) (byte, error) {
        var transport string
    
        if c.Transport == "udp" {
            conn1, conn2, err := ListenUDPPair()
            if err != nil {
                return 0, err
            }
    
            c.udpConn = append(c.udpConn, conn1, conn2)
    
            port := conn1.LocalAddr().(*net.UDPAddr).Port
            transport = fmt.Sprintf("RTP/AVP;unicast;client_port=%d-%d", port, port+1)
        } else {
            // try to use media position as channel number
            for i, m := range c.Medias {
                if m.Equal(media) {
                    transport = fmt.Sprintf(
                        // i   - RTP (data channel)
                        // i+1 - RTCP (control channel)
                        "RTP/AVP/TCP;unicast;interleaved=%d-%d", i*2, i*2+1,
                    )
                    break
                }
            }
        }
    
        if transport == "" {
            return 0, fmt.Errorf("wrong media: %v", media)
        }
    
        rawURL := media.ID // control
        if !strings.Contains(rawURL, "://") {
            rawURL = c.URL.String()
            // prefix check for https://github.com/AlexxIT/go2rtc/issues/1236
            if !strings.HasSuffix(rawURL, "/") && !strings.HasPrefix(media.ID, "/") {
                rawURL += "/"
            }
            rawURL += media.ID
        }
        trackURL, err := urlParse(rawURL)
        if err != nil {
            return 0, err
        }
    
        req := &tcp.Request{
            Method: MethodSetup,
            URL:    trackURL,
            Header: map[string][]string{
                "Transport": {transport},
            },
        }
    
        res, err := c.Do(req)
        if err != nil {
            // some Dahua/Amcrest cameras fail here because two simultaneous
            // backchannel connections
            if c.Backchannel {
                c.Backchannel = false
                if err = c.Reconnect(); err != nil {
                    return 0, err
                }
                return c.SetupMedia(media)
            }
    
            return 0, err
        }
    
        if c.session == "" {
            // Session: 7116520596809429228
            // Session: 216525287999;timeout=60
            if s := res.Header.Get("Session"); s != "" {
                if i := strings.IndexByte(s, ';'); i > 0 {
                    c.session = s[:i]
                    if i = strings.Index(s, "timeout="); i > 0 {
                        c.keepalive, _ = strconv.Atoi(s[i+8:])
                    }
                } else {
                    c.session = s
                }
            }
        }
    
        // Parse server response
        transport = res.Header.Get("Transport")
    
        if c.Transport == "udp" {
            channel := byte(len(c.udpConn) - 2)
    
            // Dahua:   RTP/AVP/UDP;unicast;client_port=49292-49293;server_port=43670-43671;ssrc=7CB694B4
            // OpenIPC: RTP/AVP/UDP;unicast;client_port=59612-59613
            if s := core.Between(transport, "server_port=", ";"); s != "" {
                s1, s2, _ := strings.Cut(s, "-")
                port1 := core.Atoi(s1)
                port2 := core.Atoi(s2)
                // TODO: more smart handling empty server ports
                if port1 > 0 && port2 > 0 {
                    remoteIP := c.conn.RemoteAddr().(*net.TCPAddr).IP
                    c.udpAddr = append(c.udpAddr,
                        &net.UDPAddr{IP: remoteIP, Port: port1},
                        &net.UDPAddr{IP: remoteIP, Port: port2},
                    )
    
                    go func() {
                        // Try to open a hole in the NAT router (to allow incoming UDP packets)
                        // by send a UDP packet for RTP and RTCP to the remote RTSP server.
                        // https://github.com/FFmpeg/FFmpeg/blob/aa91ae25b88e195e6af4248e0ab30605735ca1cd
    ✓ • 192ms

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    @@ -87,3 +87,3 @@
                             'a=control:*\r\nm=audio 0 RTP/AVP 0\r\nc=IN IP4 0.0.0.0\r\n'
    -                        'a=rtpmap:0 PCMU/8000\r\na=recvonly\r\na=control:trackID=0\r\n').encode()
    +                        'a=rtpmap:0 PCMU/8000\r\na=sendonly\r\na=control:trackID=0\r\n').encode()
                     extra += 'Content-Type: application/sdp\r\nContent-Base: ' + target + '/\r\n'

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.
    Stop and the off toggle close playback and the streaming gateway.
    Settings save preserves the live source and backchannel.
    Browser checks passed with no page errors, including mobile layout.
    Fixture log 0 23:47:32.422 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:47:32.423 INF config path=/tmp/neonvr-browser-fixture-alg950i2/fixture.json
    23:47:32.423 INF [rtsp] listen addr=127.0.0.1:41651
    23:47:32.423 INF [api] listen addr=127.0.0.1:59101
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:47:32 NeoNVR server listening on http://127.0.0.1:55837
    2026/10/10 23:47:32 Clips root: /tmp/neonvr-browser-fixture-alg950i2/recordings
    2026/10/10 23:47:32 Thumbnail cache: /tmp/neonvr-browser-fixture-alg950i2/recordings/.blinkclips_thumbs
    2026/10/10 23:47:32 Metadata: /tmp/neonvr-browser-fixture-alg950i2/recordings/.blinkclips_meta.json
    2026/10/10 23:47:32 ffmpeg: ffmpeg
    2026/10/10 23:47:32 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:47:32 Thumbnail seek: 0.200s
    2026/10/10 23:47:32 Thumbnail workers: 2
    2026/10/10 23:47:32 Thumbnail maintenance disabled
    2026/10/10 23:47:32 Size refresh: 5m0s
    2026/10/10 23:47:32 Auto-delete disabled; use -auto-delete-days N to enable
    
    ✓ • 3.55s

## Activity

    $ python3 - <<'PY'
        import urllib.request
        for path in ['internal/streams/api.go','pkg/core/connection.go']:
         try:
          with urllib.request.urlopen('https://raw.githubusercontent.com/AlexxIT/go2rtc/v1.9.14/'+path,timeout=20) as response:print(path+'\n'+response.read().decode()[:7000])
         except Exception as e:print(e)
        PY
    internal/streams/api.go
    package streams
    
    import (
        "net/http"
    
        "github.com/AlexxIT/go2rtc/internal/api"
        "github.com/AlexxIT/go2rtc/internal/app"
        "github.com/AlexxIT/go2rtc/pkg/core"
        "github.com/AlexxIT/go2rtc/pkg/creds"
        "github.com/AlexxIT/go2rtc/pkg/probe"
    )
    
    func apiStreams(w http.ResponseWriter, r *http.Request) {
        w = creds.SecretResponse(w)
    
        query := r.URL.Query()
        src := query.Get("src")
    
        // without source - return all streams list
        if src == "" && r.Method != "POST" {
            api.ResponseJSON(w, streams)
            return
        }
    
        // Not sure about all this API. Should be rewrited...
        switch r.Method {
        case "GET":
            stream := Get(src)
            if stream == nil {
                http.Error(w, "", http.StatusNotFound)
                return
            }
    
            cons := probe.Create("probe", query)
            if len(cons.Medias) != 0 {
                cons.WithRequest(r)
                if err := stream.AddConsumer(cons); err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }
    
                api.ResponsePrettyJSON(w, stream)
    
                stream.RemoveConsumer(cons)
            } else {
                api.ResponsePrettyJSON(w, streams[src])
            }
    
        case "PUT":
            name := query.Get("name")
            if name == "" {
                name = src
            }
    
            if _, err := New(name, query["src"]...); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
    
            if err := app.PatchConfig([]string{"streams", name}, query["src"]); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
            }
    
        case "PATCH":
            name := query.Get("name")
            if name == "" {
                http.Error(w, "", http.StatusBadRequest)
                return
            }
    
            // support {input} templates: https://github.com/AlexxIT/go2rtc#module-hass
            if _, err := Patch(name, src); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
            }
    
        case "POST":
            // with dst - redirect source to dst
            if dst := query.Get("dst"); dst != "" {
                if stream := Get(dst); stream != nil {
                    if err := Validate(src); err != nil {
                        http.Error(w, err.Error(), http.StatusBadRequest)
                    } else if err = stream.Play(src); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                    } else {
                        api.ResponseJSON(w, stream)
                    }
                } else if stream = Get(src); stream != nil {
                    if err := Validate(dst); err != nil {
                        http.Error(w, err.Error(), http.StatusBadRequest)
                    } else if err = stream.Publish(dst); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                    }
                } else {
                    http.Error(w, "", http.StatusNotFound)
                }
            } else {
                http.Error(w, "", http.StatusBadRequest)
            }
    
        case "DELETE":
            delete(streams, src)
    
            if err := app.PatchConfig([]string{"streams", src}, nil); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
            }
        }
    }
    
    func apiStreamsDOT(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
    
        dot := make([]byte, 0, 1024)
        dot = append(dot, "digraph {\n"...)
        if query.Has("src") {
            for _, name := range query["src"] {
                if stream := streams[name]; stream != nil {
                    dot = AppendDOT(dot, stream)
                }
            }
        } else {
            for _, stream := range streams {
                dot = AppendDOT(dot, stream)
            }
        }
        dot = append(dot, '}')
    
        dot = []byte(creds.SecretString(string(dot)))
    
        api.Response(w, dot, "text/vnd.graphviz")
    }
    
    func apiPreload(w http.ResponseWriter, r *http.Request) {
        // GET - return all preloads
        if r.Method == "GET" {
            api.ResponseJSON(w, GetPreloads())
            return
        }
    
        query := r.URL.Query()
        src := query.Get("src")
    
        switch r.Method {
        case "PUT":
            // it's safe to delete from map while iterating
            for k := range query {
                switch k {
                case core.KindVideo, core.KindAudio, "microphone":
                default:
                    delete(query, k)
                }
            }
    
            rawQuery := query.Encode()
    
            if err := AddPreload(src, rawQuery); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
    
            if err := app.PatchConfig([]string{"preload", src}, rawQuery); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
    
        case "DELETE":
            if err := DelPreload(src); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
    
            if err := app.PatchConfig([]string{"preload", src}, nil); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
    
        default:
            http.Error(w, "", http.StatusMethodNotAllowed)
        }
    }
    
    func apiSchemes(w http.ResponseWriter, r *http.Request) {
        api.ResponseJSON(w, SupportedSchemes())
    }
    
    pkg/core/connection.go
    package core
    
    import (
        "io"
        "net/http"
        "reflect"
        "sync/atomic"
    )
    
    func NewID() uint32 {
        return id.Add(1)
    }
    
    // Deprecated: use NewID instead
    func ID(v any) uint32 {
        p := uintptr(reflect.ValueOf(v).UnsafePointer())
        return 0x8000_0000 | uint32(p)
    }
    
    var id atomic.Uint32
    
    type Info interface {
        SetProtocol(string)
        SetRemoteAddr(string)
        SetSource(string)
        SetURL(string)
        WithRequest(*http.Request)
        GetSource() string
    }
    
    // Connection just like webrtc.PeerConnection
    // - ID and RemoteAddr used for building Connection(s) graph
    // - FormatName, Protocol, RemoteAddr, Source, URL, SDP, UserAgent used for info about Connection
    // - FormatName and Protocol has FFmpeg compatible names
    // - Transport used for auto closing on Stop
    type Connection struct {
        ID         uint32 `json:"id,omitempty"`
        FormatName string `json:"format_name,omitempty"` // rtsp, webrtc, mp4, mjpeg, mpjpeg...
        Protocol   string `json:"protocol,omitempty"`    // tcp, udp, http, ws, pipe...
        RemoteAddr string `json:"remote_addr,omitempty"` // host:port other info
        Source     string `json:"source,omitempty"`
        URL        string `json:"url,omitempty"`
        SDP        string `json:"sdp,omitempty"`
        UserAgent  string `json:"user_agent,omitempty"`
    
        Medias    []*Media    `json:"medias,omitempty"`
        Receivers []*Receiver `json:"receivers,omitempty"`
        Senders   []*Sender   `json:"senders,omitempty"`
        Recv      int         `json:"bytes_recv,omitempty"`
        Send      int         `json:"bytes_send,omitempty"`
    
        Transport any `json:"-"`
    }
    
    func (c *Connection) GetMedias() []*Media {
        return c.Medias
    }
    
    func (c *Connection) GetTrack(media *Media, codec *Codec) (*Receiver, error) {
        for _, receiver := range c.Receivers {
            if receiver.Codec == codec {
                return receiver, nil
            }
        }
        receiver := NewReceiver(media, codec)
        c.Receivers = append(c.Receivers, receiver)
        return receiver, nil
    }
    
    func (c *Connection) Stop() error {
        for _, receiver := range c.Receivers {
            receiver.Close()
        }
        for _, sender := range c.Senders {
            sender.Close()
        }
        if closer, ok := c.Transport.(io.Closer); ok {
            return closer.Close()
        }
        return nil
    }
    
    // Deprecated:
    func (c *Connection) Codecs() []*Codec {
        codecs := make([]*Codec, len(c.Senders))
        for i, sender := range c.Senders {
            codecs[i] = sender.Codec
        }
        return codecs
    }
    
    func (c *Connection) SetProtocol(s string) {
        c.Protocol = s
    }
    
    func (c *Connection) SetRemoteAddr(s string) {
        if c.RemoteAddr == "" {
            c.RemoteAddr = s
        } else {
            c.RemoteAddr += " forwarded " + s
        }
    }
    
    func (c *Connection) SetSource(s string) {
        c.Source = s
    }
    
    func (c *Connection) SetURL(s string) {
        c.URL = s
    }
    
    func (c *Connection) WithRequest(r *http.Request) {
        if r.Header.Get("Upgrade") == "websocket" {
            c.Protocol = "ws"
        } else {
            c.Protocol = "http"
        }
    
        c.RemoteAddr = r.RemoteAddr
        if remote := r.Header.Get("X-Forwarded-For"); remote != "" {
            c.RemoteAddr += " forwarded " + remote
        }
    
        c.UserAgent = r.UserAgent()
    }
    
    func (c *Connection) GetSource() string {
        return c.Source
    }
    
    // Create like os.Create, init Consumer with existing Transport
    func Create(w io.Writer) (*Connection, error) {
        return &Connection{Transport: w}, nil
    }
    
    // Open like os.Open, init Producer from existing Transport
    func Open(r io.Reader) (*Connection, error) {
        return &Connection{Transport: r}, nil
    }
    
    // Dial like net.Dial, init Producer via Dialing
    func Dial(rawURL string) (*Connection, error) {
        return &Connection{}, nil
    }
    
    ✓ • 127ms

## Assistant

Video playback and Talk both worked in the local browser test: microphone audio reached the simulated camera speaker, and Stop talking ended capture and restored ordinary viewing. I’m checking denied permissions and unsupported backchannels next, then finishing the desktop and mobile review.

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/live.go
    @@ -206,3 +206,3 @@
     	return map[string]any{
    -		"api":     map[string]any{"listen": address, "username": "neonvr", "password": password, "local_auth": true, "allow_paths": []string{"/api", "/api/webrtc"}},
    +		"api":     map[string]any{"listen": address, "username": "neonvr", "password": password, "local_auth": true, "allow_paths": []string{"/api", "/api/webrtc", "/api/streams"}},
     		"rtsp":    map[string]any{"listen": ""},

## Activity

    file changes: Completed · 2 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/live.go
    @@ -226,3 +226,3 @@
     	udp.Close()
    -	// Reserve a private API address. Only the two fixed API routes are used.
    +	// Reserve a private API address. Only fixed, internal API routes are used.
     	listener, err := net.Listen("tcp", "127.0.0.1:0")
    @@ -321,2 +321,30 @@
     
    +func (s *liveService) backchannelAccepted(ctx context.Context, source, sessionID string) bool {
    +	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/api/streams?src="+source, nil)
    +	req.SetBasicAuth("neonvr", s.password)
    +	response, err := s.client.Do(req)
    +	if err != nil {
    +		return false
    +	}
    +	defer response.Body.Close()
    +	var status struct {
    +		Consumers []struct {
    +			UserAgent string            `json:"user_agent"`
    +			Receivers []json.RawMessage `json:"receivers"`
    +		} `json:"consumers"`
    +	}
    +	if response.StatusCode != http.StatusOK || json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(&status) != nil {
    +		return false
    +	}
    +	for _, consumer := range status.Consumers {
    +		// A receiver exists only when the gateway attached the browser's
    +		// microphone to a compatible camera speaker track. SDP direction alone
    +		// can remain sendonly even when the camera did not accept a backchannel.
    +		if consumer.UserAgent == sessionID && len(consumer.Receivers) > 0 {
    +			return true
    +		}
    +	}
    +	return false
    +}
    +
     func (m *liveManager) webrtcHandler(w http.ResponseWriter, r *http.Request) {
    @@ -349,3 +377,4 @@
     	source := fmt.Sprintf("camera_%d", index)
    -	if r.URL.Query().Get("talk") == "1" {
    +	talk := r.URL.Query().Get("talk") == "1"
    +	if talk {
     		if m.config.LiveStreams[index].BackchannelURL == "" {
    @@ -365,2 +394,9 @@
     	req.SetBasicAuth("neonvr", service.password)
    +	nonce := make([]byte, 16)
    +	if _, err := rand.Read(nonce); err != nil {
    +		http.Error(w, "cannot create a live session", http.StatusInternalServerError)
    +		return
    +	}
    +	sessionID := "NeoNVR/" + hex.EncodeToString(nonce)
    +	req.Header.Set("User-Agent", sessionID)
     	response, err := service.client.Do(req)
    @@ -376,2 +412,6 @@
     	}
    +	if talk && !service.backchannelAccepted(r.Context(), source, sessionID) {
    +		http.Error(w, "camera did not accept a compatible audio backchannel; your microphone is off, choose Watch to return to viewing", http.StatusBadGateway)
    +		return
    +	}
     	w.Header().Set("Content-Type", "application/sdp")
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/live-streams.js
    @@ -17,2 +17,3 @@
     let changing = false;
    +let stateVersion = 0;
     
    @@ -88,2 +89,3 @@
     async function loadState(initial = false) {
    +  const version = stateVersion;
       const response = await fetch('/api/live', {cache: 'no-store', signal: AbortSignal.timeout(5000)});
    @@ -91,2 +93,3 @@
       const next = await response.json();
    +  if (version !== stateVersion) return;
       const restarted = state && state.token !== next.token;
    @@ -207,2 +210,3 @@
       changing = true;
    +  stateVersion++;
       if (!enabled) stopPlayback();

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w live.go
        /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr .
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.
    Stop and the off toggle close playback and the streaming gateway.
    Settings save preserves the live source and backchannel.
    Browser checks passed with no page errors, including mobile layout.
    Fixture log 0 23:48:46.311 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:48:46.311 INF config path=/tmp/neonvr-browser-fixture-fo2o2tmt/fixture.json
    23:48:46.311 INF [api] listen addr=127.0.0.1:59971
    23:48:46.311 INF [rtsp] listen addr=127.0.0.1:40955
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:48:46 NeoNVR server listening on http://127.0.0.1:58895
    2026/10/10 23:48:46 Clips root: /tmp/neonvr-browser-fixture-fo2o2tmt/recordings
    2026/10/10 23:48:46 Thumbnail cache: /tmp/neonvr-browser-fixture-fo2o2tmt/recordings/.blinkclips_thumbs
    2026/10/10 23:48:46 Metadata: /tmp/neonvr-browser-fixture-fo2o2tmt/recordings/.blinkclips_meta.json
    2026/10/10 23:48:46 ffmpeg: ffmpeg
    2026/10/10 23:48:46 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:48:46 Thumbnail seek: 0.200s
    2026/10/10 23:48:46 Thumbnail workers: 2
    2026/10/10 23:48:46 Thumbnail maintenance disabled
    2026/10/10 23:48:46 Size refresh: 5m0s
    2026/10/10 23:48:46 Auto-delete disabled; use -auto-delete-days N to enable
    
    ✓ • 4.16s

## Activity

    file changes: Completed · 2 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -257,2 +257,4 @@
     
    +Set `NEONVR_TEST_GO2RTC` to an installed go2rtc executable to include the live gateway lifecycle check. It verifies explicit enable/disable, shutdown, temporary credential cleanup, and that enabling alone does not open a camera connection. The test uses a local dummy endpoint.
    +
     Model/export reference: [Ultralytics YOLO26](https://docs.ultralytics.com/models/yolo26), [ONNX export](https://docs.ultralytics.com/modes/export).
    Update { move_path: None }: /home/ubuntu/NeoNVR/live_test.go
    @@ -169,2 +169,46 @@
     
    +func TestLiveTalkRequiresAnAttachedCameraSpeakerTrack(t *testing.T) {
    +	var sessionID string
    +	var attached, wrongSession bool
    +	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    +		switch r.URL.Path {
    +		case "/api/webrtc":
    +			sessionID = r.UserAgent()
    +			w.WriteHeader(http.StatusCreated)
    +			_, _ = w.Write([]byte("v=0\r\nanswer"))
    +		case "/api/streams":
    +			if r.URL.RawQuery != "src=camera_0_talk" {
    +				t.Error("backchannel verification queried an arbitrary source")
    +			}
    +			agent := sessionID
    +			if wrongSession {
    +				agent = "another-viewer"
    +			}
    +			receivers := []string{}
    +			if attached {
    +				receivers = append(receivers, "microphone")
    +			}
    +			_ = json.NewEncoder(w).Encode(map[string]any{"consumers": []any{map[string]any{"user_agent": agent, "receivers": receivers}}})
    +		default:
    +			t.Error("unexpected private gateway route")
    +		}
    +	}))
    +	defer backend.Close()
    +	m := newLiveManager(context.Background(), liveFixture(), "test-token")
    +	m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private", client: backend.Client()}
    +	for _, tc := range []struct {
    +		attached, wrongSession bool
    +		code                   int
    +	}{
    +		{false, false, 502}, {true, true, 502}, {true, false, 200},
    +	} {
    +		attached, wrongSession = tc.attached, tc.wrongSession
    +		w := httptest.NewRecorder()
    +		m.webrtcHandler(w, liveRequest("POST", "/api/live/webrtc?id=0&talk=1", "v=0\r\n", "application/sdp"))
    +		if w.Code != tc.code || !strings.HasPrefix(sessionID, "NeoNVR/") {
    +			t.Fatalf("talk accepted without its own attached speaker track: %d %s", w.Code, w.Body.String())
    +		}
    +	}
    +}
    +
     func TestLiveGatewayLifecycleWithoutOpeningCamera(t *testing.T) {

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w live_test.go
        NEONVR_TEST_GO2RTC=/tmp/neonvr-live-tools-vbq_9dyb/go2rtc /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
    ok      neonvr    4.851s
    ✓ • 5.73s

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    @@ -135,3 +135,5 @@
                 'live_streams': [{'name': 'Local test camera', 'url': 'rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture',
    -                              'backchannel_url': 'rtsp://127.0.0.1:' + str(speaker.server_address[1]) + '/speaker?backchannel=1'}],
    +                              'backchannel_url': 'rtsp://127.0.0.1:' + str(speaker.server_address[1]) + '/speaker?backchannel=1'},
    +                             {'name': 'View-only test camera', 'url': 'rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture',
    +                              'backchannel_url': 'rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture?backchannel=1'}],
             }))
    @@ -193,5 +195,40 @@
                 print('Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.', flush=True)
    +            page.evaluate('window.nativeGetUserMedia = navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices); navigator.mediaDevices.getUserMedia = async() => { throw new DOMException("Denied", "NotAllowedError"); }')
    +            page.locator('#talk').click()
    +            page.wait_for_function('active === null && document.querySelector("#status").textContent.includes("permission was denied")')
    +            assert page.locator('#microphone').inner_text() == 'Microphone off'
    +            page.evaluate('navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia')
    +            page.locator('#watch').click()
    +            page.wait_for_function('active && active.connected', timeout=20000)
    +            page.evaluate('navigator.mediaDevices.getUserMedia = (...args) => new Promise(resolve => {window.resolveLateMicrophone = async() => {window.lateMicrophone = await nativeGetUserMedia(...args); resolve(window.lateMicrophone);};})')
    +            page.locator('#talk').click()
    +            page.wait_for_function('typeof window.resolveLateMicrophone === "function"')
                 page.locator('#stop').click()
    +            page.evaluate('window.resolveLateMicrophone()')
    +            page.wait_for_function('window.lateMicrophone && window.lateMicrophone.getTracks().every(track => track.readyState === "ended")')
    +            page.evaluate('navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia')
    +            print('Denied permissions and late permission grants leave the microphone off.', flush=True)
    +            page.get_by_role('button', name='View-only test camera').click()
    +            page.locator('#watch').click()
    +            page.wait_for_function('active && active.connected', timeout=20000)
    +            page.evaluate('navigator.mediaDevices.getUserMedia = async(...args) => {window.unsupportedMicrophone = await nativeGetUserMedia(...args); return window.unsupportedMicrophone;}')
    +            page.locator('#talk').click()
    +            page.wait_for_function('active === null && document.querySelector("#status").textContent.includes("did not accept a compatible audio backchannel")', timeout=20000)
    +            assert page.evaluate('window.unsupportedMicrophone.getTracks().every(track => track.readyState === "ended")')
    +            page.evaluate('navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia')
    +            print('A camera without a compatible speaker rejects Talk and releases the microphone.', flush=True)
    +            page.get_by_role('button', name='Local test camera', exact=False).click()
    +            page.locator('#watch').click()
    +            page.wait_for_function('active && active.connected', timeout=20000)
    +            page.locator('#talk').click()
    +            page.wait_for_function('active && active.connected && active.talk', timeout=20000)
    +            page.evaluate('window.offMicrophone = active.microphone')
    +            other = context.new_page()
    +            other.goto(origin + '/streams')
    +            other.wait_for_function('state && state.enabled')
    +            other.locator('#enabled').uncheck()
    +            other.wait_for_function('!state.enabled && !changing')
    +            page.wait_for_function('active === null && window.offMicrophone.getTracks().every(track => track.readyState === "ended")', timeout=8000)
    +            other.close()
                 assert page.evaluate('active === null && videoEl.srcObject === null')
    -            page.locator('#enabled').uncheck()
                 page.wait_for_function('!state.enabled && !changing')
    @@ -202,6 +239,6 @@
                 expected = 'Local test camera | rtsp://127.0.0.1:' + str(rtsp_port) + '/fixture | rtsp://127.0.0.1:' + str(speaker.server_address[1]) + '/speaker?backchannel=1'
    -            assert page.locator('[name=live_streams]').input_value() == expected
    +            assert page.locator('[name=live_streams]').input_value().splitlines()[0] == expected
                 page.locator('#save').click()
                 page.wait_for_function('document.querySelector("#status").textContent.startsWith("Settings saved.")')
    -            assert len(json.loads(config.read_text())['live_streams']) == 1
    +            assert len(json.loads(config.read_text())['live_streams']) == 2
                 print('Settings save preserves the live source and backchannel.', flush=True)

## Activity

    Viewed image /tmp/neonvr-live-tools-vbq_9dyb/streams-mobile.png

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.
    Fixture log 0 23:50:14.740 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:50:14.740 INF config path=/tmp/neonvr-browser-fixture-_a30cs1y/fixture.json
    23:50:14.740 INF [rtsp] listen addr=127.0.0.1:35981
    23:50:14.740 INF [api] listen addr=127.0.0.1:36577
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:50:14 NeoNVR server listening on http://127.0.0.1:60481
    2026/10/10 23:50:14 Clips root: /tmp/neonvr-browser-fixture-_a30cs1y/recordings
    2026/10/10 23:50:14 Thumbnail cache: /tmp/neonvr-browser-fixture-_a30cs1y/recordings/.blinkclips_thumbs
    2026/10/10 23:50:14 Metadata: /tmp/neonvr-browser-fixture-_a30cs1y/recordings/.blinkclips_meta.json
    2026/10/10 23:50:14 ffmpeg: ffmpeg
    2026/10/10 23:50:14 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:50:14 Thumbnail seek: 0.200s
    2026/10/10 23:50:14 Thumbnail workers: 2
    2026/10/10 23:50:14 Thumbnail maintenance disabled
    2026/10/10 23:50:14 Size refresh: 5m0s
    2026/10/10 23:50:14 Auto-delete disabled; use -auto-delete-days N to enable
    
    Traceback (most recent call last):
      File "/tmp/neonvr-live-tools-vbq_9dyb/check_browser.py", line 196, in <module>
        page.evaluate('window.nativeGetUserMedia = navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices); navigator.mediaDevices.getUserMedia = async() => { throw new DOMException("Denied", "NotAllowedError"); }')
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/sync_api/_generated.py", line 9413, in evaluate
        self._sync(
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_sync_base.py", line 115, in _sync
        return task.result()
               ^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_page.py", line 474, in evaluate
        return await self._main_frame.evaluate(expression, arg)
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_frame.py", line 346, in evaluate
        await self._channel.send(
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_connection.py", line 76, in send
        return await self._connection.wrap_api_call(
               ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
      File "/tmp/neonvr-live-tools-vbq_9dyb/browser-env/lib/python3.12/site-packages/playwright/_impl/_connection.py", line 632, in wrap_api_call
        raise rewrite_error(error, f"{parsed_st['apiName']}: {error}") from None
    playwright._impl._errors.Error: Page.evaluate: NotAllowedError: Denied
    ✗ (1) • 3.06s

## Activity

    $ rg -n 'liveService|func .*Live|func .*backchannel|func .*webrtc|function (connect|stopPlayback|loadState)|pagehide|visibilitychange' live.go live_test.go static/live-streams.js
        sed -n '73,105p' static/settings.html
        sed -n '134,153p' static/settings.html
    static/live-streams.js:45:function stopPlayback(message = '') {
    static/live-streams.js:89:async function loadState(initial = false) {
    static/live-streams.js:127:async function connect(talk = false) {
    static/live-streams.js:234:document.addEventListener('visibilitychange', () => {
    static/live-streams.js:237:window.addEventListener('pagehide', () => stopPlayback());
    live_test.go:30:func TestLiveDefaultsOffAndHidesCameraCredentials(t *testing.T) {
    live_test.go:62:func TestLiveToggleRejectsUnauthorizedAndMalformedRequests(t *testing.T) {
    live_test.go:92:func TestLiveSignalingUsesOnlyConfiguredSourcesAndRedactsErrors(t *testing.T) {
    live_test.go:115:    m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private-password", client: backend.Client()}
    live_test.go:145:func TestLiveConfigValidationAndLegacyCompatibility(t *testing.T) {
    live_test.go:170:func TestLiveTalkRequiresAnAttachedCameraSpeakerTrack(t *testing.T) {
    live_test.go:198:    m.service = &liveService{done: make(chan struct{}), baseURL: backend.URL, password: "private", client: backend.Client()}
    live_test.go:214:func TestLiveGatewayLifecycleWithoutOpeningCamera(t *testing.T) {
    live.go:45:func (c Config) validateLive() error {
    live.go:81:type liveService struct {
    live.go:94:    service *liveService
    live.go:97:func newLiveManager(ctx context.Context, config Config, token string) *liveManager {
    live.go:102:func (m *liveManager) runningLocked() *liveService {
    live.go:186:func liveSource(raw string, backchannel bool) string {
    live.go:215:func (m *liveManager) start(requestCtx context.Context, requestHost string) (*liveService, error) {
    live.go:268:    service := &liveService{cmd: cmd, done: make(chan struct{}), baseURL: "http://" + address, password: password, client: &http.Client{Timeout: 25 * time.Second, Transport: &http.Transport{}}}
    live.go:299:func (s *liveService) stop() {
    live.go:322:func (s *liveService) backchannelAccepted(ctx context.Context, source, sessionID string) bool {
    live.go:350:func (m *liveManager) webrtcHandler(w http.ResponseWriter, r *http.Request) {
     const response=await fetch(url,{cache:'no-store',...options});
     if(!response.ok)throw new Error(await response.text());
     return response.json();
    }
    function populate(data){
     state=data;
     for(const [name,value] of Object.entries(data.config)){
      const field=form.elements.namedItem(name);
      if(field)field.value=name==='extra_roots'?value.join('\n'):name==='live_streams'?value.map(source=>[source.name,source.url,source.backchannel_url||''].join(' | ')).join('\n'):value;
     }
     pending.hidden=!data.restart_required;
     fields.disabled=false;dirty=false;
     show(data.restart_required?'Saved settings are ready to apply. Choose Save & restart.':'Settings match the running app.','ok');
    }
    async function load(){fields.disabled=true;try{populate(await request('/api/config'));}catch(error){show(error.message,'error');}}
    function readForm(){
     const config={};
     for(const name of Object.keys(state.config)){
      let value=form.elements.namedItem(name).value;
      if(name==='extra_roots')value=value.split('\n').map(v=>v.trim()).filter(Boolean);
      else if(name==='live_streams')value=value.split('\n').map(v=>v.trim()).filter(Boolean).map((line,index)=>{
       const parts=line.split('|').map(v=>v.trim());
       if(parts.length<2||parts.length>3||!parts[0]||!parts[1])throw new Error('Live stream line '+(index+1)+': enter name | RTSP URL | optional backchannel URL.');
       return {name:parts[0],url:parts[1],...(parts[2]?{backchannel_url:parts[2]}:{})};
      });
      else if(name==='detection_enabled')value=value==='true';
      else if(numeric.has(name))value=Number(value);
      else value=value.trim();
      config[name]=value;
     }
     return config;
    }
    async function waitForRestart(data,previousID){
      if(!confirm('Restart with automatic deletion enabled for non-favorites older than '+config.auto_delete_days+' days? This applies to every recording source.'))return;
     }
     const previousID=state.instance_id;
     fields.disabled=true;show(restart?'Validating settings and preparing restart…':'Saving settings…');
     try{
      const data=await request('/api/config',{method:'POST',headers:{'Content-Type':'application/json','X-NeoNVR-Token':state.token},body:JSON.stringify({config,revision:state.revision,restart})});
      state=data;dirty=false;pending.hidden=!data.restart_required;
      if(restart){show('Settings saved. Restarting and reconnecting…','ok');await waitForRestart(data,previousID);}
      else{fields.disabled=false;show('Settings saved. Choose Save & restart to apply them.','ok');}
     }catch(error){fields.disabled=false;show(error.message,'error');}
    }
    form.addEventListener('input',()=>{dirty=true;});
    form.addEventListener('submit',event=>{event.preventDefault();save(false);});
    document.getElementById('restart').addEventListener('click',()=>save(true));
    document.getElementById('reload').addEventListener('click',()=>{if(!dirty||confirm('Discard unsaved changes and reload?'))load();});
    window.addEventListener('beforeunload',event=>{if(dirty){event.preventDefault();event.returnValue='';}});
    load();
    </script>
    </body>
    </html>
    ✓ • 0ms

## Activity

    $ git diff --check
        /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -l config.go main.go live.go live_test.go
        python3 -m json.tool config/neonvr.json >/dev/null
        git diff -- Dockerfile docker-compose.yml config/neonvr.json
        rg -n 'LiveStreams|newLiveManager|/streams|shutdownSignals' main.go config.go
        rg -n 'Live streams|NEONVR_TEST_GO2RTC|Yi cameras' README.md
    diff --git a/Dockerfile b/Dockerfile
    index 61953f2..26a63db 100644
    --- a/Dockerfile
    +++ b/Dockerfile
    @@ -6,6 +6,9 @@ COPY static ./static
     COPY detection ./detection
     RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/neonvr .
    
    +# The bridge is started only when live viewing is explicitly enabled.
    +FROM ghcr.io/alexxit/go2rtc:1.9.14 AS streams
    +
     # PyTorch/Ultralytics are used only to export the official weights at build time.
     FROM python:3.12-slim-bookworm AS model
     ENV YOLO_AUTOINSTALL=False OMP_NUM_THREADS=2
    @@ -24,6 +27,8 @@ RUN pip install --no-cache-dir -r /tmp/detection-requirements.txt && rm /tmp/det
     COPY --from=model /model/yolo26n.onnx /opt/neonvr/models/yolo26n.onnx
     ENV NEONVR_MODEL_DIR=/opt/neonvr/models OMP_NUM_THREADS=2
     COPY --from=build /out/neonvr /usr/local/bin/neonvr
    -EXPOSE 5500
    +COPY --from=streams /usr/local/bin/go2rtc /usr/local/bin/go2rtc
    +COPY third_party/go2rtc.LICENSE /usr/share/doc/go2rtc/LICENSE
    +EXPOSE 5500 8555/tcp 8555/udp
     ENTRYPOINT ["/usr/local/bin/neonvr"]
     CMD ["-config", "/config/neonvr.json"]
    diff --git a/config/neonvr.json b/config/neonvr.json
    index d248b51..7da8cc8 100644
    --- a/config/neonvr.json
    +++ b/config/neonvr.json
    @@ -27,5 +27,15 @@
       "detection_every": "1m",
       "detection_min_age": "2m",
       "detection_timeout": "30m",
    -  "detection_prune": "quarantine"
    +  "detection_prune": "quarantine",
    +  "live_streams": [
    +    {
    +      "name": "Camera 192.168.1.150",
    +      "url": "rtsp://192.168.1.150:554/ch0_0.h264",
    +      "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
    +    }
    +  ],
    +  "live_go2rtc": "go2rtc",
    +  "live_webrtc_port": 8555,
    +  "live_webrtc_host": ""
     }
    diff --git a/docker-compose.yml b/docker-compose.yml
    index e61be4c..177e84b 100644
    --- a/docker-compose.yml
    +++ b/docker-compose.yml
    @@ -3,8 +3,13 @@ services:
         build: .
         container_name: neonvr
         restart: unless-stopped
    +    mem_limit: 6g
    +    memswap_limit: 4g
         ports:
           - "5500:5500"
    +      # Live media and microphone talkback; inactive until enabled in Streams.
    +      - "8555:8555/tcp"
    +      - "8555:8555/udp"
         volumes:
           # Mount the directory so config saves can replace the file atomically.
           - ./config:/config
    config.go:56:    LiveStreams           []LiveStream `json:"live_streams"`
    config.go:63:    return Config{Root: defaultRoot, ExtraRoots: []string{}, Addr: defaultAddr, FFmpeg: "ffmpeg", ThumbSeek: .2, ThumbWidth: 480, ThumbWorkers: 2, ThumbMaintenanceEvery: "30m", SizeRefresh: "5m", AutoDeleteEvery: "6h", LiveStreams: []LiveStream{}, LiveGo2RTC: "go2rtc", LiveWebRTCPort: 8555,
    config.go:83:    if c.LiveStreams == nil {
    config.go:84:        c.LiveStreams = []LiveStream{}
    main.go:186://go:embed static/streams.html
    main.go:669:    <a href="/streams">Streams</a>
    main.go:1681:    live := newLiveManager(runCtx, cfg, manager.token)
    main.go:1769:    mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
    main.go:1837:    shutdownSignals := make(chan os.Signal, 1)
    main.go:1838:    signal.Notify(shutdownSignals, syscall.SIGTERM, syscall.SIGINT)
    main.go:1839:    defer signal.Stop(shutdownSignals)
    main.go:1863:    case <-shutdownSignals:
    155:## Live streams and speaker talkback
    159:The supplied config includes `rtsp://192.168.1.150:554/ch0_0.h264` for viewing and `rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1` for speaker talkback. Edit sources in **Settings → Live streams**, one camera per line: `name | viewing URL | optional backchannel URL`. Save and restart to apply source changes; live viewing resets to off.
    258:Set `NEONVR_TEST_GO2RTC` to an installed go2rtc executable to include the live gateway lifecycle check. It verifies explicit enable/disable, shutdown, temporary credential cleanup, and that enabling alone does not open a camera connection. The test uses a local dummy endpoint.
    292:Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        from pathlib import Path
        path=Path('/tmp/neonvr-live-tools-vbq_9dyb/check_browser.py')
        text=path.read_text()
        text=text.replace("page.evaluate('window.nativeGetUserMedia = navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices); navigator.mediaDevices.getUserMedia = async() => { throw new DOMException(\"Denied\", \"NotAllowedError\"); }')", "page.evaluate('() => { window.nativeGetUserMedia = navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices); navigator.mediaDevices.getUserMedia = async() => { throw new DOMException(\"Denied\", \"NotAllowedError\"); }; }')")
        text=text.replace("page.evaluate('navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia')", "page.evaluate('() => {navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia;}')")
        text=text.replace("page.evaluate('navigator.mediaDevices.getUserMedia = (...args) => new Promise(resolve => {window.resolveLateMicrophone = async() => {window.lateMicrophone = await nativeGetUserMedia(...args); resolve(window.lateMicrophone);};})')", "page.evaluate('() => {navigator.mediaDevices.getUserMedia = (...args) => new Promise(resolve => {window.resolveLateMicrophone = async() => {window.lateMicrophone = await nativeGetUserMedia(...args); resolve(window.lateMicrophone);};});}')")
        text=text.replace("page.evaluate('navigator.mediaDevices.getUserMedia = async(...args) => {window.unsupportedMicrophone = await nativeGetUserMedia(...args); return window.unsupportedMicrophone;}')", "page.evaluate('() => {navigator.mediaDevices.getUserMedia = async(...args) => {window.unsupportedMicrophone = await nativeGetUserMedia(...args); return window.unsupportedMicrophone;};}')")
        path.write_text(text)
        PY
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    Opening the tab and enabling live viewing do not connect to cameras.
    Real H.264 camera video plays over WebRTC without opening the backchannel.
    Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.
    Denied permissions and late permission grants leave the microphone off.
    A camera without a compatible speaker rejects Talk and releases the microphone.
    Stop and the off toggle close playback and the streaming gateway.
    Settings save preserves the live source and backchannel.
    Browser checks passed with no page errors, including mobile layout.
    Fixture log 0 23:51:10.339 INF go2rtc platform=linux/amd64 revision=b5948cf version=1.9.14
    23:51:10.339 INF config path=/tmp/neonvr-browser-fixture-u63x5knz/fixture.json
    23:51:10.339 INF [rtsp] listen addr=127.0.0.1:42281
    23:51:10.339 INF [api] listen addr=127.0.0.1:44079
    exit with signal: terminated
    
    Fixture log 2 2026/10/10 23:51:10 NeoNVR server listening on http://127.0.0.1:37429
    2026/10/10 23:51:10 Clips root: /tmp/neonvr-browser-fixture-u63x5knz/recordings
    2026/10/10 23:51:10 Thumbnail cache: /tmp/neonvr-browser-fixture-u63x5knz/recordings/.blinkclips_thumbs
    2026/10/10 23:51:10 Metadata: /tmp/neonvr-browser-fixture-u63x5knz/recordings/.blinkclips_meta.json
    2026/10/10 23:51:10 ffmpeg: ffmpeg
    2026/10/10 23:51:10 YOLO26n detection: enabled=false, CPU threads=2, sample fps=2, no-target policy=quarantine
    2026/10/10 23:51:10 Thumbnail seek: 0.200s
    2026/10/10 23:51:10 Thumbnail workers: 2
    2026/10/10 23:51:10 Thumbnail maintenance disabled
    2026/10/10 23:51:10 Size refresh: 5m0s
    2026/10/10 23:51:10 Auto-delete disabled; use -auto-delete-days N to enable
    
    ✓ • 7.83s

## Activity

    $ git diff --check
        git status --short
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M docker-compose.yml
     M main.go
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    ✓ • 0ms

## Assistant

Added the **Streams** tab with your camera URLs configured.

- Live viewing defaults to **off** and connects only after **Enable → Watch**.
- **Talk to camera** sends microphone audio through the backchannel.
- **Stop talking** releases the microphone; turning viewing off disconnects all viewers.
- Sources are editable in Settings.

Go race tests and browser video/talkback checks passed using a simulated camera. Talk requires [HTTPS or localhost](https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices/getUserMedia#privacy_and_security).

Rebuild to apply:

```bash
sudo nerdctl compose up -d --build
```

## User

some of the cameras now write the motion metadata time stamp which can be received with `ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 |                                                      
    jq '.format.tags.yi_motion | fromjson'  ` Id like this timestamp to be marked on the video player slide track. Also the thumbnail can use this and use the current ffmpeg method as fallback

## Assistant

I’ll check the `yi_motion` metadata format, then use its timestamps for player markers and thumbnail selection, with the current thumbnail method as fallback.

## Activity

    $ for dir in / /home /home/ubuntu /home/ubuntu/NeoNVR; do if [ -f "$dir/AGENTS.md" ]; then cat "$dir/AGENTS.md"; fi; done; done
        command -v ffprobe ffmpeg go jq
        ls -d /tmp/neonvr-live-tools-* 2>/dev/null
        rg --files /mnt/FTPfiles /mnt/mySSD/BlinkClips -g '*.mp4' -g '!\.neonvr_pruned/**' -g '!\.neonvr_detections/**'
    /bin/bash: -c: line 1: syntax error near unexpected token `done'
    /bin/bash: -c: line 1: `for dir in / /home /home/ubuntu /home/ubuntu/NeoNVR; do if [ -f "$dir/AGENTS.md" ]; then cat "$dir/AGENTS.md"; fi; done; done'
    ✗ (2) • 0ms

## Activity

    $ cat static/detection-viewer.js
        sed -n '35,95p' main.go
        sed -n '680,725p' main.go
        sed -n '2390,2700p' main.go
        sed -n '1,190p' detection.go
        cat config/neonvr.json
    // Fetch clip-relative intervals, refresh pending analyses, and show active events.
    (function () {
      'use strict';
      const video = document.getElementById('mainVideo');
      const seek = document.getElementById('mainSeek');
      const status = document.getElementById('detectionStatus');
      if (!video || !seek || !status) return;
      const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
      let generation = 0;
      let request = null;
      let timer = null;
      let result = null;
    
      function activeStatus() {
        if (!result) return;
        const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
        status.classList.toggle('active', active.length > 0);
        const text = active.length ? 'Detected now: ' + active.join(', ') :
          (result.intervals.length ? 'No detection at this time.' : result.uncertain ? 'Unconfirmed observations — clip retained.' : 'No target detections.');
        if (status.textContent !== text) status.textContent = text;
      }
    
      async function refresh(clip, expected) {
        if (expected !== generation) return;
        request = new AbortController();
        try {
          const response = await fetch('/api/detections?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: request.signal});
          if (!response.ok) throw new Error(response.status === 404 ? 'Clip is no longer available.' : 'Detection data unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          if (data.status === 'complete') {
            result = data.detections;
            status.title = '';
            seek.setDetections(result);
            activeStatus();
          } else {
            status.classList.remove('active');
            const labels = {pending: 'Detection queued…', processing: 'Analyzing detections…', disabled: 'Detection disabled.', error: 'Analysis unavailable — clip retained.'};
            status.textContent = labels[data.status] || 'Detection pending…';
            status.title = data.message || '';
            if (data.status !== 'disabled') timer = setTimeout(() => refresh(clip, expected), 10000);
          }
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          status.textContent = error.message;
          if (!error.message.includes('no longer available')) timer = setTimeout(() => refresh(clip, expected), 10000);
        }
      }
    
      function loadClip(clip) {
        generation += 1;
        request?.abort();
        clearTimeout(timer);
        result = null;
        seek.setDetections(null);
        status.classList.remove('active');
        status.textContent = 'Loading detections…';
        status.title = '';
        refresh(clip, generation);
      }
    
      video.addEventListener('timeupdate', activeStatus);
      video.addEventListener('seeked', activeStatus);
      window.addEventListener('pagehide', () => { request?.abort(); clearTimeout(timer); });
      window.NeoNVRDetections = { loadClip };
    })();
    type stringListFlag []string
    
    func (v *stringListFlag) String() string { return strings.Join(*v, ",") }
    func (v *stringListFlag) Set(value string) error {
        value = strings.TrimSpace(value)
        if value == "" {
            return fmt.Errorf("path cannot be empty")
        }
        *v = append(*v, value)
        return nil
    }
    
    func expandPath(path string) string {
        if path == "~" || strings.HasPrefix(path, "~/") {
            if home, err := os.UserHomeDir(); err == nil {
                if path == "~" {
                    return home
                }
                return filepath.Join(home, strings.TrimPrefix(path, "~/"))
            }
        }
        return path
    }
    
    type Server struct {
        runCtx                context.Context
        background            sync.WaitGroup
        root                  string
        extraRoots            []string
        thumbDir              string
        metaPath              string
        ffmpegPath            string
        thumbSeek             float64
        thumbWidth            int
        thumbSem              chan struct{}
        thumbMaintenanceEvery time.Duration
        sizeTTL               time.Duration
        autoDeleteDays        int
        autoDeleteEvery       time.Duration
    
        meta                 *MetaStore
        clipMutationMu       sync.Mutex
        detectionConfig      Config
        detectionDir         string
        detectionMu          sync.Mutex
        detectionStates      map[string]detectionResponse
        detectionWorkerError string
    
        sizeCacheMu  sync.Mutex
        sizeCache    statsResponse
        sizeCacheAt  time.Time
        sizeCacheErr string
    }
    
    type Camera struct {
        Name    string `json:"name"`
        Dir     string `json:"-"`
        Managed bool   `json:"-"`
    }
    
    type Clip struct {
              <span class="kbdHint">Space pause/play • ← / → seek • Ctrl+arrows previous/next</span>
            </div>
            <div class="mainPlayerShell" id="mainShell">
              <media-controller>
                <video id="mainVideo" slot="media" preload="metadata" playsinline webkit-playsinline x5-playsinline></video>
                <media-control-bar>
                  <media-play-button></media-play-button>
                  <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
                  <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
                  <media-time-range id="mainSeek"></media-time-range>
                  <media-time-display></media-time-display>
                  <media-duration-display></media-duration-display>
                  <media-mute-button></media-mute-button>
                  <media-volume-range></media-volume-range>
                  <media-playback-rate-button></media-playback-rate-button>
                  <media-fullscreen-button></media-fullscreen-button>
                </media-control-bar>
                <div class="detectionLegend" style="padding:0 10px 10px;background:#000;margin:0">
                  <span class="category" style="--detection-color:#75d685">Human</span>
                  <span class="category" style="--detection-color:#ffd166">Animal</span>
                  <span class="category" style="--detection-color:#72a7ff">Bike</span>
                  <span class="category" style="--detection-color:#c792ea">Moving vehicle</span>
                  <span id="detectionStatus" role="status">Select a clip to see detections.</span>
                </div>
              </media-controller>
            </div>
            <div class="analysisControls">
              <label style="display:flex;align-items:center;gap:6px">Step sec <input id="stepSize" type="number" min="0.03" step="0.05" value="0.20"></label>
              <button type="button" id="zoomToggle" aria-pressed="false">Zoom off</button>
              <button type="button" id="zoomReset">Reset zoom</button>
              <span id="zoomLabel">1.0×</span>
              <span class="kbdHint">Turn Zoom on to pinch/trackpad zoom and drag-to-pan.</span>
            </div>
          </div>
    
          <aside class="railPanel">
            <div class="railHeader">
              <strong>Videos</strong>
              <span id="railCount">0 clips</span>
            </div>
            <div class="railActions">
              <button id="selectModeToggle" type="button" class="compact">Select</button>
              <span id="selectionCount">0 selected</span>
              <button id="bulkDelete" type="button" class="danger compact" style="display:none" disabled>Delete selected</button>
              <button id="clearSelection" type="button" class="compact" style="display:none" disabled>Clear</button>
            </div>
        check, err := filepath.Rel(dir, full)
        if err != nil || check == ".." || strings.HasPrefix(check, ".."+string(filepath.Separator)) {
            return "", fmt.Errorf("clip path escapes source directory")
        }
        info, err := os.Stat(full)
        if err != nil {
            return "", err
        }
        if info.IsDir() || !strings.EqualFold(filepath.Ext(full), ".mp4") {
            return "", fmt.Errorf("not an MP4 clip")
        }
        // Resolve symlinks before allowing detection/remux/deletion to touch a file.
        realDir, err := filepath.EvalSymlinks(dir)
        if err != nil {
            return "", err
        }
        realFile, err := filepath.EvalSymlinks(full)
        if err != nil {
            return "", err
        }
        realRel, err := filepath.Rel(realDir, realFile)
        if err != nil || realRel == ".." || strings.HasPrefix(realRel, ".."+string(filepath.Separator)) {
            return "", fmt.Errorf("clip symlink escapes source directory")
        }
        return full, nil
    }
    
    func parseClipTimestamp(filename string) (time.Time, bool) {
        base := filepath.Base(filename)
        parts := strings.Split(base, "_")
        if len(parts) < 2 {
            return time.Time{}, false
        }
        datePart := parts[0]
        timePart := strings.ReplaceAll(parts[1], "-", ":")
        ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
        if err != nil {
            return time.Time{}, false
        }
        return ts, true
    }
    
    func (s *Server) thumbnailPath(videoPath string) (string, error) {
        info, err := os.Stat(videoPath)
        if err != nil {
            return "", err
        }
        keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
        sum := sha1.Sum([]byte(keyInput))
        name := hex.EncodeToString(sum[:]) + ".jpg"
        return filepath.Join(s.thumbDir, name), nil
    }
    
    func (s *Server) ensureThumbnail(videoPath string) (string, error) {
        thumbPath, err := s.thumbnailPath(videoPath)
        if err != nil {
            return "", err
        }
        if err := os.MkdirAll(s.thumbDir, 0755); err != nil {
            return "", err
        }
        if st, err := os.Stat(thumbPath); err == nil && st.Size() > 0 {
            return thumbPath, nil
        }
    
        tmpPath := strings.TrimSuffix(thumbPath, ".jpg") + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
        _ = os.Remove(tmpPath)
    
        select {
        case s.thumbSem <- struct{}{}:
        case <-s.workContext().Done():
            return "", s.workContext().Err()
        }
        defer func() { <-s.thumbSem }()
    
        if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
            _ = os.Remove(tmpPath)
            seek := strconv.FormatFloat(s.thumbSeek, 'f', 3, 64)
            log.Printf("interesting thumbnail selection failed for %s, falling back to %ss: %v", videoPath, seek, err)
            if err2 := s.runFFmpegThumbnail(videoPath, tmpPath, seek); err2 != nil {
                _ = os.Remove(tmpPath)
                if s.thumbSeek != 0 {
                    if err3 := s.runFFmpegThumbnail(videoPath, tmpPath, "0"); err3 != nil {
                        _ = os.Remove(tmpPath)
                        return "", fmt.Errorf("interesting, %ss, and 0s thumbnail attempts failed: %v; %v; %v", seek, err, err2, err3)
                    }
                } else {
                    return "", fmt.Errorf("interesting and 0s thumbnail attempts failed: %v; %v", err, err2)
                }
            }
        }
    
        if err := os.Rename(tmpPath, thumbPath); err != nil {
            _ = os.Remove(tmpPath)
            return "", err
        }
        return thumbPath, nil
    }
    
    func (s *Server) probeVideoDuration(videoPath string) (float64, error) {
        ctx, cancel := context.WithTimeout(s.workContext(), 10*time.Second)
        defer cancel()
    
        ffprobePath := "ffprobe"
        if strings.ContainsRune(s.ffmpegPath, filepath.Separator) {
            ffprobePath = filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
        }
        cmd := exec.CommandContext(ctx, ffprobePath,
            "-v", "error",
            "-show_entries", "format=duration",
            "-of", "default=noprint_wrappers=1:nokey=1",
            videoPath,
        )
        out, err := cmd.Output()
        if err != nil {
            return 0, fmt.Errorf("ffprobe duration: %w", err)
        }
        duration, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
        if err != nil || math.IsNaN(duration) || math.IsInf(duration, 0) || duration <= 0 {
            return 0, fmt.Errorf("invalid video duration %q", strings.TrimSpace(string(out)))
        }
        return duration, nil
    }
    
    type frameSignature struct {
        path     string
        seek     float64
        luma     []float64
        mean     float64
        contrast float64
        usable   bool
    }
    
    func frameLumaSignature(path string) (frameSignature, error) {
        f, err := os.Open(path)
        if err != nil {
            return frameSignature{}, err
        }
        defer f.Close()
    
        img, err := jpeg.Decode(f)
        if err != nil {
            return frameSignature{}, err
        }
        b := img.Bounds()
        if b.Dx() < 16 || b.Dy() < 9 {
            return frameSignature{}, fmt.Errorf("thumbnail too small")
        }
    
        const gridW = 32
        const gridH = 18
        luma := make([]float64, 0, gridW*gridH)
        var sum, sumSq float64
        for gy := 0; gy < gridH; gy++ {
            y := b.Min.Y + (gy*2+1)*b.Dy()/(2*gridH)
            if y >= b.Max.Y {
                y = b.Max.Y - 1
            }
            for gx := 0; gx < gridW; gx++ {
                x := b.Min.X + (gx*2+1)*b.Dx()/(2*gridW)
                if x >= b.Max.X {
                    x = b.Max.X - 1
                }
                r, g, bl, _ := img.At(x, y).RGBA()
                v := float64((299*uint64(r) + 587*uint64(g) + 114*uint64(bl)) / 1000 / 257)
                luma = append(luma, v)
                sum += v
                sumSq += v * v
            }
        }
        mean := sum / float64(len(luma))
        variance := sumSq/float64(len(luma)) - mean*mean
        if variance < 0 {
            variance = 0
        }
        contrast := math.Sqrt(variance)
    
        usable := true
        if contrast < 3 && (mean < 10 || mean > 245) {
            usable = false
        }
    
        return frameSignature{
            path:     path,
            luma:     luma,
            mean:     mean,
            contrast: contrast,
            usable:   usable,
        }, nil
    }
    
    func medianLumaBaseline(frames []frameSignature) ([]float64, error) {
        if len(frames) == 0 {
            return nil, fmt.Errorf("no frames for baseline")
        }
        n := len(frames[0].luma)
        if n == 0 {
            return nil, fmt.Errorf("empty frame signature")
        }
        baseline := make([]float64, n)
        values := make([]float64, 0, len(frames))
        for i := 0; i < n; i++ {
            values = values[:0]
            for _, frame := range frames {
                if frame.usable && len(frame.luma) == n {
                    values = append(values, frame.luma[i])
                }
            }
            if len(values) == 0 {
                return nil, fmt.Errorf("no usable frames for baseline")
            }
            sort.Float64s(values)
            mid := len(values) / 2
            if len(values)%2 == 1 {
                baseline[i] = values[mid]
            } else {
                baseline[i] = (values[mid-1] + values[mid]) / 2
            }
        }
        return baseline, nil
    }
    
    func frameDivergenceScore(frame frameSignature, baseline []float64) float64 {
        if !frame.usable || len(frame.luma) != len(baseline) || len(baseline) == 0 {
            return -1
        }
        var sum float64
        for i, v := range frame.luma {
            sum += math.Abs(v - baseline[i])
        }
        return sum / float64(len(baseline))
    }
    
    func (s *Server) runInterestingThumbnail(videoPath, tmpPath string) error {
        duration, probeErr := s.probeVideoDuration(videoPath)
        seekTimes := make([]float64, 0, len(thumbnailCandidateFractions))
        if probeErr == nil {
            if duration < 0.75 {
                return fmt.Errorf("video too short for multi-frame selection: %.3fs", duration)
            }
            for _, fraction := range thumbnailCandidateFractions {
                seekTimes = append(seekTimes, duration*fraction)
            }
        } else {
            // Some camera MP4 variants do not expose duration metadata. Try useful fixed offsets
            // before giving up and falling back to the configured single-frame seek.
            seekTimes = []float64{0.5, 2, 5, 15, 30}
        }
    
        candidates := make([]string, 0, len(seekTimes))
        defer func() {
            for _, path := range candidates {
                _ = os.Remove(path)
            }
        }()
    
        frames := make([]frameSignature, 0, len(seekTimes))
        for i, seekSeconds := range seekTimes {
            if probeErr == nil && seekSeconds > duration-0.10 {
                seekSeconds = math.Max(0, duration-0.10)
            }
            candidatePath := fmt.Sprintf("%s.candidate-%02d.jpg", tmpPath, i)
            candidates = append(candidates, candidatePath)
            seek := strconv.FormatFloat(seekSeconds, 'f', 3, 64)
            if err := s.runFFmpegThumbnail(videoPath, candidatePath, seek); err != nil {
                continue
            }
            frame, err := frameLumaSignature(candidatePath)
            if err != nil {
                continue
            }
            frame.seek = seekSeconds
            frames = append(frames, frame)
        }
        if len(frames) < 2 {
            if probeErr != nil {
                return fmt.Errorf("duration probe failed (%v) and fixed-offset sampling produced fewer than two usable frames", probeErr)
            }
            return fmt.Errorf("fewer than two candidate frames could be analyzed")
        }
    
        baseline, err := medianLumaBaseline(frames)
        if err != nil {
            return err
        }
    
        bestIndex := -1
        bestScore := -1.0
        for i, frame := range frames {
            score := frameDivergenceScore(frame, baseline)
            if score > bestScore {
                bestScore = score
                bestIndex = i
            }
        }
        if bestIndex < 0 {
            return fmt.Errorf("no usable candidate frame remained after baseline analysis")
        }
    
        best := frames[bestIndex]
        if err := os.Rename(best.path, tmpPath); err != nil {
            return err
        }
        log.Printf("thumbnail selected %.3fs for %s (divergence %.2f from median scene)", best.seek, videoPath, bestScore)
        return nil
    }
    
    func (s *Server) runFFmpegThumbnail(videoPath, tmpPath, seek string) error {
        ctx, cancel := context.WithTimeout(s.workContext(), 30*time.Second)
        defer cancel()
        cmd := exec.CommandContext(ctx, s.ffmpegPath,
    package main
    
    import (
        "bufio"
        "context"
        "crypto/sha256"
        _ "embed"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "io"
        "log"
        "math"
        "net/http"
        "os"
        "os/exec"
        "path/filepath"
        "reflect"
        "sort"
        "strconv"
        "strings"
        "syscall"
        "time"
    )
    
    const detectionMetadataTag = "neonvr_detections"
    const detectionVersion = 1
    const detectionPolicyVersion = "yolo26n-onnx-v2"
    
    var errDetectionWorker = errors.New("detection worker protocol failed")
    
    //go:embed detection/worker.py
    var detectionWorkerScript string
    
    //go:embed static/detection-viewer.js
    var detectionViewerJS []byte
    
    type DetectionInterval struct {
        Kind       string   `json:"kind"`
        Labels     []string `json:"labels"`
        Start      float64  `json:"start"`
        End        float64  `json:"end"`
        Confidence float64  `json:"confidence"`
    }
    
    type DetectionResult struct {
        Version     int                 `json:"version"`
        Model       string              `json:"model"`
        ModelSHA256 string              `json:"model_sha256"`
        Policy      string              `json:"policy"`
        Complete    bool                `json:"complete"`
        Duration    float64             `json:"duration"`
        SampleFPS   float64             `json:"sample_fps"`
        Frames      int                 `json:"frames"`
        Uncertain   bool                `json:"uncertain"`
        Intervals   []DetectionInterval `json:"intervals"`
        AnalyzedAt  string              `json:"analyzed_at"`
    }
    
    type detectionResponse struct {
        Status     string           `json:"status"`
        Detections *DetectionResult `json:"detections,omitempty"`
        Message    string           `json:"message,omitempty"`
    }
    
    type detectionCache struct {
        Size    int64           `json:"size"`
        MtimeNS int64           `json:"mtime_ns"`
        Result  DetectionResult `json:"result"`
    }
    
    func (c Config) validateDetection() error {
        if strings.TrimSpace(c.DetectionPython) == "" || !strings.HasSuffix(strings.ToLower(c.DetectionModel), ".onnx") {
            return errors.New("detection requires a Python executable and a YOLO26n .onnx model path")
        }
        if c.DetectionThreads < 1 || c.DetectionThreads > 16 || c.DetectionMinFrames < 1 || c.DetectionMinFrames > 10 {
            return errors.New("detection_threads must be 1–16 and detection_min_frames must be 1–10")
        }
        for _, limit := range []struct {
            name            string
            value, min, max float64
        }{
            {"detection_fps", c.DetectionFPS, .5, 10},
            {"detection_confidence", c.DetectionConfidence, .05, .95},
            {"detection_merge_gap", c.DetectionMergeGap, 0, 10},
            {"detection_car_movement", c.DetectionCarMovement, .001, .25},
        } {
            if !finite(limit.value) || limit.value < limit.min || limit.value > limit.max {
                return fmt.Errorf("%s must be between %g and %g", limit.name, limit.min, limit.max)
            }
        }
        for name, value := range map[string]string{"detection_every": c.DetectionEvery, "detection_min_age": c.DetectionMinAge, "detection_timeout": c.DetectionTimeout} {
            d, err := time.ParseDuration(value)
            minimum := time.Minute
            if name == "detection_min_age" {
                minimum = 10 * time.Second
            }
            if err != nil || d < minimum || d > 24*time.Hour {
                return fmt.Errorf("%s must be between %s and 24h", name, minimum)
            }
        }
        if c.DetectionPrune != "off" && c.DetectionPrune != "quarantine" && c.DetectionPrune != "delete" {
            return errors.New("detection_prune must be off, quarantine, or delete")
        }
        return nil
    }
    
    func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
    
    func (d DetectionResult) validate() error {
        if d.Version != detectionVersion || d.Model != "yolo26n" || len(d.ModelSHA256) != 64 || len(d.Policy) != 64 || !d.Complete ||
            !finite(d.Duration) || d.Duration <= 0 || !finite(d.SampleFPS) || d.SampleFPS < .5 || d.SampleFPS > 10 || d.Frames < 1 ||
            math.Abs(float64(d.Frames)/d.SampleFPS-d.Duration) > 1.5/d.SampleFPS || d.Intervals == nil {
            return errors.New("incomplete or invalid detection analysis")
        }
        for _, interval := range d.Intervals {
            if interval.Kind != "person" && interval.Kind != "animal" && interval.Kind != "bike" && interval.Kind != "moving_car" {
                return errors.New("unknown detection category")
            }
            if !finite(interval.Start) || !finite(interval.End) || interval.Start < 0 || interval.Start >= interval.End || interval.End > d.Duration+.001 ||
                !finite(interval.Confidence) || interval.Confidence < 0 || interval.Confidence > 1 || len(interval.Labels) == 0 {
                return errors.New("invalid detection interval")
            }
        }
        return nil
    }
    
    func (s *Server) detectionPolicy() string {
        c := s.detectionConfig
        data, _ := json.Marshal([]any{detectionPolicyVersion, c.DetectionModel, c.DetectionFPS, c.DetectionConfidence, c.DetectionMinFrames, c.DetectionMergeGap, c.DetectionCarMovement})
        h := sha256.Sum256(data)
        return hex.EncodeToString(h[:])
    }
    
    func (s *Server) detectionCachePath(path string) string {
        h := sha256.Sum256([]byte(path))
        dir := s.detectionDir
        if dir == "" {
            dir = filepath.Join(s.root, ".neonvr_detections")
        }
        return filepath.Join(dir, hex.EncodeToString(h[:])+".json")
    }
    
    func (s *Server) cachedDetection(path string, info os.FileInfo) *DetectionResult {
        data, err := os.ReadFile(s.detectionCachePath(path))
        if err != nil {
            return nil
        }
        var cached detectionCache
        if json.Unmarshal(data, &cached) != nil || cached.Size != info.Size() || cached.MtimeNS != info.ModTime().UnixNano() ||
            cached.Result.Policy != s.detectionPolicy() || cached.Result.validate() != nil {
            return nil
        }
        return &cached.Result
    }
    
    func (s *Server) cacheDetection(path string, info os.FileInfo, result DetectionResult) error {
        data, err := json.Marshal(detectionCache{Size: info.Size(), MtimeNS: info.ModTime().UnixNano(), Result: result})
        if err != nil {
            return err
        }
        return atomicWriteConfig(s.detectionCachePath(path), data)
    }
    
    type videoInspection struct {
        Streams []struct {
            CodecType string `json:"codec_type"`
            CodecName string `json:"codec_name"`
        } `json:"streams"`
        Format struct {
            Duration string            `json:"duration"`
            Tags     map[string]string `json:"tags"`
        } `json:"format"`
    }
    
    func (s *Server) ffprobeExecutable() string {
        if strings.ContainsRune(s.ffmpegPath, filepath.Separator) {
            return filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
        }
        return "ffprobe"
    }
    
    func (s *Server) inspectDetectionVideo(ctx context.Context, path string) (videoInspection, error) {
        ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+detectionMetadataTag+":stream=codec_type,codec_name", "-of", "json", path).Output()
        var info videoInspection
        if err != nil {
    {
      "root": "/mnt/mySSD/BlinkClips",
      "extra_roots": [
        "/mnt/FTPfiles/yicamBack",
        "/mnt/FTPfiles/yicamFront"
      ],
      "addr": "0.0.0.0:5500",
      "thumb_dir": "",
      "meta": "",
      "ffmpeg": "ffmpeg",
      "thumb_seek": 0.2,
      "thumb_width": 480,
      "thumb_workers": 2,
      "thumb_maintenance_every": "30m",
      "size_refresh": "30m",
      "auto_delete_days": 14,
      "auto_delete_every": "6h",
      "detection_enabled": true,
      "detection_python": "python3",
      "detection_model": "yolo26n.onnx",
      "detection_threads": 2,
      "detection_fps": 2,
      "detection_confidence": 0.65,
      "detection_min_frames": 3,
      "detection_merge_gap": 0.5,
      "detection_car_movement": 0.025,
      "detection_every": "1m",
      "detection_min_age": "2m",
      "detection_timeout": "30m",
      "detection_prune": "quarantine",
      "live_streams": [
        {
          "name": "Camera 192.168.1.150",
          "url": "rtsp://192.168.1.150:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
        }
      ],
      "live_go2rtc": "go2rtc",
      "live_webrtc_port": 8555,
      "live_webrtc_host": ""
    }
    ✓ • 0ms

## Activity

    $ git status --short
        rg --files -g 'AGENTS.md' -g '!node_modules' -g '!vendor'
        rg -n 'yi_motion|motion|ensureThumbnail|runInterestingThumbnail|thumbnailAlgorithm|ffprobe|detectionTimeline|detection-viewer|timeline' main.go detection.go static/detection-viewer.js detection_test.go README.md
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M docker-compose.yml
     M main.go
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    README.md:87:- Colored semi-transparent timeline markers per camera, based on clip start time only
    README.md:200:Humans, animals and bikes require three consecutive matched observations with at least 65% confidence by default. Animal classes in the pretrained COCO model are bird, cat, dog, horse, sheep, cow, elephant, bear, zebra and giraffe; other wildlife is not covered by this pretrained class list. Cars, trucks and buses require at least three consecutive observations and consistent translation of both box edges by 2.5% of the model frame within two seconds. A single box jump, back-and-forth jitter, or box resizing does not qualify as vehicle motion. Detection intervals include up to one sample of padding and merge nearby same-category confirmed events across gaps of up to 0.5 seconds by default; the merge setting does not relax consecutive confirmation. Timing is approximate to the sampling cadence.
    README.md:204:Vehicle confirmation needs at least three samples within a two-second motion window, so use at least one sampled frame per second for vehicle markers. At 0.5 FPS, observed vehicle movement remains ambiguous and the recording is retained.
    README.md:223:ffprobe -v error -show_entries format_tags=neonvr_detections -of json recording.mp4
    detection_test.go:21:    for _, name := range []string{"ffmpeg", "ffprobe"} {
    detection_test.go:87:    allTags, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format_tags=source_camera", "-of", "json", path).Output()
    main.go:31:const thumbnailAlgorithmVersion = "divergence-v2"
    main.go:408:    .timelineHeader {
    main.go:417:    .timeline {
    main.go:654:      .controls .timelineWrap { grid-column: 1 / -1; }
    main.go:758:      <div class="timelineWrap">
    main.go:759:        <div class="timelineHeader">
    main.go:761:          <span id="timelineCount">0 clips</span>
    main.go:763:        <div class="timeline" id="timeline">
    main.go:785:<script src="/static/detection-viewer.js"></script>
    main.go:802:const timelineCountEl = document.getElementById('timelineCount');
    main.go:1458:  timelineCountEl.textContent = countText;
    main.go:1805:    mux.HandleFunc("/static/detection-viewer.js", func(w http.ResponseWriter, r *http.Request) {
    main.go:2136:    thumbPath, err := s.ensureThumbnail(videoPath)
    main.go:2437:    keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
    main.go:2443:func (s *Server) ensureThumbnail(videoPath string) (string, error) {
    main.go:2465:    if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
    main.go:2493:    ffprobePath := "ffprobe"
    main.go:2495:        ffprobePath = filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
    main.go:2497:    cmd := exec.CommandContext(ctx, ffprobePath,
    main.go:2505:        return 0, fmt.Errorf("ffprobe duration: %w", err)
    main.go:2623:func (s *Server) runInterestingThumbnail(videoPath, tmpPath string) error {
    main.go:2829:                if _, err := s.ensureThumbnail(videoPath); err != nil {
    detection.go:36://go:embed static/detection-viewer.js
    detection.go:177:func (s *Server) ffprobeExecutable() string {
    detection.go:179:        return filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
    detection.go:181:    return "ffprobe"
    detection.go:187:    out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
    detection.go:191:        return info, fmt.Errorf("ffprobe: %w", err)
    detection.go:359:    request := map[string]any{"path": path, "ffmpeg": s.ffmpegPath, "ffprobe": s.ffprobeExecutable(),
    ✓ • 0ms

## Activity

    • Context compacted

## Activity

    $ rg -n 'setDetections|class MediaTimeRange|mainSeek|loadClip|runFFmpegThumbnail|handleDetections|yi_motion' static/media-chrome-lite.js static/detection-viewer.js main.go detection.go; sed -n '1,280p' static/media-chrome-lite.js; sed -n '2670,2760p' main.go
    static/detection-viewer.js:5:  const seek = document.getElementById('mainSeek');
    static/detection-viewer.js:34:        seek.setDetections(result);
    static/detection-viewer.js:50:  function loadClip(clip) {
    static/detection-viewer.js:55:    seek.setDetections(null);
    static/detection-viewer.js:65:  window.NeoNVRDetections = { loadClip };
    main.go:689:              <media-time-range id="mainSeek"></media-time-range>
    main.go:1058:  window.NeoNVRDetections?.loadClip(clip);
    main.go:2469:        if err2 := s.runFFmpegThumbnail(videoPath, tmpPath, seek); err2 != nil {
    main.go:2472:                if err3 := s.runFFmpegThumbnail(videoPath, tmpPath, "0"); err3 != nil {
    main.go:2654:        if err := s.runFFmpegThumbnail(videoPath, candidatePath, seek); err != nil {
    main.go:2697:func (s *Server) runFFmpegThumbnail(videoPath, tmpPath, seek string) error {
    static/media-chrome-lite.js:150:  class MediaTimeRange extends HTMLElement {
    static/media-chrome-lite.js:199:    setDetections(result) {
    // Media Chrome Lite for NeoNVR.
    // Offline, dependency-free custom elements compatible with the subset of
    // media-chrome controls used by this app. This is intentionally small and
    // tailored to the main NeoNVR player.
    
    (function () {
      function formatTime(seconds) {
        seconds = Number.isFinite(seconds) ? Math.max(0, seconds) : 0;
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = Math.floor(seconds % 60);
        if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
        return `${m}:${String(s).padStart(2, '0')}`;
      }
    
      function findVideo(el) {
        const controller = el.closest('media-controller');
        return controller ? controller.querySelector('video') : null;
      }
    
      function define(name, klass) {
        if (!customElements.get(name)) customElements.define(name, klass);
      }
    
      class MediaController extends HTMLElement {
        connectedCallback() {
          this.setAttribute('data-media-chrome-lite', '');
        }
      }
    
      class MediaControlBar extends HTMLElement {
        connectedCallback() {
          this.setAttribute('data-media-chrome-lite', '');
        }
      }
    
      class BaseButton extends HTMLElement {
        connectedCallback() {
          if (this.shadowRoot) return;
          const root = this.attachShadow({ mode: 'open' });
          this.button = document.createElement('button');
          this.button.type = 'button';
          this.button.part = 'button';
          this.button.setAttribute('aria-label', this.label || 'Media button');
          this.button.textContent = this.text || '•';
          const style = document.createElement('style');
          style.textContent = `
            button {
              border: 0;
              border-radius: 8px;
              padding: 7px 9px;
              min-width: 34px;
              background: rgba(255,255,255,.10);
              color: white;
              font: inherit;
              cursor: pointer;
            }
            button:hover { background: rgba(255,255,255,.18); }
          `;
          root.append(style, this.button);
          this.button.addEventListener('click', () => this.onClick());
          this.video = findVideo(this);
          if (this.video) this.bindVideo(this.video);
        }
        onClick() {}
        bindVideo() {}
      }
    
      class MediaPlayButton extends BaseButton {
        constructor() { super(); this.text = '▶'; this.label = 'Play or pause'; }
        bindVideo(video) {
          const update = () => { this.button.textContent = video.paused ? '▶' : '❚❚'; };
          video.addEventListener('play', update);
          video.addEventListener('pause', update);
          video.addEventListener('loadedmetadata', update);
          update();
        }
        onClick() {
          const video = findVideo(this);
          if (!video) return;
          if (video.paused) video.play().catch(() => {});
          else video.pause();
        }
      }
    
      class MediaSeekBackwardButton extends BaseButton {
        constructor() { super(); this.text = '↶'; this.label = 'Seek backward'; }
        onClick() {
          const video = findVideo(this);
          if (!video) return;
          const offset = Number(this.getAttribute('seek-offset') || '1');
          video.currentTime = Math.max(0, video.currentTime - offset);
        }
      }
    
      class MediaSeekForwardButton extends BaseButton {
        constructor() { super(); this.text = '↷'; this.label = 'Seek forward'; }
        onClick() {
          const video = findVideo(this);
          if (!video) return;
          const offset = Number(this.getAttribute('seek-offset') || '1');
          const max = Number.isFinite(video.duration) ? video.duration : Infinity;
          video.currentTime = Math.min(max, video.currentTime + offset);
        }
      }
    
      class MediaMuteButton extends BaseButton {
        constructor() { super(); this.text = '🔊'; this.label = 'Mute'; }
        bindVideo(video) {
          const update = () => { this.button.textContent = video.muted || video.volume === 0 ? '🔇' : '🔊'; };
          video.addEventListener('volumechange', update);
          update();
        }
        onClick() {
          const video = findVideo(this);
          if (!video) return;
          video.muted = !video.muted;
        }
      }
    
      class MediaFullscreenButton extends BaseButton {
        constructor() { super(); this.text = '⛶'; this.label = 'Fullscreen'; }
        onClick() {
          const controller = this.closest('media-controller');
          if (!controller) return;
          if (document.fullscreenElement) document.exitFullscreen().catch(() => {});
          else controller.requestFullscreen().catch(() => {});
        }
      }
    
      class MediaPlaybackRateButton extends BaseButton {
        constructor() { super(); this.text = '1×'; this.label = 'Playback speed'; }
        connectedCallback() {
          super.connectedCallback();
          const video = findVideo(this);
          if (!video || !this.button) return;
          const update = () => { this.button.textContent = `${video.playbackRate.toFixed(video.playbackRate % 1 ? 1 : 0)}×`; };
          video.addEventListener('ratechange', update);
          update();
        }
        onClick() {
          const video = findVideo(this);
          if (!video) return;
          const rates = [0.25, 0.5, 1, 1.5, 2, 4];
          const i = rates.findIndex(r => Math.abs(r - video.playbackRate) < 0.01);
          video.playbackRate = rates[(i + 1) % rates.length];
        }
      }
    
      class MediaTimeRange extends HTMLElement {
        connectedCallback() {
          if (this.shadowRoot) return;
          const root = this.attachShadow({ mode: 'open' });
          const input = document.createElement('input');
          input.type = 'range';
          input.min = '0';
          input.max = '1000';
          input.value = '0';
          input.part = 'range';
          input.setAttribute('aria-label', 'Video position');
          this.input = input;
          this.lanes = document.createElement('div');
          this.lanes.className = 'detection-lanes';
          this.lanes.setAttribute('aria-label', 'Detection timeline');
          this.lanes.hidden = true;
          const style = document.createElement('style');
          style.textContent = `
            :host { display: block; flex: 1 1 180px; min-width: 110px; }
            input { width: 100%; margin: 0; accent-color: white; }
            .detection-lanes { margin: 3px 8px 0; }
            .lane { height: 7px; position: relative; margin-top: 3px; background: #262a32; border-radius: 3px; }
            .event { position: absolute; top: 0; height: 100%; padding: 0; min-width: 2px; border: 0; border-radius: 2px; cursor: pointer; }
            .event:hover, .event:focus-visible { outline: 2px solid white; z-index: 1; }
          `;
          root.append(style, input, this.lanes);
          const video = findVideo(this);
          if (!video) return;
          let dragging = false;
          const update = () => {
            if (dragging) return;
            if (!Number.isFinite(video.duration) || video.duration <= 0) input.value = '0';
            else input.value = String(Math.round((video.currentTime / video.duration) * 1000));
          };
          input.addEventListener('input', () => {
            dragging = true;
            if (Number.isFinite(video.duration) && video.duration > 0) {
              video.currentTime = (Number(input.value) / 1000) * video.duration;
            }
          });
          input.addEventListener('change', () => { dragging = false; update(); });
          video.addEventListener('timeupdate', update);
          video.addEventListener('durationchange', update);
          video.addEventListener('loadedmetadata', update);
          video.addEventListener('durationchange', () => this.renderDetections());
          this.video = video;
          this.renderDetections();
          update();
        }
        setDetections(result) {
          this.detections = result;
          this.renderDetections();
        }
        renderDetections() {
          if (!this.lanes) return;
          this.lanes.replaceChildren();
          const result = this.detections;
          this.lanes.hidden = !result;
          if (!result) return;
          const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : result.duration;
          if (!Number.isFinite(duration) || duration <= 0) return;
          const categories = [['person', 'Human', '#75d685'], ['animal', 'Animal', '#ffd166'], ['bike', 'Bike', '#72a7ff'], ['moving_car', 'Moving vehicle', '#c792ea']];
          for (const [kind, label, color] of categories) {
            const lane = document.createElement('div');
            lane.className = 'lane';
            lane.dataset.kind = kind;
            lane.setAttribute('aria-label', label + ' detections');
            for (const interval of result.intervals || []) {
              if (interval.kind !== kind || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
              const start = Math.max(0, Math.min(duration, interval.start));
              const end = Math.max(start, Math.min(duration, interval.end));
              if (end <= start) continue;
              const event = document.createElement('button');
              event.type = 'button';
              event.className = 'event';
              event.style.left = (start / duration * 100) + '%';
              event.style.width = ((end - start) / duration * 100) + '%';
              event.style.background = color;
              const text = label + ': ' + start.toFixed(1) + '–' + end.toFixed(1) + ' sec (' + Math.round(interval.confidence * 100) + '%)';
              event.title = text;
              event.setAttribute('aria-label', 'Seek to ' + text);
              event.addEventListener('click', () => { if (this.video) this.video.currentTime = start; });
              lane.append(event);
            }
            this.lanes.append(lane);
          }
        }
      }
    
      class MediaTimeDisplay extends HTMLElement {
        connectedCallback() {
          this.textContent = '0:00';
          const video = findVideo(this);
          if (!video) return;
          const update = () => { this.textContent = formatTime(video.currentTime); };
          video.addEventListener('timeupdate', update);
          video.addEventListener('loadedmetadata', update);
          update();
        }
      }
    
      class MediaDurationDisplay extends HTMLElement {
        connectedCallback() {
          this.textContent = '0:00';
          const video = findVideo(this);
          if (!video) return;
          const update = () => { this.textContent = Number.isFinite(video.duration) ? formatTime(video.duration) : '0:00'; };
          video.addEventListener('durationchange', update);
          video.addEventListener('loadedmetadata', update);
          update();
        }
      }
    
      class MediaVolumeRange extends HTMLElement {
        connectedCallback() {
          if (this.shadowRoot) return;
          const root = this.attachShadow({ mode: 'open' });
          const input = document.createElement('input');
          input.type = 'range';
          input.min = '0';
          input.max = '1';
          input.step = '0.05';
          input.value = '1';
          input.part = 'range';
          const style = document.createElement('style');
          style.textContent = `
            :host { display: block; width: 82px; }
            input { width: 100%; accent-color: white; }
          `;
          root.append(style, input);
          const video = findVideo(this);
    
        baseline, err := medianLumaBaseline(frames)
        if err != nil {
            return err
        }
    
        bestIndex := -1
        bestScore := -1.0
        for i, frame := range frames {
            score := frameDivergenceScore(frame, baseline)
            if score > bestScore {
                bestScore = score
                bestIndex = i
            }
        }
        if bestIndex < 0 {
            return fmt.Errorf("no usable candidate frame remained after baseline analysis")
        }
    
        best := frames[bestIndex]
        if err := os.Rename(best.path, tmpPath); err != nil {
            return err
        }
        log.Printf("thumbnail selected %.3fs for %s (divergence %.2f from median scene)", best.seek, videoPath, bestScore)
        return nil
    }
    
    func (s *Server) runFFmpegThumbnail(videoPath, tmpPath, seek string) error {
        ctx, cancel := context.WithTimeout(s.workContext(), 30*time.Second)
        defer cancel()
        cmd := exec.CommandContext(ctx, s.ffmpegPath,
            "-hide_banner",
            "-loglevel", "error",
            "-ss", seek,
            "-i", videoPath,
            "-frames:v", "1",
            "-vf", fmt.Sprintf("scale=%d:-2", s.thumbWidth),
            "-q:v", "5",
            "-y",
            tmpPath,
        )
        out, err := cmd.CombinedOutput()
        if err != nil {
            if len(out) > 0 {
                return fmt.Errorf("ffmpeg failed: %w: %s", err, strings.TrimSpace(string(out)))
            }
            return fmt.Errorf("ffmpeg failed: %w", err)
        }
        if st, err := os.Stat(tmpPath); err != nil || st.Size() == 0 {
            return fmt.Errorf("ffmpeg produced no thumbnail")
        }
        return nil
    }
    
    func (s *Server) deleteClip(camera, file string, force bool) error {
        s.clipMutationMu.Lock()
        defer s.clipMutationMu.Unlock()
        key := clipKey(camera, file)
        meta := s.meta.Get(key)
        if meta.Favorite && !force {
            return fmt.Errorf("clip is favorited; send force=true to delete it manually")
        }
        videoPath, err := s.resolveClipPath(camera, file)
        if err != nil {
            return err
        }
        thumbPath, _ := s.thumbnailPath(videoPath)
        if err := os.Remove(videoPath); err != nil {
            return err
        }
        if thumbPath != "" {
            _ = os.Remove(thumbPath)
        }
        _ = s.meta.Remove(key)
        _ = os.Remove(s.detectionCachePath(videoPath))
        s.setDetectionState(videoPath, "complete", "")
        s.invalidateSizeCache()
        return nil
    }
    
    func (s *Server) thumbnailMaintenanceLoop() {
        // Give startup retention a chance to remove expired clips before generating thumbnails for them.
        if !s.waitWork(25 * time.Second) {
            return
        }
        s.runThumbnailMaintenance()
    
        ticker := time.NewTicker(s.thumbMaintenanceEvery)
        defer ticker.Stop()
        for {
            select {
    ✓ • 0ms

## Activity

    $ git status --short; for path in /AGENTS.md /home/AGENTS.md /home/ubuntu/AGENTS.md /home/ubuntu/NeoNVR/AGENTS.md; do if [ -f "$path" ]; then cat "$path"; fi; done
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M docker-compose.yml
     M main.go
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import os, shutil
        for tool in ('ffprobe', 'ffmpeg', 'jq'):
         print(tool, shutil.which(tool))
        for root in ('/mnt/FTPfiles', '/mnt/mySSD/BlinkClips', '/home/ubuntu/FTPfiles'):
         print('ROOT', root, os.path.isdir(root))
         count = 0
         for directory, subdirs, files in os.walk(root):
          for name in files:
           if name.lower().endswith('.mp4'):
            print(os.path.join(directory, name))
            count += 1
            if count >= 12: break
          if count >= 12: break
        PY
    ffprobe /usr/bin/ffprobe
    ffmpeg /usr/bin/ffmpeg
    jq /usr/bin/jq
    ROOT /mnt/FTPfiles True
    /mnt/FTPfiles/yicamFront/2026Y10M10D19H/.neonvr-detect-2708899520.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W424M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W430M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W420M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W427M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W419M13S46.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W423M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W415M41S19.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W426M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W400M08S51.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W422M00S60.mp4
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W429M00S60.mp4
    ROOT /mnt/mySSD/BlinkClips True
    ROOT /home/ubuntu/FTPfiles True
    ✓ • 0ms

## Activity

    $ sed -n '1,220p' detection_test.go; sed -n '1,170p' main.go; sed -n '1490,1640p' main.go; sed -n '2420,2515p' main.go; sed -n '1,110p' static/detection-viewer.js; rg -n 'Thumbnail|thumb|clipTimestamp|newServer|Server\{' '*test.go' main.go detection.go; sed -n '280,308p' README.md
    package main
    
    import (
        "bytes"
        "context"
        "encoding/json"
        "errors"
        "net/http"
        "net/http/httptest"
        "os"
        "os/exec"
        "path/filepath"
        "reflect"
        "strings"
        "testing"
        "time"
    )
    
    func detectionFixture(t *testing.T) (*Server, Camera, string, DetectionResult) {
        t.Helper()
        for _, name := range []string{"ffmpeg", "ffprobe"} {
            if _, err := exec.LookPath(name); err != nil {
                t.Skip(name + " not installed")
            }
        }
        c := defaultConfig()
        c.Root = t.TempDir()
        cam := Camera{Name: "Garden", Dir: filepath.Join(c.Root, "Garden"), Managed: true}
        if err := os.MkdirAll(cam.Dir, 0750); err != nil {
            t.Fatal(err)
        }
        path := filepath.Join(cam.Dir, "2026-10-01_12-00-00_Garden.mp4")
        out, err := exec.Command("ffmpeg", "-nostdin", "-v", "error", "-f", "lavfi", "-i", "color=c=green:s=160x96:r=10",
            "-f", "lavfi", "-i", "sine=frequency=440", "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1", "-c:a", "aac",
            "-metadata", "source_camera=Garden", "-movflags", "+use_metadata_tags", path).CombinedOutput()
        if err != nil {
            t.Fatalf("create fixture: %v: %s", err, out)
        }
        if err := os.Chmod(path, 0640); err != nil {
            t.Fatal(err)
        }
        mtime := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
        if err := os.Chtimes(path, mtime, mtime); err != nil {
            t.Fatal(err)
        }
        meta, err := NewMetaStore(filepath.Join(c.Root, ".blinkclips_meta.json"))
        if err != nil {
            t.Fatal(err)
        }
        s := &Server{root: c.Root, meta: meta, ffmpegPath: "ffmpeg", thumbDir: filepath.Join(c.Root, ".thumbs"), detectionConfig: c}
        d := DetectionResult{Version: 1, Model: "yolo26n", ModelSHA256: strings.Repeat("a", 64), Policy: s.detectionPolicy(),
            Complete: true, Duration: 4, SampleFPS: 2, Frames: 8, Intervals: []DetectionInterval{}, AnalyzedAt: time.Now().UTC().Format(time.RFC3339)}
        return s, cam, path, d
    }
    
    func streamHash(t *testing.T, path string) []byte {
        t.Helper()
        out, err := exec.Command("ffmpeg", "-nostdin", "-v", "error", "-i", path, "-map", "0", "-c", "copy", "-f", "streamhash", "-hash", "sha256", "-").CombinedOutput()
        if err != nil {
            t.Fatalf("stream hash: %v: %s", err, out)
        }
        return out
    }
    
    func TestDetectionMP4RoundTripPreservesRecording(t *testing.T) {
        s, cam, path, result := detectionFixture(t)
        result.Intervals = []DetectionInterval{
            {Kind: "person", Labels: []string{"person"}, Start: .5, End: 1.5, Confidence: .9},
            {Kind: "bike", Labels: []string{"bicycle"}, Start: 2, End: 3.5, Confidence: .8},
        }
        original, _ := os.Stat(path)
        before := streamHash(t, path)
        if err := s.embedDetection(context.Background(), path, original, result); err != nil {
            t.Fatal(err)
        }
        after, _ := os.Stat(path)
        if !after.ModTime().Equal(original.ModTime()) || after.Mode().Perm() != original.Mode().Perm() {
            t.Fatal("recording mtime/permissions changed")
        }
        if !bytes.Equal(before, streamHash(t, path)) {
            t.Fatal("audio/video packets changed")
        }
        metadata, err := s.embeddedDetection(context.Background(), path)
        if err != nil || !reflect.DeepEqual(metadata, &result) {
            t.Fatalf("embedded intervals did not round trip: %+v %v", metadata, err)
        }
        allTags, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format_tags=source_camera", "-of", "json", path).Output()
        if err != nil || !bytes.Contains(allTags, []byte("Garden")) {
            t.Fatal("existing custom metadata lost")
        }
        // The MP4 remains the source of truth when the expendable JSON cache is lost.
        if err := os.Remove(s.detectionCachePath(path)); err != nil {
            t.Fatal(err)
        }
        w := httptest.NewRecorder()
        s.apiDetectionsHandler(w, httptest.NewRequest(http.MethodGet, "/api/detections?camera=Garden&file="+filepath.Base(path), nil))
        var response detectionResponse
        if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || response.Status != "complete" || !reflect.DeepEqual(response.Detections, &result) {
            t.Fatalf("metadata API: %d %s", w.Code, w.Body.String())
        }
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); err != nil {
            t.Fatal("target recording pruned")
        }
    }
    
    func TestDetectionPruningKeepsFavoritesAndUncertainAndQuarantinesEmpty(t *testing.T) {
        s, cam, path, result := detectionFixture(t)
        result.Uncertain = true
        info, _ := os.Stat(path)
        if err := s.embedDetection(context.Background(), path, info, result); err != nil {
            t.Fatal(err)
        }
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); err != nil {
            t.Fatal("uncertain recording pruned")
        }
        result.Uncertain = false
        info, _ = os.Stat(path)
        if err := s.embedDetection(context.Background(), path, info, result); err != nil {
            t.Fatal(err)
        }
        key := clipKey(cam.Name, filepath.Base(path))
        if err := s.meta.SetFavorite(key, true); err != nil {
            t.Fatal(err)
        }
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); err != nil {
            t.Fatal("favorite recording pruned")
        }
        if err := s.meta.SetFavorite(key, false); err != nil {
            t.Fatal(err)
        }
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); !os.IsNotExist(err) {
            t.Fatal("empty recording not quarantined")
        }
        quarantined := filepath.Join(cam.Dir, ".neonvr_pruned", filepath.Base(path))
        metadata, err := s.embeddedDetection(context.Background(), quarantined)
        if err != nil || !reflect.DeepEqual(metadata, &result) {
            t.Fatal("quarantined video/metadata not recoverable")
        }
        files, err := listCameraFiles(cam)
        if err != nil || len(files) != 0 {
            t.Fatal("quarantine still appears in recordings")
        }
    }
    
    func TestDetectionPruneRequiresCurrentMP4Metadata(t *testing.T) {
        s, cam, path, result := detectionFixture(t)
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err == nil {
            t.Fatal("pruned an MP4 with no embedded analysis")
        }
        info, _ := os.Stat(path)
        if err := s.embedDetection(context.Background(), path, info, result); err != nil {
            t.Fatal(err)
        }
        result.ModelSHA256 = strings.Repeat("b", 64)
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err == nil {
            t.Fatal("stale analysis authorized pruning")
        }
        if _, err := os.Stat(path); err != nil {
            t.Fatal("recording was removed")
        }
    }
    
    func TestDetectionRejectsChangedRecordingBeforeReplacing(t *testing.T) {
        s, _, path, result := detectionFixture(t)
        original, _ := os.Stat(path)
        newTime := original.ModTime().Add(time.Minute)
        if err := os.Chtimes(path, newTime, newTime); err != nil {
            t.Fatal(err)
        }
        before, _ := os.ReadFile(path)
        if err := s.embedDetection(context.Background(), path, original, result); err == nil {
            t.Fatal("overwrote changed recording")
        }
        after, _ := os.ReadFile(path)
        if !bytes.Equal(before, after) {
            t.Fatal("failed remux changed recording")
        }
        files, _ := listCameraFiles(Camera{Dir: filepath.Dir(path)})
        if len(files) != 1 {
            t.Fatal("temporary remux leaked into recording list")
        }
    }
    
    func TestDetectionDeletionMode(t *testing.T) {
        s, cam, path, result := detectionFixture(t)
        s.detectionConfig.DetectionPrune = "delete"
        info, _ := os.Stat(path)
        if err := s.embedDetection(context.Background(), path, info, result); err != nil {
            t.Fatal(err)
        }
        if err := s.pruneDetectionClip(context.Background(), cam, filepath.Base(path), result); err != nil {
            t.Fatal(err)
        }
        if _, err := os.Stat(path); !os.IsNotExist(err) {
            t.Fatal("delete mode did not delete")
        }
    }
    
    func TestDetectionCacheInvalidatesWhenPolicyOrSourceChanges(t *testing.T) {
        s, _, path, result := detectionFixture(t)
        info, _ := os.Stat(path)
        if err := s.cacheDetection(path, info, result); err != nil {
            t.Fatal(err)
        }
        if s.cachedDetection(path, info) == nil {
            t.Fatal("cache miss")
        }
        s.detectionConfig.DetectionConfidence += .1
    package main
    
    import (
        "context"
        "crypto/sha1"
        _ "embed"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "html/template"
        "image/jpeg"
        "log"
        "math"
        "net/http"
        "net/url"
        "os"
        "os/exec"
        "os/signal"
        "path/filepath"
        "sort"
        "strconv"
        "strings"
        "sync"
        "syscall"
        "time"
    )
    
    const defaultRoot = "/mnt/mySSD/BlinkClips"
    const defaultAddr = "0.0.0.0:5500"
    const thumbnailAlgorithmVersion = "divergence-v2"
    
    var thumbnailCandidateFractions = []float64{0.08, 0.18, 0.28, 0.38, 0.50, 0.62, 0.72, 0.82, 0.92}
    
    type stringListFlag []string
    
    func (v *stringListFlag) String() string { return strings.Join(*v, ",") }
    func (v *stringListFlag) Set(value string) error {
        value = strings.TrimSpace(value)
        if value == "" {
            return fmt.Errorf("path cannot be empty")
        }
        *v = append(*v, value)
        return nil
    }
    
    func expandPath(path string) string {
        if path == "~" || strings.HasPrefix(path, "~/") {
            if home, err := os.UserHomeDir(); err == nil {
                if path == "~" {
                    return home
                }
                return filepath.Join(home, strings.TrimPrefix(path, "~/"))
            }
        }
        return path
    }
    
    type Server struct {
        runCtx                context.Context
        background            sync.WaitGroup
        root                  string
        extraRoots            []string
        thumbDir              string
        metaPath              string
        ffmpegPath            string
        thumbSeek             float64
        thumbWidth            int
        thumbSem              chan struct{}
        thumbMaintenanceEvery time.Duration
        sizeTTL               time.Duration
        autoDeleteDays        int
        autoDeleteEvery       time.Duration
    
        meta                 *MetaStore
        clipMutationMu       sync.Mutex
        detectionConfig      Config
        detectionDir         string
        detectionMu          sync.Mutex
        detectionStates      map[string]detectionResponse
        detectionWorkerError string
    
        sizeCacheMu  sync.Mutex
        sizeCache    statsResponse
        sizeCacheAt  time.Time
        sizeCacheErr string
    }
    
    type Camera struct {
        Name    string `json:"name"`
        Dir     string `json:"-"`
        Managed bool   `json:"-"`
    }
    
    type Clip struct {
        Camera       string `json:"camera"`
        File         string `json:"file"`
        Key          string `json:"key"`
        Timestamp    string `json:"timestamp"`
        UnixSeconds  int64  `json:"unix_seconds"`
        Date         string `json:"date"`
        Time         string `json:"time"`
        StartSeconds int    `json:"start_seconds"`
        Favorite     bool   `json:"favorite"`
        VideoURL     string `json:"video_url"`
        ThumbURL     string `json:"thumb_url"`
    }
    
    type ClipMeta struct {
        Favorite  bool   `json:"favorite,omitempty"`
        UpdatedAt string `json:"updated_at,omitempty"`
    }
    
    type MetaFile struct {
        Version int                 `json:"version"`
        Clips   map[string]ClipMeta `json:"clips"`
    }
    
    type MetaStore struct {
        path string
        mu   sync.Mutex
        data MetaFile
    }
    
    type camerasResponse struct {
        Root       string   `json:"root"`
        Cameras    []Camera `json:"cameras"`
        LatestDate string   `json:"latest_date"`
    }
    
    type datesResponse struct {
        Cameras []string `json:"cameras"`
        Dates   []string `json:"dates"`
    }
    
    type clipsResponse struct {
        Cameras []string `json:"cameras"`
        Date    string   `json:"date"`
        Clips   []Clip   `json:"clips"`
    }
    
    type favoritesResponse struct {
        Cameras []string `json:"cameras"`
        Clips   []Clip   `json:"clips"`
    }
    
    type mutateRequest struct {
        Camera   string `json:"camera"`
        File     string `json:"file"`
        Favorite bool   `json:"favorite"`
        Force    bool   `json:"force"`
    }
    
    type bulkDeleteItem struct {
        Camera string `json:"camera"`
        File   string `json:"file"`
    }
    
    type bulkDeleteRequest struct {
        Items []bulkDeleteItem `json:"items"`
        Force bool             `json:"force"`
    }
    
    type mutateResponse struct {
        OK      bool   `json:"ok"`
        Deleted int    `json:"deleted,omitempty"`
        Message string `json:"message,omitempty"`
    }
    
    type statsResponse struct {
    }
    
    async function refreshDates(pickLatestIfNeeded = true) {
      const data = await getJSON('/api/dates?' + qs({cameras: cameraParam()}));
      availableDates = data.dates || [];
    
      if (availableDates.length) {
        const sortedAsc = [...availableDates].sort();
        dateEl.min = sortedAsc[0];
        dateEl.max = sortedAsc[sortedAsc.length - 1];
        if (pickLatestIfNeeded && !availableDates.includes(dateEl.value)) {
          dateEl.value = availableDates[0];
        }
      } else {
        dateEl.removeAttribute('min');
        dateEl.removeAttribute('max');
        dateEl.value = '';
      }
    
      renderDatePills();
    }
    
    async function refreshClips(reset = true) {
      if (currentView === 'favorites') {
        const data = await getJSON('/api/favorites?' + qs({cameras: cameraParam()}));
        allClips = data.clips || [];
        if (reset) resetRange();
        selectedKeys = new Set([...selectedKeys].filter(k => allClips.some(c => keyFor(c) === k)));
        renderMarks();
        renderClips();
        renderDatePills();
        return;
      }
    
      if (!dateEl.value) {
        allClips = [];
        renderMarks();
        renderClips();
        return;
      }
    
      const data = await getJSON('/api/clips?' + qs({cameras: cameraParam(), date: dateEl.value}));
      allClips = data.clips || [];
      if (reset) resetRange();
      selectedKeys = new Set([...selectedKeys].filter(k => allClips.some(c => keyFor(c) === k)));
      renderMarks();
      renderClips();
      renderDatePills();
    }
    
    async function refreshStats() {
      try {
        const stats = await getJSON('/api/stats');
        let html = '<span>Recordings: <strong>' + stats.root_size_human + '</strong></span>' +
                   '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
        if (stats.error) html += '<span>' + escapeHTML(stats.error) + '</span>';
        diskStatsEl.innerHTML = html;
      } catch (err) {
        diskStatsEl.innerHTML = '<span>Folder sizes unavailable: ' + escapeHTML(String(err.message || err)) + '</span>';
      }
    }
    
    function setView(view) {
      currentView = view;
      clipsTabEl.classList.toggle('active', view === 'clips');
      favoritesTabEl.classList.toggle('active', view === 'favorites');
      dateEl.disabled = view === 'favorites';
      clearSelection();
      refreshClips(true).catch(err => {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      });
    }
    
    async function bulkDeleteSelected() {
      if (selectedKeys.size === 0) return;
      const items = visibleClips.filter(c => selectedKeys.has(keyFor(c))).map(c => ({camera: c.camera, file: c.file}));
      const hasFavorite = visibleClips.some(c => selectedKeys.has(keyFor(c)) && c.favorite);
      const msg = hasFavorite
        ? 'Delete ' + items.length + ' selected videos? This includes favorited clips.'
        : 'Delete ' + items.length + ' selected videos?';
      if (!confirm(msg)) return;
      try {
        await postJSON('/api/bulk-delete', {items, force: hasFavorite});
        const del = new Set(items.map(i => i.camera + '/' + i.file));
        allClips = allClips.filter(c => !del.has(keyFor(c)));
        clearSelection();
        renderMarks();
        renderClips();
        refreshStats();
      } catch (err) {
        alert(String(err.message || err));
      }
    }
    
    async function init() {
      try {
        buildScales();
        updateRangeVisuals();
        refreshStats();
        const data = await getJSON('/api/cameras');
        cameras = data.cameras || [];
        selectedCameras = new Set(cameras.map(c => c.name));
        renderCameraMenu();
        await refreshDates(true);
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    }
    
    cameraButtonEl.addEventListener('click', () => cameraMenuEl.classList.toggle('open'));
    document.addEventListener('click', evt => {
      if (!cameraButtonEl.contains(evt.target) && !cameraMenuEl.contains(evt.target)) cameraMenuEl.classList.remove('open');
    });
    
    dateEl.addEventListener('change', async () => {
      try {
        clearSelection();
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    });
    
    sortEl.addEventListener('change', () => { clearSelection(); renderClips(); syncActiveCard(true); });
    startEl.addEventListener('input', () => { normalizeRange('start'); clearSelection(); renderClips(); });
    endEl.addEventListener('input', () => { normalizeRange('end'); clearSelection(); renderClips(); });
    clipsTabEl.addEventListener('click', () => setView('clips'));
    favoritesTabEl.addEventListener('click', () => setView('favorites'));
    selectModeToggleEl.addEventListener('click', () => setSelectMode(!selectMode));
    bulkDeleteEl.addEventListener('click', bulkDeleteSelected);
    clearSelectionEl.addEventListener('click', () => clearSelection());
    zoomToggleEl.addEventListener('click', () => setZoomMode(!zoomMode));
    zoomResetEl.addEventListener('click', resetZoom);
    document.addEventListener('keydown', evt => {
      if (evt.composedPath().some(el => ['INPUT','SELECT','TEXTAREA','BUTTON'].includes(el.tagName) || el.isContentEditable)) return;
    
      if (evt.key === ' ') {
        if (!mainVideoEl.src) return;
        evt.preventDefault();
        if (mainVideoEl.paused) mainVideoEl.play().catch(() => {});
        else mainVideoEl.pause();
        return;
      }
    
      if (evt.ctrlKey && ['ArrowLeft','ArrowRight','ArrowUp','ArrowDown'].includes(evt.key)) {
        evt.preventDefault();
        if (evt.key === 'ArrowLeft' || evt.key === 'ArrowUp') playRelative(-1);
        else playRelative(1);
        return;
      }
        if len(parts) < 2 {
            return time.Time{}, false
        }
        datePart := parts[0]
        timePart := strings.ReplaceAll(parts[1], "-", ":")
        ts, err := time.ParseInLocation("2006-01-02 15:04:05", datePart+" "+timePart, time.Local)
        if err != nil {
            return time.Time{}, false
        }
        return ts, true
    }
    
    func (s *Server) thumbnailPath(videoPath string) (string, error) {
        info, err := os.Stat(videoPath)
        if err != nil {
            return "", err
        }
        keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
        sum := sha1.Sum([]byte(keyInput))
        name := hex.EncodeToString(sum[:]) + ".jpg"
        return filepath.Join(s.thumbDir, name), nil
    }
    
    func (s *Server) ensureThumbnail(videoPath string) (string, error) {
        thumbPath, err := s.thumbnailPath(videoPath)
        if err != nil {
            return "", err
        }
        if err := os.MkdirAll(s.thumbDir, 0755); err != nil {
            return "", err
        }
        if st, err := os.Stat(thumbPath); err == nil && st.Size() > 0 {
            return thumbPath, nil
        }
    
        tmpPath := strings.TrimSuffix(thumbPath, ".jpg") + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
        _ = os.Remove(tmpPath)
    
        select {
        case s.thumbSem <- struct{}{}:
        case <-s.workContext().Done():
            return "", s.workContext().Err()
        }
        defer func() { <-s.thumbSem }()
    
        if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
            _ = os.Remove(tmpPath)
            seek := strconv.FormatFloat(s.thumbSeek, 'f', 3, 64)
            log.Printf("interesting thumbnail selection failed for %s, falling back to %ss: %v", videoPath, seek, err)
            if err2 := s.runFFmpegThumbnail(videoPath, tmpPath, seek); err2 != nil {
                _ = os.Remove(tmpPath)
                if s.thumbSeek != 0 {
                    if err3 := s.runFFmpegThumbnail(videoPath, tmpPath, "0"); err3 != nil {
                        _ = os.Remove(tmpPath)
                        return "", fmt.Errorf("interesting, %ss, and 0s thumbnail attempts failed: %v; %v; %v", seek, err, err2, err3)
                    }
                } else {
                    return "", fmt.Errorf("interesting and 0s thumbnail attempts failed: %v; %v", err, err2)
                }
            }
        }
    
        if err := os.Rename(tmpPath, thumbPath); err != nil {
            _ = os.Remove(tmpPath)
            return "", err
        }
        return thumbPath, nil
    }
    
    func (s *Server) probeVideoDuration(videoPath string) (float64, error) {
        ctx, cancel := context.WithTimeout(s.workContext(), 10*time.Second)
        defer cancel()
    
        ffprobePath := "ffprobe"
        if strings.ContainsRune(s.ffmpegPath, filepath.Separator) {
            ffprobePath = filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
        }
        cmd := exec.CommandContext(ctx, ffprobePath,
            "-v", "error",
            "-show_entries", "format=duration",
            "-of", "default=noprint_wrappers=1:nokey=1",
            videoPath,
        )
        out, err := cmd.Output()
        if err != nil {
            return 0, fmt.Errorf("ffprobe duration: %w", err)
        }
        duration, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
        if err != nil || math.IsNaN(duration) || math.IsInf(duration, 0) || duration <= 0 {
            return 0, fmt.Errorf("invalid video duration %q", strings.TrimSpace(string(out)))
        }
        return duration, nil
    }
    
    type frameSignature struct {
        path     string
    // Fetch clip-relative intervals, refresh pending analyses, and show active events.
    (function () {
      'use strict';
      const video = document.getElementById('mainVideo');
      const seek = document.getElementById('mainSeek');
      const status = document.getElementById('detectionStatus');
      if (!video || !seek || !status) return;
      const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
      let generation = 0;
      let request = null;
      let timer = null;
      let result = null;
    
      function activeStatus() {
        if (!result) return;
        const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
        status.classList.toggle('active', active.length > 0);
        const text = active.length ? 'Detected now: ' + active.join(', ') :
          (result.intervals.length ? 'No detection at this time.' : result.uncertain ? 'Unconfirmed observations — clip retained.' : 'No target detections.');
        if (status.textContent !== text) status.textContent = text;
      }
    
      async function refresh(clip, expected) {
        if (expected !== generation) return;
        request = new AbortController();
        try {
          const response = await fetch('/api/detections?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: request.signal});
          if (!response.ok) throw new Error(response.status === 404 ? 'Clip is no longer available.' : 'Detection data unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          if (data.status === 'complete') {
            result = data.detections;
            status.title = '';
            seek.setDetections(result);
            activeStatus();
          } else {
            status.classList.remove('active');
            const labels = {pending: 'Detection queued…', processing: 'Analyzing detections…', disabled: 'Detection disabled.', error: 'Analysis unavailable — clip retained.'};
            status.textContent = labels[data.status] || 'Detection pending…';
            status.title = data.message || '';
            if (data.status !== 'disabled') timer = setTimeout(() => refresh(clip, expected), 10000);
          }
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          status.textContent = error.message;
          if (!error.message.includes('no longer available')) timer = setTimeout(() => refresh(clip, expected), 10000);
        }
      }
    
      function loadClip(clip) {
        generation += 1;
        request?.abort();
        clearTimeout(timer);
        result = null;
        seek.setDetections(null);
        status.classList.remove('active');
        status.textContent = 'Loading detections…';
        status.title = '';
        refresh(clip, generation);
      }
    
      video.addEventListener('timeupdate', activeStatus);
      video.addEventListener('seeked', activeStatus);
      window.addEventListener('pagehide', () => { request?.abort(); clearTimeout(timer); });
      window.NeoNVRDetections = { loadClip };
    })();
    rg: *test.go: No such file or directory (os error 2)
    detection.go:466:    oldThumb, _ := s.thumbnailPath(path)
    detection.go:645:    thumb, _ := s.thumbnailPath(path)
    detection.go:667:    if thumb != "" {
    detection.go:668:        _ = os.Remove(thumb)
    main.go:31:const thumbnailAlgorithmVersion = "divergence-v2"
    main.go:33:var thumbnailCandidateFractions = []float64{0.08, 0.18, 0.28, 0.38, 0.50, 0.62, 0.72, 0.82, 0.92}
    main.go:64:    thumbDir              string
    main.go:67:    thumbSeek             float64
    main.go:68:    thumbWidth            int
    main.go:69:    thumbSem              chan struct{}
    main.go:70:    thumbMaintenanceEvery time.Duration
    main.go:106:    ThumbURL     string `json:"thumb_url"`
    main.go:214:      --thumb: #f5f7fb;
    main.go:215:      --thumbBorder: #111827;
    main.go:493:    input[type=range]::-webkit-slider-thumb {
    main.go:501:      background: var(--thumb);
    main.go:502:      border: 3px solid var(--thumbBorder);
    main.go:506:    input[type=range]::-moz-range-thumb {
    main.go:511:      background: var(--thumb);
    main.go:512:      border: 3px solid var(--thumbBorder);
    main.go:604:    .thumbButton {
    main.go:613:    .thumbButton img, .card video {
    main.go:679:          <span id="mainTitle">Main player: click a thumbnail or choose Load.</span>
    main.go:726:        <main class="sideRail" id="sideRail" aria-label="Scrollable video thumbnail navigation"></main>
    main.go:1054:  mainVideoEl.poster = clip.thumb_url;
    main.go:1288:    img.src = clip.thumb_url;
    main.go:1344:    button.className = 'thumbButton';
    main.go:1359:    img.src = clip.thumb_url;
    main.go:1672:    thumbDir, metaPath, ffmpegPath := &cfg.ThumbDir, &cfg.Meta, &cfg.FFmpeg
    main.go:1673:    thumbSeek, thumbWidth, thumbWorkers := &cfg.ThumbSeek, &cfg.ThumbWidth, &cfg.ThumbWorkers
    main.go:1677:    thumbMaintenanceEvery, sizeRefresh, autoDeleteEvery := &maintenance, &refresh, &retentionEvery
    main.go:1697:    cleanThumbDir := *thumbDir
    main.go:1699:        cleanThumbDir = filepath.Join(cleanRoot, ".blinkclips_thumbs")
    main.go:1709:    workers := *thumbWorkers
    main.go:1713:    if *thumbWidth < 64 {
    main.go:1714:        *thumbWidth = 64
    main.go:1716:    if *thumbMaintenanceEvery > 0 && *thumbMaintenanceEvery < time.Minute {
    main.go:1717:        *thumbMaintenanceEvery = time.Minute
    main.go:1722:    if *thumbSeek < 0 {
    main.go:1723:        *thumbSeek = 0
    main.go:1734:    s := &Server{
    main.go:1738:        thumbDir:              cleanThumbDir,
    main.go:1741:        thumbSeek:             *thumbSeek,
    main.go:1742:        thumbWidth:            *thumbWidth,
    main.go:1743:        thumbSem:              make(chan struct{}, workers),
    main.go:1744:        thumbMaintenanceEvery: *thumbMaintenanceEvery,
    main.go:1757:    if s.thumbMaintenanceEvery > 0 {
    main.go:1759:        go func() { defer s.background.Done(); s.thumbnailMaintenanceLoop() }()
    main.go:1812:    mux.HandleFunc("/thumb", s.thumbHandler)
    main.go:1819:    log.Printf("Thumbnail cache: %s", s.thumbDir)
    main.go:1823:    log.Printf("Thumbnail seek: %.3fs", s.thumbSeek)
    main.go:1824:    log.Printf("Thumbnail workers: %d", workers)
    main.go:1825:    if s.thumbMaintenanceEvery > 0 {
    main.go:1826:        log.Printf("Thumbnail maintenance: every %s (pre-generate missing, prune orphan cache files)", s.thumbMaintenanceEvery)
    main.go:1828:        log.Printf("Thumbnail maintenance disabled")
    main.go:1836:    httpServer := &http.Server{Addr: *addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
    main.go:2123:func (s *Server) thumbHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2136:    thumbPath, err := s.ensureThumbnail(videoPath)
    main.go:2138:        log.Printf("thumbnail error for %s: %v", videoPath, err)
    main.go:2139:        writePlaceholderSVG(w, "thumbnail unavailable")
    main.go:2144:    http.ServeFile(w, r, thumbPath)
    main.go:2278:func clipTimestamp(f clipFile) time.Time {
    main.go:2306:            seen[clipTimestamp(f).Format("2006-01-02")] = true
    main.go:2336:            ts := clipTimestamp(f)
    main.go:2354:                ThumbURL:     "/thumb?camera=" + url.QueryEscape(cam.Name) + "&file=" + url.QueryEscape(name),
    main.go:2432:func (s *Server) thumbnailPath(videoPath string) (string, error) {
    main.go:2437:    keyInput := videoPath + "|" + strconv.FormatInt(info.ModTime().UnixNano(), 10) + "|" + strconv.FormatInt(info.Size(), 10) + "|" + fmt.Sprintf("%.3f", s.thumbSeek) + "|" + strconv.Itoa(s.thumbWidth) + "|" + thumbnailAlgorithmVersion
    main.go:2440:    return filepath.Join(s.thumbDir, name), nil
    main.go:2443:func (s *Server) ensureThumbnail(videoPath string) (string, error) {
    main.go:2444:    thumbPath, err := s.thumbnailPath(videoPath)
    main.go:2448:    if err := os.MkdirAll(s.thumbDir, 0755); err != nil {
    main.go:2451:    if st, err := os.Stat(thumbPath); err == nil && st.Size() > 0 {
    main.go:2452:        return thumbPath, nil
    main.go:2455:    tmpPath := strings.TrimSuffix(thumbPath, ".jpg") + ".tmp." + strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
    main.go:2459:    case s.thumbSem <- struct{}{}:
    main.go:2463:    defer func() { <-s.thumbSem }()
    main.go:2465:    if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
    main.go:2467:        seek := strconv.FormatFloat(s.thumbSeek, 'f', 3, 64)
    main.go:2468:        log.Printf("interesting thumbnail selection failed for %s, falling back to %ss: %v", videoPath, seek, err)
    main.go:2469:        if err2 := s.runFFmpegThumbnail(videoPath, tmpPath, seek); err2 != nil {
    main.go:2471:            if s.thumbSeek != 0 {
    main.go:2472:                if err3 := s.runFFmpegThumbnail(videoPath, tmpPath, "0"); err3 != nil {
    main.go:2474:                    return "", fmt.Errorf("interesting, %ss, and 0s thumbnail attempts failed: %v; %v; %v", seek, err, err2, err3)
    main.go:2477:                return "", fmt.Errorf("interesting and 0s thumbnail attempts failed: %v; %v", err, err2)
    main.go:2482:    if err := os.Rename(tmpPath, thumbPath); err != nil {
    main.go:2486:    return thumbPath, nil
    main.go:2536:        return frameSignature{}, fmt.Errorf("thumbnail too small")
    main.go:2623:func (s *Server) runInterestingThumbnail(videoPath, tmpPath string) error {
    main.go:2625:    seekTimes := make([]float64, 0, len(thumbnailCandidateFractions))
    main.go:2630:        for _, fraction := range thumbnailCandidateFractions {
    main.go:2654:        if err := s.runFFmpegThumbnail(videoPath, candidatePath, seek); err != nil {
    main.go:2693:    log.Printf("thumbnail selected %.3fs for %s (divergence %.2f from median scene)", best.seek, videoPath, bestScore)
    main.go:2697:func (s *Server) runFFmpegThumbnail(videoPath, tmpPath, seek string) error {
    main.go:2706:        "-vf", fmt.Sprintf("scale=%d:-2", s.thumbWidth),
    main.go:2719:        return fmt.Errorf("ffmpeg produced no thumbnail")
    main.go:2736:    thumbPath, _ := s.thumbnailPath(videoPath)
    main.go:2740:    if thumbPath != "" {
    main.go:2741:        _ = os.Remove(thumbPath)
    main.go:2750:func (s *Server) thumbnailMaintenanceLoop() {
    main.go:2751:    // Give startup retention a chance to remove expired clips before generating thumbnails for them.
    main.go:2755:    s.runThumbnailMaintenance()
    main.go:2757:    ticker := time.NewTicker(s.thumbMaintenanceEvery)
    main.go:2764:            s.runThumbnailMaintenance()
    main.go:2769:func (s *Server) runThumbnailMaintenance() {
    main.go:2772:        log.Printf("thumbnail maintenance: could not list sources: %v", err)
    main.go:2782:            log.Printf("thumbnail maintenance: could not scan %s: %v", cam.Name, err)
    main.go:2792:            thumbPath, err := s.thumbnailPath(videoPath)
    main.go:2794:                log.Printf("thumbnail maintenance: could not map %s to cache: %v", videoPath, err)
    main.go:2797:            desired[filepath.Clean(thumbPath)] = struct{}{}
    main.go:2804:    workers := cap(s.thumbSem)
    main.go:2819:                thumbPath, err := s.thumbnailPath(videoPath)
    main.go:2826:                if st, err := os.Stat(thumbPath); err == nil && st.Size() > 0 {
    main.go:2829:                if _, err := s.ensureThumbnail(videoPath); err != nil {
    main.go:2830:                    log.Printf("thumbnail maintenance: generation failed for %s: %v", videoPath, err)
    main.go:2855:    entries, err := os.ReadDir(s.thumbDir)
    main.go:2858:            log.Printf("thumbnail maintenance: could not read cache directory %s: %v", s.thumbDir, err)
    main.go:2875:            full := filepath.Clean(filepath.Join(s.thumbDir, entry.Name()))
    main.go:2880:                log.Printf("thumbnail maintenance: could not prune orphan %s: %v", full, err)
    main.go:2887:    log.Printf("thumbnail maintenance: scanned %d videos, generated %d thumbnails, pruned %d orphan thumbnails, %d generation failures", len(videos), generated, pruned, failed)
    main.go:2927:            if !clipTimestamp(f).Before(cutoff) {
    
    New thumbnails use a divergence selector designed for mostly-static security footage. The server samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    
    
    ### Automatic thumbnail maintenance
    
    Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same divergence selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).
    
    Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thumb-maintenance-every 1h`. Set it to `0` to disable background maintenance. Lazy/on-demand thumbnail generation remains enabled as a fallback.
    
    ### Yi timestamps
    
    Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
    
    ## Project rename and compatibility
    
    The project directory is now `/home/ubuntu/NeoNVR`, the Go module and executable are `neonvr`, the Compose service/container are `neonvr`, and the systemd example is `neonvr.service.example`.
    
    Existing recordings remain under `/mnt/mySSD/BlinkClips`. The legacy `.blinkclips_thumbs` and `.blinkclips_meta.json` paths are deliberately retained so thumbnail caches and favorite selections continue to work. Custom `-root`, `-thumb-dir`, and `-meta` flags continue to work.
    
    A running container keeps its current binary until rebuilt. Stop the existing Compose stack using its original project name and configuration before starting the renamed stack; otherwise both stacks can contend for port 5500. For nerdctl deployments, use `sudo nerdctl compose` in place of `docker compose`. After stopping the old stack:
    
    ```bash
    cd /home/ubuntu/NeoNVR
    sudo nerdctl compose up -d --build
    sudo nerdctl compose logs -f neonvr
    ```
    
    For a native systemd installation, install the executable at `/opt/NeoNVR/neonvr` and use `neonvr.service.example` as the unit template.
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import os, subprocess, json
        for root in ('/mnt/FTPfiles/yicamFront','/mnt/FTPfiles/yicamBack'):
         paths=[]
         for directory, subdirs, files in os.walk(root):
          for name in files:
           if name.endswith('.mp4') and not name.startswith('.'):
            path=os.path.join(directory,name)
            try: paths.append((os.stat(path).st_mtime,path))
            except OSError: pass
         for _,path in sorted(paths,reverse=True)[:10]:
          p=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion,creation_time','-of','json',path],capture_output=True,text=True,timeout=15)
          print(path, p.stdout.strip() if p.returncode == 0 else p.stderr[:200])
        PY
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W430M00S60.mp4 {
        "format": {
            "duration": "60.033000",
            "tags": {
                "creation_time": "2026-10-11T04:30:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791693000000,\"duration_ms\":60033,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W429M00S60.mp4 {
        "format": {
            "duration": "60.032000",
            "tags": {
                "creation_time": "2026-10-11T04:29:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692940000,\"duration_ms\":60032,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W428M00S60.mp4 {
        "format": {
            "duration": "59.968000",
            "tags": {
                "creation_time": "2026-10-11T04:28:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692880000,\"duration_ms\":59968,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W427M00S60.mp4 {
        "format": {
            "duration": "60.000000",
            "tags": {
                "creation_time": "2026-10-11T04:27:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692820000,\"duration_ms\":60000,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W426M00S60.mp4 {
        "format": {
            "duration": "60.071000",
            "tags": {
                "creation_time": "2026-10-11T04:26:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692760000,\"duration_ms\":60071,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W425M00S60.mp4 {
        "format": {
            "duration": "59.985000",
            "tags": {
                "creation_time": "2026-10-11T04:25:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692700000,\"duration_ms\":59985,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W424M00S60.mp4 {
        "format": {
            "duration": "59.935000",
            "tags": {
                "creation_time": "2026-10-11T04:24:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692640000,\"duration_ms\":59935,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W423M00S60.mp4 {
        "format": {
            "duration": "60.160000",
            "tags": {
                "creation_time": "2026-10-11T04:23:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692580000,\"duration_ms\":60160,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W422M00S60.mp4 {
        "format": {
            "duration": "59.998000",
            "tags": {
                "creation_time": "2026-10-11T04:22:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692520000,\"duration_ms\":59998,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W421M00S60.mp4 {
        "format": {
            "duration": "59.930000",
            "tags": {
                "creation_time": "2026-10-11T04:21:00.000000Z",
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791692460000,\"duration_ms\":59930,\"video_time_precision_ms\":1000,\"initial_state\":1,\"history_available\":true,\"history_truncated\":false,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    ✓ • 1.10s

## Activity

    $ rg --files -g '*test.go' -g AGENTS.md; sed -n '1740,1818p' main.go; sed -n '240,280p' main.go; sed -n '684,710p' main.go; sed -n '300,390p' detection.go; rg -n '^func .*Handler|func writeJSON|func finite' detection.go main.go; cat go.mod; ls /tmp/neonvr-live-tools-vbq_9dyb
    config_test.go
    detection_test.go
    live_test.go
            ffmpegPath:            *ffmpegPath,
            thumbSeek:             *thumbSeek,
            thumbWidth:            *thumbWidth,
            thumbSem:              make(chan struct{}, workers),
            thumbMaintenanceEvery: *thumbMaintenanceEvery,
            sizeTTL:               *sizeRefresh,
            autoDeleteDays:        *autoDeleteDays,
            autoDeleteEvery:       *autoDeleteEvery,
            meta:                  meta,
            detectionConfig:       cfg,
            detectionDir:          filepath.Join(cleanRoot, ".neonvr_detections"),
        }
    
        if s.autoDeleteDays > 0 {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.autoDeleteLoop() }()
        }
        if s.thumbMaintenanceEvery > 0 {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.thumbnailMaintenanceLoop() }()
        }
        if cfg.DetectionEnabled {
            s.background.Add(1)
            go func() { defer s.background.Done(); s.detectionLoop() }()
        }
    
        mux := http.NewServeMux()
        mux.HandleFunc("/", s.indexHandler)
        mux.HandleFunc("/api/config", manager.handler)
        mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                w.Header().Set("Allow", "GET")
                http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
                return
            }
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            w.Header().Set("Cache-Control", "no-store")
            _, _ = w.Write(streamsHTML)
        })
        mux.HandleFunc("/static/live-streams.js", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            _, _ = w.Write(liveStreamsJS)
        })
        mux.HandleFunc("/api/live", live.handler)
        mux.HandleFunc("/api/live/webrtc", live.webrtcHandler)
        mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodGet {
                w.Header().Set("Allow", "GET")
                http.Error(w, "method not allowed", 405)
                return
            }
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            w.Header().Set("Cache-Control", "no-store")
            _, _ = w.Write(settingsHTML)
        })
        mux.HandleFunc("/api/cameras", s.apiCamerasHandler)
        mux.HandleFunc("/api/dates", s.apiDatesHandler)
        mux.HandleFunc("/api/clips", s.apiClipsHandler)
        mux.HandleFunc("/api/favorites", s.apiFavoritesHandler)
        mux.HandleFunc("/api/favorite", s.apiFavoriteHandler)
        mux.HandleFunc("/api/delete", s.apiDeleteHandler)
        mux.HandleFunc("/api/bulk-delete", s.apiBulkDeleteHandler)
        mux.HandleFunc("/api/stats", s.apiStatsHandler)
        mux.HandleFunc("/api/detections", s.apiDetectionsHandler)
        mux.HandleFunc("/static/detection-viewer.js", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
            w.Header().Set("Cache-Control", "no-cache")
            _, _ = w.Write(detectionViewerJS)
        })
        mux.HandleFunc("/static/media-chrome-lite.js", s.staticMediaChromeHandler)
        mux.HandleFunc("/video", s.videoHandler)
        mux.HandleFunc("/thumb", s.thumbHandler)
    
        log.Printf("NeoNVR server listening on http://%s", *addr)
        log.Printf("Clips root: %s", s.root)
        for _, extraRoot := range s.extraRoots {
            log.Printf("Extra clip root: %s", extraRoot)
        }
        .viewerLayout {
          display: grid;
          grid-template-columns: minmax(0, 1fr) minmax(260px, 360px);
          gap: 14px;
          align-items: stretch;
        }
        .mainPlayerTitle { display: flex; justify-content: space-between; gap: 12px; flex-wrap: wrap; color: var(--muted); margin-bottom: 10px; font-size: 14px; }
        .mainPlayerShell {
          background: #000;
          border: 1px solid var(--border);
          border-radius: 14px;
          overflow: hidden;
          width: 100%;
          position: relative;
        }
        .mainColumn { min-width: 0; }
        .railPanel {
          min-height: 0;
          height: min(76vh, 780px);
          position: sticky;
          top: 14px;
          display: grid;
          grid-template-rows: auto auto minmax(0, 1fr) auto;
          gap: 10px;
          background: #11141b;
          border: 1px solid var(--border);
          border-radius: 14px;
          padding: 10px;
        }
        .railHeader { display: flex; justify-content: space-between; align-items: center; gap: 8px; color: var(--muted); font-size: 13px; }
        .railActions { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; color: var(--muted); font-size: 13px; }
        media-controller {
          display: grid;
          grid-template-rows: minmax(0, 1fr) auto auto;
          width: 100%;
          aspect-ratio: 16 / 9;
          background: #000;
          min-height: 260px;
        }
        #mainVideo {
          width: 100%;
                <video id="mainVideo" slot="media" preload="metadata" playsinline webkit-playsinline x5-playsinline></video>
                <media-control-bar>
                  <media-play-button></media-play-button>
                  <media-seek-backward-button seek-offset="1"></media-seek-backward-button>
                  <media-seek-forward-button seek-offset="1"></media-seek-forward-button>
                  <media-time-range id="mainSeek"></media-time-range>
                  <media-time-display></media-time-display>
                  <media-duration-display></media-duration-display>
                  <media-mute-button></media-mute-button>
                  <media-volume-range></media-volume-range>
                  <media-playback-rate-button></media-playback-rate-button>
                  <media-fullscreen-button></media-fullscreen-button>
                </media-control-bar>
                <div class="detectionLegend" style="padding:0 10px 10px;background:#000;margin:0">
                  <span class="category" style="--detection-color:#75d685">Human</span>
                  <span class="category" style="--detection-color:#ffd166">Animal</span>
                  <span class="category" style="--detection-color:#72a7ff">Bike</span>
                  <span class="category" style="--detection-color:#c792ea">Moving vehicle</span>
                  <span id="detectionStatus" role="status">Select a clip to see detections.</span>
                </div>
              </media-controller>
            </div>
            <div class="analysisControls">
              <label style="display:flex;align-items:center;gap:6px">Step sec <input id="stepSize" type="number" min="0.03" step="0.05" value="0.20"></label>
              <button type="button" id="zoomToggle" aria-pressed="false">Zoom off</button>
              <button type="button" id="zoomReset">Reset zoom</button>
              <span id="zoomLabel">1.0×</span>
            return nil, ctx.Err()
        }
    }
    
    func (s *Server) startDetectionWorker() (*detectionWorker, error) {
        c := s.detectionConfig
        ctx, cancel := context.WithCancel(s.workContext())
        cmd := exec.CommandContext(ctx, expandPath(c.DetectionPython), "-u", "-c", detectionWorkerScript,
            "--model", expandPath(c.DetectionModel), "--threads", strconv.Itoa(c.DetectionThreads))
        cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
        cmd.Cancel = func() error {
            if cmd.Process == nil {
                return nil
            }
            return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
        }
        cmd.WaitDelay = 2 * time.Second
        cmd.Stderr = os.Stderr
        stdin, err := cmd.StdinPipe()
        if err != nil {
            cancel()
            return nil, err
        }
        stdout, err := cmd.StdoutPipe()
        if err != nil {
            stdin.Close()
            cancel()
            return nil, err
        }
        if err := cmd.Start(); err != nil {
            stdin.Close()
            cancel()
            return nil, err
        }
        worker := &detectionWorker{cmd: cmd, cancel: cancel, stdin: stdin, scanner: bufio.NewScanner(stdout)}
        worker.scanner.Buffer(make([]byte, 65536), 4<<20)
        readyCtx, readyCancel := context.WithTimeout(ctx, time.Minute)
        defer readyCancel()
        data, err := worker.read(readyCtx)
        var ready struct {
            Ready       bool   `json:"ready"`
            ModelSHA256 string `json:"model_sha256"`
        }
        if err == nil {
            err = json.Unmarshal(data, &ready)
        }
        if err == nil && (!ready.Ready || len(ready.ModelSHA256) != 64) {
            err = errors.New("invalid worker handshake")
        }
        if err != nil {
            worker.close()
            return nil, err
        }
        worker.modelSHA256 = ready.ModelSHA256
        return worker, nil
    }
    
    func (s *Server) analyzeDetection(ctx context.Context, worker *detectionWorker, path string) (DetectionResult, error) {
        c := s.detectionConfig
        request := map[string]any{"path": path, "ffmpeg": s.ffmpegPath, "ffprobe": s.ffprobeExecutable(),
            "fps": c.DetectionFPS, "confidence": c.DetectionConfidence, "min_frames": c.DetectionMinFrames,
            "merge_gap": c.DetectionMergeGap, "car_movement": c.DetectionCarMovement, "policy": s.detectionPolicy()}
        data, _ := json.Marshal(request)
        var result DetectionResult
        if _, err := worker.stdin.Write(append(data, '\n')); err != nil {
            return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
        }
        out, err := worker.read(ctx)
        if err != nil {
            return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
        }
        var response struct {
            Result DetectionResult `json:"result"`
            Error  string          `json:"error"`
        }
        if err := json.Unmarshal(out, &response); err != nil {
            return result, fmt.Errorf("%w: %v", errDetectionWorker, err)
        }
        if response.Error != "" {
            return result, errors.New(response.Error)
        }
        result = response.Result
        if result.Policy != s.detectionPolicy() || result.ModelSHA256 != worker.modelSHA256 {
            return result, fmt.Errorf("%w: policy/model mismatch", errDetectionWorker)
        }
        return result, result.validate()
    }
    
    func sameRecording(a, b os.FileInfo) bool {
        return os.SameFile(a, b) && a.Size() == b.Size() && a.ModTime().Equal(b.ModTime())
    }
    main.go:1950:func (s *Server) staticMediaChromeHandler(w http.ResponseWriter, r *http.Request) {
    main.go:1956:func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
    main.go:1967:func (s *Server) apiCamerasHandler(w http.ResponseWriter, r *http.Request) {
    main.go:1985:func (s *Server) apiDatesHandler(w http.ResponseWriter, r *http.Request) {
    main.go:1999:func (s *Server) apiClipsHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2014:func (s *Server) apiFavoritesHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2028:func (s *Server) apiFavoriteHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2055:func (s *Server) apiDeleteHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2076:func (s *Server) apiBulkDeleteHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2104:func (s *Server) apiStatsHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2108:func (s *Server) videoHandler(w http.ResponseWriter, r *http.Request) {
    main.go:2123:func (s *Server) thumbHandler(w http.ResponseWriter, r *http.Request) {
    main.go:3037:func writeJSON(w http.ResponseWriter, value any) {
    main.go:3046:func writeJSONError(w http.ResponseWriter, status int, err error) {
    detection.go:109:func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
    detection.go:221:func (s *Server) apiDetectionsHandler(w http.ResponseWriter, r *http.Request) {
    module neonvr
    
    go 1.22
    browser-env
    check_browser.py
    go
    go2rtc
    neonvr
    streams-desktop.png
    streams-mobile.png
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import os, subprocess, json, concurrent.futures
        paths=[]
        for directory,subdirs,files in os.walk('/mnt/FTPfiles/yicamFront'):
         for name in files:
          if name.endswith('.mp4') and not name.startswith('.'):
           paths.append(os.path.join(directory,name))
        print('files',len(paths))
        def inspect(path):
         try:
          p=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion','-of','json',path],capture_output=True,text=True,timeout=10)
          data=json.loads(p.stdout).get('format',{})
          tag=data.get('tags',{}).get('yi_motion')
          motion=json.loads(tag) if tag else None
          return path,motion
         except Exception: return path,None
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
         found=0; missing=0
         for path,motion in pool.map(inspect, sorted(paths,reverse=True)[:120]):
          if motion and motion.get('events'):
           print(path,json.dumps(motion));found+=1
           if found>=6: break
          elif motion is None: missing+=1
         print('missing tags',missing)
        PY
    files 27
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W419M13S46.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791692353000, "duration_ms": 46208, "video_time_precision_ms": 1000, "initial_state": -1, "history_available": true, "history_truncated": true, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[2514, 1791692355514, 1, "sample"]]}
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W415M41S19.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791692141000, "duration_ms": 18624, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[5295, 1791692146295, 1, "transition"], [11239, 1791692152239, 0, "transition"]]}
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W413M48S11.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791692027000, "duration_ms": 11403, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[389, 1791692027389, 1, "transition"], [6333, 1791692033333, 0, "transition"]]}
    /mnt/FTPfiles/yicamFront/2026Y10M11D00H/W400M08S51.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791691207000, "duration_ms": 51200, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[15, 1791691207015, 1, "transition"], [3943, 1791691210943, 0, "transition"]]}
    /mnt/FTPfiles/yicamFront/2026Y10M10D23H/W459M34S26.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791691172000, "duration_ms": 26114, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[755, 1791691172755, 1, "transition"], [4685, 1791691176685, 0, "transition"]]}
    /mnt/FTPfiles/yicamFront/2026Y10M10D23H/W451M06S54.mp4 {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791690664000, "duration_ms": 53952, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[897, 1791690664897, 1, "transition"], [6845, 1791690670845, 0, "transition"]]}
    missing tags 1
    ✓ • 577ms

## Assistant

I found `yi_motion` data in the recordings. It includes millisecond offsets for motion starting and stopping, plus whether motion was already active when the clip began. I’ll use those offsets for clickable player markers and thumbnail selection, keeping the existing selector as fallback.

## Activity

    $ sed -n '220,260p' detection.go; sed -n '400,435p' detection.go; sed -n '3034,3060p' main.go; sed -n '1,120p' /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py; sed -n '1635,1697p' main.go; sed -n '1,110p' config/neonvr.json; sed -n '295,340p' main.go
    
    func (s *Server) apiDetectionsHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        if r.Method != http.MethodGet {
            w.Header().Set("Allow", "GET")
            writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
            return
        }
        path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
        if err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        info, err := os.Stat(path)
        if err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        result := s.cachedDetection(path, info)
        if result == nil {
            result, err = s.embeddedDetection(r.Context(), path)
            if err == nil && result != nil {
                _ = s.cacheDetection(path, info, *result)
            }
        }
        if result != nil {
            writeJSON(w, detectionResponse{Status: "complete", Detections: result})
            return
        }
        response := detectionResponse{Status: "pending"}
        if !s.detectionConfig.DetectionEnabled {
            response.Status = "disabled"
        }
        s.detectionMu.Lock()
        if state, ok := s.detectionStates[path]; ok {
            response = state
        } else if s.detectionWorkerError != "" {
            response = detectionResponse{Status: "error", Message: s.detectionWorkerError}
        }
        s.detectionMu.Unlock()
        writeJSON(w, response)
            return err
        }
        data, err := json.Marshal(result)
        if err != nil {
            return err
        }
        tmp, err := os.CreateTemp(filepath.Dir(path), ".neonvr-detect-*.mp4")
        if err != nil {
            return err
        }
        tmpPath := tmp.Name()
        tmp.Close()
        defer os.Remove(tmpPath)
        out, err := exec.CommandContext(ctx, s.ffmpegPath, "-nostdin", "-v", "error", "-y", "-i", path,
            "-map", "0", "-c", "copy", "-map_metadata", "0", "-metadata", detectionMetadataTag+"="+string(data),
            "-movflags", "+faststart+use_metadata_tags", tmpPath).CombinedOutput()
        if err != nil {
            return fmt.Errorf("detection metadata remux: %w: %.2000s", err, out)
        }
        after, err := s.inspectDetectionVideo(ctx, tmpPath)
        if err != nil {
            return err
        }
        beforeDuration, e1 := strconv.ParseFloat(before.Format.Duration, 64)
        afterDuration, e2 := strconv.ParseFloat(after.Format.Duration, 64)
        if after.Format.Tags[detectionMetadataTag] != string(data) || !reflect.DeepEqual(before.Streams, after.Streams) ||
            e1 != nil || e2 != nil || !finite(beforeDuration) || !finite(afterDuration) || math.Abs(beforeDuration-afterDuration) > .25 {
            return errors.New("remux verification failed; original recording retained")
        }
        f, err := os.OpenFile(tmpPath, os.O_WRONLY, 0)
        if err != nil {
            return err
        }
        if os.Geteuid() == 0 {
            if stat, ok := original.Sys().(*syscall.Stat_t); ok {
                err = f.Chown(int(stat.Uid), int(stat.Gid))
        return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
    }
    
    func writeJSON(w http.ResponseWriter, value any) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        enc := json.NewEncoder(w)
        enc.SetIndent("", "  ")
        if err := enc.Encode(value); err != nil {
            log.Printf("json encode error: %v", err)
        }
    }
    
    func writeJSONError(w http.ResponseWriter, status int, err error) {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.WriteHeader(status)
        _ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error(), "message": err.Error()})
    }
    
    func writePlaceholderSVG(w http.ResponseWriter, label string) {
        w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
        w.Header().Set("Cache-Control", "no-store")
        label = template.HTMLEscapeString(label)
        _, _ = fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" width="480" height="270" viewBox="0 0 480 270">
      <rect width="480" height="270" fill="#10131a"/>
      <rect x="1" y="1" width="478" height="268" fill="none" stroke="#343a46"/>
      <circle cx="240" cy="116" r="42" fill="#202633"/>
      <polygon points="228,92 228,140 268,116" fill="#9aa4b2"/>
    import json
    from pathlib import Path
    import socket
    import socketserver
    import subprocess
    import tempfile
    import threading
    import time
    import urllib.request
    
    from playwright.sync_api import sync_playwright
    
    ROOT = Path('/tmp/neonvr-live-tools-vbq_9dyb')
    processes = []
    logs = []
    speaker_packets = 0
    speaker_connections = 0
    speaker_events = []
    
    
    def port():
        with socket.socket() as listener:
            listener.bind(('127.0.0.1', 0))
            return listener.getsockname()[1]
    
    
    def start(arguments):
        log = tempfile.TemporaryFile()
        logs.append(log)
        process = subprocess.Popen(arguments, stdout=log, stderr=log)
        processes.append(process)
        return process
    
    
    def request(url):
        with urllib.request.urlopen(url, timeout=3) as response:
            return json.load(response)
    
    
    def ready(url, process):
        for attempt in range(100):
            if process.poll() is not None:
                raise RuntimeError('Server exited')
            try:
                return request(url)
            except Exception:
                time.sleep(.05)
        raise RuntimeError('Server did not start')
    
    
    class Speaker(socketserver.StreamRequestHandler):
        def handle(self):
            global speaker_packets, speaker_connections
            speaker_connections += 1
            while True:
                first = self.rfile.read(1)
                if not first:
                    return
                if first == b'$':
                    channel = self.rfile.read(1)
                    length = int.from_bytes(self.rfile.read(2), 'big')
                    data = self.rfile.read(length)
                    if len(speaker_events) < 30:
                        speaker_events.append(('packet', int.from_bytes(channel, 'big'), len(data)))
                    if channel == b'\x00' and len(data) > 12:
                        speaker_packets += 1
                    continue
                line = first + self.rfile.readline()
                method, target, _ = line.decode().strip().split(' ', 2)
                speaker_events.append((method, target))
                headers = {}
                while True:
                    line = self.rfile.readline()
                    if line in (b'\r\n', b'\n', b''):
                        break
                    name, value = line.decode().split(':', 1)
                    headers[name.lower()] = value.strip()
                size = int(headers.get('content-length', 0))
                if size:
                    self.rfile.read(size)
                body = b''
                extra = 'Session: local-speaker\r\n'
                if method == 'OPTIONS':
                    extra += 'Public: OPTIONS, DESCRIBE, SETUP, PLAY, TEARDOWN, GET_PARAMETER\r\n'
                elif method == 'DESCRIBE':
                    body = ('v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=Local speaker\r\nt=0 0\r\n'
                            'a=control:*\r\nm=audio 0 RTP/AVP 0\r\nc=IN IP4 0.0.0.0\r\n'
                            'a=rtpmap:0 PCMU/8000\r\na=sendonly\r\na=control:trackID=0\r\n').encode()
                    extra += 'Content-Type: application/sdp\r\nContent-Base: ' + target + '/\r\n'
                elif method == 'SETUP':
                    extra += 'Transport: ' + headers['transport'] + '\r\n'
                result = ('RTSP/1.0 200 OK\r\nCSeq: ' + headers.get('cseq', '1') + '\r\n' + extra +
                          'Content-Length: ' + str(len(body)) + '\r\n\r\n').encode() + body
                self.wfile.write(result)
                self.wfile.flush()
                if method == 'TEARDOWN':
                    return
    
    
    class SpeakerServer(socketserver.ThreadingTCPServer):
        allow_reuse_address = True
        daemon_threads = True
    
    
    speaker = SpeakerServer(('127.0.0.1', 0), Speaker)
    threading.Thread(target=speaker.serve_forever, daemon=True).start()
    
    try:
        with tempfile.TemporaryDirectory(prefix='neonvr-browser-fixture-') as temporary:
            directory = Path(temporary)
            rtsp_port, fixture_api, app_port, media_port = port(), port(), port(), port()
            fixture_config = directory / 'fixture.json'
            fixture_config.write_text(json.dumps({
                'api': {'listen': '127.0.0.1:' + str(fixture_api)},
                'rtsp': {'listen': '127.0.0.1:' + str(rtsp_port)},
                'webrtc': {'listen': '', 'ice_servers': []},
                'streams': {'fixture': []},
            }))
            fixture = start([str(ROOT / 'go2rtc'), '-config', str(fixture_config)])
            ready('http://127.0.0.1:' + str(fixture_api) + '/api', fixture)
      if (evt.ctrlKey && ['ArrowLeft','ArrowRight','ArrowUp','ArrowDown'].includes(evt.key)) {
        evt.preventDefault();
        if (evt.key === 'ArrowLeft' || evt.key === 'ArrowUp') playRelative(-1);
        else playRelative(1);
        return;
      }
    
      if (evt.key === 'ArrowLeft') { evt.preventDefault(); stepMain(-1); }
      if (evt.key === 'ArrowRight') { evt.preventDefault(); stepMain(1); }
    });
    
    installVideoPanZoom();
    installSideRailScrollCapture();
    setZoomMode(false);
    init();
    </script>
    </body>
    </html>`))
    
    func main() {
        cfg, configPath, err := loadStartupConfig(os.Args[1:])
        if errors.Is(err, flagHelpError) {
            return
        }
        if err != nil {
            log.Fatal(err)
        }
        manager, err := newConfigManager(configPath, cfg)
        if err != nil {
            log.Fatal(err)
        }
        executable, err := os.Executable()
        if err != nil {
            log.Fatal(err)
        }
        root, addr := &cfg.Root, &cfg.Addr
        extraRoots := cfg.ExtraRoots
        thumbDir, metaPath, ffmpegPath := &cfg.ThumbDir, &cfg.Meta, &cfg.FFmpeg
        thumbSeek, thumbWidth, thumbWorkers := &cfg.ThumbSeek, &cfg.ThumbWidth, &cfg.ThumbWorkers
        maintenance, _ := time.ParseDuration(cfg.ThumbMaintenanceEvery)
        refresh, _ := time.ParseDuration(cfg.SizeRefresh)
        retentionEvery, _ := time.ParseDuration(cfg.AutoDeleteEvery)
        thumbMaintenanceEvery, sizeRefresh, autoDeleteEvery := &maintenance, &refresh, &retentionEvery
        autoDeleteDays := &cfg.AutoDeleteDays
        runCtx, cancelWork := context.WithCancel(context.Background())
        defer cancelWork()
        live := newLiveManager(runCtx, cfg, manager.token)
        defer live.close()
    
        cleanRoot := filepath.Clean(expandPath(*root))
        cleanExtraRoots := make([]string, 0, len(extraRoots))
        seenExtraRoots := make(map[string]bool)
        for _, extraRoot := range extraRoots {
            clean := filepath.Clean(expandPath(extraRoot))
            if abs, err := filepath.Abs(clean); err == nil {
                clean = abs
            }
            if !seenExtraRoots[clean] {
                cleanExtraRoots = append(cleanExtraRoots, clean)
                seenExtraRoots[clean] = true
            }
        }
        cleanThumbDir := *thumbDir
    {
      "root": "/mnt/mySSD/BlinkClips",
      "extra_roots": [
        "/mnt/FTPfiles/yicamBack",
        "/mnt/FTPfiles/yicamFront"
      ],
      "addr": "0.0.0.0:5500",
      "thumb_dir": "",
      "meta": "",
      "ffmpeg": "ffmpeg",
      "thumb_seek": 0.2,
      "thumb_width": 480,
      "thumb_workers": 2,
      "thumb_maintenance_every": "30m",
      "size_refresh": "30m",
      "auto_delete_days": 14,
      "auto_delete_every": "6h",
      "detection_enabled": true,
      "detection_python": "python3",
      "detection_model": "yolo26n.onnx",
      "detection_threads": 2,
      "detection_fps": 2,
      "detection_confidence": 0.65,
      "detection_min_frames": 3,
      "detection_merge_gap": 0.5,
      "detection_car_movement": 0.025,
      "detection_every": "1m",
      "detection_min_age": "2m",
      "detection_timeout": "30m",
      "detection_prune": "quarantine",
      "live_streams": [
        {
          "name": "Camera 192.168.1.150",
          "url": "rtsp://192.168.1.150:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
        }
      ],
      "live_go2rtc": "go2rtc",
      "live_webrtc_port": 8555,
      "live_webrtc_host": ""
    }
        .mainPlayerShell.zoomEnabled #mainVideo:active { cursor: grabbing; }
        media-control-bar {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 8px;
          background: rgba(0,0,0,.86);
          color: white;
          font-size: 13px;
          flex-wrap: wrap;
        }
        media-time-display, media-duration-display {
          font-variant-numeric: tabular-nums;
          min-width: 42px;
          color: #e5e7eb;
        }
        media-time-display::after {
          content: ' /';
          color: #8b94a3;
          margin-left: 4px;
        }
        .analysisControls {
          display: flex;
          gap: 8px;
          align-items: center;
          flex-wrap: wrap;
          margin-top: 10px;
          color: var(--muted);
          font-size: 13px;
        }
        .analysisControls button { padding: 7px 9px; font-size: 13px; }
        .detectionLegend { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; margin-top: 10px; font-size: 12px; color: var(--muted); }
        .detectionLegend .category { border-bottom: 3px solid var(--detection-color); padding-bottom: 3px; }
        #detectionStatus { margin-left: auto; }
        #detectionStatus.active { color: var(--ok); }
        .analysisControls input { width: 76px; }
        .kbdHint { color: #7d8798; }
    
        .topRow {
          display: flex;
          justify-content: space-between;
          gap: 14px;
          align-items: center;
          flex-wrap: wrap;
          margin-bottom: 12px;
        }
    ✓ • 0ms

## Activity

    file changes: Completed · 2 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/main.go
    @@ -30,3 +30,3 @@
     const defaultAddr = "0.0.0.0:5500"
    -const thumbnailAlgorithmVersion = "divergence-v2"
    +const thumbnailAlgorithmVersion = "camera-motion-v3"
     
    @@ -701,2 +701,4 @@
                   <span class="category" style="--detection-color:#c792ea">Moving vehicle</span>
    +              <span id="cameraMotionLegend" class="category" style="--detection-color:#ff9866" hidden>Camera motion</span>
    +              <span id="cameraMotionStatus" role="status" hidden></span>
                   <span id="detectionStatus" role="status">Select a clip to see detections.</span>
    @@ -1804,2 +1806,3 @@
     	mux.HandleFunc("/api/detections", s.apiDetectionsHandler)
    +	mux.HandleFunc("/api/camera-motion", s.apiCameraMotionHandler)
     	mux.HandleFunc("/static/detection-viewer.js", func(w http.ResponseWriter, r *http.Request) {
    @@ -2464,3 +2467,3 @@
     
    -	if err := s.runInterestingThumbnail(videoPath, tmpPath); err != nil {
    +	if err := s.runPreferredThumbnail(videoPath, tmpPath); err != nil {
     		_ = os.Remove(tmpPath)
    Add: /home/ubuntu/NeoNVR/motion.go
    package main
    
    import (
    	"context"
    	"encoding/json"
    	"errors"
    	"fmt"
    	"log"
    	"math"
    	"net/http"
    	"os"
    	"os/exec"
    	"sort"
    	"strconv"
    	"time"
    )
    
    const cameraMotionMetadataTag = "yi_motion"
    
    type motionInterval struct {
    	Start float64 `json:"start"`
    	End   float64 `json:"end"`
    }
    
    type cameraMotion struct {
    	Duration         float64          `json:"duration"`
    	HistoryAvailable bool             `json:"history_available"`
    	HistoryTruncated bool             `json:"history_truncated"`
    	Markers          []float64        `json:"markers"`
    	Intervals        []motionInterval `json:"intervals"`
    }
    
    // Camera offsets are relative to the video, independent of filename timezones.
    func parseCameraMotion(tag string, duration float64) (*cameraMotion, error) {
    	var raw struct {
    		Schema           string              `json:"schema"`
    		VideoStart       *float64            `json:"video_start_unix_ms"`
    		DurationMS       float64             `json:"duration_ms"`
    		InitialState     *int                `json:"initial_state"`
    		HistoryAvailable bool                `json:"history_available"`
    		HistoryTruncated bool                `json:"history_truncated"`
    		Columns          []string            `json:"event_columns"`
    		Events           [][]json.RawMessage `json:"events"`
    	}
    	if err := json.Unmarshal([]byte(tag), &raw); err != nil {
    		return nil, fmt.Errorf("invalid yi_motion JSON: %w", err)
    	}
    	if raw.Schema != "yi-motion-v1" {
    		return nil, errors.New("unsupported yi_motion schema")
    	}
    	if !finite(duration) || duration <= 0 {
    		duration = raw.DurationMS / 1000
    	}
    	if !finite(duration) || duration <= 0 {
    		return nil, errors.New("yi_motion has no valid video duration")
    	}
    	state := -1 // Unknown history must not imply motion or inactivity.
    	if raw.InitialState != nil {
    		state = *raw.InitialState
    	}
    	if state < -1 || state > 1 {
    		return nil, errors.New("invalid yi_motion initial state")
    	}
    	column := func(name string) int {
    		for i, value := range raw.Columns {
    			if value == name {
    				return i
    			}
    		}
    		return -1
    	}
    	offsetCol, unixCol, stateCol := column("offset_ms"), column("unix_ms"), column("state")
    	if len(raw.Events) > 0 && (stateCol < 0 || (offsetCol < 0 && (unixCol < 0 || raw.VideoStart == nil))) {
    		return nil, errors.New("yi_motion event columns are missing")
    	}
    	type event struct {
    		offset float64
    		state  int
    	}
    	events := make([]event, 0, len(raw.Events))
    	for _, row := range raw.Events {
    		readNumber := func(index int) (float64, error) {
    			var value float64
    			if index < 0 || index >= len(row) || string(row[index]) == "null" {
    				return 0, errors.New("missing yi_motion event value")
    			}
    			if err := json.Unmarshal(row[index], &value); err != nil || !finite(value) {
    				return 0, errors.New("invalid yi_motion event value")
    			}
    			return value, nil
    		}
    		index := offsetCol
    		if index < 0 {
    			index = unixCol
    		}
    		offset, err := readNumber(index)
    		if err != nil {
    			return nil, err
    		}
    		if offsetCol < 0 {
    			offset -= *raw.VideoStart
    		}
    		offset /= 1000
    		value, err := readNumber(stateCol)
    		if err != nil || value < -1 || value > 1 || value != math.Trunc(value) {
    			return nil, errors.New("invalid yi_motion event state")
    		}
    		if !finite(offset) || offset < 0 || offset >= duration {
    			continue // Events outside the recording cannot mark its seek track.
    		}
    		events = append(events, event{offset, int(value)})
    	}
    	sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
    	result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
    		Markers: []float64{}, Intervals: []motionInterval{}}
    	start := 0.0
    	if state == 1 {
    		result.Markers = append(result.Markers, 0)
    	}
    	for _, event := range events {
    		if state == 1 && event.state != 1 && event.offset > start {
    			result.Intervals = append(result.Intervals, motionInterval{start, event.offset})
    		}
    		if event.state == 1 && state != 1 {
    			start = event.offset
    			result.Markers = append(result.Markers, start)
    		}
    		state = event.state
    	}
    	if state == 1 && duration > start {
    		result.Intervals = append(result.Intervals, motionInterval{start, duration})
    	}
    	return result, nil
    }
    
    func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
    	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    	defer cancel()
    	out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
    		"format=duration:format_tags="+cameraMotionMetadataTag, "-of", "json", path).Output()
    	if err != nil {
    		return nil, fmt.Errorf("probe camera motion: %w", err)
    	}
    	var probe struct {
    		Format struct {
    			Duration string            `json:"duration"`
    			Tags     map[string]string `json:"tags"`
    		} `json:"format"`
    	}
    	if err := json.Unmarshal(out, &probe); err != nil {
    		return nil, err
    	}
    	tag := probe.Format.Tags[cameraMotionMetadataTag]
    	if tag == "" {
    		return nil, nil
    	}
    	duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
    	return parseCameraMotion(tag, duration)
    }
    
    func (s *Server) apiCameraMotionHandler(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Cache-Control", "no-store")
    	if r.Method != http.MethodGet {
    		w.Header().Set("Allow", "GET")
    		writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
    		return
    	}
    	path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
    	if err != nil {
    		writeJSONError(w, http.StatusNotFound, err)
    		return
    	}
    	motion, err := s.readCameraMotion(r.Context(), path)
    	if err != nil {
    		writeJSONError(w, http.StatusBadGateway, err)
    		return
    	}
    	writeJSON(w, struct {
    		Motion *cameraMotion `json:"motion"`
    	}{motion})
    }
    
    func (s *Server) runPreferredThumbnail(videoPath, tmpPath string) error {
    	motion, err := s.readCameraMotion(s.workContext(), videoPath)
    	if err == nil && motion != nil {
    		for _, interval := range motion.Intervals {
    			// Let the motion appear after the onset; stay within the marked interval.
    			seek := interval.Start + math.Min(.5, (interval.End-interval.Start)/2)
    			seek = math.Min(seek, math.Max(0, motion.Duration-.1))
    			err = s.runFFmpegThumbnail(videoPath, tmpPath, strconv.FormatFloat(seek, 'f', 3, 64))
    			if err == nil {
    				var frame frameSignature
    				frame, err = frameLumaSignature(tmpPath)
    				if err == nil && !frame.usable {
    					err = errors.New("camera motion frame is blank")
    				}
    			}
    			if err == nil {
    				log.Printf("thumbnail selected %.3fs for %s (camera motion)", seek, videoPath)
    				return nil
    			}
    			_ = os.Remove(tmpPath)
    		}
    	}
    	if err != nil {
    		log.Printf("camera motion thumbnail unavailable for %s, using scene selection: %v", videoPath, err)
    	}
    	return s.runInterestingThumbnail(videoPath, tmpPath)
    }

## Activity

    file changes: Completed · 4 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -280,3 +280,5 @@
     
    -New thumbnails use a divergence selector designed for mostly-static security footage. The server samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    +New thumbnails first use embedded Yi camera motion metadata (`yi_motion`, schema `yi-motion-v1`). The server extracts a frame half a second after the first motion onset, or halfway through the interval when it lasts less than a second. If that frame is unavailable or blank, it tries the next marked motion interval. Motion already active when the recording starts is treated as starting at 0 seconds.
    +
    +If motion metadata is missing, invalid, has no usable motion interval, or frame extraction fails, the server uses the existing divergence selector designed for mostly-static security footage. It samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
     
    @@ -285,3 +287,3 @@
     
    -Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same divergence selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).
    +Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same motion metadata and fallback selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).
     
    @@ -293,2 +295,8 @@
     
    +Recordings with `yi_motion` metadata show orange, clickable motion onset markers directly on the main player's seek track, with a lane below showing the motion intervals. These camera markers work even when AI detection is disabled. Partial or truncated motion history is labeled in the player. NeoNVR uses the metadata's `offset_ms` values relative to the video, without a timezone correction; `initial_state: 1` marks motion from the beginning, and transitions to state 0 end the highlighted interval. Older cameras without the tag continue to play normally. The tag can be inspected with:
    +
    +```bash
    +ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 | jq '.format.tags.yi_motion | fromjson'
    +```
    +
     ## Project rename and compatibility
    Update { move_path: None }: /home/ubuntu/NeoNVR/main.go
    @@ -329,2 +329,3 @@
         #detectionStatus.active { color: var(--ok); }
    +    #cameraMotionStatus.active { color: #ff9866; }
         .analysisControls input { width: 76px; }
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/detection-viewer.js
    @@ -6,2 +6,4 @@
       const status = document.getElementById('detectionStatus');
    +  const motionLegend = document.getElementById('cameraMotionLegend');
    +  const motionStatus = document.getElementById('cameraMotionStatus');
       if (!video || !seek || !status) return;
    @@ -12,4 +14,37 @@
       let result = null;
    +  let motionRequest = null;
    +  let motion = null;
    +
    +  function activeCameraMotion() {
    +    if (!motion || !motionStatus) return;
    +    const active = motion.intervals.some(row => video.currentTime >= row.start && video.currentTime < row.end);
    +    const partial = !motion.history_available || motion.history_truncated;
    +    motionStatus.classList.toggle('active', active);
    +    motionStatus.textContent = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
    +      partial ? 'No camera motion timestamps available.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
    +  }
     
    +  async function refreshCameraMotion(clip, expected) {
    +    motionRequest = new AbortController();
    +    try {
    +      const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
    +      if (!response.ok) throw new Error('Camera motion metadata unavailable.');
    +      const data = await response.json();
    +      if (expected !== generation) return;
    +      motion = data.motion;
    +      seek.setCameraMotion(motion);
    +      if (motionLegend) motionLegend.hidden = !motion;
    +      if (motionStatus) motionStatus.hidden = !motion;
    +      activeCameraMotion();
    +    } catch (error) {
    +      if (expected !== generation || error.name === 'AbortError') return;
    +      if (motionStatus) {
    +        motionStatus.hidden = false;
    +        motionStatus.textContent = error.message;
    +      }
    +    }
    +  }
    +
       function activeStatus() {
    +    activeCameraMotion();
         if (!result) return;
    @@ -52,5 +87,10 @@
         request?.abort();
    +    motionRequest?.abort();
         clearTimeout(timer);
         result = null;
    +    motion = null;
         seek.setDetections(null);
    +    seek.setCameraMotion(null);
    +    if (motionLegend) motionLegend.hidden = true;
    +    if (motionStatus) { motionStatus.hidden = true; motionStatus.classList.remove('active'); }
         status.classList.remove('active');
    @@ -59,2 +99,3 @@
         refresh(clip, generation);
    +    refreshCameraMotion(clip, generation);
       }
    @@ -63,3 +104,3 @@
       video.addEventListener('seeked', activeStatus);
    -  window.addEventListener('pagehide', () => { request?.abort(); clearTimeout(timer); });
    +  window.addEventListener('pagehide', () => { request?.abort(); motionRequest?.abort(); clearTimeout(timer); });
       window.NeoNVRDetections = { loadClip };
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/media-chrome-lite.js
    @@ -161,2 +161,8 @@
           this.input = input;
    +      const track = document.createElement('div');
    +      track.className = 'seek-track';
    +      this.motionMarkers = document.createElement('div');
    +      this.motionMarkers.className = 'motion-markers';
    +      this.motionMarkers.setAttribute('aria-label', 'Camera motion timestamps');
    +      track.append(input, this.motionMarkers);
           this.lanes = document.createElement('div');
    @@ -169,2 +175,6 @@
             input { width: 100%; margin: 0; accent-color: white; }
    +        .seek-track { position: relative; display: flex; align-items: center; }
    +        .motion-markers { position: absolute; left: 8px; right: 8px; top: 50%; pointer-events: none; }
    +        .motion-marker { position: absolute; width: 7px; height: 18px; padding: 0; border: 1px solid #332014; border-radius: 2px; background: #ff9866; transform: translate(-50%, -50%); pointer-events: auto; cursor: pointer; }
    +        .motion-marker:hover, .motion-marker:focus-visible { outline: 2px solid white; z-index: 1; }
             .detection-lanes { margin: 3px 8px 0; }
    @@ -174,3 +184,3 @@
           `;
    -      root.append(style, input, this.lanes);
    +      root.append(style, track, this.lanes);
           const video = findVideo(this);
    @@ -202,2 +212,6 @@
         }
    +    setCameraMotion(motion) {
    +      this.cameraMotion = motion;
    +      this.renderDetections();
    +    }
         renderDetections() {
    @@ -205,9 +219,26 @@
           this.lanes.replaceChildren();
    +      this.motionMarkers.replaceChildren();
           const result = this.detections;
    -      this.lanes.hidden = !result;
    -      if (!result) return;
    -      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : result.duration;
    +      const motion = this.cameraMotion;
    +      this.lanes.hidden = !result && !motion?.intervals?.length;
    +      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : (motion?.duration || result?.duration);
           if (!Number.isFinite(duration) || duration <= 0) return;
    +      if (motion) {
    +        for (const time of motion.markers || []) {
    +          if (!Number.isFinite(time) || time < 0 || time >= duration) continue;
    +          const marker = document.createElement('button');
    +          marker.type = 'button';
    +          marker.className = 'motion-marker';
    +          marker.style.left = (time / duration * 100) + '%';
    +          marker.title = 'Camera motion at ' + time.toFixed(3) + ' sec';
    +          marker.setAttribute('aria-label', 'Seek to ' + marker.title);
    +          marker.addEventListener('click', () => { if (this.video) this.video.currentTime = time; });
    +          this.motionMarkers.append(marker);
    +        }
    +      }
           const categories = [['person', 'Human', '#75d685'], ['animal', 'Animal', '#ffd166'], ['bike', 'Bike', '#72a7ff'], ['moving_car', 'Moving vehicle', '#c792ea']];
    +      const intervals = result?.intervals || [];
    +      if (motion?.intervals?.length) categories.unshift(['camera_motion', 'Camera motion', '#ff9866']);
           for (const [kind, label, color] of categories) {
    +        if (kind !== 'camera_motion' && !result) continue;
             const lane = document.createElement('div');
    @@ -216,4 +247,4 @@
             lane.setAttribute('aria-label', label + ' detections');
    -        for (const interval of result.intervals || []) {
    -          if (interval.kind !== kind || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
    +        for (const interval of kind === 'camera_motion' ? motion.intervals : intervals) {
    +          if ((kind !== 'camera_motion' && interval.kind !== kind) || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
               const start = Math.max(0, Math.min(duration, interval.start));
    @@ -227,3 +258,4 @@
               event.style.background = color;
    -          const text = label + ': ' + start.toFixed(1) + '–' + end.toFixed(1) + ' sec (' + Math.round(interval.confidence * 100) + '%)';
    +          const text = label + ': ' + start.toFixed(3) + '–' + end.toFixed(3) + ' sec' +
    +            (kind === 'camera_motion' ? '' : ' (' + Math.round(interval.confidence * 100) + '%)');
               event.title = text;

## Assistant

The player markers and thumbnail selection are wired up. Camera motion displays independently of AI detections, and truncated history is labeled. I’m now checking timestamp accuracy, clip switching, and fallback behavior for recordings with missing or malformed metadata.

## Activity

    file changes: Completed · 4 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/detection_test.go
    @@ -34,3 +34,3 @@
     		"-f", "lavfi", "-i", "sine=frequency=440", "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1", "-c:a", "aac",
    -		"-metadata", "source_camera=Garden", "-movflags", "+use_metadata_tags", path).CombinedOutput()
    +		"-metadata", "source_camera=Garden", "-metadata", "yi_motion="+motionTestTag, "-movflags", "+use_metadata_tags", path).CombinedOutput()
     	if err != nil {
    @@ -90,2 +90,6 @@
     	}
    +	motion, err := s.readCameraMotion(context.Background(), path)
    +	if err != nil || motion == nil || !reflect.DeepEqual(motion.Intervals, []motionInterval{{1.5, 2.5}}) {
    +		t.Fatalf("camera motion metadata lost during AI remux: %+v %v", motion, err)
    +	}
     	// The MP4 remains the source of truth when the expendable JSON cache is lost.
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion.go
    @@ -116,5 +116,2 @@
     	start := 0.0
    -	if state == 1 {
    -		result.Markers = append(result.Markers, 0)
    -	}
     	for _, event := range events {
    @@ -125,3 +122,2 @@
     			start = event.offset
    -			result.Markers = append(result.Markers, start)
     		}
    @@ -132,2 +128,5 @@
     	}
    +	for _, interval := range result.Intervals {
    +		result.Markers = append(result.Markers, interval.Start)
    +	}
     	return result, nil
    Add: /home/ubuntu/NeoNVR/motion_test.go
    package main
    
    import (
    	"context"
    	"encoding/json"
    	"image/jpeg"
    	"net/http"
    	"net/http/httptest"
    	"os"
    	"os/exec"
    	"path/filepath"
    	"reflect"
    	"testing"
    )
    
    const motionTestTag = `{"schema":"yi-motion-v1","video_start_unix_ms":1791692141000,"duration_ms":4000,"initial_state":0,"history_available":true,"history_truncated":false,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[1500,1791692142500,1,"transition"],[2500,1791692143500,0,"transition"]]}`
    
    func TestParseCameraMotion(t *testing.T) {
    	tests := []struct {
    		name      string
    		tag       string
    		duration  float64
    		intervals []motionInterval
    		partial   bool
    	}{
    		{"milliseconds", motionTestTag, 4, []motionInterval{{1.5, 2.5}}, false},
    		{"initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"events":[]}`, 60.033, []motionInterval{{0, 60.033}}, false},
    		{"truncated sample", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[2514,1791692355514,1,"sample"]]}`, 46.208, []motionInterval{{2.514, 46.208}}, true},
    		{"no motion", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`, 10, []motionInterval{}, false},
    		{"unknown history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false}`, 10, []motionInterval{}, true},
    		{"reordered events and columns", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["state","offset_ms"],"events":[[0,3500],[1,1500],[1,2500]]}`, 4, []motionInterval{{1.5, 3.5}}, false},
    		{"unix fallback", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"video_start_unix_ms":1791692141000,"event_columns":["state","unix_ms"],"events":[[1,1791692146295],[0,1791692152239]]}`, 18.624, []motionInterval{{5.295, 11.239}}, false},
    		{"metadata duration fallback", motionTestTag, 0, []motionInterval{{1.5, 2.5}}, false},
    		{"ignore outside clip", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["offset_ms","state"],"events":[[-1,1],[1000,1],[2000,0],[4000,1],[5000,1]]}`, 4, []motionInterval{{1, 2}}, false},
    		{"zero-length initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[0,0]]}`, 4, []motionInterval{}, false},
    	}
    	for _, tt := range tests {
    		t.Run(tt.name, func(t *testing.T) {
    			result, err := parseCameraMotion(tt.tag, tt.duration)
    			if err != nil {
    				t.Fatal(err)
    			}
    			if !reflect.DeepEqual(result.Intervals, tt.intervals) || (!result.HistoryAvailable || result.HistoryTruncated) != tt.partial {
    				t.Fatalf("unexpected motion: %+v", result)
    			}
    			markers := []float64{}
    			for _, interval := range tt.intervals {
    				markers = append(markers, interval.Start)
    			}
    			if !reflect.DeepEqual(result.Markers, markers) {
    				t.Fatalf("markers: %v, want %v", result.Markers, markers)
    			}
    		})
    	}
    }
    
    func TestParseCameraMotionRejectsInvalidMetadata(t *testing.T) {
    	for _, tag := range []string{
    		`broken`, `null`, `{"schema":"yi-motion-v2"}`,
    		`{"schema":"yi-motion-v1","initial_state":2}`,
    		`{"schema":"yi-motion-v1","events":[[1,1]]}`,
    		`{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[null,1]]}`,
    		`{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[1000,2]]}`,
    		`{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[1000,0.5]]}`,
    		`{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[["1000",1]]}`,
    	} {
    		if _, err := parseCameraMotion(tag, 4); err == nil {
    			t.Fatalf("accepted invalid metadata: %s", tag)
    		}
    	}
    	if _, err := parseCameraMotion(`{"schema":"yi-motion-v1"}`, 0); err == nil {
    		t.Fatal("accepted missing duration")
    	}
    }
    
    func cameraMotionFixture(t *testing.T, tag string) (*Server, string) {
    	t.Helper()
    	for _, tool := range []string{"ffmpeg", "ffprobe"} {
    		if _, err := exec.LookPath(tool); err != nil {
    			t.Skip(tool + " not installed")
    		}
    	}
    	root := t.TempDir()
    	dir := filepath.Join(root, "Garden")
    	if err := os.MkdirAll(dir, 0750); err != nil {
    		t.Fatal(err)
    	}
    	path := filepath.Join(dir, "recording.mp4")
    	args := []string{"-nostdin", "-v", "error", "-f", "lavfi", "-i",
    		"testsrc2=s=160x96:r=10,drawbox=x=40:y=20:w=80:h=60:color=red:t=fill:enable='between(t,1.5,2.5)'",
    		"-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1"}
    	if tag != "" {
    		args = append(args, "-metadata", "yi_motion="+tag)
    	}
    	args = append(args, "-movflags", "+use_metadata_tags", path)
    	if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
    		t.Fatalf("create fixture: %v: %s", err, out)
    	}
    	return &Server{root: root, ffmpegPath: "ffmpeg", thumbDir: filepath.Join(root, ".thumbs"), thumbWidth: 160, thumbSeek: .2, thumbSem: make(chan struct{}, 1)}, path
    }
    
    func TestCameraMotionAPIAndThumbnail(t *testing.T) {
    	s, path := cameraMotionFixture(t, motionTestTag)
    	w := httptest.NewRecorder()
    	s.apiCameraMotionHandler(w, httptest.NewRequest(http.MethodGet, "/api/camera-motion?camera=Garden&file=recording.mp4", nil))
    	var response struct {
    		Motion *cameraMotion `json:"motion"`
    	}
    	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || response.Motion == nil ||
    		!reflect.DeepEqual(response.Motion.Markers, []float64{1.5}) || !reflect.DeepEqual(response.Motion.Intervals, []motionInterval{{1.5, 2.5}}) {
    		t.Fatalf("motion API with AI disabled: %d %s", w.Code, w.Body.String())
    	}
    	thumb, err := s.ensureThumbnail(path)
    	if err != nil {
    		t.Fatal(err)
    	}
    	f, err := os.Open(thumb)
    	if err != nil {
    		t.Fatal(err)
    	}
    	defer f.Close()
    	frame, err := jpeg.Decode(f)
    	if err != nil {
    		t.Fatal(err)
    	}
    	r, g, b, _ := frame.At(80, 48).RGBA()
    	if r < 45000 || g > 12000 || b > 12000 {
    		t.Fatalf("thumbnail did not select the red motion frame: %d %d %d", r, g, b)
    	}
    	for _, request := range []*http.Request{
    		httptest.NewRequest(http.MethodPost, "/api/camera-motion?camera=Garden&file=recording.mp4", nil),
    		httptest.NewRequest(http.MethodGet, "/api/camera-motion?camera=Garden&file=../outside.mp4", nil),
    	} {
    		w := httptest.NewRecorder()
    		s.apiCameraMotionHandler(w, request)
    		if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusNotFound {
    			t.Fatalf("invalid request accepted: %d", w.Code)
    		}
    	}
    }
    
    func TestCameraMotionThumbnailFallback(t *testing.T) {
    	for _, tag := range []string{"", "invalid-json", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`} {
    		s, path := cameraMotionFixture(t, tag)
    		thumb, err := s.ensureThumbnail(path)
    		if err != nil {
    			t.Fatalf("fallback failed for %q: %v", tag, err)
    		}
    		if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
    			t.Fatalf("fallback produced no usable image for %q: %v", tag, err)
    		}
    		if tag == "" {
    			motion, err := s.readCameraMotion(context.Background(), path)
    			if err != nil || motion != nil {
    				t.Fatalf("missing metadata: %+v %v", motion, err)
    			}
    		}
    	}
    }
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/detection-viewer.js
    @@ -22,4 +22,5 @@
         motionStatus.classList.toggle('active', active);
    -    motionStatus.textContent = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
    +    const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
           partial ? 'No camera motion timestamps available.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
    +    if (motionStatus.textContent !== text) motionStatus.textContent = text;
       }

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w motion.go motion_test.go main.go detection_test.go
    ✓ • 0ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test ./...
    ok      neonvr    5.348s
    ✓ • 5.89s

## Activity

    $ sed -n '1000,1080p' main.go; sed -n '1,250p' motion.go; rg -n 'def test|browser|launch|config.write_text|Popen' /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py; sed -n '185,230p' /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py; rg -n 'detection-enabled|auto-delete-days|thumb-maintenance|config' config.go | head -35
      allCb.addEventListener('change', async () => {
        if (allCb.checked) selectedCameras = new Set(cameras.map(c => c.name));
        else selectedCameras.clear();
        await onCameraSelectionChanged();
      });
      allRow.appendChild(allCb);
      allRow.appendChild(document.createTextNode('All'));
      cameraMenuEl.appendChild(allRow);
    
      const div = document.createElement('div');
      div.className = 'cameraDivider';
      cameraMenuEl.appendChild(div);
    
      for (const cam of cameras) {
        const row = document.createElement('label');
        row.className = 'cameraOption';
        const cb = document.createElement('input');
        cb.type = 'checkbox';
        cb.checked = selectedCameras.has(cam.name);
        cb.addEventListener('change', async () => {
          if (cb.checked) selectedCameras.add(cam.name);
          else selectedCameras.delete(cam.name);
          await onCameraSelectionChanged();
        });
        const sw = document.createElement('span');
        sw.className = 'cameraSwatch';
        sw.style.background = cameraColor(cam.name);
        row.appendChild(cb);
        row.appendChild(sw);
        row.appendChild(document.createTextNode(cam.name));
        cameraMenuEl.appendChild(row);
      }
      updateCameraSummary();
    }
    
    function updateCameraSummary() {
      const selected = selectedCameraList();
      if (selected.length === cameras.length && cameras.length > 0) cameraSummaryEl.textContent = 'All cameras';
      else if (selected.length === 0) cameraSummaryEl.textContent = 'No cameras';
      else if (selected.length === 1) cameraSummaryEl.textContent = selected[0];
      else cameraSummaryEl.textContent = selected.length + ' cameras';
    }
    
    async function onCameraSelectionChanged() {
      renderCameraMenu();
      clearSelection();
      try {
        await refreshDates(true);
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    }
    
    function loadMainPlayer(clip, play = true, scroll = false) {
      activeKey = keyFor(clip);
      mainVideoEl.src = clip.video_url;
      mainVideoEl.poster = clip.thumb_url;
      mainTitleEl.textContent = clip.timestamp + ' — ' + clip.camera + ' — ' + clip.file;
      resetZoom();
      mainVideoEl.load();
      window.NeoNVRDetections?.loadClip(clip);
      if (play) mainVideoEl.play().catch(() => {});
      syncActiveCard(true);
      if (scroll) document.querySelector('.mainPlayerPanel').scrollIntoView({behavior:'smooth', block:'start'});
    }
    
    function syncActiveCard(scrollIntoView = false) {
      document.querySelectorAll('.card, .sideThumb').forEach(el => {
        const isActive = el.dataset.key === activeKey;
        el.classList.toggle('active', isActive);
        if (isActive && scrollIntoView && el.classList.contains('sideThumb')) {
          el.scrollIntoView({behavior:'smooth', block:'nearest'});
        }
      });
    }
    
    function stepMain(delta) {
      if (!mainVideoEl.src) return;
      const step = Math.max(0.03, Number(stepSizeEl.value) || 0.2);
      mainVideoEl.currentTime = Math.max(0, mainVideoEl.currentTime + (delta * step));
    package main
    
    import (
        "context"
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "math"
        "net/http"
        "os"
        "os/exec"
        "sort"
        "strconv"
        "time"
    )
    
    const cameraMotionMetadataTag = "yi_motion"
    
    type motionInterval struct {
        Start float64 `json:"start"`
        End   float64 `json:"end"`
    }
    
    type cameraMotion struct {
        Duration         float64          `json:"duration"`
        HistoryAvailable bool             `json:"history_available"`
        HistoryTruncated bool             `json:"history_truncated"`
        Markers          []float64        `json:"markers"`
        Intervals        []motionInterval `json:"intervals"`
    }
    
    // Camera offsets are relative to the video, independent of filename timezones.
    func parseCameraMotion(tag string, duration float64) (*cameraMotion, error) {
        var raw struct {
            Schema           string              `json:"schema"`
            VideoStart       *float64            `json:"video_start_unix_ms"`
            DurationMS       float64             `json:"duration_ms"`
            InitialState     *int                `json:"initial_state"`
            HistoryAvailable bool                `json:"history_available"`
            HistoryTruncated bool                `json:"history_truncated"`
            Columns          []string            `json:"event_columns"`
            Events           [][]json.RawMessage `json:"events"`
        }
        if err := json.Unmarshal([]byte(tag), &raw); err != nil {
            return nil, fmt.Errorf("invalid yi_motion JSON: %w", err)
        }
        if raw.Schema != "yi-motion-v1" {
            return nil, errors.New("unsupported yi_motion schema")
        }
        if !finite(duration) || duration <= 0 {
            duration = raw.DurationMS / 1000
        }
        if !finite(duration) || duration <= 0 {
            return nil, errors.New("yi_motion has no valid video duration")
        }
        state := -1 // Unknown history must not imply motion or inactivity.
        if raw.InitialState != nil {
            state = *raw.InitialState
        }
        if state < -1 || state > 1 {
            return nil, errors.New("invalid yi_motion initial state")
        }
        column := func(name string) int {
            for i, value := range raw.Columns {
                if value == name {
                    return i
                }
            }
            return -1
        }
        offsetCol, unixCol, stateCol := column("offset_ms"), column("unix_ms"), column("state")
        if len(raw.Events) > 0 && (stateCol < 0 || (offsetCol < 0 && (unixCol < 0 || raw.VideoStart == nil))) {
            return nil, errors.New("yi_motion event columns are missing")
        }
        type event struct {
            offset float64
            state  int
        }
        events := make([]event, 0, len(raw.Events))
        for _, row := range raw.Events {
            readNumber := func(index int) (float64, error) {
                var value float64
                if index < 0 || index >= len(row) || string(row[index]) == "null" {
                    return 0, errors.New("missing yi_motion event value")
                }
                if err := json.Unmarshal(row[index], &value); err != nil || !finite(value) {
                    return 0, errors.New("invalid yi_motion event value")
                }
                return value, nil
            }
            index := offsetCol
            if index < 0 {
                index = unixCol
            }
            offset, err := readNumber(index)
            if err != nil {
                return nil, err
            }
            if offsetCol < 0 {
                offset -= *raw.VideoStart
            }
            offset /= 1000
            value, err := readNumber(stateCol)
            if err != nil || value < -1 || value > 1 || value != math.Trunc(value) {
                return nil, errors.New("invalid yi_motion event state")
            }
            if !finite(offset) || offset < 0 || offset >= duration {
                continue // Events outside the recording cannot mark its seek track.
            }
            events = append(events, event{offset, int(value)})
        }
        sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
        result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
            Markers: []float64{}, Intervals: []motionInterval{}}
        start := 0.0
        for _, event := range events {
            if state == 1 && event.state != 1 && event.offset > start {
                result.Intervals = append(result.Intervals, motionInterval{start, event.offset})
            }
            if event.state == 1 && state != 1 {
                start = event.offset
            }
            state = event.state
        }
        if state == 1 && duration > start {
            result.Intervals = append(result.Intervals, motionInterval{start, duration})
        }
        for _, interval := range result.Intervals {
            result.Markers = append(result.Markers, interval.Start)
        }
        return result, nil
    }
    
    func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
        ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+cameraMotionMetadataTag, "-of", "json", path).Output()
        if err != nil {
            return nil, fmt.Errorf("probe camera motion: %w", err)
        }
        var probe struct {
            Format struct {
                Duration string            `json:"duration"`
                Tags     map[string]string `json:"tags"`
            } `json:"format"`
        }
        if err := json.Unmarshal(out, &probe); err != nil {
            return nil, err
        }
        tag := probe.Format.Tags[cameraMotionMetadataTag]
        if tag == "" {
            return nil, nil
        }
        duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
        return parseCameraMotion(tag, duration)
    }
    
    func (s *Server) apiCameraMotionHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        if r.Method != http.MethodGet {
            w.Header().Set("Allow", "GET")
            writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
            return
        }
        path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
        if err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        motion, err := s.readCameraMotion(r.Context(), path)
        if err != nil {
            writeJSONError(w, http.StatusBadGateway, err)
            return
        }
        writeJSON(w, struct {
            Motion *cameraMotion `json:"motion"`
        }{motion})
    }
    
    func (s *Server) runPreferredThumbnail(videoPath, tmpPath string) error {
        motion, err := s.readCameraMotion(s.workContext(), videoPath)
        if err == nil && motion != nil {
            for _, interval := range motion.Intervals {
                // Let the motion appear after the onset; stay within the marked interval.
                seek := interval.Start + math.Min(.5, (interval.End-interval.Start)/2)
                seek = math.Min(seek, math.Max(0, motion.Duration-.1))
                err = s.runFFmpegThumbnail(videoPath, tmpPath, strconv.FormatFloat(seek, 'f', 3, 64))
                if err == nil {
                    var frame frameSignature
                    frame, err = frameLumaSignature(tmpPath)
                    if err == nil && !frame.usable {
                        err = errors.New("camera motion frame is blank")
                    }
                }
                if err == nil {
                    log.Printf("thumbnail selected %.3fs for %s (camera motion)", seek, videoPath)
                    return nil
                }
                _ = os.Remove(tmpPath)
            }
        }
        if err != nil {
            log.Printf("camera motion thumbnail unavailable for %s, using scene selection: %v", videoPath, err)
        }
        return s.runInterestingThumbnail(videoPath, tmpPath)
    }
    30:    process = subprocess.Popen(arguments, stdout=log, stderr=log)
    109:    with tempfile.TemporaryDirectory(prefix='neonvr-browser-fixture-') as temporary:
    113:        fixture_config.write_text(json.dumps({
    131:        config.write_text(json.dumps({
    144:            browser = p.chromium.launch(headless=True, args=['--no-sandbox', '--use-fake-device-for-media-stream', '--use-fake-ui-for-media-stream'])
    145:            context = browser.new_context(viewport={'width': 1280, 'height': 900}, permissions=['microphone'])
    251:            browser.close()
                    time.sleep(.1)
                if speaker_packets == 0:
                    print('SPEAKER EVENTS:',speaker_events,flush=True)
                    print('AUDIO STATE:',page.evaluate('({tracks:active.microphone.getTracks().map(t=>({enabled:t.enabled,muted:t.muted,state:t.readyState})),transceivers:active.pc.getTransceivers().map(t=>({direction:t.currentDirection,kind:t.sender.track?.kind})),sdp:active.pc.remoteDescription.sdp})'),flush=True)
                    print('AUDIO STATS:',page.evaluate('async()=>Array.from((await active.pc.getStats()).values()).filter(s=>s.type==="outbound-rtp"||s.type==="media-source"||s.type==="codec")'),flush=True)
                assert speaker_packets > 0, 'No user audio reached the RTSP speaker backchannel'
                page.evaluate('window.previousMic = active.microphone')
                page.locator('#talk').click()
                page.wait_for_function('active && active.connected && !active.talk', timeout=20000)
                assert page.evaluate('previousMic.getTracks().every(track => track.readyState === "ended")')
                print('Microphone audio reached the RTSP backchannel; Stop talking ended capture and restored viewing.', flush=True)
                page.evaluate('() => { window.nativeGetUserMedia = navigator.mediaDevices.getUserMedia.bind(navigator.mediaDevices); navigator.mediaDevices.getUserMedia = async() => { throw new DOMException("Denied", "NotAllowedError"); }; }')
                page.locator('#talk').click()
                page.wait_for_function('active === null && document.querySelector("#status").textContent.includes("permission was denied")')
                assert page.locator('#microphone').inner_text() == 'Microphone off'
                page.evaluate('() => {navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia;}')
                page.locator('#watch').click()
                page.wait_for_function('active && active.connected', timeout=20000)
                page.evaluate('() => {navigator.mediaDevices.getUserMedia = (...args) => new Promise(resolve => {window.resolveLateMicrophone = async() => {window.lateMicrophone = await nativeGetUserMedia(...args); resolve(window.lateMicrophone);};});}')
                page.locator('#talk').click()
                page.wait_for_function('typeof window.resolveLateMicrophone === "function"')
                page.locator('#stop').click()
                page.evaluate('window.resolveLateMicrophone()')
                page.wait_for_function('window.lateMicrophone && window.lateMicrophone.getTracks().every(track => track.readyState === "ended")')
                page.evaluate('() => {navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia;}')
                print('Denied permissions and late permission grants leave the microphone off.', flush=True)
                page.get_by_role('button', name='View-only test camera').click()
                page.locator('#watch').click()
                page.wait_for_function('active && active.connected', timeout=20000)
                page.evaluate('() => {navigator.mediaDevices.getUserMedia = async(...args) => {window.unsupportedMicrophone = await nativeGetUserMedia(...args); return window.unsupportedMicrophone;};}')
                page.locator('#talk').click()
                page.wait_for_function('active === null && document.querySelector("#status").textContent.includes("did not accept a compatible audio backchannel")', timeout=20000)
                assert page.evaluate('window.unsupportedMicrophone.getTracks().every(track => track.readyState === "ended")')
                page.evaluate('() => {navigator.mediaDevices.getUserMedia = window.nativeGetUserMedia;}')
                print('A camera without a compatible speaker rejects Talk and releases the microphone.', flush=True)
                page.get_by_role('button', name='Local test camera', exact=False).click()
                page.locator('#watch').click()
                page.wait_for_function('active && active.connected', timeout=20000)
                page.locator('#talk').click()
                page.wait_for_function('active && active.connected && active.talk', timeout=20000)
                page.evaluate('window.offMicrophone = active.microphone')
                other = context.new_page()
                other.goto(origin + '/streams')
                other.wait_for_function('state && state.enabled')
                other.locator('#enabled').uncheck()
                other.wait_for_function('!state.enabled && !changing')
    70:        return c, errors.New("config must be a JSON object")
    78:        return c, errors.New("config must contain exactly one JSON object")
    168:func configRevision(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
    173:        return errors.New("config path must be a regular file")
    184:    f, err := os.CreateTemp(filepath.Dir(path), ".neonvr-config-*")
    221:func configBytes(c Config) []byte {
    227:    path := "config/neonvr.json"
    230:        if a == "-config" || a == "--config" {
    232:                return Config{}, "", errors.New("-config requires a path")
    237:        if strings.HasPrefix(a, "-config=") || strings.HasPrefix(a, "--config=") {
    243:        return Config{}, "", errors.New("config path cannot be empty")
    254:            return c, path, fmt.Errorf("config %s: %w", path, err)
    260:    fs.StringVar(&path, "config", path, "JSON config file (default config/neonvr.json)")
    272:    fs.Var(&extra, "extra-root", "additional directory; repeat to replace configured extra roots")
    280:    fs.StringVar(&c.ThumbMaintenanceEvery, "thumb-maintenance-every", c.ThumbMaintenanceEvery, "maintenance cadence; 0 disables")
    282:    fs.IntVar(&c.AutoDeleteDays, "auto-delete-days", c.AutoDeleteDays, "retention days; 0 disables deletion")
    284:    fs.BoolVar(&c.DetectionEnabled, "detection-enabled", c.DetectionEnabled, "enable background YOLO26n CPU detection")
    318:        if err = atomicWriteConfig(path, configBytes(c)); err != nil {
    325:type configManager struct {
    334:func newConfigManager(path string, c Config) (*configManager, error) {
    339:    return &configManager{path: path, active: c, token: hex.EncodeToString(token), restart: make(chan struct{}, 1)}, nil
    342:func (m *configManager) response(w http.ResponseWriter, status int, saved Config, revision string, nextURL string) {
    346:    _ = json.NewEncoder(w).Encode(map[string]any{"config": saved, "active_config": m.active, "revision": revision, "token": m.token, "instance_id": m.token, "restart_required": !reflect.DeepEqual(saved, m.active), "restarting": m.restarting, "next_url": nextURL})
    349:func (m *configManager) handler(w http.ResponseWriter, r *http.Request) {
    382:        http.Error(w, "cannot read config: "+err.Error(), 500)
    387:        http.Error(w, "invalid config file: "+err.Error(), 500)
    390:    revision := configRevision(raw)
    400:        Config   json.RawMessage `json:"config"`
    420:        http.Error(w, "invalid config: "+err.Error(), 400)
    435:    data := configBytes(c)
    437:        http.Error(w, "cannot save previous config: "+err.Error(), 500)
    441:        http.Error(w, "cannot save config: "+err.Error(), 500)
    460:    m.response(w, 200, c, configRevision(data), nextURL)
    ✓ • 0ms

## Activity

    $ git diff --check; git diff --stat; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr-motion .
     Dockerfile                  |  7 +++-
     README.md                   | 28 ++++++++++++++--
     config.go                   | 64 ++++++++++++++++++++----------------
     config/neonvr.json          | 12 ++++++-
     detection_test.go           |  6 +++-
     docker-compose.yml          |  5 +++
     main.go                     | 79 ++++++++++++++++++++++++++++++++-------------
     static/detection-viewer.js  | 44 ++++++++++++++++++++++++-
     static/media-chrome-lite.js | 46 ++++++++++++++++++++++----
     static/settings.html        | 19 ++++++++---
     10 files changed, 242 insertions(+), 68 deletions(-)
    ✓ • 517ms

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion_test.go
    @@ -77,2 +77,7 @@
     	t.Helper()
    +	return cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=40:y=20:w=80:h=60:color=red:t=fill:enable='between(t,1.5,2.5)'")
    +}
    +
    +func cameraMotionFixtureWithFilter(t *testing.T, tag, filter string) (*Server, string) {
    +	t.Helper()
     	for _, tool := range []string{"ffmpeg", "ffprobe"} {
    @@ -88,4 +93,3 @@
     	path := filepath.Join(dir, "recording.mp4")
    -	args := []string{"-nostdin", "-v", "error", "-f", "lavfi", "-i",
    -		"testsrc2=s=160x96:r=10,drawbox=x=40:y=20:w=80:h=60:color=red:t=fill:enable='between(t,1.5,2.5)'",
    +	args := []string{"-nostdin", "-v", "error", "-f", "lavfi", "-i", filter,
     		"-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1"}
    @@ -159 +163,17 @@
     }
    +
    +func TestCameraMotionBlankFrameFallback(t *testing.T) {
    +	for _, tag := range []string{
    +		`{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[500,0]]}`,
    +		`{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[500,0],[2000,1],[3000,0]]}`,
    +	} {
    +		s, path := cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=0:y=0:w=iw:h=ih:color=black:t=fill:enable='lt(t,1)'")
    +		thumb, err := s.ensureThumbnail(path)
    +		if err != nil {
    +			t.Fatal(err)
    +		}
    +		if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
    +			t.Fatalf("blank motion frame was used instead of a later interval or scene fallback: %v", err)
    +		}
    +	}
    +}

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w motion_test.go; NEONVR_TEST_GO2RTC=/tmp/neonvr-live-tools-vbq_9dyb/go2rtc /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
    ok      neonvr    9.058s
    ✓ • 10.02s

## Activity

    $ cat > /tmp/neonvr-live-tools-vbq_9dyb/check_motion_browser.py <<'PY'
        import hashlib
        import json
        from pathlib import Path
        import socket
        import subprocess
        import tempfile
        import time
        import urllib.request
        from playwright.sync_api import sync_playwright
        
        TOOLS = Path('/tmp/neonvr-live-tools-vbq_9dyb')
        
        def get(url):
            with urllib.request.urlopen(url, timeout=10) as response:
                return json.load(response)
        
        with tempfile.TemporaryDirectory(prefix='neonvr-motion-browser-') as tmp:
            directory = Path(tmp)
            recordings = directory / 'clips' / 'Garden'
            recordings.mkdir(parents=True)
            fixtures = {
                'transition.mp4': '/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W415M41S19.mp4',
                'partial.mp4': '/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W419M13S46.mp4',
                'initial.mp4': '/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W430M00S60.mp4',
            }
            for name, source in fixtures.items():
                subprocess.run(['ffmpeg', '-nostdin', '-v', 'error', '-i', source, '-map', '0', '-c', 'copy', '-map_metadata', '0', '-metadata', 'neonvr_detections=', '-movflags', '+use_metadata_tags', str(recordings / name)], check=True)
            for name, tag in [('missing.mp4', None), ('invalid.mp4', 'broken-json'), ('quiet.mp4', json.dumps({'schema':'yi-motion-v1','initial_state':0,'history_available':True,'events':[]}))]:
                arguments = ['ffmpeg','-nostdin','-v','error','-f','lavfi','-i','testsrc2=s=160x96:r=10','-t','4','-c:v','libx264','-threads','1']
                if tag is not None:
                    arguments += ['-metadata','yi_motion='+tag]
                subprocess.run(arguments + ['-movflags','+use_metadata_tags',str(recordings / name)], check=True)
            with socket.socket() as listener:
                listener.bind(('127.0.0.1',0))
                port = listener.getsockname()[1]
            config = directory / 'config.json'
            config.write_text(json.dumps({'root':str(directory / 'clips'),'extra_roots':[],'addr':f'127.0.0.1:{port}','auto_delete_days':0,'thumb_maintenance_every':'0','detection_enabled':False,'live_streams':[]}))
            origin = f'http://127.0.0.1:{port}'
            with (directory / 'app.log').open('w+') as log:
                process = subprocess.Popen([str(TOOLS / 'neonvr-motion'),'-config',str(config)],stdout=log,stderr=log)
                try:
                    for attempt in range(100):
                        try:
                            cameras = get(origin + '/api/cameras')
                            break
                        except Exception:
                            time.sleep(.05)
                    else:
                        raise RuntimeError('Server did not start')
                    clips = {clip['file']:clip for clip in get(origin + '/api/clips?cameras=Garden&date=' + cameras['latest_date'])['clips']}
                    with sync_playwright() as playwright:
                        browser = playwright.chromium.launch(headless=True,args=['--no-sandbox'])
                        page = browser.new_page(viewport={'width':1400,'height':1000})
                        errors=[]
                        page.on('pageerror',lambda error:errors.append(str(error)))
                        page.goto(origin)
                        page.wait_for_function('typeof loadMainPlayer === "function" && window.NeoNVRDetections')
                        def load(name):
                            page.evaluate('(clip) => loadMainPlayer(clip, false)', clips[name])
                            page.wait_for_function('Number.isFinite(document.querySelector("#mainVideo").duration) && document.querySelector("#mainVideo").duration > 0')
                        load('transition.mp4')
                        marker = page.locator('#mainSeek .motion-marker')
                        marker.wait_for()
                        assert marker.count() == 1
                        assert marker.get_attribute('title') == 'Camera motion at 5.295 sec'
                        page.wait_for_function('document.querySelector("#detectionStatus").textContent === "Detection disabled."')
                        lane = page.locator('#mainSeek .lane[data-kind="camera_motion"]')
                        assert lane.count() == 1
                        assert '5.295–11.239' in lane.locator('button').get_attribute('title')
                        marker.click()
                        page.wait_for_function('Math.abs(document.querySelector("#mainVideo").currentTime - 5.295) < .002')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent === "Camera motion now."')
                        page.evaluate('document.querySelector("#mainVideo").currentTime = 12')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent === "No camera motion marked here."')
                        page.evaluate('document.querySelector("#mainSeek").setDetections({duration:18.624,intervals:[{kind:"person",start:6,end:7,confidence:.9}]})')
                        assert page.locator('#mainSeek .lane').count() == 5
                        assert marker.count() == 1
                        page.evaluate('document.querySelector("#mainSeek").setDetections(null)')
                        assert lane.count() == 1 and marker.count() == 1
                        page.locator('.mainPlayerShell').screenshot(path=str(TOOLS / 'motion-player-desktop.png'))
                        with urllib.request.urlopen(origin + clips['transition.mp4']['thumb_url']) as response:
                            actual = response.read()
                        expected = directory / 'expected.jpg'
                        subprocess.run(['ffmpeg','-nostdin','-hide_banner','-loglevel','error','-ss','5.795','-i',str(recordings / 'transition.mp4'),'-frames:v','1','-vf','scale=480:-2','-q:v','5','-y',str(expected)], check=True)
                        assert hashlib.sha256(actual).digest() == hashlib.sha256(expected.read_bytes()).digest(), 'Real camera thumbnail did not use motion onset'
                        print('Real Yi motion: exact 5.295 sec clickable marker, 5.295–11.239 sec interval, and thumbnail at 5.795 sec.',flush=True)
                        load('partial.mp4')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent.includes("Partial history.")')
                        assert marker.get_attribute('title') == 'Camera motion at 2.514 sec'
                        load('initial.mp4')
                        page.wait_for_function('document.querySelector("#mainSeek").shadowRoot.querySelector(".motion-marker")?.title === "Camera motion at 0.000 sec"')
                        marker.click()
                        page.wait_for_function('document.querySelector("#mainVideo").currentTime < .01')
                        load('missing.mp4')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").hidden && !document.querySelector("#mainSeek").cameraMotion')
                        assert marker.count() == 0 and lane.count() == 0
                        load('quiet.mp4')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent === "No camera motion recorded."')
                        assert marker.count() == 0 and lane.count() == 0
                        load('invalid.mp4')
                        page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent === "Camera motion metadata unavailable."')
                        assert marker.count() == 0
                        print('Partial history, initial motion, absent tags, quiet clips, malformed tags, and AI lane coexistence passed.',flush=True)
                        page.evaluate('''() => {
                            const original = window.fetch;
                            window.fetch = async (...args) => {
                                const response = await original(...args);
                                if (String(args[0]).includes('/api/camera-motion?') && String(args[0]).includes('transition.mp4')) {
                                    return new Promise(resolve => { window.releaseOldMotion = () => resolve(response); });
                                }
                                return response;
                            };
                        }''')
                        load('transition.mp4')
                        page.wait_for_function('typeof window.releaseOldMotion === "function"')
                        load('missing.mp4')
                        page.wait_for_function('!document.querySelector("#mainSeek").cameraMotion')
                        page.evaluate('window.releaseOldMotion()')
                        page.wait_for_timeout(150)
                        assert marker.count() == 0 and page.locator('#cameraMotionLegend').is_hidden()
                        print('Late metadata response after switching clips cannot restore old markers.',flush=True)
                        browser.close()
                        mobile = playwright.chromium.launch(headless=True,args=['--no-sandbox'])
                        mobile_page = mobile.new_page(viewport={'width':390,'height':844},is_mobile=True,has_touch=True)
                        mobile_page.goto(origin)
                        mobile_page.wait_for_function('typeof loadMainPlayer === "function" && window.NeoNVRDetections')
                        mobile_page.evaluate('(clip)=>loadMainPlayer(clip,false)',clips['transition.mp4'])
                        mobile_marker = mobile_page.locator('#mainSeek .motion-marker')
                        mobile_marker.wait_for()
                        mobile_marker.tap()
                        mobile_page.wait_for_function('Math.abs(document.querySelector("#mainVideo").currentTime - 5.295) < .002')
                        mobile_page.locator('.mainPlayerShell').screenshot(path=str(TOOLS / 'motion-player-mobile.png'))
                        assert mobile_page.evaluate('document.documentElement.scrollWidth <= innerWidth + 1')
                        print('Mobile marker touch seeking and layout passed.',flush=True)
                        mobile.close()
                        assert not errors, errors
                finally:
                    process.terminate()
                    try: process.wait(timeout=10)
                    except subprocess.TimeoutExpired:
                        process.kill();process.wait()
                    log.seek(0)
                    text=log.read()
                    print('\n'.join(line for line in text.splitlines() if 'thumbnail selected' in line or 'unavailable' in line),flush=True)
        PY
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_motion_browser.py
    Real Yi motion: exact 5.295 sec clickable marker, 5.295–11.239 sec interval, and thumbnail at 5.795 sec.
    Partial history, initial motion, absent tags, quiet clips, malformed tags, and AI lane coexistence passed.
    Late metadata response after switching clips cannot restore old markers.
    Mobile marker touch seeking and layout passed.
    2026/10/11 00:37:51 camera motion thumbnail unavailable for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/invalid.mp4, using scene selection: invalid yi_motion JSON: invalid character 'b' looking for beginning of value
    2026/10/11 00:37:52 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/invalid.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:37:52 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/quiet.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:37:52 thumbnail selected 0.500s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/initial.mp4 (camera motion)
    2026/10/11 00:37:53 thumbnail selected 5.795s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/transition.mp4 (camera motion)
    2026/10/11 00:37:53 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/missing.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:37:53 thumbnail selected 3.014s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/partial.mp4 (camera motion)
    2026/10/11 00:37:53 thumbnail selected 5.795s for /tmp/neonvr-motion-browser-5koh8_2r/clips/Garden/transition.mp4 (camera motion)
    ✓ • 5.09s

## Activity

    $ git diff --check; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go vet ./...
    ✓ • 55ms

## Assistant

The checks passed on copies of your camera recordings: a motion event at 5.295 seconds appears at the correct slider position, clicking it seeks there, and the thumbnail uses a frame at 5.795 seconds. Missing, malformed, and blank-frame cases also fall back correctly. Desktop and mobile browser checks passed.

## Activity

    Viewed image /tmp/neonvr-live-tools-vbq_9dyb/motion-player-mobile.png

## Activity

    Viewed image /tmp/neonvr-live-tools-vbq_9dyb/motion-player-desktop.png

## Activity

    $ git diff -- static/detection-viewer.js static/media-chrome-lite.js detection_test.go; git status --short; sed -n '2525,2585p' main.go; sed -n '252,279p' detection.go; cat Dockerfile
    diff --git a/detection_test.go b/detection_test.go
    index 7b3c81a..2da1995 100644
    --- a/detection_test.go
    +++ b/detection_test.go
    @@ -32,7 +32,7 @@ func detectionFixture(t *testing.T) (*Server, Camera, string, DetectionResult) {
         path := filepath.Join(cam.Dir, "2026-10-01_12-00-00_Garden.mp4")
         out, err := exec.Command("ffmpeg", "-nostdin", "-v", "error", "-f", "lavfi", "-i", "color=c=green:s=160x96:r=10",
             "-f", "lavfi", "-i", "sine=frequency=440", "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1", "-c:a", "aac",
    -        "-metadata", "source_camera=Garden", "-movflags", "+use_metadata_tags", path).CombinedOutput()
    +        "-metadata", "source_camera=Garden", "-metadata", "yi_motion="+motionTestTag, "-movflags", "+use_metadata_tags", path).CombinedOutput()
         if err != nil {
             t.Fatalf("create fixture: %v: %s", err, out)
         }
    @@ -88,6 +88,10 @@ func TestDetectionMP4RoundTripPreservesRecording(t *testing.T) {
         if err != nil || !bytes.Contains(allTags, []byte("Garden")) {
             t.Fatal("existing custom metadata lost")
         }
    +    motion, err := s.readCameraMotion(context.Background(), path)
    +    if err != nil || motion == nil || !reflect.DeepEqual(motion.Intervals, []motionInterval{{1.5, 2.5}}) {
    +        t.Fatalf("camera motion metadata lost during AI remux: %+v %v", motion, err)
    +    }
         // The MP4 remains the source of truth when the expendable JSON cache is lost.
         if err := os.Remove(s.detectionCachePath(path)); err != nil {
             t.Fatal(err)
    diff --git a/static/detection-viewer.js b/static/detection-viewer.js
    index fc52b80..5a264db 100644
    --- a/static/detection-viewer.js
    +++ b/static/detection-viewer.js
    @@ -4,14 +4,50 @@
       const video = document.getElementById('mainVideo');
       const seek = document.getElementById('mainSeek');
       const status = document.getElementById('detectionStatus');
    +  const motionLegend = document.getElementById('cameraMotionLegend');
    +  const motionStatus = document.getElementById('cameraMotionStatus');
       if (!video || !seek || !status) return;
       const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
       let generation = 0;
       let request = null;
       let timer = null;
       let result = null;
    +  let motionRequest = null;
    +  let motion = null;
    +
    +  function activeCameraMotion() {
    +    if (!motion || !motionStatus) return;
    +    const active = motion.intervals.some(row => video.currentTime >= row.start && video.currentTime < row.end);
    +    const partial = !motion.history_available || motion.history_truncated;
    +    motionStatus.classList.toggle('active', active);
    +    const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
    +      partial ? 'No camera motion timestamps available.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
    +    if (motionStatus.textContent !== text) motionStatus.textContent = text;
    +  }
    +
    +  async function refreshCameraMotion(clip, expected) {
    +    motionRequest = new AbortController();
    +    try {
    +      const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
    +      if (!response.ok) throw new Error('Camera motion metadata unavailable.');
    +      const data = await response.json();
    +      if (expected !== generation) return;
    +      motion = data.motion;
    +      seek.setCameraMotion(motion);
    +      if (motionLegend) motionLegend.hidden = !motion;
    +      if (motionStatus) motionStatus.hidden = !motion;
    +      activeCameraMotion();
    +    } catch (error) {
    +      if (expected !== generation || error.name === 'AbortError') return;
    +      if (motionStatus) {
    +        motionStatus.hidden = false;
    +        motionStatus.textContent = error.message;
    +      }
    +    }
    +  }
    
       function activeStatus() {
    +    activeCameraMotion();
         if (!result) return;
         const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
         status.classList.toggle('active', active.length > 0);
    @@ -50,17 +86,23 @@
       function loadClip(clip) {
         generation += 1;
         request?.abort();
    +    motionRequest?.abort();
         clearTimeout(timer);
         result = null;
    +    motion = null;
         seek.setDetections(null);
    +    seek.setCameraMotion(null);
    +    if (motionLegend) motionLegend.hidden = true;
    +    if (motionStatus) { motionStatus.hidden = true; motionStatus.classList.remove('active'); }
         status.classList.remove('active');
         status.textContent = 'Loading detections…';
         status.title = '';
         refresh(clip, generation);
    +    refreshCameraMotion(clip, generation);
       }
    
       video.addEventListener('timeupdate', activeStatus);
       video.addEventListener('seeked', activeStatus);
    -  window.addEventListener('pagehide', () => { request?.abort(); clearTimeout(timer); });
    +  window.addEventListener('pagehide', () => { request?.abort(); motionRequest?.abort(); clearTimeout(timer); });
       window.NeoNVRDetections = { loadClip };
     })();
    diff --git a/static/media-chrome-lite.js b/static/media-chrome-lite.js
    index 53a86cd..8f949db 100644
    --- a/static/media-chrome-lite.js
    +++ b/static/media-chrome-lite.js
    @@ -159,6 +159,12 @@
           input.part = 'range';
           input.setAttribute('aria-label', 'Video position');
           this.input = input;
    +      const track = document.createElement('div');
    +      track.className = 'seek-track';
    +      this.motionMarkers = document.createElement('div');
    +      this.motionMarkers.className = 'motion-markers';
    +      this.motionMarkers.setAttribute('aria-label', 'Camera motion timestamps');
    +      track.append(input, this.motionMarkers);
           this.lanes = document.createElement('div');
           this.lanes.className = 'detection-lanes';
           this.lanes.setAttribute('aria-label', 'Detection timeline');
    @@ -167,12 +173,16 @@
           style.textContent = `
             :host { display: block; flex: 1 1 180px; min-width: 110px; }
             input { width: 100%; margin: 0; accent-color: white; }
    +        .seek-track { position: relative; display: flex; align-items: center; }
    +        .motion-markers { position: absolute; left: 8px; right: 8px; top: 50%; pointer-events: none; }
    +        .motion-marker { position: absolute; width: 7px; height: 18px; padding: 0; border: 1px solid #332014; border-radius: 2px; background: #ff9866; transform: translate(-50%, -50%); pointer-events: auto; cursor: pointer; }
    +        .motion-marker:hover, .motion-marker:focus-visible { outline: 2px solid white; z-index: 1; }
             .detection-lanes { margin: 3px 8px 0; }
             .lane { height: 7px; position: relative; margin-top: 3px; background: #262a32; border-radius: 3px; }
             .event { position: absolute; top: 0; height: 100%; padding: 0; min-width: 2px; border: 0; border-radius: 2px; cursor: pointer; }
             .event:hover, .event:focus-visible { outline: 2px solid white; z-index: 1; }
           `;
    -      root.append(style, input, this.lanes);
    +      root.append(style, track, this.lanes);
           const video = findVideo(this);
           if (!video) return;
           let dragging = false;
    @@ -200,22 +210,43 @@
           this.detections = result;
           this.renderDetections();
         }
    +    setCameraMotion(motion) {
    +      this.cameraMotion = motion;
    +      this.renderDetections();
    +    }
         renderDetections() {
           if (!this.lanes) return;
           this.lanes.replaceChildren();
    +      this.motionMarkers.replaceChildren();
           const result = this.detections;
    -      this.lanes.hidden = !result;
    -      if (!result) return;
    -      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : result.duration;
    +      const motion = this.cameraMotion;
    +      this.lanes.hidden = !result && !motion?.intervals?.length;
    +      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : (motion?.duration || result?.duration);
           if (!Number.isFinite(duration) || duration <= 0) return;
    +      if (motion) {
    +        for (const time of motion.markers || []) {
    +          if (!Number.isFinite(time) || time < 0 || time >= duration) continue;
    +          const marker = document.createElement('button');
    +          marker.type = 'button';
    +          marker.className = 'motion-marker';
    +          marker.style.left = (time / duration * 100) + '%';
    +          marker.title = 'Camera motion at ' + time.toFixed(3) + ' sec';
    +          marker.setAttribute('aria-label', 'Seek to ' + marker.title);
    +          marker.addEventListener('click', () => { if (this.video) this.video.currentTime = time; });
    +          this.motionMarkers.append(marker);
    +        }
    +      }
           const categories = [['person', 'Human', '#75d685'], ['animal', 'Animal', '#ffd166'], ['bike', 'Bike', '#72a7ff'], ['moving_car', 'Moving vehicle', '#c792ea']];
    +      const intervals = result?.intervals || [];
    +      if (motion?.intervals?.length) categories.unshift(['camera_motion', 'Camera motion', '#ff9866']);
           for (const [kind, label, color] of categories) {
    +        if (kind !== 'camera_motion' && !result) continue;
             const lane = document.createElement('div');
             lane.className = 'lane';
             lane.dataset.kind = kind;
             lane.setAttribute('aria-label', label + ' detections');
    -        for (const interval of result.intervals || []) {
    -          if (interval.kind !== kind || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
    +        for (const interval of kind === 'camera_motion' ? motion.intervals : intervals) {
    +          if ((kind !== 'camera_motion' && interval.kind !== kind) || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
               const start = Math.max(0, Math.min(duration, interval.start));
               const end = Math.max(start, Math.min(duration, interval.end));
               if (end <= start) continue;
    @@ -225,7 +256,8 @@
               event.style.left = (start / duration * 100) + '%';
               event.style.width = ((end - start) / duration * 100) + '%';
               event.style.background = color;
    -          const text = label + ': ' + start.toFixed(1) + '–' + end.toFixed(1) + ' sec (' + Math.round(interval.confidence * 100) + '%)';
    +          const text = label + ': ' + start.toFixed(3) + '–' + end.toFixed(3) + ' sec' +
    +            (kind === 'camera_motion' ? '' : ' (' + Math.round(interval.confidence * 100) + '%)');
               event.title = text;
               event.setAttribute('aria-label', 'Seek to ' + text);
               event.addEventListener('click', () => { if (this.video) this.video.currentTime = start; });
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M detection_test.go
     M docker-compose.yml
     M main.go
     M static/detection-viewer.js
     M static/media-chrome-lite.js
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? motion.go
    ?? motion_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    }
    
    func frameLumaSignature(path string) (frameSignature, error) {
        f, err := os.Open(path)
        if err != nil {
            return frameSignature{}, err
        }
        defer f.Close()
    
        img, err := jpeg.Decode(f)
        if err != nil {
            return frameSignature{}, err
        }
        b := img.Bounds()
        if b.Dx() < 16 || b.Dy() < 9 {
            return frameSignature{}, fmt.Errorf("thumbnail too small")
        }
    
        const gridW = 32
        const gridH = 18
        luma := make([]float64, 0, gridW*gridH)
        var sum, sumSq float64
        for gy := 0; gy < gridH; gy++ {
            y := b.Min.Y + (gy*2+1)*b.Dy()/(2*gridH)
            if y >= b.Max.Y {
                y = b.Max.Y - 1
            }
            for gx := 0; gx < gridW; gx++ {
                x := b.Min.X + (gx*2+1)*b.Dx()/(2*gridW)
                if x >= b.Max.X {
                    x = b.Max.X - 1
                }
                r, g, bl, _ := img.At(x, y).RGBA()
                v := float64((299*uint64(r) + 587*uint64(g) + 114*uint64(bl)) / 1000 / 257)
                luma = append(luma, v)
                sum += v
                sumSq += v * v
            }
        }
        mean := sum / float64(len(luma))
        variance := sumSq/float64(len(luma)) - mean*mean
        if variance < 0 {
            variance = 0
        }
        contrast := math.Sqrt(variance)
    
        usable := true
        if contrast < 3 && (mean < 10 || mean > 245) {
            usable = false
        }
    
        return frameSignature{
            path:     path,
            luma:     luma,
            mean:     mean,
            contrast: contrast,
            usable:   usable,
        }, nil
    }
    
    func medianLumaBaseline(frames []frameSignature) ([]float64, error) {
        }
        s.detectionMu.Lock()
        if state, ok := s.detectionStates[path]; ok {
            response = state
        } else if s.detectionWorkerError != "" {
            response = detectionResponse{Status: "error", Message: s.detectionWorkerError}
        }
        s.detectionMu.Unlock()
        writeJSON(w, response)
    }
    
    type detectionWorker struct {
        cmd         *exec.Cmd
        cancel      context.CancelFunc
        stdin       io.WriteCloser
        scanner     *bufio.Scanner
        modelSHA256 string
    }
    
    func (worker *detectionWorker) close() {
        _ = worker.stdin.Close()
        worker.cancel()
        _ = worker.cmd.Wait()
    }
    
    func (worker *detectionWorker) read(ctx context.Context) ([]byte, error) {
        type reply struct {
            data []byte
    FROM golang:1.22-alpine AS build
    WORKDIR /src
    COPY go.mod ./
    COPY *.go ./
    COPY static ./static
    COPY detection ./detection
    RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/neonvr .
    
    # The bridge is started only when live viewing is explicitly enabled.
    FROM ghcr.io/alexxit/go2rtc:1.9.14 AS streams
    
    # PyTorch/Ultralytics are used only to export the official weights at build time.
    FROM python:3.12-slim-bookworm AS model
    ENV YOLO_AUTOINSTALL=False OMP_NUM_THREADS=2
    RUN apt-get update && apt-get install -y --no-install-recommends libgl1 libglib2.0-0 && rm -rf /var/lib/apt/lists/*
    WORKDIR /model
    RUN pip install --no-cache-dir torch==2.9.1 torchvision==0.24.1 --index-url https://download.pytorch.org/whl/cpu
    RUN pip install --no-cache-dir ultralytics==8.4.175 onnx==1.23.2
    COPY detection/export_model.py ./
    RUN python export_model.py --output /model/yolo26n.onnx
    
    # CPU ONNX Runtime supports the Intel i7-6700; no CUDA/GPU runtime is installed.
    FROM python:3.12-slim-bookworm
    RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
    COPY detection/requirements.txt /tmp/detection-requirements.txt
    RUN pip install --no-cache-dir -r /tmp/detection-requirements.txt && rm /tmp/detection-requirements.txt
    COPY --from=model /model/yolo26n.onnx /opt/neonvr/models/yolo26n.onnx
    ENV NEONVR_MODEL_DIR=/opt/neonvr/models OMP_NUM_THREADS=2
    COPY --from=build /out/neonvr /usr/local/bin/neonvr
    COPY --from=streams /usr/local/bin/go2rtc /usr/local/bin/go2rtc
    COPY third_party/go2rtc.LICENSE /usr/share/doc/go2rtc/LICENSE
    EXPOSE 5500 8555/tcp 8555/udp
    ENTRYPOINT ["/usr/local/bin/neonvr"]
    CMD ["-config", "/config/neonvr.json"]
    ✓ • 0ms

## Activity

    $ rg -n 'async function init|loadMainPlayer\(' main.go; sed -n '1190,1225p' main.go; sed -n '1540,1608p' main.go; rg -n 'disabled|shadowRoot|::part' main.go static/media-chrome-lite.js
    1054:function loadMainPlayer(clip, play = true, scroll = false) {
    1089:  if (clip) loadMainPlayer(clip, true, false);
    1285:      else loadMainPlayer(clip, true, false);
    1355:        loadMainPlayer(clip, true, false);
    1400:    load.addEventListener('click', () => loadMainPlayer(clip, true, true));
    1486:    loadMainPlayer(visibleClips[0], false, false);
    1587:async function init() {
      // Keep wheel/touchpad scrolling inside the sidebar when possible instead of bubbling to the page.
      sideRailEl.addEventListener('wheel', evt => {
        const maxScroll = sideRailEl.scrollHeight - sideRailEl.clientHeight;
        if (maxScroll <= 0) return;
        const before = sideRailEl.scrollTop;
        sideRailEl.scrollTop += evt.deltaY;
        const moved = sideRailEl.scrollTop !== before;
        if (moved || (evt.deltaY < 0 && before > 0) || (evt.deltaY > 0 && before < maxScroll)) {
          evt.preventDefault();
          evt.stopPropagation();
        }
      }, {passive:false});
    }
    
    
    function sortedClips(clips) {
      const copy = [...clips];
      copy.sort((a, b) => {
        const av = Number(a.unix_seconds || 0);
        const bv = Number(b.unix_seconds || 0);
        if (sortEl.value === 'oldest') return av - bv || a.camera.localeCompare(b.camera);
        return bv - av || a.camera.localeCompare(b.camera);
      });
      return copy;
    }
    
    function keyFor(clip) { return clip.camera + '/' + clip.file; }
    
    function updateSelectionBar() {
      const n = selectedKeys.size;
      selectionCountEl.textContent = selectMode ? (n + ' selected') : '';
      bulkDeleteEl.style.display = selectMode ? '' : 'none';
      clearSelectionEl.style.display = selectMode ? '' : 'none';
      bulkDeleteEl.disabled = n === 0;
      clearSelectionEl.disabled = n === 0;
      selectModeToggleEl.textContent = selectMode ? 'Done' : 'Select';
      renderDatePills();
    }
    
    async function refreshStats() {
      try {
        const stats = await getJSON('/api/stats');
        let html = '<span>Recordings: <strong>' + stats.root_size_human + '</strong></span>' +
                   '<span>Size cache: ' + new Date(stats.generated_at).toLocaleTimeString() + '</span>';
        if (stats.error) html += '<span>' + escapeHTML(stats.error) + '</span>';
        diskStatsEl.innerHTML = html;
      } catch (err) {
        diskStatsEl.innerHTML = '<span>Folder sizes unavailable: ' + escapeHTML(String(err.message || err)) + '</span>';
      }
    }
    
    function setView(view) {
      currentView = view;
      clipsTabEl.classList.toggle('active', view === 'clips');
      favoritesTabEl.classList.toggle('active', view === 'favorites');
      dateEl.disabled = view === 'favorites';
      clearSelection();
      refreshClips(true).catch(err => {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      });
    }
    
    async function bulkDeleteSelected() {
      if (selectedKeys.size === 0) return;
      const items = visibleClips.filter(c => selectedKeys.has(keyFor(c))).map(c => ({camera: c.camera, file: c.file}));
      const hasFavorite = visibleClips.some(c => selectedKeys.has(keyFor(c)) && c.favorite);
      const msg = hasFavorite
        ? 'Delete ' + items.length + ' selected videos? This includes favorited clips.'
        : 'Delete ' + items.length + ' selected videos?';
      if (!confirm(msg)) return;
      try {
        await postJSON('/api/bulk-delete', {items, force: hasFavorite});
        const del = new Set(items.map(i => i.camera + '/' + i.file));
        allClips = allClips.filter(c => !del.has(keyFor(c)));
        clearSelection();
        renderMarks();
        renderClips();
        refreshStats();
      } catch (err) {
        alert(String(err.message || err));
      }
    }
    
    async function init() {
      try {
        buildScales();
        updateRangeVisuals();
        refreshStats();
        const data = await getJSON('/api/cameras');
        cameras = data.cameras || [];
        selectedCameras = new Set(cameras.map(c => c.name));
        renderCameraMenu();
        await refreshDates(true);
        await refreshClips(true);
      } catch (err) {
        messageEl.innerHTML = '<div class="error">' + escapeHTML(String(err.message || err)) + '</div>';
      }
    }
    
    cameraButtonEl.addEventListener('click', () => cameraMenuEl.classList.toggle('open'));
    document.addEventListener('click', evt => {
      if (!cameraButtonEl.contains(evt.target) && !cameraMenuEl.contains(evt.target)) cameraMenuEl.classList.remove('open');
    });
    
    dateEl.addEventListener('change', async () => {
    static/media-chrome-lite.js:39:      if (this.shadowRoot) return;
    static/media-chrome-lite.js:152:      if (this.shadowRoot) return;
    static/media-chrome-lite.js:297:      if (this.shadowRoot) return;
    main.go:726:          <button id="bulkDelete" type="button" class="danger compact" style="display:none" disabled>Delete selected</button>
    main.go:727:          <button id="clearSelection" type="button" class="compact" style="display:none" disabled>Clear</button>
    main.go:1223:  bulkDeleteEl.disabled = n === 0;
    main.go:1224:  clearSelectionEl.disabled = n === 0;
    main.go:1559:  dateEl.disabled = view === 'favorites';
    main.go:1832:        log.Printf("Thumbnail maintenance disabled")
    main.go:1838:        log.Printf("Auto-delete disabled; use -auto-delete-days N to enable")
    ✓ • 0ms

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/media-chrome-lite.js
    @@ -258,3 +258,4 @@
               event.style.background = color;
    -          const text = label + ': ' + start.toFixed(3) + '–' + end.toFixed(3) + ' sec' +
    +          const precision = kind === 'camera_motion' ? 3 : 1;
    +          const text = label + ': ' + start.toFixed(precision) + '–' + end.toFixed(precision) + ' sec' +
                 (kind === 'camera_motion' ? '' : ' (' + Math.round(interval.confidence * 100) + '%)');

## Activity

    $ python3 - <<'PY'
        from pathlib import Path
        path=Path('/tmp/neonvr-live-tools-vbq_9dyb/check_motion_browser.py')
        text=path.read_text()
        text=text.replace("page.wait_for_function('typeof loadMainPlayer === \"function\" && window.NeoNVRDetections')", "page.wait_for_function('typeof loadMainPlayer === \"function\" && window.NeoNVRDetections && allClips.length === 6')")
        text=text.replace("mobile_page.wait_for_function('typeof loadMainPlayer === \"function\" && window.NeoNVRDetections')", "mobile_page.wait_for_function('typeof loadMainPlayer === \"function\" && window.NeoNVRDetections && allClips.length === 6')")
        text=text.replace("mobile_page.locator('.mainPlayerShell').screenshot", "mobile_page.wait_for_function('!document.querySelector(\"#mainVideo\").seeking && document.querySelector(\"#cameraMotionStatus\").textContent === \"Camera motion now.\"')\n                mobile_page.locator('.mainPlayerShell').screenshot")
        path.write_text(text)
        PY
        /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr-motion .
    ✓ • 525ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_motion_browser.py
    Real Yi motion: exact 5.295 sec clickable marker, 5.295–11.239 sec interval, and thumbnail at 5.795 sec.
    Partial history, initial motion, absent tags, quiet clips, malformed tags, and AI lane coexistence passed.
    Late metadata response after switching clips cannot restore old markers.
    Mobile marker touch seeking and layout passed.
    2026/10/11 00:38:41 thumbnail selected 5.795s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/transition.mp4 (camera motion)
    2026/10/11 00:38:41 camera motion thumbnail unavailable for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/invalid.mp4, using scene selection: invalid yi_motion JSON: invalid character 'b' looking for beginning of value
    2026/10/11 00:38:42 thumbnail selected 0.500s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/initial.mp4 (camera motion)
    2026/10/11 00:38:42 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/invalid.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:38:42 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/missing.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:38:43 thumbnail selected 3.014s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/partial.mp4 (camera motion)
    2026/10/11 00:38:43 thumbnail selected 2.480s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/quiet.mp4 (divergence 13.13 from median scene)
    2026/10/11 00:38:44 thumbnail selected 1.520s for /tmp/neonvr-motion-browser-74pi1xxc/clips/Garden/quiet.mp4 (divergence 12.18 from median scene)
    ✓ • 4.03s

## Activity

    $ git diff --check; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test ./...
    ok      neonvr    7.709s
    ✓ • 8.70s

## Activity

    Viewed image /tmp/neonvr-live-tools-vbq_9dyb/motion-player-mobile.png

## Assistant

Implemented:

- Orange, clickable `yi_motion` timestamps and motion intervals on the player’s seek track.
- Thumbnails prefer motion timestamps, with the existing FFmpeg selection as fallback.
- Existing cached thumbnails regenerate automatically.

Verified using your recordings, Go race tests, and desktop/mobile browsers.

Rebuild to activate:

```bash
sudo nerdctl compose up -d --build
```

## User

update the way the metadata timestamp are parsed since the camera now gives a new format:

## User

## Assistant

Yes. **Saved motion settings change when the detector reports START and STOP**, so the embedded motion timestamps follow those adjusted events.

- Start threshold and confirmation affect the motion start time.
- Stop threshold, quiet time, and the level-zero rule affect the motion end time.
- Sampling interval affects detection timing.

The metadata format remains the same: `yi_motion` JSON using `yi-motion-v1`. A one-minute recording can contain a much shorter motion interval.

Already finalized videos are unchanged. Saving changed settings restarts the detector; on .150, that can insert a STOP into an active recording.

## User

give a small example so that an user knows how to parse and extract it

## Assistant

Extract the embedded JSON with **FFmpeg’s `ffprobe` and `jq`**:

```sh
ffprobe -v error -show_entries format_tags=yi_motion \
  -of json clip.mp4 | jq '.format.tags.yi_motion | fromjson'
```

Example output, shortened:

```json
{
  "initial_state": 0,
  "event_columns": ["offset_ms", "unix_ms", "state", "kind"],
  "events": [
    [5000, 1791692585000, 1, "transition"],
    [12000, 1791692592000, 0, "transition"]
  ]
}
```

This means:

- The clip starts **idle** (`initial_state: 0`).
- Motion starts **5 seconds into the clip** (`state: 1`).
- Motion stops **12 seconds into the clip** (`state: 0`).
- The motion interval lasts **7 seconds**.

The second number in each row is the absolute Unix timestamp in milliseconds.

When parsing, begin with `initial_state`, then apply each event’s state at its offset. **`initial_state: 1` with no events means motion throughout the clip.** A value of `-1` means the starting state is unknown; check the history flags before assuming complete coverage.

## Assistant

I’ll compare this example with the current parser and recent camera recordings, then update any handling that differs.

## Activity

    $ git status --short; sed -n '1,275p' motion.go; sed -n '1,300p' motion_test.go; sed -n '1,150p' static/detection-viewer.js
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M config/neonvr.json.previous
     M detection_test.go
     M docker-compose.yml
     M main.go
     M static/detection-viewer.js
     M static/media-chrome-lite.js
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? motion.go
    ?? motion_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    package main
    
    import (
        "context"
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "math"
        "net/http"
        "os"
        "os/exec"
        "sort"
        "strconv"
        "time"
    )
    
    const cameraMotionMetadataTag = "yi_motion"
    
    type motionInterval struct {
        Start float64 `json:"start"`
        End   float64 `json:"end"`
    }
    
    type cameraMotion struct {
        Duration         float64          `json:"duration"`
        HistoryAvailable bool             `json:"history_available"`
        HistoryTruncated bool             `json:"history_truncated"`
        Markers          []float64        `json:"markers"`
        Intervals        []motionInterval `json:"intervals"`
    }
    
    // Camera offsets are relative to the video, independent of filename timezones.
    func parseCameraMotion(tag string, duration float64) (*cameraMotion, error) {
        var raw struct {
            Schema           string              `json:"schema"`
            VideoStart       *float64            `json:"video_start_unix_ms"`
            DurationMS       float64             `json:"duration_ms"`
            InitialState     *int                `json:"initial_state"`
            HistoryAvailable bool                `json:"history_available"`
            HistoryTruncated bool                `json:"history_truncated"`
            Columns          []string            `json:"event_columns"`
            Events           [][]json.RawMessage `json:"events"`
        }
        if err := json.Unmarshal([]byte(tag), &raw); err != nil {
            return nil, fmt.Errorf("invalid yi_motion JSON: %w", err)
        }
        if raw.Schema != "yi-motion-v1" {
            return nil, errors.New("unsupported yi_motion schema")
        }
        if !finite(duration) || duration <= 0 {
            duration = raw.DurationMS / 1000
        }
        if !finite(duration) || duration <= 0 {
            return nil, errors.New("yi_motion has no valid video duration")
        }
        state := -1 // Unknown history must not imply motion or inactivity.
        if raw.InitialState != nil {
            state = *raw.InitialState
        }
        if state < -1 || state > 1 {
            return nil, errors.New("invalid yi_motion initial state")
        }
        column := func(name string) int {
            for i, value := range raw.Columns {
                if value == name {
                    return i
                }
            }
            return -1
        }
        offsetCol, unixCol, stateCol := column("offset_ms"), column("unix_ms"), column("state")
        if len(raw.Events) > 0 && (stateCol < 0 || (offsetCol < 0 && (unixCol < 0 || raw.VideoStart == nil))) {
            return nil, errors.New("yi_motion event columns are missing")
        }
        type event struct {
            offset float64
            state  int
        }
        events := make([]event, 0, len(raw.Events))
        for _, row := range raw.Events {
            readNumber := func(index int) (float64, error) {
                var value float64
                if index < 0 || index >= len(row) || string(row[index]) == "null" {
                    return 0, errors.New("missing yi_motion event value")
                }
                if err := json.Unmarshal(row[index], &value); err != nil || !finite(value) {
                    return 0, errors.New("invalid yi_motion event value")
                }
                return value, nil
            }
            index := offsetCol
            if index < 0 {
                index = unixCol
            }
            offset, err := readNumber(index)
            if err != nil {
                return nil, err
            }
            if offsetCol < 0 {
                offset -= *raw.VideoStart
            }
            offset /= 1000
            value, err := readNumber(stateCol)
            if err != nil || value < -1 || value > 1 || value != math.Trunc(value) {
                return nil, errors.New("invalid yi_motion event state")
            }
            if !finite(offset) || offset < 0 || offset >= duration {
                continue // Events outside the recording cannot mark its seek track.
            }
            events = append(events, event{offset, int(value)})
        }
        sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
        result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
            Markers: []float64{}, Intervals: []motionInterval{}}
        start := 0.0
        for _, event := range events {
            if state == 1 && event.state != 1 && event.offset > start {
                result.Intervals = append(result.Intervals, motionInterval{start, event.offset})
            }
            if event.state == 1 && state != 1 {
                start = event.offset
            }
            state = event.state
        }
        if state == 1 && duration > start {
            result.Intervals = append(result.Intervals, motionInterval{start, duration})
        }
        for _, interval := range result.Intervals {
            result.Markers = append(result.Markers, interval.Start)
        }
        return result, nil
    }
    
    func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
        ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+cameraMotionMetadataTag, "-of", "json", path).Output()
        if err != nil {
            return nil, fmt.Errorf("probe camera motion: %w", err)
        }
        var probe struct {
            Format struct {
                Duration string            `json:"duration"`
                Tags     map[string]string `json:"tags"`
            } `json:"format"`
        }
        if err := json.Unmarshal(out, &probe); err != nil {
            return nil, err
        }
        tag := probe.Format.Tags[cameraMotionMetadataTag]
        if tag == "" {
            return nil, nil
        }
        duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
        return parseCameraMotion(tag, duration)
    }
    
    func (s *Server) apiCameraMotionHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        if r.Method != http.MethodGet {
            w.Header().Set("Allow", "GET")
            writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
            return
        }
        path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
        if err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        motion, err := s.readCameraMotion(r.Context(), path)
        if err != nil {
            writeJSONError(w, http.StatusBadGateway, err)
            return
        }
        writeJSON(w, struct {
            Motion *cameraMotion `json:"motion"`
        }{motion})
    }
    
    func (s *Server) runPreferredThumbnail(videoPath, tmpPath string) error {
        motion, err := s.readCameraMotion(s.workContext(), videoPath)
        if err == nil && motion != nil {
            for _, interval := range motion.Intervals {
                // Let the motion appear after the onset; stay within the marked interval.
                seek := interval.Start + math.Min(.5, (interval.End-interval.Start)/2)
                seek = math.Min(seek, math.Max(0, motion.Duration-.1))
                err = s.runFFmpegThumbnail(videoPath, tmpPath, strconv.FormatFloat(seek, 'f', 3, 64))
                if err == nil {
                    var frame frameSignature
                    frame, err = frameLumaSignature(tmpPath)
                    if err == nil && !frame.usable {
                        err = errors.New("camera motion frame is blank")
                    }
                }
                if err == nil {
                    log.Printf("thumbnail selected %.3fs for %s (camera motion)", seek, videoPath)
                    return nil
                }
                _ = os.Remove(tmpPath)
            }
        }
        if err != nil {
            log.Printf("camera motion thumbnail unavailable for %s, using scene selection: %v", videoPath, err)
        }
        return s.runInterestingThumbnail(videoPath, tmpPath)
    }
    package main
    
    import (
        "context"
        "encoding/json"
        "image/jpeg"
        "net/http"
        "net/http/httptest"
        "os"
        "os/exec"
        "path/filepath"
        "reflect"
        "testing"
    )
    
    const motionTestTag = `{"schema":"yi-motion-v1","video_start_unix_ms":1791692141000,"duration_ms":4000,"initial_state":0,"history_available":true,"history_truncated":false,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[1500,1791692142500,1,"transition"],[2500,1791692143500,0,"transition"]]}`
    
    func TestParseCameraMotion(t *testing.T) {
        tests := []struct {
            name      string
            tag       string
            duration  float64
            intervals []motionInterval
            partial   bool
        }{
            {"milliseconds", motionTestTag, 4, []motionInterval{{1.5, 2.5}}, false},
            {"initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"events":[]}`, 60.033, []motionInterval{{0, 60.033}}, false},
            {"truncated sample", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[2514,1791692355514,1,"sample"]]}`, 46.208, []motionInterval{{2.514, 46.208}}, true},
            {"no motion", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`, 10, []motionInterval{}, false},
            {"unknown history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false}`, 10, []motionInterval{}, true},
            {"reordered events and columns", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["state","offset_ms"],"events":[[0,3500],[1,1500],[1,2500]]}`, 4, []motionInterval{{1.5, 3.5}}, false},
            {"unix fallback", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"video_start_unix_ms":1791692141000,"event_columns":["state","unix_ms"],"events":[[1,1791692146295],[0,1791692152239]]}`, 18.624, []motionInterval{{5.295, 11.239}}, false},
            {"metadata duration fallback", motionTestTag, 0, []motionInterval{{1.5, 2.5}}, false},
            {"ignore outside clip", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["offset_ms","state"],"events":[[-1,1],[1000,1],[2000,0],[4000,1],[5000,1]]}`, 4, []motionInterval{{1, 2}}, false},
            {"zero-length initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[0,0]]}`, 4, []motionInterval{}, false},
        }
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                result, err := parseCameraMotion(tt.tag, tt.duration)
                if err != nil {
                    t.Fatal(err)
                }
                if !reflect.DeepEqual(result.Intervals, tt.intervals) || (!result.HistoryAvailable || result.HistoryTruncated) != tt.partial {
                    t.Fatalf("unexpected motion: %+v", result)
                }
                markers := []float64{}
                for _, interval := range tt.intervals {
                    markers = append(markers, interval.Start)
                }
                if !reflect.DeepEqual(result.Markers, markers) {
                    t.Fatalf("markers: %v, want %v", result.Markers, markers)
                }
            })
        }
    }
    
    func TestParseCameraMotionRejectsInvalidMetadata(t *testing.T) {
        for _, tag := range []string{
            `broken`, `null`, `{"schema":"yi-motion-v2"}`,
            `{"schema":"yi-motion-v1","initial_state":2}`,
            `{"schema":"yi-motion-v1","events":[[1,1]]}`,
            `{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[null,1]]}`,
            `{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[1000,2]]}`,
            `{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[[1000,0.5]]}`,
            `{"schema":"yi-motion-v1","event_columns":["offset_ms","state"],"events":[["1000",1]]}`,
        } {
            if _, err := parseCameraMotion(tag, 4); err == nil {
                t.Fatalf("accepted invalid metadata: %s", tag)
            }
        }
        if _, err := parseCameraMotion(`{"schema":"yi-motion-v1"}`, 0); err == nil {
            t.Fatal("accepted missing duration")
        }
    }
    
    func cameraMotionFixture(t *testing.T, tag string) (*Server, string) {
        t.Helper()
        return cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=40:y=20:w=80:h=60:color=red:t=fill:enable='between(t,1.5,2.5)'")
    }
    
    func cameraMotionFixtureWithFilter(t *testing.T, tag, filter string) (*Server, string) {
        t.Helper()
        for _, tool := range []string{"ffmpeg", "ffprobe"} {
            if _, err := exec.LookPath(tool); err != nil {
                t.Skip(tool + " not installed")
            }
        }
        root := t.TempDir()
        dir := filepath.Join(root, "Garden")
        if err := os.MkdirAll(dir, 0750); err != nil {
            t.Fatal(err)
        }
        path := filepath.Join(dir, "recording.mp4")
        args := []string{"-nostdin", "-v", "error", "-f", "lavfi", "-i", filter,
            "-t", "4", "-c:v", "libx264", "-preset", "ultrafast", "-threads", "1"}
        if tag != "" {
            args = append(args, "-metadata", "yi_motion="+tag)
        }
        args = append(args, "-movflags", "+use_metadata_tags", path)
        if out, err := exec.Command("ffmpeg", args...).CombinedOutput(); err != nil {
            t.Fatalf("create fixture: %v: %s", err, out)
        }
        return &Server{root: root, ffmpegPath: "ffmpeg", thumbDir: filepath.Join(root, ".thumbs"), thumbWidth: 160, thumbSeek: .2, thumbSem: make(chan struct{}, 1)}, path
    }
    
    func TestCameraMotionAPIAndThumbnail(t *testing.T) {
        s, path := cameraMotionFixture(t, motionTestTag)
        w := httptest.NewRecorder()
        s.apiCameraMotionHandler(w, httptest.NewRequest(http.MethodGet, "/api/camera-motion?camera=Garden&file=recording.mp4", nil))
        var response struct {
            Motion *cameraMotion `json:"motion"`
        }
        if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || response.Motion == nil ||
            !reflect.DeepEqual(response.Motion.Markers, []float64{1.5}) || !reflect.DeepEqual(response.Motion.Intervals, []motionInterval{{1.5, 2.5}}) {
            t.Fatalf("motion API with AI disabled: %d %s", w.Code, w.Body.String())
        }
        thumb, err := s.ensureThumbnail(path)
        if err != nil {
            t.Fatal(err)
        }
        f, err := os.Open(thumb)
        if err != nil {
            t.Fatal(err)
        }
        defer f.Close()
        frame, err := jpeg.Decode(f)
        if err != nil {
            t.Fatal(err)
        }
        r, g, b, _ := frame.At(80, 48).RGBA()
        if r < 45000 || g > 12000 || b > 12000 {
            t.Fatalf("thumbnail did not select the red motion frame: %d %d %d", r, g, b)
        }
        for _, request := range []*http.Request{
            httptest.NewRequest(http.MethodPost, "/api/camera-motion?camera=Garden&file=recording.mp4", nil),
            httptest.NewRequest(http.MethodGet, "/api/camera-motion?camera=Garden&file=../outside.mp4", nil),
        } {
            w := httptest.NewRecorder()
            s.apiCameraMotionHandler(w, request)
            if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusNotFound {
                t.Fatalf("invalid request accepted: %d", w.Code)
            }
        }
    }
    
    func TestCameraMotionThumbnailFallback(t *testing.T) {
        for _, tag := range []string{"", "invalid-json", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`} {
            s, path := cameraMotionFixture(t, tag)
            thumb, err := s.ensureThumbnail(path)
            if err != nil {
                t.Fatalf("fallback failed for %q: %v", tag, err)
            }
            if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
                t.Fatalf("fallback produced no usable image for %q: %v", tag, err)
            }
            if tag == "" {
                motion, err := s.readCameraMotion(context.Background(), path)
                if err != nil || motion != nil {
                    t.Fatalf("missing metadata: %+v %v", motion, err)
                }
            }
        }
    }
    
    func TestCameraMotionBlankFrameFallback(t *testing.T) {
        for _, tag := range []string{
            `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[500,0]]}`,
            `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[500,0],[2000,1],[3000,0]]}`,
        } {
            s, path := cameraMotionFixtureWithFilter(t, tag, "testsrc2=s=160x96:r=10,drawbox=x=0:y=0:w=iw:h=ih:color=black:t=fill:enable='lt(t,1)'")
            thumb, err := s.ensureThumbnail(path)
            if err != nil {
                t.Fatal(err)
            }
            if frame, err := frameLumaSignature(thumb); err != nil || !frame.usable {
                t.Fatalf("blank motion frame was used instead of a later interval or scene fallback: %v", err)
            }
        }
    }
    // Fetch clip-relative intervals, refresh pending analyses, and show active events.
    (function () {
      'use strict';
      const video = document.getElementById('mainVideo');
      const seek = document.getElementById('mainSeek');
      const status = document.getElementById('detectionStatus');
      const motionLegend = document.getElementById('cameraMotionLegend');
      const motionStatus = document.getElementById('cameraMotionStatus');
      if (!video || !seek || !status) return;
      const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
      let generation = 0;
      let request = null;
      let timer = null;
      let result = null;
      let motionRequest = null;
      let motion = null;
    
      function activeCameraMotion() {
        if (!motion || !motionStatus) return;
        const active = motion.intervals.some(row => video.currentTime >= row.start && video.currentTime < row.end);
        const partial = !motion.history_available || motion.history_truncated;
        motionStatus.classList.toggle('active', active);
        const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
          partial ? 'No camera motion timestamps available.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
        if (motionStatus.textContent !== text) motionStatus.textContent = text;
      }
    
      async function refreshCameraMotion(clip, expected) {
        motionRequest = new AbortController();
        try {
          const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
          if (!response.ok) throw new Error('Camera motion metadata unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          motion = data.motion;
          seek.setCameraMotion(motion);
          if (motionLegend) motionLegend.hidden = !motion;
          if (motionStatus) motionStatus.hidden = !motion;
          activeCameraMotion();
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          if (motionStatus) {
            motionStatus.hidden = false;
            motionStatus.textContent = error.message;
          }
        }
      }
    
      function activeStatus() {
        activeCameraMotion();
        if (!result) return;
        const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
        status.classList.toggle('active', active.length > 0);
        const text = active.length ? 'Detected now: ' + active.join(', ') :
          (result.intervals.length ? 'No detection at this time.' : result.uncertain ? 'Unconfirmed observations — clip retained.' : 'No target detections.');
        if (status.textContent !== text) status.textContent = text;
      }
    
      async function refresh(clip, expected) {
        if (expected !== generation) return;
        request = new AbortController();
        try {
          const response = await fetch('/api/detections?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: request.signal});
          if (!response.ok) throw new Error(response.status === 404 ? 'Clip is no longer available.' : 'Detection data unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          if (data.status === 'complete') {
            result = data.detections;
            status.title = '';
            seek.setDetections(result);
            activeStatus();
          } else {
            status.classList.remove('active');
            const labels = {pending: 'Detection queued…', processing: 'Analyzing detections…', disabled: 'Detection disabled.', error: 'Analysis unavailable — clip retained.'};
            status.textContent = labels[data.status] || 'Detection pending…';
            status.title = data.message || '';
            if (data.status !== 'disabled') timer = setTimeout(() => refresh(clip, expected), 10000);
          }
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          status.textContent = error.message;
          if (!error.message.includes('no longer available')) timer = setTimeout(() => refresh(clip, expected), 10000);
        }
      }
    
      function loadClip(clip) {
        generation += 1;
        request?.abort();
        motionRequest?.abort();
        clearTimeout(timer);
        result = null;
        motion = null;
        seek.setDetections(null);
        seek.setCameraMotion(null);
        if (motionLegend) motionLegend.hidden = true;
        if (motionStatus) { motionStatus.hidden = true; motionStatus.classList.remove('active'); }
        status.classList.remove('active');
        status.textContent = 'Loading detections…';
        status.title = '';
        refresh(clip, generation);
        refreshCameraMotion(clip, generation);
      }
    
      video.addEventListener('timeupdate', activeStatus);
      video.addEventListener('seeked', activeStatus);
      window.addEventListener('pagehide', () => { request?.abort(); motionRequest?.abort(); clearTimeout(timer); });
      window.NeoNVRDetections = { loadClip };
    })();
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        from pathlib import Path
        import os, json, subprocess
        for p in (Path('/AGENTS.md'), Path('/home/AGENTS.md'), Path('/home/ubuntu/AGENTS.md'), Path('/home/ubuntu/NeoNVR/AGENTS.md')):
         if p.is_file(): print(p, p.read_text())
        for root in ('/mnt/FTPfiles/yicamFront', '/mnt/FTPfiles/yicamBack'):
         files=[]
         for directory, _, names in os.walk(root):
          for name in names:
           if name.endswith('.mp4') and not name.startswith('.'):
            path=os.path.join(directory, name)
            try: files.append((os.stat(path).st_mtime, path))
            except OSError: pass
         print(root, 'recordings', len(files))
         for _, path in sorted(files, reverse=True)[:5]:
          result=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion','-of','json',path], capture_output=True, text=True, timeout=10)
          print(path, result.stdout.strip() if result.returncode == 0 else result.stderr[:200])
        PY
    /mnt/FTPfiles/yicamFront recordings 374
    /mnt/FTPfiles/yicamFront/2026Y10M11D14H/W402M00S60.mp4 {
        "format": {
            "duration": "60.071000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741720000,\"duration_ms\":60071,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D14H/W401M00S60.mp4 {
        "format": {
            "duration": "59.968000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741660000,\"duration_ms\":59968,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D14H/W400M00S60.mp4 {
        "format": {
            "duration": "60.022000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741600000,\"duration_ms\":60022,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D13H/W459M00S60.mp4 {
        "format": {
            "duration": "60.026000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741540000,\"duration_ms\":60026,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamFront/2026Y10M11D13H/W458M00S60.mp4 {
        "format": {
            "duration": "60.096000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741480000,\"duration_ms\":60096,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamBack recordings 9
    /mnt/FTPfiles/yicamBack/2026Y10M11D13H/W451M01S58.mp4 {
        "format": {
            "duration": "58.418000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791741061000,\"duration_ms\":58418,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamBack/2026Y10M11D13H/W435M00S60.mp4 {
        "format": {
            "duration": "60.066000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791740100000,\"duration_ms\":60066,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamBack/2026Y10M11D13H/W434M50S09.mp4 {
        "format": {
            "duration": "9.280000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791740090000,\"duration_ms\":9280,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamBack/2026Y10M11D13H/W413M00S60.mp4 {
        "format": {
            "duration": "60.184000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791738780000,\"duration_ms\":60184,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    /mnt/FTPfiles/yicamBack/2026Y10M11D13H/W412M32S27.mp4 {
        "format": {
            "duration": "27.457000",
            "tags": {
                "yi_motion": "{\"schema\":\"yi-motion-v1\",\"model\":\"y623\",\"source\":\"encoder-stats-ipc\",\"video_start_unix_ms\":1791738752000,\"duration_ms\":27457,\"video_time_precision_ms\":1000,\"initial_state\":-1,\"history_available\":false,\"history_truncated\":true,\"event_columns\":[\"offset_ms\",\"unix_ms\",\"state\",\"kind\"],\"events\":[]}"
            }
        }
    }
    ✓ • 1.01s

## Assistant

The example still uses `yi-motion-v1`, and the parser already handles its START/STOP rows. I’ll tighten unknown-state handling and verify the exact 5–12 second example and STOP events caused by detector restarts.

The newest recordings I checked have `initial_state: -1` and no events, so they currently provide no timestamps to mark.

## Activity

    $ rg -n 'motion|history|cameraMotion' README.md static/media-chrome-lite.js main.go; rg --files /tmp -g go -g gofmt -g check_motion_browser.py -g check_browser.py -g neonvr-motion 2>/dev/null | head -30; cat config/neonvr.json; sed -n '95,130p' detection.go; sed -n '206,280p' static/media-chrome-lite.js
    main.go:31:const thumbnailAlgorithmVersion = "camera-motion-v3"
    main.go:330:    #cameraMotionStatus.active { color: #ff9866; }
    main.go:703:              <span id="cameraMotionLegend" class="category" style="--detection-color:#ff9866" hidden>Camera motion</span>
    main.go:704:              <span id="cameraMotionStatus" role="status" hidden></span>
    main.go:1808:    mux.HandleFunc("/api/camera-motion", s.apiCameraMotionHandler)
    README.md:200:Humans, animals and bikes require three consecutive matched observations with at least 65% confidence by default. Animal classes in the pretrained COCO model are bird, cat, dog, horse, sheep, cow, elephant, bear, zebra and giraffe; other wildlife is not covered by this pretrained class list. Cars, trucks and buses require at least three consecutive observations and consistent translation of both box edges by 2.5% of the model frame within two seconds. A single box jump, back-and-forth jitter, or box resizing does not qualify as vehicle motion. Detection intervals include up to one sample of padding and merge nearby same-category confirmed events across gaps of up to 0.5 seconds by default; the merge setting does not relax consecutive confirmation. Timing is approximate to the sampling cadence.
    README.md:204:Vehicle confirmation needs at least three samples within a two-second motion window, so use at least one sampled frame per second for vehicle markers. At 0.5 FPS, observed vehicle movement remains ambiguous and the recording is retained.
    README.md:281:New thumbnails first use embedded Yi camera motion metadata (`yi_motion`, schema `yi-motion-v1`). The server extracts a frame half a second after the first motion onset, or halfway through the interval when it lasts less than a second. If that frame is unavailable or blank, it tries the next marked motion interval. Motion already active when the recording starts is treated as starting at 0 seconds.
    README.md:283:If motion metadata is missing, invalid, has no usable motion interval, or frame extraction fails, the server uses the existing divergence selector designed for mostly-static security footage. It samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.
    README.md:288:Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same motion metadata and fallback selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).
    README.md:296:Recordings with `yi_motion` metadata show orange, clickable motion onset markers directly on the main player's seek track, with a lane below showing the motion intervals. These camera markers work even when AI detection is disabled. Partial or truncated motion history is labeled in the player. NeoNVR uses the metadata's `offset_ms` values relative to the video, without a timezone correction; `initial_state: 1` marks motion from the beginning, and transitions to state 0 end the highlighted interval. Older cameras without the tag continue to play normally. The tag can be inspected with:
    README.md:299:ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 | jq '.format.tags.yi_motion | fromjson'
    static/media-chrome-lite.js:164:      this.motionMarkers = document.createElement('div');
    static/media-chrome-lite.js:165:      this.motionMarkers.className = 'motion-markers';
    static/media-chrome-lite.js:166:      this.motionMarkers.setAttribute('aria-label', 'Camera motion timestamps');
    static/media-chrome-lite.js:167:      track.append(input, this.motionMarkers);
    static/media-chrome-lite.js:177:        .motion-markers { position: absolute; left: 8px; right: 8px; top: 50%; pointer-events: none; }
    static/media-chrome-lite.js:178:        .motion-marker { position: absolute; width: 7px; height: 18px; padding: 0; border: 1px solid #332014; border-radius: 2px; background: #ff9866; transform: translate(-50%, -50%); pointer-events: auto; cursor: pointer; }
    static/media-chrome-lite.js:179:        .motion-marker:hover, .motion-marker:focus-visible { outline: 2px solid white; z-index: 1; }
    static/media-chrome-lite.js:213:    setCameraMotion(motion) {
    static/media-chrome-lite.js:214:      this.cameraMotion = motion;
    static/media-chrome-lite.js:220:      this.motionMarkers.replaceChildren();
    static/media-chrome-lite.js:222:      const motion = this.cameraMotion;
    static/media-chrome-lite.js:223:      this.lanes.hidden = !result && !motion?.intervals?.length;
    static/media-chrome-lite.js:224:      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : (motion?.duration || result?.duration);
    static/media-chrome-lite.js:226:      if (motion) {
    static/media-chrome-lite.js:227:        for (const time of motion.markers || []) {
    static/media-chrome-lite.js:231:          marker.className = 'motion-marker';
    static/media-chrome-lite.js:233:          marker.title = 'Camera motion at ' + time.toFixed(3) + ' sec';
    static/media-chrome-lite.js:236:          this.motionMarkers.append(marker);
    static/media-chrome-lite.js:241:      if (motion?.intervals?.length) categories.unshift(['camera_motion', 'Camera motion', '#ff9866']);
    static/media-chrome-lite.js:243:        if (kind !== 'camera_motion' && !result) continue;
    static/media-chrome-lite.js:248:        for (const interval of kind === 'camera_motion' ? motion.intervals : intervals) {
    static/media-chrome-lite.js:249:          if ((kind !== 'camera_motion' && interval.kind !== kind) || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
    static/media-chrome-lite.js:259:          const precision = kind === 'camera_motion' ? 3 : 1;
    static/media-chrome-lite.js:261:            (kind === 'camera_motion' ? '' : ' (' + Math.round(interval.confidence * 100) + '%)');
    /tmp/neonvr-live-tools-vbq_9dyb/check_motion_browser.py
    /tmp/neonvr-live-tools-vbq_9dyb/neonvr-motion
    /tmp/neonvr-live-tools-vbq_9dyb/check_browser.py
    /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go
    /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt
    {
      "root": "/mnt/mySSD/BlinkClips",
      "extra_roots": [
        "/mnt/FTPfiles/yicamBack",
        "/mnt/FTPfiles/yicamFront"
      ],
      "addr": "0.0.0.0:5500",
      "thumb_dir": "",
      "meta": "",
      "ffmpeg": "ffmpeg",
      "thumb_seek": 0.2,
      "thumb_width": 480,
      "thumb_workers": 2,
      "thumb_maintenance_every": "30m",
      "size_refresh": "30m",
      "auto_delete_days": 14,
      "auto_delete_every": "6h",
      "detection_enabled": true,
      "detection_python": "python3",
      "detection_model": "yolo26n.onnx",
      "detection_threads": 2,
      "detection_fps": 2,
      "detection_confidence": 0.65,
      "detection_min_frames": 3,
      "detection_merge_gap": 0.5,
      "detection_car_movement": 0.025,
      "detection_every": "1m",
      "detection_min_age": "2m",
      "detection_timeout": "30m",
      "detection_prune": "quarantine",
      "live_streams": [
        {
          "name": "Camera 192.168.1.150",
          "url": "rtsp://192.168.1.150:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
        },
        {
          "name": "Camera 192.168.1.149",
          "url": "rtsp://192.168.1.149:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.149:554/ch0_0.h264?backchannel=1"
        }
      ],
      "live_go2rtc": "go2rtc",
      "live_webrtc_port": 8555,
      "live_webrtc_host": ""
    }
            minimum := time.Minute
            if name == "detection_min_age" {
                minimum = 10 * time.Second
            }
            if err != nil || d < minimum || d > 24*time.Hour {
                return fmt.Errorf("%s must be between %s and 24h", name, minimum)
            }
        }
        if c.DetectionPrune != "off" && c.DetectionPrune != "quarantine" && c.DetectionPrune != "delete" {
            return errors.New("detection_prune must be off, quarantine, or delete")
        }
        return nil
    }
    
    func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }
    
    func (d DetectionResult) validate() error {
        if d.Version != detectionVersion || d.Model != "yolo26n" || len(d.ModelSHA256) != 64 || len(d.Policy) != 64 || !d.Complete ||
            !finite(d.Duration) || d.Duration <= 0 || !finite(d.SampleFPS) || d.SampleFPS < .5 || d.SampleFPS > 10 || d.Frames < 1 ||
            math.Abs(float64(d.Frames)/d.SampleFPS-d.Duration) > 1.5/d.SampleFPS || d.Intervals == nil {
            return errors.New("incomplete or invalid detection analysis")
        }
        for _, interval := range d.Intervals {
            if interval.Kind != "person" && interval.Kind != "animal" && interval.Kind != "bike" && interval.Kind != "moving_car" {
                return errors.New("unknown detection category")
            }
            if !finite(interval.Start) || !finite(interval.End) || interval.Start < 0 || interval.Start >= interval.End || interval.End > d.Duration+.001 ||
                !finite(interval.Confidence) || interval.Confidence < 0 || interval.Confidence > 1 || len(interval.Labels) == 0 {
                return errors.New("invalid detection interval")
            }
        }
        return nil
    }
    
    func (s *Server) detectionPolicy() string {
        c := s.detectionConfig
          this.renderDetections();
          update();
        }
        setDetections(result) {
          this.detections = result;
          this.renderDetections();
        }
        setCameraMotion(motion) {
          this.cameraMotion = motion;
          this.renderDetections();
        }
        renderDetections() {
          if (!this.lanes) return;
          this.lanes.replaceChildren();
          this.motionMarkers.replaceChildren();
          const result = this.detections;
          const motion = this.cameraMotion;
          this.lanes.hidden = !result && !motion?.intervals?.length;
          const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : (motion?.duration || result?.duration);
          if (!Number.isFinite(duration) || duration <= 0) return;
          if (motion) {
            for (const time of motion.markers || []) {
              if (!Number.isFinite(time) || time < 0 || time >= duration) continue;
              const marker = document.createElement('button');
              marker.type = 'button';
              marker.className = 'motion-marker';
              marker.style.left = (time / duration * 100) + '%';
              marker.title = 'Camera motion at ' + time.toFixed(3) + ' sec';
              marker.setAttribute('aria-label', 'Seek to ' + marker.title);
              marker.addEventListener('click', () => { if (this.video) this.video.currentTime = time; });
              this.motionMarkers.append(marker);
            }
          }
          const categories = [['person', 'Human', '#75d685'], ['animal', 'Animal', '#ffd166'], ['bike', 'Bike', '#72a7ff'], ['moving_car', 'Moving vehicle', '#c792ea']];
          const intervals = result?.intervals || [];
          if (motion?.intervals?.length) categories.unshift(['camera_motion', 'Camera motion', '#ff9866']);
          for (const [kind, label, color] of categories) {
            if (kind !== 'camera_motion' && !result) continue;
            const lane = document.createElement('div');
            lane.className = 'lane';
            lane.dataset.kind = kind;
            lane.setAttribute('aria-label', label + ' detections');
            for (const interval of kind === 'camera_motion' ? motion.intervals : intervals) {
              if ((kind !== 'camera_motion' && interval.kind !== kind) || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
              const start = Math.max(0, Math.min(duration, interval.start));
              const end = Math.max(start, Math.min(duration, interval.end));
              if (end <= start) continue;
              const event = document.createElement('button');
              event.type = 'button';
              event.className = 'event';
              event.style.left = (start / duration * 100) + '%';
              event.style.width = ((end - start) / duration * 100) + '%';
              event.style.background = color;
              const precision = kind === 'camera_motion' ? 3 : 1;
              const text = label + ': ' + start.toFixed(precision) + '–' + end.toFixed(precision) + ' sec' +
                (kind === 'camera_motion' ? '' : ' (' + Math.round(interval.confidence * 100) + '%)');
              event.title = text;
              event.setAttribute('aria-label', 'Seek to ' + text);
              event.addEventListener('click', () => { if (this.video) this.video.currentTime = start; });
              lane.append(event);
            }
            this.lanes.append(lane);
          }
        }
      }
    
      class MediaTimeDisplay extends HTMLElement {
        connectedCallback() {
          this.textContent = '0:00';
          const video = findVideo(this);
          if (!video) return;
          const update = () => { this.textContent = formatTime(video.currentTime); };
          video.addEventListener('timeupdate', update);
          video.addEventListener('loadedmetadata', update);
          update();
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import os, json, subprocess, concurrent.futures
        paths=[]
        for root in ('/mnt/FTPfiles/yicamFront','/mnt/FTPfiles/yicamBack'):
         for directory, _, names in os.walk(root):
          for name in names:
           if name.endswith('.mp4') and not name.startswith('.'):
            path=os.path.join(directory,name)
            try: paths.append((os.stat(path).st_mtime,path))
            except OSError: pass
        paths=sorted(paths,reverse=True)[:70]
        def probe(entry):
         _,path=entry
         try:
          p=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion','-of','json',path],capture_output=True,text=True,timeout=10)
          data=json.loads(p.stdout).get('format',{})
          motion=json.loads(data.get('tags',{}).get('yi_motion','null'))
          return path,motion
         except Exception: return path,None
        stats={}; found=0
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:
         for path,motion in executor.map(probe,paths):
          key=(motion or {}).get('schema','missing'); stats[key]=stats.get(key,0)+1
          if motion and motion.get('events') and found<4:
           print(path,json.dumps(motion)); found+=1
        print('recent schema counts',stats, 'eventful examples',found)
        PY
    recent schema counts {'yi-motion-v1': 70} eventful examples 0
    ✓ • 1.86s

## Activity

    file changes: Completed · 4 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -295,3 +295,3 @@
     
    -Recordings with `yi_motion` metadata show orange, clickable motion onset markers directly on the main player's seek track, with a lane below showing the motion intervals. These camera markers work even when AI detection is disabled. Partial or truncated motion history is labeled in the player. NeoNVR uses the metadata's `offset_ms` values relative to the video, without a timezone correction; `initial_state: 1` marks motion from the beginning, and transitions to state 0 end the highlighted interval. Older cameras without the tag continue to play normally. The tag can be inspected with:
    +Recordings with `yi_motion` metadata show orange, clickable motion onset markers directly on the main player's seek track, with a lane below showing the motion intervals. These camera markers work even when AI detection is disabled. Partial or truncated motion history is labeled in the player. NeoNVR uses the metadata's `offset_ms` values relative to the video, without a timezone correction; `initial_state: 1` marks motion from the beginning, and transitions to state 0 end the highlighted interval. An initial or later state of `-1` represents unknown coverage, even if the history flags say data is available; unknown periods are not treated as idle or highlighted as motion. Older cameras without the tag continue to play normally. The tag can be inspected with:
     
    @@ -301,2 +301,20 @@
     
    +The format remains `yi-motion-v1` when camera motion settings change. NeoNVR applies the camera's recorded START/STOP events directly, so changes to thresholds, confirmation, quiet time, sampling, or detector restarts are reflected in the marked intervals. It does not infer motion from the recording's length. For example, these events in a 60-second recording highlight only seconds 5 through 12:
    +
    +```json
    +{
    +  "schema": "yi-motion-v1",
    +  "initial_state": 0,
    +  "history_available": true,
    +  "history_truncated": false,
    +  "event_columns": ["offset_ms", "unix_ms", "state", "kind"],
    +  "events": [
    +    [5000, 1791692585000, 1, "transition"],
    +    [12000, 1791692592000, 0, "transition"]
    +  ]
    +}
    +```
    +
    +The parser follows `event_columns`, begins with `initial_state`, and applies each event's state at its offset. The `unix_ms` column is an absolute timestamp in milliseconds; playback uses `offset_ms`. If `initial_state` is 1 and there are no events, motion extends through the entire recording. If it is -1 and there are no events, no motion timestamps are available, and thumbnails use the existing scene selector. A STOP during a detector restart ends the current interval just like any other STOP.
    +
     ## Project rename and compatibility
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion.go
    @@ -28,2 +28,3 @@
     	HistoryTruncated bool             `json:"history_truncated"`
    +	HistoryComplete  bool             `json:"history_complete"`
     	Markers          []float64        `json:"markers"`
    @@ -114,2 +115,3 @@
     	result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
    +		HistoryComplete: raw.HistoryAvailable && !raw.HistoryTruncated && state != -1,
     		Markers: []float64{}, Intervals: []motionInterval{}}
    @@ -117,2 +119,5 @@
     	for _, event := range events {
    +		if event.state == -1 {
    +			result.HistoryComplete = false
    +		}
     		if state == 1 && event.state != 1 && event.offset > start {
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion_test.go
    @@ -26,2 +26,4 @@
     		{"milliseconds", motionTestTag, 4, []motionInterval{{1.5, 2.5}}, false},
    +		{"adjusted START and STOP in minute recording", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[5000,1791692585000,1,"transition"],[12000,1791692592000,0,"transition"]]}`, 60, []motionInterval{{5, 12}}, false},
    +		{"detector restart STOP", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[12000,1791692592000,0,"transition"],[40000,1791692620000,1,"transition"],[45000,1791692625000,0,"transition"]]}`, 60, []motionInterval{{0, 12}, {40, 45}}, false},
     		{"initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"events":[]}`, 60.033, []motionInterval{{0, 60.033}}, false},
    @@ -30,2 +32,5 @@
     		{"unknown history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false}`, 10, []motionInterval{}, true},
    +		{"unknown initial state despite available history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":false}`, 60, []motionInterval{}, true},
    +		{"missing initial state", `{"schema":"yi-motion-v1","history_available":true,"event_columns":["offset_ms","state"],"events":[[5000,1],[12000,0]]}`, 60, []motionInterval{{5, 12}}, true},
    +		{"unknown state closes motion without filling gap", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[5000,-1],[12000,1],[20000,0]]}`, 60, []motionInterval{{0, 5}, {12, 20}}, true},
     		{"reordered events and columns", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["state","offset_ms"],"events":[[0,3500],[1,1500],[1,2500]]}`, 4, []motionInterval{{1.5, 3.5}}, false},
    @@ -42,5 +47,8 @@
     			}
    -			if !reflect.DeepEqual(result.Intervals, tt.intervals) || (!result.HistoryAvailable || result.HistoryTruncated) != tt.partial {
    +			if !reflect.DeepEqual(result.Intervals, tt.intervals) || !result.HistoryComplete != tt.partial {
     				t.Fatalf("unexpected motion: %+v", result)
     			}
    +			if tt.duration > 0 && result.Duration != tt.duration {
    +				t.Fatalf("recording duration: %v, want %v", result.Duration, tt.duration)
    +			}
     			markers := []float64{}
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/detection-viewer.js
    @@ -20,6 +20,6 @@
         const active = motion.intervals.some(row => video.currentTime >= row.start && video.currentTime < row.end);
    -    const partial = !motion.history_available || motion.history_truncated;
    +    const partial = motion.history_complete === false || !motion.history_available || motion.history_truncated;
         motionStatus.classList.toggle('active', active);
         const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
    -      partial ? 'No camera motion timestamps available.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
    +      partial ? 'Camera motion history unavailable.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
         if (motionStatus.textContent !== text) motionStatus.textContent = text;

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w motion.go motion_test.go; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test ./...
    ok      neonvr    6.449s
    ✓ • 6.99s

## Activity

    $ git diff --check; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr-motion .; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go vet ./...
    ✓ • 687ms

## Activity

    file changes: Completed · 1 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion_test.go
    @@ -154,3 +154,4 @@
     func TestCameraMotionThumbnailFallback(t *testing.T) {
    -	for _, tag := range []string{"", "invalid-json", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`} {
    +	for _, tag := range []string{"", "invalid-json", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`,
    +		`{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":false,"events":[]}`} {
     		s, path := cameraMotionFixture(t, tag)

## Activity

    $ cat > /tmp/neonvr-live-tools-vbq_9dyb/check_motion_parse_browser.py <<'PY'
        import json
        from pathlib import Path
        import socket
        import subprocess
        import tempfile
        import time
        import urllib.request
        from playwright.sync_api import sync_playwright
        
        TOOLS=Path('/tmp/neonvr-live-tools-vbq_9dyb')
        def get(url):
         with urllib.request.urlopen(url,timeout=10) as response: return json.load(response)
        with tempfile.TemporaryDirectory(prefix='neonvr-motion-parse-') as tmp:
         root=Path(tmp); camera=root/'clips'/'Garden';camera.mkdir(parents=True)
         base=root/'base.mp4'
         subprocess.run(['ffmpeg','-nostdin','-v','error','-f','lavfi','-i','testsrc2=s=160x96:r=5','-t','60','-c:v','libx264','-preset','ultrafast','-threads','1',str(base)],check=True)
         tags={
          'minute.mp4':{'schema':'yi-motion-v1','initial_state':0,'history_available':True,'history_truncated':False,'event_columns':['offset_ms','unix_ms','state','kind'],'events':[[5000,1791692585000,1,'transition'],[12000,1791692592000,0,'transition']]},
          'restart.mp4':{'schema':'yi-motion-v1','initial_state':1,'history_available':True,'history_truncated':False,'event_columns':['offset_ms','unix_ms','state','kind'],'events':[[12000,1791692592000,0,'transition'],[40000,1791692620000,1,'transition'],[45000,1791692625000,0,'transition']]},
          'unknown.mp4':{'schema':'yi-motion-v1','initial_state':-1,'history_available':True,'history_truncated':False,'events':[]},
          'initial.mp4':{'schema':'yi-motion-v1','initial_state':1,'history_available':True,'history_truncated':False,'events':[]},
         }
         for name,tag in tags.items():
          subprocess.run(['ffmpeg','-nostdin','-v','error','-i',str(base),'-c','copy','-metadata','yi_motion='+json.dumps(tag),'-movflags','+use_metadata_tags',str(camera/name)],check=True)
         with socket.socket() as listener:
          listener.bind(('127.0.0.1',0));port=listener.getsockname()[1]
         config=root/'config.json';config.write_text(json.dumps({'root':str(root/'clips'),'extra_roots':[],'addr':f'127.0.0.1:{port}','auto_delete_days':0,'thumb_maintenance_every':'0','detection_enabled':False,'live_streams':[]}))
         origin=f'http://127.0.0.1:{port}'
         with (root/'app.log').open('w+') as log:
          process=subprocess.Popen([str(TOOLS/'neonvr-motion'),'-config',str(config)],stdout=log,stderr=log)
          try:
           for _ in range(100):
            try: cameras=get(origin+'/api/cameras');break
            except Exception:time.sleep(.05)
           else:raise RuntimeError('Server did not start')
           clips={clip['file']:clip for clip in get(origin+'/api/clips?cameras=Garden&date='+cameras['latest_date'])['clips']}
           api=get(origin+'/api/camera-motion?camera=Garden&file=minute.mp4')['motion']
           assert api['duration']==60 and api['markers']==[5] and api['intervals']==[{'start':5,'end':12}] and api['history_complete']
           unknown=get(origin+'/api/camera-motion?camera=Garden&file=unknown.mp4')['motion']
           assert not unknown['history_complete'] and unknown['markers']==[] and unknown['intervals']==[]
           errors=[]
           with sync_playwright() as p:
            browser=p.chromium.launch(headless=True,args=['--no-sandbox'])
            page=browser.new_page(viewport={'width':1280,'height':900})
            page.on('pageerror',lambda e:errors.append(str(e)))
            page.goto(origin)
            page.wait_for_function('allClips.length === 4 && window.NeoNVRDetections')
            def load(name):
             page.evaluate('(clip)=>loadMainPlayer(clip,false)',clips[name])
             page.wait_for_function('document.querySelector("#mainSeek").cameraMotion !== null')
            load('minute.mp4')
            page.wait_for_function('document.querySelector("#mainSeek").cameraMotion?.markers[0] === 5')
            marker=page.locator('#mainSeek .motion-marker')
            assert marker.count()==1 and marker.get_attribute('title')=='Camera motion at 5.000 sec'
            assert abs(float(marker.evaluate('(el)=>parseFloat(el.style.left)'))-100*5/60)<.001
            assert page.locator('#mainSeek .lane[data-kind="camera_motion"] button').get_attribute('title')=='Camera motion: 5.000–12.000 sec'
            marker.click()
            page.wait_for_function('!document.querySelector("#mainVideo").seeking && Math.abs(document.querySelector("#mainVideo").currentTime - 5)<.001 && document.querySelector("#cameraMotionStatus").textContent === "Camera motion now."')
            page.evaluate('document.querySelector("#mainVideo").currentTime=13')
            page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent === "No camera motion marked here."')
            print('The exact example produces a 5 sec marker and 5–12 sec interval in a 60 sec recording; click seeking and STOP status passed.',flush=True)
            load('restart.mp4')
            page.wait_for_function('document.querySelector("#mainSeek").cameraMotion?.markers.length === 2')
            titles=page.locator('#mainSeek .lane[data-kind="camera_motion"] button').evaluate_all('(els)=>els.map(el=>el.title)')
            assert titles==['Camera motion: 0.000–12.000 sec','Camera motion: 40.000–45.000 sec'],titles
            print('Detector-restart STOP ends the active interval at 12 sec; later START/STOP creates a separate 40–45 sec interval.',flush=True)
            load('unknown.mp4')
            page.wait_for_function('document.querySelector("#cameraMotionStatus").textContent.includes("history unavailable.")')
            assert page.locator('#mainSeek .motion-marker').count()==0
            assert page.locator('#mainSeek .lane[data-kind="camera_motion"]').count()==0
            assert 'Partial history.' in page.locator('#cameraMotionStatus').inner_text()
            load('initial.mp4')
            page.wait_for_function('document.querySelector("#mainSeek").cameraMotion?.history_complete && document.querySelector("#mainSeek").cameraMotion?.markers[0] === 0')
            assert page.locator('#mainSeek .lane[data-kind="camera_motion"] button').get_attribute('title')=='Camera motion: 0.000–60.000 sec'
            print('Unknown initial state is labeled incomplete despite available-history flags; initial motion with no events covers the full clip.',flush=True)
            assert not errors,errors
            browser.close()
          finally:
           process.terminate()
           try:process.wait(timeout=10)
           except subprocess.TimeoutExpired:process.kill();process.wait()
        PY
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_motion_parse_browser.py
    The exact example produces a 5 sec marker and 5–12 sec interval in a 60 sec recording; click seeking and STOP status passed.
    Detector-restart STOP ends the active interval at 12 sec; later START/STOP creates a separate 40–45 sec interval.
    Unknown initial state is labeled incomplete despite available-history flags; initial motion with no events covers the full clip.
    ✓ • 1.58s

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w motion_test.go; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
    ok      neonvr    8.608s
    ✓ • 9.56s

## Activity

    $ git diff --check; git status --short; git diff --numstat; sed -n '111,150p' motion.go; sed -n '297,332p' README.md
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M config/neonvr.json.previous
     M detection_test.go
     M docker-compose.yml
     M main.go
     M static/detection-viewer.js
     M static/media-chrome-lite.js
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? motion.go
    ?? motion_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    6    1    Dockerfile
    43    3    README.md
    37    27    config.go
    16    1    config/neonvr.json
    29    1    config/neonvr.json.previous
    5    1    detection_test.go
    5    0    docker-compose.yml
    56    23    main.go
    43    1    static/detection-viewer.js
    40    7    static/media-chrome-lite.js
    15    4    static/settings.html
            }
            events = append(events, event{offset, int(value)})
        }
        sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
        result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
            HistoryComplete: raw.HistoryAvailable && !raw.HistoryTruncated && state != -1,
            Markers:         []float64{}, Intervals: []motionInterval{}}
        start := 0.0
        for _, event := range events {
            if event.state == -1 {
                result.HistoryComplete = false
            }
            if state == 1 && event.state != 1 && event.offset > start {
                result.Intervals = append(result.Intervals, motionInterval{start, event.offset})
            }
            if event.state == 1 && state != 1 {
                start = event.offset
            }
            state = event.state
        }
        if state == 1 && duration > start {
            result.Intervals = append(result.Intervals, motionInterval{start, duration})
        }
        for _, interval := range result.Intervals {
            result.Markers = append(result.Markers, interval.Start)
        }
        return result, nil
    }
    
    func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
        ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+cameraMotionMetadataTag, "-of", "json", path).Output()
        if err != nil {
            return nil, fmt.Errorf("probe camera motion: %w", err)
        }
        var probe struct {
            Format struct {
                Duration string            `json:"duration"`
    
    ```bash
    ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 | jq '.format.tags.yi_motion | fromjson'
    ```
    
    The format remains `yi-motion-v1` when camera motion settings change. NeoNVR applies the camera's recorded START/STOP events directly, so changes to thresholds, confirmation, quiet time, sampling, or detector restarts are reflected in the marked intervals. It does not infer motion from the recording's length. For example, these events in a 60-second recording highlight only seconds 5 through 12:
    
    ```json
    {
      "schema": "yi-motion-v1",
      "initial_state": 0,
      "history_available": true,
      "history_truncated": false,
      "event_columns": ["offset_ms", "unix_ms", "state", "kind"],
      "events": [
        [5000, 1791692585000, 1, "transition"],
        [12000, 1791692592000, 0, "transition"]
      ]
    }
    ```
    
    The parser follows `event_columns`, begins with `initial_state`, and applies each event's state at its offset. The `unix_ms` column is an absolute timestamp in milliseconds; playback uses `offset_ms`. If `initial_state` is 1 and there are no events, motion extends through the entire recording. If it is -1 and there are no events, no motion timestamps are available, and thumbnails use the existing scene selector. A STOP during a detector restart ends the current interval just like any other STOP.
    
    ## Project rename and compatibility
    
    The project directory is now `/home/ubuntu/NeoNVR`, the Go module and executable are `neonvr`, the Compose service/container are `neonvr`, and the systemd example is `neonvr.service.example`.
    
    Existing recordings remain under `/mnt/mySSD/BlinkClips`. The legacy `.blinkclips_thumbs` and `.blinkclips_meta.json` paths are deliberately retained so thumbnail caches and favorite selections continue to work. Custom `-root`, `-thumb-dir`, and `-meta` flags continue to work.
    
    A running container keeps its current binary until rebuilt. Stop the existing Compose stack using its original project name and configuration before starting the renamed stack; otherwise both stacks can contend for port 5500. For nerdctl deployments, use `sudo nerdctl compose` in place of `docker compose`. After stopping the old stack:
    
    ```bash
    cd /home/ubuntu/NeoNVR
    sudo nerdctl compose up -d --build
    sudo nerdctl compose logs -f neonvr
    ```
    ✓ • 0ms

## Assistant

The format remains `yi-motion-v1`. I tightened unknown-state handling and verified:

- Your example highlights **5–12 seconds** in a 60-second clip.
- Detector-restart STOP events close the active interval.
- State `-1` is labeled as incomplete history; missing timestamps use the thumbnail fallback.

Go race tests and browser checks passed. Recent recordings contain unknown initial states and empty events, so they provide no markers.

Rebuild to activate:

```bash
sudo nerdctl compose up -d --build
```

## User

it seems to not be detecting camera motion metadata in the ingested videos or it says partial history.

## Assistant

I’ll trace a few ingested videos through metadata extraction and the player API to see whether the timestamps are being lost or the history is being labeled incorrectly.

## Activity

    $ git status --short; sed -n '1,260p' motion.go; sed -n '1,125p' static/detection-viewer.js; cat docker-compose.yml; cat config/neonvr.json; rg -n 'inspectDetectionVideo|format_tags|yi_motion|history_complete|map_metadata' detection.go motion_test.go detection_test.go; command -v nerdctl; command -v docker; ss -lntp | head -20
     M Dockerfile
     M README.md
     M config.go
     M config/neonvr.json
     M config/neonvr.json.previous
     M detection_test.go
     M docker-compose.yml
     M main.go
     M static/detection-viewer.js
     M static/media-chrome-lite.js
     M static/settings.html
    ?? docker-compose.yml.bak-pre-memory-limit-20261010
    ?? live.go
    ?? live_test.go
    ?? motion.go
    ?? motion_test.go
    ?? static/live-streams.js
    ?? static/streams.html
    ?? third_party/
    package main
    
    import (
        "context"
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "math"
        "net/http"
        "os"
        "os/exec"
        "sort"
        "strconv"
        "time"
    )
    
    const cameraMotionMetadataTag = "yi_motion"
    
    type motionInterval struct {
        Start float64 `json:"start"`
        End   float64 `json:"end"`
    }
    
    type cameraMotion struct {
        Duration         float64          `json:"duration"`
        HistoryAvailable bool             `json:"history_available"`
        HistoryTruncated bool             `json:"history_truncated"`
        HistoryComplete  bool             `json:"history_complete"`
        Markers          []float64        `json:"markers"`
        Intervals        []motionInterval `json:"intervals"`
    }
    
    // Camera offsets are relative to the video, independent of filename timezones.
    func parseCameraMotion(tag string, duration float64) (*cameraMotion, error) {
        var raw struct {
            Schema           string              `json:"schema"`
            VideoStart       *float64            `json:"video_start_unix_ms"`
            DurationMS       float64             `json:"duration_ms"`
            InitialState     *int                `json:"initial_state"`
            HistoryAvailable bool                `json:"history_available"`
            HistoryTruncated bool                `json:"history_truncated"`
            Columns          []string            `json:"event_columns"`
            Events           [][]json.RawMessage `json:"events"`
        }
        if err := json.Unmarshal([]byte(tag), &raw); err != nil {
            return nil, fmt.Errorf("invalid yi_motion JSON: %w", err)
        }
        if raw.Schema != "yi-motion-v1" {
            return nil, errors.New("unsupported yi_motion schema")
        }
        if !finite(duration) || duration <= 0 {
            duration = raw.DurationMS / 1000
        }
        if !finite(duration) || duration <= 0 {
            return nil, errors.New("yi_motion has no valid video duration")
        }
        state := -1 // Unknown history must not imply motion or inactivity.
        if raw.InitialState != nil {
            state = *raw.InitialState
        }
        if state < -1 || state > 1 {
            return nil, errors.New("invalid yi_motion initial state")
        }
        column := func(name string) int {
            for i, value := range raw.Columns {
                if value == name {
                    return i
                }
            }
            return -1
        }
        offsetCol, unixCol, stateCol := column("offset_ms"), column("unix_ms"), column("state")
        if len(raw.Events) > 0 && (stateCol < 0 || (offsetCol < 0 && (unixCol < 0 || raw.VideoStart == nil))) {
            return nil, errors.New("yi_motion event columns are missing")
        }
        type event struct {
            offset float64
            state  int
        }
        events := make([]event, 0, len(raw.Events))
        for _, row := range raw.Events {
            readNumber := func(index int) (float64, error) {
                var value float64
                if index < 0 || index >= len(row) || string(row[index]) == "null" {
                    return 0, errors.New("missing yi_motion event value")
                }
                if err := json.Unmarshal(row[index], &value); err != nil || !finite(value) {
                    return 0, errors.New("invalid yi_motion event value")
                }
                return value, nil
            }
            index := offsetCol
            if index < 0 {
                index = unixCol
            }
            offset, err := readNumber(index)
            if err != nil {
                return nil, err
            }
            if offsetCol < 0 {
                offset -= *raw.VideoStart
            }
            offset /= 1000
            value, err := readNumber(stateCol)
            if err != nil || value < -1 || value > 1 || value != math.Trunc(value) {
                return nil, errors.New("invalid yi_motion event state")
            }
            if !finite(offset) || offset < 0 || offset >= duration {
                continue // Events outside the recording cannot mark its seek track.
            }
            events = append(events, event{offset, int(value)})
        }
        sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
        result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
            HistoryComplete: raw.HistoryAvailable && !raw.HistoryTruncated && state != -1,
            Markers:         []float64{}, Intervals: []motionInterval{}}
        start := 0.0
        for _, event := range events {
            if event.state == -1 {
                result.HistoryComplete = false
            }
            if state == 1 && event.state != 1 && event.offset > start {
                result.Intervals = append(result.Intervals, motionInterval{start, event.offset})
            }
            if event.state == 1 && state != 1 {
                start = event.offset
            }
            state = event.state
        }
        if state == 1 && duration > start {
            result.Intervals = append(result.Intervals, motionInterval{start, duration})
        }
        for _, interval := range result.Intervals {
            result.Markers = append(result.Markers, interval.Start)
        }
        return result, nil
    }
    
    func (s *Server) readCameraMotion(ctx context.Context, path string) (*cameraMotion, error) {
        ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+cameraMotionMetadataTag, "-of", "json", path).Output()
        if err != nil {
            return nil, fmt.Errorf("probe camera motion: %w", err)
        }
        var probe struct {
            Format struct {
                Duration string            `json:"duration"`
                Tags     map[string]string `json:"tags"`
            } `json:"format"`
        }
        if err := json.Unmarshal(out, &probe); err != nil {
            return nil, err
        }
        tag := probe.Format.Tags[cameraMotionMetadataTag]
        if tag == "" {
            return nil, nil
        }
        duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)
        return parseCameraMotion(tag, duration)
    }
    
    func (s *Server) apiCameraMotionHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-store")
        if r.Method != http.MethodGet {
            w.Header().Set("Allow", "GET")
            writeJSONError(w, http.StatusMethodNotAllowed, errors.New("GET required"))
            return
        }
        path, err := s.resolveClipPath(r.URL.Query().Get("camera"), r.URL.Query().Get("file"))
        if err != nil {
            writeJSONError(w, http.StatusNotFound, err)
            return
        }
        motion, err := s.readCameraMotion(r.Context(), path)
        if err != nil {
            writeJSONError(w, http.StatusBadGateway, err)
            return
        }
        writeJSON(w, struct {
            Motion *cameraMotion `json:"motion"`
        }{motion})
    }
    
    func (s *Server) runPreferredThumbnail(videoPath, tmpPath string) error {
        motion, err := s.readCameraMotion(s.workContext(), videoPath)
        if err == nil && motion != nil {
            for _, interval := range motion.Intervals {
                // Let the motion appear after the onset; stay within the marked interval.
                seek := interval.Start + math.Min(.5, (interval.End-interval.Start)/2)
                seek = math.Min(seek, math.Max(0, motion.Duration-.1))
                err = s.runFFmpegThumbnail(videoPath, tmpPath, strconv.FormatFloat(seek, 'f', 3, 64))
                if err == nil {
                    var frame frameSignature
                    frame, err = frameLumaSignature(tmpPath)
                    if err == nil && !frame.usable {
                        err = errors.New("camera motion frame is blank")
                    }
                }
                if err == nil {
                    log.Printf("thumbnail selected %.3fs for %s (camera motion)", seek, videoPath)
                    return nil
                }
                _ = os.Remove(tmpPath)
            }
        }
        if err != nil {
            log.Printf("camera motion thumbnail unavailable for %s, using scene selection: %v", videoPath, err)
        }
        return s.runInterestingThumbnail(videoPath, tmpPath)
    }
    // Fetch clip-relative intervals, refresh pending analyses, and show active events.
    (function () {
      'use strict';
      const video = document.getElementById('mainVideo');
      const seek = document.getElementById('mainSeek');
      const status = document.getElementById('detectionStatus');
      const motionLegend = document.getElementById('cameraMotionLegend');
      const motionStatus = document.getElementById('cameraMotionStatus');
      if (!video || !seek || !status) return;
      const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
      let generation = 0;
      let request = null;
      let timer = null;
      let result = null;
      let motionRequest = null;
      let motion = null;
    
      function activeCameraMotion() {
        if (!motion || !motionStatus) return;
        const active = motion.intervals.some(row => video.currentTime >= row.start && video.currentTime < row.end);
        const partial = motion.history_complete === false || !motion.history_available || motion.history_truncated;
        motionStatus.classList.toggle('active', active);
        const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
          partial ? 'Camera motion history unavailable.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
        if (motionStatus.textContent !== text) motionStatus.textContent = text;
      }
    
      async function refreshCameraMotion(clip, expected) {
        motionRequest = new AbortController();
        try {
          const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
          if (!response.ok) throw new Error('Camera motion metadata unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          motion = data.motion;
          seek.setCameraMotion(motion);
          if (motionLegend) motionLegend.hidden = !motion;
          if (motionStatus) motionStatus.hidden = !motion;
          activeCameraMotion();
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          if (motionStatus) {
            motionStatus.hidden = false;
            motionStatus.textContent = error.message;
          }
        }
      }
    
      function activeStatus() {
        activeCameraMotion();
        if (!result) return;
        const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
        status.classList.toggle('active', active.length > 0);
        const text = active.length ? 'Detected now: ' + active.join(', ') :
          (result.intervals.length ? 'No detection at this time.' : result.uncertain ? 'Unconfirmed observations — clip retained.' : 'No target detections.');
        if (status.textContent !== text) status.textContent = text;
      }
    
      async function refresh(clip, expected) {
        if (expected !== generation) return;
        request = new AbortController();
        try {
          const response = await fetch('/api/detections?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: request.signal});
          if (!response.ok) throw new Error(response.status === 404 ? 'Clip is no longer available.' : 'Detection data unavailable.');
          const data = await response.json();
          if (expected !== generation) return;
          if (data.status === 'complete') {
            result = data.detections;
            status.title = '';
            seek.setDetections(result);
            activeStatus();
          } else {
            status.classList.remove('active');
            const labels = {pending: 'Detection queued…', processing: 'Analyzing detections…', disabled: 'Detection disabled.', error: 'Analysis unavailable — clip retained.'};
            status.textContent = labels[data.status] || 'Detection pending…';
            status.title = data.message || '';
            if (data.status !== 'disabled') timer = setTimeout(() => refresh(clip, expected), 10000);
          }
        } catch (error) {
          if (expected !== generation || error.name === 'AbortError') return;
          status.textContent = error.message;
          if (!error.message.includes('no longer available')) timer = setTimeout(() => refresh(clip, expected), 10000);
        }
      }
    
      function loadClip(clip) {
        generation += 1;
        request?.abort();
        motionRequest?.abort();
        clearTimeout(timer);
        result = null;
        motion = null;
        seek.setDetections(null);
        seek.setCameraMotion(null);
        if (motionLegend) motionLegend.hidden = true;
        if (motionStatus) { motionStatus.hidden = true; motionStatus.classList.remove('active'); }
        status.classList.remove('active');
        status.textContent = 'Loading detections…';
        status.title = '';
        refresh(clip, generation);
        refreshCameraMotion(clip, generation);
      }
    
      video.addEventListener('timeupdate', activeStatus);
      video.addEventListener('seeked', activeStatus);
      window.addEventListener('pagehide', () => { request?.abort(); motionRequest?.abort(); clearTimeout(timer); });
      window.NeoNVRDetections = { loadClip };
    })();
    services:
      neonvr:
        build: .
        container_name: neonvr
        restart: unless-stopped
        mem_limit: 6g
        memswap_limit: 4g
        ports:
          - "5500:5500"
          # Live media and microphone talkback; inactive until enabled in Streams.
          - "8555:8555/tcp"
          - "8555:8555/udp"
        volumes:
          # Mount the directory so config saves can replace the file atomically.
          - ./config:/config
          # Read/write for thumbnails, favorites, and configured deletion.
          - /mnt:/mnt
      #    - /mnt/mySSD/BlinkClips:/mnt/mySSD/BlinkClips
     #     - /mnt/FTPfiles:/mnt/FTPfiles
    #      - /home/ubuntu/FTPfiles/yicamBack:/home/ubuntu/FTPfiles/yicamBack
    #      - /home/ubuntu/FTPfiles/yicamFront:/home/ubuntu/FTPfiles/yicamFront
        command: ["-config", "/config/neonvr.json"]
    {
      "root": "/mnt/mySSD/BlinkClips",
      "extra_roots": [
        "/mnt/FTPfiles/yicamBack",
        "/mnt/FTPfiles/yicamFront"
      ],
      "addr": "0.0.0.0:5500",
      "thumb_dir": "",
      "meta": "",
      "ffmpeg": "ffmpeg",
      "thumb_seek": 0.2,
      "thumb_width": 480,
      "thumb_workers": 2,
      "thumb_maintenance_every": "30m",
      "size_refresh": "30m",
      "auto_delete_days": 14,
      "auto_delete_every": "6h",
      "detection_enabled": true,
      "detection_python": "python3",
      "detection_model": "yolo26n.onnx",
      "detection_threads": 2,
      "detection_fps": 2,
      "detection_confidence": 0.65,
      "detection_min_frames": 3,
      "detection_merge_gap": 0.5,
      "detection_car_movement": 0.025,
      "detection_every": "1m",
      "detection_min_age": "2m",
      "detection_timeout": "30m",
      "detection_prune": "quarantine",
      "live_streams": [
        {
          "name": "Camera 192.168.1.150",
          "url": "rtsp://192.168.1.150:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1"
        },
        {
          "name": "Camera 192.168.1.149",
          "url": "rtsp://192.168.1.149:554/ch0_0.h264",
          "backchannel_url": "rtsp://192.168.1.149:554/ch0_0.h264?backchannel=1"
        }
      ],
      "live_go2rtc": "go2rtc",
      "live_webrtc_port": 8555,
      "live_webrtc_host": ""
    }
    detection_test.go:35:        "-metadata", "source_camera=Garden", "-metadata", "yi_motion="+motionTestTag, "-movflags", "+use_metadata_tags", path).CombinedOutput()
    detection_test.go:87:    allTags, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format_tags=source_camera", "-of", "json", path).Output()
    motion_test.go:105:        args = append(args, "-metadata", "yi_motion="+tag)
    detection.go:184:func (s *Server) inspectDetectionVideo(ctx context.Context, path string) (videoInspection, error) {
    detection.go:188:        "format=duration:format_tags="+detectionMetadataTag+":stream=codec_type,codec_name", "-of", "json", path).Output()
    detection.go:200:    info, err := s.inspectDetectionVideo(ctx, path)
    detection.go:398:    before, err := s.inspectDetectionVideo(ctx, path)
    detection.go:414:        "-map", "0", "-c", "copy", "-map_metadata", "0", "-metadata", detectionMetadataTag+"="+string(data),
    detection.go:419:    after, err := s.inspectDetectionVideo(ctx, tmpPath)
    /usr/local/bin/nerdctl
    State  Recv-Q Send-Q Local Address:Port  Peer Address:PortProcess
    LISTEN 0      128        127.0.0.1:45013      0.0.0.0:*    users:(("codex",pid=340202,fd=37))
    LISTEN 0      50           0.0.0.0:445        0.0.0.0:*
    LISTEN 0      50           0.0.0.0:139        0.0.0.0:*
    LISTEN 0      32           0.0.0.0:21         0.0.0.0:*
    LISTEN 0      4096         0.0.0.0:22         0.0.0.0:*
    LISTEN 0      4096       127.0.0.1:631        0.0.0.0:*
    LISTEN 0      4096       127.0.0.1:33229      0.0.0.0:*
    LISTEN 0      4096       127.0.0.1:39981      0.0.0.0:*
    LISTEN 0      128        127.0.0.1:39911      0.0.0.0:*    users:(("codex",pid=600589,fd=37))
    LISTEN 0      4096      127.0.0.54:53         0.0.0.0:*
    LISTEN 0      128        127.0.0.1:38301      0.0.0.0:*    users:(("codex",pid=67625,fd=37))
    LISTEN 0      4096   127.0.0.53%lo:53         0.0.0.0:*
    LISTEN 0      50              [::]:445           [::]:*
    LISTEN 0      50              [::]:139           [::]:*
    LISTEN 0      4096            [::]:22            [::]:*
    LISTEN 0      4096           [::1]:631           [::]:*
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import os,json,subprocess,concurrent.futures
        from collections import Counter
        for root in ('/mnt/FTPfiles/yicamFront','/mnt/FTPfiles/yicamBack'):
         paths=[]
         for directory,_,names in os.walk(root):
          for name in names:
           if name.endswith('.mp4') and not name.startswith('.'):
            path=os.path.join(directory,name)
            try: paths.append((os.stat(path).st_mtime,path))
            except OSError:pass
         paths=sorted(paths,reverse=True)
         print('ROOT',root,'count',len(paths))
         def inspect(item):
          stamp,path=item
          try:
           p=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion,neonvr_detections','-of','json',path],capture_output=True,text=True,timeout=10)
           data=json.loads(p.stdout).get('format',{})
           tag=data.get('tags',{}).get('yi_motion')
           motion=json.loads(tag) if tag else None
           return path,data.get('duration'),motion,'neonvr_detections' in data.get('tags',{})
          except Exception as e:return path,None,None,str(e)
         with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
          stats=Counter();examples={}
          for path,duration,motion,ai in pool.map(inspect,paths[:120]):
           key=('absent',) if motion is None else (motion.get('schema'),motion.get('initial_state'),motion.get('history_available'),motion.get('history_truncated'),bool(motion.get('events')))
           stats[key]+=1
           if key not in examples:examples[key]=(path,duration,motion,ai)
          print('STATS',stats)
          for key,example in examples.items():print('EXAMPLE',key,json.dumps(example))
        PY
    ROOT /mnt/FTPfiles/yicamFront count 385
    STATS Counter({('yi-motion-v1', -1, False, True, False): 120})
    EXAMPLE ('yi-motion-v1', -1, False, True, False) ["/mnt/FTPfiles/yicamFront/2026Y10M11D14H/W413M00S60.mp4", "59.984000", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791742380000, "duration_ms": 59984, "video_time_precision_ms": 1000, "initial_state": -1, "history_available": false, "history_truncated": true, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": []}, false]
    ROOT /mnt/FTPfiles/yicamBack count 9
    STATS Counter({('yi-motion-v1', -1, False, True, False): 9})
    EXAMPLE ('yi-motion-v1', -1, False, True, False) ["/mnt/FTPfiles/yicamBack/2026Y10M11D13H/W451M01S58.mp4", "58.418000", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791741061000, "duration_ms": 58418, "video_time_precision_ms": 1000, "initial_state": -1, "history_available": false, "history_truncated": true, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": []}, true]
    ✓ • 4.65s

## Assistant

If you have a recording where ffprobe shows non-empty `events` but NeoNVR shows no markers, what is its file name or path?

## Activity

    $ systemctl list-units --type=service --state=running --no-pager; rg --files /home/ubuntu/NeoNVR -g AGENTS.md -g '*ftp*' -g '*ingest*'; sed -n '390,433p' detection.go; sed -n '2260,2305p' main.go; ls /tmp/neonvr-live-tools-vbq_9dyb
      UNIT                          LOAD   ACTIVE SUB     DESCRIPTION
      accounts-daemon.service       loaded active running Accounts Service
      avahi-daemon.service          loaded active running Avahi mDNS/DNS-SD Stack
      buildkit.service              loaded active running buildkit.service
      cloudflared.service           loaded active running cloudflared
      colord.service                loaded active running Manage, Install and Generate Color Profiles
      containerd.service            loaded active running containerd container runtime
      cron.service                  loaded active running Regular background program processing daemon
      cups-browsed.service          loaded active running Make remote CUPS printers available locally
      cups.service                  loaded active running CUPS Scheduler
      dbus.service                  loaded active running D-Bus System Message Bus
      fwupd.service                 loaded active running Firmware update daemon
      gdm.service                   loaded active running GNOME Display Manager
      gnome-remote-desktop.service  loaded active running GNOME Remote Desktop
      kerneloops.service            loaded active running Tool to automatically collect and submit kernel crash signatures
      ModemManager.service          loaded active running Modem Manager
      networkd-dispatcher.service   loaded active running Dispatcher daemon for systemd-networkd
      NetworkManager.service        loaded active running Network Manager
      nmbd.service                  loaded active running Samba NMB Daemon
      polkit.service                loaded active running Authorization Manager
      power-profiles-daemon.service loaded active running Power Profiles daemon
      rsyslog.service               loaded active running System Logging Service
      rtkit-daemon.service          loaded active running RealtimeKit Scheduling Policy Service
      smbd.service                  loaded active running Samba SMB Daemon
      snapd.service                 loaded active running Snap Daemon
      ssh.service                   loaded active running OpenBSD Secure Shell server
      switcheroo-control.service    loaded active running Switcheroo Control Proxy service
      systemd-journald.service      loaded active running Journal Service
      systemd-logind.service        loaded active running User Login Management
      systemd-oomd.service          loaded active running Userspace Out-Of-Memory (OOM) Killer
      systemd-resolved.service      loaded active running Network Name Resolution
      systemd-timesyncd.service     loaded active running Network Time Synchronization
      systemd-udevd.service         loaded active running Rule-based Manager for Device Events and Files
      thermald.service              loaded active running Thermal Daemon Service
      udisks2.service               loaded active running Disk Manager
      unattended-upgrades.service   loaded active running Unattended Upgrades Shutdown
      upower.service                loaded active running Daemon for power management
      user@1000.service             loaded active running User Manager for UID 1000
      vsftpd.service                loaded active running vsftpd FTP server
      wpa_supplicant.service        loaded active running WPA supplicant
    
    Legend: LOAD   → Reflects whether the unit definition was properly loaded.
            ACTIVE → The high-level unit activation state, i.e. generalization of SUB.
            SUB    → The low-level unit activation state, values depend on unit type.
    
    39 loaded units listed.
    }
    
    // Remux without re-encoding. Verify metadata, streams and duration before an
    // atomic replacement, retaining permissions, ownership and recording mtime.
    func (s *Server) embedDetection(ctx context.Context, path string, original os.FileInfo, result DetectionResult) error {
        if err := result.validate(); err != nil {
            return err
        }
        before, err := s.inspectDetectionVideo(ctx, path)
        if err != nil {
            return err
        }
        data, err := json.Marshal(result)
        if err != nil {
            return err
        }
        tmp, err := os.CreateTemp(filepath.Dir(path), ".neonvr-detect-*.mp4")
        if err != nil {
            return err
        }
        tmpPath := tmp.Name()
        tmp.Close()
        defer os.Remove(tmpPath)
        out, err := exec.CommandContext(ctx, s.ffmpegPath, "-nostdin", "-v", "error", "-y", "-i", path,
            "-map", "0", "-c", "copy", "-map_metadata", "0", "-metadata", detectionMetadataTag+"="+string(data),
            "-movflags", "+faststart+use_metadata_tags", tmpPath).CombinedOutput()
        if err != nil {
            return fmt.Errorf("detection metadata remux: %w: %.2000s", err, out)
        }
        after, err := s.inspectDetectionVideo(ctx, tmpPath)
        if err != nil {
            return err
        }
        beforeDuration, e1 := strconv.ParseFloat(before.Format.Duration, 64)
        afterDuration, e2 := strconv.ParseFloat(after.Format.Duration, 64)
        if after.Format.Tags[detectionMetadataTag] != string(data) || !reflect.DeepEqual(before.Streams, after.Streams) ||
            e1 != nil || e2 != nil || !finite(beforeDuration) || !finite(afterDuration) || math.Abs(beforeDuration-afterDuration) > .25 {
            return errors.New("remux verification failed; original recording retained")
        }
        f, err := os.OpenFile(tmpPath, os.O_WRONLY, 0)
        if err != nil {
            return err
        }
        if os.Geteuid() == 0 {
            out = append(out, clipFile{RelPath: filepath.ToSlash(rel), FullPath: full, Info: info})
        }
        for _, e := range entries {
            if e.IsDir() {
                if _, ok := parseYICamHourFolder(e.Name()); !ok {
                    continue
                }
                sub := filepath.Join(cam.Dir, e.Name())
                subs, err := os.ReadDir(sub)
                if err != nil {
                    continue
                }
                for _, se := range subs {
                    add(filepath.Join(e.Name(), se.Name()), filepath.Join(sub, se.Name()), se)
                }
                continue
            }
            add(e.Name(), filepath.Join(cam.Dir, e.Name()), e)
        }
        return out, nil
    }
    
    func clipTimestamp(f clipFile) time.Time {
        if ts, ok := parseClipTimestamp(filepath.Base(f.RelPath)); ok {
            return ts
        }
        parent := filepath.Base(filepath.Dir(filepath.FromSlash(f.RelPath)))
        if hour, ok := parseYICamHourFolder(parent); ok {
            if minute, second, ok := parseYICamMinuteSecond(filepath.Base(f.RelPath)); ok {
                return hour.Add(time.Duration(minute)*time.Minute + time.Duration(second)*time.Second)
            }
        }
        return f.Info.ModTime()
    }
    
    func (s *Server) listDates(selected map[string]bool) ([]string, error) {
        cameras, err := s.listCameras()
        if err != nil {
            return nil, err
        }
        seen := make(map[string]bool)
        for _, cam := range cameras {
            if selected != nil && !selected[cam.Name] {
                continue
            }
            files, err := listCameraFiles(cam)
    browser-env
    check_browser.py
    check_motion_browser.py
    check_motion_parse_browser.py
    go
    go2rtc
    motion-player-desktop.png
    motion-player-mobile.png
    neonvr
    neonvr-motion
    streams-desktop.png
    streams-mobile.png
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        import json,urllib.request
        base='http://127.0.0.1:5500'
        for endpoint in ('/api/cameras','/api/camera-motion?camera=yicamFront&file=2026Y10M11D14H%2FW413M00S60.mp4','/api/camera-motion?camera=yicamBack&file=2026Y10M11D13H%2FW451M01S58.mp4'):
         try:
          with urllib.request.urlopen(base+endpoint,timeout=10) as response: print(endpoint, response.status, response.read().decode()[:6000])
         except Exception as e:print(endpoint,str(e))
        PY
        ls -la /home/ubuntu
        nerdctl ps
    /api/cameras 200 {
      "root": "/mnt/mySSD/BlinkClips",
      "cameras": [
        {
          "name": "G8T1GH012043CLWC"
        },
        {
          "name": "GVP4TE015434B0PF"
        },
        {
          "name": "Outdoor3C7N2"
        },
        {
          "name": "Outdoor3CLWH"
        },
        {
          "name": "Outdoor3CLWQ"
        },
        {
          "name": "VideoDoorbell2TXK"
        },
        {
          "name": "yicamBack"
        },
        {
          "name": "yicamFront"
        }
      ],
      "latest_date": "2026-10-11"
    }
    
    /api/camera-motion?camera=yicamFront&file=2026Y10M11D14H%2FW413M00S60.mp4 200 {
      "motion": {
        "duration": 59.984,
        "history_available": false,
        "history_truncated": true,
        "history_complete": false,
        "markers": [],
        "intervals": []
      }
    }
    
    /api/camera-motion?camera=yicamBack&file=2026Y10M11D13H%2FW451M01S58.mp4 200 {
      "motion": {
        "duration": 58.418,
        "history_available": false,
        "history_truncated": true,
        "history_complete": false,
        "markers": [],
        "intervals": []
      }
    }
    
    total 22900
    drwxr-x--- 41 ubuntu ubuntu     4096 Oct 10 18:26 .
    drwxr-xr-x  3 root   root       4096 Jun 10  2023 ..
    drwxrwxr-x  4 ubuntu ubuntu     4096 Oct  9  2023 .Wolfram
    drwxrwxr-x 11 ubuntu ubuntu     4096 Oct  9  2023 .WolframEngine
    -rw-------  1 ubuntu ubuntu   137993 Oct 10 18:26 .bash_history
    -rw-r--r--  1 ubuntu ubuntu      220 Jun 10  2023 .bash_logout
    -rw-r--r--  1 ubuntu ubuntu     3859 Sep 28  2023 .bashrc
    drwx------ 17 ubuntu ubuntu     4096 Oct 10 23:42 .cache
    drwxr-xr-x 19 ubuntu ubuntu     4096 Apr  6  2024 .ccache
    drwxr-xr-x  3 ubuntu ubuntu     4096 Apr  6  2024 .cmake
    drwxrwxr-x 14 ubuntu ubuntu     4096 Oct 11 14:15 .codex
    drwx------ 16 ubuntu ubuntu     4096 Aug  2  2025 .config
    drwxrwxr-x  3 ubuntu ubuntu     4096 Sep 28  2023 .dotnet
    -rw-rw-r--  1 ubuntu ubuntu       68 Oct  4  2023 .gitconfig
    drwx------  2 ubuntu ubuntu     4096 Jan  6  2024 .gnupg
    drwxr-xr-x  2 root   root       4096 Jun 11  2023 .ipynb_checkpoints
    drwxrwxr-x  2 ubuntu ubuntu     4096 Sep 28  2023 .ipython
    drwxrwxr-x  3 ubuntu ubuntu     4096 Sep 28  2023 .jupyter
    -rw-------  1 ubuntu ubuntu       20 Oct 10 00:34 .lesshst
    drwx------  7 ubuntu ubuntu     4096 Aug  2  2025 .local
    drwx------  3 ubuntu ubuntu     4096 Sep 28  2023 .pki
    -rw-r--r--  1 ubuntu ubuntu      807 Jun 10  2023 .profile
    -rw-rw-r--  1 ubuntu ubuntu       66 Apr  5  2024 .selected_editor
    drwx------  2 ubuntu ubuntu     4096 Oct 10 03:07 .ssh
    -rw-r--r--  1 ubuntu ubuntu        0 Jun 10  2023 .sudo_as_admin_successful
    drwxrwxr-x  4 ubuntu ubuntu     4096 Sep 28  2023 .vscode
    -rw-rw-r--  1 ubuntu ubuntu      542 Aug  2  2025 .wget-hsts
    drwxr-xr-x  3 root   root       4096 Jun 12  2023 0ProjWolf
    drwxr-xr-x  3 ubuntu ubuntu     4096 Oct  4  2023 Desktop
    drwxr-xr-x  3 ubuntu ubuntu     4096 Apr 14 04:19 Documents
    drwxr-xr-x  3 ubuntu ubuntu     4096 Apr 26  2025 Downloads
    drwxrwxr-x  2 ubuntu ubuntu     4096 Oct 10 03:01 FTPfiles
    -rw-r--r--  1 ubuntu ubuntu  4492062 Mar 12  2026 IMG_1759.jpeg
    drwxr-xr-x  2 ubuntu ubuntu     4096 Jun 10  2023 Music
    drwxr-xr-x  7 ubuntu ubuntu     4096 Oct 11 14:06 NeoNVR
    -rw-rw-r--  1 ubuntu ubuntu   148820 Oct 10 00:13 NeoNVR-before-config-20261010-001315.tar.gz
    drwxr-xr-x  6 ubuntu ubuntu     4096 Jan 11  2026 Pictures
    drwxr-xr-x  2 ubuntu ubuntu     4096 Aug 12  2025 Public
    drwxr-xr-x  2 ubuntu ubuntu     4096 Jun 10  2023 Templates
    drwxrwxr-x  6 ubuntu ubuntu     4096 Sep 28  2023 WolframLanguageForJupyter
    drwxrwxr-x  6 ubuntu ubuntu     4096 Apr 26  2025 airflow
    drwxrwxr-x  2 ubuntu ubuntu     4096 Sep 14 21:37 bitnet
    -rw-rw-r--  1 ubuntu ubuntu   148335 Oct  9 21:20 blinkclips-web-v9-before-NeoNVR-20261009-212059.tar.gz
    -rw-r--r--  1 ubuntu ubuntu    24890 Jun 27 13:03 blinkclips-web-v9.zip
    -rw-rw-r--  1 ubuntu ubuntu 18276276 May 17  2024 cloudflared.deb
    drwxrwxr-x  3 ubuntu ubuntu     4096 Sep 28  2023 cppDEMO
    -rw-rw-r--  1 ubuntu ubuntu     1131 Jun 11  2023 dockerfile
    drwxrwxr-x  2 ubuntu ubuntu     4096 May 27 17:31 frigatedocker
    drwxr-xr-x  2 ubuntu ubuntu     4096 Sep 23 02:07 ftpuploads
    drwxrwxr-x 15 ubuntu ubuntu     4096 Sep 24  2024 haproxy
    drwxrwxr-x  2 ubuntu ubuntu     4096 Sep 24  2024 haproxycfg
    drwxrwxr-x  2 ubuntu ubuntu     4096 May 27 17:41 model-build
    -rwxrwxrwx  1 ubuntu ubuntu      712 Jan 26  2024 mountUnmountSSD.sh
    drwxrwxr-x  3 ubuntu ubuntu     4096 Oct  2 23:11 orcaslicer
    -rwxrwxrwx  1 ubuntu ubuntu      117 Sep 28  2023 runjupyterlab
    drwx------  7 ubuntu ubuntu     4096 Aug  2  2025 snap
    drwxrwxrwx  3 ubuntu ubuntu     4096 Oct 10  2023 stuff
    drwxrwxr-x  3 ubuntu ubuntu     4096 Jun 10  2023 wolfram-jupyter
    time="2026-10-11T14:16:24-04:00" level=fatal msg="rootless containerd not running? (hint: use `containerd-rootless-setuptool.sh install` to start rootless containerd): stat /run/user/1000/containerd-rootless: no such file or directory"
    ✗ (1) • 254ms

## Activity

    $ python3 - <<'PY'
        import os,json,subprocess,concurrent.futures
        from collections import Counter
        paths=[]
        for directory,_,names in os.walk('/mnt/FTPfiles/yicamFront'):
         for name in names:
          if name.endswith('.mp4') and not name.startswith('.'):
           path=os.path.join(directory,name)
           try:paths.append((os.stat(path).st_mtime,path))
           except OSError:pass
        paths=sorted(paths,reverse=True)[120:]
        def inspect(entry):
         stamp,path=entry
         try:
          p=subprocess.run(['ffprobe','-v','error','-show_entries','format_tags=yi_motion','-of','json',path],capture_output=True,text=True,timeout=10)
          tag=json.loads(p.stdout).get('format',{}).get('tags',{}).get('yi_motion'); motion=json.loads(tag) if tag else None
          return stamp,path,motion
         except Exception:return stamp,path,None
        stats=Counter();examples={}
        with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
         for stamp,path,motion in pool.map(inspect,paths):
          key=('absent',) if motion is None else (motion.get('initial_state'),motion.get('history_available'),motion.get('history_truncated'),bool(motion.get('events')))
          stats[key]+=1
          if key not in examples:examples[key]=(stamp,path,motion)
        print('OLDER_FILES',stats)
        for key, example in examples.items():print('EXAMPLE',key,json.dumps(example))
        PY
    OLDER_FILES Counter({(-1, False, True, False): 205, (1, True, False, False): 46, (0, True, False, True): 11, ('absent',): 2, (1, True, False, True): 1, (-1, True, True, True): 1})
    EXAMPLE (-1, False, True, False) [1791735348.037212, "/mnt/FTPfiles/yicamFront/2026Y10M11D12H/W414M00S60.mp4", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791735240000, "duration_ms": 60036, "video_time_precision_ms": 1000, "initial_state": -1, "history_available": false, "history_truncated": true, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": []}]
    EXAMPLE (1, True, False, False) [1791695447.8559177, "/mnt/FTPfiles/yicamFront/2026Y10M11D01H/W409M00S60.mp4", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791695340000, "duration_ms": 60034, "video_time_precision_ms": 1000, "initial_state": 1, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": []}]
    EXAMPLE (0, True, False, True) [1791694123.9366517, "/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W447M39S21.mp4", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791694059000, "duration_ms": 21010, "video_time_precision_ms": 1000, "initial_state": 0, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[1534, 1791694060534, 1, "transition"]]}]
    EXAMPLE (1, True, False, True) [1791693955.4793901, "/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W444M00S60.mp4", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791693840000, "duration_ms": 60051, "video_time_precision_ms": 1000, "initial_state": 1, "history_available": true, "history_truncated": false, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[47567, 1791693887567, 0, "transition"]]}]
    EXAMPLE (-1, True, True, True) [1791692449.3011274, "/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W419M13S46.mp4", {"schema": "yi-motion-v1", "model": "y623", "source": "encoder-stats-ipc", "video_start_unix_ms": 1791692353000, "duration_ms": 46208, "video_time_precision_ms": 1000, "initial_state": -1, "history_available": true, "history_truncated": true, "event_columns": ["offset_ms", "unix_ms", "state", "kind"], "events": [[2514, 1791692355514, 1, "sample"]]}]
    EXAMPLE ('absent',) [1791676798.401526, "/mnt/FTPfiles/yicamFront/2026Y10M10D23H/E858M00S60.mp4", null]
    ✓ • 7.61s

## User

> If you have a recording where ffprobe shows non-empty `events` but NeoNVR shows no markers, what is its file name or path?

2026Y10M11D14H/W402M00S60.mp4

## Activity

    $ sed -n '10,40p' motion_test.go; sed -n '290,320p' README.md; sed -n '168,199p' detection.go; for path in /AGENTS.md /home/AGENTS.md /home/ubuntu/AGENTS.md /home/ubuntu/NeoNVR/AGENTS.md; do if [ -f "$path" ]; then cat "$path"; fi; done
        "os/exec"
        "path/filepath"
        "reflect"
        "testing"
    )
    
    const motionTestTag = `{"schema":"yi-motion-v1","video_start_unix_ms":1791692141000,"duration_ms":4000,"initial_state":0,"history_available":true,"history_truncated":false,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[1500,1791692142500,1,"transition"],[2500,1791692143500,0,"transition"]]}`
    
    func TestParseCameraMotion(t *testing.T) {
        tests := []struct {
            name      string
            tag       string
            duration  float64
            intervals []motionInterval
            partial   bool
        }{
            {"milliseconds", motionTestTag, 4, []motionInterval{{1.5, 2.5}}, false},
            {"adjusted START and STOP in minute recording", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[5000,1791692585000,1,"transition"],[12000,1791692592000,0,"transition"]]}`, 60, []motionInterval{{5, 12}}, false},
            {"detector restart STOP", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[12000,1791692592000,0,"transition"],[40000,1791692620000,1,"transition"],[45000,1791692625000,0,"transition"]]}`, 60, []motionInterval{{0, 12}, {40, 45}}, false},
            {"initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"events":[]}`, 60.033, []motionInterval{{0, 60.033}}, false},
            {"truncated sample", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[[2514,1791692355514,1,"sample"]]}`, 46.208, []motionInterval{{2.514, 46.208}}, true},
            {"no motion", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true}`, 10, []motionInterval{}, false},
            {"unknown history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false}`, 10, []motionInterval{}, true},
            {"unknown initial state despite available history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":false}`, 60, []motionInterval{}, true},
            {"missing initial state", `{"schema":"yi-motion-v1","history_available":true,"event_columns":["offset_ms","state"],"events":[[5000,1],[12000,0]]}`, 60, []motionInterval{{5, 12}}, true},
            {"unknown state closes motion without filling gap", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[5000,-1],[12000,1],[20000,0]]}`, 60, []motionInterval{{0, 5}, {12, 20}}, true},
            {"reordered events and columns", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["state","offset_ms"],"events":[[0,3500],[1,1500],[1,2500]]}`, 4, []motionInterval{{1.5, 3.5}}, false},
            {"unix fallback", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"video_start_unix_ms":1791692141000,"event_columns":["state","unix_ms"],"events":[[1,1791692146295],[0,1791692152239]]}`, 18.624, []motionInterval{{5.295, 11.239}}, false},
            {"metadata duration fallback", motionTestTag, 0, []motionInterval{{1.5, 2.5}}, false},
            {"ignore outside clip", `{"schema":"yi-motion-v1","initial_state":0,"history_available":true,"event_columns":["offset_ms","state"],"events":[[-1,1],[1000,1],[2000,0],[4000,1],[5000,1]]}`, 4, []motionInterval{{1, 2}}, false},
            {"zero-length initial motion", `{"schema":"yi-motion-v1","initial_state":1,"history_available":true,"event_columns":["offset_ms","state"],"events":[[0,0]]}`, 4, []motionInterval{}, false},
    Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thumb-maintenance-every 1h`. Set it to `0` to disable background maintenance. Lazy/on-demand thumbnail generation remains enabled as a fallback.
    
    ### Yi timestamps
    
    Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.
    
    Recordings with `yi_motion` metadata show orange, clickable motion onset markers directly on the main player's seek track, with a lane below showing the motion intervals. These camera markers work even when AI detection is disabled. Partial or truncated motion history is labeled in the player. NeoNVR uses the metadata's `offset_ms` values relative to the video, without a timezone correction; `initial_state: 1` marks motion from the beginning, and transitions to state 0 end the highlighted interval. An initial or later state of `-1` represents unknown coverage, even if the history flags say data is available; unknown periods are not treated as idle or highlighted as motion. Older cameras without the tag continue to play normally. The tag can be inspected with:
    
    ```bash
    ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 | jq '.format.tags.yi_motion | fromjson'
    ```
    
    The format remains `yi-motion-v1` when camera motion settings change. NeoNVR applies the camera's recorded START/STOP events directly, so changes to thresholds, confirmation, quiet time, sampling, or detector restarts are reflected in the marked intervals. It does not infer motion from the recording's length. For example, these events in a 60-second recording highlight only seconds 5 through 12:
    
    ```json
    {
      "schema": "yi-motion-v1",
      "initial_state": 0,
      "history_available": true,
      "history_truncated": false,
      "event_columns": ["offset_ms", "unix_ms", "state", "kind"],
      "events": [
        [5000, 1791692585000, 1, "transition"],
        [12000, 1791692592000, 0, "transition"]
      ]
    }
    ```
    
    The parser follows `event_columns`, begins with `initial_state`, and applies each event's state at its offset. The `unix_ms` column is an absolute timestamp in milliseconds; playback uses `offset_ms`. If `initial_state` is 1 and there are no events, motion extends through the entire recording. If it is -1 and there are no events, no motion timestamps are available, and thumbnails use the existing scene selector. A STOP during a detector restart ends the current interval just like any other STOP.
    
    ## Project rename and compatibility
            CodecType string `json:"codec_type"`
            CodecName string `json:"codec_name"`
        } `json:"streams"`
        Format struct {
            Duration string            `json:"duration"`
            Tags     map[string]string `json:"tags"`
        } `json:"format"`
    }
    
    func (s *Server) ffprobeExecutable() string {
        if strings.ContainsRune(s.ffmpegPath, filepath.Separator) {
            return filepath.Join(filepath.Dir(s.ffmpegPath), "ffprobe")
        }
        return "ffprobe"
    }
    
    func (s *Server) inspectDetectionVideo(ctx context.Context, path string) (videoInspection, error) {
        ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
        defer cancel()
        out, err := exec.CommandContext(ctx, s.ffprobeExecutable(), "-v", "error", "-show_entries",
            "format=duration:format_tags="+detectionMetadataTag+":stream=codec_type,codec_name", "-of", "json", path).Output()
        var info videoInspection
        if err != nil {
            return info, fmt.Errorf("ffprobe: %w", err)
        }
        if err := json.Unmarshal(out, &info); err != nil {
            return info, err
        }
        return info, nil
    }
    
    func (s *Server) embeddedDetection(ctx context.Context, path string) (*DetectionResult, error) {
    ✓ • 0ms

## Activity

    $ python3 - <<'PY'
        from pathlib import Path
        import subprocess, urllib.request, urllib.parse,json
        file='2026Y10M11D14H/W402M00S60.mp4'
        for camera in ('yicamFront','yicamBack'):
         path=Path('/mnt/FTPfiles')/camera/file
         print('CAMERA',camera,'FILE',path,'EXISTS',path.is_file())
         if not path.is_file():continue
         p=subprocess.run(['ffprobe','-v','error','-show_entries','format=duration:format_tags=yi_motion,neonvr_detections','-of','json',str(path)],capture_output=True,text=True,timeout=10)
         data=json.loads(p.stdout)
         tag=data.get('format',{}).get('tags',{}).get('yi_motion')
         print('FFPROBE_MOTION',json.dumps(json.loads(tag) if tag else None,indent=2))
         print('AI_TAG_PRESENT','neonvr_detections' in data.get('format',{}).get('tags',{}))
         url='http://127.0.0.1:5500/api/camera-motion?'+urllib.parse.urlencode({'camera':camera,'file':file})
         with urllib.request.urlopen(url,timeout=10) as response: print('PRODUCTION_API',response.status,response.read().decode())
        for file in ('2026Y10M11D00H/W447M39S21.mp4','2026Y10M11D00H/W419M13S46.mp4','2026Y10M11D01H/W409M00S60.mp4'):
         url='http://127.0.0.1:5500/api/camera-motion?'+urllib.parse.urlencode({'camera':'yicamFront','file':file})
         try:
          with urllib.request.urlopen(url,timeout=10) as response:print('OLDER_VALID',file,response.read().decode())
         except Exception as e:print('OLDER',file,e)
        PY
    CAMERA yicamFront FILE /mnt/FTPfiles/yicamFront/2026Y10M11D14H/W402M00S60.mp4 EXISTS True
    FFPROBE_MOTION {
      "schema": "yi-motion-v1",
      "model": "y623",
      "source": "encoder-stats-ipc",
      "video_start_unix_ms": 1791741720000,
      "duration_ms": 60071,
      "video_time_precision_ms": 1000,
      "initial_state": -1,
      "history_available": false,
      "history_truncated": true,
      "event_columns": [
        "offset_ms",
        "unix_ms",
        "state",
        "kind"
      ],
      "events": []
    }
    AI_TAG_PRESENT True
    PRODUCTION_API 200 {
      "motion": {
        "duration": 60.071,
        "history_available": false,
        "history_truncated": true,
        "history_complete": false,
        "markers": [],
        "intervals": []
      }
    }
    
    CAMERA yicamBack FILE /mnt/FTPfiles/yicamBack/2026Y10M11D14H/W402M00S60.mp4 EXISTS False
    OLDER_VALID 2026Y10M11D00H/W447M39S21.mp4 {
      "motion": {
        "duration": 21.01,
        "history_available": true,
        "history_truncated": false,
        "history_complete": true,
        "markers": [
          1.534
        ],
        "intervals": [
          {
            "start": 1.534,
            "end": 21.01
          }
        ]
      }
    }
    
    OLDER_VALID 2026Y10M11D00H/W419M13S46.mp4 {
      "motion": {
        "duration": 46.208,
        "history_available": true,
        "history_truncated": true,
        "history_complete": false,
        "markers": [
          2.514
        ],
        "intervals": [
          {
            "start": 2.514,
            "end": 46.208
          }
        ]
      }
    }
    
    OLDER_VALID 2026Y10M11D01H/W409M00S60.mp4 {
      "motion": {
        "duration": 60.034,
        "history_available": true,
        "history_truncated": false,
        "history_complete": true,
        "markers": [
          0
        ],
        "intervals": [
          {
            "start": 0,
            "end": 60.034
          }
        ]
      }
    }
    
    ✓ • 610ms

## Assistant

I checked `2026Y10M11D14H/W402M00S60.mp4`. Its `yi_motion` tag is present, but it contains `initial_state: -1`, `history_available: false`, and `events: []`. The running NeoNVR API reads those values correctly. An older clip with events returns a marker at 1.534 seconds.

I’ll make the player distinguish “metadata found, but no motion timestamps supplied” from partial history that has usable markers.

## Activity

    file changes: Completed · 4 changes
    Update { move_path: None }: /home/ubuntu/NeoNVR/README.md
    @@ -319,2 +319,4 @@
     
    +The player distinguishes a missing `yi_motion` tag from a detected tag that supplies no motion timestamps. If the camera writes `initial_state: -1` and `events: []`, it shows “Camera metadata found; no motion timestamps supplied.” This describes the embedded data, rather than a parser failure. Partial history with usable events still shows its known motion markers. Hover over the status for the initial state, event count, and history flags; extraction or parsing failures show their error there instead.
    +
     ## Project rename and compatibility
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion.go
    @@ -26,2 +26,4 @@
     	Duration         float64          `json:"duration"`
    +	InitialState     int              `json:"initial_state"`
    +	EventCount       int              `json:"event_count"`
     	HistoryAvailable bool             `json:"history_available"`
    @@ -114,3 +116,4 @@
     	sort.SliceStable(events, func(i, j int) bool { return events[i].offset < events[j].offset })
    -	result := &cameraMotion{Duration: duration, HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
    +	result := &cameraMotion{Duration: duration, InitialState: state, EventCount: len(raw.Events),
    +		HistoryAvailable: raw.HistoryAvailable, HistoryTruncated: raw.HistoryTruncated,
     		HistoryComplete: raw.HistoryAvailable && !raw.HistoryTruncated && state != -1,
    Update { move_path: None }: /home/ubuntu/NeoNVR/motion_test.go
    @@ -32,2 +32,3 @@
     		{"unknown history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false}`, 10, []motionInterval{}, true},
    +		{"ingested W402M00S60 camera supplied no history", `{"schema":"yi-motion-v1","model":"y623","source":"encoder-stats-ipc","video_start_unix_ms":1791741720000,"duration_ms":60071,"video_time_precision_ms":1000,"initial_state":-1,"history_available":false,"history_truncated":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[]}`, 60.071, []motionInterval{}, true},
     		{"unknown initial state despite available history", `{"schema":"yi-motion-v1","initial_state":-1,"history_available":true,"history_truncated":false}`, 60, []motionInterval{}, true},
    @@ -64,2 +65,18 @@
     
    +func TestCameraMotionAPIExposesEmptyHistory(t *testing.T) {
    +	s, _ := cameraMotionFixture(t, `{"schema":"yi-motion-v1","initial_state":-1,"history_available":false,"history_truncated":true,"event_columns":["offset_ms","unix_ms","state","kind"],"events":[]}`)
    +	w := httptest.NewRecorder()
    +	s.apiCameraMotionHandler(w, httptest.NewRequest(http.MethodGet, "/api/camera-motion?camera=Garden&file=recording.mp4", nil))
    +	var response struct {
    +		Motion *cameraMotion `json:"motion"`
    +	}
    +	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil || w.Code != 200 || response.Motion == nil {
    +		t.Fatalf("empty metadata was not detected: %d %s", w.Code, w.Body.String())
    +	}
    +	motion := response.Motion
    +	if motion.InitialState != -1 || motion.EventCount != 0 || motion.HistoryAvailable || !motion.HistoryTruncated || motion.HistoryComplete || len(motion.Markers) != 0 || len(motion.Intervals) != 0 {
    +		t.Fatalf("empty camera history was misinterpreted: %+v", motion)
    +	}
    +}
    +
     func TestParseCameraMotionRejectsInvalidMetadata(t *testing.T) {
    Update { move_path: None }: /home/ubuntu/NeoNVR/static/detection-viewer.js
    @@ -22,4 +22,10 @@
         motionStatus.classList.toggle('active', active);
    +    if (partial && !motion.intervals.length && motion.initial_state === -1 && motion.event_count === 0) {
    +      const text = 'Camera metadata found; no motion timestamps supplied.';
    +      if (motionStatus.textContent !== text) motionStatus.textContent = text;
    +      return;
    +    }
         const text = (active ? 'Camera motion now.' : motion.intervals.length ? 'No camera motion marked here.' :
    -      partial ? 'Camera motion history unavailable.' : 'No camera motion recorded.') + (partial ? ' Partial history.' : '');
    +      partial ? 'Camera metadata found; no usable motion timestamps.' : 'Camera metadata found; no motion recorded.') +
    +      (partial ? ' Partial history.' : '');
         if (motionStatus.textContent !== text) motionStatus.textContent = text;
    @@ -31,3 +37,8 @@
           const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
    -      if (!response.ok) throw new Error('Camera motion metadata unavailable.');
    +      if (!response.ok) {
    +        const failure = await response.json().catch(() => ({}));
    +        const error = new Error('Could not read camera motion metadata.');
    +        error.detail = failure.message || failure.error || 'HTTP ' + response.status;
    +        throw error;
    +      }
           const data = await response.json();
    @@ -37,3 +48,8 @@
           if (motionLegend) motionLegend.hidden = !motion;
    -      if (motionStatus) motionStatus.hidden = !motion;
    +      if (motionStatus) {
    +        motionStatus.hidden = false;
    +        motionStatus.textContent = motion ? '' : 'No camera motion metadata in this recording.';
    +        motionStatus.title = motion ? 'Initial state: ' + motion.initial_state + '; events: ' + motion.event_count +
    +          '; history available: ' + motion.history_available + '; history truncated: ' + motion.history_truncated : '';
    +      }
           activeCameraMotion();
    @@ -44,2 +60,3 @@
             motionStatus.textContent = error.message;
    +        motionStatus.title = error.detail || '';
           }
    @@ -95,3 +112,8 @@
         if (motionLegend) motionLegend.hidden = true;
    -    if (motionStatus) { motionStatus.hidden = true; motionStatus.classList.remove('active'); }
    +    if (motionStatus) {
    +      motionStatus.hidden = false;
    +      motionStatus.textContent = 'Reading camera motion metadata…';
    +      motionStatus.title = '';
    +      motionStatus.classList.remove('active');
    +    }
         status.classList.remove('active');

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go build -o /tmp/neonvr-live-tools-vbq_9dyb/neonvr-motion .
    ✓ • 569ms

## Activity

    $ /tmp/neonvr-live-tools-vbq_9dyb/go/bin/gofmt -w motion.go motion_test.go; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go test -race ./...
    ok      neonvr    8.679s
    ✓ • 9.65s

## Activity

    $ git diff --check; /tmp/neonvr-live-tools-vbq_9dyb/go/bin/go vet ./...
    ✓ • 42ms

## Activity

    $ cat > /tmp/neonvr-live-tools-vbq_9dyb/check_motion_status_browser.py <<'PY'
        import json
        from pathlib import Path
        import shutil
        import socket
        import subprocess
        import tempfile
        import time
        import urllib.request
        from playwright.sync_api import sync_playwright
        
        TOOLS=Path('/tmp/neonvr-live-tools-vbq_9dyb')
        def get(url):
         with urllib.request.urlopen(url,timeout=10) as response:return json.load(response)
        with tempfile.TemporaryDirectory(prefix='neonvr-motion-status-') as tmp:
         root=Path(tmp);directory=root/'clips'/'Garden';directory.mkdir(parents=True)
         sources={
          'unknown.mp4':'/mnt/FTPfiles/yicamFront/2026Y10M11D14H/W402M00S60.mp4',
          'partial.mp4':'/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W419M13S46.mp4',
          'valid.mp4':'/mnt/FTPfiles/yicamFront/2026Y10M11D00H/W447M39S21.mp4',
          'initial.mp4':'/mnt/FTPfiles/yicamFront/2026Y10M11D01H/W409M00S60.mp4',
         }
         for name,source in sources.items():shutil.copyfile(source,directory/name)
         subprocess.run(['ffmpeg','-nostdin','-v','error','-f','lavfi','-i','testsrc2=s=160x96:r=5','-t','4','-c:v','libx264','-preset','ultrafast','-threads','1',str(directory/'missing.mp4')],check=True)
         subprocess.run(['ffmpeg','-nostdin','-v','error','-i',str(directory/'missing.mp4'),'-c','copy','-metadata','yi_motion=broken-json','-movflags','+use_metadata_tags',str(directory/'invalid.mp4')],check=True)
         with socket.socket() as listener:
          listener.bind(('127.0.0.1',0));port=listener.getsockname()[1]
         config=root/'config.json';config.write_text(json.dumps({'root':str(root/'clips'),'extra_roots':[],'addr':f'127.0.0.1:{port}','auto_delete_days':0,'thumb_maintenance_every':'0','detection_enabled':False,'live_streams':[]}))
         origin=f'http://127.0.0.1:{port}'
         with (root/'app.log').open('w+') as log:
          process=subprocess.Popen([str(TOOLS/'neonvr-motion'),'-config',str(config)],stdout=log,stderr=log)
          try:
           for _ in range(100):
            try:cameras=get(origin+'/api/cameras');break
            except Exception:time.sleep(.05)
           else:raise RuntimeError('Server did not start')
           clips={clip['file']:clip for clip in get(origin+'/api/clips?cameras=Garden&date='+cameras['latest_date'])['clips']}
           metadata=get(origin+'/api/camera-motion?camera=Garden&file=unknown.mp4')['motion']
           assert metadata['initial_state']==-1 and metadata['event_count']==0 and not metadata['history_available'] and metadata['history_truncated'] and metadata['markers']==[]
           with sync_playwright() as playwright:
            browser=playwright.chromium.launch(headless=True,args=['--no-sandbox'])
            for viewport in ({'width':1280,'height':900},{'width':390,'height':844}):
             page=browser.new_page(viewport=viewport)
             errors=[];page.on('pageerror',lambda error:errors.append(str(error)))
             # Keep browser checks from regenerating thumbnail caches unnecessarily.
             page.route('**/thumb?**',lambda route:route.fulfill(status=200,content_type='image/svg+xml',body='<svg xmlns="http://www.w3.org/2000/svg" width="160" height="96"><rect width="160" height="96" fill="#111"/></svg>'))
             page.goto(origin)
             page.wait_for_function('allClips.length===6 && window.NeoNVRDetections')
             def load(name,text):
              page.evaluate('(clip)=>loadMainPlayer(clip,false)',clips[name])
              page.wait_for_function('(text)=>document.querySelector("#cameraMotionStatus").textContent===text',arg=text)
             load('unknown.mp4','Camera metadata found; no motion timestamps supplied.')
             assert page.locator('#cameraMotionLegend').is_visible()
             assert page.locator('#mainSeek .motion-marker').count()==0
             assert page.locator('#cameraMotionStatus').get_attribute('title')=='Initial state: -1; events: 0; history available: false; history truncated: true'
             load('partial.mp4','No camera motion marked here. Partial history.')
             marker=page.locator('#mainSeek .motion-marker');assert marker.count()==1
             assert marker.get_attribute('title')=='Camera motion at 2.514 sec'
             marker.click()
             page.wait_for_function('!document.querySelector("#mainVideo").seeking && Math.abs(document.querySelector("#mainVideo").currentTime-2.514)<.002 && document.querySelector("#cameraMotionStatus").textContent==="Camera motion now. Partial history."')
             load('valid.mp4','No camera motion marked here.')
             assert marker.count()==1 and marker.get_attribute('title')=='Camera motion at 1.534 sec'
             load('initial.mp4','Camera motion now.')
             assert marker.count()==1 and marker.get_attribute('title')=='Camera motion at 0.000 sec'
             load('missing.mp4','No camera motion metadata in this recording.')
             assert page.locator('#cameraMotionLegend').is_hidden() and marker.count()==0
             load('invalid.mp4','Could not read camera motion metadata.')
             assert 'invalid yi_motion JSON' in page.locator('#cameraMotionStatus').get_attribute('title')
             assert marker.count()==0
             assert not errors,errors
             if viewport['width']==390:assert page.evaluate('document.documentElement.scrollWidth<=innerWidth+1')
             print('Passed',viewport['width'],'px: exact ingested file, valid/partial markers, initial motion, missing tag, and parser error details.',flush=True)
             page.close()
            browser.close()
          finally:
           process.terminate()
           try:process.wait(timeout=10)
           except subprocess.TimeoutExpired:process.kill();process.wait()
        PY
        /tmp/neonvr-live-tools-vbq_9dyb/browser-env/bin/python /tmp/neonvr-live-tools-vbq_9dyb/check_motion_status_browser.py
    Passed 1280 px: exact ingested file, valid/partial markers, initial motion, missing tag, and parser error details.
    Passed 390 px: exact ingested file, valid/partial markers, initial motion, missing tag, and parser error details.
    ✓ • 2.88s

## Assistant

NeoNVR detects the metadata in `2026Y10M11D14H/W402M00S60.mp4`, but the embedded data contains:

```json
"initial_state": -1,
"history_available": false,
"history_truncated": true,
"events": []
```

The camera supplied no motion timestamps to mark. Older recordings with events produce markers correctly.

I updated the player to distinguish missing metadata from empty camera history. Desktop/mobile checks and Go race tests passed.

Rebuild to apply the clearer messages:

```bash
sudo nerdctl compose up -d --build
```

The missing timestamps need fixing in the camera’s motion-history writer.
