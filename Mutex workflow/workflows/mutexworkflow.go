package workflows

import (
	"go.temporal.io/sdk/workflow"
)

func MutexWorkflow(ctx workflow.Context, resourceID string) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Mutex workflow started for resource: ", resourceID)

	currentUser := ""
	hasResource := false

	requestResourceCh := workflow.GetSignalChannel(ctx, "request-resource")
	releaseResourceCh := workflow.GetSignalChannel(ctx, "release_resource")

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(requestResourceCh, func(c workflow.ReceiveChannel, more bool) {
			var userID string
			c.Receive(ctx, &userID)

			if hasResource {
				logger.Info("Resource already held by: ", "UserID", currentUser)
				return
			}

			hasResource = true
			currentUser = userID
			logger.Info("Resource held by: ", "UserID", currentUser)
		})

		selector.AddReceive(releaseResourceCh, func(c workflow.ReceiveChannel, more bool) {
			var userID string
			c.Receive(ctx, &userID)

			if currentUser != userID || !hasResource {
				logger.Info("Can not release resource ", "UserID", currentUser)
				return
			}

			hasResource = false
			currentUser = ""
			logger.Info("Resource released by: ", "UserID", userID)
		})

		selector.Select(ctx)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
