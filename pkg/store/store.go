package registries

import (
	"fmt"

	"github.com/vorteil/direkcli/pkg/util"
	"github.com/vorteil/direktiv/pkg/ingress"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// CreateRegistry a new registry
func CreateRegistry(conn *grpc.ClientConn, namespace string, key string, value string) (string, error) {
	client, ctx, cancel := util.CreateClient(conn)
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

// ListRegistries returns a list of registries
func ListRegistries(conn *grpc.ClientConn, namespace string) ([]*ingress.GetRegistriesResponse_Registry, error) {
	client, ctx, cancel := util.CreateClient(conn)
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

// DeleteRegistry removes a registry from a namespace
func DeleteRegistry(conn *grpc.ClientConn, namespace string, key string) (string, error) {
	client, ctx, cancel := util.CreateClient(conn)
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

// CreateSecret creates a new secret within a namespace
func CreateSecret(conn *grpc.ClientConn, namespace string, secret string, value string) (string, error) {
	client, ctx, cancel := util.CreateClient(conn)
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

// ListSecrets returns a list of secrets
func ListSecrets(conn *grpc.ClientConn, namespace string) ([]*ingress.GetSecretsResponse_Secret, error) {
	client, ctx, cancel := util.CreateClient(conn)
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

// DeleteSecret removes a secret from a namespace
func DeleteSecret(conn *grpc.ClientConn, namespace string, secret string) (string, error) {
	client, ctx, cancel := util.CreateClient(conn)
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
