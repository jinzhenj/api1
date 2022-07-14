package all

import (
	"io/ioutil"
	"path/filepath"

	"github.com/api1/pkg/utils"
)

const (
	defaultFilePerm = 0644
)

type CodeFile struct {
	Name    string
	Content string
}

func (f *CodeFile) WriteFile() error {
	dir := filepath.Dir(f.Name)
	if err := utils.MayCreateDir(dir); err != nil {
		return err
	}
	return ioutil.WriteFile(f.Name,
		[]byte(f.Content), defaultFilePerm)
}
