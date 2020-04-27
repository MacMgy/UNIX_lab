package main

import (
	"fmt"
	"sync"
	"time"
)

type SyncPack struct {
	mux *sync.Cond
	processed bool
}

func RunConsumer(wg *sync.WaitGroup, sync *SyncPack, ch chan string) {
	defer func() {
		fmt.Println("*** Successfully completed")
		wg.Done()
	}()
	for {
		// Выполнение goroutine останавливается
		// до поступления очередного сообщения
		select {
		case msg, ok := <-ch:
			if ok {
				sync.mux.L.Lock()

				fmt.Printf("Got massage:  %q\n", msg)
				sync.processed = true

				// Сигнализируем ожидающей функции, что можно продолжать работу
				sync.mux.Signal()
				sync.mux.L.Unlock()
			} else {
				// Закрытие канала сообщений сигнализирует о завершении работы
				return
			}
		}
	}
}

func RunProvider(wg *sync.WaitGroup, msg *SyncPack, ch chan string) {
	defer func() {
		close(ch)
		wg.Done()
	}()

	massages := []string{"context", "go.uuid", "exist", "models", "service"}
	for i := range massages {
		time.Sleep(time.Second * 1)

		msg.mux.L.Lock()
		if !msg.processed {
			// Если сообщение еще в процессе обработки,
			// ожидаем сигнала о завершении
			msg.mux.Wait()
		}

		msg.processed = false
		fmt.Printf("Send massage: %q\n", massages[i])

		msg.mux.L.Unlock()
		ch <- massages[i]
	}
	return
}

func main() {
	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)
	msg := SyncPack{
		processed: true,
		mux:       sync.NewCond(&mux),
	}
	wg.Add(2)

	// Буферизированный канал хранит поступающие значения в виде очереди
	ch := make(chan string, 10)
	go RunConsumer(&wg, &msg, ch)
	go RunProvider(&wg, &msg, ch)

	// Ожидаем завершения наших goroutine
	wg.Wait()
	return
}
