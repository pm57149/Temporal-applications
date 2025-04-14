package main

import (
	"Signals_in_workflow/workflows"
	"context"
	"go.temporal.io/sdk/client"
	"log"
	"time"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "movie-booking-123",
		TaskQueue: "SignalsGoDemoTaskQueue",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflows.MovieBookingWorkflow, "A12", 12.50)
	if err != nil {
		log.Fatalln("Unable to start workflow", err)
	}
	log.Printf("Started workflow: %s", we.GetID())

	// =============== Simulate user cancelling after a few seconds ================
	time.Sleep(3 * time.Second)
	err = c.SignalWorkflow(context.Background(), "movie-booking-123", "", "cancel-booking", "User cancelled")
	if err != nil {
		log.Fatalln("Unable to send signal", err)
	}
	log.Println("Cancel signal sent successfully")
}
