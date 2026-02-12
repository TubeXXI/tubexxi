package dto

type SetUserRoleRequest struct {
	UserID       string `json:"user_id" validate:"omitempty,uuid"`
	Email        string `json:"email" validate:"omitempty,email"`
	Role         string `json:"role" validate:"required,oneof=user admin superadmin"`
	SyncFirebase bool   `json:"sync_firebase" validate:"omitempty"`
}
