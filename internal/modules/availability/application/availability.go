// File: internal/modules/availability/application/service.go
package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"dispatch/internal/modules/availability/application/dto"
	availabilitydomain "dispatch/internal/modules/availability/domain"
	platformdb "dispatch/internal/platform/db"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateShift(ctx context.Context, req dto.CreateShiftRequest) (availabilitydomain.UserShift, error) {
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return availabilitydomain.UserShift{}, err
	}
	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		return availabilitydomain.UserShift{}, err
	}
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = "SCHEDULED"
	}
	in := availabilitydomain.UserShift{
		ID:         uuid.NewString(),
		UserID:     req.UserID,
		ShiftDate:  req.ShiftDate,
		StartsAt:   startsAt,
		EndsAt:     endsAt,
		ShiftType:  req.ShiftType,
		DistrictID: req.DistrictID,
		FacilityID: req.FacilityID,
		Status:     status,
		CreatedBy:  req.CreatedBy,
	}
	return s.repo.CreateShift(ctx, in)
}

func (s *Service) UpdateShift(ctx context.Context, id string, req dto.UpdateShiftRequest) (availabilitydomain.UserShift, error) {
	if req.Status != nil {
		v := strings.ToUpper(strings.TrimSpace(*req.Status))
		req.Status = &v
	}
	return s.repo.UpdateShift(ctx, id, req)
}

func (s *Service) GetShiftByID(ctx context.Context, id string) (availabilitydomain.UserShift, error) {
	return s.repo.GetShiftByID(ctx, id)
}

func (s *Service) ListShifts(ctx context.Context, params dto.ListShiftsParams) (platformdb.PageResult[availabilitydomain.UserShift], error) {
	items, total, err := s.repo.ListShifts(ctx, params)
	if err != nil {
		return platformdb.PageResult[availabilitydomain.UserShift]{}, err
	}
	return platformdb.PageResult[availabilitydomain.UserShift]{Items: items, Meta: platformdb.NewPageMeta(params.Pagination, total)}, nil
}

func (s *Service) UpsertAvailability(ctx context.Context, req dto.UpsertAvailabilityRequest) (availabilitydomain.UserAvailability, error) {
	status := strings.ToUpper(strings.TrimSpace(req.AvailabilityStatus))
	if req.Source == "" {
		req.Source = "MANUAL"
	}
	var lastSeenAt *time.Time
	if req.LastSeenAt != nil && *req.LastSeenAt != "" {
		parsed, err := time.Parse(time.RFC3339, *req.LastSeenAt)
		if err != nil {
			return availabilitydomain.UserAvailability{}, err
		}
		lastSeenAt = &parsed
		if strings.ToUpper(strings.TrimSpace(req.Source)) == "" {
			req.Source = "APP"
		}
	}
	in := availabilitydomain.UserAvailability{
		ID:                          uuid.NewString(),
		UserID:                      req.UserID,
		AvailabilityStatus:          status,
		Dispatchable:                req.Dispatchable,
		CurrentIncidentID:           req.CurrentIncidentID,
		CurrentDispatchAssignmentID: req.CurrentDispatchAssignmentID,
		CurrentAmbulanceID:          req.CurrentAmbulanceID,
		LastSeenAt:                  lastSeenAt,
		Source:                      strings.ToUpper(strings.TrimSpace(req.Source)),
		Notes:                       req.Notes,
		UpdatedBy:                   req.UpdatedBy,
		UpdatedAt:                   time.Now().UTC(),
	}
	return s.repo.UpsertAvailability(ctx, in)
}

func (s *Service) GetAvailabilityByUserID(ctx context.Context, userID string) (availabilitydomain.UserAvailability, error) {
	return s.repo.GetAvailabilityByUserID(ctx, userID)
}

func (s *Service) ListAvailability(ctx context.Context, params dto.ListAvailabilityParams) (platformdb.PageResult[availabilitydomain.UserAvailability], error) {
	items, total, err := s.repo.ListAvailability(ctx, params)
	if err != nil {
		return platformdb.PageResult[availabilitydomain.UserAvailability]{}, err
	}
	return platformdb.PageResult[availabilitydomain.UserAvailability]{Items: items, Meta: platformdb.NewPageMeta(params.Pagination, total)}, nil
}

func (s *Service) CreatePresenceLog(ctx context.Context, req dto.CreatePresenceLogRequest) (availabilitydomain.UserPresenceLog, error) {
	in := availabilitydomain.UserPresenceLog{
		ID:        uuid.NewString(),
		UserID:    req.UserID,
		Channel:   strings.ToUpper(strings.TrimSpace(req.Channel)),
		SeenAt:    time.Now().UTC(),
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}
	return s.repo.CreatePresenceLog(ctx, in)
}

func (s *Service) ListPresenceLogs(ctx context.Context, params dto.ListPresenceParams) (platformdb.PageResult[availabilitydomain.UserPresenceLog], error) {
	items, total, err := s.repo.ListPresenceLogs(ctx, params)
	if err != nil {
		return platformdb.PageResult[availabilitydomain.UserPresenceLog]{}, err
	}
	return platformdb.PageResult[availabilitydomain.UserPresenceLog]{Items: items, Meta: platformdb.NewPageMeta(params.Pagination, total)}, nil
}
