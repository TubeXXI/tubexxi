import { json, type RequestHandler } from '@sveltejs/kit';
import { ApiClientHandler } from '$lib/helpers/api_helpers';

type SessionRequestBody = {
	idToken?: string;
	rememberMe?: boolean;
};

type SessionResponseData = {
	redirect_url: string;
	user: User | null;
};

function getDestination(user: User | null): string {
	const roleName = user?.role?.name ?? '';
	const roleLevel = user?.role?.level ?? 0;
	const isAdmin = roleName === 'admin' || roleLevel === 1 || roleLevel === 2;
	if (isAdmin) return '/admin/dashboard';
	if (user) return '/user/account';
	return '/';
}

export const POST: RequestHandler = async (event) => {
	const body = (await event.request.json().catch(() => ({}))) as SessionRequestBody;
	const idToken = typeof body.idToken === 'string' ? body.idToken.trim() : '';
	const rememberMe = !!body.rememberMe;

	if (!idToken) {
		return json(
			{
				status: 400,
				success: false,
				message: 'Missing idToken'
			} satisfies ApiResponse<SessionResponseData>,
			{ status: 400 }
		);
	}

	const api = new ApiClientHandler(event);
	const login = await api.publicRequest<FirebaseAuthResponse>('POST', '/auth/login', { idToken });
	if (!login.success) {
		event.locals.session?.delete('access_token');
		return json(
			{
				status: login.status || 401,
				success: false,
				message: login.message || 'Login failed',
				error: login.error
			} satisfies ApiResponse<SessionResponseData>,
			{ status: login.status || 401 }
		);
	}

	const user = login.data?.user ?? null;
	const maxAge = rememberMe ? 60 * 60 * 24 * 7 : 60 * 60 * 24;
	event.locals.session?.set('access_token', idToken, maxAge);

	return json(
		{
			status: 200,
			success: true,
			message: 'Session created',
			data: {
				redirect_url: getDestination(user),
				user
			}
		} satisfies ApiResponse<SessionResponseData>,
		{ status: 200 }
	);
};
