<script lang="ts">
	import { page } from '$app/state';
	import { invalidateAll, goto } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import * as Field from '$lib/components/ui/field/index.js';
	import type { ZxcvbnResult } from '@zxcvbn-ts/core';
	import * as Password from '$lib/components/ui-extras/password';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { UserHeader, UserContent } from '@/components/user/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import { toast } from '@/stores/toast';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	let currentPasswordInput = $state<string | undefined>('');
	let newPasswordInput = $state<string | undefined>('');
	let confirmPasswordInput = $state<string | undefined>('');
	let strength = $state<ZxcvbnResult>();

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.form, {
		resetForm: false,
		onSubmit: (input) => {
			invalidateAll();
		},
		async onUpdate({ result }) {
			if (result.type === 'failure') {
				handleSubmitLoading(false);
				toast.error(result.data.message);
				return;
			}
			if (result.type === 'success') {
				handleSubmitLoading(false);
				toast.success(result.data.message);

				await invalidateAll();
				setTimeout(async () => {
					await goto(`/login?redirect=${page.url.pathname}`);
				}, 2000);
			}
		},
		onError: ({ result }) => {
			handleSubmitLoading(false);
			toast.error(result.error?.message || 'Validation error');
		}
	});
</script>

<MetaTags {...metaTags} />

<div class="mx-auto max-w-5xl space-y-6 px-4 py-10">
	<UserHeader user={data.user} />
	<UserContent title={i18n.change_password()} description={i18n.change_password_description()}>
		<form method="POST" class="space-y-4" use:enhance>
			<Field.Group>
				<Field.Set>
					<Field.Field>
						<Field.Label for="current_password">
							{i18n.current_password()}
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
									placeholder={i18n.current_password()}
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
							{i18n.new_password()}
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
									placeholder={i18n.new_password()}
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
							{i18n.confirm_password()}
							<span class="text-red-500 dark:text-red-400">*</span>
						</Field.Label>
						<div class="relative">
							<Icon icon="material-symbols:key" class="absolute top-1/2 left-3 -translate-y-1/2" />
							<Password.Root minScore={2}>
								<Password.Input
									id="confirm_password"
									bind:value={confirmPasswordInput}
									name="confirm_password"
									class="ps-10 pe-10"
									disabled={$submitting}
									placeholder={i18n.confirm_password()}
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
						{$submitting ? i18n.please_wait() : i18n.change_password()}
					</Button>
				</Field.Field>
			</Field.Group>
		</form>
	</UserContent>
</div>
