package workflows

import (
	"Signals_in_workflow/activities"
	"go.temporal.io/sdk/workflow"
	"time"
)

func MovieBookingWorkflow(ctx workflow.Context, seatNumber string, amount float64) error {
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var reserveResult string
	err := workflow.ExecuteActivity(ctx, activities.ReserveSeat, seatNumber).Get(ctx, &reserveResult)
	if err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Activity result: %s", reserveResult)

	cancelCh := workflow.GetSignalChannel(ctx, "cancel-booking")
	isCancelled := false

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(cancelCh, func(c workflow.ReceiveChannel, more bool) {
		var signalData string
		c.Receive(ctx, &signalData)
		isCancelled = true
		workflow.GetLogger(ctx).Info("Received cancel signal: %s", signalData)
	})

	workflow.Sleep(ctx, 5*time.Second)
	selector.Select(ctx)

	if isCancelled {
		workflow.GetLogger(ctx).Info("Booking cancelled, skipping payment")
		return nil
	}

	var paymentResult string
	err = workflow.ExecuteActivity(ctx, activities.ProcessPayment, amount).Get(ctx, &paymentResult)
	if err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Activity result: %s", paymentResult)

	workflow.GetLogger(ctx).Info("Booking completed successfully!")
	return nil
}
