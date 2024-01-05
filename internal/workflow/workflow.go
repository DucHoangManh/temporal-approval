package workflow

import (
	"time"

	activity2 "approval-demo/internal/activity"
	"go.temporal.io/sdk/workflow"
)

func ApprovalRequiredWorkflow(ctx workflow.Context, definition ApprovalDefinition, payload activity2.PostApproveActionPayload) error {
	logger := workflow.GetLogger(ctx)
	if err := approvalWorkflow(ctx, &definition); err != nil {
		logger.Error("approvalWorkflow failed.", "Error", err)
		return err
	}
	if definition.Status != ApprovalStatusApproved {
		logger.Info("approvalWorkflow rejected", "payload", payload)
		return nil
	}
	ctx = workflow.WithActivityOptions(
		ctx, workflow.ActivityOptions{
			StartToCloseTimeout: 2 * time.Minute,
			WaitForCancellation: true,
		},
	)
	// should handle workflow cancellation
	if err := workflow.ExecuteActivity(ctx, activity2.PostApproveActivity, payload).Get(ctx, nil); err != nil {
		logger.Error("ExecuteActivity failed.", "Error", err)
		return err
	}
	return nil
}

func approvalWorkflow(ctx workflow.Context, definition *ApprovalDefinition) error {
	logger := workflow.GetLogger(ctx)
	err := workflow.SetQueryHandler(
		ctx, "getApprovalDefinition", func() (ApprovalDefinition, error) {
			return *definition, nil
		},
	)
	if err != nil {
		logger.Error("SetQueryHandler failed.", "Error", err)
		return err
	}
	loop := true
	for loop {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(
			workflow.GetSignalChannel(ctx, ApprovalSignal), func(c workflow.ReceiveChannel, _ bool) {
				actionSignal := ApprovalUser{}
				c.Receive(ctx, &actionSignal)
				if actionSignal.Email == "" {
					logger.Info("empty email")
					return
				}
				logger.Info("received signal", "Email", actionSignal.Email, "Approved", actionSignal.Approved)
				if !definition.ContainsEmail(actionSignal.Email) {
					logger.Info("irrelevant email", "Email", actionSignal.Email)
					return
				}
				if actionSignal.Approved {
					definition.HandleApprove(actionSignal.Email)
					if definition.Status == ApprovalStatusApproved {
						logger.Info("approved")
						loop = false
					}
				} else {
					definition.HandleReject()
					logger.Info("rejected")
					loop = false
				}
			},
		)
		// waiting for cancel signal
		selector.AddReceive(
			ctx.Done(), func(_ workflow.ReceiveChannel, _ bool) {
				definition.Status = ApprovalStatusRejected
				logger.Info("canceled")
				loop = false
				return
			},
		)
		selector.Select(ctx)
	}
	return nil
}
