package definers

/*
  * Defines the different kind of ussers that can be able to use odin
*/

import (
  "github.com/alphamystic/odin/lib/utils"
)

type User struct{
  UserID string `json: "userid"`
  OwnerID string  `json: "ownerid"`
  UserName string `json: "username"`
  Email string  `json: "email"`
  Password string `json: "password"`
  Active bool `json: "active"`
  Anonymous bool  `json: "anonymous"`
  Verify bool `json: "verify"`
  Admin bool  `json: admin`
  utils.TimeStamps
}

type UserHash struct {
  UserID string `json: userid`
  Hash string `json: "hash"`
  utils.TimeStamps
}
