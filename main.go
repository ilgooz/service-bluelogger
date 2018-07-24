package main

import (
	"encoding/json"
	"log"

	"github.com/ilgooz/mesg-go/service"
)

func main() {
	s, err := service.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Listen(
		service.NewTask("log", handler),
	); err != nil {
		log.Fatal(err)
	}
}

func handler(req *service.Request) service.Response {
	var data logRequest
	if err := req.Decode(&data); err != nil {
		return service.Response{
			"error": errorResponse{err.Error()},
		}
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		return service.Response{
			"error": errorResponse{err.Error()},
		}
	}

	log.Printf("%s: %s", data.ServiceID, string(bytes))

	return service.Response{
		"success": successResponse{"ok"},
	}
}

type logRequest struct {
	ServiceID string      `json:"serviceID"`
	Data      interface{} `json:"data"`
}

type successResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}
