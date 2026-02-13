<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Checkbox } from '@/components/ui/checkbox';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime';
	import { firebaseClient } from '@/client/firebase_client';
	import { FirebaseMFAHelper } from '@/client/firebase_mfa';
	import { getMultiFactorResolver, type MultiFactorResolver } from 'firebase/auth';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let passwordType = $state('password');
	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let isProcessing = $state(false);
	let mfaResolver = $state<MultiFactorResolver | null>(null);
	let mfaCode = $state('');
	let mfaPhone = $state<string | undefined>(undefined);
	let mfaStep = $state(false);

	const mfa = new FirebaseMFAHelper();

	async function createSession(idToken: string): Promise<{ redirect_url: string } | null> {
		const res = await fetch('/api/auth/session', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
			body: JSON.stringify({ idToken, rememberMe: $form.remember_me })
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

	async function finishLogin(idToken: string) {
		const data = await createSession(idToken);
		await invalidateAll();
		isProcessing = false;
		await goto(localizeHref(data?.redirect_url || '/'));
	}

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.loginForm, {
		dataType: 'json',
		async onSubmit(input) {
			errorMessage = undefined;
			successMessage = undefined;
			isProcessing = true;
			input.cancel();
			try {
				if (mfaStep) {
					if (!mfaResolver) {
						throw new Error(i18n.page_sign_in_error_message());
					}
					if (!mfaCode.trim()) {
						throw new Error(i18n.invalid_input());
					}
					const credential = await mfa.completeSignIn(mfaResolver, mfaCode.trim());
					const idToken = await credential.user.getIdToken();
					await finishLogin(idToken);
					return;
				}

				const firebaseUser = await firebaseClient?.signInWithEmail($form.email, $form.password);
				if (!firebaseUser?.user) {
					throw new Error(i18n.page_sign_in_error_message());
				}

				const idToken = await firebaseUser.user.getIdToken();
				if (!idToken) {
					throw new Error(i18n.page_sign_in_error_message());
				}

				await finishLogin(idToken);
			} catch (error) {
				const errAny = error as any;
				if (errAny?.code === 'auth/multi-factor-auth-required') {
					try {
						const auth = firebaseClient?.getAuth();
						if (!auth) {
							throw new Error(i18n.page_sign_in_error_message());
						}
						mfaResolver = getMultiFactorResolver(auth, errAny);
						mfaPhone = (mfaResolver?.hints?.[0] as any)?.phoneNumber;
						await mfa.sendSignInVerificationCode(mfaResolver);
						mfaStep = true;
						successMessage = i18n.success();
						errorMessage = undefined;
					} catch (mfaErr) {
						errorMessage =
							mfaErr instanceof Error ? mfaErr.message : i18n.page_sign_in_error_message();
						mfaResolver = null;
						mfaStep = false;
					}
					isProcessing = false;
					return;
				}
				errorMessage = error instanceof Error ? error.message : i18n.page_sign_in_error_message();
				isProcessing = false;
			}
		},
		onError() {
			isProcessing = false;
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
					idToken,
					rememberMe: $form.remember_me
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
							idToken,
							rememberMe: $form.remember_me
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
				{i18n.page_sign_in()}
			</h2>
			<p class="text-sm text-muted-foreground">
				{i18n.page_sign_in_description()}
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
		{/if}
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
							type="text"
							class="ps-10"
							placeholder={i18n.email()}
							aria-invalid={!!$errors.email}
							autocomplete="email"
							disabled={isProcessing || $submitting || mfaStep}
						/>
					</div>
					{#if $errors.email}
						<Field.Error>{$errors.email}</Field.Error>
					{/if}
				</Field.Field>
				<Field.Field>
					<div class="flex items-center">
						<Field.Label for="password" class="capitalize">
							{i18n.password()}
							<span class="text-red-500 dark:text-red-400">*</span>
						</Field.Label>
						<a
							href={localizeHref('/auth/reset-password')}
							class="ml-auto text-xs underline-offset-4 hover:underline"
						>
							{i18n.forgot_password()}?
						</a>
					</div>
					<div class="relative">
						<Icon icon="material-symbols:key" class="absolute top-1/2 left-3 -translate-y-1/2" />
						<Input
							bind:value={$form.password}
							name="password"
							type={passwordType}
							class="ps-10"
							placeholder={i18n.password()}
							aria-invalid={!!$errors.password}
							autocomplete="new-password"
							disabled={isProcessing || $submitting || mfaStep}
						/>

						<Button
							variant="ghost"
							size="icon"
							class="absolute top-1/2 right-1 size-8 -translate-y-1/2 cursor-pointer"
							onclick={() => (passwordType = passwordType === 'password' ? 'text' : 'password')}
							disabled={isProcessing || $submitting || mfaStep}
						>
							<Icon
								icon={passwordType === 'password' ? 'mdi:eye' : 'mdi:eye-off'}
								class="absolute top-1/2 right-3 -translate-y-1/2 cursor-pointer"
							/>
						</Button>
					</div>
					{#if $errors.password}
						<Field.Error>{$errors.password}</Field.Error>
					{/if}
				</Field.Field>
				{#if mfaStep}
					<Field.Field>
						<Field.Label for="mfa_code" class="capitalize">
							Kode OTP
							{#if mfaPhone}
								<span class="ml-2 text-xs text-muted-foreground">({mfaPhone})</span>
							{/if}
						</Field.Label>
						<Input
							bind:value={mfaCode}
							name="mfa_code"
							type="text"
							inputmode="numeric"
							autocomplete="one-time-code"
							placeholder="123456"
							disabled={isProcessing || $submitting}
						/>
					</Field.Field>
					<Field.Field>
						<Button
							type="button"
							variant="ghost"
							disabled={isProcessing || $submitting}
							onclick={() => {
								mfaStep = false;
								mfaResolver = null;
								mfaCode = '';
								mfaPhone = undefined;
							}}
						>
							Batal verifikasi
						</Button>
					</Field.Field>
				{/if}
				<Field.Field orientation="horizontal">
					<Checkbox
						id="remember_me"
						name="remember_me"
						bind:checked={$form.remember_me}
						onCheckedChange={(value) => ($form.remember_me = value)}
						disabled={isProcessing || $submitting || mfaStep}
					/>
					<Field.Label for="remember_me" class="font-normal capitalize">
						{i18n.remember_me()}
					</Field.Label>
				</Field.Field>
				<Field.Field>
					<Button type="submit" disabled={$submitting || isProcessing}>
						{#if $submitting || isProcessing}
							<Spinner />
						{/if}
						{$submitting || isProcessing
							? i18n.please_wait()
							: mfaStep
								? 'Verifikasi'
								: i18n.sign_in()}
					</Button>
				</Field.Field>
				<Field.Separator>OR</Field.Separator>
				<Field.Field class="grid gap-4 sm:grid-cols-2">
					<Button
						type="button"
						variant="outline"
						disabled={$submitting || isProcessing || mfaStep}
						onclick={() => handleSocialLogin('google')}
					>
						<Icon icon="devicon:google" class="text-xl" />
						{i18n.sign_in_with_google()}
					</Button>
					<Button
						type="button"
						variant="outline"
						disabled={$submitting || isProcessing || mfaStep}
						onclick={() => handleSocialLogin('facebook')}
					>
						<Icon icon="devicon:facebook" class="text-xl" />
						{i18n.sign_in_with_facebook()}
					</Button>
				</Field.Field>
			</Field.Group>
			<div id="recaptcha-container" class="hidden"></div>
		</form>
		<div class="flex w-full flex-col items-center gap-2 pt-2">
			<div class="text-sm text-muted-foreground">
				{i18n.dont_have_account()}
			</div>
			<Button
				href={localizeHref('/auth/register')}
				type="button"
				variant="outline"
				class="w-full text-sm"
				disabled={$submitting || isProcessing}
			>
				{i18n.create_account()}
			</Button>
		</div>
	</div>

	<div class="rounded-lg px-5 py-8 text-neutral-300">
		<div class="flex w-full flex-col items-start">
			<p class="text-sm">
				{i18n.page_sign_in_terms()}
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
