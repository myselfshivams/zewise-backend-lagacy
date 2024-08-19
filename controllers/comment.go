package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
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
		reqBody := new(types.CommentCreatBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)
		if claims == nil {
            return ctx.Status(200).JSON(
                serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not avaliable"),
            )
        }

		// 创建评论
		comment := &models.CommentInfo{
			Content:  reqBody.Content,
			Username: reqBody.Username,
		}

		// 验证参数合法性
		if comment.Username == nil || comment.Content == nil {
			// 如果Username或Content为空，返回参数错误
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "Username or Content is missing"),
			)
		}

		// 调用服务方法创建评论
		err = controller.commentService.NewCommentService(claims.UID,comment)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 成功时返回响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "Comment created successfully"),
		)
	}
}
