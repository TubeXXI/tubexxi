<script lang="ts">
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { handleSubmitLoading } from '@/stores';
	import {
		AdminSidebarLayout,
		AdminDashboardFilterToolbar,
		AdminDateToggle,
		AdminHeading
	} from '@/components/admin';
	import { PlaceholderPattern, CardSpotlight } from '@/components';
	import { updateUrlParams, createDashboardManager } from '@/stores/query.js';
	import { formatToPostgresTimestampV2 } from '@/utils/time.js';
	import Icon from '@iconify/svelte';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	const queryManager = createDashboardManager();
	let query = $state(queryManager.parse(page.url));

	$effect(() => {
		query = queryManager.parse(page.url);
	});

	let dateRange = $derived({
		start: query.date_from,
		end: query.date_to
	});
	let rangeDescription = $derived(
		query.date_from && query.date_to
			? `From ${new Date(query.date_from).toLocaleDateString('en-US', {
					month: '2-digit',
					day: '2-digit',
					year: 'numeric'
				})} to ${new Date(query.date_to).toLocaleDateString('en-US', {
					month: '2-digit',
					day: '2-digit',
					year: 'numeric'
				})}`
			: 'for the selected range'
	);

	async function updateQuery(updates: Partial<typeof query>, resetPage = false) {
		handleSubmitLoading(true);
		await updateUrlParams(goto, page.url, updates, {
			resetPage,
			replaceState: true,
			invalidateAll: true
		});
		handleSubmitLoading(false);
	}

	async function resetDateRange() {
		await updateQuery({
			date_from: query.date_from,
			date_to: query.date_to,
			page: 1
		});
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
	async function setLastWeek() {
		const end = new Date();
		const start = new Date();
		start.setDate(end.getDate() - 7);

		dateRange = {
			start: formatToPostgresTimestampV2(start),
			end: formatToPostgresTimestampV2(end)
		};
		await updateQuery(
			{
				date_from: dateRange.start,
				date_to: dateRange.end,
				page: 1
			},
			true
		);
	}
	async function setLastMonth() {
		const end = new Date();
		const start = new Date();
		start.setMonth(end.getMonth() - 1);

		dateRange = {
			start: formatToPostgresTimestampV2(start),
			end: formatToPostgresTimestampV2(end)
		};
		await updateQuery(
			{
				date_from: dateRange.start,
				date_to: dateRange.end,
				page: 1
			},
			true
		);
	}
	async function setLastYear() {
		const end = new Date();
		const start = new Date();
		start.setFullYear(end.getFullYear() - 1);

		dateRange = {
			start: formatToPostgresTimestampV2(start),
			end: formatToPostgresTimestampV2(end)
		};
		await updateQuery(
			{
				date_from: dateRange.start,
				date_to: dateRange.end,
				page: 1
			},
			true
		);
	}
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Dashboard" user={data.user} setting={data.settings}>
	<div
		class="scrollbar-thin flex min-h-[calc(100vh-160px)] flex-col overflow-hidden overflow-y-auto scroll-smooth scrollbar-thumb-foreground scrollbar-track-accent"
	>
		<div class="flex-none px-4 py-4 sm:px-6">
			<div class="space-y-1">
				<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Admin Dashboard</h1>
				<p class="text-sm text-muted-foreground">
					Welcome back,
					<span class="font-semibold text-indigo-600 uppercase dark:text-indigo-400">
						{data.user?.full_name || ''!}
					</span>
				</p>
			</div>
		</div>
		<div class="flex-none border-b px-4 py-4 lg:px-6">
			<div class="space-y-4">
				<div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
					<AdminDashboardFilterToolbar
						bind:dateRange
						onReset={resetDateRange}
						{updateQuery}
						{refresh}
					/>
					<AdminDateToggle
						onReset={resetFilters}
						{setLastWeek}
						{setLastMonth}
						{setLastYear}
						{updateQuery}
					/>
				</div>
			</div>
		</div>
		<div class="flex h-full flex-1 flex-col gap-4 overflow-x-auto rounded-xl p-4">
			<div class="grid auto-rows-min gap-4 md:grid-cols-3">
				<div
					class="relative aspect-video overflow-hidden rounded-xl border border-sidebar-border/70 dark:border-sidebar-border"
				>
					<PlaceholderPattern />
					<div class="flex h-full w-full flex-col items-center justify-center gap-2 p-2">
						<div
							class="flex w-full items-center justify-between rounded-xl bg-blue-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_apps || 0} -->
								</p>
								<p class="text-sm">Apps</p>
							</div>
							<Icon icon="tdesign:app-filled" class="h-6 w-6" />
						</div>
						<div
							class="flex w-full items-center justify-between rounded-xl bg-cyan-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_platforms || 0} -->
								</p>
								<p class="text-sm">Platforms</p>
							</div>
							<Icon icon="tdesign:control-platform-filled" class="h-6 w-6" />
						</div>
					</div>
				</div>
				<div
					class="relative aspect-video overflow-hidden rounded-xl border border-sidebar-border/70 dark:border-sidebar-border"
				>
					<PlaceholderPattern />
					<div class="flex h-full w-full flex-col items-center justify-center gap-2 p-2">
						<div
							class="flex w-full items-center justify-between rounded-xl bg-green-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_downloads || 0} -->
								</p>
								<p class="text-sm">Downloads</p>
							</div>
							<Icon icon="ic:baseline-download" class="h-6 w-6" />
						</div>
						<div
							class="flex w-full items-center justify-between rounded-xl bg-red-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_subscriptions || 0} -->
								</p>
								<p class="text-sm">Subscriptions</p>
							</div>
							<Icon icon="mdi:payment-clock" class="h-6 w-6" />
						</div>
					</div>
				</div>
				<div
					class="relative aspect-video overflow-hidden rounded-xl border border-sidebar-border/70 dark:border-sidebar-border"
				>
					<PlaceholderPattern />
					<div class="flex h-full w-full flex-col items-center justify-center gap-2 p-2">
						<div
							class="flex w-full items-center justify-between rounded-xl bg-yellow-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_transactions || 0} -->
								</p>
								<p class="text-sm">Transactions</p>
							</div>
							<Icon icon="hugeicons:transaction-history" class="h-6 w-6" />
						</div>
						<div
							class="flex w-full items-center justify-between rounded-xl bg-purple-600 p-4 text-white shadow-md"
						>
							<div>
								<p class="text-xl font-bold">
									<!-- {dashboardData.stats?.total_users || 0} -->
								</p>
								<p class="text-sm">Users</p>
							</div>
							<Icon icon="mdi:account" class="h-6 w-6" />
						</div>
					</div>
				</div>
			</div>
			<div
				class="relative min-h-screen flex-1 space-y-4 rounded-xl border border-sidebar-border/70 p-4 md:min-h-min dark:border-sidebar-border"
			>
				<!-- <AdminDashboardChartBarInteractive analytics={dashboardData.analytics} {rangeDescription} /> -->
				<CardSpotlight variant="info" shadow="large" spotlightIntensity="medium" spotlight>
					<AdminHeading title="Recent Downloads" description={rangeDescription} />
					<!-- <AdminDashboardRecentDownload downloads={dashboardData.recent_downloads} /> -->
				</CardSpotlight>
			</div>
		</div>
	</div>
</AdminSidebarLayout>
