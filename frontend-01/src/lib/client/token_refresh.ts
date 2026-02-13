import { goto } from '$app/navigation';
import { browser } from '$app/environment';
import { firebaseClient } from './firebase_client';
import { localizeHref } from '@/paraglide/runtime';

/**
 * Check token status dan auto refresh if expired
 */
export async function checkAndRefreshToken(): Promise<boolean> {
	try {
		if (!firebaseClient) return false;

		const newToken = await firebaseClient.getIdToken(true);
		if (!newToken) {
			console.error('Failed to get refreshed token');
			return false;
		}

		const updateResponse = await fetch('/api/auth/refresh', {
			method: 'PUT',
			headers: { 'Content-Type': 'application/json', 'X-Platform': 'web' },
			credentials: 'include',
			body: JSON.stringify({ idToken: newToken })
		});
		if (!updateResponse.ok) {
			const errorData = await updateResponse.json();
			console.error('Failed to update token on server:', errorData);

			if (updateResponse.status === 401) {
				await firebaseClient.signOut();
				await goto(localizeHref('/auth/login?session=expired'));
			}
			return false;
		}

		return true;
	} catch (error) {
		console.error('Error checking and refreshing token:', error);
		return false;
	}
}
/**
 * Setup periodic token check (every 30 minutes)
 */
export function setupTokenRefreshInterval(): () => void {
	if (!browser) return () => {};

	const interval = setInterval(
		async () => {
			await checkAndRefreshToken();
		},
		30 * 60 * 1000
	);

	return () => clearInterval(interval);
}
/**
 * Setup visibility change event listener to check token when page becomes visible
 */
export function setupVisibilityRefresh(): () => void {
	if (!browser) {
		return () => {};
	}

	const handleVisibilityChange = async () => {
		if (document.visibilityState === 'visible') {
			// console.log('Page visible, checking token...');
			await checkAndRefreshToken();
		}
	};

	document.addEventListener('visibilitychange', handleVisibilityChange);

	// Return cleanup function
	return () => {
		document.removeEventListener('visibilitychange', handleVisibilityChange);
	};
}
