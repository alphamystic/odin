package apihandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

// CreateRmmTask - Admin dispatches a new command to target minions
func (api_hnd *APIHandler) CreateRmmTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can dispatch RMM tasks.")
		return
	}

	var task dfn.RmmTask
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	ud := req.Context().Value("userData").(*UserData)
	if ud == nil {
		api_hnd.Unauthorized(res, "User data not found")
		return
	}
	task.OwnerID = ud.UserId

	// Initialize Metadata
	task.TaskID = utils.GenerateUUID()
	task.Touch()

	ctx := context.Background()
	if err := api_hnd.Dom.CreateRmmTask(ctx, task); err != nil {
		utils.Warning(fmt.Sprintf("Error creating RMM task: %v", err))
		api_hnd.InternalServerError(res, "Failed to create RMM task.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "RMM task dispatched successfully.",
		"task_id": task.TaskID,
	})
    return
}

// UpdateRmmTask - Admin updates an existing task (e.g., modifying command or targets)
func (api_hnd *APIHandler) UpdateRmmTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	if !api_hnd.IsAdmin(req) {
		api_hnd.Unauthorized(res, "Only admins can update RMM tasks.")
		return
	}

	var task dfn.RmmTask
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !utils.CheckifStringIsEmpty(task.TaskID) {
		api_hnd.BadRequest(res, "TaskID is required.")
		return
	}

	ud := req.Context().Value("userData").(*UserData)
	task.OwnerID = ud.UserId

	ctx := context.Background()
	if err := api_hnd.Dom.UpdateRmmTask(ctx, task); err != nil {
		utils.Warning(fmt.Sprintf("Error updating RMM task: %v", err))
		api_hnd.InternalServerError(res, "Failed to update RMM task.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "RMM task updated successfully.",
	})
    return
}

// CreateRmmExecution - Minion calls this to report command output
func (api_hnd *APIHandler) CreateRmmExecution(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		api_hnd.MethodNotAllowed(res, "POST")
		return
	}

	var exec dfn.RmmExecution
	if err := json.NewDecoder(req.Body).Decode(&exec); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Basic Validation
	if !utils.CheckifStringIsEmpty(exec.TaskID) || !utils.CheckifStringIsEmpty(exec.MinionID) {
		api_hnd.BadRequest(res, "TaskID and MinionID are required.")
		return
	}

	exec.ExecutionID = utils.GenerateUUID()
	exec.Touch()

	ctx := context.Background()
	if err := api_hnd.Dom.CreateRmmExecution(ctx, exec); err != nil {
		utils.Warning(fmt.Sprintf("Error logging RMM execution: %v", err))
		api_hnd.InternalServerError(res, "Failed to log execution result.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status":  "success",
		"message": "Execution result logged.",
	})
    return
}

// ViewRmmTask - Retrieves task details and its associated execution logs
func (api_hnd *APIHandler) ViewRmmTask(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	taskID := req.URL.Query().Get("task_id")
	if !utils.CheckifStringIsEmpty(taskID) {
		api_hnd.BadRequest(res, "task_id is required.")
		return
	}

	ud := req.Context().Value("userData").(*UserData)
	if ud == nil {
		api_hnd.Unauthorized(res, "User data not found")
		return
	}

	ctx := context.Background()
	task, execs, err := api_hnd.Dom.ViewRMMTask(ctx, taskID, ud.UserId)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error viewing RMM task: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch task details.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"task":       task,
			"executions": execs,
		},
	})
    return
}

// ListRmmTasks - Admin lists all their dispatched tasks
func (api_hnd *APIHandler) ListRmmTasks(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	ud := req.Context().Value("userData").(*UserData)
	if ud == nil {
		api_hnd.Unauthorized(res, "User data not found")
		return
	}

	ctx := context.Background()
	tasks, err := api_hnd.Dom.ListRmmTasks(ctx, ud.UserId)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing RMM tasks: %v", err))
		api_hnd.InternalServerError(res, "Failed to list tasks.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
        "status":  "success",
        "message": "RMM Tasks listed successfully.",
        "data":    tasks,
    })
    return
}

// MinionListRmmTasks - Minion polls for tasks based on its OS type
func (api_hnd *APIHandler) MinionListRmmTasks(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	ostype := req.URL.Query().Get("ostype")
	if !utils.CheckifStringIsEmpty(ostype) {
		api_hnd.BadRequest(res, "ostype is required.")
		return
	}

	ctx := context.Background()
	tasks, err := api_hnd.Dom.MinionListRMMTask(ctx, ostype)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error minion polling RMM tasks: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch tasks.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
        "status":  "success",
        "message": "RMM Tasks listed successfully.",
        "data":    tasks,
    })
    return
}

// ListRmmExecutions - View execution logs for a specific task (optionally filtered by verified status)
func (api_hnd *APIHandler) ListRmmExecutions(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		api_hnd.MethodNotAllowed(res, "GET")
		return
	}

	taskID := req.URL.Query().Get("task_id")
	verified := req.URL.Query().Get("verified") // e.g., "success" or "failed"

	if !utils.CheckifStringIsEmpty(taskID) {
		api_hnd.BadRequest(res, "task_id is required.")
		return
	}

	ctx := context.Background()
	execs, err := api_hnd.Dom.ListRmmExecutuions(ctx, taskID, verified)
	if err != nil {
		utils.Warning(fmt.Sprintf("Error listing RMM executions: %v", err))
		api_hnd.InternalServerError(res, "Failed to fetch executions.")
		return
	}

	api_hnd.DynamicResponse(res, map[string]interface{}{
        "status":  "success",
        "message": "RMM Executions listed successfully.",
        "data":    execs,
    })
    return
}