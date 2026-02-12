
export const load = async ({ locals, parent }) => {
	const { user, settings, lang } = locals;

	const canonicalUrl = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);
	return {
		user,
		settings,
		lang,
		canonicalUrl,
		alternates,
	};
}
