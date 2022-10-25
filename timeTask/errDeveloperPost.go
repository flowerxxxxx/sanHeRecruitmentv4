package timeTask

import (
	"os"
	"path/filepath"
)

func getFileSize(filename string) int64 {
	var result int64
	_ = filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}
