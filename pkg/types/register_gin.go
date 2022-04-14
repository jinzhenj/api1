package types

import (
	"fmt"
	"path/filepath"
	"strings"
)

type RegisterResource struct {
	Resource          string
	ImportModuleList  []string
	Entries           []RouterEntry
	RegisterFunctions []RouterRegisterFunc
}

func (o RegisterResource) Render() string {
	return o.RenderImportedModules() + `


` + o.RenderRouterEntries() + "\n" + o.RenderRouterRegisterFunctions()
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
		kind := o.Req.Kind.GetKind()
		var req string
		if o.Req.Kind.IsBuiltin() {
			req = fmt.Sprintf(`    var req %s`, kind)
		} else {
			module := filepath.Base(o.Req.Kind.GetModule())
			req = fmt.Sprintf(`    var req %s.%s`, module, kind)
		}
		content = append(content, req)
	}
	if o.Res != nil {
		var res string
		var composeKind string
		if o.Res.Kind.IsBuiltin() {
			composeKind = o.Res.Kind.GetKind()
		} else {
			composeKind = filepath.Base(o.Res.Kind.GetModule()) + "." + o.Res.Kind.GetKind()
		}
		if o.Res.Kind.IsArray() {
			res = fmt.Sprintf(`    res := make([]%s, 0)`, composeKind)
		} else {
			res = fmt.Sprintf(`    var res %s`, composeKind)
		}

		content = append(content, res)
	}
	// render functions
	if o.Res != nil {
		content = append(content, fmt.Sprintf(`    svc.%s(c *gin.Context, &req, &res)`, strings.Title(o.Name)))
	} else {
		content = append(content, fmt.Sprintf(`    svc.%s(c *gin.Context, &req)`, strings.Title(o.Name)))
	}

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
