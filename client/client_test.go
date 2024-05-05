package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestNewClients(t *testing.T) {
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int) {
			client, err := New("localhost:5001")
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
}

func TestNewClient(t *testing.T) {
	client, err := New("localhost:5001") // TODO: mock client server
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		fmt.Println("set this: ", fmt.Sprintf("bar %d", i))
		if err := client.Set(context.TODO(), fmt.Sprintf("foo %d", i), fmt.Sprintf("bar %d", i)); err != nil {
			log.Fatal(err)
		}
		val, err := client.Get(context.TODO(), fmt.Sprintf("foo %d", i))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("got this back: ", val)
	}
}
