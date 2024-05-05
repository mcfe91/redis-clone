package client

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNewClient1(t *testing.T) {
	client, err := New("localhost:5001") // TODO: mock client server
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if err := client.Set(context.TODO(), "foo", "1"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
	val, err := client.Get(context.TODO(), "foo")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GET =>", val)
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
