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

const (
  insertApiKey = `INSERT INTO odin.apikey (apikey,ownerid,active,created_at,updated_at) VALUES(?,?,?,?,?);`
  listApiKeys = `SELECT * FROM odin.apikey WHERE (active = ? ) ORDER BY updated_at DESC;`
  updateApiKey = `UPDATE odin.apikey SET (apikey = ? AND updated_at = ?) WHERE (ownerid = ?);`
  deactivateKey = `UPDATE odin.apikey SET (active = ? updated_at = ?) WHERE (apikey = ? AND ownerid = ?);`
  viewApikey = `SELECT * FROM odin.apikey WHERE apikey	 = ?;`
  checkIfApiKey = `SELECT apikey,ownerid FROM odin.apikey WHERE (apikey= ? AND ownerid = ?);`
)

//  	apikey 	ownerid 	active 	created_at 	updated_at
func (d *Domain) CreateApiKey(a dfn.Api,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err = conn.PrepareContext(ctx,insertApiKey)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error preparing to create api key statement: %w",err),})
    return errors.New("Server encountered an error while preparing to create apikey. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&a.ApiKey,&a.OwnerID,&a.Active,&a.CreatedAt,&a.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("more than one row affected: %w",err),})
    return fmt.Errorf("Server encountered an error while creating API Key. %v",err)
  }
  return nil
}


func (d *Domain) ListApiKeys(active bool,ctx context.Context)([]dfn.Api,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listApiKeys,active)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error listing api keys: %w",err),})
    return nil,errors.New("Server encountered an error while listing all api keys.")
  }
  defer rows.Close()
  var keys []dfn.Api
  for rows.Next(){
    var a dfn.Api
    err = rows.Scan(&a.ApiKey,&a.OwnerID,&a.Active,&a.CreatedAt,&a.UpdatedAt)
    if err != nil {
      d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error scanning for api keys: %w",err),})
      return nil,errors.New("Error listing api keys.")
    }
    keys = append(keys,a)
  }
  return keys,nil
}

/*func (d *Domain) ValidateApiKey(suppliedKey,ownerId string)bool{
  var key,user string
  stmt := "SELECT apikey,ownerid FROM `odin`.`apikey` WHERE (`apikey`= ? AND `ownerid` = ?);"
  row := conn.QueryRow(suppliedKey,ownerId)
  err := row.Scan(&key,&user)
  if err != nil{
    d.LogToFile("sql",fmt.Sprintf("Error scanning apikey rows %s",err))
    utils.LogError(e)
    return false
  }
  return true
}*/

func (d *Domain) UpdateKey(ownerId,apiKey string,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  stmt,err := conn.PrepareContext(ctx,updateApiKey)
  if err != nil{
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error updating api key: %w",err),})
    return errors.New("Server encountered an error while preparing to update API Key.")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,apiKey,utils.GetCurrentTime(),ownerId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error executing update key: %w",err),})
    return errors.New("Server encountered an error while executing update apikey.")
  }
  return nil
}

func (d *Domain) DeactivateKey(ownerId,apiKey string,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  stmt,err := conn.PrepareContext(ctx,deactivateKey)
  if err != nil {
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error preparing to update apikey: %w",err),})
    return errors.New("Server encountered an error while preparing to deactivate API Key.")
  }
  defer stmt.Close()
  var res sql.Result
  res,err = stmt.ExecContext(ctx,false,utils.GetCurrentTime(),apiKey,ownerId)
  rowsAffec,_ := res.RowsAffected()
  if err != nil || rowsAffec != 1 {
    _ = utils.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error deactivating api key, more than one row affected: %s",err),})
    return errors.New("Server encountered an error while executing deactivate apikey.")
  }
  return nil
}


func (d *Domain) ViewApiKey(keyId string,ctx context.Context)(*dfn.Api,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var a dfn.Api
  row := conn.QueryRowContext(ctx,viewApikey,keyId)
  err = row.Scan(&a.ApiKey,&a.OwnerID,&a.Active,&a.CreatedAt,&a.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"apikey_danger_sql",Text:fmt.Sprintf("Api key doen't exist: %s", keyId),})
      return nil,errors.New("Requested Apikey doesn't exist")
    }
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error viwing apikey %s .ERROR:%w",keyId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing apikey of %s",keyId))
  }
  return &a,nil
}


func (d *Domain) CheckIfApiKey(apiKey,ownerId string,ctx context.Context) bool {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error getting db connection: %w",err))
    return false
  }
  defer conn.Close()
  var key,user string
  row := conn.QueryRowContext(ctx,checkIfApiKey,apiKey,ownerId)
  err = row.Scan(&key,&user)
  if err != nil {
    d.LogToFile(utils.Logger{Name:"apikey_sql",Text:fmt.Sprintf("Error scanning apikey rows %w",err)})
    return false
  }
  if key == apiKey && ownerId == user {
    return true
  }
  return false
}
