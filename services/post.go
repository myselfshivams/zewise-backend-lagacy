/*
Package services provides the implementation of services
responsible for managing posts in the micro-blog backend.
Import the necessary package for the services.
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/

package services

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

// PostService 是一个结构体，代表微博后端中负责管理帖子的服务。
// 它包含一个指向 PostStore 的引用，允许与与帖子相关的底层存储进行交互。
type PostService struct {
	postStore *stores.PostStore
}

// NewPostService 是一个工厂方法，用于创建 PostService 的新实例。
//
// 返回值
// 它初始化并返回一个 PostService，并关联了相应的 PostStore。
func (factory *Factory) NewPostService() *PostService {
	return &PostService{
		postStore: factory.storeFactory.NewPostStore(),
	}
}

func (service *PostService) GetPosts() ([]models.PostInfo, error) {
	posts := []models.PostInfo{}
	if err := service.postStore.PostFindStore(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}
