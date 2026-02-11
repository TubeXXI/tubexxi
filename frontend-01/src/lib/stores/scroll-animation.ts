
import { writable } from 'svelte/store';
import { browser } from '$app/environment';


interface AnimationConfig {
	element: HTMLElement;
	isVisible: boolean;
	animationType: keyof typeof ANIMATIONS;
	progress: number;
}

interface ScrollAnimationState {
	animations: Map<string, AnimationConfig>;
	scrollY: number;
	viewportHeight: number;
	isIntersecting: boolean;
	thresholds: Map<string, number>;
}

interface AnimationPreset {
	hidden: Partial<CSSStyleDeclaration> & Record<string, string>;
	visible: Partial<CSSStyleDeclaration> & Record<string, string>;
	transition: string;
}

interface AnimationOptions {
	animationType?: keyof typeof ANIMATIONS;
	threshold?: number;
	rootMargin?: string;
	delay?: number;
	duration?: number;
	stagger?: number;
	once?: boolean;
	offset?: number;
}

interface RegisterElementReturn {
	id: string;
	destroy: () => void;
}

interface ParallaxOptions {
	speed?: number;
	direction?: 'vertical' | 'horizontal' | 'both';
	startOffset?: number;
	endOffset?: number;
}

interface TypewriterOptions {
	text?: string;
	speed?: number;
	delay?: number;
	cursor?: boolean;
	cursorChar?: string;
	infinite?: boolean;
}

interface ProgressAnimationOptions {
	property?: keyof CSSStyleDeclaration;
	startValue?: string;
	endValue?: string;
	duration?: number;
	easing?: string;
}

interface StaggerOptions extends Omit<AnimationOptions, 'stagger'> {
	stagger?: number;
}


const ANIMATIONS = {
	fadeIn: {
		hidden: { opacity: '0', transform: 'translateY(20px)' },
		visible: { opacity: '1', transform: 'translateY(0)' },
		transition: 'all 0.6s ease-out'
	},
	fadeOut: {
		hidden: { opacity: '1', transform: 'translateY(0)' },
		visible: { opacity: '0', transform: 'translateY(20px)' },
		transition: 'all 0.6s ease-out'
	},
	slideInLeft: {
		hidden: { opacity: '0', transform: 'translateX(-50px)' },
		visible: { opacity: '1', transform: 'translateX(0)' },
		transition: 'all 0.5s ease-out'
	},
	slideInRight: {
		hidden: { opacity: '0', transform: 'translateX(50px)' },
		visible: { opacity: '1', transform: 'translateX(0)' },
		transition: 'all 0.5s ease-out'
	},
	zoomIn: {
		hidden: { opacity: '0', transform: 'scale(0.9)' },
		visible: { opacity: '1', transform: 'scale(1)' },
		transition: 'all 0.4s ease-out'
	},
	flipIn: {
		hidden: { opacity: '0', transform: 'rotateX(90deg)' },
		visible: { opacity: '1', transform: 'rotateX(0)' },
		transition: 'all 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55)'
	},
	bounceIn: {
		hidden: { opacity: '0', transform: 'scale(0.3)' },
		visible: { opacity: '1', transform: 'scale(1)' },
		transition: 'all 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55)'
	},
	blurIn: {
		hidden: { opacity: '0', filter: 'blur(10px)' },
		visible: { opacity: '1', filter: 'blur(0)' },
		transition: 'all 0.5s ease-out'
	}
} as const satisfies Record<string, AnimationPreset>;

