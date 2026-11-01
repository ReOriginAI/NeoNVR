import io
import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path
from types import SimpleNamespace
from unittest.mock import patch

from worker import Detector, HardwareDecodeError, video_frames


class IntelDecodeTests(unittest.TestCase):
    def test_hardware_options_are_before_input_and_frames_return_to_cpu(self):
        process = SimpleNamespace(stdout=io.BytesIO(bytes(2 * 2 * 3)),
                                  wait=lambda timeout=None: 0, poll=lambda: 0, kill=lambda: None)
        request = {'ffmpeg': 'ffmpeg', 'path': 'clip.mp4', 'intel_acceleration': True,
                   'vaapi_device': '/dev/dri/renderD129'}
        with patch('worker.subprocess.Popen', return_value=process) as popen:
            frames = list(video_frames(request, 1, 2, 2, 5, 6))
        command = popen.call_args.args[0]
        self.assertLess(command.index('-hwaccel'), command.index('-i'))
        self.assertEqual(command[command.index('-hwaccel_device') + 1], '/dev/dri/renderD129')
        self.assertEqual(command[command.index('-hwaccel_output_format') + 1], 'vaapi')
        filters = command[command.index('-vf') + 1]
        self.assertLess(filters.index('fps='), filters.index('hwdownload'))
        self.assertLess(filters.index('hwdownload'), filters.index('scale='))
        self.assertIn('format=nv12,', filters)
        self.assertEqual(len(frames), 1)

    def test_partial_hardware_failure_restarts_both_passes(self):
        detector = Detector.__new__(Detector)
        attempts = []
        def analyze(request):
            attempts.append(dict(request))
            if request['intel_acceleration']:
                raise HardwareDecodeError('failed after some ROI frames')
            return {'decode_backend': 'software', 'regions': [{'start': 0, 'end': 4}], 'frames': 8}
        detector._analyze = analyze
        request = {'intel_acceleration': True, 'path': 'clip.mp4'}
        result = detector.analyze(request)
        self.assertEqual([a['intel_acceleration'] for a in attempts], [True, False])
        self.assertEqual(result['decode_backend'], 'software')
        self.assertIn('failed after', result['acceleration_fallback'])
        self.assertTrue(request['intel_acceleration'])  # don't mutate later jobs

    def test_software_failure_after_hardware_failure_is_not_a_completed_result(self):
        detector = Detector.__new__(Detector)
        def analyze(request):
            if request['intel_acceleration']:
                raise HardwareDecodeError('no GPU')
            raise ValueError('corrupt video')
        detector._analyze = analyze
        with self.assertRaisesRegex(ValueError, 'corrupt video'):
            detector.analyze({'intel_acceleration': True})

    def test_disabled_acceleration_does_not_retry_software_failure(self):
        detector = Detector.__new__(Detector)
        with patch.object(detector, '_analyze', side_effect=ValueError('corrupt video')) as analyze:
            with self.assertRaisesRegex(ValueError, 'corrupt video'):
                detector.analyze({'intel_acceleration': False})
        self.assertEqual(analyze.call_count, 1)


@unittest.skipUnless(shutil.which('ffmpeg') and shutil.which('ffprobe'), 'FFmpeg required')
class IntelFallbackIntegrationTests(unittest.TestCase):
    def test_missing_render_device_uses_software_for_complete_analysis(self):
        detector = Detector.__new__(Detector)
        detector.width = detector.height = 64
        detector.model_name = 'yolo26m'
        detector.model_sha256 = 'a' * 64
        detector.predict = lambda raw, cutoff: []
        with tempfile.TemporaryDirectory() as folder:
            path = str(Path(folder) / 'clip.mp4')
            subprocess.run(['ffmpeg', '-nostdin', '-v', 'error', '-f', 'lavfi', '-i',
                            'color=c=black:s=96x64:r=10', '-t', '2', '-c:v', 'libx264',
                            '-threads', '1', path], check=True, capture_output=True, timeout=30)
            result = detector.analyze({'path': path, 'ffmpeg': shutil.which('ffmpeg'),
                'ffprobe': shutil.which('ffprobe'), 'intel_acceleration': True,
                'vaapi_device': str(Path(folder) / 'missing-render-device'), 'fps': 2,
                'confidence': .65, 'min_frames': 3, 'merge_gap': .5, 'car_movement': .025, 'policy': 'b' * 64})
        self.assertTrue(result['complete'])
        self.assertEqual(result['decode_backend'], 'software')
        self.assertTrue(result['acceleration_fallback'])
        self.assertEqual(result['motion_frames'], 4)
        self.assertEqual(result['frames'], 4)
