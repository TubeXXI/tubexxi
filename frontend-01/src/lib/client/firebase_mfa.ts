import { browser } from '$app/environment';
import {
	multiFactor,
	PhoneAuthProvider,
	PhoneMultiFactorGenerator,
	RecaptchaVerifier,
	type MultiFactorResolver,
	type MultiFactorInfo,
	type UserCredential
} from 'firebase/auth';
import { firebaseClient } from './firebase_client';
import * as i18n from '@/paraglide/messages.js';

export class FirebaseMFAHelper {
	private recaptchaVerifier: RecaptchaVerifier | null = null;
	private verificationId: string | null = null;

	/**
	 * Initialize reCAPTCHA verifier
	 */
	initRecaptcha(containerId: string = 'recaptcha-container'): RecaptchaVerifier | null {
		if (!browser || !firebaseClient) return null;

		const auth = firebaseClient.getAuth();
		if (!auth) return null;

		try {
			this.recaptchaVerifier = new RecaptchaVerifier(auth, containerId, {
				size: 'invisible',
				callback: () => {
					console.log('✅ reCAPTCHA solved');
				},
				'expired-callback': () => {
					console.log('⚠️ reCAPTCHA expired');
				}
			});

			return this.recaptchaVerifier;
		} catch (error) {
			console.error('❌ Error initializing reCAPTCHA:', error);
			return null;
		}
	}

	/**
	 * Enroll phone number for MFA
	 * Call this from settings page when user wants to enable 2FA
	 */
	async enrollPhoneNumber(phoneNumber: string, displayName?: string): Promise<boolean> {
		if (!browser || !firebaseClient) {
			throw new Error('Client-side only');
		}

		try {
			const auth = firebaseClient.getAuth();
			if (!auth) return false;

			const user = firebaseClient?.getCurrentUser();

			if (!user) {
				throw new Error(i18n.unauthorized_error_message());
			}

			// Initialize reCAPTCHA
			if (!this.recaptchaVerifier) {
				this.recaptchaVerifier = this.initRecaptcha();
			}

			if (!this.recaptchaVerifier) {
				throw new Error(i18n.error_failed_to_initialize_recaptcha());
			}

			// Get multi-factor session
			const multiFactorSession = await multiFactor(user).getSession();

			// Send verification code
			const phoneAuthProvider = new PhoneAuthProvider(auth);
			this.verificationId = await phoneAuthProvider.verifyPhoneNumber(
				{
					phoneNumber,
					session: multiFactorSession
				},
				this.recaptchaVerifier
			);

			console.log('✅ Verification code sent to:', phoneNumber);
			return true;
		} catch (error: any) {
			console.error('❌ Error enrolling phone:', error);
			throw this.handleMFAError(error);
		}
	}

	/**
	 * Complete phone enrollment with verification code
	 */
	async completeEnrollment(verificationCode: string, displayName?: string): Promise<boolean> {
		if (!browser || !firebaseClient) {
			throw new Error('Client-side only');
		}

		try {
			const auth = firebaseClient.getAuth();
			if (!auth) return false;

			const user = firebaseClient?.getCurrentUser();

			if (!user) {
				throw new Error(i18n.unauthorized_error_message());
			}

			if (!this.verificationId) {
				throw new Error(i18n.error_no_verification_id());
			}

			// Create credential
			const credential = PhoneAuthProvider.credential(this.verificationId, verificationCode);

			// Create multi-factor assertion
			const multiFactorAssertion = PhoneMultiFactorGenerator.assertion(credential);

			// Enroll
			await multiFactor(user).enroll(multiFactorAssertion, displayName || 'Phone Number');

			console.log('✅ Phone number enrolled successfully');

			// Clear verification ID
			this.verificationId = null;

			return true;
		} catch (error: any) {
			console.error('❌ Error completing enrollment:', error);
			throw this.handleMFAError(error);
		}
	}

	/**
	 * Get enrolled factors for current user
	 */
	getEnrolledFactors(): MultiFactorInfo[] {
		if (!browser || !firebaseClient) return [];

		const auth = firebaseClient.getAuth();
		if (!auth) return [];

		const user = firebaseClient?.getCurrentUser();

		if (!user) return [];

		return multiFactor(user).enrolledFactors;
	}

	/**
	 * Check if user has MFA enrolled
	 */
	isMFAEnrolled(): boolean {
		const factors = this.getEnrolledFactors();
		return factors.length > 0;
	}

	/**
	 * Unenroll a factor
	 */
	async unenrollFactor(factorUid: string): Promise<boolean> {
		if (!browser || !firebaseClient) {
			throw new Error('Client-side only');
		}

		try {
			const auth = firebaseClient.getAuth();
			if (!auth) return false;

			const user = firebaseClient?.getCurrentUser();

			if (!user) {
				throw new Error(i18n.unauthorized_error_message());
			}

			// Find the factor
			const factor = multiFactor(user).enrolledFactors.find((f) => f.uid === factorUid);

			if (!factor) {
				throw new Error(i18n.error_factor_not_found());
			}

			// Unenroll
			await multiFactor(user).unenroll(factor);

			console.log('✅ Factor unenrolled successfully');
			return true;
		} catch (error: any) {
			console.error('❌ Error unenrolling factor:', error);
			throw this.handleMFAError(error);
		}
	}

