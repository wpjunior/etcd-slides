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
		log.Fatal("my-key is not defined")
	}

	fmt.Printf("Initial value: %q", resp.Kvs[0].Value)

	watcher := cli.Watch(
		context.Background(),
		"my-key",
		clientv3.WithRev(resp.Header.Revision),
	)
	for resp := range watcher {
		for _, ev := range resp.Events {
			fmt.Printf("New value: %q", ev.Kv.Value)
		}
	}
}
