# NeoNVR

Docker-first security-camera MP4 browser with shared-stream MOG2 + CPU-only YOLO26m event detection, including Blink clips stored as:

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
https://SERVER_IP:5443/streams
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

Rebuild after changes:

```bash
docker compose up -d --build
```

The compose file mounts `/mnt/mySSD/BlinkClips` read/write into the container so thumbnails, favorites metadata, Delete, and optional auto-delete work against the real host folder.

## Features

- Default listen address: `0.0.0.0:5500`
- Configurable clips root with `-root`
- Native date/calendar selector
- Camera multi-select dropdown with an `All` toggle
- Dual time-range slider with hour labels and half-hour `½` labels
- Colored semi-transparent timeline markers per camera, based on clip start time only
- Newest-first default sorting, with oldest-first option
- Main player at the top with a scrollable one-column thumbnail rail on the right
- Full mini-video/card grid remains below the controls for browsing/actions
- Clicking a sidebar thumbnail or mini-video thumbnail loads it into the main player
- Active sidebar thumbnail and mini-video card stay highlighted/synced with the main player
- Sidebar captures wheel/touchpad scrolling so the page does not jump while browsing the rail
- `Ctrl` + arrow keys loads previous/next visible video
- Space bar pauses/resumes the main player
- `←` / `→` steps within the current video by the configured Step sec value
- Zoom mode toggle on the main player: when Zoom is on, pinch/trackpad zoom and drag-to-pan are enabled; when Zoom is off, tapping/clicking the video pauses/plays normally
- `Fave`/`Unfave` preserves clips from auto-delete
- Manual `Delete`
- All tab combining active and quarantined clips for the selected date, with purple outlines for quarantine
- Right-click clip actions and mobile ⋮ menus for download, delete, quarantine, and recovery
- Explicit Select mode for bulk deletion; Shift-click or touch-and-hold selects every visible clip in between
- Favorites tab
- Optional built-in auto-delete of non-favorite clips older than X days
- Two-pass MOG2 motion windows followed by YOLO26m detection of humans, COCO animals, bicycles, motorcycles, and moving cars/buses/trucks
- Clip-relative detection intervals embedded in the MP4, with colored interval brackets on a single timeline below the video scrubber and a live detection indicator
- Optional pruning of fully analyzed no-target clips, using recoverable quarantine by default
- Displays total size of the BlinkClips folder only, avoiding misleading parent-folder stats inside Docker

## Docker image

The Go binary is built on Alpine. A separate Python build stage downloads the official `yolo26m.pt` weights and exports a fixed 640-pixel, end-to-end ONNX model. The Debian Bookworm runtime includes FFmpeg, Python, NumPy, headless OpenCV, and **CPU-only ONNX Runtime**. PyTorch, Ultralytics and GPU libraries stay in the model build stage. The model is included in the image, so analysis and the browser work offline after building. The first build needs network access for packages and model weights.

## Docker configuration

The included Compose service starts with:

```yaml
command: ["-config", "/config/neonvr.json"]
volumes:
  - ./config:/config
```

Runtime options live in `config/neonvr.json`, rather than in the Compose command. The included recording mounts are unchanged. Retention is configured by `auto_delete_days` (14 in the supplied file; 0 disables age-based deletion) and `auto_delete_every` (`6h`). It skips favorites in every recording source. Separately, `storage_min_free_percent` defaults to 10 (0 disables space-pressure deletion). Every minute, sources below this free disk percentage are cleaned oldest first: quarantined videos, then fully analyzed vehicle-only videos (moving cars/trucks/buses or bikes). Favorites and active clips with confirmed human or animal detections are preserved by pressure cleanup; age-based retention still applies independently. Cleanup stops when the reserve is restored, or when no eligible footage remains. Independently, above 50% disk use, complete analyses with only unconfirmed observations (no confirmed intervals) move to quarantine, unless favorited. This is checked every minute even if space-pressure deletion is disabled; it also applies to existing analyses. Once quarantined, these files enter the first deletion tier if free space falls below the reserve.

After changing source code or adopting the config mount, rebuild and start the container once:

```bash
cd /home/ubuntu/NeoNVR
sudo nerdctl compose up -d --build
```

Subsequent settings changes use **Save & restart** in the browser and do not require a container rebuild.

## Native run without Docker

