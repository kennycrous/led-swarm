<script>
  import { Power, Edit2, Check, Palette, Sparkles, Sun, Pin, Maximize2, Cpu, MoreVertical } from 'lucide-svelte';
  import CyberSelect from './CyberSelect.svelte';

  let { 
    device, 
    effects = [], 
    palettes = [], 
    isPinned = true,
    cardSize = 'normal',
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

<div class="glass-panel rounded-2xl p-4 flex flex-col justify-between space-y-3 relative group border h-full transition-all duration-300 {isMenuOpen ? 'z-50' : 'hover:z-30 z-10'} {device.isOnline ? 'border-cyan-500/20 hover:border-cyan-500/50' : 'border-rose-500/20 opacity-75'}">
  
  <!-- Top Accent Status Pill Bar (Sits cleanly inside card padding) -->
  <div class="w-full h-1 rounded-full bg-gradient-to-r {device.isOnline ? 'from-cyan-500 via-purple-500 to-amber-500' : 'from-rose-600 to-rose-900'}"></div>

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
            <button onclick={() => isEditingName = true} class="opacity-0 group-hover/edit:opacity-100 p-0.5 text-slate-500 hover:text-slate-300">
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
        class="p-2.5 rounded-xl transition-all duration-200 {device.state?.on ? 'bg-cyan-500/20 text-cyan-400 border border-cyan-500/40 glow-cyan' : 'bg-[#090e17]/80 text-slate-600 border border-slate-800'} cursor-pointer"
        title="Toggle Power"
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
          class="absolute top-12 right-0 z-50 w-48 bg-[#090e17]/95 border border-cyan-500/30 rounded-2xl p-1.5 shadow-[0_0_25px_rgba(6,182,212,0.25)] backdrop-blur-xl space-y-1 text-xs font-mono"
        >
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

  <!-- Color Swatches & Controls Area -->
  <div class="space-y-3 pt-1">
    <!-- Quick Color Swatches -->
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

    <!-- Custom Glassmorphic CyberSelect Dropdowns -->
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
  </div>
</div>
