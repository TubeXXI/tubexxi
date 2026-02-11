import type { Component } from 'svelte';

type ModalConfig = {
	id: string;
	component?: Component;
	props?: Record<string, any>;
	title?: string;
	description?: string;
	content?: string;
	footer?: any;
	size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
	customClass?: string;
	transparent?: boolean;
	hasBackground?: boolean;
	disableScroll?: boolean;
	clickOutside?: boolean;
	animation?: boolean;
	showCloseButton?: boolean;
	closeButtonPosition?: 'in-top-left' | 'in-top-right' | 'out-top-left' | 'out-top-right';
	onClose?: () => void;
	onOpen?: () => void;
};

type ModalState = {
	modals: ModalConfig[];
	activeModal: ModalConfig | null;
};

class ModalStore {
	private state = $state<ModalState>({
		modals: [],
		activeModal: null
	});

	get modals() {
		return this.state.modals;
	}

	get activeModal() {
		return this.state.activeModal;
	}

	get isOpen() {
		return this.state.activeModal !== null;
	}

	/**
	 * Open a modal with the provided configuration.
	 * @param config - Partial modal configuration.
	 * @returns The ID of the opened modal.
	 */
	open(config: Omit<ModalConfig, 'id'>) {
		const id = `modal-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
		const modalConfig: ModalConfig = {
			id,
			size: 'md',
			transparent: false,
			hasBackground: true,
			disableScroll: true,
			clickOutside: true,
			animation: true,
			showCloseButton: true,
			closeButtonPosition: 'in-top-right',
			...config
		};

		this.state.modals = [...this.state.modals, modalConfig];
		this.state.activeModal = modalConfig;

		// Disable scroll if needed
		if (modalConfig.disableScroll) {
			document.body.style.overflow = 'hidden';
		}

		// Call onOpen callback
		if (modalConfig.onOpen) {
			modalConfig.onOpen();
		}

		return id;
	}

	/**
	 * Close a modal with the specified ID.
	 * @param id - The ID of the modal to close.
	 */
	close(id?: string) {
		const modalToClose = id
			? this.state.modals.find(m => m.id === id)
			: this.state.activeModal;

		if (!modalToClose) return;

		// Call onClose callback
		if (modalToClose.onClose) {
			modalToClose.onClose();
		}

		// Remove modal from list
		this.state.modals = this.state.modals.filter(m => m.id !== modalToClose.id);

		// Update active modal
		this.state.activeModal = this.state.modals[this.state.modals.length - 1] || null;

		// Re-enable scroll if no modals are open
		if (this.state.modals.length === 0) {
			document.body.style.overflow = '';
		}
	}

	/**
	 * Close all open modals.
	 */
	closeAll() {
		this.state.modals.forEach(modal => {
			if (modal.onClose) {
				modal.onClose();
			}
		});
		this.state.modals = [];
		this.state.activeModal = null;
		document.body.style.overflow = '';
	}

	/**
	 * Update a modal with the specified ID.
	 * @param id - The ID of the modal to update.
	 * @param updates - Partial modal configuration to apply.
	 */
	update(id: string, updates: Partial<ModalConfig>) {
		const modalIndex = this.state.modals.findIndex(m => m.id === id);
		if (modalIndex !== -1) {
			this.state.modals[modalIndex] = { ...this.state.modals[modalIndex], ...updates };
			if (this.state.activeModal?.id === id) {
				this.state.activeModal = this.state.modals[modalIndex];
			}
		}
	}
}

export const modalStore = new ModalStore();
export type { ModalConfig };
