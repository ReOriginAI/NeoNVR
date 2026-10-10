import unittest

from worker import LABELS, Timeline, kind


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
        intervals, uncertain = timeline.finish(1)
        self.assertFalse(uncertain)
        self.assertEqual(intervals[0]['kind'], 'moving_car')

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


if __name__ == '__main__':
    unittest.main()
