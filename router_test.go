package glam

import (
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
