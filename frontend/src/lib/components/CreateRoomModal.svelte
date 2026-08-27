<script>
  import { X, LayoutGrid, Check } from 'lucide-svelte';

  let { isOpen = false, availableDevices = [], onClose = () => {}, onCreate = () => {} } = $props();

  let nameInput = $state('');
  let descInput = $state('');
  let selectedDeviceIds = $state([]);

  function toggleDevice(id) {
    if (selectedDeviceIds.includes(id)) {
      selectedDeviceIds = selectedDeviceIds.filter((devId) => devId !== id);
    } else {
      selectedDeviceIds = [...selectedDeviceIds, id];
    }
  }

  function handleSubmit(e) {
    e.preventDefault();
    if (nameInput.trim()) {
      onCreate(nameInput.trim(), descInput.trim(), selectedDeviceIds);
      nameInput = '';
      descInput = '';
      selectedDeviceIds = [];
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
    <div class="glass-panel w-full max-w-lg rounded-2xl p-6 border border-cyan-500/40 shadow-2xl space-y-4">
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <h3 class="font-mono text-base font-bold text-slate-100 flex items-center gap-2">
          <LayoutGrid class="w-5 h-5 text-cyan-400" />
          Create 2D Room Canvas
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
          <label for="room-name-input" class="block text-xs font-mono text-slate-400 mb-1">Room Canvas Name</label>
          <input
            id="room-name-input"
            type="text"
            placeholder="e.g. Living Room, Office, Gaming Den"
            bind:value={nameInput}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-cyan-500/60 placeholder-slate-600"
          />
        </div>

        <div>
          <label for="room-desc-input" class="block text-xs font-mono text-slate-400 mb-1">Description (Optional)</label
          >
          <input
            id="room-desc-input"
            type="text"
            placeholder="e.g. Main TV and ceiling strip 2D spatial layout"
            bind:value={descInput}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-cyan-500/60 placeholder-slate-600"
          />
        </div>

        <!-- Devices Selection Checklist -->
        <div>
          <span class="block text-xs font-mono text-slate-400 mb-2">Assign WLED Light Strips to Room</span>
          {#if availableDevices.length === 0}
            <p class="text-xs text-slate-500 italic p-3 bg-slate-950/50 rounded-xl border border-slate-800">
              No discovered WLED devices available. You can add strips later.
            </p>
          {:else}
            <div class="max-h-48 overflow-y-auto space-y-2 pr-1">
              {#each availableDevices as dev (dev.id)}
                {@const isChecked = selectedDeviceIds.includes(dev.id)}
                <button
                  type="button"
                  onclick={() => toggleDevice(dev.id)}
                  class="w-full flex items-center justify-between p-3 rounded-xl border text-left text-xs transition-all cursor-pointer {isChecked
                    ? 'bg-cyan-500/10 border-cyan-500/50 text-cyan-200 shadow-[0_0_10px_rgba(6,182,212,0.15)]'
                    : 'bg-slate-950/60 border-slate-800 text-slate-400 hover:border-slate-700'}"
                >
                  <div class="flex items-center gap-2.5">
                    <div
                      class="w-4 h-4 rounded border flex items-center justify-center transition-colors {isChecked
                        ? 'border-cyan-400 bg-cyan-500 text-slate-950'
                        : 'border-slate-700 bg-slate-900'}"
                    >
                      {#if isChecked}
                        <Check class="w-3 h-3 stroke-[3]" />
                      {/if}
                    </div>
                    <span class="font-medium">{dev.name}</span>
                  </div>
                  <span class="font-mono text-[11px] text-slate-500">{dev.ipAddress} ({dev.ledCount || 0} LEDs)</span>
                </button>
              {/each}
            </div>
          {/if}
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
            class="px-4 py-2 rounded-xl bg-gradient-to-r from-cyan-500 to-purple-600 hover:from-cyan-400 hover:to-purple-500 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(6,182,212,0.4)] transition-all duration-200"
          >
            Create 2D Room
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
