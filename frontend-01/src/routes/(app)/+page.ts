export const load = async ({ data }) => {
	const { user, settings, lang, movies, pageMetaTags } = data;

	const theme = settings?.SYSTEM.theme || 'default';

	const component = (await import(`@/components/themes/${theme}/layout/MainLayout.svelte`)).default;

	return {
		user,
		settings,
		lang,
		movies,
		pageMetaTags,
		theme,
		component
	};
};
