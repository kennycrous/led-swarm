// Svelte 5 Reactive Device Store ($state runes)

let devices = $state([]);
let effects = $state([]);
let palettes = $state([]);
let isScanning = $state(false);
let isConnected = $state(false);
let ws = null;
let brightnessTimers = {};

const isWails = typeof window !== 'undefined' && window.runtime !== undefined;

export function getDeviceStore() {
  return {
    get devices() {
      return devices;
    },
    get onlineDevices() {
      return (devices || []).filter((d) => d && d.isOnline);
    },
    get effects() {
      return effects;
    },
    get palettes() {
      return palettes;
    },
    get isScanning() {
      return isScanning;
    },
    get isConnected() {
      return isConnected;
    },

    async init() {
      await loadDevices();
      await loadMetadata();

      if (!isWails) {
        connectWebSocket();
      }
    },

    async triggerScan() {
      isScanning = true;
      try {
        if (isWails && window.go?.main?.App?.TriggerScan) {
          await window.go.main.App.TriggerScan();
        } else {
          await fetch('/api/v1/scan', { method: 'POST' });
        }
      } catch (e) {
        console.error('Scan error:', e);
      }
      setTimeout(() => {
        isScanning = false;
      }, 4000);
    },

    async togglePower(deviceId) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;
      const targetState = !dev.state?.on;
      const targetBri = dev.state?.bri && dev.state.bri > 0 ? dev.state.bri : 128;
      dev.state = { ...dev.state, on: targetState, bri: targetBri };

      const statePayload = targetState
        ? { on: true, bri: targetBri, mainseg: 0, seg: [{ id: 0, on: true }] }
        : { on: false };

      try {
        if (isWails && window.go?.main?.App?.SetDevicePower) {
          await window.go.main.App.SetDevicePower(dev.ipAddress, targetState);
        } else {
          await fetch('/api/v1/devices/state', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ip: dev.ipAddress, state: statePayload })
          });
        }
      } catch (e) {
        console.error('Failed to toggle power:', e);
      }
    },

    async setBrightness(deviceId, brightness) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;
      dev.state = { ...dev.state, bri: brightness };

      if (brightnessTimers[deviceId]) {
        clearTimeout(brightnessTimers[deviceId]);
      }

      brightnessTimers[deviceId] = setTimeout(async () => {
        try {
          if (isWails && window.go?.main?.App?.SetDeviceBrightness) {
            await window.go.main.App.SetDeviceBrightness(dev.ipAddress, brightness);
          } else {
            await fetch('/api/v1/devices/state', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ ip: dev.ipAddress, state: { bri: brightness } })
            });
          }
        } catch (e) {
          console.error('Failed to set brightness:', e);
        }
      }, 75);
    },

    async setColor(deviceId, r, g, b) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;

      try {
        if (isWails && window.go?.main?.App?.SetDeviceColor) {
          await window.go.main.App.SetDeviceColor(dev.ipAddress, r, g, b);
        } else {
          await fetch('/api/v1/devices/state', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              ip: dev.ipAddress,
              state: { seg: [{ id: 0, col: [[r, g, b]] }] }
            })
          });
        }
      } catch (e) {
        console.error('Failed to set color:', e);
      }
    },

    async setEffect(deviceId, effectFx) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;

      if (!dev.state) dev.state = {};
      if (!dev.state.seg || dev.state.seg.length === 0) {
        dev.state.seg = [{ id: 0, fx: effectFx }];
      } else {
        dev.state.seg[0].fx = effectFx;
      }
      devices = [...devices];

      try {
        if (isWails && window.go?.main?.App?.SetDeviceEffect) {
          await window.go.main.App.SetDeviceEffect(dev.ipAddress, effectFx);
        } else {
          await fetch('/api/v1/devices/state', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              ip: dev.ipAddress,
              state: { seg: [{ id: 0, fx: effectFx }] }
            })
          });
        }
      } catch (e) {
        console.error('Failed to set effect:', e);
      }
    },

    async setPalette(deviceId, paletteId) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;

      if (!dev.state) dev.state = {};
      if (!dev.state.seg || dev.state.seg.length === 0) {
        dev.state.seg = [{ id: 0, pal: paletteId }];
      } else {
        dev.state.seg[0].pal = paletteId;
      }
      devices = [...devices];

      try {
        if (isWails && window.go?.main?.App?.SetDevicePalette) {
          await window.go.main.App.SetDevicePalette(dev.ipAddress, paletteId);
        } else {
          await fetch('/api/v1/devices/state', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              ip: dev.ipAddress,
              state: { seg: [{ id: 0, pal: paletteId }] }
            })
          });
        }
      } catch (e) {
        console.error('Failed to set palette:', e);
      }
    },

    async renameDevice(deviceId, newName) {
      const dev = devices.find((d) => d.id === deviceId);
      if (!dev) return;
      dev.name = newName;

      try {
        if (isWails && window.go?.main?.App?.UpdateDeviceNickname) {
          await window.go.main.App.UpdateDeviceNickname(deviceId, newName);
        } else {
          await fetch('/api/v1/devices/name', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: deviceId, name: newName })
          });
        }
      } catch (e) {
        console.error('Failed to rename device:', e);
      }
    },

    async addManualIP(ip) {
      try {
        let newDev = null;
        if (isWails && window.go?.main?.App?.AddManualDevice) {
          newDev = await window.go.main.App.AddManualDevice(ip);
        } else {
          const res = await fetch('/api/v1/devices/add', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ip })
          });
          if (res.ok) {
            newDev = await res.json();
          } else {
            const errTxt = await res.text();
            throw new Error(errTxt || 'Failed to connect to WLED device');
          }
        }
        if (newDev) {
          upsertDevice(newDev);
        }
      } catch (e) {
        console.error('Failed to add device by IP:', e);
        throw e;
      }
    }
  };
}

async function loadDevices() {
  try {
    let list = [];
    if (isWails && window.go?.main?.App?.GetDevices) {
      list = await window.go.main.App.GetDevices();
    } else {
      const res = await fetch('/api/v1/devices');
      if (res.ok) list = await res.json();
    }
    if (list && list.length > 0) {
      devices = list;
    }
  } catch (e) {
    console.log('Using initial state store');
  }
}

async function loadMetadata() {
  try {
    const [effRes, palRes] = await Promise.all([fetch('/api/v1/effects'), fetch('/api/v1/palettes')]);
    if (effRes.ok) effects = await effRes.json();
    if (palRes.ok) palettes = await palRes.json();
  } catch (e) {
    console.log('Metadata load fallback');
  }
}

function connectWebSocket() {
  const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
  const wsURL = `${protocol}//${location.host}/api/v1/ws`;

  try {
    ws = new WebSocket(wsURL);

    ws.onopen = () => {
      isConnected = true;
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === 'device_state' || msg.type === 'device_discovered') {
          if (msg.device) upsertDevice(msg.device);
        }
      } catch (e) {
        console.error('Error parsing WS message:', e);
      }
    };

    ws.onclose = () => {
      isConnected = false;
      setTimeout(connectWebSocket, 3000);
    };
  } catch (e) {
    console.error('WS init error:', e);
  }
}

function upsertDevice(newDev) {
  const idx = devices.findIndex((d) => d.id === newDev.id);
  if (idx >= 0) {
    devices[idx] = newDev;
  } else {
    devices = [...devices, newDev];
  }
}
