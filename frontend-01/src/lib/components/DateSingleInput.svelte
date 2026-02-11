<script lang="ts">
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import {
		DateFormatter,
		CalendarDate,
		type DateValue,
		getLocalTimeZone,
		toCalendarDateTime,
		toZoned,
		fromDate
	} from '@internationalized/date';
	import { cn } from '$lib/utils.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	let {
		modelValue = $bindable(),
		onchange,
		name,
		disabled,
		placeholder = 'Pick a date'
	}: {
		modelValue?: string | Date | null;
		onchange?: (value: string | null) => void;
		name?: string;
		disabled?: boolean;
		placeholder?: string;
	} = $props();

	let contentRef = $state<HTMLElement | null>(null);
	const df = new DateFormatter('en-US', {
		dateStyle: 'long'
	});
	const timezone = getLocalTimeZone();

	const createCalendarDateFromDate = (date: Date): CalendarDate => {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	};

	const dateValueToDate = (dateValue: DateValue): Date => {
		if (dateValue instanceof CalendarDate) {
			return new Date(Date.UTC(dateValue.year, dateValue.month - 1, dateValue.day));
		}

		try {
			const calendarDateTime = toCalendarDateTime(dateValue);
			const zonedDateTime = toZoned(calendarDateTime, timezone);
			return zonedDateTime.toDate();
		} catch {
			return new Date();
		}
	};

	const stringToCalendarDate = (dateString: string): CalendarDate | null => {
		try {
			const date = new Date(dateString);
			if (isNaN(date.getTime())) return null;
			return createCalendarDateFromDate(date);
		} catch {
			return null;
		}
	};

	const formatForPostgres = (date: Date): string => {
		return date.toISOString();
	};

	let value = $state<DateValue | undefined>();

	$effect(() => {
		if (!modelValue) {
			value = undefined;
		} else {
			let dateValue: CalendarDate | null = null;

			if (typeof modelValue === 'string') {
				dateValue = stringToCalendarDate(modelValue);
			} else if (modelValue instanceof Date) {
				dateValue = createCalendarDateFromDate(modelValue);
			}

			value = dateValue || undefined;
		}
	});

	const handleDateChange = (date: DateValue | undefined) => {
		if (!date) {
			modelValue = null;
			onchange?.(null);
			return;
		}

		try {
			const jsDate = dateValueToDate(date);
			const formattedDate = formatForPostgres(jsDate);

			modelValue = formattedDate;
			onchange?.(formattedDate);
		} catch (error) {
			console.error('Error converting date:', error);
			modelValue = null;
			onchange?.(null);
		}
	};

	const displayText = $derived.by(() => {
		if (!value) return placeholder;

		try {
			const jsDate = dateValueToDate(value);
			return df.format(jsDate);
		} catch {
			return placeholder;
		}
	});

	const hasValue = $derived(!!value);
</script>

{#if name && modelValue}
	<input type="hidden" {name} value={modelValue} />
{/if}
<Popover.Root>
	<Popover.Trigger
		class={cn(
			buttonVariants({
				variant: 'outline'
			}),
			'w-full justify-start text-left font-normal',
			!hasValue && 'text-muted-foreground'
		)}
		{disabled}
	>
		<CalendarIcon class="mr-2 h-4 w-4" />
		{displayText}
	</Popover.Trigger>
	<Popover.Content bind:ref={contentRef} class="w-auto p-0" align="start">
		<Calendar
			type="single"
			bind:value
			weekdayFormat="short"
			numberOfMonths={1}
			{disabled}
			onValueChange={handleDateChange}
		/>
	</Popover.Content>
</Popover.Root>
