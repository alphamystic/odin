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

	"github.com/alphamystic/odin/lib/db"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const ConnectorsPath = "../.brain/"

type Connector struct{
  *dfn.Mothership
  SessionId string //coockie/password t connect to the MS/C2
  //If IAddress and OAddress are provided the grpc is theprotocol
  //If Tls then itis SGRPC if not thne regular GRPC
  //Add a custom error such that if one fails it defaults/assumes the other to be true
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


func (cm *ConnMan) NewConnector(ms *dfn.Mothership,driver *db.Driver) error {
  cm.mu.Lock()
	defer cm.mu.Unlock()
  /*if _, ok := cm.Connections[msid]; !ok {
    utils.Logerror(fmt.Errorf("Connection with ID %s already exists",msid))
    return nil
  }*/
  for _,cn := range cm.Connections {
    if cn.MSId == ms.MSId {
      utils.Logerror(fmt.Errorf("Connection with ID %s already exists",ms.MSId))
      return fmt.Errorf("Connection with ID %s already exists",ms.MSId)
    }
  }
  con := &Connector{
    Mothership: ms,
    SessionId:  "", // Provide a default value or some meaningful initialization
  }
  cm.Connections[ms.MSId] = con
  // Change this to save to main
  err := cm.SaveConnection(con,driver)
  if err != nil {
    utils.Warning("This is really not good, just save your connector manualy (save --con khjkhkjhlklj) or probably won't be able to connect to your C2")
    utils.Logerror(err)
    return err
  }
  return nil
}

func (cm *ConnMan) DoesConnectionExist(id string) (*Connector,bool){
  cm.mu.Lock()
  defer cm.mu.Unlock()
  for _,con := range cm.Connections {
    if con.MSId == id {
      return con,true
    }
  }
  return nil,false
}

func (cm *ConnMan) SearchConnection(name string){
  cm.mu.Lock()
  defer cm.mu.Unlock()
  for _,cn := range cm.Connections {
    if strings.Contains(cn.Name,name){
      utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Name: %s",cn.Name))
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID: %s",cn.MSId))
      if cn.OAddress != ""{
        utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Operator Address:   %s",cn.OAddress))
        utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Implant Address:   %s",cn.IAddress))
      } else {
        utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   HTTP/HTTPS Server: %s",cn.Address))
      }
      utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
    }
  }
}

func (cm *ConnMan) AddConn(cnct *Connector) error{
  cm.mu.Lock()
  defer cm.mu.Unlock()
  /*if _, ok := cm.Connections[cnct.MSId]; !ok {
    return fmt.Errorf(fmt.Sprintf("Connection with ID %s already exists",cnct.MSId))
  }*/
  for _,cn := range cm.Connections {
    if cn.MSId == cnct.MSId {
      utils.Logerror(fmt.Errorf("Connection with ID %s already exists",cnct.MSId))
      return fmt.Errorf("Connection with ID %s already exists",cnct.MSId)
    }
  }
  cm.Connections[cnct.MSId] = cnct
  return nil
}

// Redundant and should probably just change the session ID
func (cm *ConnMan) UpdateConnection(cnct *Connector,driver *db.Driver) error{
  cm.mu.Lock()
  //defer cm.mu.Unlock()
  prevId := cnct.MSId
  if _, ok := cm.Connections[cnct.MSId]; !ok {
    return fmt.Errorf(fmt.Sprintf("Connector with ID %s does not exists",cnct.MSId))
  }
  for _,con := range cm.Connections {
    if con.MSId == cnct.MSId {
      cm.Connections[cnct.MSId] = cnct
      cm.mu.Unlock()
      //remove connection from db
      err := cm.RemoveConnection(prevId,driver)
      if err != nil{
        return fmt.Errorf("Updated connection but not deleted previous values from local db. %v",err)
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
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID: %s",cn.MSId))
    if cn.OAddress != "" {
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Operator Address:   %s",cn.OAddress))
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Implant Address:   %s",cn.IAddress))
    } else {
      utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   HTTP/HTTPS Server: %s",cn.Address))
    }
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

// change this to be written to online DB/Server
func (cm *ConnMan) SaveConnection(cnct *Connector,driver *db.Driver) error{
  return driver.Write("connectors",cnct.MSId,cnct)
}

// Load this from the online server
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
