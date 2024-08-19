/*
Package stores - NekoBlog backend server data access objects.
This file is for post storage accessing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"gorm.io/gorm"
)

// PostStore 博文信息数据库
type PostStore struct {
	db *gorm.DB
}

// NewPostStore 是一个工厂方法，用于创建 PostStore 的新实例。
//
// 参数
// - factory: 一个包含 gorm.DB 的 Factory 实例，用于初始化 PostStore 的数据库连接。
//
// 返回值
// 它初始化并返回一个 PostStore，并关联了相应的 gorm.DB。
func (factory *Factory) NewPostStore() *PostStore {
	return &PostStore{factory.db}
}

// GetPostLis 获取适用于用户查看的帖子信息列表。
//
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在检索过程中遇到的任何错误，如果有的话。
func (store *PostStore) GetPostList() ([]types.UserPostInfo, error) {
	var posts []models.PostInfo
	if result := store.db.Find(&posts); result.Error != nil {
		return nil, result.Error
	}
	userPosts := make([]types.UserPostInfo, len(posts))
	for i, post := range posts {
		userPosts[i] = types.UserPostInfo{
			UID:   post.ID,
			Title: post.Title,
		}
	}
	return userPosts, nil
}
