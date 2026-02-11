<script lang="ts">
	import type { Snippet } from 'svelte';
	import { modalStore } from '@/stores';
	import { onMount } from 'svelte';

	let {
		children,
		snippet
	}: {
		children?: Snippet;
		snippet?: Snippet;
	} = $props();

	let mounted = $state(false);

	onMount(() => {
		mounted = true;
		return () => {
			// Cleanup on unmount
			modalStore.closeAll();
		};
	});

	function handleBackdropClick(e: MouseEvent) {
		if (!modalStore.activeModal) return;

		const target = e.target as HTMLElement;
		if (target.classList.contains('modal-backdrop') && modalStore.activeModal.clickOutside) {
			modalStore.close();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && modalStore.activeModal) {
			modalStore.close();
		}
	}

	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl',
		full: 'max-w-full mx-4'
	};

	const closeButtonPositions = {
		'in-top-left': 'absolute top-4 left-4 z-10',
		'in-top-right': 'absolute top-4 right-4 z-10',
		'out-top-left': 'absolute -top-10 -left-10',
		'out-top-right': 'absolute -top-10 -right-10'
	};
</script>

<svelte:window onkeydown={handleKeydown} />

{#if mounted && modalStore.isOpen && modalStore.activeModal}
	{@const modal = modalStore.activeModal}
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<div
		class="modal-backdrop fixed inset-0 z-50 flex items-center justify-center p-4"
		class:bg-black={modal.hasBackground}
		class:bg-opacity-50={modal.hasBackground}
		class:backdrop-blur-sm={modal.hasBackground}
		class:animate-fade-in={modal.animation}
		onclick={handleBackdropClick}
		role="dialog"
		tabindex="-1"
		aria-modal="true"
		aria-labelledby={modal.title ? 'modal-title' : undefined}
		aria-describedby={modal.description ? 'modal-description' : undefined}
	>
		<div
			class="modal-content relative w-full {sizeClasses[modal.size || 'md']} {modal.customClass ||
				''}"
			class:animate-scale-in={modal.animation}
			role="document"
		>
			<!-- Close Button -->
			{#if modal.showCloseButton}
				<button
					type="button"
					class="{closeButtonPositions[modal.closeButtonPosition || 'in-top-right']}
						flex h-8 w-8 items-center justify-center rounded-full
						bg-neutral-200 text-neutral-700 transition-colors duration-200
						hover:bg-neutral-300 focus:ring-2
						focus:ring-blue-500 focus:ring-offset-2
						focus:outline-none dark:bg-neutral-700 dark:text-neutral-200 dark:hover:bg-neutral-600"
					onclick={() => modalStore.close()}
					aria-label="Close modal"
				>
					<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					</svg>
				</button>
			{/if}

			<!-- Modal Box -->
			<div
				class="modal-box overflow-hidden rounded-lg shadow-xl"
				class:bg-white={!modal.transparent}
				class:dark:bg-neutral-800={!modal.transparent}
				class:bg-transparent={modal.transparent}
			>
				<!-- Header -->
				{#if modal.title || modal.description}
					<div class="modal-header border-b border-neutral-200 px-6 py-4 dark:border-neutral-700">
						{#if modal.title}
							<h2 id="modal-title" class="text-2xl font-bold text-neutral-900 dark:text-white">
								{modal.title}
							</h2>
						{/if}
						{#if modal.description}
							<p id="modal-description" class="mt-1 text-sm text-neutral-600 dark:text-neutral-400">
								{modal.description}
							</p>
						{/if}
					</div>
				{/if}

				<!-- Content -->
				<div class="modal-body px-6 py-4">
					{#if modal.component}
						<modal.component {...modal.props || {}} />
					{:else if modal.content}
						<div class="text-neutral-700 dark:text-neutral-300">
							{modal.content}
						</div>
					{:else if children}
						{@render children()}
					{:else if snippet}
						{@render snippet()}
					{/if}
				</div>

				<!-- Footer -->
				{#if modal.footer}
					<div
						class="modal-footer border-t border-neutral-200 bg-neutral-50 px-6 py-4 dark:border-neutral-700 dark:bg-neutral-900"
					>
						{#if typeof modal.footer === 'function'}
							{@render modal.footer()}
						{:else}
							{modal.footer}
						{/if}
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}

<style scoped>
	@keyframes fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	@keyframes scale-in {
		from {
			opacity: 0;
			transform: scale(0.95);
		}
		to {
			opacity: 1;
			transform: scale(1);
		}
	}

	.animate-fade-in {
		animation: fade-in 0.2s ease-out;
	}

	.animate-scale-in {
		animation: scale-in 0.3s ease-out;
	}
</style>
