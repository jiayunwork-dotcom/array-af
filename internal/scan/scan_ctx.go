package scan

import "context"

func abortScanContext() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
