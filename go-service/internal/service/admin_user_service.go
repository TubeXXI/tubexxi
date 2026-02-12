package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"tubexxi/video-api/internal/dto"
	firebaseclient "tubexxi/video-api/internal/infrastructure/firebase-client"
	"tubexxi/video-api/internal/infrastructure/repository"
)

type AdminUserService struct {
	logger   *zap.Logger
	userRepo repository.UserRepository
	firebase *firebaseclient.FirebaseClient
}

func NewAdminUserService(logger *zap.Logger, userRepo repository.UserRepository, firebase *firebaseclient.FirebaseClient) *AdminUserService {
	return &AdminUserService{logger: logger, userRepo: userRepo, firebase: firebase}
}

func (s *AdminUserService) SetUserRole(ctx context.Context, req *dto.SetUserRoleRequest) error {
	roleName := strings.ToLower(strings.TrimSpace(req.Role))
	role, err := s.userRepo.FindRoleByName(ctx, roleName)
	if err != nil {
		return err
	}

	var userID uuid.UUID
	if req.UserID != "" {
		id, err := uuid.Parse(req.UserID)
		if err != nil {
			return fmt.Errorf("invalid user_id")
		}
		userID = id
	} else if req.Email != "" {
		user, err := s.userRepo.FindByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		userID = user.ID
	} else {
		return fmt.Errorf("user_id or email is required")
	}

	if err := s.userRepo.SetRoleID(ctx, userID, role.ID); err != nil {
		return err
	}

	if req.SyncFirebase {
		var email string
		if req.Email != "" {
			email = req.Email
		} else {
			user, err := s.userRepo.FindByID(ctx, userID)
			if err != nil {
				return err
			}
			email = user.Email
		}

		fbUser, err := s.firebase.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}
		claims := map[string]interface{}{
			"role":       role.Name,
			"role_level": int(role.Level),
			"admin":      role.Name == "admin" || role.Name == "superadmin",
		}
		if err := s.firebase.SetCustomUserClaims(ctx, fbUser.UID, claims); err != nil {
			return err
		}
	}

	return nil
}
