package utils

import (
	"os"
	"path"
	"path/filepath"
)

func CreatePath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	return err
}
func GetCurrentPath() string {
	dir := filepath.Dir(os.Args[0])
	dir, _ = filepath.Abs(dir)
	return dir
}
func JoinCurrentPath(p string) string {
	return path.Join(GetCurrentPath(), p)
}
