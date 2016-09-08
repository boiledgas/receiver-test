package data

import "time"

type Records struct {
	DeviceId ConfId   // идентификатор устройства
	Data     []Record // данные устройства
}

type Record struct {
	ConfId ConfId          // идентификатор устройства
	Time   time.Time       // время записи
	Values []PropertyValue // значения свойств
}
