package c2

import (
  "fmt"
  "sync"
  "encoding/json"
  "odin/lib/db"
  "odin/lib/utils"
)
/*
  This package defines all the various methods for teaming up
  We can add extra checks to the mothership to ensure only the specified operators can interact with it
  though the password does that well enough.
  So we probably just need this to log commands out or whatever the future brings
  Moving it to C2................ :)
  // should be loaded by the C2 on startup or whatever, still thinking on how I would implement it right.
*/


type Operator struct {
  Name string
  ID string
  Active bool
  CreatedAt string
  UpdatedAt string
}

type OAdmin struct{
  //Dir string
  Password string
  Driver *db.Driver
  Operators map[string]*Operator
  mu sync.RWMutex
}
//dir := "../brain/operators"
func OPManager(dir,passwd string)*OAdmin{
  driver,err := db.Old(dir,0644)
  if err != nil{
    utils.Logerror(err);return  nil
  }
  //load operators from db (init it out at C2)
  return &OAdmin{
    //Dir: "../brain/",
    Password: passwd,
    Driver: driver,
    Operators: make(map[string]*Operator),
    mu: sync.RWMutex{},
  }
}

func (opm *OAdmin) AddOperator(optr *Operator) error{
  opm.mu.Lock()
	defer opm.mu.Unlock()
  for _,opr := range opm.Operators{
    if opr.ID == optr.ID{
      return fmt.Errorf("Operator with ID already exist.")
    }
  }
  opm.Operators[optr.ID] = optr
  return nil
}

func (opm *OAdmin) RemoveOerator(id,passwd string)error{
  opm.mu.Lock()
	defer opm.mu.Unlock()
  if err := utils.CheckPasswordHash(passwd,opm.Password); err != nil{
    fmt.Errorf("Wrong admin password provided. \nERROR: ",err)
  }
  for _,ops := range opm.Operators{
    if ops.ID == id {
      if err :=opm.Driver.Delete("operators",id); err != nil{
        return fmt.Errorf("Error deleting Operator from db. \nERROR: %s",err)
      }
      delete(opm.Operators,id)
      return nil
    }
  }
  return nil
}

func (opm *OAdmin) ActivateDeactivate(id string,value bool)error{
  opm.mu.Lock()
	defer opm.mu.Unlock()
  for _,ops := range opm.Operators {
    if ops.ID == id{
      ops.Active = value
      return nil
    }
  }
  return fmt.Errorf("Error activating or deactivating operator with id %s.",id)
}
func (opm *OAdmin) LoadOperators()error{
  oprts,err := opm.Driver.ReadAll("operators")
  if err != nil{
    return fmt.Errorf("Error reading operators from DB.\nERROR: %s",err)
  }
  for _,ops := range oprts{
    var opr Operator
    if err := json.Unmarshal([]byte(ops),&opr); err != nil{
      utils.Logerror(err);continue
    }
    if err := opm.AddOperator(&opr); err != nil{
      utils.Logerror(fmt.Errorf("Error adding ooperator to operators.\nERROR: %s",err));continue
    }
  }
  return nil
}
