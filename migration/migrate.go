package main

import (
	"log"
	"poc/initializers"
	"poc/models"
)

func init() {
	log.Default().Println("Loading env...")
	initializers.LoadEnvVars()

	log.Default().Println("Connecting to the DB...")
	initializers.ConnectDB()

}

func main() {
	log.Default().Println("Migrating...")

	initializers.DB.AutoMigrate(&models.User{})
	log.Default().Println("Migration completed.")

}
