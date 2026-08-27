// Svelte 5 Reactive DDP Streamer Store ($state runes)

let streams = $state({}); // map keyed by "targetType:targetID"

const isWails = typeof window !== 'undefined' && window.runtime !== undefined;

export function getDDPStore() {
  return {
    get streams() {
      return streams;
    },

    isStreaming(targetType, targetID) {
      const key = `${targetType}:${targetID}`;
      return !!streams[key]?.active;
    },

    getStreamStatus(targetType, targetID) {
      const key = `${targetType}:${targetID}`;
      return streams[key] || { active: false, fps: 0, effect: 'rainbow_wave', speed: 1.0, intensity: 1.0 };
    },

    async init() {
      await this.fetchStatus();
    },

    async fetchStatus() {
      try {
        if (isWails && window.go?.main?.App?.GetDDPStatus) {
          const res = await window.go.main.App.GetDDPStatus();
          if (res) streams = { ...res };
        } else {
          const res = await fetch('/api/v1/ddp/status');
          if (res.ok) {
            const data = await res.json();
            streams = { ...data };
          }
        }
      } catch (e) {
        console.error('[DDPStore] Error fetching status:', e);
      }
    },

    async startStream(targetType, targetID, effect, speed = 1.0, intensity = 1.0, ips = [], ledCount = 60) {
      try {
        if (isWails && window.go?.main?.App?.StartDDPStream) {
          const res = await window.go.main.App.StartDDPStream(
            targetType,
            targetID,
            effect,
            speed,
            intensity,
            ips,
            ledCount
          );
          if (res) {
            const key = `${targetType}:${targetID}`;
            streams = { ...streams, [key]: res };
          }
        } else {
          const res = await fetch('/api/v1/ddp/start', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ targetType, targetID, effect, speed, intensity, ips, ledCount })
          });
          if (res.ok) {
            const data = await res.json();
            const key = `${targetType}:${targetID}`;
            streams = { ...streams, [key]: data };
          }
        }
      } catch (e) {
        console.error('[DDPStore] Error starting stream:', e);
      }
    },

    async stopStream(targetType, targetID) {
      try {
        if (isWails && window.go?.main?.App?.StopDDPStream) {
          const res = await window.go.main.App.StopDDPStream(targetType, targetID);
          if (res) {
            const key = `${targetType}:${targetID}`;
            const copy = { ...streams };
            delete copy[key];
            streams = copy;
          }
        } else {
          const res = await fetch('/api/v1/ddp/stop', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ targetType, targetID })
          });
          if (res.ok) {
            const key = `${targetType}:${targetID}`;
            const copy = { ...streams };
            delete copy[key];
            streams = copy;
          }
        }
      } catch (e) {
        console.error('[DDPStore] Error stopping stream:', e);
      }
    },

    handleWSMessage(msg) {
      if (msg && msg.type === 'ddp_status' && msg.data) {
        streams = { ...msg.data };
      }
    }
  };
}
