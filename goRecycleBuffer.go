package goRecycleBuffer

import (
	"container/list"
)

var Get chan []byte = nil
var Give chan []byte = nil

var bufferQueue = new(list.List)
var bufferSize = 8192
var makes int
var frees int

func makeBuffer() []byte {
	makes += 1
	return make([]byte, bufferSize)
}

func GetBufferCount() int {
	return makes
}

type queued struct {
	slice []byte
}

func Init(size int) {
	if size > 0 {
		bufferSize = size
	}

	if Get == nil && Give == nil {
		Get = make(chan []byte)
		Give = make(chan []byte)

		go func() {
			for {
				if bufferQueue.Len() == 0 {
					bufferQueue.PushFront(queued{slice: makeBuffer()})
				}

				item := bufferQueue.Front()

				select {
				case buffer := <-Give:
					bufferQueue.PushBack(queued{slice: buffer})
				case Get <- item.Value.(queued).slice:
					bufferQueue.Remove(item)
				}
			}
		}()
	}
}
