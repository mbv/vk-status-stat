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
	Id        uint `gorm:"primary_key"`
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
	ChangedAt time.Time
}

type FieldChangeModel struct {
	ID          uint `gorm:"primary_key"`
	User        UserModel
	UserID      uint
	Field       string
	ValueBefore string
	CreatedAt   time.Time
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
	dbConnection.AutoMigrate(&UserModel{}, &FriendModel{}, &OnlineModel{}, &FieldChangeModel{})
	//dbConnection.Create(&UserModel{Id:102831893, Tracked:true})
}

func CloseConnection() {
	dbConnection.Close()
}

func GetTrackedUserIds() []string {
	var users []UserModel
	dbConnection.Select("id").Find(&users, UserModel{Tracked: true})
	userIds := []string{}
	for _, user := range users {
		userIds = append(userIds, strconv.FormatUint(uint64(user.Id), 10))
	}
	return userIds
}

func SetUserOnline(user User) {
	var online OnlineModel
	last_seen_time := time.Unix(user.Last_seen.Time, 0)
	var onlineModel = OnlineModel{UserID: user.Id, Status: uint(user.Online), ChangedAt: last_seen_time}
	if !dbConnection.Order("created_at desc").First(&online, OnlineModel{UserID: user.Id}).RecordNotFound() {
		if online.Status != uint(user.Online) {
			err := dbConnection.Create(&onlineModel).Error
			fmt.Println(err)
		}
	} else {
		err := dbConnection.Create(&onlineModel).Error
		fmt.Println(err)
	}
}

func CheckFieldsChange(user User) {
	var userModel UserModel
	if !dbConnection.First(&userModel, user.Id).RecordNotFound() {
		timeCr := time.Now().UTC()
		if userModel.FirstName != user.First_name {
			dbConnection.Create(&FieldChangeModel{
				UserID:      user.Id,
				Field:       "FirstName",
				ValueBefore: userModel.FirstName,
				CreatedAt:   timeCr,
			})
			userModel.FirstName = user.First_name
		}
		if userModel.LastName != user.Last_name {
			dbConnection.Create(&FieldChangeModel{
				UserID:      user.Id,
				Field:       "LastName",
				ValueBefore: userModel.LastName,
				CreatedAt:   timeCr,
			})
			userModel.LastName = user.Last_name
		}
		dbConnection.Save(&userModel)
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
	var userVkIds []uint
	for _, friend := range needAdd.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID: user_id, FriendID: friend.(uint), Status: 1, CreatedAt: timeCr})
		fmt.Println(err)
		userVkIds = append(userVkIds, friend.(uint))
	}
	for _, friend := range needDelete.ToSlice() {
		err := dbConnection.Create(&FriendModel{UserID: user_id, FriendID: friend.(uint), Status: 0, CreatedAt: timeCr})
		fmt.Println(err)
	}
	//check users
	var userModels []UserModel
	dbConnection.Find(&userModels, userVkIds)
	usersInDb := mapset.NewSet()
	for _, user := range userModels {
		usersInDb.Add(user.Id)
	}

	needAddUsers := needAdd.Difference(usersInDb)
	for _, friend := range friends {
		if needAddUsers.Contains(friend.Id) {
			err := dbConnection.Create(&UserModel{Id: friend.Id, FirstName: friend.First_name, LastName: friend.Last_name, Tracked: false})
			fmt.Println(err)
		}
	}
	fmt.Println(usersInDb)

	fmt.Println(needAdd)
	fmt.Println(needDelete)
}
