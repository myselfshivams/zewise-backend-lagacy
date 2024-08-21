/*
Package serializers - NekoBlog backend server data serialization.
This file is for user data serialization.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package serializers

import (
	"strings"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

// UserProfileData 用户资料响应结构。
type UserProfileData struct {
	UID      uint64  `json:"uid"`        // 用户 ID
	Username string  `json:"username"`   // 用户名
	Nickname string  `json:"nickname"`   // 昵称
	Avatar   string  `json:"avatar_url"` // 头像 URL
	Birth    *int64  `json:"birth"`      // 生日
	Gender   *string `json:"gender"`     // 性别
	Level    uint64  `json:"level"`      // 等级
}

// NewUserProfileData 创建一个新的用户资料响应。
//
// 参数：
//   - model：用户资料模型
//
// 返回值：
//   - *UserProfileData：新的用户资料响应结构体。
func NewUserProfileData(model *models.UserInfo) *UserProfileData {
	// 创建用户资料响应
	profile := new(UserProfileData)
	profile.UID = uint64(model.ID)
	profile.Username = model.UserName

	// 设置昵称
	if model.NickName != nil {
		profile.Nickname = *model.NickName
	} else {
		// 如果没有昵称，则使用用户名
		profile.Nickname = model.UserName
	}

	// 设置头像
	var sb strings.Builder
	sb.WriteString("/resources/avatar/")
	sb.WriteString(model.Avatar)
	profile.Avatar = sb.String()

	// 设置生日和性别
	// 如果生日和性别为空，则设置为未知
	if model.Birth != nil {
		birth := model.Birth.Unix()
		profile.Birth = &birth
	} else {
		profile.Birth = nil
	}
	if model.Gender != nil {
		profile.Gender = model.Gender
	} else {
		profile.Gender = nil
	}
	profile.Level = model.Level

	// 返回用户资料响应
	return profile
}

// UserToken
type UserToken struct {
	Token string `json:"token"`
}

// NewUserToken 创建一个新的用户 Token 响应。
//
// 参数：
//   - token：用户 Token
//
// 返回值：
//   - *UserToken：新的用户 Token 响应结构体。
func NewUserToken(token string) *UserToken {
	return &UserToken{Token: token}
}
