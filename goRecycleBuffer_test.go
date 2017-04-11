package goRecycleBuffer

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func TestRecycleBuffer(t *testing.T) {
	rb := Init(8192)

	buf := <-rb.Get
	rb.Give <- buf
}

var doneTest2 chan bool = make(chan bool)
var endCountTest2 int = 100
var testCountTest2 int = 10
var countText2 int = 0

func testRecycleBuffer2(rb *RecycleBuffer) {
	for i := 0; i < testCountTest2; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Int63()%500))
		buf := <-rb.Get
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63()%1000))
			rb.Give <- buf
			countText2++

			if countText2 == endCountTest2*testCountTest2 {
				doneTest2 <- true
			}
		}()
	}
}
func TestRecycleBuffer2(t *testing.T) {
	rb := Init(8192)

	for i := 0; i < endCountTest2; i++ {
		go testRecycleBuffer2(rb)
	}

	<-doneTest2

	t.Logf("buffer count : %d\n", rb.GetBufferCount())
}
