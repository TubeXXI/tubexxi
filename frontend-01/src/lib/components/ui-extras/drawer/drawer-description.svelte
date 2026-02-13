<script lang="ts">
	import { Drawer as DrawerPrimitive } from 'vaul-svelte';
	import { cn } from '$lib/utils.js';

	let { class: className, ...restProps }: DrawerPrimitive.DescriptionProps = $props();
</script>

<DrawerPrimitive.Description
	data-slot="drawer-description"
	class={cn('text-sm text-muted-foreground', className)}
	{...restProps as any}
/>
