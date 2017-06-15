package main

import (
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
	"net/url"
	"strings"
	"models"
	"strconv"
)

var METHOD_USER_GET = "users.get"
var METHOD_FRIENDS_GET = "friends.get"

func main() {
	models.OpenConnection()
	defer models.CloseConnection()

	userIds := models.GetTrackedUserIds()

	fields := []string{
		"online", "last_seen",
	}

	params := url.Values{}
	params.Add("user_ids", strings.Join(userIds, ","))
	params.Add("fields", strings.Join(fields, ","))




	//url := "?user_ids=102831893,111&fields=online,last_seen&access_token=b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7&v=5.65"
	//url := "https://api.vk.com/method/friends.get?user_id=102831893&fields=first_name,last_name&access_token=b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7&v=5.65"
	var v models.ResponseUser
	err := json.Unmarshal(makeApiRequest(METHOD_USER_GET, &params), &v)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range v.Users {
		models.SetUserOnline(user.Id, user)
		updateFriendsRequest(user.Id)
	}
	fmt.Println(v)



}

func updateFriendsRequest(user_id uint) {
	fields := []string{
		"online", "last_seen",
	}

	params := url.Values{}
	params.Add("user_id", strconv.FormatUint(uint64(user_id), 10))
	params.Add("fields", strings.Join(fields, ","))

	var friendsResponse models.ResponseFriend

	err := json.Unmarshal(makeApiRequest(METHOD_FRIENDS_GET, &params), &friendsResponse)
	if err != nil {
		log.Fatal(err)
	}

	models.UpdateFriends(user_id, friendsResponse.Friends.Items)
}


func makeApiRequest(method string, params *url.Values) []byte {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiUrl := "https://api.vk.com/method/"

	baseUrlRaw := apiUrl + method

	baseUrl, err := url.Parse(baseUrlRaw)
	if err != nil {
		log.Fatal(err)
	}

	accessToken := "b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7"
	apiVersion := "5.65"

	params.Add("access_token", accessToken)
	params.Add("v", apiVersion)

	baseUrl.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return text
}