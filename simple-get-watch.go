package main

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var config atomic.Value

func main() {
	go func() {
		for {
			time.Sleep(time.Second)
			value := config.Load()
			fmt.Printf("In the memory the value is %q\n", value)
		}
	}()

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

	config.Store(resp.Kvs[0].Value)

	watcher := cli.Watch(context.Background(), "my-key", clientv3.WithRev(resp.Header.Revision))
	for resp := range watcher {
		for _, ev := range resp.Events {
			config.Store(ev.Kv.Value)
		}
	}
}
