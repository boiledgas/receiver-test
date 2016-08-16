package cache

import (
	"receiver/data"
	"receiver/repository"
	"testing"
)

func TestConfiguration_Empty(t *testing.T) {
	var updatedId interface{}
	var updateFunc UpdateFunc = func(conf data.Configuration) {
		updatedId = conf.Id
	}
	repository := repository.Configuration{}
	cache := Configuration{
		UpdateFunc: updateFunc,
		Repository: &repository,
		Cache:      make(map[data.CodeIdentity]data.Configuration),
	}
	cache.ReloadCache()

}

func TestConfiguration(t *testing.T) {
	var updatedId interface{}
	var updateFunc UpdateFunc = func(conf data.Configuration) {
		updatedId = conf.Id
	}
	repository := repository.Configuration{
		Data: make(map[data.CodeIdentity]data.Configuration),
	}
	repository.TestData()
	cache := Configuration{
		UpdateFunc: updateFunc,
		Repository: &repository,
		Cache:      make(map[data.CodeIdentity]data.Configuration),
	}
	var configuration data.Configuration
	if err := cache.Get("test1", &configuration); err != nil {
		t.Error(err)
	}
	if configuration.Id == nil {
		t.Errorf("configuration id is null")
	}
	repository.Update(&configuration)
	cache.ReloadCache()
	if updatedId == nil || updatedId != configuration.Id {
		t.Errorf("not updated %v != %v", updatedId, configuration.Id)
	}
}
