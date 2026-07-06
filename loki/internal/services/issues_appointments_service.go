package services

import (
	"context"
	"fmt"
	dfn "github.com/alphamystic/odin/lib/definers"
)

type IssuesAppointment interface {
	CreateIA(ctx context.Context, app dfn.Appointment) error
	ViewIA(ctx context.Context, appID string) (dfn.Appointment, error)
	ListIA(ctx context.Context, userID string) ([]dfn.Appointment, error)
	UpdateIA(ctx context.Context, app dfn.Appointment) error
	DeleteIA(ctx context.Context, appID string) error
}

type IssuesAppointmentService struct{ SAC *ServerAPIConnector }

func NewIAService(sac *ServerAPIConnector) *IssuesAppointmentService {
	return &IssuesAppointmentService{SAC: sac}
}

func (s *IssuesAppointmentService) CreateIA(ctx context.Context, app dfn.Appointment) error {
	_, err := s.SAC.Post("/api/appointments/create", app)
	return err
}

func (s *IssuesAppointmentService) ListIA(ctx context.Context, userID string) ([]dfn.Appointment, error) {
	data, err := s.SAC.Get(fmt.Sprintf("/api/appointments/list?userid=%s", userID))
	var apps []dfn.Appointment
	err = s.SAC.Decode(data, &apps)
	return apps, err
}