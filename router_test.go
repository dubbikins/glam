package glam

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetRouter() *Router {
	r := NewRouter()
	r.Get("/posts/{test}", func(w http.ResponseWriter, r *http.Request) {
		test, _ := GetURLParam(r, "test")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(test))
	})
	r.Get("/really/long/nested/path/that/keeps/going/on/forever", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return r
}
func GetRouterWithMiddleware() *Router {
	handler := GetRouter()
	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	handler.Use(mw)
	return handler
}
func BenchmarkNewRouter(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		b.ReportAllocs()
		b.ResetTimer()
		for p.Next() {
			NewRouter()
		}
	})
}
func BenchmarkGetRouter(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {

		b.ReportAllocs()
		b.ResetTimer()
		for p.Next() {
			r := NewRouter()
			r.Get("/posts/{test}", func(w http.ResponseWriter, r *http.Request) {
				test, _ := GetURLParam(r, "test")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test))
			})
			r.Get("/really/long/nested/path/that/keeps/going/on/forever", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		}
	})
}
func BenchmarkHandleGetWithParam(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		r := NewRouter()
		b.ReportAllocs()
		b.ResetTimer()
		for p.Next() {
			r.Get("/posts/{test}", func(w http.ResponseWriter, r *http.Request) {
				test, _ := GetURLParam(r, "test")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test))
			})
		}
	})

}
func BenchmarkHandle(b *testing.B) {

	b.RunParallel(func(p *testing.PB) {
		rtr := NewRouter()
		b.ReportAllocs()
		b.ResetTimer()
		for p.Next() {
			rtr.Handle("/test", http.MethodGet, func(w http.ResponseWriter, r *http.Request) {

			})

		}

	})
}
func BenchmarkRouterGet(b *testing.B) {
	handler := GetRouter()
	b.RunParallel(func(pb *testing.PB) {
		r, _ := http.NewRequest("GET", "/posts", nil)
		w := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			handler.ServeHTTP(w, r)
		}
	})
}
func BenchmarkRouterGetLongPath(b *testing.B) {
	handler := GetRouter()
	b.RunParallel(func(pb *testing.PB) {
		r, _ := http.NewRequest("GET", "/really/long/nested/path/that/keeps/going/on/forever", nil)
		w := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			handler.ServeHTTP(w, r)
		}
	})
}
func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := httptest.NewRecorder()
	u := r.URL
	path := u.Path
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	//for i := 0; i < b.N; i++ {
	u.Path = path
	u.RawQuery = rq
	router.ServeHTTP(w, r)
	//}
}
func BenchmarkRouterGetWithMiddleware(b *testing.B) {
	handler := GetRouter()

	r, _ := http.NewRequest("GET", "/posts/params", nil)
	benchRequest(b, handler, r)

}

func BenchmarkRouterGetWith3Middleware(b *testing.B) {
	handler := GetRouterWithMiddleware()
	b.RunParallel(func(pb *testing.PB) {
		r, _ := http.NewRequest("GET", "/posts/param", nil)
		w := httptest.NewRecorder()

		b.ReportAllocs()
		b.ResetTimer()
		for pb.Next() {
			handler.ServeHTTP(w, r)
		}
	})
}

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router.NotFoundHandler == nil {
		t.Fatal("New Router Not Found Handler Not Set")
	}
	if router.root == nil || router.root.Name != "" {
		t.Fatal("New Router root should not be nil and emtpy")
	}
}

func TestRoot(t *testing.T) {
	router := NewRouter()

	if router.root != router.Root() {
		t.Fatal("New Router Root() returned incorrect value")
	}
}

func TestDefaultNotFoundHandler(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusNotFound,
	}
	router.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
}

