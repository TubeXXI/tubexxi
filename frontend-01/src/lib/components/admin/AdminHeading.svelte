<script lang="ts">
	interface Props {
		title: string;
		description?: string;
	}
	let { title, description }: Props = $props();
</script>

<div class="mb-8 space-y-0.5">
	<div class="text-lg font-semibold tracking-tight sm:text-xl">
		{title}
	</div>
	{#if description}
		<div class="text-xs text-muted-foreground sm:text-sm">
			{description}
		</div>
	{/if}
</div>
