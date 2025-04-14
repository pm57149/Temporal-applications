package workflows

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

func CounterUpdateWorkflow(ctx workflow.Context, resourceID string, requesterID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Counter update requested by: ", requesterID)

	err := workflow.SignalExternalWorkflow(ctx, "Mutex-"+resourceID, "", "lock-request", requesterID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to signal external workflow", "error", err)
		return err
	}

	lockGrantedChan := workflow.GetSignalChannel(ctx, "lock-granted-"+resourceID)
	var granted bool
	lockGrantedChan.Receive(ctx, &granted)
	logger.Info("Lock granted to: ", requesterID)

	logger.Info("Updating counter for resource: ", resourceID)
	workflow.Sleep(ctx, 2*time.Second)

	err = workflow.SignalExternalWorkflow(ctx, "Mutex-"+resourceID, "", "lock-released-"+resourceID, true).Get(ctx, nil)

	if err != nil {
		return err
	}

	logger.Info("Counter update complete, lock released by: ", requesterID)
	return nil
}
