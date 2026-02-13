<script lang="ts">
	import { useModalSub } from './modal.svelte.js';
	import * as Dialog from '$lib/components/ui-extras/dialog/index.js';
	import * as Drawer from '$lib/components/ui-extras/drawer/index.js';
	import type { DialogContentProps } from 'bits-ui';

	const modal = useModalSub();

	let {
		hideClose = false,
		children,
		...rest
	}: DialogContentProps & { hideClose?: boolean } = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Content {...(rest as any)} {hideClose}>
		{@render children?.()}
	</Dialog.Content>
{:else}
	<Drawer.Content {...(rest as any)}>
		{@render children?.()}
	</Drawer.Content>
{/if}
