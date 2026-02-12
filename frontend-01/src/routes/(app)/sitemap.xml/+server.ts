import { redirect } from '@sveltejs/kit';
import { localizeHref, locales } from '@/paraglide/runtime.js';

export async function GET({ url, request, locals }) {
	if (request.headers.get('accept')?.includes('application/json') ||
		url.pathname.includes('__data.json')) {
		throw redirect(307, localizeHref('/sitemap.xml', { locale: locals.lang }));
	}

	let origin = url.origin || request.headers.get('origin') || '';
	if (origin.startsWith('http://')) {
		origin = origin.replace('http://', 'https://');
	}

	const sitemaps = locales
		.map(
			(l) =>
				`<sitemap>
					<loc>${origin}/sitemap-${l}.xml</loc>
				</sitemap>`
		)
		.join('');

	return new Response(
		`<?xml version="1.0" encoding="UTF-8"?>
		<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
			${sitemaps}
		</sitemapindex>`,
		{
			headers: {
				'Content-Type': 'application/xml'
			}
		}
	);
}
