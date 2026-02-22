package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

// SQL statements
const (
	createMothershipStmt = `
	INSERT INTO odin.motherships
	(ownerid, name, password, msid, address, implant_tunnel, admin_tunnel, other_motherships, description, tls, certpem, keypem, active, generate_command, machine_data, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	viewMotherShipStmt = `
	SELECT ownerid, name, password, msid, address, implant_tunnel, admin_tunnel, other_motherships, description, tls, certpem, keypem, active, generate_command, machine_data, created_at, updated_at
	FROM odin.motherships WHERE msid = ?;`

	deactivateMothershipStmt = `
	UPDATE odin.motherships SET active = ?, updated_at = ?
	WHERE msid = ? AND ownerid = ?;`

	updateMothershipStmt = `
	UPDATE odin.motherships
	SET name = ?, address = ?, description = ?, tls = ?, certpem = ?, keypem = ?, updated_at = ?
	WHERE msid = ? AND ownerid = ?;`

	listMothershipStmt = `
	SELECT ownerid, name, password, msid, address, implant_tunnel, admin_tunnel, other_motherships, description, tls, certpem, keypem, active, generate_command, machine_data, created_at, updated_at
	FROM odin.motherships
	WHERE ownerid = ? AND active = ?
	ORDER BY updated_at DESC
	LIMIT ? OFFSET ?;`

	listFilteredMothershipStmt = `
	SELECT ownerid, name, password, msid, address, implant_tunnel, admin_tunnel, other_motherships, description, tls, certpem, keypem, active, generate_command, machine_data, created_at, updated_at
	FROM odin.motherships
	WHERE ownerid = ? AND active = ? AND online = ?
	ORDER BY updated_at DESC
	LIMIT ? OFFSET ?;`
)

// CreateMothership creates a new mothership record
func (d *Domain) CreateMothership(m dfn.Mothership, ctx context.Context) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, createMothershipStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error preparing create mothership: %v", err)})
		return errors.New("server error preparing to create mothership")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, m.OwnerID, m.Name, m.Password, m.MSId, m.Address,
		m.ImplantTunnel, m.AdminTunnel, m.Motherships, m.Description, m.Tls,
		m.CertPem, m.KeyPem, m.Active, m.GenCommand, m.Machinedata, m.CreatedAt, m.UpdatedAt)

	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error executing create mothership: %v", err)})
		return fmt.Errorf("error creating mothership: %w", err)
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff != 1 {
		return errors.New("no mothership was created")
	}

	return nil
}

// UpdateMothership updates an existing mothership record
func (d *Domain) UpdateMothership(ownerID, msid string, update dfn.Mothership, ctx context.Context) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, updateMothershipStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error preparing update mothership: %v", err)})
		return errors.New("server error preparing to update mothership")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, update.Name, update.Address, update.Description, update.Tls,
		update.CertPem, update.KeyPem, utils.GetCurrentTime(), msid, ownerID)

	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error executing update mothership: %v", err)})
		return fmt.Errorf("error executing update: %w", err)
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		return fmt.Errorf("no mothership updated, check msid or ownerID")
	}

	return nil
}

// DeactivateMothership sets mothership as inactive
func (d *Domain) DeactivateMothership(msid, ownerId string, ctx context.Context) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, deactivateMothershipStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error preparing deactivate mothership: %v", err)})
		return errors.New("server error preparing to deactivate mothership")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, false, utils.GetCurrentTime(), msid, ownerId)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error executing deactivate mothership: %v", err)})
		return fmt.Errorf("error executing deactivate: %w", err)
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		return fmt.Errorf("no mothership deactivated for msid %s", msid)
	}

	return nil
}

// ListMotherships returns a paginated list of motherships for an owner
func (d *Domain) ListMotherships(ownerID string, active bool, limit, offset int, ctx context.Context) ([]dfn.Mothership, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listMothershipStmt, ownerID, active, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error listing motherships: %v", err)})
		return nil, errors.New("server error listing motherships")
	}
	defer rows.Close()

	var motherships []dfn.Mothership
	for rows.Next() {
		var m dfn.Mothership
	  var createdAtRaw, updatedAtRaw []byte
		if err := rows.Scan(&m.OwnerID, &m.Name, &m.Password, &m.MSId, &m.Address, &m.ImplantTunnel,
			&m.AdminTunnel, &m.Motherships, &m.Description, &m.Tls, &m.CertPem, &m.KeyPem,
			&m.Active, &m.GenCommand, &m.Machinedata, createdAtRaw, updatedAtRaw); err != nil {
			d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error scanning mothership: %v", err)})
			return nil, fmt.Errorf("error scanning mothership: %w", err)
		}
    _ = utils.ScanTimeStamps(&m.TimeStamps, createdAtRaw, updatedAtRaw)
		motherships = append(motherships, m)
	}
	return motherships, nil
}

// ListFilteredMotherships filters by online/active flags for a specific owner
func (d *Domain) ListFilteredMotherships(ownerID string, active, online bool, limit, offset int, ctx context.Context) ([]dfn.Mothership, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listFilteredMothershipStmt, ownerID, active, online, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error filtering motherships: %v", err)})
		return nil, errors.New("server error filtering motherships")
	}
	defer rows.Close()

	var motherships []dfn.Mothership
	for rows.Next() {
		var m dfn.Mothership
    var createdAtRaw, updatedAtRaw []byte
		if err := rows.Scan(&m.OwnerID, &m.Name, &m.Password, &m.MSId, &m.Address, &m.ImplantTunnel,
			&m.AdminTunnel, &m.Motherships, &m.Description, &m.Tls, &m.CertPem, &m.KeyPem,
			&m.Active, &m.GenCommand, &m.Machinedata, &createdAtRaw, &updatedAtRaw); err != nil {
			d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error scanning filtered mothership: %v", err)})
			return nil, fmt.Errorf("error scanning filtered mothership: %w", err)
		}
    _ = utils.ScanTimeStamps(&m.TimeStamps, createdAtRaw, updatedAtRaw)
		motherships = append(motherships, m)
	}
	return motherships, nil
}

// ViewMothership returns details of a specific mothership
func (d *Domain) ViewMothership(msid string, ctx context.Context) (*dfn.Mothership, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	var m dfn.Mothership
	row := conn.QueryRowContext(ctx, viewMotherShipStmt, msid)
  var createdAtRaw, updatedAtRaw []byte
	if err := row.Scan(&m.OwnerID, &m.Name, &m.Password, &m.MSId, &m.Address,
		&m.ImplantTunnel, &m.AdminTunnel, &m.Motherships, &m.Description, &m.Tls,
		&m.CertPem, &m.KeyPem, &m.Active, &m.GenCommand, &m.Machinedata,
		 &createdAtRaw, &updatedAtRaw); err != nil {

		d.LogToFile(utils.Logger{Name: "ms_sql", Text: fmt.Sprintf("Error viewing mothership: %v", err)})
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mothership not found: %s", msid)
		}
		return nil, fmt.Errorf("error retrieving mothership %s: %w", msid, err)
	}
  _ = utils.ScanTimeStamps(&m.TimeStamps, createdAtRaw, updatedAtRaw)
	return &m, nil
}
