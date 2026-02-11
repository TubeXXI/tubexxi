<script lang="ts">
	import { Separator } from '$lib/components/ui/separator';
	import * as Sidebar from '$lib/components/ui/sidebar';
	import { LightSwitch } from '@/components/ui-extras/light-switch';
	import { cn } from '@/utils';

	let { page }: { page?: string } = $props();
</script>

<header
	class={cn(
		'flex h-12 shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12',
		'sticky top-0 z-50 overflow-hidden rounded-t-[inherit] bg-background/50 backdrop-blur-md'
	)}
>
	<div class="flex w-full items-center justify-between px-4 lg:px-6">
		<div class="flex items-center gap-1 lg:gap-2">
			<Sidebar.Trigger class="-ml-1" />
			<Separator orientation="vertical" class="mx-2 data-[orientation=vertical]:h-4" />
			<h1 class="text-base font-medium">{page}</h1>
		</div>

		<div class="flex items-center gap-2">
			<LightSwitch />
		</div>
	</div>
</header>
