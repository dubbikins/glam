package glam

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/dubbikins/glam/logging"
	"github.com/dubbikins/glam/util"
	"github.com/xlab/treeprint"
)

type Tuple struct {
	Node *Router
	Tree treeprint.Tree
	Path string
}

type TreeConfig struct {
	WithColor bool
	WithDepth bool
}

func (router *Router) Tree(optsFunc ...func(config *TreeConfig)) treeprint.Tree {
	config := &TreeConfig{}
	for _, f := range optsFunc {
		f(config)
	}
	tree := treeprint.New()
	stack := util.NewStack[*Tuple]()
	stack.Push(&Tuple{
		Node: router,
		Tree: tree,
		Path: "",
	})
	for stack.Length() > 0 {
		next := stack.Pop()
		if config.WithDepth {
			next.Tree.AddMetaBranch("Depth", fmt.Sprintf("%d", next.Node.depth()))
		}
		if next.Node.notFound != nil {
			nfhandlerAddr := reflect.ValueOf(next.Node.notFoundHandler()).Pointer()
			file, line := runtime.FuncForPC(nfhandlerAddr).FileLine(nfhandlerAddr)
			branchName := fmt.Sprintf("NOT_FOUND Handler")
			branchValue := fmt.Sprintf("=> %s:%d", file, line)
			if config.WithColor {
				branchName = logging.Green(branchName)
				branchValue = logging.Gray(branchValue)
			}
			next.Tree.AddMetaBranch(branchName, branchValue)
		}

		for method, handler := range next.Node.Handlers {

			handlerAddr := reflect.ValueOf(handler).Pointer()

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
		for _, child := range next.Node.StaticChildren {

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

func addChildBranch(parentTuple *Tuple, child *Router, stack *util.Stack[*Tuple], withColor bool) {
	path := child.Name
	if withColor {
		path = logging.Cyan(path)
	}
	branch := parentTuple.Tree.AddBranch("/" + path)
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
