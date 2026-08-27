<script>
  import { X, Sparkles, Sun, Moon, Film, Flame, Zap, Camera } from 'lucide-svelte';

  let { 
    isOpen = false, 
    onClose = () => {}, 
    onCapture = (name, icon) => {} 
  } = $props();

  let nameInput = $state('');
  let selectedIcon = $state('Sparkles');
  let isSubmitting = $state(false);

  const availableIcons = [
    { id: 'Sparkles', icon: Sparkles, label: 'Sparkles' },
    { id: 'Sun', icon: Sun, label: 'Daylight' },
    { id: 'Moon', icon: Moon, label: 'Night' },
    { id: 'Film', icon: Film, label: 'Cinema' },
    { id: 'Flame', icon: Flame, label: 'Warm Fire' },
    { id: 'Zap', icon: Zap, label: 'Cyberpunk' }
  ];

  async function handleSubmit(e) {
    e.preventDefault();
    if (!nameInput.trim()) return;

    isSubmitting = true;
    try {
      await onCapture(nameInput.trim(), selectedIcon);
      nameInput = '';
      selectedIcon = 'Sparkles';
      onClose();
    } catch (err) {
      console.error('Scene capture error:', err);
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    
    <div class="relative w-full max-w-md bg-slate-900 border border-cyan-500/30 rounded-2xl p-6 shadow-2xl space-y-5">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-800">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl bg-cyan-500/10 border border-cyan-500/30 text-cyan-400">
            <Camera class="w-5 h-5" />
          </div>
          <h3 class="font-semibold text-slate-100 text-base">Capture Scene Snapshot</h3>
        </div>
        <button 
          onclick={onClose}
          class="p-1.5 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Form Body -->
      <form onsubmit={handleSubmit} class="space-y-4">
        <div>
          <label class="block text-xs font-mono text-slate-400 mb-1">Scene Preset Name</label>
          <input 
            type="text" 
            placeholder="e.g. Cyberpunk Neon, Movie Time, Relax" 
            bind:value={nameInput}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-cyan-500/60 placeholder-slate-600"
          />
        </div>

        <!-- Icon Picker Grid -->
        <div>
          <label class="block text-xs font-mono text-slate-400 mb-2">Preset Icon</label>
          <div class="grid grid-cols-3 gap-2">
            {#each availableIcons as item}
              {@const IconComp = item.icon}
              <button
                type="button"
                onclick={() => selectedIcon = item.id}
                class="flex flex-col items-center gap-1.5 p-3 rounded-xl border text-xs transition-all cursor-pointer {selectedIcon === item.id ? 'bg-cyan-500/10 border-cyan-500/50 text-cyan-300 shadow-[0_0_12px_rgba(6,182,212,0.3)]' : 'bg-slate-950/60 border-slate-800 text-slate-400 hover:border-slate-700'}"
              >
                <IconComp class="w-5 h-5" />
                <span class="text-[10px] font-mono">{item.label}</span>
              </button>
            {/each}
          </div>
        </div>

        <!-- Submit Button -->
        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-800">
          <button 
            type="button"
            onclick={onClose}
            class="px-4 py-2 rounded-xl text-xs font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
          >
            Cancel
          </button>
          <button 
            type="submit"
            disabled={isSubmitting || !nameInput.trim()}
            class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-medium text-xs shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all disabled:opacity-50 cursor-pointer"
          >
            <Sparkles class="w-4 h-4" />
            <span>Save Scene</span>
          </button>
        </div>
      </form>

    </div>
  </div>
{/if}
