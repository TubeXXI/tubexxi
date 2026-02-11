<script lang="ts">
	import { Button, type ButtonProps } from '$lib/components/ui/button/index.js';
	import type { Snippet } from 'svelte';

	let { ref = $bindable(null), children, ...restProps }: ButtonProps & { children?: Snippet } = $props();
</script>

<Button bind:ref type="submit" {...restProps}>
	{#if children}
		{@render children()}
	{/if}
</Button>
