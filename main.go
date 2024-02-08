package main

import (
	"cfasuite/internal/_util"
	"cfasuite/internal/_view"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	_util.ServeStaticFilesAndFavicon()
	_view.Home(db)


	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}