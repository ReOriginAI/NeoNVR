// Media Chrome Lite for NeoNVR.
// Offline, dependency-free custom elements compatible with the subset of
// media-chrome controls used by this app. This is intentionally small and
// tailored to the main NeoNVR player.

(function () {
  function formatTime(seconds) {
    seconds = Number.isFinite(seconds) ? Math.max(0, seconds) : 0;
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = Math.floor(seconds % 60);
    if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`;
    return `${m}:${String(s).padStart(2, '0')}`;
  }

  function findVideo(el) {
    const controller = el.closest('media-controller');
    return controller ? controller.querySelector('video') : null;
  }

  function define(name, klass) {
    if (!customElements.get(name)) customElements.define(name, klass);
  }

  class MediaController extends HTMLElement {
    connectedCallback() {
      this.setAttribute('data-media-chrome-lite', '');
    }
  }

  class MediaControlBar extends HTMLElement {
    connectedCallback() {
      this.setAttribute('data-media-chrome-lite', '');
    }
  }

  class BaseButton extends HTMLElement {
    connectedCallback() {
      if (this.shadowRoot) return;
      const root = this.attachShadow({ mode: 'open' });
      this.button = document.createElement('button');
      this.button.type = 'button';
      this.button.part = 'button';
      this.button.setAttribute('aria-label', this.label || 'Media button');
      this.button.textContent = this.text || '•';
      const style = document.createElement('style');
      style.textContent = `
        button {
          border: 0;
          border-radius: 8px;
          padding: 7px 9px;
          min-width: 34px;
          background: rgba(255,255,255,.10);
          color: white;
          font: inherit;
          cursor: pointer;
        }
        button:hover { background: rgba(255,255,255,.18); }
      `;
      root.append(style, this.button);
      this.button.addEventListener('click', () => this.onClick());
      this.video = findVideo(this);
      if (this.video) this.bindVideo(this.video);
    }
    onClick() {}
    bindVideo() {}
  }

  class MediaPlayButton extends BaseButton {
    constructor() { super(); this.text = '▶'; this.label = 'Play or pause'; }
    bindVideo(video) {
      const update = () => { this.button.textContent = video.paused ? '▶' : '❚❚'; };
      video.addEventListener('play', update);
      video.addEventListener('pause', update);
      video.addEventListener('loadedmetadata', update);
      update();
    }
    onClick() {
      const video = findVideo(this);
      if (!video) return;
      if (video.paused) video.play().catch(() => {});
      else video.pause();
    }
  }

  class MediaSeekBackwardButton extends BaseButton {
    constructor() { super(); this.text = '↶'; this.label = 'Seek backward'; }
    onClick() {
      const video = findVideo(this);
      if (!video) return;
      const offset = Number(this.getAttribute('seek-offset') || '1');
      video.currentTime = Math.max(0, video.currentTime - offset);
    }
  }

  class MediaSeekForwardButton extends BaseButton {
    constructor() { super(); this.text = '↷'; this.label = 'Seek forward'; }
    onClick() {
      const video = findVideo(this);
      if (!video) return;
      const offset = Number(this.getAttribute('seek-offset') || '1');
      const max = Number.isFinite(video.duration) ? video.duration : Infinity;
      video.currentTime = Math.min(max, video.currentTime + offset);
    }
  }

  class MediaMuteButton extends BaseButton {
    constructor() { super(); this.text = '🔊'; this.label = 'Mute'; }
    bindVideo(video) {
      const update = () => { this.button.textContent = video.muted || video.volume === 0 ? '🔇' : '🔊'; };
      video.addEventListener('volumechange', update);
      update();
    }
    onClick() {
      const video = findVideo(this);
      if (!video) return;
      video.muted = !video.muted;
    }
  }

  class MediaFullscreenButton extends BaseButton {
    constructor() { super(); this.text = '⛶'; this.label = 'Fullscreen'; }
    onClick() {
      const controller = this.closest('media-controller');
      if (!controller) return;
      if (document.fullscreenElement) document.exitFullscreen().catch(() => {});
      else controller.requestFullscreen().catch(() => {});
    }
  }

  class MediaPlaybackRateButton extends BaseButton {
    constructor() { super(); this.text = '1×'; this.label = 'Playback speed'; }
    connectedCallback() {
      super.connectedCallback();
      const video = findVideo(this);
      if (!video || !this.button) return;
      const update = () => { this.button.textContent = `${video.playbackRate.toFixed(video.playbackRate % 1 ? 1 : 0)}×`; };
      video.addEventListener('ratechange', update);
      update();
    }
    onClick() {
      const video = findVideo(this);
      if (!video) return;
      const rates = [0.25, 0.5, 1, 1.5, 2, 4];
      const i = rates.findIndex(r => Math.abs(r - video.playbackRate) < 0.01);
      video.playbackRate = rates[(i + 1) % rates.length];
    }
  }

  class MediaTimeRange extends HTMLElement {
    connectedCallback() {
      if (this.shadowRoot) return;
      const root = this.attachShadow({ mode: 'open' });
      const input = document.createElement('input');
      input.type = 'range';
      input.min = '0';
      input.max = '1000';
      input.value = '0';
      input.part = 'range';
      input.setAttribute('aria-label', 'Video position');
      this.input = input;
      this.lanes = document.createElement('div');
      this.lanes.className = 'detection-lanes';
      this.lanes.setAttribute('aria-label', 'Detection timeline');
      this.lanes.hidden = true;
      const style = document.createElement('style');
      style.textContent = `
        :host { display: block; flex: 1 1 180px; min-width: 110px; }
        input { width: 100%; margin: 0; accent-color: white; }
        .detection-lanes { margin: 3px 8px 0; }
        .lane { height: 7px; position: relative; margin-top: 3px; background: #262a32; border-radius: 3px; }
        .event { position: absolute; top: 0; height: 100%; padding: 0; min-width: 2px; border: 0; border-radius: 2px; cursor: pointer; }
        .event:hover, .event:focus-visible { outline: 2px solid white; z-index: 1; }
      `;
      root.append(style, input, this.lanes);
      const video = findVideo(this);
      if (!video) return;
      let dragging = false;
      const update = () => {
        if (dragging) return;
        if (!Number.isFinite(video.duration) || video.duration <= 0) input.value = '0';
        else input.value = String(Math.round((video.currentTime / video.duration) * 1000));
      };
      input.addEventListener('input', () => {
        dragging = true;
        if (Number.isFinite(video.duration) && video.duration > 0) {
          video.currentTime = (Number(input.value) / 1000) * video.duration;
        }
      });
      input.addEventListener('change', () => { dragging = false; update(); });
      video.addEventListener('timeupdate', update);
      video.addEventListener('durationchange', update);
      video.addEventListener('loadedmetadata', update);
      video.addEventListener('durationchange', () => this.renderDetections());
      this.video = video;
      this.renderDetections();
      update();
    }
    setDetections(result) {
      this.detections = result;
      this.renderDetections();
    }
    renderDetections() {
      if (!this.lanes) return;
      this.lanes.replaceChildren();
      const result = this.detections;
      this.lanes.hidden = !result;
      if (!result) return;
      const duration = Number.isFinite(this.video?.duration) && this.video.duration > 0 ? this.video.duration : result.duration;
      if (!Number.isFinite(duration) || duration <= 0) return;
      const categories = [['person', 'Human', '#75d685'], ['animal', 'Animal', '#ffd166'], ['bike', 'Bike', '#72a7ff'], ['moving_car', 'Moving vehicle', '#c792ea']];
      for (const [kind, label, color] of categories) {
        const lane = document.createElement('div');
        lane.className = 'lane';
        lane.dataset.kind = kind;
        lane.setAttribute('aria-label', label + ' detections');
        for (const interval of result.intervals || []) {
          if (interval.kind !== kind || !Number.isFinite(interval.start) || !Number.isFinite(interval.end)) continue;
          const start = Math.max(0, Math.min(duration, interval.start));
          const end = Math.max(start, Math.min(duration, interval.end));
          if (end <= start) continue;
          const event = document.createElement('button');
          event.type = 'button';
          event.className = 'event';
          event.style.left = (start / duration * 100) + '%';
          event.style.width = ((end - start) / duration * 100) + '%';
          event.style.background = color;
          const text = label + ': ' + start.toFixed(1) + '–' + end.toFixed(1) + ' sec (' + Math.round(interval.confidence * 100) + '%)';
          event.title = text;
          event.setAttribute('aria-label', 'Seek to ' + text);
          event.addEventListener('click', () => { if (this.video) this.video.currentTime = start; });
          lane.append(event);
        }
        this.lanes.append(lane);
      }
    }
  }

  class MediaTimeDisplay extends HTMLElement {
    connectedCallback() {
      this.textContent = '0:00';
      const video = findVideo(this);
      if (!video) return;
      const update = () => { this.textContent = formatTime(video.currentTime); };
      video.addEventListener('timeupdate', update);
      video.addEventListener('loadedmetadata', update);
      update();
    }
  }

  class MediaDurationDisplay extends HTMLElement {
    connectedCallback() {
      this.textContent = '0:00';
      const video = findVideo(this);
      if (!video) return;
      const update = () => { this.textContent = Number.isFinite(video.duration) ? formatTime(video.duration) : '0:00'; };
      video.addEventListener('durationchange', update);
      video.addEventListener('loadedmetadata', update);
      update();
    }
  }

  class MediaVolumeRange extends HTMLElement {
    connectedCallback() {
      if (this.shadowRoot) return;
      const root = this.attachShadow({ mode: 'open' });
      const input = document.createElement('input');
      input.type = 'range';
      input.min = '0';
      input.max = '1';
      input.step = '0.05';
      input.value = '1';
      input.part = 'range';
      const style = document.createElement('style');
      style.textContent = `
        :host { display: block; width: 82px; }
        input { width: 100%; accent-color: white; }
      `;
      root.append(style, input);
      const video = findVideo(this);
      if (!video) return;
      const update = () => { input.value = String(video.muted ? 0 : video.volume); };
      input.addEventListener('input', () => {
        video.volume = Number(input.value);
        video.muted = video.volume === 0;
      });
      video.addEventListener('volumechange', update);
      update();
    }
  }

  define('media-controller', MediaController);
  define('media-control-bar', MediaControlBar);
  define('media-play-button', MediaPlayButton);
  define('media-seek-backward-button', MediaSeekBackwardButton);
  define('media-seek-forward-button', MediaSeekForwardButton);
  define('media-time-range', MediaTimeRange);
  define('media-time-display', MediaTimeDisplay);
  define('media-duration-display', MediaDurationDisplay);
  define('media-mute-button', MediaMuteButton);
  define('media-volume-range', MediaVolumeRange);
  define('media-playback-rate-button', MediaPlaybackRateButton);
  define('media-fullscreen-button', MediaFullscreenButton);
})();
