package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/deckarep/golang-set"
	"time"
	"strconv"
	"fmt"
	"log"
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
		log.Fatal("failed to connect database")
	}
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
		if online.Status != uint(user.Online) {
			err := dbConnection.Create(&OnlineModel{UserID: user_id, Status: uint(user.Online), CreatedAt: last_seen_time}).Error
			fmt.Println(err)
		}
	} else
	{
		err := dbConnection.Create(&OnlineModel{UserID: user_id, Status: uint(user.Online), CreatedAt: last_seen_time}).Error
		fmt.Println(err)
	}

}

func UpdateFriends(user_id uint, friends []Friend) {
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
	var userIds []uint
	for _, friend := range needAdd.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID: user_id, FriendID: friend.(uint), Status: 1, CreatedAt: timeCr})
		fmt.Println(err)
		userIds = append(userIds, uint(friend))
	}
	for _, friend := range needDelete.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID: user_id, FriendID: friend.(uint), Status: 0, CreatedAt: timeCr})
		fmt.Println(err)
	}
	//check users
	var users []UserModel
	dbConnection.Find(&users, userIds)

	fmt.Println(users)

	fmt.Println(needAdd)
	fmt.Println(needDelete)
}
