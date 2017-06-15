package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/deckarep/golang-set"
	"time"
	"strconv"
	"fmt"
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

type User struct {
	Id         uint
	First_name string
	Last_name  string
	Online     int64
	Last_seen struct {
		Time int64
	}
}
type Friend struct {
	Id         uint
	First_name string
	Last_name  string
}
type Friends struct {
	//Count int64
	Items []Friend
}
type ResponseUser struct {
	Users []User `json:"response"`
}
type ResponseFriend struct {
	Friends Friends `json:"response"`
}

var dbConnection *gorm.DB

func OpenConnection() {
	var err error
	dbConnection, err = gorm.Open("postgres", "host=localhost user=vk_status_stat_user dbname=vk_status_stat_dev sslmode=disable password=12345678")
	if err != nil {
		panic("failed to connect database")
	}
	//dbConnection.AutoMigrate(&UserModel{}, &FriendModel{}, &OnlineModel{})
	//dbConnection.Create(&UserModel{VkId: 102831893})
}

func CloseConnection() {
	dbConnection.Close()
}

func GetTrackedUserIds() []string {
	var users []UserModel
	dbConnection.Select("vk_id").Find(&users, UserModel{Tracked: true})
	userIds := []string{}
	for _, user := range users {
		userIds = append(userIds, strconv.FormatUint(uint64(user.VkId), 10))
	}
	return userIds
}

func SetUserOnline(user_id uint, user User) {
	var online OnlineModel
	last_seen_time := time.Unix(user.Last_seen.Time, 0)
	if !dbConnection.Order("created_at desc").First(&online, OnlineModel{UserID: user_id}).RecordNotFound() {
		if online.Status == 1 && user.Online == 0 {
			err := dbConnection.Create(&OnlineModel{UserID: user_id, Status: 0, CreatedAt: last_seen_time}).Error
			fmt.Println(err)
		} else if online.Status == 0 && user.Online == 1 {
			err := dbConnection.Create(&OnlineModel{UserID: user_id, Status: 1, CreatedAt: last_seen_time}).Error
			fmt.Println(err)
		}
	} else
	{
		var status uint = 1
		if user.Online == 0 {
			status = 0
		}
		err := dbConnection.Create(&OnlineModel{UserID: user_id, Status: status, CreatedAt: last_seen_time}).Error
		fmt.Println(err)
	}

}

func UpdateFriends(user_id uint, friends []Friend)  {
	var friendModels []FriendModel
	oldFriends := mapset.NewSet()
	if !dbConnection.Select("friend_id").Order("created_at desc").Group("friend_id, created_at").Find(&friendModels, FriendModel{UserID: user_id}).RecordNotFound() {
		fmt.Print()
		for _, friend := range friendModels {
			oldFriends.Add(friend.FriendID)
		}
	}
	newFriends := mapset.NewSet()
	for _, friend := range friends {
		newFriends.Add(friend.Id)
	}
	needAdd := newFriends.Difference(oldFriends)
	needDelete := oldFriends.Difference(newFriends)

	timeCr := time.Now().UTC()

	for _, friend := range needAdd.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID:user_id, FriendID:friend.(uint), Status:1, CreatedAt:timeCr})
		fmt.Println(err)
	}
	for _, friend := range needDelete.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID:user_id, FriendID:friend.(uint), Status:0, CreatedAt:timeCr})
		fmt.Println(err)
	}

	fmt.Println(needAdd)
	fmt.Println(needDelete)
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
