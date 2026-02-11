import { writable, derived, get } from 'svelte/store';
import { translate, type SingleResponse, type BatchResponse } from '@siamf/google-translate';
import { langStore } from "./page_data";

interface TranslateState {
	isLoading: boolean;
	error: string | null;
	lastTranslation: SingleResponse | null;
	lastBatchTranslation: BatchResponse | null;
	cache: Map<string, SingleResponse>;
	batchCache: Map<string, BatchResponse>;
}

interface TranslateOptions {
	useCache?: boolean;
	targetLang?: string;
	fallbackOnError?: boolean;
}

function createTranslateStore(): {
	subscribe: (run: (value: TranslateState) => void, invalidate?: (value?: TranslateState) => void) => () => void;
	singleTranslate: (text: string, options?: TranslateOptions) => Promise<SingleResponse>;
	batchTranslate: (texts: string[], options?: TranslateOptions) => Promise<BatchResponse>;
	autoTranslate: (text: string, options?: Omit<TranslateOptions, 'targetLang'>) => Promise<SingleResponse>;
	translateToCurrentLang: (text: string, options?: Omit<TranslateOptions, 'targetLang'>) => Promise<SingleResponse>;
	clearCache: () => void;
	clearCacheEntry: (text: string, targetLang: string) => void;
	reset: () => void;
	isLoading: import('svelte/store').Readable<boolean>;
	error: import('svelte/store').Readable<string | null>;
	lastTranslation: import('svelte/store').Readable<SingleResponse | null>;
	lastBatchTranslation: import('svelte/store').Readable<BatchResponse | null>;
	cacheSize: import('svelte/store').Readable<number>;
	batchCacheSize: import('svelte/store').Readable<number>;
	getState: () => TranslateState;
} {
	const baseStore = writable<TranslateState>({
		isLoading: false,
		error: null,
		lastTranslation: null,
		lastBatchTranslation: null,
		cache: new Map(),
		batchCache: new Map()
	});

	const { subscribe, set, update } = baseStore;

	const generateCacheKey = (text: string, targetLang: string): string => {
		return `${text}|${targetLang}`;
	};


	const singleTranslate = async (
		text: string,
		options: TranslateOptions = {}
	): Promise<SingleResponse> => {
		const {
			useCache = true,
			targetLang = get(langStore || 'auto') || 'auto',
			fallbackOnError = true
		} = options;

		update(state => ({ ...state, isLoading: true, error: null }));

		try {
			const cacheKey = generateCacheKey(text, targetLang);
			if (useCache) {
				const cached = get(baseStore).cache.get(cacheKey);
				if (cached) {
					update(state => ({
						...state,
						isLoading: false,
						lastTranslation: cached
					}));
					return cached;
				}
			}

			const result = await translate.single(text, targetLang);

			update(state => ({
				...state,
				isLoading: false,
				lastTranslation: result,
				cache: useCache ? new Map(state.cache).set(cacheKey, result) : state.cache
			}));

			return result;

		} catch (error) {
			console.error('Translation error:', error);
			const errorMessage = error instanceof Error ? error.message : 'Unknown translation error';

			update(state => ({
				...state,
				isLoading: false,
				error: errorMessage
			}));

			if (fallbackOnError) {
				try {
					const fallbackResult = await translate.single(errorMessage, 'auto');
					update(state => ({
						...state,
						lastTranslation: fallbackResult
					}));
					return fallbackResult;
				} catch (fallbackError) {
					throw fallbackError;
				}
			}

			throw error;
		}
	};

	const batchTranslate = async (
		texts: string[],
		options: TranslateOptions = {}
	): Promise<BatchResponse> => {
		const {
			useCache = true,
			targetLang = [get(langStore || 'auto') || 'auto'],
			fallbackOnError = true
		} = options;

		update(state => ({ ...state, isLoading: true, error: null }));

		try {
			const cacheKey = generateCacheKey(texts.join('|'), Array.isArray(targetLang) ? targetLang.join(',') : targetLang);
			if (useCache) {
				const cached = get(baseStore).batchCache.get(cacheKey);
				if (cached) {
					update(state => ({
						...state,
						isLoading: false,
						lastBatchTranslation: cached
					}));
					return cached;
				}
			}


			const result = await translate.batch(texts.join(',').trim(), Array.isArray(targetLang) ? targetLang : [targetLang]);


			update(state => ({
				...state,
				isLoading: false,
				lastBatchTranslation: result,
				batchCache: useCache ? new Map(state.batchCache).set(cacheKey, result) : state.batchCache
			}));

			return result;

		} catch (error) {
			console.error('Batch translation error:', error);
			const errorMessage = error instanceof Error ? error.message : 'Unknown batch translation error';

			update(state => ({
				...state,
				isLoading: false,
				error: errorMessage
			}));

			if (fallbackOnError) {
				try {
					const fallbackResult = await translate.batch(errorMessage, ['auto']);
					update(state => ({
						...state,
						lastBatchTranslation: fallbackResult
					}));
					return fallbackResult;
				} catch (fallbackError) {
					throw fallbackError;
				}
			}

			throw error;
		}
	};

	const autoTranslate = async (
		text: string,
		options: Omit<TranslateOptions, 'targetLang'> = {}
	): Promise<SingleResponse> => {
		return singleTranslate(text, { ...options, targetLang: 'auto' });
	};

	const translateToCurrentLang = async (
		text: string,
		options: Omit<TranslateOptions, 'targetLang'> = {}
	): Promise<SingleResponse> => {
		const currentLang = get(langStore || 'auto') || 'en';
		return singleTranslate(text, { ...options, targetLang: currentLang });
	};

	const clearCache = (): void => {
		update(state => ({
			...state,
			cache: new Map(),
			batchCache: new Map()
		}));
	};

	const clearCacheEntry = (text: string, targetLang: string): void => {
		const cacheKey = generateCacheKey(text, targetLang);
		update(state => {
			const newCache = new Map(state.cache);
			newCache.delete(cacheKey);
			return { ...state, cache: newCache };
		});
	};

	const reset = (): void => {
		set({
			isLoading: false,
			error: null,
			lastTranslation: null,
			lastBatchTranslation: null,
			cache: new Map(),
			batchCache: new Map()
		});
	};

	const isLoading = derived(baseStore, $store => $store.isLoading);
	const error = derived(baseStore, $store => $store.error);
	const lastTranslation = derived(baseStore, $store => $store.lastTranslation);
	const lastBatchTranslation = derived(baseStore, $store => $store.lastBatchTranslation);
	const cacheSize = derived(baseStore, $store => $store.cache.size);
	const batchCacheSize = derived(baseStore, $store => $store.batchCache.size);

	return {
		subscribe,
		// Actions
		singleTranslate,
		batchTranslate,
		autoTranslate,
		translateToCurrentLang,
		clearCache,
		clearCacheEntry,
		reset,
		// Derived stores (read-only)
		isLoading,
		error,
		lastTranslation,
		lastBatchTranslation,
		cacheSize,
		batchCacheSize,
		// Utility
		getState: () => get(baseStore)
	};
}

export const translateStore = createTranslateStore();
