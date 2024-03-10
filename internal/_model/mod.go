package _model

import (
	"math/rand"
	"time"

	"gorm.io/gorm"
)

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
	Archived 	   bool
	Photo          []byte
	ManufacturerID uint
	Manufacturer   Manufacturer `gorm:"constraint:OnDelete:CASCADE;"`
	Tickets        []Ticket `gorm:"constraint:OnDelete:CASCADE;"` // Add the foreign key constraint
}

type TicketStatus string

const (
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
	Cost 	    float64
	RepairNotes string
}

func CreateTickets(db *gorm.DB, count int, creator, item, problem, location string, priority TicketPriority, status TicketStatus, notes, owner string, equipmentID uint, cost float64, repairNotes string) error {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		ticket := Ticket{
			Creator:     creator,
			Item:        item,
			Problem:     problem,
			Location:    location,
			Priority:    priority,
			Status:      status,
			Notes:       notes,
			Owner:       owner,
			EquipmentID: &equipmentID,
			Cost:        cost,
			RepairNotes: repairNotes,
		}
		if err := db.Create(&ticket).Error; err != nil {
			return err
		}
	}
	return nil
}