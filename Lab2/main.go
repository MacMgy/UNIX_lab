package main

import (
	"fmt"
	"sync"
	"time"
)

type Massage struct {
	data      string
	processed bool

	mux *sync.Cond
}

func RunConsumer(wg *sync.WaitGroup, msg *Massage, ch chan int) {
	defer func() {
		fmt.Println("*** Successfully completed")
		wg.Done()
	}()
	for {
		select {
		case _, ok := <-ch:
			if ok {
				msg.mux.L.Lock()

				fmt.Printf("Got massage:  %q\n", msg.data)
				processMassage(msg)

				msg.mux.Signal()
				msg.mux.L.Unlock()
			} else {
				return
			}
		}
	}
}

func processMassage(msg *Massage) {
	msg.data = ""
	msg.processed = true
}

func RunProvider(wg *sync.WaitGroup, msg *Massage, ch chan int) {
	defer func() {
		close(ch)
		wg.Done()
	}()

	massages := []string{"context", "go.uuid", "exist", "models", "service"}
	for i := range massages {
		// just sleep
		time.Sleep(time.Second * 1)

		msg.mux.L.Lock()
		if !msg.processed {
			msg.mux.Wait()
		}

		prepareMassage(msg, massages[i])
		fmt.Printf("Send massage: %q\n", msg.data)

		msg.mux.L.Unlock()
		ch <- 0
	}
	return
}

func prepareMassage(msg *Massage, data string) {
	msg.data = data
	msg.processed = false
}

func main() {
	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)
	msg := Massage{
		processed: true,
		mux:       sync.NewCond(&mux),
	}
	wg.Add(2)

	ch := make(chan int, 10)
	go RunConsumer(&wg, &msg, ch)
	go RunProvider(&wg, &msg, ch)

	wg.Wait()
	return
}
