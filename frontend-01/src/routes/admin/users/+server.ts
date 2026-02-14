import { json } from '@sveltejs/kit';

export const DELETE = async ({ locals, request }) => {
	try {
		const { ids } = await request.json();
		if (!ids || ids.length === 0) {
			return json({
				success: false,
				message: 'User ID is required'
			}, { status: 400 });
		}
		const response = await locals.deps.adminService.BulkDeleteUser(ids);
		if (response instanceof Error) {
			throw response;
		}
		return json({
			success: true,
			message: 'Users deleted successfully'
		}, { status: 200 });

	} catch (error) {
		return json({
			success: false,
			message: error instanceof Error ? error.message : 'Failed to delete users'
		}, { status: 400 });
	}
}
