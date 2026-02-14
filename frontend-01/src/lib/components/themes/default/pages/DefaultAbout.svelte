<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import * as i18n from '@/paraglide/messages.js';
	import { localizeHref } from '@/paraglide/runtime';

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
</script>

<DefaultSidebarLayout {user} {setting} {lang}>
	<div class="scrollbar-primary flex min-h-[calc(100vh-160px)] flex-col">
		<section class="space-y-6 px-2 py-16 md:px-4">
			<div class="container mx-auto px-4 text-center">
				<h1 class="mb-6 text-5xl font-bold">{i18n.about_us()}</h1>
				<p class="mx-auto max-w-3xl text-lg">
					{i18n.about_us_description({ site_name: webSetting?.site_name || 'our site' })}
				</p>
			</div>
			<div class="mx-auto flex justify-center">
				<Card.Root class="w-full shadow-lg sm:max-w-sm md:max-w-5xl">
					<Card.Content class="flex flex-col items-center justify-center rounded-lg p-8">
						<!-- Mission Statement -->
						<section class="mb-16">
							<div class="mb-12 text-center">
								<h2 class="mb-4 text-4xl font-bold">
									{i18n.about_mission_statement({ site_name: webSetting?.site_name || 'our site' })}
								</h2>
								<div class="mx-auto h-1 w-24 bg-red-600"></div>
							</div>

							<div class="grid grid-cols-1 items-center gap-12 md:grid-cols-2">
								<ul class="list-disc space-y-4 ps-2">
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_mission_statement_point_1({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_mission_statement_point_2({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_mission_statement_point_3({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_mission_statement_point_4({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_mission_statement_point_5({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
								</ul>
								<div>
									<img
										src="/images/pages/about_mission.avif"
										alt="Newsroom"
										decoding="auto"
										loading="lazy"
										class="rounded-lg shadow-lg"
										style="color: transparent"
									/>
								</div>
							</div>
						</section>
						<!-- History -->
						<section class="mb-16">
							<div class="mb-12 text-center">
								<h2 class="mb-4 text-4xl font-bold">
									{i18n.about_history({ site_name: webSetting?.site_name || 'our site' })}
								</h2>
								<div class="mx-auto h-1 w-24 bg-red-600"></div>
							</div>
							<div class="grid grid-cols-1 items-center gap-12 md:grid-cols-2">
								<ul class="list-disc space-y-4 ps-2">
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_history_description({
												year: webSetting?.site_created_at?.substring(0, 4) || '2023',
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_history_description_2({
												year: webSetting?.site_created_at?.substring(0, 4) || '2023',
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_history_description_3({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_history_description_4({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_history_description_5({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
								</ul>
								<div>
									<img
										src="/images/pages/about_history.jpg"
										alt="Newsroom"
										decoding="auto"
										loading="lazy"
										class="rounded-lg shadow-lg"
										style="color: transparent"
									/>
								</div>
							</div>
						</section>
						<!-- Values -->
						<section class="mb-16">
							<div class="mb-12 text-center">
								<h2 class="mb-4 text-4xl font-bold">
									{i18n.about_values({ site_name: webSetting?.site_name || 'our site' })}
								</h2>
								<div class="mx-auto h-1 w-24 bg-red-600"></div>
							</div>
							<div class="grid grid-cols-1 items-center gap-12 md:grid-cols-2">
								<ul class="list-disc space-y-4 ps-2">
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_values_point_1({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_values_point_2({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_values_point_3({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_values_point_4({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
									<li>
										<p class="leading-relaxed text-muted-foreground">
											{i18n.about_values_point_5({
												site_name: webSetting?.site_name || 'our site'
											})}
										</p>
									</li>
								</ul>
								<div>
									<img
										src="/images/pages/about_values.webp"
										alt="Values"
										decoding="auto"
										loading="lazy"
										class="rounded-lg shadow-lg"
										style="color: transparent"
									/>
								</div>
							</div>
						</section>
						<section class="text-center">
							<div class="rounded-lg p-12">
								<h2 class="mb-4 text-3xl font-bold">
									{i18n.about_get_in_touch({ site_name: webSetting?.site_name || 'our site' })}
								</h2>
								<p class="mb-8 text-lg text-muted-foreground">
									{i18n.terms_contact_description()}
								</p>
								<Button href={localizeHref('/contact')} variant="destructive" class="text-white">
									{i18n.contact_us()}
								</Button>
							</div>
						</section>
					</Card.Content>
				</Card.Root>
			</div>
		</section>
	</div>
</DefaultSidebarLayout>
