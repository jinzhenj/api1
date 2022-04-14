package types

import (
	"fmt"
	"strings"

	"github.com/go-swagger/pkg/utils"
)

type RegisterResource struct {
	Resource          string
	ImportModuleList  []string
	Entries           []RouterEntry
	RegisterFunctions []RouterRegisterFunc
}

func (o RegisterResource) Render() string {
	return o.RenderImportedModules() + `


` + o.RenderRouterEntries() + o.RenderRouterRegisterFunctions()
}

func (o RegisterResource) RenderImportedModules() string {
	content := make([]string, 0)
	content = append(content, "package router\n\n")
	content = append(content, "import (")
	for _, module := range o.ImportModuleList {
		content = append(content, fmt.Sprintf(`    "%s"`, module))
	}

	content = append(content, ")")
	return strings.Join(content, "\n")
}

func (o RegisterResource) RenderRouterEntries() string {
	head := fmt.Sprintf(`func Register%sApi(r *gin.RouterGroup) {`, strings.Title(o.Resource))
	tail := "}\n"
	content := make([]string, 0)
	content = append(content, head)
	for _, entry := range o.Entries {
		content = append(content, "    "+entry.Render("r", "svc"))
	}

	content = append(content, tail)
	return strings.Join(content, "\n")
}

func (o RegisterResource) RenderRouterRegisterFunctions() string {
	content := make([]string, len(o.RegisterFunctions))
	for idx, val := range o.RegisterFunctions {
		content[idx] = val.Render()
	}
	return strings.Join(content, "\n\n")
}

type RouterRegisterFunc struct {
	Name string
	Res  *HandlerBodyParams
	Req  *HandlerBodyParams
}

func (o RouterRegisterFunc) Render() string {
	content := make([]string, 0)
	header := fmt.Sprintf(`func %s(c *gin.Context) {`, o.Name)
	content = append(content, header)

	// declare request and response variable
	if o.Req != nil {
		m := strings.Split(o.Req.Name, ".")
		l := len(m)
		req := fmt.Sprintf(`    var req %s.%s`, m[l-2], m[l-1])
		content = append(content, req)
	}
	if o.Res != nil {
		ns := strings.Trim(o.Res.Name, "[]")
		m := strings.Split(ns, ".")
		l := len(m)
		var res string
		if strings.HasPrefix(o.Res.Name, "[]") {
			if utils.IsGoBuiltinTypes(ns) {
				res = fmt.Sprintf(`    res := make(%s, 0)`, o.Res.Name)
			} else {
				res = fmt.Sprintf(`    res := make([]%s.%s, 0)`, m[l-2], m[l-1])
			}
		} else {
			if utils.IsGoBuiltinTypes(ns) {
				res = fmt.Sprintf(`    var res %s`, o.Res.Name)
			} else {
				res = fmt.Sprintf(`    var res %s.%s`, m[l-2], m[l-1])
			}
		}
		content = append(content, res)
	}
	// render functions
	if o.Res != nil {
		content = append(content, fmt.Sprintf(`    svc.%s(c *gin.Context, &req, &res)`, strings.Title(o.Name)))
	} else {
		content = append(content, fmt.Sprintf(`    svc.%s(c *gin.Context, &req)`, strings.Title(o.Name)))
	}

	//TODO
	content = append(content, "}")
	return strings.Join(content, "\n")
}

type RouterEntry struct {
	Method  string
	Path    string
	Handler string
}

func (o RouterEntry) Render(routerName, svc string) string {
	return fmt.Sprintf(`%s.%s, %s)`, routerName, o.RenderPathInfo(), o.Handler)
}

func (o RouterEntry) RenderPathInfo() string {
	return fmt.Sprintf(`%s("%s"`, strings.ToUpper(o.Method), o.Path)
}
