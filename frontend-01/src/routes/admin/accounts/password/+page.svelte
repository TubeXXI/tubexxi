<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import type { ZxcvbnResult } from '@zxcvbn-ts/core';
	import { handleSubmitLoading } from '@/stores';
	import { AdminSidebarLayout, AdminAccountLayout } from '@/components/admin';
	import { AppAlertDialog } from '@/components/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import Icon from '@iconify/svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Password from '$lib/components/ui-extras/password';
	import { localizeHref } from '@/paraglide/runtime.js';
	import { firebaseClient } from '@/client/firebase_client';

	let { data } = $props();

	let metaTags = $derived(data.pageMetaTags);
	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);
	let currentPasswordInput = $state<string | undefined>('');
	let newPasswordInput = $state<string | undefined>('');
	let confirmPasswordInput = $state<string | undefined>('');
	let strength = $state<ZxcvbnResult>();

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.passwordForm, {
		resetForm: true,
		onSubmit: (input) => {
			successMessage = null;
			errorMessage = null;
			invalidateAll();
		},
		async onUpdate({ result }) {
			if (result.type === 'failure') {
				handleSubmitLoading(false);
				errorMessage = result.data.message;
				return;
			}
			if (result.type === 'success') {
				handleSubmitLoading(false);
				successMessage = result.data.message;

				await firebaseClient?.updatePassword($form.new_password);
				setTimeout(async () => {
					successMessage = null;
					errorMessage = null;
					await firebaseClient?.signOut();
					await invalidateAll();
					await goto(localizeHref('/auth/login?redirect=${page.url.pathname}'));
				}, 1000);
			}
		},
		onError: ({ result }) => {
			handleSubmitLoading(false);
			errorMessage = result.error?.message || 'Validation error';
		}
	});
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Password" user={data.user} setting={data.settings}>
	<AdminAccountLayout fixed title="Change Password" description="Change your password">
		{#snippet children()}
			<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field>
							<Field.Label for="current_password">
								Current Password
								<span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon
									icon="material-symbols:key"
									class="absolute top-1/2 left-3 -translate-y-1/2 text-neutral-900 dark:text-neutral-50"
								/>
								<Password.Root minScore={2}>
									<Password.Input
										id="current_password"
										bind:value={currentPasswordInput}
										name="current_password"
										class="ps-10 pe-10"
										disabled={$submitting}
										placeholder="Current Password"
										autocomplete="current-password"
										oninput={(e) => {
											$form.current_password = (e.target as HTMLInputElement).value;
										}}
									>
										<Password.ToggleVisibility />
									</Password.Input>
								</Password.Root>
							</div>
							{#if $errors.current_password}
								<Field.Error>{$errors.current_password}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="new_password">
								New Password
								<span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="material-symbols:key" class="absolute top-2.5 left-3 " />
								<Password.Root minScore={2}>
									<Password.Input
										id="new_password"
										bind:value={newPasswordInput}
										name="new_password"
										class="ps-10 pe-10"
										disabled={$submitting}
										placeholder="New Password"
										autocomplete="new-password"
										oninput={(e) => {
											$form.new_password = (e.target as HTMLInputElement).value;
										}}
									>
										<Password.ToggleVisibility />
									</Password.Input>
									<div class="flex flex-col gap-1">
										<Password.Strength bind:strength />
									</div>
								</Password.Root>
							</div>
							{#if $errors.new_password}
								<Field.Error>{$errors.new_password}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="confirm_password">
								Confirm Password
								<span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon
									icon="material-symbols:key"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Password.Root minScore={2}>
									<Password.Input
										id="confirm_password"
										bind:value={confirmPasswordInput}
										name="confirm_password"
										class="ps-10 pe-10"
										disabled={$submitting}
										placeholder="Confirm Password"
										autocomplete="new-password"
										oninput={(e) => {
											$form.confirm_password = (e.target as HTMLInputElement).value;
										}}
									>
										<Password.ToggleVisibility />
									</Password.Input>
								</Password.Root>
							</div>
							{#if $errors.confirm_password}
								<Field.Error>{$errors.confirm_password}</Field.Error>
							{/if}
						</Field.Field>
					</Field.Set>
					<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
						<Button type="submit" class="w-full" disabled={$submitting}>
							{#if $submitting}
								<Spinner class="mr-2 size-5" />
							{/if}
							{$submitting ? 'Submitting' : 'Change Password'}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		{/snippet}
	</AdminAccountLayout>
	{#if successMessage}
		<AppAlertDialog
			open={true}
			title="Success"
			message={successMessage}
			type="success"
			onclose={() => (successMessage = null)}
		/>
	{/if}
	{#if errorMessage}
		<AppAlertDialog
			open={true}
			title="Error"
			message={errorMessage}
			type="error"
			onclose={() => (errorMessage = null)}
		/>
	{/if}
</AdminSidebarLayout>
