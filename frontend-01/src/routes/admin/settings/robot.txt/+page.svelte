<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { MetaTags } from 'svelte-meta-tags';
	import { superForm } from 'sveltekit-superforms';
	import { handleSubmitLoading } from '@/stores';
	import { AppAlertDialog } from '@/components/index.js';
	import { AdminSidebarLayout, AdminSettingLayout } from '@/components/admin';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { Textarea } from '@/components/ui/textarea/index.js';

	let { data } = $props();
	let metaTags = $derived(data.pageMetaTags);

	let successMessage = $state<string | null>(null);
	let errorMessage = $state<string | null>(null);
	// svelte-ignore state_referenced_locally
	const { form, errors, submitting, enhance } = superForm(data.form, {
		resetForm: false,
		onSubmit(input) {
			handleSubmitLoading(true);
			successMessage = null;
			errorMessage = null;
		},
		async onUpdate(event) {
			if (event.result.type === 'failure') {
				errorMessage = event.result.data.message;
				return;
			}
			if (event.result.type === 'success') {
				handleSubmitLoading(false);
				successMessage = event.result.data.message;
				await invalidateAll();
			}
		},
		onError(event) {
			errorMessage = event.result.error.message;
			handleSubmitLoading(false);
		}
	});
</script>

<MetaTags {...metaTags} />
<AdminSidebarLayout page="robots.txt Setting" user={data.user} setting={data.settings}>
	<AdminSettingLayout fixed title="robots.txt Setting" description="Update robots.txt setting">
		{#snippet children()}
			<form method="POST" class="min-h-0 flex-1 space-y-4 pb-14" use:enhance>
				<Field.Group>
					<Field.Set>
						<Field.Field>
							<Field.Label for="content">robots.txt Content</Field.Label>
							<Textarea
								bind:value={$form.content}
								name="content"
								placeholder="Enter robots.txt content"
								aria-invalid={!!$errors.content}
								autocomplete="on"
								disabled={$submitting}
							/>
							{#if $errors.content}
								<Field.Error>{$errors.content}</Field.Error>
							{/if}
						</Field.Field>
					</Field.Set>
					<Field.Field orientation="horizontal" class="mt-6 justify-end pb-4">
						<Button type="submit" class="w-full" disabled={$submitting}>
							{#if $submitting}
								<Spinner class="mr-2 size-5" />
							{/if}
							{$submitting ? 'Submitting' : 'Update'}
						</Button>
					</Field.Field>
				</Field.Group>
			</form>
		{/snippet}
	</AdminSettingLayout>
	{#if successMessage}
		<AppAlertDialog
			open={true}
			title="Success"
			message={successMessage}
			type="success"
			onclose={() => (successMessage = null)}
		/>
	{/if}
	{#if errorMessage}
		<AppAlertDialog
			open={true}
			title="Error"
			message={errorMessage}
			type="error"
			onclose={() => (errorMessage = null)}
		/>
	{/if}
</AdminSidebarLayout>
