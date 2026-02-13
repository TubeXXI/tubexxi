import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingSystem } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update System Settings', lang)) as SingleResponse;
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
		theme: settings?.SYSTEM.theme ?? 'default',
		enable_documentation: settings?.SYSTEM.enable_documentation ?? true,
		maintenance_mode: settings?.SYSTEM.maintenance_mode ?? false,
		maintenance_message: settings?.SYSTEM.maintenance_message ?? '',
		source_logo_favicon: settings?.SYSTEM.source_logo_favicon ?? 'local',
		histats_tracking_code: settings?.SYSTEM.histats_tracking_code ?? '',
		google_analytics_code: settings?.SYSTEM.google_analytics_code ?? '',
		play_store_app_url: settings?.SYSTEM.play_store_app_url ?? '',
		app_store_app_url: settings?.SYSTEM.app_store_app_url ?? ''
	}, zod4(updateSettingSystem));

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
		const form = await superValidate(request, zod4(updateSettingSystem));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const settingsToUpdate = [
			{ key: 'theme', value: form.data.theme, group_name: 'SYSTEM' },
			{ key: 'enable_documentation', value: String(form.data.enable_documentation), group_name: 'SYSTEM' },
			{ key: 'maintenance_mode', value: String(form.data.maintenance_mode), group_name: 'SYSTEM' },
			{ key: 'maintenance_message', value: form.data.maintenance_message, group_name: 'SYSTEM' },
			{ key: 'source_logo_favicon', value: form.data.source_logo_favicon, group_name: 'SYSTEM' },
			{ key: 'histats_tracking_code', value: form.data.histats_tracking_code, group_name: 'SYSTEM' },
			{ key: 'google_analytics_code', value: form.data.google_analytics_code, group_name: 'SYSTEM' },
			{ key: 'play_store_app_url', value: form.data.play_store_app_url, group_name: 'SYSTEM' },
			{ key: 'app_store_app_url', value: form.data.app_store_app_url, group_name: 'SYSTEM' }
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
