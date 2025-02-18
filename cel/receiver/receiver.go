package receiver

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

func NewPartyProvider(attrs []AttributeSchema) cel.EnvOption {
	return func(env *cel.Env) (*cel.Env, error) {
		p := &Party{
			baseProvider: env.CELTypeProvider(),
			baseAdapter:  env.CELTypeAdapter(),
			AttrSchema:   attrs,
		}
		env, err := cel.CustomTypeAdapter(p)(env)
		if err != nil {
			return nil, err
		}

		return cel.CustomTypeProvider(p)(env)
	}
}

var PartyType = cel.ObjectType(
	"party",
	traits.ComparerType,
	traits.ReceiverType,
)

type Party struct {
	Attributes   []Attribute
	baseProvider types.Provider
	baseAdapter  types.Adapter
	AttrSchema   []AttributeSchema
}

// HasTrait implements ref.Type.
func (p *Party) HasTrait(trait int) bool {
	return trait == traits.ReceiverType
}

// TypeName implements ref.Type.
func (p *Party) TypeName() string {
	return PartyType.TypeName()
}

// NativeToValue implements ref.TypeAdapter.
func (p *Party) NativeToValue(value any) ref.Val {
	return p.baseAdapter.NativeToValue(value)
}

// EnumValue implements types.Provider.
func (p *Party) EnumValue(enumName string) ref.Val {
	return p.baseProvider.EnumValue(enumName)
}

// FindIdent implements types.Provider.
func (p *Party) FindIdent(identName string) (ref.Val, bool) {
	return p.baseProvider.FindIdent(identName)
}

// FindStructFieldNames implements types.Provider.
func (p *Party) FindStructFieldNames(structType string) ([]string, bool) {
	names, found := p.baseProvider.FindStructFieldNames(structType)
	if structType != "party" {
		return names, found
	}

	if !found {
		return []string{}, false
	}

	for _, attr := range p.Attributes {
		names = append(names, attr.Name)
	}

	return names, true
}

// FindStructFieldType implements types.Provider.
func (p *Party) FindStructFieldType(structType string, fieldName string) (*types.FieldType, bool) {
	if tp, ok := p.baseProvider.FindStructFieldType(structType, fieldName); ok {
		return tp, ok
	}

	var schema AttributeSchema
	for _, attr := range p.AttrSchema {
		if attr.Name == fieldName {
			schema = attr
			break
		}
	}

	if schema.Name == "" {
		return p.baseProvider.FindStructFieldType(structType, fieldName)
	}

	if structType == "party" {
		return &types.FieldType{
			Type: schema.Type.CelType(),
			IsSet: func(val any) bool {
				if pt, ok := val.(*Party); ok {
					for _, attr := range pt.Attributes {
						if attr.Name == fieldName {
							return true
						}
					}
				}
				return false
			},
			GetFrom: func(val any) (any, error) {
				if pt, ok := val.(*Party); ok {
					for _, attr := range pt.Attributes {
						if attr.Name == fieldName {
							switch attr.Type {
							case "float":
								return types.DefaultTypeAdapter.NativeToValue(attr.Value), nil
							case "int":
								return types.DefaultTypeAdapter.NativeToValue(attr.Value), nil
							case "string":
								return types.DefaultTypeAdapter.NativeToValue(attr.Value), nil
							case "bool":
								return types.DefaultTypeAdapter.NativeToValue(attr.Value), nil
							default:
								return nil, fmt.Errorf("unsupported attribute type: %s", attr.Type)
							}
						}
					}
				}
				return nil, fmt.Errorf("cannot get attr list")
			},
		}, true
	}

	return nil, false
}

// FindStructType implements types.Provider.
func (p *Party) FindStructType(structType string) (*types.Type, bool) {
	if structType == "party" {
		return PartyType, true
	}
	return p.baseProvider.FindStructType(structType)
}

// NewValue implements types.Provider.
func (p *Party) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	return p.baseProvider.NewValue(structType, fields)
}

type Attribute struct {
	Name    string
	Type    string
	Value   any
	Default any
}

var _ types.Provider = (*Party)(nil)
var _ types.Adapter = (*Party)(nil)
var _ ref.Type = (*Party)(nil)
