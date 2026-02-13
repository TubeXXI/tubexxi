import { writable, derived, type Readable } from 'svelte/store';
import { browser } from '$app/environment';
import { firebaseClient } from '@/client/firebase_client';

function createAuthStore() {
	const initialState: AuthState = {
		user: null,
		loading: true,
		error: null
	};

	const { subscribe, set, update } = writable<AuthState>(initialState);

	// Setup auth state listener saat di browser
	if (browser && firebaseClient) {
		const auth = firebaseClient.getAuth();

		if (auth) {
			// Listen to auth state changes
			auth.onAuthStateChanged(
				(user) => {
					update((state) => ({
						...state,
						user: user
							? {
									id: user.uid,
									email: user.email || '',
									full_name: user.displayName || '',
									avatar_url: user.photoURL || '',
									role_id: 'user',
									email_verified_at: user.emailVerified ? user.metadata.creationTime : null,
									is_active: true,
									is_verified: user.emailVerified,
									created_at: user.metadata.creationTime || '',
									updated_at: user.metadata.lastSignInTime || ''
								}
							: null,
						loading: false,
						error: null
					}));
				},
				(error) => {
					update((state) => ({
						...state,
						loading: false,
						error: error.message
					}));
				}
			);
		}
	}

	return {
		subscribe,
		/**
		 * Sign in dengan social provider
		 */
		signInWithSocial: async (provider: SocialProvider, usePopup = true) => {
			if (!firebaseClient) {
				update((state) => ({
					...state,
					error: 'Firebase not initialized'
				}));
				return null;
			}

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				const result = await firebaseClient.signInWithSocial(provider, usePopup);

				if (result) {
					// Get ID token dan simpan ke cookie
					const idToken = await result.user.getIdToken();

					// Send token ke server untuk create session
					await fetch('/api/auth/session', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
						credentials: 'include',
						body: JSON.stringify({ idToken, rememberMe: false })
					});

					update((state) => ({
						...state,
						user: result.user
							? {
									id: result.user.uid,
									email: result.user.email || '',
									full_name: result.user.displayName || '',
									avatar_url: result.user.photoURL || '',
									role_id: 'user',
									email_verified_at: result.user.emailVerified
										? result.user.metadata.creationTime
										: null,
									is_active: true,
									is_verified: result.user.emailVerified,
									created_at: result.user.metadata.creationTime || '',
									updated_at: result.user.metadata.lastSignInTime || ''
								}
							: null,
						loading: false,
						error: null
					}));

					return result;
				}

				update((state) => ({ ...state, loading: false }));
				return null;
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Register dengan email dan password
		 */
		registerWithEmail: async (email: string, password: string, displayName?: string) => {
			if (!firebaseClient) {
				update((state) => ({
					...state,
					error: 'Firebase not initialized'
				}));
				return null;
			}

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				const result = await firebaseClient.registerWithEmail(email, password, displayName);

				// Get ID token dan simpan ke cookie
				// const idToken = await result.user.getIdToken();

				// Send token ke server untuk create session
				// await fetch('/api/auth/session', {
				// 	method: 'POST',
				// 	headers: { 'Content-Type': 'application/json' },
				// 	body: JSON.stringify({ idToken })
				// });

				update((state) => ({
					...state,
					user: result.user
						? {
								id: result.user.uid,
								email: result.user.email || '',
								full_name: result.user.displayName || '',
								avatar_url: result.user.photoURL || '',
								role_id: 'user',
								email_verified_at: result.user.emailVerified
									? result.user.metadata.creationTime
									: null,
								is_active: true,
								is_verified: result.user.emailVerified,
								created_at: result.user.metadata.creationTime || '',
								updated_at: result.user.metadata.lastSignInTime || ''
							}
						: null,
					loading: false,
					error: null
				}));

				return result;
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Login dengan email dan password
		 */
		signInWithEmail: async (email: string, password: string) => {
			if (!firebaseClient) {
				update((state) => ({
					...state,
					error: 'Firebase not initialized'
				}));
				return null;
			}

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				const result = await firebaseClient.signInWithEmail(email, password);

				// Get ID token dan simpan ke cookie
				const idToken = await result.user.getIdToken();

				// Send token ke server untuk create session
				await fetch('/api/auth/session', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
					credentials: 'include',
					body: JSON.stringify({ idToken, rememberMe: false })
				});

				update((state) => ({
					...state,
					user: result.user
						? {
								id: result.user.uid,
								email: result.user.email || '',
								full_name: result.user.displayName || '',
								avatar_url: result.user.photoURL || '',
								role_id: 'user',
								email_verified_at: result.user.emailVerified
									? result.user.metadata.creationTime
									: null,
								is_active: true,
								is_verified: result.user.emailVerified,
								created_at: result.user.metadata.creationTime || '',
								updated_at: result.user.metadata.lastSignInTime || ''
							}
						: null,
					loading: false,
					error: null
				}));

				return result;
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Send verification email
		 */
		sendVerificationEmail: async () => {
			if (!firebaseClient) return;

			try {
				await firebaseClient.sendVerificationEmail();
			} catch (error: any) {
				update((state) => ({
					...state,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Send password reset email
		 */
		sendPasswordReset: async (email: string) => {
			if (!firebaseClient) return;

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				await firebaseClient.sendPasswordReset(email);
				update((state) => ({ ...state, loading: false }));
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Confirm password reset
		 */
		confirmPasswordReset: async (code: string, newPassword: string) => {
			if (!firebaseClient) return;

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				await firebaseClient.confirmPasswordReset(code, newPassword);
				update((state) => ({ ...state, loading: false }));
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Update password
		 */
		updatePassword: async (newPassword: string) => {
			if (!firebaseClient) return;

			try {
				await firebaseClient.updatePassword(newPassword);
			} catch (error: any) {
				update((state) => ({
					...state,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Update profile
		 */
		updateProfile: async (profile: { displayName?: string; photoURL?: string }) => {
			if (!firebaseClient) return;

			try {
				await firebaseClient.updateUserProfile(profile);

				// Refresh user state
				const currentUser = firebaseClient.getCurrentUser();
				update((state) => ({
					...state,
					user: currentUser
						? {
								id: currentUser.uid,
								email: currentUser.email || '',
								full_name: currentUser.displayName || '',
								avatar_url: currentUser.photoURL || '',
								role_id: 'user',
								email_verified_at: currentUser.emailVerified
									? currentUser.metadata.creationTime
									: null,
								is_active: true,
								is_verified: currentUser.emailVerified,
								created_at: currentUser.metadata.creationTime || '',
								updated_at: currentUser.metadata.lastSignInTime || ''
							}
						: null
				}));
			} catch (error: any) {
				update((state) => ({
					...state,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Re-authenticate dengan password
		 */
		reauthenticateWithPassword: async (password: string) => {
			if (!firebaseClient) return null;

			try {
				const result = await firebaseClient.reauthenticateWithPassword(password);
				return result;
			} catch (error: any) {
				update((state) => ({
					...state,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Handle redirect result (untuk redirect flow)
		 */
		handleRedirectResult: async () => {
			if (!firebaseClient) return null;

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				const result = await firebaseClient.getRedirectResult();

				if (result) {
					const idToken = await result.user.getIdToken();

					await fetch('/api/auth/session', {
						method: 'POST',
						headers: { 'Content-Type': 'application/json' },
						body: JSON.stringify({ idToken })
					});

					update((state) => ({
						...state,
						user: result.user
							? {
									id: result.user.uid,
									email: result.user.email || '',
									full_name: result.user.displayName || '',
									avatar_url: result.user.photoURL || '',
									role_id: 'user',
									email_verified_at: result.user.emailVerified
										? result.user.metadata.creationTime
										: null,
									is_active: true,
									is_verified: result.user.emailVerified,
									created_at: result.user.metadata.creationTime || '',
									updated_at: result.user.metadata.lastSignInTime || ''
								}
							: null,
						loading: false,
						error: null
					}));

					return result;
				}

				update((state) => ({ ...state, loading: false }));
				return null;
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				return null;
			}
		},

		/**
		 * Sign out
		 */
		signOut: async () => {
			if (!firebaseClient) return;

			update((state) => ({ ...state, loading: true, error: null }));

			try {
				await firebaseClient.signOut();

				// Delete session di server
				await fetch('/api/auth/session', {
					method: 'DELETE'
				});

				update((state) => ({
					...state,
					user: null,
					loading: false,
					error: null
				}));
			} catch (error: any) {
				update((state) => ({
					...state,
					loading: false,
					error: error.message
				}));
				throw error;
			}
		},

		/**
		 * Refresh ID token
		 */
		refreshToken: async () => {
			if (!firebaseClient) return null;

			const user = firebaseClient.getCurrentUser();
			if (!user) return null;

			try {
				const idToken = await user.getIdToken(true);

				// Update session di server
				await fetch('/api/auth/session', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ idToken })
				});

				return idToken;
			} catch (error) {
				console.error('Error refreshing token:', error);
				return null;
			}
		},

		/**
		 * Clear error
		 */
		clearError: () => {
			update((state) => ({ ...state, error: null }));
		},

		/**
		 * Reset store
		 */
		reset: () => {
			set(initialState);
		}
	};
}

export const authStore = createAuthStore();

// Derived stores untuk kemudahan akses
export const currentUser: Readable<User | null> = derived(
	authStore,
	($authStore) => $authStore.user
);

export const isAuthenticated: Readable<boolean> = derived(
	authStore,
	($authStore) => $authStore.user !== null
);

export const authLoading: Readable<boolean> = derived(
	authStore,
	($authStore) => $authStore.loading
);

export const authError: Readable<string | null> = derived(
	authStore,
	($authStore) => $authStore.error
);

export const isEmailVerified: Readable<boolean> = derived(
	authStore,
	($authStore) => $authStore.user?.email_verified_at !== null
);
