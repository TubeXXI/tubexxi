<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AdminSidebarLayout } from '@/components/admin';
	import { AppAlertDialog } from '@/components/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '@/components/ui/spinner/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '@/components/ui/switch/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import Icon from '@iconify/svelte';
	import { localizeHref } from '@/paraglide/runtime';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let errorMessage = $state<string | null>(null);

	// svelte-ignore state_referenced_locally
	const { form, errors, submitting, enhance } = superForm(data.form, {
		id: `create-application-form`,
		dataType: 'json',
		resetForm: false,
		onSubmit() {
			handleSubmitLoading(true);
			errorMessage = null;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			if (event.result.type === 'success') {
				handleSubmitLoading(false);
				await goto(localizeHref('/application'));
				await invalidateAll();
			}
		},
		onError(event) {
			errorMessage = event.result.error.message;
			handleSubmitLoading(false);
		}
	});

	function onPackageNameChange(val?: string) {
		if (!val) return;
		const packageNameRegex = /^[a-zA-Z0-9._-]+$/;
		if (!packageNameRegex.test(val)) {
			$errors.package_name = ['Invalid package name format'];
			return;
		}
	}

	const typeOptions = [
		{
			label: 'Android',
			value: 'android'
		},
		{
			label: 'iOS',
			value: 'ios'
		}
	];
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout
	page={`Update Application ${data.application.config.name}`}
	user={data.user}
	setting={data.settings}
