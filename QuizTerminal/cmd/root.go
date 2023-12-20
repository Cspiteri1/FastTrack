/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	str "strconv"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:   "QuizTerminal",
		Short: "This is the terminal to be used to start the quiz, retreive players detail and check the quiz questions",
		Long:  `This command will initialise`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, this is the Quiz Terminal.")
			fmt.Println("This terminal supports the following commands:")
			fmt.Println("startQuiz -- This command starts the process of the quiz")
			fmt.Println("getPlayers -- This command retreive the existing players data")
			fmt.Println("getQuestions -- This command retreives the exising questions and answers provided during the quiz")
		},
	}
	startQuizCmd = &cobra.Command{
		Use:   "startQuiz",
		Short: "Starts the process of allowing the user to participate the quiz",
		Long: `This command allows the user to start the quiz process. Initially the user will be asked for ID no. which will be verified if the user is already existing.
		If the user already exists, the existing user's information will be loaded for later use. If the Id of the user does not match with the current existing IDs, the user
		will be asked to input name and age. After this, the user will be asked to answer a set of question after which, the user will be finally be given a score and rank.
		The user's score will be updated once the quiz is finalised.`,
		Run: func(cmd *cobra.Command, args []string) {
			startQuiz()
		},
	}
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

var players []player
var emptyPlayer player

func createPlayer(id string, name string, age int, score int) (player, bool) {
	apiUrl := "http://localhost:8080/players/"
	var returnedPlayer player

	client := resty.New()
	payload := map[string]interface{}{
		"id":    id,
		"name":  name,
		"age":   age,
		"score": score,
	}

	resp, err := client.R().
		SetBody(payload).
		Post(apiUrl)

	if err != nil {
		log.Fatal("Encountered error while updating player details: ", err)
	}

	if resp.StatusCode() == 201 {
		json.Unmarshal([]byte(resp.Body()), &returnedPlayer)
		return returnedPlayer, true
	} else {
		log.Fatal("Encountered error while updating player details: ", err)
		return returnedPlayer, false
	}
}

func createAnswer(text string, valid bool) answer {
	a := answer{
		AnswerText: text,
		Valid:      valid,
	}
	return a
}

func startQuiz() {
	var idInput string
	var nameInput string
	var ageInput string
	var newPlayer player
	var retry bool = false
	var succesfulResponse bool

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
			convAge, err := str.Atoi(ageInput)
			if err != nil {
				fmt.Println("Your age input was incorrect.")
			} else {
				newPlayer, succesfulResponse = createPlayer(idInput, nameInput, convAge, 0)
				retry = true
			}
		}

		if !succesfulResponse {
			fmt.Println("There is an issue adding your user to the system. Please try again later.")
			return
		}
	}
	fmt.Println("Good luck on your Quiz!")
	fmt.Println()
	questionQuiz(newPlayer)

}

func questionQuiz(currentplayer player) {
	var questions []question
	var answersGroup []answer
	var answerInput answer
	var input int
	var check bool = true
	var score int

	questions = getQuestions()

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
	if updatePlayerScore(currentplayer, score) {
		fmt.Printf("Your score is %v and your rank is %v.", score, getRank(currentplayer.Id))
		fmt.Println()
		fmt.Println("Your answers are : ")
		for i := 0; i < len(answersGroup); i++ {
			fmt.Println(answersGroup[i].AnswerText + "," + str.FormatBool(answersGroup[i].Valid))
		}
	} else {
		fmt.Println("There was an error submitting your Score. Please try again later.")
	}

}

func updatePlayerScore(playerinput player, score int) bool {

	apiUrl := "http://localhost:8080/players/"

	client := resty.New()
	payload := map[string]interface{}{
		"id":    playerinput.Id,
		"name":  playerinput.Name,
		"age":   playerinput.Age,
		"score": score,
	}

	resp, err := client.R().
		SetBody(payload).
		Patch(apiUrl)

	if err != nil {
		log.Fatal("Encountered error while updating player details: ", err)
	}

	if resp.StatusCode() == 200 {
		return true
	} else {
		log.Fatal("Encountered error while updating player details: ", err)
		return false
	}
}

func getRank(playerId string) int {
	apiUrl := "http://localhost:8080/players-rank/" + playerId
	response, err := http.Get(apiUrl)
	var rank int
	if err != nil {
		log.Fatal("An error was encountered during the retreival of the player's ranking: ", err)
		return 0
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Unable to retreive rankings data")
			return 0
		} else {
			json.Unmarshal([]byte(data), &rank)
			return rank
		}
	}
}

func getQuestions() []question {
	var questions []question
	apiUrl := "http://localhost:8080/questions/"
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("An error was encountered during the retreival of the questions: ", err)
		return nil
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Unable to retreive questions data")
			return nil
		} else {
			json.Unmarshal([]byte(data), &questions)
			return questions
		}
	}
}

func checkExistingPlayer(id string) (bool, player) {
	apiUrl := "http://localhost:8080/players/" + id
	var returnedPlayer player
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("Quiz Api is unavailable.Try again later")
		return false, returnedPlayer
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Unable to retreive player data")
			return false, returnedPlayer
		} else {
			json.Unmarshal([]byte(data), &returnedPlayer)
			if returnedPlayer.Id == "" {
				return false, returnedPlayer
			} else {
				return true, returnedPlayer
			}
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(startQuizCmd)
}
