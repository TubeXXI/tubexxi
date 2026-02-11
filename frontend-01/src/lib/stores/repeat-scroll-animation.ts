import { writable } from 'svelte/store';
import { browser } from '$app/environment';

interface RepeatAnimationState {
	animations: Map<string, RepeatAnimationConfig>;
	activeElements: Set<string>;
	scrollDirection: 'up' | 'down';
	lastScrollY: number;
}

interface RepeatAnimationConfig {
	element: HTMLElement;
	animationType: keyof typeof REPEAT_ANIMATIONS;
	isVisible: boolean;
	threshold: number;
	offset: number;
	direction: 'both' | 'enter' | 'exit';
}

interface RepeatAnimationOptions {
	animationType?: keyof typeof REPEAT_ANIMATIONS;
	threshold?: number;
	offset?: number;
	direction?: 'both' | 'enter' | 'exit';
	delay?: number;
	duration?: number;
}

const REPEAT_ANIMATIONS = {
	fadeInOut: {
		hidden: { opacity: '0', transform: 'translateY(20px)' },
		visible: { opacity: '1', transform: 'translateY(0)' },
		transition: 'all 0.5s ease-out'
	},
	scaleInOut: {
		hidden: { opacity: '0', transform: 'scale(0.8)' },
		visible: { opacity: '1', transform: 'scale(1)' },
		transition: 'all 0.4s ease-out'
	},
	slideInOutLeft: {
		hidden: { opacity: '0', transform: 'translateX(-50px)' },
		visible: { opacity: '1', transform: 'translateX(0)' },
		transition: 'all 0.5s ease-out'
	},
	slideInOutRight: {
		hidden: { opacity: '0', transform: 'translateX(50px)' },
		visible: { opacity: '1', transform: 'translateX(0)' },
		transition: 'all 0.5s ease-out'
	},
	blurInOut: {
		hidden: { opacity: '0', filter: 'blur(10px)' },
		visible: { opacity: '1', filter: 'blur(0)' },
		transition: 'all 0.5s ease-out'
	},
	rotateInOut: {
		hidden: { opacity: '0', transform: 'rotate(-10deg)' },
		visible: { opacity: '1', transform: 'rotate(0deg)' },
		transition: 'all 0.5s ease-out'
	},
	bounceInOut: {
		hidden: { opacity: '0', transform: 'scale(0.8)' },
		visible: { opacity: '1', transform: 'scale(1)' },
		transition: 'all 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55)'
	}
} as const;

