package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/entity"
	helpers "tubexxi/video-api/internal/helper"
	firebaseclient "tubexxi/video-api/internal/infrastructure/firebase-client"
	"tubexxi/video-api/internal/infrastructure/repository"
	"tubexxi/video-api/pkg/utils"
)

type AuthService struct {
	logger   *zap.Logger
	firebase *firebaseclient.FirebaseClient
	userRepo repository.UserRepository
	mail     *helpers.MailHelper
}

func NewAuthService(logger *zap.Logger, firebase *firebaseclient.FirebaseClient, userRepo repository.UserRepository, mail *helpers.MailHelper) *AuthService {
	return &AuthService{logger: logger, firebase: firebase, userRepo: userRepo, mail: mail}
}

func (s *AuthService) LoginWithIDToken(ctx context.Context, idToken string) (*dto.FirebaseAuthResponse, error) {
	tok, err := s.firebase.VerifyIDToken(ctx, idToken, true)
	if err != nil {
		return nil, fmt.Errorf("invalid firebase token")
	}

	email, _ := tok.Claims["email"].(string)
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("email claim not found")
	}

	name, _ := tok.Claims["name"].(string)
	picture, _ := tok.Claims["picture"].(string)
	phone, _ := tok.Claims["phone_number"].(string)
	emailVerified, _ := tok.Claims["email_verified"].(bool)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		randomPassword := utils.GenerateRandomPassword()
		hash, hashErr := utils.HashPassword(randomPassword)
		if hashErr != nil {
			return nil, fmt.Errorf("failed to hash password")
		}

		nm := name
		if nm == "" {
			nm = "User"
		}

		u := &entity.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: hash,
			FullName:     nm,
			IsActive:     true,
			IsVerified:   emailVerified,
			CreatedAt:    time.Now(),
		}
		u.Phone = entity.NewNullString(phone)
		u.AvatarURL = entity.NewNullString(picture)

		if createErr := s.userRepo.CreateWithRecovery(ctx, u); createErr != nil {
			return nil, createErr
		}
		user = u
	}

	if emailVerified && !user.IsVerified {
		_ = s.userRepo.SetEmailVerified(ctx, user.ID, true)
	}

	return &dto.FirebaseAuthResponse{
		User:          user,
		FirebaseUID:   tok.UID,
		Email:         email,
		EmailVerified: emailVerified,
	}, nil
}

func (s *AuthService) RegisterWithEmail(ctx context.Context, req *dto.FirebaseRegisterRequest) (*dto.FirebaseAuthResponse, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	if !utils.IsStrongPassword(req.Password) {
		return nil, fmt.Errorf("password is not strong")
	}

	displayName := strings.TrimSpace(req.FullName)
	_, err := s.firebase.CreateUser(ctx, req.Email, req.Password, displayName)
	if err != nil {
		return nil, err
	}

	verifyLink, err := s.firebase.EmailVerificationLink(ctx, req.Email)
	if err == nil {
		oobCode := extractOobCode(verifyLink)
		if oobCode != "" {
			_ = s.mail.SendEmail(&dto.SendMailMetaData{
				Token: oobCode,
				Type:  dto.EmailVerification,
				To:    req.Email,
			})
		}
	}

	user, findErr := s.userRepo.FindByEmail(ctx, req.Email)
	if findErr != nil {
		hash, hashErr := utils.HashPassword(utils.GenerateRandomPassword())
		if hashErr != nil {
			return nil, fmt.Errorf("failed to hash password")
		}
		nm := displayName
		if nm == "" {
			nm = "User"
		}
		u := &entity.User{
			ID:           uuid.New(),
			Email:        req.Email,
			PasswordHash: hash,
			FullName:     nm,
			IsActive:     true,
			IsVerified:   false,
			CreatedAt:    time.Now(),
		}
		u.Phone = entity.NewNullString(req.Phone)
		u.AvatarURL = entity.NewNullString(req.AvatarURL)
		if createErr := s.userRepo.CreateWithRecovery(ctx, u); createErr != nil {
			return nil, createErr
		}
		user = u
	}

	return &dto.FirebaseAuthResponse{User: user, Email: user.Email, EmailVerified: user.IsVerified}, nil
}

func (s *AuthService) SendResetPassword(ctx context.Context, email string) error {
	link, err := s.firebase.PasswordResetLink(ctx, email)
	if err != nil {
		return err
	}
	oobCode := extractOobCode(link)
	if oobCode == "" {
		return fmt.Errorf("failed to extract oobCode")
	}
	return s.mail.SendEmail(&dto.SendMailMetaData{Token: oobCode, Type: dto.ResetPassword, To: email})
}

func (s *AuthService) SendVerifyEmail(ctx context.Context, email string) error {
	link, err := s.firebase.EmailVerificationLink(ctx, email)
	if err != nil {
		return err
	}
	oobCode := extractOobCode(link)
	if oobCode == "" {
		return fmt.Errorf("failed to extract oobCode")
	}
	return s.mail.SendEmail(&dto.SendMailMetaData{Token: oobCode, Type: dto.EmailVerification, To: email})
}

func (s *AuthService) ChangePassword(ctx context.Context, firebaseUID string, userID string, newPassword string) error {
	if !utils.IsStrongPassword(newPassword) {
		return fmt.Errorf("password is not strong")
	}

	_, err := s.firebase.UpdateUserPassword(ctx, firebaseUID, newPassword)
	if err != nil {
		return err
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}
	return s.userRepo.UpdatePassword(ctx, uid, passwordHash)
}

func extractOobCode(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}
	return u.Query().Get("oobCode")
}
