package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func OrderFulfillmentWorkflow(ctx workflow.Context, orderID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting OrderFulfillmentWorkflow", "orderID", orderID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// ========== Step 1: Check Order Status ==========
	var stockResult bool
	err := workflow.ExecuteActivity(ctx, "CheckStock", orderID).Get(ctx, &stockResult)
	if err != nil {
		logger.Error("Failed to check stock", "Error", err)
		return err
	}
	if !stockResult {
		logger.Info("Stock unavailable, cancelling order", "OrderID", orderID)
		return nil
	}

	// ========== Step 2: Reserve the product in inventory ==========
	err = workflow.ExecuteActivity(ctx, "ReserveProduct", orderID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to reserve product", "Error", err)
		return err
	}

	/*
		// ========== Setp 3: Introduce Non-DeterministicLogic ==========

		// Randomly decide whether to process payment or skip to shipping
		rand.Seed(time.Now().UnixNano()) // Seed random generator (BAD IDEA in workflows!)
		if rand.Intn(2) == 0 {           // 50% chance to process payment
			logger.Info("Randomly chose to process payment", "OrderID", orderID)
			var paymentResult bool
			err = workflow.ExecuteActivity(ctx, "ProcessPayment", orderID).Get(ctx, &paymentResult)
			if err != nil {
				logger.Error("Failed to process payment", "Error", err)
				return err
			}
			if !paymentResult {
				logger.Info("Payment failed, cancelling order", "OrderID", orderID)
				return nil
			}
		} else { // 50% chance to skip payment and ship directly
			logger.Info("Randomly chose to skip payment", "OrderID", orderID)
		}
	*/

	// ========== Step 3: Process payment ==========
	var paymentResult bool
	err = workflow.ExecuteActivity(ctx, "ProcessPayment", orderID).Get(ctx, &paymentResult)
	if err != nil {
		logger.Error("Failed to process payment", "Error", err)
		return err
	}
	if !paymentResult {
		logger.Info("Payment failed, cancelling order", "OrderID", orderID)
		return nil
	}

	// ========== Step 4: Ship the product ==========
	err = workflow.ExecuteActivity(ctx, "ShipProduct", orderID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to ship product", "Error", err)
		return err
	}

	logger.Info("Order fulfillment successfully", "OrderID", orderID)

	return nil
}

func CheckStock(ctx context.Context, orderID string) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Checking stock for order", "OrderID", orderID)
	time.Sleep(1 * time.Second)
	return true, nil
}

func ReserveProduct(ctx context.Context, orderID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Reserving product", "OrderID", orderID)
	time.Sleep(1 * time.Second)
	return nil
}

func ProcessPayment(ctx context.Context, orderID string) (bool, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Processing payment", "OrderID", orderID)
	time.Sleep(1 * time.Second)
	return true, nil
}

func ShipProduct(ctx context.Context, orderID string) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Shipping product", "OrderID", orderID)
	time.Sleep(1 * time.Second)
	return nil
}

func main() {
	// Create Temporal Client
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	w := worker.New(c, "order-fulfillment-task-queue", worker.Options{})

	w.RegisterWorkflow(OrderFulfillmentWorkflow) // tells the worker when you get the "OrderFulfillmentWorkflow" task then run this
	w.RegisterActivity(CheckStock)
	w.RegisterActivity(ReserveProduct)
	w.RegisterActivity(ProcessPayment)
	w.RegisterActivity(ShipProduct)

	go func() {
		err := w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "order-fulfillment-workflow-123",
		TaskQueue: "order-fulfillment-task-queue",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, OrderFulfillmentWorkflow, "order-123")
	if err != nil {
		log.Fatalf("Failed to start workflow: %v", err)
	}

	log.Printf("Started Workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())

	var result interface{}

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}

	log.Println("Workflow completed successfully", "Result", result)
}
