/*
Package cmd provides command-line functionality for the QuizTerminal application.
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
	"github.com/spf13/viper"
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
	getPlayersCmd = &cobra.Command{
		Use:   "getPlayers",
		Short: "Returns all information regards the existing players",
		Long:  `Calls the API to return all information regarding the existing players in memory`,
		Run: func(cmd *cobra.Command, args []string) {
			getPlayers()
		},
	}
	getQuestionsCmd = &cobra.Command{
		Use:   "getQuestions",
		Short: "Returns all information regarding the Quiz questions",
		Long:  `Calls the API to return all information regarding the Quiz questions`,
		Run: func(cmd *cobra.Command, args []string) {
			retieveQuestions()
		},
	}
	baseURL string
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

// createPlayer creates a new player and sends a POST request to the API.
func createPlayer(id string, name string, age int, score int) (player, bool) {
	apiUrl := baseURL + "players/"
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

// createAnswer creates a new answer object.
func createAnswer(text string, valid bool) answer {
	a := answer{
		AnswerText: text,
		Valid:      valid,
	}
	return a
}

// startQuiz initiates the quiz process, allowing the user to participate.
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
			log.Println("There is an issue adding your user to the system. Please try again later.")
			return
		}
	}
	fmt.Println("Good luck on your Quiz!")
	fmt.Println()
	questionQuiz(newPlayer)

}

// questionQuiz presents questions to the user and collects their answers.
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
		log.Println("There was an error submitting your Score. Please try again later.")
	}

}

func getPlayers() {
	apiUrl := baseURL + "players"
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Println("An error was encountered during the retreival of the player's ranking: ", err)
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Unable to retreive rankings data")
		} else {
			json.Unmarshal([]byte(data), &players)
			fmt.Println(players)
		}
	}
}

// updatePlayerScore updates the player's score by sending a PATCH request to the API.
func updatePlayerScore(playerinput player, score int) bool {

	apiUrl := baseURL + "players/"

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
		log.Println("Encountered error while updating player details: ", err)
	}

	if resp.StatusCode() == 200 {
		return true
	} else {
		log.Println("Encountered error while updating player details: ", err)
		return false
	}
}

// getRank retrieves the player's rank from the API.
func getRank(playerId string) int {
	apiUrl := baseURL + "players-rank/" + playerId
	response, err := http.Get(apiUrl)
	var rank int
	if err != nil {
		log.Println("An error was encountered during the retreival of the player's ranking: ", err)
		return 0
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Unable to retreive rankings data")
			return 0
		} else {
			json.Unmarshal([]byte(data), &rank)
			return rank
		}
	}
}

func retieveQuestions() {
	retrievedQuestions := getQuestions()

	for i := 0; i < len(retrievedQuestions); i++ {
		fmt.Println(retrievedQuestions[i].QuestionText)
		fmt.Println()
		for b := 0; b < len(retrievedQuestions[i].Answers); b++ {
			fmt.Println(b+1, ".", retrievedQuestions[i].Answers[b].AnswerText)
		}
	}
}

// getQuestions retrieves the list of questions from the API.
func getQuestions() []question {
	var questions []question
	apiUrl := baseURL + "questions/"
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Println("An error was encountered during the retreival of the questions: ", err)
		return nil
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Unable to retreive questions data")
			return nil
		} else {
			json.Unmarshal([]byte(data), &questions)
			return questions
		}
	}
}

// checkExistingPlayer checks if a player with the given ID exists in the API.
func checkExistingPlayer(id string) (bool, player) {
	apiUrl := baseURL + "players/" + id
	var returnedPlayer player
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Println("Quiz Api is unavailable.Try again later")
		return false, returnedPlayer
	} else {
		data, err := io.ReadAll(response.Body)
		if err != nil {
			log.Println("Unable to retreive player data")
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

// main function is the entry point of the application.
func main() {

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
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

// init function initializes the configuration and sets up the API base URL.
func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("api.base_url", "http://localhost:8080/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}

	baseURL = viper.GetString("api.base_url")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(startQuizCmd)
	rootCmd.AddCommand(getQuestionsCmd)
	rootCmd.AddCommand(getPlayersCmd)
}
