<script lang="ts">
	import { modal } from '@/stores';

	/**
	 * Advanced Example 1: Nested Modals
	 *
	 * This example demonstrates how to open a second modal on top of the first one.
	 * Click the "Open Second Modal" button in the first modal to see it in action.
	 */
	function openNestedModal() {
		modal.open({
			title: 'First Modal',
			content:
				'This is the first modal. Click the button below to open a second modal on top of this one.',
			size: 'md',
			footer: () => `
				<div class="flex justify-between">
					<button
						onclick="modal.close()"
						class="px-4 py-2 bg-neutral-200 dark:bg-neutral-700 text-neutral-700 dark:text-neutral-300 rounded-lg hover:bg-neutral-300 dark:hover:bg-neutral-600 transition-colors"
					>
						Close
					</button>
					<button
						onclick="openSecondModal()"
						class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Open Second Modal
					</button>
				</div>
			`
		});
	}

	/**
	 * Make function available globally for onclick
	 */
	(window as any).openSecondModal = () => {
		modal.open({
			title: 'Second Modal',
			description: 'This modal is opened on top of the first one',
			content:
				'You can have multiple modals open at the same time. Close this to go back to the first modal.',
			size: 'sm',
			footer: () => `
				<div class="flex justify-end">
					<button
						onclick="modal.close()"
						class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
					>
						Close This Modal
					</button>
				</div>
			`
		});
	};

	/**
	 * Advanced Example 2: Form with Validation
	 *
	 * This example demonstrates a registration form with client-side validation.
	 * Click the "Submit" button to see the validation in action.
	 */
	function openFormModal() {
		modal.open({
			title: 'Registration Form',
			description: 'Please fill in all required fields',
			content: `
				<form id="regForm" class="space-y-4">
					<div>
						<label class="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">
							Full Name *
						</label>
						<input
							type="text"
							name="fullname"
							required
							class="w-full px-3 py-2 border border-neutral-300 dark:border-neutral-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-neutral-700 dark:text-white"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">
							Email *
						</label>
						<input
							type="email"
							name="email"
							required
							class="w-full px-3 py-2 border border-neutral-300 dark:border-neutral-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-neutral-700 dark:text-white"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">
							Phone Number
						</label>
						<input
							type="tel"
							name="phone"
							class="w-full px-3 py-2 border border-neutral-300 dark:border-neutral-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-neutral-700 dark:text-white"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">
							Message
						</label>
						<textarea
							name="message"
							rows="3"
							class="w-full px-3 py-2 border border-neutral-300 dark:border-neutral-600 rounded-lg focus:ring-2 focus:ring-blue-500 dark:bg-neutral-700 dark:text-white resize-none"
						></textarea>
					</div>
				</form>
			`,
			size: 'lg',
			footer: () => `
				<div class="flex justify-end gap-2">
					<button
						type="button"
						onclick="modal.close()"
						class="px-4 py-2 bg-neutral-200 dark:bg-neutral-700 text-neutral-700 dark:text-neutral-300 rounded-lg hover:bg-neutral-300 dark:hover:bg-neutral-600 transition-colors"
					>
						Cancel
					</button>
					<button
						type="button"
						onclick="submitForm()"
						class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
					>
						Submit
					</button>
				</div>
			`
		});
	}

	(window as any).submitForm = () => {
		const form = document.getElementById('regForm') as HTMLFormElement;
		if (form.checkValidity()) {
			const formData = new FormData(form);
			console.log('Form submitted:', Object.fromEntries(formData));
			modal.close();
			modal.alert('Success!', 'Your registration has been submitted successfully.');
		} else {
			form.reportValidity();
		}
	};

	/**
	 * Advanced Example 3: Image Gallery Modal
	 *
	 * This example demonstrates an image gallery modal with thumbnail previews.
	 * Click on any thumbnail to open a larger preview in a new modal.
	 */
	function openGalleryModal() {
		const images = [
			'https://picsum.photos/800/600?random=1',
			'https://picsum.photos/800/600?random=2',
			'https://picsum.photos/800/600?random=3',
			'https://picsum.photos/800/600?random=4'
		];

		modal.open({
			title: 'Image Gallery',
			content: `
				<div class="grid grid-cols-2 gap-4">
					${images
						.map(
							(img, i) => `
						<div class="relative group cursor-pointer" onclick="openImagePreview('${img}')">
							<img
								src="${img}"
								alt="Gallery image ${i + 1}"
								class="w-full h-48 object-cover rounded-lg group-hover:opacity-75 transition-opacity"
							/>
							<div class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
								<svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
								</svg>
							</div>
						</div>
					`
						)
						.join('')}
				</div>
			`,
			size: 'xl'
		});
	}

	(window as any).openImagePreview = (imageUrl: string) => {
		modal.open({
			content: `
				<div class="flex items-center justify-center">
					<img src="${imageUrl}" alt="Preview" class="max-w-full max-h-[70vh] rounded-lg" />
				</div>
			`,
			size: 'full',
			transparent: true,
			closeButtonPosition: 'out-top-right'
		});
	};

	/**
	 * Advanced Example 4: Loading Modal
	 *
	 * This example demonstrates a loading modal with a spinning animation.
	 * Click the "Show Loading Modal" button to see it in action.
	 */
	function openLoadingModal() {
		const id = modal.open({
			content: `
				<div class="flex flex-col items-center justify-center py-8">
					<div class="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600"></div>
					<p class="mt-4 text-neutral-600 dark:text-neutral-400">Loading, please wait...</p>
				</div>
			`,
			showCloseButton: false,
			clickOutside: false,
			animation: true
		});

		// Simulate loading
		setTimeout(() => {
			modal.close(id);
			modal.alert('Complete!', 'Loading finished successfully.');
		}, 3000);
	}

	/**
	 * Advanced Example 5: Confirmation with Input
	 *
	 * This example demonstrates a confirmation modal with user input validation.
	 * Click the "Delete Forever" button to see it in action.
	 */
	function openDeleteConfirmation() {
		modal.open({
			title: 'Confirm Deletion',
			description: 'This action cannot be undone',
			content: `
				<div class="space-y-4">
					<p class="text-neutral-700 dark:text-neutral-300">
						You are about to delete <strong>important-file.pdf</strong>.
						This action is permanent and cannot be undone.
					</p>
					<div>
						<label class="block text-sm font-medium text-neutral-700 dark:text-neutral-300 mb-1">
							Type "DELETE" to confirm:
						</label>
						<input
							type="text"
							id="deleteConfirm"
							class="w-full px-3 py-2 border border-neutral-300 dark:border-neutral-600 rounded-lg focus:ring-2 focus:ring-red-500 dark:bg-neutral-700 dark:text-white"
							placeholder="DELETE"
						/>
					</div>
				</div>
			`,
			footer: () => `
				<div class="flex justify-end gap-2">
					<button
						onclick="modal.close()"
						class="px-4 py-2 bg-neutral-200 dark:bg-neutral-700 text-neutral-700 dark:text-neutral-300 rounded-lg hover:bg-neutral-300 dark:hover:bg-neutral-600 transition-colors"
					>
						Cancel
					</button>
					<button
						onclick="confirmDelete()"
						class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
					>
						Delete Forever
					</button>
				</div>
			`,
			size: 'md'
		});
	}

	(window as any).confirmDelete = () => {
		const input = document.getElementById('deleteConfirm') as HTMLInputElement;
		if (input && input.value === 'DELETE') {
			modal.close();
			modal.alert('Deleted', 'The file has been permanently deleted.');
		} else {
			alert('Please type DELETE to confirm');
		}
	};

	/**
	 * Advanced Example 6: Terms and Conditions
	 *
	 * This example demonstrates a modal with terms and conditions.
	 * Click the "Accept" button to see it in action.
	 */
	function openTermsModal() {
		modal.open({
			title: 'Terms and Conditions',
			description: 'Please read carefully before accepting',
			content: `
				<div class="prose dark:prose-invert max-h-96 overflow-y-auto space-y-4 text-sm">
					<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.</p>
					<h3 class="text-lg font-semibold">1. User Agreement</h3>
					<p>Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
					<h3 class="text-lg font-semibold">2. Privacy Policy</h3>
					<p>Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.</p>
					<h3 class="text-lg font-semibold">3. Data Collection</h3>
					<p>Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</p>
					<h3 class="text-lg font-semibold">4. Limitation of Liability</h3>
					<p>Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium.</p>
					<div class="mt-6 flex items-start gap-2">
						<input type="checkbox" id="termsAccept" class="mt-1" />
						<label for="termsAccept" class="text-sm">
							I have read and agree to the Terms and Conditions
						</label>
					</div>
				</div>
			`,
			footer: () => `
				<div class="flex justify-between items-center">
					<a href="#" class="text-sm text-blue-600 hover:underline">Download PDF</a>
					<div class="flex gap-2">
						<button
							onclick="modal.close()"
							class="px-4 py-2 bg-neutral-200 dark:bg-neutral-700 text-neutral-700 dark:text-neutral-300 rounded-lg hover:bg-neutral-300 dark:hover:bg-neutral-600 transition-colors"
						>
							Decline
						</button>
						<button
							onclick="acceptTerms()"
							class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
						>
							Accept
						</button>
					</div>
				</div>
			`,
			size: 'xl',
			clickOutside: false
		});
	}

	(window as any).acceptTerms = () => {
		const checkbox = document.getElementById('termsAccept') as HTMLInputElement;
		if (checkbox && checkbox.checked) {
			modal.close();
			modal.alert('Thank You!', 'You have accepted the terms and conditions.');
		} else {
			alert('Please check the agreement checkbox');
		}
	};
