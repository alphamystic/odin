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
  createScanStmt = `INSERT INTO odin.scans (scan_id,name,scan_type,owner_id,created_at,updated_at) VALUES(?,?,?,?,?,?);`
  createWebDataStmt = `INSERT INTO odin.webdata (target_id, 	directory_path, 	parameter_path, 	file_path, created_at, updated_at) VALUES(?,?,?,?,?,?);`
  listScansStmt = "SELECT `scan_id`,`name`,`scan_type`,`owner_id`,`created_at`,`updated_at` FROM `odin`.`scans`WHERE `owner_id` = ? ORDER BY updated_at DESC;"
  listScanByTypeStmt = "SELECT `scan_id`,`name`,`scan_type`,`owner_id`,`created_at`,`updated_at`  FROM `odin`.`scans` WHERE `scan_type` = ? AND `owner_id` = ? ORDER BY updated_at DESC;"
  viewScanStmt = `SELECT * FROM odin.scans   WHERE (scan_id = ?);`
  getScanStmt = `SELECT * FROM odin.scans   WHERE (scan_id = ? AND owner_id = ?);`
  getWebDataStmt =  "SELECT directory_path,parameter_path,file_path FROM `odin`.`webdata` WHERE (`target_id` = ?);"
)

// scan_id 	name 	scan_type 	owner_id 	created_at 	updated_at
//target_id 	directory_path 	parameter_path 	file_path 	created_at 	updated_at
func (d *Domain) CreateWebdata(ctx context.Context,target_id, 	directory_path, 	parameter_path, 	file_path string) error {
  var ins *sql.Stmt
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  ins,err = conn.PrepareContext(ctx,createWebDataStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error preparing to create webdata: %s",err),})
    return errors.New("Server encountered an error while preparing to create webdata. Try again later :).")
  }
  defer ins.Close()
  var tt = new(utils.TimeStamps)
  tt.Touch()
  res,err := ins.ExecContext(ctx,target_id, 	directory_path, 	parameter_path, 	file_path,&tt.CreatedAt,&tt.UpdatedAt)
  if err != nil {
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create webdata: %s",err),})
    return fmt.Errorf("Error executing create webdata: %w",err)
  }
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create webdata: %s",err),})
    return errors.New("Server encountered an error while creating webdata.")
  }
  return nil
}

// Add owner id to the scan
//  	scan_id 	name 	scan_type 	created_at 	updated_at
func (d *Domain) CreateScan(ctx context.Context, name, scan_type, owner_id string) (string,error) {
  var ins *sql.Stmt
  var tt utils.TimeStamps
  tt.Touch()
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return "",fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  ins,err = conn.PrepareContext(ctx,createScanStmt)
  if err !=  nil{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error preparing to create scan: %s",err),})
    return "",errors.New("Server encountered an error while preparing to create scan. Try again later :).")
  }
  defer ins.Close()
  scanid := utils.Md5Hash(utils.GenerateUUID())
  res,err := ins.ExecContext(ctx,scanid,name,scan_type,owner_id,&tt.CreatedAt,&tt.UpdatedAt)
  if err != nil {
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create scan: %s",err),})
    return "",fmt.Errorf("Error executing create scan: %w",err)
  }
  rowsAffec, _  := res.RowsAffected()
  if err != nil || rowsAffec != 1{
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error executing create user hash: %s",err),})
    return "",errors.New("Server encountered an error while creating user.")
  }
  return scanid,nil
}

