package values

import ()

type valueBase struct {
	mask byte
}

type gpsValue struct {
	valueBase
	latitude  float32
	longitude float32
	altitude  float32
	speed     float32
	course    float32
}

type GpsValue interface {
	SetLatitude(float32)
	GetLatitude() (float32, bool)
	SetLongitude(float32)
	GetLongitude() (float32, bool)
	SetAltitude(float32)
	GetAltitude() (float32, bool)
	SetSpeed(float32)
	GetSpeed() (float32, bool)
	SetCourse(float32)
	GetCourse() (float32, bool)
}

func (g gpsValue) SetLatitude(v float32) {
	g.latitude = v
	g.mask |= 0x01
}

func (g gpsValue) GetLatitude() (v float32, ok bool) {
	v = g.latitude
	ok = g.mask&0x01 > 0
	return
}

func (g gpsValue) SetLongitude(v float32) {
	g.longitude = v
	g.mask |= 0x02
}

func (g gpsValue) GetLongitude() (v float32, ok bool) {
	v = g.longitude
	ok = g.mask&0x02 > 0
	return
}

func (g gpsValue) SetAltitude(v float32) {
	g.altitude = v
	g.mask |= 0x04
}

func (g gpsValue) GetAltitude() (v float32, ok bool) {
	v = g.altitude
	ok = g.mask&0x04 > 0
	return
}

func (g gpsValue) SetSpeed(v float32) {
	g.speed = v
	g.mask |= 0x08
}

func (g gpsValue) GetSpeed() (v float32, ok bool) {
	v = g.speed
	ok = g.mask&0x08 > 0
	return
}

func (g gpsValue) SetCourse(v float32) {
	g.course = v
	g.mask |= 0x10
}

func (g gpsValue) GetCourse() (v float32, ok bool) {
	v = g.course
	ok = g.mask&0x10 > 0
	return
}
