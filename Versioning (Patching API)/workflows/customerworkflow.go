package workflows

import (
	"Versioning_Patching_API/activities"
	"Versioning_Patching_API/model"
	"go.temporal.io/sdk/workflow"
	"time"
)

func CustomerWorkflow(ctx workflow.Context, customer model.Customer) (model.Account, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("CustomerWorkflow Info workflow started.", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 6,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	bonus := 100

	var account model.Account
	err := workflow.ExecuteActivity(ctx, activities.GetCustomerAccount, customer).Get(ctx, &account)
	if err != nil {
		logger.Error("GetCustomerAccount failed.", "Error", err)
		return model.Account{}, err
	}

	_ = workflow.Sleep(ctx, customer.DemoWaitDuration)

	//// CHANGE 1
	//v1 := workflow.GetVersion(ctx, "addedCheck", workflow.DefaultVersion, 1)
	//if v1 == 1 {
	//	var checkPassed bool
	//	err = workflow.ExecuteActivity(ctx, activities.CheckCustomerAccount, customer).Get(ctx, &checkPassed)
	//	if err != nil {
	//		logger.Error("GetCustomerAccount failed.", "Error", err)
	//		return model.Account{}, err
	//	}
	//	if !checkPassed {
	//		return model.Account{}, err
	//	}
	//	_ = workflow.Sleep(ctx, customer.DemoWaitDuration)
	//}
	//// END CHANGE 1

	// CHANGE 2
	//v21 := workflow.GetVersion(ctx, "addedBonus", workflow.DefaultVersion, 1)
	//if v21 == 1 {
	//	bonus = 200
	//}
	// END CHANGE 2

	err = workflow.ExecuteActivity(ctx, activities.UpdateCustomerAccount, customer, bonus).Get(ctx, &account)
	if err != nil {
		logger.Error("UpdateCustomerAccount failed.", "Error", err)
		return model.Account{}, err
	}

	return account, err
}
