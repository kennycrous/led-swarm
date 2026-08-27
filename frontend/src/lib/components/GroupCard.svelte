<script>
  import { Power, Sun, Palette, Sparkles, Trash2, Layers, Cpu, Pin, Maximize2, MoreVertical } from 'lucide-svelte';
  import CyberSelect from './CyberSelect.svelte';

  let { 
    group, 
    allDevices = [],
    effects = [], 
    palettes = [],
    isPinned = true,
    cardSize = 'normal',
    onTogglePower = () => {},
    onSetBrightness = () => {},
    onSetColor = () => {},
    onSetEffect = () => {},
    onSetPalette = () => {},
    onDelete = () => {},
    onTogglePin = () => {},
    onToggleSize = () => {}
  } = $props();

  let groupPower = $state(true);
  let groupBrightness = $state(180);
  let selectedColor = $state('#06b6d4');
  let selectedEffect = $state(0);
  let selectedPalette = $state(0);
  let isMenuOpen = $state(false);

  const presetColors = [
    { name: 'Cyan Neon', hex: '#06b6d4', r: 6, g: 182, b: 212 },
    { name: 'Magenta Glow', hex: '#a855f7', r: 168, g: 85, b: 247 },
    { name: 'Cyber Amber', hex: '#f59e0b', r: 245, g: 158, b: 11 },
    { name: 'Void White', hex: '#ffffff', r: 255, g: 255, b: 255 },
    { name: 'Emerald Online', hex: '#10b981', r: 16, g: 185, b: 129 },
    { name: 'Crimson Red', hex: '#ef4444', r: 239, g: 68, b: 68 }
  ];

  const assignedDevices = $derived(
    allDevices.filter(d => group.deviceIds?.includes(d.id))
  );

  const onlineCount = $derived(
    assignedDevices.filter(d => d.isOnline).length
  );

  function applyColor(c) {
    selectedColor = c.hex;
    onSetColor(group.id, c.r, c.g, c.b);
  }
</script>

<div class="glass-panel rounded-2xl p-4 flex flex-col justify-between space-y-3 relative group border h-full transition-all duration-300 {isMenuOpen ? 'z-50' : 'hover:z-30 z-10'} border-purple-500/20 hover:border-purple-500/40">
  
  <!-- Top Accent Status Pill Bar -->
  <div class="w-full h-1 rounded-full bg-gradient-to-r from-purple-500 via-cyan-500 to-amber-500"></div>

  <!-- FUNCTIONAL GROUP HEADER & OPTIONS MENU -->
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-2 flex-1">
      <div class="p-2 rounded-xl bg-purple-500/10 border border-purple-500/30 text-purple-400">
        <Layers class="w-4 h-4" />
      </div>
      <div>
        <h3 class="font-semibold text-slate-100 text-sm tracking-wide">{group.name}</h3>
        <p class="text-[11px] font-mono text-slate-400">
          {group.deviceIds?.length ?? 0} Strips ({onlineCount} Online)
        </p>
      </div>
    </div>

    <!-- Power, Delete & Card Options Menu -->
    <div class="flex items-center gap-1.5 relative">
      <button 
        onclick={() => {
          groupPower = !groupPower;
          onTogglePower(group.id, groupPower);
        }}
        class="p-2.5 rounded-xl transition-all duration-200 {groupPower ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan' : 'bg-[#090e17]/80 text-slate-600 border border-slate-800'} cursor-pointer"
        title="Toggle Group Power"
      >
        <Power class="w-4 h-4" />
      </button>

      <!-- Card Options Dropdown Trigger -->
      <button 
        onclick={() => isMenuOpen = !isMenuOpen}
        class="p-2.5 rounded-xl bg-[#090e17]/80 hover:bg-slate-800 text-slate-400 hover:text-slate-200 border border-slate-800 transition-all cursor-pointer"
        title="Card Options"
      >
        <MoreVertical class="w-4 h-4" />
      </button>

      <!-- Glassmorphic Cyber Context Dropdown Menu -->
      {#if isMenuOpen}
        <div 
          class="absolute top-12 right-0 z-50 w-48 bg-[#090e17]/95 border border-purple-500/30 rounded-2xl p-1.5 shadow-[0_0_25px_rgba(168,85,247,0.25)] backdrop-blur-xl space-y-1 text-xs font-mono"
        >
          <button 
            onclick={() => {
              onToggleSize(group.id);
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-purple-300 hover:bg-purple-500/15 transition-all text-left cursor-pointer"
          >
            <Maximize2 class="w-3.5 h-3.5 text-purple-400" />
            <span>{cardSize === 'wide' ? 'Normal Width' : 'Expand Full Width'}</span>
          </button>

          <button 
            onclick={() => {
              onTogglePin(group.id, 'group');
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-purple-300 hover:bg-purple-500/15 transition-all text-left cursor-pointer"
          >
            <Pin class="w-3.5 h-3.5 text-purple-400" />
            <span>{isPinned ? 'Unpin from Canvas' : 'Pin to Dashboard'}</span>
          </button>

          <button 
            onclick={() => {
              onDelete(group.id);
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-rose-400 hover:bg-rose-500/15 transition-all text-left cursor-pointer border-t border-slate-800/80 mt-1 pt-1.5"
          >
            <Trash2 class="w-3.5 h-3.5 text-rose-400" />
            <span>Delete Group</span>
          </button>
        </div>
      {/if}
    </div>
  </div>

  <!-- Group Strip Members Badges -->
  <div class="flex flex-wrap gap-1.5 -mt-1">
    {#each assignedDevices as dev (dev.id)}
      <span class="inline-flex items-center gap-1 text-[10px] font-mono px-2 py-0.5 rounded-md border {dev.isOnline ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400' : 'bg-slate-800/50 border-slate-700 text-slate-500'}">
        <Cpu class="w-3 h-3" />
        {dev.name}
      </span>
    {/each}
  </div>

  <!-- Controls Area -->
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

    <!-- Group Brightness Slider -->
    <div class="flex items-center gap-3 bg-[#090e17]/60 px-3 py-1.5 rounded-xl border border-slate-800/80 h-9">
      <Sun class="w-4 h-4 text-purple-400" />
      <input 
        type="range" 
        min="0" 
        max="255" 
        value={groupBrightness}
        oninput={(e) => {
          groupBrightness = parseInt(e.target.value);
          onSetBrightness(group.id, groupBrightness);
        }}
        disabled={!groupPower}
        class="flex-1 accent-purple-400 cursor-pointer disabled:opacity-30"
      />
      <span class="text-xs font-mono text-purple-300 w-8 text-right">
        {groupBrightness}
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
        onChange={(fxId) => {
          selectedEffect = fxId;
          onSetEffect(group.id, fxId);
        }}
      />

      <!-- Palette CyberSelect -->
      <CyberSelect 
        value={selectedPalette}
        options={palettes}
        icon={Palette}
        iconColor="text-purple-400"
        hoverBorder="hover:border-purple-500/40"
        onChange={(palId) => {
          selectedPalette = palId;
          onSetPalette(group.id, palId);
          if (selectedEffect === 0) {
            selectedEffect = 2; // Breathe
            onSetEffect(group.id, 2);
          }
        }}
      />
    </div>
  </div>
</div>
