/*
Package stores - NekoBlog backend server data access objects.
This file is for comment storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"errors"

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
//
// 参数 ：- uid：用户id，- username: 用户名，- postID: 博文id，- content: 博文内容
//
// 返回：
//
//	-error 正确返回nil
func (store *CommentStore) CreateComment(uid uint64, username string, postID uint64, content string) error {
	newComment := &models.CommentInfo{
		PostID:   postID,
		Username: username,
		Content:  content,
		UID:      uid,
		Like:     nil,
		Dislike:  nil,
		IsPublic: true,
	}

	result := store.db.Create(newComment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//	ValidateCommentExistence 判断评论是否存在
//
//	参数：
//	- commentID: 评论ID
//
// 返回值：
//   - error：如果评论存在返回true，不存在判断具体的错误类型返回false
func (store *CommentStore) ValidateCommentExistence(commentID uint64) (bool, error) {
	var comment models.CommentInfo
	result := store.db.Where("id = ?", commentID).First(&comment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	// 返回错误类型
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// UpdateComment 修改评论
//
//	参数：
//	- commentID: 评论ID
//	- content: 修改内容
//
// 返回值：
//   - error：如果评论存在返回true，不存在判断具体的错误类型返回false
func (store *CommentStore) UpdateComment(commentID uint64, content string) error {
	commentInfo := new(models.CommentInfo)
	result := store.db.Where("id = ?", commentID).First(commentInfo)
	if result.Error != nil {
		return result.Error
	}

	commentInfo.Content = content
	result = store.db.Save(commentInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteComment 删除评论
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - error：返回删除处理的成功与否
func (store *CommentStore) DeleteComment(commentID uint64) error {
	return store.db.Where("id = ?", commentID).Unscoped().Delete(&models.CommentInfo{}).Error
}

// GetCommentList 获取评论列表
//
// 返回值：
//   - 成功则返回评论列表
//   - 失败返回nil
func (store *CommentStore) GetCommentList() ([]models.CommentInfo, error) {
	var userComments []models.CommentInfo
	result := store.db.Find(&userComments)
	if result.Error != nil {
		return nil, result.Error
	}
	return userComments, nil
}

// GetCommentInfo 获取评论信息
//
// 参数：
//   - commentID：评论ID
//
// 返回值：
//   - models.CommentInfo：成功返回评论信息
//   - error：失败返回error
func (store *CommentStore) GetCommentInfo(commentID uint64) (models.CommentInfo, error) {
	comment := models.CommentInfo{}
	result := store.db.Where("id = ?", commentID).First(&comment)
	return comment, result.Error
}
