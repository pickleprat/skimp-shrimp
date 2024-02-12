package main

import (
	"cfasuite/internal/_model"
	"cfasuite/internal/_util"
	"cfasuite/internal/_view"
	"fmt"
	"net/http"

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

	_util.ServeStaticFilesAndFavicon(mux)
	_view.Home(mux, db)
	_view.Login(mux, db)
	_view.App(mux, db)


	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)

}