package activity

import (
	"context"

	"go.temporal.io/sdk/activity"
)

func PostApproveActivity(ctx context.Context, data PostApproveActionPayload) error {
	logger := activity.GetLogger(ctx)
	logger.Info("PostApproveActivity", "data", data)
	return nil
}
