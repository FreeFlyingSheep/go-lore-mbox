package mbox

import "strings"

type style int

const (
	undefined style = iota
	text
	start
	end
	before
	after
	change
	diff
	index
	add
	del
	quote
)

var styles = map[style]string{
	undefined: "undefined",
	text:      "text",
	start:     "git-start",
	end:       "git-end",
	before:    "git-before",
	after:     "git-after",
	change:    "git-change",
	diff:      "git-diff",
	index:     "git-index",
	add:       "git-add",
	del:       "git-del",
	quote:     "quote",
}

func parseStyle(line string) style {
	if line == "---" {
		return start
	} else if line == "--" || line == "-- " {
		return end
	} else if strings.HasPrefix(line, "--- ") {
		return before
	} else if strings.HasPrefix(line, "+++ ") {
		return after
	} else if strings.HasPrefix(line, "@@ ") {
		return change
	} else if strings.HasPrefix(line, "diff ") {
		return diff
	} else if strings.HasPrefix(line, "index ") {
		return index
	} else if strings.HasPrefix(line, ">") {
		return quote
	} else if strings.HasPrefix(line, "+") {
		return add
	} else if strings.HasPrefix(line, "-") {
		return del
	} else {
		return text
	}
}
