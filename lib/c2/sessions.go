package c2

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alphamystic/odin/lib/db"
	dfn "github.com/alphamystic/odin/lib/definers"
	"github.com/alphamystic/odin/lib/utils"
)

var (
	SNF = errors.New("Session not found")
)

const SessionsPath = "../.brain/"

type Session struct {
	Min       *dfn.Minion `json:"min"`
	Persisted bool        `json:"persisted"`
}

type SessionManager struct {
	Sessions map[string]*Session
	mu       sync.RWMutex
}

// InitializeNewSessionManager initializes a new empty data type of SessionManager
func InitializeNewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
		mu:       sync.RWMutex{},
	}
}

// NewSession creates a new session
func (sm *SessionManager) NewSession(minion *dfn.Minion) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	ssid := utils.Md5Hash(utils.RandString(10))
	s := &Session{
		Min: minion,
	}
	sm.Sessions[ssid] = s
	return s
}

// Add adds a new session to the manager
// func (sm *SessionManager) Add(session *Session) error {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()
// 	for _, ses := range sm.Sessions {
// 		if ses.Min.MinionID == session.Min.MinionID {
// 			return errors.New("Session with ID already exists")
// 		}
// 	}
// 	// FIXED: Changed from invalid `sm.Min.Sessions[...]` to storing in the manager's map directly via MinionID
// 	sm.Sessions[session.Min.MinionID] = session
// 	return nil
// }
// Add adds a new session to the manager safely
func (sm *SessionManager) Add(session *Session) error {
	if session == nil {
		return errors.New("cannot add a nil session reference")
	}

	// FIXED: Defensive safeguard to catch malformed database records
	if session.Min == nil {
		return errors.New("malformed database record: nested minion data structure configuration is missing")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, ses := range sm.Sessions {
		// FIXED: Check that the existing map entry contains a valid Minion payload pointer
		if ses != nil && ses.Min != nil {
			if ses.Min.MinionID == session.Min.MinionID {
				return errors.New("Session with ID already exists")
			}
		}
	}

	sm.Sessions[session.Min.MinionID] = session
	return nil
}

// UpdateSession updates a session in sessions
func (sm *SessionManager) UpdateSession(session *Session, driver *db.Driver) error {
	sm.mu.Lock()
	for _, ses := range sm.Sessions {
		if ses.Min.MinionID == session.Min.MinionID {
			err := sm.DeleteSession(session.Min.MinionID, driver)
			if err != nil {
				sm.mu.Unlock()
				return err
			}
			sm.mu.Unlock() // Unlock to prevent deadlock in sm.Add
			err = sm.Add(session)
			if err != nil {
				return err
			}
			return sm.SaveSession(session, driver)
		}
	}
	sm.mu.Unlock()
	return nil
}

// MarkSessionAsActiveInactive marks a given session as active or inactive
func (sm *SessionManager) MarkSessionAsActiveInactive(sid string, val bool, driver *db.Driver) error {
	ses, err := sm.GetSession(sid)
	if err != nil {
		return err
	}
	// FIXED: Accessing 'Active' via the nested 'Min' struct pointer
	ses.Min.Active = val
	err = sm.UpdateSession(ses, driver)
	if err != nil {
		return err
	}
	return nil
}

// GetSession retrieves a session by its id
func (sm *SessionManager) GetSession(id string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.Sessions[id]
	if !ok {
		return nil, SNF
	}
	return s, nil
}

// GetSessionsFromMS gets sessions from a specific MotherShip
func (sm *SessionManager) GetSessionsFromMS(msid string) ([]Session, error) {
	var sessions []Session
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, session := range sm.Sessions {
		// FIXED: Accessing 'MothershipID' safely via nested 'Min' field
		if session.Min != nil && session.Min.MothershipID == msid {
			sessions = append(sessions, *session)
		}
	}
	if len(sessions) == 0 { // FIXED: logic check for empty slices
		return nil, errors.New("No session with specified MotherShip ID.")
	}
	return sessions, nil
}

// ListFromMS prints and lists active minions connected to a target Mothership
func (sm *SessionManager) ListFromMS(msid string) error {
	var sessions []Session
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, session := range sm.Sessions {
		// FIXED: Safely parsing items to slice matching MSID
		if session.Min != nil && session.Min.MothershipID == msid {
			sessions = append(sessions, *session)
		}
	}
	if len(sessions) == 0 {
		return errors.New("No session with specified MotherShip ID.")
	}
	utils.PrintTextInASpecificColorInBold("yellow", fmt.Sprintf("    ********** CURRENT  Minion Sessions  From %s  ********** ", msid))
	for _, s := range sessions {
		utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
		utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   MinionID: %s", s.Min.MinionID))
		utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Mothership ID:    %s", s.Min.MothershipID))
		utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Owner ID:  %s", s.Min.OwnerID)) // FIXED: Field 'ID' doesn't exist on Minion; swapped with OwnerID
		if s.Min.Active {
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   True"))
		} else {
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   False"))
		}
		utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   OS:   %s", s.Min.Os))
		utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Last Seen:    %s", s.Min.LastSeen))
		utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
	}
	return nil
}

// DeleteSession removes a session from the manager
func (sm *SessionManager) DeleteSession(id string, driver *db.Driver) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.Sessions[id]; !ok {
		return SNF
	}
	err := driver.Delete("sessions", id)
	if err != nil {
		return fmt.Errorf("Error deleting session from db.\nERROR: %s", err)
	}
	delete(sm.Sessions, id)
	return nil
}

