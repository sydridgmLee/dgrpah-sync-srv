package mstruct

import "encoding/json"

type User struct {
	UID   string `json:"uid,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (user *User) GetItemFromBody(body []byte) (err error) {

	err = json.Unmarshal(body, user)

	return
}
