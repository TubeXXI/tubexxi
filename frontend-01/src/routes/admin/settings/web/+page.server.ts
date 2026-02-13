import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingWeb } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update Website Settings', lang)) as SingleResponse;
	const siteName = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_name || '',
		lang
	)) as SingleResponse;
	const tagline = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_tagline || '',
		lang
	)) as SingleResponse;
	const description = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_description || '',
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
			title: `${capitalizeFirstLetter(title.data.target.text || '')} - ${capitalizeFirstLetter(siteName.data.target.text || '')}`,
			tagline: capitalizeFirstLetter(tagline.data.target.text || ''),
			description: capitalizeFirstLetter(description.data.target.text || ''),
			keywords: keywords.map((keyword: SingleResponse) =>
				capitalizeFirstLetter(keyword.data.target.text || '')
			),
			robots: 'noindex, nofollow',
			canonical: defaultOrigin,
			alternates,
			graph_type: 'website'
		},
		settings
	);

	const form = await superValidate({
		site_name: settings?.WEBSITE?.site_name || '',
		site_tagline: settings?.WEBSITE?.site_tagline || '',
		site_description: settings?.WEBSITE?.site_description || '',
		site_keywords: settings?.WEBSITE?.site_keywords || '',
		site_email: settings?.WEBSITE?.site_email || '',
		site_phone: settings?.WEBSITE?.site_phone || '',
		site_url: defaultOrigin,
	}, zod4(updateSettingWeb));

	return {
		pageMetaTags,
		form,
		settings,
		user,
		lang
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(updateSettingWeb));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const settingsToUpdate = [
			{ key: 'site_name', value: form.data.site_name, group_name: 'WEBSITE' },
			{ key: 'site_tagline', value: form.data.site_tagline, group_name: 'WEBSITE' },
			{ key: 'site_description', value: form.data.site_description, group_name: 'WEBSITE' },
			{ key: 'site_keywords', value: form.data.site_keywords, group_name: 'WEBSITE' },
			{ key: 'site_email', value: form.data.site_email, group_name: 'WEBSITE' },
			{ key: 'site_phone', value: form.data.site_phone, group_name: 'WEBSITE' },
			{ key: 'site_url', value: form.data.site_url, group_name: 'WEBSITE' }
		].filter(s => s.value !== undefined) as { key: string; value: string; group_name: string }[];

		const updateResponse = await locals.deps.settingService.updateBulkSetting(settingsToUpdate);
		if (updateResponse instanceof Error) {
			return fail(500, {
				form,
				message: updateResponse.message || 'Failed to update settings.'
			});
		}

		return {
			form,
			message: 'Settings updated successfully.'
		};
	}
};
