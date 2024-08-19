/*
Package controllers - NekoBlog backend server controllers.
This file is for factory of controllers.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
- sjyhlxysybzdhxd<2023122308@jou.edu.cn>
*/
package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
)

// Factory 控制器工厂
type Factory struct {
	serviceFactory *services.Factory
}

// NewFactory 创建控制器工厂
func NewFactory(serviceFactory *services.Factory) *Factory {
	return &Factory{serviceFactory: serviceFactory}
}



