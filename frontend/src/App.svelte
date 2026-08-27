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
    Plus,
    Camera,
    Play,
    Sparkles,
    Sun,
    Moon,
    Film,
    Flame,
    Maximize2,
    Grid,
    Trash2,
    X,
    PanelTop
  } from 'lucide-svelte';
  import { getDeviceStore } from '$lib/stores/deviceStore.svelte.js';
  import { getGroupStore } from '$lib/stores/groupStore.svelte.js';
  import { getDashboardStore } from '$lib/stores/dashboardStore.svelte.js';
  import DeviceCard from '$lib/components/DeviceCard.svelte';
  import ManualIpModal from '$lib/components/ManualIpModal.svelte';
  import GroupCard from '$lib/components/GroupCard.svelte';
  import SceneCard from '$lib/components/SceneCard.svelte';
  import CreateGroupModal from '$lib/components/CreateGroupModal.svelte';
  import CaptureSceneModal from '$lib/components/CaptureSceneModal.svelte';
  import StripsManagement from '$lib/components/StripsManagement.svelte';

  const store = getDeviceStore();
  const groupStore = getGroupStore();
  const dashboardStore = getDashboardStore();

  let activeTab = $state('dashboard');
  let masterPower = $state(true);
  let masterBrightness = $state(200);
  let isAddModalOpen = $state(false);
  let isCreateGroupModalOpen = $state(false);
  let isCaptureSceneModalOpen = $state(false);
  let isAddPanelModalOpen = $state(false);
  let newPanelTitle = $state('');

  // Drag and Drop active target tracking
  let draggedItemId = $state(null);
  let dragOverPanelId = $state(null);

  const pinnedCards = $derived.by(() => {
    let list = [];

    // Add pinned groups
    groupStore.groups.forEach(g => {
      if (dashboardStore.isPinned(g.id)) {
        list.push({
          id: g.id,
          type: 'group',
          data: g,
          size: dashboardStore.getSize(g.id),
          panelId: dashboardStore.getPanelId(g.id)
        });
      }
    });

    // Add pinned scene presets
    groupStore.scenes.forEach(s => {
      if (dashboardStore.isPinned(s.id)) {
        list.push({
          id: s.id,
          type: 'scene',
          data: s,
          size: dashboardStore.getSize(s.id),
          panelId: dashboardStore.getPanelId(s.id)
        });
      }
    });

    // Add pinned individual devices
    store.devices.forEach(d => {
      if (dashboardStore.isPinned(d.id)) {
        if (!groupStore.activeGroupId || isDeviceInGroup(d.id, groupStore.activeGroupId)) {
          list.push({
            id: d.id,
            type: 'device',
            data: d,
            size: dashboardStore.getSize(d.id),
            panelId: dashboardStore.getPanelId(d.id)
          });
        }
      }
    });

    return list;
  });

  // Cards directly on the main dashboard (unassigned to any custom panel)
  const unassignedCards = $derived(
    pinnedCards.filter(c => !c.panelId)
  );

  // Custom user-created panels
  const customPanels = $derived(
    dashboardStore.panels.filter(p => p.id !== 'default')
  );

  function isDeviceInGroup(deviceId, groupId) {
    const g = groupStore.groups.find(group => group.id === groupId);
    return g?.deviceIds?.includes(deviceId);
  }

  onMount(() => {
    store.init();
    groupStore.init();
    dashboardStore.init();
  });

  function toggleMasterPower() {
    masterPower = !masterPower;
    store.devices.filter(d => d.isOnline).forEach(d => store.togglePower(d.id));
  }

  function handleMasterBrightness(val) {
    masterBrightness = val;
    store.devices.filter(d => d.isOnline).forEach(d => store.setBrightness(d.id, val));
  }

  async function cycleCardSize(id, currentSize) {
    const nextSize = currentSize === 'normal' ? 'wide' : 'normal';
    await dashboardStore.setSize(id, nextSize);
  }

  function getCardSpanClass(size) {
    if (size === 'wide') return 'col-span-1 md:col-span-2 lg:col-span-3';
    return 'col-span-1';
  }

  function handleAddPanel(e) {
    e.preventDefault();
    if (newPanelTitle.trim()) {
      dashboardStore.addPanel(newPanelTitle.trim());
      newPanelTitle = '';
      isAddPanelModalOpen = false;
    }
  }

  // HTML5 Drag and Drop handlers
  function handleDragStart(e, id) {
    if (e.target.closest('button, input, select, textarea, [role="button"]')) {
      e.preventDefault();
      return;
    }
    draggedItemId = id;
    e.dataTransfer.setData('text/plain', id);
    e.dataTransfer.effectAllowed = 'move';
  }

  function handleDragOver(e, panelId) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    dragOverPanelId = panelId;
  }

  function handleDragLeave(panelId) {
    if (dragOverPanelId === panelId) {
      dragOverPanelId = null;
    }
  }

  async function handleDrop(e, targetPanelId) {
    e.preventDefault();
    const id = e.dataTransfer.getData('text/plain') || draggedItemId;
    if (id) {
      await dashboardStore.setPanelId(id, targetPanelId);
    }
    draggedItemId = null;
    dragOverPanelId = null;
  }
