package main

func filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func findFirst[T any](ss *[]T, test func(T) bool) *T {
	for _, s := range *ss {
		if test(s) {
			return &s
		}
	}
	return nil
}
