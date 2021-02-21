package instance

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logs returns all logs associated with the workflow instance ID
func Logs(conn *grpc.ClientConn, id string) ([]*ingress.GetWorkflowInstanceLogsResponse_WorkflowInstanceLog, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	offset := int32(0)
	limit := int32(10000)

	// prepare request
	request := ingress.GetWorkflowInstanceLogsRequest{
		InstanceId: &id,
		Offset:     &offset,
		Limit:      &limit,
	}

	// send grpc request
	resp, err := client.GetWorkflowInstanceLogs(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.GetWorkflowInstanceLogs(), nil
}

// List workflow instances
func List(conn *grpc.ClientConn, namespace string) ([]*ingress.GetWorkflowInstancesResponse_WorkflowInstance, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetWorkflowInstancesRequest{
		Namespace: &namespace,
	}

	// send grpc request
	resp, err := client.GetWorkflowInstances(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.WorkflowInstances, nil
}

// Get returns a workflow instance.
func Get(conn *grpc.ClientConn, id string) (*ingress.GetWorkflowInstanceResponse, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetWorkflowInstanceRequest{
		Id: &id,
	}

	// send grpc request
	resp, err := client.GetWorkflowInstance(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp, nil
}
