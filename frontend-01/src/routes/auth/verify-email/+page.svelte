<script lang="ts">
	import { onMount } from 'svelte';
	import { MetaTags } from 'svelte-meta-tags';
	// import { ModalResendVerification } from '@/components/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { Spinner } from '@/components/ui/spinner/index.js';
	import Icon from '@iconify/svelte';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '$lib/paraglide/runtime';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let errorMessage = $state<string | undefined>(undefined);
	let successMessage = $state<string | undefined>(undefined);
	let progressValue = $state<number>(0);

	const validateToken = async () => {
		progressValue = 20;
		try {
			const response = await fetch('/api/auth/verify-email', {
				method: 'POST',
				credentials: 'include',
				body: JSON.stringify({ email: data.email })
			});
			const json = await response.json();
			if (!response.ok) {
				errorMessage = json.error.message || i18n.invalid_input();
			} else {
				successMessage = json.message || i18n.page_verify_email_success_message();
			}
		} catch (error) {
			errorMessage =
				error instanceof Error ? error.message : i18n.page_verify_email_error_message();
		} finally {
			progressValue = 100;
		}
	};

	onMount(async () => {
		await validateToken();
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
					<Button type="button" variant="default" href={localizeHref('/auth/login')}>
						{i18n.button_go_back_to_sign_in()}
					</Button>
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
