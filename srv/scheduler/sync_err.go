package scheduler

import (
	"dgrpah-sync-srv/srv/kits"
	"dgrpah-sync-srv/srv/mstruct"
	"encoding/json"
	"fmt"

	"github.com/robfig/cron"
)

func FixSyncErr() {
	c := cron.New()
	c.AddFunc("@every 10s", fixSyncErr)
	c.Start()
}

func fixSyncErr() {
	fmt.Println("fixing sync err...")

	syncErrs, err := loadSyncErrs()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, syncErr := range syncErrs {
		if syncErr.Action == "create" || syncErr.Action == "update" {
			esSyncUpdate(&syncErr)
		} else if syncErr.Action == "delete" {
			esSyncDelete(&syncErr)
		}
	}

	fmt.Println("Done")
}

func esSyncUpdate(syncErr *mstruct.SyncError) {
	if len(syncErr.Task) == 1 {
		task := syncErr.Task[0]

		err := task.ESUpdate()

		if err != nil {
			if syncErr.TryTimes == 2 {
				// TODO send notification to admin
			}

			// increase sync err try times
			err = syncErr.IncTryTimes()

			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			// delete sync err
			err = syncErr.DBDelete()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else if len(syncErr.Task) == 0 {
		// delete sync err
		err := syncErr.DBDelete()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func esSyncDelete(syncErr *mstruct.SyncError) {
	task := syncErr.Task[0]

	err := task.ESDelete()

	if err != nil {
		if syncErr.TryTimes == 2 {
			// TODO send notification to admin
		}

		// increase sync err try times
		err = syncErr.IncTryTimes()
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		// delete sync err
		err = syncErr.DBDelete()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func loadSyncErrs() ([]mstruct.SyncError, error) {
	query := `{
		q(func: has(sync_action), orderasc: create_date) {
		  	uid
			expand(_all_) {
				uid
				expand(_all_)
			}
		}
	}`

	resp, err := kits.DBQuery(query)

	fmt.Println(string(resp.Json))

	if err != nil {
		return nil, err
	}

	r := &syncErrQueryResult{}
	err = json.Unmarshal(resp.Json, &r)

	if err != nil {
		return nil, err
	}

	return r.SyncErrs, nil
}

type syncErrQueryResult struct {
	SyncErrs []mstruct.SyncError `json:"q"`
}
