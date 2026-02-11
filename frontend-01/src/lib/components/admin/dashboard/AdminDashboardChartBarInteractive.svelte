<script lang="ts">
	import { CardSpotlight } from '@/components';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Chart from '$lib/components/ui/chart/index.js';
	import { BarChart, Highlight } from 'layerchart';
	import { cubicInOut } from 'svelte/easing';
	import { scaleUtc } from 'd3-scale';

	let { analytics, rangeDescription }: { analytics?: AnalyticsDaily[]; rangeDescription?: string } =
		$props();
	let chartData = $derived(analytics || []);

	const chartConfig = {
		total_downloads: { label: 'Total Downloads', color: 'var(--chart-1)' },
		total_revenue: { label: 'Total Revenue', color: 'var(--chart-2)' },
		total_users: { label: 'Total Users', color: 'var(--chart-3)' },
		active_users: { label: 'Active Users', color: 'var(--chart-4)' }
	} satisfies Chart.ChartConfig;

	let context = $state<any>();
	let activeChart = $state<keyof typeof chartConfig>('total_downloads');

	const activeSeries = $derived([
		{
			key: activeChart,
			label: chartConfig[activeChart].label,
			color: chartConfig[activeChart].color
		}
	]);

	function safeScaleTicks(scale: any) {
		try {
			if (!scale) return [];
			if (typeof scale.ticks === 'function') {
				return scale.ticks();
			}
			if (typeof scale.domain === 'function' && typeof scale.range === 'function') {
				const domain = scale.domain?.();
				const range = scale.range?.();
				if (domain && domain.length >= 2 && range && range.length >= 2) {
					return scaleUtc(domain, range).ticks();
				}
			}
			return [];
		} catch {
			return [];
		}
	}
</script>

<CardSpotlight variant="success" shadow="large" spotlightIntensity="medium" spotlight>
	<Card.Root class="@container/card bg-white/20 dark:bg-black/20">
		<Card.Header>
			<Card.Title>Daily Analytics</Card.Title>
			<Card.Description>
				<span class="hidden @[540px]/card:block">
					{rangeDescription || 'Total for the selected range'}
				</span>
			</Card.Description>
		</Card.Header>
		<Card.Content class="px-2 pt-4 sm:px-6 sm:pt-6">
			<Chart.Container config={chartConfig} class="aspect-auto h-62.5 w-full">
				<BarChart
					bind:context
					data={chartData}
					x="date"
					axis="x"
					series={activeSeries}
					props={{
						bars: {
							stroke: 'none',
							rounded: 'none',
							initialY: context?.height || 0,
							initialHeight: 0,
							motion: {
								y: { type: 'tween', duration: 500, easing: cubicInOut },
								height: { type: 'tween', duration: 500, easing: cubicInOut }
							}
						},
						highlight: { area: { fill: 'none' } },
						xAxis: {
							format: (d: unknown) => {
								const date = typeof d === 'string' ? new Date(d) : (d as Date);
								if (!date || isNaN(date.getTime())) return 'Invalid Date';
								return date.toLocaleDateString('en-US', { month: 'short', day: '2-digit' });
							},
							ticks: (scale: any) => safeScaleTicks(scale)
						},
						yAxis: {
							format: (v: number) => (isNaN(v) ? '0' : v.toLocaleString())
						}
					}}
				>
					{#snippet belowMarks()}
						<Highlight area={{ class: 'fill-muted' }} />
					{/snippet}
					{#snippet tooltip()}
						<Chart.Tooltip
							nameKey={activeChart}
							labelFormatter={(v: unknown) => {
								const date = typeof v === 'string' ? new Date(v) : (v as Date);
								if (!date || isNaN(date.getTime())) return 'Invalid Date';
								return date.toLocaleDateString('en-US', {
									month: 'short',
									day: 'numeric',
									year: 'numeric'
								});
							}}
						/>
					{/snippet}
				</BarChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
</CardSpotlight>
