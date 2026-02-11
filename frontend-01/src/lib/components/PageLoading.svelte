<script>
	import { customPageLoading as loading } from '$lib/stores/page-loader';
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';

	$effect(() => {
		if (typeof document !== 'undefined') {
			if ($loading) {
				document.body.style.overflow = 'hidden';
			} else {
				document.body.style.overflow = '';
			}
		}
	});

	onMount(() => {
		return () => {
			if (typeof document !== 'undefined') {
				document.body.style.overflow = '';
			}
		};
	});
</script>

{#if $loading}
	<div
		class="fixed inset-0 z-9999 flex items-center justify-center bg-black/50 backdrop-blur-sm"
		transition:fade={{ duration: 200 }}
	>
		<div class="flex flex-col items-center gap-4">
			<div class="relative h-16 w-16">
				<div
					class="absolute h-full w-full rounded-full border-4 border-neutral-300 opacity-25"
				></div>
				<div
					class="absolute h-full w-full animate-spin rounded-full border-4 border-transparent border-t-white"
				></div>
			</div>

			<p class="text-lg font-medium text-white">Loading...</p>
		</div>
	</div>
{/if}
