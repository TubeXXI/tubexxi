<script lang="ts">
	import type { Snippet } from 'svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { DefaultSidebar } from '$lib/components/themes/index.js';
	import { localizeHref } from '@/paraglide/runtime';
	import { DefaultSidebarSearch, DefaultSidebarUserMenu } from '$lib/components/themes/index.js';

	let {
		user,
		setting,
		children
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
		children?: Snippet<[]>;
	} = $props();

	let isOpen = $state(false);

	let webSetting = $derived(setting?.WEBSITE);

	let logo = $derived(
		setting?.SYSTEM?.source_logo_favicon === 'remote' &&
			setting?.WEBSITE?.site_logo?.startsWith('http')
			? setting?.WEBSITE?.site_logo
			: '/images/logo.png'
	);
</script>

<main class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider bind:open={isOpen} class="flex flex-col">
		<header class="sticky top-0 z-50 flex w-full items-center border-b bg-background">
			<div class="flex h-(--header-height) w-full items-center justify-between gap-2 px-2">
				<div class="flex items-center gap-1 lg:gap-2 lg:pl-2">
					<Sidebar.Trigger class="size-8" />
					<a href={localizeHref('/')} class="hidden md:inline-block">
						<img src={logo} alt={webSetting?.site_name ?? 'Indoxxi'} class="h-6 w-auto" />
					</a>
				</div>
				<DefaultSidebarSearch />
				<DefaultSidebarUserMenu {user} {setting} />
			</div>
		</header>
		<div class="flex flex-1">
			<DefaultSidebar {user} {setting} />
			<Sidebar.Inset>
				<div class="h-full bg-muted p-4">
					{@render children?.()}
				</div>
			</Sidebar.Inset>
		</div>
	</Sidebar.Provider>
</main>
