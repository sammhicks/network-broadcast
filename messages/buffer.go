package messages

type element struct {
	value Message
	next  *element
}

type queue struct {
	frontElem *element
	backElem  *element
}

func (q *queue) Empty() bool {
	return q.frontElem == nil
}

func (q *queue) Front() Message {
	if !q.Empty() {
		return q.frontElem.value
	}
	panic("Empty queue has no front!")
}

func (q *queue) Pop() {
	if !q.Empty() {
		q.frontElem = q.frontElem.next
	}
}

func (q *queue) Push(m Message) {
	e := &element{
		next:  nil,
		value: m}

	if q.Empty() {
		q.frontElem = e
	} else {
		q.backElem.next = e
	}

	q.backElem = e
}

func bufferMessage(in <-chan Message) <-chan Message {
	out := make(chan Message)

	var buffer queue

	go func() {
		for {
			if buffer.Empty() {
				buffer.Push(<-in)
			} else {
				select {
				case m, ok := <-in:
					if ok {
						buffer.Push(m)
					} else {
						for !buffer.Empty() {
							out <- buffer.Front()
							buffer.Pop()
						}

						close(out)
						return
					}
				case out <- buffer.Front():
					buffer.Pop()
				}
			}
		}
	}()

	return out
}
