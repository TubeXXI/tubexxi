import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';
import type {
	RegisterAppSchema,
	ResetPasswordSchema,
	VerifyEmailSchema,
	ChangePasswordSchema
} from '@/utils/schema';

export class AuthServiceImpl extends BaseService implements AuthService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient) {
		super(event);
	}
	async Login(idToken: string): Promise<FirebaseAuthResponse | Error> {
		try {
			const response = await this.api.publicRequest<FirebaseAuthResponse>('POST', '/auth/login', { idToken });
			if (!response.success) {
				throw new Error(response.error?.message || response.message || 'Failed to login');
			}
			if (!response.data) {
				throw new Error(response.error?.message || response.message || 'Failed to login');
			}
			return response.data;
		} catch (error) {
			return error instanceof Error ? error : new Error('Failed to login');
		}
	}
	async Register(data: RegisterAppSchema): Promise<FirebaseAuthResponse | Error> {
		try {
			const response = await this.api.publicRequest<FirebaseAuthResponse>('POST', '/auth/register', data);
			if (!response.success) {
				throw new Error(response.error?.message || response.message || 'Failed to register');
			}
			if (!response.data) {
				throw new Error(response.error?.message || response.message || 'Failed to register');
			}
			return response.data;
		} catch (error) {
			return error instanceof Error ? error : new Error('Failed to register');
		}
	}
	async ResetPassword(data: ResetPasswordSchema): Promise<string | Error> {
		try {
			const response = await this.api.publicRequest<string>('POST', '/auth/reset-password', data);
			if (!response.success) {
				throw new Error(response.error?.message || response.message || 'Failed to reset password');
			}
			if (!response.data) {
				throw new Error(response.error?.message || response.message || 'Failed to reset password');
			}
			return response.message || 'Password reset successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Failed to reset password');
		}
	}
	async VerifyEmail(data: VerifyEmailSchema): Promise<string | Error> {
		try {
			const response = await this.api.authRequest<string>('POST', '/auth/verify-email', data);
			if (!response.success) {
				throw new Error(response.error?.message || response.message || 'Failed to verify email');
			}
			if (!response.data) {
				throw new Error(response.error?.message || response.message || 'Failed to verify email');
			}
			return response.message || 'Email verified successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Failed to verify email');
		}
	}
	async ChangePassword(data: ChangePasswordSchema): Promise<string | Error> {
		try {
			const response = await this.api.authRequest<string>('POST', '/auth/change-password', data);
			if (!response.success) {
				throw new Error(response.error?.message || response.message || 'Failed to change password');
			}
			if (!response.data) {
				throw new Error(response.error?.message || response.message || 'Failed to change password');
			}
			return response.message || 'Password changed successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Failed to change password');
		}
	}
}
