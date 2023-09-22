package c2

/*
  This are connections to different/active motherships.
  When the cli need to send a command to a given ms, A client (Connector)
  is created and sends the message return ning the output
*/

import (
  "fmt"
  "sync"
  "errors"
  "strings"
//  "io/ioutil"
//  "path/filepath"
  "encoding/json"

	"odin/lib/db"
  "odin/lib/utils"
)

const ConnectorsPath = "../.brain/"

type Connector struct{
  SessionId string
  IAddress string
  OAddress string
  Name string
  Tls bool
  Protocol string /// Limit to grpc,http/s,dns,tls
  CreatedAt string
}

type ConnMan struct{
  Connections map[string]*Connector
	mu sync.RWMutex
}

func InitializeNewConnectionMan() *ConnMan {
	return &ConnMan{
		Connections: make(map[string]*Connector),
		mu: sync.RWMutex{},
	}
}

func (cm *ConnMan) GetAllConnectors()[]*Connector{
  cm.mu.Lock()
	defer cm.mu.Unlock()
  var cons []*Connector
  for _,con := range cm.Connections {
    cons = append(cons,con)
  }
  return cons
}


func (cm *ConnMan) NewConnector(msid,iaddr,oaddr,name string,driver *db.Driver)*Connector{
  cm.mu.Lock()
	defer cm.mu.Unlock()
  /*if _, ok := cm.Connections[msid]; !ok {
    utils.Logerror(fmt.Errorf("Connection with ID %s already exists",msid))
    return nil
  }*/
  for _,cn := range cm.Connections {
    if cn.SessionId == msid {
      utils.Logerror(fmt.Errorf("Connection with ID %s already exists",msid))
      return nil
    }
  }
  cnct := &Connector{
    SessionId: msid,// reverting to msid as the indicator for conn to it
    IAddress: iaddr,
    OAddress: oaddr,
    Name: name,
  }
  cm.Connections[cnct.SessionId] = cnct
  err := cm.SaveConnection(cnct,driver)
  if err != nil{
    utils.Warning("This is really not good, just save your connector manualy (save --con khjkhkjhlklj) or probably won't be able to connect to your C2")
    utils.Logerror(err)
  }
  return cnct
}

func (cm *ConnMan) DoesConnectionExist(id string) (string,bool){
  cm.mu.Lock()
  defer cm.mu.Unlock()
  for _,con := range cm.Connections {
    if con.SessionId == id {
      return con.OAddress,true
    }
  }
  return "",false
}

func (cm *ConnMan) SearchConnection(name string){
  cm.mu.Lock()
  defer cm.mu.Unlock()
  for _,cn := range cm.Connections {
    if strings.Contains(cn.Name,name){
      utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Name: %s",cn.Name))
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID: %s",cn.SessionId))
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Operator Address:   %s",cn.OAddress))
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Implant Address:   %s",cn.IAddress))
      utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
    }
  }
}

func (cm *ConnMan) AddConn(cnct *Connector) error{
  cm.mu.Lock()
  defer cm.mu.Unlock()
  /*if _, ok := cm.Connections[cnct.SessionId]; !ok {
    return fmt.Errorf(fmt.Sprintf("Connection with ID %s already exists",cnct.SessionId))
  }*/
  for _,cn := range cm.Connections {
    if cn.SessionId == cnct.SessionId {
      utils.Logerror(fmt.Errorf("Connection with ID %s already exists",cnct.SessionId))
      return nil
    }
  }
  cm.Connections[cnct.SessionId] = cnct
  return nil
}

func (cm *ConnMan) UpdateConnection(cnct *Connector,driver *db.Driver) error{
  cm.mu.Lock()
  //defer cm.mu.Unlock()
  prevId := cnct.SessionId
  if _, ok := cm.Connections[cnct.SessionId]; !ok {
    return fmt.Errorf(fmt.Sprintf("Connection with ID %s does not exists",cnct.SessionId))
  }
  for _,con := range cm.Connections {
    if con.SessionId == cnct.SessionId {
      cm.Connections[cnct.SessionId] = cnct
      cm.mu.Unlock()
      //remove connection from db
      err := cm.RemoveConnection(prevId,driver)
      if err != nil{
        return fmt.Errorf("Updated connection but not deleted previous values from db. %v",err)
      }
      return cm.SaveConnection(cnct,driver)//same as returning nil (bad progamming dude.)
    }
  }
  return nil
}

func (cm *ConnMan) RemoveConnection(cid string,driver *db.Driver) error {
  cm.mu.Lock()
  defer cm.mu.Unlock()
  if _,ok := cm.Connections[cid]; !ok{
    return fmt.Errorf(fmt.Sprintf("Error no connection with id %s",cid))
  }
  err := driver.Delete("connectors",cid)
  if err != nil{
    return fmt.Errorf("Error deleting connector from db. \nERROR: %v",err)
  }
  delete(cm.Connections,cid)
  return nil
}

func (cm *ConnMan) ListConnectors(){
  cm.mu.RLock()
	defer cm.mu.RUnlock()
  utils.PrintTextInASpecificColorInBold("yellow","    ***********    CURRENT CONNECTIONS    *********** ")
  for _,cn := range cm.Connections {
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Name: %s",cn.Name))
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID: %s",cn.SessionId))
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Operator Address:   %s",cn.OAddress))
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Implant Address:   %s",cn.IAddress))
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
  }
}

func (cm *ConnMan) GetConn(cid string) (*Connector,error){
  cm.mu.RLock()
  defer cm.mu.RUnlock()
  con,ok := cm.Connections[cid]
  if !ok {
    return nil,errors.New(fmt.Sprintf("Error getting connection with ID: %s. Does not Exists.",cid))
  }
  return con,nil
}

func (cm *ConnMan) SaveConnection(cnct *Connector,driver *db.Driver) error{
  return driver.Write("connectors",cnct.SessionId,cnct)
}

func (cm *ConnMan) LoadConnectors(driver *db.Driver)(error){
  connections,err := driver.ReadAll("connectors")
  if err != nil { return err }
  for _,cnct := range connections {
    var con Connector
    if err := json.Unmarshal([]byte(cnct),&con); err != nil{
      utils.Warning(fmt.Sprintf("%s",err))
      continue
    }
    if err := cm.AddConn(&con);err != nil{
      utils.Warning(fmt.Sprintf("%s",err))
      continue
    }
  }
	return nil
}

/*
func (cm *ConnMan) LoadConnectors(driver *db.Driver) error{
	dir := filepath.Join(driver.Dir,"connectors")
	if _,err := driver.ExportedStat(dir); err != nil{
    return errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n","sessions",err))
  }
	files,err := ioutil.ReadDir(dir)
  if err != nil{
    return fmt.Errorf("Error reading connections directory: %v",err)
  }
	for _,file := range files{
		var con Connector
		b,err := ioutil.ReadFile(filepath.Join(dir,file.Name()))
		if err != nil{
			utils.NoticeError(fmt.Sprintf("No session with such name: %s\nERROR: %s",file.Name(),err))
			continue
		}
		err = json.Unmarshal(b,&con)
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error unmarshalling to session: %s",err))
			continue
		}
		err = cm.AddConn(&con)
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error adding to session manager: %s",err))
			continue
		}
	}
  return nil
}
*/
