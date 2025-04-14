package workflows

import (
	"Using_Signals/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

func SampleWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	if err := workflow.ExecuteActivity(ctx, activities.ActivityA1).Get(ctx, &result); err != nil {
		logger.Error("Activity A1 failed", "Error", err)
		return err
	}

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA2).Get(ctx, &result); err != nil {
		logger.Error("Activity A2 failed", "Error", err)
		return err
	}

	logger.Error("About to execute next activity...")
	workflow.Sleep(ctx, 1*time.Minute)

	//version := workflow.GetVersion(ctx, "add-a2.5", workflow.DefaultVersion, 1)
	//if version == 1 {
	//	if err := workflow.ExecuteActivity(ctx, activities.ActivityA2_5).Get(ctx, &result); err != nil {
	//		logger.Error("Activity A2.5 failed", "Error", err)
	//		return err
	//	}
	//} else {
	//	logger.Info("Skipping A2.5 for default version")
	//}

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA3).Get(ctx, &result); err != nil {
		logger.Error("Activity A3 failed", "Error", err)
		return err
	}

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA4).Get(ctx, &result); err != nil {
		logger.Error("Activity A4 failed", "Error", err)
		return err
	}

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA5).Get(ctx, &result); err != nil {
		logger.Error("Activity A5 failed", "Error", err)
		return err
	}

	logger.Info("Workflow completed successfully")
	return nil
}
