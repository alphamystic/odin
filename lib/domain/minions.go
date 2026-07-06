package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"encoding/json"

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
                           (minionid, name, uname, userid, groupid, homedir, ostype, description, installed, mothershipid, address, motherships, tunnel_address, tls, active, describers, ownerid, lastseen, is_dropper, generate_command, created_at, updated_at)
                           VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	checkMinionRegStmt           = `SELECT minionid FROM odin.minion WHERE minionid = ?`
	listInstalledMinionsStmt     = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at FROM odin.minion WHERE installed = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listMyInstalledMinionsStmt   = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at FROM odin.minion WHERE installed = ? AND ownerid = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listInstalledMinionFromMS    = `SELECT minionid,name,uname,userid,groupid,homedir,ostype,description,installed,mothershipid,address,motherships,tunnel_address,tls,ownerid,lastseen,is_dropper,generate_command,created_at,updated_at FROM odin.minion WHERE mothershipid = ? AND installed = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	markMinionAsInstalled        = `UPDATE odin.minion SET installed = ?, updated_at = ? WHERE minionid = ?;`
	listAllMinionStmt = `SELECT minionid, name, uname, userid, groupid, homedir, ostype, description, installed, mothershipid, address, motherships, tunnel_address, tls, active, describers, ownerid, lastseen, is_dropper, generate_command, created_at, updated_at
        FROM odin.minion WHERE ownerid = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	listAllMinionByMS = `SELECT minionid, name, uname, userid, groupid, homedir, ostype, description, installed, mothershipid, address, motherships, tunnel_address, tls, active, describers, ownerid, lastseen, is_dropper, generate_command, created_at, updated_at
        FROM odin.minion WHERE mothershipid = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
	viewMinionStmt = `SELECT minionid, name, uname, userid, groupid, homedir, ostype, description, installed, mothershipid, address, motherships, tunnel_address, tls, active, describers, ownerid, lastseen, is_dropper, generate_command, created_at, updated_at
        FROM odin.minion WHERE minionid = ?;`
	updateMinionStmt = `UPDATE odin.minion
        SET name=?, uname=?, userid=?, groupid=?, homedir=?, ostype=?, description=?, installed=?, mothershipid=?, address=?, motherships=?, tunnel_address=?, tls=?, active=?, describers=?, ownerid=?, lastseen=?, is_dropper=?, generate_command=?, updated_at=?
        WHERE minionid = ?;`
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
	// Handle JSON data for describers
    descData, _ := json.Marshal(m.Describer)
    // You mentioned encoding motherships/describers in base64 in previous logic
    descEncoded := utils.Base64Encode(string(descData))

	stmt, err := conn.PrepareContext(ctx, createMinionStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error preparing create: %s", err)})
		return errors.New("server error preparing to create minion")
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx,
              m.MinionID,        // 1
              m.Name,            // 2
              m.UName,           // 3
              m.UserID,          // 4
              m.GroupID,         // 5
              m.HomeDir,         // 6
              m.Os,              // 7
              m.Description,     // 8
              m.Installed,       // 9
              m.MothershipID,    // 10
              m.Address,         // 11
              m.Motherships,     // 12 (The string field)
              m.TunnelAddress,   // 13
              m.Tls,             // 14
              m.Active,          // 15
              descEncoded,       // 16 (The JSON/Base64 field)
              m.OwnerID,         // 17
              m.LastSeen,        // 18
              m.IsDropper,       // 19
              m.GenCommand,      // 20
              m.CreatedAt,       // 21
              m.UpdatedAt,       // 22
           )
	if err != nil {
		d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error executing create: %s", err)})
		return errors.New("server encountered an error while creating minion")
	}
    rowsAff, err := res.RowsAffected()
	if rowsAff != 1 || err != nil {
	    d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error, more than 1 row affected: %s", err)})
		return errors.New("no mothership was created")
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

func (d *Domain) ListAllMinions(ctx context.Context,limit, offset int, owner_id string) ([]dfn.Minion, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listAllMinionStmt, owner_id, limit, offset)
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
    var rawDescribers string // Buffer for the JSON blob
    row := conn.QueryRowContext(ctx, viewMinionStmt, minionId)
    var createdAtRaw, updatedAtRaw []byte

    // Fixed: Added m.Active and rawDescribers to reach 22 arguments
    err = row.Scan(
       &m.MinionID, &m.Name, &m.UName, &m.UserID, &m.GroupID, &m.HomeDir,
       &m.Os, &m.Description, &m.Installed, &m.MothershipID, &m.Address,
       &m.Motherships, &m.TunnelAddress, &m.Tls, &m.Active, &rawDescribers,
       &m.OwnerID, &m.LastSeen, &m.IsDropper, &m.GenCommand, &createdAtRaw, &updatedAtRaw)

    if err != nil {
       d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error viewing minion %s: %s", minionId, err)})
       if err == sql.ErrNoRows {
          return nil, errors.New("minion not found")
       }
       return nil, errors.New("server error viewing minion")
    }

    // Decode JSON describers
    decoded := utils.Base64Decode(rawDescribers)
    json.Unmarshal([]byte(decoded), &m.Describer)

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
        var rawDescribers string
        var created, updated []byte

        // Ensure this list exactly matches your SELECT statements (22 columns)
        err := rows.Scan(
            &m.MinionID, &m.Name, &m.UName, &m.UserID, &m.GroupID, &m.HomeDir,
            &m.Os, &m.Description, &m.Installed, &m.MothershipID, &m.Address,
            &m.Motherships, &m.TunnelAddress, &m.Tls, &m.Active, &rawDescribers,
            &m.OwnerID, &m.LastSeen, &m.IsDropper, &m.GenCommand, &created, &updated)

        if err != nil {
            d.LogToFile(utils.Logger{Name: "minion_sql", Text: fmt.Sprintf("Error scanning minion row: %s", err)})
            return nil, errors.New("error reading minion data")
        }

        decoded := utils.Base64Decode(rawDescribers)
        json.Unmarshal([]byte(decoded), &m.Describer)

        _ = utils.ScanTimeStamps(&m.TimeStamps, created, updated)
        minions = append(minions, m)
    }
    return minions, nil
}

// ====================== UPDATE ======================
func (d *Domain) UpdateMinion(ctx context.Context, m dfn.Minion) error {
    if !d.CheckIfMinonIsRegistered(m.MinionID, ctx) {
        return errors.New("minion not found")
    }

    conn, err := d.GetConnection(ctx)
    if err != nil {
        return fmt.Errorf("error getting db connection: %w", err)
    }
    defer conn.Close()

    // Handle JSON data for describers
    descData, _ := json.Marshal(m.Describer)
    descEncoded := utils.Base64Encode(string(descData))

    stmt, err := conn.PrepareContext(ctx, updateMinionStmt)
    if err != nil {
        d.LogToFile(utils.Logger{
            Name: "minion_sql",
            Text: fmt.Sprintf("Error preparing update for %s: %s", m.MinionID, err),
        })
        return errors.New("server error preparing to update minion")
    }
    defer stmt.Close()

    // Execute the update with 20 parameters to match the SET clause + 1 for WHERE
    res, err := stmt.ExecContext(ctx,
        m.Name, m.UName, m.UserID, m.GroupID, m.HomeDir, m.Os,        // 1-6
        m.Description, m.Installed, m.MothershipID, m.Address,       // 7-10
        m.Motherships, m.TunnelAddress, m.Tls, m.Active,             // 11-14
        descEncoded, m.OwnerID, m.LastSeen, m.IsDropper,             // 15-18
        m.GenCommand, utils.GetCurrentTime(),                        // 19-20
        m.MinionID,                                                  // 21 (WHERE clause)
    )

    if err != nil {
        d.LogToFile(utils.Logger{
            Name: "minion_sql",
            Text: fmt.Sprintf("Error executing update for %s: %s", m.MinionID, err),
        })
        return errors.New("server encountered an error while updating minion")
    }

    rowsAffec, _ := res.RowsAffected()
    if rowsAffec == 0 {
        return errors.New("no changes were made to the minion")
    }

    return nil
}