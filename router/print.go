package router

import (
	"fmt"
	"github.com/xlab/treeprint"
	"github.com/dubbikins/glam/logging"
	"reflect"
	"runtime"
)

type Tuple struct {
	Node *Node
	Tree treeprint.Tree
	Path string
}

type TreeConfig struct {
	WithColor bool
}
func (router *Router) Tree(config *TreeConfig) treeprint.Tree {
	tree := treeprint.New()
	stack := NewStack[*Tuple]()
	stack.Push(&Tuple{
		Node: router.root,
		Tree: tree,
		Path: "",
	})
	for stack.Length() > 0 {
		next := stack.Pop()
		for method, handler := range next.Node.Handlers {
			handlerAddr := reflect.ValueOf(handler).Pointer()
			file, line := runtime.FuncForPC(handlerAddr).FileLine(handlerAddr)
			branchName := fmt.Sprintf("%s handler", method)
			branchValue :=  fmt.Sprintf("=> %s:%d", file, line)
			if config.WithColor{
				branchName = logging.Green(branchName)
				branchValue = logging.Gray(branchValue)
			}
			next.Tree.AddMetaBranch(branchName, branchValue)
		}
		for _,middleware := range next.Node.Middleware {
			middlewareAddr := reflect.ValueOf(middleware).Pointer()
			file, line := runtime.FuncForPC(middlewareAddr).FileLine(middlewareAddr)
			branchName := fmt.Sprintf("middleware")
			branchValue :=  fmt.Sprintf("=> %s:%d", file, line)
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

func addChildBranch (parentTuple *Tuple, child *Node, stack *Stack[*Tuple], withColor bool) {
	path := child.Name
	if withColor {
		path = logging.Cyan(path)
	}
	branch := parentTuple.Tree.AddBranch(path)
	nodeType := getNodeType(child.Name)
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
