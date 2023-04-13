package util

func Batch[T any, V any](fn func(T) V, ts []T) []V {
	if ts == nil {
		return nil
	}
	res := make([]V, len(ts))
	for i := range ts {
		res = append(res, fn(ts[i]))
	}
	return res
}

func AppendNotNil[T any](values ...*T) []*T {
	res := make([]*T, 0, len(values))
	for i, v := range values {
		if v != nil {
			res = append(res, values[i])
		}
	}
	return res
}
