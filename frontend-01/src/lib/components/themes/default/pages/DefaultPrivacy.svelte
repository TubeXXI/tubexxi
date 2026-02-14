<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { localizeHref } from '@/paraglide/runtime';
	import * as i18n from '@/paraglide/messages.js';
	import { smoothScroll } from '$lib/stores';

	let {
		user,
		setting,
		lang = 'en'
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
		lang?: string;
	} = $props();

	let webSetting = $derived(setting?.WEBSITE);
	let isScrolling = $state(false);

	const handleScroll = (id: string, offset: number = 100) => {
		smoothScroll.scrollToAnchor(id, offset);
	};

	onMount(() => {
		const unsubscribe = smoothScroll.subscribe((state) => {
			isScrolling = state.isScrolling;
		});

		return () => unsubscribe();
	});
</script>

<DefaultSidebarLayout {user} {setting} {lang}>
	<div class="scrollbar-primary flex min-h-[calc(100vh-160px)] flex-col">
		<section class="space-y-6 px-0 py-16 md:px-4">
			<div class="container mx-auto px-4 text-center">
				<h1 class="mb-6 text-5xl font-bold">{i18n.privacy_policy()}</h1>
				<p class="mx-auto max-w-3xl text-lg">
					{i18n.privacy_policy_description({ site_name: webSetting?.site_name || 'our site' })}
				</p>
			</div>
			<div class="mx-auto flex justify-center">
				<Card.Root class="w-full shadow-lg sm:max-w-sm md:max-w-6xl">
					<Card.Content class="flex flex-col items-center justify-center rounded-lg p-2 md:p-6">
						<div class="grid grid-cols-1 gap-12 lg:grid-cols-4">
							<!-- Table of Contents -->
							<div class="lg:col-span-1">
								<div class="sticky top-6 rounded-lg bg-sky-100 p-6 shadow-lg dark:bg-sky-900/30">
									<h3 class="mb-4 text-lg font-bold">{i18n.privacy_policy_table_of_contents()}</h3>
									<nav class="grid gap-2">
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#privacy_policy_introduction')}
										>
											{i18n.privacy_policy_introduction()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#personal_information')}
										>
											{i18n.privacy_policy_personal_information()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#usage_data')}
										>
											{i18n.privacy_policy_usage_data()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#cookies')}
										>
											{i18n.privacy_policy_cookies()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#cookies_types')}
										>
											{i18n.privacy_policy_cookies_types()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#data_storage')}
										>
											{i18n.privacy_policy_data_storage()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#advertisement')}
										>
											{i18n.privacy_policy_advertisement()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#data_retention')}
										>
											{i18n.privacy_policy_data_retention()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#your_right')}
										>
											{i18n.privacy_policy_your_right()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#children')}
										>
											{i18n.privacy_policy_children()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#privacy_policy_international')}
										>
											{i18n.privacy_policy_international()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#privacy_policy_changes')}
										>
											{i18n.privacy_policy_changes()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#contact')}
										>
											{i18n.privacy_policy_contact()}
										</button>
									</nav>
								</div>
							</div>
							<!-- Privacy Content -->
							<div class=" lg:col-span-3">
								<div
									class="privacy-content space-y-8 rounded-lg bg-white p-8 shadow-lg dark:bg-neutral-950"
								>
									<div id="privacy_policy_introduction">
										<h2 class="text-3xl font-bold">{i18n.privacy_policy_introduction()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_introduction_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="personal_information">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_personal_information()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_personal_information_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="usage_data">
										<h2 class=" text-3xl font-bold">
											{i18n.privacy_policy_usage_data()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_usage_data_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="cookies">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_cookies()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_cookies_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="cookies_types">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_cookies_types()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_cookies_types_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_cookies_types_description_2()}
											</p>
										</div>
									</div>
									<div id="data_storage">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_data_storage()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_2()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_3()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_4()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_5()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_6({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_storage_description_7()}
											</p>
										</div>
									</div>
									<div id="advertisement">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_advertisement()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_advertisement_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="data_retention">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_data_retention()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description_2({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description_3({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description_4()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description_5({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_data_retention_description_6()}
											</p>
										</div>
									</div>
									<div id="your_right">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_your_right()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_your_right_description()}
											</p>
										</div>
									</div>
									<div id="children">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_children()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_children_description()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_children_description_2()}
											</p>
										</div>
									</div>
									<div id="privacy_policy_international">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_international()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_international_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_international_description_2({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_international_description_3({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="privacy_policy_changes">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_changes()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_changes_description()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_changes_description_2()}
											</p>
										</div>
									</div>
									<div id="contact">
										<h2 class="text-3xl font-bold">
											{i18n.privacy_policy_contact()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.privacy_policy_contact_description()}
											</p>
										</div>
									</div>
									<section class="text-center">
										<div class="rounded-lg p-12">
											<h2 class="mb-4 text-3xl font-bold">
												{i18n.about_get_in_touch({
													site_name: webSetting?.site_name || 'our site'
												})}
											</h2>
											<p class="mb-8 text-lg text-muted-foreground">
												{i18n.terms_contact_description()}
											</p>
											<Button
												href={localizeHref('/contact')}
												variant="destructive"
												class="text-white"
											>
												{i18n.contact_us()}
											</Button>
										</div>
									</section>
								</div>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</section>
	</div>
</DefaultSidebarLayout>
