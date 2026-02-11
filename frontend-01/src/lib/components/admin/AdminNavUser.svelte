<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar';
	import Icon from '@iconify/svelte';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Button } from '../ui/button';
	import { goto } from '$app/navigation';
	import { localizeHref } from '@/paraglide/runtime';

	let { user }: { user?: User | null } = $props();
	const sidebar = Sidebar.useSidebar();

	async function handleLogout() {
		await fetch(`/accounts/logout`, {
			method: 'POST',
			headers: {
				'X-Platform': 'web'
			},
			credentials: 'include'
		});

		await goto(localizeHref('/login'));
	}
</script>

<Sidebar.Menu>
	<Sidebar.MenuItem>
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Sidebar.MenuButton
						{...props}
						size="lg"
						class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
					>
						<Avatar.Root class="size-8 rounded-lg grayscale-25">
							<Avatar.Image src={user?.avatar_url} alt={user?.full_name} />
							<Avatar.Fallback class="rounded-lg">
								{user?.full_name?.split(' ')[0].slice(0, 2).toUpperCase() || 'CN'}
							</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{user?.full_name}</span>
							<span class="truncate text-xs text-muted-foreground">
								{user?.email}
							</span>
						</div>
						<Icon icon="mage:dots" class="ml-auto size-4" />
					</Sidebar.MenuButton>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
				side={sidebar.isMobile ? 'bottom' : 'right'}
				align="end"
				sideOffset={4}
			>
				<DropdownMenu.Label class="p-0 font-normal">
					<div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Image src={user?.avatar_url} alt={user?.full_name} />
							<Avatar.Fallback class="rounded-lg">
								{user?.full_name?.split(' ')[0].slice(0, 2).toUpperCase() || 'CN'}
							</Avatar.Fallback>
						</Avatar.Root>
						<div class="grid flex-1 text-left text-sm leading-tight">
							<span class="truncate font-medium">{user?.full_name}</span>
							<span class="truncate text-xs text-muted-foreground">
								{user?.email}
							</span>
						</div>
					</div>
				</DropdownMenu.Label>
				<DropdownMenu.Separator />
				<DropdownMenu.Group>
					<DropdownMenu.Item class="ps-5">
						{#snippet child()}
							<a
								href={localizeHref('/accounts/profile')}
								class="flex items-center gap-2 ps-3 text-sm transition-all active:scale-95"
							>
								<Icon icon="solar:user-circle-linear" />
								Profile
							</a>
						{/snippet}
					</DropdownMenu.Item>
				</DropdownMenu.Group>
				<DropdownMenu.Separator />
				<DropdownMenu.Item>
					{#snippet child()}
						<Button
							type="button"
							variant="ghost"
							class="text-sm text-red-600 dark:text-red-500"
							onclick={handleLogout}
						>
							<Icon icon="ic:outline-logout" />
							Log Out
						</Button>
					{/snippet}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	</Sidebar.MenuItem>
</Sidebar.Menu>
