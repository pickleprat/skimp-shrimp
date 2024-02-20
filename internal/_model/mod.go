package _model

import "gorm.io/gorm"

type Manufacturer struct {
    gorm.Model
    Name        string
    Email       string
    Phone       string
    Equipment   []Equipment `gorm:"foreignKey:ManufacturerID"`
}

type Equipment struct {
    gorm.Model
    Nickname       string
    SerialNumber   string
    Photo          []byte
    ManufacturerID uint
    Manufacturer   Manufacturer `gorm:"constraint:OnDelete:CASCADE;"`
    QRCodeToken    string
    Tickets        []Ticket `gorm:"constraint:OnDelete:CASCADE;"` // Add the foreign key constraint
}

type Ticket struct {
    gorm.Model
    Creator       string
    Item          string
    Problem       string
    Location      string
    Photo         []byte
    EquipmentID   *uint
    Equipment     Equipment
}
