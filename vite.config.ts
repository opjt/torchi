import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';
import path from 'path';
import { fileURLToPath } from 'url';
import { readFileSync } from 'fs';

const file = fileURLToPath(new URL('package.json', import.meta.url));
const json = readFileSync(file, 'utf8');
const pkg = JSON.parse(json);

export default defineConfig({
	resolve: { alias: { $lib: path.resolve('./src/lib'), $src: path.resolve('./src') } },
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:25565',
				changeOrigin: true,
				secure: false,
				// rewrite: (path) => path.replace(/^\/api/, ''),
			},
		},
	},
	plugins: [
		tailwindcss(),
		sveltekit(),

		SvelteKitPWA({
			srcDir: './src',
			mode: 'production', // 'development' or 'production',
			strategies: 'injectManifest',
			filename: 'service-worker.ts',

			// Manifest 설정
			manifest: {
				name: 'torchi',
				short_name: 'torchi',
				description: 'On-demand hook push notification platform',
				background_color: '#ffffff',
				theme_color: 'transparent',
				display: 'fullscreen',
				scope: '/app',
				start_url: '/app',
				orientation: 'portrait',
				id: '/app',
				icons: [
					{
						src: '/logo/icon-192x192.png',
						sizes: '192x192',
						type: 'image/png',
						purpose: 'any',
					},
					{
						src: '/logo/icon-512x512.png',
						sizes: '512x512',
						type: 'image/png',
						purpose: 'any maskable',
					},
				],
			},
			// Service Worker 설정
			injectManifest: {
				swDest: 'service-worker.js',
				// globPatterns: ['**/*.{js,css,html,ico,png,svg,webp,woff,woff2}']
				globPatterns: ['**/*.{js,css,html,ico,png,svg,webmanifest,wasm}'],
			},

			// 개발 모드에서도 PWA 테스트
			devOptions: { enabled: true, type: 'module', navigateFallback: '/' },
		}),
	],
	define: {
		__APP_VERSION__: JSON.stringify(pkg.version),
	},
});
