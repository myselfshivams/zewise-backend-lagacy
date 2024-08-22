/*
Package models - NekoBlog backend server database models
This file is for comment related models.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// CommentInfo 评论信息模型
type CommentInfo struct {
	gorm.Model               // 基本模型
	PostID     uint64        `gorm:"column:post_id"`                // 博文ID
	UID        uint64        `gorm:"column:uid"`                    // 用户ID
	Username   string        `gorm:"column:username"`               // 用户名
	Content    string        `gorm:"column:content"`                // 内容
	Like       pq.Int64Array `gorm:"column:like;type:bigint[]"`     // 点赞数 记录UID
	Dislike    pq.Int64Array `gorm:"column:dislike;type:bigint[]"`  // 踩数 记录UID
	IsPublic   bool          `gorm:"column:is_public;default:true"` // 是否公开
	// Share   uint64 `gorm:"column:share"`                         // 分享数 暂时不实现
}


