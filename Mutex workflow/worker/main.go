package main

import (
	"Mutex_workflow/activities"
	"Mutex_workflow/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "resource-task-queue", worker.Options{})

	w.RegisterWorkflow(workflows.ClientWorkflow)

	w.RegisterWorkflow(workflows.MutexWorkflow)

	w.RegisterActivity(activities.ActivityA1)
	w.RegisterActivity(activities.ActivityA2)
	w.RegisterActivity(activities.ActivityA2_5)
	w.RegisterActivity(activities.ActivityA3)
	w.RegisterActivity(activities.ActivityA4)
	w.RegisterActivity(activities.ActivityA5)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