>
	<div class="@container/main flex flex-col gap-4">
		<div class="flex-none px-4 py-4 sm:px-6">
			<div class="space-y-1">
				<h1 class="text-2xl font-bold tracking-tight sm:text-3xl">Update Application</h1>
				<p class="text-sm text-muted-foreground">
					Update application {data.application.config.name}.
				</p>
			</div>
		</div>
		<div class="space-y-4 rounded-md border border-neutral-300 px-3 py-5 dark:border-neutral-700">
			<form method="POST" class="relative space-y-6" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="is_active">
									Enable Application : <Badge variant={$form.is_active ? 'default' : 'destructive'}>
										{$form.is_active ? 'Enabled' : 'Disabled'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.is_active
										? 'Application will be active and available for use.'
										: 'Application will be inactive and not available for use.'}
								</Field.Description>
							</Field.Content>
							<input type="hidden" name="is_active" value={$form.is_active} />
							<Switch
								id="is_active"
								bind:checked={$form.is_active}
								name="is_active"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.is_active = val)}
							/>
						</Field.Field>
						<Field.Field>
							<Field.Label for="name">Name</Field.Label>
							<div class="relative">
								<Icon
									icon="material-symbols:title"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Input
									bind:value={$form.name}
									id="name"
									name="name"
									type="text"
									class="ps-10"
									placeholder="Enter application name"
									autocomplete="name"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.name}
								<Field.Error>{$errors.name}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="package_name">Package Name</Field.Label>
							<div class="relative">
								<Icon
									icon="material-symbols:package-2-outline-sharp"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Input
									bind:value={$form.package_name}
									id="package_name"
									name="package_name"
									type="text"
									class="ps-10"
									placeholder="Enter application package name"
									autocomplete="on"
									disabled={$submitting}
									oninput={(e) => onPackageNameChange((e.target as HTMLInputElement).value)}
								/>
							</div>
							{#if $errors.package_name}
								<Field.Error>{$errors.package_name}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="version">Version</Field.Label>
							<div class="relative">
								<Icon icon="ix:version-history" class="absolute top-1/2 left-3 -translate-y-1/2" />
								<Input
									bind:value={$form.version}
									id="version"
									name="version"
									type="text"
									class="ps-10"
									placeholder="Enter application version"
									autocomplete="on"
									disabled={$submitting}
								/>
							</div>
							{#if $errors.version}
								<Field.Error>{$errors.version}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Field.Label for="type">Type</Field.Label>
							<div class="relative">
								<Icon
									icon="tabler:device-mobile-code"
									class="absolute top-1/2 left-3 -translate-y-1/2"
								/>
								<Select.Root
									bind:value={$form.type}
									type="single"
									name="type"
									disabled={$submitting}
								>
									<Select.Trigger class="w-full ps-10 capitalize">
										{$form.type.replaceAll('-', ' ')}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											<Select.Label>Types</Select.Label>
											{#each typeOptions as option}
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
							</div>
							{#if $errors.type}
								<Field.Error>{$errors.type}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="enable_monetize">
									Enable Monetization : <Badge
										variant={$form.enable_monetize ? 'default' : 'destructive'}
									>
										{$form.enable_monetize ? 'Enabled' : 'Disabled'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.enable_monetize
										? 'Application will be monetized.'
										: 'Application will not be monetized.'}
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

						<Card.Root class="w-full">
							<Card.Header>
								<Card.Title>AdMob Configuration</Card.Title>
								<Card.Description>Enter your AdMob configuration below.</Card.Description>
							</Card.Header>
							<Card.Content>
								<Field.Group>
									<Field.Set>
										<Field.Field orientation="horizontal">
											<Field.Content>
												<Field.Label for="enable_admob">
													Enable AdMob : <Badge
														variant={$form.enable_admob ? 'default' : 'destructive'}
													>
														{$form.enable_admob ? 'Enabled' : 'Disabled'}
													</Badge>
												</Field.Label>
												<Field.Description>
													{$form.enable_admob
														? 'Application will use AdMob for monetization.'
														: 'Application will not use AdMob for monetization.'}
												</Field.Description>
											</Field.Content>
											<input type="hidden" name="enable_admob" value={$form.enable_admob} />
											<Switch
												id="enable_admob"
												bind:checked={$form.enable_admob}
												name="enable_admob"
												class="cursor-pointer"
												onCheckedChange={(val) => ($form.enable_admob = val)}
											/>
										</Field.Field>
										{#if $form.enable_admob}
											<Field.Set>
												<Field.Field>
													<Field.Label for="admob_id">AdMob Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_id}
															id="admob_id"
															name="admob_id"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_id}
														<Field.Error>{$errors.admob_id}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="admob_auto_ad">AdMob Auto Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_auto_ad}
															id="admob_auto_ad"
															name="admob_auto_ad"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Auto Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_auto_ad}
														<Field.Error>{$errors.admob_auto_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="admob_banner_ad">AdMob Banner Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_banner_ad}
															id="admob_banner_ad"
															name="admob_banner_ad"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Banner Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_banner_ad}
														<Field.Error>{$errors.admob_banner_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="admob_interstitial_ad"
														>AdMob Interstitial Ad Unit ID</Field.Label
													>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_interstitial_ad}
															id="admob_interstitial_ad"
															name="admob_interstitial_ad"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Interstitial Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_interstitial_ad}
														<Field.Error>{$errors.admob_interstitial_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="admob_native_ad">AdMob Native Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_native_ad}
															id="admob_native_ad"
															name="admob_native_ad"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Native Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_native_ad}
														<Field.Error>{$errors.admob_native_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="admob_rewarded_ad"
														>AdMob Rewarded Ad Unit ID</Field.Label
													>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.admob_rewarded_ad}
															id="admob_rewarded_ad"
															name="admob_rewarded_ad"
															type="text"
															class="ps-10"
															placeholder="Enter AdMob Rewarded Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.admob_rewarded_ad}
														<Field.Error>{$errors.admob_rewarded_ad}</Field.Error>
													{/if}
												</Field.Field>
											</Field.Set>
										{/if}
									</Field.Set>
								</Field.Group>
							</Card.Content>
						</Card.Root>

						<Card.Root class="w-full">
							<Card.Header>
								<Card.Title>Unity Ad Configuration</Card.Title>
								<Card.Description>Enter your Unity Ad configuration below.</Card.Description>
							</Card.Header>
							<Card.Content>
								<Field.Group>
									<Field.Set>
										<Field.Field orientation="horizontal">
											<Field.Content>
												<Field.Label for="enable_unity_ad">
													Enable Unity Ad : <Badge
														variant={$form.enable_unity_ad ? 'default' : 'destructive'}
													>
														{$form.enable_unity_ad ? 'Enabled' : 'Disabled'}
													</Badge>
												</Field.Label>
												<Field.Description>
													{$form.enable_unity_ad
														? 'Application will use Unity Ad for monetization.'
														: 'Application will not use Unity Ad for monetization.'}
												</Field.Description>
											</Field.Content>
											<input type="hidden" name="enable_unity_ad" value={$form.enable_unity_ad} />
											<Switch
												id="enable_unity_ad"
												bind:checked={$form.enable_unity_ad}
												name="enable_unity_ad"
												class="cursor-pointer"
												onCheckedChange={(val) => ($form.enable_unity_ad = val)}
											/>
										</Field.Field>
										{#if $form.enable_unity_ad}
											<Field.Set>
												<Field.Field>
													<Field.Label for="unity_ad_id">Unity Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.unity_ad_id}
															id="unity_ad_id"
															name="unity_ad_id"
															type="text"
															class="ps-10"
															placeholder="Enter Unity Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.unity_ad_id}
														<Field.Error>{$errors.unity_ad_id}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="unity_banner_ad">Unity Banner Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.unity_banner_ad}
															id="unity_banner_ad"
															name="unity_banner_ad"
															type="text"
															class="ps-10"
															placeholder="Enter Unity Banner Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.unity_banner_ad}
														<Field.Error>{$errors.unity_banner_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="unity_interstitial_ad"
														>Unity Interstitial Ad Unit ID</Field.Label
													>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.unity_interstitial_ad}
															id="unity_interstitial_ad"
															name="unity_interstitial_ad"
															type="text"
															class="ps-10"
															placeholder="Enter Unity Interstitial Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.unity_interstitial_ad}
														<Field.Error>{$errors.unity_interstitial_ad}</Field.Error>
													{/if}
												</Field.Field>
												<Field.Field>
													<Field.Label for="unity_rewarded_ad"
														>Unity Rewarded Ad Unit ID</Field.Label
													>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.unity_rewarded_ad}
															id="unity_rewarded_ad"
															name="unity_rewarded_ad"
															type="text"
															class="ps-10"
															placeholder="Enter Unity Rewarded Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.unity_rewarded_ad}
														<Field.Error>{$errors.unity_rewarded_ad}</Field.Error>
													{/if}
												</Field.Field>
											</Field.Set>
										{/if}
									</Field.Set>
								</Field.Group>
							</Card.Content>
						</Card.Root>

						<Card.Root class="w-full">
							<Card.Header>
								<Card.Title>Start App Ad Configuration</Card.Title>
								<Card.Description>Enter your Start App Ad configuration below.</Card.Description>
							</Card.Header>
							<Card.Content>
								<Field.Group>
									<Field.Set>
										<Field.Field orientation="horizontal">
											<Field.Content>
												<Field.Label for="enable_star_io_ad">
													Enable Start App : <Badge
														variant={$form.enable_star_io_ad ? 'default' : 'destructive'}
													>
														{$form.enable_star_io_ad ? 'Enabled' : 'Disabled'}
													</Badge>
												</Field.Label>
												<Field.Description>
													{$form.enable_star_io_ad
														? 'Application will use Start App for monetization.'
														: 'Application will not use Start App for monetization.'}
												</Field.Description>
											</Field.Content>
											<input
												type="hidden"
												name="enable_star_io_ad"
												value={$form.enable_star_io_ad}
											/>
											<Switch
												id="enable_star_io_ad"
												bind:checked={$form.enable_star_io_ad}
												name="enable_star_io_ad"
												class="cursor-pointer"
												onCheckedChange={(val) => ($form.enable_star_io_ad = val)}
											/>
										</Field.Field>
										{#if $form.enable_star_io_ad}
											<Field.Set>
												<Field.Field>
													<Field.Label for="star_io_ad_id">Start App Ad Unit ID</Field.Label>
													<div class="relative">
														<Icon
															icon="mdi:code"
															class="absolute top-1/2 left-3 -translate-y-1/2"
														/>
														<Input
															bind:value={$form.star_io_ad_id}
															id="star_io_ad_id"
															name="star_io_ad_id"
															type="text"
															class="ps-10"
															placeholder="Enter Start App Ad Unit ID"
															autocomplete="on"
															disabled={$submitting}
														/>
													</div>
													{#if $errors.star_io_ad_id}
														<Field.Error>{$errors.star_io_ad_id}</Field.Error>
													{/if}
												</Field.Field>
											</Field.Set>
										{/if}
									</Field.Set>
								</Field.Group>
							</Card.Content>
						</Card.Root>

						<Card.Root class="w-full">
							<Card.Header>
								<Card.Title>One Signal Configuration</Card.Title>
								<Card.Description>Enter your One Signal configuration below.</Card.Description>
							</Card.Header>
							<Card.Content>
								<Field.Group>
									<Field.Set>
										<Field.Field>
											<Field.Label for="one_signal_id">One Signal ID</Field.Label>
											<div class="relative">
												<Icon
													icon="icon-park-outline:signal-one"
													class="absolute top-1/2 left-3 -translate-y-1/2"
												/>
												<Input
													bind:value={$form.one_signal_id}
													id="one_signal_id"
													name="one_signal_id"
													type="text"
													class="ps-10"
													placeholder="Enter One Signal ID"
													autocomplete="on"
													disabled={$submitting}
												/>
											</div>
											{#if $errors.one_signal_id}
												<Field.Error>{$errors.one_signal_id}</Field.Error>
											{/if}
										</Field.Field>
									</Field.Set>
								</Field.Group>
							</Card.Content>
						</Card.Root>

						<Field.Field orientation="horizontal">
							<Field.Content>
								<Field.Label for="enable_in_app_purchase">
									Enable In App Purchase : <Badge
										variant={$form.enable_in_app_purchase ? 'default' : 'destructive'}
									>
										{$form.enable_in_app_purchase ? 'Enabled' : 'Disabled'}
									</Badge>
								</Field.Label>
								<Field.Description>
									{$form.enable_in_app_purchase
										? 'Application will use In App Purchase for monetization.'
										: 'Application will not use In App Purchase for monetization.'}
								</Field.Description>
							</Field.Content>
							<input
								type="hidden"
								name="enable_in_app_purchase"
								value={$form.enable_in_app_purchase}
							/>
							<Switch
								id="enable_in_app_purchase"
								bind:checked={$form.enable_in_app_purchase}
								name="enable_in_app_purchase"
								class="cursor-pointer"
								onCheckedChange={(val) => ($form.enable_in_app_purchase = val)}
							/>
						</Field.Field>
					</Field.Set>
					<Field.Field>
						<Button type="submit" disabled={$submitting} class="w-full">
							{#if $submitting}
								<Spinner class="mr-2" />
								Please wait
							{:else}
								Update application
							{/if}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		</div>
	</div>
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
