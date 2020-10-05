package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type quizPart struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func main() {
	csvFile, _ := os.Open("questions.csv")
	reader := csv.NewReader(csvFile)
	var questions []quizPart
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		questions = append(questions, quizPart{
			Question: line[0],
			Answer:   line[1],
		})
	}
	var correctAnswerCount int
	var userAnswer string
	const questionCount = 10
	for i := 0; i < questionCount; i++ {
		var answerPos = strconv.Itoa(i + 1)
		rand.Seed(time.Now().UTC().UnixNano())
		var index = rand.Intn(19)
		fmt.Println("Вопрос " + answerPos + "/10:" + questions[index].Question)
		fmt.Fscan(os.Stdin, &userAnswer)
		if userAnswer == questions[index].Answer {
			correctAnswerCount++
		}
	}
	fmt.Println("Количество правильных ответов: " + strconv.Itoa(correctAnswerCount))
}
