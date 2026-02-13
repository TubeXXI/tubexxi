<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import {
		AdminSidebarLayout,
		AdminAccountLayout,
		AdminAccountUploadAvatar
	} from '@/components/admin';
	import { AppAlertDialog } from '@/components/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import Icon from '@iconify/svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Badge } from '@/components/ui/badge/index.js';
	import { Mail, Calendar, Clock, Info } from '@lucide/svelte';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
	// svelte-ignore state_referenced_locally
	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(data.profileForm, {
		resetForm: false,
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

				await invalidateAll();
				setTimeout(async () => {
					successMessage = null;
					errorMessage = null;
				}, 2000);
			}
		},
		onError: ({ result }) => {
			handleSubmitLoading(false);
			errorMessage = result.error?.message || 'Validation error';
		}
	});

	function formatDate(dateString?: Date | string): string {
		if (!dateString) {
			return 'N/A';
		}
		if (dateString instanceof String) {
			return new Date(dateString).toLocaleString('en-US', {
				year: 'numeric',
				month: 'long',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		}
		return dateString.toLocaleString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="Update Profile" user={data.user} setting={data.settings}>
	<AdminAccountLayout title="Update Profile" description="Update your profile information">
		{#snippet children()}
			<div class="space-y-4">
				<Card.Root class="mb-5">
					<Card.Content>
						<div class="flex flex-col items-start gap-6 md:flex-row md:items-center">
							<div class="relative">
								<Avatar.Root class="h-24 w-24">
									<Avatar.Image
										src={data.user?.avatar_url || ''}
										alt={data.user?.full_name || 'User Avatar'}
									/>
									<Avatar.Fallback
										>{data.user?.full_name?.split(' ')[0]?.slice(0, 2).toUpperCase() ||
											'User'}</Avatar.Fallback
									>
								</Avatar.Root>
								<AdminAccountUploadAvatar user={data.user} />
							</div>
							<div class="flex-1 space-y-2">
								<div class="flex flex-col gap-2 md:flex-row md:items-center">
									<h1 class="text-2xl font-bold text-neutral-900 capitalize dark:text-white">
										{data.user?.full_name || 'N/A'}
									</h1>
								</div>
								<div class="mt-1 flex items-center gap-2">
									<Badge variant="default" class="px-4 text-xs font-semibold">
										{data.user?.role?.name || 'N/A'}
									</Badge>
								</div>
								<div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
									<Tooltip.Provider>
										<Tooltip.Root>
											<Tooltip.Trigger
												class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
											>
												<Mail class="size-4" />
												{data.user?.email || 'N/A'}
											</Tooltip.Trigger>
											<Tooltip.Content>
												<p>Email address</p>
											</Tooltip.Content>
										</Tooltip.Root>
									</Tooltip.Provider>
									<Tooltip.Provider>
										<Tooltip.Root>
											<Tooltip.Trigger
												class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
											>
												<Calendar class="size-4" />
												{data.user?.created_at ? formatDate(data.user?.created_at) : 'N/A'}
											</Tooltip.Trigger>
											<Tooltip.Content>
												<p>
													Member since {data.user?.created_at
														? formatDate(data.user?.created_at)
														: 'N/A'}
												</p>
											</Tooltip.Content>
										</Tooltip.Root>
									</Tooltip.Provider>
									<Tooltip.Provider>
										<Tooltip.Root>
											<Tooltip.Trigger
												class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
											>
												<Clock class="size-4" />
												{data.user?.last_login_at ? formatDate(data.user?.last_login_at) : 'N/A'}
											</Tooltip.Trigger>
											<Tooltip.Content>
												<p>
													Last login on {data.user?.last_login_at
														? formatDate(data.user?.last_login_at)
														: 'N/A'}
												</p>
											</Tooltip.Content>
										</Tooltip.Root>
									</Tooltip.Provider>
								</div>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
				<form method="POST" class="space-y-4" use:enhance>
					<Field.Group>
						<Field.Set>
							<Field.Group>
								<Field.Set>
									<Field.Legend>Personal Information</Field.Legend>
									<Field.Description>Update your personal information</Field.Description>
									<Field.Group>
										<Field.Field>
											<Field.Label for="full_name">
												Full Name
												<span class="text-red-500 dark:text-red-400">*</span>
											</Field.Label>
											<div class="relative">
												<Icon icon="mdi:account" class="absolute top-1/2 left-3 -translate-y-1/2" />
												<Input
													bind:value={$form.full_name}
													name="full_name"
													type="text"
													class="ps-10"
													placeholder="Enter your full name"
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
												Email
												<span class="text-red-500 dark:text-red-400">*</span>
											</Field.Label>
											<div class="relative">
												<Icon icon="mdi:email" class="absolute top-1/2 left-3 -translate-y-1/2" />
												<Input
													bind:value={$form.email}
													name="email"
													type="email"
													class="ps-10"
													placeholder="Enter your email"
													aria-invalid={!!$errors.email}
													autocomplete="email"
													disabled={$submitting}
												/>
											</div>
											{#if $errors.email}
												<Field.Error>{$errors.email}</Field.Error>
											{/if}
										</Field.Field>
									</Field.Group>
								</Field.Set>
								<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
									<Button type="submit" class="w-full" disabled={$submitting}>
										{#if $submitting}
											<Spinner class="mr-2 size-5" />
										{/if}
										{$submitting ? 'Please wait' : 'Update profile'}
									</Button>
								</Field.Field>
							</Field.Group>
						</Field.Set>
					</Field.Group>
				</form>
			</div>
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
