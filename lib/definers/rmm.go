package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type RmmTask struct {
    TaskID         string   `json:"task_id"`
    OwnerID        string   `json:"ownerid"`
    Command string   `json:"command"`
    TargetTpe string `json:"target_type"` //custom, linux, windows,android
    MinionIDs      []string `json:"minion_ids"` // Target list Store this in base 64 in in db
    utils.TimeStamps
}

type RmmExecution struct {
    ExecutionID      string `json:"execution_id"`
    TaskID           string `json:"task_id"`
    MinionID         string `json:"minionid"`
    Status           string `json:"status"` // success, failed
    Output           string `json:"output"`
    VerificationHash string `json:"verification_hash"` // CommandID + MinionID hash
    utils.TimeStamps
}