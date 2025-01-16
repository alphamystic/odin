package services

import (
  dom"github.com/alphamystic/odin/lib/domain"
  "sync"
)

type (
  Notify interface {
    CreateNotification()
    ListNotification()
  }
  NotificationService struct{
    Dom *dom.Domain
    Notifications map[string]*Notification
    mu sync.RWMutex
  }
  Notification struct{
    UserID string `json:"userid"`
    Title string  `json:"title"`
    NType NotificationType  `json:"ntype"`
    Handled bool  `json:"handled"`
  }
)

type NotificationType int

const (
  EVENT NotificationType = iota
)

func CreateNotifyer(domain *dom.Domain) *NotificationService {
  return &NotificationService{
    Dom: domain,
    Notifications: make(map[string]*Notification),
    mu: sync.RWMutex{},
  }
}

func (ns *NotificationService) WriteTOFile() error{
  //purged handled issues
  //write to file
  return nil
}

// load from file
func (ns *NotificationService) LoadFromFile() error {
  return nil
}

func (ns *NotificationService) ListUserNotifications(uid string)([]Notification,error){
  return nil,nil
}
