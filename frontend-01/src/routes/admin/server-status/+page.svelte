<script lang="ts">
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { handleSubmitLoading } from '@/stores';
	import { AdminSidebarLayout } from '@/components/admin';
	import { AppAlertDialog, CardSpotlight } from '@/components/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';
	import Icon from '@iconify/svelte';
	import { toast } from '@/stores';
	import { formatTimeAgo } from '@/utils/time.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	let serverHealth = $derived(data.serverHealth);
	let serverLogs = $derived(data.serverLogs);

	let selectedLevel = $state('all');
	let searchQuery = $state('');
	// svelte-ignore state_referenced_locally
	let currentPage = $state(serverLogs?.pagination?.current_page || 1);
	// svelte-ignore state_referenced_locally
	let limit = $state(serverLogs?.pagination?.limit || 50);

	let logsData = $derived(serverLogs?.data || []);
	let paginationMeta = $derived(serverLogs?.pagination);
	let openClearLogDialog = $state(false);

	const filteredLogs = $derived(() => {
		if (!logsData) return [];

		let logs = logsData;

		if (selectedLevel !== 'all') {
			logs = logs.filter((log) => log.level === selectedLevel);
		}

		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			logs = logs.filter(
				(log) =>
					log.message?.toLowerCase().includes(query) ||
					log.level?.toLowerCase().includes(query) ||
					log.command?.toLowerCase().includes(query) ||
					log.sql?.toLowerCase().includes(query)
			);
		}

		return logs;
	});

	const logLevelCounts = $derived(() => {
		if (!logsData) return { all: 0, info: 0, debug: 0, warn: 0, error: 0 };

		const counts = { all: logsData.length, info: 0, debug: 0, warn: 0, error: 0 };
		logsData.forEach((log) => {
			if (log.level in counts) {
				if (
					log.level === 'info' ||
					log.level === 'debug' ||
					log.level === 'warn' ||
					log.level === 'error'
				) {
					counts[log.level]++;
				}
			}
		});
		return counts;
	});

	function getStatusColor(status: string) {
		return status === 'up'
			? 'text-green-600 dark:text-green-400'
			: 'text-red-600 dark:text-red-400';
	}

	function getStatusBgColor(status: string) {
		return status === 'up' ? 'bg-green-50 dark:bg-green-900' : 'bg-red-50 dark:bg-red-900';
	}

	function getLogLevelVariant(level: string): any {
		const variants: Record<string, any> = {
			info: 'default',
			debug: 'secondary',
			warn: 'outline',
			error: 'destructive'
		};
		return variants[level] || 'default';
	}

	function getLogLevelIcon(level: string) {
		const icons: Record<string, string> = {
			info: 'mdi:information-outline',
			debug: 'mdi:bug-outline',
			warn: 'mdi:alert-outline',
			error: 'mdi:alert-circle-outline'
		};
		return icons[level] || 'mdi:circle-outline';
	}

	function formatFullTime(time: string) {
		try {
			const date = new Date(time);
			return new Intl.DateTimeFormat('id-ID', {
				day: '2-digit',
				month: 'short',
				year: 'numeric',
				hour: '2-digit',
				minute: '2-digit',
				second: '2-digit',
				hour12: false
			}).format(date);
		} catch {
			return time;
		}
	}

	async function refreshData() {
		await goto(`?page=${currentPage}&limit=${limit}`, { keepFocus: true, replaceState: true });
		await invalidateAll();
	}

	let autoRefresh = $state(true);
	let refreshInterval: number;

	$effect(() => {
		if (autoRefresh) {
			refreshInterval = window.setInterval(refreshData, 30000);
		} else {
			clearInterval(refreshInterval);
		}

		return () => clearInterval(refreshInterval);
	});

	$effect(() => {
		const url = new URL(window.location.href);
		const pageParam = Number(url.searchParams.get('page')) || 1;
		const limitParam = Number(url.searchParams.get('limit')) || limit;
		if (pageParam !== currentPage || limitParam !== limit) {
			goto(`?page=${currentPage}&limit=${limit}`, { keepFocus: true, replaceState: true });
		}
	});

	async function handleClearLog() {
		try {
			handleSubmitLoading(true);
			const response = await fetch('/admin/server-status/logs', {
				method: 'DELETE'
			});
			const data = await response.json();
			if (!response.ok || !data.success) {
				throw new Error(data.message || 'Failed to clear server logs');
			}
			toast.success('Logs cleared successfully');
		} catch (error) {
			toast.error(error instanceof Error ? error.message : 'Unknown error');
		} finally {
			await refreshData();
			openClearLogDialog = false;
			handleSubmitLoading(false);
		}
	}
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Server Status" user={data.user} setting={data.settings}>
	<div class="@container/main flex flex-col gap-4 md:gap-6">
		<div class="flex-none px-4 py-4 sm:px-6">
			<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div class="space-y-1">
					<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Server Status</h1>
					<p class="text-sm text-muted-foreground">Monitor server status and logs in real-time.</p>
				</div>
				<div class="flex items-center justify-center gap-2">
					<Button variant="outline" size="sm" onclick={refreshData}>
						<Icon icon="mdi:refresh" />
						<span class="sr-only lg:not-sr-only"> Refresh </span>
					</Button>
					<Button
						variant={autoRefresh ? 'default' : 'outline'}
						size="sm"
						onclick={() => (autoRefresh = !autoRefresh)}
					>
						<Icon icon={autoRefresh ? 'mdi:autorenew' : 'mdi:refresh-off'} />
						<span class="sr-only lg:not-sr-only"> Auto Refresh </span>
					</Button>
					<Button variant="destructive" size="sm" onclick={() => (openClearLogDialog = true)}>
						<Icon icon="mdi:trash-can-outline" />
						<span class="sr-only lg:not-sr-only"> Clear Logs </span>
					</Button>
				</div>
			</div>
		</div>
		<!-- Server Health Status -->
		<div class="px-4 sm:px-6">
			<CardSpotlight variant="success" shadow="large" spotlightIntensity="medium" spotlight>
				<Card.Root class="bg-white/40 dark:bg-black/40">
					<Card.Header>
						<div class="flex items-center justify-between">
							<Card.Title>Server Health</Card.Title>
							{#if serverHealth?.time}
								<span class="text-xs text-muted-foreground">
									Last updated: {formatFullTime(serverHealth.time)}
								</span>
							{/if}
						</div>
					</Card.Header>
					<Card.Content>
						<div class="grid gap-4 md:grid-cols-3">
							<div
								class="rounded-lg border {getStatusBgColor(serverHealth?.database || 'down')} p-4"
							>
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-3">
										<div class="rounded-full bg-white p-2.5 shadow-sm dark:bg-neutral-900">
											<Icon
												icon="mdi:database"
												class="h-5 w-5 {getStatusColor(serverHealth?.database || 'down')}"
											/>
										</div>
										<div>
											<p class="text-sm font-medium text-muted-foreground">Database</p>
											<p
												class="text-2xl font-bold capitalize {getStatusColor(
													serverHealth?.database || 'down'
												)}"
											>
												{serverHealth?.database || 'Unknown'}
											</p>
										</div>
									</div>
								</div>
							</div>
							<div class="rounded-lg border {getStatusBgColor(serverHealth?.redis || 'down')} p-4">
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-3">
										<div class="rounded-full bg-white p-2.5 shadow-sm dark:bg-neutral-900">
											<Icon
												icon="mdi:memory"
												class="h-5 w-5 {getStatusColor(serverHealth?.redis || 'down')}"
											/>
										</div>
										<div>
											<p class="text-sm font-medium text-muted-foreground">Redis</p>
											<p
												class="text-2xl font-bold capitalize {getStatusColor(
													serverHealth?.redis || 'down'
												)}"
											>
												{serverHealth?.redis || 'Unknown'}
											</p>
										</div>
									</div>
								</div>
							</div>
							<div
								class="rounded-lg border {getStatusBgColor(
									serverHealth?.database === 'up' && serverHealth?.redis === 'up' ? 'up' : 'down'
								)} p-4"
							>
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-3">
										<div class="rounded-full bg-white p-2.5 shadow-sm dark:bg-neutral-900">
											<Icon
												icon="mdi:server"
												class="h-5 w-5 {getStatusColor(
													serverHealth?.database === 'up' && serverHealth?.redis === 'up'
														? 'up'
														: 'down'
												)}"
											/>
										</div>
										<div>
											<p class="text-sm font-medium text-muted-foreground">Overall</p>
											<p
												class="text-2xl font-bold capitalize {getStatusColor(
													serverHealth?.database === 'up' && serverHealth?.redis === 'up'
														? 'up'
														: 'down'
												)}"
											>
												{serverHealth?.database === 'up' && serverHealth?.redis === 'up'
													? 'Healthy'
													: 'Issues Detected'}
											</p>
										</div>
									</div>
								</div>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</CardSpotlight>
		</div>

		<!-- Server Logs -->
		<div class="px-4 sm:px-6">
			<CardSpotlight variant="info" shadow="large" spotlightIntensity="medium" spotlight>
				<Card.Root class="bg-white/40 dark:bg-black/40">
					<Card.Header>
						<Card.Title>Server Logs</Card.Title>
						<Card.Description>Real-time server activity and events</Card.Description>
					</Card.Header>
					<Card.Content>
						<!-- Log Filters -->
						<div class="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
							<!-- Level Filter -->
							<div class="flex flex-wrap gap-2">
								<Button
									variant={selectedLevel === 'all' ? 'default' : 'outline'}
									size="sm"
									onclick={() => (selectedLevel = 'all')}
								>
									All ({logLevelCounts().all})
								</Button>
								<Button
									variant={selectedLevel === 'info' ? 'default' : 'outline'}
									size="sm"
									onclick={() => (selectedLevel = 'info')}
								>
									<Icon icon="mdi:information-outline" class="mr-1 h-3.5 w-3.5" />
									Info ({logLevelCounts().info})
								</Button>
								<Button
									variant={selectedLevel === 'debug' ? 'default' : 'outline'}
									size="sm"
									onclick={() => (selectedLevel = 'debug')}
								>
									<Icon icon="mdi:bug-outline" class="mr-1 h-3.5 w-3.5" />
									Debug ({logLevelCounts().debug})
								</Button>
								<Button
									variant={selectedLevel === 'warn' ? 'default' : 'outline'}
									size="sm"
									onclick={() => (selectedLevel = 'warn')}
								>
									<Icon icon="mdi:alert-outline" class="mr-1 h-3.5 w-3.5" />
									Warn ({logLevelCounts().warn})
								</Button>
								<Button
									variant={selectedLevel === 'error' ? 'default' : 'outline'}
									size="sm"
									onclick={() => (selectedLevel = 'error')}
								>
									<Icon icon="mdi:alert-circle-outline" class="mr-1 h-3.5 w-3.5" />
									Error ({logLevelCounts().error})
								</Button>
							</div>

							<!-- Search -->
							<div class="relative w-full sm:w-64">
								<Icon
									icon="mdi:magnify"
									class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-muted-foreground"
								/>
								<input
									type="text"
									placeholder="Search logs..."
									bind:value={searchQuery}
									class="h-9 w-full rounded-md border border-input bg-background px-9 py-1 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:ring-1 focus-visible:ring-ring focus-visible:outline-none"
								/>
							</div>
						</div>

						<!-- Logs Table -->
						{#if filteredLogs().length > 0}
							<ScrollArea class="h-[calc(100vh-300px)] rounded-md border p-4">
								<Table.Root class="w-full caption-bottom text-sm">
									<Table.Header class="sticky top-0 z-10 bg-background/95 backdrop-blur-sm">
										<Table.Row>
											<Table.Head class="w-25">Time</Table.Head>
											<Table.Head class="w-20">Level</Table.Head>
											<Table.Head>Message</Table.Head>
											<Table.Head class="hidden md:table-cell">Details</Table.Head>
										</Table.Row>
									</Table.Header>

									<Table.Body>
										{#each filteredLogs() as log}
											<Table.Row>
												<Table.Cell class="font-mono text-xs">
													{formatTimeAgo(log.timestamp)}
												</Table.Cell>
												<Table.Cell>
													<Badge variant={getLogLevelVariant(log.level)}>
														<Icon icon={getLogLevelIcon(log.level)} class="mr-1 h-3 w-3" />
														{log.level}
													</Badge>
												</Table.Cell>
												<Table.Cell class="font-medium">
													{log.message}
													{#if log.error}
														<p class="mt-1 text-xs text-destructive">{log.error}</p>
													{/if}
												</Table.Cell>
												<Table.Cell class="hidden md:table-cell">
													<div class="space-y-1 text-xs text-muted-foreground">
														{#if log.port}
															<div><span class="font-semibold">Port:</span> {log.port}</div>
														{/if}
														{#if log.command}
															<div><span class="font-semibold">Command:</span> {log.command}</div>
														{/if}
														{#if log.duration}
															<div><span class="font-semibold">Duration:</span> {log.duration}</div>
														{/if}
														{#if log.sql}
															<div class="max-w-md truncate">
																<span class="font-semibold">SQL:</span>
																{log.sql}
															</div>
														{/if}
														{#if log.status}
															<div><span class="font-semibold">Status:</span> {log.status}</div>
														{/if}
														{#if log.method}
															<div>
																<span class="font-semibold">Method:</span>
																{log.method}
																{log.path}
															</div>
														{/if}
														{#if log.ip}
															<div><span class="font-semibold">IP:</span> {log.ip}</div>
														{/if}
														{#if log.count !== undefined}
															<div><span class="font-semibold">Count:</span> {log.count}</div>
														{/if}
														{#if log.pipeline_size}
															<div>
																<span class="font-semibold">Pipeline Size:</span>
																{log.pipeline_size}
															</div>
														{/if}
														{#if log.caller}
															<div><span class="font-semibold">Caller:</span> {log.caller}</div>
														{/if}
														{#if log.latency}
															<div><span class="font-semibold">Latency:</span> {log.latency}</div>
														{/if}
														{#if log.args}
															<div><span class="font-semibold">Args:</span> {log.args}</div>
														{/if}
														{#if log.user_agent}
															<div>
																<span class="font-semibold">User Agent:</span>
																{log.user_agent}
															</div>
														{/if}
													</div>
												</Table.Cell>
											</Table.Row>
										{/each}
									</Table.Body>
								</Table.Root>
							</ScrollArea>
						{:else}
							<Empty.Root>
								<Empty.Header>
									<Empty.Media>
										<Icon icon="mdi:text-box-search-outline" class="h-10 w-10" />
									</Empty.Media>
									<Empty.Title>No logs found</Empty.Title>
									<Empty.Description>
										{selectedLevel !== 'all' || searchQuery
											? 'Try adjusting your filters or search query'
											: 'Server logs will appear here'}
									</Empty.Description>
								</Empty.Header>
							</Empty.Root>
						{/if}
					</Card.Content>
					{#if paginationMeta && paginationMeta.total_pages > 1}
						<Card.Footer>
							<div class="flex w-full justify-end">
								<Pagination.Root
									count={paginationMeta.total_items}
									perPage={limit}
									bind:page={currentPage}
								>
									{#snippet children({ pages, range })}
										<Pagination.Content>
											<Pagination.Item>
												<Pagination.Previous />
											</Pagination.Item>
											{#each pages as page (page.key)}
												{#if page.type === 'ellipsis'}
													<Pagination.Item>
														<Pagination.Ellipsis />
													</Pagination.Item>
												{:else}
													<Pagination.Item>
														<Pagination.Link {page} isActive={currentPage === page.value}>
															{page.value}
														</Pagination.Link>
													</Pagination.Item>
												{/if}
											{/each}
											<Pagination.Item>
												<Pagination.Next />
											</Pagination.Item>
										</Pagination.Content>
									{/snippet}
								</Pagination.Root>
							</div>
						</Card.Footer>
					{/if}
				</Card.Root>
			</CardSpotlight>
		</div>
	</div>
	<AppAlertDialog
		bind:open={openClearLogDialog}
		type="warning"
		title="Warning"
		message="Are you sure you want to clear server logs?"
		labelClose="Cancel"
		labelAction="Confirm"
		onaction={handleClearLog}
		onclose={() => {
			openClearLogDialog = false;
		}}
	/>
</AdminSidebarLayout>
