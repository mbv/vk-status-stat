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
)

func main() {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiUrl := "https://api.vk.com/method/"
	methodUsers := "users.get"

	baseUrlRaw := apiUrl + methodUsers

	baseUrl, err := url.Parse(baseUrlRaw)
	if err != nil {
		log.Fatal(err)
	}

	userIds := []string{
		"102831893", "87854589", "69971049",
	}
	fields := []string{
		"online", "last_seen",
	}

	accessToken := "b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7"
	apiVersion := "5.65"


	params := url.Values{}
	params.Add("user_ids", strings.Join(userIds, ","))
	params.Add("fields", strings.Join(fields, ","))
	params.Add("access_token", accessToken)
	params.Add("v", apiVersion)

	baseUrl.RawQuery = params.Encode()

	//url := "?user_ids=102831893,111&fields=online,last_seen&access_token=b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7&v=5.65"
	//url := "https://api.vk.com/method/friends.get?user_id=102831893&fields=first_name,last_name&access_token=b6c7d77ab6c7d77ab6c7d77a55b69b4be3bb6c7b6c7d77aeffb36e5c2305e428ee15bd7&v=5.65"

	req, err := http.NewRequest("GET", baseUrl.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	type User struct {
		Id         int64
		First_name string
		Last_name  string
		Online     int64
		Last_seen struct {
			Time int64
		}
	}
	type Friend struct {
		Id         int64
		First_name string
		Last_name  string
	}
	type Friends struct {
		//Count int64
		Items []Friend
	}
	type response struct {
		Users []User `json:"response"`
	}
	type response1 struct {
		Friends Friends `json:"response"`
	}
	var v response
	err = json.Unmarshal(text, &v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(v)

}
