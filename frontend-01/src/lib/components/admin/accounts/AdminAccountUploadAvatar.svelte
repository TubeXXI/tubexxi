<script lang="ts" module>
	type UploadedFile = {
		name: string;
		type: string;
		size: number;
		uploadedAt: number;
		url: Promise<string>;
	};
</script>

<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { onDestroy } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';
	import { Button, buttonVariants } from '@/components/ui/button';
	import * as Dialog from '@/components/ui/dialog';
	import { CameraIcon, XIcon } from '@lucide/svelte';
		import * as FileDropZone from '$lib/components/ui-extras/file-drop-zone';
	import {
		displaySize,
		MEGABYTE,
	} from '@/components/ui-extras/file-drop-zone';
	import { Progress } from '@/components/ui/progress';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from '@/stores';
	import { sleep } from '@/utils/sleep.js';

	let {
		user
	}: {
		user?: User | null;
	} = $props();

	const onUpload: FileDropZone.FileDropZoneRootProps['onUpload'] = async (files) => {
		await Promise.allSettled(files.map((file) => uploadFile(file)));
	};
	const onFileRejected: FileDropZone.FileDropZoneRootProps['onFileRejected'] = async ({ reason, file }) => {
		toast.error(`${file.name} failed to upload! ${reason}`);
	};
	const uploadFile = async (file: File) => {
		if (files.find((f) => f.name === file.name)) return;
		const urlPromise = new Promise<string>((resolve) => {
			sleep(1000).then(() => resolve(URL.createObjectURL(file)));
		});

		files.push({
			name: `${new Date().getTime()}_${user?.full_name?.replaceAll(/\s+/g, '').toLowerCase() || 'user'}`,
			type: file.type,
			size: file.size,
			uploadedAt: Date.now(),
			url: urlPromise
		});
		avatarFile = file;
		await urlPromise;
	};

	let open = $state(false);
	let files = $state<UploadedFile[]>([]);
	let date = new SvelteDate();
	let avatarFile = $state<File | null>(null);
	let isUploading = $state(false);


	async function uploadToServer() {
		if (!avatarFile) return;
		try {
			isUploading = true;
			const formData = new FormData();
			formData.append('file', avatarFile);

			// Corrected endpoint: /accounts (handled by src/routes/(admin)/accounts/+server.ts)
			const response = await fetch('/accounts', {
				method: 'POST',
				body: formData
			});
			const data = await response.json();

			if (!response.ok || !data.success) {
				toast.error(data.message || 'Failed to upload avatar!');
				files = [];
				return;
			}
			if (data.success) {
				toast.success(data.message || 'Avatar uploaded successfully!');
				await invalidateAll();
				files = [];
			}
		} catch (error) {
			toast.error(
				(error instanceof Error ? error.message : 'Failed to upload avatar!')
		)
		} finally {
			isUploading = false;
			await invalidateAll();
			open = false;
			files = [];
			avatarFile = null;
		}
	}

	onDestroy(async () => {
		for (const file of files) {
			URL.revokeObjectURL(await file.url);
		}
	});

	$effect(() => {
		const interval = setInterval(() => {
			date.setTime(Date.now());
		}, 10);
		return () => {
			clearInterval(interval);
		};
	});
</script>

<Dialog.Root bind:open onOpenChange={(val) => val && invalidateAll()}>
	<Dialog.Trigger
		class={buttonVariants({
			variant: 'outline',
			size: 'icon',
			className: 'absolute right-0 bottom-0 h-8 w-8 cursor-pointer rounded-full'
		})}
	>
		<CameraIcon />
	</Dialog.Trigger>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Upload Avatar</Dialog.Title>
			<Dialog.Description>
				Upload a new avatar image for the user.
			</Dialog.Description>
		</Dialog.Header>
		<div class="flex w-full flex-col gap-2 p-6">
			{#if isUploading}
				<Empty.Root class="w-full">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Spinner />
						</Empty.Media>
						<Empty.Title>
							Please wait
						</Empty.Title>
						<Empty.Description>
							Processing request...
						</Empty.Description>
					</Empty.Header>
				</Empty.Root>
			{:else}
					<FileDropZone.Root
					{onUpload}
					{onFileRejected}
					maxFileSize={5 * MEGABYTE}
					fileCount={files.length}
					accept="image/*"
					maxFiles={1}
					disabled={files.length > 0 || isUploading}
				>
				<FileDropZone.Trigger />
						</FileDropZone.Root>
				<div class="flex flex-col gap-2">
					{#each files as file, i (file.name)}
						<div class="flex place-items-center justify-between gap-2">
							<div class="flex place-items-center gap-2">
								{#await file.url then src}
									<div class="relative size-9 overflow-clip">
										<img
											{src}
											alt={file.name}
											class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 overflow-clip"
										/>
									</div>
								{/await}
								<div class="flex flex-col">
									<span>{file.name}</span>
									<span class="text-xs text-muted-foreground">{displaySize(file.size)}</span>
								</div>
							</div>
							{#await file.url}
								<Progress
									class="h-2 w-full grow"
									value={((date.getTime() - file.uploadedAt) / 1000) * 100}
									max={100}
								/>
							{:then url}
								<Button
									variant="outline"
									size="icon"
									onclick={() => {
										URL.revokeObjectURL(url);
										files = [...files.slice(0, i), ...files.slice(i + 1)];
									}}
								>
									<XIcon />
								</Button>
							{/await}
						</div>
					{/each}
				</div>
			{/if}
		</div>
		<Dialog.Footer>
			<Button type="button" onclick={uploadToServer}>
				Upload
			</Button>
			<Dialog.Close>
				<Button variant="destructive" size="default">
					Close
				</Button>
			</Dialog.Close>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
