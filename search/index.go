package search

import "fmt"
import "strings"
import "sort"

type Index struct {
	terms map[string]TermInfo
}

type TermInfo struct {
	documentFrequency int
	postings          []string
}

type TermInfos []TermInfo

func (slice TermInfos) Len() int {
	return len(slice)
}

func (slice TermInfos) Less(i int, j int) bool {
	return slice[i].documentFrequency < slice[j].documentFrequency
}

func (slice TermInfos) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
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

func tokenise(text string) []string {
	replacer := strings.NewReplacer(".", "", ",", "", ":", "", ";", "")
	return strings.Fields(replacer.Replace(strings.ToLower(text)))
}

func (index Index) Index(docId string, text string) {
	tokens := tokenise(text)
	for _, term := range tokens {
		index.insertPosting(term, docId)
	}
}

func NewIndex() Index {
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
