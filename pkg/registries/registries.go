package registries

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Create a new registry
func Create(conn *grpc.ClientConn, namespace string, key string, value string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.StoreRegistryRequest{
		Namespace: &namespace,
		Name:      &key,
		Data:      []byte(value),
	}

	// send grpc request
	_, err := client.StoreRegistry(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully created registry '%s'.", key), nil
}

// List returns a list of registries
func List(conn *grpc.ClientConn, namespace string) ([]*ingress.GetRegistriesResponse_Registry, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetRegistriesRequest{
		Namespace: &namespace,
	}

	// send grpc request
	resp, err := client.GetRegistries(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Registries, nil
}

// Delete removes a registry from a namespace
func Delete(conn *grpc.ClientConn, namespace string, key string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.DeleteRegistryRequest{
		Namespace: &namespace,
		Name:      &key,
	}

	// send grpc request
	_, err := client.DeleteRegistry(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully removed registry '%s'.", key), nil
}
