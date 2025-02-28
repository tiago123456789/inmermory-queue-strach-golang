package linkedlist

type Node struct {
	data interface{}
	next *Node
}

type LinkedList struct {
	head *Node
	tail *Node
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

func (l *LinkedList) Add(value interface{}) {
	if l.head == nil {
		head := Node{
			data: value,
			next: nil,
		}
		l.head = &head
		l.tail = &head
		return
	}

	if l.head.next == nil {
		item := Node{
			data: value,
			next: nil,
		}
		l.head.next = &item
		l.tail = &item
		return
	}

	oldTail := l.tail
	l.tail = &Node{
		data: value,
		next: nil,
	}
	oldTail.next = l.tail
}

func (l *LinkedList) Get() interface{} {
	if l.head == nil {
		return nil
	}

	oldHead := *l.head

	if oldHead.data == nil {
		return nil
	}

	newHead := l.head.next
	l.head = newHead
	return oldHead.data
}
