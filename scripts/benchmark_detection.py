"""Read-only benchmark inside the running container.
Run from the repository: sudo nerdctl exec -i neonvr python3 - < scripts/benchmark_detection.py
Optional arguments after python3 -: --clip /recordings/example.mp4 --seconds 30
Does not call analysis persistence, pruning, or any server API.
"""
import argparse
import gc
import hashlib
import json
import math
import os
from pathlib import Path
import subprocess
import time
import types
from contextlib import closing


def emit(value):
    print(json.dumps(value), flush=True)


def runtime_worker():
    for path in Path('/proc').glob('[0-9]*/cmdline'):
        try:
            args = path.read_bytes().decode().split('\0')
            if '-c' in args and '--model' in args:
                source = args[args.index('-c') + 1]
                if 'class Detector:' not in source or 'motion_regions' not in source:
                    continue
                module = types.ModuleType('neonvr_benchmark_worker')
                exec(compile(source, '<running-neonvr-worker>', 'exec'), module.__dict__)
                model = Path(args[args.index('--model') + 1]).expanduser()
                if not model.is_absolute():
                    model = Path(os.environ.get('NEONVR_MODEL_DIR', '.')) / model
                return module, str(model)
        except (OSError, UnicodeError, ValueError):
            continue
    raise RuntimeError('No active two-pass NeoNVR Python worker found; check container logs.')


def sample_clip(config):
    for root in [config['root']] + config.get('extra_roots', []):
        examined = 0
        for folder, dirs, files in os.walk(root):
            dirs[:] = sorted(d for d in dirs if not d.startswith('.'))
            for name in sorted(files):
                examined += 1
                path = Path(folder) / name
                if name.lower().endswith('.mp4') and not path.is_symlink() and time.time() - path.stat().st_mtime > 120:
                    return str(path)
            if examined > 5000:
                break
    raise RuntimeError('No mature recording found; pass --clip with a container video path.')


