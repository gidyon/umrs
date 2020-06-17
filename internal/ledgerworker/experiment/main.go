package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	v1(20)
	v2(20)
}

func v1(count int) {
	start := time.Now()

	wg := &sync.WaitGroup{}

	for i := 0; i < count; i++ {
		wg.Add(1)
		go process1(wg)
	}

	wg.Wait()

	dur := time.Now().Sub(start).Milliseconds()

	fmt.Printf("v1(mutex+waitgroup) took %d\n", dur)
}

func v2(count int) {
	start := time.Now()

	for i := 0; i < count; i++ {
		process2()
	}

	dur := time.Now().Sub(start).Milliseconds()

	fmt.Printf("took(single) %d\n", dur)
}

func process1(wg *sync.WaitGroup) {
	defer wg.Done()

	time.Sleep(time.Second)
}

func process2() {
	time.Sleep(time.Second)
}
