import { redirect, type RequestEvent } from '@sveltejs/kit';
import { BaseHelper } from './base_helper';

export class AuthHelper extends BaseHelper {
	constructor(event: RequestEvent) {
		super(event);
	}
	clearAuthCookies() {
		if (!this.event.locals.session) return;

		this.event.locals.session?.delete('access_token');
		this.event.locals.session?.delete('refresh_token');
		this.event.locals.session?.delete('csrf');
		this.event.locals.session?.delete('cookie');
		this.event.request.headers?.delete('Authorization');
		this.event.request.headers?.delete('X-CSRF-Token');
		this.event.request.headers?.delete('Cookie');
	}

	handleUnauthorized() {
		this.clearAuthCookies();
		throw redirect(302, '/login');
	}

	getSessionByTokenType(tokenType: SessionType) {
		return this.event.locals.session?.get(tokenType);
	}

	parsePostgresTimestamp(timestamp: string | Date) {
		return typeof timestamp === 'string' ? new Date(timestamp) : timestamp;
	}

	getIpAddress() {
		return (
			this.event.request.headers.get('x-forwarded-for') ||
			this.event.request.headers.get('remote-address') ||
			'unknown'
		);
	}

	getUserAgent() {
		return this.event.request.headers.get('user-agent') || 'unknown';
	}
}
