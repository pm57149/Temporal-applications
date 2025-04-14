package workflows

import (
	"Versioning_Patching_API/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

const (
	Version1 workflow.Version = 0
	Version2 workflow.Version = 1
)

func StorageManagementWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Check available space
	err := workflow.ExecuteActivity(ctx, activities.CheckAvailableSpace).Get(ctx, nil)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Minute)

	// Verify cluster availability
	err = workflow.ExecuteActivity(ctx, activities.VerifyClusterAvailability).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Create storage aggregate
	err = workflow.ExecuteActivity(ctx, activities.CreateStorageAggregate).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Create volume
	err = workflow.ExecuteActivity(ctx, activities.CreateVolume).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Configure access and security
	err = workflow.ExecuteActivity(ctx, activities.ConfigureAccessAndSecurity).Get(ctx, nil)
	if err != nil {
		return err
	}

	// Enable deduplication
	err = workflow.ExecuteActivity(ctx, activities.EnableDeduplication).Get(ctx, nil)
	if err != nil {
		return err
	}

	/*
		// Use workflow versioning to conditionally execute new activities
		if workflow.GetVersion(ctx, "enableCompression", workflow.DefaultVersion, Version2) == Version2 {
			// Enable compression (new activity in version 2.0)
			err = workflow.ExecuteActivity(ctx, activities.EnableCompression).Get(ctx, nil)
			if err != nil {
				return err
			}
		}
	*/

	// Enable compression (new activity in version 2.0)
	err = workflow.ExecuteActivity(ctx, activities.EnableCompression).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
