package utils

import (
        "os"
)
func CreatePath(path string)error  {
	_, err := os.Stat(path)
	if err == nil {
			return nil
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm);err != nil {
			return  err
		}
		return nil
	}
	return err
}