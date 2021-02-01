package namespace

import (
	"errors"
	"fmt"

	"github.com/segmentio/ksuid"
	log "github.com/vorteil/direkcli/pkg/log"
	"github.com/vorteil/direkcli/pkg/util"
	"github.com/vorteil/direktiv/pkg/direktiv"
	"github.com/vorteil/vorteil/pkg/elog"
)

var logger elog.View

func init() {
	log := log.GetLogger()
	logger = log
}

// List returns a list of namespaces
func List() ([]direktiv.CmdGetNamespacesResponse, error) {
	logger.Printf("Fetching namespaces...")

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return nil, err
	}
	defer n.Conn.Close()

	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.GetNamespaces

	var da direktiv.CmdGetNamespaces
	cmd.Cmd = da

	// send request to direktiv
	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return nil, err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return nil, errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		// Response was successful :)
		namespaces := make([]direktiv.CmdGetNamespacesResponse, 0)
		if err := n.DirektivUnmarshal(resp, &namespaces); err != nil {
			return nil, err
		}
		return namespaces, nil
	}

	return nil, errors.New("An unexpected error occurred")
}

// Delete removes a namespace previously added
func Delete(name string) (string, error) {
	logger.Printf("Deleting namespace '%s'...", name)

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	defer n.Conn.Close()

	// Delete namespace
	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.DeleteNamespace

	var da direktiv.CmdAddDeleteNamespace
	da.Name = name
	cmd.Cmd = da

	// Send Request to direktiv
	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		// Response was successful :)
		return fmt.Sprintf("Namespace '%s' was successfully deleted.", name), nil
	}

	return "", errors.New("An unexpected error occurred")
}

// Create adds a new namespace and returns with a successful message or an error
func Create(name string) (string, error) {
	logger.Printf("Creating namespace '%s'...", name)

	// Open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	// Close nats at end of function
	defer n.Conn.Close()

	// Create namespace
	// Construct direktiv request
	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.AddNamespace

	var da direktiv.CmdAddDeleteNamespace
	da.Name = name
	cmd.Cmd = da

	// Send request
	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		// Response was successful :)
		return fmt.Sprintf("Namespace '%s' was successfully created.", name), nil
	}

	return "", errors.New("An unexpected error occurred")
}
