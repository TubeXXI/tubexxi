import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateApplicationSchema, type UpdateApplicationSchema } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';
import { redirect } from '@sveltejs/kit';
import { localizeHref } from '@/paraglide/runtime.js';

export const load = async ({ locals, parent, params }) => {
	const { user, settings, deps, lang } = locals;

	const { id } = params;
	if (!id) {
		throw redirect(302, localizeHref('/admin/applications'));
	}

	const application = await deps.applicationService.GetByPackageName(id) as ApplicationResponse | null;
	if (!application) {
		throw redirect(302, localizeHref('/admin/applications'));
	}

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Create Application', lang)) as SingleResponse;
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

	const devaultValue = {
		package_name: application.config.package_name,
		name: application.config.name,
		version: application.config.version,
		type: application.config.type || 'android',
		is_active: application.config.is_active,
		store_url: application.config.store_url,
		enable_monetize: application.monetize.enable_monetize,
		enable_admob: application.monetize.enable_admob,
		enable_unity_ad: application.monetize.enable_unity_ad,
		enable_star_io_ad: application.monetize.enable_star_io_ad,
		enable_in_app_purchase: application.monetize.enable_in_app_purchase,
		admob_id: application.monetize.admob_id,
		unity_ad_id: application.monetize.unity_ad_id,
		star_io_ad_id: application.monetize.star_io_ad_id,
		admob_auto_ad: application.monetize.admob_auto_ad,
		admob_banner_ad: application.monetize.admob_banner_ad,
		admob_interstitial_ad: application.monetize.admob_interstitial_ad,
		admob_native_ad: application.monetize.admob_native_ad,
		admob_rewarded_ad: application.monetize.admob_rewarded_ad,
		unity_banner_ad: application.monetize.unity_banner_ad,
		unity_interstitial_ad: application.monetize.unity_interstitial_ad,
		unity_rewarded_ad: application.monetize.unity_rewarded_ad,
		one_signal_id: application.monetize.one_signal_id,
	} as UpdateApplicationSchema;

	const form = await superValidate(devaultValue, zod4(updateApplicationSchema));


	return {
		pageMetaTags,
		settings,
		user,
		lang,
		form,
		application
	};
};
export const actions = {
	default: async ({ locals, request }) => {
		const { deps } = locals;
		const form = await superValidate(request, zod4(updateApplicationSchema));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', ')
			});
		}

		const getValue = (value: string | undefined | null, defaultValue = '') => {
			if (value === undefined || value === null) {
				return defaultValue;
			}
			if (typeof value === 'string' && value.trim() === '') {
				return defaultValue;
			}
			return value;
		};

		const appsToUpdate = [
			{ key: 'name', value: getValue(form.data.name, ''), group_name: 'CONFIG' },
			{ key: 'package_name', value: getValue(form.data.package_name, ''), group_name: 'CONFIG' },
			{ key: 'version', value: getValue(form.data.version, '1.0.0'), group_name: 'CONFIG' },
			{ key: 'type', value: getValue(form.data.type, 'android'), group_name: 'CONFIG' },
			{ key: 'is_active', value: form.data.is_active ? 'true' : 'false', group_name: 'CONFIG' },
			{ key: 'store_url', value: getValue(form.data.store_url, ''), group_name: 'CONFIG' },
			{ key: 'enable_monetize', value: form.data.enable_monetize ? 'true' : 'false', group_name: 'MONETIZE' },
			{ key: 'enable_admob', value: form.data.enable_admob ? 'true' : 'false', group_name: 'MONETIZE' },
			{ key: 'enable_unity_ad', value: form.data.enable_unity_ad ? 'true' : 'false', group_name: 'MONETIZE' },
			{ key: 'enable_star_io_ad', value: form.data.enable_star_io_ad ? 'true' : 'false', group_name: 'MONETIZE' },
			{ key: 'enable_in_app_purchase', value: form.data.enable_in_app_purchase ? 'true' : 'false', group_name: 'MONETIZE' },
			{ key: 'admob_id', value: getValue(form.data.admob_id, ''), group_name: 'MONETIZE' },
			{ key: 'admob_auto_ad', value: getValue(form.data.admob_auto_ad, ''), group_name: 'MONETIZE' },
			{ key: 'admob_banner_ad', value: getValue(form.data.admob_banner_ad, ''), group_name: 'MONETIZE' },
			{ key: 'admob_interstitial_ad', value: getValue(form.data.admob_interstitial_ad, ''), group_name: 'MONETIZE' },
			{ key: 'admob_native_ad', value: getValue(form.data.admob_native_ad, ''), group_name: 'MONETIZE' },
			{ key: 'admob_rewarded_ad', value: getValue(form.data.admob_rewarded_ad, ''), group_name: 'MONETIZE' },
			{ key: 'unity_ad_id', value: getValue(form.data.unity_ad_id, ''), group_name: 'MONETIZE' },
			{ key: 'unity_banner_ad', value: getValue(form.data.unity_banner_ad, ''), group_name: 'MONETIZE' },
			{ key: 'unity_interstitial_ad', value: getValue(form.data.unity_interstitial_ad, ''), group_name: 'MONETIZE' },
			{ key: 'unity_rewarded_ad', value: getValue(form.data.unity_rewarded_ad, ''), group_name: 'MONETIZE' },
			{ key: 'star_io_ad_id', value: getValue(form.data.star_io_ad_id, ''), group_name: 'MONETIZE' },
			{ key: 'one_signal_id', value: getValue(form.data.one_signal_id, ''), group_name: 'MONETIZE' },
		].filter(s => s.value !== undefined) as { key: string; value: string; group_name: string }[];

		const success = await deps.applicationService.UpdateApplication(form.data.package_name, appsToUpdate);
		if (!success) {
			return fail(400, {
				form,
				message: 'Failed to update application'
			});
		}
		redirect(303, localizeHref('/admin/applications'));
	}
};
