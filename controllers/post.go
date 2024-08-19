/*
Package controllers - NekoBlog backend server controllers.
This file is for post controller, which is used to create handlee post related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
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

// NewListHandler 博文列表函数
//
// 返回值：
// - fiber.Handle：新的博文列表函数
func (controller *PostController) NewListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		posts, err := controller.postService.GetPostList()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(posts)
	}
}
