package utils

import (
	"io/ioutil"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathExistsOrCreate(path string) error {
	exist, err := PathExists(path)
	if err != nil {
		return err
	}

	if !exist {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func EmptyDir(path string) bool {
	dir, _ := ioutil.ReadDir(path)
	return len(dir) == 0
}