</script>

<div
	class="min-h-screen bg-linear-to-br from-purple-50 to-blue-50 px-4 py-12 dark:from-neutral-900 dark:to-neutral-800"
>
	<div class="mx-auto max-w-6xl">
		<h1 class="mb-2 text-4xl font-bold text-neutral-900 dark:text-white">
			Advanced Modal Examples
		</h1>
		<p class="mb-8 text-neutral-600 dark:text-neutral-400">
			Complex use cases dan real-world scenarios
		</p>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			<!-- Nested Modals -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900"
					>
						<svg
							class="h-6 w-6 text-blue-600 dark:text-blue-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Nested Modals</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Open multiple modals on top of each other
				</p>
				<button
					onclick={openNestedModal}
					class="w-full rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
				>
					Try Nested Modals
				</button>
			</div>

			<!-- Form Modal -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-green-100 dark:bg-green-900"
					>
						<svg
							class="h-6 w-6 text-green-600 dark:text-green-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Form Modal</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Registration form dengan validasi
				</p>
				<button
					onclick={openFormModal}
					class="w-full rounded-lg bg-green-600 px-4 py-2 text-white transition-colors hover:bg-green-700"
				>
					Open Form
				</button>
			</div>

			<!-- Gallery Modal -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-purple-100 dark:bg-purple-900"
					>
						<svg
							class="h-6 w-6 text-purple-600 dark:text-purple-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Image Gallery</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Image gallery dengan preview
				</p>
				<button
					onclick={openGalleryModal}
					class="w-full rounded-lg bg-purple-600 px-4 py-2 text-white transition-colors hover:bg-purple-700"
				>
					View Gallery
				</button>
			</div>

			<!-- Loading Modal -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-yellow-100 dark:bg-yellow-900"
					>
						<svg
							class="h-6 w-6 animate-spin text-yellow-600 dark:text-yellow-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Loading Modal</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Loading state untuk async operations
				</p>
				<button
					onclick={openLoadingModal}
					class="w-full rounded-lg bg-yellow-600 px-4 py-2 text-white transition-colors hover:bg-yellow-700"
				>
					Show Loading
				</button>
			</div>

			<!-- Delete Confirmation -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-red-100 dark:bg-red-900"
					>
						<svg
							class="h-6 w-6 text-red-600 dark:text-red-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Delete Confirm</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Konfirmasi dengan input validasi
				</p>
				<button
					onclick={openDeleteConfirmation}
					class="w-full rounded-lg bg-red-600 px-4 py-2 text-white transition-colors hover:bg-red-700"
				>
					Delete File
				</button>
			</div>

			<!-- Terms Modal -->
			<div
				class="rounded-xl bg-white p-6 shadow-lg transition-shadow hover:shadow-xl dark:bg-neutral-800"
			>
				<div class="mb-4 flex items-center gap-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-100 dark:bg-indigo-900"
					>
						<svg
							class="h-6 w-6 text-indigo-600 dark:text-indigo-400"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
							/>
						</svg>
					</div>
					<h3 class="text-lg font-semibold text-neutral-900 dark:text-white">Terms & Conditions</h3>
				</div>
				<p class="mb-4 text-sm text-neutral-600 dark:text-neutral-400">
					Long content dengan scrollable area
				</p>
				<button
					onclick={openTermsModal}
					class="w-full rounded-lg bg-indigo-600 px-4 py-2 text-white transition-colors hover:bg-indigo-700"
				>
					Read Terms
				</button>
			</div>
		</div>
	</div>
</div>
