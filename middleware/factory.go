package middleware

import "github.com/Kirisakiii/neko-micro-blog-backend/stores"

type Factory struct {
	store *stores.Factory
}

func NewFactory(store *stores.Factory) *Factory {
	return &Factory{store: store}
}
