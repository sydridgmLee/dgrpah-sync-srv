// Package main
package main

import (
	"dgrpah-sync-srv/srv/handler"
	"dgrpah-sync-srv/srv/scheduler"

	proto_task "dgrpah-sync-srv/srv/proto/task"
	proto_user "dgrpah-sync-srv/srv/proto/user"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"
)

func main() {

	// sync err scheduler
	scheduler.FixSyncErr()

	service := micro.NewService(
		micro.Name("go.micro.srv.dgraph_sync"),
	)

	// optionally setup command line usage
	service.Init()

	// Register User Handlers
	proto_user.RegisterUserHandler(service.Server(), new(handler.UserHandler))

	// Register Task Handlers
	proto_task.RegisterTaskHandler(service.Server(), new(handler.TaskHandler))

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
