package mbox

import (
	"sort"
	"time"
)

// Message thread node contains a message and its children.
type ThreadNode struct {
	Mesg  *Message
	Child []*ThreadNode
}

// Message thread
type Thread struct {
	Name string
	Node *ThreadNode
}

// Create a thread from messages.
func Create(name string, messages []*Message) (*Thread, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	id, thread, err := check(name, messages)
	if err != nil {
		return nil, err
	}

	queue := []*ThreadNode{thread.Node}
	delete(id, thread.Node.Mesg.MessageId)
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, n := range id {
			// Add messages to its parent message
			if n.Mesg.InReplyTo == node.Mesg.MessageId {
				node.Child = append(node.Child, n)
				queue = append(queue, n)
				delete(id, n.Mesg.MessageId)
			}
		}
	}
	sortThread(thread.Node)
	return thread, nil
}

func check(name string, messages []*Message) (map[string]*ThreadNode, *Thread, error) {
	id := map[string]*ThreadNode{}
	for _, message := range messages {
		id[message.MessageId] = &ThreadNode{message, []*ThreadNode{}}
	}

	// The first email is always the head of the thread
	node := id[messages[0].MessageId]
	for _, message := range messages[1:] {
		// Fix missing messages
		if _, ok := id[message.InReplyTo]; !ok {
			m := &Message{
				MessageId: message.InReplyTo,
				InReplyTo: node.Mesg.MessageId,
				Subject:   "[not found]",
				Date:      time.Now(),
				Exist:     false,
			}
			id[m.MessageId] = &ThreadNode{m, []*ThreadNode{}}
		}
	}

	thread := &Thread{name, node}
	return id, thread, nil
}

func sortThread(node *ThreadNode) {
	if len(node.Child) == 0 {
		return
	}

	sort.Slice(node.Child, func(i, j int) bool {
		return node.Child[i].Mesg.Date.Before(node.Child[j].Mesg.Date)
	})

	for _, n := range node.Child {
		sortThread(n)
	}
}
