package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"receiver"
	"receiver/cache"
	"receiver/client"
	"receiver/config"
	"receiver/data"
	"receiver/logic"
	_ "receiver/parser"
	"receiver/repository"
	"sync"
	"time"
)

func Serve() {
	cfg := config.Service{} // heap
	var err error
	var bytes []byte
	if bytes, err = ioutil.ReadFile(".\\config.yaml"); err != nil {
		log.Print("read file: %v", err)
		return
	} else {
		if err = yaml.Unmarshal(bytes, &cfg); err != nil {
			log.Printf("parse yaml config: %v", err)
		}
	}

	go func() {
		for _, endpointCfg := range cfg.Receiver {
			for i := 0; i < 1000; i++ {
				// heap
				c := client.Client{Host: endpointCfg.Host, Port: endpointCfg.Port}
				if err := c.Start(); err != nil {
				}
			}

		}
	}()

	repository := repository.Configuration{}
	repository.Init()
	repository.TestData()
	cache := &cache.Configuration{Repository: &repository, Index: make(map[data.CodeId]data.ConfigurationId), Cache: make(map[data.ConfigurationId]data.Configuration)}
	contextProvider := &logic.ContextProvider{Cache: cache}
	cache.UpdateFunc = contextProvider.UpdateConfiguration

	// heap
	service := receiver.Service{
		Config:   cfg,
		Provider: contextProvider,
	}
	service.ListenAndServe()
	time.Sleep(time.Second * 5)
	if err = service.Stop("Test1"); err != nil {
		log.Printf("stop: %v", err)
	}
	time.Sleep(time.Second * 5)
	if err = service.Start("Test1"); err != nil {
		log.Printf("start: %v", err)
	}
	time.Sleep(time.Second * 10)
	if err = service.Stop("Test1"); err != nil {
		log.Printf("stop: %v", err)
	}
}

var wg sync.WaitGroup

func main() {
	Serve()
	//device := data.Device{Code: "device"}
	//device.Module(data.Module{Code: "module"})
	//device.Property("module", data.Property{Code: "property", Type: values.DATATYPE_GPS})
	//devices[device.Code] = device
	//
	//load := make(chan core.DeviceRequest)
	//go loadDevice(load)
	//go receiveRecords()
	//
	//context := core.Context{
	//	Load: load,
	//}
	//context.Id("device")
	//v := context.Value("module", "property", time.Now().UnixNano()).(values.GpsValue)
	//v.SetLatitude(54.55)
	//v.SetLongitude(35.55)
	//wg.Add(1)
	//context.Flush()
	//wg.Wait()
}
