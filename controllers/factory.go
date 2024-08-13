package controllers

import (
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
)

type Factory struct {
	serviceFactory *services.Factory
}

func NewFactory(serviceFactory *services.Factory) *Factory {
	return &Factory{serviceFactory: serviceFactory}
}
