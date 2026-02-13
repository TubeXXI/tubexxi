import { json } from "@sveltejs/kit";

export const POST = async ({ locals }) => {
	await locals.deps.userService.Logout();

	return json({ success: true, message: 'Logged out successfully' });
};
