package domain


import (
  "fmt"
  "errors"
  "context"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

// minionid 	name 	uname 	userid 	groupid 	homedir 	ostype 	description 	installed 	mothershipid 	address 	motherships
// tunnel_address 	tls 	ownerid 	lastseen 	is_dropper 	generate_command 	created_at 	updated_at
const (
  createMinionStmt = "INSERT INTO `odin`.`minion` (minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
  checkMinionRegStmt = "SELECT `minionid` FROM `odin`.`minion` WHERE minionid = ?"
  listInstalledMinionsStmt = "SELECT * FROM `odin`.`minion` WHERE (`installed` = ? ) ORDER BY updated_at ASC;"
  listMyInstalledMinionsStmt = "SELECT * FROM `odin`.`minion` WHERE (`installed` = ? AND ownerid = ?) ORDER BY updated_at ASC;"
  listIsntalledMinionFromMS = "SELECT * FROM `odin`.`minion` WHERE (`msid` = ? AND `installed` = ? ) ORDER BY updated_at ASC;"
  markMinionAsInstalled = "UPDATE `odin`.`minion` SET `installed` = ? AND `updated_at` = ? WHERE (`miniond` = ?);";
  listAllMinionStmt = "SELECT * FROM `odin`.`minion` ORDER BY updated_at ASC;"
  listAllMinionByMS = "SELECT * FROM `odin`.`minion` WHERE (`msid` = ? ) ORDER BY updated_at ASC;"
  viewMinionStmt = "SELECT * FROM `odin`.`minion` WHERE minionid	 = ?;"
)

func (d *Domain) CreateMinion(m dfn.Minion,ctx context.Context) error {
  if d.CheckIfMinonIsRegistered(m.MinionID,ctx){
    return errors.New("Minon is already registered")
  }
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,createMinionStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error preparing to create minion: %s",err),})
    return errors.New("Server encountered an error while creating minion, Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error executing create minion: %s",err),})
    return errors.New("Server encountered an error while creating minion.")
  }
  return nil
}

//a transaction checking if a minion is registered more than twice then deletes one of the record
func (d *Domain) SanitizeMinions()error{
  //get the mid of all minions
  //range through all if any exists more than once
  //delete that record
  return nil
}

func (d *Domain) CheckIfMinonIsRegistered(minionId string,ctx context.Context) bool {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error db connection: %s",err))
    return true
  }
  defer conn.Close()
  var mid string
  row := conn.QueryRowContext(ctx,checkMinionRegStmt,minionId)
  err = row.Scan(&mid)
  if err != nil{
    if err == sql.ErrNoRows{
      return false
    }
    utils.NoticeError(fmt.Sprintf("%s",err))
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error checking if minion is registered: %s",err),})
    // there's an error so we default to true that way It should be redone
    return true
  }
  return true
}

func (d *Domain) ListInstalledMinions(mine bool,ownerId string,ctx context.Context) ([]dfn.Minion,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var rows *sql.Rows
  if mine {
    rows,err = conn.QueryContext(ctx,listMyInstalledMinionsStmt,true,ownerId)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing installed minions: %s",err),})
      return nil,errors.New("Server encountered an error while listing all installed minions.")
    }
  } else {
    rows,err = conn.QueryContext(ctx,listInstalledMinionsStmt,true)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing installed minions: %s",err),})
      return nil,errors.New("Server encountered an error while listing all installed minions.")
    }
  }
  defer rows.Close()
  var minions []dfn.Minion
  for rows.Next(){
    var m dfn.Minion
    err = rows.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error scanning for installed minions: %s",err),})
      return nil,errors.New("Server encountered an error while listing all installed minions.")
    }
    minions = append(minions,m)
  }
  return minions,nil
}

func (d *Domain) ListAllMinionsFromASpecificMotherShip(msid string,ctx context.Context)([]dfn.Minion,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listAllMinionByMS,msid)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing minions from a specific mothership: %s",err),})
    return nil,errors.New("Server encountered an error while listing all minions from specified mother ship.")
  }
  defer rows.Close()
  var minions []dfn.Minion
  for rows.Next(){
    var m dfn.Minion
    err = rows.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error scanning minions from a specific mothership id %S.ERROR: %s",msid,err),})
      return nil,errors.New("Server encountered an error while listing all minions from specified mothership.")
    }
    minions = append(minions,m)
  }
  return minions,nil
}

func (d *Domain) ListInstalledMinionsFromAMotherShip(msid string,ctx context.Context)([]dfn.Minion,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listIsntalledMinionFromMS,msid,true)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing for installed minions from a ms %s. ERROR: %s",msid,err),})
    return nil,errors.New("Server encountered an error while listing all my minions from the specified mother ship.")
  }
  defer rows.Close()
  var minions []dfn.Minion
  for rows.Next(){
    var m dfn.Minion
    err = rows.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error scanning for installed minions from a ms %s. ERROR: %s",msid,err),})
      return nil,errors.New("Server encountered an error while listing all minions from the specified mothership.")
    }
    minions = append(minions,m)
  }
  return minions,nil
}

func (d *Domain) MarkMinionAsInstalled(minionId string,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  stmt,err := conn.PrepareContext(ctx,markMinionAsInstalled)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error marking minion id %s as installed: %s",minionId,err),})
    return errors.New("Server encountered an error while marking minion as installed.")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,true,utils.GetCurrentTime(),minionId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error marking minion %s as installed. ERROR: %s",minionId,err),})
    return errors.New("Server encountered an error while making marking minion as installed.")
  }
  return nil
}

func (d *Domain) ListAllMinions(ctx context.Context)([]dfn.Minion,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listAllMinionStmt)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing all minions: %s",err),})
    return nil,errors.New("Server encountered an error while listing all my minions.")
  }
  defer rows.Close()
  var minions []dfn.Minion
  for rows.Next(){
    var m dfn.Minion
    err = rows.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error listing all minions: %s",err),})
      return nil,errors.New("Server encountered an error while listing all my minions.")
    }
    minions = append(minions,m)
  }
  return minions,nil
}

func (d *Domain) ViewMinion(minionId string,ctx context.Context)(*dfn.Minion,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var m dfn.Minion
  row := conn.QueryRowContext(ctx,viewMinionStmt,minionId)
  err = row.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error viewing minon %s %s",minionId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing minon with id of %s",minionId))
  }
  return &m,nil
}

/*
func (d *Domain) CheckIfMinonIsRegistered(minionId string,ctx context.Context)(*dfn.Minion,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var m dfn.Minion
  row := conn.QueryRowContext(ctx,viewMinionStmt,minionId)
  err := row.Scan(&m.MinionID,&m.Name,&m.UName,&m.UserID,&m.GroupID,&m.HomeDir,&m.Os,&m.Description,&m.Installed,utils.GetCurrentTime(),&m.Address,&m.Motherships,&m.TunnelAddress,&m.Tls,&m.OwnerID,&m.LastSeen,&m.IsDropper,&m.GenCommand,&m.CreatedAt,&m.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      return nil,dfn.ErrNoMinion
    }
    d.LogToFile(utils.Logger{Name:"minion_sql",Text:fmt.Sprintf("Error viewing minon %s %s",minionId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing minon with id of %s",minionId))
  }
  return &m,nil
}
*/
