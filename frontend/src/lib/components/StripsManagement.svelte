<script>
  import { Cpu, Plus, RefreshCw, Pin, Edit2, Check, Trash2, Sliders } from 'lucide-svelte';

  let {
    devices = [],
    isScanning = false,
    dashboardStore,
    onDiscover = () => {},
    onAddManual = () => {},
    onRename = () => {}
  } = $props();

  let editingId = $state(null);
  let editName = $state('');

  function startRename(dev) {
    editingId = dev.id;
    editName = dev.name;
  }

  function saveRename(id) {
    if (editName.trim()) {
      onRename(id, editName.trim());
    }
    editingId = null;
  }
</script>

<div class="space-y-6">
  <!-- Header Banner -->
  <div
    class="glass-panel rounded-3xl p-6 flex flex-col md:flex-row md:items-center justify-between gap-4 border border-cyan-500/20"
  >
    <div>
      <h2 class="text-xl font-bold text-slate-100 flex items-center gap-2">
        <Sliders class="w-6 h-6 text-cyan-400" />
        Strips & Devices Management
      </h2>
      <p class="text-xs font-mono text-slate-400 mt-1">
        Complete network WLED strip inventory. Control dashboard visibility, rename strips, add direct IPs, or trigger
        mDNS discovery.
      </p>
    </div>

    <!-- Action Buttons -->
    <div class="flex items-center gap-2">
      <button
        onclick={onDiscover}
        disabled={isScanning}
        class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-slate-800/80 hover:bg-slate-800 text-slate-200 border border-slate-700 font-mono text-xs transition-all cursor-pointer disabled:opacity-50"
      >
        <RefreshCw class="w-4 h-4 text-cyan-400 {isScanning ? 'animate-spin' : ''}" />
        <span>{isScanning ? 'Scanning...' : 'Discover Strips'}</span>
      </button>

      <button
        onclick={onAddManual}
        class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-cyan-500/20 hover:bg-cyan-500/30 text-cyan-300 border border-cyan-500/40 font-mono text-xs shadow-neonCyan transition-all cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>Add Direct IP</span>
      </button>
    </div>
  </div>

  <!-- Devices Inventory Table / Grid -->
  <div class="glass-panel rounded-3xl p-6 border border-slate-800/80">
    {#if !devices || devices.length === 0}
      <div class="text-center py-12">
        <Cpu class="w-12 h-12 text-slate-600 mx-auto mb-3" />
        <p class="text-sm font-semibold text-slate-300">No WLED Light Strips Discovered</p>
        <p class="text-xs text-slate-500 font-mono mt-1">
          Click "Discover Strips" or "Add Direct IP" to register strips on your network.
        </p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each devices as dev (dev.id)}
          {@const isPinned = dashboardStore ? dashboardStore.isPinned(dev.id) : true}
          <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-4 p-4 rounded-2xl bg-slate-900/60 border border-slate-800 hover:border-cyan-500/30 transition-all duration-200"
          >
            <!-- Left Info Section -->
            <div class="flex items-center gap-4">
              <!-- Online Health Pill -->
              <span
                class="w-3 h-3 rounded-full {dev.isOnline
                  ? 'bg-emerald-400 shadow-[0_0_10px_rgba(52,211,153,0.5)] animate-pulse'
                  : 'bg-rose-500'}"
                title={dev.isOnline ? 'Online' : 'Offline'}
              ></span>

              <div>
                {#if editingId === dev.id}
                  <div class="flex items-center gap-2">
                    <input
                      type="text"
                      bind:value={editName}
                      class="bg-slate-950 border border-cyan-500/50 rounded-lg px-2 py-1 text-xs text-slate-100 font-semibold focus:outline-none"
                    />
                    <button onclick={() => saveRename(dev.id)} class="p-1 text-cyan-400 hover:text-cyan-200">
                      <Check class="w-4 h-4" />
                    </button>
                  </div>
                {:else}
                  <div class="flex items-center gap-2">
                    <h4 class="font-semibold text-slate-100 text-sm">{dev.name}</h4>
                    <button onclick={() => startRename(dev)} class="p-1 text-slate-500 hover:text-slate-300">
                      <Edit2 class="w-3 h-3" />
                    </button>
                  </div>
                {/if}

                <div class="flex items-center gap-3 text-[11px] font-mono text-slate-400 mt-0.5">
                  <span>{dev.ipAddress}</span>
                  <span>•</span>
                  <span>{dev.macAddress || 'MAC Unknown'}</span>
                  <span>•</span>
                  <span>{dev.ledCount} LEDs</span>
                </div>
              </div>
            </div>

            <!-- Right Controls: Pin & Delete -->
            <div class="flex items-center gap-3">
              <!-- Dashboard Pin Toggle -->
              {#if dashboardStore}
                <button
                  onclick={() => dashboardStore.togglePin(dev.id, 'device')}
                  class="flex items-center gap-2 px-3 py-1.5 rounded-xl border text-xs font-mono transition-all cursor-pointer {isPinned
                    ? 'bg-purple-500/20 text-purple-300 border-purple-500/40'
                    : 'bg-slate-800/60 text-slate-500 border-slate-700 hover:text-slate-300'}"
                >
                  <Pin class="w-3.5 h-3.5 {isPinned ? 'fill-current' : ''}" />
                  <span>{isPinned ? 'Pinned to Dashboard' : 'Hidden from Dashboard'}</span>
                </button>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
