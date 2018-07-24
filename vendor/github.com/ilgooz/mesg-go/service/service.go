// Package service is a service client for mesg-core.
// For more information please visit https://mesg.com.
package service

import (
	"errors"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"encoding/json"

	"context"

	"github.com/ilgooz/mesg-go/service/servicetest"
	"github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
)

const tcpEndpointEnv = "MESG_ENDPOINT_TCP"
const tokenEnv = "MESG_TOKEN"

// Service represents a MESG service.
type Service struct {
	endpoint string
	token    string

	// Client is the gRPC service client of MESG.
	client service.ServiceClient
	conn   *grpc.ClientConn

	callTimeout time.Duration

	tasks map[string]Task
	mt    sync.RWMutex

	log       *log.Logger
	logOutput io.Writer
}

// Option is the configuration func of Service.
type Option func(*Service)

// New starts a new Service with options.
func New(options ...Option) (*Service, error) {
	s := &Service{
		endpoint:    os.Getenv(tcpEndpointEnv),
		token:       os.Getenv(tokenEnv),
		tasks:       map[string]Task{},
		callTimeout: time.Second * 10,
		logOutput:   os.Stdout,
	}
	for _, option := range options {
		option(s)
	}
	s.log = log.New(s.logOutput, "mesg", log.LstdFlags)
	if s.client == nil {
		if s.endpoint == "" {
			return nil, errors.New("endpoint is not set")
		}
		if s.token == "" {
			return nil, errors.New("token is not set")
		}
		if err := s.setupServiceClient(); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// EndpointOption receives the TCP endpoint of MESG.
func EndpointOption(address string) Option {
	return func(s *Service) {
		s.endpoint = address
	}
}

// TokenOption receives token which is the unique id of this service.
func TokenOption(token string) Option {
	return func(s *Service) {
		s.token = token
	}
}

// TimeoutOption receives d to use while dialing mesg-core and making requests.
func TimeoutOption(d time.Duration) Option {
	return func(s *Service) {
		s.callTimeout = d
	}
}

// LogOutputOption uses out as a log destination.
func LogOutputOption(out io.Writer) Option {
	return func(s *Service) {
		s.logOutput = out
	}
}

func MockOption(client *servicetest.Client) Option {
	return func(s *Service) {
		s.client = client
	}
}

func (s *Service) setupServiceClient() error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	s.conn, err = grpc.DialContext(ctx, s.endpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.client = service.NewServiceClient(s.conn)
	return nil
}

// Listen blocks while listening for tasks.
func (s *Service) Listen(task Task, tasks ...Task) error {
	s.mt.Lock()
	if len(s.tasks) > 0 {
		s.mt.Unlock()
		return errors.New("tasks already set")
	}
	s.tasks[task.name] = task
	for _, task := range tasks {
		s.tasks[task.name] = task
	}
	s.mt.Unlock()
	if err := s.validateTasks(); err != nil {
		return err
	}
	return s.listenTasks()
}

func (s *Service) validateTasks() error { return nil }

func (s *Service) listenTasks() error {
	stream, err := s.client.ListenTask(context.Background(), &service.ListenTaskRequest{
		Token: s.token,
	})
	if err != nil {
		return err
	}
	for {
		data, err := stream.Recv()
		if err != nil {
			return err
		}
		s.executeTask(data)
	}
}

func (s *Service) executeTask(data *service.TaskData) {
	s.mt.RLock()
	fn, ok := s.tasks[data.TaskKey]
	s.mt.RUnlock()
	if !ok {
		return
	}
	req := &Request{
		executionID: data.ExecutionID,
		key:         data.TaskKey,
		data:        data.InputData,
		service:     s,
	}
	if err := req.reply(fn.handler(req)); err != nil {
		s.log.Println(err)
	}
}

// Emit emits a MESG event with given data for name.
func (s *Service) Emit(event string, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.callTimeout)
	defer cancel()
	_, err = s.client.EmitEvent(ctx, &service.EmitEventRequest{
		Token:     s.token,
		EventKey:  event,
		EventData: string(dataBytes),
	})
	return err
}

// Close gracefully closes underlying connections and stops listening for tasks.
func (s *Service) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
