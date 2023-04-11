package syncer

func ToAnySlice[T any](ts []T) []any {
	if ts == nil {
		return nil
	}
	res := make([]any, len(ts))
	for i := range ts {
		res = append(res, ts[i])
	}
	return res
}
