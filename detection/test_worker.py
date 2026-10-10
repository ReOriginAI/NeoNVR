import shutil
import subprocess
import tempfile
import unittest
from pathlib import Path

from worker import Detector, LABELS, Timeline, kind


def observation(class_id, x=.2, y=.2, score=.85):
    return (class_id, [x, y, x + .12, y + .15], score)


class TimelineTests(unittest.TestCase):
    def timeline(self):
        return Timeline(fps=2, min_frames=2, merge_gap=.5, car_movement=.015)

    def test_target_categories(self):
        for class_id in LABELS:
            with self.subTest(class_id=class_id):
                timeline = self.timeline()
                for frame in range(4):
                    timeline.observe(frame / 2, [observation(class_id, x=.2 + frame * .02)])
                intervals, uncertain = timeline.finish(2)
                self.assertFalse(uncertain)
                self.assertEqual([row['kind'] for row in intervals], [kind(class_id)])
                self.assertIn(LABELS[class_id], intervals[0]['labels'])

    def test_parked_car_and_detector_jitter_do_not_qualify(self):
        timeline = self.timeline()
        for frame in range(20):
            timeline.observe(frame / 2, [observation(2, x=.2 + (frame % 3) * .002)])
        self.assertEqual(timeline.finish(10), ([], False))

    def test_moving_car_with_two_observations(self):
        timeline = self.timeline()
        timeline.observe(0, [observation(2)])
        timeline.observe(.5, [observation(2, x=.23)])
        self.assertEqual(timeline.finish(1), ([], True))

    def test_moving_car_with_consistent_translation(self):
        timeline = self.timeline()
        for frame in range(3):
            timeline.observe(frame / 2, [observation(2, x=.2 + frame * .02)])
        intervals, uncertain = timeline.finish(1.5)
        self.assertFalse(uncertain)
        self.assertEqual(intervals[0]['kind'], 'moving_car')

    def test_single_parked_car_box_jump_is_not_motion(self):
        timeline = self.timeline()
        for frame in range(8):
            timeline.observe(frame / 2, [observation(2, x=.23 if frame >= 4 else .2)])
        self.assertEqual(timeline.finish(4), ([], True))

    def test_parked_car_box_spike_is_not_motion(self):
        timeline = self.timeline()
        for frame in range(8):
            timeline.observe(frame / 2, [observation(2, x=.23 if frame == 4 else .2)])
        self.assertEqual(timeline.finish(4), ([], True))

    def test_parked_car_box_resizing_is_not_motion(self):
        timeline = self.timeline()
        for frame in range(5):
            # One edge stays fixed while the detector expands the box.
            box = [.2, .2, .32 + frame * .02, .35]
            timeline.observe(frame / 2, [(2, box, .85)])
        self.assertEqual(timeline.finish(2.5), ([], True))

    def test_low_confidence_false_positives_do_not_confirm_or_prune(self):
        for class_id in (0, 1, 2, 16):
            with self.subTest(class_id=class_id):
                timeline = self.timeline()
                for frame in range(6):
                    timeline.observe(frame / 2, [observation(class_id, x=.2 + frame * .02, score=.5)])
                self.assertEqual(timeline.finish(3), ([], True))

    def test_one_confident_frame_cannot_promote_a_weak_track(self):
        timeline = Timeline(2, 3, .5, .025, confidence=.65)
        for frame in range(6):
            timeline.observe(frame / 2, [observation(0, score=.9 if frame == 3 else .5)])
        self.assertEqual(timeline.finish(3), ([], True))

    def test_two_confident_person_frames_remain_unconfirmed_by_default(self):
        timeline = Timeline(2, 3, .5, .025)
        for frame in range(2):
            timeline.observe(frame / 2, [observation(0, score=.7628)])
        self.assertEqual(timeline.finish(1), ([], True))

    def test_low_confidence_sample_breaks_confirmation(self):
        timeline = Timeline(2, 3, .5, .025)
        for frame in range(5):
            timeline.observe(frame / 2, [observation(0, score=.5 if frame == 2 else .85)])
        self.assertEqual(timeline.finish(2.5), ([], True))

    def test_merge_gap_does_not_confirm_intermittent_observations(self):
        timeline = Timeline(2, 3, 10, .025)
        for frame in range(6):
            timeline.observe(frame / 2, [observation(0)] if frame % 2 == 0 else [])
        self.assertEqual(timeline.finish(3), ([], True))

    def test_confident_presence_does_not_require_motion(self):
        for class_id in (0, 1, 16):
            with self.subTest(class_id=class_id):
                timeline = Timeline(2, 3, .5, .025)
                for frame in range(4):
                    timeline.observe(frame / 2, [observation(class_id)])
                intervals, uncertain = timeline.finish(2)
                self.assertFalse(uncertain)
                self.assertEqual([row['kind'] for row in intervals], [kind(class_id)])

    def test_confidence_setting_can_be_lowered_explicitly(self):
        timeline = Timeline(2, 3, .5, .025, confidence=.45)
        for frame in range(3):
            timeline.observe(frame / 2, [observation(0, score=.5)])
        intervals, uncertain = timeline.finish(1.5)
        self.assertFalse(uncertain)
        self.assertEqual([row['kind'] for row in intervals], ['person'])

    def test_large_confirmation_count_still_allows_local_vehicle_motion(self):
        timeline = Timeline(2, 10, .5, .025)
        for frame in range(12):
            timeline.observe(frame / 2, [observation(2, x=.2 + frame * .01)])
        intervals, uncertain = timeline.finish(6)
        self.assertFalse(uncertain)
        self.assertEqual([row['kind'] for row in intervals], ['moving_car'])

    def test_vehicle_motion_at_low_sample_rate_remains_ambiguous(self):
        timeline = Timeline(.5, 3, .5, .025)
        for frame in range(4):
            timeline.observe(frame * 2, [observation(2, x=.2 + frame * .04)])
        self.assertEqual(timeline.finish(8), ([], True))

    def test_movement_does_not_mark_entire_parked_presence(self):
        timeline = self.timeline()
        for frame in range(30):
            x = .2 if frame < 10 else .2 + min(3, frame - 10) * .03
            timeline.observe(frame / 2, [observation(2, x=x)])
        intervals, _ = timeline.finish(15)
        self.assertGreater(intervals[0]['start'], 0)
        self.assertLess(intervals[-1]['end'], 10)

    def test_brief_observations_prevent_pruning(self):
        for class_id in (0, 1, 2, 16):
            timeline = self.timeline()
            timeline.observe(0, [observation(class_id)])
            self.assertEqual(timeline.finish(.5), ([], True))
        timeline = self.timeline()
        timeline.observe(0, [observation(2)])
        timeline.observe(.5, [observation(2)])
        self.assertEqual(timeline.finish(1), ([], True))

    def test_multiple_timeframes_and_overlapping_categories(self):
        timeline = self.timeline()
        for frame in range(14):
            boxes = []
            if frame in (0, 1, 10, 11):
                boxes.append(observation(0))
            if frame in (0, 1):
                boxes.append(observation(16, x=.6))
            timeline.observe(frame / 2, boxes)
        intervals, uncertain = timeline.finish(7)
        self.assertFalse(uncertain)
        self.assertEqual(sum(row['kind'] == 'person' for row in intervals), 2)
        self.assertEqual(sum(row['kind'] == 'animal' for row in intervals), 1)
        self.assertTrue(all(0 <= row['start'] < row['end'] <= 7 for row in intervals))

    def test_unrelated_cars_do_not_become_a_moving_track(self):
        timeline = self.timeline()
        timeline.observe(0, [observation(2, x=.1)])
        timeline.observe(.5, [observation(2, x=.8)])
        self.assertEqual(timeline.finish(1), ([], True))

    def test_empty_frames_have_no_targets(self):
        timeline = self.timeline()
        for frame in range(10):
            timeline.observe(frame / 2, [])
        self.assertEqual(timeline.finish(5), ([], False))


