package main

import (
	"context"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer c.Close()

	workflowID := "newworkflow4"
	namespace := "default"

	// Step 1: Upsert search attribute to tag as version 34
	handle, err := c.UpdateWorkflow(context.Background(), client.UpdateWorkflowOptions{
		WorkflowID: workflowID,
		UpdateName: "upsert-search-attributes",
		Args: []interface{}{
			map[string]interface{}{
				"TemporalChangeVersion": "version-34",
			},
		},
	})
	if err != nil {
		log.Fatalln("Failed to initiate workflow update:", err)
	}

	// Wait for update completion
	err = handle.Get(context.Background(), nil)
	if err != nil {
		log.Fatalln("Failed to complete workflow update:", err)
	}

	// Step 2: Reset workflow to start from the beginning
	_, err = c.ResetWorkflowExecution(context.Background(), &workflowservice.ResetWorkflowExecutionRequest{
		Namespace: namespace,
		// WorkflowId:                workflowID,
		Reason:                    "Rollback to version 34",
		WorkflowTaskFinishEventId: 0, // Reset to start
	})
	if err != nil {
		log.Fatalln("Failed to reset workflow:", err)
	}

	log.Println("Successfully reset newworkflow4 to follow version 34")
}
