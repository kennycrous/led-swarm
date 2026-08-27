<script>
  import { Sparkles, Play, Trash2, Sun, Moon, Film, Flame, Zap } from 'lucide-svelte';

  let { 
    scene, 
    onApply = () => {}, 
    onDelete = () => {} 
  } = $props();

  let isApplying = $state(false);

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

<div class="relative group bg-slate-900/60 backdrop-blur-xl border border-cyan-500/20 hover:border-cyan-500/50 rounded-2xl p-4 shadow-xl transition-all duration-300">
  
  <div class="flex items-center justify-between">
    <div class="flex items-center gap-3">
      <div class="p-3 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-400">
        <SelectedIcon class="w-6 h-6" />
      </div>
      <div>
        <h4 class="font-semibold text-slate-100 text-sm tracking-wide">{scene.name}</h4>
        <p class="text-[11px] font-mono text-slate-400">Multi-Strip Snapshot</p>
      </div>
    </div>

    <div class="flex items-center gap-2">
      <!-- Apply Scene Button -->
      <button 
        onclick={handleApply}
        disabled={isApplying}
        class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-medium text-xs shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all duration-200 active:scale-95 disabled:opacity-50 cursor-pointer"
      >
        <Play class="w-3.5 h-3.5 fill-current" />
        <span>{isApplying ? 'Applying...' : 'Apply'}</span>
      </button>

      <!-- Delete Button -->
      <button 
        onclick={() => onDelete(scene.id)}
        class="p-2 rounded-xl bg-slate-800/60 hover:bg-rose-500/20 text-slate-400 hover:text-rose-400 border border-slate-700 hover:border-rose-500/40 transition-all duration-200 cursor-pointer"
        title="Delete Scene"
      >
        <Trash2 class="w-4 h-4" />
      </button>
    </div>
  </div>
</div>
