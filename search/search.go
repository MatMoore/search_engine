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

func mergeTermInfos(binaryOperation func(TermInfo, TermInfo) TermInfo, termInfos TermInfos) TermInfo {
	first := termInfos[0]
	rest := termInfos[1:]
	result := first
	for len(rest) > 0 {
		first = rest[0]
		rest = rest[1:]
		result = binaryOperation(first, result)
	}

	return result
}

func union(termInfo1 TermInfo, termInfo2 TermInfo) TermInfo {
	var result TermInfo
	posting1 := termInfo1.postings
	posting2 := termInfo2.postings

	for len(posting1) > 0 && len(posting2) > 0 {
		var next string

		if posting1[0] < posting2[0] {
			next = posting1[0]
			posting1 = posting1[1:]
		} else if posting1[0] == posting2[0] {
			next = posting1[0]
			posting1 = posting1[1:]
			posting2 = posting2[1:]
		} else {
			next = posting2[0]
			posting2 = posting1[1:]
		}

		result.postings = append(result.postings, next)
		result.documentFrequency += 1
	}

	result.postings = append(result.postings, posting1...)
	result.postings = append(result.postings, posting2...)

	return result
}

// TODO: it's inefficient to process the OR subqueries first, because we fetch
// the maximum number of documents right away, even if they would be filtered out
// by the ANDs.

// Assuming that the collections returned by the subqueries are very different
// in size, we should probably convert to ORs of ANDs instead (See Manning
// Exercise 1.6 on disjunctive normal form). This way we fetch less
// non-matching documents when the size of the subqueries varies, because the
// number of documents examined by the intersect algorithm is bound by the
// smallest of the two sets.
func (index Index) runQuery(query parsing.QueryNode) TermInfo {
	switch typedQuery := query.(type) {
	case parsing.AndQuery:
		var termInfos TermInfos
		for _, operand := range typedQuery.Operands {
			termInfos = append(termInfos, index.runQuery(operand))
		}
		return intersectTermInfos(termInfos)
	case parsing.OrQuery:
		var termInfos TermInfos
		for _, operand := range typedQuery.Operands {
			termInfos = append(termInfos, index.runQuery(operand))
		}
		return mergeTermInfos(union, termInfos)
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
