package api1

import (
	"strings"

	"github.com/pkg/errors"
)

var builtinTypes []ScalarType = []ScalarType{
	{HasName: HasName{Name: "int"}},
	{HasName: HasName{Name: "float"}},
	{HasName: HasName{Name: "string"}},
	{HasName: HasName{Name: "boolean"}},
	{HasName: HasName{Name: "object"}},
	{HasName: HasName{Name: "any"}},
}

func (schema *Schema) Check() error {
	names := make(map[string]bool)
	types := make(map[string]interface{})

	addName := func(name string) error {
		if _, ok := names[name]; ok {
			return errors.Errorf("Type [%s] defined more than once", name)
		}
		names[name] = true
		return nil
	}
	addType := func(name string, t interface{}) error {
		if err := addName(name); err != nil {
			return err
		}
		types[name] = t
		return nil
	}
	var checkType func(t *TypeRef, nullable bool) error
	checkType = func(t *TypeRef, nullable bool) error {
		if t == nil {
			if nullable {
				return nil
			}
			return errors.New("Type cannot be empty")
		}
		if t.Name != "" {
			_, ok := types[t.Name]
			if !ok {
				return errors.Errorf("Type [%s] cannot be found", t.Name)
			}
			return nil
		}
		if t.ItemType != nil {
			return checkType(t.ItemType, false)
		}
		return nil
	}

	for _, t := range builtinTypes {
		addType(t.Name, t)
	}

	// check no duplications
	for _, g := range schema.Groups {
		for _, sc := range g.ScalarTypes {
			if err := addType(sc.Name, sc); err != nil {
				return err
			}
		}
		for _, en := range g.EnumTypes {
			if err := addType(en.Name, en); err != nil {
				return err
			}
			var optionNames []string
			for _, option := range en.Options {
				optionNames = append(optionNames, option.Name)
			}
			dup := duplicated(optionNames)
			if len(dup) > 0 {
				return errors.Errorf("Enum [%s] has duplicated options [%s]",
					en.Name, strings.Join(dup, ", "))
			}
		}
		for _, st := range g.StructTypes {
			if err := addType(st.Name, st); err != nil {
				return err
			}
			var fieldNames []string
			for _, field := range st.Fields {
				fieldNames = append(fieldNames, field.Name)
			}
			dup := duplicated(fieldNames)
			if len(dup) > 0 {
				return errors.Errorf("Struct [%s] has duplicated fields [%s]",
					st.Name, strings.Join(dup, ", "))
			}
		}
		for _, iface := range g.Ifaces {
			if err := addName(iface.Name); err != nil {
				return err
			}
			var funNames []string
			for _, fun := range iface.Funs {
				funNames = append(funNames, fun.Name)
			}
			dup := duplicated(funNames)
			if len(dup) > 0 {
				return errors.Errorf("Interface [%s] has duplicated functions [%s]",
					iface.Name, strings.Join(dup, ", "))
			}
			for _, fun := range iface.Funs {
				var paramNames []string
				for _, param := range fun.Params {
					paramNames = append(paramNames, param.Name)
				}
				dup := duplicated(paramNames)
				if len(dup) > 0 {
					return errors.Errorf("Function [%s.%s] has duplicated params [%s]",
						iface.Name, fun.Name, strings.Join(dup, ", "))
				}
			}
		}
	}

	// check referenced types can all be found
	for _, g := range schema.Groups {
		for _, st := range g.StructTypes {
			for _, field := range st.Fields {
				if err := checkType(field.Type, false); err != nil {
					return err
				}
			}
		}
		for _, iface := range g.Ifaces {
			for _, fun := range iface.Funs {
				if err := checkType(fun.Type, true); err != nil {
					return err
				}
				for _, param := range fun.Params {
					if err := checkType(param.Type, false); err != nil {
						return err
					}
				}
			}
		}
	}

	// check enum options' value have same type
	for _, g := range schema.Groups {
		for _, en := range g.EnumTypes {
			var findNoVal, findIntVal, findStrVal bool
			for _, option := range en.Options {
				if option.Value == nil {
					findNoVal = true
				} else if option.Value.IntVal != nil {
					findIntVal = true
				} else {
					findStrVal = true
				}
			}
			if findIntVal && (findNoVal || findStrVal) {
				return errors.Errorf("Enum [%s] has mixed value types", en.Name)
			}
		}
	}

	// group duplicated???
	// check pkg not empty
	// check struct, enum, interface is not empty
	return nil
}

func duplicated(items []string) []string {
	m := make(map[string]int)
	for _, item := range items {
		if _, ok := m[item]; !ok {
			m[item] = 0
		}
		m[item]++
	}
	var duplicated []string
	for item, count := range m {
		if count > 1 {
			duplicated = append(duplicated, item)
		}
	}
	return duplicated
}
