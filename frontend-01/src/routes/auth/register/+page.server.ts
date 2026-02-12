import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from "@siamf/google-translate";
import { capitalizeFirstLetter } from "@/utils/format.js";
import { superValidate } from 'sveltekit-superforms';
import { registerSchema } from '$lib/utils/schema';
import { fail } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import * as i18n from '@/paraglide/messages.js';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = await deps.languageHelper.singleTranslate('Register', lang) as SingleResponse;
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

	const registerForm = await superValidate(zod4(registerSchema));

	return {
		pageMetaTags,
		registerForm,
		settings,
		user,
		lang
	};
}
export const actions = {
	default: async ({ request, locals }) => {
		const { deps, session } = locals;

		const formData = await request.formData();

		const idToken = formData.get('idToken') as string;
		const email = formData.get('email') as string;
		const password = formData.get('password') as string;
		const full_name = formData.get('full_name') as string;
		const phone = formData.get('phone') as string;
		const avatar_url = formData.get('avatar_url') as string;

		const form = await superValidate(formData, zod4(registerSchema));

		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: i18n.invalid_input(),
				error: form.errors
			});
		}

		try {
			const registerResponse = await deps.authService.Register({
				idToken,
				email,
				password,
				full_name,
				phone,
				avatar_url
			}) as FirebaseAuthResponse | Error;

			if (registerResponse instanceof Error) {
				return fail(400, {
					form,
					success: false,
					message: registerResponse.message || i18n.page_sign_up_error_message(),
					error: registerResponse
				});
			}

			session?.set('access_token', idToken, 60 * 60 * 24);

			locals.user = registerResponse.user;

			return {
				form,
				success: true,
				message: i18n.page_sign_in_success_message()
			}

		} catch (error) {
			console.error('Error registering user:', error);
			return fail(500, {
				form,
				success: false,
				message: error instanceof Error ? error.message : i18n.page_sign_up_error_message(),
				error: null
			});
		}
	}
}
