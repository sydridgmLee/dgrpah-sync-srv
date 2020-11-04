package handler

import (
	"context"
	"dgrpah-sync-srv/srv/mstruct"
	proto_task "dgrpah-sync-srv/srv/proto/task"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-log/log"
)

type TaskHandler struct{}

func (taskHandler *TaskHandler) Create(ctx context.Context, req *proto_task.Request, rsp *proto_task.Response) error {
	log.Log("Received User.Create request")

	// get reporter
	reporter := &mstruct.User{
		Email: req.Reporter.Email,
	}

	err := reporter.DBGet()
	if err != nil {
		rsp.Msg = "error: " + err.Error()
		return nil
	}

	if reporter.UID == "" {
		// user exist
		rsp.Msg = "error: reporter not exist"
		return nil
	}

	// get asignee
	assignee := &mstruct.User{
		Email: req.Assignee.Email,
	}

	err = assignee.DBGet()
	if err != nil {
		rsp.Msg = "error: " + err.Error()
		return nil
	}

	if assignee.UID == "" {
		// user exist
		rsp.Msg = "error: assignee not exist"
		return nil
	}

	// create task
	task := &mstruct.Task{
		UID:         "_:task",
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DType:       []string{"Task"},
		Assignee: []mstruct.User{
			mstruct.User{
				UID: assignee.UID,
			},
		},
		Reporter: []mstruct.User{
			mstruct.User{
				UID: reporter.UID,
			},
		},
	}

	// add task to db
	err = task.DBCreate()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		//sync ES
		go esSynCreate(task)

		// send notification
		go sendNotification("task created", true, task)

		rsp.Msg = "success"
	}

	return nil
}

func (taskHandler *TaskHandler) Get(ctx context.Context, req *proto_task.Request, rsp *proto_task.Response) error {
	log.Log("Received User.Get request")

	task := &mstruct.Task{
		UID: req.Uid,
	}

	err := task.DBGet()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		if task.UID == "" {
			rsp.Msg = "not found"
		} else {
			b, _ := json.Marshal(task)
			rsp.Msg = string(b)
		}
	}

	return nil
}

func (taskHandler *TaskHandler) Delete(ctx context.Context, req *proto_task.Request, rsp *proto_task.Response) error {
	log.Log("Received User.Delete request")

	task := &mstruct.Task{
		UID: req.Uid,
	}

	err := task.DBDelete()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		// sync es
		go esSyncDelete(task)

		rsp.Msg = "success"
	}

	return nil
}

func (taskHandler *TaskHandler) Update(ctx context.Context, req *proto_task.Request, rsp *proto_task.Response) error {
	log.Log("Received User.Update request")

	task := &mstruct.Task{
		UID:         req.Uid,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
	}

	err := task.DBUpdate()

	if err != nil {
		rsp.Msg = "error: " + err.Error()
	} else {
		// sync es
		go esSyncUpdate(task)

		// send notification
		go sendNotification("task updated", false, task)

		rsp.Msg = "success"
	}

	return nil
}

// Sync ES
func esSynCreate(task *mstruct.Task) {
	err := task.ESUpdate()

	if err != nil {
		addSyncErr2DB("create", task.UID)
	}
}

func esSyncDelete(task *mstruct.Task) {
	err := task.ESDelete()

	if err != nil {
		addSyncErr2DB("delete", task.UID)
	}
}

func esSyncUpdate(task *mstruct.Task) {
	err := task.DBGet() // read updated task first
	if err != nil {
		addSyncErr2DB("update", task.UID)
	}

	err = task.ESUpdate()
	if err != nil {
		addSyncErr2DB("update", task.UID)
	}
}

func addSyncErr2DB(action, taskUid string) {
	// add to sync error to db
	syncErr := &mstruct.SyncError{
		Action: action,
		Task: []mstruct.Task{
			mstruct.Task{
				UID: taskUid,
			},
		},
		TryTimes:   1,
		DType:      []string{"SynError"},
		CreateDate: time.Now().UTC().Unix(),
	}

	err := syncErr.DBCreate()

	if err != nil {
		//TODO sync error locally and send notification to admin
	}
}

// send notification
func sendNotification(message string, isNewTask bool, task *mstruct.Task) {
	if !isNewTask {
		err := task.DBGet()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	sendTo := append(task.Assignee, task.Reporter...)

	notification := &mstruct.Notification{
		Message: message,
		Task: []mstruct.Task{
			mstruct.Task{
				UID: task.UID,
			},
		},
		SendTo:     sendTo,
		DType:      []string{"Notification"},
		CreateDate: time.Now().UTC().Unix(),
	}

	err := notification.DBCreate()

	if err != nil {
		fmt.Println(err.Error())
	}
}
