package receiver

import "github.com/google/cel-go/cel"

type AttributeType string

const (
	TypeString AttributeType = "string"
	TypeFloat  AttributeType = "float"
	TypeInt    AttributeType = "int"
)

type AttributeSchema struct {
	Name     string
	Type     AttributeType
	Default  string
	Optional bool
}

func (t AttributeType) CelType() *cel.Type {
	switch t {
	case TypeString:
		return cel.StringType
	case TypeFloat:
		return cel.DoubleType
	case TypeInt:
		return cel.IntType
	default:
		return cel.DynType
	}
}
