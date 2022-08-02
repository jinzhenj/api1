package golang

import (
	"fmt"
	"strings"

	"github.com/jinzhenj/api1/pkg/utils"
)

var (
	sprintf = fmt.Sprintf
	indent  = utils.Indent
)

type CodeFile struct {
	Name    string
	Content string
}

func (t *GoTypeDef) Code() string {
	return sprintf("type %s %s\n", t.Name, t.Type.Code())
}

func CodeComments(c []string) string {
	code := ""
	for _, line := range c {
		code += sprintf("// %s\n", line)
	}
	return code
}

func (o *GoEnumOption) Code() string {
	code := ""
	code += CodeComments(o.Comments)
	// o.name := b.Name + pascalCase(o.Name)
	if o.Value.IntVal != nil {
		code += sprintf("%s %s = %d\n", o.Name, o.TypeName, *o.Value.IntVal)
	} else {
		code += sprintf("%s %s = \"%s\"\n", o.Name, o.TypeName, *o.Value.StrVal)
	}
	return code
}

func CodeEnumOptions(options []GoEnumOption) string {
	code := "\n"
	if len(options) > 0 {
		code += "const (\n"
		for _, o := range options {
			code += indent(o.Code())
		}
		code += ")\n"
	}
	return code
}

func (b *GoEnumCodeBlock) Code() string {
	code := ""
	baseType := b.BaseType.Code()
	code += CodeComments(b.Comments)
	code += sprintf("type %s %s\n", b.Name, baseType)
	code += CodeEnumOptions(b.Options)
	code += "\n"
	code += sprintf("func (o %s) IsValid() bool {\n", b.Name)
	if len(b.Options) > 0 {
		code += indent(sprintf("switch %s(o) {\n", baseType))
		code += indent("case\n")
		for i, opt := range b.Options {
			code += indent(indent(sprintf("%s(%s)", baseType, opt.Name)))
			if i < len(b.Options)-1 {
				code += ",\n"
			} else {
				code += ":\n"
			}
		}
		code += indent(indent("return true\n"))
		code += indent("}\n")
	}
	code += indent("return false\n")
	code += "}\n"
	return code
}

func CodeStructFields(fields []GoStructField) string {
	code := ""
	for _, field := range fields {
		code += CodeComments(field.Comments)
		code += sprintf("%s %s", field.Name, field.Type.Code())
		tags := CodeTags(field.Tags)
		if len(tags) > 0 {
			code += " " + tags
		}
		code += "\n"
	}
	return code
}

func (t *GoType) Code() string {
	if t == nil ||
		(len(t.Name) == 0 &&
			t.ItemType == nil) {
		return ""
	}
	var code string
	if len(t.Name) > 0 {
		code = t.Name
	} else if t.KeyType == nil {
		code = "[]" + t.ItemType.Code()
	} else {
		code = sprintf("map[%s]%s",
			t.KeyType.Code(), t.ItemType.Code())
	}
	if t.IsPointer {
		code = "*" + code
	}
	return code
}

func CodeTags(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}
	var a []string
	for k, v := range tags {
		a = append(a, sprintf("%s:\"%s\"", k, v))
	}
	return sprintf("`%s`", strings.Join(a, " "))
}

func (s *GoStructType) Code() string {
	code := ""
	code += CodeComments(s.Comments)
	code += sprintf("type %s struct {\n", s.Name)
	code += indent(CodeStructFields(s.Fields))
	code += "}\n"
	return code
}

func (p *GoParam) Code() string {
	return sprintf("%s %s", p.Name, p.Type.Code())
}

func CodeParams(params []GoParam) string {
	var codes []string
	for _, p := range params {
		codes = append(codes, p.Code())
	}
	return strings.Join(codes, ", ")
}

func (fun *GoFunction) Code() string {
	code := ""
	code += CodeComments(fun.Comments)
	if !fun.InIface {
		code += "func "
	}
	if fun.Receiver != nil {
		code += sprintf("(%s) ", CodeParams([]GoParam{*fun.Receiver}))
	}
	code += sprintf("%s(%s)", fun.Name, CodeParams(fun.Params))
	if len(fun.RetTypes) > 0 {
		var types []string
		for _, t := range fun.RetTypes {
			types = append(types, t.Code())
		}
		typesCode := strings.Join(types, ", ")
		if len(types) == 1 {
			code += sprintf(" %s", typesCode)
		} else {
			code += sprintf(" (%s)", typesCode)
		}
	}
	if fun.InIface {
		code += "\n"
		return code
	}
	code += " {\n"
	for _, stmt := range fun.Statements {
		code += "\n"
		code += indent(stmt.Code())
	}
	code += "}\n"
	return code
}

func (iface *GoInterface) Code() string {
	code := ""
	code += CodeComments(iface.Comments)
	code += sprintf("type %s interface {\n", iface.Name)
	for _, f := range iface.Functions {
		code += "\n"
		code += indent(f.Code())
	}
	code += "}\n"
	return code
}

func (route *RouteStatement) Code() string {
	code := sprintf("_r.Router.%s(\"%s\", ", strings.ToUpper(route.Method), route.Path)

	for _, middleware := range route.Middlewares {
		code += middleware + ", "
	}

	code += "_wrap(func(_c *gin.Context) error {\n"

	if len(route.PathParams) > 0 {
		block := "\n"
		block += "var _path struct {\n"
		block += indent(CodeStructFields(route.PathParams))
		block += "}\n"
		block += "if _err := _c.ShouldBindUri(&_path); _err != nil {\n"
		block += indent("return _err\n")
		block += "}\n"
		code += indent(block)
	}

	if len(route.QueryParams) > 0 {
		block := "\n"
		block += "var _query struct {\n"
		block += indent(CodeStructFields(route.QueryParams))
		block += "}\n"
		block += "if _err := _c.ShouldBindQuery(&_query); _err != nil {\n"
		block += indent("return _err\n")
		block += "}\n"
		code += indent(block)
	}

	if route.BodyParam != nil {
		param := route.BodyParam
		block := "\n"
		block += sprintf("var %s %s\n", param.Name, param.Type.Code())
		block += sprintf("if _err := _c.ShouldBindJSON(&%s); _err != nil {\n", param.Name)
		block += "  return _err\n"
		block += "}\n"
		code += indent(block)
	}

	block := "\n"
	if route.HasRet {
		block += "_ret, "
	}
	block += sprintf("_err := _o.%s(%s)\n", route.Name,
		strings.Join(append(route.ParamExprs, "_c"), ", "))
	block += "if _err != nil {\n"
	block += "  return _err\n"
	block += "}\n"
	if route.HasRet {
		block += "_c.Set(\"ret\", _ret)\n"
	}
	block += "return nil\n"

	code += indent(block)
	code += "}))\n"
	return code
}

func (file *GoFile) Code() string {
	code := "// Code generated by api1; DO NOT EDIT.\n"
	code += sprintf("package %s\n", file.Package)

	if len(file.Imports) > 0 {
		code += "\n"
		code += "import (\n"
		for _, imp := range file.Imports {
			code += indent(sprintf("\"%s\"\n", imp))
		}
		code += ")\n"
	}

	for _, codeGen := range file.CodeGens {
		code += "\n"
		code += codeGen.Code()
	}
	return code
}

func (r *RawCode) Code() string {
	return r.code
}
