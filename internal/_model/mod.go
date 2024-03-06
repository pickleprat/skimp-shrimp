package _model

import "gorm.io/gorm"

type Manufacturer struct {
	gorm.Model
	Name      string
	Email     string
	Phone     string
	Equipment []Equipment `gorm:"foreignKey:ManufacturerID"`
}

type Equipment struct {
	gorm.Model
	Nickname       string
	SerialNumber   string
	ModelNumber    string
	Photo          []byte
	ManufacturerID uint
	Manufacturer   Manufacturer `gorm:"constraint:OnDelete:CASCADE;"`
	Tickets        []Ticket `gorm:"constraint:OnDelete:CASCADE;"` // Add the foreign key constraint
}

type TicketStatus string

const (
	TicketStatusNew      TicketStatus = "new"
	TicketStatusActive   TicketStatus = "active"
	TicketStatusComplete TicketStatus = "complete"
	TicketStatusOnHold   TicketStatus = "onhold"
)

type TicketPriority string

const (
	TicketPriorityUrgent       TicketPriority = "urgent"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityLow          TicketPriority = "low"
)

type Ticket struct {
	gorm.Model
	Creator     string
	Item        string
	Problem     string
	Location    string
	Priority    TicketPriority
	Status      TicketStatus
	Notes       string
	Owner       string
	EquipmentID *uint
	Equipment   Equipment
}
