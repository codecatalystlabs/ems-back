package application

import (
	"context"

	devicedomain "dispatch/internal/modules/device_tokens/domain"
)

type Repository interface {
	Register(ctx context.Context, in devicedomain.DeviceToken) (devicedomain.DeviceToken, error)
	Update(ctx context.Context, id string, req UpdateDeviceTokenRequest) (devicedomain.DeviceToken, error)
	Deactivate(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (devicedomain.DeviceToken, error)
	List(ctx context.Context, params ListDeviceTokensParams) ([]devicedomain.DeviceToken, int64, error)
	GetPushTokensByUserID(ctx context.Context, userID string) ([]string, error)
}
