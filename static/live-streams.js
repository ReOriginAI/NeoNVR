'use strict';

const videoEl = document.getElementById('video');
const placeholderEl = document.getElementById('placeholder');
const sourcesEl = document.getElementById('sources');
const watchEl = document.getElementById('watch');
const stopEl = document.getElementById('stop');
const cameraAudioEl = document.getElementById('cameraAudio');
const talkEl = document.getElementById('talk');
const microphoneEl = document.getElementById('microphone');
const statusEl = document.getElementById('status');

fetch('/api/https', {cache: 'no-store'}).then(response => response.ok ? response.json() : null).then(https => {
  if (https && https.url) {
    document.getElementById('httpsHint').hidden = false;
    const link = document.getElementById('httpsLink');
    link.href = https.url;
    if (location.protocol === 'https:') link.textContent = 'HTTPS is active';
  }
}).catch(() => {});

let state = null;
let selectedID = null;
let active = null;
let generation = 0;
let changing = false;
let stateVersion = 0;
let activation = null;
let listenToCamera = false;

function show(message, type = '') {
  statusEl.textContent = message;
  statusEl.className = 'status ' + type;
}

function selectedSource() {
  return state && state.streams.find(source => source.id === selectedID);
}

function updateControls() {
  const source = selectedSource();
  watchEl.disabled = changing || !source;
  stopEl.disabled = !active;
  const hasAudio = !!(active && active.connected && active.audioTrack && active.audioTrack.readyState === 'live');
  const listening = hasAudio && !videoEl.muted && videoEl.volume > 0;
  cameraAudioEl.disabled = !hasAudio;
  cameraAudioEl.textContent = active && active.connected && !hasAudio ? 'No camera audio' : listening ? 'Mute camera audio' : 'Unmute camera audio';
  cameraAudioEl.setAttribute('aria-pressed', String(listening));
  cameraAudioEl.title = active && active.connected && !hasAudio ? 'The camera did not provide a playable audio track.' : '';
  talkEl.disabled = changing || !active || !active.connected || !source || !source.talk || !window.isSecureContext || !navigator.mediaDevices;
  talkEl.textContent = active && active.talk ? 'Stop talking' : 'Talk to camera';
  talkEl.setAttribute('aria-pressed', String(!!(active && active.talk && active.connected)));
  microphoneEl.textContent = active && active.talk && active.connected ? 'Microphone live' : 'Microphone off';
  microphoneEl.className = 'mic' + (active && active.talk && active.connected ? ' active' : '');
  placeholderEl.hidden = !!(active && active.connected);
  placeholderEl.querySelector('strong').textContent = active ? 'Connecting…' : 'Live viewing is off';
  placeholderEl.querySelector('p').textContent = active ? 'Waiting for the camera.' : 'Select a camera and choose Watch, or double-click a camera.';
}

function stopPlayback(message = '') {
  generation++;
  const previous = active;
  active = null;
  if (previous) {
    clearTimeout(previous.timer);
    previous.abort.abort();
    if (previous.microphone) previous.microphone.getTracks().forEach(track => track.stop());
    previous.pc.close();
  }
  videoEl.pause();
  videoEl.srcObject = null;
  updateControls();
  if (message) show(message);
}

function renderSources() {
  sourcesEl.replaceChildren();
  if (!state.streams.length) {
    const empty = document.createElement('p');
    empty.textContent = 'Add an RTSP camera in Settings to get started.';
    sourcesEl.append(empty);
    return;
  }
  for (const source of state.streams) {
    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'source';
    button.setAttribute('aria-pressed', String(source.id === selectedID));
    button.title = 'Select ' + source.name + '; double-click to watch';
    button.append(document.createTextNode(source.name));
    const detail = document.createElement('span');
    detail.textContent = source.host + (source.talk ? ' · Talk available' : '');
    button.append(detail);
    const select = () => {
      if (selectedID === source.id) return;
      stopPlayback('Selected ' + source.name + '. Choose Watch or double-click to connect.');
      selectedID = source.id;
      // Keep these buttons in place so the first click does not interrupt dblclick.
      for (const candidate of sourcesEl.querySelectorAll('button.source')) {
        candidate.setAttribute('aria-pressed', String(candidate === button));
      }
      updateControls();
    };
    button.addEventListener('click', select);
    button.addEventListener('dblclick', () => {
      select();
      connect(false);
    });
    sourcesEl.append(button);
  }
}

