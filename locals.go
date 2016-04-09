package main

import (
	"strings"
)

func isLocalWord(w word) bool {
	return strings.HasPrefix(string(w), "!") ||
		strings.HasPrefix(string(w), "$") ||
		w == "get-locals" ||
		w == "drop-locals"
}

func (m *machine) executeLocalsWord(w word) {
	switch w {
	case "get-locals":
		m.locals = append(m.locals, make(map[string]stackEntry))
	case "drop-locals":
		m.locals = m.locals[:len(m.locals)-1]
	}
}
