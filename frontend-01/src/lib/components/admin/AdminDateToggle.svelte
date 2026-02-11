<script lang="ts" module>
	type QueryOption = {
		page: number;
		limit: number;
		sort_by: string;
		order_by: 'asc' | 'desc';
		date_from: string;
		date_to: string;
	};
</script>

<script lang="ts">
	import * as ToggleGroup from '$lib/components/ui/toggle-group/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import * as Select from '@/components/ui/select';

	let {
		onReset,
		setLastWeek,
		setLastMonth,
		setLastYear,
		updateQuery
	}: {
		onReset: () => void;
		setLastWeek(): Promise<void>;
		setLastMonth(): Promise<void>;
		setLastYear(): Promise<void>;
		updateQuery(updates: Partial<QueryOption>, resetPage?: boolean): Promise<void>;
	} = $props();

	let selectedValue = $state('');
	let selectedLimit = $state<number>(10);

	async function handleLimitChange() {
		await updateQuery({ limit: selectedLimit }, true);
	}
</script>

<div class="flex flex-col lg:flex-row items-center gap-4 lg:gap-2">
<Select.Root
	type="single"
	value={String(selectedLimit)}
	onValueChange={(value) => {
		if (value) {
			selectedLimit = Number(value);
			handleLimitChange();
		}
	}}
>
	<Select.Trigger class="h-8 w-full lg:w-auto">
		{String(selectedLimit)}
	</Select.Trigger>
	<Select.Content side="top">
		{#each [10, 20, 30, 40, 50, 100] as pageSize (pageSize)}
			<Select.Item value={String(pageSize)}>
				{pageSize}
			</Select.Item>
		{/each}
	</Select.Content>
</Select.Root>

<ToggleGroup.Root
	type="single"
	variant="outline"
	onValueChange={(value) => (selectedValue = value)}
	class="gap-1 *:data-[slot=toggle-group-item]:px-4! @[767px]/card:flex"
>
	<ToggleGroup.Item
		value="week"
		aria-label="Toggle week"
		class={buttonVariants({ variant: 'ghost', class: 'cursor-pointer' })}
		onclick={setLastWeek}
	>
		Week
	</ToggleGroup.Item>
	<ToggleGroup.Item
		value="month"
		aria-label="Toggle month"
		class={buttonVariants({ variant: 'ghost', class: 'cursor-pointer' })}
		onclick={setLastMonth}
	>
		Month
	</ToggleGroup.Item>
	<ToggleGroup.Item
		value="year"
		aria-label="Toggle year"
		class={buttonVariants({ variant: 'ghost', class: 'cursor-pointer' })}
		onclick={setLastYear}
	>
		Year
	</ToggleGroup.Item>
	<ToggleGroup.Item
		value="reset"
		aria-label="Toggle reset"
		class={buttonVariants({ variant: 'destructive', class: 'cursor-pointer' })}
		onclick={() => {
			onReset();
			selectedValue = '';
		}}
	>
		Reset
	</ToggleGroup.Item>
</ToggleGroup.Root>
</div>
