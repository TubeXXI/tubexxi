import { json } from "@sveltejs/kit";
import * as i18n from '@/paraglide/messages.js';

export const POST = async ({ locals, request }) => {
	const { deps } = locals;

	try {
		const { email } = await request.json()
		if (!email) {
			return json(
				{
					success: false,
					message: i18n.page_verify_email_error_message(),
				},
				{ status: 400 }
			);
		}
		const response = await deps.authService.VerifyEmail(email) as string | Error
		if (response instanceof Error) {
			return json(
				{
					success: false,
					message: response.message || i18n.page_verify_email_error_message(),
				},
				{ status: 400 }
			);
		}

		return json({
			success: false,
			message: i18n.page_verify_email_success_message(),
		}, {
			status: 200
		})

	} catch (error) {
		return json(
			{
				success: false,
				message: error instanceof Error ? error.message : i18n.page_verify_email_error_message(),
			},
			{ status: 500 }
		);
	}

}
