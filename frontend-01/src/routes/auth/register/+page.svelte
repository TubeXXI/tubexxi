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
	import { FirebaseMFAHelper } from '@/client/firebase_mfa';
	import { getMultiFactorResolver, type MultiFactorResolver } from 'firebase/auth';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	let settingState = $derived(data.settings);

	let showConfirmPassword = $state(false);
	let phoneInput = $state<E164Number | undefined>('');
	let passwordInput = $state<string | undefined>('');
	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let isProcessing = $state(false);
	const SCORE_NAMING = ['Poor', 'Weak', 'Average', 'Strong', 'Secure'];
	let strength = $state<ZxcvbnResult>();

	const mfa = new FirebaseMFAHelper();

	async function createSession(idToken: string): Promise<{ redirect_url: string } | null> {
		const res = await fetch('/api/auth/session', {
			method: 'POST',
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
	async function finishRegister(idToken: string) {
		const data = await createSession(idToken);
		await invalidateAll();
		isProcessing = false;
	}

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.registerForm, {
		resetForm: true,
		dataType: 'json',
		async onSubmit(input) {
			errorMessage = undefined;
			successMessage = undefined;
			isProcessing = true;
			input.cancel();
			try {
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

				input.formData.append('idToken', idToken);
				input.formData.append('firebaseUid', result.user.uid);
				input.formData.append('full_name', $form.full_name);
				input.formData.append('email', $form.email);
				input.formData.append('password', $form.password);
				if ($form.phone) {
					input.formData.append('phone', $form.phone);
				}
				if (result.user.photoURL) {
					input.formData.append('avatar_url', result.user.photoURL);
				}
			} catch (error) {
				input.cancel();
				errorMessage = error instanceof Error ? error.message : i18n.invalid_input();
				isProcessing = false;
			}
		},
		async onUpdate(event) {
			isProcessing = false;

			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;

				try {
					const currentUser = firebaseClient?.getCurrentUser();
					if (currentUser) {
						await currentUser.delete();
					}
				} catch (error) {
					console.error('Error cleanup current user:', error);
				}
				return;
			}
			if (event.result.type === 'success') {
				successMessage = event.result.data.message;
				await invalidateAll();
			}
		},
		onError(event) {
			errorMessage = event.result.error.message;
		}
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
		<form method="POST" class="w-full" use:enhance>
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
							aria-invalid={!!$errors.full_name}
							autocomplete="given-name"
						/>
					</div>
					{#if $errors.full_name}
						<Field.Error>{$errors.full_name}</Field.Error>
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
							aria-invalid={!!$errors.email}
							autocomplete="email"
						/>
					</div>
					{#if $errors.email}
						<Field.Error>{$errors.email}</Field.Error>
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
						country={(data.lang as CountryCode) || 'US'}
						placeholder={i18n.phone()}
						disabled={$submitting || isProcessing}
					/>
					{#if $errors.phone}
						<Field.Error>{$errors.phone}</Field.Error>
					{/if}
				</Field.Field>
				<div class="grid grid-cols-2 gap-4">
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
									disabled={$submitting || isProcessing}
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

						{#if $errors.password}
							<Field.Error>{$errors.password}</Field.Error>
						{/if}
					</Field.Field>
				</div>
				<Field.Field>
					<Button type="submit" disabled={$submitting || isProcessing}>
						{#if $submitting || isProcessing}
							<Spinner />
						{/if}
						{$submitting || isProcessing ? i18n.please_wait() : i18n.create_account()}
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
