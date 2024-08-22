/*
Package controllers - NekoBlog backend server controllers.
This file is for user controller, which is used to create handlee user related requests.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 返回一个新的 UserController 实例。
//
// 返回值：
//   - *UserController：新的 UserController 实例。
func (factory *Factory) NewUserController() *UserController {
	return &UserController{
		userService: factory.serviceFactory.NewUserService(),
	}
}

// NewProfileHandler 返回获取用户资料的处理函数。
//
// 返回值：
//   - fiber.Handler：新的获取用户资料的处理函数。
func (controller *UserController) NewProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 判断传入的查询参数是UID还是Username
		uidString := ctx.Query("uid")
		username := ctx.Query("username")
		if uidString == "" && username == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "parameter uid or username is required"),
			)
		}

		// 获取用户信息
		var (
			user *models.UserInfo
			err  error
		)
		switch uidString {
		// 根据用户名获取用户信息
		case "":
			user, err = controller.userService.GetUserInfoByUsername(username)

		// 根据UID获取用户信息
		default:
			var uid uint64
			// 将UID String转换为uint64
			if uid, err = strconv.ParseUint(uidString, 10, 64); err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
				)
			}
			user, err = controller.userService.GetUserInfoByUID(uid)
		}

		// 若用户不存在则返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "user does not exist"),
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
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewUserProfileData(user)),
		)
	}
}

// NewRegisterHandler 返回注册用户的处理函数。
//
// 返回值：
//   - fiber.Handler：新的注册用户的处理函数。
func (controller *UserController) NewRegisterHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.Username == "" || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username or password is required"),
			)
		}

		// 注册用户
		err = controller.userService.RegisterUser(reqBody.Username, reqBody.Password)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewLoginHandler 返回登陆用户的处理函数。
//
// 返回值：
//   - fiber.Handler：新的登陆用户的处理函数。
func (controller *UserController) NewLoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.Username == "" || reqBody.Password == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username or password is required"),
			)
		}

		// 解析 UA
		userAgentString := ctx.Get("User-Agent")
		ua := useragent.New(userAgentString)
		// 获取浏览器信息
		browser, version := ua.Browser()
		var sb strings.Builder
		sb.WriteString(browser)
		sb.WriteString(" ")
		sb.WriteString(version)
		browserInfo := sb.String()
		// 获取操作系统信息
		os := ua.OSInfo().FullName

		// 登陆
		token, err := controller.userService.LoginUser(reqBody.Username, reqBody.Password, ctx.IP(), browserInfo, os)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewUserToken(token)),
		)
	}
}

// NewUploadAvatarHandler 返回上传头像的处理函数。
//
// 返回值：
//   - fiber.Handler：新的上传头像的处理函数。
func (controller *UserController) NewUploadAvatarHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取文件
		form, err := ctx.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["avatar"]
		// 检查表单中文件的数量
		if len(files) != 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "required 1 file, but got more or less"),
			)
		}
		fileHeader := files[0]

		// 保存头像
		err = controller.userService.UserUploadAvatar(claims.UID, fileHeader)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回成功的 JSON 响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

//	NewUserUpdatePassword 修改密码的函数
//
// 返回值：
//   - fiber.Handler：新密码的处理函数。
func (controller *UserController) NewUpdatePasswordHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserUpdatePasswordBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.Username == "" || reqBody.Password == "" || reqBody.NewPassword == ""{
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "username, password or new password is required"),
			)
		}

		// 修改密码
		err = controller.userService.UserUpdatePassword(
			reqBody.Username,
			reqBody.Password,
			reqBody.NewPassword,
		)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		// 返回成功的 JSON 响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "password updated successfully"),
		)
	}
}

// NewUserUpdateProfileHandler 返回更新用户资料的处理函数。
//
// 返回值：
//   - fiber.Handler：新的更新用户资料的处理函数。
func (controller *UserController) NewUpdateProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserUpdateProfileBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 校验参数
		if reqBody.NickName == nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "nickname is required"),
			)
		}

		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 更新用户资料
		err = controller.userService.UpdateUserInfo(claims.UID, reqBody)
		if err != nil {
			return ctx.Status(500).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, "failed to update profile"),
			)
		}

		// 返回成功的 JSON 响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "profile updated successfully"),
		)
	}
}
