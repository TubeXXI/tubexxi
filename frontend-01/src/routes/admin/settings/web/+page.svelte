<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AppAlertDialog } from '@/components/index.js';
	import {
		AdminSidebarLayout,
		AdminSettingLayout,
		AdminSettingUploadFavicon,
		AdminSettingUploadLogo
	} from '@/components/admin';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Textarea } from '@/components/ui/textarea/index.js';
	import { TagsInput } from '$lib/components/ui-extras/tags-input/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import Icon from '@iconify/svelte';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	// svelte-ignore state_referenced_locally
	let keywordsInput = $state<string[]>(data.settings?.WEBSITE.site_keywords?.split(',') || []);
	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);
	let tabSelected = $state('update-web-setting');

	// svelte-ignore state_referenced_locally
	const { form, errors, submitting, enhance } = superForm(data.form, {
		id: 'setting-update-web',
		dataType: 'json',
		resetForm: false,
		onSubmit(input) {
			handleSubmitLoading(true);
			successMessage = null;
			errorMessage = null;
			$form.site_keywords = keywordsInput.join(', ');
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

	const formatPhoneInput = (event: Event) => {
		let value = (event.target as HTMLInputElement).value;
		value = value.replace(/[^\d+]/g, '');

		if (value.startsWith('+')) {
			value = '+' + value.slice(1).replace(/\+/g, '');
		} else {
			value = value.replace(/\+/g, '');
		}

		(event.target as HTMLInputElement).value = value;
		$form.site_phone = value;
	};

	const tabItems = [
		{
			label: 'Update Web Setting',
			value: 'update-web-setting'
		},
		{
			label: 'Update Web Logo',
			value: 'update-web-logo'
		},
		{
			label: 'Update Web Favicon',
			value: 'update-web-favicon'
		}
	];
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Web Setting" user={data.user} setting={data.settings}>
	<AdminSettingLayout fixed title="Web Setting" description="Update web setting">
		{#snippet children()}
			<Tabs.Root
				bind:value={tabSelected}
				class="space-y-6 bg-transparent"
				onValueChange={(val) => (tabSelected = val)}
			>
				<Tabs.List class="bg-transparent">
					{#snippet child()}
						<div class="flex items-center justify-center">
							{#each tabItems as item}
								<Tabs.Trigger
									value={item.value}
									class="hidden  cursor-pointer data-[state=active]:bg-sky-600 data-[state=active]:text-white lg:block lg:w-auto dark:data-[state=active]:bg-sky-700 dark:data-[state=active]:text-white"
									onclick={() => (tabSelected = item.value)}
								>
									{item.label}
								</Tabs.Trigger>
							{/each}
						</div>
						<RadioGroup.Root
							bind:value={tabSelected}
							class="flex lg:hidden"
							orientation="horizontal"
							onValueChange={(val) => (tabSelected = val)}
						>
							{#each tabItems as item}
								<div class="flex items-center space-x-2">
									<RadioGroup.Item value={item.value} id={item.value} />
									<Label for={item.value}>{item.label}</Label>
								</div>
							{/each}
						</RadioGroup.Root>
					{/snippet}
				</Tabs.List>
				<Tabs.Content value="update-web-setting">
					<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
						<Field.Group>
							<Field.Set>
								<Field.Field>
									<Field.Label for="site_name">Site Name</Field.Label>
									<div class="relative">
										<Icon icon="lucide:globe" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.site_name}
											name="site_name"
											type="text"
											class="ps-10"
											placeholder="Enter site name"
											aria-invalid={!!$errors.site_name}
											autocomplete="name"
											disabled={$submitting}
										/>
									</div>
									{#if $errors.site_name}
										<Field.Error>{$errors.site_name}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_tagline">Site Tagline</Field.Label>
									<div class="relative">
										<Icon icon="lucide:tag" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.site_tagline}
											name="site_tagline"
											type="text"
											class="ps-10"
											placeholder="Enter site tagline"
											aria-invalid={!!$errors.site_tagline}
											autocomplete="name"
											disabled={$submitting}
										/>
									</div>
									{#if $errors.site_tagline}
										<Field.Error>{$errors.site_tagline}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_description">Site Description</Field.Label>
									<Textarea
										bind:value={$form.site_description}
										name="site_description"
										placeholder="Enter site description"
										aria-invalid={!!$errors.site_description}
										disabled={$submitting}
									/>
									{#if $errors.site_description}
										<Field.Error>{$errors.site_description}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_keywords">Site Keywords</Field.Label>
									<div class="relative">
										<Icon icon="lucide:tag" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<TagsInput
											bind:value={keywordsInput}
											name="site_keywords"
											type="text"
											class="pl-10"
											placeholder="Enter site keywords"
											aria-invalid={!!$errors.site_keywords}
											autocomplete="on"
											disabled={$submitting}
											onValueChange={(val) => {
												$form.site_keywords = val.join(',');
											}}
										/>
									</div>
									{#if $errors.site_keywords}
										<Field.Error>{$errors.site_keywords}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_email">Site Email</Field.Label>
									<div class="relative">
										<Icon icon="lucide:mail" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.site_email}
											name="site_email"
											type="email"
											class="ps-10"
											placeholder="Enter site email"
											aria-invalid={!!$errors.site_email}
											autocomplete="name"
											disabled={$submitting}
										/>
									</div>
									{#if $errors.site_email}
										<Field.Error>{$errors.site_email}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_phone">Site Phone</Field.Label>
									<div class="relative">
										<Icon icon="lucide:phone" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.site_phone}
											name="site_phone"
											type="text"
											class="ps-10"
											placeholder="Enter site phone"
											aria-invalid={!!$errors.site_phone}
											autocomplete="tel"
											disabled={$submitting}
											oninput={formatPhoneInput}
										/>
									</div>
									{#if $errors.site_phone}
										<Field.Error>{$errors.site_phone}</Field.Error>
									{/if}
								</Field.Field>
								<Field.Field>
									<Field.Label for="site_url">Site URL</Field.Label>
									<div class="relative">
										<Icon icon="lucide:link" class="absolute top-1/2 left-3 -translate-y-1/2" />
										<Input
											bind:value={$form.site_url}
											name="site_url"
											type="url"
											class="ps-10"
											placeholder="Enter site URL"
											aria-invalid={!!$errors.site_url}
											autocomplete="url"
											disabled={$submitting}
										/>
									</div>
									{#if $errors.site_url}
										<Field.Error>{$errors.site_url}</Field.Error>
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
				</Tabs.Content>
				<Tabs.Content value="update-web-logo">
					<AdminSettingUploadLogo web={data.settings?.WEBSITE} />
				</Tabs.Content>
				<Tabs.Content value="update-web-favicon">
					<AdminSettingUploadFavicon web={data.settings?.WEBSITE} />
				</Tabs.Content>
			</Tabs.Root>
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
