package main

import "bufio"
import "fmt"
import "os"
import "github.com/matmoore/search_engine/search"

const doc1 = "I did enact Julius Caesar: I was killed i’ the Capitol; Brutus killed me."

const doc2 = "So let it be with Caesar. The noble Brutus hath told you Caesar was ambitious:"

func main() {
	index := search.NewIndex()
	index.Index("doc1", doc1)
	index.Index("doc2", doc2)

	fmt.Println("WELCOME TO GOOGLE 2.0 🔍 ")
	fmt.Println("------------------------\n")
	fmt.Print(index.Show())

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		query := scanner.Text()
		fmt.Printf("Results: %#v\n", index.Search(query))
		fmt.Print("> ")
	}

	fmt.Println("Bye!")
}
