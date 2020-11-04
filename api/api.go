package main

import (
	"log"

	"dgrpah-sync-srv/api/handler"
	proto_task "dgrpah-sync-srv/srv/proto/task"
	proto_user "dgrpah-sync-srv/srv/proto/user"

	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.dgraph_sync"),
	)

	// parse command line flags
	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(
			&handler.Dgraph_sync{
				UserClient: proto_user.NewUserService("go.micro.srv.dgraph_sync", service.Client()),
				TaskClient: proto_task.NewTaskService("go.micro.srv.dgraph_sync", service.Client()),
			},
		),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
