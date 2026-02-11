
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface ScrollProgressConfig {
	element: HTMLElement;
	animationType: keyof typeof PROGRESS_ANIMATIONS;
	startThreshold: number;
	endThreshold: number;
	progress: number;
	direction: 'vertical' | 'horizontal';
	easing?: (t: number) => number;
}

interface ScrollProgressOptions {
	animationType?: keyof typeof PROGRESS_ANIMATIONS;
	startThreshold?: number;
	endThreshold?: number;
	direction?: 'vertical' | 'horizontal';
	easing?: (t: number) => number;
}


const PROGRESS_ANIMATIONS = {
	fadeProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			transform: `translateY(${20 * (1 - progress)}px)`
		})
	},
	scaleProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			transform: `scale(${0.8 + 0.2 * progress})`
		})
	},
	rotateProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			transform: `rotate(${-10 + 10 * progress}deg)`
		})
	},
	blurProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			filter: `blur(${10 * (1 - progress)}px)`
		})
	},
	slideLeftProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			transform: `translateX(${-50 * (1 - progress)}px)`
		})
	},
	slideRightProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			transform: `translateX(${50 * (1 - progress)}px)`
		})
	},
	colorProgress: {
		calculate: (progress: number) => ({
			opacity: progress.toString(),
			backgroundColor: `rgb(
        ${Math.floor(255 * progress)},
        ${Math.floor(100 + 155 * progress)},
        ${Math.floor(50 + 205 * progress)}
      )`
		})
	}
} as const;

function createScrollProgressAnimationStore() {
	const { subscribe, set, update } = writable({
		animations: new Map<string, ScrollProgressConfig>(),
		scrollY: 0,
		viewportHeight: 0,
		viewportWidth: 0
	});

	let scrollHandler: (() => void) | null = null;
	let resizeHandler: (() => void) | null = null;

	if (browser) {
		set({
			animations: new Map(),
			scrollY: window.scrollY,
			viewportHeight: window.innerHeight,
			viewportWidth: window.innerWidth
		});

		scrollHandler = () => {
			const scrollY = window.scrollY;
			const viewportHeight = window.innerHeight;

			update(state => {
				const newAnimations = new Map(state.animations);

				newAnimations.forEach((config, id) => {
					const rect = config.element.getBoundingClientRect();
					const elementTop = rect.top + scrollY;
					const elementBottom = rect.bottom + scrollY;
					const elementHeight = rect.height;

					let progress = 0;

					if (config.direction === 'vertical') {
						const viewportTop = scrollY;
						const viewportBottom = scrollY + viewportHeight;

						const startPoint = viewportBottom - (viewportHeight * config.startThreshold);
						const endPoint = viewportTop + (viewportHeight * config.endThreshold);

						if (elementTop < startPoint && elementBottom > endPoint) {
							const totalDistance = startPoint - endPoint;
							const elementDistance = startPoint - elementTop;
							progress = Math.max(0, Math.min(1, elementDistance / totalDistance));

							if (config.easing) {
								progress = config.easing(progress);
							}
						} else if (elementTop >= startPoint) {
							progress = 0;
						} else if (elementBottom <= endPoint) {
							progress = 1;
						}
					}

					const animation = PROGRESS_ANIMATIONS[config.animationType];
					const styles = animation.calculate(progress);
					Object.assign(config.element.style, styles);

					newAnimations.set(id, { ...config, progress });
				});

				return {
					...state,
					scrollY,
					animations: newAnimations
				};
			});
		};


		resizeHandler = () => {
			update(state => ({
				...state,
				viewportHeight: window.innerHeight,
				viewportWidth: window.innerWidth
			}));
			if (scrollHandler) scrollHandler();
		};

		window.addEventListener('scroll', scrollHandler, { passive: true });
		window.addEventListener('resize', resizeHandler);
	}

	const registerElement = (
		element: HTMLElement,
		options: ScrollProgressOptions = {}
	): { id: string; destroy: () => void } | null => {
		if (!browser || !element) return null;

		const {
			animationType = 'fadeProgress',
			startThreshold = 0.25,
			endThreshold = 0.5,
			direction = 'vertical',
			easing
		} = options;

		const id = `progress-anim-${Math.random().toString(36).substring(2, 11)}`;

		const animation = PROGRESS_ANIMATIONS[animationType];
		const initialStyles = animation.calculate(0);
		Object.assign(element.style, {
			...initialStyles,
			willChange: 'opacity, transform',
			transition: 'opacity 0.1s, transform 0.1s, filter 0.1s'
		});

		update(state => ({
			...state,
			animations: new Map(state.animations.set(id, {
				element,
				animationType,
				startThreshold,
				endThreshold,
				progress: 0,
				direction,
				easing
			}))
		}));

		if (scrollHandler) setTimeout(() => scrollHandler(), 100);

		return {
			id,
			destroy: () => {
				update(state => {
					const newAnimations = new Map(state.animations);
					newAnimations.delete(id);


					Object.assign(element.style, {
						opacity: '',
						transform: '',
						filter: '',
						backgroundColor: '',
						willChange: '',
						transition: ''
					});

					return { ...state, animations: newAnimations };
				});
			}
		};
	};

	const getProgress = (elementId: string): number => {
		let progress = 0;

		subscribe(state => {
			const config = state.animations.get(elementId);
			if (config) progress = config.progress;
		})();

		return progress;
	};

	const cleanup = (): void => {
		if (scrollHandler) {
			window.removeEventListener('scroll', scrollHandler);
		}
		if (resizeHandler) {
			window.removeEventListener('resize', resizeHandler);
		}

		set({
			animations: new Map(),
			scrollY: 0,
			viewportHeight: 0,
			viewportWidth: 0
		});
	};

	return {
		subscribe,
		registerElement,
		getProgress,
		cleanup,
		animations: PROGRESS_ANIMATIONS
	};
}

export const scrollProgressAnimation = createScrollProgressAnimationStore();
