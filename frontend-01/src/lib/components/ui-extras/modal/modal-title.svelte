<script lang="ts">
	import { useModalSub } from './modal.svelte.js';
	import * as Dialog from '$lib/components/ui-extras/dialog/index.js';
	import * as Drawer from '$lib/components/ui-extras/drawer/index.js';
	import type { DialogTitleProps } from 'bits-ui';

	const modal = useModalSub();

	let { ...rest }: DialogTitleProps = $props();
</script>

{#if modal.view === 'desktop'}
	<Dialog.Title {...rest as any} />
{:else}
	<Drawer.Title {...rest as any} />
{/if}
