package core


import (
  "fmt"
  "sync"
  "errors"
  "github.com/alphamystic/odin/lib/utils"
)
/*
  * The cave just create a channelfor everyone to manipulate or access data returned  from it.
  * Basic rules,
      Read and write to your tunnel as much you want BUT only close it from the cave
*/
type Cave struct{
  mu sync.RWMutex
  Tunnels map[string]*Tunnel
}

type Tunnel struct{
  Id string
  Data chan interface{}
}

func InitializeTunnelMan()*Cave{
  return &Cave{
    mu: sync.RWMutex{},
    Tunnels: make(map[string]*Tunnel),
  }
}

func (cave *Cave) CreateTunnel()(*Tunnel){
  cave.mu.RLock()
	defer cave.mu.RUnlock()
  tunnel := &Tunnel{
    Id: utils.RandString(10),
    Data: make(chan interface{}),
  }
  for _,tun := range cave.Tunnels{
    if tun.Id == tunnel.Id{
      tunnel.Id = utils.RandString(11)
    }
  }
  cave.Tunnels[tunnel.Id] = tunnel
  return tunnel
}

func (cave *Cave) GetTunnel(id string) (*Tunnel,error){
  cave.mu.RLock()
	defer cave.mu.RUnlock()
  tunnel,ok := cave.Tunnels[id]
  if !ok {
    return nil,fmt.Errorf("No tunnnel with such id as %s",id)
  }
  return tunnel,nil
}

func (cave *Cave) CloseTunnel(id string) error{
  cave.mu.RLock()
	defer cave.mu.RUnlock()
  for _,tun := range cave.Tunnels{
    if tun.Id == id{
      close(tun.Data)
      delete(cave.Tunnels,tun.Id)
    }
  }
  return errors.New("Error, no channel with such an id")
}
