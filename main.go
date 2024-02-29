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
	_view.Home(mux, db)
	_view.ManufacturerForm(mux, db)
	_view.Manufacturer(mux, db)
	_view.Equipment(mux, db)
	_view.GetEquipmentQRCode(mux, db)
	_view.EquipmentTicket(mux, db)
	_view.TicketForm(mux, db)
	_view.Tickets(mux, db)
	_view.Ticket(mux, db)

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
	_api.UpdateTicketPublicDetails(mux, db)
	_api.AssignTicket(mux, db)
	_api.TicketResetEquipment(mux, db)

	_partial.EquipmentSelectionList(mux, db)
	_partial.EquipmentByManufacturer(mux, db)
	_partial.Div(mux)

	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), mux)

}
