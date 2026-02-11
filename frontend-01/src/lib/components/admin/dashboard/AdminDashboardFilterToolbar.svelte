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
import { Button } from '$lib/components/ui/button/index.js';
import { DateRangeInput } from '@/components';
import Icon from '@iconify/svelte';
import { IsMobile } from '@/hooks/is-mobile.svelte';

let {
	dateRange = $bindable(),
	updateQuery,
	onReset,
	refresh
}: {
	dateRange: {
		start: string;
		end: string;
	};
	updateQuery(updates: Partial<QueryOption>, resetPage?: boolean): Promise<void>;
	onReset(): Promise<void>;
	refresh(): Promise<void>;
} = $props();

const mobileHook = new IsMobile();
const isMobile = $derived(mobileHook.current);

async function handleDateChange(range: { start: string; end: string } | null) {
	if (range) {
		await updateQuery({ date_from: range.start, date_to: range.end }, true);
	}
}
</script>

<div class="flex w-auto flex-col items-center gap-2 lg:flex-row">
	<DateRangeInput
		bind:modelValue={dateRange}
		onchange={handleDateChange}
		class="w-full lg:w-auto"
	/>
	<div class="flex items-center gap-2">
		<Button variant="outline" size={isMobile ? 'sm' : 'icon'} onclick={refresh}>
			<Icon icon="material-symbols:refresh" />
			<span class="not-sr-only text-xs lg:sr-only">Refresh data</span>
		</Button>
		<Button variant="destructive" size={isMobile ? 'sm' : 'icon'} onclick={onReset}>
			<Icon icon="material-symbols:clear-all" />
			<span class="not-sr-only text-xs lg:sr-only">Reset date</span>
		</Button>
	</div>
</div>
