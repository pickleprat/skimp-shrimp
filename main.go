package main

import (
	"cfasuite/internal/_form"
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

	// loading env
	_ = godotenv.Load()
	
	// db connection
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&_model.Manufacturer{})
	db.AutoMigrate(&_model.Equipment{})
	db.AutoMigrate(&_model.Ticket{})
	// _model.CreateTickets(
	// 	db,
	// 	100,
	// 	"admin",
	// 	"laptop",
	// 	"broken screen",
	// 	"Southroads",
	// 	"urgent",
	// 	"complete",
	// 	"none",
	// 	"admin",
	// 	1,
	// 	100.00,
	// 	"none",
	// )
	
	// setting up server and serving static files
	mux := http.NewServeMux()
	_view.ServeStaticFilesAndFavicon(mux)
	
	// public views
	_view.Login(mux, db)
	_view.PublicCreateTicket(mux, db)
	_view.PublicViewTickets(mux, db)

	// admin views
	_view.AdminHome(mux, db)
	_view.UpdateEquipment(mux, db)
	_view.DeleteEquipment(mux, db)
	_view.Manufacturer(mux, db)
	_view.UpdateManufacturer(mux, db)
	_view.DeleteManufacturer(mux, db)
	_view.CreateManufacturers(mux, db)
	_view.EquipmentArchive(mux, db)
	_view.AdminViewTicket(mux, db)
	_view.AdminCreateTickets(mux, db)
	_view.AdminUpdateTicket(mux, db)
	_view.AdminDeleteTicket(mux, db)
	_view.AdminCompleteTicket(mux, db)

	// forms
	_form.Login(mux, db)
	_form.Logout(mux, db)
	_form.CreateManufacturer(mux, db)
	_form.DeleteManufacturer(mux, db)
	_form.UpdateManufacturer(mux, db)
	_form.CreateEquipment(mux, db)
	_form.UpdateEquipment(mux, db)
	_form.DeleteEquipment(mux, db)
	_form.CreateTicketPublic(mux, db)
	_form.CreateTicketAdmin(mux, db)
	_form.UpdateTicket(mux, db)
	_form.DeleteTicket(mux, db)
	_form.AssignTicket(mux, db)
	_form.TicketResetEquipment(mux, db)
	_form.CompleteTicket(mux, db)

	// partials
	_partial.EquipmentSelectionList(mux, db)
	_partial.EquipmentByManufacturer(mux, db)
	_partial.ManufactuerList(mux, db)
	_partial.TicketList(mux, db)
	_partial.ResetEquipmentLink(mux, db)
	_partial.PublicTicketList(mux, db)
	_partial.AuthWarning(mux, db)
	_partial.CompletedTicketList(mux, db)

	// serving
	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		fmt.Println(err)
	}

}
