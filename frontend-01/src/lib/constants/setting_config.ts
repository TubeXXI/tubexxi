export const defaultSettings: SettingsValue = {
	WEBSITE: {
		site_name: 'Video Downloader',
		site_tagline: 'Download any Videos for Free',
		site_description:
			'Discover a vast collection of videos available for free downloading at Video Downloader. Enjoy the latest blockbusters and timeless classics without any cost.',
		site_keywords: 'Video Downloader, Videos, Free Downloading',
		site_logo: '/images/icon.png',
		site_favicon: '/images/icon.png',
		site_email: 'contact@idvideodownloader.com',
		site_phone: '+1 323 456 7890',
		site_url: 'localhost:5173'
	},
	EMAIL: {
		smtp_service: 'gmail',
		smtp_host: 'smtp.gmail.com',
		smtp_port: 587,
		smtp_user: 'contact@idvideodownloader.com',
		smtp_password: '1234567890abcdef',
		from_email: 'contact@idvideodownloader.com',
		from_name: 'Video Downloader'
	},
	SYSTEM: {
		maintenance_mode: true,
		maintenance_message: 'Video Downloader is currently under maintenance. We will be back soon!',
		source_logo_favicon: 'local'
	},
	MONETIZE: {
		enable_monetize: false,
		type_monetize: 'adsense',
		auto_ad_code: '',
		popup_ad_code: '',
		socialbar_ad_code: '',
		banner_rectangle_ad_code: '',
		banner_horizontal_ad_code: '',
		banner_vertical_ad_code: '',
		native_ad_code: '',
		direct_link_ad_code: ''
	}
};
