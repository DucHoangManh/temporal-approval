package main

import (
	"log"
	"log/slog"
	"os"

	"approval-demo/internal/activity"
	workflow2 "approval-demo/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	c, err := client.Dial(
		client.Options{
			HostPort:  client.DefaultHostPort,
			Namespace: client.DefaultNamespace,
			Logger:    logger,
		},
	)
	if err != nil {
		log.Fatalln("Unable to create Temporal Client.", err)
	}
	defer c.Close()
	// create Worker
	w := worker.New(c, workflow2.TaskQueueName, worker.Options{})
	// register Activity and Workflow
	w.RegisterWorkflow(workflow2.ApprovalRequiredWorkflow)
	w.RegisterActivity(activity.PostApproveActivity)

	log.Println("Worker is starting.")
	// Listen to Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Worker.", err)
	}
}
