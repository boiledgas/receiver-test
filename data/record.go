package data

type Record struct {
	DeviceId interface{}     // идентификатор устройства
	Time     int64           // время записи
	Values   []PropertyValue // значения свойств
}
