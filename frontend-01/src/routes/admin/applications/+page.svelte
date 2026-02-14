<script lang="ts">
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { handleSubmitLoading } from '@/stores';
	import { AdminSidebarLayout, AdminApplicationTable } from '@/components/admin';
	import { DateRangeInput, AppAlertDialog, CardSpotlight } from '@/components/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import { updateUrlParams, createApplicationManager } from '@/stores/query.js';
	import { toast } from '@/stores';
	import { localizeHref } from '@/paraglide/runtime';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	let openSingleDeleteDialog = $state(false);
	let dataEntityAction = $state<ApplicationResponse | null | undefined>(null);

	const queryManager = createApplicationManager();
	let query = $state(queryManager.parse(page.url));

	$effect(() => {
		query = queryManager.parse(page.url);
	});

	let dateRange = $derived({
		start: query.date_from,
		end: query.date_to
	});

	async function updateQuery(updates: Partial<typeof query>, resetPage = false) {
		handleSubmitLoading(true);
		await updateUrlParams(goto, page.url, updates, {
			resetPage,
			replaceState: true,
			invalidateAll: true
		});
		handleSubmitLoading(false);
	}

	async function handleDateChange(range: { start: string; end: string } | null) {
		if (range) {
			await updateQuery({ date_from: range.start, date_to: range.end }, true);
		}
	}

	async function resetFilters() {
		const url = new URL(page.url);
		url.search = '';
		await goto(url.toString(), { replaceState: true, invalidateAll: true });
	}

	async function refresh() {
		handleSubmitLoading(true);
		await goto(page.url.toString(), { replaceState: true, invalidateAll: true });
		handleSubmitLoading(false);
		await invalidateAll();
	}

	async function handleSingleDelete() {
		if (!dataEntityAction) return;

		try {
			const response = await fetch(`/admin/applications/${dataEntityAction.config.package_name}`, {
				method: 'DELETE'
			});

			const result = await response.json();
			if (!response.ok) throw new Error(result.message || 'Failed to delete application');
			toast.success('Application deleted successfully');
			openSingleDeleteDialog = false;
			dataEntityAction = null;
			await refresh();
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Unknown error');
		}
	}

	async function handleBulkDelete(ids: string[]) {
		if (ids.length === 0) return;

		try {
			const response = await fetch('/admin/applications/bulk', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ ids })
			});

			const result = await response.json();
			if (!response.ok) throw new Error(result.message || 'Failed to delete applications');
			toast.success('Applications deleted successfully');
			await refresh();
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Unknown error');
		}
	}
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Application" user={data.user} setting={data.settings}>
	<div class="@container/main flex flex-col gap-4 md:gap-6">
		<div class="flex-none px-4 py-4 sm:px-6">
			<div class="space-y-1">
				<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Application</h1>
				<p class="text-sm text-muted-foreground">Manage your application data.</p>
			</div>
		</div>
		<div class="flex-none border-b px-4 py-4 sm:px-6">
			<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
				<div class="w-full lg:w-auto">
					<DateRangeInput
						bind:modelValue={dateRange}
						onchange={handleDateChange}
						class="w-full lg:w-auto"
					/>
				</div>
				<div class="flex flex-wrap items-center justify-center gap-2 lg:justify-end">
					<Button href={localizeHref('/admin/applications/create')} variant="default" size="sm">
						<Icon icon="material-symbols:add" />
						<span class="sr-only lg:not-sr-only"> Add Application </span>
					</Button>
					<Button variant="outline" size="sm" onclick={() => refresh()}>
						<Icon icon="material-symbols:refresh" />
						<span class="sr-only lg:not-sr-only"> Refresh Data </span>
					</Button>
					<Button variant="destructive" size="sm" onclick={resetFilters}>
						<Icon icon="ic:sharp-clear-all" />
						<span class="sr-only lg:not-sr-only"> Reset Filters </span>
					</Button>
				</div>
			</div>
		</div>
		<div class="relative flex flex-col gap-4 overflow-auto p-4">
			<CardSpotlight
				variant="primary"
				shadow="large"
				spotlightIntensity="medium"
				spotlight
				class="p-2"
			>
				<div class="overflow-hidden rounded-lg border bg-white/60 px-1 py-4 dark:bg-black/60">
					<AdminApplicationTable
						data={data.applications}
						{updateQuery}
						class="min-w-full"
						onreset={() => resetFilters()}
						onEdit={(val) => {
							if (val) goto(localizeHref(`/admin/applications/${val.config.package_name}`));
						}}
						onDelete={(application) => {
							if (application) {
								dataEntityAction = application;
								openSingleDeleteDialog = true;
							}
						}}
						onbulkdelete={(ids) => handleBulkDelete(ids)}
					/>
				</div>
			</CardSpotlight>
		</div>
	</div>
	<AppAlertDialog
		bind:open={openSingleDeleteDialog}
		type="warning"
		title="Warning"
		message={`Are you sure you want to delete ${dataEntityAction?.config.name || ''}?`}
		labelClose="Cancel"
		labelAction="Confirm"
		onaction={() => {
			if (dataEntityAction) {
				handleSingleDelete();
			}
			dataEntityAction = null;
			openSingleDeleteDialog = false;
		}}
		onclose={() => {
			dataEntityAction = null;
			openSingleDeleteDialog = false;
		}}
	/>
</AdminSidebarLayout>
