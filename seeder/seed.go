package main

import (
	"fmt"
	"log"
	"math/rand"
	"poc/initializers"
	"poc/models"
	"time"

	"github.com/bxcodec/faker/v4"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
}

func RandomTimeWithinLastNDays(days int) time.Time {
	now := time.Now()
	// Generate a random duration within N days
	duration := time.Duration(rand.Intn(days*24)) * time.Hour
	return now.Add(-duration)
}

func main() {
	count := 1000 //os.Getenv("USER_SEED_COUNT")
	rand.Seed(time.Now().UnixNano())
	users := make([]models.User, 0, 1000)

	for i := 0; i < count; i++ {
		timestamp := RandomTimeWithinLastNDays(20)

		user := models.User{
			// ID:        uint(i),
			Email:     faker.Email(), // You might want to make it unique
			Username:  faker.Username(),
			Phone:     faker.Phonenumber(),
			Address:   faker.Word(),
			CreatedAt: timestamp,
			UpdatedAt: timestamp.Add(time.Duration(rand.Intn(6*60)) * time.Minute),
		}
		users = append(users, user)
	}

	if err := initializers.DB.Create(&users).Error; err != nil {
		log.Fatal("Failed to seed users:", err)
	}

	fmt.Printf("Seeded %d users successfully.\n", count)

}
