package mbox

import (
	"errors"
)

type ThreadNode struct {
	Mesg  *Message
	Child []*Message
}

func Create(messages []*Message) (*ThreadNode, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	id, thread, err := check(messages)
	if err != nil {
		return nil, err
	}

	queue := []*ThreadNode{thread}
	delete(id, thread.Mesg.MessageId)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, m := range id {
			if m.InReplyTo == node.Mesg.MessageId {
				node.Child = append(node.Child, m)
				queue = append(queue, &ThreadNode{m, []*Message{}})
				delete(id, m.MessageId)
			}
		}
	}
	return thread, nil
}

func check(messages []*Message) (map[string]*Message, *ThreadNode, error) {
	id := map[string]*Message{}
	for _, message := range messages {
		if _, ok := id[message.MessageId]; ok {
			return nil, nil, errors.New("mbox: duplicate Message-Id")
		}
		id[message.MessageId] = message
	}

	thread := &ThreadNode{messages[0], []*Message{}}
	for _, message := range messages[1:] {
		if _, ok := id[message.InReplyTo]; !ok {
			m := &Message{
				MessageId: message.InReplyTo,
				InReplyTo: thread.Mesg.MessageId,
				Exist:     false,
			}
			id[m.MessageId] = m
		}
	}
	return id, thread, nil
}
