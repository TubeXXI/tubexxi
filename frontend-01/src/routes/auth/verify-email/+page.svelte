<script lang="ts">
	import { onMount } from 'svelte';
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Spinner } from '@/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime';
	import { firebaseClient } from '@/client/firebase_client.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let progressValue = $state<number>(0);

	const sendVerificationEmail = async () => {
		progressValue = 20;
		try {
			const response = await fetch('/api/auth/verify-email', {
				method: 'POST',
				credentials: 'include',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email: data.email })
			});
			const json = (await response.json().catch(() => null)) as ApiResponse<unknown> | null;
			if (!response.ok) {
				errorMessage = json?.message || i18n.invalid_input();
			} else {
				successMessage = json?.message || i18n.page_verify_email_success_message();
			}
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : i18n.page_verify_email_error_message();
		} finally {
			progressValue = 100;
		}
	};

	const confirmEmailWithToken = async (token: string) => {
		progressValue = 20;
		try {
			if (!firebaseClient) {
				throw new Error(i18n.error_firebase_auth_not_initialized());
			}
			await firebaseClient.applyEmailVerificationCode(token);
			const currentUser = firebaseClient.getCurrentUser();
			if (currentUser) {
				const idToken = await currentUser.getIdToken(true);
				await fetch('/api/auth/refresh', {
					method: 'PUT',
					credentials: 'include',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ idToken })
				});
			}
			await invalidateAll();
			successMessage = i18n.page_verify_email_success_message();
			errorMessage = undefined;
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : i18n.page_verify_email_error_message();
		} finally {
			progressValue = 100;
		}
	};

	onMount(async () => {
		const token = ((data as any).token as string) || '';
		if (token) {
			await confirmEmailWithToken(token);
			return;
		}
		if (data.email) {
			await sendVerificationEmail();
		}
	});
</script>

<MetaTags {...metaTags} />
<data.component setting={data.settings}>
	<div class="flex w-full flex-col items-start gap-y-6 rounded-lg bg-accent px-5 py-8">
		<Card.Root class="w-full">
			<Card.Header>
				<Card.Title>
					{i18n.page_verify_email()}
				</Card.Title>
				<Card.Description>{i18n.page_verify_email_description()}</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-4">
				{#if errorMessage}
					<Alert.Root variant="destructive">
						<Icon icon="mingcute:warning-line" class="size-4" />
						<Alert.Title>{i18n.error()}</Alert.Title>
						<Alert.Description>{errorMessage}</Alert.Description>
					</Alert.Root>
					<!-- <ModalResendVerification
						onclose={(msg) => {
							if (msg) {
								successMessage = msg;
								errorMessage = undefined;
							}
						}}
					/> -->
				{:else if successMessage}
					<Alert.Root variant="default">
						<Icon icon="mingcute:check-line" class="size-4" />
						<Alert.Title>{i18n.success()}</Alert.Title>
						<Alert.Description>{successMessage}</Alert.Description>
					</Alert.Root>
					<div class="flex flex-col gap-2">
						<Button type="button" variant="default" href={localizeHref('/auth/login')}>
							{i18n.button_go_back_to_sign_in()}
						</Button>
						{#if data.email && !(data as any).token}
							<Button
								type="button"
								variant="secondary"
								onclick={async () => {
									errorMessage = undefined;
									successMessage = undefined;
									await sendVerificationEmail();
								}}
							>
								{i18n.resend_verification_email()}
							</Button>
						{/if}
					</div>
				{:else}
					<Progress value={progressValue} class="w-full" />
					<div class="flex items-center justify-center gap-2">
						<Spinner class="size-4" />
						<p class="text-sm opacity-70">
							{i18n.please_wait()}
						</p>
					</div>
				{/if}
			</Card.Content>
		</Card.Root>
	</div>
</data.component>
