package application

import (
	"context"

	"dispatch/internal/modules/availability/application/dto"

	availabilitydomain "dispatch/internal/modules/availability/domain"
)

type Repository interface {
	CreateShift(ctx context.Context, in availabilitydomain.UserShift) (availabilitydomain.UserShift, error)
	UpdateShift(ctx context.Context, id string, req dto.UpdateShiftRequest) (availabilitydomain.UserShift, error)
	GetShiftByID(ctx context.Context, id string) (availabilitydomain.UserShift, error)
	ListShifts(ctx context.Context, params dto.ListShiftsParams) ([]availabilitydomain.UserShift, int64, error)

	UpsertAvailability(ctx context.Context, in availabilitydomain.UserAvailability) (availabilitydomain.UserAvailability, error)
	GetAvailabilityByUserID(ctx context.Context, userID string) (availabilitydomain.UserAvailability, error)
	ListAvailability(ctx context.Context, params dto.ListAvailabilityParams) ([]availabilitydomain.UserAvailability, int64, error)

	CreatePresenceLog(ctx context.Context, in availabilitydomain.UserPresenceLog) (availabilitydomain.UserPresenceLog, error)
	ListPresenceLogs(ctx context.Context, params dto.ListPresenceParams) ([]availabilitydomain.UserPresenceLog, int64, error)
}
