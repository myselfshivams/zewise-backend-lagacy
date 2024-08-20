package serializers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

type PostListResponse struct {
	IDs []uint64 `json:"ids"`
}

func NewPostListResponse(postInfos []models.PostInfo) *PostListResponse {
	ids := make([]uint64, 0)
	for _, postInfo := range postInfos {
		ids = append(ids, uint64(postInfo.ID))
	}
	return &PostListResponse{IDs: ids}
}

// PostDetailResponse 文章信息响应结构
type PostDetailResponse struct {
	CommentID    uint64   `json:"comment_id"`   //
	UID          uint64   `json:"uid"`          // 用户ID
	Title        string   `json:"title"`        // 标题
	Content      string   `json:"content"`      // 内容
	ParentPostID *uint64  `json:"ParentPostID"` // 转发自文章ID
	Images       []string `json:"images"`       // 图片
	Like         int      `json:"like"`         // 点赞数
	Favorite     int      `json:"favorite"`     // 收藏数
	Farward      int      `json:"farward"`      // 转发数
}

// NewPostDetailResponse 创建新的文章信息响应
//
// 参数：
//   - model：文章信息模型
//
// 返回值：
//   - *PostProfileData：新的文章信息响应结构
func NewPostDetailResponse(post models.PostInfo) *PostDetailResponse {
	// 创建一个新的 PostProfileData 实例
	profileData := &PostDetailResponse{
		CommentID:    uint64(post.ID),
		UID:          post.UID,
		Title:        post.Title,
		Content:      post.Content,
		ParentPostID: post.ParentPostID,
		Images:       post.Images,
		Like:         len(post.Like),
		Favorite:     len(post.Farward),
		Farward:      len(post.Farward),
	}

	return profileData
}

// CreatePostResponse 用于将 PostInfo 转换为 JSON 格式的结构体
type CreatePostResponse struct {
	ID uint64 `json:"id"`
}

// NewPostResponse 用于创建 PostResponse 实例
func NewCreatePostResponse(postInfo models.PostInfo) CreatePostResponse {
	var resp = CreatePostResponse{
		ID: uint64(postInfo.ID),
	}
	return resp
}
