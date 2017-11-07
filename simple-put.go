package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	value := fmt.Sprintf("my-value-%d", rand.Int())
	resp, err := cli.Put(context.Background(), "my-key", value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Set value of my-key to: %s\n", value)
	fmt.Printf("The revision created: %d\n", resp.Header.Revision)
}
