package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)

type IOC struct {
  IocID int
  VirusID string
  Type string
  Value string
  Source string
  utils.TimeStamps
}
