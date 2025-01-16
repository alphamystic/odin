package domain

import (
  "fmt"
  "errors"
  "context"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/alphamystic/odin/lib/utils"
  png_hnd"github.com/alphamystic/odin/lib/handlers"
)

const (
  createServiceStmt = `INSERT INTO odin.services (service_id,target_id,service_name,port,protocol,state,version,AT,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?);`
  getServiceStmt = `SELECT * FROM odin.services WHERE service_id = ? AND target_id, = ?;`
  listServicesStmt = `SELECT * FROM odin.services WHERE (target_id = ? ) ORDER BY updated_at DESC;`
)

// service_id 	target_id 	service_name 	port 	protocol 	state 	version 	AT 	created_at 	updated_at data
// this should be a statement
func (d *Domain) CreateService(ctx context.Context,targetId int,data string,srvc png_hnd.Service) error {
  //var data string
  var err error
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  data,err = utils.MultipleToToken(&srvc.Data)
  if err != nil {
    return fmt.Errorf("Error creating token output: %q",err)
  }
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,createServiceStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error preparing to create service: %s"),})
    return errors.New("Server encountered an error while preparing to create service. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&srvc.ServiceID,&srvc.TargetID,&srvc.ServiceName,&srvc.Port,&srvc.Protocol,&srvc.State,&srvc.Version,&srvc.CreatedAt,&srvc.UpdatedAt,&data)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create service: %s",err),})
    return errors.New("Server encountered an error while creating service.")
  }
  return nil
}

func (d *Domain) GetService(ctx context.Context,targetId,serviceId int) (*png_hnd.Service,error){
  var srvc *png_hnd.Service
  var token string
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  type output []png_hnd.Output
  row := conn.QueryRowContext(ctx,getServiceStmt,serviceId,targetId)
  err = row.Scan(&srvc.ServiceID,&srvc.TargetID,&srvc.ServiceName,&srvc.Port,&srvc.Protocol,&srvc.State,&srvc.Version,&srvc.CreatedAt,&srvc.UpdatedAt,&token)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Hash for userid %s does not exist: ERROR: %",serviceId,err),})
      return nil,errors.New(fmt.Sprintf("Service id of %s is non existance.",serviceId))
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing service %s. ERROR: %s",serviceId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing service with id of %s",serviceId))
  }
  data,_ := utils.TokenToData(token)
  srvc.Data = data.(output)
  return srvc,nil
}

func (d *Domain) GetServices(ctx context.Context,targetId string) ([]*png_hnd.Service,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var svcs []*png_hnd.Service
  type output []png_hnd.Output
  rows,err := conn.QueryContext(ctx,listServicesStmt,targetId)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error listing api keys: %s",err),})
    return nil,errors.New("Server encountered an error while listing all targets services.")
  }
  defer rows.Close()
  for rows.Next() {
    var srvc png_hnd.Service
    var token string
    err = rows.Scan(&srvc.ServiceID,&srvc.TargetID,&srvc.ServiceName,&srvc.Port,&srvc.Protocol,&srvc.State,&srvc.Version,&srvc.CreatedAt,&srvc.UpdatedAt,&token)
    if err != nil {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error scanning for targets services: %s",err),})
      return nil,errors.New("Error listing targets services.")
    }
    data,_ := utils.TokenToData(token)
    srvc.Data = data.(output)
    svcs = append(svcs,&srvc)
  }
  return svcs,nil
}
