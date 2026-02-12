import { redirect } from '@sveltejs/kit';
import { localizeHref } from '@/paraglide/runtime';

export const load = async ({ locals }) => {
	const { user } = locals;

	if (!user) {
		throw redirect(302, localizeHref('/auth/login'));
	}
	if (user.role?.level !== 1 && user.role?.level !== 2) {
		throw redirect(302, localizeHref('/user/account'));
	}

	throw redirect(302, localizeHref('/admin/dashboard'));
};
