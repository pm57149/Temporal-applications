package activities

import (
	"context"
	"go.temporal.io/sdk/activity"
	"time"
)

func ReserveSeat(ctx context.Context, seatNumber string) (string, error) {
	activity.GetLogger(ctx).Info("Reserving seat: %s", seatNumber)
	time.Sleep(2 * time.Second)
	return "Seat " + seatNumber + " reserved", nil
}

func ProcessPayment(ctx context.Context, amount float64) (string, error) {
	activity.GetLogger(ctx).Info("Processing payment of $%.2f", amount)
	time.Sleep(3 * time.Second)
	return "Payment successful", nil
}
