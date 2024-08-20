package serializers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

type CommentListResponse struct {
	IDs []uint64 `json:"ids"`
}

// NewCommentListResponse 创建评论列表的响应
//
// 参数：
//   - 评论体列表
//
// 返回值：
//   - 评论列表的响应
func NewCommentListResponse(commentInfos []models.CommentInfo) CommentListResponse {
	var ids []uint64
	for _, commentInfos := range commentInfos {
		ids = append(ids, uint64(commentInfos.ID))
	}
	return CommentListResponse{IDs: ids}
}

// CommentDetailResponse 文章信息响应结构
// TODO: 实现 Replies 字段
type CommentDetailResponse struct {
	CommentID     uint64 `json:"comment_id"`     // 评论ID
	PostID        uint64 `json:"post_id"`        // 博文ID
	PosterUID     uint64 `json:"poster_uid"`     // 发布者UID
	PostTimestamp int64  `json:"post_timestamp"` // 博文发布时间戳
	Content       string `json:"content"`        // 内容
	Likes         int    `json:"likes"`          // 点赞数
	Replies       int    `json:"replies"`        // 回复数
	Is_liked      bool   `json:"is_liked"`       // 是否点赞
	Is_disliked   bool   `json:"is_disliked"`    // 是否点踩
}

// NewCommentDetailResponse 创建评论实例
func NewCommentDetailResponse(comment models.CommentInfo) *CommentDetailResponse {
	// 创建一个新的 CommentProfileData 实例
	// TODO: 实现 Is_liked / Is_disliked
	profileData := &CommentDetailResponse{
		CommentID:     uint64(comment.ID),
		PostID:        comment.PostID,
		PosterUID:     comment.UID,
		PostTimestamp: comment.CreatedAt.Unix(),
		Content:       comment.Content,
		Likes:         len(comment.Like),
	}

	return profileData
}

// CreateCommentResponse 用于将 CommentInfo 转换为 JSON 格式的结构体
type CreateCommentResponse struct {
	ID uint64 `json:"id"`
}

// NewCommentResponse 用于创建 CommentResponse 实例
func NewCreateCommentResponse(commentInfo models.CommentInfo) CreateCommentResponse {
	var resp = CreateCommentResponse{
		ID: uint64(commentInfo.ID),
	}
	return resp
}
