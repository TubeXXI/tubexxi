<script lang="ts">
	import { useModalSub } from './modal.svelte.js';
	import * as Dialog from '$lib/components/ui-extras/dialog/index.js';
	import * as Drawer from '$lib/components/ui-extras/drawer/index.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	const modal = useModalSub();

	let {
		ref = $bindable(null),
		children,
		...rest
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Footer {...rest as any}>
		{@render children?.()}
	</Dialog.Footer>
{:else}
	<Drawer.Footer {...rest as any}>
		{@render children?.()}
	</Drawer.Footer>
{/if}
