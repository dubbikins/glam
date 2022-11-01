package glam

import "testing"

func TestIsParamWhenParam(t *testing.T) {
	nodeType := Param
	if !nodeType.IsParam() {
		t.Fail()
	}
}
func TestIsParamWhenRegexp(t *testing.T) {
	nodeType := Regexp
	if !nodeType.IsParam() {
		t.Fail()
	}
}
func TestIsParamWhenStrict(t *testing.T) {
	nodeType := Strict
	if nodeType.IsParam() {
		t.Fail()
	}
}

func TestMatchesParam(t *testing.T) {
	nodeType := Param
	if !nodeType.Matches("{id}", "test") {
		t.Fail()
	}
}

func TestMatchesRegexp(t *testing.T) {
	nodeType := Regexp
	if !nodeType.Matches("`[tes]`", "test") {
		t.Fail()
	}
}

func TestMatchesStrict(t *testing.T) {
	nodeType := Strict
	if nodeType.Matches("test", "test") {
		t.Fail()
	}
}

func TestIsWrappedBy(t *testing.T) {
	if !isWrappedBy("{", "}", "{test}") {
		t.Fail()
	}
}
func TestIsWrappedByTooShort(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	isWrappedBy("`", "`", "`")
}
func TestIsWrappedByEmptyBody(t *testing.T) {
	if !isWrappedBy("{", "}", "{}") {
		t.Fail()
	}
}
func TestIsWrappedByRegex(t *testing.T) {
	if !isWrappedBy("`", "`", "``") {
		t.Fail()
	}
}

func TestIsWrappedByEmpty(t *testing.T) {
	if !isWrappedBy("`", "`", "``") {
		t.Fail()
	}
}

func TestIsWrappedWithMissingSuffix(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	isWrappedBy("`", "`", "`test}")
}
func TestIsWrappedWithMissingPrefix(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	isWrappedBy("{", "}", "test}")
}

func TestToStringParam(t *testing.T) {
	nodeType := Param
	if nodeType.ToString() != "Param" {
		t.Fail()
	}
}

func TestToStringRegexp(t *testing.T) {
	nodeType := Regexp
	if nodeType.ToString() != "Regexp" {
		t.Fail()
	}
}
func TestToStringStrict(t *testing.T) {
	nodeType := Strict
	if nodeType.ToString() != "Strict" {
		t.Fail()
	}
}

func TestGetNodeTypeStrict(t *testing.T) {

	if getNodeType("test") != Strict {
		t.Fail()
	}
}

func TestGetNodeTypeParam(t *testing.T) {

	if getNodeType("{id}") != Param {
		t.Fail()
	}
}

func TestGetNodeTypeRegexp(t *testing.T) {

	if getNodeType("`[abc]`") != Regexp {
		t.Fail()
	}
}
