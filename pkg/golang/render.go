package golang

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhenj/api1/pkg/api1"
	"github.com/jinzhenj/api1/pkg/utils"
)

const (
	defaultOutputDir = "pkg/api"
	defaultPackage   = "api"
	ginImportPath    = "github.com/gin-gonic/gin"
)

type Render struct {
	rParser   *api1.RouteParser
	imports   map[string]bool
	outputDir string
	scalars   map[string]scalarInfo
}

type scalarInfo struct {
	typ string
	pkg string
	def bool
}

func (r *Render) addImport(s string) {
	if r.imports == nil {
		r.imports = make(map[string]bool)
	}
	r.imports[s] = true
}

func (r *Render) popImports() []string {
	if r.imports == nil {
		return []string{}
	}
	var arr []string
	for key := range r.imports {
		arr = append(arr, key)
	}
	sort.Strings(arr)
	r.imports = nil
	return arr
}

// func isGoBuiltinType(typ string) bool {
// 	switch typ {
// 	case "bool", "string", "byte", "rune",
// 		"int", "int8", "int16", "int32", "int64",
// 		"uint", "uint8", "uint16", "uint32", "uint64",
// 		"float32", "float64", "complex64", "complex128":
// 		return true
// 	}
// 	return false
// }

func getScalarInfo(semComments map[string]interface{}) *scalarInfo {
	typ, ok := semComments["go.type"].(string)
	if !ok {
		return nil
	}
	_, defined := semComments["go.typeDef"]
	if pkg, ok := semComments["go.typePkg"].(string); ok {
		return &scalarInfo{typ: typ, pkg: pkg, def: defined}
	}
	parts := strings.Split(typ, "/")
	if len(parts) < 2 {
		return &scalarInfo{typ: typ, def: defined}
	}
	return &scalarInfo{
		typ: sprintf("%s.%s", parts[len(parts)-2], parts[len(parts)-1]),
		pkg: strings.Join(parts[:len(parts)-1], "/"),
		def: defined,
	}
}

func (r *Render) addScalar(sc *api1.ScalarType) {
	if r.scalars == nil {
		r.scalars = make(map[string]scalarInfo)
	}
	if s := getScalarInfo(sc.SemComments); s != nil && !s.def {
		r.scalars[sc.Name] = *s
	}
}

func (r *Render) renderScalar(sc *api1.ScalarType) *GoTypeDef {
	if s := getScalarInfo(sc.SemComments); s != nil && s.def {
		if s.pkg != "" {
			r.addImport(s.pkg)
		}
		return &GoTypeDef{
			Name: sc.Name,
			Type: &GoType{Name: s.typ},
		}
	}
	return nil
}

func (r *Render) getOutputDir() string {
	trimed := strings.TrimLeft(r.outputDir, "/")
	trimed = strings.TrimRight(trimed, "/")
	if len(trimed) == 0 {
		return defaultOutputDir
	}
	return trimed
}

func (r *Render) getPackage() string {
	outputDir := r.getOutputDir()
	parts := strings.Split(outputDir, "/")
	if len(parts) == 0 {
		return defaultPackage
	}
	return parts[len(parts)-1]
}

func (r *Render) Render(schema *api1.Schema) ([]GoFile, error) {
	r.rParser = &api1.RouteParser{}
	r.rParser.LoadSchema(schema)
	r.popImports()

	outputDir := r.getOutputDir()
	packageName := r.getPackage()

	// load scalars
	r.scalars = nil
	for _, g := range schema.Groups {
		for _, sc := range g.ScalarTypes {
			r.addScalar(&sc)
		}
	}

	// generate files
	var files []GoFile
	for _, g := range schema.Groups {
		file := GoFile{
			Name:    fmt.Sprintf("%s/%s.go", outputDir, g.Name),
			Package: packageName,
		}
		for _, sc := range g.ScalarTypes {
			if typeDef := r.renderScalar(&sc); typeDef != nil {
				file.CodeGens = append(file.CodeGens, typeDef)
			}
		}
		for _, en := range g.EnumTypes {
			file.CodeGens = append(file.CodeGens, r.renderEnum(&en))
		}
		for _, st := range g.StructTypes {
			file.CodeGens = append(file.CodeGens, r.renderStruct(&st))
		}
		for _, iface := range g.Ifaces {
			file.CodeGens = append(file.CodeGens, r.renderIface(&iface))
		}
		file.Imports = r.popImports()
		files = append(files, file)

		file2 := GoFile{
			Name:    fmt.Sprintf("%s/%s_route.go", outputDir, g.Name),
			Package: packageName,
		}
		for _, iface := range g.Ifaces {
			fun, err := r.renderRoutes(&iface)
			if err != nil {
				return nil, err
			}
			file2.CodeGens = append(file2.CodeGens, fun)
		}
		file2.Imports = r.popImports()
		files = append(files, file2)
	}
	files = append(files, r.renderHelperFile())
	return files, nil
}

