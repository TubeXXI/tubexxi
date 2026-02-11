import { goto } from '$app/navigation';
import { type Locale } from '@/paraglide/runtime.js';

export const LanguageLabels: Partial<Record<Locale, string>> = {
	en: 'English',
	es: 'Español',
	de: 'German',
	pt: 'Português',
	fr: 'Français',
	id: 'Bahasa Indonesia',
	hi: 'हिन्दी',
	ar: 'العربية',
	zh: '中文',
	ru: 'Русский',
	ja: '日本語',
	tr: 'Türkçe',
	vi: 'Tiếng Việt',
	th: 'ไทย',
	el: 'Ελληνικά',
	it: 'Italiano'
};


export function localizePath(path: string, lang: string) {
	return `/${lang}${path.startsWith('/') ? path : `/${path}`}`;
}

export function gotoLocale(path: string, lang: string) {
	const href = `/${lang}${path.startsWith('/') ? path : `/${path}`}`;
	goto(href);
}
