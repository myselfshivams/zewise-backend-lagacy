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

// CommentCreatebody 创建评论请求体
type UserCommentCreateBody struct {
	PostID  *uint64 `json:"post_id" form:"post_id"` // 博文ID
	Content string  `json:"content" form:"content"` // 内容
}

// UserCommentUpdateBody 更新评论请求体
type UserCommentUpdateBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"` // 评论ID
	Content   string  `json:"content" form:"content"`       // 内容
}

// UserPostInfo 创建博文请求体
type UserPostInfo struct {
	UID   uint   `json:"id"`    // 用户ID
	Title string `json:"title"` // 标题

}

// PostCreateBody 创建博文请求体
type PostCreateBody struct {
	Title   string `json:"title" form:"title"`     //标题
	Content string `json:"content" form:"content"` //内容
}

// UserPostInfo 创建博文请求体
type UserCommentDeleteBody struct {
	CommentID *uint64 `json:"comment_id" form:"comment_id"` // 评论ID
}
