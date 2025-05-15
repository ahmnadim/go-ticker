package repositories

import (
	"os"
	"poc/initializers"
	"poc/models"
	"strconv"
)

// fetchDataFromAPI makes a GET request and returns the response body as string
func FetchDataFromDB(offset int) ([]models.User, error) {
	limit, err := strconv.Atoi(os.Getenv("DATA_LIMIT"))
	if err != nil {
		limit = 10
	}
	var users []models.User
	if err := initializers.DB.
		Select("id, username, email, phone, address, created_at, updated_at").
		Offset(offset).
		Limit(limit).
		Order("created_at asc, id asc").
		Find(&users).Error; err != nil {
		return nil, err
	}
	// fmt.Println("users: ", offset, users)
	return users, nil
}
