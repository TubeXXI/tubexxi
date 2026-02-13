<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AppAlertDialog } from '@/components/index.js';
	import { AdminSidebarLayout, AdminSettingLayout } from '@/components/admin';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Switch } from '@/components/ui/switch/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { Input } from '@/components/ui/input/index.js';
	import Icon from '@iconify/svelte';
	import { Badge } from '@/components/ui/badge/index.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);
	// svelte-ignore state_referenced_locally
	const { form, errors, submitting, enhance } = superForm(data.form, {
		resetForm: false,
		onSubmit(input) {
			handleSubmitLoading(true);
			successMessage = null;
			errorMessage = null;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			if (event.result.type === 'success') {
				handleSubmitLoading(false);
				successMessage = event.result.data.message;
				await invalidateAll();
			}
		},
		onError(event) {
			errorMessage = event.result.error.message;
			handleSubmitLoading(false);
		}
	});
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Email Setting" user={data.user} setting={data.settings}>
	<AdminSettingLayout fixed title="Email Setting" description="Update email setting">
		{#snippet children()}
			<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="smtp_enabled">
									SMTP Enabled : <Badge variant={$form.smtp_enabled ? 'default' : 'destructive'}>
										{$form.smtp_enabled ? 'Active' : 'Inactive'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.smtp_enabled
										? 'SMTP is enabled. Emails will be sent using the configured SMTP server.'
										: 'SMTP is disabled. Emails will be sent using the default mail server.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="smtp_enabled" value={$form.smtp_enabled} />
							<Switch
								id="smtp_enabled"
								bind:checked={$form.smtp_enabled}
								name="smtp_enabled"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.smtp_enabled = val)}
							/>
						</Field.Field>
						<Field.Field>
							<Field.Label for="smtp_service">
								SMTP Service <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.smtp_service}
									name="smtp_service"
									type="text"
									class="ps-10"
									placeholder="Enter SMTP service (e.g., smtp.gmail.com)"
									aria-invalid={!!$errors.smtp_service}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.smtp_service}
								<Field.Error>{$errors.smtp_service}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="smtp_host">
								SMTP Host <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.smtp_host}
									name="smtp_host"
									type="text"
									class="ps-10"
									placeholder="Enter SMTP host (e.g., smtp.gmail.com)"
									aria-invalid={!!$errors.smtp_host}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.smtp_host}
								<Field.Error>{$errors.smtp_host}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="smtp_port">
								SMTP Port <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.smtp_port}
									name="smtp_port"
									type="number"
									class="ps-10"
									placeholder="Enter SMTP port (e.g., 587)"
									aria-invalid={!!$errors.smtp_port}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.smtp_port}
								<Field.Error>{$errors.smtp_port}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="smtp_user">
								SMTP Username <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.smtp_user}
									name="smtp_user"
									type="text"
									class="ps-10"
									placeholder="Enter SMTP username (e.g., your-email@gmail.com)"
									aria-invalid={!!$errors.smtp_user}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.smtp_user}
								<Field.Error>{$errors.smtp_user}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="smtp_password">
								SMTP Password <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.smtp_password}
									name="smtp_password"
									type="password"
									class="ps-10"
									placeholder="Enter SMTP password"
									aria-invalid={!!$errors.smtp_password}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.smtp_password}
								<Field.Error>{$errors.smtp_password}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="from_email">
								From Email <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.from_email}
									name="from_email"
									type="text"
									class="ps-10"
									placeholder="Enter from email (e.g., your-email@gmail.com)"
									aria-invalid={!!$errors.from_email}
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.from_email}
								<Field.Error>{$errors.from_email}</Field.Error>
							{/if}
						</Field.Field>
					</Field.Set>
					<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
						<Button type="submit" class="w-full" disabled={$submitting}>
							{#if $submitting}
								<Spinner class="mr-2 size-5" />
							{/if}
							{$submitting ? 'Submitting' : 'Update'}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		{/snippet}
	</AdminSettingLayout>
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
