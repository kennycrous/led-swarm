<script>
  import { X, Layers, Cpu, Check } from 'lucide-svelte';

  let { 
    isOpen = false, 
    allDevices = [], 
    onClose = () => {}, 
    onCreate = () => {} 
  } = $props();

  let nameInput = $state('');
  let descInput = $state('');
  let selectedDeviceIds = $state([]);

  function toggleDevice(id) {
    if (selectedDeviceIds.includes(id)) {
      selectedDeviceIds = selectedDeviceIds.filter(dId => dId !== id);
    } else {
      selectedDeviceIds = [...selectedDeviceIds, id];
    }
  }

  function handleSubmit(e) {
    e.preventDefault();
    if (nameInput.trim() && selectedDeviceIds.length > 0) {
      onCreate(nameInput.trim(), descInput.trim(), selectedDeviceIds);
      // Reset form
      nameInput = '';
      descInput = '';
      selectedDeviceIds = [];
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
    <div class="glass-panel w-full max-w-md rounded-2xl p-6 border border-purple-500/40 shadow-2xl space-y-4">
      
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-slate-800 pb-3">
        <h3 class="font-mono text-base font-bold text-slate-100 flex items-center gap-2">
          <Layers class="w-5 h-5 text-purple-400" />
          Create Virtual Strip Group
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
          <label for="create-group-name" class="block text-xs font-mono text-slate-400 mb-1">Group Name</label>
          <input 
            id="create-group-name"
            type="text" 
            placeholder="e.g. Desk Lights, TV Setup" 
            bind:value={nameInput}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-purple-500/60 placeholder-slate-600"
          />
        </div>

        <div>
          <label for="create-group-desc" class="block text-xs font-mono text-slate-400 mb-1">Description (Optional)</label>
          <input 
            id="create-group-desc"
            type="text" 
            placeholder="e.g. All monitors and ambient backlights" 
            bind:value={descInput}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-purple-500/60 placeholder-slate-600"
          />
        </div>

        <!-- Devices Selection Checkboxes -->
        <div>
          <span class="block text-xs font-mono text-slate-400 mb-2">Assign WLED Light Strips</span>
          <div class="max-h-40 overflow-y-auto space-y-1.5 pr-1">
            {#each allDevices as dev}
              <button
                type="button"
                onclick={() => toggleDevice(dev.id)}
                class="w-full flex items-center justify-between px-3 py-2 rounded-xl border text-xs transition-all text-left cursor-pointer {selectedDeviceIds.includes(dev.id) ? 'bg-purple-500/10 border-purple-500/50 text-purple-200' : 'bg-slate-950/60 border-slate-800 text-slate-400'}"
              >
                <div class="flex items-center gap-2">
                  <Cpu class="w-3.5 h-3.5 text-purple-400" />
                  <span class="font-medium">{dev.name}</span>
                  <span class="text-[10px] font-mono text-slate-500">({dev.ledCount} LEDs)</span>
                </div>
                {#if selectedDeviceIds.includes(dev.id)}
                  <Check class="w-4 h-4 text-purple-400" />
                {/if}
              </button>
            {/each}
          </div>
          {#if selectedDeviceIds.length === 0}
            <p class="text-[10px] text-rose-400 font-mono mt-1">* Select at least 1 strip to create a group</p>
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
            disabled={!nameInput.trim() || selectedDeviceIds.length === 0}
            class="px-4 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-cyan-500 hover:from-purple-500 hover:to-cyan-400 text-white font-mono text-xs font-semibold shadow-[0_0_15px_rgba(168,85,247,0.4)] transition-all duration-200 disabled:opacity-40"
          >
            Create Group
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
