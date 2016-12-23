package parsing

import "fmt"
import "strings"

// A search term or operator
type QueryNode interface {
}

type AndQuery struct {
	Operands []QueryNode
}

func ParseQuery(query string) QueryNode {
	terms := strings.Split(query, " AND ")
	if len(terms) > 1 {
		query := AndQuery{}
		for _, term := range terms {
			query.Operands = append(query.Operands, term)
		}
		return query
	}

	return query
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
	case string:
		result += indentStr + fmt.Sprintf("Term \"%s\"\n", typedQuery)
	}

	return result
}
