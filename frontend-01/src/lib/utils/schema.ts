import { z } from 'zod';

const isLetter = (char: string) => /^[a-zA-Z]$/.test(char);

export const loginSchema = z.object({
	email: z
		.string({ error: 'Email is required' })
		.email({ error: 'Invalid email address' })
		.refine((value) => isLetter(value[0]), {
			message: 'Email must start with a letter'
		})
		.nonempty({ error: 'Email is required' }),
	password: z
		.string({ error: 'Password is required' })
		.min(6)
		.transform((value) => value.replaceAll(/\s+/g, '')),
	remember_me: z.boolean().default(false)
});
export const registerSchema = z.object({
	id_token: z.string().nonempty({ error: 'Please is required' }),
	email: z
		.string()
		.email({ error: 'Email is not valid' })
		.nonempty({ error: 'Please is required' }),
	password: z
		.string()
		.min(6, 'Password must be at least 6 characters')
		.nonempty({ error: 'Please is required' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	full_name: z.string().nonempty({ error: 'Please is required' }),
	phone: z.string().optional(),
	avatar_url: z.string().optional()
});
export const resetPasswordSchema = z.object({
	email: z
		.string({ error: 'Email is required' })
		.email('Email is not valid')
		.min(3, 'Email must be at least 3 characters long')
		.nonempty('Email is required')
});
export const changePasswordSchema = z
	.object({
		new_password: z
			.string()
			.min(6, 'Password must be at least 6 characters')
			.transform((value) => value.replaceAll(/\s+/g, '')),
		confirm_password: z
			.string()
			.nonempty('Confirm password is required')
			.transform((value) => value.replaceAll(/\s+/g, ''))
	})
	.superRefine((data, ctx) => {
		if (data.new_password != data.confirm_password) {
			ctx.addIssue({
				path: ['confirm_password'],
				code: z.ZodIssueCode.custom,
				message: 'Password and confirm password must be the same'
			});
		}
	});
export const verifyEmailSchema = z.object({
	email: z
		.string({ error: 'Email is required' })
		.email('Email is not valid')
		.min(3, 'Email must be at least 3 characters long')
		.nonempty('Email is required')
});

export const contactSchema = z.object({
	name: z.string().nonempty('Name is required'),
	email: z.string().email('Email is not valid').nonempty('Email is required'),
	subject: z.string().nonempty('Subject is required'),
	message: z.string().nonempty('Message is required')
});

export const downloadVideoSchema = z.object({
	url: z.string().nonempty('URL is required'),
	type: z
		.enum([
			'youtube',
			'tiktok',
			'instagram',
			'facebook',
			'twitter',
			'vimeo',
			'dailymotion',
			'rumble',
			'any-video-downloader',
			'snackvideo',
			'linkedin',
			'baidu',
			'pinterest',
			'snapchat',
			'twitch',
			'youtube-to-mp3',
			'facebook-to-mp3',
			'tiktok-to-mp3',
			'linkedin-to-mp3',
			'snackvideo-to-mp3',
			'twitch-to-mp3',
			'baidu-to-mp3',
			'pinterest-to-mp3',
			'snapchat-to-mp3',
			'instagram-to-mp3',
			'twitter-to-mp3',
			'vimeo-to-mp3',
			'dailymotion-to-mp3',
			'rumble-to-mp3'
		])
		.default('any-video-downloader'),
	user_id: z.string().optional().or(z.literal('')),
	platform_id: z.string().optional().or(z.literal('')),
	app_id: z.string().optional().or(z.literal(''))
});

// Admin
export const setRoleSchema = z.object({
	user_id: z.string().nonempty('User ID is required'),
	email: z.string().email('Email is not valid').nonempty('Email is required'),
	role: z.enum(['user', 'admin', 'superadmin']).default('user'),
	sync_firebase: z.boolean().optional().default(true)
});

// Website
export const updateSettingWeb = z.object({
	site_name: z.string().optional(),
	site_tagline: z.string().optional(),
	site_description: z.string().optional(),
	site_keywords: z.string().optional(),
	site_email: z.string().email().optional(),
	site_phone: z.string().optional(),
	site_url: z.string().url().optional()
});
export const updateSettingEmail = z.object({
	smtp_enabled: z.boolean().optional().default(true),
	smtp_service: z.string().optional().or(z.literal('')).default('gmail'),
	smtp_host: z.string().optional().or(z.literal('')).default('smtp.gmail.com'),
	smtp_port: z.number().optional().default(587),
	smtp_user: z.string().optional().or(z.literal('')),
	smtp_password: z.string().optional().or(z.literal('')),
	from_email: z.string().email().optional().or(z.literal('')),
	from_name: z.string().optional().or(z.literal(''))
});
export const updateSettingSystem = z.object({
	theme: z.enum(['default', 'red-dark', 'green-dark', 'blue-dark']).optional().default('default'),
	enable_documentation: z.boolean().optional().default(true),
	maintenance_mode: z.boolean().optional().default(false),
	maintenance_message: z.string().optional().or(z.literal('')),
	source_logo_favicon: z.enum(['local', 'remote']).optional().default('local'),
	histats_tracking_code: z.string().optional().or(z.literal('')),
	google_analytics_code: z.string().optional().or(z.literal('')),
	play_store_app_url: z.string().optional().or(z.literal('')),
	app_store_app_url: z.string().optional().or(z.literal(''))
});
export const updateSettingMonetization = z.object({
	enable_monetize: z.boolean().optional().default(false),
	type_monetize: z.enum(['adsense', 'revenuecat', 'adsterra']).optional().default('adsense'),
	publisher_id: z.string().optional().or(z.literal('')),
	enable_popup_ad: z.boolean().optional().default(false),
	auto_ad_code: z.string().optional().or(z.literal('')),
	popup_ad_code: z.string().optional().or(z.literal('')),
	socialbar_ad_code: z.string().optional().or(z.literal('')),
	banner_rectangle_ad_code: z.string().optional().or(z.literal('')),
	banner_horizontal_ad_code: z.string().optional().or(z.literal('')),
	banner_vertical_ad_code: z.string().optional().or(z.literal('')),
	native_ad_code: z.string().optional().or(z.literal('')),
	direct_link_ad_code: z.string().optional().or(z.literal(''))
});
export const updateSettingAdsTxt = z.object({
	content: z.string().optional().or(z.literal(''))
});
export const updateSettingRobotTxt = z.object({
	content: z.string().optional().or(z.literal(''))
});
export const updateSettingCookie = z.object({
	cookies: z.string().optional().or(z.literal(''))
});

// Account
export const updateProfileSchema = z.object({
	full_name: z
		.string({ error: 'Name is required' })
		.min(3, 'Name must be at least 3 characters long')
		.nonempty('Name is required'),
	email: z
		.string({ error: 'Email is required' })
		.email('Email is not valid')
		.nonempty('Email is required'),
	phone: z.string().optional().or(z.literal('')).default('')
});
export const updatePasswordSchema = z.object({
	current_password: z
		.string({ error: 'Current password is required' })
		.min(1, { message: 'Current password is required' })
		.min(6, { message: 'Current password must be at least 6 characters long' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	new_password: z
		.string({ error: 'New password is required' })
		.min(6, { message: 'New password must be at least 6 characters long' })
		.transform((value) => value.replaceAll(/\s+/g, '')),
	confirm_password: z
		.string({ error: 'Confirm password is required' })
		.nonempty({ message: 'Confirm password is required' })
		.transform((value) => value.replaceAll(/\s+/g, ''))
});

// Application
export const registerAppSchema = z
	.object({
		name: z.string({ error: 'Name is required' }).nonempty('Name is required'),
		package_name: z
			.string({ error: 'Package name is required' })
			.nonempty('Package name is required'),
		version: z.string({ error: 'Version is required' }).nonempty('Version is required'),
		type: z.enum(['android', 'ios']).default('android'),
		store_url: z.string().optional().or(z.literal('')),
		is_active: z.boolean().default(true),
		enable_monetize: z.boolean().default(false),
		enable_admob: z.boolean().default(false),
		enable_unity_ad: z.boolean().default(false),
		enable_star_io_ad: z.boolean().default(false),
		enable_in_app_purchase: z.boolean().default(false),
		admob_id: z.string().optional(),
		unity_ad_id: z.string().optional(),
		star_io_ad_id: z.string().optional(),
		admob_auto_ad: z.string().optional(),
		admob_banner_ad: z.string().optional(),
		admob_interstitial_ad: z.string().optional(),
		admob_native_ad: z.string().optional(),
		admob_rewarded_ad: z.string().optional(),
		unity_banner_ad: z.string().optional(),
		unity_interstitial_ad: z.string().optional(),
		unity_rewarded_ad: z.string().optional(),
		one_signal_id: z.string().optional(),
	})

export const updateApplicationSchema = registerAppSchema;

// Blog Post
export const PostSchema = z.object({
	title: z.string(),
	slug: z.string().nonempty('Slug is required'),
	thumbnail: z.string().nonempty('Thumbnail is required'),
	description: z.string(),
	publishedDate: z.string(),
	lastUpdatedDate: z.string().optional(),
	tags: z.array(z.string()).optional(),
	status: z.enum(['draft', 'published']),
	series: z
		.object({
			order: z.number(),
			title: z.string()
		})
		.optional()
});

export const webErrorReportSchema = z.object({
	error: z.string().optional(),
	message: z.string().optional(),
	platform_id: z.string().optional(),
	user_id: z.string().optional(),
	ip_address: z.string().optional(),
	user_agent: z.string().optional(),
	url: z.string().optional(),
	method: z.string().optional(),
	request: z.string().optional(),
	status: z.number().optional(),
	level: z.string().optional(),
	locale: z.string().optional(),
	timestamp_ms: z.string().optional()
});

export type LoginSchema = z.infer<typeof loginSchema>;
export type RegisterSchema = z.infer<typeof registerSchema>;
export type ResetPasswordSchema = z.infer<typeof resetPasswordSchema>;
export type ChangePasswordSchema = z.infer<typeof changePasswordSchema>;
export type VerifyEmailSchema = z.infer<typeof verifyEmailSchema>;
export type ContactSchema = z.infer<typeof contactSchema>;
export type DownloadVideoSchema = z.infer<typeof downloadVideoSchema>;

// Admin
export type SetRoleSchema = z.infer<typeof setRoleSchema>;

// Settings
export type UpdateSettingWebSchema = z.infer<typeof updateSettingWeb>;
export type UpdateSettingEmailSchema = z.infer<typeof updateSettingEmail>;
export type UpdateSettingSystemSchema = z.infer<typeof updateSettingSystem>;
export type UpdateSettingMonetizationSchema = z.infer<typeof updateSettingMonetization>;
export type UpdateSettingRobotTxtSchema = z.infer<typeof updateSettingRobotTxt>;
export type UpdateSettingAdsTxtSchema = z.infer<typeof updateSettingAdsTxt>;
export type UpdateSettingCookieSchema = z.infer<typeof updateSettingCookie>;

// Accounts
export type UpdateProfileSchema = z.infer<typeof updateProfileSchema>;
export type UpdatePasswordSchema = z.infer<typeof updatePasswordSchema>;
// Applications
export type RegisterAppSchema = z.infer<typeof registerAppSchema>;
export type UpdateApplicationSchema = z.infer<typeof updateApplicationSchema>;

// Blog Post
export type PostSchema = z.infer<typeof PostSchema>;
export type WebErrorReportSchema = z.infer<typeof webErrorReportSchema>;
