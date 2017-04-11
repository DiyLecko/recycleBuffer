package goRecycleBuffer

import (
	"container/list"
)

type RecycleBuffer struct {
	Get         chan []byte
	Give        chan []byte
	bufferSize  int
	bufferQueue *list.List
	makes       int
}

type queued struct {
	slice []byte
}

func (recycleBuffer *RecycleBuffer) makeBuffer() []byte {
	recycleBuffer.makes += 1
	return make([]byte, recycleBuffer.bufferSize)
}

func (recycleBuffer *RecycleBuffer) GetBufferCount() int {
	return recycleBuffer.makes
}

func Init(size int) *RecycleBuffer {
	recycleBuffer := new(RecycleBuffer)
	recycleBuffer.bufferSize = size
	recycleBuffer.bufferQueue = new(list.List)
	recycleBuffer.makes = 0

	recycleBuffer.Get = make(chan []byte)
	recycleBuffer.Give = make(chan []byte)

	go func() {
		for {
			if recycleBuffer.bufferQueue.Len() == 0 {
				recycleBuffer.bufferQueue.PushFront(queued{slice: recycleBuffer.makeBuffer()})
			}

			item := recycleBuffer.bufferQueue.Front()

			select {
			case buffer := <-recycleBuffer.Give:
				recycleBuffer.bufferQueue.PushBack(queued{slice: buffer})
			case recycleBuffer.Get <- item.Value.(queued).slice:
				recycleBuffer.bufferQueue.Remove(item)
			}
		}
	}()

	return recycleBuffer
}
