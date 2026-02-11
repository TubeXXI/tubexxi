import { writable, derived, get } from 'svelte/store';
import { page } from '$app/stores';

export interface PopupAdState {
	open_popup_ad: boolean;
	open_socialbar_ad: boolean;
}

export interface BannerAdsterraData {
	key: string;
	format: string;
	height: number;
	width: number;
	params: Record<string, any>;
	src: string;
}

export interface NativeAdsterraData {
	id: string;
	src: string;
	dataCfasync: boolean;
}

const createAdsStore = () => {
	const adsData = writable<SettingMonetize | null>(null);
	const popupAd = writable<PopupAdState>({
		open_popup_ad: false,
		open_socialbar_ad: false
	});

	const isEnableMonetize = derived(adsData, $ads => $ads?.enable_monetize ?? false);
	const typeMonetization = derived(adsData, $ads => $ads?.type_monetize ?? 'adsterra');
	const isEnablePopupAd = derived(
		adsData,
		$ads => ($ads?.enable_popup_ad && $ads?.popup_ad_code !== undefined) ?? false
	);
	const isEnableSocialbarAd = derived(
		adsData,
		$ads => ($ads?.enable_socialbar_ad && $ads?.socialbar_ad_code !== undefined) ?? false
	);

	const excludeRoutes = ['/about', '/faq', '/contact', '/privacy', '/terms'];
	const isExcludeRoute = derived(
		page,
		$page => excludeRoutes.includes($page.url.pathname)
	);

	let visibilityHandler: (() => void) | null = null;

	const handleVisibilityChange = () => {
		if (document.visibilityState !== 'visible') return;

		const shouldOpen = !get(isExcludeRoute) &&
			get(isEnableMonetize) &&
			(get(isEnablePopupAd) || get(isEnableSocialbarAd));

		if (shouldOpen) {
			popupAd.set({
				open_popup_ad: get(isEnablePopupAd),
				open_socialbar_ad: get(isEnableSocialbarAd)
			});
		}
	};

	const initAutoReopen = () => {
		if (typeof document === 'undefined') return;

		cleanupAutoReopen();
		visibilityHandler = handleVisibilityChange;
		document.addEventListener('visibilitychange', visibilityHandler);
	};

	const cleanupAutoReopen = () => {
		if (visibilityHandler) {
			document.removeEventListener('visibilitychange', visibilityHandler);
			visibilityHandler = null;
		}
	};

	const splitBannerValue = (banner?: string | null): BannerAdsterraData | null => {
		if (!banner) return null;

		const result: BannerAdsterraData = {
			key: '',
			format: '',
			height: 0,
			width: 0,
			params: {},
			src: ''
		};

		const keyMatch = banner.match(/'key'\s*:\s*'([^']+)'/);
		if (keyMatch) result.key = keyMatch[1];

		const formatMatch = banner.match(/'format'\s*:\s*'([^']+)'/);
		if (formatMatch) result.format = formatMatch[1];

		const heightMatch = banner.match(/'height'\s*:\s*(\d+)/);
		if (heightMatch) result.height = parseInt(heightMatch[1]);

		const widthMatch = banner.match(/'width'\s*:\s*(\d+)/);
		if (widthMatch) result.width = parseInt(widthMatch[1]);

		const srcMatch = banner.match(/src=["']([^"']+)["']/);
		if (srcMatch) result.src = srcMatch[1];

		return result;
	};

	const splitNativeAdsValue = (ad?: string | null): NativeAdsterraData | null => {
		if (!ad) return null;

		const result: NativeAdsterraData = {
			id: '',
			src: '',
			dataCfasync: false
		};

		const idMatch = ad.match(/id=["']([^"']+)["']/);
		if (idMatch) result.id = idMatch[1];

		const srcMatch = ad.match(/src=["']([^"']+)["']/);
		if (srcMatch) result.src = srcMatch[1];

		const dataCfasyncMatch = ad.match(/data-cfasync=["']([^"']+)["']/);
		if (dataCfasyncMatch) {
			result.dataCfasync = dataCfasyncMatch[1] === 'true';
		}

		return result;
	};

	const waitForDOMReady = (): Promise<void> => {
		return new Promise((resolve) => {
			if (typeof document === 'undefined') {
				resolve();
				return;
			}

			if (document.readyState === 'complete') {
				resolve();
			} else {
				const handler = () => {
					if (document.readyState === 'complete') {
						document.removeEventListener('readystatechange', handler);
						resolve();
					}
				};

				document.addEventListener('readystatechange', handler);
				// Fallback timeout
				setTimeout(resolve, 500);
			}
		});
	};

	// Initialize popup state based on ads data
	adsData.subscribe(($ads) => {
		if (!$ads) return;

		popupAd.set({
			open_popup_ad: ($ads.enable_popup_ad && $ads.popup_ad_code !== undefined) ?? false,
			open_socialbar_ad: ($ads.enable_socialbar_ad && $ads.socialbar_ad_code !== undefined) ?? false
		});
	});

	return {
		// Main store
		ads: adsData,
		popupAd,

		// Derived states
		isEnableMonetize,
		typeMonetization,
		isEnablePopupAd,
		isEnableSocialbarAd,
		isExcludeRoute,

		// Methods
		initAutoReopen,
		cleanupAutoReopen,
		splitBannerValue,
		splitNativeAdsValue,
		waitForDOMReady,

		// Convenience getters
		getAutoAdCode: () => {
			const data = get(adsData);
			return data?.auto_ad_code ?? null;
		},
		getBannerHorizonatlCode: () => {
			const data = get(adsData);
			return data?.banner_horizontal_ad_code ?? null;
		},
		getBannerRectangleCode: () => {
			const data = get(adsData);
			return data?.banner_rectangle_ad_code ?? null;
		},
		getBannerVerticalCode: () => {
			const data = get(adsData);
			return data?.banner_vertical_ad_code ?? null;
		},
		getPopupCode: () => {
			const data = get(adsData);
			return data?.popup_ad_code ?? null;
		},
		getSocialbarCode: () => {
			const data = get(adsData);
			return data?.socialbar_ad_code ?? null;
		},
		getNativeAdsCode: () => {
			const data = get(adsData);
			return data?.native_ad_code ?? null;
		}
	};
};

export const adsStore = createAdsStore();
