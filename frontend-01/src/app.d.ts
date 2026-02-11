import type { Dependencies } from '$lib/dependencies';

declare global {
	namespace App {
		interface Error {
			code?: string;
			retryAfter?: number;
			redirect?: string;
		}
		interface Locals {
			safeGetUser: () => Promise<User | null | undefined>;
			safeGetSettings: () => Promise<SettingsValue | null | undefined>;
			user?: User | null;
			settings?: SettingsValue | null;
			deps: Dependencies; // Dependencies are mandatory
			session?: Session;
			lang: string;
		}
		interface PageData {
			settings?: SettingsValue | null;
			user?: User | null;
			lang: string;
			errors?: {
				code: string;
				message: string;
				details?: any;
			};
		}
		interface PageState {
			settings?: SettingsValue | null;
			user?: User | null;
			lang: string;
			errors?: {
				code: string;
				message: string;
				details?: any;
			};
		}
		// interface Platform {}
	}
	interface Session {
		set(key: SessionType, value: string, maxAge?: number): void;
		get(key: SessionType): string | null;
		delete(key: SessionType): void;
	}
	type SessionType = 'access_token' | 'refresh_token' | 'cookie' | 'csrf';
}
declare module '*.md' {
	import type { SvelteComponent } from 'svelte'

	export default class Comp extends SvelteComponent { }

	export const metadata: Record<string, unknown>
}
export { };
