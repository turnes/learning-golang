package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const FILE = "problems.csv"
const TIME int = 30

type Quiz struct {
	Questions         []question
	NumberOfQuestions int
	Correct           int
	Timer             *time.Timer
}

type question struct {
	question string
	answer   string
}

func main() {
	csvFile := flag.String("file", FILE, "an csv file that has questions,answer structure")
	timer := flag.Int("timer", int(TIME), "time to finish the quiz")
	flag.Parse()
	content := openFile(csvFile)
	quiz := Quiz{
		Timer: time.NewTimer(time.Duration(*timer) * time.Second),
	}
	quiz.loadQuestions(content)
	quiz.start()

}

func openFile(filepath *string) [][]string {
	file, err := os.Open(FILE)
	if err != nil {
		log.Panicln(err)
	}
	defer file.Close()
	csv_reader := csv.NewReader(file)
	lines, err := csv_reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}
	return lines
}

func (q *Quiz) loadQuestions(lines [][]string) {
	q.NumberOfQuestions = len(lines)
	q.Questions = make([]question, len(lines))
	q.Correct = 0
	for i, line := range lines {
		q.Questions[i].question = line[0]
		q.Questions[i].answer = strings.TrimSpace(line[1])
	}
}

func (q *Quiz) start() {
	fmt.Print("Are you ready ?")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
quesionsLoop:
	for _, question := range q.Questions {
		answerCh := make(chan string)
		fmt.Println(question.question)
		fmt.Print("-> ")
		go func() {
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			answerCh <- answer
		}()

		select {
		case <-q.Timer.C:
			fmt.Println()
			break quesionsLoop
		case answer := <-answerCh:
			if answer == question.answer {
				q.Correct++
			}
		}
	}
	fmt.Println(q.Correct, "of", q.NumberOfQuestions)
}
