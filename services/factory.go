/*
Package services - NekoBlog backend server services.
This file is for factory of services.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package services

import "github.com/Kirisakiii/neko-micro-blog-backend/stores"

// Factory 服务工厂
type Factory struct {
	storeFactory *stores.Factory
}

// NewFactory 创建服务工厂
//
// 参数：
// storeFactory *stores.Factory - 存储工厂
//
// 返回值：
// *Factory - 服务工厂
func NewFactory(storeFactory *stores.Factory) *Factory {
	return &Factory{storeFactory: storeFactory}
}
