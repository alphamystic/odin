package domain

import (
  "fmt"
  "sync"
  "errors"
  "context"
  "database/sql"
  "github.com/alphamystic/odin/lib/utils"
  png_hnd"github.com/alphamystic/odin/lib/handlers"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const (
  createScanStmt =`INSERT INTO odin.scans (scanid,name,scan_type,created_at,updated_at) VALUES(?,?,?,?,?);`
  listScansStmt = "SELECT * FROM `odin`.`scans` ORDER BY updated_at DESC;"
  listScanByTypeStmt = "SELECT * FROM `odin`.`scans` WHERE (`scan_type` = ? ) ORDER BY updated_at DESC;"
  viewScanStmt = `SELECT * FROM odin.scans   WHERE (scan_id = ?);`
  getWebDataStmt =  "SELECT directory_path,parameter_path,file_path FROM `odin`.`scans` WHERE (`target_id` = ?);"
)

//  	scan_id 	name 	scan_type 	created_at 	updated_at
func (d *Domain) CreateScan(name,scan_type /*domIp*/ string,ctx context.Context) error {
  var ins *sql.Stmt
  var tt utils.TimeStamps
  tt.Touch()
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  ins,err = conn.PrepareContext(ctx,createScanStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error preparing to create scan: %s",err),})
    return errors.New("Server encountered an error while preparing to create scan. Try again later :).")
  }
  defer ins.Close()
  res,err := ins.ExecContext(ctx,utils.Md5Hash(utils.GenerateUUID()),scan_type,name,&tt.CreatedAt,&tt.UpdatedAt)
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create user hash: %s",err),})
    return errors.New("Server encountered an error while creating hash.")
  }
  return nil
}

func (d *Domain) ListScan(ctx context.Context,scan_type string) ([]png_hnd.Scans,error) {
  var rows *sql.Rows
  var err error
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  if scan_type == "ALL"{
    rows,err = conn.QueryContext(ctx,listScansStmt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error listing all scans: %s",err),})
      return nil,errors.New("Server encountered an error while listing all scans.")
    }
  } else {
    rows,err = conn.QueryContext(ctx,listScanByTypeStmt,scan_type)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("ELAU: %s",err),})
      return nil,errors.New("Server encountered an error while listing all specified users.")
    }
  }
  defer rows.Close()
  var scans []png_hnd.Scans
  for rows.Next(){
    var scan png_hnd.Scans
    err = rows.Scan(&scan.ScanID,&scan.Name,&scan.ScanType,&scan.CreatedAt,&scan.UpdatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error scanning list scans: %s",err),})
      continue
    }
    scans = append(scans,scan)
  }
  return scans,nil
}


func (d *Domain) DoesScanExists(scanId,userId string,ctx context.Context) bool {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    utils.Warning(fmt.Sprintf("Error getting db connection: %q",err))
    return true
  }
  defer conn.Close()
  var scan png_hnd.Scans
  row := conn.QueryRowContext(ctx,getHash,userId)
  err = row.Scan(&scan.ScanID,&scan.Name,&scan.ScanType,&scan.CreatedAt,&scan.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      return false
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing scan %s. ERROR: %s",scanId,err),})
    utils.Warning(fmt.Sprintf("Server encountered an error while viewing scan with id of %s",scanId))
    return false
  }
  return true
}

// should return a scan and any data availble on it(recon data,vulnerabilities and exploits)
func (d *Domain) ViewScan(scanId,userId string,ctx context.Context) (*png_hnd.Scans,error) {
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var scan png_hnd.Scans
  row := conn.QueryRowContext(ctx,getHash,userId)
  err = row.Scan(&scan.ScanID,&scan.Name,&scan.ScanType,&scan.CreatedAt,&scan.UpdatedAt)
  if err != nil {
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Scan %s does not exist: ERROR: %",scanId,err),})
      return nil,errors.New(fmt.Sprintf("Scan id %s is non existant.",scanId))
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing scan %s. ERROR: %s",scanId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing scan with id of %s",scanId))
  }
  return &scan,nil
}


// target_id 	directory_path 	parameter_path 	file_path 	created_at 	updated_at
func (d *Domain)  GetWebData(ctx context.Context,targetID string) (*png_hnd.WebData,error) {
  var dirUnsn,parUnsn,filesUnsn string
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  row := conn.QueryRowContext(ctx,getWebDataStmt,targetID)
  err = row.Scan(dirUnsn,parUnsn,filesUnsn)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Target %s does not exist: ERROR: %",targetID,err),})
      return nil,errors.New(fmt.Sprintf("Target id %s is non existance.",targetID))
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing webdata for target %s. ERROR: %s",targetID,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing target with id of %s",targetID))
  }
  var dir,par,fls []string
  var wg sync.WaitGroup
  go func(){
    dir,err = utils.TokenToArray(dirUnsn)
  }()
  go func(){
    par,err = utils.TokenToArray(parUnsn)
  }()
  go func(){
    fls,err = utils.TokenToArray(filesUnsn)
  }()
  wg.Wait()
  return &png_hnd.WebData{
    Directories: dir,
    Parameters:  par,
    Files:       fls,
  }, err
}


// deprecated with the new DB
func (d *Domain) GetWebData1(targetID int, ctx context.Context) (*png_hnd.WebData, error) {
  dirCtx, cancelDir := context.WithCancel(ctx)
  parCtx, cancelPar := context.WithCancel(ctx)
  flsCtx, cancelFls := context.WithCancel(ctx)
  defer cancelDir()
  defer cancelPar()
  defer cancelFls()
  var wg sync.WaitGroup
  var dir, par, fls []string
  var err error
  wg.Add(3)
  go func() {
    defer wg.Done()
    dir, err = d.GetSpecificWebData(1, 1, dirCtx)
  }()
  go func() {
    defer wg.Done()
    par, err = d.GetSpecificWebData(1, 2, parCtx)
  }()
  go func() {
    defer wg.Done()
    fls, err = d.GetSpecificWebData(1, 3, flsCtx)
  }()
  wg.Wait()
  if err != nil {
    return nil, fmt.Errorf("Error getting webdata: %q",err)
  }
  return &png_hnd.WebData{
    Directories: dir,
    Parameters:  par,
    Files:       fls,
  }, nil
}


// deprecated
// target_id 	webdata 	directory_path 	parameter_path 	file_path 	created_at 	updated_at
func (d *Domain) GetSpecificWebData(targetId,wdType int,ctx context.Context) ([]string,error) {
  var data []string
  var stmt string
  switch wdType {
  case 1:
    stmt = "statement for getting directories."
  case 2:
    stmt = "statement for parameters."
  case 3:
    stmt = "statement for files"
  default:
    return data,dfn.Undefined
  }
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var h dfn.UserHash
  row := conn.QueryRowContext(ctx,stmt,targetId)
  err = row.Scan(&h.UserID,&h.Hash,&h.CreatedAt,&h.UpdatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Web Data %s does not exist: ERROR: %",targetId,err),})
      return nil,errors.New(fmt.Sprintf("Target ID of %s is non existance.",targetId))
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing Web Data  %s. ERROR: %s",targetId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing web data with id of %s",targetId))
  }
  return data,nil
}
