package utils

import "context"

func CancelAndReBuildCtx(
	buildFunc func(ctx context.Context) context.Context,
	cancel ...context.CancelFunc,
) context.Context {
	for _, cf := range cancel {
		cf()
	}
	return buildFunc(nil)
}
