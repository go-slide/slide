package ferry

import (
	"io/ioutil"
	"mime"
	"path"
	"path/filepath"
)

func getAllPaths(dirPath string, paths *[]string) error {
	fileInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			nextPath := path.Join(dirPath, file.Name())
			if err := getAllPaths(nextPath, paths); err != nil {
				return err
			}
		} else {
			extension := file.Name()
			nextPath := path.Join(dirPath, extension)
			*paths = append(*paths, nextPath)
		}
	}
	return nil
}

func getFileContentType(path string) (string, error) {
	m := mime.TypeByExtension(filepath.Ext(path))
	return m, nil
}
