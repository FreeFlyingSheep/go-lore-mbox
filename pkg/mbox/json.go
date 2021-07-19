package mbox

import (
	"encoding/json"
	"net/mail"
)

// Message thread node for generating JSON.
type ThreadNodeJSON struct {
	Mesg  *MessageJSON
	Child []*ThreadNodeJSON
}

// Message thread for generating JSON.
type ThreadJSON struct {
	Name string
	Node *ThreadNodeJSON
}

// AddressJSON represents a single mail address for generating JSON.
type AddressJSON struct {
	Name    string
	Address string
}

// BodyJSON represents a parsed mail body for generating JSON.
type BodyJSON struct {
	Class string
	Data  []string
}

// MessageJSON represents a parsed message for generating JSON.
type MessageJSON struct {
	MessageId string
	InReplyTo string
	From      *AddressJSON
	To        []*AddressJSON
	Cc        []*AddressJSON
	Subject   string
	Date      string
	Body      []*BodyJSON
	Exist     bool
}

// ParseJSON parses the thread for generating JSON.
func (t *Thread) ParseJSON() ([]byte, error) {
	thread := ThreadJSON{
		Name: t.Name,
		Node: &ThreadNodeJSON{},
	}
	convertNode(t.Node, thread.Node)
	return json.MarshalIndent(thread, "", "    ")
}

func convertNode(t1 *ThreadNode, t2 *ThreadNodeJSON) {
	t2.Mesg = convertMessage(t1.Mesg)
	t2.Child = []*ThreadNodeJSON{}

	if len(t1.Child) > 0 {
		for _, n := range t1.Child {
			node := &ThreadNodeJSON{}
			t2.Child = append(t2.Child, node)
			convertNode(n, node)
		}
	}
}

func convertMessage(m *Message) *MessageJSON {
	message := &MessageJSON{
		MessageId: m.MessageId,
		InReplyTo: m.InReplyTo,
		Subject:   m.Subject,
		Date:      m.Date.String(),
		Exist:     m.Exist,
	}

	if !m.Exist {
		return message
	}

	message.From = convertAddress(m.From)

	message.To = []*AddressJSON{}
	for _, a := range m.To {
		message.To = append(message.To, convertAddress(a))
	}

	message.Cc = []*AddressJSON{}
	for _, a := range m.Cc {
		message.Cc = append(message.Cc, convertAddress(a))
	}

	message.Body = parseBodyJSON(m.Body)

	return message
}

func convertAddress(a *mail.Address) *AddressJSON {
	address := &AddressJSON{
		Name:    a.Name,
		Address: a.Address,
	}
	return address
}

func parseBodyJSON(lines []string) []*BodyJSON {
	body := []*BodyJSON{}
	content := &BodyJSON{styles[undefined], []string{}}
	s, last := undefined, undefined

	for _, line := range lines {
		s = parseStyle(line)
		if last != undefined && s != last {
			body = append(body, content)
			content = &BodyJSON{styles[undefined], []string{}}
		}
		for k, v := range styles {
			if last != k && s == k {
				content.Class = v
				break
			}
		}
		last = s
		content.Data = append(content.Data, line)
	}

	return body
}
