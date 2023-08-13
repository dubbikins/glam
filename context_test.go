package glam

import (
	"net/http"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Log("TestQuery")
	r, _ := http.NewRequest("GET", "http://localhost:8080/?test=1", nil)
	value, ok := Query(r, "test")
	if !ok || value != "1" {
		t.Error("Query() failed")
	}
}

func TestGetParam(t *testing.T) {
	t.Log("TestGetParam")
	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	value, ok := GetParam(r, "test")
	if ok || value != "" {
		t.Error("GetParam() failed")
	}
}

func TestWithParam(t *testing.T) {
	t.Log("TestWithParam")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	r = withParam(r, "{test}", "1")
	value, ok := GetParam(r, "test")
	if !ok || value != "1" {
		t.Error("withParam() failed")
	}
}

func TestWithRegex(t *testing.T) {
	t.Log("TestWithRegex")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	r = withRegex(r, "(test:.*)", "1")
	value, ok := GetParam(r, "test")
	if !ok || value != "1" {
		t.Error("withRegex() failed")
	}
}

func TestWithParamAndRegex(t *testing.T) {
	t.Log("TestWithParamAndRegex")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	r = withParam(r, "{test}", "1")
	r = withRegex(r, "(test2:.*)", "2")
	value, ok := GetParam(r, "test")
	if !ok || value != "1" {
		t.Error("withParam() failed")
	}
	value, ok = GetParam(r, "test2")
	if !ok || value != "2" {
		t.Error("withRegex() failed")
	}
}

func TestWithParamAndRegexAndQuery(t *testing.T) {
	t.Log("TestWithParamAndRegexAndQuery")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test?test2=2", nil)
	r = withParam(r, "{test}", "1")
	r = withRegex(r, "(test2:.*)", "2")
	value, ok := GetParam(r, "test")
	if !ok || value != "1" {
		t.Error("withParam() failed")
	}
	value, ok = GetParam(r, "test2")
	if !ok || value != "2" {
		t.Error("withRegex() failed")
	}
	value, ok = Query(r, "test2")
	if !ok || value != "2" {
		t.Error("Query() failed")
	}
}

func TestDuplicateParam(t *testing.T) {
	t.Log("TestDuplicateParam")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	r = withParam(r, "{test}", "1")
	r = withParam(r, "{test}", "2")
	value, ok := GetParam(r, "test")
	if !ok || value != "2" {
		t.Error("withParam() failed")
	}
}

func TestDuplicateRegex(t *testing.T) {
	t.Log("TestDuplicateRegex")
	r, _ := http.NewRequest("GET", "http://localhost:8080/test", nil)
	r = withRegex(r, "(test:.*)", "1")
	r = withRegex(r, "(test:.*)", "2")
	value, ok := GetParam(r, "test")
	if !ok || value != "2" {
		t.Error("withRegex() failed")
	}
}
