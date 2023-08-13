package glam

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"net/http"
)

type regexNode struct {
	node
	matcher *regexp.Regexp
}

func (node *regexNode) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")

}

type node struct {
	Name           string                  //16
	ParamChild     *node                   //8
	Children       children                //8
	RegexpChildren children                //8
	StaticChildren children                //8
	Handlers       map[string]http.Handler //8
	Middleware     []Middleware            //24
	NotFound       http.Handler
	regexMatcher   *regexp.Regexp
	parent         *node
}

type Middleware func(next http.Handler) http.Handler

type children map[string]*node

func (c children) String() string {
	keys := make([]string, len(c))
	i := 0
	for k := range c {
		keys[i] = k
		i++
	}
	return fmt.Sprintf("%v", keys)
}

func newChildren() map[string]*node {
	return make(children)
}
func (n *node) notFoundHandler() http.Handler {
	if n.NotFound != nil {
		return n.NotFound
	} else if n.parent != nil && n.parent.notFoundHandler() != nil {
		return n.parent.notFoundHandler()
	} else {
		return defaultNotFoundHandler
	}
}
func newNode(path string, parent *node) *node {
	return &node{
		Name:           path,
		Children:       newChildren(),
		StaticChildren: newChildren(),
		RegexpChildren: newChildren(),
		parent:         parent,
	}
}
func newRoot(router *Router) *node {
	return newNode("", router.root)
}
func (n *node) Type() nodeType {
	return getNodeType(n.Name)
}
func (n *node) applyMiddleware(handler http.Handler) http.Handler {
	if n.Middleware != nil && len(n.Middleware) > 0 {
		for i := len(n.Middleware) - 1; i >= 0; i-- {
			handler = n.Middleware[i](handler)
		}
	}
	return handler
}

func (n *node) depth() int {
	if n.parent == nil {
		return 0
	} else {
		return n.parent.depth() + 1
	}
}
func (n *node) matchPrefix(r *http.Request) (start, end, remainder int) {
	found := false
	count := 0
	path := r.URL.Path

	for i, char := range path {
		if char == '/' {
			count += 1
		} else if count == n.depth()+1 {
			if !found {
				start = i
				end = i
				found = true
			}
			end += 1
		}
	}
	remainder = count - n.depth()
	return
}
func (n *node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start, end, remainder := n.matchPrefix(r)
	if r.URL.Path[start:end] == "" && remainder == 0 || n.Type() == Static {
		if r.Method == http.MethodOptions {
			methods := []string{}
			for method, _ := range n.Handlers {
				methods = append(methods, method)
			}
			w.Header().Set("Allow", strings.Join(methods, ", "))
			w.WriteHeader(http.StatusOK)
			return
		}
		handler, found := n.Handlers[r.Method]
		if !found {
			handler = n.notFoundHandler()
		}
		n.applyMiddleware(handler).ServeHTTP(w, r)
	} else {
		child, inChildren := n.Children[r.URL.Path[start:end]]
		staticChild, inStaticChildren := n.StaticChildren[r.URL.Path[start:end]+"~"]
		if inChildren {
			n.applyMiddleware(child).ServeHTTP(w, r)
		} else if inStaticChildren {
			n.applyMiddleware(staticChild).ServeHTTP(w, r)
		} else {
			for _, child := range n.RegexpChildren {
				if child.regexMatcher.MatchString(r.URL.Path[start:end]) {
					n.applyMiddleware(child).ServeHTTP(w, withRegex(r, child.Name, r.URL.Path[start:end]))
					return
				}
			}
			if n.ParamChild != nil && end != start {
				n.applyMiddleware(n.ParamChild).ServeHTTP(w, withParam(r, n.ParamChild.Name, r.URL.Path[start:end]))
				return
			} else {
				n.applyMiddleware(n.notFoundHandler()).ServeHTTP(w, r)
			}
		}
	}
}

