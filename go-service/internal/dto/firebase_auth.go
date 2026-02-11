package dto

import "tubexxi/video-api/internal/entity"

type FirebaseLoginRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type FirebaseRegisterRequest struct {
	IDToken   string `json:"id_token" validate:"omitempty"`
	Email     string `json:"email" validate:"omitempty,email"`
	Password  string `json:"password" validate:"omitempty,min=8"`
	FullName  string `json:"full_name" validate:"omitempty,min=2"`
	Phone     string `json:"phone" validate:"omitempty"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}

type FirebaseResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type FirebaseChangePasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type FirebaseVerifyEmailRequest struct {
	Email string `json:"email" validate:"omitempty,email"`
}

type FirebaseAuthResponse struct {
	User          *entity.User `json:"user"`
	FirebaseUID   string       `json:"firebase_uid"`
	Email         string       `json:"email"`
	EmailVerified bool         `json:"email_verified"`
}
