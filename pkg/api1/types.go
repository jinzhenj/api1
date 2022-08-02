package api1

type HasName struct {
	Name string `json:"name,omitempty"`
}

type HasComments struct {
	Comments     []string               `json:"comments,omitempty"`
	PostComments []string               `json:"postComments,omitempty"`
	SemComments  map[string]interface{} `json:"semComments,omitempty"`
}

// Name & ItemType cannot be both set
type TypeRef struct {
	HasName
	ItemType *TypeRef `json:"itemType,omitempty"`
	Nullable bool     `json:"nullable"`
}

type ScalarType struct {
	HasName
	HasComments
}

type IntOrString struct {
	IntVal *int64  `json:"intVal,omitempty"`
	StrVal *string `json:"strVal,omitempty"`
}

type EnumOption struct {
	HasName
	HasComments
	Value *IntOrString `json:"value,omitempty"`
}

type EnumType struct {
	HasName
	HasComments
	Options []EnumOption `json:"options"`
}

type StructField struct {
	HasName
	HasComments
	Type *TypeRef `json:"type"`
}

type StructType struct {
	HasName
	HasComments
	Fields []StructField `json:"fields"`
}

type Param struct {
	HasName
	HasComments
	Type *TypeRef `json:"type"`
}

type Fun struct {
	HasName
	HasComments
	Params []Param  `json:"params"`
	Type   *TypeRef `json:"type"`
}

type Iface struct {
	HasName
	HasComments
	Funs []Fun `json:"funs"`
}

type ApiGroup struct {
	HasName
	HasComments
	ScalarTypes []ScalarType `json:"scalarTypes"`
	EnumTypes   []EnumType   `json:"enumTypes"`
	StructTypes []StructType `json:"structTypes"`
	Ifaces      []Iface      `json:"ifaces"`
}

type Schema struct {
	Groups []ApiGroup `json:"groups"`
}
