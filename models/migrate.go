package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&UserInfo{})
	db.AutoMigrate(&UserAuthInfo{})
	db.AutoMigrate(&UserLoginLog{})
	db.AutoMigrate(&UserAuthToken{})
}
