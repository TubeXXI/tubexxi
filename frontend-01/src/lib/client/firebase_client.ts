import { browser } from '$app/environment';
import { initializeApp, type FirebaseApp, getApps, getApp } from 'firebase/app';
import {
	getAuth,
	signInWithPopup,
	signOut as firebaseSignOut,
	onAuthStateChanged,
	GoogleAuthProvider,
	FacebookAuthProvider,
	GithubAuthProvider,
	TwitterAuthProvider,
	type Auth,
	type User,
	type UserCredential,
	signInWithRedirect,
	getRedirectResult,
	createUserWithEmailAndPassword,
	signInWithEmailAndPassword,
	sendEmailVerification,
	sendPasswordResetEmail,
	confirmPasswordReset,
	updatePassword,
	updateProfile,
	EmailAuthProvider,
	reauthenticateWithCredential,
	applyActionCode
} from 'firebase/auth';
import { firebaseClientConfig } from '@/constants/firebase.js';
import * as i18n from '@/paraglide/messages.js';

export class FirebaseClientHelper {
	private static instance: FirebaseClientHelper;
	private app: FirebaseApp | null = null;
	private auth: Auth | null = null;
	private currentUser: User | null = null;

	private constructor() {
		if (browser) {
			this.initializeFirebase();
		}
	}

	static getInstance(): FirebaseClientHelper {
		if (!FirebaseClientHelper.instance) {
			FirebaseClientHelper.instance = new FirebaseClientHelper();
		}
		return FirebaseClientHelper.instance;
	}

	private initializeFirebase(): void {
		try {
			if (getApps().length === 0) {
				this.app = initializeApp(firebaseClientConfig);
			} else {
				this.app = getApp();
			}
			this.auth = getAuth(this.app);

			// Setup auth state listener
			this.setupAuthStateListener();
		} catch (error) {
			console.error('Error initializing Firebase:', error);
			throw error;
		}
	}

	private setupAuthStateListener(): void {
		if (!this.auth) return;

		onAuthStateChanged(this.auth, (user) => {
			this.currentUser = user;
		});
	}

	/**
	 * Get provider berdasarkan nama
	 */
	private getProvider(providerName: SocialProvider) {
		switch (providerName) {
			case 'google':
				const googleProvider = new GoogleAuthProvider();
				googleProvider.addScope('profile');
				googleProvider.addScope('email');
				return googleProvider;
			case 'facebook':
				const facebookProvider = new FacebookAuthProvider();
				facebookProvider.addScope('email');
				facebookProvider.addScope('public_profile');
				return facebookProvider;
			case 'github':
				const githubProvider = new GithubAuthProvider();
				githubProvider.addScope('user:email');
				return githubProvider;
			case 'twitter':
				return new TwitterAuthProvider();
			default:
				throw new Error(i18n.error_unsupported_provider({ provider: providerName }));
		}
	}

