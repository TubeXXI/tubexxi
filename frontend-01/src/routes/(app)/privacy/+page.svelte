<script lang="ts">
	import { MetaTags } from 'svelte-meta-tags';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);
</script>

<MetaTags {...metaTags} />

<data.component user={data.user} setting={data.settings} lang={data.lang} />
