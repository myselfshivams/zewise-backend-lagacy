package stores

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"gorm.io/gorm"
)

// Comment 评论信息数据库
type CommentStore struct {
	db *gorm.DB
}

// NewCommentStore 返回一个新的用户存储实例。
// 返回：
//   - *CommentStore: 返回一个指向新的用户存储实例的指针。
func (factory *Factory) NewCommentStore() *CommentStore {
	return &CommentStore{factory.db}
}

// NewCommentStore 存储comment
// 参数 ：comment实例
// 返回：
//
//	-error
func (store *CommentStore) NewCommentStore(uid uint64, comment *models.CommentInfo) error {

	newComment := &models.CommentInfo{
		Username: comment.Username,
		Content:  comment.Content,
		UID:      uid,
	}

	result := store.db.Create(newComment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
