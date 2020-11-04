package mstruct

import "encoding/json"

type Task struct {
	UID         string   `json:"uid,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Priority    string   `json:"priority,omitempty"`
	Status      string   `json:"status,omitempty"`
	Assignee    User     `json:"assignee,omitempty"`
	Reporter    User     `json:"reporter,omitempty"`
	DType       []string `json:"dgraph.type,omitempty"`
}

func (task *Task) GetItemFromBody(body []byte) (err error) {

	err = json.Unmarshal(body, task)

	return
}
