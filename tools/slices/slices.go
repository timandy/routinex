package slices

func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	result := make(S, len(s))
	copy(result, s)
	return result
}

func Equal[S ~[]E, E comparable](s1, s2 S) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func Filter[S ~[]E, E any](s S, predicate func(E) bool) S {
	result := make(S, 0, len(s))
	for _, item := range s {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
