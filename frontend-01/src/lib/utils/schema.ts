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
export const authGoogleSchema = z.object({
	credential: z.string().nonempty('Credential is required')
});
export const forgotSchema = z.object({
	email: z
		.string({ error: 'Email is required' })
		.email('Email is not valid')
		.min(3, 'Email must be at least 3 characters long')
		.nonempty('Email is required')
});
export const resetPasswordSchema = z
	.object({
		new_password: z
			.string()
			.min(6, 'Password must be at least 6 characters')
			.transform((value) => value.replaceAll(/\s+/g, '')),
		confirm_password: z
			.string()
			.nonempty('Confirm password is required')
			.transform((value) => value.replaceAll(/\s+/g, '')),
		token: z.string().nonempty('Token is required')
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

export const contactSchema = z.object({
	name: z.string().nonempty('Name is required'),
	email: z.string().email('Email is not valid').nonempty('Email is required'),
	subject: z.string().nonempty('Subject is required'),
	message: z.string().nonempty('Message is required')
});

export const downloadVideoSchema = z.object({
	url: z.string().nonempty('URL is required'),
	type: z.enum(['youtube', 'tiktok', 'instagram', 'facebook', 'twitter', 'vimeo', 'dailymotion', 'rumble', 'any-video-downloader', 'snackvideo', 'linkedin', 'baidu', 'pinterest', 'snapchat', 'twitch', 'youtube-to-mp3', 'facebook-to-mp3', 'tiktok-to-mp3', 'linkedin-to-mp3', 'snackvideo-to-mp3', 'twitch-to-mp3', 'baidu-to-mp3', 'pinterest-to-mp3', 'snapchat-to-mp3', 'instagram-to-mp3', 'twitter-to-mp3', 'vimeo-to-mp3', 'dailymotion-to-mp3', 'rumble-to-mp3']).default('any-video-downloader'),
	user_id: z.string().optional().or(z.literal('')),
	platform_id: z.string().optional().or(z.literal('')),
	app_id: z.string().optional().or(z.literal(''))
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
})
export const updateSettingAdsTxt = z.object({
	content: z.string().optional().or(z.literal(''))
})
export const updateSettingRobotTxt = z.object({
	content: z.string().optional().or(z.literal(''))
})
export const updateSettingCookie = z.object({
	cookies: z.string().optional().or(z.literal(''))
})


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

// Platforms
export const createPlatformSchema = z.object({
	name: z.string({ error: 'Name is required' }).nonempty('Name is required'),
	slug: z.string({ error: 'Slug is required' }).nonempty('Slug is required'),
	thumbnail_url: z.string({ error: 'Thumbnail URL is required' }).nonempty('Thumbnail URL is required'),
	url_pattern: z.string().optional(),
	is_active: z.boolean().default(true),
	is_premium: z.boolean().default(false),
	config: z.record(z.string(), z.any()).optional().default({})
});

export const updatePlatformSchema = z.object({
	id: z.string({ error: 'ID is required' }).nonempty('ID is required'),
	name: z.string({ error: 'Name is required' }).nonempty('Name is required'),
	slug: z.string({ error: 'Slug is required' }).nonempty('Slug is required'),
	type: z.enum(['youtube', 'tiktok', 'instagram', 'facebook', 'twitter', 'vimeo', 'dailymotion', 'rumble', 'any-video-downloader', 'youtube-to-mp3', 'snackvideo']).optional().default('youtube'),
	url_pattern: z.string().optional(),
	is_active: z.boolean().default(true),
	is_premium: z.boolean().default(false),
	config: z.record(z.string(), z.any()).optional().default({})
});

// Application
export const registerAppSchema = z.object({
	name: z.string({ error: 'Name is required' }).nonempty('Name is required'),
	package_name: z.string({ error: 'Package name is required' }).nonempty('Package name is required'),
	version: z.string({ error: 'Version is required' }).nonempty('Version is required'),
	platform: z.string({ error: 'Platform is required' }).nonempty('Platform is required'),
	enable_monetization: z.boolean().default(false),
	enable_admob: z.boolean().default(false),
	enable_unity_ad: z.boolean().default(false),
	enable_start_app: z.boolean().default(false),
	enable_in_app_purchase: z.boolean().default(false),
	admob_ad_unit_id: z.string().optional(),
	unity_ad_unit_id: z.string().optional(),
	start_app_ad_unit_id: z.string().optional(),
	admob_banner_ad_unit_id: z.string().optional(),
	admob_interstitial_ad_unit_id: z.string().optional(),
	admob_native_ad_unit_id: z.string().optional(),
	admob_rewarded_ad_unit_id: z.string().optional(),
	unity_banner_ad_unit_id: z.string().optional(),
	unity_interstitial_ad_unit_id: z.string().optional(),
	unity_native_ad_unit_id: z.string().optional(),
	unity_rewarded_ad_unit_id: z.string().optional(),
	one_signal_id: z.string().optional(),
	is_active: z.boolean().default(true),
}).superRefine((data, ctx) => {
	if (data.enable_monetization) {
		if (!data.enable_admob && !data.enable_unity_ad && !data.enable_start_app) {
			ctx.addIssue({
				code: z.ZodIssueCode.custom,
				message: 'At least one ad network must be enabled when monetization is enabled',
			});
		}
		if (data.enable_admob) {
			if (!data.admob_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Admob ad unit ID is required when admob is enabled',
				});
			}
			if (!data.admob_banner_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Admob banner ad unit ID is required when admob is enabled',
				});
			}
			if (!data.admob_interstitial_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Admob interstitial ad unit ID is required when admob is enabled',
				});
			}
			if (!data.admob_native_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Admob native ad unit ID is required when admob is enabled',
				});
			}
			if (!data.admob_rewarded_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Admob rewarded ad unit ID is required when admob is enabled',
				});
			}
		}
		if (data.enable_unity_ad) {
			if (!data.unity_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Unity ad unit ID is required when unity ad is enabled',
				});
			}
			if (!data.unity_banner_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Unity banner ad unit ID is required when unity ad is enabled',
				});
			}
			if (!data.unity_interstitial_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Unity interstitial ad unit ID is required when unity ad is enabled',
				});
			}
			if (!data.unity_rewarded_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Unity rewarded ad unit ID is required when unity ad is enabled',
				});
			}
		}
		if (data.enable_start_app) {
			if (!data.start_app_ad_unit_id) {
				ctx.addIssue({
					code: z.ZodIssueCode.custom,
					message: 'Start App ad unit ID is required when Start App is enabled',
				});
			}
		}
	}
});
export const updateApplicationSchema = registerAppSchema.extend({
	id: z.string({ error: 'ID is required' }).nonempty('ID is required'),
});

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
	timestamp_ms: z.string().optional(),
});


export type LoginSchema = z.infer<typeof loginSchema>;
export type AuthGoogleSchema = z.infer<typeof authGoogleSchema>;
export type ForgotSchema = z.infer<typeof forgotSchema>;
export type ResetPasswordSchema = z.infer<typeof resetPasswordSchema>;
export type ContactSchema = z.infer<typeof contactSchema>;
export type DownloadVideoSchema = z.infer<typeof downloadVideoSchema>;

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

// Platforms
export type CreatePlatformSchema = z.infer<typeof createPlatformSchema>;
export type UpdatePlatformSchema = z.infer<typeof updatePlatformSchema>;

// Applications
export type RegisterAppSchema = z.infer<typeof registerAppSchema>;
export type UpdateApplicationSchema = z.infer<typeof updateApplicationSchema>;

// Blog Post
export type PostSchema = z.infer<typeof PostSchema>;
export type WebErrorReportSchema = z.infer<typeof webErrorReportSchema>;
