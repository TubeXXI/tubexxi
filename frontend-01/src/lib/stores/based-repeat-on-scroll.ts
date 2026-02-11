
import { browser } from '$app/environment';
import type { Action } from 'svelte/action';

interface RepeatOnScrollOptions {
	animation?: 'fade' | 'slide' | 'scale' | 'blur' | 'rotate';
	threshold?: number;
	offset?: number;
	direction?: 'both' | 'enter' | 'exit';
	delay?: number;
	duration?: number;
	once?: boolean;
}

export const repeatOnScroll: Action<HTMLElement, RepeatOnScrollOptions> = (node, options = {}) => {
	if (!browser) return {};

	const {
		animation = 'fade',
		threshold = 0.1,
		offset = 0,
		direction = 'both',
		delay = 0,
		duration = 0.5,
		once = false
	} = options;

	const animations = {
		fade: {
			hidden: { opacity: '0', transform: 'translateY(20px)' },
			visible: { opacity: '1', transform: 'translateY(0)' }
		},
		slide: {
			hidden: { opacity: '0', transform: 'translateX(-50px)' },
			visible: { opacity: '1', transform: 'translateX(0)' }
		},
		scale: {
			hidden: { opacity: '0', transform: 'scale(0.8)' },
			visible: { opacity: '1', transform: 'scale(1)' }
		},
		blur: {
			hidden: { opacity: '0', filter: 'blur(10px)' },
			visible: { opacity: '1', filter: 'blur(0)' }
		},
		rotate: {
			hidden: { opacity: '0', transform: 'rotate(-10deg)' },
			visible: { opacity: '1', transform: 'rotate(0deg)' }
		}
	};

	const selectedAnimation = animations[animation];


	Object.assign(node.style, {
		...selectedAnimation.hidden,
		transition: `all ${duration}s ease ${delay}s`,
		willChange: 'opacity, transform, filter'
	});

	let timeoutId: NodeJS.Timeout;

	const observer = new IntersectionObserver(
		(entries) => {
			entries.forEach(entry => {
				const shouldAnimate =
					direction === 'both' ||
					(direction === 'enter' && entry.isIntersecting) ||
					(direction === 'exit' && !entry.isIntersecting);

				if (shouldAnimate) {
					clearTimeout(timeoutId);

					timeoutId = setTimeout(() => {
						if (entry.isIntersecting) {
							Object.assign((entry.target as HTMLElement).style, selectedAnimation.visible);
						} else {
							Object.assign((entry.target as HTMLElement).style, selectedAnimation.hidden);
						}

						if (once) {
							observer.unobserve(entry.target as HTMLElement);
						}
					}, delay * 1000);
				}
			});
		},
		{
			threshold,
			rootMargin: `${offset}px 0px`
		}
	);

	observer.observe(node);

	return {
		destroy() {
			clearTimeout(timeoutId);
			observer.disconnect();
		},
		update(newOptions: RepeatOnScrollOptions) {
			observer.disconnect();
			return repeatOnScroll(node, newOptions);
		}
	};
};

export const scrollProgress: Action<HTMLElement, {
	property?: 'opacity' | 'transform' | 'filter' | 'backgroundColor';
	startValue?: string;
	endValue?: string;
	startThreshold?: number;
	endThreshold?: number;
}> = (node, options = {}) => {
	if (!browser) return {};

	const {
		property = 'opacity',
		startValue = '0',
		endValue = '1',
		startThreshold = 0.2,
		endThreshold = 0.8
	} = options;

	const handleScroll = () => {
		const rect = node.getBoundingClientRect();
		const viewportHeight = window.innerHeight;
		const scrollY = window.scrollY;

		const elementTop = rect.top + scrollY;
		const elementBottom = rect.bottom + scrollY;

		const viewportTop = scrollY;
		const viewportBottom = scrollY + viewportHeight;

		const startPoint = viewportBottom - (viewportHeight * startThreshold);
		const endPoint = viewportTop + (viewportHeight * endThreshold);

		let progress = 0;
		if (elementTop < startPoint && elementBottom > endPoint) {
			const totalDistance = startPoint - endPoint;
			const elementDistance = startPoint - elementTop;
			progress = Math.max(0, Math.min(1, elementDistance / totalDistance));
		} else if (elementTop >= startPoint) {
			progress = 0;
		} else if (elementBottom <= endPoint) {
			progress = 1;
		}

		if (property === 'opacity') {
			node.style.opacity = `calc(${startValue} + (${endValue} - ${startValue}) * ${progress})`;
		} else if (property === 'transform') {
			const translateY = 20 * (1 - progress);
			node.style.transform = `translateY(${translateY}px)`;
		} else if (property === 'filter') {
			const blur = 10 * (1 - progress);
			node.style.filter = `blur(${blur}px)`;
		}
	};

	window.addEventListener('scroll', handleScroll, { passive: true });
	handleScroll();

	return {
		destroy() {
			window.removeEventListener('scroll', handleScroll);
		}
	};
};
