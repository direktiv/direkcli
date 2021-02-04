package instance

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Logs returns all logs associated with the workflow instance ID
func Logs(conn *grpc.ClientConn, id string) ([]*protocol.GetWorkflowInstanceLogsResponse_WorkflowInstanceLog, error) {
	client := protocol.NewDirektivClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	offset := int32(0)
	limit := int32(10000)

	request := protocol.GetWorkflowInstanceLogsRequest{
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
func List(conn *grpc.ClientConn, namespace string) ([]*protocol.GetWorkflowInstancesResponse_WorkflowInstance, error) {
	client := protocol.NewDirektivClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := protocol.GetWorkflowInstancesRequest{
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
func Get(conn *grpc.ClientConn, id string) (*protocol.GetWorkflowInstanceResponse, error) {
	client := protocol.NewDirektivClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := protocol.GetWorkflowInstanceRequest{
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
