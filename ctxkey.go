package wenex

import "context"

type ctxKey uint8

const (
	ctxRun ctxKey = iota
)

func GetRun(ctx context.Context) *Run {
	runInterface := ctx.Value(ctxRun)

	if runInterface == nil {
		return nil
	}

	if run, ok := runInterface.(*Run); ok {
		return run
	}

	return nil
}
