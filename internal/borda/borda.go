package borda

import (
	"borda/internal/controller"
	"borda/pkg/postgres"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	db, err := postgres.NewPostgresDatabase(postgres.Options{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal("postgres.NewPostgresDatabase:", err)
	}

	fmt.Printf("%+v", db.Stats())

	// TODO: Close connection
	router := gin.Default()
	controller.NewController(router)
}
