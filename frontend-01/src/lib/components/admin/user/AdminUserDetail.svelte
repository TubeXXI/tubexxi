<script lang="ts">
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { Badge } from '@/components/ui/badge/index.js';
	import { Mail, Calendar, Clock, Info } from '@lucide/svelte';

	let {
		open = $bindable(false),
		user,
		onclose
	}: {
		open?: boolean;
		user?: User | null;
		onclose?: () => void;
	} = $props();

	const handleClose = () => {
		open = false;
		onclose?.();
	};

	function formatDate(dateString?: Date | string): string {
		if (!dateString) {
			return 'N/A';
		}
		if (dateString instanceof String) {
			return new Date(dateString).toLocaleString('en-US', {
				year: 'numeric',
				month: 'long',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		}
		return dateString.toLocaleString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content class="w-full">
		<Card.Root class="-my-4 mt-5 w-full min-w-md">
			<Card.Header>
				<Card.Title>User detail</Card.Title>
				<Card.Description></Card.Description>
			</Card.Header>
			<div class="flex flex-col items-start gap-6 px-3 md:flex-row md:items-center lg:px-6">
				<div class="relative">
					<Avatar.Root class="h-24 w-24">
						<Avatar.Image src={user?.avatar_url || ''} alt={user?.full_name || 'User Avatar'} />
						<Avatar.Fallback
							>{user?.full_name?.split(' ')[0]?.slice(0, 2).toUpperCase() ||
								'User'}</Avatar.Fallback
						>
					</Avatar.Root>
				</div>
				<div class="flex-1 space-y-2">
					<div class="flex flex-col gap-2 md:flex-row md:items-center">
						<h1 class="text-2xl font-bold text-neutral-900 capitalize dark:text-white">
							{user?.full_name || 'N/A'}
						</h1>
					</div>
					<div class="mt-1 flex items-center gap-2">
						<Badge variant="default" class="px-4 text-xs font-semibold">
							{user?.role?.name || 'N/A'}
						</Badge>
					</div>
					<div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger
									class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
								>
									<Mail class="size-4" />
									{user?.email || 'N/A'}
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>Email address</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger
									class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
								>
									<Calendar class="size-4" />
									{user?.created_at ? formatDate(user?.created_at) : 'N/A'}
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>
										Member since {user?.created_at ? formatDate(user?.created_at) : 'N/A'}
									</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
						<Tooltip.Provider>
							<Tooltip.Root>
								<Tooltip.Trigger
									class="flex cursor-pointer items-center gap-1 text-sm font-semibold text-muted-foreground capitalize"
								>
									<Clock class="size-4" />
									{user?.last_login_at ? formatDate(user?.last_login_at) : 'N/A'}
								</Tooltip.Trigger>
								<Tooltip.Content>
									<p>
										Last login on {user?.last_login_at ? formatDate(user?.last_login_at) : 'N/A'}
									</p>
								</Tooltip.Content>
							</Tooltip.Root>
						</Tooltip.Provider>
					</div>
				</div>
			</div>
			<Card.Footer class="flex-col gap-2">
				<Button type="button" class="w-full" onclick={handleClose}>Close</Button>
			</Card.Footer>
		</Card.Root>
	</Dialog.Content>
</Dialog.Root>
