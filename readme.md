
# BlueLogger

BlueLogger is a logger service for [MESG](https://mesg.tech).

  
## Installation
`$ mesg-core service deploy https://github.com/ilgooz/service-bluelogger`


## API
API of the logger service

### Tasks

#### log
Logs input data to standard output.
##### Input Data
* ServiceID `(string)`
ID of the service that `Data` received from.
* Data `(object)`
Actual log message

##### Output Data
* Success: OK `(bool)`
* Error: Message `(string)`


## Sample usage in your MESG Application

```go
package main

import (
    "log"

    mesg "github.com/mesg-foundation/go-application"
)

const (
    discordServiceID = "fill here with the id of the service to log it's task outputs"
    loggerServiceID  = "fill here with the logger service's id"
)

type logRequest struct {
    ServiceID string      `json:"serviceID"`
    Data      interface{} `json:"data"`
}

func main() {
    app, err := mesg.New()
    if err != nil {
        log.Fatal(err)
    }

    // 1- wait for send task's results from Discord service.
    // 2- send the results to log service with service id of Discord.
    _, err := app.
        WhenResult(discordServiceID, mesg.TaskFilterOption("send")).
        Filter(func(r *mesg.Result) bool {
            var resp interface{}
            return r.Data(&resp) == nil
        }).
        Map(func(r *mesg.Result) mesg.Data {
            var resp interface{}
            r.Data(&resp)
            return logRequest{
                ServiceID: discordServiceID,
                Data:      resp,
            }
        }).
        Execute(loggerServiceID, "log")
        
    if err != nil {
        log.Fatal(err)
    }
}
```