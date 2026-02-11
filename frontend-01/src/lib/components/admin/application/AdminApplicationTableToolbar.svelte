<script lang="ts" module>
	type QueryOption = {
		page: number;
		limit: number;
		search: string;
		status: string;
		sort_by: string;
		order_by: 'asc' | 'desc';
		date_from: string;
		date_to: string;
	};
</script>

<script lang="ts" generics="Application">
	import type { Table } from '@tanstack/table-core';
	import { SvelteSet } from 'svelte/reactivity';
	import { DataTableViewOption, CommandSelectInput } from '@/components';
	import { Button } from '@/components/ui/button';
	import { Input } from '@/components/ui/input';
	import { XIcon } from '@lucide/svelte';

	let {
		table,
		statusOptions,
		updateQuery,
		onreset
	}: {
		table: Table<Application>;
		statusOptions: { label: string; value: string }[];
		updateQuery: (updates: Partial<QueryOption>, resetPage: boolean) => Promise<void>;
		onreset: () => Promise<void>;
	} = $props();

	const filtered = $derived(table.getState().columnFilters.length > 0);
	let selectedStatus = $derived(new SvelteSet<string>([]));

	let searchTerms = $state<string | undefined>('');
	let releasedDate = $state<string | undefined>(undefined);

	async function handleReset() {
		table.resetColumnFilters();
		await onreset();
		searchTerms = undefined;
		releasedDate = undefined;
		selectedStatus.clear();
	}

	let searchTimer: ReturnType<typeof setTimeout>;
	async function handleSearch(value: string) {
		clearTimeout(searchTimer);
		searchTimer = setTimeout(async () => {
			await updateQuery?.({ search: value || '' }, true);
		}, 500);
	}
</script>

<div class="flex flex-col items-center gap-4 p-2 lg:flex-row lg:justify-between">
	<div class="flex w-full flex-col items-center gap-2 gap-x-4 lg:flex-row">
		<Input
			bind:value={searchTerms}
			placeholder="Search applications..."
			oninput={(e) => {
				handleSearch(e.currentTarget.value);
			}}
			class="h-8 w-full lg:w-auto"
		/>
	</div>

	{#if filtered || releasedDate || (searchTerms && searchTerms !== '') || selectedStatus.size > 0}
		<Button variant="ghost" class="h-8 px-2 lg:px-3" onclick={handleReset}>
			Reset filters
			<XIcon />
		</Button>
	{/if}

	<CommandSelectInput
		bind:selectedValue={selectedStatus}
		title="Status"
		options={statusOptions}
		onchange={(value) => updateQuery?.({ status: [...value].join(',') }, true)}
	/>

	<DataTableViewOption {table} />
</div>
