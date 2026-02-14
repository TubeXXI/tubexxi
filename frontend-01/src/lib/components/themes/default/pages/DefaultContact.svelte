<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';
	import type { SuperValidated } from 'sveltekit-superforms';
	import { invalidateAll } from '$app/navigation';
	import { superForm } from 'sveltekit-superforms';
	import type { ContactSchema } from '@/utils/schema';
	import { handleSubmitLoading } from '@/stores';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '@/components/ui/textarea/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import Icon from '@iconify/svelte';
	import * as i18n from '@/paraglide/messages.js';
	import { toast } from '@/stores/toast';
	import { localizeHref } from '@/paraglide/runtime';

	let {
		user,
		setting,
		lang = 'en',
		form: formData
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
		lang?: string;
		form: SuperValidated<ContactSchema>;
	} = $props();

	let webSetting = $derived(setting?.WEBSITE);

	// svelte-ignore state_referenced_locally
	const { form, enhance, errors, submitting } = superForm(formData!, {
		resetForm: true,
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
			}
		},
		onError: ({ result }) => {
			handleSubmitLoading(false);
			toast.error(result.error?.message || i18n.contact_error());
		}
	});

	const socialLinks = [
		{
			name: 'Facebook',
			icon: 'mdi:facebook',
			href: 'https://www.facebook.com'
		},
		{
			name: 'Twitter',
			icon: 'mdi:twitter',
			href: 'https://twitter.com'
		},
		{
			name: 'Instagram',
			icon: 'mdi:instagram',
			href: 'https://www.instagram.com'
		},
		{
			name: 'YouTube',
			icon: 'mdi:youtube',
			href: 'https://www.youtube.com'
		},
		{ name: 'LinkedIn', icon: 'mdi:linkedin', href: 'https://www.linkedin.com' }
	];
</script>