```bash
sudo apt update
sudo apt install -y golang-go ffmpeg

go build -o neonvr .
./neonvr -root /mnt/mySSD/BlinkClips
```

Useful flags:

```bash
./neonvr \
  -root /mnt/mySSD/BlinkClips \
  -addr 0.0.0.0:5500 \
  -thumb-workers 2 \
  -thumb-width 480 \
  -thumb-seek 0.2 \
  -size-refresh 30m
```

## Live streams and speaker talkback

Open the **Streams** tab to view live RTSP cameras. Live viewing is **off after every app start**. Loading the recordings page or Streams tab does not connect to any camera. Select a camera and choose **Watch**, or double-click a camera in the sidebar to connect directly. The streaming gateway starts automatically when needed. A single click selects the camera without starting a stream. **Stop**, leaving the page, or hiding the browser tab closes that viewer's connection and microphone.

The supplied config includes `rtsp://192.168.1.150:554/ch0_0.h264` for viewing and `rtsp://192.168.1.150:554/ch0_0.h264?backchannel=1` for speaker talkback. Edit sources in **Settings → Live streams**, one camera per line: `name | viewing URL | optional backchannel URL`. Save and restart to apply source changes; live viewing resets to off.

**Talk to camera** explicitly requests microphone permission and sends your audio to the camera speaker using the backchannel URL. **Stop talking** stops microphone capture and returns to the ordinary viewing URL. Microphone access requires HTTPS or localhost. Camera backchannel and browser audio codecs must be compatible; unsupported backchannels produce an error and stop microphone capture.

### Built-in HTTPS for LAN microphone access

The supplied configuration enables **HTTPS on port 5443**, directly in the NeoNVR container. No reverse proxy, extra container, external certificate service, or internet access is needed. Rebuild with `sudo nerdctl compose up -d --build`, then open `https://192.168.1.130:5443/streams` (or use **Open HTTPS for microphone access** on the HTTP Streams page). Accept the local certificate warning and allow microphone access when choosing **Talk to camera**. HTTP remains available on port 5500 for recordings and setup.

NeoNVR creates its own self-signed local certificate authority and a server certificate at first HTTPS startup. The certificate identities, including private keys, persist in `config/tls/`, inside the existing `/config` mount. Files use mode 0600 and are excluded from Git. They are generated at runtime, never included in the image. Rebuilds and app restarts reuse them. Server certificates renew on startup when nearing expiry or when configured hostnames change; the local CA stays the same, so devices do not need to trust it again.

