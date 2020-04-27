package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Massage struct {
	data      string
	processed bool

	mux *sync.Cond
}

func GetMassage(msg *Massage) {
	msg.data = ""
	msg.processed = true
}

func SendMassage(msg *Massage, data string) {
	msg.data = data
	msg.processed = false
}

func RunConsumer(ctx context.Context, msg *Massage, ch chan int) {
	defer func() {
		fmt.Println("*** Successfully completed")
	}()
	for {
		select {
		case _, ok := <-ch:
			if ok {
				msg.mux.L.Lock()

				fmt.Printf("Got massage:  %q\n", msg.data)
				GetMassage(msg)

				msg.mux.Signal()
				msg.mux.L.Unlock()
			} else {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func RunProvider(ctx context.Context, msg *Massage, ch chan int) {
	defer close(ch)
	massages := []string{"context", "go.uuid", "exist", "models", "service"}

	for i := range massages {
		// just sleep
		time.Sleep(time.Second * 1)

		msg.mux.L.Lock()
		if !msg.processed {
			msg.mux.Wait()
		}

		SendMassage(msg, massages[i])
		fmt.Printf("Send massage: %q\n", msg.data)

		msg.mux.L.Unlock()
		ch <- 0

		if ctx.Err() == context.DeadlineExceeded {
			break
		}
	}
	return
}

func main() {
	ch := make(chan int, 10)
	mux := sync.Mutex{}
	msg := Massage{
		processed: true,
		mux:       sync.NewCond(&mux),
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)

	go RunConsumer(ctx, &msg, ch)
	go RunProvider(ctx, &msg, ch)

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Printf("Too long, canceling process...\n")
		}
		return
	}
}
