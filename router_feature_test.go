package glam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type response struct {
	Header     map[string]string `json:"header"`
	Body       string            `json:"body"`
	StatusCode int               `json:"statusCode"`
}
type request struct {
	Header http.Header `json:"header"`
	Body   io.Reader   `json:"body"`
	Path   string      `json:"path"`
	Method string      `json:"method"`
}

func (r *request) UnmarshalJSON(data []byte) error {

	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	var ok bool
	if methodRaw, found := v["method"]; found {
		r.Method, ok = methodRaw.(string)
		if !ok {
			return errors.New("cannot unmarshal method")
		}
	}
	if bodyRaw, found := v["body"]; found {
		body, ok := bodyRaw.(string)
		if !ok {
			return errors.New("cannot unmarshal method")
		}
		r.Body = strings.NewReader(body)
	}
	r.Header = http.Header{}
	if headerRaw, found := v["header"]; found {
		header, ok := headerRaw.(map[string]interface{})
		if !ok {
			return errors.New("cannot unmarshal header")
		}
		for k, v := range header {
			r.Header.Add(k, v.(string))
		}
	}
	r.Path, ok = v["path"].(string)
	if !ok {
		return errors.New("cannot unmarshal path")
	}
	return nil
}

type rootRouterContextKey struct{}
type requestContextKey struct{}
type responseContextKey struct{}

func thereIsARootRouter(ctx context.Context) (context.Context, error) {
	router := NewRouter()
	return context.WithValue(ctx, rootRouterContextKey{}, router), nil
}
func theRouterHasAHandlerForPathThatRespondsWith(ctx context.Context, method, path string, respDoc *godog.DocString) error {
	router, ok := ctx.Value(rootRouterContextKey{}).(*Router)
	if !ok {
		return errors.New("A router is not available")
	}
	expectedResponse := &response{}
	json.Unmarshal([]byte(respDoc.Content), expectedResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		for k, v := range expectedResponse.Header {
			w.Header().Add(k, v)
		}
		w.WriteHeader(expectedResponse.StatusCode)
		w.Write([]byte(expectedResponse.Body))
	}
	switch method {
	case http.MethodGet:
		router.Get(path, handler)
	case http.MethodPut:
		router.Put(path, handler)
	case http.MethodPost:
		router.Post(path, handler)
	case http.MethodDelete:
		router.Delete(path, handler)
	case http.MethodPatch:
		router.Patch(path, handler)
	case http.MethodHead:
		router.Head(path, handler)
	case http.MethodConnect:
		router.Connect(path, handler)
	// case http.MethodOptions:
	// 	router.Options(path, handler)
	case http.MethodTrace:
		router.Trace(path, handler)
	}

	return nil
}

func theRouterResponds(ctx context.Context) (context.Context, error) {
	router, ok := ctx.Value(rootRouterContextKey{}).(*Router)
	if !ok {
		return ctx, errors.New("Root router is not available")
	}
	r, ok := ctx.Value(requestContextKey{}).(*http.Request)
	if !ok {
		return ctx, errors.New("request is not available")
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	return context.WithValue(ctx, responseContextKey{}, w.Result()), nil
}

func aRequestIsMade(ctx context.Context, body *godog.DocString) (context.Context, error) {
	req := &request{}
	err := json.Unmarshal([]byte(body.Content), req)
	if err != nil {
		panic(err.Error())
	}

	r := httptest.NewRequest(req.Method, req.Path, req.Body)
	r.Header = req.Header
	return context.WithValue(ctx, requestContextKey{}, r), nil
}

func theResponseShouldMatch(ctx context.Context, respDoc *godog.DocString) error {
	resp, ok := ctx.Value(responseContextKey{}).(*http.Response)
	expected := &response{}
	json.Unmarshal([]byte(respDoc.Content), expected)
	if !ok {
		return errors.New("response is not available")
	}

	haveBody, _ := ioutil.ReadAll(resp.Body)

	haveStatusCode := resp.StatusCode
	if want := expected.StatusCode; want != haveStatusCode {
		return errors.New(fmt.Sprintf("expected a %d, instead got: %d", want, haveStatusCode))
	}
	if want := expected.Body; want != string(haveBody) {
		return errors.New(fmt.Sprintf("expected a %s, instead got: %s", want, haveBody))
	}

	for header, want := range expected.Header {
		if have := resp.Header.Get(header); have != want {
			return errors.New(fmt.Sprintf("expected a %s:%s, instead got: %s:%s", header, want, header, have))
		}
	}
	return nil
}
func theRouterHasMiddlewareThatAddsTheFollowingToTheResponseHeader(ctx context.Context, rawHeaders *godog.DocString) (err error) {

	parent, ok := ctx.Value(rootRouterContextKey{}).(*Router)
	if !ok {
		err = errors.New("A router is not available")
		return
	}
	headers := make(map[string]interface{})
	err = json.Unmarshal([]byte(rawHeaders.Content), &headers)
	if err != nil {
		return
	}
	parent.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				if value, ok := v.(string); !ok {
					err = errors.New("header unmarshalling error")
				} else {
					w.Header().Add(k, value)
				}
			}
			next.ServeHTTP(w, r)
		})
	})
	return nil
}
func isASubrouterWithAHandlerHandlerForPathThatRespondsWithAndIsMountedAtPath(ctx context.Context, method, prefix, path string, respDoc *godog.DocString) error {
	parent, ok := ctx.Value(rootRouterContextKey{}).(*Router)
	if !ok {
		return errors.New("A router is not available")
	}
	router := NewRouter()
	expectedResponse := &response{}
	json.Unmarshal([]byte(respDoc.Content), expectedResponse)
	handler := func(w http.ResponseWriter, r *http.Request) {
		for k, v := range expectedResponse.Header {
			w.Header().Add(k, v)
		}
		w.WriteHeader(expectedResponse.StatusCode)
		w.Write([]byte(expectedResponse.Body))
	}
	switch method {
	case http.MethodGet:
		router.Get(path, handler)
	case http.MethodPut:
		router.Put(path, handler)
	case http.MethodPost:
		router.Post(path, handler)
	case http.MethodDelete:
		router.Delete(path, handler)
	case http.MethodPatch:
		router.Patch(path, handler)
	case http.MethodHead:
		router.Head(path, handler)
	case http.MethodConnect:
		router.Connect(path, handler)
	// case http.MethodOptions:
	// 	router.Options(path, handler)
	case http.MethodTrace:
		router.Trace(path, handler)
	}
	parent.Mount(prefix, router)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is a root Router$`, thereIsARootRouter)
	ctx.Step(`^the router has a "([^"]*)" handler for path "([^"]*)" that responds with:$`, theRouterHasAHandlerForPathThatRespondsWith)
	ctx.Step(`^a subrouter with a "([^"]*)" handler mounted at "([^"]*)" for path "([^"]*)" that responds with:$`, isASubrouterWithAHandlerHandlerForPathThatRespondsWithAndIsMountedAtPath)
	ctx.Step(`^the router has middleware that adds the following to the response header:$`, theRouterHasMiddlewareThatAddsTheFollowingToTheResponseHeader)
	ctx.Step(`^a request is made:$`, aRequestIsMade)
	ctx.Step(`^the Router responds$`, theRouterResponds)
	ctx.Step(`^the response should match:$`, theResponseShouldMatch)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
