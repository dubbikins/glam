package glam

import (
	"bytes"
	"fmt"

	"net/http"
)

type Node struct {
	Name            string                  //16
	Router          *Router                 //8
	ParamChild      *Node                   //8
	Children        Children                //8
	RegexpChildren  Children                //8
	Handlers        map[string]http.Handler //8
	Middleware      []Middleware            //24
	NotFoundHandler http.Handler
}

type Middleware func(next http.Handler) http.Handler

type Children map[string]*Node

func NewChildren() map[string]*Node {
	return make(Children)
}

func NewNode(path string, router *Router) *Node {
	return &Node{
		Router:         router,
		Name:           path,
		Children:       NewChildren(),
		RegexpChildren: NewChildren(),
	}
}
func NewRoot(router *Router) *Node {
	return NewNode("", router)
}
func (n *Node) Type() NodeType {
	return getNodeType(n.Name)
}
func (n *Node) ApplyMiddleware(handler http.Handler) http.Handler {

	if n.Middleware != nil && len(n.Middleware) > 0 {

		for i := len(n.Middleware) - 1; i >= 0; i-- {
			fmt.Println("Applying middleware")
			handler = n.Middleware[i](handler)
		}
	}
	return handler
}

func (n *Node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	prefix := trimRequestPath(r)
	if prefix == "" && r.URL.Path == "" {
		handler, found := n.Handlers[r.Method]
		if !found {
			handler = n.Router.NotFoundHandler
		}
		if n.Middleware != nil && len(n.Middleware) > 0 {
			for i := len(n.Middleware) - 1; i >= 0; i-- {
				handler = n.Middleware[i](handler)
			}
		}
		handler.ServeHTTP(w, r)
	} else {
		child, inChildren := n.Children[prefix]
		if inChildren {
			n.ApplyMiddleware(child).ServeHTTP(w, r)
		} else {
			for _, child := range n.RegexpChildren {
				if child.Type().Matches(child.Name, prefix) {
					n.ApplyMiddleware(child).ServeHTTP(w, requestWithURLParam(r, child.Name, prefix))
					return
				}
			}
			if n.ParamChild != nil && prefix != "" {
				n.ApplyMiddleware(n.ParamChild).ServeHTTP(w, requestWithURLParam(r, n.ParamChild.Name, prefix))

			} else {
				n.Router.NotFoundHandler.ServeHTTP(w, r)
			}
		}
	}
}
func (n *Node) insertHandler(path []string, method string, handler http.Handler) {
	if len(path) == 0 {
		if n.Handlers == nil {
			n.Handlers = make(map[string]http.Handler)
		}
		if method == "NOTFOUND" {
			n.NotFoundHandler = handler
		}
		n.Handlers[method] = handler
	} else {
		nodeType := getNodeType(path[0])
		if nodeType == Strict {
			child, inChidren := n.Children[path[0]]
			if !inChidren {
				child = NewNode(path[0], n.Router)
				n.Children[path[0]] = child
			}
			child.insertHandler(path[1:], method, handler)
		} else if nodeType == Param {
			if n.ParamChild == nil {
				n.ParamChild = NewNode(path[0], n.Router)
				n.ParamChild.insertHandler(path[1:], method, handler)
			} else if n.ParamChild.Name == path[0] {
				n.ParamChild.insertHandler(path[1:], method, handler)
			} else {
				panic("Can't have multiple param prefixes assigned to node")
			}
		} else {
			child, in := n.RegexpChildren[path[0]]
			if !in {
				child = NewNode(path[0], n.Router)
				n.RegexpChildren[path[0]] = child
			}
			child.insertHandler(path[1:], method, handler)
		}
	}
}

func (n *Node) insertMiddleware(path []string, middleware []Middleware) {
	if len(path) == 0 {
		n.Middleware = middleware
	} else {
		fmt.Println(path[0])
		if n.Type() == Strict {
			fmt.Println("strict")
			child, inChidren := n.Children[path[0]]
			if !inChidren {
				child = NewNode(path[0], n.Router)
				n.Children[path[0]] = child
			}

			child.insertMiddleware(path[1:], middleware)
		} else if n.Type() == Param {
			if n.ParamChild == nil {
				n.ParamChild = NewNode(path[0], n.Router)
				n.ParamChild.insertMiddleware(path[1:], middleware)
			} else if n.ParamChild.Name == path[0] {
				n.ParamChild.insertMiddleware(path[1:], middleware)
			} else {
				panic("Can't have multiple param prefixes assigned to node")
			}
		} else {
			child, in := n.RegexpChildren[path[0]]
			if !in {
				child = NewNode(path[0], n.Router)
				n.RegexpChildren[path[0]] = child
			}
			child.insertMiddleware(path[1:], middleware)
		}
	}
}

func trimRequestPath(r *http.Request) string {
	buff := bytes.NewBufferString("")
	for len(r.URL.Path) > 0 && r.URL.Path[0] == '/' {
		r.URL.Path = r.URL.Path[1:]
	}
	for len(r.URL.Path) > 0 && r.URL.Path[0] != '/' {
		buff.WriteByte(r.URL.Path[0])
		r.URL.Path = r.URL.Path[1:]
	}
	return buff.String()
}

func (n *Node) toLeafNode(prefix string) {
	n.Name = prefix
}

func (n *Node) InsertNodeAt(prefix string, newNode *Node) {
	nodeType := getNodeType(prefix)
	newNode.toLeafNode(prefix)
	if nodeType == Strict {
		n.Children[prefix] = newNode
	} else if nodeType == Param {
		n.ParamChild = newNode
	} else {
		n.RegexpChildren[prefix] = newNode
	}

}
