<script>
  import { Sparkles, Play, Trash2, Sun, Moon, Film, Flame, Zap, Pin, Maximize2, MoreVertical } from 'lucide-svelte';

  let {
    scene,
    isPinned = true,
    cardSize = 'normal',
    onApply = () => {},
    onDelete = () => {},
    onTogglePin = () => {},
    onToggleSize = () => {}
  } = $props();

  let isApplying = $state(false);
  let isMenuOpen = $state(false);

  async function handleApply() {
    isApplying = true;
    try {
      await onApply(scene.id);
    } finally {
      setTimeout(() => {
        isApplying = false;
      }, 300);
    }
  }

  const iconMap = {
    Sparkles,
    Sun,
    Moon,
    Film,
    Flame,
    Zap
  };

  const SelectedIcon = $derived(iconMap[scene.icon] || Sparkles);
</script>

<div
  class="glass-panel rounded-2xl p-4 flex flex-col justify-between space-y-3 relative group border h-full border-cyan-500/20 hover:border-cyan-500/50 transition-all duration-300"
>
  <!-- Top Accent Status Pill Bar -->
  <div class="w-full h-1 rounded-full bg-gradient-to-r from-cyan-500 via-purple-500 to-amber-500"></div>

  <!-- FUNCTIONAL SCENE HEADER & OPTIONS MENU -->
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-3">
      <div class="p-2.5 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-400">
        <SelectedIcon class="w-5 h-5" />
      </div>
      <div>
        <h4 class="font-semibold text-slate-100 text-sm tracking-wide">{scene.name}</h4>
        <p class="text-[11px] font-mono text-slate-400">Multi-Strip Snapshot</p>
      </div>
    </div>

    <!-- Options Menu & Delete Button -->
    <div class="flex items-center gap-1.5 relative">
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
              onToggleSize(scene.id);
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-cyan-300 hover:bg-cyan-500/15 transition-all text-left cursor-pointer"
          >
            <Maximize2 class="w-3.5 h-3.5 text-cyan-400" />
            <span>{cardSize === 'wide' ? 'Normal Width' : 'Expand Full Width'}</span>
          </button>

          <button
            onclick={() => {
              onTogglePin(scene.id, 'scene');
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-slate-300 hover:text-purple-300 hover:bg-purple-500/15 transition-all text-left cursor-pointer"
          >
            <Pin class="w-3.5 h-3.5 text-purple-400" />
            <span>{isPinned ? 'Unpin from Canvas' : 'Pin to Dashboard'}</span>
          </button>

          <button
            onclick={() => {
              onDelete(scene.id);
              isMenuOpen = false;
            }}
            class="w-full flex items-center gap-2.5 px-3 py-2 rounded-xl text-rose-400 hover:bg-rose-500/15 transition-all text-left cursor-pointer border-t border-slate-800/80 mt-1 pt-1.5"
          >
            <Trash2 class="w-3.5 h-3.5 text-rose-400" />
            <span>Delete Scene</span>
          </button>
        </div>
      {/if}
    </div>
  </div>

  <!-- Bottom Apply Action Bar -->
  <div class="pt-2">
    <button
      onclick={handleApply}
      disabled={isApplying}
      class="w-full flex items-center justify-center gap-2 py-2 px-4 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-medium text-xs shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all duration-200 active:scale-95 disabled:opacity-50 cursor-pointer h-9"
    >
      <Play class="w-4 h-4 fill-current" />
      <span>{isApplying ? 'Applying Scene Snapshot...' : 'Apply Scene Snapshot'}</span>
    </button>
  </div>
</div>
