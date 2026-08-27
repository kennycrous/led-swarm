<script>
  import {
    LayoutGrid,
    Power,
    Sun,
    Palette,
    Sparkles,
    Pin,
    Maximize2,
    MoreVertical,
    Trash2,
    Edit3,
    Zap,
    Volume2
  } from 'lucide-svelte';
  import CyberSelect from './CyberSelect.svelte';

  let {
    room = { id: 'default', title: 'Main Room Canvas', width: 2000, height: 1200 },
    devices = [],
    placements = [],
    effects = [],
    palettes = [],
    isPinned = true,
    cardSize = 'normal',
    showSizeToggle = false,
    ddpStore,
    audioStore,
    onEditLayout = () => {},
    onRename = () => {},
    onTogglePin = () => {},
    onToggleSize = () => {},
    onDelete = null
  } = $props();

  let isMenuOpen = $state(false);
  let isEditingTitle = $state(false);
  let editedTitle = $state('');
  let roomPower = $state(true);
  let roomBrightness = $state(200);
  let selectedColor = $state('#06b6d4');
  let selectedEffect = $state(0);
  let selectedPalette = $state(0);

  let isDDPActive = $derived(ddpStore?.isStreaming('room', room.id) || false);
  let ddpStatus = $derived(ddpStore?.getStreamStatus('room', room.id) || {});

  let cardMode = $derived.by(() => {
    if (!isDDPActive) return 'native';
    if (ddpStatus.effect && ddpStatus.effect.startsWith('audio_')) return 'music';
    return 'ddp';
  });

  const ddpSpatialEffects = [
    { id: 'spatial_ripple', name: '🎯 2D DDP: Radial Room Ripple' },
    { id: 'spatial_sweep', name: '↔️ 2D DDP: Directional Room Sweep' },
    { id: 'rainbow_wave', name: '⚡ DDP: Rainbow Wave' },
    { id: 'digital_rain', name: '⚡ DDP: Matrix Rain' },
    { id: 'pulse_beads', name: '⚡ DDP: Neon Pulse' },
    { id: 'cyber_fire', name: '⚡ DDP: Cyber Flame' }
  ];

  const audioSpatialEffects = [
    { id: 'audio_bass_pulse', name: '🎵 Sound: 2D Bass Kick Pulse' },
    { id: 'audio_spectrum_waterfall', name: '🎵 Sound: 2D 16-Bin Spectrum' },
    { id: 'audio_vu_meter', name: '🎵 Sound: 2D Stereo VU Meter' },
    { id: 'audio_treble_sparkle', name: '🎵 Sound: 2D Treble Sparkles' }
  ];

  function handleSaveTitle() {
    if (editedTitle.trim() && editedTitle.trim() !== room.title) {
      onRename(room.id, editedTitle.trim());
    }
    isEditingTitle = false;
  }

  const presetColors = [
    { name: 'Cyan Neon', hex: '#06b6d4', r: 6, g: 182, b: 212 },
    { name: 'Magenta Glow', hex: '#a855f7', r: 168, g: 85, b: 247 },
    { name: 'Cyber Amber', hex: '#f59e0b', r: 245, g: 158, b: 11 },
    { name: 'Void White', hex: '#ffffff', r: 255, g: 255, b: 255 },
    { name: 'Emerald Online', hex: '#10b981', r: 16, g: 185, b: 129 },
    { name: 'Crimson Red', hex: '#ef4444', r: 239, g: 68, b: 68 }
  ];

  // Filter placements for this room
  let roomPlacements = $derived(placements.filter((p) => p.roomId === room.id));

  // Compute auto-centered & scaled preview coordinates for room placements
  let previewPlacements = $derived.by(() => {
    if (roomPlacements.length === 0) return [];

    let minX = Infinity,
      maxX = -Infinity;
    let minY = Infinity,
      maxY = -Infinity;

    roomPlacements.forEach((p) => {
      if (p.posX < minX) minX = p.posX;
      if (p.posX > maxX) maxX = p.posX;
      if (p.posY < minY) minY = p.posY;
      if (p.posY > maxY) maxY = p.posY;
    });

    const centerX = (minX + maxX) / 2;
    const centerY = (minY + maxY) / 2;
    const spanX = maxX - minX;
    const spanY = maxY - minY;

    const rangeX = Math.max(spanX * 1.6, 600);
    const rangeY = Math.max(spanY * 1.6, 360);

    return roomPlacements.map((p) => {
      const dev = devices.find((d) => d.id === p.deviceId);
      const relX = 50 + ((p.posX - centerX) / rangeX) * 100;
      const relY = 50 + ((p.posY - centerY) / rangeY) * 100;
      const posX = Math.max(12, Math.min(88, relX));
      const posY = Math.max(15, Math.min(85, relY));

      const rot = p.rotation || 0;
      const isOn = dev?.state?.on ?? false;
      const color = dev?.state?.seg?.[0]?.col?.[0]
        ? `rgb(${dev.state.seg[0].col[0][0]}, ${dev.state.seg[0].col[0][1]}, ${dev.state.seg[0].col[0][2]})`
        : '#06b6d4';

      return {
        ...p,
        dev,
        posX,
        posY,
        rot,
        isOn,
        color
      };
    });
  });

  // Devices in this room
  let roomDevices = $derived(devices.filter((dev) => roomPlacements.some((p) => p.deviceId === dev.id)));
  let onlineCount = $derived(roomDevices.filter((d) => d.isOnline).length);

  let isAnyOn = $derived(roomDevices.some((d) => d.state?.on));

  function toggleRoomPower() {
    const nextOn = !isAnyOn;
    roomPower = nextOn;
    roomDevices.forEach((dev) => {
      if (dev && dev.isOnline) {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: dev.ipAddress,
            state: { on: nextOn }
          })
        }).catch(() => {});
      }
    });
  }

  function handleBriChange(e) {
    const val = Number(e.target.value);
    roomBrightness = val;
    roomDevices.forEach((dev) => {
      if (dev && dev.isOnline) {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: dev.ipAddress,
            state: { bri: val }
          })
        }).catch(() => {});
      }
    });
  }

  function applyColor(c) {
    selectedColor = c.hex;
    roomDevices.forEach((dev) => {
      if (dev && dev.isOnline) {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: dev.ipAddress,
            state: { seg: [{ id: 0, col: [[c.r, c.g, c.b]] }] }
          })
        }).catch(() => {});
      }
    });
  }

  function handleEffectChange(fxId) {
    selectedEffect = fxId;
    roomDevices.forEach((dev) => {
      if (dev && dev.isOnline) {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: dev.ipAddress,
            state: { seg: [{ id: 0, fx: fxId }] }
          })
        }).catch(() => {});
      }
    });
  }

  function handlePaletteChange(palId) {
    selectedPalette = palId;
    if (selectedEffect === 0) selectedEffect = 2;
    roomDevices.forEach((dev) => {
      if (dev && dev.isOnline) {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: dev.ipAddress,
            state: { seg: [{ id: 0, pal: palId, fx: selectedEffect || 2 }] }
          })
        }).catch(() => {});
      }
    });
  }
