class AudioStore {
  isCapturing = $state(false);
  fftBins = $state(new Array(16).fill(0));
  gain = $state(1.0);
  band = $state('full'); // 'bass' | 'mid' | 'treble' | 'full'
  selectedSource = $state('default');

  audioContext = null;
  analyserNode = null;
  mediaStream = null;
  animationFrameId = null;

  bassEnergy = $derived.by(() => {
    // Bins 0-3 (Bass 20Hz-250Hz)
    const sum = this.fftBins.slice(0, 4).reduce((a, b) => a + b, 0);
    return Math.min(1.0, (sum / 4) * this.gain);
  });

  midEnergy = $derived.by(() => {
    // Bins 4-9 (Mid 250Hz-4kHz)
    const sum = this.fftBins.slice(4, 10).reduce((a, b) => a + b, 0);
    return Math.min(1.0, (sum / 6) * this.gain);
  });

  trebleEnergy = $derived.by(() => {
    // Bins 10-15 (Treble 4kHz-20kHz)
    const sum = this.fftBins.slice(10, 16).reduce((a, b) => a + b, 0);
    return Math.min(1.0, (sum / 6) * this.gain);
  });

  peakEnergy = $derived.by(() => {
    switch (this.band) {
      case 'bass':
        return this.bassEnergy;
      case 'mid':
        return this.midEnergy;
      case 'treble':
        return this.trebleEnergy;
      case 'full':
      default:
        return Math.max(this.bassEnergy, this.midEnergy, this.trebleEnergy);
    }
  });

  setGain(val) {
    this.gain = Math.max(0.1, Math.min(5.0, Number(val)));
  }

  setBand(bandName) {
    if (['bass', 'mid', 'treble', 'full'].includes(bandName)) {
      this.band = bandName;
    }
  }

  updateFFTBins(bins) {
    if (Array.isArray(bins) && bins.length === 16) {
      this.fftBins = bins.map((v) => Math.max(0.0, Math.min(1.0, Number(v))));
    }
  }

  async startCapture() {
    if (typeof window === 'undefined' || !navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      console.warn('[AudioStore] Audio mediaDevices API not available in environment.');
      return;
    }

    try {
      this.mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false });
      const AudioCtx = window.AudioContext || window.webkitAudioContext;
      this.audioContext = new AudioCtx();
      const source = this.audioContext.createMediaStreamSource(this.mediaStream);

      this.analyserNode = this.audioContext.createAnalyser();
      this.analyserNode.fftSize = 64; // Gives 32 frequency bins
      this.analyserNode.smoothingTimeConstant = 0.8;
      source.connect(this.analyserNode);

      this.isCapturing = true;
      this.processAudioLoop();
    } catch (err) {
      console.error('[AudioStore] Failed to start audio capture:', err);
      this.isCapturing = false;
    }
  }

  processAudioLoop = () => {
    if (!this.isCapturing || !this.analyserNode) return;

    const bufferLength = this.analyserNode.frequencyBinCount;
    const dataArray = new Uint8Array(bufferLength);
    this.analyserNode.getByteFrequencyData(dataArray);

    // Downsample to 16 normalized bins (0.0 - 1.0)
    const bins = new Array(16).fill(0);
    const step = Math.max(1, Math.floor(bufferLength / 16));

    for (let i = 0; i < 16; i++) {
      let sum = 0;
      let count = 0;
      for (let j = 0; j < step && i * step + j < bufferLength; j++) {
        sum += dataArray[i * step + j];
        count++;
      }
      bins[i] = count > 0 ? sum / count / 255.0 : 0.0;
    }

    this.updateFFTBins(bins);

    // Send audio_fft message over WebSocket if WS connection exists
    if (typeof window !== 'undefined' && window.__ledSwarmWS && window.__ledSwarmWS.readyState === 1) {
      try {
        window.__ledSwarmWS.send(
          JSON.stringify({
            type: 'audio_fft',
            bins: this.fftBins,
            bass: this.bassEnergy,
            mid: this.midEnergy,
            treble: this.trebleEnergy,
            peak: this.peakEnergy
          })
        );
      } catch (_) {}
    }

    this.animationFrameId = requestAnimationFrame(this.processAudioLoop);
  };

  stopCapture() {
    this.isCapturing = false;
    if (this.animationFrameId) {
      cancelAnimationFrame(this.animationFrameId);
      this.animationFrameId = null;
    }

    if (this.mediaStream) {
      this.mediaStream.getTracks().forEach((track) => track.stop());
      this.mediaStream = null;
    }

    if (this.audioContext) {
      this.audioContext.close();
      this.audioContext = null;
    }

    this.fftBins = new Array(16).fill(0);
  }
}

let audioStoreInstance = null;

export function getAudioStore() {
  if (!audioStoreInstance) {
    audioStoreInstance = new AudioStore();
  }
  return audioStoreInstance;
}
