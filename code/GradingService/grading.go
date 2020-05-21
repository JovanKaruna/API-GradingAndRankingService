package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Doc struct {
	Data []Users `json:"data"`
}

type Users struct {
	Idexam int    `json:"id_exam"`
	Iduser int    `json:"id_user"`
	Sheets []Exam `json:"answersheets"`
}

type Exam struct {
	Idsession int     `json:"id_session"`
	Answers   Session `json:"answers"`
}

type Session struct {
	One   bool `json:"1"`
	Two   bool `json:"2"`
	Three bool `json:"3"`
	Four  bool `json:"4"`
	Five  bool `json:"5"`
	Six   bool `json:"6"`
	Seven bool `json:"7"`
	Eight bool `json:"8"`
	Nine  bool `json:"9"`
	Ten   bool `json:"10"`
}

type Data struct {
	Data []Score `json:"data"`
}
type Score struct {
	Iduser     int `json:"id_user"`
	EventScore int `json:"event_score"`
}

func inputScore(id int, score int) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/grading")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO grade VALUES (?, ?)", id, score)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

}

func grading(value Doc) {
	for i := 0; i < len(value.Data); i++ { //There are 500 data users and 500 examsheets
		event_score := 0
		for j := 0; j < 10; j++ { //10 session / event
			session_score := 0
			if value.Data[i].Sheets[j].Answers.One == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Two == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Three == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Four == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Five == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Six == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Seven == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Eight == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Nine == true {
				session_score += 1
			}
			if value.Data[i].Sheets[j].Answers.Ten == true {
				session_score += 1
			}
			fmt.Printf("session %d : %d \n", value.Data[i].Sheets[j].Idsession, session_score)
			event_score += session_score
		}
		fmt.Printf("iduser %d : %d \n", value.Data[i].Iduser, event_score)
		//Input id and eventscore into mysql local
		inputScore(value.Data[i].Iduser, event_score)
	}
	fmt.Println("Done Grading")
}

func gradingService() { //This is used for grading service
	var score = Doc{}
	filename := "examsheets.json"
	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer jsonFile.Close()
	examSheet, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(examSheet), &score)

	grading(score)
}

func getScore() { //API to fetch the score from local MySQL
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/grading")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM grade")
	if err != nil {
		panic(err.Error())
	}

	var scores Data

	for result.Next() {
		var score Score
		err := result.Scan(&score.Iduser, &score.EventScore)
		if err != nil {
			panic(err.Error())
		}
		scores.Data = append(scores.Data, score)

	}
	data, _ := json.MarshalIndent(scores, "", "	")
	fmt.Println(string(data))                      //data is in Json Format
	_ = ioutil.WriteFile("score.json", data, 0644) //Write the data json into a file
}
func main() {
	//If you already grade the examsheets and input the data into mysql then comment the gradingService()
	gradingService() //Takes examsheets.json, grade the scores, input the scores in mysql Local

	getScore() //get score from mysql local as a Json file named "score.json" in the same directory
}
