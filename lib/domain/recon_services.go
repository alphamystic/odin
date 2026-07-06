package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	png_hnd "github.com/alphamystic/odin/lib/handlers"
	"github.com/alphamystic/odin/lib/utils"
)

// 	service_id 	target_id 	service_name 	port 	protocol 	state 	version 	created_at 	updated_at 	data
const (
	createServiceStmt = `INSERT INTO odin.services (target_id, service_name, port, protocol, state, version, created_at, updated_at, data) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?);`
	getServiceStmt  = `SELECT service_id, target_id, service_name, port, protocol, state, version, created_at, updated_at, data FROM odin.services WHERE service_id = ? AND target_id = ?;`
	listServicesStmt = `SELECT service_id, target_id, service_name, port, protocol, state, version, created_at, updated_at, data FROM odin.services WHERE target_id = ? ORDER BY updated_at DESC;`
)

func (d *Domain) CreateService(ctx context.Context, srvc png_hnd.Service) error {
	var err error
	// Get DB connection
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()
	// Convert Data to Token]
	data := utils.Base64Encode(srvc.Data)
	// Prepare SQL Statement
	ins, err := conn.PrepareContext(ctx, createServiceStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error preparing to create service: %s", err)})
		return errors.New("server encountered an error while preparing to create service. Try again later :).")
	}
	defer ins.Close()
	// Execute Query
	var tt = new(utils.TimeStamps)
  tt.Touch()
	res, err := ins.ExecContext(ctx, srvc.TargetID, srvc.ServiceName, srvc.Port, srvc.Protocol, srvc.State, srvc.Version, tt.CreatedAt, tt.UpdatedAt, data)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error executing create service: %s", err)})
		return errors.New("server encountered an error while creating service.")
	}
	// Check affected rows
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected != 1 {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: "No rows affected while creating service."})
		return errors.New("server encountered an error while creating service.")
	}

	return nil
}

func (d *Domain) GetService(ctx context.Context, targetId string, serviceId int) (*png_hnd.Service, error) {
	var srvc png_hnd.Service
	var token string

	// Get DB connection
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()
	// Query Database
	row := conn.QueryRowContext(ctx, getServiceStmt, serviceId, targetId)
	var createdAt, updatedAt []byte
	err = row.Scan(&srvc.ServiceID, &srvc.TargetID, &srvc.ServiceName, &srvc.Port, &srvc.Protocol, &srvc.State, &srvc.Version, &createdAt, &updatedAt, &token)
	if err != nil {
		if err == sql.ErrNoRows {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Service ID %d does not exist: ERROR: %v", serviceId, err)})
			return nil, fmt.Errorf("service ID %d does not exist", serviceId)
		}
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error viewing service %d. ERROR: %s", serviceId, err)})
		return nil, fmt.Errorf("server encountered an error while viewing service with ID %d", serviceId)
	}
	if err := utils.ScanTimeStamps(&srvc.TimeStamps, createdAt, updatedAt); err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
	}
	srvc.Data = utils.Base64Decode(token)
	return &srvc, nil
}

func (d *Domain) GetServices(ctx context.Context, targetId string) ([]png_hnd.Service, error) {
	var svcs []png_hnd.Service
	// Get DB connection
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()
	// Query Database
	rows, err := conn.QueryContext(ctx, listServicesStmt, targetId)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error listing services: %s", err)})
		return nil, errors.New("server encountered an error while listing all target services.")
	}
	defer rows.Close()
	// Iterate over rows
	for rows.Next() {
		var srvc png_hnd.Service
		var token string
    var createdAt, updatedAt []byte
		err = rows.Scan(&srvc.ServiceID, &srvc.TargetID, &srvc.ServiceName, &srvc.Port, &srvc.Protocol, &srvc.State, &srvc.Version, &createdAt, &updatedAt, &token)
		if err != nil {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Error scanning services: %s", err)})
			continue
		}

		// Use ScanTimeStamps for timestamp handling
		if err := utils.ScanTimeStamps(&srvc.TimeStamps, createdAt, updatedAt); err != nil {
			d.LogToFile(utils.Logger{Name: "recon_sql", Text: fmt.Sprintf("Timestamp parsing error: %s", err)})
      continue
		}

		if !utils.CheckifStringIsEmpty(token) {
			srvc.Data = ""
		} else{
			srvc.Data = utils.Base64Decode(token)
		}
		// Append to list
		svcs = append(svcs, srvc)
	}
	return svcs, nil
}
