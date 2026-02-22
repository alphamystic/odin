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

// minion table columns:
// minionid, name, uname, userid, groupid, homedir, ostype, description,
// installed, mothershipid, address, motherships, tunnel_address, tls,
// ownerid, lastseen, is_dropper, generate_command, created_at, updated_at

const (
	createMinionStmt = `INSERT INTO odin.minion
	(minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	checkMinionRegStmt           = `SELECT minionid FROM odin.minion WHERE minionid = ?`
	listInstalledMinionsStmt     = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at FROM odin.minion WHERE installed = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listMyInstalledMinionsStmt   = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_atFROM odin.minion WHERE installed = ? AND ownerid = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listInstalledMinionFromMS    = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_atFROM odin.minion WHERE mothershipid = ? AND installed = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	markMinionAsInstalled        = `UPDATE odin.minion SET installed = ?, updated_at = ? WHERE minionid = ?;`
	listAllMinionStmt            = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_atFROM odin.minion ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listAllMinionByMS            = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_atFROM odin.minion WHERE mothershipid = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	viewMinionStmt               = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_atFROM odin.minion WHERE minionid = ?;`
)

// ====================== CREATE ======================

func (d *Domain) CreateMinion(ctx context.Context, m dfn.Minion) error {
	if d.CheckIfMinonIsRegistered(m.MinionID, ctx) {
		return errors.New("minion is already registered")
	}

	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, createMinionStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error preparing create: %s", err)})
		return errors.New("server error preparing to create minion")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		m.MinionID, m.Name, m.UName, m.UserID, m.GroupID, m.HomeDir,
		m.Os, m.Description, m.Installed, m.MothershipID, m.Address,
		m.Motherships, m.TunnelAddress, m.Tls, m.OwnerID,
		m.LastSeen, m.IsDropper, m.GenCommand, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error executing create: %s", err)})
		return errors.New("server encountered an error while creating minion")
	}
	return nil
}

// ====================== CHECK ======================

func (d *Domain) CheckIfMinonIsRegistered(minionId string, ctx context.Context) bool {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error db connection: %s", err))
		return true
	}
	defer conn.Close()

	var mid string
	row := conn.QueryRowContext(ctx, checkMinionRegStmt, minionId)
	if err := row.Scan(&mid); err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error checking registration: %s", err)})
		return true
	}
	return true
}

// ====================== LIST FUNCTIONS ======================

func (d *Domain) ListInstalledMinions(mine bool, ownerId string, limit, offset int, ctx context.Context) ([]dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	var rows *sql.Rows
	if mine {
		rows, err = conn.QueryContext(ctx, listMyInstalledMinionsStmt, true, ownerId, limit, offset)
	} else {
		rows, err = conn.QueryContext(ctx, listInstalledMinionsStmt, true, limit, offset)
	}
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error listing installed: %s", err)})
		return nil, errors.New("server error listing installed minions")
	}
	defer rows.Close()

	return scanMinions(rows, d)
}

func (d *Domain) ListAllMinions(limit, offset int, ctx context.Context) ([]dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listAllMinionStmt, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error listing all: %s", err)})
		return nil, errors.New("server error listing minions")
	}
	defer rows.Close()

	return scanMinions(rows, d)
}

func (d *Domain) ListAllMinionsFromASpecificMotherShip(msid string, limit, offset int, ctx context.Context) ([]dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listAllMinionByMS, msid, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error listing from mothership %s: %s", msid, err)})
		return nil, errors.New("server error listing minions by mothership")
	}
	defer rows.Close()

	return scanMinions(rows, d)
}

func (d *Domain) ListInstalledMinionsFromAMotherShip(msid string, limit, offset int, ctx context.Context) ([]dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listInstalledMinionFromMS, msid, true, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error listing installed from ms %s: %s", msid, err)})
		return nil, errors.New("server error listing installed minions by mothership")
	}
	defer rows.Close()

	return scanMinions(rows, d)
}

// ====================== VIEW ======================

func (d *Domain) ViewMinion(ctx context.Context, minionId string) (*dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	var m dfn.Minion
	row := conn.QueryRowContext(ctx, viewMinionStmt, minionId)
  var createdAtRaw, updatedAtRaw []byte
	err = row.Scan(
		&m.MinionID, &m.Name, &m.UName, &m.UserID, &m.GroupID, &m.HomeDir,
		&m.Os, &m.Description, &m.Installed, &m.MothershipID, &m.Address,
		&m.Motherships, &m.TunnelAddress, &m.Tls, &m.OwnerID, &m.LastSeen,
		&m.IsDropper, &m.GenCommand, &createdAtRaw, &updatedAtRaw)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error viewing minion %s: %s", minionId, err)})
		if err == sql.ErrNoRows {
			return nil, errors.New("minion not found")
		}
		return nil, errors.New("server error viewing minion")
	}
  _ = utils.ScanTimeStamps(&m.TimeStamps, createdAtRaw, updatedAtRaw)
	return &m, nil
}

// ====================== UPDATE ======================

func (d *Domain) MarkMinionAsInstalled(minionId string, ctx context.Context) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, markMinionAsInstalled)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error preparing update %s: %s", minionId, err)})
		return errors.New("server error marking minion as installed")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, true, utils.GetCurrentTime(), minionId)
	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error updating minion %s: %s", minionId, err)})
		return errors.New("failed to mark minion as installed")
	}
	return nil
}

// ====================== HELPERS ======================

func scanMinions(rows *sql.Rows, d *Domain) ([]dfn.Minion, error) {
	var minions []dfn.Minion
	for rows.Next() {
		var m dfn.Minion
    var createdAtRaw, updatedAtRaw []byte
		err := rows.Scan(
			&m.MinionID, &m.Name, &m.UName, &m.UserID, &m.GroupID, &m.HomeDir,
			&m.Os, &m.Description, &m.Installed, &m.MothershipID, &m.Address,
			&m.Motherships, &m.TunnelAddress, &m.Tls, &m.OwnerID, &m.LastSeen,
			&m.IsDropper, &m.GenCommand, &createdAtRaw, &updatedAtRaw)
		if err != nil {
			d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error scanning minion row: %s", err)})
			return nil, errors.New("error reading minion data")
		}
    _ = utils.ScanTimeStamps(&m.TimeStamps, createdAtRaw, updatedAtRaw)
		minions = append(minions, m)
	}
	return minions, nil
}