@unittest.skipUnless(shutil.which('ffmpeg') and shutil.which('ffprobe'), 'FFmpeg and FFprobe required')
class AnalysisTests(unittest.TestCase):
    def test_raising_confirmation_confidence_keeps_weak_candidates(self):
        # Exercise real decoding with controlled model observations. Filtering
        # weak candidates before tracking would wrongly permit clip pruning.
        detector = Detector.__new__(Detector)
        detector.width = detector.height = 64
        detector.model_sha256 = 'a' * 64
        detector.predict = lambda raw, cutoff: [observation(0, score=.5)] if cutoff <= .5 else []
        with tempfile.TemporaryDirectory() as folder:
            path = str(Path(folder) / 'recording.mp4')
            subprocess.run(['ffmpeg', '-nostdin', '-v', 'error', '-f', 'lavfi', '-i',
                            'color=c=black:s=96x64:r=10', '-t', '2', '-c:v', 'libx264',
                            '-threads', '1', path], check=True, capture_output=True, timeout=30)
            result = detector.analyze({'path': path, 'ffmpeg': shutil.which('ffmpeg'),
                                       'ffprobe': shutil.which('ffprobe'), 'fps': 2,
                                       'confidence': .65, 'min_frames': 3, 'merge_gap': .5,
                                       'car_movement': .025, 'policy': 'b' * 64})
        self.assertTrue(result['complete'])
        self.assertEqual(result['frames'], 4)
        self.assertEqual(result['intervals'], [])
        self.assertTrue(result['uncertain'])


if __name__ == '__main__':
    unittest.main()
