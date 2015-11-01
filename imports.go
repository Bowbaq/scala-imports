package scalaimports

import (
	"sort"
	"strings"
)

type imports []string

func (i *imports) Add(line string) {
	*i = append(*i, strings.TrimSpace(strings.TrimPrefix(line, "import")))
}

func (i *imports) Expand() {
	expanded := make(imports, 0)

	for _, imp := range *i {
		if openBrace := strings.Index(imp, "{"); openBrace != -1 && !strings.Contains(imp, "=> _") {
			prefix, group := imp[:openBrace], imp[openBrace+1:len(imp)-1]

			for _, suffix := range strings.Split(group, ",") {
				if strings.Contains(suffix, "=>") {
					expanded = append(expanded, prefix+"{"+strings.TrimSpace(suffix)+"}")
				} else {
					expanded = append(expanded, prefix+strings.TrimSpace(suffix))
				}
			}
		} else {
			expanded = append(expanded, imp)
		}
	}

	*i = expanded
}

func (i *imports) Rewrite(from, to string) {
	for x, imp := range *i {
		if strings.HasPrefix(imp, from) {
			(*i)[x] = strings.Replace(imp, from, to, 1)
		}
	}
}

func (i *imports) Remove(needle string) {
	for x, imp := range *i {
		if imp == needle {
			*i = append((*i)[:x], (*i)[x+1:]...)
		}
	}
}

func (i *imports) Sort() {
	sort.Sort(i)
}

func (i *imports) Len() int {
	return len((*i))
}

func (i *imports) Less(x, y int) bool {
	a, b := (*i)[x], (*i)[y]

	for _, comp := range config.Comparators {
		switch comp(a, b) {
		case -1:
			return true
		case 1:
			return false
		}
	}

	return false
}

func (i *imports) Swap(x, y int) {
	(*i)[x], (*i)[y] = (*i)[y], (*i)[x]
}

func (i *imports) Dedup() {
	var (
		unique []string
		seen   = make(map[string]bool)
	)

	for _, imp := range *i {
		if _, in := seen[imp]; !in {
			unique = append(unique, imp)
			seen[imp] = true
		}
	}

	(*i) = unique
}

func (i *imports) RemoveUnused(file *ScalaFile) {
	var used []string

	for _, imp := range *i {
		if file.IsImportUsed(imp) {
			used = append(used, imp)
		}
	}

	(*i) = used
}

func (i *imports) Group() {
	var (
		prefixes []string
		groups   = make(map[string][]string)
	)

	// Group imports
	for _, imp := range *i {
		i := strings.LastIndex(imp, ".")
		prefix, suffix := imp[:i+1], imp[i+1:]

		if suffixes, ok := groups[prefix]; ok {
			groups[prefix] = append(suffixes, suffix)
		} else {
			groups[prefix] = []string{suffix}
		}

		if !contains(prefixes, prefix) {
			prefixes = append(prefixes, prefix)
		}
	}

	var (
		grouped  []string
		previous string
	)

	for _, prefix := range prefixes {
		// Insert newline to separate groups if needed
		if previous != "" && compareInternal(previous, prefix) == -1 {
			grouped = append(grouped, "")
			previous = ""
		}

		if previous != "" && compareLang(previous, prefix) == -1 {
			grouped = append(grouped, "")
		}

		previous = prefix

		// Compact group
		group := groups[prefix]
		if len(group) == 1 {
			grouped = append(grouped, prefix+group[0])
		} else {
			items, renames := make([]string, 0), make([]string, 0)
			for _, g := range group {
				if strings.Contains(g, "=>") {
					renames = append(renames, g[1:len(g)-1])
				} else {
					items = append(items, g)
				}
			}

			items = append(renames, items...)

			formatted := prefix + "{" + strings.Join(items, ", ") + "}"
			if len(formatted) > config.MaxLineLength || contains(items, "_") {
				formatted = prefix + "_"
			}
			grouped = append(grouped, formatted)
		}
	}

	(*i) = grouped
}
