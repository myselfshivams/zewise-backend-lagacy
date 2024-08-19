/*
Package controllers - NekoBlog backend server controllers.
This file is for comment controller, which is used to create handlee comment related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/gofiber/fiber/v2"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *services.CommentService
}

// NewCommentController 创建一个新的评论控制器实例。
//
// 返回：
//   - *CommentController: 返回一个新的评论控制器实例。
func (factory *Factory) NewCommentController() *CommentController {
	return &CommentController{
		commentService: factory.serviceFactory.NewCommentService(),
	}
}

// NewCreateCommentHandler 返回一个处理创建评论请求的Handler函数。
//
// 返回：
//   - 返回处理请求
func (controller *CommentController) NewCreateCommentHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.CommentCreateBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.Content == "" || reqBody.Username == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "content or user_id is required"),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)
		if claims == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not avaliable"),
			)
		}

		// 调用服务方法创建评论
		err = controller.commentService.NewCommentService(claims.UID, reqBody.Content, reqBody.Username)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.	SERVER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "Comment created successfully"),
		)
	}
}
