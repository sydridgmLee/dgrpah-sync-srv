package mstruct

import (
	"dgrpah-sync-srv/srv/kits"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type SyncError struct {
	UID        string   `json:"uid"`
	Action     string   `json:"sync_action"`
	Task       []Task   `json:"err_task"`
	TryTimes   int      `json:"try_times"`
	DType      []string `json:"dgraph.type,omitempty"`
	CreateDate int64    `json:"create_date"`
}

func (syncError *SyncError) DBCreate() error {
	b, err := json.Marshal(syncError)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	_, err = kits.DBMutate(b)

	return err
}

func (syncError *SyncError) DBDelete() error {
	uid := map[string]string{"uid": syncError.UID}
	b, err := json.Marshal(uid)
	if err != nil {
		log.Fatal(err)
	}

	return kits.DBDelete(b)
}

func (syncError *SyncError) IncTryTimes() error {
	times := syncError.TryTimes + 1
	fmt.Println(times)
	set := `<` + syncError.UID + `> <try_times> "` + strconv.Itoa(times) + `" .`
	return kits.DBUpdate(set)
}
