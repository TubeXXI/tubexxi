import { defaultMetaTags } from '@/utils/meta-tags.js';
import type { SingleResponse } from '@siamf/google-translate';
import { capitalizeFirstLetter } from '@/utils/format.js';
import { superValidate } from 'sveltekit-superforms';
import { fail } from '@sveltejs/kit';
import { updateSettingRobotTxt } from '$lib/utils/schema';
import { zod4 } from 'sveltekit-superforms/adapters';
import { readFileSync, writeFileSync, existsSync } from 'fs';
import { join } from 'path';

export const load = async ({ locals, parent }) => {
	const { user, settings, deps, lang } = locals;

	const defaultOrigin = await parent().then((data) => data.canonicalUrl || '');
	const alternates = await parent().then((data) => data.alternates || []);

	const title = (await deps.languageHelper.singleTranslate('Admin - Update robots.txt Settings', lang)) as SingleResponse;
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

	const robotTxt = await GetRobotTxtContent();
	const form = await superValidate(
		{
			content: robotTxt
		},
		zod4(updateSettingRobotTxt)
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
		const form = await superValidate(request, zod4(updateSettingRobotTxt));
		if (!form.valid) {
			return fail(400, {
				form,
				message: Object.values(form.errors).flat().join(', '),
			});
		}

		const error = await UpdateRobotTxtContent(form.data.content || '');
		if (error instanceof Error) {
			return fail(500, {
				form,
				message: error.message || 'Failed to update robots.txt settings.'
			});
		}

		return {
			form,
			message: 'robots.txt updated successfully.'
		};
	}
};
async function GetRobotTxtContent(): Promise<string> {
	try {

		const staticPath = join(process.cwd(), 'static', 'robots.txt');
		let robotContent: string;


		if (existsSync(staticPath)) {
			robotContent = readFileSync(staticPath, 'utf-8');
		} else {
			robotContent = generateDefaultRobotTxt();

			writeFileSync(staticPath, robotContent, 'utf-8');
		}
		return robotContent;
	} catch (error) {
		console.error('Error reading robots.txt:', error);
		const fallbackContent = generateDefaultRobotTxt();
		return fallbackContent;
	}
}
async function UpdateRobotTxtContent(content: string): Promise<string | Error> {
	try {
		const staticPath = join(process.cwd(), 'static', 'robots.txt');
		writeFileSync(staticPath, content, 'utf-8');
		return 'success';
	} catch (error) {
		console.error('Error writing robots.txt:', error);
		return error instanceof Error ? error : new Error('Unknown error');
	}
}
function generateDefaultRobotTxt(): string {
	return `# Default robots.txt for ${process.env.ORIGIN || 'your-domain.com'}
User-agent: *
Disallow: /
# Add your own robots.txt entries here
# Contact: admin@your-domain.com`;
}
