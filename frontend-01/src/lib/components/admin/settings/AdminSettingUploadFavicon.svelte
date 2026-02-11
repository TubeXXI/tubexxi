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
	import { Button } from '@/components/ui/button';
	import * as FileDropZone from '$lib/components/ui-extras/file-drop-zone';
	import {
		displaySize,
		MEGABYTE,
	} from '@/components/ui-extras/file-drop-zone';
	import { Progress } from '@/components/ui/progress';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { toast } from '@/stores';
	import { sleep } from '@/utils/sleep.js';
	import { XIcon } from '@lucide/svelte';


	let {
		web
	}: {
		web?: SettingWeb | null;
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
			name: `${new Date().getTime()}_${web?.site_name?.replaceAll(/\s+/g, '').toLowerCase() || 'logo'}`,
			type: file.type,
			size: file.size,
			uploadedAt: Date.now(),
			url: urlPromise
		});
		faviconFile = file;
		await urlPromise;
	};

	let files = $state<UploadedFile[]>([]);
	let date = new SvelteDate();
	let faviconFile = $state<File | null>(null);
	let isUploading = $state(false);


	async function uploadToServer() {
		if (!faviconFile) return;
		try {
			isUploading = true;
			const formData = new FormData();
			formData.append('file', faviconFile);

			const response = await fetch('/settings/web/favicon', {
				method: 'POST',
				body: formData
			});
			const data = await response.json();

			if (!response.ok || !data.success) {
				toast.error(data.message || 'Failed to upload favicon!');
				files = [];
				return;
			}
			if (data.success) {
				toast.success(data.message || 'Favicon uploaded successfully!');
				await invalidateAll();
				files = [];
			}
		} catch (error) {
			toast.error(
				(error instanceof Error ? error.message : 'Failed to upload favicon!')
			);
		} finally {
			isUploading = false;
			await invalidateAll();
			files = [];
			faviconFile = null;
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

<Field.Group>
	<Field.Set>
		<Field.Legend>Upload Favicon</Field.Legend>
		<Field.Description
			>Upload a favicon image to be used on the site.</Field.Description
		>
		<Field.Group>
			<Field.Content>
				<div class="flex w-full flex-col gap-2">
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
							id="site_favicon_image"
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
			</Field.Content>
		</Field.Group>
	</Field.Set>
	<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
		<Button type="button" class="w-full" disabled={isUploading} onclick={uploadToServer}>
			{#if isUploading}
				<Spinner class="mr-2 size-5" />
			{/if}
			{isUploading ? 'Submitting' : 'Upload Favicon'}
		</Button>
	</Field.Field>
</Field.Group>
