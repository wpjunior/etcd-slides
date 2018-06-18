package main

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/fatih/color"
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

	session, err := concurrency.NewSession(
		cli,
		concurrency.WithTTL(10),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	mutex := concurrency.NewMutex(session, "/my-mutex-id/")

	for {
		color.Red("Esperando a oportunidade de ter o lock")
		err = mutex.Lock(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		color.Green("Consegui o lock é meu por 5 segundos ....")
		time.Sleep(time.Second * 5)
		err = mutex.Unlock(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		color.Yellow("Lock liberado")
	} // slide-end
}
