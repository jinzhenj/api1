package types

type SourceInfo struct {
	FileName string
}

type Tag struct {
	Json    string
	Bson    string
	Binding bool
	Form    string
}

type Field struct {
	Name     string
	Tag      *Tag
	IsArray  bool
	Kind     string
	Comments string
}

type StructRecord struct {
	Name   string
	Fields []Field
	SInfo  SourceInfo
}
