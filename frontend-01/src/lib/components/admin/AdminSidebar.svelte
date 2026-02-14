<script lang="ts">
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	import { AdminNavMain, AdminNavSetting, AdminNavUser, AdminNavBottom } from '@/components/admin';
	import { cn } from '@/utils';
	import { localizeHref } from '@/paraglide/runtime';

	let {
		user,
		setting,
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & {
		user?: User | null;
		setting?: SettingsValue | null;
	} = $props();

	let webSetting = $derived(setting?.WEBSITE);
	let logo = $derived(
		setting?.SYSTEM?.source_logo_favicon === 'remote' &&
			setting?.WEBSITE?.site_logo?.startsWith('http')
			? setting?.WEBSITE?.site_logo
			: '/images/logo.png'
	);
	let favicon = $derived(
		setting?.SYSTEM?.source_logo_favicon === 'remote' &&
			setting?.WEBSITE?.site_favicon?.startsWith('http')
			? setting?.WEBSITE?.site_favicon
			: '/images/icon.png'
	);

	const sidebar = useSidebar();

	const data = {
		navMain: [
			{
				id: 1,
				title: 'Dashboard',
				url: localizeHref('/admin/dashboard'),
				icon: 'material-symbols:home'
			},
			{
				id: 2,
				title: 'Applications',
				url: localizeHref('/admin/applications'),
				icon: 'tdesign:app-filled'
			},
			{
				id: 3,
				title: 'Users',
				url: localizeHref('/admin/users'),
				icon: 'material-symbols:account-circle'
			}
		] satisfies MenuItem[],
		navSetting: [
			{
				id: 4,
				title: 'Settings',
				url: '#',
				icon: 'mingcute:settings-2-fill',
				child: [
					{
						title: 'Website',
						url: localizeHref('/admin/settings/web'),
						icon: 'mdi:web'
					},
					{
						title: 'Email',
						url: localizeHref('/admin/settings/email'),
						icon: 'ri:mail-settings-fill'
					},
					{
						title: 'System',
						url: localizeHref('/admin/settings/system'),
						icon: 'solar:settings-linear'
					},
					{
						title: 'Monetization',
						url: localizeHref('/admin/settings/monetization'),
						icon: 'streamline:dollar-coin-remix'
					},
					{
						title: 'Ads.txt',
						url: localizeHref('/admin/settings/ads.txt'),
						icon: 'lsicon:file-txt-filled'
					},
					{
						title: 'Robot.txt',
						url: localizeHref('/admin/settings/robot.txt'),
						icon: 'lsicon:file-txt-filled'
					}
				]
			},
			{
				id: 5,
				title: 'Account',
				url: '#',
				icon: 'mdi:account-cog',
				child: [
					{
						title: 'Profile',
						url: localizeHref('/admin/accounts/profile'),
						icon: 'ic:outline-account-circle'
					},
					{
						title: 'Password',
						url: localizeHref('/admin/accounts/password'),
						icon: 'ic:round-key'
					}
				]
			}
		] satisfies MenuItem[],
		navBottom: [
			{
				id: 6,
				title: 'Server Status',
				url: localizeHref('/admin/server-status'),
				icon: 'ic:outline-monitor-heart'
			},
			{
				id: 7,
				title: 'User Panel',
				url: localizeHref('/user/account'),
				icon: 'material-symbols:account-circle'
			},
			{
				id: 8,
				title: 'Home',
				url: localizeHref('/'),
				icon: 'material-symbols:home'
			}
		] satisfies MenuItem[]
	};
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton class="data-[slot=sidebar-menu-button]:p-1.5!">
					{#snippet child({ props })}
						<a
							href={localizeHref('/admin/dashboard')}
							{...props}
							class="flex items-start gap-3 rounded-lg"
						>
							<img
								src={sidebar.state === 'expanded' ? logo : favicon}
								alt={webSetting?.site_name}
								class={cn(
									' rounded-lg',
									sidebar.state === 'expanded' ? 'h-8 w-auto' : 'aspect-square size-7 object-cover'
								)}
							/>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content
		class="scrollbar-thumb-cyan scrollbar-thin overflow-hidden overflow-y-auto scrollbar-thumb-foreground scrollbar-track-accent"
	>
		<AdminNavMain items={data.navMain} />
		<AdminNavSetting items={data.navSetting} />
		<AdminNavBottom items={data.navBottom} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<AdminNavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
