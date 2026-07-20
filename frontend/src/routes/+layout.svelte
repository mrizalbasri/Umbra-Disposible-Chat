<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/logo.ico';
	import { fade } from 'svelte/transition';
	import { page } from '$app/stores';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<div class="min-h-screen relative" in:fade={{ duration: 250, delay: 50 }} out:fade={{ duration: 150 }}>
	{@render children()}
</div>
