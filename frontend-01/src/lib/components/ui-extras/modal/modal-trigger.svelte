<script lang="ts">
	import { useModalSub } from './modal.svelte.js';
	import * as Dialog from '$lib/components/ui-extras/dialog/index.js';
	import * as Drawer from '$lib/components/ui-extras/drawer/index.js';
	import type { DialogTriggerProps } from 'bits-ui';

	const modal = useModalSub();

	let { children, ...rest }: DialogTriggerProps = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Trigger {...(rest as any)}>
		{@render children?.()}
	</Dialog.Trigger>
{:else}
	<Drawer.Trigger {...(rest as any)}>
		{@render children?.()}
	</Drawer.Trigger>
{/if}
