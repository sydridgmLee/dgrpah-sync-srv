package handler

import (
	"context"
	"dgrpah-sync-srv/srv/mstruct"
	proto_user "dgrpah-sync-srv/srv/proto/user"
	"encoding/json"

	"github.com/go-log/log"
)

type UserHandler struct{}

func (userHandler *UserHandler) Create(ctx context.Context, req *proto_user.Request, rsp *proto_user.Response) error {
	log.Log("Received User.Create request")

	user := &mstruct.User{
		Email: req.Email,
		Name:  req.Name,
		DType: []string{"User"},
	}

	// check if email already exist
	err := user.DBGet()
	if err != nil {
		rsp.Msg = "error: " + err.Error()
		return nil
	}

	if user.UID != "" {
		// user exist
		rsp.Msg = "error: user already exist"
		return nil
	}

	// insert to db
	err = user.DBCreate()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		rsp.Msg = "success"
	}

	return nil
}

func (userHandler *UserHandler) Get(ctx context.Context, req *proto_user.Request, rsp *proto_user.Response) error {
	log.Log("Received User.Get request")

	user := &mstruct.User{
		Email: req.Email,
	}

	err := user.DBGet()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		if user.UID == "" {
			rsp.Msg = "not found"
		} else {
			b, _ := json.Marshal(user)
			rsp.Msg = string(b)
		}
	}

	return nil
}

func (userHandler *UserHandler) Delete(ctx context.Context, req *proto_user.Request, rsp *proto_user.Response) error {
	log.Log("Received User.Delete request")

	user := &mstruct.User{
		Email: req.Email,
	}

	// check if user exist
	err := user.DBGet()
	if err != nil {
		rsp.Msg = "error: " + err.Error()
		return nil
	}

	if user.UID == "" {
		// user exist
		rsp.Msg = "error: user not exist"
		return nil
	}

	err = user.DBDelete()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		rsp.Msg = "success"
	}

	return nil
}

func (userHandler *UserHandler) Update(ctx context.Context, req *proto_user.Request, rsp *proto_user.Response) error {
	log.Log("Received User.Update request")

	user := &mstruct.User{
		Email: req.Email,
		Name:  req.Name,
	}

	err := user.DBUpdate()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		rsp.Msg = "success"
	}

	return nil
}
