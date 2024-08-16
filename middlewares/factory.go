/*
Package middlewares - NekoBlog backend server middlewares.
This file is for middlewares factory.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package middlewares

import "github.com/Kirisakiii/neko-micro-blog-backend/stores"

type Factory struct {
	store *stores.Factory
}

func NewFactory(store *stores.Factory) *Factory {
	return &Factory{store: store}
}
