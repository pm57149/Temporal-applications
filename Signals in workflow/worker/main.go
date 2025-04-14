package main

import (
	"Signals_in_workflow/activities"
	"Signals_in_workflow/workflows"
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

	w := worker.New(c, "SignalsGoDemoTaskQueue", worker.Options{})

	w.RegisterWorkflow(workflows.MovieBookingWorkflow)
	w.RegisterActivity(activities.ReserveSeat)
	w.RegisterActivity(activities.ProcessPayment)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}
}
