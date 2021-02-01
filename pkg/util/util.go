package util

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/vorteil/direkcli/pkg/log"
	"github.com/vorteil/direktiv/pkg/direktiv"
	"github.com/vorteil/vorteil/pkg/elog"
)

var logger elog.View

func init() {
	log := log.GetLogger()
	logger = log
}

const (
	// NatsRequestRetries ...
	NatsRequestRetries = 3
	// NatsRequestRetryCooldown ..
	NatsRequestRetryCooldown = 5 * time.Millisecond
	// NatsRequestDefaultTimeout ...
	NatsRequestDefaultTimeout = 500 * time.Millisecond
)

type NatsHandler struct {
	Conn *nats.Conn
}

func JsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func CreateNatsConnection(ip string) (*NatsHandler, error) {
	nc, err := nats.Connect(ip, nats.Token("pwd")) // todo token flag for nats?
	if err != nil {
		return nil, err
	}
	return &NatsHandler{
		Conn: nc,
	}, err
}

func (n *NatsHandler) DirektivUnmarshal(request *direktiv.CmdRequest, target interface{}) error {
	data, err := json.Marshal(request.Cmd)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

func (n *NatsHandler) DirektivRequest(subj string, id string, cmd interface{}) (*direktiv.CmdRequest, error) {
	m := new(nats.Msg)

	b, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	for i := 0; i < NatsRequestRetries; i++ {
		logger.Debugf("Attempting direktiv request... (id=%s, subject=%s, attempt=%v)\n", id, subj, i+1)
		m, err = n.Conn.Request(subj, b, NatsRequestDefaultTimeout)
		if err == nil {
			break
		}
		time.Sleep(NatsRequestRetryCooldown)
	}

	resp := new(direktiv.CmdRequest)
	err = json.Unmarshal(m.Data, &resp)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func CmdErrorCheck(resp *direktiv.CmdRequest) *direktiv.CmdErrorResponse {
	if resp.CmdType != direktiv.Error {
		return nil
	}

	val := new(direktiv.CmdErrorResponse)
	b, err := json.Marshal(resp.Cmd)
	if err != nil {
		panic(err) // This should never happen
	}
	err = json.Unmarshal(b, &val)
	if err != nil {
		panic(err) // this should never happen
	}

	return val
}
