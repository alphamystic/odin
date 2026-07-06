package domain

import (
  "fmt"
  "errors"
  "context"
  //"database/sql"
  "encoding/json"
  _ "github.com/go-sql-driver/mysql"

  "github.com/alphamystic/odin/lib/utils"
  dfn"github.com/alphamystic/odin/lib/definers"
)

const (
    insertApiKey = `INSERT INTO odin.apikey (key_id, tool_name, ownerid, credential_type, encrypted_data, active, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?);`
    listApiKeys  = `SELECT key_id, tool_name, ownerid, credential_type, encrypted_data, active, created_at, updated_at FROM odin.apikey WHERE ownerid = ? AND active = ? ORDER BY updated_at DESC LIMIT ? OFFSET ?;`
    viewApikey   = `SELECT key_id, tool_name, ownerid, credential_type, encrypted_data, active, created_at, updated_at FROM odin.apikey WHERE key_id = ?;`
    updateApiKey = `UPDATE odin.apikey SET tool_name = ?, credential_type = ?, encrypted_data = ?, active = ?, updated_at = ? WHERE key_id = ?;`
    checkApiKey    = `SELECT key_id, ownerid FROM odin.apikey WHERE encrypted_data = ? AND ownerid = ? AND active = 1;`
)

//	key_id	tool_name	ownerid	credential_type	encrypted_data	active	created_at	updated_at
// CreateApiKey handles encryption and verbose error reporting for new credentials
func (d *Domain) CreateApiKey(ctx context.Context, a dfn.Api) error {
    conn, err := d.GetConnection(ctx)
    if err != nil { return fmt.Errorf("error getting db connection: %w", err) }
    defer conn.Close()

    // Serialize and Encrypt
    jsonData, _ := json.Marshal(a.Data)
    encrypted, err := d.UE.EncryptSecret(string(jsonData))
    if err != nil { return fmt.Errorf("encryption error: %w", err) }

    res, err := conn.ExecContext(ctx, insertApiKey, a.KeyID, a.ToolName, a.OwnerID, a.CredentialType, encrypted, a.Active, a.CreatedAt, a.UpdatedAt)
    if err != nil {
        d.LogToFile(utils.Logger{Name: "apikey_sql", Text: fmt.Sprintf("Error executing create API Key: %v", err)})
        return errors.New("server encountered an error while creating API Key")
    }

    rowsAff, err := res.RowsAffected()
    if err != nil || rowsAff != 1 {
        d.LogToFile(utils.Logger{Name: "apikey_sql", Text: fmt.Sprintf("Unexpected rows affected: %d", rowsAff)})
        return errors.New("API Key creation failed - persistence error")
    }
    return nil
}

// ListApiKeys retrieves and decrypts paginated results
func (d *Domain) ListApiKeys(ctx context.Context, ownerID string, active bool, limit, offset int) ([]dfn.Api, error) {
    conn, err := d.GetConnection(ctx)
    if err != nil { return nil, err }
    defer conn.Close()

    rows, err := conn.QueryContext(ctx, listApiKeys, ownerID, active, limit, offset)
    if err != nil {
        d.LogToFile(utils.Logger{Name: "apikey_sql", Text: fmt.Sprintf("Error listing keys: %v", err)})
        return nil, errors.New("failed to retrieve API Key list")
    }
    defer rows.Close()

    var keys []dfn.Api
    for rows.Next() {
        var a dfn.Api
        var encrypted string
        var created, updated []byte

        err := rows.Scan(&a.KeyID, &a.ToolName, &a.OwnerID, &a.CredentialType, &encrypted, &a.Active, &created, &updated)
        if err != nil {
            utils.NoticeError(fmt.Sprintf("%q",err))
            continue
           }
        _ = utils.ScanTimeStamps(&a.TimeStamps, created, updated)

        // Decrypt for application use
        decrypted, _ := d.UE.DecryptSecret(encrypted)
        a.Data.KeyID = a.KeyID
        _ = json.Unmarshal([]byte(decrypted), &a.Data)
        keys = append(keys, a)
    }
    return keys, nil
}

// UpdateKey modifies the record and re-encrypts updated data
func (d *Domain) UpdateKey(ctx context.Context, a dfn.Api) error {
    conn, err := d.GetConnection(ctx)
    if err != nil { return err }
    defer conn.Close()

    jsonData, _ := json.Marshal(a.Data)
    encrypted, _ := d.UE.EncryptSecret(string(jsonData))

    res, err := conn.ExecContext(ctx, updateApiKey, a.ToolName, a.CredentialType, encrypted, a.Active, a.UpdatedAt, a.KeyID)
    if err != nil {
        d.LogToFile(utils.Logger{Name: "apikey_sql", Text: fmt.Sprintf("Update failed for %s: %v", a.KeyID, err)})
        return errors.New("failed to update credential")
    }

    rows, _ := res.RowsAffected()
    if rows == 0 { return errors.New("no API Key found with that ID") }
    return nil
}

func (d *Domain) ViewApiKey(ctx context.Context, keyID string) (*dfn.Api, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil { return nil, err }
	defer conn.Close()

	var a dfn.Api
	var encrypted string
	var created, updated []byte
	row := conn.QueryRowContext(ctx, viewApikey, keyID)
	err = row.Scan(&a.KeyID, &a.ToolName, &a.OwnerID, &a.CredentialType, &encrypted, &a.Active, &created, &updated)
	if err != nil { return nil, err }

	_ = utils.ScanTimeStamps(&a.TimeStamps, created, updated)
	decrypted, _ := d.UE.DecryptSecret(encrypted)
	_ = json.Unmarshal([]byte(decrypted), &a.Data)
	return &a, nil
}
func (d *Domain) CheckIfApiKey(ctx context.Context, apiKey, ownerID string) bool {
	conn, err := d.GetConnection(ctx)
	if err != nil { return false }
	defer conn.Close()

	// Since data is encrypted, we encrypt the input to check for a match
	encryptedInput, _ := d.UE.EncryptSecret(apiKey)
	var kid, oid string
	err = conn.QueryRowContext(ctx, checkApiKey, encryptedInput, ownerID).Scan(&kid, &oid)
	return err == nil
}