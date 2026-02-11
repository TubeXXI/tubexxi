
import { browser } from '$app/environment';
import type { Action } from 'svelte/action';

interface ScrollActionOptions {
	threshold?: number;
	delay?: number;
	duration?: number;
	offset?: number;
	once?: boolean;
}

export const fadeIn: Action<HTMLElement, ScrollActionOptions> = (node, options = {}) => {
	if (!browser) return {};

	const {
		threshold = 0.1,
		delay = 0,
		duration = 0.6,
		offset = 0,
		once = true
	} = options;

	Object.assign(node.style, {
		opacity: '0',
		transform: 'translateY(20px)',
		transition: `opacity ${duration}s ease ${delay}s, transform ${duration}s ease ${delay}s`,
		willChange: 'opacity, transform'
	});

	const observer = new IntersectionObserver(
		(entries) => {
			entries.forEach(entry => {
				if (entry.isIntersecting) {
					Object.assign((entry.target as HTMLElement).style, {
						opacity: '1',
						transform: 'translateY(0)'
					});

					if (once) observer.unobserve(entry.target);
				} else if (!once) {
					Object.assign((entry.target as HTMLElement).style, {
						opacity: '0',
						transform: 'translateY(20px)'
					});
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
			observer.unobserve(node);
			observer.disconnect();
		},
		update(newOptions: ScrollActionOptions) {
			observer.unobserve(node);
			observer.disconnect();

			return fadeIn(node, newOptions);
		}
	};
};

export const slideInLeft: Action<HTMLElement, ScrollActionOptions> = (node, options = {}) => {
	if (!browser) return {};

	const { threshold = 0.1, delay = 0, duration = 0.5 } = options;

	Object.assign(node.style, {
		opacity: '0',
		transform: 'translateX(-50px)',
		transition: `opacity ${duration}s ease ${delay}s, transform ${duration}s ease ${delay}s`
	});

	const observer = new IntersectionObserver(
		(entries) => {
			entries.forEach(entry => {
				if (entry.isIntersecting) {
					Object.assign((entry.target as HTMLElement).style, {
						opacity: '1',
						transform: 'translateX(0)'
					});
					observer.unobserve(entry.target);
				}
			});
		},
		{ threshold }
	);

	observer.observe(node);

	return {
		destroy: () => observer.disconnect(),
		update: (newOptions: ScrollActionOptions) => {
			observer.unobserve(node);
			observer.disconnect();
			return slideInLeft(node, newOptions);
		}
	};
};
