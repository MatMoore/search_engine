package main

import "fmt"
import "strings"
import "sort"

const doc1 = "I did enact Julius Caesar: I was killed i’ the Capitol; Brutus killed me."

const doc2 = "So let it be with Caesar. The noble Brutus hath told you Caesar was ambitious:"

type Index struct {
	terms map[string][]string
}

func tokenise(text string) []string {
	return strings.Fields(text)
}

func uniq(docs []string) []string {
	result := make([]string, 0, len(docs))
	var last string
	for _, doc := range docs {
		if doc != last {
			result = append(result, doc)
		}
		last = doc
	}

	return result
}

func (index Index) insertPosting(term string, docId string) {
	existing := index.terms[term]
	withDoc := append(existing, docId)
	sort.Strings(withDoc)
	index.terms[term] = uniq(withDoc)
}

func (index Index) Index(docId string, text string) {
	tokens := tokenise(text)
	for _, term := range tokens {
		index.insertPosting(term, docId)
	}
}

func New() Index {
	return Index{terms: make(map[string][]string)}
}

func (index Index) Show() string {
	var result string
	for key, value := range index.terms {
		result += fmt.Sprintf("%s -> %s\n", key, strings.Join(value, ", "))
	}
	return result
}

func main() {
	index := New()
	index.Index("doc1", doc1)
	index.Index("doc2", doc2)
	fmt.Println("WELCOME TO GOOGLE 2.0 🔍 ")
	fmt.Println("------------------------")
	fmt.Print(index.Show())
}
