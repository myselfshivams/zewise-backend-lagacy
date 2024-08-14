package serializers

import (
	"fmt"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
)

// UserProfileData 用户资料响应结构。
type UserProfileData struct {
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
	NickName string `json:"nickname"`
	Avatar   string `json:"avatar_url"`
	Birth    string `json:"birth"`
	Gender   string `json:"gender"`
	Level    uint64 `json:"level"`
}

// NewUserProfileData 创建一个新的用户资料响应。
//
// 参数：
//   - model：用户资料模型
//
// 返回值：
//   - *UserProfileData：新的用户资料响应结构体。
func NewUserProfileData(model *models.UserInfo) *UserProfileData {
	profile := new(UserProfileData)
	profile.UID = uint64(model.ID)
	profile.Username = model.UserName
	if model.NickName != nil {
		profile.NickName = *model.NickName
	} else {
		profile.NickName = model.UserName
	}
	profile.Avatar = fmt.Sprintf("/resources/avatar/%s.jpg", model.Avatar)
	if model.Birth != nil {
		profile.Birth = model.Birth.Format("2006-01-02")
	} else {
		profile.Birth = "未知"
	}
	if model.Gender != nil {
		profile.Gender = *model.Gender
	} else {
		profile.Gender = "未知"
	}
	profile.Level = model.Level

	return profile
}
