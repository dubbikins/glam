package glam

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTree(t *testing.T) {
	t.Log("TestTree")
	router := NewRouter()
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})
	tree := router.Tree() // just make sure it doesn't panic
	fmt.Println(tree.String())
	if tree != nil {
		branch := tree.FindByValue("test")
		branch = branch.FindByValue("[type]")
		fmt.Println(branch.String())
		t.Error("Tree() failed")
	}

}
