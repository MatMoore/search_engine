package search

import "github.com/matmoore/search_engine/parsing"
import "sort"

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

func (index Index) runQuery(query parsing.QueryNode) TermInfo {
	switch typedQuery := query.(type) {
	case parsing.AndQuery:
		var termInfos TermInfos
		for _, operand := range typedQuery.Operands {
			termInfos = append(termInfos, index.runQuery(operand))
		}
		return intersectTermInfos(termInfos)
	case string:
		return index.terms[typedQuery]
	}

	return TermInfo{}
}

func (index Index) Search(query string) []string {
	parsedQuery := parsing.ParseQuery(query)
	result := index.runQuery(parsedQuery)
	return result.postings
}
