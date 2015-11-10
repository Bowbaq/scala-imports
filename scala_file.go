package scalaimports

import (
	"os"
	"strings"
)

type ScalaFile struct {
	Path string

	PackageLines packageLines
	ImportLines  imports
	CodeLines    code
}

func NewScalaFile(path string) *ScalaFile {
	return &ScalaFile{
		Path: path,
	}
}

func (f *ScalaFile) FormattedImports() string {
	f.ImportLines.Expand()
	for from, to := range config.Rewrites {
		f.ImportLines.Rewrite(from, to)
	}
	for _, r := range config.Remove {
		f.ImportLines.Remove(r)
	}
	f.ImportLines.Sort()
	f.ImportLines.Dedup()
	f.ImportLines.RemoveUnused(f)
	f.ImportLines.Group()

	for i, line := range f.ImportLines {
		if len(line) > 0 {
			f.ImportLines[i] = "import " + line
		}
	}

	return strings.Join(f.ImportLines, "\n")
}

// IsImportUsed returns false if the canonical import is not referenced in this file. A canonical import has
// the form `my.package.name.MyImport` or `my.package.name.{MyImport => MyAlias}`. Passing anything else to
// this method will result in undefined behavior.
func (f *ScalaFile) IsImportUsed(canonical string) bool {
	// Pull out the last part of the import
	i := strings.LastIndex(canonical, ".")
	suffix := canonical[i+1:]

	// Wildcard imports are presumed to be used
	if strings.Contains(suffix, "_") {
		return true
	}

	if strings.Contains(suffix, "=>") {
		parts := strings.Split(suffix[1:len(suffix)-1], "=>")
		alias := strings.TrimSpace(parts[1])

		if f.CodeLines.Contains(alias) {
			return true
		}
	}

	if f.CodeLines.Contains(suffix) {
		return true
	}

	for _, ignore := range config.Ignore {
		if strings.HasPrefix(canonical, ignore) {
			return true
		}
	}

	return false
}

func (f *ScalaFile) Rewrite() error {
	file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write package line + newline
	if len(f.PackageLines) > 0 {
		file.WriteString("package " + strings.Join(f.PackageLines, ".") + "\n\n")
	}

	// Write clean imports + newline
	if len(f.ImportLines) > 0 {
		file.WriteString(f.FormattedImports() + "\n\n")
	}

	// Write back rest of file, end with newline
	file.WriteString(strings.Join(f.CodeLines, "\n") + "\n")

	return nil
}

type packageLines []string

func (l *packageLines) Add(line string) {
	*l = append(*l, strings.TrimSpace(line))
}

type code []string

func (c *code) Add(line string) {
	*c = append(*c, line)
}

func (c *code) Contains(needle string) bool {
	for _, line := range *c {
		if strings.Contains(line, needle) {
			return true
		}
	}

	return false
}
