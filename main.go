package main

import (
	"cfasuite/internal/_api"
	"cfasuite/internal/_model"
	"cfasuite/internal/_partial"
	"cfasuite/internal/_view"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	_ = godotenv.Load()
	mux := http.NewServeMux()

	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&_model.Manufacturer{})
	db.AutoMigrate(&_model.Equipment{})
	db.AutoMigrate(&_model.Ticket{})

	_view.ServeStaticFilesAndFavicon(mux)
	
	// public views
	_view.Login(mux, db)
	_view.PublicCreateTicket(mux, db)
	_view.PublicViewTickets(mux, db)
	
	// admin home
	_view.AdminHome(mux, db)

	// admin equipment views
	_view.UpdateEquipment(mux, db)
	_view.DeleteEquipment(mux, db)

	// admin manufacturer views
	_view.Manufacturer(mux, db)
	_view.UpdateManufacturer(mux, db)
	_view.DeleteManufacturer(mux, db)
	_view.CreateManufacturers(mux, db)
	_view.EquipmentArchive(mux, db)

	// admin ticket views
	_view.AdminViewTicket(mux, db)
	_view.AdminCreateTickets(mux, db)
	_view.AdminUpdateTicket(mux, db)
	_view.AdminDeleteTicket(mux, db)
	_view.AdminCompleteTicket(mux, db)

	_api.Login(mux, db)
	_api.Logout(mux, db)
	_api.CreateManufacturer(mux, db)
	_api.DeleteManufacturer(mux, db)
	_api.UpdateManufacturer(mux, db)
	_api.CreateEquipment(mux, db)
	_api.UpdateEquipment(mux, db)
	_api.DeleteEquipment(mux, db)
	_api.CreateTicketPublic(mux, db)
	_api.CreateTicketAdmin(mux, db)
	_api.UpdateTicket(mux, db)
	_api.DeleteTicket(mux, db)
	_api.AssignTicket(mux, db)
	_api.TicketResetEquipment(mux, db)
	_api.CompleteTicket(mux, db)

	// partials
	_partial.EquipmentSelectionList(mux, db)
	_partial.EquipmentByManufacturer(mux, db)
	_partial.Div(mux)
	_partial.ManufactuerList(mux, db)
	_partial.TicketList(mux, db)
	_partial.ResetEquipmentLink(mux, db)
	_partial.PublicTicketList(mux, db)

	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
