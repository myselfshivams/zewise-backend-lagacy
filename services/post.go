/*
Package services - NekoBlog backend server services.
This file is for user related services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/

package services

import (
	"mime/multipart"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/converters"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/validers"
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
func (service *PostService) GetPostList() ([]models.PostInfo, error) {
	userPosts, err := service.postStore.GetPostList()
	if err != nil {
		return nil, err
	}
	return userPosts, nil
}

// GetPostInfoByUsername 根据用户名获取用户信息。
//
// 参数：
//   - UID：用户ID
//
// 返回值：
//   - *models.postInfo：用户信息模型。
func (service *PostService) GetPostInfo(postID uint64) (models.PostInfo, error) {
	post, err := service.postStore.GetPostInfo(postID)
	if err != nil {
		return models.PostInfo{}, err
	}
	return post, nil
}

// CreatePost 根据用户提交的帖子信息创建帖子。
//
// 参数：
//   - userID：用户ID，用于关联帖子与用户。
//   - ipAddr：IP地址
//   - postReqInfo：帖子信息，包含标题、内容等。
//   - postImages:帖子图片
//
// 返回值：
//   - error：如果在创建过程中发生错误，则返回相应的错误信息，否则返回nil。
func (service *PostService) CreatePost(uid uint64, ipAddr string, postReqInfo types.PostCreateBody, postImages []*multipart.FileHeader) (models.PostInfo, error) {
	var (
		converteredImages [][]byte
		imageFile         multipart.File
		err               error
	)
	// 校验并处理图片
	for _, image := range postImages {
		imageFile, err = image.Open()
		if err != nil {
			imageFile.Close()
			return models.PostInfo{}, err
		}

		// 验证图像文件的有效性，包括尺寸和文件类型等
		fileType, err := validers.ValidImageFile(
			image,
			&imageFile,
			consts.POST_IMAGE_MIN_WIDTH,
			consts.POST_IMAGE_MIN_HEIGHT,
			consts.MAX_AVATAR_FILE_SIZE,
		)
		if err != nil {
			imageFile.Close()
			return models.PostInfo{}, err
		}

		//调整帖子图片的大小
		converteredImage, err := converters.ResizePostImage(fileType, &imageFile)
		if err != nil {
			imageFile.Close()
			return models.PostInfo{}, err
		}
		converteredImages = append(converteredImages, converteredImage)
		imageFile.Close()
	}

	// 调用存储层的方法创建帖子
	postInfo, err := service.postStore.CreatePost(uid, ipAddr, postReqInfo, converteredImages)
	if err != nil {
		return models.PostInfo{}, err
	}
	return postInfo, nil
}

// DeletePost 是用于删除博文的服务方法
//
// 参数：
// - postID uint64：待删除博文的ID
//
// 返回值：
// - error：如果发生错误，返回相应错误信息；否则返回 nil
func (service *PostService) DeletePost(postID uint64) error {
	// 调用post存储中的删除post方法
	return service.postStore.DeletePost(postID)
}
