import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { loginSchema } from '$lib/utils/schema';
import { fail } from '@sveltejs/kit';
import { zod4 } from 'sveltekit-superforms/adapters';
import { localizeHref } from '@/paraglide/runtime';
import * as i18n from '@/paraglide/messages.js';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Login', lang)) as SingleResponse;
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
			robots: 'index, follow',
			canonical: defaultOrigin,
			alternates,
			graph_type: 'website'
		},
		settings
	);

	const loginForm = await superValidate(zod4(loginSchema));

	return {
		pageMetaTags,
		loginForm,
		settings,
		user,
		lang
	};
};

export const actions = {
	default: async ({ locals, request }) => {
		const { deps, session } = locals;

		const formData = await request.formData();

		const idToken = formData.get('idToken') as string;
		const rememberMe = formData.get('rememberMe') === 'true';

		const form = await superValidate(formData, zod4(loginSchema));

		if (!form.valid) {
			return fail(400, {
				form,
				success: false,
				message: i18n.invalid_input(),
				error: form.errors
			});
		}

		const loginResponse = (await deps.authService.Login(idToken)) as FirebaseAuthResponse | Error;

		if (loginResponse instanceof Error) {
			return fail(400, {
				form,
				success: false,
				message: loginResponse.message || i18n.page_sign_in_error_message(),
				error: loginResponse
			});
		}
		const maxAge = rememberMe ? 60 * 60 * 24 * 7 : 60 * 60 * 24;
		session?.set('access_token', idToken, maxAge);

		locals.user = loginResponse.user;

		const roleName = loginResponse.user?.role?.name || '';
		const roleLevel = loginResponse.user?.role?.level || 0;
		const isAdmin = roleName === 'admin' || roleLevel === 1 || roleLevel === 2;
		const destination = isAdmin ? '/admin/dashboard' : '/user/account';

		return {
			success: true,
			message: 'Login success',
			redirect_url: localizeHref(destination)
		};
	}
};
