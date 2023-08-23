package lib

import (
  "fmt"
  "sync"
)

type Event struct{
  Name string
  ID string
  EType EventType
  Handled bool
  CreatedAt string
  UpdatedAt string
}

// added to allow for future additions like compromised crojob/service
type EventType int
const (
  MUTANT EventType = iota
)

type EventManager struct{
  Events map[string]*Event
  mu sync.RWMutex
}

func InitilizeEventManager()*EventManager{
  return &EventManager{
    Events: make(map[string]*Event),
    mu: sync.RWMutex{},
  }
}

func (em *EventManager) AddEvent(event *Event) error{
  em.mu.Lock()
	defer em.mu.Unlock()
  for _, evts := range em.Events{
    if evts.ID == event.ID{
      return fmt.Errorf("Error event with id %s already exists.",event.ID)
    }
  }
  em.Events[event.ID] = event
  return nil
}

func (em *EventManager) RemoveEvent(id string)error{
  em.mu.Lock()
	defer em.mu.Unlock()
  for _,evt := range em.Events{
    if evt.ID == id{
      delete(em.Events,id)
      return nil
    }
  }
  return fmt.Errorf("No event with an id %s.",id)
}

func (em *EventManager) ListEvent(eventType EventType) []*Event{
  em.mu.Lock()
	defer em.mu.Unlock()
  var events []*Event
  for _, event := range em.Events {
    if event.EType == eventType{
      events = append(events,event)
    }
  }
  return events
}
