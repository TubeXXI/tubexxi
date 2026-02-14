import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { registerAppSchema } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';
import { redirect } from '@sveltejs/kit';
import { localizeHref } from '@/paraglide/runtime.js';
import { generateEncodeBase64 } from '$lib/utils/random.js';

export const load = async ({ locals, parent, url }) => {
	const { user, settings, deps, lang } = locals;

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

	const form = await superValidate(zod4(registerAppSchema));


	return {
		pageMetaTags,
		settings,
		user,
		lang,
		form
	};
};
export const actions = {
	default: async ({ locals, request }) => {
		const { deps } = locals;
		const form = await superValidate(request, zod4(registerAppSchema));
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

		const appsToRegister = [
			// CONFIG group - selalu ada
			{ package_name: form.data.package_name, key: 'name', value: getValue(form.data.name, ''), group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'package_name', value: form.data.package_name, group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'version', value: getValue(form.data.version, '1.0.0'), group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'type', value: getValue(form.data.type, 'android'), group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'is_active', value: form.data.is_active ? 'true' : 'false', group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'store_url', value: getValue(form.data.store_url, ''), group_name: 'CONFIG' },
			{ package_name: form.data.package_name, key: 'api_key', value: generateEncodeBase64(getValue(form.data.package_name, '')), group_name: 'CONFIG' },
			// MONETIZE group - conditional
			...(form.data.enable_monetize ? [
				{ package_name: form.data.package_name, key: 'enable_monetize', value: 'true', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_admob', value: form.data.enable_admob ? 'true' : 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_unity_ad', value: form.data.enable_unity_ad ? 'true' : 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_star_io_ad', value: form.data.enable_star_io_ad ? 'true' : 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_in_app_purchase', value: form.data.enable_in_app_purchase ? 'true' : 'false', group_name: 'MONETIZE' },
				// Admob IDs - hanya jika enable_admob true
				...(form.data.enable_admob ? [
					{ package_name: form.data.package_name, key: 'admob_id', value: getValue(form.data.admob_id, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'admob_auto_ad', value: getValue(form.data.admob_auto_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'admob_banner_ad', value: getValue(form.data.admob_banner_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'admob_interstitial_ad', value: getValue(form.data.admob_interstitial_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'admob_native_ad', value: getValue(form.data.admob_native_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'admob_rewarded_ad', value: getValue(form.data.admob_rewarded_ad, ''), group_name: 'MONETIZE' },
				] : []),
				// Unity IDs - hanya jika enable_unity_ad true
				...(form.data.enable_unity_ad ? [
					{ package_name: form.data.package_name, key: 'unity_ad_id', value: getValue(form.data.unity_ad_id, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'unity_banner_ad', value: getValue(form.data.unity_banner_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'unity_interstitial_ad', value: getValue(form.data.unity_interstitial_ad, ''), group_name: 'MONETIZE' },
					{ package_name: form.data.package_name, key: 'unity_rewarded_ad', value: getValue(form.data.unity_rewarded_ad, ''), group_name: 'MONETIZE' },
				] : []),
				// Star IO IDs - hanya jika enable_star_io_ad true
				...(form.data.enable_star_io_ad ? [
					{ package_name: form.data.package_name, key: 'star_io_ad_id', value: getValue(form.data.star_io_ad_id, ''), group_name: 'MONETIZE' },
				] : []),

				// One Signal ID - selalu ada jika monetize enabled
				{ package_name: form.data.package_name, key: 'one_signal_id', value: getValue(form.data.one_signal_id, ''), group_name: 'MONETIZE' },
			] : [
				// Jika monetize disabled, tetap insert dengan nilai false
				{ package_name: form.data.package_name, key: 'enable_monetize', value: 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_admob', value: 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_unity_ad', value: 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_star_io_ad', value: 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'enable_in_app_purchase', value: 'false', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_id', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'unity_ad_id', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'star_io_ad_id', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_auto_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_banner_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_interstitial_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_native_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'admob_rewarded_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'unity_banner_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'unity_interstitial_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'unity_rewarded_ad', value: '', group_name: 'MONETIZE' },
				{ package_name: form.data.package_name, key: 'one_signal_id', value: '', group_name: 'MONETIZE' },
			])
		];

		// console.log('Apps to register:', appsToRegister.length);

		const success = await deps.applicationService.RegisterApplication(appsToRegister);
		if (!success) {
			return fail(400, {
				form,
				message: 'Failed to register application'
			});
		}
		redirect(303, localizeHref('/admin/applications'));
	}
};
