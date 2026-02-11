<script lang="ts">
import type { Snippet } from 'svelte';
import * as Sidebar from '$lib/components/ui/sidebar';
import { AdminSidebar, AdminSidebarHeader } from '@/components/admin';
import { cn } from '@/utils';

let {
	user,
	setting,
	page,
	children
}: {
	user?: User | null;
	setting?: SettingsValue | null;
	page?: string;
	children?: Snippet<[]>;
} = $props();
</script>

<Sidebar.Provider open={true}>
	<AdminSidebar variant="sidebar" user={user} setting={setting} />
	<Sidebar.Inset
		class={cn(
			'mx-auto! lg:max-w-full',
			'max-[113rem]:peer-data-[variant=inset]:mr-2! min-[101rem]:peer-data-[variant=inset]:peer-data-[state=collapsed]:mr-auto!'
		)}
	>
		<AdminSidebarHeader page={page} />
		<main class="h-full p-4 md:p-6">
			{@render children?.()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
