export const load = async ({ data }) => {
	const { user, settings, lang, registerForm, pageMetaTags } = data;

	const theme = settings?.SYSTEM.theme || 'default';

	const component = (await import(`@/components/themes/${theme}/layout/AuthLayout.svelte`)).default;

	return {
		user,
		settings,
		lang,
		registerForm,
		pageMetaTags,
		theme,
		component
	};
};
