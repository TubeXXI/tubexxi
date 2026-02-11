package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"
	"sync"
	"time"
	"tubexxi/video-api/config"
	"tubexxi/video-api/internal/dto"
	"tubexxi/video-api/internal/entity"
	helpers "tubexxi/video-api/internal/helper"
	"tubexxi/video-api/internal/infrastructure/contextpool"
	minioclient "tubexxi/video-api/internal/infrastructure/minio-client"
	"tubexxi/video-api/internal/infrastructure/repository"
	"tubexxi/video-api/pkg/utils"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo      repository.UserRepository
	userHelper    *helpers.UserHelper
	sessionHelper *helpers.SessionHelper
	logger        *zap.Logger
	minio         *minioclient.MinioClient
	cfg           *config.Config
}

func NewUserService(
	userRepo repository.UserRepository,
	userHelper *helpers.UserHelper,
	sessionHelper *helpers.SessionHelper,
	logger *zap.Logger,
	minio *minioclient.MinioClient,
	cfg *config.Config,
) *UserService {
	return &UserService{
		userRepo:      userRepo,
		userHelper:    userHelper,
		sessionHelper: sessionHelper,
		logger:        logger,
		minio:         minio,
		cfg:           cfg,
	}
}
func (s *UserService) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userIDUUID)
	if err != nil {
		s.logger.Error("Failed to get user",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return nil, dto.ErrUserNotFound
	}
	s.logger.Info("User found",
		zap.String("user_id", userID),
	)

	return user.ToResponse(), nil
}
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByEmail(subCtx, email)
	if err != nil {
		return nil, dto.ErrUserNotFound
	}

	s.logger.Info("User found",
		zap.String("email", email),
	)

	return user.ToResponse(), nil
}
func (s *UserService) Logout(ctx context.Context, userID string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.sessionHelper.DeleteAllUserSessions(subCtx, user.ID); err != nil {
			errChan <- fmt.Errorf("failed to clear auth cache: %w", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.userHelper.ClearAllSessionCache(subCtx, userID); err != nil {
			errChan <- fmt.Errorf("failed to clear user cache: %w", err)
		}
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to logout user: %w", err)
		}
	}

	return nil
}
func (s *UserService) ChangeAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return "", dto.ErrUserNotFound
	}

	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("invalid file type: %s", contentType)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	bucketName := s.cfg.MinIO.MinioBucketName
	objectName := fmt.Sprintf("users/%s-%s", user.ID, file.Filename)

	if user.AvatarURL.Valid && user.AvatarURL.String != "" {
		oldObject := user.AvatarURL.String
		parsedURL, err := url.Parse(oldObject)
		if err == nil {
			path := parsedURL.Path
			path = strings.TrimPrefix(path, "/")

			if strings.HasPrefix(path, s.cfg.MinIO.MinioBucketName+"/") {
				objectName := strings.TrimPrefix(path, s.cfg.MinIO.MinioBucketName+"/")

				if err := s.minio.DeleteFile(subCtx, s.cfg.MinIO.MinioBucketName, objectName); err != nil {
					s.logger.Error("[UserService.ChangeAvatar]", zap.Error(err), zap.String("object", objectName))
				}
			}
		}
	}

	uploadedPath, err := s.minio.UploadFile(subCtx, bucketName, objectName, src, file.Size, contentType)
	if err != nil {
		return "", err
	}
	newAvatarUrl, err := s.userRepo.UpdateAvatar(subCtx, user.ID, uploadedPath)
	if err != nil {
		return "", err
	}

	return newAvatarUrl, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req *dto.UpdateProfileRequest) (*entity.User, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return nil, dto.ErrUserNotFound
	}

	payload := &entity.User{
		ID:              user.ID,
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           entity.NewNullString(req.Phone),
		AvatarURL:       entity.NewNullString(user.AvatarURL.String),
		IsActive:        user.IsActive,
		IsVerified:      user.IsVerified,
		EmailVerifiedAt: user.EmailVerifiedAt,
	}
	userUpdate, err := s.userRepo.Update(subCtx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	userUpdate.PasswordHash = "[Hidden]"
	userUpdate.TwoFaSecret = entity.NewNullString(user.TwoFaSecret.String)

	return userUpdate, nil
}
func (s *UserService) UpdatePassword(ctx context.Context, userID string, req *dto.ChangePasswordRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	if utils.ComparePasswords(req.NewPassword, user.PasswordHash) {
		return fmt.Errorf("new password cannot be the same as the old password")
	}

	if !utils.IsStrongPassword(req.NewPassword) {
		return fmt.Errorf("new password is not strong enough")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(subCtx, userUUID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
func (s *UserService) EnableTwoFactor(ctx context.Context, userID string, req *dto.EnableTwoFactorRequest) (string, string, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return "", "", dto.ErrUserNotFound
	}

	switch req.Status {
	case true:
		if user.TwoFaSecret.Valid && user.TwoFaSecret.String != "" {
			return "", "", fmt.Errorf("two-factor authentication is already enabled")
		}
		s.userHelper.ClearTemp2FASecret(subCtx, userUUID)

		qrcode, secret, err := s.userHelper.Generate2FAQRCode(user.ID, user.FullName)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate two-factor authentication secret: %w", err)
		}
		payload := map[string]interface{}{
			"qr_url":  qrcode,
			"secret":  secret,
			"user_id": user.ID,
		}
		if err := s.userHelper.SetTemp2FASecret(subCtx, userUUID, payload); err != nil {
			return "", "", fmt.Errorf("failed to set two-factor authentication status: %w", err)
		}

		return qrcode, secret, nil
	case false:
		if !user.TwoFaSecret.Valid || user.TwoFaSecret.String == "" {
			return "", "", fmt.Errorf("two-factor authentication is not enabled")
		}
		if err := s.userRepo.UpdateTwoFaSecret(subCtx, userUUID, nil); err != nil {
			return "", "", fmt.Errorf("failed to disable two-factor authentication: %w", err)
		}
		if err := s.userHelper.ClearTemp2FASecret(subCtx, userUUID); err != nil {
			return "", "", fmt.Errorf("failed to clear two-factor authentication status: %w", err)
		}
		return "", "", nil
	default:
		return "", "", fmt.Errorf("invalid status value")
	}
}
func (s *UserService) ActivateTwoFactor(ctx context.Context, userID string, req *dto.ActivateTwoFactorRequest) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	user, err := s.userRepo.FindByID(subCtx, userUUID)
	if err != nil || user == nil {
		return dto.ErrUserNotFound
	}

	if user.TwoFaSecret.Valid || user.TwoFaSecret.String != "" {
		return fmt.Errorf("two-factor authentication already enabled")
	}

	payload, err := s.userHelper.GetTemp2FASecret(subCtx, userUUID)
	if err != nil {
		return fmt.Errorf("failed to get two-factor authentication status: %w", err)
	}
	secret, ok := payload["secret"].(string)
	if !ok || secret == "" {
		return fmt.Errorf("two-factor authentication secret not found")
	}

	valid, err := s.userHelper.Verify2FA(subCtx, userUUID, req.Code, secret)
	if err != nil {
		return fmt.Errorf("failed to verify two-factor authentication token: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid two-factor authentication token")
	}

	if err := s.userRepo.UpdateTwoFaSecret(subCtx, userUUID, &secret); err != nil {
		return fmt.Errorf("failed to update two-factor authentication secret: %w", err)
	}
	if err := s.userHelper.ClearTemp2FASecret(subCtx, userUUID); err != nil {
		return fmt.Errorf("failed to clear two-factor authentication status: %w", err)
	}
	return nil
}
