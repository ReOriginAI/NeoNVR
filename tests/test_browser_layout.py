"""Browser checks with mocked APIs; no recording directories or running server needed.
Run: python -m unittest discover -s tests -p 'test_browser_*.py'
Requires playwright and its Chromium browser, plus FFmpeg.
"""
import json
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest
from urllib.parse import urlsplit

try:
    from playwright.sync_api import sync_playwright
except ImportError:
    sync_playwright = None

ROOT = Path(__file__).resolve().parents[1]


@unittest.skipUnless(sync_playwright and shutil.which('ffmpeg'), 'Playwright and FFmpeg required')
class PlayerLayoutTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.folder = tempfile.TemporaryDirectory()
        video = Path(cls.folder.name) / 'clip.mp4'
        subprocess.run(['ffmpeg', '-nostdin', '-v', 'error', '-f', 'lavfi', '-i',
                        'color=c=green:s=160x96:r=10', '-t', '4', '-c:v', 'libx264',
                        '-threads', '1', '-movflags', '+faststart', str(video)], check=True)
        cls.video = video.read_bytes()
        source = (ROOT / 'main.go').read_text()
        cls.html = '<!doctype html>' + source.split('Parse(`<!doctype html>', 1)[1].split('</html>`))', 1)[0] + '</html>'
        cls.html = cls.html.replace('{{.Root}}', 'Test recordings')
        cls.clips = [dict(camera='Garden', file=f'clip-{i}.mp4', favorite=False, quarantined=False,
                          start_seconds=28800 + i * 180 + 17, unix_seconds=28800 + i * 180 + 17,
                          date='2026-11-01', time=f'Clip {i}', timestamp=f'2026-11-01 clip {i}',
                          video_url=f'/video?clip={i}', thumb_url='/thumb') for i in range(40)]
        cls.playwright = sync_playwright().start()
        cls.browser = cls.playwright.chromium.launch()

    @classmethod
    def tearDownClass(cls):
        cls.browser.close()
        cls.playwright.stop()
        cls.folder.cleanup()

    def route(self, route):
        path = urlsplit(route.request.url).path
        if path == '/':
            route.fulfill(body=self.html, content_type='text/html')
        elif path.startswith('/static/'):
            route.fulfill(body=(ROOT / path.lstrip('/')).read_bytes(), content_type='text/javascript')
        elif path == '/video':
            request_range = route.request.headers.get('range')
            if request_range:
                begin, end = request_range.removeprefix('bytes=').split('-', 1)
                begin, end = int(begin or 0), int(end) if end else len(self.video) - 1
                end = min(end, len(self.video) - 1)
                route.fulfill(status=206, body=self.video[begin:end+1], content_type='video/mp4',
                              headers={'Accept-Ranges': 'bytes', 'Content-Range': f'bytes {begin}-{end}/{len(self.video)}'})
            else:
                route.fulfill(body=self.video, content_type='video/mp4', headers={'Accept-Ranges': 'bytes'})
        elif path == '/thumb':
            route.fulfill(body='<svg xmlns="http://www.w3.org/2000/svg" width="160" height="90"/>', content_type='image/svg+xml')
        else:
            responses = {
                '/api/cameras': {'cameras': [{'name': 'Garden'}]},
                '/api/dates': {'dates': ['2026-11-01']},
                '/api/clips': {'clips': self.clips},
                '/api/stats': {'root_size_human': '1 MB', 'generated_at': '2026-11-01T12:00:00Z'},
                '/api/detections': {'status': 'complete', 'detections': {
                    'duration': 4, 'regions': [{'start': 0, 'end': 4}],
                    'motion_frames': 20, 'frames': 8, 'intervals': [{'kind': 'person', 'start': 1, 'end': 3, 'confidence': .9}]}},
                '/api/camera-motion': {'motion': {'duration': 4, 'creation_time': '2026-11-01T12:00:00Z',
                    'data': [], 'invalid_rows': 0, 'outside_clip_rows': 0,
                    'samples': [{'time': 1.5, 'bin_img_ratio': 25.0, 'moving_level': 3, 'scene': 1, 'move': 1, 'moving_th': 2, 'unix_ms': 0}]}},
            }
            route.fulfill(body=json.dumps(responses.get(path, {})), content_type='application/json')

    def exercise_layout(self, mobile):
        context = self.browser.new_context(viewport={'width': 390 if mobile else 1280, 'height': 844 if mobile else 900},
                                           is_mobile=mobile, has_touch=mobile)
        page = context.new_page()
        page.set_default_timeout(5000)
        errors = []
        page.on('pageerror', lambda error: errors.append(str(error)))
        page.route('**/*', self.route)
        page.goto('http://neonvr.test/', wait_until='domcontentloaded')
        page.wait_for_function('document.querySelectorAll(".sideThumb").length === 40 && document.getElementById("mainVideo").readyState >= 1')
        seek = page.locator('media-time-range')
        self.assertEqual(seek.locator('.lane').count(), 1)
        self.assertEqual(seek.locator('.roi-bracket').count(), 2)
        before = page.locator('#mainVideo').bounding_box()
        shell = page.locator('#mainShell').bounding_box()
        page.evaluate('mainVideoEl.currentTime = 1.5')
        page.wait_for_function('document.getElementById("cameraMotionStatus").textContent.startsWith("Camera sample at")')
        after = page.locator('#mainVideo').bounding_box()
        self.assertAlmostEqual(before['width'], after['width'], delta=.5)
        self.assertAlmostEqual(before['height'], after['height'], delta=.5)
        page.evaluate('document.getElementById("cameraMotionStatus").textContent = "Camera sample at " + "long details ".repeat(80)')
        self.assertAlmostEqual(shell['height'], page.locator('#mainShell').bounding_box()['height'], delta=.5)
        # Selecting an exact video moves the point to that five-minute bin.
        page.locator('.sideThumb[data-key="Garden/clip-30.mp4"]').click()
        self.assertEqual(page.locator('#timeCenter').input_value(), '34200')
        # Point selection chooses a representative nearby clip without snapping.
        page.evaluate('timeCenterEl.value=29820; timeCenterEl.dispatchEvent(new Event("input"))')
        self.assertEqual(page.locator('#timeCenter').input_value(), '29700')
        self.assertIn(page.locator('.sideThumb.active').get_attribute('data-key'), ['Garden/clip-5.mp4', 'Garden/clip-6.mp4'])
        page.locator('#nextVideo').scroll_into_view_if_needed()
        page.wait_for_timeout(100)
        scroll_y = page.evaluate('window.scrollY')
        for _ in range(4):
            page.locator('#nextVideo').click()
        self.assertAlmostEqual(scroll_y, page.evaluate('window.scrollY'), delta=1)
        # An enabled filter follows sidebar choices, while dragging keeps the
        # requested minute even if the representative clip starts nearby.
        page.evaluate('timeFilterEnabled=true; timeCenterEl.value=29820; timeCenterEl.dispatchEvent(new Event("input"))')
        self.assertEqual(page.locator('#timeCenter').input_value(), '29700')
        self.assertIn(page.locator('.sideThumb.active').get_attribute('data-key'), ['Garden/clip-5.mp4', 'Garden/clip-6.mp4'])
        self.assertLess(page.locator('.sideThumb').count(), 40)
        page.locator('.sideThumb[data-key="Garden/clip-4.mp4"]').click()
        self.assertEqual(page.locator('#timeCenter').input_value(), '29400')
        self.assertEqual(page.locator('.sideThumb.active').get_attribute('data-key'), 'Garden/clip-4.mp4')
        if not mobile:
            page.locator('media-fullscreen-button').click()
            page.wait_for_function('document.fullscreenElement !== null')
            self.assertLessEqual(page.locator('#mainVideo').bounding_box()['height'], page.viewport_size['height'])
            self.assertLessEqual(page.locator('.detectionStatuses').bounding_box()['y'] + page.locator('.detectionStatuses').bounding_box()['height'], page.viewport_size['height'] + 1)
            page.evaluate('document.exitFullscreen()')
        self.assertEqual(errors, [])
        context.close()

    def test_desktop(self):
        self.exercise_layout(False)

    def test_mobile(self):
        self.exercise_layout(True)
