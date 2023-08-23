package main

import (
  "log"
  "fmt"
  "time"
  "sync"
  "errors"
  "odin/wheagle/server/grpcapi"
)

//instead of jacuzzi I would have just used Work []*grpcapi.Command
// but now I'm in too deap I guess..............

type Pool struct{
  OwnerId string
  Active bool
  Work []*Jacuzzi
}

type Jacuzzi struct{
  OperatorId string
  CmdIn string
  CmdOut string
  Done bool
}

type Pillar struct{
  Pools map[string]*Pool
  mu sync.RWMutex
}

var ErrNoWork = errors.New("Error: No work for user.")
func InitializePillar()*Pillar{
  return &Pillar{
    Pools: make(map[string]*Pool),
		mu: sync.RWMutex{},
  }
}

func (p *Pillar) NewPool(id string) *Pool{
  p.mu.Lock()
	defer p.mu.Unlock()
  //works := []*Jacuzzi
  pool := &Pool{
    OwnerId: id,
    Active: true,
    //Work: works,
  }
  p.Pools[pool.OwnerId] = pool
  fmt.Sprintf("Added user with id %s to pool.",id)
  return pool
}

func (p *Pillar) AddWork(id string,work *Jacuzzi) error {
  p.mu.Lock()
	defer p.mu.Unlock()
  for _,pool := range p.Pools{
    if pool.OwnerId == id {
      pool.Work = append(pool.Work,work)
      fmt.Println("Added work ",work.CmdIn)
      return nil
    } else {
      return fmt.Errorf("(ERROR_ADDING_WORK) No pool with such an id as %s",id)
    }
  }
  return nil
}

func (p *Pillar) GetWork(cmd *grpcapi.Command) (*grpcapi.Command,error){
  p.mu.Lock()
	defer p.mu.Unlock()
  fmt.Println("So we need work for ",cmd.UserId)
  var res = new(grpcapi.Command)
  for _,pool := range p.Pools{
    fmt.Println("Ranging through pools: ",pool.OwnerId)
    if pool.OwnerId == cmd.UserId{
      fmt.Println("Maybe we atleast found the pool...",pool.OwnerId,cmd.UserId)
      if len(pool.Work) < 1{
        return res,ErrNoWork
      }
      for _,jcz := range pool.Work {
        if jcz.Done == false {
          res.In = jcz.CmdIn
          res.UserId = cmd.UserId
          res.Individual = true
          return res,nil
        } else {
          return res,ErrNoWork
        }
      }
    }
  }
  return res,fmt.Errorf("(ERROR_GETTING_WORK) No pool with such an id as %s",cmd.UserId)
}

func (p *Pillar) GetMyPool(cmd *grpcapi.Command)(*Pool,error){
  p.mu.Lock()
	defer p.mu.Unlock()
  for _,pool := range p.Pools{
    if pool.OwnerId == cmd.UserId {
      //pool.In = cmd.In
      return pool,nil
    }
  }
  return nil,errors.New("No work belonging to user present.")
}

func (p *Pillar) MarkAsDone(cmd *grpcapi.Command) error {
  p.mu.Lock()
	defer p.mu.Unlock()
  fmt.Println("User ID marked as done: ",cmd.UserId)
  for _,pool := range p.Pools{
    if pool.OwnerId == cmd.UserId /*&& cmd.In == pool.In*/{
      for _, jcz := range pool.Work{
        if jcz.CmdIn == cmd.In{
          jcz.CmdOut = cmd.Out
          jcz.Done = true
          return nil
        }
      }
      return errors.New("No such a jacuzzi was found.")
    }
  }
  return errors.New("No such pool found.")
}

func (p *Pillar) GetWorkDone(cmd *grpcapi.Command)(workDone *grpcapi.Command,err error){
  p.mu.RLock()
  defer p.mu.RUnlock()
  var count int
  BEGIN:
  for _,pool := range p.Pools{
    if pool.OwnerId == cmd.UserId{
      for _, jcz := range pool.Work{
        if jcz.CmdIn == cmd.In && jcz.Done{
          workDone = &grpcapi.Command{
            In: jcz.CmdIn,
            Out: jcz.CmdOut,
            Individual: true,
            UserId: pool.OwnerId,
          }
        //purge the work done to clear mem of arrays holding the work
        //delete(pool.Work,)
        return workDone,nil
        } else {
          if !jcz.Done{
            if count >= 3{
              return nil,errors.New("Theres an issue making the work done probably.")
            }
            count += 1
            time.Sleep(3 *time.Second)
            goto BEGIN
          }
        }
      }
    }
  }
  err = errors.New("Probably their was an error marking minion as inactive so server assumed it is still alive and caches it's commands.")
  return nil,err
}
//deactivate at send output
func (p *Pillar) Deactivate(cmd *grpcapi.Command) error{
  p.mu.Lock()
	defer p.mu.Unlock()
  for _, pool := range p.Pools{
    if pool.OwnerId == cmd.UserId && cmd.In == "exit"{
      /*pool.Active = false
      pool.Done = true*/
      delete(p.Pools,cmd.UserId)
      fmt.Sprintf("Removed user with id %s from pool ",cmd.UserId)
      return nil
    }
  }
  return errors.New("No such pool found")
}

func (p *Pillar) IsActive(cmd *grpcapi.Command) bool {
  p.mu.Lock()
	defer p.mu.Unlock()
  for _, pool := range p.Pools{
    if pool.OwnerId == cmd.UserId{
      if pool.Active { return true}
    }
  }
  return false
}
/*
type Pool struct{
  OwnerId string
  Active bool
  Work []*Jacuzzi
}

type Jacuzzi struct{
  OperatorId string
  CmdIn string
  CmdOut string
  Done bool
}*/

func main(){
  pillar := InitializePillar()
  pool := pillar.NewPool("SAM")
  pool2 := pillar.NewPool("ODHIAMBO")
  fmt.Println("Printing pool.....",pool)
  fmt.Println("Printing pool.....",pool2)
  cmd1 := &grpcapi.Command{
    UserId: "SAM",
  }
  work,err := pillar.GetWork(cmd1)
  if err != nil{
    if errors.Is(err,ErrNoWork){
      log.Fatal("No work found")
    }
    log.Fatal(err)
  }
  fmt.Println("SAM WORK: ",work)
}
