package servicetest

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
)

var ErrClosedConn = errors.New("closed connection")

type Client struct {
	stream  *taskDataStream
	emitC   chan *service.EmitEventRequest
	submitC chan *service.SubmitResultRequest
}

func (c *Client) EmitEvent(ctx context.Context, in *service.EmitEventRequest,
	opts ...grpc.CallOption) (*service.EmitEventReply, error) {
	c.emitC <- in
	return nil, nil
}

func (c *Client) ListenTask(ctx context.Context,
	in *service.ListenTaskRequest,
	opts ...grpc.CallOption) (service.Service_ListenTaskClient, error) {
	return c.stream, nil
}

func (c *Client) SubmitResult(ctx context.Context,
	in *service.SubmitResultRequest,
	opts ...grpc.CallOption) (*service.SubmitResultReply, error) {
	c.submitC <- in
	return nil, nil
}

type taskDataStream struct {
	taskC chan *service.TaskData
	grpc.ClientStream
}

func (s taskDataStream) Recv() (*service.TaskData, error) {
	data, ok := <-s.taskC
	if !ok {
		return nil, ErrClosedConn
	}
	return data, nil
}
