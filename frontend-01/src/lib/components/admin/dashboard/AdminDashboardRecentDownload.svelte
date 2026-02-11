<script lang="ts">
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Avatar from '$lib/components/ui/avatar/index.js';
	import { formatTimeAgo } from '@/utils/time';

	let { downloads }: { downloads: Download[] } = $props();
</script>

<Table.Root>
	<Table.Header>
		<Table.Row>
			<Table.Head>Video Title</Table.Head>
			<Table.Head>Original URL</Table.Head>
			<Table.Head>Status</Table.Head>
			<Table.Head>Ip Address</Table.Head>
			<Table.Head>Download Date</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body>
		{#each downloads as download}
			<Table.Row>
				<Table.Cell class="flex w-40 items-center gap-2">
					<Avatar.Root>
						<Avatar.Image src={download.thumbnail_url} />
						<Avatar.Fallback>
							{download.title?.slice(0, 2).toUpperCase()}
						</Avatar.Fallback>
					</Avatar.Root>
					{download.title}
				</Table.Cell>
				<Table.Cell>{download.original_url}</Table.Cell>
				<Table.Cell>{download.status}</Table.Cell>
				<Table.Cell>{download.ip_address}</Table.Cell>
				<Table.Cell>{formatTimeAgo(download.created_at)}</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
