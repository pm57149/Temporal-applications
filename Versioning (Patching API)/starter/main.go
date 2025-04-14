package main

import (
	"Versioning_Patching_API/model"
	"Versioning_Patching_API/workflows"
	"context"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func main() {
	//// Create Temporal Client
	//c, err := client.Dial(client.Options{
	//	HostPort: "localhost:7233",
	//})
	//if err != nil {
	//	log.Fatalf("Failed to create Temporal client: %v", err)
	//}
	//defer c.Close()
	//
	//w := worker.New(c, "task-queue-name", worker.Options{})
	//
	//w.RegisterWorkflow(workflows.StorageManagementWorkflow)
	//w.RegisterActivity(activities.CheckAvailableSpace)
	//w.RegisterActivity(activities.VerifyClusterAvailability)
	//w.RegisterActivity(activities.CreateStorageAggregate)
	//w.RegisterActivity(activities.CreateVolume)
	//w.RegisterActivity(activities.ConfigureAccessAndSecurity)
	//w.RegisterActivity(activities.EnableDeduplication)
	//w.RegisterActivity(activities.EnableCompression)
	//
	//go func() {
	//	err := w.Run(worker.InterruptCh())
	//	if err != nil {
	//		log.Fatalf("Failed to start worker: %v", err)
	//	}
	//}()
	//
	//workflowOptions := client.StartWorkflowOptions{
	//	ID:        "storage-management-workflow",
	//	TaskQueue: "task-queue-name",
	//}
	//
	//we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.StorageManagementWorkflow)
	//if err != nil {
	//	log.Fatalf("Failed to execute workflow: %v", err)
	//}
	//
	//log.Printf("Started workflow %q.\n", we.GetID())
	//
	//var result interface{}
	//
	//err = we.Get(context.Background(), &result)
	//if err != nil {
	//	log.Fatalf("Workflow failed: %v", err)
	//}
	//
	//log.Println("Workflow completed successfully", "Result", result)

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	//customer := &model.Customer{
	//	AccountNum:       "c1",
	//	Name:             "Pawam",
	//	Email:            "pawan@buildcode.com",
	//	CustomerType:     "new",
	//	DemoWaitDuration: 2 * time.Minute,
	//}

	customers := []model.Customer{
		{
			AccountNum:       "c1",
			Name:             "Pawan",
			Email:            "pawan@buildcode.com",
			CustomerType:     "new",
			DemoWaitDuration: 10 * time.Second, // Short wait to finish quickly
		},
		{
			AccountNum:       "c2",
			Name:             "Alice",
			Email:            "alice@buildcode.com",
			CustomerType:     "existing",
			DemoWaitDuration: 30 * time.Second, // Medium wait
		},
		{
			AccountNum:       "c3",
			Name:             "Bob",
			Email:            "bob@buildcode.com",
			CustomerType:     "new",
			DemoWaitDuration: 60 * time.Second, // Long wait to stay running
		},
	}

	for _, customer := range customers {
		workflowOptions := client.StartWorkflowOptions{
			ID:        customer.AccountNum,
			TaskQueue: "versioningGoDemoTaskQueue",
		}

		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CustomerWorkflow, customer)
		if err != nil {
			log.Fatalln("Unable to execute workflow for", customer.AccountNum, err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		// Synchronously wait for the workflow completion.
		var result model.Account
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalln("Unable get workflow result", err)
		}
		log.Println("Workflow result:", result.Amount)
	}

	//workflowOptions := client.StartWorkflowOptions{
	//	ID:        customer.AccountNum,
	//	TaskQueue: "versioningGoDemoTaskQueue",
	//}
	//
	//we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.CustomerWorkflow, *customer)
	//if err != nil {
	//	log.Fatalln("Unable to execute workflow", err)
	//}
	// log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	//// Synchronously wait for the workflow completion.
	//var result model.Account
	//err = we.Get(context.Background(), &result)
	//if err != nil {
	//	log.Fatalln("Unable get workflow result", err)
	//}
	//log.Println("Workflow result:", result.Amount)
}
