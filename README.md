# FastTrackQuiz
This project in Golang is made up of 2 components, a quiz api, and a cli tool which interacts with the api.

The api file 'main.go'  contains a simple api which accepts the following requests:
  1) GET requests which returns all the questions and answers, players information and their rank.
  2) POST request which create a player
  3) PATCH request which updates an existing player

The main functions are then found seperately in the folder named 'handler' in the file 'handler.go'

How to run

1) First download the code and run the api server locally in the terminal. The questions for the api will be generated and stored in memory automatically. 
Run the following: 
```cmd 
cd (path)/FastTrack 
\Fastrack go run main.go
```

2) Build the cli.exe found in the cli folder. Run the following
```cmd
cd (path)FastTrack/QuizTerminal
\quiztest\cli> go build -o (name of application).exe
```

3) Once the cli.exe is created, open command prompt :
```cmd
\path\> (name of cli)exe
```
The CLI accepts the following commands:
startQuiz -- This command starts the process of the quiz
getPlayers -- This command retreive the existing players data
getQuestions -- This command retreives the exising questions and answers provided during the quiz
