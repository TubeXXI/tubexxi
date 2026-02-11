import { writable } from 'svelte/store';

export const userStore = writable<User | null>(null);
export const settingStore = writable<SettingsValue | null>(null);
export const langStore = writable<string | null>(null);
