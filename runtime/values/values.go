package values

type ValueType string

const (
	NoneValue   ValueType = "None"
	NumberValue ValueType = "Number"
)

type RtVal interface {
	GetType() ValueType
}

type NumVal struct {
	Type  ValueType `json:"type"`
	Value float64   `json:"value"`
}

func (nv *NumVal) GetType() ValueType { return nv.Type }

type NoneVal struct {
	Type  ValueType
	Value string
}

func (nov *NoneVal) GetType() ValueType { return nov.Type }
