<script lang="ts" module>
	interface MenuItem {
		id: number;
		title: string;
		url: string;
		icon: string;
	}
</script>

<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import { page } from '$app/state';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { Separator } from '@/components/ui/separator';
	import Icon from '@iconify/svelte';
	import { cn } from '@/utils';
	import { LightSwitch } from '@/components/ui-extras/light-switch';
	import { LanguageSwitcher } from '$lib/components/ui-extras/language-switcher/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import {
		localizeHref,
		getLocale,
		setLocale,
		type Locale,
		locales as availableLocales,
		isLocale
	} from '$lib/paraglide/runtime';
	import { LanguageLabels } from '@/utils/localize-path.js';

	let {
		ref = $bindable(null),
		user,
		setting,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & {
		user?: User | null;
		setting?: SettingsValue | null;
	} = $props();

	let currentLang = $derived(getLocale());
	let webSetting = $derived(setting?.WEBSITE);
	let logo = $derived(
		setting?.SYSTEM?.source_logo_favicon === 'remote' &&
			setting?.WEBSITE?.site_logo?.startsWith('http')
			? setting?.WEBSITE?.site_logo
			: '/images/logo.png'
	);

	const data = {
		navMain: [
			{
				id: 1,
				title: i18n.label_home(),
				url: localizeHref('/'),
				icon: 'bx:home-alt'
			},
			{
				id: 2,
				title: i18n.label_genre(),
				url: localizeHref('/genres'),
				icon: 'icon-park-outline:film'
			},
			{
				id: 3,
				title: i18n.label_popular(),
				url: localizeHref('/popular'),
				icon: 'ic:outline-star-purple500'
			},
			{
				id: 4,
				title: i18n.label_movie(),
				url: localizeHref('/movies'),
				icon: 'mingcute:movie-line'
			},
			{
				id: 5,
				title: i18n.label_tv_series(),
				url: localizeHref('/series'),
				icon: 'ic:baseline-ondemand-video'
			},
			{
				id: 6,

				title: i18n.label_anime(),
				url: localizeHref('/anime'),
				icon: 'simple-icons:myanimelist'
			},
			{
				id: 9,
				title: i18n.label_country(),
				url: localizeHref('/countries'),
				icon: 'material-symbols:globe'
			},
			{
				id: 10,
				title: i18n.label_year(),
				url: localizeHref('/years'),
				icon: 'material-symbols:date-range'
			}
		] satisfies MenuItem[]
	};
	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));

	function removeEndSlash(url: string) {
		return url.replace(/\/$/, '');
	}
</script>

<Sidebar.Root
	collapsible="icon"
	class="top-(--header-height) h-[calc(100svh-var(--header-height))]! group-data-[collapsible=icon]:w-16"
	{...restProps}
>
	<Sidebar.Header class="block md:hidden">
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a
							href={localizeHref('/admin/dashboard')}
							{...props}
							class="flex items-center gap-3 rounded-lg py-3"
						>
							<img src={logo} alt={webSetting?.site_name} class={cn('h-8 w-auto')} />
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Separator class="my-2 lg:hidden" />
	<Sidebar.Content
		class="scrollbar-thumb-cyan scrollbar-thin overflow-hidden overflow-y-auto scrollbar-thumb-foreground scrollbar-track-accent group-data-[collapsible=icon]:bg-background"
	>
		<Sidebar.Group>
			<Sidebar.GroupContent class="group-data-[collapsible=icon]:px-2">
				<Sidebar.Menu class="space-y-4 group-data-[collapsible=icon]:space-y-2">
					{#each data.navMain as item (item.title)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton
								tooltipContentProps={{
									hidden: false
								}}
								isActive={removeEndSlash(item.url) === removeEndSlash(page.url.pathname)}
								class="px-5 group-data-[collapsible=icon]:px-0"
							>
								{#snippet tooltipContent()}
									{item.title}
								{/snippet}
								{#snippet child({ props })}
									<a
										href={item.url}
										{...props}
										class={cn(
											'flex min-h-max flex-row items-center gap-3 pl-4 group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:gap-0 group-data-[collapsible=icon]:pl-0'
										)}
									>
										{#if item.icon}
											<Icon
												icon={item.icon}
												class={cn(
													'size-6',
													removeEndSlash(item.url) === removeEndSlash(page.url.pathname)
														? 'font-bold text-yellow-500 dark:text-yellow-400'
														: ''
												)}
											/>
											<p
												class={cn(
													'hidden text-center text-[10px] group-data-[collapsible=icon]:block',
													removeEndSlash(item.url) === removeEndSlash(page.url.pathname)
														? 'font-bold text-yellow-500 dark:text-yellow-400'
														: ''
												)}
											>
												{item.title}
											</p>
										{/if}
										<span
											class={cn(
												removeEndSlash(item.url) === removeEndSlash(page.url.pathname)
													? 'font-bold text-yellow-500 dark:text-yellow-400'
													: '',
												'group-data-[collapsible=icon]:hidden'
											)}
										>
											{item.title}
										</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
		<Separator class="my-2" />
		<Sidebar.Group class="mt-auto">
			<Sidebar.GroupContent class="group-data-[collapsible=icon]:px-2">
				<Sidebar.Menu class="space-y-4 group-data-[collapsible=icon]:space-y-2">
					<Sidebar.MenuItem>
						<Sidebar.MenuButton
							size="sm"
							tooltipContentProps={{
								hidden: false
							}}
							class="px-5 group-data-[collapsible=icon]:px-0"
						>
							{#snippet tooltipContent()}
								{i18n.label_theme()}
							{/snippet}
							{#snippet child({ props })}
								<div
									class="flex min-h-max flex-row items-center gap-3 pl-4 group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:gap-0 group-data-[collapsible=icon]:pl-0"
								>
									<LightSwitch />
									<span class={cn('group-data-[collapsible=icon]:hidden')}>
										{i18n.label_theme()}
									</span>
								</div>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
					<Sidebar.MenuItem>
						<Sidebar.MenuButton
							size="sm"
							tooltipContentProps={{
								hidden: false
							}}
							class="px-5 group-data-[collapsible=icon]:px-0"
						>
							{#snippet tooltipContent()}
								{i18n.label_language()}
							{/snippet}
							{#snippet child({ props })}
								<div
									class="flex min-h-max flex-row items-center gap-3 pl-4 group-data-[collapsible=icon]:flex-col group-data-[collapsible=icon]:gap-0 group-data-[collapsible=icon]:pl-0"
								>
									<LanguageSwitcher
										{languages}
										bind:value={currentLang}
										onChange={(code: string) => {
											if (isLocale(code)) setLocale(code);
										}}
									/>
									<span class={cn('group-data-[collapsible=icon]:hidden')}>
										{i18n.label_language()}
									</span>
								</div>
							{/snippet}
						</Sidebar.MenuButton>
					</Sidebar.MenuItem>
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer></Sidebar.Footer>
</Sidebar.Root>
