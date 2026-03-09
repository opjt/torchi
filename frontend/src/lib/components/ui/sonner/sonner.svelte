<script lang="ts">
	import CircleCheckIcon from '@lucide/svelte/icons/circle-check';
	import InfoIcon from '@lucide/svelte/icons/info';
	import Loader2Icon from '@lucide/svelte/icons/loader-2';
	import OctagonXIcon from '@lucide/svelte/icons/octagon-x';
	import TriangleAlertIcon from '@lucide/svelte/icons/triangle-alert';

	import { Toaster as Sonner, type ToasterProps as SonnerProps } from 'svelte-sonner';
	import { mode } from 'mode-watcher';

	let { ...restProps }: SonnerProps = $props();
</script>

<Sonner
	theme={mode.current}
	class="toaster group"
	toastOptions={{
		classes: {
			// toast: 'rounded-box border border-base-300 bg-base-content! text-base-200! shadow-none!',
			success: '',
			error: '',
			warning: '',
			info: ''
		}
	}}
	{...restProps}
>
	{#snippet loadingIcon()}
		<Loader2Icon class="size-4 animate-spin opacity-70" />
	{/snippet}

	{#snippet successIcon()}
		<CircleCheckIcon class="size-4 text-success" />
	{/snippet}

	{#snippet errorIcon()}
		<OctagonXIcon class="size-4 text-error" />
	{/snippet}

	{#snippet infoIcon()}
		<InfoIcon class="size-4 text-info" />
	{/snippet}

	{#snippet warningIcon()}
		<TriangleAlertIcon class="size-4 text-warning" />
	{/snippet}
</Sonner>
