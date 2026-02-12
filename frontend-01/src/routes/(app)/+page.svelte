<script lang="ts">
	import { page } from '$app/state';
	import { goto, invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { updateUrlParams, createQueryManager } from '@/stores/query.js';
	import { handleSubmitLoading } from '@/stores';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	const queryManager = createQueryManager();
	let query = $state(queryManager.parse(page.url));

	$effect(() => {
		query = queryManager.parse(page.url);
	});

	async function updateQuery(updates: Partial<typeof query>, resetPage = false) {
		handleSubmitLoading(true);
		await updateUrlParams(goto, page.url, updates, {
			resetPage,
			replaceState: true,
			invalidateAll: true
		});
		handleSubmitLoading(false);
	}
</script>

<MetaTags {...metaTags} />

<data.component user={data.user} setting={data.settings}>
	<div>Content Group</div>
</data.component>
