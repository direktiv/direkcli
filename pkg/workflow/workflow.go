package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/segmentio/ksuid"
	"github.com/vorteil/direkcli/pkg/log"
	"github.com/vorteil/direkcli/pkg/util"
	"github.com/vorteil/direktiv/pkg/direktiv"
	"github.com/vorteil/vorteil/pkg/elog"
)

var logger elog.View

func init() {
	log := log.GetLogger()
	logger = log
}

// List returns an array of workflows for a given namespace
func List(namespace string) ([]direktiv.CmdGetWorkflowsResponse, error) {
	logger.Printf("Fetching Workflow list...")

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return nil, err
	}

	defer n.Conn.Close()

	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.GetWorkflows

	var da direktiv.CmdGetWorkflows
	da.Namespace = namespace
	cmd.Cmd = da

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
		workflows := make([]direktiv.CmdGetWorkflowsResponse, 0)
		if err := n.DirektivUnmarshal(resp, &workflows); err != nil {
			return nil, err
		}
		return workflows, nil
	}

	return nil, errors.New("An unexpected error occurred")
}

// Execute runs the yaml provided from the workflow
func Execute(input string, id, namespace string) (string, error) {
	logger.Printf("Executing workflow '%s'...", id)
	var err error
	b := []byte{}

	// if provided input read file
	if input != "" {
		// read input for workflow
		b, err = ioutil.ReadFile(input)
		if err != nil {
			return "", err
		}
	}

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}

	defer n.Conn.Close()
	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.InvokeWorkflow

	var da direktiv.CmdInvokeWorkflow
	da.Workflow = id
	da.Namespace = namespace
	da.Data = b

	cmd.Cmd = da

	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		dWorkflow := new(direktiv.CmdInvokeWorkflowResponse)
		if err := n.DirektivUnmarshal(resp, &dWorkflow); err != nil {
			return "", err
		}
		return string(dWorkflow.InstanceID), nil
	}

	return "", errors.New("An unexpected error occurred")

}

// Get returns the YAML contents of the workflow
func Get(id string, namespace string) (string, error) {
	logger.Printf("Fetching YAML workflow '%s'...", id)

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}

	defer n.Conn.Close()
	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()

	var da direktiv.CmdGetWorkflow

	cmd.CmdType = direktiv.GetWorkflow
	da.ID = id
	cmd.Cmd = da

	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		dWorkflow := new(direktiv.CmdGetWorkflowResponse)
		if err := n.DirektivUnmarshal(resp, &dWorkflow); err != nil {
			return "", err
		}
		return string(dWorkflow.Workflow), nil
	}

	return "", errors.New("An unexpected error occurred")
}

// Update updates a workflow from the provided id
func Update(filepath string, id string, namespace string) (string, error) {
	logger.Printf("Updating Workflow '%s'...", id)

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}

	defer n.Conn.Close()

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.UpdateWorkflow

	var da direktiv.CmdUpdateWorkflow
	da.ID = id
	da.Workflow = data
	da.Active = true

	cmd.Cmd = da

	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		return fmt.Sprintf("Successfully updated '%s'", resp.Cmd.(string)), nil
	}

	return "", errors.New("An unexpected error occurred")
}

// Delete removes a workflow
func Delete(id, namespace string) (string, error) {
	logger.Printf("Deleting workflow '%s'...", id)

	// Open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	// Close nats at end of function
	defer n.Conn.Close()

	var cmd direktiv.CmdRequest

	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.DeleteWorkflow

	var da direktiv.CmdDeleteWorkflow
	da.ID = id

	cmd.Cmd = da

	resp, err := n.DirektivRequest(direktiv.CmdSubscription, id, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		return resp.Cmd.(string), nil
	}

	return "", errors.New("An unexpected error occurred")
}

// Add creates a new workflow on a namespace
func Add(filepath string, namespace string) (string, error) {
	logger.Printf("Adding Workflow...")

	// Open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	// Close nats at end of function
	defer n.Conn.Close()

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	var id = ksuid.New().String()
	var cmd direktiv.CmdRequest
	cmd.CmdID = id
	cmd.CmdType = direktiv.AddWorkflow

	var da direktiv.CmdAddWorkflow
	da.Namespace = namespace
	da.Active = true
	da.Workflow = data

	cmd.Cmd = da

	resp, err := n.DirektivRequest(direktiv.CmdSubscription, id, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		return resp.Cmd.(string), nil
	}

	return "", errors.New("An unexpected error occurred")

}
