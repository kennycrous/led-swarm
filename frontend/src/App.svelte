<script>
  import { onMount } from 'svelte';
  import { 
    Zap, 
    Layers, 
    LayoutGrid, 
    Radio, 
    Sliders, 
    Settings, 
    RefreshCw, 
    Wifi, 
    Power, 
    Activity,
    Plus
  } from 'lucide-svelte';
  import { getDeviceStore } from '$lib/stores/deviceStore.svelte.js';
  import DeviceCard from '$lib/components/DeviceCard.svelte';
  import ManualIpModal from '$lib/components/ManualIpModal.svelte';

  const store = getDeviceStore();

  let activeTab = $state('dashboard');
  let masterPower = $state(true);
  let masterBrightness = $state(200);
  let isAddModalOpen = $state(false);

  onMount(() => {
    store.init();
  });

  function toggleMasterPower() {
    masterPower = !masterPower;
    store.devices.filter(d => d.isOnline).forEach(d => store.togglePower(d.id));
  }

  function handleMasterBrightness(val) {
    masterBrightness = val;
    store.devices.filter(d => d.isOnline).forEach(d => store.setBrightness(d.id, val));
  }
</script>

<div class="flex h-screen w-screen bg-[#06090e] text-slate-100 font-sans antialiased overflow-hidden">
  
  <!-- Left Navigation Rail -->
  <aside class="w-16 flex flex-col items-center justify-between py-5 border-r border-cyan-500/10 bg-[#080d14]/80 backdrop-blur-xl z-20">
    <div class="flex flex-col items-center gap-6">
      <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-cyan-500 to-purple-600 flex items-center justify-center shadow-lg shadow-cyan-500/20">
        <Zap class="w-6 h-6 text-white" />
      </div>

      <nav class="flex flex-col gap-4 mt-4">
        <button 
          onclick={() => activeTab = 'dashboard'}
          class="p-2.5 rounded-xl transition-all duration-200 {activeTab === 'dashboard' ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'}"
          title="Dashboard"
        >
          <LayoutGrid class="w-5 h-5" />
        </button>

        <button 
          onclick={() => activeTab = 'canvas'}
          class="p-2.5 rounded-xl transition-all duration-200 {activeTab === 'canvas' ? 'bg-purple-500/20 text-purple-400 border border-purple-500/40 glow-magenta' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'}"
          title="2D Visual Canvas"
        >
          <Layers class="w-5 h-5" />
        </button>

        <button 
          onclick={() => activeTab = 'audio'}
          class="p-2.5 rounded-xl transition-all duration-200 {activeTab === 'audio' ? 'bg-amber-500/20 text-amber-400 border border-amber-500/40' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'}"
          title="Audio Reactivity"
        >
          <Radio class="w-5 h-5" />
        </button>

        <button 
          onclick={() => activeTab = 'scenes'}
          class="p-2.5 rounded-xl transition-all duration-200 {activeTab === 'scenes' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/40' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'}"
          title="Groups & Scenes"
        >
          <Sliders class="w-5 h-5" />
        </button>
      </nav>
    </div>

    <button 
      onclick={() => activeTab = 'settings'}
      class="p-2.5 rounded-xl text-slate-400 hover:text-slate-200 hover:bg-slate-800/40 transition-all duration-200"
      title="Settings"
    >
      <Settings class="w-5 h-5" />
    </button>
  </aside>

  <!-- Main View Area -->
  <div class="flex-1 flex flex-col overflow-hidden">
    
    <!-- Top Cyber Bar -->
    <header class="h-16 border-b border-cyan-500/10 bg-[#090e17]/60 backdrop-blur-xl px-6 flex items-center justify-between z-10">
      <div class="flex items-center gap-3">
        <h1 class="text-lg font-bold tracking-wider bg-gradient-to-r from-cyan-400 via-purple-400 to-amber-400 bg-clip-text text-transparent">
          LED SWARM ORCHESTRATOR
        </h1>
        <span class="px-2 py-0.5 text-[10px] font-mono tracking-widest bg-cyan-500/10 text-cyan-400 border border-cyan-500/20 rounded-full">
          v0.2.0-LIVE
        </span>
      </div>

      <!-- Master Swarm Controls -->
      <div class="flex items-center gap-4">
        <!-- Master Brightness -->
        <div class="flex items-center gap-3 bg-slate-900/60 border border-slate-800 px-3 py-1.5 rounded-xl">
          <span class="text-xs text-slate-400 font-mono">BRI</span>
          <input 
            type="range" 
            min="0" 
            max="255" 
            value={masterBrightness}
            oninput={(e) => handleMasterBrightness(parseInt(e.target.value))}
            class="w-24 accent-cyan-400 cursor-pointer"
          />
          <span class="text-xs text-cyan-400 font-mono w-8 text-right">{masterBrightness}</span>
        </div>

        <!-- Add IP Button -->
        <button 
          onclick={() => isAddModalOpen = true}
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-purple-500/10 border border-purple-500/30 text-purple-300 text-xs font-mono hover:bg-purple-500/20 hover:border-purple-400 transition-all duration-200"
        >
          <Plus class="w-3.5 h-3.5" />
          ADD IP
        </button>

        <!-- mDNS Scan Button -->
        <button 
          onclick={() => store.triggerScan()}
          disabled={store.isScanning}
          class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-300 text-xs font-mono hover:bg-cyan-500/20 hover:border-cyan-400 transition-all duration-200 disabled:opacity-50"
        >
          <RefreshCw class="w-3.5 h-3.5 {store.isScanning ? 'animate-spin text-cyan-400' : ''}" />
          {store.isScanning ? 'SCANNING...' : 'DISCOVER STREAMS'}
        </button>

        <!-- Master Power -->
        <button 
          onclick={toggleMasterPower}
          class="p-2 rounded-xl border transition-all duration-200 {masterPower ? 'bg-cyan-500/20 text-cyan-300 border-cyan-500/50 glow-cyan' : 'bg-slate-900/60 text-slate-500 border-slate-800'}"
          title="Master Swarm Power"
        >
          <Power class="w-4 h-4" />
        </button>
      </div>
    </header>

    <!-- Main Workspace Panel -->
    <main class="flex-1 overflow-y-auto p-8 relative">
      {#if activeTab === 'dashboard'}
        <div class="max-w-7xl mx-auto space-y-6">
          
          <!-- Swarm Metric Banner -->
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="glass-panel p-4 rounded-2xl flex items-center justify-between">
              <div>
                <p class="text-xs font-mono text-slate-400">TOTAL DEVICES</p>
                <p class="text-2xl font-bold text-cyan-400 mt-1 font-mono">{store.devices.length}</p>
              </div>
              <Wifi class="w-8 h-8 text-cyan-500/30" />
            </div>

            <div class="glass-panel p-4 rounded-2xl flex items-center justify-between">
              <div>
                <p class="text-xs font-mono text-slate-400">ACTIVE LEDS</p>
                <p class="text-2xl font-bold text-purple-400 mt-1 font-mono">
                  {store.devices.reduce((acc, d) => acc + (d.state?.on ? d.ledCount : 0), 0)}
                </p>
              </div>
              <Activity class="w-8 h-8 text-purple-500/30" />
            </div>

            <div class="glass-panel p-4 rounded-2xl flex items-center justify-between">
              <div>
                <p class="text-xs font-mono text-slate-400">ONLINE STATUS</p>
                <p class="text-2xl font-bold text-emerald-400 mt-1 font-mono">
                  {store.devices.filter(d => d.isOnline).length} / {store.devices.length}
                </p>
              </div>
              <Zap class="w-8 h-8 text-emerald-500/30" />
            </div>

            <div class="glass-panel p-4 rounded-2xl flex items-center justify-between">
              <div>
                <p class="text-xs font-mono text-slate-400">REALTIME SYNC</p>
                <p class="text-2xl font-bold text-amber-400 mt-1 font-mono">
                  {store.isConnected ? 'LIVE WS' : 'REST'}
                </p>
              </div>
              <Radio class="w-8 h-8 text-amber-500/30" />
            </div>
          </div>

          <!-- Devices Grid -->
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h2 class="text-sm font-mono tracking-wider text-slate-400 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
                DISCOVERED WLED STRIPS ({store.devices.length})
              </h2>
            </div>

            {#if store.devices.length === 0}
              <div class="glass-panel rounded-3xl p-12 text-center border border-dashed border-cyan-500/20">
                <Wifi class="w-12 h-12 text-cyan-500/40 mx-auto mb-3 animate-pulse" />
                <h3 class="text-base font-bold text-slate-200">No WLED Devices Found Yet</h3>
                <p class="text-xs font-mono text-slate-500 mt-1">Click "DISCOVER STREAMS" to scan mDNS (_wled._tcp) or click "ADD IP" to add manually.</p>
              </div>
            {:else}
              <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                {#each store.devices as dev (dev.id)}
                  <DeviceCard 
                    device={dev}
                    effects={store.effects}
                    palettes={store.palettes}
                    onTogglePower={(id) => store.togglePower(id)}
                    onSetBrightness={(id, bri) => store.setBrightness(id, bri)}
                    onSetColor={(id, r, g, b) => store.setColor(id, r, g, b)}
                    onSetEffect={(id, fx) => store.setEffect(id, fx)}
                    onSetPalette={(id, pal) => store.setPalette(id, pal)}
                    onRename={(id, name) => store.renameDevice(id, name)}
                  />
                {/each}
              </div>
            {/if}
          </div>
        </div>

      {:else if activeTab === 'canvas'}
        <div class="h-full flex items-center justify-center border border-dashed border-purple-500/20 rounded-3xl bg-slate-950/40">
          <div class="text-center space-y-3">
            <Layers class="w-12 h-12 text-purple-400/40 mx-auto" />
            <h2 class="text-lg font-bold text-slate-300">2D Visual Layout Canvas</h2>
            <p class="text-xs font-mono text-slate-500">Arrange WLED light strips in 2D physical spatial room layout.</p>
          </div>
        </div>

      {:else}
        <div class="h-full flex items-center justify-center border border-dashed border-cyan-500/20 rounded-3xl bg-slate-950/40">
          <div class="text-center space-y-3">
            <Zap class="w-12 h-12 text-cyan-400/40 mx-auto" />
            <h2 class="text-lg font-bold text-slate-300">Swarm Module Active</h2>
            <p class="text-xs font-mono text-slate-500">Module controls coming online.</p>
          </div>
        </div>
      {/if}
    </main>
  </div>
</div>

<ManualIpModal 
  isOpen={isAddModalOpen} 
  onClose={() => isAddModalOpen = false} 
  onAdd={(ip) => store.addManualIP(ip)} 
/>
