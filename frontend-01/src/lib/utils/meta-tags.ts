import type { MetaTagsProps, MetaTag, LinkTag } from 'svelte-meta-tags';
import { defaultSettings } from '@/constants';

export const defaultMetaTags = (
	options?: PageMetaProps,
	setting?: SettingsValue | null
): MetaTagsProps => ({
	title: `${options?.title || setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name} ${options?.use_tagline === true
		? `- ${options?.tagline || setting?.WEBSITE.site_tagline || defaultSettings.WEBSITE.site_tagline}`
		: ''
		}`,
	description:
		options?.description ||
		setting?.WEBSITE.site_description ||
		defaultSettings.WEBSITE.site_description,
	keywords:
		options?.keywords ||
		(setting?.WEBSITE?.site_keywords || defaultSettings.WEBSITE.site_keywords)
			?.split(', ')
			.map((k) => k.trim()),
	robots: options?.robots || 'index, follow',
	twitter: {
		cardType: 'summary_large_image',
		site: '@idtubexxi',
		image:
			setting?.SYSTEM.source_logo_favicon === 'remote'
				? (isValidUrl(setting?.WEBSITE.site_logo || '') ? setting?.WEBSITE.site_logo || '/images/cover.png' : '/images/cover.png')
				: '/images/cover.png',
		title: options?.title || setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name
	},
	additionalMetaTags: [
		{
			name: 'viewport',
			content: 'width=device-width, initial-scale=1.0'
		},
		{
			property: 'dc:creator',
			content: options?.title || setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name
		},
		{
			name: 'application-name',
			content: options?.title || setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name
		},
		{
			httpEquiv: 'x-ua-compatible',
			content: 'IE=edge'
		},
		{
			name: 'description',
			content:
				options?.description ||
				setting?.WEBSITE.site_description ||
				defaultSettings.WEBSITE.site_description
		},
		{
			name: 'mobile-web-app-capable',
			content: 'yes'
		},
		{
			name: 'mobile-web-app-status-bar-style',
			content: 'black-translucent'
		},
		{
			name: 'mobile-web-app-title',
			content: options?.title || setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name
		},
		{
			name: 'mobile-web-app-icon',
			content: '/mobile-web-app-icon.png'
		}
	] as MetaTag[],
	additionalLinkTags: [
		{
			rel: 'canonical',
			href: options?.canonical || ''
		},
		{
			rel: 'alternate',
			hreflang: 'x-default',
			href: options?.canonical || ''
		},
		...(options?.alternates || []).map((alt) => ({
			rel: 'alternate',
			hreflang: alt.lang,
			href: alt.href
		})),
		{
			rel: 'alternate',
			hreflang: 'application/rss+xml',
			href: `/rss-${options?.language || 'en'}.xml`
		},
		{
			rel: 'icon',
			type: 'image/x-icon',
			sizes: '96x96',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/images/logo.png' : '/images/logo.png')
					: '/images/icon.png'
		},
		{
			rel: 'icon',
			type: 'image/png',
			sizes: '32x32',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/favicon-32x32.png' : '/favicon-32x32.png')
					: '/images/icon.png'
		},
		{
			rel: 'icon',
			type: 'image/png',
			sizes: '16x16',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/favicon-16x16.png' : '/favicon-16x16.png')
					: '/images/icon.png'
		},
		{
			rel: 'icon',
			type: 'image/png',
			sizes: '192x192',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/favicon-192x192.png' : '/favicon-192x192.png')
					: '/images/icon.png'
		},
		{
			rel: 'icon',
			type: 'image/png',
			sizes: '512x512',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/favicon-512x512.png' : '/favicon-512x512.png')
					: '/images/icon.png'
		},
		{
			rel: 'apple-touch-icon',
			type: 'image/png',
			sizes: '180x180',
			href:
				setting?.SYSTEM.source_logo_favicon === 'remote'
					? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/apple-touch-icon.png' : '/apple-touch-icon.png')
					: '/images/icon.png'
		},
	] as LinkTag[],
	openGraph: {
		type: options?.graph_type || 'website',
		url: options?.canonical || '',
		title: options?.title || '',
		description: options?.description || '',
		locale: 'en_IE',
		siteName: setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name,
		images: [
			{
				url:
					setting?.SYSTEM.source_logo_favicon === 'remote'
						? (isValidUrl(setting?.WEBSITE.site_logo || '') ? setting?.WEBSITE.site_logo || '/images/icon.png' : '/images/icon.png')
						: '/images/icon.png',
				width: 800,
				height: 600,
				alt: 'Tube XXI Cover Image',
				type: 'image/png'
			},
			{
				url:
					setting?.SYSTEM.source_logo_favicon === 'remote'
						? (isValidUrl(setting?.WEBSITE.site_favicon || '') ? setting?.WEBSITE.site_favicon || '/images/logo.png' : '/images/logo.png')
						: '/images/icon.png',
				width: 512,
				height: 512,
				alt: 'Tube XXI Android Chrome Icon',
				type: 'image/x-icon'
			}
		],
		profile: {
			firstName: setting?.WEBSITE.site_name || defaultSettings.WEBSITE.site_name,
			lastName: setting?.WEBSITE.site_tagline || defaultSettings.WEBSITE.site_tagline,
			username: 'idtubexxi'
		}
	}
});


function isValidUrl(url: string): boolean {
	try {
		new URL(url);
		return true;
	} catch (_) {
		return false;
	}
}
