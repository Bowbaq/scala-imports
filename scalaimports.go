package scalaimports

import (
	"bufio"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/Bowbaq/pool"
)

type Config struct {
	// Imports starting with these prefixes are considered internal, and grouped on top
	Internal []string

	// Imports starting with these prefixes are considered standard library, and grouped at the bottom
	Lang []string

	// Imports prefixed by one of the keys are rewritten to be prefixed by the corresponding value
	Rewrites map[string]string

	// Imports with prefixes in this list are always considered to be used, and never removed
	Ignore []string

	// Imports in this list are spurious and always removed
	Remove []string

	Comparators []Comparator

	MaxLineLength int
}

type Comparator func(string, string) int

var (
	Verbose = false
	// Parallelism = uint(1)
	Parallelism = uint(runtime.NumCPU())

	config = Config{
		Internal: []string{"ai", "common", "dataImport", "emailService", "workflowEngine", "mailgunWebhookService"},
		Lang:     []string{"scala", "java", "javax"},
		Rewrites: map[string]string{
			"_root_.util":        "util",
			"Tap._":              "util.Tap._",
			"MustMatchers._":     "org.scalatest.MustMatchers._",
			"DataPoint.DataType": "models.DataPoint.DataType",
			"action.":            "controllers.action.",
			"helpers.":           "controllers.helpers.",
			"concurrent.":        "scala.concurrent.",
			"collection.":        "scala.collection.",
			"Keys._":             "sbt.Keys._",
		},
		Ignore: []string{
			"scala.collection.JavaConversions",
			"scala.collection.JavaConverters",
			"scala.concurrent.ExecutionContext.Implicits",
			"scala.language.implicitConversions",
			"scala.sys.process",
			"play.api.Play.current",
			"ai.somatix.data.csv.CanBuildFromCsv",
		},
		Remove: []string{
			"import scala.Some",
		},

		MaxLineLength: 110,
	}
)

func init() {
	config.Comparators = []Comparator{compareInternal, compareLang, lexicographical}
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
		debug("Formatting", root)
		scalaFiles = []string{root}
	}

	cleanFiles(scalaFiles, Parallelism)
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

func cleanFiles(paths []string, parallelism uint) {
	var wg sync.WaitGroup
	wg.Add(len(paths))

	p := pool.NewPool(parallelism, func(id uint, payload interface{}) interface{} {
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
