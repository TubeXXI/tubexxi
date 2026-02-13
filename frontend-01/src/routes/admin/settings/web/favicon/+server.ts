import { json } from "@sveltejs/kit";


export const POST = async ({ locals, request }) => {
	const { deps } = locals;
	try {
		const formData = await request.formData();
		const file = formData.get('file') as File;
		if (!file) {
			return json({
				success: false,
				message: 'No file uploaded'
			}, { status: 400 });
		}

		const response = await deps.settingService.updateFavicon(file);
		if (response instanceof Error) {
			throw response;
		}
		return json({
			success: true,
			message: 'Favicon updated successfully',
			data: {
				url: response
			}
		}, { status: 200 });


	} catch (error) {
		return json({
			success: false,
			message: error instanceof Error ? error.message : 'Failed to upload favicon'
		}, { status: 500 });
	}
}
