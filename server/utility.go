package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func jsonResponse(w http.ResponseWriter, body map[string]any, err error) {
	data := map[string]any{
		"body":   body,
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		data["status"] = "error"
		data["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	}

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		panic(err)
	}
}

func getPath(doc map[string]any, parts []string) (any, bool) {
	var docSegment any = doc
	for _, part := range parts {
		m, ok := docSegment.(map[string]any)
		if !ok {
			return nil, false
		}
		if docSegment, ok = m[part]; !ok {
			return nil, false
		}
	}
	return docSegment, true
}

func parseQuery(q string) (*query, error) {
	if q == "" {
		return &query{}, nil
	}
	i := 0
	var parsed query
	var qRune = []rune(q)
	for i < len(qRune) {
		for unicode.IsSpace(qRune[i]) {
			i++
		}
		key, nextIndex, err := lexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("expected valid key, got [%s]: `%q`", err, q[nextIndex])
		}
		if q[nextIndex] != ':' {
			return nil, fmt.Errorf("expected `:`, got [%d]: `%q`", nextIndex, q[nextIndex])
		}
		i = nextIndex + 1

		op := "="
		if q[i] == '>' || q[i] == '<' {
			i++
			op = string(q[i])
		}
		value, nextIndex, err := lexString(qRune, i)
		if err != nil {
			return nil, fmt.Errorf("expected valid value, got [%s]: `%q`", err, q[nextIndex])
		}
		i = nextIndex
		argument := queryComparison{key: strings.Split(key, "."), value: value, op: op}
		parsed.ands = append(parsed.ands, argument)
	}
	return &parsed, nil
}

func lexString(input []rune, index int) (string, int, error) {
	if index >= len(input) {
		return "", index, nil
	}
	if input[index] == '"' {
		index++
		foundEnd := false
		var s []rune
		for index < len(input) {
			if input[index] == '"' {
				foundEnd = true
				break
			}
			s = append(s, input[index])
			index++
		}
		if !foundEnd {
			return "", index, fmt.Errorf("expected end of string")
		}
		return string(s), index + 1, nil
	}

	var s []rune
	var c rune
	for index < len(input) {
		c = input[index]
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '.') {
			break
		}
		s = append(s, c)
		index++
	}
	if len(s) == 0 {
		return "", index, fmt.Errorf("no string found")
	}
	return string(s), index, nil
}

func parseCommandLine(q []string) (string, string) {
	q = q[1:]

	if len(q) < 1 {
		return "data.doc", "8080"
	} else if len(q) < 2 {
		return q[0], "8080"
	} else {
		return q[0], q[1]
	}
}
