package mbox

import (
	"os"
	"strconv"
	"strings"
)

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

var gid int

// Parse parses the thread for generating HTML.
func (t *Thread) Parse(css, js string) []string {
	content := []string{}

	content = append(content, "<!DOCTYPE html>")
	content = append(content, "<html>")

	content = append(content, "<head>")

	content = append(content, "<meta charset=\"utf-8\">")

	content = append(content, parseTitle(t))

	c := parseCSS(css)
	content = append(content, c...)

	content = append(content, "</header>")

	content = append(content, "<body>")

	content = append(content, "<div class=\"thread\">")
	parseThread(t.Node, &content, 1)
	content = append(content, "</div>")

	content = append(content, "<div class=\"content\">")
	parseData(t.Node, &content)
	content = append(content, "</div>")

	j := parseJS(js)
	content = append(content, j...)

	content = append(content, "</body>")

	content = append(content, "</html>")

	return content
}

func parseTitle(t *Thread) string {
	return "<title>" + t.Name + "</title>"
}

func parseCSS(css string) []string {
	content := []string{}

	data, err := os.ReadFile(css)
	if err != nil {
		return content
	}

	content = append(content, "<style>")
	c := strings.Split(string(data), "\n")
	content = append(content, c...)
	content = append(content, "</style>")
	return content
}

func parseJS(js string) []string {
	content := []string{}

	data, err := os.ReadFile(js)
	if err != nil {
		return content
	}

	content = append(content, "<script>")
	j := strings.Split(string(data), "\n")
	content = append(content, j...)
	content = append(content, "</script>")
	return content
}

func parseThread(node *ThreadNode, content *[]string, depth int) {
	title := "<a href=\"#" + node.Mesg.MessageId + "\">"
	if depth == 1 || !node.Mesg.Exist {
		title += node.Mesg.Subject + "</a>"
	} else {
		title += node.Mesg.From.Name + "</a>"
	}
	*content = append(*content, title)

	if len(node.Child) > 0 {
		*content = append(*content, "<ul>")
		for _, n := range node.Child {
			*content = append(*content, "<li>")
			parseThread(n, content, depth+1)
			*content = append(*content, "</li>")
		}
		*content = append(*content, "</ul>")
	}
}

func parseData(node *ThreadNode, content *[]string) {
	message := parseMessage(node.Mesg)
	*content = append(*content, message...)

	for _, n := range node.Child {
		parseData(n, content)
	}
}

func parseMessage(m *Message) []string {
	content := []string{}

	message := "<div id=\"" + m.MessageId + "\" class=\"message\">"
	content = append(content, message)

	if !m.Exist {
		content = append(content, "<div class=\"not-found\">")
		content = append(content, "<div>[not found]</div>")
		content = append(content, "</div>") // not-found
		content = append(content, "</div>") // message
		return content
	}

	header := parseHeader(m)
	content = append(content, header...)

	body := parseBody(m.Body)
	content = append(content, body...)

	content = append(content, "</div>") // message

	return content
}

func parseHeader(m *Message) []string {
	content := []string{}

	content = append(content, "<div class=\"found\">")

	subject := "<div class=\"subject\">" + m.Subject + "</div>"
	content = append(content, subject)

	date := "<div class=\"date\">" + m.Date.String() + "</div>"
	content = append(content, date)

	content = append(content, "</div>") // found

	button := parseButton("message-header")
	content = append(content, button...)

	button = parseButton("from")
	content = append(content, button...)
	content = append(content, "From:")
	content = append(content, "<ul>")
	from := "<li>" + m.From.Name +
		" <a href=\"mailto:" + m.From.Address +
		"\">&lt;" + m.From.Address + "&gt;</a></li>"
	content = append(content, from)
	content = append(content, "</ul>")
	content = append(content, "</div>") // from

	button = parseButton("to")
	content = append(content, button...)
	content = append(content, "To:")
	content = append(content, "<ul>")
	for _, t := range m.To {
		to := "<li>" + t.Name +
			" <a href=\"mailto:" + t.Address + "\">&lt;" +
			t.Address + "&gt;</a></li>"
		content = append(content, to)
	}
	content = append(content, "</ul>")
	content = append(content, "</div>") // to

	button = parseButton("cc")
	content = append(content, button...)
	content = append(content, "Cc:")
	content = append(content, "<ul>")
	for _, c := range m.Cc {
		cc := "<li>" + c.Name +
			" <a href=\"mailto:" + c.Address + "\">&lt;" +
			c.Address + "&gt;</a></li>"
		content = append(content, cc)
	}
	content = append(content, "</ul>")
	content = append(content, "</div>") // cc

	content = append(content, "</div>") // message-header

	return content
}

func parseBody(lines []string) []string {
	content := []string{}

	button := parseButton("message-body")
	content = append(content, button...)

	lines = parseLines(lines)
	content = append(content, lines...)

	content = append(content, "</div>") // message-body

	return content
}

func parseButton(class string) []string {
	content := []string{}
	id := strconv.Itoa(gid)

	button := "<div id=\"button-" + id + "\" class=\"button\">"
	content = append(content, button)

	button = "<a href=\"javascript:fold('" + id + "')\">[-]</a>"
	content = append(content, button)

	content = append(content, "</div>") // button

	div := "<div id=\"fold-" + id + "\" class=\"" + class + " fold\">"
	content = append(content, div)

	gid++
	return content
}

func parseLines(lines []string) []string {
	modes := map[mode]string{
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

	content := []string{}
	m, last := undefined, undefined

	for _, line := range lines {
		m = parseMode(line)
		if last != undefined && m != last {
			content = append(content, "</div>")
		}
		for k, v := range modes {
			if last != k && m == k {
				button := parseButton(v)
				content = append(content, button...)
				break
			}
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
	line = strings.ReplaceAll(line, "&", "&amp;")
	line = strings.ReplaceAll(line, "\"", "&quot;")
	line = strings.ReplaceAll(line, "<", "&lt;")
	line = strings.ReplaceAll(line, ">", "&gt;")
	line = "<pre>" + line + "</pre>"
	return line
}
