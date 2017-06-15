package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"strconv"
)

type UserModel struct {
	ID        uint `gorm:"primary_key"`
	VkId      uint
	FirstName string
	LastName  string
	Tracked   bool
	CreatedAt time.Time
}

type FriendModel struct {
	ID        uint `gorm:"primary_key"`
	User      UserModel
	UserID    uint
	Friend    UserModel
	FriendID  uint
	Status    uint
	CreatedAt time.Time
}
type OnlineModel struct {
	ID        uint `gorm:"primary_key"`
	User      UserModel
	UserID    uint
	Status    uint
	CreatedAt time.Time
}

var dbConnection *gorm.DB

func OpenConnection() {
	var err error
	dbConnection, err = gorm.Open("postgres", "host=localhost user=vk_status_stat_user dbname=vk_status_stat_dev sslmode=disable password=12345678")
	if err != nil {
		panic("failed to connect database")
	}
	dbConnection.AutoMigrate(&UserModel{}, &FriendModel{}, &OnlineModel{})
}

func CloseConnection()  {
	dbConnection.Close()
}

func GetTrackedUserIds() []string {
	var users []UserModel
	dbConnection.Select("vk_id").Find(&users, UserModel{Tracked: true})
	userIds := []string{}
	for _, user:= range users {
		userIds = append(userIds, strconv.FormatUint(uint64(user.VkId), 10))
	}
	return userIds
}


func SetUserOnline(user_id uint, user main.User) {
	var online OnlineModel
	dbConnection.Order("created_at").Find(&online, OnlineModel{UserID:user_id})
	if online {

	} else	if user.Online == 1 {
		dbConnection.Create(&OnlineModel{UserID:user_id,Status:1,CreatedAt:time.Now()})
	}

}


	/*// Migrate the schema
	db.AutoMigrate(&UserModel{}, &FriendModel{})





	// Create
	db.Create(&UserModel{FirstName: "L1212"})



	// Read
	var product UserModel
	db.First(&product, 1)                         // find product with id 1
	db.First(&product, "first_name = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	//db.Delete(&product)*/
