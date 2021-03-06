package main

import (
	"github.com/boiledgas/receiver-test"
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/data/cache"
	"github.com/boiledgas/receiver-test/data/repository"
	rec "github.com/boiledgas/receiver-test/receiver"
	_ "github.com/boiledgas/receiver-test/parser"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"time"
)

func Serve() {
	var err error
	var bytes []byte
	if bytes, err = ioutil.ReadFile("/go/src/github.com/boiledgas/receiver-test/config.yaml"); err != nil {
		log.Printf("read file: %v", err)
		return
	}
	cfg := receiver.Config{} // heap
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		log.Printf("parse yaml config: %v", err)
		return
	}

	repository := repository.Configuration{}
	repository.Init()
	repository.TestData()
	cache := &cache.Configuration{
		Repository: &repository,
		Index:      make(map[data.CodeId]data.ConfId),
		Cache:      make(map[data.ConfId]data.Conf)}
	contextProvider := &rec.ContextProvider{
		Cache:    cache,
		Contexts: make(map[data.CodeId]*rec.Context),
	}
	cache.UpdateFunc = contextProvider.UpdateConfiguration

	// heap
	service := receiver.Service{
		Config:   cfg,
		Cache:    cache,
		Provider: contextProvider,
	}
	if err = service.ListenAndServe(); err != nil {
		log.Printf("service.ListenAndServe: %v", err)
		return
	}
	time.Sleep(time.Second * 60 * 5)
	if err = service.Stop("Test"); err != nil {
		log.Printf("stop: %v", err)
		return
	}
	time.Sleep(time.Second * 5)
	if err = service.Start("Test"); err != nil {
		log.Printf("start: %v", err)
		return
	}
	time.Sleep(time.Second * 10)
	if err = service.Stop("Test"); err != nil {
		log.Printf("stop: %v", err)
		return
	}
}

func main() {
	Serve()
}
