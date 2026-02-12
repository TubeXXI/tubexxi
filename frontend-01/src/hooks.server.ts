import { redirect, type RequestEvent, type Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { env } from '$env/dynamic/private';
import { Dependencies } from '@/server/deps';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { locales as SUPPORTED_LOCALES, type Locale } from '@/paraglide/runtime';
import { ApiClientHandler } from '$lib/helpers/api_helpers';

const NODE_ENV = env.NODE_ENV || 'development';

const paraglideHandleWithAutoDetectedLocale: Handle = ({ event, resolve }) => {
	const { request } = event;
	const pathname = event.url.pathname;

	if (
		pathname.startsWith('/api') ||
		pathname.startsWith('/_app') ||
		pathname.includes('.')
	) {
		return resolve(event);
	}

	const ua = request.headers.get('user-agent');
	const isBot = !!ua && /bot|crawl|spider|facebookexternalhit|twitterbot/i.test(ua);
	const pathLocale = getLocaleFromPath(pathname);
	const detectedLocale = detectLocale(event);

	const resolveWithParaglide = (locale: Locale) => {
		return paraglideMiddleware(event.request, ({ request: localizedRequest }) => {
			event.request = localizedRequest;
			try {
				(event as any).url = new URL(localizedRequest.url);
			} catch {
				try {
					Object.defineProperty(event, 'url', { value: new URL(localizedRequest.url) });
				} catch {
					// ignore
				}
			}
			return resolve(event, {
				transformPageChunk: ({ html }) => {
					return html.replace('%lang%', event.locals.lang);
				}
			});
		});
	};

	if (request.method !== 'GET' && request.method !== 'HEAD') {
		event.locals.lang = (pathLocale ?? detectedLocale) as Locale;
		setCookie(event, event.locals.lang as Locale);
		return resolveWithParaglide(event.locals.lang as Locale);
	}

	if (isBot) {
		event.locals.lang = (pathLocale ?? 'en') as Locale;
		setCookie(event, event.locals.lang as Locale);
		return resolveWithParaglide(event.locals.lang as Locale);
	}

	if (pathname === '/en' || pathname.startsWith('/en/')) {
		const stripped = pathname === '/en' ? '/' : pathname.slice(3);
		throw redirect(302, `${stripped}${event.url.search}`);
	}

	if (!pathLocale) {
		event.locals.lang = detectedLocale;
		setCookie(event, event.locals.lang as Locale);
		if (event.locals.lang !== 'en') {
			throw redirect(302, `/${event.locals.lang}${pathname === '/' ? '' : pathname}${event.url.search}`);
		}
		return resolveWithParaglide(event.locals.lang as Locale);
	}

	event.locals.lang = pathLocale;
	setCookie(event, pathLocale);
	return resolveWithParaglide(event.locals.lang as Locale);
};

const dependenciesInject: Handle = async ({ event, resolve }) => {
	event.locals.deps = new Dependencies(event);
	event.locals.session = {
		set: (key, value, maxAge?: number) => {
			const isProduction = NODE_ENV === 'production';

			event.cookies.set(key, value, {
				httpOnly: true,
				secure: isProduction,
				sameSite: 'lax',
				path: '/',
				maxAge: maxAge ?? 60 * 60 * 24 * 7 // default 1 week
			});
		},
		get: (key) => {
			const val = event.cookies.get(key);
			return val ?? null;
		},
		delete: (key) => {
			const isProduction = NODE_ENV === 'production';
			event.cookies.delete(key, {
				path: '/',
				httpOnly: true,
				secure: isProduction,
				sameSite: 'lax'
			});
		}
	};
	event.locals.safeGetUser = async () => {
		try {
			const accessToken = event.locals.session?.get('access_token');

			if (!accessToken) {
				return null;
			}

			const user = await event.locals.deps.userService.CurrentUser();
			if (user instanceof Error) {
				return null;
			}
			return user;
		} catch (e) {
			console.error("Error fetching user:", e);
			return null;
		}
	};
	event.locals.safeGetSettings = async () => {
		try {
			const settings = await event.locals.deps.settingService.getPublicSettings();
			if (settings instanceof Error) {
				return null;
			}
			return settings;
		} catch (e) {
			console.error("Error fetching settings:", e);
			return null;
		}
	};

	const response = await resolve(event);
	return response;
};

const authHandle: Handle = async ({ event, resolve }) => {
	const { locals } = event;

	try {
		if (!locals.settings) {
			locals.settings = await event.locals.safeGetSettings();
		}
	} catch (error: any) {
		console.error('Error fetching settings:', error);
	}

	try {
		if (!locals.user) {
			locals.user = await event.locals.safeGetUser();
		}

		return resolve(event);
	} catch (error: any) {
		console.error('Auth Handle Error:', error);
		if (error?.status === 302 || error?.status === 301 || error?.status === 303) {
			throw error;
		}
		return resolve(event);
	}
}

const adminMiddleware: Handle = async ({ event, resolve }) => {
	const { url, locals } = event;


	const isProtected = url.pathname.startsWith('/admin');

	if (isProtected) {
		if (!locals.user) {
			throw redirect(303, `/auth/login?redirect=${encodeURIComponent(url.pathname)}`);
		}

		const user = locals.user;
		const isAdmin = user.role?.name === 'admin' || user.role?.level === 1 || user.role?.level === 2;

		if (!isAdmin) {
			throw redirect(303, '/user/account');
		}
	}

	const isAuthRoute = url.pathname.startsWith("/auth");

	if (isAuthRoute) {
		if (locals.user && locals.user.role?.name === 'admin') {
			throw redirect(303, '/admin/dashboard');
		} else if (locals.user && locals.user.role?.name === 'user') {
			throw redirect(303, '/user/account');
		}
	}

	return resolve(event);
};

const errorHandling: Handle = async ({ event, resolve }) => {
	try {
		const response = await resolve(event);

		if (response.status === 404) {
			const locale = getLocaleFromPath(event.url.pathname);
			let redirectPath = '/';
			if (locale) {
				redirectPath = `/${locale}`;
			}


			const isStaticAsset = /\.(css|js|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$/i.test(event.url.pathname);
			const isApiRoute = event.url.pathname.startsWith('/api/');
			const isImagesRoute = event.url.pathname.startsWith('/images/');
			const isInternalRoute = event.url.pathname.startsWith('/_') || event.url.pathname.includes('__');

			if (!isStaticAsset && !isApiRoute && !isInternalRoute && !isImagesRoute) {
				return new Response(null, {
					status: 302,
					headers: {
						'Location': redirectPath,
						'Cache-Control': 'no-cache'
					}
				});
			}
		}



		if (event.url.pathname.startsWith('/auth/')) {
			response.headers.set('Cross-Origin-Opener-Policy', 'unsafe-none');
		}
		return response;
	} catch (error) {
		console.error('Server error:', error);

		if (
			(error as any)?.status === 301 ||
			(error as any)?.status === 302 ||
			(error as any)?.status === 303 ||
			(error as any)?.status === 307 ||
			(error as any)?.status === 308
		) {
			throw error;
		}

		if (event.url.pathname.startsWith('/api/')) {
			return new Response(
				JSON.stringify({
					error: "Internal Server Error",
					message: error instanceof Error ? error.message : "An unknown error occurred"
				}),
				{
					status: 500,
					headers: { 'Content-Type': 'application/json' }
				}
			);
		}

		throw error;
	}
};

export const handle: Handle = sequence(
	paraglideHandleWithAutoDetectedLocale,
	dependenciesInject,
	authHandle,
	adminMiddleware,
	errorHandling
);

let lastSsrReportAt = 0;
let lastSsrKey = '';

export const handleError = async ({ error, event, status, message }) => {
	try {
		if (!status || status < 500) {
			return;
		}
		const now = Date.now();
		const url = event.url?.toString?.() ? event.url.toString() : '';
		const routeId = (event.route as any)?.id ? String((event.route as any).id) : '';
		const pathname = event.url?.pathname ?? '';
		const search = event.url?.search ?? '';
		const errMsg = error instanceof Error ? error.message : String(message ?? error);
		const stack = error instanceof Error ? error.stack ?? '' : '';
		const key = `${status}|${routeId}|${pathname}|${errMsg}`.slice(0, 500);
		if (key === lastSsrKey && now - lastSsrReportAt < 15_000) {
			return;
		}
		lastSsrKey = key;
		lastSsrReportAt = now;

		const api = new ApiClientHandler(event);
		const ip = (() => {
			try {
				return event.getClientAddress?.() ?? '';
			} catch {
				return '';
			}
		})();
		const locale = (event.locals.lang as string | undefined) ?? '';
		const userId = (event.locals as any)?.user?.id ? String((event.locals as any).user.id) : '';
		const reqMeta = {
			route_id: routeId,
			pathname,
			search,
			params: (event.params ? JSON.stringify(event.params).slice(0, 500) : ''),
			referer: event.request.headers.get('referer') ?? '',
			x_request_id: event.request.headers.get('x-request-id') ?? '',
			cf_ray: event.request.headers.get('cf-ray') ?? ''
		};
		// await api.publicRequest('POST', '/web-client/report/errors', {
		// 	error: stack ? stack.slice(0, 1200) : String(error).slice(0, 1200),
		// 	message: errMsg.slice(0, 500),
		// 	platform_id: 'frontend-ssr',
		// 	user_id: userId,
		// 	ip_address: ip,
		// 	user_agent: event.request.headers.get('user-agent') ?? '',
		// 	url,
		// 	method: event.request.method,
		// 	request: JSON.stringify(reqMeta).slice(0, 800),
		// 	status,
		// 	level: 'error',
		// 	locale,
		// 	timestamp_ms: now
		// }, true);
	} catch {
		return;
	}
};

function hasLocalePrefix(path: string): boolean {
	return SUPPORTED_LOCALES.some(
		(l) => path === `/${l}` || path.startsWith(`/${l}/`)
	);
}

function getLocaleFromPath(pathname: string): Locale | null {
	const match = pathname.match(/^\/(en|id|es|ru|pt|fr|de|zh|hi|ar|ja|tr|vi|th|el|it)(\/|$)/);
	return match ? (match[1] as Locale) : null;
}

function detectLocale(event: RequestEvent): Locale {
	// const cookie = event.cookies.get('PARAGLIDE_LOCALE') as Locale | null;
	// if (cookie && SUPPORTED_LOCALES.includes(cookie)) return cookie;

	const accept = event.request.headers.get('accept-language');
	const l = accept?.split(',')[0].split('-')[0] as Locale;
	const supported = SUPPORTED_LOCALES.includes(l);

	return supported ? l : 'en';
}

function setCookie(event: RequestEvent, locale: Locale) {
	event.cookies.set('PARAGLIDE_LOCALE', locale, {
		httpOnly: true,
		secure: NODE_ENV === 'production',
		sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7
	});
}