</script>

<div class="flex h-screen bg-[#06090e] text-slate-100 font-sans antialiased overflow-hidden selection:bg-cyan-500 selection:text-black">
  
  <!-- CYBERPUNK SIDE NAVIGATION BAR -->
  <aside class="w-64 bg-[#090e17]/80 backdrop-blur-xl border-r border-cyan-500/20 flex flex-col justify-between p-4 z-30 shadow-[4px_0_24px_rgba(0,0,0,0.5)]">
    <div class="space-y-6">
      <!-- App Brand Logo -->
      <div class="flex items-center gap-3 px-2 pt-2">
        <div class="relative flex items-center justify-center w-10 h-10 rounded-xl bg-gradient-to-tr from-cyan-500 to-purple-600 shadow-[0_0_15px_rgba(6,182,212,0.5)]">
          <Zap class="w-6 h-6 text-black fill-current animate-pulse" />
        </div>
        <div>
          <h1 class="font-mono text-base font-bold tracking-wider bg-gradient-to-r from-cyan-400 via-purple-400 to-amber-400 bg-clip-text text-transparent">
            LED SWARM
          </h1>
          <p class="text-[10px] font-mono text-cyan-500/70 tracking-widest uppercase">Orchestrator v0.1</p>
        </div>
      </div>

      <!-- Navigation Menu -->
      <nav class="space-y-1">
        <button 
          onclick={() => activeTab = 'dashboard'}
          class="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-mono text-xs transition-all duration-200 cursor-pointer {activeTab === 'dashboard' ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/40 shadow-[0_0_15px_rgba(6,182,212,0.2)]' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'}"
        >
          <LayoutGrid class="w-4 h-4" />
          <span>Dashboard Canvas</span>
        </button>

        <button 
          onclick={() => activeTab = 'strips'}
          class="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-mono text-xs transition-all duration-200 cursor-pointer {activeTab === 'strips' ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/40 shadow-[0_0_15px_rgba(6,182,212,0.2)]' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'}"
        >
          <Sliders class="w-4 h-4" />
          <span>Strips & Devices</span>
        </button>

        <button 
          onclick={() => activeTab = 'canvas'}
          class="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-mono text-xs transition-all duration-200 cursor-pointer {activeTab === 'canvas' ? 'bg-cyan-500/10 text-cyan-400 border border-cyan-500/40 shadow-[0_0_15px_rgba(6,182,212,0.2)]' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'}"
        >
          <Grid class="w-4 h-4" />
          <span>2D Room Canvas</span>
        </button>

        <button 
          onclick={() => activeTab = 'groups'}
          class="w-full flex items-center gap-3 px-3.5 py-2.5 rounded-xl font-mono text-xs transition-all duration-200 cursor-pointer {activeTab === 'groups' ? 'bg-purple-500/10 text-purple-400 border border-purple-500/40 shadow-[0_0_15px_rgba(168,85,247,0.2)]' : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'}"
        >
          <Layers class="w-4 h-4" />
          <span>Groups & Scenes</span>
        </button>
      </nav>
    </div>

    <!-- Active Strips Telemetry Footer -->
    <div class="glass-panel p-3.5 rounded-2xl space-y-2 border border-slate-800">
      <div class="flex items-center justify-between text-xs font-mono">
        <span class="text-slate-400 flex items-center gap-1.5">
          <Activity class="w-3.5 h-3.5 text-emerald-400" />
          Active Swarm
        </span>
        <span class="text-cyan-400 font-bold">{(store.onlineDevices?.length || 0)} / {(store.devices?.length || 0)}</span>
      </div>
      <div class="w-full bg-slate-900 h-1.5 rounded-full overflow-hidden border border-slate-800">
        <div 
          class="h-full bg-gradient-to-r from-cyan-500 to-purple-500 transition-all duration-500"
          style="width: {store.devices?.length ? ((store.onlineDevices?.length || 0) / store.devices.length) * 100 : 0}%"
        ></div>
      </div>
    </div>
  </aside>

  <!-- MAIN APPLICATION CONTENT WORKSPACE -->
  <main class="flex-1 flex flex-col min-w-0 overflow-hidden bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-slate-900/40 via-[#06090e] to-[#06090e]">
    
    <!-- TOP CONTROL BAR -->
    <header class="h-16 border-b border-cyan-500/10 bg-[#090e17]/60 backdrop-blur-md px-6 flex items-center justify-between z-20">
      <!-- Search & Discovery Actions -->
      <div class="flex items-center gap-3">
        <button 
          onclick={() => store.discover()}
          disabled={store.isScanning}
          class="flex items-center gap-2 px-3.5 py-1.5 rounded-xl bg-slate-900/80 hover:bg-cyan-500/10 border border-cyan-500/30 text-cyan-400 font-mono text-xs transition-all duration-200 disabled:opacity-50 cursor-pointer"
        >
          <RefreshCw class="w-3.5 h-3.5 {store.isScanning ? 'animate-spin' : ''}" />
          <span>{store.isScanning ? 'Scanning mDNS...' : 'Scan WLED Strips'}</span>
        </button>

        <button 
          onclick={() => isAddModalOpen = true}
          class="flex items-center gap-2 px-3.5 py-1.5 rounded-xl bg-slate-900/80 hover:bg-purple-500/10 border border-purple-500/30 text-purple-400 font-mono text-xs transition-all duration-200 cursor-pointer"
        >
          <Plus class="w-3.5 h-3.5" />
          <span>Manual IP</span>
        </button>
      </div>

      <!-- Global Master Hardware Control Bar -->
      <div class="flex items-center gap-6 glass-panel px-4 py-2 rounded-2xl border border-cyan-500/20">
        <!-- Master Power Switch -->
        <button 
          onclick={toggleMasterPower}
          class="flex items-center gap-2 px-3 py-1 rounded-xl transition-all duration-200 font-mono text-xs cursor-pointer {masterPower ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/50 glow-cyan' : 'bg-slate-900 text-slate-500 border border-slate-800'}"
        >
          <Power class="w-4 h-4" />
          <span>MASTER {masterPower ? 'ON' : 'OFF'}</span>
        </button>

        <!-- Master Brightness Slider -->
        <div class="flex items-center gap-3 w-48">
          <Sun class="w-4 h-4 text-purple-400" />
          <input 
            type="range" 
            min="0" 
            max="255" 
            value={masterBrightness}
            oninput={(e) => handleMasterBrightness(parseInt(e.target.value))}
            class="w-full accent-cyan-400 cursor-pointer"
          />
          <span class="text-xs font-mono text-cyan-400 w-8 text-right">{masterBrightness}</span>
        </div>
      </div>
    </header>

    <!-- CONTENT BODY VIEW CONTAINER -->
    <div class="flex-1 overflow-y-auto p-6 space-y-6 scrollbar-thin">
      
      {#if activeTab === 'dashboard'}
        <!-- CUSTOMIZABLE DASHBOARD CANVAS VIEW -->
        <div class="space-y-6">
          <!-- Quick Action Controls & Header -->
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold font-mono tracking-wide text-slate-100">DASHBOARD CANVAS</h2>
              <p class="text-xs font-mono text-slate-400">Drag and drop cards into custom section panels or position them directly on the main canvas.</p>
            </div>

            <!-- Add Custom Panel Trigger Button -->
            <button 
              onclick={() => isAddPanelModalOpen = true}
              class="flex items-center gap-2 px-4 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-cyan-500 hover:from-purple-500 hover:to-cyan-400 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(168,85,247,0.3)] transition-all duration-200 cursor-pointer"
            >
              <Plus class="w-4 h-4" />
              <span>Add Dashboard Panel</span>
            </button>
          </div>

          <!-- 1. MAIN DASHBOARD CANVAS DROP ZONE (Cards sit directly on Dashboard!) -->
          <div 
            role="region"
            aria-label="Main Dashboard Canvas"
            ondragover={(e) => handleDragOver(e, '')}
            ondragleave={() => handleDragLeave('')}
            ondrop={(e) => handleDrop(e, '')}
            class="space-y-3 rounded-3xl p-4 transition-all duration-300 border border-transparent {dragOverPanelId === '' ? 'border-cyan-400/80 bg-cyan-500/10 shadow-[0_0_30px_rgba(6,182,212,0.3)]' : ''}"
          >
            <h2 class="text-xs font-mono tracking-wider text-slate-400 flex items-center gap-2">
              <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse"></span>
              DASHBOARD CANVAS ({unassignedCards.length})
            </h2>

            {#if pinnedCards.length === 0}
              <div class="glass-panel rounded-3xl p-12 text-center border border-dashed border-cyan-500/20">
                <Wifi class="w-12 h-12 text-cyan-500/40 mx-auto mb-3 animate-pulse" />
                <h3 class="text-base font-bold text-slate-200">No Items Pinned To Dashboard Canvas</h3>
                <p class="text-xs font-mono text-slate-500 mt-1">Go to Settings -> Strips & Devices or Groups & Scenes to pin cards directly to your dashboard.</p>
              </div>
            {:else if unassignedCards.length === 0}
              <div class="glass-panel rounded-2xl p-6 text-center border border-dashed border-slate-800">
                <p class="text-xs font-mono text-slate-500">No unassigned cards on main canvas. Drag cards here from any panel to move them back to the main dashboard.</p>
              </div>
            {:else}
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 items-stretch">
                {#each unassignedCards as item (item.id)}
                  <div 
                    role="group"
                    aria-label="Draggable card"
                    draggable="true"
                    ondragstart={(e) => handleDragStart(e, item.id)}
                    class="{getCardSpanClass(item.size)} transition-all duration-300 h-full cursor-grab active:cursor-grabbing"
                  >
                    {#if item.type === 'device'}
                      <DeviceCard 
                        device={item.data}
                        effects={store.effects}
                        palettes={store.palettes}
                        isPinned={true}
                        cardSize={item.size}
                        onTogglePower={(id) => store.togglePower(id)}
                        onSetBrightness={(id, bri) => store.setBrightness(id, bri)}
                        onSetColor={(id, r, g, b) => store.setColor(id, r, g, b)}
                        onSetEffect={(id, fx) => store.setEffect(id, fx)}
                        onSetPalette={(id, pal) => store.setPalette(id, pal)}
                        onRename={(id, name) => store.renameDevice(id, name)}
                        onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                        onToggleSize={(id) => cycleCardSize(id, item.size)}
                      />
                    {:else if item.type === 'group'}
                      <GroupCard 
                        group={item.data}
                        allDevices={store.devices}
                        effects={store.effects}
                        palettes={store.palettes}
                        isPinned={true}
                        cardSize={item.size}
                        onTogglePower={(id, pwr) => groupStore.setGroupState(id, pwr ? { on: true, bri: 200, mainseg: 0, seg: [{ id: 0, on: true }] } : { on: false })}
                        onSetBrightness={(id, bri) => groupStore.setGroupState(id, { bri })}
                        onSetColor={(id, r, g, b) => groupStore.setGroupState(id, { seg: [{ id: 0, col: [[r, g, b]] }] })}
                        onSetEffect={(id, fx) => groupStore.setGroupState(id, { seg: [{ id: 0, fx }] })}
                        onSetPalette={(id, pal) => groupStore.setGroupState(id, { seg: [{ id: 0, pal, fx: 2 }] })}
                        onDelete={(id) => groupStore.deleteGroup(id)}
                        onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                        onToggleSize={(id) => cycleCardSize(id, item.size)}
                      />
                    {:else if item.type === 'scene'}
                      <SceneCard 
                        scene={item.data}
                        isPinned={true}
                        cardSize={item.size}
                        onApply={(id) => groupStore.applyScene(id)}
                        onDelete={(id) => groupStore.deleteScene(id)}
                        onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                        onToggleSize={(id) => cycleCardSize(id, item.size)}
                      />
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          </div>

          <!-- 2. OPTIONAL CUSTOM DASHBOARD PANELS (DROP TARGETS) -->
          {#each customPanels as panel (panel.id)}
            {@const panelCards = pinnedCards.filter(c => c.panelId === panel.id)}
            <div 
              role="region"
              aria-label="{panel.title} Panel"
              ondragover={(e) => handleDragOver(e, panel.id)}
              ondragleave={() => handleDragLeave(panel.id)}
              ondrop={(e) => handleDrop(e, panel.id)}
              class="glass-panel rounded-3xl p-6 border transition-all duration-300 space-y-4 {dragOverPanelId === panel.id ? 'border-purple-400/80 bg-purple-500/10 shadow-[0_0_30px_rgba(168,85,247,0.3)]' : 'border-slate-800/80'}"
            >
              <!-- Panel Header -->
              <div class="flex items-center justify-between border-b border-slate-800/60 pb-3">
                <h2 class="text-sm font-mono tracking-wider text-slate-200 flex items-center gap-2">
                  <PanelTop class="w-4 h-4 text-purple-400" />
                  {panel.title.toUpperCase()} ({panelCards.length})
                </h2>
                <button 
                  onclick={() => dashboardStore.deletePanel(panel.id)}
                  class="p-1 text-slate-500 hover:text-rose-400 transition-colors cursor-pointer"
                  title="Remove Panel"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>

              <!-- Panel Grid Cards -->
              {#if panelCards.length === 0}
                <div class="text-center py-8 border border-dashed border-slate-800/80 rounded-2xl">
                  <p class="text-xs font-mono text-slate-500">Drag and drop cards here to move them into this panel.</p>
                </div>
              {:else}
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 items-stretch">
                  {#each panelCards as item (item.id)}
                    <div 
                      role="group"
                      aria-label="Draggable card"
                      draggable="true"
                      ondragstart={(e) => handleDragStart(e, item.id)}
                      class="{getCardSpanClass(item.size)} transition-all duration-300 h-full cursor-grab active:cursor-grabbing"
                    >
                      {#if item.type === 'device'}
                        <DeviceCard 
                          device={item.data}
                          effects={store.effects}
                          palettes={store.palettes}
                          isPinned={true}
                          cardSize={item.size}
                          onTogglePower={(id) => store.togglePower(id)}
                          onSetBrightness={(id, bri) => store.setBrightness(id, bri)}
                          onSetColor={(id, r, g, b) => store.setColor(id, r, g, b)}
                          onSetEffect={(id, fx) => store.setEffect(id, fx)}
                          onSetPalette={(id, pal) => store.setPalette(id, pal)}
                          onRename={(id, name) => store.renameDevice(id, name)}
                          onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                          onToggleSize={(id) => cycleCardSize(id, item.size)}
                        />
                      {:else if item.type === 'group'}
                        <GroupCard 
                          group={item.data}
                          allDevices={store.devices}
                          effects={store.effects}
                          palettes={store.palettes}
                          isPinned={true}
                          cardSize={item.size}
                          onTogglePower={(id, pwr) => groupStore.setGroupState(id, pwr ? { on: true, bri: 200, mainseg: 0, seg: [{ id: 0, on: true }] } : { on: false })}
                          onSetBrightness={(id, bri) => groupStore.setGroupState(id, { bri })}
                          onSetColor={(id, r, g, b) => groupStore.setGroupState(id, { seg: [{ id: 0, col: [[r, g, b]] }] })}
                          onSetEffect={(id, fx) => groupStore.setGroupState(id, { seg: [{ id: 0, fx }] })}
                          onSetPalette={(id, pal) => groupStore.setGroupState(id, { seg: [{ id: 0, pal, fx: 2 }] })}
                          onDelete={(id) => groupStore.deleteGroup(id)}
                          onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                          onToggleSize={(id) => cycleCardSize(id, item.size)}
                        />
                      {:else if item.type === 'scene'}
                        <SceneCard 
                          scene={item.data}
                          isPinned={true}
                          cardSize={item.size}
                          onApply={(id) => groupStore.applyScene(id)}
                          onDelete={(id) => groupStore.deleteScene(id)}
                          onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                          onToggleSize={(id) => cycleCardSize(id, item.size)}
                        />
                      {/if}
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {:else if activeTab === 'strips'}
        <!-- STRIPS MANAGEMENT TAB -->
        <StripsManagement 
          devices={store.devices}
          isScanning={store.isScanning}
          dashboardStore={dashboardStore}
          onDiscover={() => store.discover()}
          onAddManual={() => isAddModalOpen = true}
          onRename={(id, name) => store.renameDevice(id, name)}
        />
      {:else if activeTab === 'canvas'}
        <!-- 2D ROOM VISUAL CANVAS PLACEHOLDER -->
        <div class="glass-panel rounded-3xl p-12 text-center border border-dashed border-cyan-500/20 space-y-4">
          <Grid class="w-12 h-12 text-cyan-400/40 mx-auto animate-pulse" />
          <h3 class="text-lg font-bold font-mono text-slate-200">2D Room Canvas Workspace</h3>
          <p class="text-xs font-mono text-slate-400 max-w-md mx-auto">Visual 2D room canvas mapping for drag-and-drop spatial LED positioning, spatial light effects, and live pixel mirroring.</p>
        </div>
      {:else if activeTab === 'groups'}
        <!-- GROUPS & SCENES MANAGEMENT WORKSPACE -->
        <div class="space-y-6">
          <!-- Groups Header & Actions -->
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-xl font-bold font-mono tracking-wide text-slate-100">GROUPS & MULTI-ZONE SCENES</h2>
              <p class="text-xs font-mono text-slate-400">Combine multiple WLED strips into virtual synchronized groups and save multi-zone scene presets.</p>
            </div>
            
            <div class="flex items-center gap-3">
              <button 
                onclick={() => isCaptureSceneModalOpen = true}
                class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-cyan-500/10 hover:bg-cyan-500/20 border border-cyan-500/40 text-cyan-300 font-mono text-xs font-semibold transition-all cursor-pointer"
              >
                <Camera class="w-4 h-4" />
                <span>Capture Current Scene</span>
              </button>

              <button 
                onclick={() => isCreateGroupModalOpen = true}
                class="flex items-center gap-2 px-3.5 py-2 rounded-xl bg-purple-500/20 hover:bg-purple-500/30 border border-purple-500/50 text-purple-300 font-mono text-xs font-semibold transition-all cursor-pointer shadow-[0_0_15px_rgba(168,85,247,0.3)]"
              >
                <Plus class="w-4 h-4" />
                <span>Create Strip Group</span>
              </button>
            </div>
          </div>

          <!-- Scenes Grid Section -->
          {#if groupStore.scenes.length > 0}
            <div class="space-y-3">
              <h3 class="text-xs font-mono text-slate-400 uppercase tracking-wider flex items-center gap-2">
                <Sparkles class="w-3.5 h-3.5 text-cyan-400" />
                Saved Multi-Strip Scene Snapshots ({groupStore.scenes.length})
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {#each groupStore.scenes as scene (scene.id)}
                  <SceneCard 
                    scene={scene}
                    isPinned={dashboardStore.isPinned(scene.id)}
                    cardSize={dashboardStore.getSize(scene.id)}
                    onApply={(id) => groupStore.applyScene(id)}
                    onDelete={(id) => groupStore.deleteScene(id)}
                    onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                    onToggleSize={(id) => cycleCardSize(id, dashboardStore.getSize(scene.id))}
                  />
                {/each}
              </div>
            </div>
          {/if}

          <!-- Groups Grid Section -->
          <div class="space-y-3">
            <h3 class="text-xs font-mono text-slate-400 uppercase tracking-wider flex items-center gap-2">
              <Layers class="w-3.5 h-3.5 text-purple-400" />
              Virtual WLED Strip Groups ({groupStore.groups.length})
            </h3>
            
            {#if groupStore.groups.length === 0}
              <div class="glass-panel rounded-3xl p-12 text-center border border-dashed border-purple-500/20">
                <Layers class="w-12 h-12 text-purple-500/40 mx-auto mb-3 animate-pulse" />
                <h3 class="text-base font-bold text-slate-200">No Virtual Groups Created</h3>
                <p class="text-xs font-mono text-slate-500 mt-1">Click "Create Strip Group" to group WLED light strips together into synchronized zones.</p>
              </div>
            {:else}
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {#each groupStore.groups as group (group.id)}
                  <GroupCard 
                    group={group}
                    allDevices={store.devices}
                    effects={store.effects}
                    palettes={store.palettes}
                    isPinned={dashboardStore.isPinned(group.id)}
                    cardSize={dashboardStore.getSize(group.id)}
                    onTogglePower={(id, pwr) => groupStore.setGroupState(id, pwr ? { on: true, bri: 200, mainseg: 0, seg: [{ id: 0, on: true }] } : { on: false })}
                    onSetBrightness={(id, bri) => groupStore.setGroupState(id, { bri })}
                    onSetColor={(id, r, g, b) => groupStore.setGroupState(id, { seg: [{ id: 0, col: [[r, g, b]] }] })}
                    onSetEffect={(id, fx) => groupStore.setGroupState(id, { seg: [{ id: 0, fx }] })}
                    onSetPalette={(id, pal) => groupStore.setGroupState(id, { seg: [{ id: 0, pal, fx: 2 }] })}
                    onDelete={(id) => groupStore.deleteGroup(id)}
                    onTogglePin={(id, type) => dashboardStore.togglePin(id, type)}
                    onToggleSize={(id) => cycleCardSize(id, dashboardStore.getSize(group.id))}
                  />
                {/each}
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </main>
</div>

<!-- MODAL DIALOGS -->
<ManualIpModal 
  isOpen={isAddModalOpen}
  onClose={() => isAddModalOpen = false}
  onAddDevice={async (ip) => {
    await store.addManualDevice(ip);
    isAddModalOpen = false;
  }}
/>

<CreateGroupModal 
  isOpen={isCreateGroupModalOpen}
  allDevices={store.devices}
  onClose={() => isCreateGroupModalOpen = false}
  onCreate={(name, desc, devIds) => {
    groupStore.createGroup(name, desc, devIds);
    isCreateGroupModalOpen = false;
  }}
/>

<CaptureSceneModal 
  isOpen={isCaptureSceneModalOpen}
  onClose={() => isCaptureSceneModalOpen = false}
  onCapture={(name, icon) => {
    groupStore.captureScene(name, icon);
    isCaptureSceneModalOpen = false;
  }}
/>

<!-- ADD DASHBOARD PANEL MODAL -->
{#if isAddPanelModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
    <div class="glass-panel w-full max-w-md rounded-2xl p-6 border border-cyan-500/40 shadow-2xl space-y-4">
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <h3 class="font-mono text-base font-bold text-slate-100 flex items-center gap-2">
          <PanelTop class="w-5 h-5 text-cyan-400" />
          Add Custom Panel Section
        </h3>
        <button 
          onclick={() => isAddPanelModalOpen = false}
          class="p-1.5 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <form onsubmit={handleAddPanel} class="space-y-4">
        <div>
          <label for="panel-title-input" class="block text-xs font-mono text-slate-400 mb-1">Panel Section Title</label>
          <input 
            id="panel-title-input"
            type="text" 
            placeholder="e.g. Living Room Zone, Office Setup" 
            bind:value={newPanelTitle}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-cyan-500/60 placeholder-slate-600"
          />
        </div>

        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-800">
          <button 
            type="button"
            onclick={() => isAddPanelModalOpen = false}
            class="px-4 py-2 rounded-xl text-xs font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
          >
            Cancel
          </button>
          <button 
            type="submit"
            class="px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all duration-200"
          >
            Create Panel
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
