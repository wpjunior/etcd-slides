package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var (
	discount atomic.Value
)

func main() {
	go loadConfig()

	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		price := 50.0
		fmt.Fprintf(w, "Temos o produto x\n")

		if discount.Load().(string) == "blackfriday" {
			price = price * 0.7
		}

		fmt.Fprintf(w, "Preço %0.2f", price)
	})
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func loadConfig() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	key := "discount"

	resp, err := cli.Get(context.Background(), key)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) == 0 {
		log.Fatalf("%s is not defined", key)
	}

	discount.Store(string(resp.Kvs[0].Value))

	watcher := cli.Watch(context.Background(), key, clientv3.WithRev(resp.Header.Revision))
	for resp := range watcher {
		for _, ev := range resp.Events {
			log.Printf("Discount now is %s", ev.Kv.Value)
			discount.Store(string(ev.Kv.Value))
		}
	}
}
