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

func createPlayer(name string, age int, score int) player {
	p := player{
		Name:  name,
		Age:   age,
		Score: score,
	}

	return p
}

var players []player

func main() {
	players = append(players, createPlayer("Clayton", 28, 60))
	players = append(players, createPlayer("Giancarl", 27, 40))
	players = append(players, createPlayer("Emma", 33, 100))
	players = append(players, createPlayer("Paulinha", 30, 80))

	router := gin.Default()
	router.GET("/players", getPlayers)
	router.POST("/players", setPlayer)
	router.GET("/books/:name", getPlayer)
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

	if reflect.TypeOf(inputName) == reflect.TypeOf(returnedPlayer.Name) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "The value you provided is not valid"})
		return
	}

	for i := 0; i < len(players); i++ {
		if players[i].Name == inputName {
			returnedPlayer = players[i]
		}
	}

	c.IndentedJSON(http.StatusOK, returnedPlayer)
}
