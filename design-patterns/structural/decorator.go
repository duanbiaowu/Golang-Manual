package structural

type Operator func(int) int

func OpDecorator(fn Operator) Operator {
	return func(n int) int {
		result := fn(n)
		return result
	}
}
