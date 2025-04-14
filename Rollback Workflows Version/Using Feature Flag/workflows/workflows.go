package workflows

import (
	"Using_Feature_Flag/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

type VersionSignal struct {
	TargetVersion int `json:"targetVersion"`
}

type ActivityConfig struct {
	ActivityName string
}

func OrderWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Order workflow started...")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	var targetVersion workflow.Version = -2
	const maxSupportedVersion = 58
	const minSupportedVersion = workflow.DefaultVersion

	versionSignalCh := workflow.GetSignalChannel(ctx, "VersionOverrideSignal")

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA1).Get(ctx, &result); err != nil {
		logger.Error("Activity A1 failed", "Error", err)
		return err
	}

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA2).Get(ctx, &result); err != nil {
		logger.Error("Activity A2 failed", "Error", err)
		return err
	}

	logger.Info("About to execute versioned activity...")
	workflow.Sleep(ctx, 1*time.Minute)

	// =================== Versioning Logic ===================
	var signal VersionSignal
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(versionSignalCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &signal)
		if signal.TargetVersion < int(minSupportedVersion) || signal.TargetVersion > maxSupportedVersion {
			logger.Error("Invalid target version", "TargetVersion", signal.TargetVersion)
		} else {
			logger.Info("Received version override signal", "TargetVersion", signal.TargetVersion)
			targetVersion = workflow.Version(signal.TargetVersion)
		}
	})
	selector.AddDefault(func() {
		logger.Info("No version override signal received, using default")
	})
	selector.Select(ctx)

	version := workflow.GetVersion(ctx, "OrderWorkflowVersion", minSupportedVersion, maxSupportedVersion)

	effectiveVersion := targetVersion
	if effectiveVersion == -2 {
		effectiveVersion = version
	}

	activityConfig := getActivityConfig(int(effectiveVersion))
	if activityConfig == nil {
		logger.Error("No activity configured for version", "Version", effectiveVersion)
		return nil
	}

	logger.Info("Executing versioned activity", "Activity", activityConfig.ActivityName, "Version", effectiveVersion)
	if err := workflow.ExecuteActivity(ctx, activityConfig.ActivityName).Get(ctx, &result); err != nil {
		logger.Error("Versioned activity failed", "Activity", activityConfig.ActivityName, "Error", err)
		return err
	}
	// ============================================================

	if err := workflow.ExecuteActivity(ctx, activities.ActivityA3).Get(ctx, &result); err != nil {
		logger.Error("Activity A3 failed", "Error", err)
		return err
	}

	return nil
}

func getActivityConfig(version int) *ActivityConfig {
	if version < 0 {
		return nil
	}
	
	if version == 58 {
		return &ActivityConfig{
			ActivityName: "New_Activity58",
		}
	} else if version == 38 {
		return &ActivityConfig{
			ActivityName: "New_Activity38",
		}
	} else {
		return &ActivityConfig{
			ActivityName: "Default_Activity",
		}
	}
}
