<script>
  import {
    Power,
    Edit2,
    Check,
    Palette,
    Sparkles,
    Sun,
    Pin,
    Maximize2,
    Cpu,
    MoreVertical,
    Zap,
    Volume2
  } from 'lucide-svelte';
  import CyberSelect from './CyberSelect.svelte';

  let {
    device,
    effects = [],
    palettes = [],
    isPinned = true,
    cardSize = 'normal',
    showSizeToggle = false,
    ddpStore,
    audioStore,
    onTogglePower = () => {},
    onSetBrightness = () => {},
    onSetColor = () => {},
    onSetEffect = () => {},
    onSetPalette = () => {},
    onRename = () => {},
    onTogglePin = () => {},
    onToggleSize = () => {}
  } = $props();

  let isEditingName = $state(false);
  let nickname = $state('');
  let selectedColor = $state('#06b6d4');
  let isMenuOpen = $state(false);

  let isDDPActive = $derived(ddpStore?.isStreaming('device', device.id) || false);
  let ddpStatus = $derived(ddpStore?.getStreamStatus('device', device.id) || {});

  let cardMode = $derived.by(() => {
    if (!isDDPActive) return 'native';
    if (ddpStatus.effect && ddpStatus.effect.startsWith('audio_')) return 'music';
    return 'ddp';
  });

  const ddpEffects = [
    { id: 'rainbow_wave', name: '⚡ DDP: Rainbow Wave' },
    { id: 'digital_rain', name: '⚡ DDP: Matrix Rain' },
    { id: 'pulse_beads', name: '⚡ DDP: Neon Pulse' },
    { id: 'cyber_fire', name: '⚡ DDP: Cyber Flame' }
  ];

  const audioEffects = [
    { id: 'audio_bass_pulse', name: '🎵 Sound: Bass Kick Pulse' },
    { id: 'audio_spectrum_waterfall', name: '🎵 Sound: 16-Bin Spectrum' },
    { id: 'audio_vu_meter', name: '🎵 Sound: Stereo VU Meter' },
    { id: 'audio_treble_sparkle', name: '🎵 Sound: Treble Sparkles' }
  ];

  function startEditing() {
    nickname = device.name;
    isEditingName = true;
  }

  const presetColors = [
    { name: 'Cyan Neon', hex: '#06b6d4', r: 6, g: 182, b: 212 },
    { name: 'Magenta Glow', hex: '#a855f7', r: 168, g: 85, b: 247 },
    { name: 'Cyber Amber', hex: '#f59e0b', r: 245, g: 158, b: 11 },
    { name: 'Void White', hex: '#ffffff', r: 255, g: 255, b: 255 },
    { name: 'Emerald Online', hex: '#10b981', r: 16, g: 185, b: 129 },
    { name: 'Crimson Red', hex: '#ef4444', r: 239, g: 68, b: 68 }
  ];

  function handleSaveName() {
    if (nickname.trim() && nickname !== device.name) {
      onRename(device.id, nickname.trim());
    }
    isEditingName = false;
  }

  function applyColor(c) {
    selectedColor = c.hex;
    onSetColor(device.id, c.r, c.g, c.b);
  }
</script>

<div
  class="glass-panel rounded-2xl p-4 flex flex-col justify-between space-y-3 relative group border h-full transition-all duration-300 {isMenuOpen
    ? 'z-50'
    : 'hover:z-30 z-10'} {device.isOnline
    ? 'border-cyan-500/20 hover:border-cyan-500/50'
    : 'border-rose-500/20 opacity-75'}"
