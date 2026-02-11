export default {
	async fetch(req: Request) {
		if (req.method !== 'GET') return fetch(req);

		const url = new URL(req.url);

		if (req.headers.get('accept')?.includes('text/html') === false) {
			return fetch(req);
		}

		const imagesRegex = /\.(png|jpg|jpeg|gif|webp|svg)$/;

		if (url.pathname.startsWith('/_app') ||
			url.pathname.startsWith('/api') ||
			url.pathname.startsWith('/favicon.ico') ||
			imagesRegex.test(url.pathname) ||
			url.pathname.startsWith('/images/') ||
			url.pathname.startsWith('/robots.txt') ||
			url.pathname.includes('.')) {
			return fetch(req);
		}

		if (/^\/(en|id|es)(\/|$)/.test(url.pathname)) {
			return fetch(req);
		}

		if (url.pathname !== '/' && url.pathname !== '') {
			return fetch(req);
		}

		const ua = req.headers.get('user-agent') || '';
		if (/bot|crawl|spider|facebookexternalhit|twitterbot/i.test(ua)) {
			return fetch(req);
		}

		const country = req.headers.get('cf-ipcountry') || 'US';
		const locale = country === 'ID' ? 'id' : 'en';

		return Response.redirect(`${url.origin}/${locale}`, 302);
	}
};