// Add this t factor in the owner id
func (d *Domain) ListScan(ctx context.Context,scan_type,owner_id string) ([]png_hnd.Scans,error) {
  var rows *sql.Rows
  var err error
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  if scan_type == "ALL"{
    rows,err = conn.QueryContext(ctx,listScansStmt,owner_id)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error listing all scans: %s",err),})
      return nil,errors.New("Server encountered an error while listing all scans.")
    }
  } else {
    rows,err = conn.QueryContext(ctx,listScanByTypeStmt,scan_type,owner_id)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("ELAU: %s",err),})
      return nil,errors.New("Server encountered an error while listing all specified users.")
    }
  }
  defer rows.Close()
  var scans []png_hnd.Scans
  for rows.Next(){
    var scan png_hnd.Scans
    var createdAt, updatedAt []byte
    err = rows.Scan(&scan.ScanID,&scan.Name,&scan.ScanType,&scan.OwnerID,&createdAt, &updatedAt)
    if err != nil{
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error scanning list scans: %s",err),})
      continue
    }
    if err := utils.ScanTimeStamps(&scan.TimeStamps, createdAt, updatedAt); err != nil {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
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
  row := conn.QueryRowContext(ctx,getScanStmt,scanId,userId)
  var createdAt, updatedAt []byte
  err = row.Scan(&scan.ScanID,&scan.Name,&scan.ScanType,&scan.OwnerID,&createdAt, &updatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      return false
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing (getting) scan %s for user %s. ERROR: %s",scanId,userId,err),})
    utils.Warning(fmt.Sprintf("Server encountered an error while viewing scan with id of %s for user %s",scanId,userId))
    return false
  }
  // if err := utils.ScanTimeStamps(&scan.TimeStamps, createdAt, updatedAt); err != nil {
  //   d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
  // }
  return true
}

// should return a scan and any data availble on it(recon data,vulnerabilities and exploits)
func (d *Domain) GetWebData(ctx context.Context, targetID string) (*png_hnd.WebData, int, error) {
  var dirUnsn, parUnsn, filesUnsn string
  conn, err := d.GetConnection(ctx)
  if err != nil {
    return nil, 0, fmt.Errorf("error getting DB connection: %q", err)
  }
  defer conn.Close()
  row := conn.QueryRowContext(ctx, getWebDataStmt, targetID)
  err = row.Scan(&dirUnsn, &parUnsn, &filesUnsn)
  if err != nil {
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name: "recon_sql",Text: fmt.Sprintf("Web Data for Target %s does not exist. ERROR: %v", targetID, err),})
      utils.Notice(fmt.Sprintf("Web Data for  target ID %s does not exist", targetID))
      return nil, 0, dfn.WebDataForTargetDoesNotExist
    }
    d.LogToFile(utils.Logger{Name: "recon_sql",Text: fmt.Sprintf("Error retrieving web data for target %s. ERROR: %v", targetID, err),})
    return nil, 0,fmt.Errorf("server encountered an error retrieving target with ID %s", targetID)
  }
  var (
    dir, par, fls []string
    count = 0
    wg            sync.WaitGroup
    errCh         = make(chan error, 3) // Buffer for three potential errors
  )
  // Goroutine for converting directories
  wg.Add(1)
  go func() {
    defer wg.Done()
    var convertErr error
    dir, convertErr = utils.TokenToArray(dirUnsn)
    if convertErr != nil {
      count ++
      errCh <- fmt.Errorf("error converting directories: %w", convertErr)
      d.LogToFile(utils.Logger{Name: "recon_conversion", Text: fmt.Sprintf("%s",convertErr)})
    }
  }()
  // Goroutine for converting parameters
  wg.Add(1)
  go func() {
    defer wg.Done()
    var convertErr error
    par, convertErr = utils.TokenToArray(parUnsn)
    if convertErr != nil {
      count ++
      errCh <- fmt.Errorf("error converting parameters: %w", convertErr)
      d.LogToFile(utils.Logger{Name: "recon_conversion", Text: fmt.Sprintf("%s",convertErr)})
    }
  }()
  // Goroutine for converting files
  wg.Add(1)
  go func() {
    defer wg.Done()
    var convertErr error
    fls, convertErr = utils.TokenToArray(filesUnsn)
    if convertErr != nil {
      count ++
      errCh <- fmt.Errorf("error converting files: %w", convertErr)
      d.LogToFile(utils.Logger{Name: "recon_conversion", Text: fmt.Sprintf("%s",convertErr)})
    }
  }()
  wg.Wait()
  close(errCh) // Close the error channel after all goroutines complete
  // Check for errors from the channel
  for convertErr := range errCh {
    if convertErr != nil {
      utils.Danger(convertErr)
      //return nil, convertErr
    }
  }
  // Return the parsed data
  return &png_hnd.WebData{
    Directories: dir,
    Parameters:  par,
    Files:       fls,
  }, count, nil
}


/*
// deprecated with the new DB
func (d *Domain) GetWebData1(targetID int, ctx context.Context) (*png_hnd.WebData, error) {
  // whats the effect of using a single ctx here
  dirCtx, cancelDir := context.WithCancel(ctx)
  parCtx, cancelPar := context.WithCancel(ctx)
  flsCtx, cancelFls := context.WithCancel(ctx)
  defer cancelDir()
  defer cancelPar()
  defer cancelFls()
  var wg sync.WaitGroup
  var dir, par, fls []string
  var err error
  count := 0
  go func() {
    wg.Add(1)
    defer wg.Done()
    dir, err1 := d.GetSpecificWebData(1, 1, dirCtx)
    if err1 != nil{
      count += 1
      //write error cahnnel
    }
  }()
  go func() {
    wg.Add(1)
    defer wg.Done()
    par, err2 := d.GetSpecificWebData(1, 2, parCtx)
    if err2 != nil {
      count += 1
      // write to error channel
    }
  }()
  go func() {
    wg.Add(1)
    defer wg.Done()
    fls, err3 := d.GetSpecificWebData(1, 3, flsCtx)
    if err3 != nil{
      count +=1
      //write to error channel
    }
  }()
  //create an error channel to log and handl the errors say via a count
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
*/
