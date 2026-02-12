<script lang="ts">
	import { onMount } from 'svelte';
	import { AppModalProvider } from '@/components/modal';

	let { data, children } = $props();

	let systemSetting = $derived(data?.settings?.SYSTEM);
	let monetizeSetting = $derived(data?.settings?.MONETIZE);

	onMount(() => {
		const histatstCode = systemSetting?.histats_tracking_code;
		if (!histatstCode) {
			return;
		}

		if (!document.getElementById('histats-inline')) {
			const inline = document.createElement('script');
			inline.id = 'histats-inline';
			inline.type = 'text/javascript';
			inline.text = `
      window._Hasync = window._Hasync || [];
      window._Hasync.push(['Histats.start', '1,${histatstCode},4,0,0,0,00010000']);
      window._Hasync.push(['Histats.fasi', '1']);
      window._Hasync.push(['Histats.track_hits', '']);
    `;
			document.head.appendChild(inline);
		}

		if (!document.getElementById('histats-external')) {
			const ext = document.createElement('script');
			ext.id = 'histats-external';
			ext.src = 'https://s10.histats.com/js15_as.js';
			ext.async = true;
			document.head.appendChild(ext);
		}
	});
</script>

<svelte:head>
	{#if systemSetting?.google_analytics_code}
		<script
			async
			src="https://www.googletagmanager.com/gtag/js?id={systemSetting?.google_analytics_code}"
		></script>
		<script>
			window.dataLayer = window.dataLayer || [];
			function gtag() {
				dataLayer.push(arguments);
			}
			gtag('js', new Date());
			gtag('config', systemSetting?.google_analytics_code, {
				custom_map: { dimension1: 'lang' }
			});
			gtag('event', 'page_view', {
				lang: lang || 'en'
			});
		</script>
	{/if}
	{#if monetizeSetting?.type_monetize === 'adsense' && monetizeSetting?.publisher_id}
		<script
			async
			src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=${monetizeSetting?.publisher_id}"
			crossorigin="anonymous"
		></script>
		<meta name="google-adsense-account" content={monetizeSetting?.publisher_id} />
	{/if}
</svelte:head>

<AppModalProvider>
	{@render children()}
</AppModalProvider>
