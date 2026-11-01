# Detection optimization measurements — 2026-11-01

Measured inside the running `neonvr` nerdctl container, with the service still
processing its normal workload. These are initial measurements on one recording,
not a controlled throughput or detection-accuracy benchmark.

> The latest worker now shares one model-size decode stream between MOG2 and YOLO.
> Measurements below describe the preceding separate-decoder implementation;
> they do not measure the new shared-stream pipeline.

## Retained change: sample before GPU-to-CPU download

Previously the VA-API filter order downloaded every decoded frame, then sampled
and resized it. Sampling now happens on hardware frame references first, so
only the needed frames are downloaded. CPU resizing, letterboxing, detection
thresholds, tracking, and hardware-failure fallback remain the same.

Thirty-second sample, 320×180 output at 5 FPS:

| Trial | Original VA-API path | Sample before download |
| --- | ---: | ---: |
| Initial 1 | 4.593 s | 2.547 s |
| Initial 2 | 4.971 s | 2.461 s |
| Frame verification 1 | 5.610 s | 3.025 s |
| Frame verification 2 | 5.761 s | 2.885 s |

The verification trials produced all 150 frames and identical SHA-256 hashes
for software decoding, the original VA-API path, and the new path:
`61ba36843bfd4faf89660857118b7df2fd2af86f9316710994254bf8af67d741`.
An additional comparison checked a sought 10–14 second ROI at 640×640 and 4 FPS;
all 16 output frames and their relative timestamps matched exactly.
This is approximately a 46–50% reduction in this decode/resize stage's wall time,
not a measured doubling of total detection throughput.

## Alternatives tested and not adopted

- GPU scaling (`scale_vaapi`) failed with unsupported processing-profile errors
  for both the default driver and an attempted i965 driver selection.
- OpenVINO 2026.4.1 exposed only a CPU device. Compiling for GPU failed.
- On identical real input frames, ONNX Runtime measured 369/510 ms per frame
  across two trials, versus OpenVINO CPU FP32 at 622/650 ms per frame. OpenVINO
  also produced different raw output rows; detection parity was not established.
  ONNX Runtime remains the inference backend.

The OpenVINO experiment used an isolated `/tmp/neonvr-openvino-benchmark` virtual
environment inside the running container. It did not replace runtime packages.

## Reproduction

`scripts/benchmark_detection.py` provides the regular read-only stage benchmark.
Its `--optimize --decode-only` mode compares decoding filter orders and hashes.
The complete `--optimize` mode additionally requires OpenVINO installed in the
Python environment used for the benchmark. Existing service contention and
cache state affect timing; repeat on busy and quiet recordings before making
capacity estimates.

## Shared-stream verification

The subsequent shared-stream implementation was checked using the real YOLO26m
model inside the container. A 20.012-second segment at 640×640 and 4 FPS used one
VA-API decode: 80 MOG2 samples, 16 ROI-selected YOLO frames, and a peak pre-roll
queue of 10 frames. Analysis took 10.282 seconds under live service contention.
This is not a matched before/after throughput comparison.

A 59.968-second clip triggered a VA-API surface-sync error. The whole analysis
restarted in software and completed, with 240 shared samples and 62 selected
YOLO frames. Hardware-error recovery is an exception to the single-decode path;
there is no separate normal decode for MOG2 and YOLO.

The rebuilt running service was verified to use the shared-stream worker and
report synchronized 4 FPS settings for both stages.

## Matched before/after comparison of the shared stream

Both implementations were run inside the live container on exactly the same
stream-copied footage, with the same YOLO26m weights, six inference threads,
4 FPS YOLO sampling, thresholds, padding, and VA-API configuration. Each segment
used the order **before, after, after, before** (two measured runs per variant).
The before source includes the earlier sampling-before-download optimization.
The service was left running throughout; system load varied during the test.

The intentional pipeline differences were MOG2 resolution/rate: before used
320×180 at 5 FPS; after used the shared 640×640 stream at 4 FPS. Therefore ROI
coverage can differ, and this is an end-to-end comparison of the requested
configurations, not an isolated measure of decoding savings.

| Segment | Duration | Before mean | After mean | Reduction |
| --- | ---: | ---: | ---: | ---: |
| Quiet footage | 20.012 s | 11.505 s | 11.409 s | 0.8% |
| Vehicle activity | 16.221 s | 34.930 s | 32.409 s | 7.2% |

Quiet footage selected the same 4-second ROI and 16 YOLO frames in both
implementations. The timing difference is small relative to run-to-run variation.
The before activity run selected 15.2 seconds / 61 YOLO frames; the compared
after run selected 14.0 seconds / 56 frames. Both found the same car event at
6.625–10.250 seconds with confidence 0.9446. Fewer selected inference frames
account for a substantial part of the activity-segment improvement.

In the measured activity runs, YOLO consumed approximately 95% of processing
time. Shared decoding removed the normal second decoder, but the larger MOG2
input and dominant inference work limit end-to-end gains. Activity-heavy footage
still processed slower than real time in this test.

Source SHA-256 hashes used in the comparison:

- Before: `471800bf8d985cd8de1c113634ad5540998335e65ad49138cced080df2252a0a`
- After: `7688feac56087c8addd6f2645a5dbc68429b07db760be2adc58254b4a4ebfe9e`
