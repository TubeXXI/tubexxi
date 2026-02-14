import { json } from "@sveltejs/kit";

export const POST = async ({ request, locals }) => {
	try {
		const formData = await request.formData();
		const file = formData.get('file') as File;
		if (!file) {
			return json({ success: false, message: 'Avatar file is required' }, { status: 400 });
		}

		const result = await locals.deps.userService.UpdateAvatar(file);
		if (result instanceof Error) {
			return json({ success: false, message: result.message || 'Failed to update avatar' }, { status: 500 });
		}

		return json({
			success: true,
			message: 'Avatar updated successfully',
			data: { avatar_url: result }
		}, { status: 200 });

	} catch (error) {
		console.error('❌ [Avatar Upload] Unexpected error:', error);
		return json({ success: false, message: 'Failed to update avatar' }, { status: 500 });
	}
}
