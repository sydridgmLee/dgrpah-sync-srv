package mstruct

import (
	"dgrpah-sync-srv/srv/kits"
	"encoding/json"
	"fmt"
)

type Notification struct {
	UID        string   `json:"uid,omitempty""`
	Message    string   `json:"message"`
	Task       []Task   `json:"n_task"`
	SendTo     []User   `json:"send_to"`
	DType      []string `json:"dgraph.type,omitempty"`
	CreateDate int64    `json:"create_date"`
}

func (notification *Notification) DBCreate() error {
	b, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	_, err = kits.DBMutate(b)

	return err
}
