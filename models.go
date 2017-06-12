package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	VkId      uint
	FirstName string
	LastName  string
	Tracked   bool
	CreatedAt time.Time
}

type Friend struct {
	ID        uint `gorm:"primary_key"`
	User      User
	UserID    uint
	Friend    User
	FriendID  uint
	status    uint
	CreatedAt time.Time
}

func main() {
	db, err := gorm.Open("postgres", "host=localhost user=vk_status_stat_user dbname=vk_status_stat_dev sslmode=disable password=12345678")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Create
	db.Create(&User{FirstName: "L1212"})

	// Read
	var product User
	db.First(&product, 1)                         // find product with id 1
	db.First(&product, "first_name = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	db.Delete(&product)
}
