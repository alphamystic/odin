package definers

import (
  "github.com/alphamystic/odin/lib/utils"
)


// When caaling Minion, take what yuou need alone
type Minion struct {
  MinionID string `json:"minionid"`
  Name string `json:"name"`
  UName string `json:"uname"`
  UserID string  `json:"uid"`
  GroupID string  `json:"groupid"`
  HomeDir string `json:"homedir"`
  Os  string  `json:"ostype"`
  Description string `json:"description"`
  Installed bool `json:"installed"`
  MothershipID string `json:"mothershipid"`
  Address string `json:"address"`// this is set to string just incase the address is a domain name
  Port string `json:"port"`
  Motherships string  `json:"motherships"` // keep it as a bunch of strings and initialize it when needed
  TunnelAddress string `json:"tunnel_address"`
  Tls bool `json:"tls"`
  //Motherships []map[string]string  `json:"motherships"`
  OwnerID string `json:"owner_id"`
  LastSeen string `json:"lastseen"`
  IsDropper bool `json:"is_dropper"`
  GenCommand string `json:"generate_command"`
  utils.TimeStamps
}

/*
  * motherships are a bunch of []Server
  * tunneling address is also a Server type
*/
