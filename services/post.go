/*
Package services - NekoBlog backend server services.
This file is for user related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/

package services

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// PostService 博文服务
type PostService struct {
	postStore *stores.PostStore
}

// PostService 返回一个新的 PostService 实例
//
// 返回值：
//   - *PostService：新的 PostService 实力。
func (factory *Factory) NewPostService() *PostService {
	return &PostService{
		postStore: factory.storeFactory.NewPostStore(),
	}
}

// GetPostList 获取适用于用户查看的帖子信息列表。
// 返回值：
// - []models.UserPostInfo: 包含适用于用户查看的帖子信息的切片。
// - error: 在获取帖子信息过程中遇到的任何错误，如果有的话。
func (service *PostService) GetPostList() ([]types.UserPostInfo, error) {
	var userPosts []types.UserPostInfo
	var err error
	if userPosts, err = service.postStore.GetPostList(); err != nil {
		return nil, err
	}
	return userPosts, nil
}
