import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingAdsTxt } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';
import { readFileSync, writeFileSync, existsSync } from 'fs';
import { join } from 'path';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update Ads.txt Settings', lang)) as SingleResponse;
	const siteName = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_name || '',
		lang
	)) as SingleResponse;
	const tagline = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_tagline || '',
		lang
	)) as SingleResponse;
	const description = (await deps.languageHelper.singleTranslate(
		settings?.WEBSITE?.site_description || '',
		lang
	)) as SingleResponse;
	const keywords = await Promise.all(
		(settings?.WEBSITE?.site_keywords?.split(',') || ['']).map(
			async (keyword) =>
				(await deps.languageHelper.singleTranslate(keyword.trim(), lang)) as SingleResponse
		)
	);

	const pageMetaTags = defaultMetaTags(
		{
			path_url: defaultOrigin,
			title: `${capitalizeFirstLetter(title.data.target.text || '')} - ${capitalizeFirstLetter(siteName.data.target.text || '')}`,
			tagline: capitalizeFirstLetter(tagline.data.target.text || ''),
			description: capitalizeFirstLetter(description.data.target.text || ''),
			keywords: keywords.map((keyword: SingleResponse) =>
				capitalizeFirstLetter(keyword.data.target.text || '')
			),
			robots: 'noindex, nofollow',
			canonical: defaultOrigin,
			alternates,
			graph_type: 'website'
		},
		settings
	);

	const adsTxt = await GetAdsTxtContent();
	const form = await superValidate(
		{
			content: adsTxt
		},
		zod4(updateSettingAdsTxt)
	);

	return {
		pageMetaTags,
		form,
		settings,
		user,
		lang
	};
};
export const actions = {
	default: async ({ request, locals }) => {
		const form = await superValidate(request, zod4(updateSettingAdsTxt));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const error = await UpdateAdsTxtContent(form.data.content || '');
		if (error instanceof Error) {
			return fail(500, {
				form,
				message: error.message || 'Failed to update ads.txt settings.'
			});
		}

		return {
			form,
			message: 'ads.txt updated successfully.'
		};
	}
};
async function GetAdsTxtContent(): Promise<string> {
	try {

		const staticPath = join(process.cwd(), 'static', 'ads.txt');
		let adsContent: string;


		if (existsSync(staticPath)) {
			adsContent = readFileSync(staticPath, 'utf-8');
		} else {
			adsContent = generateDefaultAdsTxt();

			writeFileSync(staticPath, adsContent, 'utf-8');
		}
		return adsContent;
	} catch (error) {
		console.error('Error reading ads.txt:', error);
		const fallbackContent = generateDefaultAdsTxt();
		return fallbackContent;
	}
}
async function UpdateAdsTxtContent(content: string): Promise<string | Error> {
	try {
		const staticPath = join(process.cwd(), 'static', 'ads.txt');
		writeFileSync(staticPath, content, 'utf-8');
		return 'success';
	} catch (error) {
		console.error('Error writing ads.txt:', error);
		return error instanceof Error ? error : new Error('Unknown error');
	}
}
function generateDefaultAdsTxt(): string {
	return `# Default ads.txt for ${process.env.ORIGIN || 'your-domain.com'}
google.com, pub-0000000000000000, DIRECT, f08c47fec0942fa0
google.com, pub-0000000000000001, RESELLER
# Add your own ad network entries here
# Contact: admin@your-domain.com`;
}
