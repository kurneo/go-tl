package helper

import (
	"strings"
)

type PathPreFixer struct {
	prefix    string
	separator string
}

func (p PathPreFixer) PrefixPath(path string) string {
	return p.prefix + strings.TrimLeft(path, "\\/")
}

func (p PathPreFixer) StripPrefix(path string) string {
	return path[len(p.prefix):]
}

func (p PathPreFixer) StripDirectoryPrefix(path string) string {
	return strings.TrimRight(p.StripPrefix(path), "\\/")
}

func (p PathPreFixer) PrefixDirectoryPath(path string) string {
	prefixedPath := p.PrefixPath(strings.TrimRight(path, "\\/"))

	if prefixedPath == "" || prefixedPath[len(prefixedPath)-1:] == p.separator {
		return prefixedPath
	}

	return prefixedPath + p.separator
}

func (p PathPreFixer) StripTrailingSeparator(path string) string {
	return strings.TrimRight(path, p.separator)
}

func NewPreFixer(prefix, separator string) PathPreFixer {
	preFixer := PathPreFixer{
		prefix:    strings.TrimRight(prefix, "\\/"),
		separator: separator,
	}

	if preFixer.prefix != "" || preFixer.prefix == preFixer.separator {
		preFixer.prefix = preFixer.prefix + preFixer.separator
	}
	return preFixer
}
