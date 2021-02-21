package namespace

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// SendEvent takes a namespace and a cloud event file to send the event to the namespace
func SendEvent(conn *grpc.ClientConn, namespace string, filepath string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	event, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	request := ingress.BroadcastEventRequest{
		Namespace:  &namespace,
		Cloudevent: event,
	}

	_, err = client.BroadcastEvent(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully sent event to '%s'", namespace), nil
}

// List returns a list of namespaces on direktiv
func List(conn *grpc.ClientConn) ([]*ingress.GetNamespacesResponse_Namespace, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.GetNamespacesRequest{}

	resp, err := client.GetNamespaces(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Namespaces, nil
}

// Delete deletes a namespace on direktiv
func Delete(name string, conn *grpc.ClientConn) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.DeleteNamespaceRequest{
		Name: &name,
	}

	resp, err := client.DeleteNamespace(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Deleted namespace: %s", resp.GetName()), nil
}

// Create creates a new namespace on direktiv
func Create(name string, conn *grpc.ClientConn) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)
	defer conn.Close()

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.AddNamespaceRequest{
		Name: &name,
	}

	resp, err := client.AddNamespace(ctx, &request)
	if err != nil {
		// convert the error
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Created namespace: %s", resp.GetName()), nil
}
