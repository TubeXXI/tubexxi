<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AppAlertDialog } from '@/components/index.js';
	import { AdminSidebarLayout, AdminSettingLayout } from '@/components/admin';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Textarea } from '@/components/ui/textarea/index.js';
	import { Input } from '@/components/ui/input/index.js';
	import { Switch } from '@/components/ui/switch/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import { Badge } from '@/components/ui/badge/index.js';
	import Icon from '@iconify/svelte';

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

	const sourceLogoFaviconOptions = [
		{
			label: 'Local',
			value: 'local',
			description: 'Use the local logo and favicon files.'
		},
		{
			label: 'Remote',
			value: 'remote',
			description: 'Use the remote logo and favicon files.'
		}
	];
	const themeOptions = [
		{
			label: 'Default',
			value: 'default',
			description: 'Default theme.',
			active: true
		},
		{
			label: 'Red Dark',
			value: 'red-dark',
			description: 'Red dark theme.',
			active: false
		},
		{
			label: 'Green Dark',
			value: 'green-dark',
			description: 'Green dark theme.',
			active: false
		},
		{
			label: 'Blue Dark',
			value: 'blue-dark',
			description: 'Blue dark theme.',
			active: false
		}
	];
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="System Setting" user={data.user} setting={data.settings}>
	<AdminSettingLayout fixed title="Update System Setting" description="Update system setting">
		{#snippet children()}
			<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="maintenance_mode">
									Maintenance Mode : <Badge
										variant={$form.maintenance_mode ? 'default' : 'destructive'}
									>
										{$form.maintenance_mode ? 'Active' : 'Inactive'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.maintenance_mode
										? 'Maintenance mode is enabled. Only administrators can access the site.'
										: 'Maintenance mode is disabled. All users can access the site.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="maintenance_mode" value={$form.maintenance_mode} />
							<Switch
								id="maintenance_mode"
								bind:checked={$form.maintenance_mode}
								name="maintenance_mode"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.maintenance_mode = val)}
							/>
						</Field.Field>
						{#if $form.maintenance_mode}
							<Field.Field>
								<Field.Label for="maintenance_message">Maintenance Message</Field.Label>
								<Textarea
									bind:value={$form.maintenance_message}
									name="maintenance_message"
									placeholder="Enter maintenance message"
									aria-invalid={!!$errors.maintenance_message}
									disabled={$submitting}
								/>
								{#if $errors.maintenance_message}
									<Field.Error>{$errors.maintenance_message}</Field.Error>
								{/if}
							</Field.Field>
						{/if}
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="enable_documentation">
									Enable Documentation : <Badge
										variant={$form.enable_documentation ? 'default' : 'destructive'}
									>
										{$form.enable_documentation ? 'Active' : 'Inactive'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.enable_documentation
										? 'Documentation is enabled. Users can access the documentation.'
										: 'Documentation is disabled. Users cannot access the documentation.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="enable_documentation" value={$form.enable_documentation} />
							<Switch
								id="enable_documentation"
								bind:checked={$form.enable_documentation}
								name="enable_documentation"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.enable_documentation = val)}
							/>
						</Field.Field>
						<Field.Field>
							<Field.Label for="source_logo_favicon">Source Logo Favicon</Field.Label>
							<Field.Description>
								Choose where the site logo and favicon should be sourced from.
							</Field.Description>
							<RadioGroup.Root
								bind:value={$form.source_logo_favicon}
								name="source_logo_favicon"
								onValueChange={(val) => ($form.source_logo_favicon = val as 'local' | 'remote')}
							>
								{#each sourceLogoFaviconOptions as option}
									<Field.Label for={option.value}>
										<Field.Field orientation="horizontal" class="cursor-pointer">
											<Field.Content>
												<Field.Title>{option.label}</Field.Title>
												<Field.Description>
													{option.description}
												</Field.Description>
											</Field.Content>
											<RadioGroup.Item value={option.value} id={option.value} />
										</Field.Field>
									</Field.Label>
								{/each}
							</RadioGroup.Root>
						</Field.Field>
						<Field.Field>
							<Field.Label for="theme">Theme</Field.Label>
							<div class="relative">
								<Icon icon="mdi:palette-outline" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Select.Root
									bind:value={$form.theme}
									type="single"
									name="theme"
									disabled={$submitting}
								>
									<Select.Trigger class="w-full ps-10 uppercase">
										{$form.theme.replaceAll('-', ' ')}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											<Select.Label>Themes</Select.Label>
											{#each themeOptions as option}
												<Select.Item
													value={option.value}
													label={option.label}
													disabled={!option.active}
												>
													{option.label}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							</div>
							{#if $errors.app_store_app_url}
								<Field.Error>{$errors.app_store_app_url}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="google_analytics_code">Google Analytics Code</Field.Label>
							<div class="relative">
								<Icon icon="ic:round-analytics" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.google_analytics_code}
									name="google_analytics_code"
									type="text"
									class="ps-10"
									placeholder="Enter google analytics code"
									autocomplete="on"
									aria-invalid={!!$errors.google_analytics_code}
									disabled={$submitting}
								/>
							</div>
							{#if $errors.google_analytics_code}
								<Field.Error>{$errors.google_analytics_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="histats_tracking_code">Histats Tracking Code</Field.Label>
							<div class="relative">
								<Icon
									icon="ic:round-track-changes"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Input
									bind:value={$form.histats_tracking_code}
									name="histats_tracking_code"
									type="text"
									class="ps-10"
									placeholder="Enter histats tracking code"
									autocomplete="on"
									aria-invalid={!!$errors.histats_tracking_code}
									disabled={$submitting}
								/>
							</div>
							{#if $errors.histats_tracking_code}
								<Field.Error>{$errors.histats_tracking_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="play_store_app_url">Play Store App URL</Field.Label>
							<div class="relative">
								<Icon icon="mdi:google-play" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.play_store_app_url}
									name="play_store_app_url"
									type="url"
									class="ps-10"
									placeholder="Enter play store app url"
									autocomplete="url"
									aria-invalid={!!$errors.play_store_app_url}
									disabled={$submitting}
								/>
							</div>
							{#if $errors.play_store_app_url}
								<Field.Error>{$errors.play_store_app_url}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="app_store_app_url">App Store App URL</Field.Label>
							<div class="relative">
								<Icon
									icon="mingcute:appstore-fill"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Input
									bind:value={$form.app_store_app_url}
									name="app_store_app_url"
									type="url"
									class="ps-10"
									placeholder="Enter app store app url"
									autocomplete="url"
									aria-invalid={!!$errors.app_store_app_url}
									disabled={$submitting}
								/>
							</div>
							{#if $errors.app_store_app_url}
								<Field.Error>{$errors.app_store_app_url}</Field.Error>
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
