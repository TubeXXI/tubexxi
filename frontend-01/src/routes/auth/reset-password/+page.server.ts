import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from "@siamf/google-translate";
import { capitalizeFirstLetter } from "@/utils/format.js";
import { superValidate } from 'sveltekit-superforms';
import { resetPasswordSchema } from '$lib/utils/schema';
import { fail } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import * as i18n from '@/paraglide/messages.js';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = await deps.languageHelper.singleTranslate('Reset Password', lang) as SingleResponse;
	const siteName = await deps.languageHelper.singleTranslate(settings?.WEBSITE?.site_name || '', lang) as SingleResponse;
	const tagline = await deps.languageHelper.singleTranslate(settings?.WEBSITE?.site_tagline || '', lang) as SingleResponse;
	const description = await deps.languageHelper.singleTranslate(settings?.WEBSITE?.site_description || '', lang) as SingleResponse;
	const keywords = await Promise.all((settings?.WEBSITE?.site_keywords?.split(',') || ['']).map(async (keyword) => await deps.languageHelper.singleTranslate(keyword.trim(), lang) as SingleResponse));

	const pageMetaTags = defaultMetaTags({
		path_url: defaultOrigin,
		title: `${capitalizeFirstLetter(title.data.target.text || '')} - ${capitalizeFirstLetter(siteName.data.target.text || '')}`,
		tagline: capitalizeFirstLetter(tagline.data.target.text || ''),
		description: capitalizeFirstLetter(description.data.target.text || ''),
		keywords: keywords.map((keyword: SingleResponse) => capitalizeFirstLetter(keyword.data.target.text || '')),
		robots: 'index, follow',
		canonical: defaultOrigin,
		alternates,
		graph_type: 'website'
	}, settings);

	const resetForm = await superValidate(zod4(resetPasswordSchema));

	return {
		pageMetaTags,
		resetForm,
		settings,
		user,
		lang
	};
}
export const actions = {
	default: async ({ locals, request }) => {
		const { deps } = locals;

		const form = await superValidate(request, zod4(resetPasswordSchema));
		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: i18n.invalid_input(),
				error: form.errors
			});
		}

		try {

			const response = await deps.authService.ResetPassword(form.data) as string | Error
			if (response instanceof Error) {
				return fail(500, {
					form,
					success: false,
					message: response.message || i18n.page_forgot_password_error_message(),
					error: response
				});
			}

			return {
				form,
				success: true,
				message: i18n.page_forgot_password_success_message()
			}

		} catch (error) {
			console.error('Error resetting password request:', error);
			return fail(500, {
				form,
				success: false,
				message: error instanceof Error ? error.message : i18n.page_forgot_password_error_message(),
				error: null
			});
		}
	}
}
