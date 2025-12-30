package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	Term       string
	Definition string
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Input the number of cards:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	numCards, _ := strconv.Atoi(input)

	// Kartochkalarni slice ichida saqlaymiz
	var cards []Card

	// Term va definitionlarni kiritish
	for i := 1; i <= numCards; i++ {
		fmt.Printf("The term for card #%d:\n", i)
		term, _ := reader.ReadString('\n')
		term = strings.TrimSpace(term)

		fmt.Printf("The definition for card #%d:\n", i)
		definition, _ := reader.ReadString('\n')
		definition = strings.TrimSpace(definition)

		cards = append(cards, Card{Term: term, Definition: definition})
	}

	// Foydalanuvchidan javob so‘rash va tekshirish (kiritilgan tartibda)
	for _, card := range cards {
		fmt.Printf("Print the definition of \"%s\":\n", card.Term)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		if answer == card.Definition {
			fmt.Println("Correct!")
		} else {
			fmt.Printf("Wrong. The right answer is \"%s\".\n", card.Definition)
		}
	}
}
