package messages

import "container/list"

func bufferMessage(in <-chan Message) <-chan Message {
	out := make(chan Message)

	buffer := list.New()

	go func() {
		for {
			if buffer.Len() == 0 {
				buffer.PushBack(<-in)
			} else {
				select {
				case newMessage, ok := <-in:
					if ok {
						buffer.PushBack(newMessage)
					} else {
						close(out)
						return
					}
				case out <- buffer.Front().Value.(Message):
					buffer.Remove(buffer.Front())
				}
			}
		}
	}()

	return out
}
