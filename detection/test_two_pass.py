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
        self.assertEqual(result['version'], 3)
        self.assertEqual(result['motion_frames'], 48)
        self.assertEqual(result['regions'], [{'start': 0, 'end': 4, 'frames': 8}])
        self.assertEqual(result['frames'], len(calls))
        self.assertEqual(result['frames'], 8)  # 48 would be a full YOLO scan.

    def test_motion_pass_selects_padded_windows_and_yolo_seeks_only_those(self):
        detector = self.detector(lambda raw, cutoff: [(0, [.2, .2, .4, .5], .9)])
        with tempfile.TemporaryDirectory() as folder:
            request = self.make_clip(folder, True)
            with patch('worker.video_frames', wraps=video_frames) as decoder:
                result = detector.analyze(request)
        self.assertEqual(result['motion_frames'], 48)
        self.assertGreaterEqual(len(result['regions']), 2)
        self.assertLess(result['frames'], 48)
        self.assertEqual(result['frames'], sum(r['frames'] for r in result['regions']))
        self.assertEqual(decoder.call_count, 1)
        self.assertEqual(decoder.call_args.args[1:], (2, 64, 64, 0, 24.0))
        self.assertEqual(result['motion_fps'], result['sample_fps'])
        self.assertLessEqual(result['peak_buffer_frames'], 6)
        moving = result['regions'][1]
        self.assertLessEqual(moving['start'], 10.2)
        self.assertGreater(moving['end'], 15)
        for interval in result['intervals']:
            self.assertTrue(any(interval['start'] >= r['start'] and interval['end'] <= r['end'] for r in result['regions']))
        self.assertTrue(any(row['start'] >= moving['start'] for row in result['intervals']))

    def test_motion_decode_failure_does_not_return_complete_analysis(self):
        detector = self.detector(lambda raw, cutoff: [])
        with tempfile.TemporaryDirectory() as folder:
            request = self.make_clip(folder, False)
            def broken(*args):
                yield 0, bytes(64 * 64 * 3)
                raise ValueError('decode failed')
            with patch('worker.video_frames', broken):
                with self.assertRaisesRegex(ValueError, 'decode failed'):
                    detector.analyze(request)

    def test_inference_failure_does_not_publish_complete_result(self):
        def broken_prediction(raw, cutoff):
            raise ValueError('inference failed')
        detector = self.detector(broken_prediction)
        with tempfile.TemporaryDirectory() as folder:
            with self.assertRaisesRegex(ValueError, 'inference failed'):
                detector.analyze(self.make_clip(folder, False))


class SharedStreamTests(unittest.TestCase):
    def run_shared(self, events, pre=2, post=3, fps=2, duration=24):
        import json
        from types import SimpleNamespace
        selected = []
        class Gate:
            def __init__(self, request, duration, width, height, fps):
                self.pre = pre
                self.count = 0
                self.regions = []
                add_region(self.regions, 0, min(1, duration), pre, post, duration)
            def observe(self, timestamp, raw):
                if timestamp in events:
                    add_region(self.regions, timestamp, min(duration, timestamp+1/fps), pre, post, duration)
                self.count += 1
        def samples(*args):
            for i in range(round(duration*fps)):
                yield i/fps, bytes([i]) * (64*64*3)
        detector = Detector.__new__(Detector)
        detector.width = detector.height = 64
        detector.model_name = 'yolo26m'
        detector.model_sha256 = 'a'*64
        detector.predict = lambda raw, cutoff: selected.append(raw[0]) or []
        request = {'path':'unused','ffprobe':'ffprobe','fps':fps,'confidence':.65,
                   'min_frames':3,'merge_gap':.5,'car_movement':.025,'policy':'b'*64}
        with patch('worker.MotionGate',Gate), patch('worker.video_frames',samples), patch('worker.subprocess.run',return_value=SimpleNamespace(stdout=json.dumps({'streams':[{'duration':str(duration)}]}))):
            result = detector.analyze(request)
        expected = [i for i in range(round(duration*fps)) if any(r['start'] <= i/fps < r['end'] for r in result['regions'])]
        self.assertEqual(selected, expected)
        self.assertEqual(len(set(selected)), len(selected))
        self.assertLessEqual(result['peak_buffer_frames'], __import__('math').ceil(pre*fps)+2)
        return result

    def test_pre_roll_post_roll_merging_and_no_duplicate_inference(self):
        result = self.run_shared({7,12,20})
        self.assertEqual([(r['start'],r['end']) for r in result['regions']],[(0,4),(5,15.5),(18,23.5)])
        self.assertEqual(result['frames'],sum(r['frames'] for r in result['regions']))

    def test_zero_padding_and_end_of_clip_motion(self):
        result = self.run_shared({23.5},pre=0,post=0)
        self.assertEqual(result['regions'][-1],{'start':23.5,'end':24,'frames':1})

    def test_fractional_padding(self):
        self.run_shared({7,12,20},pre=1.3,post=.7)

    def test_shared_half_frame_per_second(self):
        result = self.run_shared({12,20},fps=.5)
        self.assertEqual(result['motion_fps'],.5)
        self.assertEqual(result['motion_frames'],12)
