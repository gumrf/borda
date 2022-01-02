package main

import (
	"borda/pkg/postgres"
	"fmt"
	"log"
)

func main(){
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
}