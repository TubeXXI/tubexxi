<script lang="ts">
	import { UserLayout } from '$lib/components/user/index.js';

	let { data, children } = $props();
</script>

<UserLayout setting={data.settings}>
	{@render children()}
</UserLayout>
