package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	"temporalBasic1/app/iplocate"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Must specify a name as the command-line argument")
	}
	name := os.Args[1]

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowID := "GetAddressFromIP-" + uuid.New().String()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: iplocate.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, iplocate.GetAddressFromIP, name)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get result", err)
	}

	fmt.Println(result)
}
