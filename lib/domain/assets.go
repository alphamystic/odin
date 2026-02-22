package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/alphamystic/odin/lib/utils"
	dfn "github.com/alphamystic/odin/lib/definers"
)

const (
	createAssetStmt = `
	INSERT INTO odin.assets
	(asset_id, name, description, describers, active, hardware, owner_id, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	viewAssetStmt = `
	SELECT asset_id, name, description, describers, active, hardware, owner_id, created_at, updated_at
	FROM odin.assets WHERE asset_id = ?;
	`

	listAssetFilteredStmt = `
	SELECT asset_id, name, description, describers, active, hardware, owner_id, created_at, updated_at
	FROM odin.assets
	WHERE owner_id = ? AND active = ? AND hardware = ?
	ORDER BY updated_at DESC
	LIMIT ? OFFSET ?;
	`

	updateAssetStmt = `
	UPDATE odin.assets
	SET name = ?, description = ?, describers = ?, active = ?, hardware = ?, updated_at = ?
	WHERE asset_id = ?;
	`

  listAssetsByOwner = `
		SELECT asset_id, name, description, describers, active, hardware, owner_id, created_at, updated_at
		FROM odin.assets
		WHERE owner_id = ? AND active = ? AND hardware = ?
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?;
	`
)

// CreateAsset inserts a new asset into the DB.
func (d *Domain) CreateAsset(ctx context.Context, a dfn.Asset) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, createAssetStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error preparing asset insert: %s", err)})
		return errors.New("server error while preparing to create asset")
	}
	defer stmt.Close()

	// Base64 encode describer if it's a struct
	var describerJSON string
	switch v := a.Dscbr.(type) {
	case dfn.Describer:
		raw, _ := json.Marshal(v)
		describerJSON = utils.Base64Encode(string(raw))
	case string:
		describerJSON = v
	default:
		describerJSON = utils.Base64Encode(fmt.Sprintf("%v", v))
	}

	_, err = stmt.ExecContext(
		ctx,
		a.AssetID, a.Name, a.Description, describerJSON,
		a.Active, a.Hard, a.OwnerID, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error executing create asset: %s", err)})
		return errors.New("server encountered an error while creating asset")
	}

	return nil
}

// ViewAsset retrieves one asset by its asset_id.
func (d *Domain) ViewAsset(ctx context.Context, assetID string) (*dfn.Asset, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, viewAssetStmt, assetID)
	var a dfn.Asset
	var rawDescribers string
	var createdAtRaw, updatedAtRaw []byte

	err = row.Scan(&a.AssetID, &a.Name, &a.Description, &rawDescribers, &a.Active, &a.Hard, &a.OwnerID, &createdAtRaw, &updatedAtRaw)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset with id %s does not exist", assetID)
		}
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error viewing asset %s: %s", assetID, err)})
		return nil, errors.New("error retrieving asset details")
	}
  _ = utils.ScanTimeStamps(&a.TimeStamps, createdAtRaw, updatedAtRaw)

	// Decode describer
	decoded := utils.Base64Decode(rawDescribers)
	var desc dfn.Describer
	if json.Unmarshal([]byte(decoded), &desc) == nil {
		a.Dscbr = desc
	} else {
		a.Dscbr = decoded
	}
	return &a, nil
}

// ListAssetsByFilter lists assets by owner, active, and hardware filters with pagination.
func (d *Domain) ListAssetsByFilter(ctx context.Context, ownerID string, active, hardware bool, limit, offset int) ([]dfn.Asset, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listAssetFilteredStmt, ownerID, active, hardware, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error listing assets: %s", err)})
		return nil, errors.New("server encountered an error while listing assets")
	}
	defer rows.Close()

	var assets []dfn.Asset
	for rows.Next() {
		var a dfn.Asset
		var rawDescribers string
  	var createdAtRaw, updatedAtRaw []byte
		if err := rows.Scan(&a.AssetID, &a.Name, &a.Description, &rawDescribers, &a.Active, &a.Hard, &a.OwnerID, &createdAtRaw, &updatedAtRaw); err != nil {
			d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error scanning assets: %s", err)})
			continue
		}
    _ = utils.ScanTimeStamps(&a.TimeStamps, createdAtRaw, updatedAtRaw)

		decoded := utils.Base64Decode(rawDescribers)
		var desc dfn.Describer
		if json.Unmarshal([]byte(decoded), &desc) == nil {
			a.Dscbr = desc
		} else {
			a.Dscbr = decoded
		}

		assets = append(assets, a)
	}

	return assets, nil
}

// UpdateAsset modifies asset fields.
func (d *Domain) UpdateAsset(ctx context.Context, a dfn.Asset) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("error getting db connection: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, updateAssetStmt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error preparing update asset: %s", err)})
		return errors.New("server error while preparing to update asset")
	}
	defer stmt.Close()

	var describerJSON string
	switch v := a.Dscbr.(type) {
	case dfn.Describer:
		raw, _ := json.Marshal(v)
		describerJSON = utils.Base64Encode(string(raw))
	case string:
		describerJSON = v
	default:
		describerJSON = utils.Base64Encode(fmt.Sprintf("%v", v))
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.Description, describerJSON, a.Active, a.Hard, a.UpdatedAt, a.AssetID)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Error executing update asset: %s", err)})
		return errors.New("server error while updating asset")
	}

	rowsAff, _ := res.RowsAffected()
	if rowsAff == 0 {
		return fmt.Errorf("no asset found with id %s", a.AssetID)
	}

	return nil
}



func (d *Domain) ListAssetsByOwner(ownerID string, active, hardware bool, limit, offset int, ctx context.Context) ([]dfn.Asset, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error getting DB connection: %v", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listAssetsByOwner, ownerID, active, hardware, limit, offset)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("ListByOwner error: %s", err)})
		return nil, errors.New("Error listing minions")
	}
	defer rows.Close()

	var assets []dfn.Asset
	for rows.Next() {
		var a dfn.Asset
  	var createdAtRaw, updatedAtRaw []byte
		if err := rows.Scan(&a.AssetID, &a.Name, &a.Description, &a.Dscbr, &a.Active, &a.Hard, &a.OwnerID, &createdAtRaw, &updatedAtRaw); err != nil {
			d.LogToFile(utils.Logger{Name: "assets_sql", Text: fmt.Sprintf("Scan error: %s", err)})
			continue
		}
    _ = utils.ScanTimeStamps(&a.TimeStamps, createdAtRaw, updatedAtRaw)
		assets = append(assets, a)
	}

	return assets, nil
}
