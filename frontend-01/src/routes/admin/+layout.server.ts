export const load = async ({ locals }) => {
	const { user, settings } = locals;

	return {
		user,
		settings,
	};
};
