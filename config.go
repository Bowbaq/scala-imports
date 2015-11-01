package scalaimports

import "runtime"

var (
	compareInternal comparator
	compareLang     comparator

	config = NewConfig()
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

	MaxLineLength int

	Parallelism uint

	Verbose bool

	comparators []comparator
}

func NewConfig() Config {
	c := Config{
		Internal:      []string{},
		Lang:          []string{"scala", "java", "javax"},
		Rewrites:      make(map[string]string),
		Ignore:        []string{},
		Remove:        []string{},
		MaxLineLength: 110,
		Parallelism:   uint(runtime.NumCPU()),
	}

	compareInternal = reverse(comparePrefix(c.Internal))
	compareLang = comparePrefix(c.Lang)
	c.comparators = []comparator{compareInternal, compareLang, lexicographical}

	return c
}

func SetConfig(c Config) {
	compareInternal = reverse(comparePrefix(c.Internal))
	compareLang = comparePrefix(c.Lang)
	c.comparators = []comparator{compareInternal, compareLang, lexicographical}

	if c.Parallelism == 0 {
		c.Parallelism = 1
	}

	config = c
}
