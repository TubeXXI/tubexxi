import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';
import { derived, writable } from 'svelte/store';


interface WsConfig {
	url?: string;
	userId?: string;
	username?: string;
	avatar_url?: string;
	reconnect?: boolean;
	reconnectAttempts?: number;
	reconnectDelay?: number;
	onMessage?: (message: DownloadState) => void;
}

interface WSStoreState {
	connected: boolean;
	messages: DownloadState[];
	identified: boolean;
	error: string | null;
	reconnecting: boolean;
}

export type DownloadTaskView = {
	id: string;
	status: string;
	progress: number;
	title?: string | null;
	thumbnail_url?: string | null;
	created_at?: string | null;
	file_path?: string | null;
};

export type DownloadState = {
	tasks: Record<string, DownloadTaskView>;
};

export function createWebSocketStore(config: WsConfig = {}) { }
