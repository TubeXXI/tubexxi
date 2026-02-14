<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import {
		localizeHref,
		getLocale,
		setLocale,
		type Locale,
		locales as availableLocales,
		isLocale
	} from '@/paraglide/runtime';
	import { LanguageLabels } from '@/utils/localize-path.js';
	import { LightSwitch } from '@/components/ui-extras/light-switch';
	import { LanguageSwitcher } from '$lib/components/ui-extras/language-switcher/index.js';

	let {
		children,
		setting
	}: {
		children?: Snippet;
		setting?: SettingsValue | null;
	} = $props();

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));
	let currentLang = $derived(getLocale());

	async function handleLanguageChange(code?: string) {
		if (!code || !isLocale(code)) return;

		setLocale(code);

		const rawPath = removeLocaleFromPath(page.url.pathname);
		const localized = localizeHref(rawPath, { locale: code });

		await goto(localized);
	}

	function removeLocaleFromPath(path: string) {
		const parts = path.split('/');
		if (parts[1] && availableLocales.includes(parts[1] as Locale)) {
			return '/' + parts.slice(2).join('/');
		}
		return path;
	}
</script>

<div
	class="relative flex min-h-screen w-full items-center justify-center bg-[linear-gradient(rgba(0,0,0,0.5),rgba(0,0,0,0.5)),url('/images/app-bg.jpg')] bg-cover bg-fixed bg-center bg-no-repeat"
>
	<div class="absolute top-4 right-4 z-10">
		<div class="flex items-center gap-2">
			<LightSwitch />
			<LanguageSwitcher
				{languages}
				bind:value={currentLang}
				class="cursor-pointer"
				onChange={(code: string) => {
					if (isLocale(code)) setLocale(code);
				}}
			/>
		</div>
	</div>
	<div class="m-auto w-full">
		{#if children}
			{@render children()}
		{/if}
	</div>
</div>
