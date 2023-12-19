package main

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	str "strconv"

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
	startQuiz()
	router := gin.Default()
	router.GET("/players", getPlayers)
	router.POST("/players", setPlayer)
	router.GET("/players/:name", getPlayer)
	router.Run("localhost:8080")
}

func startQuiz() {
	var idInput string
	var nameInput string
	var ageInput string
	var newPlayer player
	var retry bool = false

	fmt.Println("Welcome to the Quiz, Kindly enter your ID")
	fmt.Scan(&idInput)
	check, newPlayer := checkExistingPlayer(idInput)
	if check {
		fmt.Printf("Welcome back %v .", newPlayer.Name)
		fmt.Println()
	}
	if !check {
		fmt.Println("Kindly enter your name")
		fmt.Scan(&nameInput)
		for !retry {
			fmt.Println("Kindly enter your age")
			fmt.Scan(&ageInput)
			convertedAge, err := str.Atoi(ageInput)
			if err != nil {
				fmt.Println("Your age input was incorrect.")
			} else {
				newPlayer = createPlayer(idInput, nameInput, convertedAge, 0)
				players = append(players, newPlayer)
				retry = true
			}
		}

	}
	fmt.Println("Good luck on your Quiz!")
	fmt.Println()
	questionQuiz(newPlayer)
}

func questionQuiz(currentplayer player) {
	var answersGroup []answer
	var answerInput answer
	var input int
	var check bool = true
	var score int
	fmt.Println("Kindly select one answer for each question.")

	for i := 0; i < len(questions); i++ {
		check = true
		fmt.Println(questions[i].QuestionText)
		fmt.Println()
		for b := 0; b < len(questions[i].Answers); b++ {
			fmt.Println(b+1, ".", questions[i].Answers[b].AnswerText)
		}
		for check {
			fmt.Scan(&input)
			if input > len(questions[i].Answers) || input < 0 {
				fmt.Println("Invalid Answer. Try again")
			} else {
				answerInput = createAnswer(questions[i].Answers[input-1].AnswerText, questions[i].Answers[input-1].Valid)
				answersGroup = append(answersGroup, answerInput)
				if questions[i].Answers[input-1].Valid {
					score += (100 / len(questions))
				}
				check = false
			}
		}

		fmt.Println()
	}

	fmt.Println("End of Quiz.")
	fmt.Printf("Your score is %v and your rank is %v.", score, generateRank(score, currentplayer))
	fmt.Println()
	fmt.Println("Your answers are : ")
	for i := 0; i < len(answersGroup); i++ {
		fmt.Println(answersGroup[i].AnswerText + "," + str.FormatBool(answersGroup[i].Valid))
	}

	idx := sort.Search(len(players), func(i int) bool {
		return string(players[i].Name) >= currentplayer.Name
	})
	players[idx+1].Score = score
}

func generateRank(score int, currentplayer player) int {
	var rank int = len(players)

	for i := 0; i < len(players); i++ {
		if score > players[i].Score {
			if players[i].Id != currentplayer.Id {
				rank -= 1
			}
		}
	}
	return rank
}

func checkExistingPlayer(id string) (bool, player) {
	for i := 0; i < len(players); i++ {
		if players[i].Id == id {
			return true, players[i]
		}
	}
	return false, emptyPlayer
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
