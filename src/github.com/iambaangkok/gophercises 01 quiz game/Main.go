package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	// "sync"
)

const workDir string = "./src/github.com/iambaangkok/gophercises 01 quiz game/"

func shuffle(data [][]string) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(data); i++ {
		r := random.Intn(i + 1)
		data[i], data[r] = data[r], data[i]
	}
}

// var wg sync.WaitGroup

var currentQuestion int = 0

func countDown(ch chan struct{}, seconds int, questionId int) {
	time.Sleep(time.Duration(seconds) * time.Second)
	if (currentQuestion == questionId) {
		ch <- struct{}{}
	}
}

func readString(ch chan string, inputReader *bufio.Reader) {
	text, _ := inputReader.ReadString('\n')
	ch <- text
}

func main() {

	// Channels
	readStringCh := make(chan string)
	timeLimitCh := make(chan struct{}) 

	// Declare flags
	timeLimitPtr := flag.Int("time-limit", 30, "time limit per question (seconds)")
	randomizePtr := flag.Bool("rand", false, "randomize quiz order")
	flag.Parse()

	// Read file
	const problemsFileName = "problems.csv"
    file, err := os.Open(workDir + problemsFileName)
    if err != nil {
		log.Fatal("Error while reading the file")
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	problems, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error while reading records")
	}

	if *randomizePtr {
		shuffle(problems)
	}

	// Press Enter to start
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("QUIZ GAME: Press Enter to start! You have %v seconds for each question.", *timeLimitPtr)
	inputReader.ReadString('\n')
		

	score := 0
	problemCount := len(problems)
	invalidFormatProblemCount := 0 

	for i, problem := range problems {
		// wg.Add(1)
		currentQuestion = i

		if len(problem) != 2 {
			fmt.Printf("Problem %v is in invalid format", i)
			invalidFormatProblemCount++
			continue
		}

		question := problem[0]
		solution := problem[1]

		fmt.Printf("#%v %s: ", i, question)
		go countDown(timeLimitCh, *timeLimitPtr, i)
		go readString(readStringCh, inputReader)

		text := ""
		timedOut := false
		select {
			case t := <- readStringCh:
				text = t
			case _ = <- timeLimitCh:
				timedOut = true
		}

		textTrimmed := strings.TrimSpace(text)
		solutionTrimmed := strings.Replace(solution, " ", "", -1)

		fmt.Printf("  sol: %v - (%v) | ans: %v - (%v)\n", solution, len(solution), textTrimmed, len(textTrimmed))
		if timedOut {
			fmt.Println(">>> timed out.")
		} else if strings.Compare(strings.ToLower(solutionTrimmed), strings.ToLower(textTrimmed)) == 0 {
			fmt.Println(">>> correct!")
			score++
		} else {
			fmt.Println(">>> incorrect.")
		}
	}

	fmt.Printf("Total score: %v/%v", score, problemCount-invalidFormatProblemCount)
}