package scalaimports

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/Bowbaq/pool"
)

func init() {
	SetConfig(DefaultConfig)
}

func Format(root string) {
	var scalaFiles []string

	if isDir(root) {
		debug("Formatting all scala files in", root)
		scalaFiles = findScalaFiles(root)

		if len(scalaFiles) == 0 {
			debug(root, "contains no scala files")
			return
		}
		debug("Found scala files in", root)
		for _, path := range scalaFiles {
			debug("  ", path)
		}
	} else {
		if isScalaFile(root) {
			debug("Formatting", root)
			scalaFiles = []string{root}
		}
	}

	cleanFiles(scalaFiles)
}

func ParseFile(path string) (*ScalaFile, error) {
	var scalaFile = NewScalaFile(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isCodeSection := false

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "package"):
			scalaFile.PackageLines.Add(line)

		case strings.HasPrefix(line, "import"):
			scalaFile.ImportLines.Add(line)

		case len(strings.TrimSpace(line)) == 0 && !isCodeSection:
			// Ignore blank lines in the import section
			continue

		default:
			if !isCodeSection {
				isCodeSection = true
			}
			scalaFile.CodeLines.Add(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return scalaFile, nil
}

func cleanFiles(paths []string) {
	var wg sync.WaitGroup
	wg.Add(len(paths))

	p := pool.NewPool(config.Parallelism, func(id uint, payload interface{}) interface{} {
		path := payload.(string)

		debug("Worker", id, "processing", path)

		file, err := ParseFile(path)
		if err != nil {
			debug("Worker", id, "error parsing", path, "-", err)
			return err
		}

		err = file.Rewrite()
		if err != nil {
			debug("Worker", id, "error rewriting", path, "-", err)
			return err
		}

		wg.Done()

		return nil
	})

	for _, path := range paths {
		p.Submit(pool.NewJob(path))
	}

	wg.Wait()
}
