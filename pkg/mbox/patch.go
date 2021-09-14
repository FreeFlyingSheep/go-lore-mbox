package mbox

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type patch struct {
	subject string
	content string
}

// ParseJSON parses the thread for generating patches.
func (t *Thread) ParsePatch() ([][]string, error) {
	patches, err := findPatch(t)
	if err != nil {
		return nil, err
	}

	res := [][]string{}
	for _, p := range patches {
		subject, err := parseSubject(p.subject)
		if err != nil {
			return nil, err
		}

		newPatch := []string{subject, p.content}
		res = append(res, newPatch)
	}
	return res, nil
}

func findPatch(t *Thread) ([]*patch, error) {
	subject := t.Node.Mesg.Subject
	prefix, index, err := parsePrefix(subject)
	if err != nil {
		return nil, err
	}

	patches := []*patch{}
	if index < 0 {
		mesg := t.Node.Mesg
		pos := strings.Index(mesg.Subject, "]")
		subject = "0001" + mesg.Subject[pos+1:]
		p := &patch{subject, parseCotent(mesg)}
		patches = append(patches, p)
	} else {
		for _, node := range t.Node.Child {
			mesg := node.Mesg
			if strings.HasPrefix(mesg.Subject, prefix) {
				_, index, err := parsePrefix(mesg.Subject)
				if err != nil {
					return nil, err
				}
				pos := strings.Index(mesg.Subject, "]")
				subject = fmt.Sprintf("%04d%s", index, mesg.Subject[pos+1:])
				p := &patch{subject, parseCotent(mesg)}
				patches = append(patches, p)
			}
		}
	}
	return patches, nil
}

func parseSubject(subject string) (string, error) {
	// Escape special symbols
	reg, err := regexp.Compile(`[\/":*?<>| ]+`)
	subject = reg.ReplaceAllString(subject, "-")
	return subject, err
}

func parsePrefix(subject string) (string, int, error) {
	err := error(nil)
	index := 0
	pos := strings.Index(subject, "]")
	prefix := subject[0 : pos+1]

	segs := strings.Split(prefix, " ")
	if len(segs) > 2 {
		prefix = strings.Join(segs[0:len(segs)-1], " ")

		indexes := strings.Split(segs[len(segs)-1], "/")
		index, err = strconv.Atoi(indexes[0])
	} else {
		index = -1
	}

	return prefix, index, err
}

func parseCotent(mesg *Message) string {
	content := "From: " + mesg.From.Name + " <" + mesg.From.Address + ">\n"
	content += "Date: " + mesg.Date.Format("Mon Jan 02 15:04:05 -0700 2006") + "\n"
	content += "Subject: " + mesg.Subject + "\n\n"
	content += strings.Join(mesg.Body, "\n")
	return content
}
