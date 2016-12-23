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

type TermInfos []TermInfo

// A search term or operator
type queryNode interface {
}

type andQuery struct {
	operands []queryNode
}

func (slice TermInfos) Len() int {
	return len(slice)
}

func (slice TermInfos) Less(i int, j int) bool {
	return slice[i].documentFrequency < slice[j].documentFrequency
}

func (slice TermInfos) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
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

func intersectTermInfos(termInfos TermInfos) TermInfo {
	sort.Sort(termInfos)
	first := termInfos[0]
	rest := termInfos[1:]
	result := first.postings
	for len(result) > 0 && len(rest) > 0 {
		first = rest[0]
		rest = rest[1:]
		result = intersect(first.postings, result)
	}

	return TermInfo{postings: result, documentFrequency: len(result)}
}

func parseQuery(query string) queryNode {
	terms := strings.Split(query, " AND ")
	if len(terms) > 1 {
		query := andQuery{}
		for _, term := range terms {
			query.operands = append(query.operands, term)
		}
		return query
	}

	return query
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

func (index Index) runQuery(query queryNode) TermInfo {
	switch typedQuery := query.(type) {
	case andQuery:
		var termInfos TermInfos
		for _, operand := range typedQuery.operands {
			termInfos = append(termInfos, index.runQuery(operand))
		}
		return intersectTermInfos(termInfos)
	case string:
		return index.terms[typedQuery]
	}

	return TermInfo{}
}

func (index Index) Search(query string) []string {
	parsedQuery := parseQuery(query)
	result := index.runQuery(parsedQuery)
	return result.postings
}

func showQuery(query queryNode, indent int) string {
	var result string
	indentStr := strings.Repeat("  ", indent)

	switch typedQuery := query.(type) {
	case andQuery:
		result += indentStr + "AND (\n"

		for _, operand := range typedQuery.operands {
			result += indentStr + showQuery(operand, indent+1)
		}

		result += indentStr + ")\n"
	case string:
		result += indentStr + fmt.Sprintf("Term \"%s\"\n", typedQuery)
	}

	return result
}

func main() {
	index := New()
	index.Index("doc1", doc1)
	index.Index("doc2", doc2)
	fmt.Println("WELCOME TO GOOGLE 2.0 🔍 ")
	fmt.Println("------------------------\n")
	fmt.Print(index.Show())
	fmt.Println("\nQuery:")
	fmt.Print(showQuery(parseQuery("hello AND world"), 1))
	fmt.Printf("Results: %#v\n", index.Search("julius AND caesar"))
}
