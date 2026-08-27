<script>
  import { Power, Sun, Palette, Sparkles, Trash2, Layers, Cpu } from 'lucide-svelte';

  let { 
    group, 
    allDevices = [],
    effects = [], 
    palettes = [],
    onTogglePower = () => {},
    onSetBrightness = () => {},
    onSetColor = () => {},
    onSetEffect = () => {},
    onSetPalette = () => {},
    onDelete = () => {} 
  } = $props();

  let groupPower = $state(true);
  let groupBrightness = $state(180);
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

<div class="glass-panel rounded-2xl p-5 flex flex-col justify-between space-y-4 relative overflow-hidden group border border-purple-500/20 hover:border-purple-500/40 transition-all duration-300">
  
  <!-- Top Accent Status Bar -->
  <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-purple-500 via-cyan-500 to-amber-500"></div>

  <!-- Group Header -->
  <div class="flex items-center justify-between pt-1 pb-2 border-b border-slate-800/80">
    <div class="flex items-center gap-3">
      <div class="p-2.5 rounded-xl bg-purple-500/10 border border-purple-500/30 text-purple-400">
        <Layers class="w-5 h-5" />
      </div>
      <div>
        <h3 class="font-semibold text-slate-100 text-sm tracking-wide">{group.name}</h3>
        <p class="text-[11px] font-mono text-slate-400">
          {group.deviceIds?.length ?? 0} Strips ({onlineCount} Online)
        </p>
      </div>
    </div>

    <!-- Power & Delete Controls -->
    <div class="flex items-center gap-2">
      <button 
        onclick={() => {
          groupPower = !groupPower;
          onTogglePower(group.id, groupPower);
        }}
        class="p-2 rounded-xl transition-all duration-200 {groupPower ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan' : 'bg-slate-900/80 text-slate-600 border border-slate-800'}"
        title="Toggle Group Power"
      >
        <Power class="w-4 h-4" />
      </button>

      <button 
        onclick={() => onDelete(group.id)}
        class="p-2 rounded-xl bg-slate-900/80 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-slate-800 hover:border-rose-500/40 transition-all duration-200 cursor-pointer"
        title="Delete Group"
      >
        <Trash2 class="w-4 h-4" />
      </button>
    </div>
  </div>

  <!-- Group Strip Members Badges -->
  <div class="flex flex-wrap gap-1.5">
    {#each assignedDevices as dev}
      <span class="inline-flex items-center gap-1 text-[10px] font-mono px-2 py-0.5 rounded-md border {dev.isOnline ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400' : 'bg-slate-800/50 border-slate-700 text-slate-500'}">
        <Cpu class="w-3 h-3" />
        {dev.name}
      </span>
    {/each}
  </div>

  <!-- Controls Grid -->
  <div class="space-y-3 pt-1">
    <!-- Quick Color Swatches & Indicator -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-1.5">
        {#each presetColors as c}
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
    <div class="flex items-center gap-3 bg-slate-900/40 px-3 py-1.5 rounded-xl border border-slate-800/80">
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

    <!-- Group FX & Palette Selectors -->
    <div class="grid grid-cols-2 gap-2">
      <!-- FX Selector -->
      <div class="flex items-center gap-1.5 bg-slate-900/60 border border-slate-800 px-2 py-1 rounded-xl">
        <Sparkles class="w-3.5 h-3.5 text-cyan-400" />
        <select 
          value={selectedEffect}
          onchange={(e) => {
            selectedEffect = parseInt(e.target.value);
            onSetEffect(group.id, selectedEffect);
          }}
          class="bg-transparent text-[11px] font-mono text-slate-300 w-full focus:outline-none cursor-pointer"
        >
          {#if effects.length === 0}
            <option value="0" class="bg-slate-900 text-slate-200">Solid</option>
          {/if}
          {#each effects as fxName, index}
            <option value={index} class="bg-slate-900 text-slate-200">{fxName}</option>
          {/each}
        </select>
      </div>

      <!-- Palette Selector -->
      <div class="flex items-center gap-1.5 bg-slate-900/60 border border-slate-800 px-2 py-1 rounded-xl">
        <Palette class="w-3.5 h-3.5 text-purple-400" />
        <select 
          value={selectedPalette}
          onchange={(e) => {
            selectedPalette = parseInt(e.target.value);
            onSetPalette(group.id, selectedPalette);
            if (selectedEffect === 0) {
              selectedEffect = 2; // Breathe
              onSetEffect(group.id, 2);
            }
          }}
          class="bg-transparent text-[11px] font-mono text-slate-300 w-full focus:outline-none cursor-pointer"
        >
          {#if palettes.length === 0}
            <option value="0" class="bg-slate-900 text-slate-200">Default Palette</option>
          {/if}
          {#each palettes as palName, index}
            <option value={index} class="bg-slate-900 text-slate-200">{palName}</option>
          {/each}
        </select>
      </div>
    </div>
  </div>
</div>
