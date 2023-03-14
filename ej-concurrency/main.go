package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	go processUsers()

	time.Sleep(3 * time.Second)
}

func processUsers() {
	ch := make(chan string)

	go func() {
		for true {
			fmt.Println(fmt.Sprintf("Show message %s", <-ch))
		}
	}()

	i := 0
	for true {
		go func() {
			ch <- fmt.Sprintf("Process user %d", i)
		}()
		time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
		i++
	}
}
