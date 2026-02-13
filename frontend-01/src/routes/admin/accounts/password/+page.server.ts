import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updatePasswordSchema } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update Password', lang)) as SingleResponse;
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

	const passwordForm = await superValidate(zod4(updatePasswordSchema));

	return {
		pageMetaTags,
		passwordForm,
		settings,
		user,
		lang
	};
};

export const actions = {
	default: async ({ request, locals }) => {
		const { deps } = locals;
		const passwordForm = await superValidate(request, zod4(updatePasswordSchema));
		if (!passwordForm.valid) {
			return fail(400, {
				passwordForm,
				message: Object.values(passwordForm.errors).flat().join(', ')
			});
		}
		try {

			const response = await deps.userService.UpdatePassword(passwordForm.data);
			if (response instanceof Error) {
				return fail(500, {
					passwordForm,
					message: response.message
				});
			}

			await deps.userService.Logout();

			return {
				passwordForm,
				success: true,
				message: 'Password updated successfully'
			};

		} catch (error) {
			return fail(500, {
				passwordForm,
				message: error instanceof Error ? error.message : 'Internal server error'
			});
		}
	}
}
