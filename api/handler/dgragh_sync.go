package handler

import (
	"context"
	"dgrpah-sync-srv/api/mstruct"
	proto_task "dgrpah-sync-srv/srv/proto/task"
	proto_user "dgrpah-sync-srv/srv/proto/user"
	"encoding/json"
	"log"

	api "github.com/micro/go-micro/v2/api/proto"

	"github.com/micro/go-micro/errors"
)

type Dgraph_sync struct {
	UserClient proto_user.UserService
	TaskClient proto_task.TaskService
}

func (dgraph_sync *Dgraph_sync) User(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received User API request")

	user := &mstruct.User{}
	user.GetItemFromBody([]byte(req.Body))

	var err error
	var response *proto_user.Response

	switch req.Method {
	case "POST":
		if user.Name == "" || user.Email == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "name and email are required")
		}
		response, err = dgraph_sync.UserClient.Create(ctx, &proto_user.Request{
			Name:  user.Name,
			Email: user.Email,
		})
	case "GET":
		if user.Email == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "email is required")
		}
		response, err = dgraph_sync.UserClient.Get(ctx, &proto_user.Request{
			Email: user.Email,
		})
	case "DELETE":
		if user.Email == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "email is required")
		}
		response, err = dgraph_sync.UserClient.Delete(ctx, &proto_user.Request{
			Email: user.Email,
		})
	case "PUT":
		if user.Email == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "email is required")
		}
		response, err = dgraph_sync.UserClient.Update(ctx, &proto_user.Request{
			Name:  user.Name,
			Email: user.Email,
		})
	}

	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}

func (dgraph_sync *Dgraph_sync) Task(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Print("Received Task API request")

	task := &mstruct.Task{}
	task.GetItemFromBody([]byte(req.Body))

	var err error
	var response *proto_task.Response

	switch req.Method {
	case "POST":
		if task.Title == "" || task.Reporter.Email == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "reporter's email and title are required")
		}
		response, err = dgraph_sync.TaskClient.Create(ctx, &proto_task.Request{
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      task.Status,
			Reporter: &proto_task.User{
				Email: task.Reporter.Email,
			},
			Assignee: &proto_task.User{
				Email: task.Assignee.Email,
			},
		})
	case "GET":
		if task.UID == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "task uid is required")
		}
		response, err = dgraph_sync.TaskClient.Get(ctx, &proto_task.Request{
			Uid: task.UID,
		})
	case "DELETE":
		if task.UID == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "task uid is required")
		}
		response, err = dgraph_sync.TaskClient.Delete(ctx, &proto_task.Request{
			Uid: task.UID,
		})
	case "PUT":
		if task.UID == "" {
			return errors.BadRequest("go.micro.api.dgraph_sync", "task uid is required")
		}
		response, err = dgraph_sync.TaskClient.Update(ctx, &proto_task.Request{
			Uid:         task.UID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    task.Priority,
			Status:      task.Status,
		})
	}
	if err != nil {
		return err
	}

	rsp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": response.Msg,
	})
	rsp.Body = string(b)

	return nil
}
