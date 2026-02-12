package main

import (
	"context"
	"log"
	"time"
	"tubexxi/video-api/internal/dependencies"
	"tubexxi/video-api/internal/entity"
	"tubexxi/video-api/pkg/utils"

	"github.com/google/uuid"
)

func main() {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cont, err := dependencies.NewContainer(ctxTimeout)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}
	defer cont.Close()

	role := []entity.Role{
		{
			ID:          uuid.New(),
			Name:        entity.RoleSuperAdmin,
			Slug:        entity.RoleSuperAdminSlug,
			Level:       entity.RoleLevelSuperAdmin,
			Description: entity.NewNullString("Super admin role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleAdmin,
			Slug:        entity.RoleAdminSlug,
			Level:       entity.RoleLevelAdmin,
			Description: entity.NewNullString("Admin role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        entity.RoleUser,
			Slug:        entity.RoleUserSlug,
			Level:       entity.RoleLevelUser,
			Description: entity.NewNullString("User role"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, r := range role {
		query := `
			INSERT INTO roles (id, name, slug, description, level, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (name) DO NOTHING
		`
		_, err = cont.DBPool.Exec(ctxTimeout, query, r.ID, r.Name, r.Slug, r.Description, r.Level, r.CreatedAt, r.UpdatedAt)
		if err != nil {
			log.Fatalf("Failed to seed role %s: %v", r.Name, err)
		} else {
			log.Printf("Successfully seeded role %s", r.Name)
		}
	}
	hash, err := utils.HashPassword("admin1234")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	adminUser := entity.User{
		ID:           uuid.New(),
		FullName:     "Admin",
		Email:        "id.tubexxi@gmail.com",
		PasswordHash: hash,
		RoleID:       role[0].ID,
		IsActive:     true,
		IsVerified:   true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	insertAdminQuery := `
		INSERT INTO users (id, email, password_hash, full_name, role_id, is_active, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (email) DO NOTHING
	`
	_, err = cont.DBPool.Exec(ctxTimeout, insertAdminQuery,
		adminUser.ID,
		adminUser.Email,
		adminUser.PasswordHash,
		adminUser.FullName,
		adminUser.RoleID,
		adminUser.IsActive,
		adminUser.IsVerified,
		adminUser.CreatedAt,
		adminUser.UpdatedAt,
	)
	if err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
	}
	log.Printf("Successfully seeded admin user %s", adminUser.Email)
}
