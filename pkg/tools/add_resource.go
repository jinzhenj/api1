package tools

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

var (
	ginModules = []string{"github.com/gin-gonic/gin"}
)

func (o UpdateApiInterface) Do() error {
	groupHandler, err := o.GetGroupHandlers()
	if err != nil {
		return err
	}

	for resource, arr := range groupHandler {
		if err := o.update(resource, arr); err != nil {
			return err
		}

	}
	return nil
}

func (o UpdateApiInterface) GetGroupHandlers() (map[string][]types.HttpHandler, error) {
	records, err := o.getHttpHandlers()
	if err != nil {
		return nil, err
	}
	namesMap := make(map[string]bool)
	for _, record := range records {
		namesMap[record.Resource] = true
	}

	ret := make(map[string][]types.HttpHandler, 0)

	for resource, _ := range namesMap {
		arr := make([]types.HttpHandler, 0)
		for idx, record := range records {
			if record.Resource == resource {
				arr = append(arr, records[idx])
			}
		}
		ret[resource] = arr
	}

	return ret, nil
}

func (o UpdateApiInterface) getHttpHandlers() ([]types.HttpHandler, error) {
	filter := func(s string) bool {
		return strings.HasSuffix(s, ".api")
	}
	return ListHttpHandlers(o.log, o.apiDir, filter)
}

func (o UpdateApiInterface) update(resource string, arr []types.HttpHandler) error {
	res := types.RegisterResource{Resource: resource}

	// module related
	modules := make([]string, 0)
	modules = append(modules, ginModules...)
	modulesMap := make(map[string]bool)

	modulesMap[fmt.Sprintf("%s/pkg/svc", o.module)] = true
	for _, h := range arr {
		if h.Req != nil {
			modulesMap[fmt.Sprintf("%s/%s", o.module, filepath.Dir(strings.ReplaceAll(h.Req.Name, ".", "/")))] = true
		}
		if h.Res != nil {
			ns := strings.Trim(h.Res.Name, "[]")
			if !utils.IsGoBuiltinTypes(ns) {
				modulesMap[fmt.Sprintf("%s/%s", o.module, filepath.Dir(strings.ReplaceAll(h.Res.Name, ".", "/")))] = true
			}
		}
	}
	for key := range modulesMap {
		modules = append(modules, key)
	}
	res.ImportModuleList = modules

	// router entry related
	entries := make([]types.RouterEntry, len(arr))
	for idx, h := range arr {
		entries[idx] = types.RouterEntry{
			Method:  h.Method,
			Path:    h.Endpoint,
			Handler: h.Name,
		}
	}
	res.Entries = entries

	// register functions related
	registerFunctions := make([]types.RouterRegisterFunc, len(arr))
	for idx, h := range arr {
		registerFunctions[idx] = types.RouterRegisterFunc{Name: h.Name, Res: h.Res, Req: h.Req}
	}
	res.RegisterFunctions = registerFunctions

	// render file
	if err := utils.MayCreateDir("pkg/router"); err != nil {
		return err
	}
	fileName := filepath.Join("pkg/router", fmt.Sprintf("register_%s.go", strings.ToLower(resource)))
	return ioutil.WriteFile(fileName, []byte(res.Render()), 0644)
}
