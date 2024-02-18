package main

import (
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"cfasuite/internal/_view"
	"fmt"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	
	// err = godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// }
	mux := http.NewServeMux()

	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&_model.Manufacturer{})
	db.AutoMigrate(&_model.Equipment{})

	fmt.Println("db connected")
	fmt.Println("env variables")
	fmt.Println("PORT: " + os.Getenv("PORT"))
	fmt.Println("ADMIN_USERNAME: " + os.Getenv("ADMIN_USERNAME"))
	fmt.Println("ADMIN_PASSWORD: " + os.Getenv("ADMIN_PASSWORD"))
	fmt.Println("ADMIN_SESSION_TOKEN: " + os.Getenv("ADMIN_SESSION_TOKEN"))

	_util.ServeStaticFilesAndFavicon(mux)
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
	_view.EquipmentSettingsForm(mux, db)
	_view.ClearComponent(mux, db)
	_view.ClientRedirect(mux, db)
	_view.UpdateEquipment(mux, db)


	fmt.Println("Server is running on port " + os.Getenv("PORT"))
	http.ListenAndServe(":" + os.Getenv("PORT"), mux)

}