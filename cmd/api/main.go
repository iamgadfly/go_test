package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const url = "https://httpbin.org/delay/%d"

func main() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	done := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			rpm := time.Tick(time.Minute / time.Duration(150))
			for {
				select {
				case <-done:
					fmt.Printf("goroutine %d end\n", id)
					return
				case <-rpm:
					sendRequest(client, id)
				}
			}
		}(i)
	}
	time.Sleep(5 * time.Minute)

	close(done)
	wg.Wait()

	log.Println("end")
}

// sendRequest ..
func sendRequest(client *http.Client, id int) {
	delay := rand.Intn(6) + 1
	url := fmt.Sprintf(url, delay)
	log.Printf("goroutine %d send req %s\n", id, url)
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("goroutine %d: error: %v\n", id, err)
		return
	}
	defer resp.Body.Close()
	log.Printf("goroutine %d: resp: %s\n", id, resp.Status)
}
