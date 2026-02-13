import { readFileSync, writeFileSync, existsSync } from 'fs';
import { join } from 'path';


export const GET = async ({ url, request }) => {
	try {

		const staticPath = join(process.cwd(), 'static', 'ads.txt');
		let adsContent: string;

		if (existsSync(staticPath)) {
			adsContent = readFileSync(staticPath, 'utf-8');
		} else {
			adsContent = generateDefaultAdsTxt();

			writeFileSync(staticPath, adsContent, 'utf-8');
		}

		return new Response(adsContent, {
			headers: {
				'Content-Type': 'text/plain; charset=utf-8',
				// 'Cache-Control': 'public, max-age=3600' // Cache 1 hour
			}
		});

	} catch (error) {
		console.error('Error handling ads.txt:', error);

		const fallbackContent = generateDefaultAdsTxt(origin);

		return new Response(fallbackContent, {
			status: 500,
			headers: {
				'Content-Type': 'text/plain; charset=utf-8'
			}
		});
	}
};

function generateDefaultAdsTxt(origin: string = process.env.ORIGIN || 'your-domain.com'): string {
	return `# Default ads.txt for ${origin}
google.com, pub-0000000000000000, DIRECT, f08c47fec0942fa0
google.com, pub-0000000000000001, RESELLER
# Add your own ad network entries here
# Contact: admin@your-domain.com`;
}
