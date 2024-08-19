/*
Package controllers provides the implementation of controllers
responsible for handling and managing posts in the micro-blog backend.
Import the necessary package for the controllers.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/gofiber/fiber/v2"
)

// PostController 是一个结构体，代表微博后端中负责管理和处理帖子的控制器。
// 它包含一个指向 PostService 的引用，允许与与帖子相关的底层服务进行交互。
type PostController struct {
	postService *services.PostService
}

// NewPostController 是一个工厂方法，用于创建 PostController 的新实例。
//
// 返回值
// 它初始化并返回一个 PostController，并关联了相应的 PostService。
func (factory *Factory) NewPostController() *PostController {
	return &PostController{
		postService: factory.serviceFactory.NewPostService(),
	}
}

func (controller *PostController) NewPostListHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Logic to retrieve the list of posts
		posts, err := controller.postService.GetPosts()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(posts)
	}
}
