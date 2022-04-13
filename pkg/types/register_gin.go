package types

import (
	"fmt"
	"strings"
)

type RegisterResource struct {
	Resource         string
	ImportModuleList []string
	Entries          []RouterEntry
}

func (o RegisterResource) Render() string {
	return o.RenderImportedModules() + `


` + o.RenderRegisterFunctions()
}

func (o RegisterResource) RenderImportedModules() string {
	content := make([]string, 0)
	content = append(content, "package router")
	content = append(content, "import (")
	for _, module := range o.ImportModuleList {
		content = append(content, fmt.Sprintf(`    "%s"`, module))
	}

	content = append(content, ")")
	return strings.Join(content, "\n")
}

func (o RegisterResource) RenderRegisterFunctions() string {
	head := fmt.Sprintf(`func Register%sApi(r *gin.RouterGroup) {`, strings.Title(o.Resource))
	tail := "}\n"
	content := make([]string, 0)
	content = append(content, head)
	for _, entry := range o.Entries {
		if entry.IsNeedToUpdate() {
			content = append(content, "    "+entry.Render("r", "svc"))
		}
	}

	content = append(content, tail)
	return strings.Join(content, "\n")
}

type RouterEntry struct {
	Method         string
	Path           string
	Handler        string
	noNeedToUpdate bool
}

func (o RouterEntry) Render(routerName, svc string) string {
	return fmt.Sprintf(`%s.%s, %s.%s)`, routerName, o.RenderPathInfo(), svc, strings.Title(o.Handler))
}

func (o RouterEntry) RenderPathInfo() string {
	return fmt.Sprintf(`%s("%s"`, strings.ToUpper(o.Method), o.Path)
}

func (o RouterEntry) SetNoNeedToUpdate() {
	o.noNeedToUpdate = true
}

func (o RouterEntry) IsNeedToUpdate() bool {
	return !o.noNeedToUpdate
}
