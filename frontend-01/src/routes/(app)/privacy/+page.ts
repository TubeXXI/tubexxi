export const load = async ({ data }) => {
	const { user, settings, lang, pageMetaTags } = data;

	const theme = settings?.SYSTEM.theme || 'default';

	const component = (await import(`@/components/themes/${theme}/pages/DefaultPrivacy.svelte`)).default;

	return {
		user,
		settings,
		lang,
		pageMetaTags,
		theme,
		component
	};
};
