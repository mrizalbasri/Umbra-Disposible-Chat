<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';

	onMount(() => {
		goto(resolve('/chatroom?tab=match'));
	});
</script>

<div class="flex items-center justify-center h-screen bg-background">
	<div class="flex flex-col items-center space-y-4">
		<div class="w-10 h-10 border-4 border-[#00aeef] border-t-transparent rounded-full animate-spin"></div>
		<p class="text-sm font-mono text-gray-500 uppercase tracking-wider">Menghubungkan ke Random Match...</p>
	</div>
</div>
