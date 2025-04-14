package main

import (
	"Mutex_workflow/activities"
	"Mutex_workflow/workflows"
	"context"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"time"
)

func main() {
	//c, err := client.Dial(client.Options{
	//	HostPort: client.DefaultHostPort,
	//})
	//
	//if err != nil {
	//	log.Fatalf("Failed to create Temporal client: %v", err)
	//}
	//defer c.Close()
	//
	//resourceID := "resource-1"
	//
	//_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
	//	ID:        "MutexWorkflow",
	//	TaskQueue: "resource-task-queue",
	//}, workflows.MutexWorkflow, resourceID)
	//if err != nil {
	//	log.Fatalf("Failed to start MutexWorkflow: %v", err)
	//}
	//
	//time.Sleep(time.Minute * 1)
	//
	//users := []struct {
	//	userID  string
	//	message string
	//}{
	//	{"user1", "Request from user1"},
	//	{"user2", "Request from user2"},
	//}
	//
	//for _, user := range users {
	//	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
	//		ID:        "ClientWorkflow-" + user.userID,
	//		TaskQueue: "resource-task-queue",
	//	}, workflows.ClientWorkflow, user.userID, resourceID)
	//	if err != nil {
	//		log.Fatalf("Failed to start ClientWorkflow for %s: %v", user.userID, err)
	//	}
	//}
	//
	//time.Sleep(time.Second * 15)

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "resource-task-queue", worker.Options{})

	w.RegisterWorkflow(workflows.MutexWorkflow)
	w.RegisterWorkflow(workflows.ClientWorkflow)
	w.RegisterActivity(activities.ActivityA1)
	w.RegisterActivity(activities.ActivityA2)
	w.RegisterActivity(activities.ActivityA3)

	go func() {
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()

	time.Sleep(time.Second * 2)

	resourceID := "resource-1"
	_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        "MutexWorkflow-" + resourceID,
		TaskQueue: "resource-task-queue",
	}, workflows.MutexWorkflow, resourceID)
	if err != nil {
		log.Fatalf("Failed to start MutexWorkflow: %v", err)
	}

	time.Sleep(time.Second * 1)

	users := []struct {
		userID string
	}{
		{"user1"},
		{"user2"},
	}

	for _, user := range users {
		_, err = c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
			ID:        "ClientWorkflow-" + user.userID,
			TaskQueue: "resource-task-queue",
		}, workflows.ClientWorkflow, user.userID, resourceID)
		if err != nil {
			log.Fatalf("Failed to start ClientWorkflow for %s: %v", user.userID, err)
		}
	}

	time.Sleep(time.Second * 15)
}
