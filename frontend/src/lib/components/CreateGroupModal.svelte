<script>
  import { X, Plus, Layers, Cpu } from 'lucide-svelte';

  let { 
    isOpen = false, 
    allDevices = [], 
    onClose = () => {}, 
    onCreate = (name, description, deviceIds) => {} 
  } = $props();

  let nameInput = $state('');
  let descInput = $state('');
  let selectedDeviceIds = $state([]);
  let isSubmitting = $state(false);

  function toggleDevice(id) {
    if (selectedDeviceIds.includes(id)) {
      selectedDeviceIds = selectedDeviceIds.filter(dId => dId !== id);
    } else {
      selectedDeviceIds = [...selectedDeviceIds, id];
    }
  }

  async function handleSubmit(e) {
    e.preventDefault();
    if (!nameInput.trim()) return;

    isSubmitting = true;
    try {
      await onCreate(nameInput.trim(), descInput.trim(), selectedDeviceIds);
      nameInput = '';
      descInput = '';
      selectedDeviceIds = [];
      onClose();
    } catch (err) {
      console.error('Group creation error:', err);
    } finally {
      isSubmitting = false;
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-md">
    
    <div class="relative w-full max-w-md bg-slate-900 border border-purple-500/30 rounded-2xl p-6 shadow-2xl space-y-5">
      
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 border-b border-slate-800">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl bg-purple-500/10 border border-purple-500/30 text-purple-400">
            <Layers class="w-5 h-5" />
          </div>
          <h3 class="font-semibold text-slate-100 text-base">Create Strip Group</h3>
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
          <label class="block text-xs font-mono text-slate-400 mb-1">Group Name</label>
          <input 
            type="text" 
            placeholder="e.g. Desk Lights, TV Setup" 
            bind:value={nameInput}
            required
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-purple-500/60 placeholder-slate-600"
          />
        </div>

        <div>
          <label class="block text-xs font-mono text-slate-400 mb-1">Description (Optional)</label>
          <input 
            type="text" 
            placeholder="e.g. All monitors and ambient backlights" 
            bind:value={descInput}
            class="w-full bg-slate-950 border border-slate-800 rounded-xl px-3.5 py-2 text-sm text-slate-200 focus:outline-none focus:border-purple-500/60 placeholder-slate-600"
          />
        </div>

        <!-- Devices Selection Checkboxes -->
        <div>
          <label class="block text-xs font-mono text-slate-400 mb-2">Assign WLED Light Strips</label>
          <div class="max-h-40 overflow-y-auto space-y-1.5 pr-1">
            {#each allDevices as dev}
              <button
                type="button"
                onclick={() => toggleDevice(dev.id)}
                class="w-full flex items-center justify-between px-3 py-2 rounded-xl border text-xs transition-all text-left cursor-pointer {selectedDeviceIds.includes(dev.id) ? 'bg-purple-500/10 border-purple-500/50 text-purple-200' : 'bg-slate-950/60 border-slate-800 text-slate-400'}"
              >
                <span class="flex items-center gap-2 font-mono">
                  <Cpu class="w-3.5 h-3.5 text-cyan-400" />
                  {dev.name} ({dev.ipAddress})
                </span>
                <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-slate-900 border border-slate-800">
                  {selectedDeviceIds.includes(dev.id) ? '✓ Included' : '+ Select'}
                </span>
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
            class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-medium text-xs shadow-[0_0_15px_rgba(168,85,247,0.4)] transition-all disabled:opacity-50 cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            <span>Create Group</span>
          </button>
        </div>
      </form>

    </div>
  </div>
{/if}
