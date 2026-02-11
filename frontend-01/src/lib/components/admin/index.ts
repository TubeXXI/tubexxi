import AdminSidebar from './AdminSidebar.svelte';
import AdminSidebarLayout from './AdminSidebarLayout.svelte';
import AdminSidebarHeader from './AdminSidebarHeader.svelte';
import AdminNavMain from './AdminNavMain.svelte';
import AdminNavUser from './AdminNavUser.svelte';
import AdminNavSetting from './AdminNavSetting.svelte';
import AdminNavBottom from './AdminNavBottom.svelte';
import AdminHeading from './AdminHeading.svelte';
import AdminBulkActionFloating from './AdminBulkActionFloating.svelte';
import AdminDeleteDialog from './AdminDeleteDialog.svelte';
import AdminDateToggle from './AdminDateToggle.svelte';
import AdminDashboardFilterToolbar from './dashboard/AdminDashboardFilterToolbar.svelte';

// Settings
import AdminSettingLayout from './settings/AdminSettingLayout.svelte';
import AdminSettingSidenav from './settings/AdminSettingSidenav.svelte';
import AdminSettingUploadFavicon from './settings/AdminSettingUploadFavicon.svelte';
import AdminSettingUploadLogo from './settings/AdminSettingUploadLogo.svelte';

// Accounts
import AdminAccountLayout from './accounts/AdminAccountLayout.svelte';
import AdminAccountUploadAvatar from './accounts/AdminAccountUploadAvatar.svelte';

// Dashboard
import AdminDashboardChartBarInteractive from './dashboard/AdminDashboardChartBarInteractive.svelte';
import AdminDashboardRecentDownload from './dashboard/AdminDashboardRecentDownload.svelte';

// platforms
import AdminPlatformTable from './platform/AdminPlatformTable.svelte';
import AdminPlatformTableToolbar from './platform/AdminPlatformTableToolbar.svelte';
import AdminPlatformUploadThumbnail from './platform/AdminPlatformUploadThumbnail.svelte';

// applications
import AdminApplicationTable from './application/AdminApplicationTable.svelte';
import AdminApplicationTableToolbar from './application/AdminApplicationTableToolbar.svelte';

// downloads
import AdminDownloadTable from './download/AdminDownloadTable.svelte';
import AdminDownloadTableToolbar from './download/AdminDownloadTableToolbar.svelte';

// subscriptions
import AdminSubscriptionTable from './subscription/AdminSubscriptionTable.svelte';
import AdminSubscriptionTableToolbar from './subscription/AdminSubscriptionTableToolbar.svelte';
import AdminSubscriptionDetail from './subscription/AdminSubscriptionDetail.svelte';

// users
import AdminUserTable from './user/AdminUserTable.svelte';
import AdminUserTableToolbar from './user/AdminUserTableToolbar.svelte';
import AdminUserDetail from './user/AdminUserDetail.svelte';

export {
	AdminSidebar,
	AdminSidebarLayout,
	AdminSidebarHeader,
	AdminNavMain,
	AdminNavUser,
	AdminNavSetting,
	AdminNavBottom,
	AdminHeading,
	AdminBulkActionFloating,
	AdminDeleteDialog,
	AdminDateToggle,
	AdminDashboardFilterToolbar,
	// Settings
	AdminSettingLayout,
	AdminSettingSidenav,
	AdminSettingUploadFavicon,
	AdminSettingUploadLogo,
	// Accounts
	AdminAccountLayout,
	AdminAccountUploadAvatar,
	// Dashboard
	AdminDashboardChartBarInteractive,
	AdminDashboardRecentDownload,
	// platforms
	AdminPlatformTable,
	AdminPlatformTableToolbar,
	AdminPlatformUploadThumbnail,
	// applications
	AdminApplicationTable,
	AdminApplicationTableToolbar,
	// downloads
	AdminDownloadTable,
	AdminDownloadTableToolbar,
	// subscriptions
	AdminSubscriptionTable,
	AdminSubscriptionTableToolbar,
	AdminSubscriptionDetail,
	// users
	AdminUserTable,
	AdminUserTableToolbar,
	AdminUserDetail,
};