Some browsers/devices require explicit certificate trust before enabling the microphone. Download `http://SERVER_IP:5500/certificate` (also linked under the Streams player), install **neonvr-ca.crt** as a trusted certificate authority on that device, and reopen the HTTPS page. On iOS/iPadOS, installing the certificate profile also requires [enabling full trust](https://support.apple.com/102390) under **Settings → General → About → Certificate Trust Settings**. This download contains only the public CA certificate.

In **Settings → Server and storage**, set **HTTPS server names and IP addresses** to the LAN IP or hostname used by your browsers, one per line, and choose **Save & restart**. The supplied configuration includes this server's `192.168.1.130`; localhost and the server's network addresses are also included automatically. Container network addresses do not identify the host's LAN address, so update this setting if the host's LAN IP changes. An empty **HTTPS listen address** disables built-in HTTPS; native installs can enable it with `-https-addr 0.0.0.0:5443`. Changing the HTTPS port also requires updating the Compose mapping.

The Docker image includes [go2rtc v1.9.14](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14). Its private API listens only on loopback and is accessed through NeoNVR. It starts on demand with no preloaded camera connections. Viewing uses WebRTC with no video transcoding. Unsupported camera audio, including AAC from Yi cameras, is converted to Opus only while needed for viewing. FFmpeg reads the shared camera stream through a private loopback RTSP port, avoiding another camera connection. Camera audio starts muted; use **Unmute camera audio** or the video's volume control to listen. The button shows **No camera audio** if no playable audio track is received. Listening is separate from microphone talkback. The browser UI and LAN streaming work without a CDN or public STUN service.

The Compose file maps **8555/TCP and 8555/UDP** for WebRTC media in addition to the app's HTTP and HTTPS ports. The browser must be able to reach that media port on the server. In Settings, a blank **WebRTC server hostname or IP** uses the hostname from the URL used to enable viewing; set it explicitly if a reverse proxy uses a different address from the media server. Changing the media port also requires changing both Compose mappings. Native installations must install the [go2rtc binary](https://github.com/AlexxIT/go2rtc/releases/tag/v1.9.14) and set its executable path in Settings if it is not on `PATH`.

## Main player controls

- Click a thumbnail or `Load` to send it to the main player.
- Space pauses/resumes.
- `←` and `→` step within the loaded video.
- `Ctrl` + `←` / `↑` loads the previous visible video.
- `Ctrl` + `→` / `↓` loads the next visible video.
- Change `Step sec` for fine-grain frame/time stepping.
- Use `Zoom off` / `Zoom on` to enable inspection zoom. With Zoom off, tapping/clicking the video pauses/plays. With Zoom on, pinch/trackpad zooms the actual video display and drag pans while zoomed. `Reset zoom` restores the view. This is inspection zoom only; it does not create missing detail.

## Bulk delete

Click `Select` to enter delete-selection mode. In that mode, clicking thumbnails selects them instead of loading them. Shift-click or tap and hold another thumbnail to select every visible clip in between, then use `Delete selected`.

## Files written inside BlinkClips

By default the app writes:

```text
/mnt/mySSD/BlinkClips/.blinkclips_thumbs/
/mnt/mySSD/BlinkClips/.blinkclips_meta.json
```

The thumbnail cache can be moved with `-thumb-dir`.

Detection also writes `.neonvr_detections/` under the primary root as an expendable JSON cache. Detection metadata lives inside the MP4 itself, so deleting the cache does not lose event intervals. Quarantined recordings live in `<camera-directory>/.neonvr_pruned/`, preserving Yi hour subdirectories where applicable. The **All** tab combines the main library and quarantine for the selected date, with purple outlines and labels on quarantined cards and thumbnails. The **Quarantine** tab lists quarantined clips across all dates with camera and time filters, playback, Previous/Next navigation, favorites, and single or bulk deletion. Favorited quarantine clips also appear in Favorites and are protected from space-pressure deletion.

## Two-pass MOG2 / YOLO26m detection and false-positive pruning

Detection uses one shared decode stream with two analysis stages:

1. FFmpeg decodes and samples the original clip **once**, producing letterboxed frames at the model input dimensions (**640×640** for the bundled models). **Shared MOG2 / YOLO frames per second** controls both stages; the separate legacy `detection_motion_fps` field is normalized to `detection_fps`.
2. MOG2 examines every sampled frame, excludes shadows and small speckles, and selects temporal windows padded by **2 seconds before and 3 seconds after**. The opening second is always included, with padding, to check subjects present at startup. Overlapping windows merge.
3. A bounded rolling buffer retains enough of these exact frames to cover pre-roll. Once future motion can no longer change a frame's eligibility, YOLO either analyzes that same frame or it is discarded. Selected frames are never decoded or inferred twice; post-roll is consumed as the stream advances. Tracking restarts between distinct windows. There is no full-clip frame cache or temporary video file.

At 640×640, 4 FPS and 2 seconds of pre-roll, the raw frame queue holds at most about **12 MB** (up to 10 frames), plus MOG2/model/runtime memory. Its size depends on pre-roll and FPS, not clip duration. MOG2 itself does more pixel work at 640×640 than at the previous 320×180 size, so removing duplicate decoding does not guarantee an overall speedup for every scene.

**Settings → Object detection** exposes the shared sampling FPS, MOG2 variance threshold, minimum connected-region area, pre-roll and post-roll. Changing the analysis policy triggers reanalysis of active recordings after restart. The bundled default is `yolo26m.onnx`; the previous small and nano exports remain available. Any hardware decode failure discards partial state and restarts the complete stream in software. Decoder coverage and inference must both succeed before results can authorize pruning.

References: [OpenCV MOG2](https://docs.opencv.org/4.x/d7/d7b/classcv_1_1BackgroundSubtractorMOG2.html), [Ultralytics YOLO26](https://docs.ultralytics.com/models/yolo26).

Motion gating can miss very slow, tiny or stationary subjects outside the opening window; lowering the minimum region area or variance threshold increases motion sensitivity. Camera movement and lighting changes can select large portions of a clip, so CPU savings depend on the scene. A motion-pass or ROI-decode failure keeps the recording and does not publish a complete analysis.

References: [OpenCV MOG2](https://docs.opencv.org/4.x/d7/d7b/classcv_1_1BackgroundSubtractorMOG2.html), [Ultralytics YOLO26](https://docs.ultralytics.com/models/yolo26).

The daily activity sparkline replaces the two-ended range slider. Double-click it to enable or disable the time filter, or use the **20-minute filter** button on mobile. With filtering enabled, drag the single point to include videos starting within 10 minutes before or after that time; the window stops at midnight. Filtering starts off, showing the whole selected date. The blue line counts clips per five-minute bin, and purple shows quarantined activity. All and Clips use the date picker and date pills; All includes dates containing only quarantine. Selecting a sidebar video moves the point to its five-minute bin. Moving the point selects a nearby clip and keeps it visible in the sidebar; this is approximate browsing, not a second-accurate seek. Dragging within the same five-minute bin does not reload the video or rebuild the list.

The video rail includes Previous/Next buttons in the current sort/filter order. These buttons scroll only the thumbnail list, keeping the page steady on mobile. Shift+Left/Right also selects the previous/next video, including when a navigation button has focus; unmodified arrows still seek. On touch devices, hold a thumbnail for 550 ms to enter selection mode; tap another thumbnail to toggle it, or hold another to select the intervening range. Moving your finger cancels a pending hold so the list can scroll normally.

Detection is enabled by default. Rebuild the container with `docker compose up -d --build` (or `sudo nerdctl compose up -d --build`). Background analysis starts after 20 seconds, scans all sources every minute, and waits two minutes after the last file write before processing. Existing recordings are analyzed as well as new ones. Files are processed serially, newest first within each camera. The model stays loaded between recordings.

The defaults target an **Intel Core i7-6700 @ 3.40 GHz** without a GPU: one worker, two ONNX inference threads, two sampled frames per second, and the nano model at 640 pixels. FFmpeg uses one decoder/filter thread. ONNX threads sleep between frames instead of spinning. The analysis runs through recordings as background work, so the sample rate describes video-time coverage, not a promised real-time throughput. Lower `detection_fps` to `1` or use a 320-pixel export if the server is busy; increase the sample rate for brief events. Lower resolutions and sample rates can miss small or fast subjects. Measurements on the actual server are needed to tune throughput.

Humans, animals and bikes require three consecutive matched observations with at least 65% confidence by default. Animal classes in the pretrained COCO model are bird, cat, dog, horse, sheep, cow, elephant, bear, zebra and giraffe; other wildlife is not covered by this pretrained class list. Cars, trucks and buses require at least three consecutive observations and consistent translation of both box edges by 2.5% of the model frame within two seconds. A single box jump, back-and-forth jitter, or box resizing does not qualify as vehicle motion. Detection intervals include up to one sample of padding and merge nearby same-category confirmed events across gaps of up to 0.5 seconds by default; the merge setting does not relax consecutive confirmation. Timing is approximate to the sampling cadence.

Settings provide separate **Human confidence threshold** (`detection_person_confidence`, default `0.65`) and **Human observations to confirm** (`detection_person_min_frames`, default `3`). Other-object confidence and confirmation settings remain separate. Save & restart to apply changes and invalidate active recordings’ previous analysis. Older configurations inherit their existing general detection thresholds for humans until changed. Sampling rate and interval merging remain shared.

Candidates down to 40% confidence (or the configured confidence if lower) remain ambiguous when they fail confirmation, preserving the recording without adding event markers. Increasing confidence or the confirmation count does not make these candidates eligible for the no-target pruning policy. Above 50% disk use, clips containing only these unconfirmed candidates move to quarantine for review. This can leave noisy recordings in the browser for review. Thresholds trade event coverage for fewer confirmed false positives; review brief and distant subjects before tightening them further.

Vehicle confirmation needs at least three samples within a two-second motion window, so use at least one sampled frame per second for vehicle markers. At 0.5 FPS, observed vehicle movement remains ambiguous and the recording is retained.

Use **Settings → Object detection** to adjust confidence, sample rate, confirmation count, movement threshold, scan timing and the no-target policy:

- **Move to recoverable quarantine** (default): clips with no qualifying events leave the recording browser and move into `.neonvr_pruned` within their camera directory. This keeps the original video and its metadata on disk for review, and does **not** reclaim disk space.
- **Keep no-target recordings**: embed intervals and display markers without no-target pruning. The separate space-pressure policies still apply.
- **Permanently delete**: remove fully analyzed no-target clips and their cached thumbnails. This option reclaims disk space.

Favorites are always preserved by detection pruning. Brief unconfirmed detections stay in the main library until disk use exceeds 50%, then move to quarantine when no confirmed intervals exist. Decode errors, incomplete video coverage, missing dependencies/models, metadata-write failures, timeouts and changing recordings all keep the clip. No-target classification filters recordings caused by leaves; it is not a separate leaf-recognition model. Sampled object detection can miss events, so review quarantined clips when tuning settings before enabling permanent deletion.

Right-click a card or thumbnail, or tap its **⋮** button on mobile, to download, delete, quarantine, or recover the clip. Recovery restores the original relative path and records your review decision, so automatic no-target pruning and the 50% quarantine rule do not undo it. Favorite status follows the clip between locations. Moves refuse to overwrite existing recordings. Age-based retention and vehicle-only space-pressure deletion still apply to non-favorite recovered clips; favorite them to preserve them. The existing age-based retention policy continues to operate on active recordings; it does not purge the quarantine.

### MP4 metadata and viewer

FFmpeg adds a JSON `neonvr_detections` format tag using MP4 `mdta` metadata and stream-copy remuxing. Audio/video are not re-encoded. Existing streams and container duration are checked before atomic replacement, and file permissions, ownership and modification time are preserved. The timestamp fallback for loosely named recordings therefore stays unchanged.

The tag contains the model and its checksum, analysis policy, duration, YOLO sample rate and ROI frame count, MOG2 sample rate and whole-clip frame count, padded `regions` (`start`, `end`, `frames`), completeness/ambiguity flags, analysis timestamp and one or more `{kind, labels, start, end, confidence}` intervals. `start`/`end` are seconds relative to playback, not wall-clock times. Inspect a clip with:

```bash
ffprobe -v error -show_entries format_tags=neonvr_detections -of json recording.mp4
```

The browser uses `/api/detections?camera=NAME&file=RELATIVE_PATH`. It refreshes queued analyses, clears markers when switching clips, and shows timelines aligned with the playback scrubber: padded MOG2-selected windows use teal curly `{` and `}` brackets on the same track as camera samples and confirmed detections; human, animal, bike and moving vehicle intervals use colored opening `[` and closing `]` brackets. Click or keyboard-activate either endpoint to seek to its timestamp; no filled AI interval obscures the camera samples. The indicator reports which categories are active at the current playback time. Status messages occupy fixed-height rows below the legend; their text cannot resize the video. Changing analysis thresholds or the detection policy version invalidates cached and embedded results and causes active recordings to be reanalyzed.

### Intel hardware video decoding

In **Settings → Object detection**, enable **Intel hardware video decoding** and choose **Save & restart**. The default render device is `/dev/dri/renderD128`. This enables VA-API decoding for the shared MOG2/YOLO frame stream. Frame-rate sampling happens before downloading decoded frames to CPU memory, avoiding transfers of frames that will be discarded. Resizing remains on the CPU; GPU-side scaling was not supported by the tested driver. MOG2, image resizing and YOLO inference still run on the CPU; this option does not enable GPU neural-network inference.

For Docker, the main Compose file includes the Intel GPU device mapping. Recreate the container to apply it:

```bash
docker compose up -d --build
```

The image includes Intel's iHD VA-API driver, the legacy i965 driver, and `vainfo`. For a different **host** render device, set `NEONVR_VAAPI_DEVICE=/dev/dri/renderD129` before the Compose command; Compose still maps it to `/dev/dri/renderD128` **inside** the container, which is the path to keep in Settings. On hosts without an Intel render device, remove the Intel GPU devices block from `docker-compose.yml`.

To check driver and device access inside the container:

```bash
docker exec neonvr vainfo --display drm --device /dev/dri/renderD128
```

Native installs need a working VA-API driver and permission to access the render device (usually membership in the host's `render` group). A container configured to run as a non-root user also needs its render-device group permissions. To try the legacy driver on older Intel graphics, set `LIBVA_DRIVER_NAME=i965` in the process/container environment.

Unsupported codecs, missing devices and hardware decoding errors trigger a restart of the complete shared stream using software decoding. If software decoding also fails, the recording is retained without publishing a complete result. The detection tooltip reports the decoder and fallback reason, and the server logs any fallback. Hardware decoding benefits depend on codec, resolution and GPU-to-CPU transfer overhead; compare processing times on your recordings.

The configuration keys are `detection_intel_acceleration` (default `false`) and `detection_vaapi_device`. Enabling acceleration changes the analysis policy, so active recordings are reanalyzed. Existing quarantine and favorite behavior is preserved.

### Native CPU detection setup

Install FFmpeg, Python 3.12 and the runtime packages, then export weights once using a separate environment. Keeping the exporter separate avoids installing PyTorch in the running server's environment:

```bash
python3 -m venv .venv-detection
.venv-detection/bin/pip install -r detection/requirements.txt

python3 -m venv .venv-export
.venv-export/bin/pip install torch==2.9.1 torchvision==0.24.1 --index-url https://download.pytorch.org/whl/cpu
.venv-export/bin/pip install ultralytics==8.4.175 onnx==1.23.2
.venv-export/bin/python detection/export_model.py --output models/yolo26m.onnx

go build -o neonvr .
./neonvr -config config/neonvr.json \
  -detection-python "$PWD/.venv-detection/bin/python" \
  -detection-model "$PWD/models/yolo26m.onnx"
```

These flags persist to configuration, as other legacy flags do. To make a smaller model for a busy CPU, pass `--size 320` to `export_model.py`, or explicitly export `--weights yolo26s.pt` and select that model. To run the recording browser without detection, set `detection_enabled` to `false` or use `-detection-enabled=false`.

Detection logic and MP4 integration checks:

```bash
python3 -m unittest discover -s detection -p 'test_*.py' -v
go test -race ./...
```

To include real model inference and the background prune pipeline in the Go checks, set `NEONVR_TEST_MODEL` to an absolute ONNX path and `NEONVR_TEST_PYTHON` to the detection environment's Python executable. Ordinary Go tests do not download models; MP4 tests skip if FFmpeg/FFprobe are unavailable.

Set `NEONVR_TEST_GO2RTC` to an installed go2rtc executable to include the live gateway lifecycle check. It verifies explicit enable/disable, shutdown, temporary credential cleanup, and that enabling alone does not open a camera connection. The test uses a local dummy endpoint.

Model/export reference: [Ultralytics YOLO26](https://docs.ultralytics.com/models/yolo26), [ONNX export](https://docs.ultralytics.com/modes/export).

## Additional clip directories

The server can browse additional directories alongside the primary BlinkClips root. Repeat `-extra-root` for as many sources as needed:

```bash
./neonvr \
  -root /mnt/mySSD/BlinkClips \
  -extra-root /home/ubuntu/FTPfiles/yicamBack \
  -extra-root /path/to/another/source
```

Each extra root appears in the camera/source selector using the directory basename (for example `yicamBack`). Extra roots support MP4 files directly in the directory and Yi-style hour folders such as `2026Y09M21D16H/`. For a file such as `W440M00S60.mp4` in that folder, the server derives `2026-09-21 16:40:00`; loose MP4 files that do not contain a full timestamp fall back to their filesystem modification time.

Yi filenames without the `W4` prefix are also recognized automatically: `2020Y04M17D16H/16M36S06.mp4` means `2020-04-17 16:16:36` in server local time. In both variants, the numbers before `M` and `S` are minutes and seconds; the final number is the clip duration. This timestamp is used for date filtering, sorting, and age-based retention, even when the file was uploaded more recently.

When using Docker/nerdctl, every extra root must also be mounted into the container and listed in `extra_roots` in the config. The included Compose file mounts both `/home/ubuntu/FTPfiles/yicamBack` and `/home/ubuntu/FTPfiles/yicamFront`; the config lists them in `extra_roots`.

Automatic age-based deletion applies to the primary `-root` and all configured `-extra-root` directories. Favorited clips are preserved everywhere.

### Thumbnail selection

New thumbnails first use compact Yi camera readings from the `yi_motion` tag. The server extracts a frame at the first valid sample's exact playback timestamp, calculated from the MP4's UTC `creation_time`. If that frame is unavailable or blank, it tries subsequent retained samples. A continuing motion does not receive a fabricated onset at the beginning of the video.

If camera readings are missing, invalid, unsupported, empty, or frame extraction fails, the server uses the existing divergence selector designed for mostly-static security footage. It samples nine frames across the clip, builds a small per-pixel median baseline representing the common scene, and selects the frame whose luminance pattern differs most from that median. Obviously flat near-black or near-white frames are excluded so upload/camera glitches do not win merely by being different. If duration probing is unavailable, it tries fixed offsets at 0.5, 2, 5, 15, and 30 seconds; if candidate extraction still fails, it falls back to `-thumb-seek` and then 0 seconds. The cache key includes a thumbnail-algorithm version, so existing cached thumbnails are regenerated as they are requested.


### Automatic thumbnail maintenance

Thumbnail maintenance runs in the background every 30 minutes by default in the included compose configuration. Each pass scans all configured clip roots, pre-generates any missing thumbnails using the same motion metadata and fallback selector as browser requests, and removes cached `.jpg` thumbnails that no longer correspond to a current video (including stale thumbnails from deleted videos or older thumbnail algorithms).

Use `-thumb-maintenance-every DURATION` to change the cadence, for example `-thumb-maintenance-every 1h`. Set it to `0` to disable background maintenance. Lazy/on-demand thumbnail generation remains enabled as a fallback.

### Yi timestamps

Yi cameras write folder/file timestamps in their local timezone. NeoNVR uses those dates and times as written for display and date grouping, without a timezone correction. For example, `2026Y09M23D23H/W421M41S18.mp4` appears as `2026-09-23 23:21:41`.

### Compact Yi camera samples

New y623 SD recordings contain a compact `yi_motion` tag with a single `schema` column header and `data` arrays of raw readings. NeoNVR reads this format exclusively; the previous `yi-motion-v1` event/level parser has been removed. Older recordings remain playable, retain their original tags, and use scene-based thumbnail selection. The player labels their camera metadata as unsupported; AI interval brackets remain available independently.

```bash
ffprobe -v error -show_entries format_tags=yi_motion -of json recording.mp4 |
  jq '.format.tags.yi_motion | fromjson' > motion.json
jq -r '.schema, .data[] | @csv' motion.json
```

```json
{
  "schema": ["unix_ms", "Scene", "Move", "MovingLevel", "BinImgRatio", "MovingTh"],
  "data": [
    [1791858420000, 0, 1, 1, 2.10, 20],
    [1791858420200, 0, 2, 2, 4.25, 20],
    [1791858420400, 0, 2, 2, 3.80, 20]
  ]
}
```

The reader follows the header's column order and keeps the raw header and rows intact, including additional vendor columns. For playback it validates documented ranges: `Scene` 0–2, `Move` 0–4, `MovingLevel` 0–3, `BinImgRatio` 0–100%, and `MovingTh` 10–31. Invalid rows and readings outside the recording are excluded from displayed samples. `MovingTh` is the vendor block-difference threshold, separate from the detector's START/STOP percentages. These values are measurements, not object counts or probabilities.

Each sample's playback position is `unix_ms / 1000 - UTC creation_time`. NeoNVR uses the MP4 format's `creation_time`, or a valid video stream's creation tag if the format tag is missing or invalid. It never uses local filenames or the first sample as the origin. Creation time has one-second precision on these cameras, so the positions also reflect that limit and polling latency.

The single timeline below the video scrubber shows **colored AI interval brackets** and **camera intensity marks at actual sample timestamps**. The intensity scale runs **green at MovingLevel 0 through yellow/orange to red at MovingLevel 3**. Hover or focus a sample to see all vendor readings and its absolute timestamp; clicking seeks to its exact offset. The current sample details are shown when the playback position is within 100 ms of a retained reading.

No motion intervals or trigger events are reconstructed from measurements. Burst samples roughly 200 ms apart appear at their actual spacing, as do later roughly one-second samples. Gaps remain unpainted because they can mean idle time, invalid readings, or missing data. A clip that begins during active motion shows only its continuing readings. Empty `data` reports no valid retained active samples; a missing tag reports no camera samples. Neither case is treated as evidence of idle coverage. Metadata errors, missing UTC creation time, and invalid row counts are shown separately in the player status or its tooltip.

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
