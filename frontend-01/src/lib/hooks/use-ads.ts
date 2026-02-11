import { getContext, setContext } from 'svelte';
import { adsStore } from "@/stores";

const ADS_CONTEXT_KEY = Symbol('ads');

export function useAds() {
	const existing = getContext<ReturnType<typeof createAds>>(ADS_CONTEXT_KEY);
	if (existing) return existing;

	const ads = createAds();
	setContext(ADS_CONTEXT_KEY, ads);
	return ads;
}

function createAds() {
	const {
		ads,
		popupAd,
		isEnableMonetize,
		typeMonetization,
		isEnablePopupAd,
		isEnableSocialbarAd,
		isExcludeRoute,
		initAutoReopen,
		cleanupAutoReopen,
		splitBannerValue,
		splitNativeAdsValue,
		waitForDOMReady,
		getAutoAdCode,
		getBannerHorizonatlCode,
		getBannerRectangleCode,
		getBannerVerticalCode,
		getPopupCode,
		getSocialbarCode,
		getNativeAdsCode
	} = adsStore;

	return {
		// Stores
		ads,
		popupAd,

		// Computed properties
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

		// Getters
		getAutoAdCode,
		getBannerHorizonatlCode,
		getBannerRectangleCode,
		getBannerVerticalCode,
		getPopupCode,
		getSocialbarCode,
		getNativeAdsCode,


		// Convenience
		get autoAdData() {
			const code = getAutoAdCode();
			return splitBannerValue(code);
		},
		get bannerRectangleData() {
			const code = getBannerRectangleCode();
			return splitBannerValue(code);
		},

		get bannerHorizontalData() {
			const code = getBannerHorizonatlCode();
			return splitBannerValue(code);
		},

		get bannerVerticalData() {
			const code = getBannerVerticalCode();
			return splitBannerValue(code);
		},

		get popupData() {
			const code = getPopupCode();
			return splitNativeAdsValue(code);
		},

		get socialbarData() {
			const code = getSocialbarCode();
			return splitNativeAdsValue(code);
		},

		get nativeAdsData() {
			const code = getNativeAdsCode();
			return splitNativeAdsValue(code);
		}
	};
}
