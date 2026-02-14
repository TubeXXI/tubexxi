package service

import (
	"context"
	"tubexxi/video-api/internal/dto"
	helpers "tubexxi/video-api/internal/helper"

	"go.uber.org/zap"
)

type ClientService struct {
	mailHelper *helpers.MailHelper
	logger     *zap.Logger
}

func NewClientService(
	mailHelper *helpers.MailHelper,
	logger *zap.Logger,
) *ClientService {
	return &ClientService{
		mailHelper: mailHelper,
		logger:     logger,
	}
}

func (s *ClientService) SendContactEmail(ctx context.Context, payload *dto.ContactRequest, clientOrigin string) error {
	return s.mailHelper.SendContactEmail(ctx, payload, clientOrigin)
}
