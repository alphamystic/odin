package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type Plugin struct {
  Owner string
  Name string
  Hash string
  PType  int
  Decsription string
  Verified bool
  Active bool
  Signed bool
  utils.TimeStamps
}
