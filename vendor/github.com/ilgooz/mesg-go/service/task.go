package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/api/service"
)

// Task represents a MESG task.
type Task struct {
	name    string
	handler func(*Request) Response
}

// NewTask creates a task with name, handler executed when a matching task request received.
func NewTask(name string, handler func(*Request) Response) Task {
	t := Task{
		name:    name,
		handler: handler,
	}
	return t
}

type Key string
type Data interface{}

type Response map[Key]Data

// TaskRequest holds information about a Task request.
type Request struct {
	executionID string
	key         string
	data        string
	service     *Service
}

// Decode decodes task data input to out.
func (t *Request) Decode(out interface{}) error {
	return json.Unmarshal([]byte(t.data), out)
}

func (t *Request) reply(resp Response) error {
	for key, data := range resp {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = t.service.client.SubmitResult(context.Background(), &service.SubmitResultRequest{
			ExecutionID: t.executionID,
			OutputKey:   string(key),
			OutputData:  string(dataBytes),
		})
		return err
	}
	return nil
}
