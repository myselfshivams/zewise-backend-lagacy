package rontines

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AvatarCleanerJob 头像清理任务
type AvatarCleanerJob struct {
	logger *logrus.Logger // 日志记录器
	db     *gorm.DB       // 数据库连接
}

// NewAvatarCleanerJob 创建一个新的头像清理任务。
//
// 参数：
//   - db：数据库连接
//   - logger：日志记录器
//
// 返回值：
//   - *AvatarCleanerJob：新的头像清理任务。
func NewAvatarCleanerJob(logger *logrus.Logger, db *gorm.DB) *AvatarCleanerJob {
	return &AvatarCleanerJob{
		logger: logger,
		db:     db,
	}
}

// Run 执行头像清理任务。
func (job *AvatarCleanerJob) Run() {
	job.logger.Infoln("正在执行头像清理任务...")
	// 获取清除队列
	var waitList []models.AvatarDeletionWaitList
	result := job.db.Find(&waitList)
	if result.Error != nil {
		job.logger.Errorln("获取头像清理队列失败:", result.Error)
		return
	}

	// 清理头像
	for _, item := range waitList {
		err := os.Remove(filepath.Join("./public/avatars", item.FileName))
		// 头像文件不存在
		if errors.Is(err, os.ErrNotExist) {
			job.logger.Errorln("wtf")
			job.logger.Warnln("头像文件不存在:", item.FileName)
			result = job.db.Unscoped().Delete(&item)
			if result.Error != nil {
				job.logger.Errorln("清理头像失败:", result.Error)
			}
			continue
		}
		if err != nil {
			job.logger.Warningln("清理头像失败:", err)
			continue
		}
		// 删除数据库记录
		result = job.db.Unscoped().Delete(&item)
		if result.Error != nil {
			job.logger.Errorln("清理头像失败:", result.Error)
			continue
		}
	}
	job.logger.Infoln("头像清理任务执行完毕")
}
