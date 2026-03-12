package application

import (
	"context"
	"strings"
	"time"

	"dispatch/internal/modules/incidents/domain"
	platformdb "dispatch/internal/platform/db"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) List(ctx context.Context, p platformdb.Pagination) (platformdb.PageResult[domain.Incident], error) {
	items, total, err := s.repo.ListIncidents(ctx, p)
	if err != nil {
		return platformdb.PageResult[domain.Incident]{}, err
	}
	return platformdb.PageResult[domain.Incident]{
		Items: items,
		Meta:  platformdb.NewPageMeta(p, total),
	}, nil
}

func (s *Service) Create(ctx context.Context, req CreateIncidentRequest, createdBy *string) (domain.Incident, error) {
	now := time.Now().UTC()

	// Resolve incident type: allow UUID, code, or name.
	incidentTypeID, err := s.repo.ResolveIncidentTypeID(ctx, req.IncidentTypeID)
	if err != nil {
		return domain.Incident{}, err
	}

	// Resolve optional district and facility IDs from human-friendly values.
	var districtID *string
	if req.DistrictID != nil {
		if id, err := s.repo.ResolveDistrictID(ctx, *req.DistrictID); err == nil {
			districtID = id
		}
	}

	var facilityID *string
	if req.FacilityID != nil {
		if id, err := s.repo.ResolveFacilityID(ctx, *req.FacilityID); err == nil {
			facilityID = id
		}
	}

	incident := domain.Incident{
		ID:                 uuid.NewString(),
		IncidentNumber:     s.generateIncidentNumber(now),
		SourceChannel:      strings.ToUpper(strings.TrimSpace(req.SourceChannel)),
		CallerName:         req.CallerName,
		CallerPhone:        req.CallerPhone,
		PatientName:        req.PatientName,
		PatientPhone:       req.PatientPhone,
		PatientAgeGroup:    req.PatientAgeGroup,
		PatientSex:         req.PatientSex,
		IncidentTypeID:     incidentTypeID,
		Summary:            req.Summary,
		Description:        req.Description,
		DistrictID:         districtID,
		FacilityID:         facilityID,
		Village:            req.Village,
		Parish:             req.Parish,
		Subcounty:          req.Subcounty,
		Landmark:           req.Landmark,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		VerificationStatus: "PENDING",
		Status:             "NEW",
		ReportedAt:         now,
		CreatedByUserID:    createdBy,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return s.repo.CreateIncident(ctx, incident)
}

func (s *Service) generateIncidentNumber(t time.Time) string {
	// Simple timestamp-based incident number, e.g. INC-20260310-123456
	return t.Format("INC-20060102-150405")
}

func (s *Service) Get(ctx context.Context, id string) (domain.Incident, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateIncidentRequest) (domain.Incident, error) {
	return s.repo.UpdateIncident(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.DeleteIncident(ctx, id)
}
