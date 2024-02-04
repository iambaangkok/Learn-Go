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
)

const workDir string = "./src/github.com/iambaangkok/gophercises 01 quiz game/"

func shuffle(data [][]string) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(data); i++ {
		r := random.Intn(i + 1)
		data[i], data[r] = data[r], data[i]
	}
}

func main() {

	randomizeFlag := flag.Bool("rand", false, "randomize quiz order")
	flag.Parse()


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

	score := 0
	problemCount := len(problems)
	invalidFormatProblemCount := 0 

	inputReader := bufio.NewReader(os.Stdin)

	if *randomizeFlag {
		shuffle(problems)
	}

	for i, problem := range problems {
		if len(problem) != 2 {
			fmt.Printf("Problem %v is in invalid format", i)
			invalidFormatProblemCount++
			continue
		}

		question := problem[0]
		solution := problem[1]

		fmt.Printf("#%v %s: ", i, question)
		text, _ := inputReader.ReadString('\n')
		textTrimmed := strings.TrimSpace(text)
		solutionTrimmed := strings.Replace(solution, " ", "", -1)

		fmt.Printf("  sol: %v - (%v) | ans: %v - (%v)\n", solution, len(solution), textTrimmed, len(textTrimmed))
		if strings.Compare(strings.ToLower(solutionTrimmed), strings.ToLower(textTrimmed)) == 0 {
			fmt.Println(">>> correct!")
			score++
		} else {
			fmt.Println(">>> incorrect.")
		}
	}

	fmt.Printf("Total score: %v/%v", score, problemCount-invalidFormatProblemCount)
}