	/**
	 * Sign with social provider using popup or redirect flow
	 */
	async signInWithSocial(
		provider: SocialProvider,
		usePopup = true
	): Promise<UserCredential | null> {
		if (!this.auth) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			const authProvider = this.getProvider(provider);

			if (usePopup) {
				const result = await signInWithPopup(this.auth, authProvider);
				return result;
			} else {
				await signInWithRedirect(this.auth, authProvider);
				return null;
			}
		} catch (error: any) {
			console.error(`Error signing in with ${provider}:`, error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Get redirect result (for redirect flow)
	 */
	async getRedirectResult(): Promise<UserCredential | null> {
		if (!this.auth) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			const result = await getRedirectResult(this.auth);
			return result;
		} catch (error: any) {
			console.error('Error getting redirect result:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Get ID Token from current user (with optional force refresh)
	 */
	async getIdToken(forceRefresh = false): Promise<string | null> {
		if (!this.auth?.currentUser) {
			return null;
		}

		try {
			const tokenResult = await this.auth.currentUser.getIdTokenResult();
			const expirationTime = new Date(tokenResult.expirationTime).getTime();
			const currentTime = Date.now();
			const fiveMinutes = 5 * 60 * 1000;

			if (currentTime + fiveMinutes >= expirationTime) {
				return tokenResult.token;
			}

			const shouldRefresh = forceRefresh || expirationTime - currentTime < fiveMinutes;

			const token = await this.auth.currentUser.getIdToken(shouldRefresh);

			return token;
		} catch (error) {
			console.error('Error getting ID token:', error);
			return null;
		}
	}

	/**
	 * Sign out
	 */
	async signOut(): Promise<void> {
		if (!this.auth) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			await firebaseSignOut(this.auth);
			this.currentUser = null;
		} catch (error) {
			console.error('Error signing out:', error);
			throw error;
		}
	}

	/**
	 * Get current user
	 */
	getCurrentUser(): User | null {
		return this.auth?.currentUser ?? null;
	}

	/**
	 * Wait for auth state to be ready
	 */
	async waitForAuthReady(): Promise<User | null> {
		if (!this.auth) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		return new Promise((resolve) => {
			const unsubscribe = onAuthStateChanged(this.auth!, (user) => {
				unsubscribe();
				resolve(user);
			});
		});
	}

	/**
	 * Handle Firebase auth errors
	 */
	private handleAuthError(error: any): Error {
		const errorCode = error.code;
		const errorMessage = error.message;

		const errorMessages: Record<string, string> = {
			'auth/popup-closed-by-user': i18n.error_login_cancelled(),
			'auth/popup-blocked': i18n.error_popup_blocked(),
			'auth/cancelled-popup-request': i18n.error_login_cancelled(),
			'auth/account-exists-with-different-credential':
				i18n.error_account_exists_with_different_credential(),
			'auth/invalid-credential': i18n.error_invalid_credential(),
			'auth/operation-not-allowed': i18n.error_operation_not_allowed(),
			'auth/user-disabled': i18n.error_user_disabled(),
			'auth/user-not-found': i18n.error_user_not_found(),
			'auth/wrong-password': i18n.error_wrong_password(),
			'auth/too-many-requests': i18n.error_too_many_requests(),
			'auth/network-request-failed': i18n.error_network_request_failed(),
			// Email/Password specific errors
			'auth/email-already-in-use': i18n.error_email_already_in_use(),
			'auth/invalid-email': i18n.error_invalid_email(),
			'auth/weak-password': i18n.error_weak_password(),
			'auth/missing-password': i18n.error_missing_password(),
			'auth/requires-recent-login': i18n.error_requires_recent_login(),
			'auth/expired-action-code': i18n.error_expired_action_code(),
			'auth/invalid-action-code': i18n.error_invalid_action_code(),
			'auth/user-token-expired': i18n.error_user_token_expired()
		};

		const message = errorMessages[errorCode] || errorMessage || i18n.error_unknown_login_error();

		return new Error(message);
	}

	/**
	 * Get auth instance (for custom operations)
	 */
	getAuth(): Auth | null {
		return this.auth;
	}

	// ============================================
	// EMAIL/PASSWORD AUTHENTICATION METHODS
	// ============================================

	/**
	 * Register with email and password
	 */
	async registerWithEmail(
		email: string,
		password: string,
		displayName?: string
	): Promise<UserCredential> {
		if (!this.auth) {
			throw new Error(i18n.error_firebase_auth_not_initialized());
		}

		try {
			const userCredential = await createUserWithEmailAndPassword(this.auth, email, password);

			if (displayName && userCredential.user) {
				await updateProfile(userCredential.user, { displayName });
			}

			// await this.sendVerificationEmail();

			return userCredential;
		} catch (error: any) {
			console.error('Error registering with email:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Sign in with email and password
	 */
	async signInWithEmail(email: string, password: string): Promise<UserCredential> {
		if (!this.auth) {
			throw new Error(i18n.error_firebase_auth_not_initialized());
		}

		try {
			const userCredential = await signInWithEmailAndPassword(
				this.auth,
				email,
				password
			);
			return userCredential;
		} catch (error: any) {
			console.error('Error signing in with email:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Send email verification to current user
	 */
	async sendVerificationEmail(): Promise<void> {
		if (!this.auth?.currentUser) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			await sendEmailVerification(this.auth.currentUser);
		} catch (error: any) {
			console.error('Error sending verification email:', error);
			throw this.handleAuthError(error);
		}
	}

	async applyEmailVerificationCode(code: string): Promise<void> {
		if (!this.auth) {
			throw new Error(i18n.error_firebase_auth_not_initialized());
		}
		try {
			await applyActionCode(this.auth, code);
			if (this.auth.currentUser) {
				await this.auth.currentUser.reload();
			}
		} catch (error: any) {
			console.error('Error applying email verification code:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Send password reset email to user with identifier
	 */
	async sendPasswordReset(email: string): Promise<void> {
		if (!this.auth) {
			throw new Error(i18n.error_firebase_auth_not_initialized());
		}

		try {
			await sendPasswordResetEmail(this.auth, email);
		} catch (error: any) {
			console.error('Error sending password reset email:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Confirm password reset with code from email
	 */
	async confirmPasswordReset(code: string, newPassword: string): Promise<void> {
		if (!this.auth) {
			throw new Error(i18n.error_firebase_auth_not_initialized());
		}

		try {
			await confirmPasswordReset(this.auth, code, newPassword);
		} catch (error: any) {
			console.error('Error confirming password reset:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Update password for currently signed-in user
	 */
	async updatePassword(newPassword: string): Promise<void> {
		if (!this.auth?.currentUser) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			await updatePassword(this.auth.currentUser, newPassword);
		} catch (error: any) {
			console.error('Error updating password:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Update user profile
	 */
	async updateUserProfile(profile: { displayName?: string; photoURL?: string }): Promise<void> {
		if (!this.auth?.currentUser) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			await updateProfile(this.auth.currentUser, profile);
		} catch (error: any) {
			console.error('Error updating profile:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Re-authenticate user with password (required for sensitive operations)
	 */
	async reauthenticateWithPassword(password: string): Promise<UserCredential> {
		if (!this.auth?.currentUser?.email) {
			throw new Error(i18n.error_no_user_is_currently_signed_in());
		}

		try {
			const credential = EmailAuthProvider.credential(this.auth.currentUser.email, password);
			const userCredential = await reauthenticateWithCredential(this.auth.currentUser, credential);
			return userCredential;
		} catch (error: any) {
			console.error('Error re-authenticating:', error);
			throw this.handleAuthError(error);
		}
	}

	/**
	 * Check if current user's email is verified
	 */
	isEmailVerified(): boolean {
		return this.auth?.currentUser?.emailVerified ?? false;
	}

	/**
	 * Setup automatic token refresh
	 */
	setupTokenRefresh(): () => void {
		if (!this.auth) return () => { };

		const interval = setInterval(
			async () => {
				try {
					const user = this.auth?.currentUser;
					if (user) {
						await user.getIdToken(true); // Force refresh
						console.log('Token automatically refreshed');
					}
				} catch (error) {
					console.error('Auto token refresh failed:', error);
				}
			},
			55 * 60 * 1000
		); // 55 minutes

		return () => clearInterval(interval);
	}
}

export const firebaseClient = browser ? FirebaseClientHelper.getInstance() : null;
