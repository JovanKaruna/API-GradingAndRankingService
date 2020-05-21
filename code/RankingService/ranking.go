package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Id struct {
	Data q `json:"data"`
}

type q struct {
	Users []User `json:"Users"`
}

type User struct {
	Iduser   int    `json:"id_user"`
	Username string `json:"username"`
	Name     string `json:"name"`
	School   string `json:"school"`
	City     string `json:"city"`
	Province string `json:"province"`
}

type Rank struct {
	Data []Ranks `json:"data"`
}

type Ranks struct {
	Iduser   int    `json:"id_user"`
	Score    int    `json:"event_score"`
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	Name     string `json:"name"`
	School   string `json:"school"`
	City     string `json:"city"`
	Province string `json:"province"`
}

func getId(id int) User { //Input is the id_user and output is User Type
	url := "https://dummy.edukasystem.id/user"
	var user = Id{}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		log.Fatalln(err)
	}
	json.Unmarshal(body, &user)
	for i := 0; i < len(user.Data.Users); i++ {
		if user.Data.Users[i].Iduser == id {
			return user.Data.Users[i]
		}
	}
	//Not found
	fmt.Println("Not found")
	var not_found = User{}
	return not_found
}

func rankerService() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/grading")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM grade ORDER BY event_score DESC")
	if err != nil {
		panic(err.Error())
	}

	var data Rank
	var rank = 1
	for result.Next() {
		var value Ranks
		err := result.Scan(&value.Iduser, &value.Score)
		if err != nil {
			panic(err.Error())
		}
		var temp User = getId(value.Iduser)
		value.Rank = rank
		value.Username = temp.Username
		value.Name = temp.Name
		value.School = temp.School
		value.City = temp.City
		value.Province = temp.Province
		data.Data = append(data.Data, value)
		rank++
	}
	results, _ := json.MarshalIndent(data, "", "	")
	fmt.Println(string(results))                     //results is in Json Format
	_ = ioutil.WriteFile("rank.json", results, 0644) //Write the data json into a file
}

func findRank(rank int) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/grading")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM grade ORDER BY event_score DESC LIMIT ? , 1", rank-1)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var value Ranks
		err := result.Scan(&value.Iduser, &value.Score)
		if err != nil {
			panic(err.Error())
		}
		var temp User = getId(value.Iduser)
		value.Rank = rank
		value.Username = temp.Username
		value.Name = temp.Name
		value.School = temp.School
		value.City = temp.City
		value.Province = temp.Province
		data, _ := json.MarshalIndent(value, "", "	")
		fmt.Println(string(data)) //results is in Json Format
	}
}

func main() {
	rankerService() //Take score from grading that is stored at local mysql and then rank the score and make a jsonfile with complete user details

	//There are only 500 users so the lowest rank is 500, if you input rank > 500 it will not output anything
	findRank(1) //Change the input
}
