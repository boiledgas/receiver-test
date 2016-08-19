package main

import (
	"log"
	"receiver/cache"
	"receiver/data"
	"receiver/repository"
	"sync"
	"time"
)

var conf_cache *cache.Configuration

func reload() {
	defer time.AfterFunc(time.Second*1, reload)
	if count, err := conf_cache.ReloadCache(); err != nil {
		log.Printf("reload: %v", err)
	} else {
		log.Printf("reload: %v", count)
	}
}

func main() {
	wg := sync.WaitGroup{}
	var updateFunc cache.UpdateFunc = func(conf data.Configuration) {
		log.Printf("updated: %v %v", conf.Id, conf.ETag)
	}
	repository := repository.Configuration{
		Data: make(map[data.CodeId]data.Configuration),
	}
	conf_cache = &cache.Configuration{
		UpdateFunc: updateFunc,
		Repository: &repository,
		Index:      make(map[data.CodeId]data.Configuration),
	}
	updateConfiguration := func(code data.CodeId, pause time.Duration) {
		conf := data.Configuration{}
		for i := 0; i < 5; i ++ {
			if err := conf_cache.GetByCode(code, &conf); err != nil {
				log.Printf("find %v: %v", code, err)
			} else {
				break
			}
			time.Sleep(time.Second * 1)
		}
		for i := 0; i < 100; i++ {
			repository.Update(&conf)
			time.Sleep(pause)
		}
		wg.Done()
	}
	wg.Add(3)
	go updateConfiguration("test1", time.Millisecond * 500)
	go updateConfiguration("test2", time.Millisecond * 2000)
	go updateConfiguration("test3", time.Millisecond * 1500)
	go reload()
	time.Sleep(time.Second * 1)
	repository.TestData()
	wg.Wait()
}
