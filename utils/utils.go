package utils

func Filter[T any](ss []T, test func(T) bool) []T {
	ret := make([]T, 0)
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func FindFirst[T any](ss *[]T, test func(T) bool) *T {
	for _, s := range *ss {
		if test(s) {
			return &s
		}
	}
	return nil
}
