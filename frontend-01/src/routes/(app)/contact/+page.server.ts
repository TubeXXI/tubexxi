import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { contactSchema } from '$lib/utils/schema';
import { fail } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Contact Us', lang)) as SingleResponse;
	const description = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_description || '',
		lang
	)) as SingleResponse;
	const tagline = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_tagline || '',
		lang
	)) as SingleResponse;
	const keywords = await Promise.all(
		(settings?.WEBSITE?.site_keywords?.split(',') || ['']).map(
			async (keyword) =>
				(await deps.languageHelper.singleTranslate(keyword.trim(), lang)) as SingleResponse
		)
	);

	const pageMetaTags = defaultMetaTags(
		{
			path_url: defaultOrigin,
			title: `${capitalizeFirstLetter(title.data.target.text)} - ${settings?.WEBSITE?.site_name || ''}`,
			tagline: capitalizeFirstLetter(tagline.data.target.text),
			description: capitalizeFirstLetter(description.data.target.text),
			keywords: keywords.map((keyword: SingleResponse) =>
				capitalizeFirstLetter(keyword.data.target.text)
			),
			robots: 'index, follow',
			canonical: defaultOrigin,
			alternates,
			graph_type: 'website',
			language: lang
		},
		settings
	);

	const form = await superValidate(zod4(contactSchema));

	return {
		pageMetaTags,
		user,
		settings,
		lang,
		form
	};
};