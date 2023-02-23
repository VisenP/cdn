package main

func filter[T any](ss []T, test func(T) bool) []T {
	ret := make([]T, 0)
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func findFirst[T any](ss *[]T, test func(T) bool) *T {
	for _, s := range *ss {
		if test(s) {
			return &s
		}
	}
	return nil
}
