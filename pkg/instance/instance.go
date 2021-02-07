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

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	offset := int32(0)
	limit := int32(10000)

	request := ingress.GetWorkflowInstanceLogsRequest{
		InstanceId: &id,
		Offset:     &offset,
		Limit:      &limit,
	}

	resp, err := client.GetWorkflowInstanceLogs(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.GetWorkflowInstanceLogs(), nil
}

// List returns a list of workflow instances
func List(conn *grpc.ClientConn, namespace string) ([]*ingress.GetWorkflowInstancesResponse_WorkflowInstance, error) {
	client := ingress.NewDirektivIngressClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.GetWorkflowInstancesRequest{
		Namespace: &namespace,
	}

	resp, err := client.GetWorkflowInstances(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.WorkflowInstances, nil
}

// Get returns a pretty printed json of the workflow instance id
func Get(conn *grpc.ClientConn, id string) (*ingress.GetWorkflowInstanceResponse, error) {
	client := ingress.NewDirektivIngressClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.GetWorkflowInstanceRequest{
		Id: &id,
	}

	resp, err := client.GetWorkflowInstance(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp, nil
}
