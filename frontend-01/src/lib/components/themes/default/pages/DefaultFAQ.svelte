<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';
	import * as Accordion from '$lib/components/ui/accordion/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import Icon from '@iconify/svelte';
	import * as i18n from '@/paraglide/messages.js';

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
				<h1 class="mb-6 text-5xl font-bold">{i18n.faq()}</h1>
				<p class="mx-auto max-w-3xl text-lg">
					{i18n.faq_description({ site_name: webSetting?.site_name || 'our site' })}
				</p>
			</div>
			<div class="mx-auto flex justify-center">
				<Card.Root class="w-full shadow-lg sm:max-w-sm md:max-w-2xl lg:max-w-4xl">
					<Card.Content class="flex flex-col items-center justify-start rounded-lg p-8 ">
						<Accordion.Root type="single" class="w-full">
							<Accordion.Item value="item-1" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_1({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_1({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-2" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_2({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_2({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-3" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_3({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_3({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-4" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_4({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_4({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-5" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_5({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_5({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-6" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_6({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_6({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
							<Accordion.Item value="item-7" class="mt-2 flex flex-col items-center justify-start">
								<Accordion.Trigger class="w-full">
									<Icon icon="wpf:faq" class="mr-2 h-5 w-5" />
									{i18n.faq_question_7({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Trigger>
								<Accordion.Content class="mt-2">
									{i18n.faq_answer_7({ site_name: webSetting?.site_name || 'our site' })}
								</Accordion.Content>
							</Accordion.Item>
						</Accordion.Root>
					</Card.Content>
				</Card.Root>
			</div>
		</section>
	</div>
</DefaultSidebarLayout>
