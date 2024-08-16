/*
Package stores - NekoBlog backend server data access objects.
This file is for factory of storages.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package stores

import "gorm.io/gorm"

type Factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{db: db}
}
