package mstruct

import (
	"dgrpah-sync-srv/srv/kits"
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	UID          string         `json:"uid,omitempty"`
	Name         string         `json:"name,omitempty"`
	Email        string         `json:"email,omitempty"`
	DType        []string       `json:"dgraph.type,omitempty"`
	WorkOn       []Task         `json:"work_on,omitempty"`
	Manage       []Task         `json:"manage,omitempty"`
	Notification []Notification `json:"notification,omitempty"`
}

type userQueryResult struct {
	Users []User `json:"q"`
}

func (user *User) DBCreate() error {
	b, err := json.Marshal(user)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	_, err = kits.DBMutate(b)

	return err
}

func (user *User) DBGet() error {
	query := `
		{
			q(func: eq(email, "` + user.Email + `")) {
				uid
				name
				email
				dgraph.type
				work_on: ~assignee {
					uid
					expand(_all_) {
						expand(_all_) 
					}
				}
				manage: ~reporter {
					uid
					expand(_all_) {
						expand(_all_) 
					}
				}
				notification: ~send_to {
					uid
					expand(_all_) {
						expand(_all_) 
					}
				}
			}
		}
	`

	resp, err := kits.DBQuery(query)

	if err == nil {

		r := &userQueryResult{}
		err = json.Unmarshal(resp.Json, &r)

		if len(r.Users) == 1 {
			*user = r.Users[0]
		}
	}

	return err
}

func (user *User) DBDelete() error {
	uid := map[string]string{"uid": user.UID}
	b, err := json.Marshal(uid)
	if err != nil {
		log.Fatal(err)
	}

	return kits.DBDelete(b)
}

func (user *User) DBUpdate() error {

	set := ""
	if user.Name != "" {
		set = `uid(v) <name> "` + user.Name + `" .`
	}

	query := `
		query {
			me(func: eq(email, "` + user.Email + `")) {
				v as uid
			}
		}
	`

	return kits.DBUpdateWithQuery(query, set)
}
