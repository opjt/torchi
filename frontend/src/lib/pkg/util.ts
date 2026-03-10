// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function debugLog(...args: any[]) {
	if (import.meta.env.DEV) {
		console.log('[DEBUG]', ...args);
	}
}

export function isDev() {
	return import.meta.env.DEV;
}

export function linkify(text: string): string {
	const urlRegex = /(https?:\/\/[^\s]+)/g;
	return text
		.replace(/&/g, '&amp;')
		.replace(/</g, '&lt;')
		.replace(/>/g, '&gt;')
		.replace(
			urlRegex,
			(url) =>
				`<a href="${url}" target="_blank" rel="noopener noreferrer" class=" underline underline-offset-2 break-all hover:opacity-70 transition-opacity">${url}</a>`,
		);
}
