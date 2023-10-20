package collection

func Transform[T any, E any](src []T, f func(T) E) []E {
	seq := make([]E, len(src))
	for index, each := range src {
		seq[index] = f(each)
	}
	return seq
}
