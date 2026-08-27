<script>
  import { Power, Edit2, Check, Palette, Sparkles, Sliders } from 'lucide-svelte';

  let { 
    device, 
    effects = [], 
    palettes = [], 
    onTogglePower = () => {}, 
    onSetBrightness = () => {}, 
    onSetColor = () => {}, 
    onSetEffect = () => {}, 
    onSetPalette = () => {}, 
    onRename = () => {} 
  } = $props();

  let isEditingName = $state(false);
  let nickname = $state('');
  let selectedColor = $state('#06b6d4');

  $effect(() => {
    nickname = device.name;
  });

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

<div class="glass-panel rounded-2xl p-5 flex flex-col justify-between space-y-4 relative overflow-hidden group border transition-all duration-300 {device.isOnline ? 'border-cyan-500/20 hover:border-purple-500/40' : 'border-rose-500/20 opacity-75'}">
  
  <!-- Top Accent Status Bar -->
  <div class="absolute top-0 left-0 right-0 h-1 bg-gradient-to-r {device.isOnline ? 'from-cyan-500 via-purple-500 to-amber-500' : 'from-rose-600 to-rose-900'}"></div>

  <!-- Header & Name Editor -->
  <div class="flex items-center justify-between pt-1">
    <div class="flex items-center gap-2 flex-1">
      <!-- Status Pill -->
      <span class="w-2.5 h-2.5 rounded-full {device.isOnline ? 'bg-emerald-400 shadow-lg shadow-emerald-500/50 animate-pulse' : 'bg-rose-500'}"></span>

      {#if isEditingName}
        <div class="flex items-center gap-1 flex-1">
          <input 
            type="text" 
            bind:value={nickname}
            class="bg-slate-900 border border-cyan-500/50 rounded-lg px-2 py-0.5 text-xs text-slate-100 font-semibold focus:outline-none w-full"
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
          <button onclick={() => isEditingName = true} class="opacity-0 group-hover/edit:opacity-100 p-0.5 text-slate-500 hover:text-slate-300">
            <Edit2 class="w-3 h-3" />
          </button>
        </div>
      {/if}
    </div>

    <!-- Power Switch -->
    <button 
      onclick={() => onTogglePower(device.id)}
      class="p-2 rounded-xl transition-all duration-200 {device.state?.on ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan' : 'bg-slate-900/80 text-slate-600 border border-slate-800'}"
    >
      <Power class="w-4 h-4" />
    </button>
  </div>

  <p class="text-[11px] font-mono text-slate-400 -mt-2">
    {device.ipAddress} • {device.ledCount} LEDs
  </p>

  <!-- Color Swatches & Effect Controls -->
  <div class="space-y-3 pt-1">
    <!-- Quick Color Swatches -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-1.5">
        {#each presetColors as c}
          <button 
            onclick={() => applyColor(c)}
            style="background-color: {c.hex}"
            class="w-5 h-5 rounded-md border border-slate-900 shadow-sm transition-transform hover:scale-110 focus:outline-none"
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
    <div class="flex items-center gap-3 bg-slate-900/40 px-3 py-1.5 rounded-xl border border-slate-800/80">
      <Sliders class="w-3.5 h-3.5 text-slate-400" />
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

    <!-- WLED FX & Palette Selection -->
    <div class="grid grid-cols-2 gap-2">
      <!-- FX Selector -->
      <div class="flex items-center gap-1.5 bg-slate-900/60 border border-slate-800 px-2 py-1 rounded-xl">
        <Sparkles class="w-3.5 h-3.5 text-cyan-400" />
        <select 
          value={device.state?.seg?.[0]?.fx ?? 0}
          onchange={(e) => onSetEffect(device.id, parseInt(e.target.value))}
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
          value={device.state?.seg?.[0]?.pal ?? 0}
          onchange={(e) => {
            const palId = parseInt(e.target.value);
            onSetPalette(device.id, palId);
            // If current effect is Solid (0), switch to Breathe (2) so palette colors are visible!
            if ((device.state?.seg?.[0]?.fx ?? 0) === 0) {
              onSetEffect(device.id, 2);
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