func (n *node) insertHandler(path []string, method string, handler http.Handler) (err error) {
	nodeType := getNextNodeType(path)
	if nodeType == None {
		if _, found := n.Handlers[method]; found {
			return errors.New("method handler already exists for method " + method)
		}
		if n.Handlers == nil {
			n.Handlers = make(map[string]http.Handler)
		}
		n.Handlers[method] = handler
		return
	} else {
		if nodeType == Strict {
			child, inChidren := n.Children[path[0]]
			if !inChidren {
				child = newNode(path[0], n)
				n.Children[path[0]] = child
			}
			return child.insertHandler(path[1:], method, handler)

		} else if nodeType == Static {
			child, inChidren := n.StaticChildren[path[0]]
			if !inChidren {
				child = newNode(path[0], n)
				n.StaticChildren[path[0]] = child
			}
			return child.insertHandler(path[1:], method, handler)

		} else if nodeType == Param {
			if n.ParamChild == nil {
				n.ParamChild = newNode(path[0], n)
				n.ParamChild.insertHandler(path[1:], method, handler)
				return
			} else if n.ParamChild.Name == path[0] {
				n.ParamChild.insertHandler(path[1:], method, handler)
				return
			} else {
				panic("Can't have multiple param prefixes assigned to node")
			}
		} else {
			child, in := n.RegexpChildren[path[0]]
			if !in {
				_, sep := getRegexKeyIndices(path[0])
				var err error
				child = newNode(path[0], n)
				child.regexMatcher, err = regexp.Compile(path[0][sep+1 : len(path[0])-1])
				if err != nil {
					panic("invalid regular expression")
				}
				n.RegexpChildren[path[0]] = child
			}
			return child.insertHandler(path[1:], method, handler)
		}
	}
}

func (n *node) insertMiddleware(path []string, middleware []Middleware) {
	nodeType := getNextNodeType(path)
	if nodeType == None || nodeType == Static {
		n.Middleware = append(n.Middleware, middleware...)
	} else {
		if nodeType == Strict {
			child, inChidren := n.Children[path[0]]
			if !inChidren {
				child = newNode(path[0], n)
				n.Children[path[0]] = child
			}
			child.insertMiddleware(path[1:], middleware)
		} else if nodeType == Static {
			child, inChidren := n.StaticChildren[path[0]]
			if !inChidren {
				child = newNode(path[0], n)
				n.StaticChildren[path[0]] = child
			}
			child.insertMiddleware(path[1:], middleware)

		} else if nodeType == Param {
			if n.ParamChild == nil {
				n.ParamChild = newNode(path[0], n)
				n.ParamChild.insertMiddleware(path[1:], middleware)
			} else if n.ParamChild.Name == path[0] {
				n.ParamChild.insertMiddleware(path[1:], middleware)
			} else {
				panic("Can't have multiple param prefixes assigned to node")
			}
		} else if nodeType == Static {
			child, inChidren := n.Children[path[0][:len(path[0])-1]]
			if !inChidren {
				child = newNode(path[0], n)
				n.Children[path[0]] = child
			}
			child.insertMiddleware(path[1:], middleware)
		} else {
			child, in := n.RegexpChildren[path[0]]
			if !in {
				_, sep := getRegexKeyIndices(path[0])
				var err error
				child = newNode(path[0], n)
				child.regexMatcher, err = regexp.Compile(path[0][sep+1 : len(path[0])-1])
				if err != nil {
					panic("invalid regular expression")
				}
				n.RegexpChildren[path[0]] = child
			}
			child.insertMiddleware(path[1:], middleware)
		}
	}
}

func (n *node) toChildNode(prefix string, parent *node) {
	n.Name = prefix
	n.parent = parent
}

func (n *node) insertNodeAt(prefix string, node *node) {
	prefix = strings.TrimPrefix(prefix, "/")
	_type := getNodeType(prefix)
	node.toChildNode(prefix, n)
	if _type == Strict {
		n.Children[prefix] = node
	} else if _type == Param {
		n.ParamChild = node
	} else {
		n.RegexpChildren[prefix] = node
	}
}