function createRepeatScrollAnimationStore() {
	const { subscribe, set, update } = writable<RepeatAnimationState>({
		animations: new Map(),
		activeElements: new Set(),
		scrollDirection: 'down',
		lastScrollY: 0
	});

	const observers = new Map<string, IntersectionObserver>();
	let scrollHandler: (() => void) | null = null;
	let animationTimeouts = new Map<string, NodeJS.Timeout>();

	if (browser) {
		set({
			animations: new Map(),
			activeElements: new Set(),
			scrollDirection: 'down',
			lastScrollY: window.scrollY
		});

		scrollHandler = () => {
			update(state => {
				const currentScrollY = window.scrollY;
				const direction = currentScrollY > state.lastScrollY ? 'down' : 'up';

				return {
					...state,
					scrollDirection: direction,
					lastScrollY: currentScrollY
				};
			});
		};

		window.addEventListener('scroll', scrollHandler, { passive: true });
	}

	const registerElement = (
		element: HTMLElement,
		options: RepeatAnimationOptions = {}
	): { id: string; destroy: () => void } | null => {
		if (!browser || !element) return null;

		const {
			animationType = 'fadeInOut',
			threshold = 0.1,
			offset = 0,
			direction = 'both',
			delay = 0,
			duration = 0.5
		} = options;

		const id = `repeat-anim-${Math.random().toString(36).substring(2, 11)}`;
		const animation = REPEAT_ANIMATIONS[animationType];

		Object.assign(element.style, {
			...animation.hidden,
			transition: `${animation.transition}, opacity ${duration}s ease ${delay}s, transform ${duration}s ease ${delay}s`,
			willChange: 'opacity, transform'
		});

		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((entry) => {
					if (!(entry.target instanceof HTMLElement)) return;

					update(state => {
						const config = state.animations.get(id);
						if (!config) return state;

						const shouldAnimate =
							direction === 'both' ||
							(direction === 'enter' && entry.isIntersecting) ||
							(direction === 'exit' && !entry.isIntersecting);

						if (shouldAnimate) {
							if (animationTimeouts.has(id)) {
								clearTimeout(animationTimeouts.get(id)!);
							}

							const timeoutId = setTimeout(() => {
								if (entry.isIntersecting) {

									Object.assign((entry.target as HTMLElement).style, animation.visible);
								} else {

									Object.assign((entry.target as HTMLElement).style, animation.hidden);
								}
							}, delay * 1000);

							animationTimeouts.set(id, timeoutId);

							return {
								...state,
								animations: new Map(state.animations.set(id, {
									...config,
									isVisible: entry.isIntersecting
								})),
								activeElements: entry.isIntersecting
									? new Set(state.activeElements).add(id)
									: new Set([...state.activeElements].filter(x => x !== id))
							};
						}

						return state;
					});
				});
			},
			{
				threshold,
				rootMargin: `${offset}px 0px`
			}
		);

		observer.observe(element);
		observers.set(id, observer);

		update(state => ({
			...state,
			animations: new Map(state.animations.set(id, {
				element,
				animationType,
				isVisible: false,
				threshold,
				offset,
				direction
			}))
		}));

		return {
			id,
			destroy: () => {
				const observer = observers.get(id);
				if (observer) {
					observer.disconnect();
					observers.delete(id);
				}

				const timeoutId = animationTimeouts.get(id);
				if (timeoutId) {
					clearTimeout(timeoutId);
					animationTimeouts.delete(id);
				}

				update(state => {
					const newAnimations = new Map(state.animations);
					newAnimations.delete(id);

					const newActiveElements = new Set(state.activeElements);
					newActiveElements.delete(id);

					return {
						...state,
						animations: newAnimations,
						activeElements: newActiveElements
					};
				});
			}
		};
	};

	const triggerAnimation = (elementId: string, forceState?: 'show' | 'hide'): void => {
		update(state => {
			const config = state.animations.get(elementId);
			if (!config || !config.element) return state;

			const animation = REPEAT_ANIMATIONS[config.animationType];

			if (forceState === 'show' || (!forceState && !config.isVisible)) {
				Object.assign(config.element.style, animation.visible);
			} else {
				Object.assign(config.element.style, animation.hidden);
			}

			return {
				...state,
				animations: new Map(state.animations.set(elementId, {
					...config,
					isVisible: forceState === 'show' || (!forceState && !config.isVisible)
				}))
			};
		});
	};

	const triggerAllVisible = (): void => {
		update(state => {
			const newAnimations = new Map(state.animations);

			newAnimations.forEach((config, id) => {
				const rect = config.element.getBoundingClientRect();
				const isInViewport =
					rect.top <= window.innerHeight * (1 - config.threshold) &&
					rect.bottom >= window.innerHeight * config.threshold &&
					rect.left <= window.innerWidth &&
					rect.right >= 0;

				const animation = REPEAT_ANIMATIONS[config.animationType];

				if (isInViewport) {
					Object.assign(config.element.style, animation.visible);
					newAnimations.set(id, { ...config, isVisible: true });
				} else {
					Object.assign(config.element.style, animation.hidden);
					newAnimations.set(id, { ...config, isVisible: false });
				}
			});

			return { ...state, animations: newAnimations };
		});
	};

	const cleanup = (): void => {
		observers.forEach(observer => observer.disconnect());
		observers.clear();

		animationTimeouts.forEach(timeout => clearTimeout(timeout));
		animationTimeouts.clear();

		if (scrollHandler) {
			window.removeEventListener('scroll', scrollHandler);
		}

		set({
			animations: new Map(),
			activeElements: new Set(),
			scrollDirection: 'down',
			lastScrollY: 0
		});
	};

	return {
		subscribe,
		registerElement,
		triggerAnimation,
		triggerAllVisible,
		cleanup,
		animations: REPEAT_ANIMATIONS
	};
}

export const repeatScrollAnimation = createRepeatScrollAnimationStore();