	/**
	 * Send verification code during sign-in
	 * Call this when you get auth/multi-factor-auth-required error
	 */
	async sendSignInVerificationCode(resolver: MultiFactorResolver): Promise<string> {
		if (!browser || !firebaseClient) {
			throw new Error('Client-side only');
		}

		try {
			const auth = firebaseClient.getAuth();
			if (!auth) {
				throw new Error(i18n.error_firebase_auth_not_initialized());
			}

			// Get the first hint (phone number)
			const phoneInfoOptions = {
				multiFactorHint: resolver.hints[0],
				session: resolver.session
			};

			// Initialize reCAPTCHA if needed
			if (!this.recaptchaVerifier) {
				this.recaptchaVerifier = this.initRecaptcha();
			}

			if (!this.recaptchaVerifier) {
				throw new Error(i18n.error_failed_to_initialize_recaptcha());
			}

			// Send verification code
			const phoneAuthProvider = new PhoneAuthProvider(auth);
			this.verificationId = await phoneAuthProvider.verifyPhoneNumber(
				phoneInfoOptions,
				this.recaptchaVerifier
			);

			console.log('✅ Sign-in verification code sent');
			return this.verificationId;
		} catch (error: any) {
			console.error('❌ Error sending verification code:', error);
			throw this.handleMFAError(error);
		}
	}

	/**
	 * Complete sign-in with verification code
	 */
	async completeSignIn(
		resolver: MultiFactorResolver,
		verificationCode: string
	): Promise<UserCredential> {
		if (!browser || !firebaseClient) {
			throw new Error('Client-side only');
		}

		try {
			if (!this.verificationId) {
				throw new Error(i18n.error_no_verification_id());
			}

			// Create credential
			const credential = PhoneAuthProvider.credential(this.verificationId, verificationCode);

			// Create multi-factor assertion
			const multiFactorAssertion = PhoneMultiFactorGenerator.assertion(credential);

			// Resolve sign-in
			const userCredential = await resolver.resolveSignIn(multiFactorAssertion);

			console.log('✅ Sign-in completed with MFA');

			// Clear verification ID
			this.verificationId = null;

			return userCredential;
		} catch (error: any) {
			console.error('❌ Error completing sign-in:', error);
			throw this.handleMFAError(error);
		}
	}

	/**
	 * Handle MFA errors with user-friendly messages
	 */
	private handleMFAError(error: any): Error {
		const errorCode = error.code;
		const errorMessages: Record<string, string> = {
			'auth/invalid-verification-code': i18n.error_invalid_verification_code(),
			'auth/code-expired': i18n.error_code_expired(),
			'auth/invalid-phone-number': i18n.error_invalid_phone_number(),
			'auth/missing-phone-number': i18n.error_missing_phone_number(),
			'auth/quota-exceeded': i18n.error_quota_exceeded(),
			'auth/too-many-requests': i18n.error_too_many_requests(),
			'auth/second-factor-already-in-use': i18n.error_second_factor_already_in_use(),
			'auth/maximum-second-factor-count-exceeded':
				i18n.error_maximum_second_factor_count_exceeded(),
			'auth/unsupported-first-factor': i18n.error_unsupported_first_factor(),
			'auth/unverified-email': i18n.error_unverified_email(),
			'auth/multi-factor-auth-required': i18n.error_multi_factor_auth_required(),
			'auth/requires-recent-login': i18n.error_requires_recent_login()
		};

		const message = errorMessages[errorCode] || error.message || i18n.error_enrollment_failed();

		// Preserve the original code for consumers to branch on
		const err = new Error(message) as Error & { code?: string };
		if (errorCode) {
			err.code = errorCode;
		}
		return err;
	}

	/**
	 * Clear reCAPTCHA
	 */
	clearRecaptcha(): void {
		if (this.recaptchaVerifier) {
			try {
				this.recaptchaVerifier.clear();
				this.recaptchaVerifier = null;
			} catch (error) {
				console.error('Error clearing reCAPTCHA:', error);
			}
		}
	}

	/**
	 * Save MFA resolver to sessionStorage
	 */
	saveResolver(resolver: MultiFactorResolver): void {
		if (!browser) return;

		try {
			const resolverData = {
				hints: resolver.hints.map((hint) => ({
					factorId: hint.factorId,
					uid: hint.uid,
					displayName: hint.displayName,
					enrollmentTime: hint.enrollmentTime
				}))
				// Note: session cannot be serialized, will need to get from error again
			};

			sessionStorage.setItem('mfa_resolver', JSON.stringify(resolverData));
			sessionStorage.setItem('mfa_required', 'true');
		} catch (error) {
			console.error('Error saving resolver:', error);
		}
	}

	/**
	 * Get saved resolver data
	 */
	getSavedResolverData(): any | null {
		if (!browser) return null;

		try {
			const data = sessionStorage.getItem('mfa_resolver');
			return data ? JSON.parse(data) : null;
		} catch (error) {
			console.error('Error getting resolver:', error);
			return null;
		}
	}

	/**
	 * Clear saved resolver
	 */
	clearSavedResolver(): void {
		if (!browser) return;

		sessionStorage.removeItem('mfa_resolver');
		sessionStorage.removeItem('mfa_required');
	}

	/**
	 * Check if MFA is required (from sessionStorage)
	 */
	isMFARequired(): boolean {
		if (!browser) return false;
		return sessionStorage.getItem('mfa_required') === 'true';
	}
}

export const mfaHelper = browser ? new FirebaseMFAHelper() : null;
