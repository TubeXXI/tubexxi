import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingMonetization } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update Monetization Settings', lang)) as SingleResponse;
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
		enable_monetize: settings?.MONETIZE.enable_monetize ?? false,
		type_monetize: settings?.MONETIZE.type_monetize ?? 'adsterra',
		publisher_id: settings?.MONETIZE.publisher_id ?? '',
		enable_popup_ad: settings?.MONETIZE.enable_popup_ad ?? false,
		auto_ad_code: settings?.MONETIZE.auto_ad_code ?? '',
		popup_ad_code: settings?.MONETIZE.popup_ad_code ?? '',
		socialbar_ad_code: settings?.MONETIZE.socialbar_ad_code ?? '',
		banner_rectangle_ad_code: settings?.MONETIZE.banner_rectangle_ad_code ?? '',
		banner_horizontal_ad_code: settings?.MONETIZE.banner_horizontal_ad_code ?? '',
		banner_vertical_ad_code: settings?.MONETIZE.banner_vertical_ad_code ?? '',
		native_ad_code: settings?.MONETIZE.native_ad_code ?? '',
		direct_link_ad_code: settings?.MONETIZE.direct_link_ad_code ?? ''
	}, zod4(updateSettingMonetization));

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
		const form = await superValidate(request, zod4(updateSettingMonetization));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const settingsToUpdate = [
			{ key: 'enable_monetize', value: String(form.data.enable_monetize), group_name: 'MONETIZE' },
			{ key: 'type_monetize', value: form.data.type_monetize, group_name: 'MONETIZE' },
			{ key: 'publisher_id', value: form.data.publisher_id, group_name: 'MONETIZE' },
			{ key: 'enable_popup_ad', value: String(form.data.enable_popup_ad), group_name: 'MONETIZE' },
			{ key: 'auto_ad_code', value: form.data.auto_ad_code, group_name: 'MONETIZE' },
			{ key: 'popup_ad_code', value: form.data.popup_ad_code, group_name: 'MONETIZE' },
			{ key: 'socialbar_ad_code', value: form.data.socialbar_ad_code, group_name: 'MONETIZE' },
			{ key: 'banner_rectangle_ad_code', value: form.data.banner_rectangle_ad_code, group_name: 'MONETIZE' },
			{ key: 'banner_horizontal_ad_code', value: form.data.banner_horizontal_ad_code, group_name: 'MONETIZE' },
			{ key: 'banner_vertical_ad_code', value: form.data.banner_vertical_ad_code, group_name: 'MONETIZE' },
			{ key: 'native_ad_code', value: form.data.native_ad_code, group_name: 'MONETIZE' },
			{ key: 'direct_link_ad_code', value: form.data.direct_link_ad_code, group_name: 'MONETIZE' },
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
