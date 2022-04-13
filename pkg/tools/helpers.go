package tools

import (
	"github.com/go-logr/logr"

	"github.com/go-swagger/pkg/parser"
	"github.com/go-swagger/pkg/types"
	"github.com/go-swagger/pkg/utils"
)

func ListHttpHandlers(log logr.Logger, dir string, filter func(string) bool) ([]types.HttpHandler, error) {
	log.V(6).Info("scan dirs", "dir", dir)
	files, err := utils.ListFiles(dir, filter)
	if err != nil {
		log.Error(err, "failed to list files", "dir", dir)
		return nil, err
	}
	log.V(6).Info("build swagger file", "dir", dir, "found api definition files", files)

	// parse api definions from file
	httpHandlers := make([]types.HttpHandler, 0)

	for _, fn := range files {
		tmp, err := parser.ParsrApiDefFile(log, fn)
		if err != nil {
			log.Error(err, "failed to parser api def file", "fileName", fn)
			return nil, err
		}
		httpHandlers = append(httpHandlers, tmp...)
	}
	return httpHandlers, nil
}
