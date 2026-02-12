<script lang="ts">
	import type { Snippet } from 'svelte';
	import { MaxWidthWrapper } from '@/components';
	import { LightSwitch } from '@/components/ui-extras/light-switch';
	import { localizeHref } from '$lib/paraglide/runtime';
	import { LanguageSwitcher } from '$lib/components/ui-extras/language-switcher/index.js';
	import {
		getLocale,
		setLocale,
		locales as availableLocales,
		isLocale
	} from '$lib/paraglide/runtime';
	import { LanguageLabels } from '@/utils/localize-path.js';

	let { children, setting }: { children?: Snippet; setting?: SettingsValue } = $props();

	let webSetting = $derived(setting?.WEBSITE);
	let logo = $derived(webSetting?.site_logo || '/images/logo.png');

	const languages = availableLocales.map((code) => ({
		code,
		label: LanguageLabels[code] ?? code.toUpperCase()
	}));
	let currentLang = $derived(getLocale());
</script>

<main
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
	<MaxWidthWrapper>
		<div class="relative">
			<div class="mx-auto flex w-full max-w-lg flex-col items-start pt-4 md:pt-20">
				<div class=" flex w-full items-center justify-center py-5 md:justify-start">
					<a href={localizeHref('/')} class="flex items-center gap-x-2">
						<img src={logo} alt={webSetting?.site_name || 'Tube XXI'} class="h-14 w-auto" />
					</a>
				</div>
				{@render children?.()}
			</div>
		</div>
	</MaxWidthWrapper>
</main>
