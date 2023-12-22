package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)


type Virus struct {
  AptID string //UKNOWN for new threats
  VirusID string
  Hash []map[string]string
  VirusType string //worm,rootkit,trojan
  FileType string //lnk,zip.iso,
  CommunicationMode string //p2p or CnC
  OsType string
  Description string
  utils.TimeStamps
}
