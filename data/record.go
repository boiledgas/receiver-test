package data

type Records struct {
	DeviceId interface{} // идентификатор устройства
	Data     []Record    // данные устройства
}

type Record struct {
	DeviceId interface{}     // идентификатор устройства
	Time     int64           // время записи
	Values   []PropertyValue // значения свойств
}
