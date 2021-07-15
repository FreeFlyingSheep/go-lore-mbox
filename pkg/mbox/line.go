package mbox

import "strings"

type mode int

const (
	undefined mode = iota
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

func parseLines(lines []string) []string {
	content := []string{}
	m, last := undefined, undefined

	for _, line := range lines {
		m = parseMode(line)
		if last != undefined && m != last {
			content = append(content, "</div>")
		}
		if last != text && m == text {
			content = append(content, "<div class=\"text\">")
		} else if last != start && m == start {
			content = append(content, "<div class=\"git-start\">")
		} else if last != end && m == end {
			content = append(content, "<div class=\"git-end\">")
		} else if last != before && m == before {
			content = append(content, "<div class=\"git-before\">")
		} else if last != after && m == after {
			content = append(content, "<div class=\"git-after\">")
		} else if last != change && m == change {
			content = append(content, "<div class=\"git-change\">")
		} else if last != diff && m == diff {
			content = append(content, "<div class=\"git-diff\">")
		} else if last != index && m == index {
			content = append(content, "<div class=\"git-index\">")
		} else if last != add && m == add {
			content = append(content, "<div class=\"git-add\">")
		} else if last != del && m == del {
			content = append(content, "<div class=\"git-del\">")
		} else if last != quote && m == quote {
			content = append(content, "<div class=\"quote\">")
		}
		last = m
		content = append(content, parseLine(line))
	}

	content = append(content, "</div>")
	return content
}

func parseMode(line string) mode {
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

func parseLine(line string) string {
	// Escape special symbols
	line = strings.ReplaceAll(line, "<", "&lt;")
	line = strings.ReplaceAll(line, ">", "&gt;")
	line = "<div>" + line + "</div>"
	return line
}
