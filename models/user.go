/*
Package models - NekoBlog backend server database models
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	gorm.Model            // 基本模型
	UserName   string     `gorm:"unique;column:username"`        // 用户名
	NickName   *string    `gorm:"column:nickname"`               // 昵称
	Avatar     string     `gorm:"default:vanilla;column:avatar"` // 头像
	Birth      *time.Time `gorm:"column:birth"`                  // 生日
	Gender     *string    `gorm:"column:gender"`                 // 性别
	Authority  uint64     `gorm:"default:0;column:authority"`    // 权限等级
	Level      uint64     `gorm:"default:1;column:level"`        // 等级
}

type UserAuthInfo struct {
	gorm.Model          // 基本模型
	UID          uint64 `gorm:"unique;column:uid"`      // 用户ID
	UserName     string `gorm:"unique;column:username"` // 用户名
	Salt         string `gorm:"column:salt"`            // 盐
	PasswordHash string `gorm:"column:psw_hash"`        // 密码哈希值
}

type UserLoginLog struct {
	gorm.Model            // 基本模型
	UID         uint64    `gorm:"unique;column:uid"`                  // 用户ID
	LoginTime   time.Time `gorm:"column:login_time"`                  // 登录时间
	LoginIP     string    `gorm:"column:login_ip"`                    // 登录IP
	Device      string    `gorm:"default:unknown;column:device"`      // 登录时登陆的设备 如：Windows iOS Android
	Application string    `gorm:"default:unknown;column:application"` // 登录时使用的应用 如 Chrome 236.12
}

type UserAuthToken struct {
	gorm.Model            // 基本模型
	UID         uint64    `gorm:"unique;column:uid"`                  // 用户ID
	LoginIP     string    `gorm:"column:login_ip"`                    // 登录IP
	BearerToken string    `gorm:"column:bearer_token"`                // 令牌Token
	Device      string    `gorm:"default:unknown;column:device"`      // 登录时登陆的设备 如：Windows iOS Android
	Application string    `gorm:"default:unknown;column:application"` // 登录时使用的应用 如 Chrome 236.12
	ExpireTime  time.Time `gorm:"column:expire_time"`                 // 令牌过期时间
}
