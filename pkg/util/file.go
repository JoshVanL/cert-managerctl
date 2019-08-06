package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteFile(path string, b []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0744); err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, perm)
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat("/path/to/whatever")
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
