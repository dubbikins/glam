package glam

import (
	"fmt"
	"regexp"
	"strings"
)

type NodeType uint8

const (
	Strict NodeType = iota
	Param
	Regexp
)

func (n NodeType) IsParam() bool {
	return n == Param || n == Regexp
}

func (n NodeType) Matches(pattern, value string) bool {
	if n == Param && value != "" {
		return true
	} else if n == Regexp {
		pattern = pattern[1 : len(pattern)-1]
		r, _ := regexp.Compile(pattern)
		return r.MatchString(value)
	}
	return false
}

func (n NodeType) ToString() string {

	var nodeType string
	if n == Strict {
		nodeType = "Strict"
	} else if n == Param {
		nodeType = "Param"
	} else if n == Regexp {
		nodeType = "Regexp"
	}
	return nodeType
}

func isWrappedBy(prefix, suffix, pattern string) bool {
	hasPrefix := strings.HasPrefix(pattern, prefix)
	hasSuffix := strings.HasSuffix(pattern, suffix)
	if len(pattern) < 2 {
		panic(fmt.Sprintf("Route param is missing is invalid length"))
	}
	if hasPrefix && !hasSuffix {
		panic(fmt.Sprintf("Route param is missing '%s'", suffix))
	} else if !hasPrefix && hasSuffix {
		panic(fmt.Sprintf("Route param is missing '%s'", prefix))
	}
	return hasPrefix && hasSuffix
}

func getNodeType(prefix string) NodeType {
	if isWrappedBy("{", "}", prefix) {
		return Param
	} else if isWrappedBy("`", "`", prefix) {
		return Regexp
	} else {
		return Strict
	}
}
