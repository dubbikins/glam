package util

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/dubbikins/glam"
	"github.com/dubbikins/glam/logging"
	"github.com/xlab/treeprint"
)

type Tuple struct {
	Node *glam.Node
	Tree treeprint.Tree
	Path string
}

type TreeConfig struct {
	WithColor bool
}

func Tree(router *glam.Router, config *TreeConfig) treeprint.Tree {
	tree := treeprint.New()
	stack := NewStack[*Tuple]()
	stack.Push(&Tuple{
		Node: router.Root(),
		Tree: tree,
		Path: "",
	})
	for stack.Length() > 0 {
		next := stack.Pop()
		for method, handler := range next.Node.Handlers {

			handlerAddr := reflect.ValueOf(handler).Pointer()
			handler, ok := reflect.ValueOf(handler).Interface().(glam.Middleware)
			if ok {
				fmt.Println("middleware")
				fmt.Println(handler)
			}

			file, line := runtime.FuncForPC(handlerAddr).FileLine(handlerAddr)
			branchName := fmt.Sprintf("%s handler", method)
			branchValue := fmt.Sprintf("=> %s:%d", file, line)
			if config.WithColor {
				branchName = logging.Green(branchName)
				branchValue = logging.Gray(branchValue)
			}
			next.Tree.AddMetaBranch(branchName, branchValue)

		}
		for _, middleware := range next.Node.Middleware {
			middlewareAddr := reflect.ValueOf(middleware).Pointer()
			file, line := runtime.FuncForPC(middlewareAddr).FileLine(middlewareAddr)
			branchName := fmt.Sprintf("middleware")
			branchValue := fmt.Sprintf("=> %s:%d", file, line)
			if config.WithColor {
				branchName = logging.Magenta(branchName)
				branchValue = logging.Gray(branchValue)
			}
			next.Tree.AddMetaBranch(branchName, branchValue)
		}
		for _, child := range next.Node.Children {
			addChildBranch(next, child, stack, config.WithColor)
		}
		for _, child := range next.Node.RegexpChildren {
			addChildBranch(next, child, stack, config.WithColor)
		}
		if child := next.Node.ParamChild; child != nil {
			addChildBranch(next, child, stack, config.WithColor)
		}
	}
	return tree
}

func addChildBranch(parentTuple *Tuple, child *glam.Node, stack *Stack[*Tuple], withColor bool) {
	path := child.Name
	if withColor {
		path = logging.Cyan(path)
	}
	branch := parentTuple.Tree.AddBranch(path)
	nodeType := child.Type()
	branchName := "type"
	branchValue := nodeType.ToString()
	if withColor {
		branchName = logging.Yellow(branchName)
		branchValue = logging.Red(branchValue)
	}
	branch.AddMetaBranch(branchName, branchValue)
	stack.Push(&Tuple{
		Node: child,
		Tree: branch,
		Path: path,
	})
}
