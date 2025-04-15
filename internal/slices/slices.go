package slices

import "sort"

func IndexFunc[S ~[]E, E any](s S, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

func ContainsFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	return IndexFunc(s, f) >= 0
}

func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) {
	sort.Slice(x, func(i, j int) bool {
		return cmp(x[i], x[j]) < 0
	})
}

func Contains[S ~[]E, E comparable](s S, v E) bool {
	return Index(s, v) >= 0
}

func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	for i, v := range s {
		if del(v) {
			j := i
			for i++; i < len(s); i++ {
				v = s[i]
				if !del(v) {
					s[j] = v
					j++
				}
			}
			return s[:j]
		}
	}
	return s
}

func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j]

	return append(s[:i], s[j:]...)
}

func IsSortedFunc[S ~[]E, E any](x S, cmp func(a, b E) int) bool {
	for i := len(x) - 1; i > 0; i-- {
		if cmp(x[i], x[i-1]) < 0 {
			return false
		}
	}
	return true
}

func CompareFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, cmp func(E1, E2) int) int {
	for i, v1 := range s1 {
		if i >= len(s2) {
			return +1
		}
		v2 := s2[i]
		if c := cmp(v1, v2); c != 0 {
			return c
		}
	}
	if len(s1) < len(s2) {
		return -1
	}
	return 0
}

func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

func Clone[S ~[]E, E any](s S) S {
	if s == nil {
		return nil
	}
	return append(S([]E{}), s...)
}
