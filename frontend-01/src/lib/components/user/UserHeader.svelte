<script lang="ts">
	import { goto } from '$app/navigation';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { UserUploadAvatar } from '$lib/components/user';
	import Icon from '@iconify/svelte';
	import { formatTimeAgo, formatToYYYYMMDD } from '@/utils/time';
	import { localizeHref } from '@/paraglide/runtime';
	import * as i18n from '@/paraglide/messages.js';

	let { user }: { user?: User | null } = $props();

	async function handleLogout() {
		await fetch(`/user/account/signout`, {
			method: 'POST',
			headers: {
				'X-Platform': 'web'
			},
			credentials: 'include'
		});

		await goto('/auth/login');
	}
</script>

<Card.Root class="bg-white shadow-lg backdrop-blur-lg dark:bg-neutral-950">
	<Card.Content class="flex flex-col items-start gap-6 md:flex-row md:items-center">
		<div class="relative">
			<Avatar.Root class="h-24 w-24">
				<Avatar.Image src={user?.avatar_url || ''} alt={user?.full_name || 'User Avatar'} />
				<Avatar.Fallback
					>{user?.full_name?.split(' ')[0]?.slice(0, 2).toUpperCase() || 'User'}</Avatar.Fallback
				>
			</Avatar.Root>
			<UserUploadAvatar {user} />
		</div>
		<div class="flex-1 space-y-2">
			<div class="flex flex-col gap-2 md:flex-row md:items-center">
				<h1 class="text-2xl font-bold">{user?.full_name}</h1>
				<Badge variant="default" class="text-xs uppercase">
					{user?.role?.name}
				</Badge>
			</div>
			<p class="text-muted-foreground">
				{user?.email}
			</p>
			<Separator />
			<div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger class="cursor-pointer">
							<div class="flex items-center gap-1 text-sm">
								<Icon icon="mdi:clock-outline" />
								{user?.last_login_at ? formatTimeAgo(user?.last_login_at as string) : 'Never'}
							</div>
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>{i18n.last_login_at()}</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
				<Tooltip.Provider>
					<Tooltip.Root>
						<Tooltip.Trigger class="cursor-pointer">
							<div class="flex items-center gap-1 text-sm">
								<Icon icon="mdi:calendar" />
								{user?.last_login_at ? formatToYYYYMMDD(new Date(user?.last_login_at)) : 'Never'}
							</div>
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>{i18n.join_at()}</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			</div>
		</div>
		<div class="flex items-center gap-2">
			<Button href={localizeHref('/')} variant="default">
				<Icon icon="mdi:home-outline" />
			</Button>
			{#if user?.role?.name === 'admin' || user?.role?.level === 1 || user?.role?.level === 2}
				<Button href={localizeHref('/admin/dashboard')} variant="default">
					<Icon icon="mdi:tools" />
				</Button>
			{/if}
			<Button variant="destructive" onclick={handleLogout}>
				<Icon icon="mingcute:exit-fill" />
			</Button>
		</div>
	</Card.Content>
</Card.Root>
