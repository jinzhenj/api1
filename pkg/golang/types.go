package golang

type CodeGen interface {
	Code() string
}

type GoTypeDef struct {
	Name string
	Type *GoType
}

type IntOrString struct {
	IntVal *int64
	StrVal *string
}

type GoEnumOption struct {
	Comments []string
	Name     string
	TypeName string
	Value    IntOrString
}

type GoEnumCodeBlock struct {
	Comments []string
	Name     string
	BaseType *GoType
	Options  []GoEnumOption
}

type GoType struct {
	Name      string
	KeyType   *GoType
	ItemType  *GoType
	IsPointer bool
}

type GoStructField struct {
	Comments []string
	Name     string
	Type     *GoType
	Tags     map[string]string
}

type GoStructType struct {
	Comments []string
	Name     string
	Fields   []GoStructField
}

type GoParam struct {
	Comments []string
	Name     string
	Type     *GoType
}

type GoFunction struct {
	Comments   []string
	Name       string
	Receiver   *GoParam
	Params     []GoParam
	RetTypes   []GoType
	InIface    bool
	Statements []GoStatement
}

type GoInterface struct {
	Comments  []string
	Name      string
	Functions []GoFunction
}

type GoStatement interface {
	CodeGen
}

type RouteStatement struct {
	Comments    []string
	Name        string
	Method      string
	Path        string
	Middlewares []string
	PathParams  []GoStructField
	QueryParams []GoStructField
	BodyParam   *GoParam
	ParamExprs  []string
	HasRet      bool
}

type GoFile struct {
	Name     string
	Package  string
	Imports  []string
	CodeGens []CodeGen
}

type RawCode struct {
	code string
}
