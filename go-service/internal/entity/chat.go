package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	SenderID        uuid.UUID       `json:"sender_id" db:"sender_id"`
	ReceiverID      uuid.UUID       `json:"receiver_id" db:"receiver_id"`
	Message         *string         `json:"message,omitempty" db:"message"`
	Type            string          `json:"type" db:"type"` // text, image, video, audio, file, location, contact
	FileURL         *string         `json:"file_url,omitempty" db:"file_url"`
	FileName        *string         `json:"file_name,omitempty" db:"file_name"`
	FileSize        *int64          `json:"file_size,omitempty" db:"file_size"`
	MimeType        *string         `json:"mime_type,omitempty" db:"mime_type"`
	IsRead          bool            `json:"is_read" db:"is_read"`
	IsDelivered     bool            `json:"is_delivered" db:"is_delivered"`
	IsEdited        bool            `json:"is_edited" db:"is_edited"`
	IsDeleted       bool            `json:"is_deleted" db:"is_deleted"`
	ReplyToID       *uuid.UUID      `json:"reply_to_id,omitempty" db:"reply_to_id"`
	ForwardedFromID *uuid.UUID      `json:"forwarded_from_id,omitempty" db:"forwarded_from_id"`
	Metadata        json.RawMessage `json:"metadata,omitempty" db:"metadata"` // For location, contact, etc.
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
	ReplyTo         *Chat           `json:"reply_to,omitempty" db:"-"`
	ForwardedFrom   *Chat           `json:"forwarded_from,omitempty" db:"-"`
	Replies         []*Chat         `json:"replies,omitempty" db:"-"`
	Sender          *User           `json:"sender,omitempty" db:"-"`
	Receiver        *User           `json:"receiver,omitempty" db:"-"`
	Reactions       []*ChatReaction `json:"reactions,omitempty" db:"-"`
	ReactionsCount  int64           `json:"reactions_count" db:"reactions_count"`
	RepliesCount    int64           `json:"replies_count" db:"replies_count"`
	IsSentByMe      bool            `json:"is_sent_by_me" db:"-"`
	IsForwarded     bool            `json:"is_forwarded" db:"-"`
	SenderName      string          `json:"sender_name,omitempty" db:"sender_name"`
	SenderAvatar    *string         `json:"sender_avatar,omitempty" db:"sender_avatar"`
	ReceiverName    string          `json:"receiver_name,omitempty" db:"receiver_name"`
	ReceiverAvatar  *string         `json:"receiver_avatar,omitempty" db:"receiver_avatar"`
}

type ChatReaction struct {
	ID           uuid.UUID `json:"id" db:"id"`
	ChatID       uuid.UUID `json:"chat_id" db:"chat_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	ReactionType string    `json:"reaction_type" db:"reaction_type"` // like, love, haha, wow, sad, angry
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	User         *User     `json:"user,omitempty" db:"-"`
	Chat         *Chat     `json:"chat,omitempty" db:"-"`
}

type ChatMedia struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChatID    uuid.UUID `json:"chat_id" db:"chat_id"`
	URL       string    `json:"url" db:"url"`
	Type      string    `json:"type" db:"type"`
	Name      string    `json:"name" db:"name"`
	Size      int64     `json:"size" db:"size"`
	MimeType  string    `json:"mime_type" db:"mime_type"`
	Width     *int      `json:"width,omitempty" db:"width"`
	Height    *int      `json:"height,omitempty" db:"height"`
	Duration  *int      `json:"duration,omitempty" db:"duration"` // for video/audio
	Thumbnail *string   `json:"thumbnail,omitempty" db:"thumbnail"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ChatCreateRequest struct {
	ReceiverID      *uuid.UUID        `json:"receiver_id,omitempty"`
	GroupID         *uuid.UUID        `json:"group_id,omitempty"`
	Message         *string           `json:"message,omitempty"`
	Type            string            `json:"type" validate:"required,oneof=text image video audio file location contact"`
	ReplyToID       *uuid.UUID        `json:"reply_to_id,omitempty"`
	ForwardedFromID *uuid.UUID        `json:"forwarded_from_id,omitempty"`
	Media           *ChatMediaRequest `json:"media,omitempty"`
	Location        *LocationData     `json:"location,omitempty"`
	Contact         *ContactData      `json:"contact,omitempty"`
	Metadata        json.RawMessage   `json:"metadata,omitempty"`
}

type ChatMediaRequest struct {
	URL       string  `json:"url"`
	Name      string  `json:"name"`
	Size      int64   `json:"size"`
	MimeType  string  `json:"mime_type"`
	Width     *int    `json:"width,omitempty"`
	Height    *int    `json:"height,omitempty"`
	Duration  *int    `json:"duration,omitempty"`
	Thumbnail *string `json:"thumbnail,omitempty"`
}

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      *string `json:"name,omitempty"`
	Address   *string `json:"address,omitempty"`
}