<DefaultSidebarLayout {user} {setting} {lang}>
	<div class="scrollbar-primary flex min-h-[calc(100vh-160px)] flex-col">
		<section class="py-16">
			<div class="container mx-auto px-4 text-center">
				<h1 class="mb-6 text-5xl font-bold">{i18n.contact_us()}</h1>
				<p class="mx-auto max-w-3xl text-lg">
					{i18n.contact_us_description({ site_name: webSetting?.site_name || 'our site' })}
				</p>
			</div>
		</section>
		<div
			class="flex w-full flex-col items-center justify-center bg-neutral-400 dark:bg-neutral-700"
		>
			<Card.Root class="w-full">
				<Card.Content>
					<div class="grid gap-4">
						<div class="rounded-lg bg-white p-8 shadow-lg backdrop-blur-md dark:bg-neutral-950">
							<h2 class="mb-6 text-3xl font-bold text-neutral-800 dark:text-neutral-100">
								{i18n.send_us_a_message()}
							</h2>
							<form method="POST" class="space-y-4" use:enhance>
								<Field.Group>
									<Field.Set>
										<Field.Field>
											<Field.Label for="name">
												{i18n.name()}
												<span class="text-red-500 dark:text-red-400">*</span>
											</Field.Label>
											<div class="relative">
												<Icon icon="mdi:account" class="absolute top-1/2 left-3 -translate-y-1/2" />
												<Input
													bind:value={$form.name}
													name="name"
													type="text"
													class="ps-10"
													placeholder={i18n.name()}
													aria-invalid={!!$errors.name}
													autocomplete="name"
													disabled={$submitting}
												/>
											</div>
											{#if $errors.name}
												<Field.Error>{$errors.name}</Field.Error>
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
											<Field.Label for="subject">
												{i18n.subject()}
												<span class="text-red-500 dark:text-red-400">*</span>
											</Field.Label>
											<div class="relative">
												<Icon
													icon="mdi:format-title"
													class="absolute top-1/2 left-3 -translate-y-1/2"
												/>
												<Input
													bind:value={$form.subject}
													name="subject"
													type="text"
													class="ps-10"
													placeholder={i18n.subject()}
													aria-invalid={!!$errors.subject}
													autocomplete="on"
													disabled={$submitting}
												/>
											</div>
											{#if $errors.subject}
												<Field.Error>{$errors.subject}</Field.Error>
											{/if}
										</Field.Field>
										<Field.Field>
											<Field.Label for="message">
												{i18n.message()}
												<span class="text-red-500 dark:text-red-400">*</span>
											</Field.Label>
											<Textarea
												bind:value={$form.message}
												name="message"
												placeholder={i18n.message()}
												aria-invalid={!!$errors.message}
												autocomplete="on"
												disabled={$submitting}
											/>
											{#if $errors.message}
												<Field.Error>{$errors.message}</Field.Error>
											{/if}
										</Field.Field>
									</Field.Set>
									<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
										<Button type="submit" disabled={$submitting}>
											{#if $submitting}
												<Spinner class="mr-2 size-5" />
											{/if}
											{$submitting ? i18n.please_wait() : i18n.send_us_a_message()}
										</Button>
									</Field.Field>
								</Field.Group>
							</form>
						</div>
						<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
							<div class="mb-6 rounded-lg bg-background p-6 shadow-lg">
								<h3 class="mb-4 text-xl font-bold text-neutral-800 dark:text-neutral-100">
									{i18n.contact_us()}
								</h3>
								<div class="space-y-4">
									<div class="flex items-start space-x-3">
										<div
											class="rounded-full bg-red-100 p-2 text-red-600 dark:bg-red-700 dark:text-red-100"
										>
											<Icon icon="fa:map-marker" />
										</div>
										<div>
											<p class="font-semibold">{i18n.address()}</p>
											<p class="text-sm text-muted-foreground">
												1211 Avenue of the Americas<br />New York, NY 10036
											</p>
										</div>
									</div>
									<div class="flex items-start space-x-3">
										<div
											class="rounded-full bg-red-100 p-2 text-red-600 dark:bg-red-700 dark:text-red-100"
										>
											<Icon icon="fa:phone" />
										</div>
										<div>
											<p class="font-semibold">{i18n.phone()}</p>
											<p class="text-sm text-muted-foreground">-</p>
										</div>
									</div>
									<div class="flex items-start space-x-3">
										<div
											class="rounded-full bg-red-100 p-2 text-red-600 dark:bg-red-700 dark:text-red-100"
										>
											<Icon icon="fa:envelope" />
										</div>
										<div>
											<p class="font-semibold">{i18n.email()}</p>
											<p class="text-sm text-muted-foreground">
												{webSetting?.site_email}
											</p>
										</div>
									</div>
								</div>
							</div>
							<div
								class="mb-6 rounded-lg border border-red-200 bg-red-50 p-6 shadow-lg dark:border-red-400 dark:bg-red-900/50"
							>
								<h3 class="mb-4 text-xl font-bold text-red-800 dark:text-red-100">
									<Icon icon="fa:lightbulb-o" class="mr-2" />
									{i18n.contact_news_tips()}
								</h3>
								<p class="mb-4 text-sm text-red-700 dark:text-red-300">
									{i18n.contact_news_tips_description()}
								</p>
								<div class="space-y-2">
									<p class="font-semibold text-red-600 dark:text-red-400">
										{i18n.email()}: {webSetting?.site_email}
									</p>
									<p class="font-semibold text-red-600 dark:text-red-400">
										{i18n.phone()}: {webSetting?.site_phone}
									</p>
								</div>
							</div>
							<div class="mb-6 rounded-lg bg-background p-6 shadow-lg">
								<h3 class="mb-4 text-xl font-bold text-neutral-800 dark:text-neutral-100">
									{i18n.contact_press_inquiries()}
								</h3>
								<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
									{i18n.contact_press_inquiries_description()}
								</p>
								<div class="space-y-2">
									<p class="font-semibold text-neutral-800 dark:text-neutral-100">
										{i18n.email()}: {webSetting?.site_email}
									</p>
									<p class="font-semibold text-neutral-800 dark:text-neutral-100">
										{i18n.phone()}: -
									</p>
								</div>
							</div>
							<div class="mb-6 rounded-lg bg-background p-6 shadow-lg">
								<h3 class="mb-4 text-xl font-bold text-neutral-800 dark:text-neutral-100">
									{i18n.contact_advertising()}
								</h3>
								<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
									{i18n.contact_advertising_description({ site_name: webSetting?.site_name || '' })}
								</p>
								<div class="space-y-2">
									<p class="font-semibold text-neutral-800 dark:text-neutral-100">
										{i18n.email()}: {webSetting?.site_email}
									</p>
									<p class="font-semibold text-neutral-800 dark:text-neutral-100">
										{i18n.phone()}: {webSetting?.site_phone || '-'}
									</p>
								</div>
							</div>
							<div class="rounded-lg bg-background p-6 shadow-lg">
								<h3 class="mb-4 text-xl font-bold text-neutral-800 dark:text-neutral-100">
									{i18n.contact_follow_us()}
								</h3>
								<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
									{i18n.contact_follow_us_description({ site_name: webSetting?.site_name || '' })}
								</p>
								<div class="flex gap-2">
									{#each socialLinks as social}
										<a
											href={social.href}
											target="_blank"
											rel="noopener noreferrer"
											class="inline-block rounded-full bg-red-600 p-2 text-white dark:bg-red-700"
										>
											<Icon icon={social.icon} />
										</a>
									{/each}
								</div>
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	</div>
</DefaultSidebarLayout>