func (r *Render) renderEnum(en *api1.EnumType) *GoEnumCodeBlock {
	e := GoEnumCodeBlock{
		Comments: en.Comments,
		Name:     en.Name,
		BaseType: &GoType{Name: "string"},
	}
	for _, op := range en.Options {
		var value IntOrString
		if op.Value == nil {
			strVal := op.Name
			value.StrVal = &strVal
		} else if op.Value.IntVal != nil {
			e.BaseType = &GoType{Name: "int64"}
			value.IntVal = op.Value.IntVal
		} else {
			value.StrVal = op.Value.StrVal
		}
		e.Options = append(e.Options, GoEnumOption{
			Comments: op.Comments,
			Name:     fmt.Sprintf("%s%s", en.Name, utils.PascalCase(op.Name)),
			TypeName: en.Name,
			Value:    value,
		})
	}
	return &e
}

func (r *Render) renderStruct(st *api1.StructType) *GoStructType {
	s := GoStructType{
		Comments: st.Comments,
		Name:     st.Name,
	}
	for _, sf := range st.Fields {
		s.Fields = append(s.Fields, r.renderStructField(&sf))
	}
	return &s
}

func (r *Render) renderStructField(sf *api1.StructField) GoStructField {
	f := GoStructField{
		Comments: sf.Comments,
		Name:     utils.PascalCase(sf.Name),
		Type:     r.renderType(sf.Type, sf.SemComments),
		Tags:     map[string]string{"json": sf.Name},
	}
	if _, ok := sf.SemComments["omitempty"]; ok {
		f.Tags["json"] = sf.Name + ",omitempty"
	}
	if _, ok := sf.SemComments["ignore"]; ok {
		f.Tags["json"] = "-"
	}
	return f
}

func (r *Render) renderType(t *api1.TypeRef, semComments map[string]interface{}) *GoType {
	if t == nil {
		return nil
	}
	var typ GoType
	typ.IsPointer = t.Nullable
	if semComments != nil {
		if s := getScalarInfo(semComments); s != nil {
			typ.Name = s.typ
			if s.pkg != "" {
				r.addImport(s.pkg)
			}
			return &typ
		}
	}
	if len(t.Name) > 0 {
		if t.Name == "object" {
			typ.KeyType = &GoType{Name: "string"}
			typ.ItemType = &GoType{Name: "interface{}"}
		} else if t.Name == "any" {
			typ.Name = "interface{}"
		} else if t.Name == "float" {
			typ.Name = "float64"
		} else if t.Name == "int" {
			typ.Name = "int64"
		} else if t.Name == "boolean" {
			typ.Name = "bool"
		} else if s, ok := r.scalars[t.Name]; ok {
			typ.Name = s.typ
			if s.pkg != "" {
				r.addImport(s.pkg)
			}
		} else {
			typ.Name = t.Name
		}
	} else {
		typ.ItemType = r.renderType(t.ItemType, nil)
	}
	return &typ
}

func (r *Render) renderIface(iface *api1.Iface) *GoInterface {
	i := GoInterface{
		Comments: iface.Comments,
		Name:     iface.Name,
	}
	for _, fun := range iface.Funs {
		i.Functions = append(i.Functions, *r.renderFun(&fun))
	}
	return &i
}

func (r *Render) renderFun(fun *api1.Fun) *GoFunction {
	f := GoFunction{
		Comments: fun.Comments,
		Name:     utils.PascalCase(fun.Name),
		InIface:  true,
	}
	if t := r.renderType(fun.Type, fun.SemComments); t != nil {
		f.RetTypes = append(f.RetTypes, *t)
	}
	for _, param := range fun.Params {
		f.Params = append(f.Params, r.renderParam(&param))
	}
	if _, ok := fun.SemComments["route"].(string); ok {
		f.Params = append(f.Params, GoParam{
			Name: "c",
			Type: &GoType{
				Name:      "gin.Context",
				IsPointer: true,
			},
		})
		f.RetTypes = append(f.RetTypes, GoType{
			Name: "error",
		})
		r.addImport(ginImportPath)
	}
	return &f
}

