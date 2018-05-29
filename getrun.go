package wenex

import "context"

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
