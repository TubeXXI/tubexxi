import { paraglideVitePlugin } from '@inlang/paraglide-js';
import devtoolsJson from 'vite-plugin-devtools-json';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	logLevel: 'info',
	build: {
		minify: true
	},
	server: {
		allowedHosts: [
			'client.giuadiario.info',
			'compositely-sanguinolent-cari.ngrok-free.dev',
			'agcforge.local'
		]
	},
	plugins: [
		tailwindcss(),
		sveltekit(),
		devtoolsJson(),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide',
			strategy: ['url', 'cookie', 'baseLocale']
		})
	],
	ssr: {
		noExternal: ['svelte-motion']
	},
	optimizeDeps: {
		include: ['svelte', 'svelte/internal']
	}
});
