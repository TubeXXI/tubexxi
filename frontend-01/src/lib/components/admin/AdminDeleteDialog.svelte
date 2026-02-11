<script lang="ts" generics="TData">
	import { invalidateAll } from '$app/navigation';
	import { Input } from '@/components/ui/input';
	import * as Field from '@/components/ui/field';
	import { ConfirmDialog } from '@/components';
	import { AlertTriangle } from '@lucide/svelte';
	import * as Alert from '@/components/ui/alert';

	type DeleteDialogProps = {
		open: boolean;
		onOpenChange: (open: boolean) => void;
		data: TData[];
		onconfirm?: (ids: string[] | number[]) => Promise<void>;
		heading?: string;
		disabled?: boolean;
		isLoading?: boolean;
	};

	let {
		open = $bindable(false),
		onOpenChange,
		data,
		onconfirm,
		disabled,
		isLoading,
		heading
	}: DeleteDialogProps = $props();

	let inputValue = $state<string | undefined>(undefined);
	let errorMessage = $state<string | undefined>(undefined);

	async function handleDelete() {
		if (!inputValue) {
			errorMessage = 'Confirm by typing the word "delete"';
			return;
		}
		if (onconfirm) {
			await onconfirm(data.map((item) => (item as any).id));
			inputValue = undefined;
			errorMessage = undefined;
			onOpenChange(false);
			await invalidateAll();
		}
	}
</script>

<ConfirmDialog
	{open}
	{onOpenChange}
	handleConfirm={handleDelete}
	{disabled}
	{isLoading}
	cancelBtnText="Cancel"
	confirmText="Delete"
>
	{#snippet title()}
		<span class="text-destructive">
			<AlertTriangle class="me-1 inline-block stroke-destructive" size={18} />
			{' '}
			{`Delete ${data.length} item${data.length > 1 ? 's' : ''}`}
		</span>
	{/snippet}
	{#snippet description()}
		<div class="space-y-4">
			<Field.Group>
				<Field.Set>
					<Field.Legend>{heading}</Field.Legend>
					<Field.Group>
						<Field.Field>
							<Field.Label for="confirm-delete">Confirm by typing the word "delete"</Field.Label>
							<Input
								bind:value={inputValue}
								type="text"
								name="confirm-delete"
								placeholder="Confirm by typing the word 'delete'"
								aria-invalid={!!errorMessage}
								disabled={disabled || isLoading}
							/>
							{#if errorMessage}
								<Field.Error>{errorMessage}</Field.Error>
							{/if}
						</Field.Field>
						<Field.Field>
							<Alert.Root variant="destructive">
								<Alert.Title>Warning</Alert.Title>
								<Alert.Description>Are you sure you want to delete these items?</Alert.Description>
							</Alert.Root>
						</Field.Field>
						<Field.Field orientation="horizontal"></Field.Field>
					</Field.Group>
				</Field.Set>
			</Field.Group>
		</div>
	{/snippet}
</ConfirmDialog>
