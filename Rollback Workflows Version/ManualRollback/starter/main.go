package main

import (
	"UsingUpdateWorkflow/workflows"
	"context"
	"go.temporal.io/sdk/client"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	//for i := 1; i <= 5; i++ {
	//	workflowID := fmt.Sprintf("UsingUpdateWorkflow-%d", i)
	//	workflowOptions := client.StartWorkflowOptions{
	//		ID:        workflowID,
	//		TaskQueue: "RollbackWithUpdateWorkflow",
	//	}
	//
	//	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SampleWorkflow)
	//	if err != nil {
	//		log.Fatalf("Failed to start Workflow %d: %v", i, err)
	//	}
	//
	//	log.Printf("Started Workflow %d: %s", i, we.GetID())
	//}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "new-workflow",
		TaskQueue: "RollbackWithUpdateWorkflow",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.SampleWorkflow)
	log.Printf("Started workflow: %s", we.GetID())

	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}
}
