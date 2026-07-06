package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

const (
	insertRmmTaskStmt      = `INSERT INTO odin.rmm_tasks (task_id, ownerid, command, target_type, minion_ids, created_at, updated_at) VALUES (?,?,?,?,?,?,?);`
	insertRmmExecutionStmt = `INSERT INTO odin.rmm_executions (execution_id, task_id, minionid, status, output, verification_hash, created_at, updated_at) VALUES (?,?,?,?,?,?,?,?);`
	updateRmmTaskStmt      = `UPDATE odin.rmm_tasks SET command = ?, target_type = ?, minion_ids = ?, updated_at = ? WHERE task_id = ? AND ownerid = ?;`
	viewRmmTaskStmt        = `SELECT task_id, ownerid, command, target_type, minion_ids, created_at, updated_at FROM odin.rmm_tasks WHERE task_id = ? AND ownerid = ?;`
	viewRmmExecutionsStmt  = `SELECT execution_id, task_id, minionid, status, output, verification_hash, created_at, updated_at FROM odin.rmm_executions WHERE task_id = ?;`
	listRmmTasksStmt       = `SELECT task_id, ownerid, command, target_type, minion_ids, created_at, updated_at FROM odin.rmm_tasks WHERE ownerid = ? ORDER BY created_at DESC;`
	minionListTasksStmt    = `SELECT task_id, ownerid, command, target_type, minion_ids, created_at, updated_at FROM odin.rmm_tasks WHERE target_type = ? OR target_type = 'custom' ORDER BY created_at DESC;`
)

// 1. CreateRmmTask - Encodes MinionIDs to Base64 before insertion
func (d *Domain) CreateRmmTask(ctx context.Context, t dfn.RmmTask) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
	    d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
	    return fmt.Errorf("db connection error: %w", err)
	   }
	defer conn.Close()

	idsJson, _ := json.Marshal(t.MinionIDs)
	encodedIDs := utils.Base64Encode(string(idsJson))

	res, err := conn.ExecContext(ctx, insertRmmTaskStmt, t.TaskID, t.OwnerID, t.Command, t.TargetTpe, encodedIDs, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error inserting RMM task %s: %v", t.TaskID, err)})
		return errors.New("failed to persist RMM task")
	}

	rowsAff, err := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("SQL Error: Expected 1 row affected, got %d for task %s", rowsAff, t.TaskID)})
		return errors.New("task creation failed: row mismatch")
	}
	return nil
}

// 2. UpdateRmmTask - Updates existing task details
func (d *Domain) UpdateRmmTask(ctx context.Context, t dfn.RmmTask) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
        d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
        return fmt.Errorf("db connection error: %w", err)
    }
	defer conn.Close()

	idsJson, _ := json.Marshal(t.MinionIDs)
	encodedIDs := utils.Base64Encode(string(idsJson))

	res, err := conn.ExecContext(ctx, updateRmmTaskStmt, t.Command, t.TargetTpe, encodedIDs, utils.GetCurrentTime(), t.TaskID, t.OwnerID)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error updating RMM task %s: %v", t.TaskID, err)})
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Update failed: No rows affected for task %s", t.TaskID)})
		return errors.New("no task updated, verify task_id and owner")
	}
	return nil
}

// 3. CreateRmmExecution - Called by Minion to log execution results
func (d *Domain) CreateRmmExecution(ctx context.Context, e dfn.RmmExecution) error {
	conn, err := d.GetConnection(ctx)
	if err != nil {
        d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
        return fmt.Errorf("db connection error: %w", err)
   }
	defer conn.Close()

	res, err := conn.ExecContext(ctx, insertRmmExecutionStmt, e.ExecutionID, e.TaskID, e.MinionID, e.Status, e.Output, e.VerificationHash, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error creating execution log for minion %s: %v", e.MinionID, err)})
		return errors.New("failed to log minion execution result")
	}

	rows, err := res.RowsAffected()
	if err != nil || rows != 1 {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("SQL Error: Failed to insert execution for minion %s", e.MinionID)})
		return errors.New("execution log failed: row mismatch")
	}
	return nil
}

