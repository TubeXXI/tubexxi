package telegram

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"tubexxi/video-api/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Notifier interface {
	SendAlert(alert AlertRequest)
	Shutdown()
}
type AlertRequest struct {
	Subject  string
	Message  string
	Metadata map[string]interface{}
}
type unifiedNotifier struct {
	telegramConfig *config.TelegramConfig
	workers        int
	queue          chan AlertRequest
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	tgBot          *tgbotapi.BotAPI
	tgEnabled      bool
	logger         *zap.Logger
}

func NewUnifiedNotifier(workers int, queueSize int, Cooldown time.Duration, cfg *config.TelegramConfig, logger *zap.Logger) (Notifier, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg == nil {
		return nil, fmt.Errorf("telegram config is nil")
	}
	un := &unifiedNotifier{
		telegramConfig: cfg,
		workers:        workers,
		queue:          make(chan AlertRequest, queueSize),
		ctx:            ctx,
		cancel:         cancel,
	}

	if cfg.TeleBotToken != "" {
		bot, err := tgbotapi.NewBotAPIWithClient(cfg.TeleBotToken, tgbotapi.APIEndpoint, &http.Client{
			Timeout: 10 * time.Second,
		})
		if err != nil {
			logger.Warn("Telegram notifier disabled - initialization failed",
				zap.Error(err),
				zap.String("token_prefix", safeTokenPrefix(cfg.TeleBotToken)),
			)
			un.tgEnabled = false
		} else {
			un.tgBot = bot
			un.tgEnabled = true
			logger.Info("✅ Telegram notifier initialized",
				zap.String("bot_username", bot.Self.UserName),
			)
		}
	} else {
		logger.Warn("Telegram notifier disabled - no token provided")
	}

	if un.tgEnabled {
		un.wg.Add(workers)
		for i := 0; i < workers; i++ {
			go un.worker()
		}
	} else {
		logger.Warn("No notifiers available - alert system will be disabled")
	}

	return un, nil
}

func (un *unifiedNotifier) SendAlert(alert AlertRequest) {
	select {
	case un.queue <- alert:
	case <-un.ctx.Done():
		un.logger.Debug("Alert system is shutdown")
	}
}
func (un *unifiedNotifier) Shutdown() {
	un.cancel()
	un.wg.Wait()
	un.logger.Info("Alert manager closed successfully")
	un.logger.Info("Telegram Instance closed successfully")
}
func (un *unifiedNotifier) worker() {
	defer un.wg.Done()

	for {
		select {
		case alert := <-un.queue:
			un.processAlert(alert)
		case <-un.ctx.Done():
			return
		}
	}
}
func (un *unifiedNotifier) processAlert(alert AlertRequest) {
	var wg sync.WaitGroup
	if un.tgEnabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			un.sendTelegram(alert)
		}()
	}
	wg.Wait()
}
func (un *unifiedNotifier) sendTelegram(alert AlertRequest) {
	if !un.tgEnabled || un.tgBot == nil {
		return
	}

	msgText := fmt.Sprintf("<b>%s</b>\n%s", alert.Subject, alert.Message)
	if len(alert.Metadata) > 0 {
		msgText += "\n\n<b>Metadata:</b>"
		for k, v := range alert.Metadata {
			msgText += fmt.Sprintf("\n%s: %v", k, v)
		}
	}

	chatID, err := strconv.ParseInt(un.telegramConfig.TeleChatID, 10, 64)
	if err != nil {
		un.logger.Error("Invalid Telegram ChatID", zap.String("chat_id", un.telegramConfig.TeleChatID), zap.Error(err))
		return
	}
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ParseMode = "HTML"

	if _, err := un.tgBot.Send(msg); err != nil {
		un.logger.Error("Failed to send Telegram alert",
			zap.String("subject", alert.Subject),
			zap.Error(err))
	} else {
		un.logger.Debug("Telegram alert sent successfully",
			zap.String("subject", alert.Subject))
	}
}
func safeTokenPrefix(token string) string {
	if len(token) > 5 {
		return token[:3] + "..."
	}
	return "[redacted]"
}
