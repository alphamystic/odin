package domain

import (
  "fmt"
  "net"
  "errors"
  "context"
  "database/sql"
  "github.com/alphamystic/odin/lib/utils"
  //dfn"github.com/alphamystic/odin/lib/definers"
  png_hnd"github.com/alphamystic/odin/lib/handlers"
)

const (
  createTargetStmt = `INSERT INTO odin.targets (target_id,scan_id,host,host_ip,target_ip,firewall_name,decoys,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?);`
  listTargetsStmt = "SELECT * FROM odin.targets WHERE (scan_id = ?) ORDER BY updated_at ASC;"
  viewTargetStmt = `SELECT * FROM odin.targets WHERE target_id = ?;`
)

func (d *Domain) WriteTargetToDB(trgt *png_hnd.Target,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,createTargetStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error preparing to create user target: %s"),})
    return errors.New("Server encountered an error while preparing to create target. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&trgt.TargetID,&trgt.ScanID,&trgt.Host,&trgt.HostIp,&trgt.FireWallName,&trgt.Decoys,&trgt.CreatedAt,&trgt.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create target: %s",err),})
    return errors.New("Server encountered an error while creating target.")
  }
  return nil
}

func (d *Domain) ListTargets(scanId int,ctx context.Context) ([]png_hnd.Target,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listTargetsStmt,scanId)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error listing targets %s",err),})
    return nil,errors.New("Server encountered an error while listing targets.")
  }
  defer rows.Close()
  var tgs []png_hnd.Target
  for rows.Next(){
    var trgt png_hnd.Target
    var token string
    err = rows.Scan(&trgt.TargetID,&trgt.ScanID,&trgt.Host,&trgt.HostIp,&trgt.FireWallName,&token,&trgt.CreatedAt,&trgt.UpdatedAt)
    if err != nil {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error scanning for targets: %s",err),})
      return nil,errors.New("Error listing targets.")
    }
    decoys,_ := UnmarshalDecoys(token)
    trgt.Decoys = decoys
    tgs = append(tgs,trgt)
  }
  return tgs,nil
}


func (d *Domain) ViewTarget(targetId string,ctx context.Context)(*png_hnd.Target,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var t png_hnd.Target
  var token string
  row := conn.QueryRowContext(ctx,viewTargetStmt,targetId)
  err = row.Scan(&t.TargetID,&t.ScanID,&t.Host,&t.HostIp,&t.TargetIp,&t.FireWallName,&token,&t.CreatedAt,&t.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Hash for userid %s does not exist: ERROR: %",targetId,err),})
      return nil,errors.New(fmt.Sprintf("Target ID of %s is non existance.",targetId))
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing target id %s. ERROR: %s",targetId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing target with the id of %s",targetId))
  }
  decoys,_ := UnmarshalDecoys(token)
  t.Decoys = decoys
  return &t,nil
}



func UnmarshalDecoys(token string) ([]net.IP,error){
  var ips []string
  var decoys []net.IP
  ips,err := utils.TokenToArray(token)
  if err != nil {
    if errors.Is(err,utils.NotImplemented){
      return decoys,nil
    }
    return nil,err
  }
  for _,ip := range ips {
    decoys = append(decoys,net.ParseIP(ip))
  }
  return decoys,nil
}
