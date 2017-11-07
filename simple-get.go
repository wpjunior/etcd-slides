package main

import (
	"context"
	"fmt"
	"log"
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

	resp, err := cli.Get(context.Background(), "my-key")
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) == 0 {
		fmt.Println("my-key is not found")
		return
	}

	fmt.Printf("The revision of cluster is: %d\n", resp.Header.Revision)
	fmt.Printf("The value of my-key is %s\n", resp.Kvs[0].Value)
	fmt.Printf("The value was created in revision: %d\n", resp.Kvs[0].CreateRevision)
	fmt.Printf("The value was modified in revision: %d\n", resp.Kvs[0].ModRevision)
}
