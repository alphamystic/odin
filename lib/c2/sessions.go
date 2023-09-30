package c2

import (
	"fmt"
	"sync"
	"time"
	"errors"
	"strings"
	//"io/ioutil"
  //"path/filepath"
  "encoding/json"

	"github.com/alphamystic/odin/lib/db"
	"github.com/alphamystic/odin/lib/utils"
)

var (
	SNF = errors.New("Session not found")
)

const SessionsPath = "../.brain/"

type Session struct {
	ID string
	MotherShipID string
	Expiry string
	Active bool
	SessionID string
}

// SessionManager is responsible for managing active connections doesn't matter the protocol
// and ensuring that all connected clients have up-to-date information
/*type SessionManager struct {
	// The mutex to protect the session map
	mu sync.Mutex
	// The map of active sessions
	sessions map[net.Addr]*Session
}
*/
// change to seessions to remove address limiting for mutli agents
type SessionManager struct {
	Sessions map[string]*Session
	mu sync.RWMutex
}

//initialize a new empty data type of SessionManager (allowsus to save,load = NewSession,olsSesions)
func InitializeNewSessionManager() *SessionManager {
	return &SessionManager{
		Sessions: make(map[string]*Session),
		mu: sync.RWMutex{},
	}
}

// NewSession creates a new session
func (sm *SessionManager) NewSession(id,msid,expiry string,driver *db.Driver) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	ssid := utils.Md5Hash(utils.RandString(10))
	s := &Session{
		ID: id,
		MotherShipID: msid,
		Expiry:  expiry,
		Active: true,
		SessionID: ssid,
	}
	sm.Sessions[ssid] = s
	if err := sm.SaveSession(s,driver); err != nil{
		utils.Warning(fmt.Sprintf("Error saving newly created session %s to db\n.ERROR: %s",s.SessionID,err))
		utils.PrintTextInASpecificColor("cyan","Save manually with save --ses hlkjnkjjkjii where gibberish is Session ID.")
	}
	return s
}

// Add adds a new session to the manager
func (sm *SessionManager) Add(session *Session) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	for _,ses := range sm.Sessions{
		if ses.SessionID == session.SessionID {
			return errors.New("Session with ID already exists")
		}
	}
	sm.Sessions[session.SessionID] = session
	return nil
}
// update a session in sessions
func (sm *SessionManager) UpdateSession(session *Session,driver *db.Driver) error{
	sm.mu.Lock()
	for _,ses := range sm.Sessions {
		if ses.SessionID == session.SessionID {
			err := sm.DeleteSession(session.SessionID,driver)
			if err != nil {
				return err
			}
			// unlock to prevent dead lock
			sm.mu.Unlock()
			err = sm.Add(session)
			if err != nil{
				return err
			}
			return sm.SaveSession(session,driver)
		}
	}
	return nil
}

// mark a given session as active or inactive
func (sm *SessionManager) MarkSessionAsActiveInactive(sid string,val bool,driver *db.Driver) error{
	ses,err := sm.GetSession(sid)
	if err != nil {
		return err
	}
	ses.Active = val
	err = sm.UpdateSession(ses,driver)
	if err != nil{
		return err
	}
	return nil
}

// Get retrieves a session by its id
func (sm *SessionManager) GetSession(id string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	s, ok := sm.Sessions[id]
	if !ok {
		//change to utils Danger
		return nil,SNF
	}
	return s, nil
}

// get sessions from a specific Mother Shiip
func (sm *SessionManager) GetSessionsFromMS(msid string)([]Session,error){
	var sessions []Session
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _,session := range sm.Sessions {
		var ses Session
		session = &ses
		if session.MotherShipID == msid {
			sessions = append(sessions,ses)
		}
	}
	if len(sessions) < 0 {
		return nil,errors.New("No session with specified MotherShip ID.")
	}
	return sessions,nil
}

func (sm *SessionManager) ListFromMS(msid string)(error){
	var sessions []Session
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	for _,session := range sm.Sessions {
		var ses Session
		session = &ses
		if session.MotherShipID == msid {
			sessions = append(sessions,ses)
		}
	}
	if len(sessions) < 0 {
		return errors.New("No session with specified MotherShip ID.")
	}
	utils.PrintTextInASpecificColorInBold("yellow",fmt.Sprintf("    **********    CURRENT  Minion Sessions  From %s  ********** ",msid))
	for _,s := range sessions {
		utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID:	%s",s.SessionID))
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mothership ID:	%s",s.MotherShipID))
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mule ID:	%s",s.ID))
		if s.Active {
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	True"))
		} else {
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	False"))
		}
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Expiry:	%s",s.Expiry))
		utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
  }
	return nil
}

