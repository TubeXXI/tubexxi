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
				url: localizeHref('/dashboard'),
				icon: 'material-symbols:home'
			},
			{
				id: 2,
				title: 'Platforms',
				url: localizeHref('/platform'),
				icon: 'tdesign:control-platform-filled'
			},
			{
				id: 3,
				title: 'Applications',
				url: localizeHref('/application'),
				icon: 'tdesign:app-filled'
			},
			{
				id: 4,
				title: 'Downloads',
				url: localizeHref('/download'),
				icon: 'ic:baseline-download'
			},
			{
				id: 5,
				title: 'Subscriptions',
				url: localizeHref('/subscription'),
				icon: 'mdi:payment-clock'
			},
			{
				id: 7,
				title: 'Users',
				url: localizeHref('/users'),
				icon: 'material-symbols:account-circle'
			}
		] satisfies MenuItem[],
		navSetting: [
			{
				id: 7,
				title: 'Settings',
				url: '#',
				icon: 'mingcute:settings-2-fill',
				child: [
					{
						title: 'Website',
						url: localizeHref('/settings/web'),
						icon: 'mdi:web'
					},
					{
						title: 'Email',
						url: localizeHref('/settings/email'),
						icon: 'ri:mail-settings-fill'
					},
					{
						title: 'System',
						url: localizeHref('/settings/system'),
						icon: 'solar:settings-linear'
					},
					{
						title: 'Monetization',
						url: localizeHref('/settings/monetization'),
						icon: 'streamline:dollar-coin-remix'
					},
					{
						title: 'Ads.txt',
						url: localizeHref('/settings/ads.txt'),
						icon: 'lsicon:file-txt-filled'
					},
					{
						title: 'Robot.txt',
						url: localizeHref('/settings/robot.txt'),
						icon: 'lsicon:file-txt-filled'
					}
				]
			},
			{
				id: 8,
				title: 'Account',
				url: '#',
				icon: 'mdi:account-cog',
				child: [
					{
						title: 'Profile',
						url: localizeHref('/accounts/profile'),
						icon: 'ic:outline-account-circle'
					},
					{
						title: 'Password',
						url: localizeHref('/accounts/password'),
						icon: 'ic:round-key'
					}
				]
			},
			{
				id: 9,
				title: 'Cookies',
				url: localizeHref('/cookies'),
				icon: 'fluent:cookies-16-filled'
			}
		] satisfies MenuItem[],
		navBottom: [
			{
				id: 10,
				title: 'Server Status',
				url: localizeHref('/server-status'),
				icon: 'ic:outline-monitor-heart'
			},
			{
				id: 11,
				title: 'User Panel',
				url: localizeHref('/user'),
				icon: 'material-symbols:account-circle'
			},
			{
				id: 12,
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
							href={localizeHref('/dashboard')}
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
							<span
								class={cn('text-lg font-bold', sidebar.state === 'expanded' ? 'block' : 'hidden')}
							>
								{webSetting?.site_name}
							</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content
		class="scrollbar-thumb-cyan scrollbar-thin scrollbar-thumb-foreground scrollbar-track-accent overflow-hidden overflow-y-auto"
	>
		<AdminNavMain items={data.navMain} />
		<AdminNavSetting items={data.navSetting} />
		<AdminNavBottom items={data.navBottom} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<AdminNavUser {user} />
	</Sidebar.Footer>
</Sidebar.Root>
