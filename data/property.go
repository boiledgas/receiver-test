package data

import "github.com/boiledgas/receiver-test/data/values"

type Property struct {
	Id       uint16
	Code     CodeId
	Type     values.DataType
	ModuleId uint16
}

type PropertyValue struct {
	ModuleId   uint16      // идентификатор модуля
	PropertyId uint16      // идентификатор свойства
	Value      interface{} // значение свойства
}
