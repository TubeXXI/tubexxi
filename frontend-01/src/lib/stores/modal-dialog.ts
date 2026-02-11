import type { Component } from 'svelte';
import { modalStore, type ModalConfig } from "@/hooks/modal.svelte";

export const modal = {
	/**
	 * Show an alert modal with a title, content, and optional footer.
	 * @param title - The title of the modal.
	 * @param content - The content of the modal.
	 * @param options - Optional modal configuration.
	 * @returns The ID of the opened modal.
	 */
	alert(title: string, content: string, options?: Partial<ModalConfig>) {
		return modalStore.open({
			title,
			content,
			footer: () => `
				<div class="flex justify-end">
					<button
						onclick="modalStore.close()"
						class="bg-primary text-primary-foreground hover:bg-primary/90 shadow-xs focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 cursor-pointer active:scale-95 h-9 px-4 py-2 has-[>svg]:px-3"
					>
						OK
					</button>
				</div>
			`,
			...options
		});
	},

	/**
	 * Show a confirmation modal with a title, content, and optional footer.
	 * @param title - The title of the modal.
	 * @param content - The content of the modal.
	 * @param onConfirm - Callback function to be called when the confirm button is clicked.
	 * @param options - Optional modal configuration.
	 * @returns The ID of the opened modal.
	 */
	confirm(
		title: string,
		content: string,
		onConfirm: () => void,
		options?: Partial<ModalConfig>
	) {
		return modalStore.open({
			title,
			content,
			footer: () => `
				<div class="flex justify-end gap-2">
					<button
						onclick="modalStore.close()"
						class="focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 cursor-pointer active:scale-95 bg-background hover:bg-accent hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50 border shadow-xs h-9 px-4 py-2 has-[>svg]:px-3"
					>
						Cancel
					</button>
					<button
						onclick="handleConfirm()"
						class="bg-primary text-primary-foreground hover:bg-primary/90 shadow-xs focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive inline-flex shrink-0 items-center justify-center gap-2 rounded-md text-sm font-medium whitespace-nowrap transition-all outline-none focus-visible:ring-[3px] disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 cursor-pointer active:scale-95 h-9 px-4 py-2 has-[>svg]:px-3"
					>
						Confirm
					</button>
				</div>
			`,
			onClose: () => {
				// Cleanup event listener
			},
			...options
		});
	},

	/**
	 * Show a custom modal component with optional props and configuration.
	 * @param component - The custom modal component to be rendered.
	 * @param props - Optional props to pass to the custom modal component.
	 * @param options - Optional modal configuration.
	 * @returns The ID of the opened modal.
	 */
	custom(component: Component, props?: Record<string, any>, options?: Partial<ModalConfig>) {
		return modalStore.open({
			component,
			props,
			...options
		});
	},

	/**
	 * Open a modal with all options.
	 * @param config - Partial modal configuration.
	 * @returns The ID of the opened modal.
	 */
	open(config: Partial<ModalConfig>) {
		return modalStore.open(config);
	},

	/**
	 * Close a modal with the specified ID.
	 * @param id - The ID of the modal to close.
	 */
	close(id?: string) {
		modalStore.close(id);
	},

	/**
	 * Close all open modals.
	 */
	closeAll() {
		modalStore.closeAll();
	}
};

export { modalStore };
