import { redirect, type RequestEvent } from '@sveltejs/kit';

export class BaseHelper {
	protected event: RequestEvent;

	constructor(event: RequestEvent) {
		this.event = event;
	}

	protected static redirectToLogin() {
		throw redirect(302, '/auth/sign-in');
	}
}
