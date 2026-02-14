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
		<section class="space-y-6 px-2 py-16 md:px-4">
			<div class="container mx-auto px-4 text-center">
				<h1 class="mb-6 text-5xl font-bold">{i18n.terms_of_service()}</h1>
				<p class="mx-auto max-w-3xl text-lg">
					{i18n.terms_of_service_description({ site_name: webSetting?.site_name || 'our site' })}
				</p>
			</div>
			<div class="mx-auto flex justify-center">
				<Card.Root class="w-full shadow-lg sm:max-w-sm md:max-w-6xl">
					<Card.Content class="flex flex-col items-center justify-center rounded-lg p-2 md:p-6">
						<div class="grid grid-cols-1 gap-12 lg:grid-cols-4">
							<!-- Table of Contents -->
							<div class="lg:col-span-1">
								<div class="sticky top-6 rounded-lg bg-sky-100 p-6 shadow-lg dark:bg-sky-900/30">
									<h3 class="mb-4 text-lg font-bold">{i18n.terms_table_of_contents()}</h3>
									<nav class="grid items-start gap-2 text-start">
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_introduction')}
										>
											{i18n.terms_introduction()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_acceptance')}
										>
											{i18n.terms_acceptance()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_usage')}
										>
											{i18n.terms_usage()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_content_and_intellectual_property')}
										>
											{i18n.terms_content_and_intellectual_property()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_privacy')}
										>
											{i18n.terms_privacy()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_disclaimer')}
										>
											{i18n.terms_disclaimer()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_limitation')}
										>
											{i18n.terms_limitation()}
										</button>
										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_termination')}
										>
											{i18n.terms_termination()}
										</button>

										<button
											type="button"
											class="cursor-pointer justify-self-start text-start text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300"
											onclick={() => handleScroll('#terms_contact')}
										>
											{i18n.terms_contact()}
										</button>
									</nav>
								</div>
							</div>
							<!-- Terms Content -->
							<div class=" lg:col-span-3">
								<div
									class="privacy-content space-y-8 rounded-lg bg-white p-8 shadow-lg dark:bg-neutral-950"
								>
									<div id="terms_introduction">
										<h2 class="text-3xl font-bold">{i18n.terms_introduction()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_introduction_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="terms_acceptance">
										<h2 class="text-3xl font-bold">{i18n.terms_acceptance()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_acceptance_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="terms_usage">
										<h2 class="text-3xl font-bold">{i18n.terms_usage()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description_2({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description_3({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description_4({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description_5({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_usage_description_6({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
										</div>
									</div>
									<div id="terms_content_and_intellectual_property">
										<h2 class="text-3xl font-bold">
											{i18n.terms_content_and_intellectual_property()}
										</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_content_and_intellectual_property_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_content_and_intellectual_property_description_2()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_content_and_intellectual_property_description_3()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_content_and_intellectual_property_description_4()}
											</p>
										</div>
									</div>
									<div id="terms_privacy">
										<h2 class="text-3xl font-bold">{i18n.terms_privacy()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_privacy_description()}
											</p>
										</div>
									</div>
									<div id="terms_disclaimer">
										<h2 class="text-3xl font-bold">{i18n.terms_disclaimer()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_disclaimer_description({
													site_name: webSetting?.site_name || 'our site'
												})}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_disclaimer_description_2()}
											</p>
										</div>
									</div>
									<div id="terms_limitation">
										<h2 class="text-3xl font-bold">{i18n.terms_limitation()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_limitation_description()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_limitation_description_2()}
											</p>
										</div>
									</div>
									<div id="terms_termination">
										<h2 class="text-3xl font-bold">{i18n.terms_termination()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_termination_description()}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_termination_description_2()}
											</p>
										</div>
									</div>
									<div id="terms_gover_law">
										<h2 class="text-3xl font-bold">{i18n.terms_gover_law()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_gover_law_description({ country: 'our country' })}
											</p>
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_gover_law_description_2({ country: 'our country' })}
											</p>
										</div>
									</div>
									<div id="terms_contact">
										<h2 class="text-3xl font-bold">{i18n.terms_contact()}</h2>
										<div class="mt-4 space-y-4">
											<p class="leading-relaxed text-muted-foreground">
												{i18n.terms_contact_description()}
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
						</div></Card.Content
					>
				</Card.Root>
			</div>
		</section>
	</div>
</DefaultSidebarLayout>
