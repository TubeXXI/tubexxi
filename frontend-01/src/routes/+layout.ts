

export const load = async ({ data }) => {
	const { user, settings, lang, canonicalUrl, alternates } = data;

	return {
		user,
		settings,
		lang,
		canonicalUrl,
		alternates
	};
};
