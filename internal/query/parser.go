package query

import (
	"errors"
	"fmt"
	"strings"
)

// tokenization splits sql string into words (tokens) and parsing is the process of converting tokens into a structured query
// QueryType represents the type of SQL query
type QueryType int

const (
	QuerySelect QueryType = iota
	QueryInsert
	// Add more query types as needed
)

// Query represents a parsed SQL query
type Query struct {
	Type   QueryType
	Table  string
	Fields []string
	Values []interface{}
}

// ParseSQL parses a SQL string into a Query struct
func ParseSQL(sql string) (*Query, error) {
	tokens := strings.Fields(sql)
	fmt.Printf("Tokens: %v\n", tokens)

	if len(tokens) < 2 {
		return nil, errors.New("invalid SQL query")
	}

	switch strings.ToUpper(tokens[0]) {
	case "SELECT":
		return parseSelect(tokens)
	case "INSERT":
		return parseInsert(tokens)
	default:
		return nil, errors.New("unsupported query type")
	}
}

func parseSelect(tokens []string) (*Query, error) {
	fmt.Printf("Parsing SELECT tokens: %v\n", tokens)

	if len(tokens) < 4 {
		return nil, fmt.Errorf("invalid SELECT query: too few tokens")
	}

	if strings.ToUpper(tokens[2]) != "FROM" {
		return nil, fmt.Errorf("invalid SELECT query: expected FROM, got %s", tokens[2])
	}

	var fields []string
	if tokens[1] == "*" {
		fields = []string{"*"}
	} else {
		fields = strings.Split(tokens[1], ",")
	}

	return &Query{
		Type:   QuerySelect,
		Fields: fields,
		Table:  tokens[3],
	}, nil
}

func parseInsert(tokens []string) (*Query, error) {
	fmt.Printf("Parsing INSERT tokens: %v\n", tokens)

	if len(tokens) < 4 || strings.ToUpper(tokens[1]) != "INTO" {
		return nil, errors.New("invalid INSERT query")
	}

	return &Query{
		Type:   QueryInsert,
		Table:  tokens[2],
		Values: parseValues(tokens[3:]),
	}, nil
}

func parseValues(tokens []string) []interface{} {
	// TODO: Implement proper value parsing
	return []interface{}{}
}