func (r *Render) renderParam(param *api1.Param) GoParam {
	p := GoParam{
		Comments: param.Comments,
		Name:     param.Name,
		Type:     r.renderType(param.Type, param.SemComments),
	}
	return p
}

func (r *Render) renderRoutes(iface *api1.Iface) (*GoFunction, error) {
	f := GoFunction{
		Comments: iface.Comments,
		Name:     "Register" + utils.PascalCase(iface.Name),
		Receiver: &GoParam{
			Name: "_r",
			Type: &GoType{
				Name:      "ApiRoutes",
				IsPointer: true,
			},
		},
		Params: []GoParam{
			{
				Name: "_o",
				Type: &GoType{
					Name: iface.Name,
				},
			},
		},
	}
	for _, fun := range iface.Funs {
		stmt, err := r.renderRouteStmt(iface, &fun)
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			f.Statements = append(f.Statements, stmt)
		}
	}
	return &f, nil
}

func (r *Render) renderRouteStmt(iface *api1.Iface, fun *api1.Fun) (GoStatement, error) {
	var route string
	var ok bool
	if route, ok = fun.SemComments["route"].(string); !ok {
		return nil, nil
	}

	r.addImport(ginImportPath)
	m, path, pathParams, err := api1.ParseRoute(route, api1.PathStyleColon)
	if err != nil {
		return nil, err
	}

	stmt := RouteStatement{
		Comments: fun.Comments,
		Name:     utils.PascalCase(fun.Name),
		Method:   m,
		Path:     path,
		HasRet:   fun.Type != nil,
	}

	if middleware, ok := fun.SemComments["go.middleware"]; ok {
		if middlewareStr, ok := middleware.(string); ok {
			stmt.Middlewares = append(stmt.Middlewares, middlewareStr)
		} else if middlewares, ok := middleware.([]interface{}); ok {
			for _, mid := range middlewares {
				if midStr, ok := mid.(string); ok {
					stmt.Middlewares = append(stmt.Middlewares, midStr)
				}
			}
		}
	}

	routeParams, err := r.rParser.ParseParams(iface, fun, pathParams)
	if err != nil {
		return nil, err
	}
	for _, param := range routeParams {
		var paramExpr string
		switch param.In {
		case api1.PositionBody:
			stmt.BodyParam = &GoParam{
				Comments: param.Comments,
				Name:     param.Name,
				Type:     r.renderType(param.Type, param.SemComments),
			}
			paramExpr = param.Name
		case api1.PositionPath:
			stmt.PathParams = append(stmt.PathParams, GoStructField{
				Comments: param.Comments,
				Name:     utils.PascalCase(param.Name),
				Type:     r.renderType(param.Type, param.SemComments),
				Tags:     map[string]string{"uri": param.Name},
			})
			paramExpr = fmt.Sprintf("_path.%s", utils.PascalCase(param.Name))
		case api1.PositionQuery:
			stmt.QueryParams = append(stmt.QueryParams, GoStructField{
				Comments: param.Comments,
				Name:     utils.PascalCase(param.Name),
				Type:     r.renderType(param.Type, param.SemComments),
				Tags:     map[string]string{"form": param.Name},
			})
			paramExpr = fmt.Sprintf("_query.%s", utils.PascalCase(param.Name))
		}
		stmt.ParamExprs = append(stmt.ParamExprs, paramExpr)
	}

	return &stmt, nil
}

func (r *Render) renderHelperFile() GoFile {
	code := `type ApiRoutes struct {
	Router *gin.RouterGroup
}

type Enum interface {
	IsValid() bool
}

// TODO: enum array?
func EnumValidation(fl validator.FieldLevel) bool {
	if en, ok := fl.Field().Interface().(Enum); ok {
		return en.IsValid()
	}
	return true
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("enum", EnumValidation)
	}
}

func _wrap(f func(*gin.Context) error) func(*gin.Context) {
	return func(c *gin.Context) {
		if err := f(c); err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}
`
	outputDir := r.getOutputDir()
	packageName := r.getPackage()
	return GoFile{
		Name:    sprintf("%s/zz_helper.go", outputDir),
		Package: packageName,
		Imports: []string{
			"github.com/gin-gonic/gin",
			"github.com/gin-gonic/gin/binding",
			"github.com/go-playground/validator/v10",
		},
		CodeGens: []CodeGen{
			&RawCode{
				code: code,
			},
		},
	}
}
