package activities

import (
	"context"
	"go.temporal.io/sdk/activity"
)

func ActivityA1(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA1...")
	return "A1 done", nil
}

func ActivityA2(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA2...")
	return "A2 done", nil
}

func ActivityA2_5(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA2_5...")
	return "A2.5 done", nil
}

func ActivityA3(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA3...")
	return "A3 done", nil
}

func ActivityA4(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA4...")
	return "A4 done", nil
}

func ActivityA5(ctx context.Context) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ActivityA5...")
	return "A5 done", nil
}
