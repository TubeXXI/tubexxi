package centrifugo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"tubexxi/video-api/internal/infrastructure/contextpool"

	"github.com/centrifugal/gocent/v3"
	"go.uber.org/zap"
)

func (c *CentrifugoClient) IsUp() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isUp
}

func (c *CentrifugoClient) PublishMessage(ctx context.Context, channel string, data interface{}) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = c.client.Publish(ctx, channel, dataBytes)
	if err != nil {
		c.logger.Error("Failed to publish message to Centrifugo",
			zap.String("channel", channel),
			zap.Error(err),
		)
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}

	c.logger.Debug("Message published to Centrifugo",
		zap.String("channel", channel),
	)
	return nil
}

func (c *CentrifugoClient) PublishToUser(ctx context.Context, userID string, data interface{}) error {
	channel := fmt.Sprintf("user:%s", userID)
	return c.PublishMessage(ctx, channel, data)
}

func (c *CentrifugoClient) PublishToConversation(ctx context.Context, conversationID string, data interface{}) error {
	channel := fmt.Sprintf("conversation:%s", conversationID)
	return c.PublishMessage(ctx, channel, data)
}

func (c *CentrifugoClient) BroadcastNewMessage(ctx context.Context, conversationID string, message map[string]interface{}) error {
	payload := map[string]interface{}{
		"type":    "new_message",
		"message": message,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) BroadcastTypingIndicator(ctx context.Context, conversationID, userID string, isTyping bool) error {
	payload := map[string]interface{}{
		"type":      "typing",
		"user_id":   userID,
		"is_typing": isTyping,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) BroadcastMessageRead(ctx context.Context, conversationID, userID string, messageIDs []string) error {
	payload := map[string]interface{}{
		"type":        "message_read",
		"user_id":     userID,
		"message_ids": messageIDs,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) BroadcastConversationUpdate(ctx context.Context, conversationID string, updates map[string]interface{}) error {
	payload := map[string]interface{}{
		"type":    "conversation_update",
		"updates": updates,
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) BroadcastConversationClosed(ctx context.Context, conversationID, closedBy string) error {
	payload := map[string]interface{}{
		"type":      "conversation_closed",
		"closed_by": closedBy,
		"closed_at": time.Now().Unix(),
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) BroadcastConversationReopened(ctx context.Context, conversationID, reopenedBy string) error {
	payload := map[string]interface{}{
		"type":        "conversation_reopened",
		"reopened_by": reopenedBy,
		"reopened_at": time.Now().Unix(),
	}
	return c.PublishToConversation(ctx, conversationID, payload)
}

func (c *CentrifugoClient) Presence(ctx context.Context, channel string) (*gocent.PresenceResult, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.Presence(ctx, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to get presence for channel %s: %w", channel, err)
	}
	return &result, nil
}

func (c *CentrifugoClient) History(ctx context.Context, channel string, limit int) (*gocent.HistoryResult, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.History(ctx, channel, gocent.WithLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get history for channel %s: %w", channel, err)
	}
	return &result, nil
}

func (c *CentrifugoClient) Channels(ctx context.Context) ([]string, error) {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	result, err := c.client.Channels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channels: %w", err)
	}
	channels := make([]string, 0, len(result.Channels))
	for ch := range result.Channels {
		channels = append(channels, ch)
	}
	return channels, nil
}

func (c *CentrifugoClient) Disconnect(ctx context.Context, userID string, opts ...gocent.DisconnectOption) error {
	ctx, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(ctx, userID, opts...)
	if err != nil {
		return fmt.Errorf("failed to disconnect user %s: %w", userID, err)
	}
	return nil
}

func (c *CentrifugoClient) Refresh(ctx context.Context, userID string, opts ...gocent.DisconnectOption) error {
	_, cancel := contextpool.WithTimeoutIfNone(ctx, 5*time.Second)
	defer cancel()

	err := c.client.Pipe().AddDisconnect(userID, opts...)
	if err != nil {
		return fmt.Errorf("failed to refresh user %s: %w", userID, err)
	}
	return nil
}

func (c *CentrifugoClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isUp = false
	c.logger.Info("Centrifugo client closed")
	return nil
}
