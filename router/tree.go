package router

import (
	"errors"
	"net/http"
	"strings"
	"fmt"
	"github.com/xlab/treeprint"
)

type Node struct {
	Name string
	FullPath string
	Children Children
	ParamChild *Node
	RegexpChildren Children
	Type NodeType
	tree treeprint.Tree
	Middleware []Middleware
	Handlers map[string]http.Handler
}

type Middleware func(next http.Handler) http.Handler

type Children map[string]*Node

func NewChildren() map[string]*Node {
	return make(Children)
}

func NewNode(path string, nodeType NodeType) *Node {
	return &Node{
		Name: path,
		Children: NewChildren(),
		RegexpChildren: NewChildren(),
		Middleware: make([]Middleware, 0),
		Type: nodeType,
		Handlers: make(map[string]http.Handler),
	}
}
func NewRoot() *Node {
	return NewNode("",Strict)
}

func (n *Node) traverse(path string) *Node {
	prefix, remainder := splitPath(path)
	nodeType := getNodeType(prefix)
	if prefix == "" && remainder == "" {
		return n
	} else {
		
		if nodeType == Strict {
			child, inChidren := n.Children[prefix]
			if !inChidren {
				child = NewNode(prefix, nodeType)
				child.FullPath = join(n.FullPath, prefix)
				n.Children[prefix] = child
			}
			return child.traverse(remainder)
		}else if nodeType == Param {
			if n.ParamChild == nil {
				n.ParamChild = NewNode(prefix, nodeType)
				n.ParamChild.FullPath = join(n.FullPath, prefix)
				return n.ParamChild.traverse(remainder)
			}else if n.ParamChild.Name == prefix {
				return n.ParamChild.traverse(remainder)
			} else {
				panic("Can't have multiple param prefixes assigned to node")
			}
		} else {
			child, in := n.RegexpChildren[prefix]
			if !in {
				child = NewNode(prefix, nodeType)
				child.FullPath = join(n.FullPath, prefix)
				n.RegexpChildren[prefix] = child
			}
			return child.traverse(remainder)
		}
	}
}

type Callback func (*Node, string)

func (n *Node) traverseMatchApply(path, previousName string, callback Callback) (*Node, error) {
	prefix, remainder := splitPath(path)
	callback(n, previousName)
	if strings.HasPrefix(remainder, "/"){
		return nil, errors.New("Match Error")
	} else if prefix == "" && remainder == "" {
		return n, nil
	} else {
		child, inChidren := n.Children[prefix]
		if  inChidren {
			return child.traverseMatchApply(remainder, prefix, callback)
		} else {
			for _, child := range n.RegexpChildren {
				if child.Type.Matches(child.Name, prefix) {
					return child.traverseMatchApply(remainder, prefix, callback)
				}
			}
			if n.ParamChild != nil && prefix != "" {
				return n.ParamChild.traverseMatchApply(remainder, prefix, callback)
			}else {
				return nil, errors.New("Match Error")
			}
		}
	}
}

func join(prefix, name string) string {
	return strings.Join([]string{prefix, name}, "/")
}

func splitPath(path string) (prefix string, remainder string) {
	if strings.HasPrefix(path, "/") {
		path = path [1:]
	}
	splits := strings.SplitN(path, "/", 2)
	prefix = splits[0]
	if len(splits) > 1 {
		remainder = splits[1]
	}
	return 
}

func (n *Node) toLeafNode(prefix string) {
	n.Name = prefix
	n.Type = getNodeType(prefix)
}

func (n *Node) InsertNodeAt(path string, newNode *Node) {
	prefix, remainder := splitPath(path)
	if remainder != "" {
		panic("can't add node to nested path")
	}else {
		nodeType := getNodeType(prefix)
		newNode.toLeafNode(prefix)
		newNode.FullPath = join(n.FullPath, prefix)
		if nodeType == Strict {
			n.Children[prefix] = newNode
		} else if nodeType == Param {
			n.ParamChild = newNode
		} else {
			n.RegexpChildren[prefix] = newNode
		}
	}
}

func (n *Node) InsertRouteHandler(path, method string, handler http.HandlerFunc) {
	leaf := n.traverse(path)
	if _, in := leaf.Handlers[method]; in {
		panic("Can't overwrite an existing handler")
	}
	leaf.Handlers[method] = handler
}

func (n *Node) InsertRouteMiddleware(path string , middleware Middleware) {
	leaf := n.traverse(path)
	leaf.Middleware = append(leaf.Middleware, middleware)
}

func (n *Node) MatchRequest(r *http.Request) ([]Middleware, http.Handler, map[string]string, error) {
	var err error 
	middleware := make([]Middleware, 0)
	params := make(map[string]string)
	var onTraversal Callback = func (node *Node, prefix string) {
		for _, mw := range node.Middleware {
			middleware = append(middleware, mw)
		}
		nodeType := getNodeType(node.Name)
		if nodeType == Param || nodeType == Regexp {
			params[node.Name[1:len(node.Name)-1]] = prefix
		}
	}
	leaf, err := n.traverseMatchApply(r.URL.EscapedPath(), "", onTraversal)
	if err != nil {
		return middleware, nil, params, err
	}
	handler, in := leaf.Handlers[r.Method]
	if !in {
		err = errors.New("No Routes Matched the request")
	}
	return middleware, handler, params, err
}

func Merge(p1 map[string]string, p2 map[string]string) map[string]string {
	for k, v := range p2 {
		p1[k]= v
	}
	return p1
}

func (n *Node) PrintTree() {
	fmt.Printf(n.tree.String())
}
func (n *Node) String() string {
	return n.tree.String()
}