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
  createAssetStmt = "INSERT INTO `odin`.`assets` (name,asset_id,description,describers,active,hardware,created_at,updated_at) VALUES(?,?,?,?,?,?,?);"
  viewAssetStmt = `SELECT * FROM odin.assets WHERE (asset_id = ?);`
  listAssetStmt = "SELECT * FROM `odin`.`assets` WHERE (`active` = ? ) ORDER BY updated_at DESC;"
)

//  asset_id 	name 	description 	describers 	active 	hardware 	created_at 	updated_at
func (d *Doamin) CreateAsset(a dfn.Asset,ctx context.Context) error {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var ins *sql.Stmt
  ins,err := conn.PrepareContext(ctx,createAssetStmt)
  if err !=  nil{
    _  = utils.LogToFile(utils.Logger{Name:"assets_sql",Text:fmt.Sprintf("Error preparing to create asset: %s",err),})
    return errors.New("Server encountered an error while preparing to create asset. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,&a.Name,&a.AssetID,&a.Description,&a.Dscbr,&a.Active,&a.Hard,&a.CreatedAt,&a.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    _  = utils.LogToFile(utils.Logger{Name:"assets_sql",Text:fmt.Sprintf("Error executing create asset: %s",err),})
    return errors.New("Server encountered an error while creating asset.")
  }
  return nil
}

func (d *Doamin)  ListAssets(active bool,ctx context.Context)([]dfn.Asset,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  rows,err := conn.QueryContext(ctx,listAssetStmt,active)
  if err != nil{
    _  = utils.LogToFile(utils.Logger{Name:"assets_sql",Text:fmt.Sprintf("EQA: %s",err),})
    return nil,errors.New("Server encountered an error while listing assets.")
  }
  defer rows.Close()
  var asts []dfn.Asset
  for rows.Next(){
    var a dfn.Asset
    err = rows.Scan(&a.Name,&a.AssetID,&a.Description,&a.Dscbr,&a.Active,&a.Hard,&a.CreatedAt,&a.UpdatedAt)
    if err != nil{
      _  = utils.LogToFile(utils.Logger{Name:"assets_sql",Text:fmt.Sprintf("Error scanning for assets: %s",err),})
      return nil,errors.New("Server encountered an error while listing asset.")
    }
    asts = append(asts,a)
  }
  return asts,nil
}

func (d *Doamin)  ViewAsset(assetId string,ctx context.Context)(*dfn.Asset,error){
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var a dfn.Asset
  row := conn.QueryRowContext(ctx,viewAssetStmt,assetId)
  err := row.Scan(&a.Name,&a.AssetID,&a.Description,&a.Dscbr,&a.Active,&a.Hard,&a.CreatedAt,&a.UpdatedAt)
  if err != nil{
    _  = utils.LogToFile(utils.Logger{Name:"assets_danger_sql",Text:fmt.Sprintf("Error viewing asset with id %s.ERROR: %s",assetId,err),})
    if err == sql.ErrNoRows {
      utils.Danger(fmt.Errorf("Asset Id %s doen't exist: ", assetId))
      return nil,errors.New("Requested asset doesn't exist")
    }
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing asset of %s",assetId))
  }
  return &a,nil
}
