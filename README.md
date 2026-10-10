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
- Explicit Select mode for bulk deletion; Shift-click selects every visible clip in between
- Favorites tab
- Optional built-in auto-delete of non-favorite clips older than X days
- YOLO26n detection of humans, COCO animals, bicycles, motorcycles, and moving cars/buses/trucks
- Clip-relative detection intervals embedded in the MP4, with colored, clickable lanes below the video scrubber and a live detection indicator
- Optional pruning of fully analyzed no-target clips, using recoverable quarantine by default
- Displays total size of the BlinkClips folder only, avoiding misleading parent-folder stats inside Docker

## Docker image

The Go binary is built on Alpine. A separate Python build stage downloads the official `yolo26n.pt` weights and exports a fixed 640-pixel, end-to-end ONNX model. The Debian Bookworm runtime includes FFmpeg, Python, NumPy, and **CPU-only ONNX Runtime**. PyTorch, Ultralytics and GPU libraries stay in the model build stage. The model is included in the image, so analysis and the browser work offline after building. The first build needs network access for packages and model weights.

## Docker configuration

The included Compose service starts with:

```yaml
command: ["-config", "/config/neonvr.json"]
volumes:
  - ./config:/config
```

Runtime options live in `config/neonvr.json`, rather than in the Compose command. The included recording mounts are unchanged. Retention is configured by `auto_delete_days` (14 in the supplied file; 0 disables deletion) and `auto_delete_every` (`6h`). It skips favorites in every recording source.

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

Detection is enabled by default. Rebuild the container with `docker compose up -d --build` (or `sudo nerdctl compose up -d --build`). Background analysis starts after 20 seconds, scans all sources every minute, and waits two minutes after the last file write before processing. Existing recordings are analyzed as well as new ones. Files are processed serially, newest first within each camera. The model stays loaded between recordings.

The defaults target an **Intel Core i7-6700 @ 3.40 GHz** without a GPU: one worker, two ONNX inference threads, two sampled frames per second, and the nano model at 640 pixels. FFmpeg uses one decoder/filter thread. ONNX threads sleep between frames instead of spinning. The analysis runs through recordings as background work, so the sample rate describes video-time coverage, not a promised real-time throughput. Lower `detection_fps` to `1` or use a 320-pixel export if the server is busy; increase the sample rate for brief events. Lower resolutions and sample rates can miss small or fast subjects. Measurements on the actual server are needed to tune throughput.

Humans, animals and bikes require two matched observations by default. Animal classes in the pretrained COCO model are bird, cat, dog, horse, sheep, cow, elephant, bear, zebra and giraffe; other wildlife is not covered by this pretrained class list. Cars, trucks and buses qualify only when their tracked center moves at least 1.5% of the frame within two seconds. A parked car does not keep a leaf-triggered recording. Detection intervals include up to one sample of padding and merge nearby same-category events; timing is approximate to the sampling cadence.

Use **Settings → Object detection** to adjust confidence, sample rate, confirmation count, movement threshold, scan timing and the no-target policy:

- **Move to recoverable quarantine** (default): clips with no qualifying events leave the recording browser and move into `.neonvr_pruned` within their camera directory. This keeps the original video and its metadata on disk for review, and does **not** reclaim disk space.
- **Keep all recordings**: embed intervals and display markers, retaining every clip.
- **Permanently delete**: remove fully analyzed no-target clips and their cached thumbnails. This option reclaims disk space.

Favorites are always preserved by detection pruning. Brief unconfirmed detections, decode errors, incomplete video coverage, missing dependencies/models, metadata-write failures, timeouts and changing recordings all keep the clip. No-target classification filters recordings caused by leaves; it is not a separate leaf-recognition model. Sampled object detection can miss events, so review quarantined clips when tuning settings before enabling permanent deletion.

To restore a quarantined clip, first set the no-target policy to **Keep all recordings** and restart, then move the MP4 from `.neonvr_pruned` back to its original relative path in that camera directory. The existing age-based retention policy continues to operate on active recordings; it does not purge the quarantine.

### MP4 metadata and viewer

FFmpeg adds a JSON `neonvr_detections` format tag using MP4 `mdta` metadata and stream-copy remuxing. Audio/video are not re-encoded. Existing streams and container duration are checked before atomic replacement, and file permissions, ownership and modification time are preserved. The timestamp fallback for loosely named recordings therefore stays unchanged.

The tag contains the model and its checksum, analysis policy, duration, sample rate, frame count, completeness/ambiguity flags, analysis timestamp and one or more `{kind, labels, start, end, confidence}` intervals. `start`/`end` are seconds relative to playback, not wall-clock times. Inspect a clip with:

```bash
ffprobe -v error -show_entries format_tags=neonvr_detections -of json recording.mp4
```

The browser uses `/api/detections?camera=NAME&file=RELATIVE_PATH`. It refreshes queued analyses, clears markers when switching clips, and shows four colored lanes aligned with the playback scrubber: human, animal, bike and moving vehicle. Click a colored interval to seek to its beginning. The indicator reports which categories are active at the current playback time. Changing analysis thresholds invalidates the cache and causes active recordings to be reanalyzed.

### Native CPU detection setup

Install FFmpeg, Python 3.12 and the runtime packages, then export weights once using a separate environment. Keeping the exporter separate avoids installing PyTorch in the running server's environment:

```bash
python3 -m venv .venv-detection
.venv-detection/bin/pip install -r detection/requirements.txt

python3 -m venv .venv-export
.venv-export/bin/pip install torch==2.9.1 torchvision==0.24.1 --index-url https://download.pytorch.org/whl/cpu
.venv-export/bin/pip install ultralytics==8.4.175 onnx==1.23.2
.venv-export/bin/python detection/export_model.py --output models/yolo26n.onnx

go build -o neonvr .
./neonvr -config config/neonvr.json \
  -detection-python "$PWD/.venv-detection/bin/python" \
  -detection-model "$PWD/models/yolo26n.onnx"
```

These flags persist to configuration, as other legacy flags do. To make a smaller model for a busy CPU, pass `--size 320` to `export_model.py`. To run the recording browser without detection, set `detection_enabled` to `false` or use `-detection-enabled=false`.

Detection logic and MP4 integration checks:

```bash
python3 -m unittest discover -s detection -p 'test_*.py' -v
go test -race ./...
```

To include real model inference and the background prune pipeline in the Go checks, set `NEONVR_TEST_MODEL` to an absolute ONNX path and `NEONVR_TEST_PYTHON` to the detection environment's Python executable. Ordinary Go tests do not download models; MP4 tests skip if FFmpeg/FFprobe are unavailable.

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
sudo nerdctl compose logs -f neonvr
```

For a native systemd installation, install the executable at `/opt/NeoNVR/neonvr` and use `neonvr.service.example` as the unit template.
