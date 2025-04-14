package util

import "go.temporal.io/sdk/workflow"

type WorkflowVersionOverrideFlagValue struct {
	OverriddenVersion int `json:"overriddenVersion"` // Version of the workflow logic
	RolloutWeight     int `json:"rolloutWeight"`     // Relative weight (e.g., 10, 20, etc.)
}

type WorkflowVersionOverrideFlags map[string][]WorkflowVersionOverrideFlagValue

type Config struct {
	WorkflowVersionOverrideFlags WorkflowVersionOverrideFlags `json:"WorkflowVersionOverrideFlags"`
}

func GetVersion(ctx workflow.Context, versionOverride map[string]int, changeID string, minSupported, maxSupported workflow.Version) workflow.Version {
	if overriddenVersion, ok := versionOverride[changeID]; ok {
		workflow.GetVersion(ctx, changeID, workflow.Version(overriddenVersion), workflow.Version(maxSupported))
		return workflow.Version(overriddenVersion)
	}
	return workflow.GetVersion(ctx, changeID, minSupported, maxSupported)
}