</script>

<div
  class="glass-panel rounded-2xl p-4 flex flex-col justify-between space-y-3 relative group border h-full transition-all duration-300 {isMenuOpen
    ? 'z-50'
    : 'hover:z-30 z-10'} border-cyan-500/20 hover:border-cyan-500/40"
>
  <!-- Top Accent Status Pill Bar -->
  <div class="w-full h-1 rounded-full bg-gradient-to-r from-purple-500 via-cyan-500 to-amber-500"></div>

  <!-- FUNCTIONAL ROOM HEADER & OPTIONS MENU -->
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2 flex-1">
      <div class="p-2 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-400">
        <LayoutGrid class="w-4 h-4" />
      </div>
      <div>
        {#if isEditingTitle}
          <form
            onsubmit={(e) => {
              e.preventDefault();
              handleSaveTitle();
            }}
            class="flex items-center gap-1"
          >
            <input
              type="text"
              bind:value={editedTitle}
              onblur={handleSaveTitle}
              class="bg-[#06090e] border border-cyan-500/50 rounded-lg px-2 py-0.5 text-xs font-semibold text-slate-100 focus:outline-none focus:border-cyan-400 w-36"
            />
          </form>
        {:else}
          <button
            type="button"
            onclick={() => {
              isEditingTitle = true;
              editedTitle = room.title;
            }}
            class="group/title flex items-center gap-1.5 cursor-pointer text-left border-0 bg-transparent p-0"
            title="Click to rename room"
          >
            <h3 class="font-semibold text-slate-100 text-sm tracking-wide">{room.title}</h3>
            <Edit3 class="w-3 h-3 text-slate-500 opacity-0 group-hover/title:opacity-100 transition-opacity" />
          </button>
        {/if}
        <p class="text-[11px] font-mono text-slate-400">
          {roomDevices.length} Strips ({onlineCount} Online)
        </p>
      </div>
    </div>

    <!-- Power & Card Options Menu -->
    <div class="flex items-center gap-1.5 relative">
      <button
        onclick={toggleRoomPower}
        class="p-2.5 rounded-xl transition-all duration-200 {isAnyOn
          ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan'
          : 'bg-[#090e17]/80 text-slate-600 border border-slate-800'} cursor-pointer"
        title="Toggle Room Power"
      >
        <Power class="w-4 h-4" />
      </button>

      <!-- Card Options Dropdown Trigger -->
      <button
        onclick={() => (isMenuOpen = !isMenuOpen)}
        class="p-2.5 rounded-xl bg-[#090e17]/80 hover:bg-slate-800 text-slate-400 hover:text-slate-200 border border-slate-800 transition-all cursor-pointer"
        title="Card Options"
      >
        <MoreVertical class="w-4 h-4" />
      </button>

      <!-- Glassmorphic Cyber Context Dropdown Menu -->
      {#if isMenuOpen}
        <div
          class="absolute top-12 right-0 z-50 w-48 bg-[#090e17]/95 border border-cyan-500/30 rounded-2xl p-1.5 shadow-[0_0_25px_rgba(6,182,212,0.25)] backdrop-blur-xl space-y-1 text-xs font-mono"
        >
          <button
            onclick={() => {
              isEditingTitle = true;
              editedTitle = room.title;
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-cyan-300 hover:bg-cyan-500/15 transition-all text-left cursor-pointer"
          >
            <Edit3 class="w-3.5 h-3.5 text-cyan-400" />
            <span>Rename Room</span>
          </button>

          <button
            onclick={() => {
              onEditLayout(room.id);
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-cyan-300 hover:bg-cyan-500/15 transition-all text-left cursor-pointer"
          >
            <LayoutGrid class="w-3.5 h-3.5 text-cyan-400" />
            <span>Edit 2D Layout</span>
          </button>

          {#if showSizeToggle}
            <button
              onclick={() => {
                onToggleSize(room.id);
                isMenuOpen = false;
              }}
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-cyan-300 hover:bg-cyan-500/15 transition-all text-left cursor-pointer"
            >
              <Maximize2 class="w-3.5 h-3.5 text-cyan-400" />
              <span>{cardSize === 'wide' ? 'Normal Width' : 'Expand Full Width'}</span>
            </button>
          {/if}

          <button
            onclick={() => {
              onTogglePin(room.id, 'room');
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-cyan-300 hover:bg-cyan-500/15 transition-all text-left cursor-pointer"
          >
            <Pin class="w-3.5 h-3.5 text-cyan-400" />
            <span>{isPinned ? 'Unpin from Canvas' : 'Pin to Dashboard'}</span>
          </button>

          {#if onDelete}
            <button
              onclick={() => {
                onDelete(room.id);
                isMenuOpen = false;
              }}
              class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-rose-400 hover:bg-rose-500/15 transition-all text-left cursor-pointer border-t border-slate-800/80 mt-1 pt-1.5"
            >
              <Trash2 class="w-3.5 h-3.5 text-rose-400" />
              <span>Delete Room</span>
            </button>
          {/if}
        </div>
      {/if}
    </div>
  </div>

  <!-- Mode Segmented Switch Row (Native | DDP 60 | Music) -->
  <div
    class="grid grid-cols-3 gap-1 p-1 rounded-2xl bg-[#06090e] border border-slate-800/80 font-mono text-[11px] font-bold"
  >
    <button
      type="button"
      onclick={() => {
        if (isDDPActive) ddpStore?.stopStream('room', room.id);
      }}
      class="py-1.5 rounded-xl transition-all text-center cursor-pointer {cardMode === 'native'
        ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/40 glow-cyan'
        : 'text-slate-500 hover:text-slate-300'}"
      title="Native WLED Mode (Standard WLED FX & Palettes)"
    >
      Native
    </button>

    <button
      type="button"
      onclick={() => {
        ddpStore?.startStream('room', room.id, 'spatial_ripple', 1.0, 1.0);
      }}
      class="py-1.5 rounded-xl transition-all text-center cursor-pointer flex items-center justify-center gap-1 {cardMode ===
      'ddp'
        ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40 glow-amber'
        : 'text-slate-500 hover:text-amber-400'}"
      title="DDP 60 FPS 2D Spatial Matrix Mode"
    >
      <Zap class="w-3.5 h-3.5 {cardMode === 'ddp' ? 'animate-pulse text-amber-400' : ''}" />
      <span>{cardMode === 'ddp' ? `${ddpStatus.fps || 60} FPS` : 'DDP 60'}</span>
    </button>

    <button
      type="button"
      onclick={() => {
        if (audioStore && !audioStore.isCapturing) audioStore.startCapture();
        ddpStore?.startStream('room', room.id, 'audio_bass_pulse', 1.0, 1.0);
      }}
      class="py-1.5 rounded-xl transition-all text-center cursor-pointer flex items-center justify-center gap-1 {cardMode ===
      'music'
        ? 'bg-purple-500/20 text-purple-300 border border-purple-500/40 glow-purple'
        : 'text-slate-500 hover:text-purple-300'}"
      title="Sound Sync Music Mode (FFT Audio Reactivity)"
    >
      <Volume2 class="w-3.5 h-3.5 {cardMode === 'music' ? 'animate-pulse text-purple-400' : ''}" />
      <span>Music</span>
    </button>
  </div>

  <!-- Mini 2D Room Canvas Preview Grid -->
  <div class="relative my-2 h-28 w-full overflow-hidden rounded-xl border border-slate-800 bg-[#06090e] p-2">
    <!-- Background Grid Lines -->
    <div
      class="absolute inset-0 bg-[linear-gradient(to_right,rgba(56,189,248,0.05)_1px,transparent_1px),linear-gradient(to_bottom,rgba(56,189,248,0.05)_1px,transparent_1px)] bg-[size:16px_16px]"
    ></div>

    <!-- Mini Placements Preview -->
    {#each previewPlacements as p (p.deviceId)}
      <div
        class="absolute flex items-center justify-center rounded-full transition-all duration-300"
        style="left: {p.posX}%; top: {p.posY}%; transform: translate(-50%, -50%) rotate({p.rot}deg);"
      >
        <!-- Strip Line Indicator -->
        <div
          class="h-2 w-10 rounded-full border border-cyan-400/40 shadow-sm"
          style="background-color: {p.isOn ? p.color : '#334155'}; box-shadow: {p.isOn
            ? `0 0 8px ${p.color}`
            : 'none'};"
        ></div>
      </div>
    {/each}
  </div>

  <!-- Mode, Color, Brightness & Palette Controls -->
  <div class="space-y-3 pt-1">
    <!-- Quick Color Swatches & Indicator -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-1.5">
        {#each presetColors as c (c.hex)}
          <button
            onclick={() => applyColor(c)}
            style="background-color: {c.hex}"
            class="w-5 h-5 rounded-md border border-slate-900 shadow-sm transition-transform hover:scale-110 focus:outline-none cursor-pointer"
            title={c.name}
          ></button>
        {/each}
      </div>

      <!-- Color Indicator -->
      <div
        style="background-color: {selectedColor}"
        class="w-6 h-6 rounded-lg border border-cyan-500/40 shadow-neonCyan"
      ></div>
    </div>

    <!-- Room Brightness Slider -->
    <div class="flex items-center gap-3 bg-[#090e17]/60 px-3 py-1.5 rounded-xl border border-slate-800/80 h-9">
      <Sun class="w-4 h-4 text-cyan-400" />
      <input
        type="range"
        min="0"
        max="255"
        value={roomBrightness}
        oninput={handleBriChange}
        disabled={!isAnyOn}
        class="flex-1 accent-cyan-400 cursor-pointer disabled:opacity-30"
      />
      <span class="text-xs font-mono text-cyan-300 w-8 text-right">
        {roomBrightness}
      </span>
    </div>

    <!-- Mode-Specific Controls Drawer & Dropdowns -->
    {#if cardMode === 'music'}
      <div class="space-y-2 bg-purple-500/10 p-3 rounded-2xl border border-purple-500/30 font-mono">
        <div class="flex items-center justify-between text-[11px] text-purple-300 font-bold">
          <span class="flex items-center gap-1.5">
            <Volume2 class="w-3.5 h-3.5 text-purple-400 animate-pulse" />
            2D Room Sound Sync Audio Reactivity
          </span>
          <span class="text-[10px] text-purple-400/80">{audioStore?.isCapturing ? 'MIC LIVE' : 'OFFLINE'}</span>
        </div>

        <!-- 16-Band Mini FFT Visualizer Spectrum Bar -->
        <div class="flex items-end gap-1 h-8 bg-[#06090e] p-1 rounded-xl border border-purple-500/20 overflow-hidden">
          {#each audioStore?.fftBins || new Array(16).fill(0) as binVal, idx (idx)}
            <div
              class="flex-1 rounded-t-sm transition-all duration-75 {idx < 4
                ? 'bg-purple-500'
                : idx < 10
                  ? 'bg-cyan-400'
                  : 'bg-amber-400'}"
              style="height: {Math.max(10, Math.min(100, binVal * (audioStore?.gain || 1.0) * 100))}%"
            ></div>
          {/each}
        </div>

        <!-- Gain Slider & Band Selector -->
        <div class="grid grid-cols-2 gap-2 text-[10px]">
          <div class="space-y-1">
            <div class="flex justify-between text-slate-400">
              <span>Gain:</span>
              <span class="text-purple-300 font-bold">{audioStore?.gain || 1.0}x</span>
            </div>
            <input
              type="range"
              min="0.5"
              max="4.0"
              step="0.1"
              value={audioStore?.gain || 1.0}
              oninput={(e) => audioStore?.setGain(e.target.value)}
              class="w-full accent-purple-400 cursor-pointer"
            />
          </div>

          <div class="space-y-1">
            <div class="text-slate-400">Target Band:</div>
            <div class="grid grid-cols-4 gap-0.5">
              {#each ['bass', 'mid', 'treble', 'full'] as b (b)}
                <button
                  type="button"
                  onclick={() => audioStore?.setBand(b)}
                  class="py-0.5 rounded text-[9px] font-bold uppercase transition-all cursor-pointer {audioStore?.band ===
                  b
                    ? 'bg-purple-500 text-white'
                    : 'bg-[#06090e] text-slate-400 border border-slate-800'}"
                >
                  {b.slice(0, 3)}
                </button>
              {/each}
            </div>
          </div>
        </div>

        <!-- Audio Reactive Preset Picker -->
        <select
          value={ddpStatus.effect || 'audio_bass_pulse'}
          onchange={(e) => {
            ddpStore?.startStream('room', room.id, e.target.value, ddpStatus.speed || 1.0, ddpStatus.intensity || 1.0);
          }}
          class="w-full bg-[#06090e] border border-purple-500/40 rounded-xl px-2.5 py-1 text-xs font-mono text-purple-300 focus:outline-none cursor-pointer"
        >
          {#each audioSpatialEffects as eff (eff.id)}
            <option value={eff.id}>{eff.name}</option>
          {/each}
        </select>
      </div>
    {:else if cardMode === 'ddp'}
      <div class="space-y-2 bg-amber-500/10 p-2.5 rounded-2xl border border-amber-500/30">
        <div class="flex items-center justify-between text-[11px] font-mono text-amber-400 font-bold">
          <span class="flex items-center gap-1">
            <Zap class="w-3 h-3 text-amber-400 animate-pulse" />
            2D Room Canvas Spatial DDP Stream
          </span>
          <span>{ddpStatus.fps || 60} FPS</span>
        </div>
        <select
          value={ddpStatus.effect || 'spatial_ripple'}
          onchange={(e) => {
            ddpStore?.startStream('room', room.id, e.target.value, ddpStatus.speed || 1.0, ddpStatus.intensity || 1.0);
          }}
          class="w-full bg-[#06090e] border border-amber-500/40 rounded-xl px-2.5 py-1 text-xs font-mono text-amber-300 focus:outline-none cursor-pointer"
        >
          {#each ddpSpatialEffects as eff (eff.id)}
            <option value={eff.id}>{eff.name}</option>
          {/each}
        </select>
      </div>
    {:else}
      <div class="grid grid-cols-2 gap-2">
        <!-- FX CyberSelect -->
        <CyberSelect
          value={selectedEffect}
          options={effects}
          icon={Sparkles}
          iconColor="text-cyan-400"
          hoverBorder="hover:border-cyan-500/40"
          onChange={handleEffectChange}
        />

        <!-- Palette CyberSelect -->
        <CyberSelect
          value={selectedPalette}
          options={palettes}
          icon={Palette}
          iconColor="text-purple-400"
          hoverBorder="hover:border-purple-500/40"
          onChange={handlePaletteChange}
        />
      </div>
    {/if}

    <!-- Action Controls Footer -->
    <div class="flex items-center justify-between gap-2 pt-2 border-t border-slate-800/60">
      <span class="text-[11px] font-mono text-slate-400">2D Room Canvas</span>
      <button
        type="button"
        onclick={() => onEditLayout(room.id)}
        class="flex items-center gap-1.5 rounded-xl border border-cyan-500/30 bg-cyan-950/40 px-3 py-1.5 font-mono text-xs font-semibold text-cyan-300 transition-all duration-200 hover:border-cyan-400 hover:bg-cyan-900/50 hover:text-cyan-200 cursor-pointer shadow-[0_0_10px_rgba(6,182,212,0.15)]"
      >
        <LayoutGrid class="h-3.5 w-3.5" />
        <span>Edit 2D Layout</span>
      </button>
    </div>
  </div>
</div>
