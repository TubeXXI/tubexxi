<script lang="ts">
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { localizeHref } from '@/paraglide/runtime';
	import * as i18n from '@/paraglide/messages.js';
	import { translateStore } from '@/stores';
	import Icon from '@iconify/svelte';

	let {
		setting,
		lang = 'en'
	}: {
		setting?: SettingsValue | null;
		lang?: string;
	} = $props();

	let webSetting = $derived(setting?.WEBSITE);
	let systemSetting = $derived(setting?.SYSTEM);
	// svelte-ignore state_referenced_locally
	let translatedSiteDescription = $state(webSetting?.site_description || '');
	let translateLoading = $state(false);

	let logo = $derived(
		setting?.SYSTEM?.source_logo_favicon === 'remote' &&
			setting?.WEBSITE?.site_logo?.startsWith('http') &&
			setting?.WEBSITE?.site_logo !== ''
			? setting?.WEBSITE?.site_logo
			: '/images/logo.png'
	);

	const handleTranslate = async (value: string) => {
		translateLoading = true;
		try {
			const result = await translateStore.singleTranslate(value, {
				targetLang: lang,
				useCache: true
			});
			translatedSiteDescription = result.data.target.text;
		} catch (error) {
			console.error(error);
			return value;
		} finally {
			translateLoading = false;
		}
	};

	onMount(async () => {
		if (webSetting?.site_description) {
			await handleTranslate(webSetting?.site_description);
		}
	});

	const socialLink = [
		{
			title: 'Facebook',
			url: 'https://www.facebook.com',
			icon: 'logos:facebook'
		},
		{
			title: 'Twitter',
			url: 'https://www.twitter.com',
			icon: 'logos:twitter'
		},
		{
			title: 'Instagram',
			url: 'https://www.instagram.com',
			icon: 'skill-icons:instagram'
		},
		{
			title: 'YouTube',
			url: 'https://www.youtube.com',
			icon: 'logos:youtube-icon'
		}
	];
</script>

<div class="@container/main flex flex-col gap-4 overflow-hidden py-4 pl-0 lg:pl-10">
	<div class="grid grid-cols-1 gap-6 px-4 py-4 md:grid-cols-3 lg:px-10 lg:py-6">
		<div class="col-span-1 flex flex-col items-center lg:items-start">
			<a href={localizeHref('/')} rel="noopener noreferrer">
				<img src={logo} alt={webSetting?.site_name} class="h-12 w-auto" />
			</a>
			<div class="mt-4 space-y-4 lg:pl-2">
				{#if translateLoading}
					<p class="line-clamp-4 text-sm">{webSetting?.site_description || ''}</p>
				{:else}
					<p class="text-sm">{translatedSiteDescription || ''}</p>
				{/if}
				{#if systemSetting?.play_store_app_url}
					<a
						href={systemSetting?.play_store_app_url}
						target="_blank"
						rel="noopener noreferrer"
						class="flex max-w-max rounded-md bg-neutral-100 p-2 shadow-xl/30 shadow-neutral-500 backdrop-blur-md dark:bg-neutral-800 dark:shadow-neutral-400"
					>
						<img src="/images/play-store.png" alt="" class="h-10 w-auto" />
					</a>
				{/if}
				{#if systemSetting?.app_store_app_url}
					<a
						href={systemSetting?.app_store_app_url}
						target="_blank"
						rel="noopener noreferrer"
						class="flex max-w-max rounded-md bg-neutral-100 p-2 shadow-xl/30 shadow-neutral-500 backdrop-blur-md dark:bg-neutral-800 dark:shadow-neutral-400"
					>
						<img src="/images/app-store.png" alt="" class="h-10 w-auto" />
					</a>
				{/if}
			</div>
		</div>
		<div class="col-span-1 flex flex-col items-start">
			<h3 class="font-roboto text-base font-semibold text-white">{i18n.page()}</h3>
			<ul class="list-none space-y-0">
				<li>
					<a
						href={localizeHref('/about')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_about()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/contact')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_contact()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/terms')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_terms_conditions()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/privacy')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_privacy_notice()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/faq')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_faq()}
					</a>
				</li>
				<li>
					<a
						href="/sitemap.xml"
						data-sveltekit-preload-data="off"
						rel="sitemap"
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_sitemap()}
					</a>
				</li>
			</ul>
		</div>
		<div class="col-span-1 flex flex-col items-start">
			<h3 class="font-roboto text-base font-semibold text-white">{i18n.featured()}</h3>
			<ul class="list-none space-y-0">
				<li>
					<a
						href={localizeHref('/genres')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_genre()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/popular')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_popular()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/movies')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_movie()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/series')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_tv_series()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/anime')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_anime()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/countries')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_country()}
					</a>
				</li>
				<li>
					<a
						href={localizeHref('/years')}
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_year()}
					</a>
				</li>
				<li>
					<a
						href="/sitemap.xml"
						data-sveltekit-preload-data="off"
						rel="sitemap"
						class="text-sm text-muted-foreground transition-colors hover:text-blue-600 dark:hover:text-blue-400"
					>
						{i18n.label_sitemap()}
					</a>
				</li>
			</ul>
		</div>
	</div>
	<div class="border-t border-t-neutral-200 px-4 py-4 lg:px-10 lg:py-6 dark:border-t-neutral-700">
		<div class="flex flex-col gap-4 pb-10 md:flex-row md:justify-between">
			<div class="flex flex-wrap gap-4">
				<div class="flex items-center justify-center gap-4">
					{#each socialLink as item}
						<a
							href={item.url}
							target="_blank"
							rel="noopener noreferrer"
							class="flex items-center justify-center"
						>
							<Icon icon={item.icon} class="h-6 w-6" />
						</a>
					{/each}
				</div>
			</div>
			<div class="text-sm text-muted-foreground">
				© {new Date().getFullYear()}
				{webSetting?.site_name}. {i18n.all_rights_reserved()}.
			</div>
		</div>
	</div>
</div>
