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

	resp, err := cli.Get(
		context.Background(),
		"my",
		clientv3.WithPrefix(),
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, kv := range resp.Kvs {
		fmt.Printf(
			"The value of %s is %s\n",
			kv.Key,
			kv.Value,
		)
	}

}
