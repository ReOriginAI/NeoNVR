// Fetch clip-relative intervals, refresh pending analyses, and show active events.
(function () {
  'use strict';
  const video = document.getElementById('mainVideo');
  const seek = document.getElementById('mainSeek');
  const status = document.getElementById('detectionStatus');
  const motionLegend = document.getElementById('cameraMotionLegend');
  const motionStatus = document.getElementById('cameraMotionStatus');
  const intensityLegend = document.getElementById('cameraIntensityLegend');
  if (!video || !seek || !status) return;
  const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
  let generation = 0;
  let request = null;
  let timer = null;
  let result = null;
  let motionRequest = null;
  let motion = null;

  function activeCameraMotion() {
    if (!motion || !motionStatus) return;
    const sample = seek.cameraSampleAt(video.currentTime);
    motionStatus.classList.toggle('active', !!sample);
    const text = sample ? 'Camera sample at ' + sample.time.toFixed(3) + ' sec: ' + seek.cameraSampleDetails(sample) + '.' :
      motion.samples.length ? motion.samples.length + ' camera samples; no reading at this position.' : 'Camera metadata found; no valid retained active samples.';
    if (motionStatus.textContent !== text) motionStatus.textContent = text;
  }

  async function refreshCameraMotion(clip, expected) {
    motionRequest = new AbortController();
    try {
      const response = await fetch('/api/camera-motion?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: motionRequest.signal});
      if (!response.ok) {
        const failure = await response.json().catch(() => ({}));
        const error = new Error('Could not read camera motion metadata.');
        error.detail = failure.message || failure.error || 'HTTP ' + response.status;
        throw error;
      }
      const data = await response.json();
      if (expected !== generation) return;
      motion = data.motion;
      seek.setCameraMotion(motion);
      if (motionLegend) motionLegend.hidden = !motion;
      if (intensityLegend) intensityLegend.hidden = !motion?.samples.length;
      if (motionStatus) {
        motionStatus.hidden = false;
        motionStatus.textContent = motion ? '' : data.unsupported ? 'Camera metadata format not supported.' : 'No camera samples in this recording.';
        motionStatus.title = motion ? 'UTC creation_time: ' + motion.creation_time + '; retained rows: ' + motion.data.length +
          '; valid samples: ' + motion.samples.length + '; invalid readings: ' + motion.invalid_rows + '; outside clip: ' + motion.outside_clip_rows : data.message || '';
      }
      activeCameraMotion();
    } catch (error) {
      if (expected !== generation || error.name === 'AbortError') return;
      if (motionStatus) {
        motionStatus.hidden = false;
        motionStatus.textContent = error.message;
        motionStatus.title = error.detail || '';
      }
    }
  }

  function activeStatus() {
    activeCameraMotion();
    if (!result) return;
    const active = [...new Set(result.intervals.filter(row => video.currentTime >= row.start && video.currentTime < row.end).map(row => names[row.kind]))];
    status.classList.toggle('active', active.length > 0);
    const text = active.length ? 'Detected now: ' + active.join(', ') :
      (result.intervals.length ? 'No detection at this time.' : result.uncertain ? 'Unconfirmed observations — clip retained.' : 'No target detections.');
    if (status.textContent !== text) status.textContent = text;
  }

  async function refresh(clip, expected) {
    if (expected !== generation) return;
    request = new AbortController();
    try {
      const response = await fetch('/api/detections?' + new URLSearchParams({camera: clip.camera, file: clip.file}), {cache: 'no-store', signal: request.signal});
      if (!response.ok) throw new Error(response.status === 404 ? 'Clip is no longer available.' : 'Detection data unavailable.');
      const data = await response.json();
      if (expected !== generation) return;
      if (data.status === 'complete') {
        result = data.detections;
        status.title = '';
        seek.setDetections(result);
        activeStatus();
      } else {
        status.classList.remove('active');
        const labels = {pending: 'Detection queued…', processing: 'Analyzing detections…', disabled: 'Detection disabled.', error: 'Analysis unavailable — clip retained.'};
        status.textContent = labels[data.status] || 'Detection pending…';
        status.title = data.message || '';
        if (data.status !== 'disabled') timer = setTimeout(() => refresh(clip, expected), 10000);
      }
    } catch (error) {
      if (expected !== generation || error.name === 'AbortError') return;
      status.textContent = error.message;
      if (!error.message.includes('no longer available')) timer = setTimeout(() => refresh(clip, expected), 10000);
    }
  }

  function loadClip(clip) {
    generation += 1;
    request?.abort();
    motionRequest?.abort();
    clearTimeout(timer);
    result = null;
    motion = null;
    seek.setDetections(null);
    seek.setCameraMotion(null);
    if (motionLegend) motionLegend.hidden = true;
    if (intensityLegend) intensityLegend.hidden = true;
    if (motionStatus) {
      motionStatus.hidden = false;
      motionStatus.textContent = 'Reading camera motion metadata…';
      motionStatus.title = '';
      motionStatus.classList.remove('active');
    }
    status.classList.remove('active');
    status.textContent = 'Loading detections…';
    status.title = '';
    refresh(clip, generation);
    refreshCameraMotion(clip, generation);
  }

  video.addEventListener('timeupdate', activeStatus);
  video.addEventListener('seeked', activeStatus);
  window.addEventListener('pagehide', () => { request?.abort(); motionRequest?.abort(); clearTimeout(timer); });
  window.NeoNVRDetections = { loadClip };
})();
