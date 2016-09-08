package cache

import (
	"receiver/data"
	"receiver/data/repository"
	"testing"
)

func TestConfiguration_Empty(t *testing.T) {
	var updatedId interface{}
	var updateFunc UpdateFunc = func(conf data.Conf) {
		updatedId = conf.Id
	}
	repository := repository.Configuration{}
	cache := Configuration{
		UpdateFunc: updateFunc,
		Repository: &repository,
		Index:      make(map[data.CodeId]data.Conf),
	}
	cache.ReloadCache()

}

func TestConfiguration(t *testing.T) {
	var updatedId interface{}
	var updateFunc UpdateFunc = func(conf data.Conf) {
		updatedId = conf.Id
	}
	repository := repository.Configuration{
		Data: make(map[data.CodeId]data.Conf),
	}
	repository.TestData()
	cache := Configuration{
		UpdateFunc: updateFunc,
		Repository: &repository,
		Index:      make(map[data.CodeId]data.Conf),
	}
	var configuration data.Conf
	if err := cache.GetByCode("test1", &configuration); err != nil {
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
