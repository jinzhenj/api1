package tools

import (
	"github.com/go-logr/logr"
)

func NewUpdateApiInterface(log logr.Logger, apiDir, module string) *UpdateApiInterface {
	return &UpdateApiInterface{log: log, apiDir: apiDir, module: module}
}

type UpdateApiInterface struct {
	log    logr.Logger
	apiDir string
	module string
}
