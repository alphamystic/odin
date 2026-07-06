package domain

import (
  "fmt"
  "net"
  "errors"
  "context"
  "database/sql"
  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
  png_hnd"github.com/alphamystic/odin/lib/handlers"
)

//target_id 	scan_id 	host 	host_ip 	target_ip 	firewall_name 	decoys 	created_at 	updated_at
const (
  createTargetStmt = `INSERT INTO odin.targets (target_id,scan_id,host,host_ip,target_ip,firewall_name,decoys,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?);`
  listTargetsStmt = "SELECT * FROM odin.targets WHERE (scan_id = ?) ORDER BY updated_at ASC;"
  //viewTargetStmt = `SELECT * FROM odin.targets WHERE target_id = ?;`
  viewTargetStmt = `SELECT target_id, scan_id, host, host_ip, target_ip, firewall_name, decoys, created_at, updated_at FROM odin.targets WHERE target_id = ?;`
)

func (d *Domain) WriteTargetToDB(trgt *png_hnd.Target, ctx context.Context) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("Error getting db connection: %q", err)
	}
	defer conn.Close()
	ins, err := conn.PrepareContext(ctx, createTargetStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error preparing to create target: %s", err)})
		return errors.New("Server encountered an error while preparing to create target. Try again later :).")
	}
	defer ins.Close()
	decoys := utils.IPArrayToString(trgt.Decoys)
  hostIP := utils.IPToStorageFormat(trgt.HostIp)
	targetIP := utils.IPToStorageFormat(trgt.TargetIp)
	res, err := ins.ExecContext(ctx, &trgt.TargetID, &trgt.ScanID, &trgt.Host, hostIP, targetIP, &trgt.FireWallName, decoys, &trgt.CreatedAt, &trgt.UpdatedAt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error executing create target: %s", err)})
		return errors.New("Server encountered an error while creating target.")
	}
	rowsAffec, err := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error executing create target, no rows affected: %s", err)})
		return errors.New("Server encountered an error while creating target.")
	}
	return nil
}


func (d *Domain) ListTargets(ctx context.Context, scanId string) ([]png_hnd.Target, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting db connection: %q", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listTargetsStmt, scanId)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error listing targets %s", err)})
		return nil, errors.New("Server encountered an error while listing targets.")
	}
	defer rows.Close()

	var tgs []png_hnd.Target
	for rows.Next() {
		var trgt png_hnd.Target
		var hostIP, targetIP, decoysStr string
    var createdAt, updatedAt []byte
		err = rows.Scan(&trgt.TargetID, &trgt.ScanID, &trgt.Host, &hostIP, &targetIP, &trgt.FireWallName, &decoysStr, &createdAt, &updatedAt)
		if err != nil {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error scanning for targets: %s", err)})
			return nil, errors.New("Error listing targets.")
		}
		// Convert string back to net.IP and []net.IP
		trgt.HostIp = utils.StringToIP(hostIP)
		trgt.TargetIp = utils.StringToIP(targetIP)
		trgt.Decoys = utils.StringToIPList(decoysStr)
    if err := utils.ScanTimeStamps(&trgt.TimeStamps, createdAt, updatedAt); err != nil {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
      continue
		}
		tgs = append(tgs, trgt)
	}
	return tgs, nil
}



// ad
func (d *Domain) ViewTarget(scanId,targetId,userId string,ctx context.Context)(*png_hnd.Target,error) {
  if !d.DoesScanExists(scanId,userId,ctx) {
    return nil,dfn.ScanDoesNotExists
  }
  conn,err := d.GetConnection(ctx)
  if err != nil {
    return nil,fmt.Errorf("Error getting db connection: %q",err)
  }
  defer conn.Close()
  var t png_hnd.Target
  var token string
  row := conn.QueryRowContext(ctx,viewTargetStmt,targetId)
  var createdAt, updatedAt []byte
  var hostIp, targetIp string
  err = row.Scan(&t.TargetID,&t.ScanID,&t.Host,&hostIp,&targetIp,&t.FireWallName,&token,&createdAt, &updatedAt)
  if err != nil{
    if err == sql.ErrNoRows {
      d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Target %s does not exist: ERROR: %",targetId,err),})
      return nil,dfn.TargetDoesNotExist
    }
    d.LogToFile(utils.Logger{Name:"recon_sql",Text:fmt.Sprintf("Error viewing target id %s. ERROR: %s",targetId,err),})
    return nil,errors.New(fmt.Sprintf("Server encountered an error while viewing target with the id of %s",targetId))
  }
  t.HostIp = utils.StorageFormatToIP(hostIp)
  t.TargetIp = utils.StorageFormatToIP(targetIp)
  if err := utils.ScanTimeStamps(&t.TimeStamps, createdAt, updatedAt); err != nil {
    d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
    t.Touch()
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
