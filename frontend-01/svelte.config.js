import { mdsvex } from 'mdsvex';
import adapter from '@sveltejs/adapter-node';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: { 
		adapter: adapter({
			out: 'build',
			precompress: true,
			dynamic_origin: true
		}),
		paths: {
			relative: false
		},
		csrf: {
			trustedOrigins:
				process.env.NODE_ENV === 'production'
					? ['https://simontokz.com', 'https://www.simontokz.com']
					: ['*']
		},
		alias: {
			'@': './src/lib',
			'@/*': './src/lib/*'
		}
	 },
	preprocess: [mdsvex()],
	extensions: ['.svelte', '.svx']
};

export default config;