type ContactData struct {
	Name    string  `json:"name"`
	Phone   string  `json:"phone"`
	Email   *string `json:"email,omitempty"`
	Address *string `json:"address,omitempty"`
}

type ChatUpdateRequest struct {
	Message *string `json:"message" validate:"required"`
}

type ChatFilter struct {
	UserID      uuid.UUID  `json:"user_id" validate:"required"`
	OtherUserID *uuid.UUID `json:"other_user_id,omitempty"`
	GroupID     *uuid.UUID `json:"group_id,omitempty"`
	Type        string     `json:"type,omitempty"`
	IsRead      *bool      `json:"is_read,omitempty"`
	IsDelivered *bool      `json:"is_delivered,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Search      *string    `json:"search,omitempty"`
	Limit       int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset      int        `json:"offset" validate:"omitempty,min=0"`
	SortBy      string     `json:"sort_by" validate:"omitempty,oneof=created_at is_read"`
	SortOrder   string     `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type ChatResponse struct {
	*Chat
	ReplyTo       *ChatResponse   `json:"reply_to,omitempty"`
	ForwardedFrom *ChatResponse   `json:"forwarded_from,omitempty"`
	Reactions     []*ChatReaction `json:"reactions,omitempty"`
	Media         *ChatMedia      `json:"media,omitempty"`
}

type Conversation struct {
	OtherUserID     uuid.UUID `json:"other_user_id" db:"other_user_id"`
	OtherUserName   string    `json:"other_user_name" db:"other_user_name"`
	OtherUserAvatar *string   `json:"other_user_avatar,omitempty" db:"other_user_avatar"`
	LastMessage     *Chat     `json:"last_message" db:"-"`
	UnreadCount     int64     `json:"unread_count" db:"unread_count"`
	TotalMessages   int64     `json:"total_messages" db:"total_messages"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ChatThread struct {
	RootMessage  *Chat                      `json:"root_message"`
	Replies      []*Chat                    `json:"replies"`
	Participants []*User                    `json:"participants"`
	Reactions    map[string][]*ChatReaction `json:"reactions"`
}

func (c *Chat) IsRoot() bool {
	return c.ReplyToID == nil
}

func (c *Chat) IsForwardedMessage() bool {
	return c.ForwardedFromID != nil
}

func (c *Chat) HasReplies() bool {
	return c.RepliesCount > 0 && c.Replies != nil && len(c.Replies) > 0
}

func (c *Chat) HasReactions() bool {
	return c.ReactionsCount > 0 && c.Reactions != nil && len(c.Reactions) > 0
}

func (c *Chat) AddReply(reply *Chat) {
	if c.Replies == nil {
		c.Replies = make([]*Chat, 0)
	}
	c.Replies = append(c.Replies, reply)
	c.RepliesCount++
}

func (c *Chat) AddReaction(reaction *ChatReaction) {
	if c.Reactions == nil {
		c.Reactions = make([]*ChatReaction, 0)
	}
	c.Reactions = append(c.Reactions, reaction)
	c.ReactionsCount++
}

func (c *Chat) RemoveReaction(userID uuid.UUID, reactionType string) {
	for i, reaction := range c.Reactions {
		if reaction.UserID == userID && reaction.ReactionType == reactionType {
			c.Reactions = append(c.Reactions[:i], c.Reactions[i+1:]...)
			c.ReactionsCount--
			return
		}
	}
}

func (c *Chat) GetReactionByUser(userID uuid.UUID) *ChatReaction {
	for _, reaction := range c.Reactions {
		if reaction.UserID == userID {
			return reaction
		}
	}
	return nil
}

func (c *Chat) GetReactionCountByType() map[string]int64 {
	counts := make(map[string]int64)
	for _, reaction := range c.Reactions {
		counts[reaction.ReactionType]++
	}
	return counts
}

func (c *Chat) ToResponse() *ChatResponse {
	response := &ChatResponse{
		Chat:      c,
		Reactions: c.Reactions,
	}

	if c.ReplyTo != nil {
		response.ReplyTo = c.ReplyTo.ToResponse()
	}

	if c.ForwardedFrom != nil {
		response.ForwardedFrom = c.ForwardedFrom.ToResponse()
	}

	return response
}

func (c *Chat) IsMediaMessage() bool {
	return c.Type != "text" && c.FileURL != nil
}

func (c *Chat) IsSystemMessage() bool {
	return c.Type == "system"
}
