package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/mcfe19/redis-clone/client"
)

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			client, err := client.New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer client.Close()

			key := fmt.Sprintf("client_foo_%d", it)
			value := fmt.Sprintf("client_bar_%d", it)
			if err := client.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}

			val, err := client.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("client %d got this val back: %s\n", it, val)

			wg.Done()
		}(i)
	}
	wg.Wait()

	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers but got %d", len(server.peers))
	}

}
