<script>
  import { ChevronDown, Check } from 'lucide-svelte';

  let { 
    value = 0, 
    options = [], 
    icon: IconComp = null, 
    iconColor = 'text-cyan-400',
    hoverBorder = 'hover:border-cyan-500/40',
    onChange = (val) => {} 
  } = $props();

  let isOpen = $state(false);

  const selectedLabel = $derived(
    options[value] || (options.length === 0 ? 'Default' : options[0] || 'Select...')
  );

  function handleSelect(index) {
    onChange(index);
    isOpen = false;
  }
</script>

<div class="relative w-full {isOpen ? 'z-50' : 'z-10'}">
  <!-- Trigger Button -->
  <button
    type="button"
    onclick={() => isOpen = !isOpen}
    class="w-full flex items-center justify-between gap-1.5 bg-[#090e17]/80 hover:bg-slate-900/90 border border-slate-800 {hoverBorder} rounded-xl px-2.5 py-1 h-9 text-[11px] font-mono text-slate-200 transition-all cursor-pointer shadow-sm"
  >
    <div class="flex items-center gap-1.5 min-w-0 truncate">
      {#if IconComp}
        <IconComp class="w-3.5 h-3.5 {iconColor} shrink-0" />
      {/if}
      <span class="truncate">{selectedLabel}</span>
    </div>
    <ChevronDown class="w-3.5 h-3.5 text-slate-500 shrink-0 transition-transform {isOpen ? 'rotate-180 text-cyan-400' : ''}" />
  </button>

  <!-- Cyberpop Glassmorphic Dropdown List -->
  {#if isOpen}
    <!-- Backdrop dismiss -->
    <button 
      type="button" 
      aria-label="Close dropdown"
      class="fixed inset-0 z-40 bg-transparent border-0 cursor-default p-0" 
      onclick={() => isOpen = false}
    ></button>

    <!-- Outer Rounded Container (Enforces overflow-hidden on corners to clip scrollbar) -->
    <div class="absolute top-10 left-0 right-0 z-50 rounded-2xl border border-cyan-500/30 bg-[#090e17]/95 shadow-[0_0_25px_rgba(6,182,212,0.25)] backdrop-blur-xl overflow-hidden">
      <!-- Inner Scrollable Area -->
      <div class="max-h-48 overflow-y-auto p-1.5 pr-1 space-y-0.5 text-xs font-mono custom-scrollbar">
        {#if options.length === 0}
          <div class="px-3 py-1.5 text-slate-500 text-[11px]">Solid</div>
        {:else}
          {#each options as opt, idx}
            <button
              type="button"
              onclick={() => handleSelect(idx)}
              class="w-full flex items-center justify-between px-2.5 py-1.5 rounded-xl text-[11px] text-left transition-all cursor-pointer {value === idx ? 'bg-cyan-500/20 text-cyan-300 font-semibold' : 'text-slate-300 hover:text-cyan-200 hover:bg-cyan-500/10'}"
            >
              <span class="truncate pr-1">{opt}</span>
              {#if value === idx}
                <Check class="w-3 h-3 text-cyan-400 shrink-0" />
              {/if}
            </button>
          {/each}
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(6, 182, 212, 0.35);
    border-radius: 9999px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(6, 182, 212, 0.7);
  }
</style>
