package mstruct

import (
	"dgrpah-sync-srv/srv/kits"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Task struct {
	UID          string         `json:"uid,omitempty"`
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
	Priority     string         `json:"priority,omitempty"`
	Status       string         `json:"status,omitempty"`
	Assignee     []User         `json:"assignee,omitempty"`
	Reporter     []User         `json:"reporter,omitempty"`
	Notification []Notification `json:"notification,omitempty"`
	Err          []SyncError    `json:"err,omitempty"`
	DType        []string       `json:"dgraph.type,omitempty"`
}

type ESTask struct {
	UID         string `json:"uid,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type taskQueryResult struct {
	Tasks []Task `json:"q"`
}

func (task *Task) DBCreate() error {
	b, err := json.Marshal(task)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	response, err := kits.DBMutate(b)

	if err == nil {
		task.UID = response.Uids["task"]
	}

	return err
}

func (task *Task) DBGet() error {
	query := `{
		q(func: uid("` + task.UID + `")){
			uid
			expand(_all_) {
				uid
				expand(_all_) 
			}
			notification: ~n_task {
				uid
				expand(_all_) {
					expand(_all_) 
				}
			}
			err: ~err_task {
				uid
				expand(_all_) {
					expand(_all_) 
				}
			}
		}
	}`

	fmt.Println(query)

	resp, err := kits.DBQuery(query)

	if err == nil {

		r := &taskQueryResult{}
		err = json.Unmarshal(resp.Json, &r)

		if len(r.Tasks) == 1 {
			*task = r.Tasks[0]
		}
	}

	return err
}

func (task *Task) DBDelete() error {
	uid := map[string]string{"uid": task.UID}
	b, err := json.Marshal(uid)
	if err != nil {
		log.Fatal(err)
	}

	return kits.DBDelete(b)
}

func (task *Task) DBUpdate() error {
	var sets []string
	if task.Title != "" {
		set := `<` + task.UID + `> <title> "` + task.Title + `" .`
		sets = append(sets, set)
	}

	if task.Description != "" {
		set := `<` + task.UID + `> <description> "` + task.Description + `" .`
		sets = append(sets, set)
	}

	if task.Priority != "" {
		set := `<` + task.UID + `> <priority> "` + task.Priority + `" .`
		sets = append(sets, set)
	}

	if task.Status != "" {
		set := `<` + task.UID + `> <status> "` + task.Status + `" .`
		sets = append(sets, set)
	}

	set := strings.Join(sets, "\n")

	return kits.DBUpdate(set)
}

func (task *Task) ESUpdate() error {
	endpoint := "/tasks/_doc/" + task.UID

	esTask := ESTask{
		UID:         task.UID,
		Title:       task.Title,
		Description: task.Description,
	}

	b, _ := json.Marshal(esTask)

	_, err := kits.ESReqWithJSON(b, http.MethodPut, endpoint)

	return err
}

func (task *Task) ESDelete() error {
	endpoint := "/tasks/_doc/" + task.UID

	_, err := kits.ESReq("DELETE", endpoint)

	return err
}
