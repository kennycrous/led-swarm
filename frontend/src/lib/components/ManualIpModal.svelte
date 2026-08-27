<script>
  import { X, Plus, Wifi } from 'lucide-svelte';

  let { isOpen = false, onClose = () => {}, onAdd = (ip) => {} } = $props();

  let ipInput = $state('');
  let isSubmitting = $state(false);
  let errorMsg = $state('');

  async function handleSubmit(e) {
    e.preventDefault();
    if (!ipInput.trim()) return;

    isSubmitting = true;
    errorMsg = '';

    try {
      await onAdd(ipInput.trim());
      ipInput = '';
      onClose();
    } catch (err) {
      errorMsg = err.message || 'Failed to connect to WLED device at that IP address.';
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    <div class="glass-panel w-full max-w-md p-6 rounded-2xl border border-cyan-500/30 shadow-2xl relative">
      <div class="flex items-center justify-between pb-4 border-b border-cyan-500/10">
        <div class="flex items-center gap-2">
          <Wifi class="w-5 h-5 text-cyan-400" />
          <h3 class="font-bold text-slate-100">Add WLED Device by IP</h3>
        </div>
        <button onclick={onClose} class="p-1 rounded-lg text-slate-400 hover:text-slate-200 hover:bg-slate-800/60">
          <X class="w-5 h-5" />
        </button>
      </div>

      <form onsubmit={handleSubmit} class="mt-4 space-y-4">
        <div>
          <label for="wled-ip-input" class="block text-xs font-mono text-slate-400 mb-1.5">WLED IP ADDRESS</label>
          <input
            id="wled-ip-input"
            type="text"
            placeholder="192.168.1.150"
            bind:value={ipInput}
            class="w-full px-4 py-2.5 bg-slate-900/80 border border-slate-800 rounded-xl text-slate-100 font-mono text-sm focus:outline-none focus:border-cyan-500/60"
          />
        </div>

        {#if errorMsg}
          <p class="text-xs text-rose-400 font-mono">{errorMsg}</p>
        {/if}

        <div class="flex justify-end gap-3 pt-2">
          <button
            type="button"
            onclick={onClose}
            class="px-4 py-2 rounded-xl text-xs font-mono border border-slate-800 text-slate-400 hover:bg-slate-800/50"
          >
            CANCEL
          </button>

          <button
            type="submit"
            disabled={isSubmitting || !ipInput.trim()}
            class="flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-mono bg-gradient-to-r from-cyan-500 to-purple-600 text-white font-semibold shadow-lg shadow-cyan-500/20 hover:opacity-90 disabled:opacity-50"
          >
            <Plus class="w-4 h-4" />
            {isSubmitting ? 'CONNECTING...' : 'ADD DEVICE'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
