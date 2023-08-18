package glam

import (
	"strings"
)

const StaticIdentifier string = "*"
const ParamPrefix string = "{"
const ParamSuffix string = "}"
const RegexpPrefix string = "("
const RegexpSuffix string = ")"

type nodeType uint8

const (
	Strict nodeType = iota
	Param
	Regexp
	Static
	None
)

func getRegexKeyIndices(key string) (start, sepIndex int) {
	if len(key) < 2 {
		panic("invalid pattern definition: invalid length")
	}
	start = 1
	sepIndex = strings.IndexRune(key, ':')
	if sepIndex == -1 {
		panic("invalid regular expression pattern definition")
	}
	return
}

func (n nodeType) ToString() string {
	switch n {
	case Strict:
		return "Strict"
	case Param:
		return "Param"
	case Regexp:
		return "Regexp"
	case Static:
		return "Static"
	case None:
		return "None"
	}
	return ""
}

func isWrappedBy(prefix, suffix, pattern string) bool {
	hasPrefix := strings.HasPrefix(pattern, prefix)
	hasSuffix := strings.HasSuffix(pattern, suffix)
	if len(pattern) < 2 {
		return false
	}
	if (hasPrefix && !hasSuffix) || (!hasPrefix && hasSuffix) {
		panic("Invalid param/regex pattern definition")
	}
	return hasPrefix && hasSuffix
}

func getNodeType(name string) nodeType {

	if strings.HasSuffix(name, StaticIdentifier) {
		return Static
	} else if isWrappedBy(ParamPrefix, ParamSuffix, name) {
		return Param
	} else if isWrappedBy(RegexpPrefix, RegexpSuffix, name) {
		return Regexp
	} else {
		return Strict
	}
}
func getNextNodeType(path []string) nodeType {
	if len(path) == 0 {
		return None
	}
	return getNodeType(path[0])
}
