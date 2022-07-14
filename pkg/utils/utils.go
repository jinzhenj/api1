package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ListFiles(dir string, filter func(string) bool) ([]string, error) {
	ret := make([]string, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			subFiles, err := ListFiles(filepath.Join(dir, f.Name()), filter)
			if err != nil {
				return nil, err
			}
			ret = append(ret, subFiles...)

		} else {
			if filter(f.Name()) {
				ret = append(ret, filepath.Join(dir, f.Name()))
			}
		}
	}
	return ret, nil
}

func MayCreateDir(path string) error {
	if path == "" || path == "." || path == "./" {
		return nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0777); err != nil {
			return err
		}
		return nil
	}
	return err
}

func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
