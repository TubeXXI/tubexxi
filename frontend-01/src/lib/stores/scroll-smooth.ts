import { writable } from 'svelte/store';

interface ScrollState {
	isScrolling: boolean;
	lastTarget: string | null;
	isSupported: boolean;
}

function createSmoothScroll() {
	const { subscribe, set, update } = writable<ScrollState>({
		isScrolling: false,
		lastTarget: null,
		isSupported: typeof window !== 'undefined'
	});

	const scrollToAnchor = (id: string, offset: number = 80): void => {
		if (typeof window === 'undefined' || !document) return;

		const element = document.querySelector(id);
		if (!element) {
			console.warn(`Element ${id} not found`);
			return;
		}

		const elementPosition = element.getBoundingClientRect().top;
		const offsetPosition = elementPosition + window.pageYOffset - offset;

		window.scrollTo({
			top: offsetPosition,
			behavior: 'smooth' as ScrollBehavior,
		});
	};

	const scrollToTop = (): void => {
		if (typeof window === 'undefined') return;

		window.scrollTo({
			top: 0,
			behavior: 'smooth' as ScrollBehavior
		});
	};

	return {
		subscribe,
		scrollToAnchor,
		scrollToTop
	};
}

export const smoothScroll = createSmoothScroll();
