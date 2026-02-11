<script lang="ts" module>
	interface NavigationMenuItem {
		label: string;
		name: string;
		description?: string;
		icon: string;
		to: string;
		exact?: boolean;
		target?: string;
	}
</script>

<script lang="ts">
	import { goto } from '$app/navigation';
	import type { ClassValue } from 'svelte/elements';
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button/index.js';
	import { ScrollArea } from '@/components/ui/scroll-area';
	import * as Select from '@/components/ui/select';
	import Icon from '@iconify/svelte';
	import { cn } from '@/utils';

	let {
		links,
		class: className
	}: {
		links: NavigationMenuItem[];
		class?: ClassValue;
	} = $props();

	let selectedMenu = $state<string>(page.url.pathname);

	function handleNavigate(link: string) {
		selectedMenu = link;
		goto(link);
	}

	const triggerContent = $derived(links.find((f) => f.to === selectedMenu)?.label ?? 'Select menu');
</script>

<div class="p-1 md:hidden">
	<Select.Root bind:value={selectedMenu} type="single" onValueChange={handleNavigate}>
		<Select.Trigger class="h-12 w-full">
			{triggerContent}
		</Select.Trigger>
		<Select.Content>
			{#each links as item}
				<Select.Item value={item.to}>
					<div class="flex items-center gap-4 px-4 py-2">
						<Icon icon={item.icon} class="scale-125" />
						<span class="text-sm font-medium">{item.label}</span>
					</div>
				</Select.Item>
			{/each}
		</Select.Content>
	</Select.Root>
</div>
<ScrollArea
	orientation="horizontal"
	type="always"
	class="hidden w-full min-w-40 flex-col bg-background px-1 py-2 md:flex"
>
	<nav class={cn('flex space-x-2 py-1 lg:flex-col lg:space-y-1 lg:space-x-0', className)}>
		{#each links as item}
			<Button
				href={item.to}
				variant="link"
				class={cn(
					page.url.pathname === item.to ? 'bg-muted hover:bg-accent' : 'hover:bg-accent',
					'justify-start'
				)}
			>
				<Icon icon={item.icon} class="me-1" />
				<span class="line-clamp-1 text-sm font-medium">{item.label}</span>
			</Button>
		{/each}
	</nav>
</ScrollArea>
