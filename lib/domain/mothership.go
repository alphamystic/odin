package domain

import (
  "fmt"
  "errors"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const (
  createMothershipStmt = "INSERT INTO `odin`.`motherships` (ownerid,name,password,msid,address,implant_tunnel,admin_tunnel,other_motherships,description,tls,certpem,keypem,active,generate_command,machine_data,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
  viewMotherShipStmt = "SELECT * FROM `odin`.`motherships` WHERE `msid`	 = ?;"
  deactivateMothershipStmt = "UPDATE `odin`.`motherships` SET (`active` = ? `updated_at` = ?) WHERE (`msid` = ? AND `ownerid` = ?);";
  listMothershipStmt = "SELECT * FROM `odin`.`motherships` WHERE (`active` = ? ) ORDER BY updated_at DESC;"
)

// ownerid 	name 	password 	msid 	address 	implant_tunnel 	admin_tunnel 	other_motherships 	description 	tls 	certpem 	keypem 	active 	generate_command 	machine_data 	created_at 	updated_at
func (d *Domain) CreateMothership(m dfn.Mothership,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err := conn.PrepareContext(ctx,createMothershipStmt)
  if err !=  nil{
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error preparing to create mothership: %s",err),})
    return errors.New("Server encountered an error while preparing to create mothership. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&m.OwnerID,&m.Name,&m.Password,&m.MSId,&m.Address,&m.ImplantTunnel,&m.AdminTunnel,&m.Motherships,&m.Description,&m.Tls,&m.CertPem,&m.KeyPem,&m.Active,&m.GenCommand,&m.MaichineData,&m.CreatedAt,&m.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error excuting create mothership: %s",err),})
    return fmt.Errorf("Server encountered an error while creating mothership. %v",err)
  }
  return nil
}

func (d *Domain) ListMotherShips(active bool,ctx context.Context) ([]dfn.Mothership,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listMothershipStmt,active)
  if err != nil{
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error listing mothership: %s",err),})
    return nil,errors.New("Server encountered an error while listing MotherShips.")
  }
  defer rows.Close()
  var msps []dfn.Mothership
  for rows.Next(){
    var m dfn.Mothership
    err = rows.Scan(&m.OwnerID,&m.Name,&m.Password,&m.MSId,&m.Address,&m.ImplantTunnel,&m.AdminTunnel,&m.Motherships,&m.Description,&m.Tls,&m.CertPem,&m.KeyPem,&m.Active,&m.GenCommand,&m.MaichineData,&m.CreatedAt,&m.UpdatedAt)
    if err != nil{
      _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error scanning for mothership: %s",err),})
      return nil,errors.New("Server encountered an error while listing mothership.")
    }
    msps = append(msps,m)
  }
  return msps,nil
}

func (d *Domain) DeactivateMothership(msid,ownerId string,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  stmt,err := conn.PrepareContext(ctx,deactivateMothershipStmt)
  if err != nil{
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error preparing to deactivate mothership: %s",err),})
    return errors.New("Server encountered an error while preparing to deactivate mothership.")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,false,currentTime,msid,ownerId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error executing deactivate mothership: %s",err),})
    return errors.New("Server encountered an error while executing deactivate mothership.")
  }
  return nil
}

func (d *Domain) ViewMothership(msid string,ctx context.Context)(*dfn.Mothership,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var m dfn.Mothership
  row := conn.QueryRowContext(ctx,viewMotherShipStmt,msid)
  err := row.Scan(&m.OwnerID,&m.Name,&m.Password,&m.MSId,&m.Address,&m.ImplantTunnel,&m.AdminTunnel,&m.Motherships,&m.Description,&m.Tls,&m.CertPem,&m.KeyPem,&m.Active,&m.GenCommand,&m.MaichineData,&m.CreatedAt,&m.UpdatedAt)
  if err != nil{
    _  = utils.LogToFile(utils.Logger{Name:"ms_sql",Text:fmt.Sprintf("Error viewing mothership: %s",err),})
    if err == sql.ErrNoRows {
      utils.Danger(fmt.Errorf("MS doen't exist: %s", msid))
      return nil,errors.New("Requested msid doesn't exist")
    }
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing MS of %s",msid))
  }
  return &m,nil
}
