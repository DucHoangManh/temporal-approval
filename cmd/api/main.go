package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"approval-demo/cmd/internal/activity"
	"approval-demo/cmd/internal/workflow"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

func main() {
	temporalClient, err := client.Dial(
		client.Options{
			HostPort: client.DefaultHostPort,
		},
	)
	if err != nil {
		log.Fatalln("Unable to create Temporal Client.", err)
	}
	router := gin.Default()
	router.POST(
		"/approvals", func(context *gin.Context) {
			workflowOptions := client.StartWorkflowOptions{
				ID:                                       fmt.Sprintf("approval-workflow-%v", rand.Int()),
				TaskQueue:                                workflow.TaskQueueName,
				WorkflowExecutionErrorWhenAlreadyStarted: true,
			}
			_, err := temporalClient.ExecuteWorkflow(
				context.Request.Context(),
				workflowOptions,
				workflow.ApprovalRequiredWorkflow,
				workflow.DefaultApprovalDefinition,
				activity.PostApproveActionPayload{Id: rand.Int()},
			)
			if err != nil {
				context.JSON(500, gin.H{"message": err.Error()})
				return
			}
		},
	)
	router.GET(
		"/approvals/:workflowID", func(context *gin.Context) {
			workflowID := context.Param("workflowID")
			queryResult, err := temporalClient.QueryWorkflow(
				context.Request.Context(),
				workflowID,
				"",
				"getApprovalDefinition",
			)
			if err != nil {
				context.JSON(500, gin.H{"message": err.Error()})
				return
			}
			var queryResponse workflow.ApprovalDefinition
			err = queryResult.Get(&queryResponse)
			if err != nil {
				context.JSON(500, gin.H{"message": err.Error()})
				return
			}
			context.JSON(200, queryResponse)
		},
	)
	router.POST(
		"/approvals/:workflowID/approve", func(context *gin.Context) {
			workflowID := context.Param("workflowID")
			request := ApproveRequest{}
			if err := context.ShouldBindJSON(&request); err != nil {
				context.JSON(400, gin.H{"message": err.Error()})
				return
			}
			if err := temporalClient.SignalWorkflow(
				context.Request.Context(),
				workflowID,
				"",
				workflow.ApprovalSignal,
				workflow.ApprovalUser{
					Email:    request.Email,
					Approved: request.Approved,
				},
			); err != nil {
				context.JSON(500, gin.H{"message": err.Error()})
				return
			}
			context.JSON(200, gin.H{"message": "ok"})
		},
	)
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalln("Unable to start HTTP server.", err)
	}
}
