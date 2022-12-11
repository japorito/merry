package stockings

type Stack[T any] struct {
	data []T
}

func (stack *Stack[T]) Size() int {
	return len(stack.data)
}

func (stack *Stack[T]) Push(items ...T) *Stack[T] {
	stack.data = append(stack.data, items...)

	return stack
}

func (stack *Stack[T]) Pop() T {
	return stack.Top(1)[0]
}

func (stack *Stack[T]) Top(count int) []T {
	lastidx := len(stack.data) - count
	items := stack.data[lastidx:]
	stack.data = stack.data[:lastidx]

	return items
}

func (stack *Stack[T]) Peek() T {
	return stack.PeekTop(1)[0]
}

func (stack *Stack[T]) PeekTop(count int) []T {
	lastidx := len(stack.data) - count
	return stack.data[lastidx:]
}
