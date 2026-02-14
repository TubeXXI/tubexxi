<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import type { E164Number } from 'svelte-tel-input/types';
	import { PhoneInput } from '@/components/ui-extras/phone-input/index.js';
	import { UserHeader, UserContent } from '@/components/user/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import { toast } from '@/stores/toast';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	// svelte-ignore state_referenced_locally
	let phoneInput = $state<E164Number | null>(data.user?.phone || null);
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
				setTimeout(async () => {}, 2000);
			}
		},
		onError: ({ result }) => {
			handleSubmitLoading(false);
			toast.error(result.error?.message || 'Validation error');
		}
	});

	$effect(() => {
		if (data.user?.phone && data.user.phone.trim() !== '') {
			phoneInput = data.user.phone as E164Number;
		}
	});
	$effect(() => {
		if (phoneInput && typeof phoneInput === 'string' && phoneInput.trim() !== '') {
			const cleanNumber = phoneInput.replace(/[\s-]/g, '');
			$form.phone = cleanNumber;
		} else {
			$form.phone = '';
		}
	});
</script>

<MetaTags {...metaTags} />

<div class="mx-auto max-w-5xl space-y-6 px-4 py-10">
	<UserHeader user={data.user} />
	<UserContent
		title={i18n.account_profile_setting()}
		description={i18n.account_profile_setting_description()}
	>
		<form method="POST" class="space-y-4" use:enhance>
			<Field.Group>
				<Field.Set>
					<Field.Group>
						<Field.Set>
							<Field.Legend>{i18n.personal_information()}</Field.Legend>
							<Field.Description>{i18n.personal_information_description()}</Field.Description>
							<Field.Group>
								<Field.Field>
									<Field.Label for="full_name">
										{i18n.name()}
										<span class="text-red-500 dark:text-red-400">*</span>
									</Field.Label>
									<div class="relative">
										<Icon icon="mdi:account" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.full_name}
											name="full_name"
											type="text"
											class="ps-10"
											placeholder={i18n.name()}
											aria-invalid={!!$errors.full_name}
											autocomplete="name"
											disabled={$submitting}
										/>
									</div>
									{#if $errors.full_name}
										<Field.Error>{$errors.full_name}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="email">
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
											disabled={$submitting}
										/>
									</div>
									{#if $errors.email}
										<Field.Error>{$errors.email}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="phone">
										{i18n.phone()}
										<span class="text-red-500 dark:text-red-400">*</span>
									</Field.Label>
									<div class="relative">
										<Icon icon="mdi:phone" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<PhoneInput
											bind:value={phoneInput}
											name="phone"
											placeholder={i18n.phone()}
											disabled={$submitting}
											options={{
												autoPlaceholder: true,
												spaces: false,
												format: 'international',
												invalidateOnCountryChange: true
											}}
										/>
									</div>
									{#if $errors.phone}
										<Field.Error>{$errors.phone}</Field.Error>
									{/if}
								</Field.Field>
							</Field.Group>
						</Field.Set>
						<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
							<Button type="submit" class="w-full" disabled={$submitting}>
								{#if $submitting}
									<Spinner class="mr-2 size-5" />
								{/if}
								{$submitting ? i18n.please_wait() : i18n.update_profile()}
							</Button>
						</Field.Field>
					</Field.Group>
				</Field.Set>
			</Field.Group>
		</form>
	</UserContent>
</div>
