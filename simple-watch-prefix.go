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

	watcher := cli.Watch(context.Background(), "my", clientv3.WithPrefix())
	for resp := range watcher {
		for _, ev := range resp.Events {
			fmt.Printf("operation=%s key=%q value=%q mod_revision=%d\n", ev.Type, ev.Kv.Key, ev.Kv.Value, ev.Kv.ModRevision)
		}
	}
}
