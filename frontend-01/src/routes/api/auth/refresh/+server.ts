import { json, type RequestHandler } from '@sveltejs/kit';
import { ApiClientHandler } from '$lib/helpers/api_helpers';

type RefreshRequestBody = {
	idToken?: string;
	rememberMe?: boolean;
};

type RefreshResponseData = {
	user: User | null;
};

export const PUT: RequestHandler = async (event) => {
	const body = (await event.request.json().catch(() => ({}))) as RefreshRequestBody;
	const idToken = typeof body.idToken === 'string' ? body.idToken.trim() : '';
	const rememberMe = !!body.rememberMe;

	if (!idToken) {
		return json(
			{
				status: 400,
				success: false,
				message: 'Missing idToken'
			} satisfies ApiResponse<RefreshResponseData>,
			{ status: 400 }
		);
	}

	const api = new ApiClientHandler(event);
	const currentUser = await api.authRequest<User>('GET', '/user/protected/current', undefined, {
		Authorization: `Bearer ${idToken}`
	});

	if (!currentUser.success) {
		event.locals.session?.delete('access_token');
		return json(
			{
				status: 401,
				success: false,
				message: currentUser.message || 'Unauthorized',
				error: currentUser.error
			} satisfies ApiResponse<RefreshResponseData>,
			{ status: 401 }
		);
	}

	const maxAge = rememberMe ? 60 * 60 * 24 * 7 : 60 * 60 * 24;
	event.locals.session?.set('access_token', idToken, maxAge);

	return json(
		{
			status: 200,
			success: true,
			message: 'Session refreshed',
			data: { user: currentUser.data ?? null }
		} satisfies ApiResponse<RefreshResponseData>,
		{ status: 200 }
	);
};

