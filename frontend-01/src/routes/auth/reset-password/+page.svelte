<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.resetForm, {
		async onSubmit(input) {
			errorMessage = undefined;
			successMessage = undefined;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			successMessage = event.result.data.message;
			await invalidateAll();
		},
		onError(event) {
			errorMessage = event.result.error.message;
		}
	});
</script>

<MetaTags {...metaTags} />
<data.component setting={data.settings}>
	<div class="flex w-full flex-col items-start gap-y-6 rounded-lg bg-accent px-5 py-8">
		<div class="w-full">
			<h1 class="text-2xl font-bold text-white">
				{i18n.page_forgot_password()}
			</h1>
			<p class="text-sm text-muted-foreground">
				{i18n.page_forgot_password_description()}
			</p>
		</div>
		{#if errorMessage}
			<Alert.Root variant="destructive">
				<Icon icon="mingcute:warning-line" class="size-4" />
				<Alert.Title>Error</Alert.Title>
				<Alert.Description>{errorMessage}</Alert.Description>
			</Alert.Root>
		{/if}
		{#if successMessage}
			<Alert.Root variant="default">
				<Icon icon="mingcute:check-line" class="size-4" />
				<Alert.Title>Success</Alert.Title>
				<Alert.Description>{successMessage}</Alert.Description>
			</Alert.Root>
		{:else}
			<form method="POST" class="w-full" use:enhance>
				<Field.Group class="Root">
					<Field.Field>
						<Field.Label for="email" class="capitalize">
							{i18n.email()}
							<span class="text-red-500 dark:text-red-400">*</span>
						</Field.Label>
						<div class="relative">
							<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
							<Input
								bind:value={$form.email}
								name="email"
								type="email"
								class="ps-10"
								placeholder={i18n.email()}
								aria-invalid={!!$errors.email}
								autocomplete="email"
							/>
						</div>
						{#if $errors.email}
							<Field.Error>{$errors.email}</Field.Error>
						{/if}
					</Field.Field>
					<Field.Field>
						<Button type="submit" disabled={$submitting}>
							{#if $submitting}
								<Spinner />
							{/if}
							{$submitting ? i18n.please_wait() : i18n.reset_password()}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		{/if}
		<div class="flex w-full flex-col items-center gap-2 pt-2">
			<div class="text-sm text-muted-foreground">
				{i18n.remember_me()}
			</div>
			<Button
				href={localizeHref('/auth/login')}
				type="button"
				variant="outline"
				class="w-full text-sm"
			>
				{i18n.sign_in()}
			</Button>
		</div>
	</div>
	<div class="rounded-lg px-5 py-8 text-neutral-300">
		<div class="flex w-full flex-col items-start">
			<p class="text-sm">
				{i18n.page_forgot_password_terms()}
				<a href={localizeHref('/terms')}>
					{i18n.terms_of_service()}
				</a>
				{i18n.and()}
				<a href={localizeHref('/privacy')}>
					{i18n.privacy_policy()}
				</a>
			</p>
		</div>
	</div>
</data.component>