>
  <!-- Top Accent Status Pill Bar (Sits cleanly inside card padding) -->
  <div
    class="w-full h-1 rounded-full bg-gradient-to-r {device.isOnline
      ? 'from-cyan-500 via-purple-500 to-amber-500'
      : 'from-rose-600 to-rose-900'}"
  ></div>

  <!-- FUNCTIONAL DEVICE HEADER & OPTIONS MENU -->
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2.5 flex-1">
      <div class="p-2 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-400">
        <Cpu class="w-4 h-4" />
      </div>
      <div>
        {#if isEditingName}
          <div class="flex items-center gap-1">
            <input
              type="text"
              bind:value={nickname}
              class="bg-[#090e17] border border-cyan-500/50 rounded-lg px-2 py-0.5 text-xs text-slate-100 font-semibold focus:outline-none"
            />
            <button onclick={handleSaveName} class="p-1 text-cyan-400 hover:text-cyan-200">
              <Check class="w-3.5 h-3.5" />
            </button>
          </div>
        {:else}
          <div class="flex items-center gap-1.5 group/edit">
            <h3 class="font-semibold text-slate-100 group-hover:text-cyan-300 transition-colors text-sm">
              {device.name}
            </h3>
            <button
              onclick={startEditing}
              class="opacity-0 group-hover/edit:opacity-100 p-0.5 text-slate-500 hover:text-slate-300"
            >
              <Edit2 class="w-3 h-3" />
            </button>
          </div>
        {/if}
        <p class="text-[11px] font-mono text-slate-400">
          {device.ipAddress} • {device.ledCount} LEDs
        </p>
      </div>
    </div>

    <!-- Power & Card Options Menu -->
    <div class="flex items-center gap-1.5 relative">
      <button
        onclick={() => onTogglePower(device.id)}
        class="p-2.5 rounded-xl transition-all duration-200 {device.state?.on
          ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan'
          : 'bg-[#090e17]/80 text-slate-600 border border-slate-800'} cursor-pointer"
        title="Toggle Power"
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
          {#if showSizeToggle}
            <button
              onclick={() => {
                onToggleSize(device.id);
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
              onTogglePin(device.id, 'device');
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-purple-300 hover:bg-purple-500/15 transition-all text-left cursor-pointer"
          >
            <Pin class="w-3.5 h-3.5 text-purple-400" />
            <span>{isPinned ? 'Unpin from Canvas' : 'Pin to Dashboard'}</span>
          </button>
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
        if (isDDPActive) ddpStore?.stopStream('device', device.id);
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
        ddpStore?.startStream('device', device.id, 'rainbow_wave', 1.0, 1.0, [device.ipAddress], device.ledCount);
      }}
      class="py-1.5 rounded-xl transition-all text-center cursor-pointer flex items-center justify-center gap-1 {cardMode ===
      'ddp'
        ? 'bg-amber-500/20 text-amber-300 border border-amber-500/40 glow-amber'
        : 'text-slate-500 hover:text-amber-400'}"
      title="DDP 60 FPS Procedural Matrix Mode"
    >
      <Zap class="w-3.5 h-3.5 {cardMode === 'ddp' ? 'animate-pulse text-amber-400' : ''}" />
      <span>{cardMode === 'ddp' ? `${ddpStatus.fps || 60} FPS` : 'DDP 60'}</span>
    </button>

    <button
      type="button"
      onclick={() => {
        if (audioStore && !audioStore.isCapturing) audioStore.startCapture();
        ddpStore?.startStream('device', device.id, 'audio_bass_pulse', 1.0, 1.0, [device.ipAddress], device.ledCount);
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

  <!-- Color Swatches & Controls Area -->
  <div class="space-y-3 pt-1">
    <!-- Quick Color Swatches -->
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

    <!-- Brightness Slider -->
    <div class="flex items-center gap-3 bg-[#090e17]/60 px-3 py-1.5 rounded-xl border border-slate-800/80 h-9">
      <Sun class="w-4 h-4 text-purple-400" />
      <input
        type="range"
        min="0"
        max="255"
        value={device.state?.bri ?? 128}
        oninput={(e) => onSetBrightness(device.id, parseInt(e.target.value))}
        disabled={!device.state?.on}
        class="flex-1 accent-purple-400 cursor-pointer disabled:opacity-30"
      />
      <span class="text-xs font-mono text-purple-300 w-8 text-right">
        {device.state?.bri ?? 0}
      </span>
    </div>

    <!-- Mode-Specific Controls Drawer & Dropdowns -->
    {#if cardMode === 'music'}
      <div class="space-y-2 bg-purple-500/10 p-3 rounded-2xl border border-purple-500/30 font-mono">
        <div class="flex items-center justify-between text-[11px] text-purple-300 font-bold">
          <span class="flex items-center gap-1.5">
            <Volume2 class="w-3.5 h-3.5 text-purple-400 animate-pulse" />
            Sound Sync Audio Reactivity
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
            ddpStore?.startStream(
              'device',
              device.id,
              e.target.value,
              ddpStatus.speed || 1.0,
              ddpStatus.intensity || 1.0,
              [device.ipAddress],
              device.ledCount
            );
          }}
          class="w-full bg-[#06090e] border border-purple-500/40 rounded-xl px-2.5 py-1 text-xs font-mono text-purple-300 focus:outline-none cursor-pointer"
        >
          {#each audioEffects as eff (eff.id)}
            <option value={eff.id}>{eff.name}</option>
          {/each}
        </select>
      </div>
    {:else if cardMode === 'ddp'}
      <div class="space-y-2 bg-amber-500/10 p-2.5 rounded-2xl border border-amber-500/30">
        <div class="flex items-center justify-between text-[11px] font-mono text-amber-400 font-bold">
          <span class="flex items-center gap-1">
            <Zap class="w-3 h-3 text-amber-400 animate-pulse" />
            DDP 60 FPS Procedural Engine
          </span>
          <span>{ddpStatus.fps || 60} FPS</span>
        </div>
        <select
          value={ddpStatus.effect || 'rainbow_wave'}
          onchange={(e) => {
            ddpStore?.startStream(
              'device',
              device.id,
              e.target.value,
              ddpStatus.speed || 1.0,
              ddpStatus.intensity || 1.0,
              [device.ipAddress],
              device.ledCount
            );
          }}
          class="w-full bg-[#06090e] border border-amber-500/40 rounded-xl px-2.5 py-1 text-xs font-mono text-amber-300 focus:outline-none cursor-pointer"
        >
          {#each ddpEffects as eff (eff.id)}
            <option value={eff.id}>{eff.name}</option>
          {/each}
        </select>
      </div>
    {:else}
      <div class="grid grid-cols-2 gap-2">
        <!-- FX CyberSelect -->
        <CyberSelect
          value={device.state?.seg?.[0]?.fx ?? 0}
          options={effects}
          icon={Sparkles}
          iconColor="text-cyan-400"
          hoverBorder="hover:border-cyan-500/40"
          onChange={(fxId) => onSetEffect(device.id, fxId)}
        />

        <!-- Palette CyberSelect -->
        <CyberSelect
          value={device.state?.seg?.[0]?.pal ?? 0}
          options={palettes}
          icon={Palette}
          iconColor="text-purple-400"
          hoverBorder="hover:border-purple-500/40"
          onChange={(palId) => {
            onSetPalette(device.id, palId);
            if ((device.state?.seg?.[0]?.fx ?? 0) === 0) {
              onSetEffect(device.id, 2);
            }
          }}
        />
      </div>
    {/if}
  </div>
</div>
