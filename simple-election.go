package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	election := concurrency.NewElection(
		session,
		"my-election-id",
	)

	uniqueId := fmt.Sprintf("localhost/%d", os.Getpid())
	fmt.Println("Eu sou o gorila", uniqueId)

	// Monitorar em uma nova goroutina se o gorila da bola azul mudou
	go func() {
		channel := election.Observe(context.Background())
		for {
			resp := <-channel
			if len(resp.Kvs) == 0 {
				fmt.Println(
					"Estamos sem gorila da bola azul",
				)
				continue
			}
			gorila := string(resp.Kvs[0].Value)
			if gorila == uniqueId {
				continue
			}
			fmt.Println(
				"O Gorila da bola azul é o",
				gorila,
			)
		} // end-slide
	}()

	// wait for election
	err = election.Campaign(
		context.Background(),
		uniqueId,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer election.Resign(context.Background())
	color.Blue(
		"Eu %s sou o sou o gorilão da bola azul!",
		uniqueId,
	)
	// end-slide
	// wait for os kill signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
