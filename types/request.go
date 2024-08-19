/*
Package type - NekoBlog backend server types.
This file is for user related types.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

// UserAuthBody 认证请求体
type UserAuthBody struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// UserRegisterBody 注册请求体
type UserUpdatePasswordBody struct {
	UserAuthBody        // 认证请求体
	NewPassword  string `json:"new_password"` // 新密码
}

// UserUpdateProfileBody 更新用户资料请求体
type UserUpdateProfileBody struct {
	NickName *string `json:"nickname"` // 昵称
	Birth    *uint64 `json:"birth"`    // 出生日期
	Gender   *string `json:"gender"`   // 性别
}

// Comm

type CommentCreatBody struct {
	Username *string `json:"username"` // 用户名
	Content  *string `json:"content"`  // 内容
}
