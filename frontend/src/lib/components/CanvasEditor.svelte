<script>
  import { onDestroy } from 'svelte';
  import {
    Grid,
    Save,
    Sparkles,
    Move,
    RotateCw,
    RotateCcw,
    ZoomIn,
    ZoomOut,
    Maximize2,
    Cpu,
    ArrowLeft,
    Edit3
  } from 'lucide-svelte';

  let {
    devices = [],
    placements = [],
    rooms = [],
    currentRoomId = null,
    onSavePlacements = () => {},
    onUpdatePlacement = () => {},
    onTriggerSweep = () => {},
    onSelectRoom = () => {},
    onCreateRoom = () => {},
    onRenameRoom = () => {},
    onDeleteRoom = () => {},
    onBackToGroups = () => {}
  } = $props();

  let selectedDeviceId = $state(null);
  let isDragging = $state(false);
  let dragOffset = $state({ x: 0, y: 0 });
  let snapToGrid = $state(true);
  let sweepActive = $state(false);
  let sweepTime = $state(0);
  let animationFrameId = null;
  let isEditingTitle = $state(false);
  let editedTitle = $state('');

  let currentRoom = $derived(
    rooms.find((r) => r.id === currentRoomId) || { title: '2D Room Layout', width: 2000, height: 1200 }
  );

  // Devices belonging to this specific room
  let roomDevices = $derived.by(() => {
    if (!currentRoomId) return devices;
    const placedDevIds = placements.filter((p) => p.roomId === currentRoomId).map((p) => p.deviceId);

    if (placedDevIds.length === 0) return devices;
    const filtered = devices.filter((d) => placedDevIds.includes(d.id));
    return filtered.length > 0 ? filtered : devices;
  });

  async function handleSaveLayout() {
    // Record current placements for all roomDevices into store placements array
    roomDevices.forEach((dev) => {
      const p = getPlacement(dev.id);
      onUpdatePlacement(dev.id, {
        posX: p.posX,
        posY: p.posY,
        rotation: p.rotation,
        scale: p.scale,
        geometry: p.geometry
      });
    });

    await onSavePlacements();
    onBackToGroups();
  }

  // Transient slider angle preview map (degrees readout while dragging)
  let sliderPreviewAngle = $state({});

  onDestroy(() => {
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId);
    }
  });

  // Helper to find or initialize placement for a device
  function getPlacement(devId) {
    const found = placements.find((p) => p.deviceId === devId);
    return (
      found || {
        deviceId: devId,
        posX: 100 + (devices.findIndex((d) => d.id === devId) % 4) * 220,
        posY: 100 + Math.floor(devices.findIndex((d) => d.id === devId) / 4) * 160,
        rotation: 0,
        scale: 1.0,
        geometry: 'strip'
      }
    );
  }

  function handlePointerDown(e, devId) {
    selectedDeviceId = devId;
    isDragging = true;
    const placement = getPlacement(devId);
    dragOffset = {
      x: e.clientX - placement.posX,
      y: e.clientY - placement.posY
    };
  }

  function handlePointerMove(e) {
    if (!isDragging || !selectedDeviceId) return;
    let rawX = e.clientX - dragOffset.x;
    let rawY = e.clientY - dragOffset.y;

    if (snapToGrid) {
      rawX = Math.round(rawX / 20) * 20;
      rawY = Math.round(rawY / 20) * 20;
    }

    rawX = Math.max(20, Math.min(1600, rawX));
    rawY = Math.max(20, Math.min(1000, rawY));

    onUpdatePlacement(selectedDeviceId, { posX: rawX, posY: rawY });
  }

  function handlePointerUp() {
    isDragging = false;
  }

  function resetRotation(devId) {
    delete sliderPreviewAngle[devId];
    onUpdatePlacement(devId, { rotation: 0 });
  }

  function resetAllRotations() {
    sliderPreviewAngle = {};
    devices.forEach((dev) => {
      onUpdatePlacement(dev.id, { rotation: 0 });
    });
  }

  function stepRotate(devId, currentRotation, step) {
    const nextRot = (Math.round(currentRotation + step) + 360) % 360;
    delete sliderPreviewAngle[devId];
    onUpdatePlacement(devId, { rotation: nextRot });
  }

  function handleSliderInput(devId, val) {
    sliderPreviewAngle[devId] = Number(val);
  }

  function handleSliderChange(devId, val) {
    const finalRot = Number(val);
    delete sliderPreviewAngle[devId];
    onUpdatePlacement(devId, { rotation: finalRot });
  }

  function triggerSpatialSweep() {
    if (animationFrameId) {
      cancelAnimationFrame(animationFrameId);
    }

    // 1. Snapshot original physical states
    const snapshots = devices
      .filter((dev) => dev && dev.isOnline && dev.ipAddress)
      .map((dev) => ({
        ip: dev.ipAddress,
        originalState: dev.state ? JSON.parse(JSON.stringify(dev.state)) : null
      }));

    sweepActive = true;
    sweepTime = 0;
    const startTime = Date.now();

    // Call store sweep callback
    onTriggerSweep();

    // Trigger physical WLED strips via REST API with correct IP parameter
    snapshots.forEach((item) => {
      if (typeof window !== 'undefined' && window.go?.main?.App?.SetDeviceEffect) {
        window.go.main.App.SetDeviceEffect(item.ip, 9);
      } else {
        fetch('/api/v1/devices/state', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            ip: item.ip,
            state: { on: true, bri: 255, seg: [{ id: 0, fx: 9, sx: 160 }] }
          })
        }).catch((err) => console.warn('[CanvasEditor] Failed to send sweep to WLED strip:', item.ip, err));
      }
    });

    const animate = () => {
      const elapsed = Date.now() - startTime;
      if (elapsed >= 3500) {
        sweepActive = false;
        animationFrameId = null;

        // 2. Restore original pre-sweep states to physical WLED strips
        snapshots.forEach((item) => {
          if (item.originalState) {
            if (typeof window !== 'undefined' && window.go?.main?.App?.SetDeviceEffect) {
              const origFx = item.originalState.seg?.[0]?.fx ?? 0;
              window.go.main.App.SetDeviceEffect(item.ip, origFx);
            } else {
              fetch('/api/v1/devices/state', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                  ip: item.ip,
                  state: item.originalState
                })
              }).catch(() => {});
            }
          }
        });
        return;
      }
      sweepTime = elapsed;
      animationFrameId = requestAnimationFrame(animate);
    };

    animationFrameId = requestAnimationFrame(animate);
  }

  function getPixelColors(dev) {
    const p = getPlacement(dev.id);
    const count = Math.min(dev.ledCount || 30, 20); // Render up to 20 LED dots on preview
    const colors = [];

    for (let i = 0; i < count; i++) {
      if (sweepActive) {
        // Spatial 2D coordinate phase mapping based on physical placement (p.posX, p.posY)
        const spatialPhase = (p.posX * 0.3 + p.posY * 0.2 + i * 10 + sweepTime * 0.4) % 360;
        colors.push(`hsl(${spatialPhase}, 100%, 60%)`);
      } else if (dev.isOnline) {
        colors.push('#06b6d4');
      } else {
        colors.push('#334155');
      }
    }
    return colors;
  }
