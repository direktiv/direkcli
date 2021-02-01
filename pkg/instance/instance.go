package instance

import (
	"encoding/json"
	"errors"

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

// Logs returns all logs associated with the workflow instance ID
func Logs(id string) (string, error) {
	logger.Printf("Fetching logs for '%s'...", id)

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	defer n.Conn.Close()

	var cmd direktiv.CmdRequest
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.GetInstanceAllLogs

	var da direktiv.CmdGetInstanceAllLogs
	da.InstanceID = id

	cmd.Cmd = da
	// send request to direktiv
	resp, err := n.DirektivRequest(direktiv.CmdSubscription, cmd.CmdID, cmd)
	if err != nil {
		return "", err
	}

	// Check for response if error
	if dirErr := util.CmdErrorCheck(resp); dirErr != nil {
		return "", errors.New(dirErr.Error)
	}

	if resp.CmdType == direktiv.OK {
		logs := new(direktiv.CmdGetInstanceLogsResponse)
		if err := n.DirektivUnmarshal(resp, &logs); err != nil {
			return "", err
		}

		b, err := json.Marshal(logs.Logs.Data)
		if err != nil {
			return "", err
		}
		return util.JsonPrettyPrint(string(b)), err
	}

	return "", errors.New("An unexpected error occurred")
}

// List returns a list of workflow instances
func List(namespace string) ([]direktiv.CmdGetWorkflowInstancesResponse, error) {
	logger.Printf("Fetching workflow instance list...")

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return nil, err
	}
	defer n.Conn.Close()

	var id = ksuid.New().String()
	var cmd direktiv.CmdRequest
	cmd.CmdID = id
	cmd.CmdType = direktiv.GetWorkflowInstances

	var da direktiv.CmdGetWorkflowInstances
	da.Namespace = namespace
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
		wfinstances := make([]direktiv.CmdGetWorkflowInstancesResponse, 0)
		if err := n.DirektivUnmarshal(resp, &wfinstances); err != nil {
			return nil, err
		}

		return wfinstances, nil
	}

	return nil, errors.New("An unexpected error occurred")
}

// Get returns a pretty printed json of the workflow instance id
func Get(id string) (string, error) {
	logger.Printf("Fetching workflow instance '%s'...", id)

	// open nats
	n, err := util.CreateNatsConnection("192.168.43.128:4222")
	if err != nil {
		return "", err
	}
	defer n.Conn.Close()

	var cmd direktiv.CmdRequest
	var db direktiv.CmdGetWorkflowInstance
	cmd.CmdID = ksuid.New().String()
	cmd.CmdType = direktiv.GetWorkflowInstance

	db.InstanceID = id
	cmd.Cmd = db

	// send request to direktiv
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
		wfinstance := new(direktiv.CmdGetWorkflowInstanceResponse)
		if err := n.DirektivUnmarshal(resp, &wfinstance); err != nil {
			return "", err
		}

		b, err := json.Marshal(wfinstance)
		if err != nil {
			return "", err
		}

		return util.JsonPrettyPrint(string(b)), nil
	}

	return "", errors.New("An unexpected error occurred")
}
