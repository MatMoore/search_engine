package parsing

import "fmt"
import "strings"

// A search term or operator
type QueryNode interface {
}

type AndQuery struct {
	Operands []QueryNode
}

type OrQuery struct {
	Operands []QueryNode
}

func parseDisjunction(subquery string) QueryNode {
	terms := strings.Split(subquery, " OR ")
	if len(terms) > 1 {
		node := OrQuery{}
		for _, term := range terms {
			node.Operands = append(node.Operands, term)
		}
		return node
	}

	return subquery
}

// Queries are assumed to be written in disjunctive normal form (ANDs of ORs)
func ParseQuery(query string) QueryNode {
	terms := strings.Split(query, " AND ")
	if len(terms) > 1 {
		node := AndQuery{}
		for _, term := range terms {
			node.Operands = append(node.Operands, parseDisjunction(term))
		}
		return node
	}

	return parseDisjunction(query)
}

func showQuery(query QueryNode, indent int) string {
	var result string
	indentStr := strings.Repeat("  ", indent)

	switch typedQuery := query.(type) {
	case AndQuery:
		result += indentStr + "AND (\n"

		for _, operand := range typedQuery.Operands {
			result += indentStr + showQuery(operand, indent+1)
		}

		result += indentStr + ")\n"
	case OrQuery:
		result += indentStr + "OR (\n"
		for _, operand := range typedQuery.Operands {
			result += indentStr + showQuery(operand, indent+1)
		}

		result += indentStr + ")\n"
	case string:
		result += indentStr + fmt.Sprintf("Term \"%s\"\n", typedQuery)
	}

	return result
}
