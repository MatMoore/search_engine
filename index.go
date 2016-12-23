package main

import "fmt"
import "strings"
import "sort"

const doc1 = "I did enact Julius Caesar: I was killed i’ the Capitol; Brutus killed me."

const doc2 = "So let it be with Caesar. The noble Brutus hath told you Caesar was ambitious:"

type Index struct {
	terms map[string]TermInfo
}

type TermInfo struct {
	documentFrequency int
	postings          []string
}

func tokenise(text string) []string {
	replacer := strings.NewReplacer(".", "", ",", "", ":", "", ";", "")
	return strings.Fields(replacer.Replace(strings.ToLower(text)))
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

func intersect(posting1 []string, posting2 []string) []string {
	var answer []string
	for len(posting1) > 0 && len(posting2) > 0 {
		if posting1[0] == posting2[0] {
			answer = append(answer, posting1[0])
			posting1 = posting1[1:]
			posting2 = posting2[1:]
		} else if posting1[0] < posting2[0] {
			posting1 = posting1[1:]
		} else {
			posting2 = posting2[1:]
		}
	}

	return answer
}

func (index Index) insertPosting(term string, docId string) {
	termInfo, ok := index.terms[term]
	if !ok {
		termInfo = TermInfo{}
	}

	withDoc := append(termInfo.postings, docId)
	sort.Strings(withDoc)
	termInfo.postings = uniq(withDoc)
	termInfo.documentFrequency = len(termInfo.postings)

	index.terms[term] = termInfo
}

func (index Index) Index(docId string, text string) {
	tokens := tokenise(text)
	for _, term := range tokens {
		index.insertPosting(term, docId)
	}
}

func New() Index {
	return Index{terms: make(map[string]TermInfo)}
}

func (index Index) Show() string {
	var result string
	for term, termInfo := range index.terms {
		result += fmt.Sprintf(
			"%s (%d) -> %s\n",
			term,
			termInfo.documentFrequency,
			strings.Join(termInfo.postings, ", "),
		)
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
