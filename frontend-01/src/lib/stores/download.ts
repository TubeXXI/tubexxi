import { browser } from '$app/environment';
import { PUBLIC_CENTRIFUGE_URL, PUBLIC_API_URL } from '$env/static/public';
import { writable } from 'svelte/store';
import { Centrifuge, Subscription } from 'centrifuge';

const CENTRIFUGO_URL = PUBLIC_CENTRIFUGE_URL;

export type DownloadFormat = {
	format_id?: string;
	url?: string;
	ext?: string;
	resolution?: string;
	filesize?: number | null;
	height?: number | null;
	vcodec?: string;
	acodec?: string;
	tbr?: number | null;
};

export type DownloadTaskView = {
	id: string;
	status: string;
	progress: number;
	title?: string | null;
	thumbnail_url?: string | null;
	type: string;
	created_at?: string | null;
	file_path?: string | null;
	formats?: DownloadFormat[] | null;
};

export type DownloadState = {
	tasks: Record<string, DownloadTaskView>;
};

export type WebsocketStore = ReturnType<typeof createWebsocketStore>;


export const createWebsocketStore = (userID?: string | null) => {
	let centrifuge: Centrifuge | null = null;
	let subscriptions: Record<string, Subscription> = {};

	const state = writable<DownloadState>({
		tasks: {}
	});

	function unsubscribe(taskId: string) {
		const sub = subscriptions[taskId];
		if (sub) {
			try {
				sub.unsubscribe();
			} catch (_) {
			}
			delete subscriptions[taskId];
		}
	}

	function applyEvent(data: any) {
		if (!data || !data.task_id) return;

		const id = String(data.task_id);
		const status = String(data.status ?? '');
		const progress =
			typeof data.progress === 'number'
				? data.progress
				: status === 'completed'
					? 100
					: status === 'failed'
						? 0
						: 0;

		state.update((current) => {
			const existing = current.tasks[id];
			const payload = data.payload ?? {};

			return {
				tasks: {
					...current.tasks,
					[id]: {
						...existing,
						id,
						status: status || existing?.status || 'processing',
						progress,
						created_at: data.created_at ?? existing?.created_at ?? null,
						file_path: payload.file_path ?? existing?.file_path ?? null,
						formats: payload.formats ?? existing?.formats ?? null
					}
				}
			};
		});

		if (status === 'completed' || status === 'failed') {
			unsubscribe(id);
		}
	}

	async function connect() {
		if (!browser) return;

		if (centrifuge) return;

		// Use environment variable if available (need to be added to .env) or fallback
		// For this environment, user specified ws://infrastructure-centrifugo:8000
		// But browser needs localhost or public URL.
		const url = CENTRIFUGO_URL;

		// Fetch connection token
		let token = '';
		try {
			// Ensure no double slash
			const baseUrl = PUBLIC_API_URL.endsWith('/') ? PUBLIC_API_URL.slice(0, -1) : PUBLIC_API_URL;
			const res = await fetch(`${baseUrl}/web-client/centrifugo/token`);
			if (res.ok) {
				const data = await res.json();
				token = data.data.token;
			} else {
				console.error('Failed to fetch Centrifugo token:', res.statusText);
			}
		} catch (e) {
			console.error('Error fetching Centrifugo token:', e);
		}

		centrifuge = new Centrifuge(url, {
			token: token
		});

		centrifuge.on('connected', (ctx) => {
			// console.log('Centrifugo connected', ctx);
		});

		centrifuge.on('disconnected', (ctx) => {
			// console.log('Centrifugo disconnected', ctx);
		});

		centrifuge.connect();
	}

	function subscribe(taskId: string) {
		if (!browser || !centrifuge) return;
		if (subscriptions[taskId]) return;

		const channel = `download:progress:${taskId}`;
		// console.log(`Subscribing to channel: ${channel}`);

		const sub = centrifuge.newSubscription(channel);

		sub.on('publication', (ctx) => {
			// ctx.data is the payload published from backend
			applyEvent(ctx.data);
		});

		sub.on('subscribing', (ctx) => {
			// console.log(`Subscribing to ${channel}`, ctx);
		});

		sub.on('subscribed', (ctx) => {
			// console.log(`Subscribed to ${channel}`, ctx);
		});

		sub.on('error', (ctx) => {
			console.error(`Subscription error for ${channel}`, ctx);
		});

		sub.subscribe();
		subscriptions[taskId] = sub;
	}

	function upsertTaskFromApi(task: any) {
		if (!task || !task.id) return;

		const id = String(task.id);

		state.update((current) => {
			const existing = current.tasks[id];

			return {
				tasks: {
					...current.tasks,
					[id]: {
						id,
						status: task.status ?? existing?.status ?? 'queued',
						progress: existing?.progress ?? 0,
						title: task.title ?? existing?.title ?? null,
						thumbnail_url: task.thumbnail_url ?? existing?.thumbnail_url ?? null,
						type: task.type ?? existing?.type ?? null,
						created_at: task.created_at ?? existing?.created_at ?? null,
						file_path: task.file_path ?? existing?.file_path ?? null,
						formats: task.formats ?? existing?.formats ?? null
					}
				}
			};
		});

		const status = String(task.status ?? '').toLowerCase();
		if (status === 'pending' || status === 'processing' || status === 'queued') {
			if (!centrifuge) {
				void connect();
			}
			subscribe(id);
		}
	}

	function disconnect() {
		if (centrifuge) {
			centrifuge.disconnect();
			centrifuge = null;
			subscriptions = {};
		}
		state.set({ tasks: {} });
	}

	function removeTask(taskId: string) {
		if (!taskId) return;
		unsubscribe(taskId);
		state.update((current) => {
			const updatedTasks: Record<string, DownloadTaskView> = { ...current.tasks };
			if (updatedTasks[taskId]) {
				delete updatedTasks[taskId];
			}
			return {
				...current,
				tasks: updatedTasks
			};
		});
	}

	function clearAll() {
		Object.keys(subscriptions).forEach((id) => unsubscribe(id));
		state.set({ tasks: {} });
	}

	return {
		state,
		stateSubscribe: state.subscribe,
		connect,
		upsertTaskFromApi,
		disconnect,
		subscribe,
		removeTask,
		clearAll
	};
};
