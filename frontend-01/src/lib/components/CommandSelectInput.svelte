<script lang="ts">
	import type { ClassValue } from 'svelte/elements';
	import { tick } from 'svelte';
	import { SvelteSet } from 'svelte/reactivity';
	import { Button } from '@/components/ui/button';
	import * as Command from '@/components/ui/command';
	import * as Popover from '@/components/ui/popover';
	import { Separator } from '@/components/ui/separator';
	import { Badge } from '@/components/ui/badge';
	import { CirclePlusIcon, CheckIcon } from '@lucide/svelte';
	import { cn } from '@/utils';

	let {
		options,
		selectedValue = $bindable(new SvelteSet<string>([])),
		title = $bindable(''),
		onchange,
		class: className = '',
		variant = 'outline'
	}: {
		options: {
			label: string;
			value: string;
		}[];
		selectedValue?: SvelteSet<string>;
		title?: string;
		onchange: (value: SvelteSet<string>) => Promise<void>;
		class?: ClassValue;
		variant?: 'link' | 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost';
	} = $props();

	let openFilter = $state(false);

	function closeAndFocusTrigger(triggerId: string) {
		openFilter = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}
</script>

<Popover.Root bind:open={openFilter}>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				{variant}
				size="sm"
				class={cn('h-8 w-full border-dashed px-2 lg:w-auto', className)}
			>
				<CirclePlusIcon />
				{title || 'Filter'}
				{#if selectedValue.size > 0}
					<Separator orientation="vertical" class="mx-2 h-4" />
					<Badge variant="secondary" class="rounded-sm px-1 font-normal lg:hidden">
						{selectedValue.size}
					</Badge>
					<div class="hidden space-x-1 lg:flex">
						{#if selectedValue.size > 2}
							<Badge variant="secondary" class="rounded-sm px-1 font-normal">
								{selectedValue.size}
								{'selected'}
							</Badge>
						{:else}
							{#each options.filter((opt) => selectedValue.has(opt.value)) as option (option)}
								<Badge variant="secondary" class="rounded-sm px-1 font-normal">
									{option.label}
								</Badge>
							{/each}
						{/if}
					</div>
				{/if}
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content class="w-full p-0 lg:w-50" align="center">
		<Command.Root class="w-full">
			<Command.Input placeholder={'Filter'} class="h-8 w-full" />
			<Command.List>
				<Command.Empty>{'No results found.'}</Command.Empty>
				<Command.Group>
					{#each options as option (option.value)}
						{@const isSelected = selectedValue.has(option.value)}
						<Command.Item
							value={option.value}
							onSelect={() => {
								if (isSelected) {
									selectedValue.delete(option.value);
								} else {
									selectedValue.add(option.value);
								}
								closeAndFocusTrigger('command-filter-trigger');
								onchange?.(selectedValue);
							}}
						>
							<div
								class={cn(
									'mr-2 flex size-4 items-center justify-center rounded-sm border border-primary',
									isSelected ? 'bg-primary text-primary-foreground' : 'opacity-50 [&_svg]:invisible'
								)}
							>
								<CheckIcon class="size-4" />
							</div>
							<span>{option.label}</span>
						</Command.Item>
					{/each}
				</Command.Group>
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
