<script>
  import { LayoutGrid, Edit3, Power, Sun, Pin, PinOff, Trash2 } from 'lucide-svelte';

  let {
    room = { id: 'default', title: 'Main Room Canvas', width: 2000, height: 1200 },
    devices = [],
    placements = [],
    isPinned = true,
    cardSize = 'normal',
    onEditLayout = () => {},
    onTogglePin = () => {}
  } = $props();

  // Filter placements for this room
  let roomPlacements = $derived(
    placements.filter((p) => p.roomId === room.id || (room.id === 'default' && (!p.roomId || p.roomId === 'default')))
  );

  // Devices in this room
  let roomDevices = $derived(devices.filter((dev) => roomPlacements.some((p) => p.deviceId === dev.id)));

  let isAnyOn = $derived(roomDevices.some((d) => d.state?.on));
  let avgBri = $derived(
    roomDevices.length > 0
      ? Math.round(roomDevices.reduce((sum, d) => sum + (d.state?.bri || 0), 0) / roomDevices.length)
      : 255
  );

  function toggleRoomPower() {
    const nextOn = !isAnyOn;
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
</script>

<div
  class="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-cyan-500/20 bg-slate-900/65 p-5 backdrop-blur-xl transition-all duration-300 hover:border-purple-500/40 hover:shadow-[0_0_25px_rgba(168,85,247,0.15)]"
>
  <!-- Card Header -->
  <div class="flex items-start justify-between">
    <div class="flex items-center gap-3">
      <div
        class="flex h-10 w-10 items-center justify-center rounded-xl border border-cyan-400/30 bg-cyan-950/40 text-cyan-400 shadow-[0_0_12px_rgba(6,182,212,0.2)]"
      >
        <LayoutGrid class="h-5 w-5" />
      </div>
      <div>
        <h3 class="font-semibold text-slate-100">{room.title}</h3>
        <p class="text-xs text-slate-400">{roomDevices.length} strips placed • 2D Canvas</p>
      </div>
    </div>

    <!-- Pin Toggle Button -->
    <button
      onclick={onTogglePin}
      class="rounded-lg p-2 text-slate-400 transition-colors hover:bg-slate-800/80 hover:text-cyan-400"
      title={isPinned ? 'Unpin from dashboard' : 'Pin to dashboard'}
    >
      {#if isPinned}
        <Pin class="h-4 w-4 text-cyan-400" />
      {:else}
        <PinOff class="h-4 w-4" />
      {/if}
    </button>
  </div>

  <!-- Mini 2D Room Canvas Preview Grid -->
  <div class="relative my-4 h-32 w-full overflow-hidden rounded-xl border border-slate-800 bg-[#06090e] p-2">
    <!-- Background Grid Lines -->
    <div
      class="absolute inset-0 bg-[linear-gradient(to_right,rgba(56,189,248,0.05)_1px,transparent_1px),linear-gradient(to_bottom,rgba(56,189,248,0.05)_1px,transparent_1px)] bg-[size:16px_16px]"
    ></div>

    <!-- Mini Placements Preview -->
    {#each roomPlacements as p (p.deviceId)}
      {@const dev = devices.find((d) => d.id === p.deviceId)}
      {@const posX = (p.posX / 2000) * 100}
      {@const posY = (p.posY / 1200) * 100}
      {@const rot = p.rotation || 0}
      {@const isOn = dev?.state?.on ?? false}
      {@const color = dev?.state?.seg?.[0]?.col?.[0]
        ? `rgb(${dev.state.seg[0].col[0][0]}, ${dev.state.seg[0].col[0][1]}, ${dev.state.seg[0].col[0][2]})`
        : '#06b6d4'}

      <div
        class="absolute flex items-center justify-center rounded-full transition-transform duration-300"
        style="left: {posX}%; top: {posY}%; transform: translate(-50%, -50%) rotate({rot}deg);"
      >
        <!-- Strip Line Indicator -->
        <div
          class="h-2 w-12 rounded-full border border-cyan-400/40 shadow-sm"
          style="background-color: {isOn ? color : '#334155'}; box-shadow: {isOn ? `0 0 8px ${color}` : 'none'};"
        ></div>
      </div>
    {/each}
  </div>

  <!-- Quick Room Controls & Footer -->
  <div class="flex items-center justify-between gap-3 pt-2">
    <!-- Master Room Power Toggle -->
    <button
      onclick={toggleRoomPower}
      class="flex items-center gap-2 rounded-xl border px-3 py-1.5 text-xs font-medium transition-all duration-200 {isAnyOn
        ? 'border-emerald-500/50 bg-emerald-950/40 text-emerald-300 shadow-[0_0_12px_rgba(16,185,129,0.2)]'
        : 'border-slate-800 bg-slate-900/60 text-slate-400 hover:border-slate-700'}"
    >
      <Power class="h-3.5 w-3.5 {isAnyOn ? 'text-emerald-400 animate-pulse' : ''}" />
      <span>{isAnyOn ? 'Room ON' : 'Room OFF'}</span>
    </button>

    <!-- Master Room Brightness Slider -->
    <div class="flex flex-1 items-center gap-2 px-2">
      <Sun class="h-3.5 w-3.5 text-slate-400" />
      <input
        type="range"
        min="0"
        max="255"
        value={avgBri}
        onchange={handleBriChange}
        class="h-1.5 w-full cursor-pointer appearance-none rounded-lg bg-slate-800 accent-cyan-400"
      />
    </div>

    <!-- Edit 2D Layout Button -->
    <button
      onclick={onEditLayout}
      class="flex items-center gap-1.5 rounded-xl border border-cyan-500/30 bg-cyan-950/30 px-3 py-1.5 text-xs font-semibold text-cyan-300 transition-all duration-200 hover:border-cyan-400 hover:bg-cyan-900/40 hover:text-cyan-200 shadow-sm"
    >
      <Edit3 class="h-3.5 w-3.5" />
      <span>Edit 2D Layout</span>
    </button>
  </div>
</div>
