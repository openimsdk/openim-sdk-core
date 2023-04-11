package syncer

func AnySlice[T any](ts []T) []any {
	rs := make([]any, len(ts))
	for i := range ts {
		rs = append(rs, rs[i])
	}
	return rs
}
