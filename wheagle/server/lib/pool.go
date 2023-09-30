package lib

import (
  "fmt"
  "time"
  "sync"
  "errors"
  "github.com/alphamystic/odin/lib/utils"
  "github.com/alphamystic/odin/wheagle/server/grpcapi"
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
  utils.PrintTextInASpecificColorInBold("yellow",fmt.Sprintf("Added user with id %s to pool.",id))
  return pool
}

func (p *Pillar) AddWork(id string,work *Jacuzzi) error {
  p.mu.Lock()
	defer p.mu.Unlock()
  for _,pool := range p.Pools{
    if pool.OwnerId == id {
      pool.Work = append(pool.Work,work)
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
  var res = new(grpcapi.Command)
  for _,pool := range p.Pools{
    if len(pool.Work) < 1{
      res.In = "No Work."
      return res,ErrNoWork
    }
    if pool.OwnerId == cmd.UserId{
      for _,jcz := range pool.Work {
        if jcz.Done == false {
          if jcz.CmdIn == "upload" || jcz.CmdIn == "download" || jcz.CmdIn == "screenshot"{
            jcz.Done = true
            res.In = jcz.CmdIn
            res.UserId = cmd.UserId
            res.Individual = true
            res.OperatorId = jcz.OperatorId
            res.Out = jcz.CmdOut
            return res,nil
          }
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
  for _,pool := range p.Pools{
    if pool.OwnerId == cmd.UserId /*&& cmd.In == pool.In*/{
      for _, jcz := range pool.Work{
        if jcz.CmdIn == cmd.In{
          jcz.CmdOut = cmd.Out
          jcz.Done = true
          fmt.Println("User ID marked as done: ",cmd.UserId)
          return nil
        }
      }
      return errors.New("No such a jacuzzi was found.")
    }
  }
  return errors.New("No such pool found.")
}

// Iteratinghere allows not to itterate on every given protocol
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
          fmt.Println("Returning work done.............")
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
      utils.PrintTextInASpecificColorInBold("yellow",fmt.Sprintf("Removed user with id %s from pool ",cmd.UserId))
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

func (p *Pillar) ClearOut(id string) error{
  p.mu.Lock()
	defer p.mu.Unlock()
  var updatedWork []*Jacuzzi
  for _, pool := range p.Pools{
    if pool.OwnerId == id {
      for index, jcz := range pool.Work{
        if jcz.Done{
          updatedWork = append(pool.Work[:index],pool.Work[index+1:]...)
        }
      }
      pool.Work = updatedWork
    }
    fmt.Println("[+]  Deleted all works for ",id)
  }
  return nil
}
