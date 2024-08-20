/*
Package services - NekoBlog backend server services.
This file is for comment related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package services

import (
	"errors"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

// CommentService 评论服务
type CommentService struct {
	commentStore *stores.CommentStore
}

// NewCommentService 返回一个新的评论服务实例。
//
// 返回：
//   - *CommentService: 返回一个指向新的评论服务实例的指针。
func (factory *Factory) NewCommentService() *CommentService {
	return &CommentService{
		commentStore: factory.storeFactory.NewCommentStore(),
	}
}

// NewCommentService 创建评论
//
// 参数：
//   - uid：用户ID
//   - postID: 博文编号
//   - content: 博文内容
//   - postStore，userStore：绑定post和user层来调用方法
//
// 返回值：
//
//	-error 创建失败返回创建失败时候的具体信息
func (service *CommentService) CreateComment(uid uint64, postID uint64, content string, postStore *stores.PostStore, userStore *stores.UserStore) error {
	// 校验评论是否存在
	existance, err := postStore.ValidatePostExistence(postID)
	if err != nil {
		return err
	}
	if !existance {
		return errors.New("post does not exist")
	}

	// 根据 UID 获取 Username
	user, err := userStore.GetUserByUID(uid)
	if err != nil {
		return err
	}

	// 调用存储层的方法存储评论
	err = service.commentStore.CreateComment(uid, user.UserName, postID, content)
	if err != nil {
		return err
	}
	return nil
}

// UpdateComment 修改评论
//
// 参数：
//   - comment：评论ID
//   - content: 博文内容
//
// 返回值：
//
//	-error 如果评论存在返回修改评论时候的信息
func (service *CommentService) UpdateComment(commentID uint64, content string) error {
	//检查评论是否存在
	// TODO: 改为控制器层判断
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}

	// 调用数据库或其他存储方法更新评论内容
	err = service.commentStore.UpdateComment(commentID, content)
	if err != nil {
		return err
	}

	// 如果更新成功，返回nil
	return nil
}

// DeleteComment 删除评论
//
// 参数：
//   - commentID: 评论ID
//
// 返回值：
//   - error 返回处理删除的信息
func (service *CommentService) DeleteComment(commentID uint64) error {
	// TODO: 改为控制器层判断
	exists, err := service.commentStore.ValidateCommentExistence(commentID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("comment does not exist")
	}
	// 调用评论存储中的删除评论方法
	err = service.commentStore.DeleteComment(commentID)
	if err != nil {
		// 如果发生错误，则返回错误
		return err
	}

	// 如果没有发生错误，则返回 nil
	return nil
}

// GetCommentList 获取评论列表
//
// 返回值：
//   - 成功则返回评论列表
//   - 失败返回nil
func (service *CommentService) GetCommentList() ([]models.CommentInfo, error) {
	return service.commentStore.GetCommentList()
}

// GetCommentInfo 获取评论信息
//
// 返回值：
//   - 成功返回评论体
//   - 失败返回nil
func (service *CommentService) GetCommentInfo(commentID uint64) (models.CommentInfo, error) {
	return service.commentStore.GetCommentInfo(commentID)
}
