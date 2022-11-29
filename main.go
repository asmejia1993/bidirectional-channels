package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	numbers := make(chan int)
	doubles := make(chan int)
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	wg.Add(3)
	go counter(&wg, numbers, done)
	go doubler(&wg, numbers, doubles)
	go printer(&wg, doubles)
	wg.Wait()
}

func counter(wg *sync.WaitGroup, numbers chan<- int, done <-chan bool) {
	defer wg.Done()
	totalSent := 0
	for i := 1; ; i++ {
		numbers <- i
		totalSent++
		select {
		case _, ok := <-done:
			if ok {
				close(numbers)
				fmt.Printf("Closing numbers. Total sent: %d\n", totalSent)
				return
			}
		default:
		}
	}
}

func doubler(wg *sync.WaitGroup, numbers <-chan int, doubles chan<- int) {
	defer wg.Done()
	for {
		x, more := <-numbers
		if !more {
			close(doubles)
			fmt.Println("Closing doubles")
			return
		}
		doubles <- (x * 2)
	}
}

func printer(wg *sync.WaitGroup, doubles <-chan int) {
	defer wg.Done()
	totalReceived := 0
	for {
		double, more := <-doubles
		if !more {
			fmt.Printf("Closing printer, Total received: %d\n", totalReceived)
			return
		}
		totalReceived++
		fmt.Println(double)
	}
}
