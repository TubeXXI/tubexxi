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
	import { AdminSettingSidenav } from '@/components/admin';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Separator } from '@/components/ui/separator';
	import { cn } from '@/utils';
	import { localizeHref } from '@/paraglide/runtime';
	import * as i18n from '@/paraglide/messages.js';

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
			label: i18n.account_profile_setting(),
			name: 'profile',
			description: i18n.account_profile_setting_description(),
			icon: 'mdi:account',
			to: localizeHref('/user/account'),
			exact: true
		},
		{
			label: i18n.change_password(),
			name: 'password',
			description: i18n.change_password_description(),
			icon: 'mdi:lock',
			to: localizeHref('/user/account/password'),
			exact: true
		}
	];
</script>

<Card.Root>
	<Card.Content class="bg-white shadow-lg backdrop-blur-lg dark:bg-neutral-950">
		<div class={cn('@container/main flex grow flex-col gap-4 md:gap-6', className)}>
			<div class="flex-none px-4 py-4 sm:px-6">
				<div class="space-y-1">
					<h2 class="text-2xl font-bold tracking-tight md:text-3xl">
						{title}
					</h2>
					<p class="text-sm text-muted-foreground">{description}</p>
				</div>
			</div>
			<Separator />
			<div
				class="flex flex-1 flex-col space-y-2 md:space-y-2 lg:flex-row lg:space-y-0 lg:space-x-12"
			>
				<aside class="sticky top-0 lg:w-1/5">
					<AdminSettingSidenav {links} />
				</aside>
				<Separator orientation="vertical" class="hidden lg:flex" />
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
							class="faded-bottom scrollbar-thin max-h-[calc(100vh-300px)] w-full overflow-y-auto scroll-smooth pe-4 pb-12 scrollbar-thumb-foreground scrollbar-track-accent"
						>
							<div class="-mx-1 w-full px-1.5">
								{@render children()}
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</Card.Content>
</Card.Root>
