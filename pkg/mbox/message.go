package mbox

import (
	"errors"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"
)

// Message represents a parsed mail message.
type Message struct {
	MessageId string
	InReplyTo string
	From      *mail.Address
	To        []*mail.Address
	Cc        []*mail.Address
	Subject   string
	Date      time.Time
	Body      []string
	Exist     bool
}

// Read messages from mbox data.
func Read(data []byte) ([]*Message, error) {
	messages := []*Message{}
	contents := strings.Split(string(data), "From mboxrd@z Thu Jan  1 00:00:00 1970\n")
	for _, content := range contents[1:] {
		m, err := mail.ReadMessage(strings.NewReader(content))
		if err != nil {
			return nil, fmt.Errorf("mbox: invalid Message: %v", err)
		}

		message, err := convert(m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func convert(m *mail.Message) (*Message, error) {
	message := &Message{}

	// Ignore the content after the angle brackets
	id := m.Header.Get("Message-Id")
	start := strings.Index(id, "<")
	end := strings.LastIndex(id, ">")
	if start < 0 || end < 0 {
		return nil, errors.New("mbox: invalid Message-Id")
	}
	message.MessageId = id[start : end+1]

	// Ignore the content after the angle brackets
	id = m.Header.Get("In-Reply-To")
	if id != "" {
		start = strings.Index(id, "<")
		end = strings.LastIndex(id, ">")
		if start < 0 || end < 0 {
			return nil, errors.New("mbox: invalid In-Reply-To")
		}
		message.InReplyTo = id[start : end+1]
	}

	from, err := m.Header.AddressList("From")
	if err != nil {
		return nil, fmt.Errorf("mbox: invalid From: %v", err)
	}
	message.From = from[0]

	message.To, err = m.Header.AddressList("To")
	if err != nil {
		return nil, fmt.Errorf("mbox: invalid To: %v", err)
	}

	if m.Header.Get("Cc") != "" {
		message.Cc, err = m.Header.AddressList("Cc")
		if err != nil {
			return nil, fmt.Errorf("mbox: invalid Cc: %v", err)
		}
	}

	message.Subject = m.Header.Get("Subject")

	message.Date, err = m.Header.Date()
	if err != nil {
		return nil, fmt.Errorf("mbox: invalid Date: %v", err)
	}

	body, err := io.ReadAll(m.Body)
	if err != nil {
		return nil, fmt.Errorf("mbox: invalid Body: %v", err)
	}
	message.Body = strings.Split(string(body), "\n")
	// Remove the last empty line
	message.Body = message.Body[:len(message.Body)-1]

	message.Exist = true

	return message, nil
}
