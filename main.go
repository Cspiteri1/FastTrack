package main

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type player struct {
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

func createPlayer(name string, age int, score int) player {
	p := player{
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

	players = append(players, createPlayer("Clayton", 28, 60))
	players = append(players, createPlayer("Giancarl", 27, 40))
	players = append(players, createPlayer("Emma", 33, 100))
	players = append(players, createPlayer("Paulinha", 30, 80))

	var answersGroupA []answer = []answer{createAnswer("Italian", false), createAnswer("Arabic", false), createAnswer("Maltese", true)}
	var answersGroupB []answer = []answer{createAnswer("Lira", false), createAnswer("Pound", false), createAnswer("Euro", true)}
	var answersGroupC []answer = []answer{createAnswer("Mediterranean ", true), createAnswer("Dead Sea", false), createAnswer("Sea of Samsara", false)}
	var answersGroupD []answer = []answer{createAnswer("3", true), createAnswer("1", false), createAnswer("4", false)}
	var answersGroupE []answer = []answer{createAnswer("Blue,White and Red", false), createAnswer("White and Red", true), createAnswer("Red and White", false)}

	questions = append(questions, createQuestion("What ls the national language of Malta?", answersGroupA))
	questions = append(questions, createQuestion("What currency is used in Malta?", answersGroupB))
	questions = append(questions, createQuestion("In which sea is Malta located?", answersGroupC))
	questions = append(questions, createQuestion("How many inhabited islands make up the Republic of Malta?", answersGroupD))
	questions = append(questions, createQuestion("What are the colours of Malta's National flag?", answersGroupE))

}

var players []player
var questions []question

func main() {

	initData()
	router := gin.Default()
	router.GET("/players", getPlayers)
	router.POST("/players", setPlayer)
	router.GET("/players/:name", getPlayer)
	router.Run("localhost:8080")
}

func getPlayers(c *gin.Context) {
	returnedPlayers := players

	if len(returnedPlayers) < 0 {
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

	c.IndentedJSON(http.StatusCreated, players)
}

func getPlayer(c *gin.Context) {
	var returnedPlayer player
	inputName := c.Param("name")
	playertype := reflect.TypeOf(inputName)

	if reflect.TypeOf(playertype) == reflect.TypeOf(returnedPlayer.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	for i := 0; i < len(players); i++ {
		if players[i].Name == inputName {
			returnedPlayer = players[i]
			break
		}
	}

	c.IndentedJSON(http.StatusOK, returnedPlayer)
}
