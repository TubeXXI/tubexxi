<script lang="ts">
	import { invalidateAll, goto } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Password from '$lib/components/ui-extras/password';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import type { E164Number } from 'svelte-tel-input/types';
	import type { ZxcvbnResult } from '@zxcvbn-ts/core';
	import { PhoneInput } from '@/components/ui-extras/phone-input/index.js';
	import type { CountryCode } from 'svelte-tel-input/types';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime';
	import { firebaseClient } from '@/client/firebase_client.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let phoneInput = $state<E164Number | undefined>('');
	let passwordInput = $state<string | undefined>('');
	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let isProcessing = $state(false);
	let fieldErrors = $state<{
		full_name?: string;
		email?: string;
		phone?: string;
		password?: string;
	}>({});
	let strength = $state<ZxcvbnResult>();

	async function createSession(idToken: string): Promise<{ redirect_url: string } | null> {
		const res = await fetch('/api/auth/session', {
			method: 'POST',
			credentials: 'include',
			headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
			body: JSON.stringify({ idToken })
		});
		const payload = (await res.json().catch(() => null)) as ApiResponse<{
			redirect_url: string;
			user: User | null;
		}> | null;
		if (!res.ok || !payload?.success || !payload.data?.redirect_url) {
			throw new Error(payload?.message || i18n.page_sign_in_error_message());
		}
		return payload.data;
	}

	async function handleRegisterSubmit() {
		fieldErrors = {};
		errorMessage = undefined;
		successMessage = undefined;
		isProcessing = true;
		try {
			if (!$form.full_name?.trim()) {
				fieldErrors = { ...fieldErrors, full_name: i18n.invalid_input() };
				throw new Error(i18n.invalid_input());
			}
			if (!$form.email?.trim()) {
				fieldErrors = { ...fieldErrors, email: i18n.invalid_input() };
				throw new Error(i18n.invalid_input());
			}
			const phone = `${phoneInput || ''}`.trim();
			if (!phone) {
				fieldErrors = { ...fieldErrors, phone: i18n.invalid_input() };
				throw new Error(i18n.invalid_input());
			}
			if (!$form.password?.trim()) {
				fieldErrors = { ...fieldErrors, password: i18n.invalid_input() };
				throw new Error(i18n.invalid_input());
			}

			const result = await firebaseClient?.registerWithEmail(
				$form.email,
				$form.password,
				$form.full_name
			);
			if (!result) {
				throw new Error(i18n.page_sign_up_error_message());
			}
			const idToken = await result.user.getIdToken();
			if (!idToken) {
				throw new Error(i18n.page_sign_up_error_message());
			}

			await createSession(idToken);
			await fetch('/api/auth/verify-email', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email: $form.email })
			});
			successMessage = i18n.page_sign_up_success_message();
			await invalidateAll();
			await goto(localizeHref(`/auth/verify-email?email=${encodeURIComponent($form.email)}`));
		} catch (error) {
			errorMessage = error instanceof Error ? error.message : i18n.invalid_input();
		} finally {
			isProcessing = false;
		}
	}

	// svelte-ignore state_referenced_locally
	const { form } = superForm(data.registerForm, {
		resetForm: true,
		dataType: 'json'
	});

	async function handleSocialLogin(provider: SocialProvider) {
		isProcessing = true;
		errorMessage = undefined;

		try {
			const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
				navigator.userAgent
			);

			const result = await firebaseClient?.signInWithSocial(provider, !isMobile);

			if (!result) {
				throw new Error(i18n.page_sign_up_error_message());
			}

			const idToken = await result.user.getIdToken();

			const response = await fetch('/api/auth/social-login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					idToken
				})
			});

			const responseData = await response.json();

			if (!response.ok) {
				throw new Error(responseData.message || i18n.page_sign_up_error_message());
			}

			successMessage = responseData.message;
			await invalidateAll();
			await goto(localizeHref(responseData.data?.redirect_url || '/'));
		} catch (error: any) {
			errorMessage = error.message || i18n.page_sign_up_error_message();
			isProcessing = false;
		}
	}

	$effect(() => {
		firebaseClient?.getRedirectResult().then(async (result) => {
			if (result) {
				isProcessing = true;
				try {
					const idToken = await result.user.getIdToken();

					const response = await fetch('/api/auth/social-login', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
						body: JSON.stringify({
							idToken
						})
					});

					const responseData = await response.json();

					if (response.ok) {
						successMessage = responseData.message;
						await invalidateAll();
						await goto(localizeHref(responseData.data?.redirect_url || '/'));
					} else {
						errorMessage = responseData.message;
					}
				} catch (err: any) {
					errorMessage = err.message || 'Login failed';
				} finally {
					isProcessing = false;
				}
			}
		});
	});
</script>

