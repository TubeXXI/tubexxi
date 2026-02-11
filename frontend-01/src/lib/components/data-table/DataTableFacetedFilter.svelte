<script lang="ts" generics="TData, TValue">
	import type { ClassValue } from 'svelte/elements';
	import type { Component } from 'svelte';
	import type { Column } from '@tanstack/table-core';
	import { SvelteSet } from 'svelte/reactivity';
	import { Button } from '@/components/ui/button';
	import * as Command from '@/components/ui/command';
	import * as Popover from '@/components/ui/popover';
	import { Separator } from '@/components/ui/separator';
	import { Badge } from '@/components/ui/badge';
	import { CirclePlusIcon, CheckIcon } from '@lucide/svelte';
	import { cn } from '@/utils';

	let {
		column,
		title,
		options,
		class: classNames = '',
		onreset,
		onchange
	}: {
		column: Column<TData, TValue>;
		title: string;
		options: {
			label: string;
			value: string;
			icon?: Component;
		}[];
		class?: ClassValue;
		onreset?: () => Promise<void>;
		onchange?: (value: string) => void;
	} = $props();

	const facets = $derived(column?.getFacetedUniqueValues());
	const selectedValues = $derived(new SvelteSet(column?.getFilterValue() as string[]));
</script>

<Popover.Root>
	<Popover.Trigger>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="outline"
				size="sm"
				class={cn('h-8 w-auto border-dashed', classNames)}
			>
				<CirclePlusIcon />
				{title}
				{#if selectedValues.size > 0}
					<Separator orientation="vertical" class="mx-2 h-4" />
					<Badge variant="secondary" class="rounded-sm px-1 font-normal lg:hidden">
						{selectedValues.size}
					</Badge>
					<div class="hidden space-x-1 lg:flex">
						{#if selectedValues.size > 2}
							<Badge variant="secondary" class="rounded-sm px-1 font-normal">
								{selectedValues.size}
								'Selected'
							</Badge>
						{:else}
							{#each options.filter((opt) => selectedValues.has(opt.value)) as option (option)}
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
	<Popover.Content class="w-50 px-0" align="center">
		<Command.Root class="w-full">
			<Command.Input placeholder={title} class="h-8 w-full" />
			<Command.List>
				<Command.Empty>'No results found'</Command.Empty>
				<Command.Group>
					{#each options as option (option)}
						{@const isSelected = selectedValues.has(option.value)}
						<Command.Item
							onSelect={() => {
								if (isSelected) {
									selectedValues.delete(option.value);
								} else {
									selectedValues.add(option.value);
								}
								const filterValues = Array.from(selectedValues);
								column?.setFilterValue(filterValues.length ? filterValues : undefined);
								onchange?.(filterValues.join(','));
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
							{#if option.icon}
								{@const Icon = option.icon}
								<Icon class="text-muted-foreground" />
							{/if}

							<span>{option.label}</span>
							{#if facets?.get(option.value)}
								<span class="ml-auto flex size-4 items-center justify-center font-mono text-xs">
									{facets.get(option.value)}
								</span>
							{/if}
						</Command.Item>
					{/each}
				</Command.Group>
				{#if selectedValues.size > 0}
					<Command.Separator />
					<Command.Group>
						<Command.Item
							onSelect={() => {
								column?.setFilterValue(undefined);
								onreset?.();
							}}
							class="justify-center text-center"
						>
							'Clear filters'
						</Command.Item>
					</Command.Group>
				{/if}
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
