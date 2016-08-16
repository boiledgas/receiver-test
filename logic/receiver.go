package logic

type Receiver interface {
	IsActive() bool // приемщик включен
	Start() error
	Stop() error
	Disconnect(uint32)
}
