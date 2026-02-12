<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';

	let {
		user,
		setting
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
	} = $props();
</script>

<DefaultSidebarLayout {user} {setting}>
	<div
		class="scrollbar-thin scrollbar-thumb-foreground scrollbar-track-accent flex min-h-[calc(100vh-160px)] flex-col overflow-hidden overflow-y-auto scroll-smooth"
	>
		<div>Content</div>
	</div>
</DefaultSidebarLayout>
