package services

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)


// CommentService 评论服务
type CommentService struct {
	commentStore *stores.CommentStore
}


// NewCommentService 返回一个新的评论服务实例。
// 返回：
//   - *CommentService: 返回一个指向新的评论服务实例的指针。
func (factory *Factory) NewCommentService() *CommentService {
	return &CommentService{
		commentStore: factory.storeFactory.NewCommentStore(),
	}
}

//NewCommentService 创建评论
//
// 参数： commment实例
// 返回值：
//      -error

func (service *CommentService) NewCommentService(uid uint64,comment *models.CommentInfo) error {

    // 调用存储层的方法存储评论
    err := service.commentStore.NewCommentStore(uid,comment)
    if err != nil {
        // 处理错误
        return err
    }
    return nil
}


