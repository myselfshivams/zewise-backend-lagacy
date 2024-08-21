package models

import "gorm.io/gorm"

// AvatarDeletionWaitList 头像删除等待列表
type AvatarDeletionWaitList struct {
	gorm.Model        // 基本模型
	FileName   string `gorm:"column:file_name"` // 文件名
}
