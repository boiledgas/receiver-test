package ittt

type Condition struct {
	Id   uint16
	Type uint16
}

type Aggregate struct {
	Condition
}

type Simple struct {
	Condition
	Value float32
}

type Reaction struct {

}

type Event struct {
	Conditions map[uint16]Condition
	Root       uint16
	Action     Reaction
}
