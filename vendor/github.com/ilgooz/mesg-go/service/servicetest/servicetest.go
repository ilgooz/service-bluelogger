package servicetest

// import (
// 	"encoding/json"

// 	"github.com/mesg-foundation/core/api/service"
// 	uuid "github.com/satori/go.uuid"
// )

// type Service struct {
// 	client *Client
// }

// func New() *Service {
// 	return &Service{
// 		client: &Client{
// 			emitC:   make(chan *service.EmitEventRequest, 0),
// 			submitC: make(chan *service.SubmitResultRequest, 0),
// 			stream:  &taskDataStream{taskC: make(chan *service.TaskData, 0)},
// 		},
// 	}
// }

// func (s *Service) Client() *Client {
// 	return s.client
// }

// func (s *Service) LastEmit() *Emit {
// 	e := <-s.client.emitC
// 	return &Emit{
// 		name:  e.EventKey,
// 		data:  e.EventData,
// 		token: e.Token,
// 	}
// }

// type Emit struct {
// 	name  string
// 	data  interface{}
// 	token string
// }

// func (e *Emit) Name() string {
// 	return e.name
// }

// func (e *Emit) Data() interface{} {
// 	return e.data
// }

// func (e *Emit) Token() string {
// 	return e.token
// }

// type Execution struct {
// 	id     string
// 	client *Client
// }

// func (s *Service) Execute(task string, data interface{}) (*Execution, error) {
// 	bytes, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	id := uuid.NewV4().String()
// 	s.client.stream.taskC <- &service.TaskData{
// 		ExecutionID: id,
// 		TaskKey:     task,
// 		InputData:   string(bytes),
// 	}
// 	return &Execution{
// 		id:     id,
// 		client: s.client,
// 	}, nil
// }

// func (e *Execution) ID() string {
// 	return e.id
// }

// type Response struct {
// }

// func (e *Execution) Response() error {
// 	data := <-e.client.submitC
// 	data.
// 	return
// }
