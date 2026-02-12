<script lang="ts">
	import type { HTMLFormAttributes } from 'svelte/elements';
	import type { WithElementRef } from '@/utils.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { SearchIcon } from '@lucide/svelte';
	import * as i18n from '@/paraglide/messages.js';

	let { ref = $bindable(null), ...restProps }: WithElementRef<HTMLFormAttributes> = $props();
</script>

<form {...restProps} bind:this={ref} class="w-full md:w-md lg:w-lg xl:w-xl">
	<div class="relative w-full">
		<Label for="search" class="sr-only">{i18n.search()}</Label>
		<Sidebar.Input id="search" placeholder={i18n.search()} class="h-9 w-full ps-7" />
		<SearchIcon
			class="pointer-events-none absolute start-2 top-1/2 size-4 -translate-y-1/2 opacity-50 select-none"
		/>
	</div>
</form>
