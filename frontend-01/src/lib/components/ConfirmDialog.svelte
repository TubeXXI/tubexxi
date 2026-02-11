<script lang="ts" module>
	type ConfirmDialogProps = {
		open: boolean;
		onOpenChange: (open: boolean) => void;
		title: Snippet;
		disabled?: boolean;
		description: Snippet;
		cancelBtnText?: string;
		confirmText?: string;
		destructive?: boolean;
		handleConfirm: () => void;
		isLoading?: boolean;
		class?: ClassValue;
		children?: Snippet<[]>;
	};
</script>

<script lang="ts">
	import type { ClassValue } from 'svelte/elements';
	import type { Snippet } from 'svelte';
	import { cn } from '@/utils';
	import * as AlertDialog from '@/components/ui/alert-dialog';
	import { Button } from '@/components/ui/button';

	let {
		open = $bindable(false),
		onOpenChange,
		title,
		disabled,
		description,
		cancelBtnText,
		confirmText,
		destructive,
		handleConfirm,
		isLoading,
		class: className,
		children
	}: ConfirmDialogProps = $props();
</script>

<AlertDialog.Root bind:open {onOpenChange}>
	<AlertDialog.Content class={cn('p-6', className)}>
		<AlertDialog.Header class="text-start">
			<AlertDialog.Title>
				{#snippet children()}
					{@render title?.()}
				{/snippet}
			</AlertDialog.Title>
			<AlertDialog.Description>
				{#snippet children()}
					{@render description?.()}
				{/snippet}
			</AlertDialog.Description>
		</AlertDialog.Header>
		{@render children?.()}
		<AlertDialog.Footer>
			<AlertDialog.Cancel disabled={isLoading}>
				{cancelBtnText || 'Cancel'}
			</AlertDialog.Cancel>
			<Button
				variant={destructive ? 'destructive' : 'default'}
				disabled={disabled || isLoading}
				onclick={handleConfirm}
			>
				{confirmText || 'Confirm'}
			</Button>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
