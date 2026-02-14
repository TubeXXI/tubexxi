<script lang="ts">
	import { DefaultSidebarLayout } from '$lib/components/themes/index.js';

	let {
		user,
		setting,
		lang = 'en'
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
		lang?: string;
	} = $props();
</script>

<DefaultSidebarLayout {user} {setting} {lang}>
	<div class="scrollbar-primary flex min-h-[calc(100vh-160px)] flex-col">
		<div>Content</div>
	</div>
</DefaultSidebarLayout>