def optimize_bench(worker, model, config, request, duration, decode_only=False):
    import numpy as np
    # Repeat variants to reduce the influence of cache warmup and live load.
    fps = config.get('detection_motion_fps', 5)
    common = [request['ffmpeg'], '-nostdin', '-v', 'error', '-threads', '1']
    sample = f"setpts=PTS-STARTPTS,fps=fps={fps}:start_time=0:round=up:eof_action=pass,"
    software = sample + 'scale=320:180:force_original_aspect_ratio=decrease,pad=320:180:(ow-iw)/2:(oh-ih)/2:color=0x727272'
    gpu = sample + 'scale_vaapi=w=320:h=180:force_original_aspect_ratio=decrease:force_divisible_by=2:format=nv12,hwdownload,format=nv12,pad=320:180:(ow-iw)/2:(oh-ih)/2:color=0x727272'
    for trial in range(2):
        original_hash = None
        for name, filters, hardware in [('software',software,False),('original_vaapi','hwdownload,format=nv12,'+software,True),('sample_before_download',sample+'hwdownload,format=nv12,'+software[len(sample):],True),('gpu_resize_sample',gpu,True),('gpu_resize_i965',gpu,True)]:
            command = list(common)
            if hardware:
                command += ['-hwaccel','vaapi','-hwaccel_device',request['vaapi_device'],'-hwaccel_output_format','vaapi']
            command += ['-i',request['path'],'-t',str(duration),'-an','-sn','-dn','-vf',filters,'-filter_threads','1','-threads','1','-f','rawvideo','-pix_fmt','bgr24','pipe:1']
            started=time.perf_counter()
            run=subprocess.run(command,stdout=subprocess.PIPE,stderr=subprocess.PIPE,timeout=120,env=dict(os.environ,LIBVA_DRIVER_NAME='i965') if name.endswith('i965') else None)
            output_hash = hashlib.sha256(run.stdout).hexdigest()
            if name == 'original_vaapi' and not run.returncode:
                original_hash = output_hash
            emit({'sha256':output_hash,'matches_original_vaapi':output_hash == original_hash if original_hash else None,
                  'stage':'decode_optimization','trial':trial,'variant':name,'seconds':round(time.perf_counter()-started,3),
                  'frames':len(run.stdout)//(320*180*3),'error':run.stderr.decode()[:500] if run.returncode else None})
    if decode_only:
        return
    # Use identical real frames for both runtimes, and compare raw model outputs.
    detector=worker.Detector(model,config.get('detection_threads',6))
    with closing(worker.video_frames(dict(request,intel_acceleration=False),2,detector.width,detector.height,0,min(duration,4))) as samples:
        frames=[raw for _,raw in samples]
    tensors=[np.ascontiguousarray(np.frombuffer(raw,dtype=np.uint8).reshape(detector.height,detector.width,3)[:,:,::-1].transpose(2,0,1)[None],dtype=np.float32)/255 for raw in frames]
    import openvino as ov
    core=ov.Core()
    emit({'stage':'openvino_environment','version':ov.__version__,'devices':core.available_devices})
    baseline=[detector.session.run(None,{detector.input_name:t})[0] for t in tensors]
    for device in ['ONNX','CPU','GPU']:
        try:
            if device=='ONNX':
                predict=lambda t:detector.session.run(None,{detector.input_name:t})[0]
            else:
                options={'PERFORMANCE_HINT':'LATENCY'}
                if device=='CPU':
                    options.update(INFERENCE_NUM_THREADS=config.get('detection_threads',6),INFERENCE_PRECISION_HINT='f32')
                compiled=core.compile_model(model,device,options)
                predict=lambda t:compiled([t])[compiled.output(0)].copy()
            predict(tensors[0])
            for trial in range(2):
                started=time.perf_counter()
                outputs=[predict(t) for t in tensors]
                elapsed=time.perf_counter()-started
                emit({'stage':'inference_optimization','device':device,'trial':trial,'frames':len(tensors),
                      'milliseconds_per_frame':round(elapsed/len(tensors)*1000,2),
                      'max_output_difference':float(max(np.max(np.abs(a-b)) for a,b in zip(baseline,outputs)))})
            if device!='ONNX':
                del compiled
                gc.collect()
        except Exception as error:
            emit({'stage':'inference_optimization','device':device,'error':str(error)[:1000]})
    del detector


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--config', default='/config/neonvr.json')
    parser.add_argument('--clip')
    parser.add_argument('--optimize', action='store_true')
    parser.add_argument('--decode-only', action='store_true')
    parser.add_argument('--seconds', type=float, default=30)
    args = parser.parse_args()
    if not 1 <= args.seconds <= 120:
        parser.error('--seconds must be between 1 and 120')
    config = json.loads(Path(args.config).read_text())
    worker, model = runtime_worker()
    clip = args.clip or sample_clip(config)
    ffmpeg = config.get('ffmpeg', 'ffmpeg')
    ffprobe = str(Path(ffmpeg).with_name('ffprobe')) if '/' in ffmpeg else 'ffprobe'
    probe = json.loads(subprocess.check_output([ffprobe, '-v', 'error', '-show_entries', 'format=duration', '-of', 'json', clip], timeout=30))
    duration = min(float(probe['format']['duration']), args.seconds)
    if not math.isfinite(duration) or duration <= 0:
        raise RuntimeError('Invalid clip duration')
    request = dict(path=clip, ffmpeg=ffmpeg, ffprobe=ffprobe,
                   vaapi_device=config.get('detection_vaapi_device', '/dev/dri/renderD128'),
                   fps=config.get('detection_fps', 2), motion_fps=config.get('detection_fps', 2),
                   motion_threshold=config.get('detection_motion_threshold', 16),
                   motion_min_area=config.get('detection_motion_min_area', .001),
                   roi_pre_roll=config.get('detection_roi_pre_roll', 2),
                   roi_post_roll=config.get('detection_roi_post_roll', 3))
    if args.optimize:
        optimize_bench(worker,model,config,request,duration,args.decode_only)
        return
    emit({'note': 'Live service remains running; contention affects timings.', 'clip': clip, 'sample_seconds': duration})
    for hardware in (False, True):
        request['intel_acceleration'] = hardware
        started = time.perf_counter()
        try:
            with closing(worker.video_frames(request, request['motion_fps'], 320, 180, 0, duration)) as frames:
                count = sum(1 for _ in frames)
            emit({'stage': 'decode_resize_sample', 'backend': 'vaapi' if hardware else 'software',
                  'seconds': round(time.perf_counter()-started, 3), 'frames': count})
        except Exception as error:
            emit({'stage': 'decode_resize_sample', 'backend': 'vaapi' if hardware else 'software', 'error': str(error)[:500]})
    request['intel_acceleration'] = config.get('detection_intel_acceleration', False)
    started = time.perf_counter()
    try:
        regions, _, count = worker.motion_regions(request, duration)
    except worker.HardwareDecodeError:
        request['intel_acceleration'] = False
        started = time.perf_counter()
        regions, _, count = worker.motion_regions(request, duration)
    selected = sum(region['end']-region['start'] for region in regions)
    emit({'stage': 'decode_resize_and_mog2', 'seconds': round(time.perf_counter()-started, 3),
          'frames': count, 'roi_seconds': round(selected, 3), 'roi_fraction': round(selected/duration, 3),
          'backend': 'vaapi' if request['intel_acceleration'] else 'software'})
    for threads in (2, 4, 6):
        started = time.perf_counter()
        detector = worker.Detector(model, threads)
        load_seconds = time.perf_counter()-started
        # Use original video content, software decoded here to isolate inference.
        frame_request = dict(request, intel_acceleration=False)
        with closing(worker.video_frames(frame_request, 2, detector.width, detector.height, 0, min(duration, 4))) as samples:
            raw_frames = [raw for _, raw in samples]
        started = time.perf_counter()
        for raw in raw_frames:
            detector.predict(raw, .4)
        elapsed = time.perf_counter()-started
        emit({'stage': 'yolo_preprocess_and_inference', 'threads': threads, 'frames': len(raw_frames),
              'load_and_warmup_seconds': round(load_seconds, 3),
              'milliseconds_per_frame': round(elapsed/len(raw_frames)*1000, 2),
              'inference_fps': round(len(raw_frames)/elapsed, 2)})
        del detector
        gc.collect()


if __name__ == '__main__':
    main()