func TestOverrideNotFoundHandler(t *testing.T) {
	router := NewRouter()

	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusNotImplemented,
	}
	expected_body := "Opps! Page Not Found"
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(expected_body))
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestHandleWithRegisteredSlashSuffic(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}
func TestHandleWithSlashSuffic(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}
func TestHandleWithBothSlashSuffic(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test/", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Get("/test/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestOptionsStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodOptions, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Options("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}
func TestGetStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}
func TestGetMultiStrictPathsRegistered(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.Get("/tester", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Get("/tes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestGetMultiLongStrictPathsRegistered(t *testing.T) {
	router := NewRouter()
	path := "/test/long/path"
	r, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.WriteHeader(expected.StatusCode)
	})
	router.Get("/test/different/path", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.WriteHeader(http.StatusBadRequest)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
}
func TestGetMultiMethodsRegistered(t *testing.T) {
	router := NewRouter()
	path := "/test/long/path"
	r, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.WriteHeader(expected.StatusCode)
	})
	router.Put(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.Post(path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
}

func TestGetWithMiddleware(t *testing.T) {
	router := NewRouter()
	path := "/test/long/path"
	r, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("key", "value")
			next.ServeHTTP(w, r)
		})
	})
	router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if resp.Header.Get("key") != "value" {
		t.Fatalf("HTTP GET Request Failed: Middleware did not get applied")
	}
}

func TestPutStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodPut, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Put("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestPostStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodPost, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Post("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestPatchStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodPatch, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Patch("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestHeadStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodHead, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Head("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestConnectStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodConnect, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Connect("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestDeleteStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodDelete, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Delete("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestTraceStrictPath(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodTrace, "/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Trace("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected_body))
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestRoutes(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/routes/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "hello world"
	router.Routes("routes", func(r *Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expected_body))
			w.WriteHeader(expected.StatusCode)
		})
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestRoutesParam(t *testing.T) {
	router := NewRouter()
	r, err := http.NewRequest(http.MethodGet, "/routes/test", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	expected_body := "routes_id"
	router.Routes("{id}", func(r *Router) {
		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			id, _ := GetURLParam(r, "id")
			w.Write([]byte(id + "_id"))
			w.WriteHeader(expected.StatusCode)
		})
	})
	router.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != expected.StatusCode {
		t.Fatalf("HTTP GET Request Failed: Expected status code %d but was %d", expected.StatusCode, resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Fatalf("HTTP GET Request Failed: Expected body %s but was %s", expected_body, body)
	}
}

func TestGetGlamContextWhenNotSet(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/routes/test", nil)
	if err != nil {
		t.Fail()
	}
	ctx := GetGlamContext(r)
	if ctx != nil {
		t.Fatal("expected glam context to be nil when unset by router")
	}

}

func TestGetGlamContextWhenSet(t *testing.T) {
	router := NewRouter()
	var ctx *glamContext
	r, err := http.NewRequest(http.MethodGet, "/test/123", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx = GetGlamContext(r)
		w.WriteHeader(expected.StatusCode)
	})
	router.ServeHTTP(w, r)
	if ctx == nil {
		t.Fatal("Glam Context should not be nil")
	}
}

func TestGetURLParamWhenNotSet(t *testing.T) {
	router := NewRouter()

	r, err := http.NewRequest(http.MethodGet, "/test/123", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		name, found := GetURLParam(r, "name")
		w.WriteHeader(expected.StatusCode)
		if name != "" || found == true {
			t.Fatal("URL Param should return false when param not set")
		}
	})
	router.ServeHTTP(w, r)

}

func TestGetURLParamWhenSet(t *testing.T) {
	router := NewRouter()

	r, err := http.NewRequest(http.MethodGet, "/test/123", nil)
	if err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	expected := http.Response{
		StatusCode: http.StatusOK,
	}
	router.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, found := GetURLParam(r, "id")
		w.WriteHeader(expected.StatusCode)
		if id != "123" || found == false {
			t.Fatal("URL Param should return false when param not set")
		}
	})
	router.ServeHTTP(w, r)

}

func TestGetURLParamWhenNoContextSet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/test/123", nil)
	id, found := GetURLParam(r, "id")
	if id != "" || found == true {
		t.Fatal("URL Param should return false when param not set")
	}

}
