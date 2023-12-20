package main

// main.go
// This file contains the main entry point for the FastTrack application.
// It initializes data, sets up the Gin router, and defines API endpoints.
// Author: Clayton Spiteri

import (
	"FastTrack/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	handler.InitData()                                     // Initialize sample data for players and questions.
	router := gin.Default()                                // Set up Gin router and define API endpoints.
	router.GET("/players", handler.GetPlayers)             // Retrieve information about all players.
	router.GET("/players/:id", handler.GetPlayer)          // Retrieve information about a specific player based on ID.
	router.GET("/questions", handler.GetQuestions)         // Retrieve a list of quiz questions.
	router.GET("/players-rank/:id", handler.GetPlayerRank) // Retrieve the rank of a player based on their score.
	router.POST("/players", handler.SetPlayer)             // Add a new player.
	router.PATCH("/players", handler.UpdatePlayer)         // Update an existing player.
	router.Run("localhost:8080")
}
