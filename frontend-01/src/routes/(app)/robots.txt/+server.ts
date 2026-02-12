import { readFileSync, writeFileSync, existsSync } from 'fs';
import { join } from 'path';
import { text } from '@sveltejs/kit';


export const GET = async ({ locals, url, request }) => {
	const { settings } = locals
	try {

		const staticPath = join(process.cwd(), 'static', 'robots.txt');
		let adsContent: string;

		let origin = url.origin || request.headers.get('origin') || '';
		if (origin.startsWith('http://')) {
			origin = origin.replace('http://', 'https://');
		}

		if (existsSync(staticPath)) {
			adsContent = readFileSync(staticPath, 'utf-8');
		} else {
			adsContent = generateDefaultRobotTxt(settings?.WEBSITE.site_email || 'support@agcforge.com', origin);

			writeFileSync(staticPath, adsContent, 'utf-8');
		}

		return text(adsContent, {
			headers: {
				'Content-Type': 'text/plain; charset=utf-8',
				'Cache-Control': 'public, max-age=3600' // Cache 1 hour
			}
		});

	} catch (error) {
		console.error('Error handling robots.txt:', error);

		const fallbackContent = generateDefaultRobotTxt(settings?.WEBSITE.site_email || 'support@agcforge.com', origin);

		return text(fallbackContent, {
			status: 500,
			headers: {
				'Content-Type': 'text/plain; charset=utf-8'
			}
		});
	}
};

function generateDefaultRobotTxt(siteEmail?: string, origin?: string): string {
	return `# Default robots.txt for ${origin || process.env.ORIGIN || 'agcforge.com'}
User-agent: *
Disallow: /admin/
Disallow: /private/
Disallow: /tmp/
Disallow: /uploads/
# Add your own disallow entries here
# Contact: ${siteEmail || 'admin@agcforge.com'}`;
}
