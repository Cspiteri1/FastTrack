package main

import (
	"net/http"
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
	_ "github.com/spf13/cobra"
)

type player struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Score int    `json:"score"`
}

type question struct {
	QuestionText string
	Answers      []answer
}

type answer struct {
	AnswerText string
	Valid      bool
}

func createPlayer(id string, name string, age int, score int) player {
	p := player{
		Id:    id,
		Name:  name,
		Age:   age,
		Score: score,
	}

	return p
}

func createQuestion(text string, answer []answer) question {
	q := question{
		QuestionText: text,
		Answers:      answer,
	}
	return q
}

func createAnswer(text string, valid bool) answer {
	a := answer{
		AnswerText: text,
		Valid:      valid,
	}
	return a
}

func initData() {

	players = append(players, createPlayer("000001M", "Clayton", 28, 30))
	players = append(players, createPlayer("000002M", "Giancarl", 27, 10))
	players = append(players, createPlayer("000003M", "Emma", 33, 90))
	players = append(players, createPlayer("000004M", "Paulinha", 30, 50))

	var answersGroupA []answer = []answer{createAnswer("Italian", false), createAnswer("Arabic", false), createAnswer("Maltese", true)}
	var answersGroupB []answer = []answer{createAnswer("Lira", false), createAnswer("Pound", false), createAnswer("Euro", true)}
	var answersGroupC []answer = []answer{createAnswer("Mediterranean ", true), createAnswer("Dead Sea", false), createAnswer("Sea of Samsara", false)}
	var answersGroupD []answer = []answer{createAnswer("three", true), createAnswer("one", false), createAnswer("four", false)}
	var answersGroupE []answer = []answer{createAnswer("Blue,White and Red", false), createAnswer("White and Red", true), createAnswer("Red and White", false)}

	questions = append(questions, createQuestion("What ls the national language of Malta?", answersGroupA))
	questions = append(questions, createQuestion("What currency is used in Malta?", answersGroupB))
	questions = append(questions, createQuestion("In which sea is Malta located?", answersGroupC))
	questions = append(questions, createQuestion("How many inhabited islands make up the Republic of Malta?", answersGroupD))
	questions = append(questions, createQuestion("What are the colours of Malta's National flag?", answersGroupE))

}

var players []player
var questions []question
var emptyPlayer player

func main() {
	initData()
	router := gin.Default()
	router.GET("/players", getPlayers)
	router.GET("/players/:id", getPlayer)
	router.GET("/questions/", getQuestions)
	router.GET("/players-rank/:id", getPlayerRank)
	router.POST("/players", setPlayer)
	router.PATCH("/players/", updatePlayer)
	router.Run("localhost:8080")
}

func updatePlayer(c *gin.Context) {
	var newPlayer player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
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

func getPlayerRank(c *gin.Context) {
	var returnedPlayer player
	inputId := c.Param("id")
	inputType := reflect.TypeOf(inputId)
	var rank int = len(players)

	if reflect.TypeOf(inputType) == reflect.TypeOf(returnedPlayer.Id) {
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

func getQuestions(c *gin.Context) {
	returnedQuestions := questions

	if len(questions) == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "There was an issue returning the players"})
	} else {
		c.IndentedJSON(http.StatusOK, returnedQuestions)
	}
}

func getPlayers(c *gin.Context) {
	returnedPlayers := players

	if len(returnedPlayers) <= 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "There was an issue returning the players"})
		return
	}

	c.IndentedJSON(http.StatusOK, returnedPlayers)
}

func setPlayer(c *gin.Context) {
	var newPlayer player
	currentPlayerBase := len(players)

	if err := c.BindJSON(&newPlayer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	players = append(players, newPlayer)

	if len(players) != currentPlayerBase+1 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "There was an error adding your player.Please try again later"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newPlayer)
}

func getPlayer(c *gin.Context) {
	var returnedPlayer player
	inputId := c.Param("id")
	playertype := reflect.TypeOf(inputId)

	if reflect.TypeOf(playertype) == reflect.TypeOf(returnedPlayer.Id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	for i := 0; i < len(players); i++ {
		if players[i].Id == inputId {
			returnedPlayer = players[i]
			break
		}
	}

	c.IndentedJSON(http.StatusOK, returnedPlayer)
}
