package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Term       string
	Definition string
	Mistakes   int
}

var (
	cards         []Card
	termsMap      = make(map[string]bool) // track terms
	definitionMap = make(map[string]bool) // track definitions
	reader        = bufio.NewReader(os.Stdin)
	logs          []string // store all input/output lines
)

func main() {
	rand.Seed(time.Now().UnixNano())

	for {
		printAndLog("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
		action := readInput()

		switch action {
		case "add":
			addCard()
		case "remove":
			removeCard()
		case "import":
			importCards()
		case "export":
			exportCards()
		case "ask":
			askQuestions()
		case "log":
			saveLog()
		case "hardest card":
			hardestCard()
		case "reset stats":
			resetStats()
		case "exit":
			printAndLog("Bye bye!")
			return
		default:
			printAndLog("Unknown action.")
		}
		printAndLog("")
	}
}

func printAndLog(s string) {
	fmt.Println(s)
	logs = append(logs, s)
}

func readInput() string {
	str, _ := reader.ReadString('\n')
	str = strings.TrimSpace(str)
	logs = append(logs, "> "+str)
	return str
}

func addCard() {
	var term string
	for {
		printAndLog("The card:")
		term = readInput()
		if termsMap[term] {
			printAndLog("This term already exists. Try again:")
			continue
		}
		break
	}

	var def string
	for {
		printAndLog("The definition of the card:")
		def = readInput()
		if definitionMap[def] {
			printAndLog("This definition already exists. Try again:")
			continue
		}
		break
	}

	termsMap[term] = true
	definitionMap[def] = true
	cards = append(cards, Card{Term: term, Definition: def, Mistakes: 0})
	printAndLog(fmt.Sprintf("The pair (\"%s\":\"%s\") has been added.", term, def))
}

func removeCard() {
	printAndLog("Which card?")
	term := readInput()
	if termsMap[term] {
		termsMap[term] = false
		defToRemove := ""
		newCards := []Card{}
		for _, c := range cards {
			if c.Term == term {
				defToRemove = c.Definition
				continue
			}
			newCards = append(newCards, c)
		}
		cards = newCards
		if defToRemove != "" {
			definitionMap[defToRemove] = false
		}
		printAndLog("The card has been removed.")
	} else {
		printAndLog(fmt.Sprintf("Can't remove \"%s\": there is no such card.", term))
	}
}

func importCards() {
	printAndLog("File name:")
	fileName := readInput()
	file, err := os.Open(fileName)
	if err != nil {
		printAndLog("File not found.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 2 {
			continue
		}
		term, def := parts[0], parts[1]
		mistakes := 0
		if len(parts) == 3 {
			mistakes, _ = strconv.Atoi(parts[2])
		}

		if termsMap[term] {
			for i, c := range cards {
				if c.Term == term {
					definitionMap[c.Definition] = false
					cards[i].Definition = def
					cards[i].Mistakes = mistakes
					definitionMap[def] = true
					break
				}
			}
		} else if definitionMap[def] {
			printAndLog(fmt.Sprintf("This definition already exists. Skipping \"%s\".", term))
			continue
		} else {
			termsMap[term] = true
			definitionMap[def] = true
			cards = append(cards, Card{Term: term, Definition: def, Mistakes: mistakes})
		}
		count++
	}
	printAndLog(fmt.Sprintf("%d cards have been loaded.", count))
}

func exportCards() {
	printAndLog("File name:")
	fileName := readInput()
	file, err := os.Create(fileName)
	if err != nil {
		printAndLog("Error creating file.")
		return
	}
	defer file.Close()

	for _, card := range cards {
		fmt.Fprintf(file, "%s:%s:%d\n", card.Term, card.Definition, card.Mistakes)
	}
	printAndLog(fmt.Sprintf("%d cards have been saved.", len(cards)))
}

func askQuestions() {
	if len(cards) == 0 {
		printAndLog("No cards available.")
		return
	}

	printAndLog("How many times to ask?")
	nStr := readInput()
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		printAndLog("Invalid number.")
		return
	}

	for i := 0; i < n; i++ {
		card := cards[rand.Intn(len(cards))]
		printAndLog(fmt.Sprintf("Print the definition of \"%s\":", card.Term))
		answer := readInput()

		if answer == card.Definition {
			printAndLog("Correct!")
		} else {
			card.Mistakes++
			for i := range cards {
				if cards[i].Term == card.Term {
					cards[i].Mistakes = card.Mistakes
					break
				}
			}
			otherTerm := ""
			for _, c := range cards {
				if c.Definition == answer {
					otherTerm = c.Term
					break
				}
			}
			if otherTerm != "" {
				printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\", but your definition is correct for \"%s\" card.", card.Definition, otherTerm))
			} else {
				printAndLog(fmt.Sprintf("Wrong. The right answer is \"%s\".", card.Definition))
			}
		}
	}
}

func saveLog() {
	printAndLog("File name:")
	fileName := readInput()
	file, err := os.Create(fileName)
	if err != nil {
		printAndLog("Error creating file.")
		return
	}
	defer file.Close()

	for _, line := range logs {
		fmt.Fprintln(file, line)
	}
	printAndLog("The log has been saved.")
}

func hardestCard() {
	maxMistakes := 0
	var hardest []string
	for _, c := range cards {
		if c.Mistakes > maxMistakes {
			maxMistakes = c.Mistakes
			hardest = []string{c.Term}
		} else if c.Mistakes == maxMistakes && maxMistakes > 0 {
			hardest = append(hardest, c.Term)
		}
	}
	if maxMistakes == 0 {
		printAndLog("There are no cards with errors.")
		return
	}
	if len(hardest) == 1 {
		printAndLog(fmt.Sprintf("The hardest card is \"%s\". You have %d errors answering it.", hardest[0], maxMistakes))
	} else {
		printAndLog(fmt.Sprintf("The hardest cards are \"%s\". You have %d errors answering them.", strings.Join(hardest, "\", \""), maxMistakes))
	}
}

func resetStats() {
	for i := range cards {
		cards[i].Mistakes = 0
	}
	printAndLog("Card statistics have been reset.")
}
