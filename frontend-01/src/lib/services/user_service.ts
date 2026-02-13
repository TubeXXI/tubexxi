import type { RequestEvent } from '@sveltejs/kit';
import { BaseService } from './base_service';
import type { UpdateProfileSchema, UpdatePasswordSchema } from '@/utils/schema';

export class UserServiceImpl extends BaseService implements UserService {
	constructor(
		protected readonly event: RequestEvent,
		private readonly api: ApiClient
	) {
		super(event);
	}
	async UpdateProfile(data: UpdateProfileSchema): Promise<User | Error> {
		try {
			const response = await this.api.authRequest<User>('PUT', '/user/protected/profile', data);
			if (!response.success) {
				throw new Error(response.message || 'Failed to update profile');
			}
			if (!response.data) {
				throw new Error('Failed to update profile');
			}
			return response.data;
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async CurrentUser(): Promise<User | Error> {
		try {
			const response = await this.api.authRequest<User>('GET', '/user/protected/current');
			if (!response.success) {
				throw new Error(response.message || 'Failed to get current user');
			}
			if (!response.data) {
				throw new Error('Failed to get current user');
			}
			return response.data;
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async UpdateAvatar(file: File): Promise<string | Error> {
		try {
			const formData = new FormData();
			formData.append('avatar', file);

			const response = await this.api.multipartAuthRequest<{ avatar_url?: string }>(
				'POST',
				'/user/protected/avatar',
				formData
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to update avatar');
			}
			if (!response.data) {
				throw new Error('Failed to update avatar');
			}
			return response.data.avatar_url || '';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async UpdatePassword(data: UpdatePasswordSchema): Promise<string | Error> {
		try {
			const response = await this.api.authRequest<{ message?: string }>(
				'PUT',
				'/user/protected/password',
				data
			);
			if (!response.success) {
				throw new Error(response.message || 'Failed to update password');
			}
			return response.message || 'Password updated successfully';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
	async Logout(): Promise<string | Error> {
		try {
			const response = await this.api.authRequest('POST', '/user/protected/logout');
			if (!response.success) {
				throw new Error(response.message || 'Failed to logout');
			}

			this.event.locals.deps.authHelper.clearAuthCookies();
			this.event.locals.user = null;

			return response.message || 'Logout successful';
		} catch (error) {
			return error instanceof Error ? error : new Error('Unknown server error');
		}
	}
}