// Delete removes a session from the manager
func (sm *SessionManager) DeleteSession(id string,driver *db.Driver) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.Sessions[id]; !ok {
		return SNF
	}
	err := driver.Delete("sessions",id)
	if err != nil {
		return fmt.Errorf("Error deleting session from db.\nERROR: %s",err)
	}
	delete(sm.Sessions, id)
	return nil
}

func (sm *SessionManager) SearchSession(name string){
	sm.mu.Lock()
	defer sm.mu.Unlock()
  for _,s := range sm.Sessions {
    if strings.Contains(s.ID,name){
			utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
	    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID:	%s",s.SessionID))
	    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mothership ID:	%s",s.MotherShipID))
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mule ID:	%s",s.ID))
			if s.Active {
				utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	True"))
			} else {
				utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	False"))
			}
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Expiry:	%s",s.Expiry))
	    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
    }
  }
}

func (sm *SessionManager) ListSessions(){
	utils.PrintTextInASpecificColorInBold("yellow","    **********    CURRENT  Minion Sessions    ********** ")
  for _,s := range sm.Sessions {
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   SessionID:	%s",s.SessionID))
    utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mothership ID:	%s",s.MotherShipID))
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Mule ID:	%s",s.ID))
		if s.Active {
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	True"))
		} else {
			utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Active:	False"))
		}
		utils.PrintTextInASpecificColorInBold("cyan",fmt.Sprintf("   Expiry:	%s",s.Expiry))
    utils.PrintTextInASpecificColorInBold("magenta","***********************************************************************")
  }
}

// Remove expired sessions
func (sm *SessionManager) PurgeExpiredSessions(driver *db.Driver) {
	now := time.Now()
	for id, s := range sm.Sessions {
  	sm.mu.Lock()
		expiry,_ := time.Parse(time.RFC3339,s.Expiry)
		if expiry.Before(now) {
			delete(sm.Sessions, id)
		}
		sm.mu.Unlock()
		err := sm.DeleteSession(id,driver)
		if err != nil{
			utils.Logerror(err)
			continue
		}
	}
}

// Close closes all active sessions (more like delete there presence from db. )
//Create a backup technique to store them somewhere else.
func (s *SessionManager) Close(driver *db.Driver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, session := range s.Sessions {
		err := s.DeleteSession(session.SessionID,driver)
		if err != nil{
			utils.Logerror(err)
			continue
		}
	}
}

func (s *SessionManager) SaveSession(session *Session,driver *db.Driver) error {
	return driver.Write("sessions",session.SessionID,session)
}

func (s *SessionManager) LoadSessions(driver *db.Driver) error {
	//type sessions []Session
	sessions,err := driver.ReadAll("sessions")
	if err != nil {
		return err
	}
	for _,sesn := range sessions {
		var ses Session
		if err := json.Unmarshal([]byte(sesn),&ses);err != nil{
			utils.Warning(fmt.Sprintf("%s",err));continue
		}
		err = s.Add(&ses)
		if err != nil{
			utils.Warning(fmt.Sprintf("%s",err))
			continue
		}
	}
	return nil
}

/*
func (s *SessionManager) LoadSessions(driver *db.Driver) error {
	sessions,err := driver.ReadAll("sessions")
	if err != nil {
		return err
	}
	for _,sesn := range sessions {
		var ses Session
		ses.ID = sesn.ID
		ses.MotherShipID = sesn.MotherShipID
		ses.Expiry = sesn.Expiry
		if sesn.Active == "true"{
			ses.Active = true
		} else {
			ses.Active = false
		}
		ses.SessionID = sesn.MotherShipID
		err = s.Add(ses)
		if err != nil{
			utils.Warning(fmt.Sprintf("%s",err))
			continue
		}
	}
	return nil
}

func (s *SessionManager) LoadSessions(driver *db.Driver) error{
	dir := filepath.Join(driver.Dir,"sessions")
	if _,err := driver.ExportedStat(dir); err != nil{
    return errors.New(fmt.Sprintf("Database with name %s does not exist: %s\n","sessions",err))
  }
	files,err := ioutil.ReadDir(dir)
	if err != nil{
    return fmt.Errorf("Error reading sessions directory: %v",err)
  }
	for _,file := range files{
		var ses Session
		b,err := ioutil.ReadFile(filepath.Join(dir,file.Name()))
		if err != nil{
			utils.NoticeError(fmt.Sprintf("No session with such name: %s\nERROR: %s",file.Name(),err))
			continue
		}
		err = json.Unmarshal(b,&ses)
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error unmarshalling to session: %s",err))
			continue
		}
		err = s.Add(&ses)
		if err != nil{
			utils.NoticeError(fmt.Sprintf("Error adding to session manager: %s",err))
			continue
		}
	}
	return nil
}
*/
