import { json } from '@sveltejs/kit';

export const DELETE = async ({ locals, params }) => {
	try {
		const { id } = params;
		if (!id) {
			return json({
				success: false,
				message: 'User ID is required'
			}, { status: 400 });
		}
		const response = await locals.deps.adminService.HardDeleteUser(id);
		if (response instanceof Error) {
			throw response;
		}
		return json({
			success: true,
			message: 'User deleted successfully'
		}, { status: 200 });

	} catch (error) {
		return json({
			success: false,
			message: error instanceof Error ? error.message : 'Failed to delete user'
		}, { status: 400 });
	}
}
