/*
Package controllers - NekoBlog backend server controllers.
This file is for post controller, which is used to create handlee post related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
- CBofJOU<2023122312@jou.edu.cn>
*/
package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
)

// PostController 博文控制器结构体
type PostController struct {
	postService *services.PostService
}

// NewPostController 博文控制器工厂函数。
//
// 返回值：
// - *PostController 博文控制器指针
func (factory *Factory) NewPostController() *PostController {
	return &PostController{
		postService: factory.serviceFactory.NewPostService(),
	}
}

// NewPostListHandler 博文列表函数
//
// 返回值：
// - fiber.Handle：新的博文列表函数
func (controller *PostController) NewPostListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		posts, err := controller.postService.GetPostList()
		if err != nil {
			return c.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		return c.Status(200).JSON(
			serializers.NewPostListResponse(posts),
		)
	}
}

// NewDetailHandler 获取文章信息的函数
//
// 返回值：
//   - fiber.Handler：新的获取文章信息的函数
func (controller *PostController) NewPostDetailHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 根据 UID 获取用户信息
		postIDString := ctx.Query("post-id")
		if postIDString == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post id is required"),
			)
		}

		//获取帖子的唯一标识符
		postID, err := strconv.ParseUint(postIDString, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取帖子的详细信息
		post, err := controller.postService.GetPostInfo(postID)
		// 若post不存在则返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post does not exist"),
			)
		}

		// 返回其他错误
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewPostDetailResponse(post)),
		)
	}
}

// NewDetailHandler 博文详情函数。
//
// 返回值：
// - fiber.Handler：新的博文详情函数
func (controller *PostController) NewCreatePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 提取令牌声明
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 解析用户请求
		reqBody := types.PostCreateBody{}
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 验证参数
		if reqBody.Title == "" || reqBody.Content == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "post title or post content"),
			)
		}

		// 处理图片上传
		form, err := ctx.MultipartForm()
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}
		files := form.File["images"]

		// 检查表单中文件的数量
		if len(files) > 9 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "The number of images cannot exceed 9"),
			)
		}

		// 创建博文
		postInfo, err := controller.postService.CreatePost(claims.UID, ctx.IP(), reqBody, files)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回成功响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(
				consts.SUCCESS,
				"post created successfully",
				serializers.NewCreatePostResponse(postInfo),
			),
		)
	}
}

// NewDeletePostHandler 返回一个用于处理删除博文请求的 Fiber 处理函数
//
// 参数：
// - controller *PostController：博文控制器实例
//
// 返回值：
// - fiber.Handler：新的博文删除函数
func (controller *PostController) NewDeletePostHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取PostID
		postID := ctx.Params("post")
		fmt.Println(postID)

		//验证postID是否为空
		if postID == "" {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id cannot be empty"))
		}

		// 将post ID转换为无符号整数
		postIDUint, err := strconv.ParseUint(postID, 10, 64)
		if err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.PARAMETER_ERROR, "post id must be a number"))
		}

		// 执行删除操作
		if err := controller.postService.DeletePost(postIDUint); err != nil {
			return ctx.Status(200).JSON(serializers.NewResponse(consts.SERVER_ERROR, err.Error()))
		}

		return ctx.JSON(serializers.NewResponse(consts.SUCCESS, "succeed"))
	}
}
