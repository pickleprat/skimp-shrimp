package _model

import "gorm.io/gorm"

type Manufacturer struct {
	gorm.Model
	Name string
	Email string
	Phone string
}

