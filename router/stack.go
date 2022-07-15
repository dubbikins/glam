package router


type Stack[T any] struct {
	items[]T
}
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}
func (s *Stack[T]) Pop() T {
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}
func (s *Stack[T]) Length() int {
	return len(s.items)
}
// func (stack *Stack[T]) Print() {

// 	fmt.Printf("[%s]", strings.Join(fmt.Sprintf("%v"), ","))
// }



func NewStack[T any] () *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

