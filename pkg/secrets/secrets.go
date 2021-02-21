package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// Create creates a new secret within a namespace
func Create(conn *grpc.ClientConn, namespace string, secret string, value string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.StoreSecretRequest{
		Namespace: &namespace,
		Name:      &secret,
		Data:      []byte(value),
	}

	// send grpc request
	_, err := client.StoreSecret(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully created secret '%s'.", secret), nil
}

// List returns a list of secrets
func List(conn *grpc.ClientConn, namespace string) ([]*ingress.GetSecretsResponse_Secret, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.GetSecretsRequest{
		Namespace: &namespace,
	}

	// send grpc request
	resp, err := client.GetSecrets(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return nil, fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return resp.Secrets, nil
}

// Delete removes a secret from a namespace
func Delete(conn *grpc.ClientConn, namespace string, secret string) (string, error) {
	client := ingress.NewDirektivIngressClient(conn)

	// set context with 3 second timeout
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()

	// prepare request
	request := ingress.DeleteSecretRequest{
		Namespace: &namespace,
		Name:      &secret,
	}

	// send grpc request
	_, err := client.DeleteSecret(ctx, &request)
	if err != nil {
		s := status.Convert(err)
		return "", fmt.Errorf("[%v] %v", s.Code(), s.Message())
	}

	return fmt.Sprintf("Successfully removed secret '%s'.", secret), nil
}
