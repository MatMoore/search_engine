package search

import "github.com/matmoore/search_engine/parsing"
import "sort"

func intersect(termInfo1 TermInfo, termInfo2 TermInfo) TermInfo {
	var result TermInfo
	posting1 := termInfo1.postings
	posting2 := termInfo2.postings

	for len(posting1) > 0 && len(posting2) > 0 {
		if posting1[0] == posting2[0] {
			result.postings = append(result.postings, posting1[0])
			posting1 = posting1[1:]
			posting2 = posting2[1:]
		} else if posting1[0] < posting2[0] {
			posting1 = posting1[1:]
		} else {
			posting2 = posting2[1:]
		}
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

func difference(include TermInfo, exclude TermInfo) TermInfo {
	var result TermInfo
	includePostings := include.postings
	excludePostings := exclude.postings

	for len(includePostings) > 0 && len(excludePostings) > 0 {
		if includePostings[0] < excludePostings[0] {
			result.postings = append(result.postings, includePostings[0])
			includePostings = includePostings[1:]
		} else if includePostings[0] == excludePostings[0] {
			includePostings = includePostings[1:]
			excludePostings = excludePostings[1:]
		} else {
			excludePostings = excludePostings[1:]
		}
	}

	for _, docId := range includePostings {
		result.postings = append(result.postings, docId)
		result.documentFrequency += 1
	}

	return result
}

func unionTermInfos(termInfos TermInfos) TermInfo {
	if len(termInfos) == 0 {
		return TermInfo{}
	}

	first := termInfos[0]
	rest := termInfos[1:]
	result := first
	for len(rest) > 0 {
		first = rest[0]
		rest = rest[1:]
		result = union(first, result)
	}

	return result
}

func intersectTermInfos(termInfos TermInfos) TermInfo {
	if len(termInfos) == 0 {
		return TermInfo{}
	}

	// Process in increasing order of size to fetch
	// the minimal number of documents.
	sort.Sort(termInfos)

	first := termInfos[0]
	rest := termInfos[1:]
	result := first

	// Short circuit if a subquery has zero results
	for result.documentFrequency > 0 && len(rest) > 0 {
		first = rest[0]
		rest = rest[1:]
		result = intersect(first, result)
	}

	return result
}

// TODO: this is probably not the best order of operations. I think we should
// calculate the number of documents in each operand, regardless of whether
// it's negated, and combine operands from smallest to largest, using intersect
// or difference as appropriate.
func intersectAndExcludeTermInfos(included TermInfos, excluded TermInfos) TermInfo {
	allIncluded := intersectTermInfos(included)
	allExcluded := intersectTermInfos(excluded)
	return difference(allIncluded, allExcluded)
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
		var includedTermInfos TermInfos
		var excludedTermInfos TermInfos

		// Process `foo AND NOT bar` subqueries in one step, to avoid fetching
		// *everything* just to calculate an intermediate result
		for _, operand := range typedQuery.Operands {
			switch typedSubQuery := operand.(type) {
			case parsing.NotQuery:
				excludedTermInfos = append(excludedTermInfos, index.runQuery(typedSubQuery.Operand))
			case string:
				includedTermInfos = append(includedTermInfos, index.runQuery(operand))
			}
		}

		return intersectAndExcludeTermInfos(includedTermInfos, excludedTermInfos)
	case parsing.OrQuery:
		var termInfos TermInfos
		for _, operand := range typedQuery.Operands {
			termInfos = append(termInfos, index.runQuery(operand))
		}
		return unionTermInfos(termInfos)

	case parsing.NotQuery:
		return difference(index.all(), index.runQuery(typedQuery.Operand))

	case string:
		return index.terms[typedQuery]
	}

	return TermInfo{}
}

func (index Index) all() TermInfo {
	var termInfos TermInfos
	for _, termInfo := range index.terms {
		termInfos = append(termInfos, termInfo)
	}

	return unionTermInfos(termInfos)
}

func (index Index) Search(query string) []string {
	parsedQuery := parsing.ParseQuery(query)
	result := index.runQuery(parsedQuery)
	return result.postings
}
