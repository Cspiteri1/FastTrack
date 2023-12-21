// handler.go
package handler

import (
	"net/http"
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
)

type Player struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Score int    `json:"score"`
}

type Answer struct {
	AnswerText string
	Valid      bool
}

type Question struct {
	QuestionText string
	Answers      []Answer
}

var players []Player
var questions []Question

func createPlayer(id string, name string, age int, score int) Player {
	p := Player{
		Id:    id,
		Name:  name,
		Age:   age,
		Score: score,
	}

	return p
}

func createQuestion(text string, answer []Answer) Question {
	q := Question{
		QuestionText: text,
		Answers:      answer,
	}

	return q
}

func createAnswer(text string, valid bool) Answer {
	a := Answer{
		AnswerText: text,
		Valid:      valid,
	}

	return a
}

func InitData() {

	players = append(players, createPlayer("000001M", "Clayton", 28, 30))
	players = append(players, createPlayer("000002M", "Giancarl", 27, 10))
	players = append(players, createPlayer("000003M", "Emma", 33, 90))
	players = append(players, createPlayer("000004M", "Paulinha", 30, 50))

	var answersGroupA []Answer = []Answer{createAnswer("Italian", false), createAnswer("Arabic", false), createAnswer("Maltese", true)}
	var answersGroupB []Answer = []Answer{createAnswer("Lira", false), createAnswer("Pound", false), createAnswer("Euro", true)}
	var answersGroupC []Answer = []Answer{createAnswer("Mediterranean ", true), createAnswer("Dead Sea", false), createAnswer("Sea of Samsara", false)}
	var answersGroupD []Answer = []Answer{createAnswer("three", true), createAnswer("one", false), createAnswer("four", false)}
	var answersGroupE []Answer = []Answer{createAnswer("Blue,White and Red", false), createAnswer("White and Red", true), createAnswer("Red and White", false)}

	questions = append(questions, createQuestion("What ls the national language of Malta?", answersGroupA))
	questions = append(questions, createQuestion("What currency is used in Malta?", answersGroupB))
	questions = append(questions, createQuestion("In which sea is Malta located?", answersGroupC))
	questions = append(questions, createQuestion("How many inhabited islands make up the Republic of Malta?", answersGroupD))
	questions = append(questions, createQuestion("What are the colours of Malta's National flag?", answersGroupE))

}

func UpdatePlayer(c *gin.Context) {
	var newPlayer Player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data. Check the provided player details.", "error": err.Error()})
		return
	}

	idx := sort.Search(len(players), func(i int) bool {
		return string(players[i].Id) >= newPlayer.Id
	})

	if idx != len(players) {
		players[idx] = newPlayer
		c.IndentedJSON(http.StatusOK, newPlayer)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User does not exist"})
		return
	}
}

func GetPlayerRank(c *gin.Context) {
	var returnedPlayer Player
	inputId := c.Param("id")
	rank := len(players)

	if reflect.TypeOf(inputId) != reflect.TypeOf(returnedPlayer.Id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	idx := sort.Search(len(players), func(i int) bool {
		return string(players[i].Id) >= inputId
	})

	if idx != len(players) {
		returnedPlayer = players[idx]
		for i := 0; i < len(players); i++ {
			if returnedPlayer.Score > players[i].Score {
				if players[i].Id != returnedPlayer.Id {
					rank -= 1
				}
			}
		}

		c.IndentedJSON(http.StatusOK, rank)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found"})
		return
	}

}

func GetQuestions(c *gin.Context) {
	returnedQuestions := questions

	if len(questions) == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "No questions available"})
		return
	}

	c.IndentedJSON(http.StatusOK, returnedQuestions)

}

func GetPlayers(c *gin.Context) {
	if len(players) == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "No players available"})
		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func SetPlayer(c *gin.Context) {
	var newPlayer Player
	currentPlayerBase := len(players)

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data. Check the provided player details."})
		return
	}

	players = append(players, newPlayer)

	if len(players) != currentPlayerBase+1 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "There was an error adding your player.Please try again later"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func GetPlayer(c *gin.Context) {
	var returnedPlayer Player
	inputId := c.Param("id")

	if reflect.TypeOf(inputId) != reflect.TypeOf(returnedPlayer.Id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	for _, p := range players {
		if p.Id == inputId {
			returnedPlayer = p
			break
		}
	}

	c.IndentedJSON(http.StatusOK, returnedPlayer)
}
