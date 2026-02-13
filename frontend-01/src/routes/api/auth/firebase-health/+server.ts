import { json, type RequestHandler } from '@sveltejs/kit';
import { firebaseClientConfig } from '$lib/constants/firebase';

export const GET: RequestHandler = async () => {
	try {
		const res = await fetch(
			`https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${firebaseClientConfig.apiKey}`,
			{
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({})
			}
		);

		return json({ ok: true, status: res.status });
	} catch (error) {
		return json(
			{
				ok: false,
				error: error instanceof Error ? error.message : 'Unknown error'
			},
			{ status: 500 }
		);
	}
};
