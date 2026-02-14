import { redirect } from '@sveltejs/kit';
import { localizeHref } from '@/paraglide/runtime';

export const load = async ({ locals }) => {
	const { user } = locals;

	if (!user) {
		throw redirect(302, localizeHref('/auth/login'));
	}

	throw redirect(302, localizeHref('/user/account'));
};
