<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Avatar from '@/components/ui/avatar';
	import * as Item from '$lib/components/ui/item/index.js';
	import Icon from '@iconify/svelte';
	import { localizeHref } from '$lib/paraglide/runtime';
	import { firebaseClient } from '@/client/firebase_client';
	import * as i18n from '@/paraglide/messages.js';

	let {
		user,
		setting
	}: {
		user?: User | null;
		setting?: SettingsValue | null;
	} = $props();

	async function handleLogout() {
		await fetch(`/api/auth/sign-out`, {
			method: 'POST',
			headers: {
				'X-Platform': 'web'
			},
			credentials: 'include'
		});
		await firebaseClient?.signOut();
		await goto(localizeHref(page.url.pathname));
		await invalidateAll();
	}
</script>

<div class="flex items-center gap-2">
	{#if user}
		<DropdownMenu.Root>
			<DropdownMenu.Trigger>
				{#snippet child({ props })}
					<Button variant="ghost" size="icon" {...props}>
						<Avatar.Root>
							<Avatar.Image src={user?.avatar_url ?? '/images/avatar.jpg'}></Avatar.Image>
							<Avatar.Fallback>{user?.full_name?.slice(0, 2).toUpperCase() || '?'}</Avatar.Fallback>
						</Avatar.Root>
					</Button>
				{/snippet}
			</DropdownMenu.Trigger>
			<DropdownMenu.Content
				class="mr-2 w-(--bits-dropdown-menu-anchor-width) min-w-64 space-y-3 rounded-lg"
				align="center"
			>
				<Item.Root variant="outline">
					{#snippet child({ props })}
						<div class="flex items-center gap-2">
							<Item.Media>
								<Avatar.Root>
									<Avatar.Image src={user?.avatar_url ?? '/images/avatar.jpg'}></Avatar.Image>
									<Avatar.Fallback
										>{user?.full_name?.slice(0, 2).toUpperCase() || '?'}</Avatar.Fallback
									>
								</Avatar.Root>
							</Item.Media>
							<Item.Content>
								<Item.Title class="text-sm font-medium">{user?.full_name || '?'}</Item.Title>
								<Item.Description
									class="line-clamp-1 text-xs text-neutral-500 dark:text-neutral-400"
								>
									{user?.email || '?'}
								</Item.Description>
							</Item.Content>
						</div>
					{/snippet}
				</Item.Root>
				<DropdownMenu.Separator class="mb-2" />
				<DropdownMenu.Item>
					{#snippet child()}
						<a
							href={localizeHref('/user/account')}
							class="flex items-center gap-2 ps-3 text-sm transition-all active:scale-95"
						>
							<Icon icon="solar:user-circle-linear" />
							{i18n.page_account_profile()}
						</a>
					{/snippet}
				</DropdownMenu.Item>
				{#if user.role?.level === 1 || user.role?.level === 2}
					<DropdownMenu.Item>
						{#snippet child()}
							<a
								href={localizeHref('/admin/dashboard')}
								class="flex items-center gap-2 ps-3 text-sm transition-all active:scale-95"
							>
								<Icon icon="material-symbols:admin-panel-settings-outline" />
								{i18n.admin()}
							</a>
						{/snippet}
					</DropdownMenu.Item>
				{/if}
				<DropdownMenu.Separator />
				<DropdownMenu.Item>
					{#snippet child()}
						<Button
							type="button"
							variant="ghost"
							class="text-sm text-red-600 dark:text-red-500"
							onclick={handleLogout}
						>
							<Icon icon="ic:outline-logout" />
							{i18n.log_out()}
						</Button>
					{/snippet}
				</DropdownMenu.Item>
			</DropdownMenu.Content>
		</DropdownMenu.Root>
	{:else}
		<Button href={localizeHref('/auth/login')} type="button" variant="outline" class="text-xs">
			{i18n.sign_in()}
		</Button>
		<Button href={localizeHref('/auth/register')} type="button" variant="outline" class="text-xs">
			{i18n.sign_up()}
		</Button>
	{/if}
</div>

<style scoped>
</style>
