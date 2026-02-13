import { redirect } from '@sveltejs/kit';
import { localizeHref, type Locale } from '@/paraglide/runtime.js';

let pages: string[] = ['about', 'contact', 'faq', 'privacy', 'terms', 'movies', 'series', 'anime'];

const sitemap = (site: string, pages: string[]) => `<?xml version="1.0" encoding="UTF-8" ?>
<urlset
	xmlns="https://www.sitemaps.org/schemas/sitemap/0.9"
	xmlns:news="https://www.google.com/schemas/sitemap-news/0.9"
	xmlns:xhtml="https://www.w3.org/1999/xhtml"
	xmlns:mobile="https://www.google.com/schemas/sitemap-mobile/1.0"
	xmlns:image="https://www.google.com/schemas/sitemap-image/1.1"
	xmlns:video="https://www.google.com/schemas/sitemap-video/1.1"
>
	${pages
		.map(
			(page) => `
	<url>
		<loc>${site}/${page}</loc>
		<changefreq>daily</changefreq>
		<priority>0.5</priority>
	</url>
	`
		)
		.join('')}
</urlset>`;

export async function GET({ url, request, locals, params }) {
	if (
		request.headers.get('accept')?.includes('application/json') ||
		url.pathname.includes('__data.json')
	) {
		throw redirect(307, localizeHref('/sitemap.xml', { locale: locals.lang }));
	}

	const locale = params.locale as Locale;
	let origin = url.origin || request.headers.get('origin') || '';
	if (origin.startsWith('http://')) {
		origin = origin.replace('http://', 'https://');
	}

	const body = sitemap(`${origin}/${locale}`, pages);
	const response = new Response(body);
	response.headers.set('Cache-Control', 'max-age=0, s-maxage=3600');
	response.headers.set('Content-Type', 'application/xml');
	return response;
}
