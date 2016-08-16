package ittt

import "receiver/errors"

// switch activation
type toggle struct {
	state   bool   // on/off state
	life    uint16 // switch on-lifetime
	onTime  uint64 // time of last turn on
	offTime uint64 // time of last turn off
}

type aggregationOp byte

// aggregation activation
type bunch struct {
	op       aggregationOp      // aggregation operation
	children [10]aggregationKey // children of aggregation

}

type aggregationKey struct {
	id        uint32 // identity
	aggregate bool   // true = aggregation; false = switch
}

type SwitchNetwork struct {
	Switches     map[uint16]toggle            // switches
	Relations    map[uint16][3]aggregationKey // parents of switch
	Aggregations map[uint16]bunch             // bunches
}

func (n *SwitchNetwork) Register(id uint16)

func (n *SwitchNetwork) Turn(id uint16, state bool) (err error) {
	var t toggle
	var ok bool
	if t, ok = n.Switches[id]; !ok {
		err = errors.New("not found")
		return
	}

	if t.state == state {
		return
	}

	return
}
