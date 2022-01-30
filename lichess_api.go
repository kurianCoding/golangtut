package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	Name  string `json:"username"`
	Perfs struct {
		Blitz struct {
			Games int64 `json:"games"`
		} `json:"blitz"`
		Puzzle struct {
			Games int64 `json:"games"`
		} `json:"puzzle"`
		Bullet struct {
			Games int64 `json:"games"`
		} `json:"bullet"`
		Correspondence struct {
			Games int64 `json:"games"`
		} `json:"correspondence"`
		Classical struct {
			Games int64 `json:"games"`
		} `json:"classical"`
		Rapid struct {
			Games int64 `json:"games"`
		} `json:"rapid"`
	} `json:"perfs"`
	TotalGames int64
}

var token string  // TODO make this part of struct
var APP_NM string //TODO make this constant and in code
func init() {
	token := os.Getenv("LCHSSTKN")
	if token == "" {
		panic("no token given")
	}
	APP_NM = os.Getenv("APP_NM")
}

func GetApi(uri string) (response []byte, err error) {
	client := http.Client{} // TODO set time out
	request, err := http.NewRequest("GET", uri, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Add("User-Agent", APP_NM)
	if err != nil {
		return
	}

	res, err := client.Do(request)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("non ok status %d", res.StatusCode)) //TODO log error and framework to handle error
		return
	}

	if res.Body == nil {
		err = errors.New("empty body") //TODO log error and framework to handle error
		return
	}

	defer res.Body.Close()
	response, err = ioutil.ReadAll(res.Body)

	return
}

//THE FUNCTION GETS DETAILS OF A PLAYER
func GetUser(name string) (u User) {
	res, err := GetApi(fmt.Sprintf("https://lichess.org/api/user/%s", name))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res, &u)
	p := u.Perfs
	u.TotalGames = p.Blitz.Games + p.Puzzle.Games + p.Bullet.Games + p.Correspondence.Games + p.Classical.Games + p.Rapid.Games
	return
}
