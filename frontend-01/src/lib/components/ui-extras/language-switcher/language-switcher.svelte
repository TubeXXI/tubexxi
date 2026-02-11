<script lang="ts">
	import GlobeIcon from '@lucide/svelte/icons/globe';
	import * as DropdownMenu from '$lib/components/ui-extras/dropdown-menu';
	import { buttonVariants } from '$lib/components/ui-extras/button';
	import { cn } from '$lib/utils.js';
	import { LanguageLabels } from '@/utils/localize-path.js';
	import type { LanguageSwitcherProps } from './types';

	let {
		languages = [],
		value = $bindable(''),
		align = 'end',
		variant = 'outline',
		onChange,
		class: className
	}: LanguageSwitcherProps = $props();

	function setFlagUrl(code?: string): string {
		if (!code) return '';
		switch (code) {
			case 'en':
				return 'https://flagicons.lipis.dev/flags/4x3/us.svg';
			case 'ja':
				return 'https://flagicons.lipis.dev/flags/4x3/jp.svg';
			case 'ar':
				return 'https://flagicons.lipis.dev/flags/4x3/sa.svg';
			case 'zh':
				return 'https://flagicons.lipis.dev/flags/4x3/cn.svg';
			case 'hi':
				return 'https://flagicons.lipis.dev/flags/4x3/in.svg';
			case 'el':
				return 'https://flagicons.lipis.dev/flags/4x3/gr.svg';
			default:
				return 'https://flagicons.lipis.dev/flags/4x3/' + code + '.svg';
		}
	}
	// set default code if there isn't one selected
	// svelte-ignore state_referenced_locally
	if (value === '') {
		value = languages[0].code;
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={cn(buttonVariants({ variant, size: 'icon' }), className)}
		aria-label="Change language"
	>
		{#if value}
			<img
				src={setFlagUrl(value)}
				alt={LanguageLabels[value as keyof typeof LanguageLabels] || value.toUpperCase()}
				class="size-4"
			/>
		{:else}
			<GlobeIcon class="size-4" />
		{/if}
		<span class="sr-only">Change language</span>
	</DropdownMenu.Trigger>
	<DropdownMenu.Content {align}>
		<DropdownMenu.RadioGroup bind:value onValueChange={onChange}>
			{#each languages as language (language.code)}
				<DropdownMenu.RadioItem value={language.code} class="flex items-center gap-2">
					<img src={setFlagUrl(language.code)} alt={language.label} class="size-4" />
					{language.label}
				</DropdownMenu.RadioItem>
			{/each}
		</DropdownMenu.RadioGroup>
	</DropdownMenu.Content>
</DropdownMenu.Root>
