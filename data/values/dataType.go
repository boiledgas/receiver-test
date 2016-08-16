package values

type DataType byte

const (
	DATATYPE_GPS DataType = 1
)

func (t DataType) GetValue() (v interface{}) {
	switch t {
	case DATATYPE_GPS:
		v = gpsValue{}
	}
	return
}