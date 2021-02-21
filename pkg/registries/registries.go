package registries

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Create creates a new registry on a namespace
func Create(conn *grpc.ClientConn, namespace string, key string, value string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.StoreRegistryRequest{
		Namespace: &namespace,
		Name:      &key,
		Data:      []byte(value),
	}

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

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.GetRegistriesRequest{
		Namespace: &namespace,
	}

	resp, err := client.GetRegistries(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Registries, nil
}

// Delete remvoes a registry from a namespace
func Delete(conn *grpc.ClientConn, namespace string, key string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	request := ingress.DeleteRegistryRequest{
		Namespace: &namespace,
		Name:      &key,
	}

	_, err := client.DeleteRegistry(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully removed registry '%s'.", key), nil
}
