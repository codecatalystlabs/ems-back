package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"dispatch/internal/platform/events"
	"dispatch/internal/shared/constants"
)

type BroadcastService struct {
	bus events.Publisher
}

func NewBroadcastService(bus events.Publisher) *BroadcastService {
	return &BroadcastService{bus: bus}
}

func (s *BroadcastService) RequestBroadcast(
	ctx context.Context,
	bloodRequisitionID string,
	incidentID *string,
	bloodGroupCode, bloodProductCode string,
	unitsRequested int,
	urgencyLevel string,
	destinationFacilityID *string,
	destinationFacilityName string,
	clinicalSummary string,
	requestedByUserID *string,
) error {
	msg := fmt.Sprintf(
		"Urgent blood request: %d unit(s) of %s %s needed for %s. Case: %s",
		unitsRequested, bloodProductCode, bloodGroupCode, destinationFacilityName, clinicalSummary,
	)

	return s.bus.Publish(ctx, constants.TopicBloodBroadcastRequested, events.Event{
		ID:          uuid.NewString(),
		Topic:       constants.TopicBloodBroadcastRequested,
		AggregateID: bloodRequisitionID,
		Type:        constants.TopicBloodBroadcastRequested,
		OccurredAt:  time.Now().UTC(),
		Payload: map[string]any{
			"blood_requisition_id":      bloodRequisitionID,
			"incident_id":               incidentID,
			"blood_group_code":          bloodGroupCode,
			"blood_product_code":        bloodProductCode,
			"units_requested":           unitsRequested,
			"urgency_level":             urgencyLevel,
			"destination_facility_id":   destinationFacilityID,
			"destination_facility_name": destinationFacilityName,
			"clinical_summary":          clinicalSummary,
			"request_message":           msg,
			"requested_by_user_id":      requestedByUserID,
			"requested_at":              time.Now().UTC(),
		},
	})
}
