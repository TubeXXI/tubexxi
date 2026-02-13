import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';
import { derived, writable } from 'svelte/store';

interface WsConfig {
	url?: string;
	userId?: string;
	email?: string;
	avatar_url?: string;
	reconnect?: boolean;
	reconnectAttempts?: number;
	reconnectDelay?: number;
	onMessage?: (message: any) => void;
}

interface WSStoreState {
	connected: boolean;
	messages: any[];
	identified: boolean;
	error: string | null;
	reconnecting: boolean;
}

export function createWebSocketStore(config: WsConfig = {}) {}
