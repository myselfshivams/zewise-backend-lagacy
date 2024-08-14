package services

import "github.com/Kirisakiii/neko-micro-blog-backend/stores"

type Factory struct {
	storeFactory *stores.Factory
}

func NewFactory(storeFactory *stores.Factory) *Factory {
	return &Factory{storeFactory: storeFactory}
}
