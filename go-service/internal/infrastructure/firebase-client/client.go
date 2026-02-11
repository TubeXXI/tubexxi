package firebaseclient

import (
	"context"
	"fmt"
	"os"
	"time"
	"tubexxi/video-api/config"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	App        *firebase.App
	AuthClient *auth.Client
	logger     *zap.Logger
	config     *config.AppConfig
}

func NewFirebaseClient(ctx context.Context, logger *zap.Logger, cfg *config.AppConfig) (*FirebaseClient, error) {
	if cfg.FirebaseProjectID == "" {
		return nil, fmt.Errorf("missing FirebaseProjectID in config")
	}

	opts, err := getCredentialsOptions(cfg, logger)
	if err != nil {
		return nil, err
	}

	firebaseConfig := &firebase.Config{
		ProjectID: cfg.FirebaseProjectID,
	}

	app, err := firebase.NewApp(ctx, firebaseConfig, opts...)
	if err != nil {
		logger.Error("Failed to initialize Firebase app", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Error("Failed to initialize Firebase Auth", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize Firebase Auth: %w", err)
	}

	logger.Info("Firebase client initialized successfully",
		zap.String("project_id", cfg.FirebaseProjectID),
		zap.Bool("is_development", cfg.IsDevelopment()),
	)

	return &FirebaseClient{
		App:        app,
		AuthClient: authClient,
		logger:     logger,
		config:     cfg,
	}, nil
}

func getCredentialsOptions(cfg *config.AppConfig, logger *zap.Logger) ([]option.ClientOption, error) {
	var path string
	if cfg.IsDevelopment() {
		possiblePaths := []string{
			"./service-account.json",                    // Root project
			"../service-account.json",                   // Satu level di atas
			"/go-service/service-account.json",          // Docker development
			os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"), // Standard Google env
		}

		for _, p := range possiblePaths {
			if p != "" {
				if _, err := os.Stat(p); err == nil {
					path = p
					logger.Info("Found credentials file", zap.String("path", path))
					break
				}
			}
		}

		if path != "" {
			return []option.ClientOption{option.WithCredentialsFile(path)}, nil
		}
	}

	if !cfg.IsDevelopment() {
		productionPaths := []string{
			"/app/service-account.json",         // Docker production
			"/secrets/service-account.json",     // Kubernetes secrets
			"/run/secrets/service-account.json", // Docker swarm secrets
		}

		for _, p := range productionPaths {
			if _, err := os.Stat(p); err == nil {
				logger.Info("Using production credentials", zap.String("path", p))
				return []option.ClientOption{option.WithCredentialsFile(p)}, nil
			}
		}
	}

	logger.Warn("No explicit credentials found, falling back to Application Default Credentials")
	return []option.ClientOption{}, nil
}

func (fc *FirebaseClient) VerifyIDToken(ctx context.Context, idToken string, checkRevoked bool) (*auth.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var (
		token *auth.Token
		err   error
	)
	if checkRevoked {
		token, err = fc.AuthClient.VerifyIDTokenAndCheckRevoked(ctx, idToken)
	} else {
		token, err = fc.AuthClient.VerifyIDToken(ctx, idToken)
	}
	if err != nil {
		fc.logger.Error("Failed to verify ID token", zap.Error(err))
		return nil, err
	}
	return token, nil
}

func (fc *FirebaseClient) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	user, err := fc.AuthClient.GetUser(ctx, uid)
	if err != nil {
		fc.logger.Error("Failed to get user", zap.String("uid", uid), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (fc *FirebaseClient) CreateUser(ctx context.Context, email string, password string, displayName string) (*auth.UserRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	params := (&auth.UserToCreate{}).Email(email).Password(password)
	if displayName != "" {
		params = params.DisplayName(displayName)
	}

	user, err := fc.AuthClient.CreateUser(ctx, params)
	if err != nil {
		fc.logger.Error("Failed to create Firebase user", zap.String("email", email), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (fc *FirebaseClient) UpdateUserPassword(ctx context.Context, uid string, newPassword string) (*auth.UserRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	user, err := fc.AuthClient.UpdateUser(ctx, uid, (&auth.UserToUpdate{}).Password(newPassword))
	if err != nil {
		fc.logger.Error("Failed to update Firebase user password", zap.String("uid", uid), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (fc *FirebaseClient) EmailVerificationLink(ctx context.Context, email string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	link, err := fc.AuthClient.EmailVerificationLink(ctx, email)
	if err != nil {
		fc.logger.Error("Failed to generate email verification link", zap.String("email", email), zap.Error(err))
		return "", err
	}
	return link, nil
}

func (fc *FirebaseClient) PasswordResetLink(ctx context.Context, email string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	link, err := fc.AuthClient.PasswordResetLink(ctx, email)
	if err != nil {
		fc.logger.Error("Failed to generate password reset link", zap.String("email", email), zap.Error(err))
		return "", err
	}
	return link, nil
}

func (fc *FirebaseClient) Close() error {
	fc.logger.Info("Closing Firebase client")
	return nil
}
