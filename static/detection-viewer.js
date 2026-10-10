// Fetch clip-relative intervals, refresh pending analyses, and show active events.
(function () {
  'use strict';
  const video = document.getElementById('mainVideo');
  const seek = document.getElementById('mainSeek');
  const status = document.getElementById('detectionStatus');
  if (!video || !seek || !status) return;
  const names = { person: 'Human', animal: 'Animal', bike: 'Bike', moving_car: 'Moving vehicle' };
  let generation = 0;
  let request = null;
  let timer = null;
  let result = null;

  function activeStatus() {
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
    clearTimeout(timer);
    result = null;
    seek.setDetections(null);
    status.classList.remove('active');
    status.textContent = 'Loading detections…';
    status.title = '';
    refresh(clip, generation);
  }

  video.addEventListener('timeupdate', activeStatus);
  video.addEventListener('seeked', activeStatus);
  window.addEventListener('pagehide', () => { request?.abort(); clearTimeout(timer); });
  window.NeoNVRDetections = { loadClip };
})();
