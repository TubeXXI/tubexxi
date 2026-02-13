<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AppAlertDialog } from '@/components/index.js';
	import { AdminSidebarLayout, AdminSettingLayout } from '@/components/admin';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Textarea } from '@/components/ui/textarea/index.js';
	import { Switch } from '@/components/ui/switch/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Badge } from '@/components/ui/badge/index.js';
	import Icon from '@iconify/svelte';
	import { isValidUrl } from '@/utils/urls.js';

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

	function onDirectLinkChange(value?: string) {
		if (!value) return;

		if (!isValidUrl(value)) {
			$errors.direct_link_ad_code = ['Invalid direct link ad code'];
			return;
		}
	}

	const typeMonetizeOptions = [
		{ value: 'adsense', label: 'AdSense' },
		{ value: 'revenuecat', label: 'RevenueCat' },
		{ value: 'adsterra', label: 'Adsterra' }
	];
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Monetization Setting" user={data.user} setting={data.settings}>
	<AdminSettingLayout fixed title="Monetization Setting" description="Update monetization setting">
		{#snippet children()}
			<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="enable_monetize">
									Enable Monetization : <Badge
										variant={$form.enable_monetize ? 'default' : 'destructive'}
									>
										{$form.enable_monetize ? 'Active' : 'Inactive'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.enable_monetize
										? 'Monetization is enabled. Ads will be shown on the video player.'
										: 'Monetization is disabled. Ads will not be shown on the video player.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="enable_monetize" value={$form.enable_monetize} />
							<Switch
								id="enable_monetize"
								bind:checked={$form.enable_monetize}
								name="enable_monetize"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.enable_monetize = val)}
							/>
						</Field.Field>
						<Field.Field>
							<Field.Label for="type_monetize">Type Monetize</Field.Label>
							<Select.Root
								type="single"
								name="type_monetize"
								bind:value={$form.type_monetize}
								disabled={$submitting || !$form.enable_monetize}
							>
								<Select.Trigger class="w-full capitalize">
									{$form.type_monetize ? $form.type_monetize : 'Select Type Monetize'}
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										<Select.Label>Type Monetize</Select.Label>
										{#each typeMonetizeOptions as option (option.value)}
											<Select.Item
												value={option.value}
												label={option.label}
												disabled={option.value === 'grapes'}
											>
												{option.label}
											</Select.Item>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
							{#if $errors.type_monetize}
								<Field.Error>{$errors.type_monetize}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="publisher_id">
								Publisher ID <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:account" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.publisher_id}
									name="publisher_id"
									type="text"
									class="ps-10"
									placeholder="Example: ca-pub-1234567890123456"
									aria-invalid={!!$errors.publisher_id}
									autocomplete="on"
									disabled={$submitting || !$form.enable_monetize}
								/>
							</div>
							{#if $errors.publisher_id}
								<Field.Error>{$errors.publisher_id}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="enable_popup_ad">
									Enable Popup Ad : <Badge
										variant={$form.enable_popup_ad ? 'default' : 'destructive'}
									>
										{$form.enable_popup_ad ? 'Active' : 'Inactive'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.enable_popup_ad
										? 'Popup ad is enabled. Ads will be shown on the video player.'
										: 'Popup ad is disabled. Ads will not be shown on the video player.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="enable_popup_ad" value={$form.enable_popup_ad} />
							<Switch
								id="enable_popup_ad"
								bind:checked={$form.enable_popup_ad}
								name="enable_popup_ad"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.enable_popup_ad = val)}
							/>
						</Field.Field>
						<Field.Field>
							<Field.Label for="popup_ad_code">Popup Ad Code</Field.Label>
							<Textarea
								bind:value={$form.popup_ad_code}
								name="popup_ad_code"
								placeholder="Enter popup ad code"
								aria-invalid={!!$errors.popup_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.popup_ad_code}
								<Field.Error>{$errors.popup_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="socialbar_ad_code">Socialbar Ad Code</Field.Label>
							<Textarea
								bind:value={$form.socialbar_ad_code}
								name="socialbar_ad_code"
								placeholder="Enter socialbar ad code"
								aria-invalid={!!$errors.socialbar_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.socialbar_ad_code}
								<Field.Error>{$errors.socialbar_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="banner_rectangle_ad_code">Banner Rectangle Ad Code</Field.Label>
							<Textarea
								bind:value={$form.banner_rectangle_ad_code}
								name="banner_rectangle_ad_code"
								placeholder="Enter banner rectangle ad code"
								aria-invalid={!!$errors.banner_rectangle_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.banner_rectangle_ad_code}
								<Field.Error>{$errors.banner_rectangle_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="banner_horizontal_ad_code">Banner Horizontal Ad Code</Field.Label>
							<Textarea
								bind:value={$form.banner_horizontal_ad_code}
								name="banner_horizontal_ad_code"
								placeholder="Enter banner horizontal ad code"
								aria-invalid={!!$errors.banner_horizontal_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.banner_horizontal_ad_code}
								<Field.Error>{$errors.banner_horizontal_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="banner_vertical_ad_code">Banner Vertical Ad Code</Field.Label>
							<Textarea
								bind:value={$form.banner_vertical_ad_code}
								name="banner_vertical_ad_code"
								placeholder="Enter banner vertical ad code"
								aria-invalid={!!$errors.banner_vertical_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.banner_vertical_ad_code}
								<Field.Error>{$errors.banner_vertical_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="native_ad_code">Native Ad Code</Field.Label>
							<Textarea
								bind:value={$form.native_ad_code}
								name="native_ad_code"
								placeholder="Enter native ad code"
								aria-invalid={!!$errors.native_ad_code}
								autocomplete="on"
								disabled={$submitting || !$form.enable_monetize}
							/>
							{#if $errors.native_ad_code}
								<Field.Error>{$errors.native_ad_code}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="direct_link_ad_code">
								Direct Link Ad Code <span class="text-red-500 dark:text-red-400">*</span>
							</Field.Label>
							<div class="relative">
								<Icon icon="mdi:link" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.direct_link_ad_code}
									name="direct_link_ad_code"
									type="text"
									class="ps-10"
									placeholder="Enter direct link ad code"
									aria-invalid={!!$errors.direct_link_ad_code}
									autocomplete="on"
									disabled={$submitting || !$form.enable_monetize}
									oninput={(e) => onDirectLinkChange((e.target as HTMLInputElement)?.value)}
								/>
							</div>
							{#if $errors.direct_link_ad_code}
								<Field.Error>{$errors.direct_link_ad_code}</Field.Error>
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
