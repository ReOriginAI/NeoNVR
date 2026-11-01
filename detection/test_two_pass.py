import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

from worker import Detector, add_region, video_frames


class RegionTests(unittest.TestCase):
    def test_padding_clamping_and_merging(self):
        regions = []
        for start, end in [(0, 1), (3, 4), (20, 21), (29, 30)]:
            add_region(regions, start, end, 2, 3, 30)
        self.assertEqual(regions, [{'start': 0, 'end': 7}, {'start': 18, 'end': 24}, {'start': 27, 'end': 30}])


@unittest.skipUnless(shutil.which('ffmpeg') and shutil.which('ffprobe'), 'FFmpeg and FFprobe required')
class TwoPassTests(unittest.TestCase):
    def make_clip(self, folder, moving):
        path = str(Path(folder) / 'clip.mp4')
        filters = "drawbox=x=50:y=30:w=40:h=40:color=red:t=fill:enable='between(t,12,13)'" if moving else 'null'
        subprocess.run(['ffmpeg', '-nostdin', '-v', 'error', '-f', 'lavfi', '-i',
                        'color=c=black:s=160x96:r=10', '-vf', filters, '-t', '24',
                        '-c:v', 'libx264', '-threads', '1', path], check=True, capture_output=True, timeout=30)
        return {'path': path, 'ffmpeg': shutil.which('ffmpeg'), 'ffprobe': shutil.which('ffprobe'),
                'fps': 2, 'motion_fps': 5, 'motion_threshold': 16, 'motion_min_area': .001,
                'roi_pre_roll': 2, 'roi_post_roll': 3, 'confidence': .65,
                'min_frames': 3, 'merge_gap': .5, 'car_movement': .025, 'policy': 'b' * 64}

    def detector(self, predictions):
        detector = Detector.__new__(Detector)
        detector.width = detector.height = 64
        detector.model_name = 'yolo26m'
        detector.model_sha256 = 'a' * 64
        detector.predict = predictions
        return detector

    def test_static_clip_only_runs_yolo_in_opening_window(self):
        calls = []
        detector = self.detector(lambda raw, cutoff: calls.append(raw) or [])
        with tempfile.TemporaryDirectory() as folder:
            result = detector.analyze(self.make_clip(folder, False))
        self.assertEqual(result['version'], 2)
        self.assertEqual(result['motion_frames'], 120)
        self.assertEqual(result['regions'], [{'start': 0, 'end': 4, 'frames': 8}])
        self.assertEqual(result['frames'], len(calls))
        self.assertEqual(result['frames'], 8)  # 48 would be a full YOLO scan.

    def test_motion_pass_selects_padded_windows_and_yolo_seeks_only_those(self):
        detector = self.detector(lambda raw, cutoff: [(0, [.2, .2, .4, .5], .9)])
        with tempfile.TemporaryDirectory() as folder:
            request = self.make_clip(folder, True)
            with patch('worker.video_frames', wraps=video_frames) as decoder:
                result = detector.analyze(request)
        self.assertEqual(result['motion_frames'], 120)
        self.assertGreaterEqual(len(result['regions']), 2)
        self.assertLess(result['frames'], 48)
        self.assertEqual(result['frames'], sum(r['frames'] for r in result['regions']))
        yolo_calls = [call.args for call in decoder.call_args_list if call.args[2] == 64]
        self.assertEqual([(args[-2], args[-1]) for args in yolo_calls],
                         [(r['start'], r['end']) for r in result['regions']])
        moving = result['regions'][1]
        self.assertLessEqual(moving['start'], 10.2)
        self.assertGreater(moving['end'], 15)
        for interval in result['intervals']:
            self.assertTrue(any(interval['start'] >= r['start'] and interval['end'] <= r['end'] for r in result['regions']))
        self.assertTrue(any(row['start'] >= moving['start'] for row in result['intervals']))

    def test_motion_decode_failure_does_not_return_complete_analysis(self):
        detector = self.detector(lambda raw, cutoff: self.fail('YOLO ran before motion pass succeeded'))
        with tempfile.TemporaryDirectory() as folder:
            request = self.make_clip(folder, False)
            def broken(*args):
                yield 0, bytes(320 * 180 * 3)
                raise ValueError('decode failed')
            with patch('worker.video_frames', broken):
                with self.assertRaisesRegex(ValueError, 'decode failed'):
                    detector.analyze(request)

    def test_yolo_decode_failure_does_not_return_complete_analysis(self):
        detector = self.detector(lambda raw, cutoff: [])
        with tempfile.TemporaryDirectory() as folder:
            request = self.make_clip(folder, False)
            def broken_yolo(req, fps, width, height, start, end):
                if width == 320:
                    yield from video_frames(req, fps, width, height, start, end)
                else:
                    raise ValueError('ROI decode failed')
            with patch('worker.video_frames', broken_yolo):
                with self.assertRaisesRegex(ValueError, 'ROI decode failed'):
                    detector.analyze(request)
