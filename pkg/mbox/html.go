package mbox

import (
	"os"
	"strings"
)

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

	j := parseJS(js)
	content = append(content, j...)

	content = append(content, "</header>")

	content = append(content, "<body>")

	content = append(content, "<div class=\"thread\">")
	parseThread(t.Node, &content)
	content = append(content, "</div>")

	content = append(content, "<div class=\"content\">")
	parseData(t.Node, &content)
	content = append(content, "</div>")

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

func parseThread(node *ThreadNode, content *[]string) {
	title := "<a href=\"#" + node.Mesg.MessageId +
		"\">" + node.Mesg.Subject + "</a>"
	if node.Mesg.Exist {
		title += " " + node.Mesg.From.Name
	}
	*content = append(*content, title)

	if len(node.Child) > 0 {
		*content = append(*content, "<ul>")
		for _, n := range node.Child {
			*content = append(*content, "<li>")
			parseThread(n, content)
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
		content = append(content, "<br /><br />")
		content = append(content, "</div>")
		content = append(content, "</div>")
		return content
	}

	header := parseHeader(m)
	content = append(content, header...)

	body := parseBody(m.Body)
	content = append(content, body...)

	content = append(content, "</div>")

	return content
}

func parseHeader(m *Message) []string {
	content := []string{}

	content = append(content, "<div class=\"message-header\">")

	subject := "<div class=\"subject\">" + m.Subject + "</div>"
	content = append(content, subject)

	date := "<div class=\"date\">" + m.Date.String() + "</div>"
	content = append(content, date)

	content = append(content, "<div class=\"from\">From:")
	content = append(content, "<ul>")
	from := "<li>" + m.From.Name +
		" <a href=\"mailto:" + m.From.Address +
		"\">&lt;" + m.From.Address + "&gt;</a></li>"
	content = append(content, from)
	content = append(content, "</ul>")
	content = append(content, "</div>")

	content = append(content, "<div class=\"to\">To:")
	content = append(content, "<ul>")
	for _, t := range m.To {
		to := "<li>" + t.Name +
			" <a href=\"mailto:" + t.Address + "\">&lt;" +
			t.Address + "&gt;</a></li>"
		content = append(content, to)
	}
	content = append(content, "</ul>")
	content = append(content, "</div>")

	content = append(content, "<div class=\"cc\">Cc:")
	content = append(content, "<ul>")
	for _, c := range m.Cc {
		cc := "<li>" + c.Name +
			" <a href=\"mailto:" + c.Address + "\">&lt;" +
			c.Address + "&gt;</a></li>"
		content = append(content, cc)
	}
	content = append(content, "</ul>")
	content = append(content, "</div>")

	content = append(content, "</div>")

	return content
}

func parseBody(lines []string) []string {
	content := []string{}

	content = append(content, "<div class=\"message-body\">")

	for _, line := range lines {
		content = append(content, parseLine(line))
	}

	content = append(content, "</div>")

	return content
}

func parseLine(line string) string {
	// Escape special symbols
	line = strings.ReplaceAll(line, "<", "&lt;")
	line = strings.ReplaceAll(line, ">", "&gt;")

	if line == "---" {
		line = "<div class=\"git-start\">" + line
	} else if line == "--" || line == "-- " {
		line = "<div class=\"git-end\">" + line
	} else if strings.HasPrefix(line, "--- ") {
		line = "<div class=\"git-before\">" + line
	} else if strings.HasPrefix(line, "+++ ") {
		line = "<div class=\"git-after\">" + line
	} else if strings.HasPrefix(line, "@@ ") {
		line = "<div class=\"git-change\">" + line
	} else if strings.HasPrefix(line, "diff ") {
		line = "<div class=\"git-diff\">" + line
	} else if strings.HasPrefix(line, "index ") {
		line = "<div class=\"git-index\">" + line
	} else if strings.HasPrefix(line, "&gt;") {
		line = "<div class=\"quote\">" + line
	} else if strings.HasPrefix(line, "-") {
		line = "<div class=\"git-delete\">" + line
	} else if strings.HasPrefix(line, "+") {
		line = "<div class=\"git-add\">" + line
	} else {
		line = "<div class=\"text\">" + line
	}

	line += "</div>"
	return line
}
