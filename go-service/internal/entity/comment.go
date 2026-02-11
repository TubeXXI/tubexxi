package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	UserID       uuid.UUID      `json:"user_id" db:"user_id"`
	PageURL      string         `json:"page_url" db:"page_url"`
	Name         *string        `json:"name,omitempty" db:"name"`
	Email        *string        `json:"email,omitempty" db:"email"`
	Comment      string         `json:"comment" db:"comment"`
	Type         string         `json:"type" db:"type"`
	ReplyToID    *uuid.UUID     `json:"reply_to_id,omitempty" db:"reply_to_id"`
	IsEdited     bool           `json:"is_edited" db:"is_edited"`
	IsDeleted    bool           `json:"is_deleted" db:"is_deleted"`
	MediaURL     *string        `json:"media_url,omitempty" db:"media_url"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
	Parent       *Comment       `json:"parent,omitempty" db:"-"`
	Replies      []*Comment     `json:"replies,omitempty" db:"-"`
	User         *User          `json:"user,omitempty" db:"-"`
	Likes        []*LikeComment `json:"likes,omitempty" db:"-"`
	LikesCount   int64          `json:"likes_count" db:"likes_count"`
	RepliesCount int64          `json:"replies_count" db:"replies_count"`
	IsLikedByMe  bool           `json:"is_liked_by_me" db:"-"`
	Depth        int            `json:"depth,omitempty" db:"depth"`
	Path         string         `json:"path,omitempty" db:"path"`
}

type LikeComment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CommentID uuid.UUID `json:"comment_id" db:"comment_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	User      *User     `json:"user,omitempty" db:"-"`
	Comment   *Comment  `json:"comment,omitempty" db:"-"`
}

type CommentCreateRequest struct {
	PageURL   string     `json:"page_url" validate:"required"`
	Comment   string     `json:"comment" validate:"required"`
	Type      string     `json:"type" validate:"omitempty,oneof=text image video link document other"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
	MediaURL  *string    `json:"media_url,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Email     *string    `json:"email,omitempty" validate:"omitempty,email"`
}
type CommentUpdateRequest struct {
	Comment  *string `json:"comment,omitempty"`
	MediaURL *string `json:"media_url,omitempty"`
}
type CommentResponse struct {
	*Comment
	Replies []*CommentResponse `json:"replies,omitempty"`
	Parent  *CommentResponse   `json:"parent,omitempty"`
}
type CommentFilter struct {
	PageURL        string     `json:"page_url"`
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	ParentID       *uuid.UUID `json:"parent_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	SortBy         string     `json:"sort_by" validate:"omitempty,oneof=created_at likes_count replies_count"`
	SortOrder      string     `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Limit          int        `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset         int        `json:"offset" validate:"omitempty,min=0"`
	IncludeReplies bool       `json:"include_replies"`
	Depth          int        `json:"depth" validate:"omitempty,min=1,max=10"`
}
type CommentTree struct {
	Comment  *Comment       `json:"comment"`
	Children []*CommentTree `json:"children"`
	Level    int            `json:"level"`
}
type CommentStats struct {
	TotalComments  int64            `json:"total_comments"`
	TotalReplies   int64            `json:"total_replies"`
	TotalLikes     int64            `json:"total_likes"`
	UniqueUsers    int64            `json:"unique_users"`
	CommentsByType map[string]int64 `json:"comments_by_type"`
	CommentsByHour map[int]int64    `json:"comments_by_hour"`
}
type CommentWithUserInfo struct {
	Comment
	UserName      string  `json:"user_name" db:"user_name"`
	UserAvatar    *string `json:"user_avatar,omitempty" db:"user_avatar"`
	UserEmail     *string `json:"user_email,omitempty" db:"user_email"`
	ParentAuthor  *string `json:"parent_author,omitempty" db:"parent_author"`
	ParentContent *string `json:"parent_content,omitempty" db:"parent_content"`
}

func (c *Comment) IsRoot() bool {
	return c.ReplyToID == nil
}

func (c *Comment) HasReplies() bool {
	return c.RepliesCount > 0 && c.Replies != nil && len(c.Replies) > 0
}
func (c *Comment) AddReply(reply *Comment) {
	if c.Replies == nil {
		c.Replies = make([]*Comment, 0)
	}
	c.Replies = append(c.Replies, reply)
	c.RepliesCount++
}
func (c *Comment) RemoveReply(replyID uuid.UUID) {
	for i, reply := range c.Replies {
		if reply.ID == replyID {
			c.Replies = append(c.Replies[:i], c.Replies[i+1:]...)
			c.RepliesCount--
			return
		}
	}
}
func (c *Comment) AddLike(like *LikeComment) {
	if c.Likes == nil {
		c.Likes = make([]*LikeComment, 0)
	}
	c.Likes = append(c.Likes, like)
	c.LikesCount++
}
func (c *Comment) RemoveLike(userID uuid.UUID) {
	for i, like := range c.Likes {
		if like.UserID == userID {
			c.Likes = append(c.Likes[:i], c.Likes[i+1:]...)
			c.LikesCount--
			return
		}
	}
}
func (c *Comment) IsLikedBy(userID uuid.UUID) bool {
	for _, like := range c.Likes {
		if like.UserID == userID {
			return true
		}
	}
	return false
}
func (c *Comment) ToResponse() *CommentResponse {
	response := &CommentResponse{
		Comment: c,
	}

	if c.HasReplies() {
		response.Replies = make([]*CommentResponse, len(c.Replies))
		for i, reply := range c.Replies {
			response.Replies[i] = reply.ToResponse()
		}
	}

	if c.Parent != nil {
		response.Parent = c.Parent.ToResponse()
	}

	return response
}

func BuildTree(comments []*Comment) []*CommentTree {
	// Create map of comments by ID
	commentMap := make(map[uuid.UUID]*CommentTree)
	roots := make([]*CommentTree, 0)

	// Initialize trees
	for _, comment := range comments {
		commentMap[comment.ID] = &CommentTree{
			Comment:  comment,
			Children: make([]*CommentTree, 0),
			Level:    0,
		}
	}

	// Build tree structure
	for _, comment := range comments {
		tree := commentMap[comment.ID]

		if comment.ReplyToID != nil {
			if parent, exists := commentMap[*comment.ReplyToID]; exists {
				tree.Level = parent.Level + 1
				parent.Children = append(parent.Children, tree)
			} else {
				roots = append(roots, tree)
			}
		} else {
			roots = append(roots, tree)
		}
	}

	return roots
}

func FlattenTree(trees []*CommentTree, maxDepth int) []*Comment {
	comments := make([]*Comment, 0)

	var flatten func(tree *CommentTree, depth int)
	flatten = func(tree *CommentTree, depth int) {
		if depth > maxDepth {
			return
		}

		tree.Comment.Depth = depth
		comments = append(comments, tree.Comment)

		for _, child := range tree.Children {
			flatten(child, depth+1)
		}
	}

	for _, tree := range trees {
		flatten(tree, 0)
	}

	return comments
}
