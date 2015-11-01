package scalaimports

import (
	"os"
	"path/filepath"
	"strings"
)

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	check(err)
	return fileInfo.IsDir()
}

func findScalaFiles(root string) (paths []string) {
	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".scala") && !strings.HasSuffix(path, "package.scala") {
			paths = append(paths, path)
		}
		return nil
	})

	return paths
}