function createScrollAnimationStore() {
	const { subscribe, set, update } = writable<ScrollAnimationState>({
		animations: new Map(),
		scrollY: 0,
		viewportHeight: 0,
		isIntersecting: false,
		thresholds: new Map()
	});

	const observers = new Map<string, IntersectionObserver>();
	const animationCallbacks = new Map<string, () => void>();
	const scrollListeners = new Map<string, () => void>();

	if (browser) {
		set({
			animations: new Map(),
			scrollY: window.pageYOffset,
			viewportHeight: window.innerHeight,
			isIntersecting: false,
			thresholds: new Map()
		});

		const handleScroll = () => {
			update(state => ({
				...state,
				scrollY: window.pageYOffset
			}));
		};

		window.addEventListener('scroll', handleScroll, { passive: true });

		const handleResize = () => {
			update(state => ({
				...state,
				viewportHeight: window.innerHeight
			}));
		};

		window.addEventListener('resize', handleResize);

		scrollListeners.set('scroll', handleScroll);
		scrollListeners.set('resize', handleResize);
	}

	const registerElement = (
		element: HTMLElement,
		options: AnimationOptions = {}
	): RegisterElementReturn | null => {
		if (!browser || !element) return null;

		const {
			animationType = 'fadeIn',
			threshold = 0.1,
			rootMargin = '0px',
			delay = 0,
			duration = 0.6,
			stagger = 0,
			once = true,
			offset = 0
		} = options;

		const id = `anim-${Math.random().toString(36).substring(2, 11)}`;
		const animation = ANIMATIONS[animationType] || ANIMATIONS.fadeIn;

		Object.assign(element.style, {
			...animation.hidden,
			transition: `${animation.transition}, opacity ${duration}s ease ${delay}s, transform ${duration}s ease ${delay}s`,
			willChange: 'opacity, transform'
		});

		const observer = new IntersectionObserver(
			(entries) => {
				entries.forEach((entry) => {
					if (!(entry.target instanceof HTMLElement)) return;

					if (entry.isIntersecting) {
						setTimeout(() => {
							Object.assign((entry.target as HTMLElement).style, animation.visible);

							if (once) {
								observer.unobserve(entry.target);
								observers.delete(id);
							}
						}, stagger);

						update(state => ({
							...state,
							animations: new Map(state.animations.set(id, {
								element: entry.target as HTMLElement,
								isVisible: true,
								animationType,
								progress: 1
							}))
						}));
					} else if (!once) {
						Object.assign(entry.target.style, animation.hidden);

						update(state => ({
							...state,
							animations: new Map(state.animations.set(id, {
								element: entry.target as HTMLElement,
								isVisible: false,
								animationType,
								progress: 0
							}))
						}));
					}
				});
			},
			{
				threshold,
				rootMargin: `${offset}px 0px ${rootMargin}`
			}
		);

		observer.observe(element);
		observers.set(id, observer);

		animationCallbacks.set(id, () => {
			observer.unobserve(element);
			observer.disconnect();
		});

		return {
			id,
			destroy: () => {
				const cleanup = animationCallbacks.get(id);
				if (cleanup) cleanup();

				observers.delete(id);
				animationCallbacks.delete(id);

				update(state => {
					const newAnimations = new Map(state.animations);
					newAnimations.delete(id);
					return { ...state, animations: newAnimations };
				});
			}
		};
	};

	const registerElements = (
		elements: HTMLElement[],
		options: StaggerOptions = {}
	): (RegisterElementReturn | null)[] => {
		if (!browser || !elements.length) return [];

		const {
			stagger = 100,
			...elementOptions
		} = options;

		return elements.map((element, index) =>
			registerElement(element, {
				...elementOptions,
				delay: (index * stagger) / 1000,
				stagger: index * stagger
			})
		);
	};

	const createParallax = (
		element: HTMLElement,
		options: ParallaxOptions = {}
	): { id: string; destroy: () => void } | null => {
		if (!browser || !element) return null;

		const {
			speed = 0.5,
			direction = 'vertical'
		} = options;

		const id = `parallax-${Math.random().toString(36).substring(2, 11)}`;

		const handleScroll = () => {
			const rect = element.getBoundingClientRect();
			const viewportHeight = window.innerHeight;

			const elementTop = rect.top + window.pageYOffset;
			const scrollY = window.pageYOffset;
			const elementCenter = elementTop + rect.height / 2;
			const viewportCenter = scrollY + viewportHeight / 2;

			const distanceFromCenter = elementCenter - viewportCenter;
			const maxDistance = viewportHeight + rect.height;
			const progress = Math.max(-1, Math.min(1, distanceFromCenter / maxDistance));

			let transform = '';
			if (direction === 'vertical') {
				const translateY = progress * speed * 100;
				transform = `translateY(${translateY}px)`;
			} else if (direction === 'horizontal') {
				const translateX = progress * speed * 100;
				transform = `translateX(${translateX}px)`;
			} else if (direction === 'both') {
				const translate = progress * speed * 100;
				transform = `translate(${translate}px, ${translate}px)`;
			}

			element.style.transform = transform;
		};

		window.addEventListener('scroll', handleScroll, { passive: true });
		handleScroll();

		const destroy = () => {
			window.removeEventListener('scroll', handleScroll);
			element.style.transform = '';
			scrollListeners.delete(id);
		};

		scrollListeners.set(id, handleScroll);

		return { id, destroy };
	};

	const createTypewriter = (
		element: HTMLElement,
		options: TypewriterOptions = {}
	): {
		id: string;
		start: () => void;
		stop: () => void;
		reset: () => void;
		destroy: () => void;
	} | null => {
		if (!browser || !element) return null;

		const {
			text = element.textContent || '',
			speed = 50,
			delay = 0,
			cursor = true,
			cursorChar = '|',
			infinite = false
		} = options;

		const id = `typewriter-${Math.random().toString(36).substring(2, 11)}`;
		let currentText = '';
		let currentIndex = 0;
		let isTyping = false;
		let timeoutId: NodeJS.Timeout | null = null;
		let cursorInterval: NodeJS.Timeout | null = null;

		const type = () => {
			if (currentIndex < text.length) {
				currentText += text[currentIndex];
				element.textContent = currentText + (cursor ? cursorChar : '');
				// element.style = 'color: rgba(0, 0, 0, 0); background: linear-gradient(to right, #ff7e5f, #feb47b);'
				currentIndex++;
				timeoutId = setTimeout(type, speed);
			} else if (infinite) {
				setTimeout(() => {
					currentText = '';
					currentIndex = 0;
					type();
				}, 2000);
			} else if (cursor) {
				let cursorVisible = true;
				cursorInterval = setInterval(() => {
					element.textContent = currentText + (cursorVisible ? cursorChar : '');
					cursorVisible = !cursorVisible;
				}, 500);
			}
		};

		const start = () => {
			if (!isTyping) {
				isTyping = true;
				element.textContent = cursor ? cursorChar : '';
				timeoutId = setTimeout(type, delay);
			}
		};

		const stop = () => {
			isTyping = false;
			if (timeoutId) {
				clearTimeout(timeoutId);
				timeoutId = null;
			}
		};

		const reset = () => {
			stop();
			if (cursorInterval) {
				clearInterval(cursorInterval);
				cursorInterval = null;
			}
			currentText = '';
			currentIndex = 0;
			element.textContent = text;
		};

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting) {
					start();
					if (!infinite) observer.unobserve(element);
				}
			},
			{ threshold: 0.1 }
		);

		observer.observe(element);

		const destroy = () => {
			stop();
			reset();
			observer.unobserve(element);
			observer.disconnect();
		};

		return { id, start, stop, reset, destroy };
	};

	const createProgressAnimation = (
		element: HTMLElement,
		options: ProgressAnimationOptions = {}
	): { id: string; destroy: () => void } | null => {
		if (!browser || !element) return null;

		const {
			property = 'width',
			startValue = '0%',
			endValue = '100%',
			duration = 1000,
			easing = 'ease-out'
		} = options;

		const id = `progress-${Math.random().toString(36).substring(2, 11)}`;

		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting) {
					element.style.transition = `${property.toString()} ${duration}ms ${easing}`;
					(element.style as any)[property] = endValue;
					observer.unobserve(element);
				}
			},
			{ threshold: 0.1 }
		);

		observer.observe(element);
		(element.style as any)[property] = startValue;

		const destroy = () => {
			observer.unobserve(element);
			observer.disconnect();
		};

		return { id, destroy };
	};

	const getAnimationState = (elementId: string): AnimationConfig | undefined => {
		let animationState: AnimationConfig | undefined;

		subscribe(state => {
			animationState = state.animations.get(elementId);
		})();

		return animationState;
	};

	const isElementVisible = (elementId: string): boolean => {
		const state = getAnimationState(elementId);
		return state?.isVisible || false;
	};

	const cleanup = (): void => {
		observers.forEach(observer => {
			observer.disconnect();
		});
		observers.clear();

		animationCallbacks.forEach(callback => callback());
		animationCallbacks.clear();

		if (browser) {
			scrollListeners.forEach((listener, type) => {
				if (type === 'scroll') {
					window.removeEventListener('scroll', listener);
				} else if (type === 'resize') {
					window.removeEventListener('resize', listener);
				}
			});
			scrollListeners.clear();
		}

		update(state => ({
			...state,
			animations: new Map()
		}));
	};

	return {
		subscribe,
		registerElement,
		registerElements,
		createParallax,
		createTypewriter,
		createProgressAnimation,
		getAnimationState,
		isElementVisible,
		cleanup,
		animations: ANIMATIONS
	};
}

export const scrollAnimation = createScrollAnimationStore();

export type {
	AnimationOptions,
	ParallaxOptions,
	TypewriterOptions,
	ProgressAnimationOptions,
	RegisterElementReturn,
	AnimationPreset
};
