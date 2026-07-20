<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/logo.ico';
	import { fade } from 'svelte/transition';
	import { page } from '$app/stores';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

{#key $page.url.pathname}
	<div in:fade={{ duration: 220, delay: 60 }} out:fade={{ duration: 160 }}>
		{@render children()}
	</div>
{/key}
