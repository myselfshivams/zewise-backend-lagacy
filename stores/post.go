/*
Package stores provides the implementation of stores
responsible for handling posts in the micro-blog backend.
Import the necessary package for the store implementation.
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package stores

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"gorm.io/gorm"
)

// PostStore 是一个结构体，代表微博后端中负责处理帖子的存储。
// 它包含一个指向 gorm.DB 的引用，允许与与帖子相关的数据库进行交互。
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

func (store *PostStore) PostFindStore(posts *[]models.PostInfo) error {
	if result := store.db.Find(posts); result.Error != nil {
		return result.Error
	}
	return nil
}
