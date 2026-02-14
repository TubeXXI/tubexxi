export const load = async ({ data }) => {
	const { user, settings, lang, pageMetaTags, form } = data;

	const theme = settings?.SYSTEM.theme || 'default';

	const component = (await import(`@/components/themes/${theme}/pages/DefaultContact.svelte`)).default;

	return {
		user,
		settings,
		lang,
		pageMetaTags,
		theme,
		component,
		form
	};
};
