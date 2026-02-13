<script lang="ts">
	import { useModalSub } from './modal.svelte.js';
	import * as Dialog from '$lib/components/ui-extras/dialog/index.js';
	import * as Drawer from '$lib/components/ui-extras/drawer/index.js';
	import type { DialogDescriptionProps } from 'bits-ui';

	const modal = useModalSub();

	let { ...rest }: DialogDescriptionProps = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Description {...rest as any} />
{:else}
	<Drawer.Description {...rest as any} />
{/if}
