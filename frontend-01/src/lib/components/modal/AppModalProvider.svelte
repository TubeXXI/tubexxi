<script lang="ts">
	import type { Snippet } from 'svelte';
	import AppModal from './AppModal.svelte';

	let {
		children
	}: {
		children: Snippet;
	} = $props();
</script>

{@render children()}
<AppModal />
