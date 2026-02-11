<script lang="ts">
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import Icon from '@iconify/svelte';
	import * as Collapsible from '$lib/components/ui/collapsible';

	let {
		items
	}: {
		items: MenuItem[];
	} = $props();

	let selectedParentIndex = $state<number | null>(null);
	let selectedChildIndex = $state<number | null>(null);

	const toggleCollapse = (index: number) => {
		if (selectedParentIndex === index) {
			selectedParentIndex = null;
			selectedChildIndex = null;
		} else {
			selectedParentIndex = index;
			selectedChildIndex = null;
		}
	};

	const selectChild = (parentIndex: number, childIndex: number) => {
		selectedParentIndex = parentIndex;
		selectedChildIndex = childIndex;
	};
</script>

<Sidebar.Group>
	<Sidebar.GroupLabel>Settings</Sidebar.GroupLabel>
	<Sidebar.Menu>
		{#each items as item (item.id)}
			<Collapsible.Root open={selectedParentIndex === item.id} class="group/collapsible">
				{#snippet child({ props })}
					<Sidebar.MenuItem {...props}>
						<Collapsible.Trigger>
							{#snippet child({ props })}
								<Sidebar.MenuButton
									{...props}
									tooltipContent={item.title}
									onclick={() => toggleCollapse(item.id)}
									class={selectedParentIndex === item.id ? 'bg-accent' : ''}
								>
									{#snippet child({ props })}
										<a {...props} href={item.url}>
											{#if item.icon}
												<Icon
													icon={item.icon}
													class={item.url === page.url.pathname
														? 'font-bold text-sky-400 dark:text-yellow-500'
														: ''}
												/>
											{/if}
											<span
												class={item.url === page.url.pathname
													? 'font-bold text-sky-400 dark:text-yellow-500'
													: ''}
											>
												{item.title}
											</span>
											{#if item.child?.length}
												<Icon
													icon="lucide:chevron-right"
													class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"
												/>
											{/if}
										</a>
									{/snippet}
								</Sidebar.MenuButton>
							{/snippet}
						</Collapsible.Trigger>
						{#if item.child && item.child.length > 0}
							<Collapsible.Content>
								<Sidebar.MenuSub>
									{#each item.child as sub, childIndex (sub.url)}
										<Sidebar.MenuSubItem>
											<Sidebar.MenuSubButton
												class={selectedChildIndex === childIndex ? 'bg-accent' : ''}
												onclick={() => selectChild(item.id, childIndex)}
											>
												{#snippet child({ props })}
													<a href={sub.url} {...props}>
														{#if sub.icon}
															<Icon
																icon={sub.icon}
																class={page.url.pathname.startsWith(sub.url)
																	? 'font-bold text-sky-400 dark:text-yellow-500'
																	: ''}
															/>
														{/if}
														<span
															class={page.url.pathname.startsWith(sub.url)
																? 'font-bold text-sky-400 dark:text-yellow-500'
																: ''}
														>
															{sub.title}
														</span>
													</a>
												{/snippet}
											</Sidebar.MenuSubButton>
										</Sidebar.MenuSubItem>
									{/each}
								</Sidebar.MenuSub>
							</Collapsible.Content>
						{/if}
					</Sidebar.MenuItem>
				{/snippet}
			</Collapsible.Root>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
