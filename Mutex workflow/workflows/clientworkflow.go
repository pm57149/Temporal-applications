package workflows

import (
	"Mutex_workflow/activities"
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func ClientWorkflow(ctx workflow.Context, userID string, resourceID string) error {
	logger := workflow.GetLogger(ctx)

	logger.Info("Client workflow started...")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA1).Get(ctx, &result); err != nil {
		logger.Error("Activity A1 failed", "Error", err)
		return err
	}

	// ========== Resource Locking ============================= //
	mutexWorkflowID := "MutexWorkflow-" + resourceID
	logger.Info("Requesting resource for user: ", userID)

	var err error
	for i := 0; i < 3; i++ {
		err = workflow.SignalExternalWorkflow(ctx, mutexWorkflowID, "", "request-resource", userID).Get(ctx, nil)
		if err == nil {
			break
		}
		logger.Info("Failed to signal MutexWorkflow, retrying...", "Attempt", i+1, "Error", err)
		workflow.Sleep(ctx, time.Second*1)
	}
	if err != nil {
		return fmt.Errorf("could not signal external workflow after retries: %w", err)
	}

	workflow.Sleep(ctx, time.Second*2)

	logger.Info("Processing request for user: ", userID)

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA2).Get(ctx, &result); err != nil {
		_ = workflow.SignalExternalWorkflow(
			ctx,
			mutexWorkflowID,
			"",
			"release_resource",
			userID).Get(ctx, nil)
		logger.Error("Activity A2 failed", "Error", err)
		return err
	}

	logger.Info("Releasing resource ", "UserID", userID)
	err = workflow.SignalExternalWorkflow(
		ctx,
		mutexWorkflowID,
		"",
		"release_resource",
		userID).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to release resource: %w", err)
	}
	// =================================================================//

	workflow.Sleep(ctx, 5*time.Second)

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA3).Get(ctx, &result); err != nil {
		logger.Error("Activity A3 failed", "Error", err)
		return err
	}

	return nil
}
