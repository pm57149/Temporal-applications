package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"net/http"
	"temporalBasic1/app/iplocate"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}

	defer c.Close()

	w := worker.New(c, iplocate.TaskQueueName, worker.Options{})

	activities := &iplocate.IPActivities{
		HTTPClient: http.DefaultClient,
	}

	w.RegisterWorkflow(iplocate.GetAddressFromIP)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
