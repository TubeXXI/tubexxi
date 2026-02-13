import { json } from '@sveltejs/kit';

export const GET = async ({ url, locals }) => {
	const page = Number(url.searchParams.get('page')) || 1;
	const limit = Number(url.searchParams.get('limit')) || 50;

	const { deps } = locals;
	const logs = await deps.serverStatusService.GetServerLogs(page, limit);

	return json(logs);
};

export const DELETE = async ({ locals }) => {
	const { deps } = locals;
	try {
		const response = await deps.serverStatusService.ClearServerLogs();
		if (response instanceof Error) {
			throw response;
		}
		return json({
			success: true,
			message: 'Server logs cleared successfully'
		},
			{ status: 200 }
		);
	} catch (error) {
		return json({
			success: false,
			message: error instanceof Error ? error.message : 'Failed to clear server logs'
		},
			{ status: 500 }
		);
	}
}
