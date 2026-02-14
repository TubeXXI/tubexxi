export const load = async ({ locals }) => {
	const { user, settings, lang } = locals;

	return {
		user,
		settings,
		lang
	}
};