// 4. ViewRMMTask - Retrieves task and its associated execution logs
func (d *Domain) ViewRMMTask(ctx context.Context, rmm_id, owner_id string) (*dfn.RmmTask, []dfn.RmmExecution, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
    	    d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
    	    return nil, nil, fmt.Errorf("db connection error: %w", err)
    }
	defer conn.Close()

	var t dfn.RmmTask
	var tc, tu []byte
	var encodedIDs string

	err = conn.QueryRowContext(ctx, viewRmmTaskStmt, rmm_id, owner_id).Scan(&t.TaskID, &t.OwnerID, &t.Command, &t.TargetTpe, &encodedIDs, &tc, &tu)
	if err != nil {
	    if errors.Is(err,sql.ErrNoRows) {
	        return nil, nil, dfn.RMMTaskNotFound
	    }
        d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error viewing (Scan) RMM task: %v", err)})
	    return nil,nil,err
	}

	_ = utils.ScanTimeStamps(&t.TimeStamps, tc, tu)
	json.Unmarshal([]byte(utils.Base64Decode(encodedIDs)), &t.MinionIDs)

	rows, err := conn.QueryContext(ctx, viewRmmExecutionsStmt, rmm_id)
	if err != nil {
	    d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error listing RMM Executions RMM task: %v", err)})
	    return nil, nil, err
	}
	defer rows.Close()

	var execs []dfn.RmmExecution
	for rows.Next() {
		var e dfn.RmmExecution
		var ec, eu []byte
		if err := rows.Scan(&e.ExecutionID, &e.TaskID, &e.MinionID, &e.Status, &e.Output, &e.VerificationHash, &ec, &eu); err == nil {
			_ = utils.ScanTimeStamps(&e.TimeStamps, ec, eu)
			execs = append(execs, e)
		}
	}
	return &t, execs, nil
}

// 5. ListRmmTasks - Lists all tasks for an admin
func (d *Domain) ListRmmTasks(ctx context.Context, owner_id string) ([]dfn.RmmTask, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
        d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
        return nil, fmt.Errorf("db connection error: %w", err)
    }
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, listRmmTasksStmt, owner_id)
	if err != nil { return nil, err }
	defer rows.Close()

	var tasks []dfn.RmmTask
	for rows.Next() {
		var t dfn.RmmTask
		var tc, tu []byte
		var encodedIDs string
		if err := rows.Scan(&t.TaskID, &t.OwnerID, &t.Command, &t.TargetTpe, &encodedIDs, &tc, &tu); err == nil {
			_ = utils.ScanTimeStamps(&t.TimeStamps, tc, tu)
			json.Unmarshal([]byte(utils.Base64Decode(encodedIDs)), &t.MinionIDs)
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

// 6. MinionListRMMTask - Filtering by OS type
func (d *Domain) MinionListRMMTask(ctx context.Context, ostype string) ([]dfn.RmmTask, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
    	d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
    	return nil, fmt.Errorf("db connection error: %w", err)
    }
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, minionListTasksStmt, ostype)
	if err != nil { return nil, err }
	defer rows.Close()

	var tasks []dfn.RmmTask
	for rows.Next() {
		var t dfn.RmmTask
		var tc, tu []byte
		var encodedIDs string
		if err := rows.Scan(&t.TaskID, &t.OwnerID, &t.Command, &t.TargetTpe, &encodedIDs, &tc, &tu); err == nil {
			_ = utils.ScanTimeStamps(&t.TimeStamps, tc, tu)
			json.Unmarshal([]byte(utils.Base64Decode(encodedIDs)), &t.MinionIDs)
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

// 7. ListRmmExecutions - Lists results for a specific task
func (d *Domain) ListRmmExecutuions(ctx context.Context, task_id, verified string) ([]dfn.RmmExecution, error) {
	conn, err := d.GetConnection(ctx)
	if err != nil {
        d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error getting connection for RMM task: %v", err)})
    	return nil, fmt.Errorf("db connection error: %w", err)
    }
	defer conn.Close()

	query := viewRmmExecutionsStmt
	var args []interface{}
	args = append(args, task_id)

	if verified != "" {
		query += " AND status = ?"
		args = append(args, verified)
	}

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		d.LogToFile(utils.Logger{Name: "rmm_sql", Text: fmt.Sprintf("Error listing executions for task %s: %v", task_id, err)})
		return nil, err
	}
	defer rows.Close()

	var execs []dfn.RmmExecution
	for rows.Next() {
		var e dfn.RmmExecution
		var ec, eu []byte
		if err := rows.Scan(&e.ExecutionID, &e.TaskID, &e.MinionID, &e.Status, &e.Output, &e.VerificationHash, &ec, &eu); err == nil {
			_ = utils.ScanTimeStamps(&e.TimeStamps, ec, eu)
			execs = append(execs, e)
		}
	}
	return execs, nil
}