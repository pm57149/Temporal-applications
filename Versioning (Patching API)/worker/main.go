package main

import (
	"Versioning_Patching_API/activities"
	"Versioning_Patching_API/workflows"
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

	w := worker.New(c, "versioningGoDemoTaskQueue", worker.Options{})

	w.RegisterWorkflow(workflows.CustomerWorkflow)
	w.RegisterActivity(activities.CheckCustomerAccount)
	w.RegisterActivity(activities.GetCustomerAccount)
	w.RegisterActivity(activities.UpdateCustomerAccount)
	w.RegisterActivity(activities.SendBonusEmail)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
