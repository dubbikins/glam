package glam

import (
	"testing"
)

func BenchmarkApplyMiddleware(b *testing.B) {
	// router := NewRouter()
	// node := newNode("test", router.root)
	// mw := func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		next.ServeHTTP(w, r)
	// 	})
	// }
	// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// })
	// node.Middleware = append(node.Middleware, mw, mw)
	// b.RunParallel(func(pb *testing.PB) {

	// 	b.ReportAllocs()
	// 	b.ResetTimer()
	// 	for pb.Next() {
	// 		node.applyMiddleware(handler)
	// 	}
	// })
}
