package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	var err error

	// User
	if err = db.AutoMigrate(&UserInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserAuthInfo{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserLoginLog{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&UserAvaliableToken{}); err != nil {
		return err
	}

	return nil
}