</script>

<div
  class="space-y-6"
  role="region"
  aria-label="2D Room Layout Canvas Editor"
  onpointermove={handlePointerMove}
  onpointerup={handlePointerUp}
>
  <!-- Header Banner & Action Bar -->
  <div
    class="glass-panel rounded-3xl p-6 flex flex-col md:flex-row md:items-center justify-between gap-4 border border-cyan-500/20 shadow-2xl"
  >
    <div>
      {#if isEditingTitle}
        <form
          onsubmit={(e) => {
            e.preventDefault();
            if (editedTitle.trim() && editedTitle.trim() !== currentRoom.title) {
              onRenameRoom(currentRoom.id, editedTitle.trim());
            }
            isEditingTitle = false;
          }}
          class="flex items-center gap-2"
        >
          <Grid class="w-6 h-6 text-cyan-400" />
          <input
            type="text"
            bind:value={editedTitle}
            onblur={() => {
              if (editedTitle.trim() && editedTitle.trim() !== currentRoom.title) {
                onRenameRoom(currentRoom.id, editedTitle.trim());
              }
              isEditingTitle = false;
            }}
            class="bg-[#06090e] border border-cyan-500/50 rounded-lg px-2.5 py-0.5 text-lg font-bold text-slate-100 focus:outline-none focus:border-cyan-400 w-64"
          />
        </form>
      {:else}
        <button
          type="button"
          onclick={() => {
            isEditingTitle = true;
            editedTitle = currentRoom.title;
          }}
          class="group/roomTitle flex items-center gap-2 text-xl font-bold text-slate-100 hover:text-cyan-300 transition-colors cursor-pointer text-left border-0 bg-transparent p-0"
          title="Click to rename room"
        >
          <Grid class="w-6 h-6 text-cyan-400" />
          <span>{currentRoom.title} 2D Layout Canvas</span>
          <Edit3 class="w-4 h-4 text-slate-500 opacity-0 group-hover/roomTitle:opacity-100 transition-opacity" />
        </button>
      {/if}
      <p class="text-xs font-mono text-slate-400 mt-1">
        Interactive 2D room map for dragging, rotating, and positioning WLED light strips with live pixel mirroring.
      </p>
    </div>

    <!-- Controls Toolbar -->
    <div class="flex flex-wrap items-center gap-3">
      <!-- Grid Snap Toggle -->
      <button
        type="button"
        onclick={() => (snapToGrid = !snapToGrid)}
        class="flex items-center gap-2 px-3 py-2 rounded-xl border text-xs font-mono transition-all cursor-pointer {snapToGrid
          ? 'bg-cyan-500/20 border-cyan-500/40 text-cyan-300'
          : 'bg-slate-800/80 border-slate-700 text-slate-400'}"
      >
        <Grid class="w-4 h-4" />
        <span>Snap: {snapToGrid ? 'ON' : 'OFF'}</span>
      </button>

      <!-- Reset All Rotations -->
      <button
        type="button"
        onclick={resetAllRotations}
        class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-slate-800/80 hover:bg-slate-700 border border-slate-700 text-slate-300 font-mono text-xs transition-all cursor-pointer"
        title="Reset all strip rotations to 0°"
      >
        <RotateCcw class="w-4 h-4 text-cyan-400" />
        <span>Reset</span>
      </button>

      <!-- Spatial Sweep Button -->
      <button
        type="button"
        onclick={triggerSpatialSweep}
        class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-cyan-500 hover:from-purple-500 hover:to-cyan-400 text-white font-mono text-xs font-semibold shadow-neonPurple transition-all cursor-pointer"
      >
        <Sparkles class="w-4 h-4 {sweepActive ? 'animate-spin' : ''}" />
        <span>{sweepActive ? 'Sweeping Wave...' : 'Spatial Sweep'}</span>
      </button>

      <!-- Save Placements & Return to Groups -->
      <button
        type="button"
        onclick={handleSaveLayout}
        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all cursor-pointer"
      >
        <Save class="w-4 h-4" />
        <span>Save</span>
      </button>
    </div>
  </div>

  <!-- Interactive 2D Room Grid Area -->
  <div
    class="glass-panel rounded-3xl border border-slate-800/80 relative overflow-hidden min-h-[600px] select-none"
    style="background-color: #06090e; background-image: radial-gradient(rgba(6, 182, 212, 0.12) 1px, transparent 0); background-size: 20px 20px;"
  >
    <!-- Background Canvas Info Watermark -->
    <div class="absolute top-4 left-4 pointer-events-none opacity-40 font-mono text-xs text-cyan-500/60 space-y-1">
      <p>ROOM CANVAS GRID (1200x700 px)</p>
      <p>ACTIVE STRIPS: {roomDevices.length}</p>
    </div>

    <!-- Placed Strip Elements -->
    {#each roomDevices as dev (dev.id)}
      {@const p = getPlacement(dev.id)}
      {@const isSelected = selectedDeviceId === dev.id}
      {@const pixelColors = getPixelColors(dev)}
      {@const displayAngle = sliderPreviewAngle[dev.id] ?? p.rotation}

      <div
        role="button"
        tabindex="0"
        aria-label="Drag element for {dev.name}"
        style="transform: translate({p.posX}px, {p.posY}px) rotate({p.rotation}deg) scale({p.scale});"
        class="absolute cursor-grab active:cursor-grabbing transition-transform duration-200 ease-out group z-10"
        onpointerdown={(e) => handlePointerDown(e, dev.id)}
        onkeydown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') selectedDeviceId = dev.id;
        }}
      >
        <!-- Strip Container Box -->
        <div
          class="glass-panel rounded-2xl p-3 border border-slate-800 hover:border-cyan-500/60 shadow-lg min-w-[240px] space-y-2.5 bg-slate-900/90 backdrop-blur-md {isSelected
            ? 'border-cyan-400 ring-2 ring-cyan-500/30'
            : ''}"
        >
          <!-- Header Bar -->
          <div class="flex items-center justify-between text-xs font-mono">
            <span class="font-bold text-slate-100 truncate max-w-[140px]">{dev.name}</span>
            <span class="text-[10px] text-slate-400">{dev.ledCount || 30} LEDs</span>
          </div>

          <!-- Live LED Pixel Mirroring Line -->
          <div
            class="flex items-center gap-1.5 p-2 rounded-xl bg-slate-950/80 border border-slate-800/80 overflow-x-auto"
          >
            {#each pixelColors as color, idx (idx)}
              <div
                class="w-3 h-3 rounded-full flex-shrink-0 transition-colors duration-200"
                style="background-color: {color}; box-shadow: {dev.isOnline || sweepActive
                  ? `0 0 6px ${color}`
                  : 'none'};"
              ></div>
            {/each}
          </div>

          <!-- Controls Overlay on Selected (event propagation stopped) -->
          {#if isSelected}
            <div
              role="presentation"
              class="space-y-2 pt-2 border-t border-slate-800/60 text-[10px] font-mono"
              onpointerdown={(e) => e.stopPropagation()}
            >
              <!-- Rotation Step Buttons & Reset -->
              <div class="flex items-center justify-between gap-1">
                <div class="flex items-center gap-1">
                  <button
                    type="button"
                    title="Rotate -45°"
                    onclick={() => stepRotate(dev.id, p.rotation, -45)}
                    class="px-1.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-cyan-300 font-mono transition-colors cursor-pointer"
                  >
                    -45°
                  </button>
                  <button
                    type="button"
                    title="Rotate +45°"
                    onclick={() => stepRotate(dev.id, p.rotation, 45)}
                    class="px-1.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-cyan-300 font-mono transition-colors cursor-pointer"
                  >
                    +45°
                  </button>
                  <button
                    type="button"
                    title="Rotate +90°"
                    onclick={() => stepRotate(dev.id, p.rotation, 90)}
                    class="px-1.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-cyan-300 font-mono transition-colors cursor-pointer"
                  >
                    +90°
                  </button>
                </div>

                <!-- Explicit Reset Rotation to 0° -->
                <button
                  type="button"
                  title="Reset rotation to 0°"
                  onclick={() => resetRotation(dev.id)}
                  class="px-2 py-1 rounded bg-cyan-500/20 hover:bg-cyan-500/30 text-cyan-300 border border-cyan-500/40 flex items-center gap-1 font-bold transition-all cursor-pointer"
                >
                  <RotateCcw class="w-3 h-3 text-cyan-400" />
                  <span>0° Reset</span>
                </button>
              </div>

              <!-- Rotation Slider & Geometry Selector -->
              <div class="flex items-center justify-between gap-2">
                <div class="flex items-center gap-1.5">
                  <RotateCw class="w-3.5 h-3.5 text-cyan-400 flex-shrink-0" />
                  <input
                    type="range"
                    min="0"
                    max="360"
                    step="15"
                    value={displayAngle}
                    oninput={(e) => handleSliderInput(dev.id, e.target.value)}
                    onchange={(e) => handleSliderChange(dev.id, e.target.value)}
                    class="w-20 accent-cyan-400 cursor-pointer"
                  />
                  <span class="text-slate-400 min-w-[28px]">{displayAngle}°</span>
                </div>

                <select
                  value={p.geometry || 'strip'}
                  onchange={(e) => onUpdatePlacement(dev.id, { geometry: e.target.value })}
                  class="bg-slate-950 text-slate-300 border border-slate-800 rounded px-1.5 py-0.5"
                >
                  <option value="strip">Strip</option>
                  <option value="matrix">Matrix</option>
                  <option value="ring">Ring</option>
                </select>
              </div>
            </div>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>
