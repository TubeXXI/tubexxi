<script lang="ts" module>
	interface NavigationMenuItem {
		label: string;
		name: string;
		description?: string;
		icon: string;
		to: string;
		exact?: boolean;
		target?: string;
	}
</script>

<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { ClassValue } from 'svelte/elements';
	import { cn } from '@/utils';
	import { AdminSettingSidenav } from '@/components/admin';
	import { Separator } from '@/components/ui/separator';
	import { localizeHref } from '@/paraglide/runtime';

	let {
		children,
		class: className,
		fixed,
		fluid,
		title,
		description
	}: {
		children: Snippet;
		class?: ClassValue;
		fixed?: boolean;
		fluid?: boolean;
		title?: string;
		description?: string;
	} = $props();

	const links: NavigationMenuItem[] = [
		{
			label: 'Web',
			name: 'Web',
			description: 'Web settings',
			icon: 'mdi:web',
			to: localizeHref('/settings/web'),
			exact: true
		},
		{
			label: 'Email',
			name: 'Email',
			description: 'Email settings',
			icon: 'lucide-mail',
			to: localizeHref('/settings/email')
		},
		{
			label: 'System',
			name: 'System',
			description: 'System settings',
			icon: 'solar:settings-linear',
			to: localizeHref('/settings/system')
		},
		{
			label: 'Monetization',
			name: 'Monetization',
			description: 'Monetization settings',
			icon: 'streamline:dollar-coin-remix',
			to: localizeHref('/settings/monetization')
		},
		{
			label: 'Ads.txt',
			name: 'Ads.txt',
			description: 'Ads.txt settings',
			icon: 'lsicon:file-txt-filled',
			to: localizeHref('/settings/ads.txt')
		},
		{
			label: 'Robot.txt',
			name: 'Robot.txt',
			description: 'Robot.txt settings',
			icon: 'lsicon:file-txt-filled',
			to: localizeHref('/settings/robot.txt')
		}
	];
</script>

<div class={cn('@container/main flex grow flex-col gap-4 md:gap-6', className)}>
	<div class="flex-none px-4 py-4 sm:px-6">
		<div class="space-y-1">
			<h2 class="text-2xl font-bold tracking-tight md:text-3xl">{title}</h2>
			<p class="text-sm text-muted-foreground">{description}</p>
		</div>
	</div>
	<Separator />
	<div class="flex flex-1 flex-col space-y-2 md:space-y-2 lg:flex-row lg:space-y-0 lg:space-x-12">
		<aside class="top-0 lg:sticky lg:w-1/5">
			<AdminSettingSidenav {links} />
		</aside>
		<Separator orientation="vertical" class="hidden lg:block" />
		<div class="flex w-full p-1">
			<div class="flex flex-1 flex-col">
				<div class="flex-none">
					<h3 class="text-lg font-medium">{title}</h3>
					<p class="text-sm text-muted-foreground">
						{description}
					</p>
				</div>
				<Separator class="my-4 flex-none" />
				<div
					class="faded-bottom scrollbar-thin scrollbar-thumb-foreground scrollbar-track-accent max-h-[calc(100vh-300px)] w-full overflow-y-auto scroll-smooth pe-4 pb-12"
				>
					<div class="-mx-1 w-full px-1.5">
						{@render children()}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