async function loadState(initial = false) {
  const version = stateVersion;
  const response = await fetch('/api/live', {cache: 'no-store', signal: AbortSignal.timeout(5000)});
  if (!response.ok) throw new Error(await response.text());
  const next = await response.json();
  if (version !== stateVersion) return;
  const restarted = state && state.token !== next.token;
  const wasEnabled = state && state.enabled;
  state = next;
  if (selectedID === null || !selectedSource()) selectedID = state.streams.length ? state.streams[0].id : null;
  if (restarted || !state.enabled && (active || wasEnabled)) {
    stopPlayback(restarted ? 'Server restarted. Live viewing is off.' : 'Live viewing was turned off. All camera connections are closed.');
  }
  if (initial || restarted) {
    renderSources();
    show(state.streams.length ? 'Select a camera and choose Watch, or double-click it to connect.' : 'Add a camera using Edit sources.');
  }
  updateControls();
}

async function enableLiveViewing() {
  if (state.enabled) return;
  if (!activation) {
    changing = true;
    stateVersion++;
    updateControls();
    // Share startup across camera switches. Only the current viewer may connect
    // when startup completes; stopping during startup never opens a camera.
    activation = (async () => {
      const response = await fetch('/api/live', {
        method: 'POST', headers: {'Content-Type': 'application/json', 'X-NeoNVR-Token': state.token},
        body: JSON.stringify({enabled: true}), signal: AbortSignal.timeout(10000)
      });
      if (!response.ok) throw new Error(await response.text());
      state = await response.json();
      if (!state.enabled) throw new Error('Live viewing could not start. Choose Watch to retry.');
    })().finally(() => {
      activation = null;
      changing = false;
      updateControls();
    });
  }
  await activation;
}

function waitForICE(pc, signal) {
  if (pc.iceGatheringState === 'complete') return Promise.resolve();
  return new Promise((resolve, reject) => {
    const finish = error => {
      clearTimeout(timer);
      pc.removeEventListener('icegatheringstatechange', changed);
      signal.removeEventListener('abort', aborted);
      if (error) reject(error); else resolve();
    };
    const changed = () => { if (pc.iceGatheringState === 'complete') finish(); };
    const aborted = () => finish(new DOMException('Connection stopped', 'AbortError'));
    const timer = setTimeout(() => finish(), 5000);
    pc.addEventListener('icegatheringstatechange', changed);
    signal.addEventListener('abort', aborted, {once: true});
    if (signal.aborted) aborted();
  });
}

