package main

import (
	"cfasuite/internal/_model"
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

	_view.ServeStaticFilesAndFavicon(mux)
	_view.Home(mux, db)
	_view.Login(mux, db)
	_view.App(mux, db)
	_view.Logout(mux, db)
	_view.CreateManufacturer(mux, db)
	_view.Manufacturer(mux, db)
	_view.DeleteManufacturer(mux, db)
	_view.UpdateManufacturer(mux, db)
	_view.CreateEquipment(mux, db)
	_view.Equipment(mux, db)
	_view.UpdateEquipment(mux, db)
	_view.DeleteEquipment(mux, db)
	_view.GetEquipmentQRCode(mux, db)
	_view.EquipmentTicket(mux, db)


	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	http.ListenAndServe(":" + os.Getenv("PORT"), mux)

}