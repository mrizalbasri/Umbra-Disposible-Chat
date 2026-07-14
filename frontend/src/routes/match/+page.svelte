<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  onMount(() => {
    goto('/chatroom?tab=match');
  });
</script>

<div class="h-screen flex flex-col items-center justify-center bg-[#F8FAFC] space-y-4">
  <div class="w-10 h-10 border-4 border-[#00AEEF] border-t-transparent rounded-full animate-spin"></div>
  <p class="text-sm font-sans text-slate-500 font-medium">Menghubungkan ke Matchmaking...</p>
</div>
