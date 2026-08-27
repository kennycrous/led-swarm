<script>
  import { X, Camera, Sparkles, Sun, Moon, Film, Flame, Zap } from 'lucide-svelte';

  let { isOpen = false, onClose = () => {}, onCapture = () => {} } = $props();

  let nameInput = $state('');
  let selectedIcon = $state('Sparkles');

  const availableIcons = [
    { id: 'Sparkles', icon: Sparkles, label: 'Sparkles' },
    { id: 'Sun', icon: Sun, label: 'Daytime' },
    { id: 'Moon', icon: Moon, label: 'Night' },
    { id: 'Film', icon: Film, label: 'Movie' },
    { id: 'Flame', icon: Flame, label: 'Fire' },
    { id: 'Zap', icon: Zap, label: 'Cyber' }
  ];

  function handleSubmit(e) {
    e.preventDefault();
    if (nameInput.trim()) {
      onCapture(nameInput.trim(), selectedIcon);
      nameInput = '';
      selectedIcon = 'Sparkles';
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
    <div class="glass-panel w-full max-w-md rounded-2xl p-6 border border-cyan-500/40 shadow-2xl space-y-4">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <h3 class="font-mono text-base font-bold text-slate-100 flex items-center gap-2">
          <Camera class="w-5 h-5 text-cyan-400" />
          Capture Scene Preset
        </h3>
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
          <label for="capture-scene-name" class="block text-xs font-mono text-slate-400 mb-1">Scene Preset Name</label>
          <input
            id="capture-scene-name"
            type="text"
            placeholder="e.g. Cyberpunk Neon, Movie Time, Relax"
            bind:value={nameInput}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-cyan-500/60 placeholder-slate-600"
          />
        </div>

        <!-- Icon Picker Grid -->
        <div>
          <span class="block text-xs font-mono text-slate-400 mb-2">Preset Icon</span>
          <div class="grid grid-cols-3 gap-2">
            {#each availableIcons as item (item.id)}
              {@const IconComp = item.icon}
              <button
                type="button"
                onclick={() => (selectedIcon = item.id)}
                class="flex flex-col items-center gap-1.5 p-3 rounded-xl border text-xs transition-all cursor-pointer {selectedIcon ===
                item.id
                  ? 'bg-cyan-500/10 border-cyan-500/50 text-cyan-300 shadow-[0_0_12px_rgba(6,182,212,0.3)]'
                  : 'bg-slate-950/60 border-slate-800 text-slate-400 hover:border-slate-700'}"
              >
                <IconComp class="w-5 h-5" />
                <span class="font-mono text-[11px]">{item.label}</span>
              </button>
            {/each}
          </div>
        </div>

        <!-- Action Buttons -->
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
            disabled={!nameInput.trim()}
            class="px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all duration-200 disabled:opacity-40"
          >
            Snapshot Scene
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
