package stores

import "gorm.io/gorm"

type Factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{db: db}
}