// DeleteNoneInstalledSessions handles cleanup logic
func (sm *SessionManager) DeleteNoneInstalledSessions(id string, driver *db.Driver) error { // FIXED: Added missing driver parameter
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.Sessions[id]; !ok {
		return SNF
	}
	if err := driver.Delete("sessions", id); err != nil {
		return fmt.Errorf("Error deleting session from db.\nERROR: %s", err)
	}
	delete(sm.Sessions, id)
	return nil
}

// SearchSession filters running sessions based on names/ids
func (sm *SessionManager) SearchSession(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _, s := range sm.Sessions {
		if s.Min != nil && strings.Contains(s.Min.MinionID, name) {
			utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   MinionID: %s", s.Min.MinionID))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Mothership ID:    %s", s.Min.MothershipID))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Owner ID:  %s", s.Min.OwnerID))
			if s.Min.Active {
				utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   True"))
			} else {
				utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   False"))
			}
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   OS:   %s", s.Min.Os))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Last Seen:    %s", s.Min.LastSeen))
			utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
		}
	}
}

// ListSessions loops and renders tracking details
func (sm *SessionManager) ListSessions(persisted bool) {
	utils.PrintTextInASpecificColorInBold("yellow", "    ********** CURRENT  Minion Sessions    ********** ")
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _, s := range sm.Sessions {
		if s.Persisted == persisted { // FIXED: Field corrected from 'Persist' to 'Persisted'
			utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   MinionID: %s", s.Min.MinionID))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Mothership ID:    %s", s.Min.MothershipID))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Owner ID:  %s", s.Min.OwnerID))
			if s.Min.Active {
				utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   True"))
			} else {
				utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Active:   False"))
			}
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   OS:   %s", s.Min.Os))
			utils.PrintTextInASpecificColorInBold("cyan", fmt.Sprintf("   Last Seen:    %s", s.Min.LastSeen))
			utils.PrintTextInASpecificColorInBold("magenta", "***********************************************************************")
		}
	}
}

// PurgeExpiredSessions drops inactive systems older than 6 months
func (sm *SessionManager) PurgeExpiredSessions(driver *db.Driver) {
	now := time.Now()
	// To prevent runtime deadlocks during iteration/mutations, find targets under read lock first
	sm.mu.RLock()
	var targets []string
	for id, s := range sm.Sessions {
		if s.Min != nil {
			expiry, err := time.Parse(time.RFC3339, s.Min.LastSeen)
			if err != nil {
				continue
			}
			limit := expiry.AddDate(0, 6, 0)
			if now.After(limit) {
				targets = append(targets, id)
			}
		}
	}
	sm.mu.RUnlock()

	// Execute removals cleanly using a write lock
	if len(targets) > 0 {
		sm.mu.Lock()
		for _, id := range targets {
			_ = driver.Delete("sessions", id)
			delete(sm.Sessions, id)
		}
		sm.mu.Unlock()
	}
}

// Close gracefully updates disk states on teardown
func (s *SessionManager) Close(driver *db.Driver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, session := range s.Sessions {
		if session.Min != nil {
			err := s.DeleteSession(session.Min.MinionID, driver)
			if err != nil {
				utils.Logerror(err)
				continue
			}
		}
	}
}

func (s *SessionManager) SaveSession(session *Session, driver *db.Driver) error {
	if session.Min == nil {
		return errors.New("cannot save a session with an empty Minion configuration")
	}
	return driver.Write("sessions", session.Min.MinionID, session)
}

func (s *SessionManager) LoadSessions(driver *db.Driver) error {
	sessions, err := driver.ReadAll("sessions")
	if err != nil {
		return err
	}
	for _, sesn := range sessions {
		var ses Session
		if err := json.Unmarshal([]byte(sesn), &ses); err != nil {
			utils.Warning(fmt.Sprintf("%s", err))
			continue
		}
		err = s.Add(&ses)
		if err != nil {
			utils.Warning(fmt.Sprintf("%s", err))
			continue
		}
	}
	return nil
}