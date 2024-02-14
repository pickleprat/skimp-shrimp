package _model

import "gorm.io/gorm"

type Manufacturer struct {
    gorm.Model
    Name      string
    Email     string
    Phone     string
    Equipment []Equipment `gorm:"foreignKey:ManufacturerID"` // Define the foreign key relationship
}


type Equipment struct {
    gorm.Model
    Nickname     string
    SerialNumber string
    Photo        string // Assuming the photo is stored as a string (e.g., file path or URL)
    ManufacturerID uint // Foreign key
    Manufacturer Manufacturer `gorm:"constraint:OnDelete:CASCADE;"` // Define the foreign key constraint and cascade delete
}

