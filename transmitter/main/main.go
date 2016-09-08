package main

import (
	"receiver/data"
	"receiver/data/cache"
	"receiver/data/repository"
	_ "receiver/parser"
	"receiver/source"
	"receiver/transmitter"
	"log"
	"time"
)

func main() {
	transmitterConfig := transmitter.Config{
		Parser:   "telematics",
		Protocol: "tcp",
		Server:   "localhost:777",
	}
	repository := repository.Configuration{
		Data: make(map[data.CodeId]data.Conf),
	}
	repository.TestData()
	cache := cache.Configuration{
		Repository: &repository,
		Index:      make(map[data.CodeId]data.ConfId),
		Cache:      make(map[data.ConfId]data.Conf),
	}
	var factory transmitter.ParserFactory
	var err error
	if factory, err = transmitter.Factory.Create("Telematics"); err != nil {
		log.Printf("factory: %v", err)
		return
	}
	transmitter := transmitter.SingleClientTransmitter{
		Config:  transmitterConfig,
		Cache:   &cache,
		Clients: make(map[data.CodeId]*transmitter.Client),
		Factory: factory,
	}
	if err := transmitter.Start(); err != nil {
		log.Printf("transmitter: %v", err)
		return
	}

	receiver := Receiver{Port: 777}
	if err := receiver.Start(); err != nil {
		log.Printf("receiver: %v", err)
		return
	}

	sourceConfig := source.Config{
		Ids: []data.ConfId{1},
	}
	source := source.Random{
		Cache:   &cache,
		Clients: make(map[data.CodeId]source.RandomClient),
		Config:  sourceConfig,
	}
	if err := transmitter.RegisterSource(&source); err != nil {
		log.Printf("source: %v", err)
		return
	}
	time.Sleep(time.Second * 100)
}
