import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingEmail } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update Email Settings', lang)) as SingleResponse;
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
		smtp_enabled: settings?.EMAIL.smtp_enabled ?? false,
		smtp_service: settings?.EMAIL.smtp_service ?? '',
		smtp_host: settings?.EMAIL.smtp_host ?? '',
		smtp_port: settings?.EMAIL.smtp_port ?? 0,
		smtp_user: settings?.EMAIL.smtp_user ?? '',
		smtp_password: settings?.EMAIL.smtp_password ?? '',
		from_email: settings?.EMAIL.from_email ?? '',
		from_name: settings?.EMAIL.from_name ?? '',
	}, zod4(updateSettingEmail));

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
		const form = await superValidate(request, zod4(updateSettingEmail));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const settingsToUpdate = [
			{ key: 'smtp_enabled', value: String(form.data.smtp_enabled), group_name: 'EMAIL' },
			{ key: 'smtp_service', value: form.data.smtp_service, group_name: 'EMAIL' },
			{ key: 'smtp_host', value: form.data.smtp_host, group_name: 'EMAIL' },
			{ key: 'smtp_port', value: String(form.data.smtp_port), group_name: 'EMAIL' },
			{ key: 'smtp_user', value: form.data.smtp_user, group_name: 'EMAIL' },
			{ key: 'smtp_password', value: form.data.smtp_password, group_name: 'EMAIL' },
			{ key: 'from_email', value: form.data.from_email, group_name: 'EMAIL' },
			{ key: 'from_name', value: form.data.from_name, group_name: 'EMAIL' }
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
