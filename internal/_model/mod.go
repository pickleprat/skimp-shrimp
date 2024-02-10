package _model

import "gorm.io/gorm"

type Manufacturer struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
	Email string
	Phone string
}

