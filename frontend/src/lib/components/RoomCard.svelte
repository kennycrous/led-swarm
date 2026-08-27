<script>
  import { LayoutGrid, Power, Sun, Palette, Sparkles, Pin, Maximize2, MoreVertical, Trash2 } from 'lucide-svelte';
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
    onEditLayout = () => {},
    onTogglePin = () => {},
    onToggleSize = () => {},
    onDelete = null
  } = $props();

  let isMenuOpen = $state(false);
  let roomPower = $state(true);
  let roomBrightness = $state(200);
  let selectedColor = $state('#06b6d4');
  let selectedEffect = $state(0);
  let selectedPalette = $state(0);

  const presetColors = [
    { name: 'Cyan Neon', hex: '#06b6d4', r: 6, g: 182, b: 212 },
    { name: 'Magenta Glow', hex: '#a855f7', r: 168, g: 85, b: 247 },
    { name: 'Cyber Amber', hex: '#f59e0b', r: 245, g: 158, b: 11 },
    { name: 'Void White', hex: '#ffffff', r: 255, g: 255, b: 255 },
    { name: 'Emerald Online', hex: '#10b981', r: 16, g: 185, b: 129 },
    { name: 'Crimson Red', hex: '#ef4444', r: 239, g: 68, b: 68 }
  ];

  // Filter placements for this room
  let roomPlacements = $derived(
    placements.filter((p) => p.roomId === room.id || (room.id === 'default' && (!p.roomId || p.roomId === 'default')))
  );

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
        <h3 class="font-semibold text-slate-100 text-sm tracking-wide">{room.title}</h3>
        <p class="text-[11px] font-mono text-slate-400">
          {roomDevices.length} Strips ({onlineCount} Online)
        </p>
      </div>
    </div>

    <!-- Power, Options Menu & Controls -->
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

    <!-- Custom Glassmorphic CyberSelect Dropdowns -->
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