async function connect(talk = false) {
  const source = selectedSource();
  if (!source || document.hidden) return;
  if (!window.RTCPeerConnection) {
    show('This browser does not support live WebRTC playback.', 'error');
    return;
  }
  if (talk && (!window.isSecureContext || !navigator.mediaDevices)) {
    show('Use Open HTTPS for microphone access below the player, accept or trust the local certificate, and allow microphone access.', 'error');
    return;
  }
  stopPlayback();
  const currentGeneration = generation;
  const session = {
    pc: new RTCPeerConnection({iceServers: [], bundlePolicy: 'max-bundle'}),
    abort: new AbortController(), microphone: null, audioTrack: null, talk, connected: false, timer: null
  };
  active = session;
  show(talk ? 'Requesting microphone access…' : 'Connecting to ' + source.name + '…');
  updateControls();
  session.timer = setTimeout(() => {
    if (active === session) stopPlayback('Connection timed out. Check the camera and that the server WebRTC port is reachable.');
  }, 25000);
  const received = new MediaStream();
  videoEl.muted = !listenToCamera;
  videoEl.srcObject = received;
  session.pc.addEventListener('track', event => {
    if (active !== session) return;
    received.addTrack(event.track);
    if (event.track.kind === 'audio') {
      session.audioTrack = event.track;
      event.track.addEventListener('ended', () => { if (active === session) updateControls(); });
      updateControls();
    }
    videoEl.play().catch(() => show('Connected. Press Play in the video controls.', 'live'));
  });
  session.pc.addEventListener('connectionstatechange', () => {
    if (active !== session) return;
    const connection = session.pc.connectionState;
    if (connection === 'connected') {
      clearTimeout(session.timer);
      session.connected = true;
      if (session.microphone) session.microphone.getTracks().forEach(track => { track.enabled = true; });
      updateControls();
      show(talk ? 'Talking to ' + source.name + '. Choose Stop talking to turn your microphone off.' :
        'Watching ' + source.name + (session.audioTrack ? '.' : '. No playable camera audio was received.'), 'live');
    } else if (connection === 'failed' || connection === 'disconnected' || connection === 'closed') {
      stopPlayback('Camera disconnected. Choose Watch to reconnect.');
    }
  });
  try {
    await enableLiveViewing();
    if (active !== session) return;
    let sender = null;
    if (talk) {
      const microphone = await navigator.mediaDevices.getUserMedia({audio: {echoCancellation: true, noiseSuppression: true}});
      if (generation !== currentGeneration || active !== session) {
        microphone.getTracks().forEach(track => track.stop());
        return;
      }
      session.microphone = microphone;
      for (const track of microphone.getAudioTracks()) {
        track.enabled = false;
        sender = session.pc.addTransceiver(track, {direction: 'sendonly'});
      }
    }
    session.pc.addTransceiver('video', {direction: 'recvonly'});
    session.pc.addTransceiver('audio', {direction: 'recvonly'});
    await session.pc.setLocalDescription(await session.pc.createOffer());
    await waitForICE(session.pc, session.abort.signal);
    if (active !== session) return;
    show('Connecting to ' + source.name + (talk ? ' with microphone audio…' : '…'));
    const response = await fetch('/api/live/webrtc?id=' + source.id + (talk ? '&talk=1' : ''), {
      method: 'POST', headers: {'Content-Type': 'application/sdp', 'X-NeoNVR-Token': state.token},
      body: session.pc.localDescription.sdp, signal: session.abort.signal
    });
    if (!response.ok) throw new Error(await response.text());
    const answer = await response.text();
    if (active !== session) return;
    await session.pc.setRemoteDescription({type: 'answer', sdp: answer});
    if (talk && (!sender || !['sendonly', 'sendrecv'].includes(sender.currentDirection))) {
      throw new Error('The camera did not accept a compatible audio backchannel. Your microphone is off; choose Watch to return to viewing.');
    }
  } catch (error) {
    if (active !== session) return;
    stopPlayback();
    show(error.name === 'NotAllowedError' ? 'Microphone permission was denied. Choose Watch to view without talking.' : error.message, 'error');
  }
}

watchEl.addEventListener('click', () => connect(false));
stopEl.addEventListener('click', () => stopPlayback('Stopped. The camera connection and microphone are closed.'));
cameraAudioEl.addEventListener('click', () => {
  if (!active || !active.audioTrack) return;
  const unmute = videoEl.muted || videoEl.volume === 0;
  if (unmute && videoEl.volume === 0) videoEl.volume = 1;
  videoEl.muted = !unmute;
  listenToCamera = unmute;
  updateControls();
  if (unmute) videoEl.play().catch(() => show('Press Play in the video controls to hear camera audio.', 'live'));
});
videoEl.addEventListener('volumechange', () => {
  listenToCamera = !videoEl.muted && videoEl.volume > 0;
  updateControls();
});
talkEl.addEventListener('click', () => connect(!(active && active.talk)));
document.addEventListener('visibilitychange', () => {
  if (document.hidden && active) stopPlayback('Stopped while this tab was hidden. Choose Watch to reconnect.');
});
window.addEventListener('pagehide', () => stopPlayback());
loadState(true).catch(error => show(error.message, 'error'));
setInterval(() => {
  if (document.hidden || changing) return;
  loadState().catch(() => {
    if (active) stopPlayback('Server unavailable. The camera connection and microphone are closed.');
  });
}, 3000);
