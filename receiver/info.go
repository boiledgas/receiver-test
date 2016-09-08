package receiver

type ConnectionInfo interface {
	// время создания подключения
	ConnectionTime() int64
	// время последнего пакета
	LastPacketDate() int64
	// количество отосланных пакетов
	PacketCount() uint32
}
