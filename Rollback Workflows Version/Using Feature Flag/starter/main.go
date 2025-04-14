package main

import (
	"Using_Feature_Flag/workflows"
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"log"
)

type VersionSignal struct {
	TargetVersion int `json:"targetVersion"`
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	for i := 1; i <= 5; i++ {
		workflowID := fmt.Sprintf("RollbackUsingSignalWorkflow-%d", i)
		workflowOptions := client.StartWorkflowOptions{
			ID:        workflowID,
			TaskQueue: "RollbackUsingSignal",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.OrderWorkflow)
		if err != nil {
			log.Fatalf("Failed to start Workflow %d: %v", i, err)
		}

		log.Printf("Started Workflow %d: %s", i, we.GetID())
	}

	// ======== Sending Signal to Roll back to Version 38 ========
	workflowOptions := client.StartWorkflowOptions{
		ID:        "new-workflow",
		TaskQueue: "RollbackUsingSignal",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.OrderWorkflow)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}
	log.Printf("Started workflow: %s", we.GetID())

	workflowID := "new-workflow"

	signal := VersionSignal{TargetVersion: 38}

	err = c.SignalWorkflow(context.Background(), workflowID, "", "VersionOverrideSignal", signal)
	if err != nil {
		log.Fatalf("Failed to send signal to workflow %s: %v", workflowID, err)
	}

	log.Printf("Successfully sent signal to workflow %s with target version %d", workflowID, signal.TargetVersion)
	// =====================================================================
}
