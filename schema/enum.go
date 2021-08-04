package schema

import "fmt"

func resolveEnumInit(t *Type, init interface{}) error {
	if t == nil || init == nil || t.Kind != KindEnum {
		return nil
	}
	enum := t.Enum
	if enum == nil {
		return nil
	}

	var option *EnumOption
	var name string
	switch v := init.(type) {
	case Nil:
		t.Init = v
	case *EnumOption:
		option = v
	case Expression:
		name = string(v)
	case string:
		name = v
	}
	if option != nil {
		t.Init = option
		return nil
	}
	option = enum.GetOption(name)
	if option != nil {
		return nil
	}
	return fmt.Errorf("%s:%d invalid enum option: %s:%d %s does not have an option named: %s",
		t.File.Path, t.Line, t.Enum.Type.File.Path, t.Enum.Type.Line, t.Enum.Name, name)
}