<MetaTags {...metaTags} />
<data.component setting={data.settings}>
	<div class="flex w-full flex-col items-start gap-y-6 rounded-lg bg-accent px-5 py-8">
		<div class="w-full">
			<h2 class="text-2xl font-semibold">
				{i18n.page_sign_up()}
			</h2>
			<p class="text-sm text-muted-foreground">
				{i18n.page_sign_up_description()}
			</p>
		</div>
		{#if errorMessage}
			<Alert.Root variant="destructive">
				<Icon icon="mingcute:warning-line" class="size-4" />
				<Alert.Title>{i18n.error()}</Alert.Title>
				<Alert.Description>{errorMessage}</Alert.Description>
			</Alert.Root>
		{/if}
		{#if successMessage}
			<Alert.Root
				variant="default"
				class="w-full rounded-md border border-green-500 dark:border-green-400"
			>
				<Icon icon="mingcute:check-line" class="size-4 text-green-500 dark:text-green-400" />
				<Alert.Title class="text-green-500 dark:text-green-400">{i18n.success()}</Alert.Title>
				<Alert.Description class="text-green-500 dark:text-green-400">
					{successMessage}
				</Alert.Description>
			</Alert.Root>
			<Button type="button" variant="default" href="/auth/login">
				{i18n.button_go_back_to_sign_in()}
			</Button>
		{/if}
		<form
			class="w-full"
			onsubmit={async (e) => {
				e.preventDefault();
				await handleRegisterSubmit();
			}}
		>
			<Field.Group class="Root">
				<Field.Field>
					<Field.Label for="full_name" class="capitalize">
						{i18n.label_full_name()}
						<span class="text-red-500 dark:text-red-400">*</span>
					</Field.Label>
					<div class="relative">
						<Icon icon="mdi:account" class="absolute top-1/2 left-3 -translate-y-1/2" />
						<Input
							bind:value={$form.full_name}
							name="name"
							type="text"
							class="ps-10"
							placeholder={i18n.label_full_name()}
							aria-invalid={!!fieldErrors.full_name}
							autocomplete="given-name"
						/>
					</div>
					{#if fieldErrors.full_name}
						<Field.Error>{fieldErrors.full_name}</Field.Error>
					{/if}
				</Field.Field>
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
							aria-invalid={!!fieldErrors.email}
							autocomplete="email"
						/>
					</div>
					{#if fieldErrors.email}
						<Field.Error>{fieldErrors.email}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="phone" class="capitalize">
						{i18n.phone()}
						<span class="text-red-500 dark:text-red-400">*</span>
					</Field.Label>
					<PhoneInput
						bind:value={phoneInput}
						name="phone"
						country={(data.lang.toUpperCase() as CountryCode) || 'US'}
						placeholder={i18n.phone()}
						disabled={isProcessing}
					/>
					{#if fieldErrors.phone}
						<Field.Error>{fieldErrors.phone}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Field.Label for="password" class="capitalize">
						{i18n.password()}
						<span class="text-red-500 dark:text-red-400">*</span>
					</Field.Label>
					<div class="relative">
						<Icon icon="material-symbols:key" class="absolute top-4.5 left-3 -translate-y-1/2" />
						<Password.Root minScore={2}>
							<Password.Input
								bind:value={passwordInput}
								name="password"
								class="ps-10 pe-10"
								disabled={isProcessing}
								placeholder={i18n.password()}
								autocomplete="new-password"
								oninput={(e) => {
									$form.password = (e.target as HTMLInputElement).value;
								}}
							>
								<Password.ToggleVisibility />
							</Password.Input>
							<div class="flex flex-col gap-1">
								<Password.Strength bind:strength />
							</div>
						</Password.Root>
					</div>

					{#if fieldErrors.password}
						<Field.Error>{fieldErrors.password}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<Button type="submit" disabled={isProcessing}>
						{#if isProcessing}
							<Spinner />
						{/if}
						{isProcessing ? i18n.please_wait() : i18n.create_account()}
					</Button>
				</Field.Field>
				<Field.Separator>OR</Field.Separator>
				<Field.Field class="grid gap-4 sm:grid-cols-2">
					<Button
						type="button"
						variant="outline"
						disabled={isProcessing}
						onclick={() => handleSocialLogin('google')}
					>
						<Icon icon="devicon:google" class="text-xl" />
						{i18n.sign_up_with_google()}
					</Button>
					<Button
						type="button"
						variant="outline"
						disabled={isProcessing}
						onclick={() => handleSocialLogin('facebook')}
					>
						<Icon icon="devicon:facebook" class="text-xl" />
						{i18n.sign_up_with_facebook()}
					</Button>
				</Field.Field>
			</Field.Group>
		</form>
		<div class="flex w-full flex-col items-center gap-2 pt-2">
			<div class="text-sm text-muted-foreground">
				{i18n.already_have_account()}
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
				{i18n.page_sign_up_terms()}
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
