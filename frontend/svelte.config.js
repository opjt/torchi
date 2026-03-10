import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		adapter: adapter({
			fallback: 'index.html', // 404 대신 index.html을 서빙하도록 함
		}),
		prerender: {
			handleHttpError: ({ path, message }) => {
				if (path === '/manifest.webmanifest') return;
				throw new Error(message);
			},
		},
		paths: {
			// 이 설정을 비워두거나 '/'로 인식되게 하면
			// 모든 자원(JS, CSS)을 절대 경로인 '/_app/...'으로 불러옵니다.
			base: '',
		},
	},
};

export default config;
