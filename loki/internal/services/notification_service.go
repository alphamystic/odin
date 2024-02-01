package services

import (
  "sync"
)

type (
  Notify interface {
    CreateNotification()
    ListNotification()
  }
  NotificationService struct{
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

func CreateNotifyer() *NotificationService {
  return &NotificationService{
    Notifications: male(map[string]*Notification),
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
