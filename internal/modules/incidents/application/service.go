package application

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	incidentdomain "dispatch/internal/modules/incidents/domain"
	platformdb "dispatch/internal/platform/db"
	"dispatch/internal/platform/events"
)

type Service struct {
	repo Repository
	bus  events.Publisher
	log  *zap.Logger
}

func NewService(repo Repository, bus events.Publisher, log *zap.Logger) *Service {
	return &Service{repo: repo, bus: bus, log: log}
}

func (s *Service) CreateIncident(ctx context.Context, req CreateIncidentRequest) (CreateIncidentResponse, error) {
	incidentNumber, err := s.repo.NextIncidentNumber(ctx)
	if err != nil {
		return CreateIncidentResponse{}, err
	}
	status := "NEW"
	verificationStatus := "PENDING"
	now := time.Now()

	inc := incidentdomain.Incident{
		ID:                 uuid.NewString(),
		IncidentNumber:     incidentNumber,
		SourceChannel:      strings.ToUpper(strings.TrimSpace(req.SourceChannel)),
		CallerName:         req.CallerName,
		CallerPhone:        req.CallerPhone,
		PatientName:        req.PatientName,
		PatientPhone:       req.PatientPhone,
		PatientAgeGroup:    req.PatientAgeGroup,
		PatientSex:         strings.ToUpper(strings.TrimSpace(req.PatientSex)),
		IncidentTypeID:     req.IncidentTypeID,
		SeverityLevelID:    req.SeverityLevelID,
		PriorityLevelID:    req.PriorityLevelID,
		Summary:            req.Summary,
		Description:        req.Description,
		DistrictID:         req.DistrictID,
		FacilityID:         req.FacilityID,
		Village:            req.Village,
		Parish:             req.Parish,
		Subcounty:          req.Subcounty,
		Landmark:           req.Landmark,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		VerificationStatus: verificationStatus,
		Status:             status,
		ReportedAt:         now,
		CreatedByUserID:    req.CreatedByUserID,
	}
	created, err := s.repo.CreateIncident(ctx, inc)
	if err != nil {
		return CreateIncidentResponse{}, err
	}
	_ = s.repo.CreateIncidentUpdate(ctx, created.ID, "COMMENT", "", "", "incident created", req.CreatedByUserID)

	resp := CreateIncidentResponse{
		Incident:                   created,
		AutoDispatchEligible:       false,
		DispatchRecommendationHint: "manual review",
	}

	if len(req.TriageResponses) > 0 {
		questionnaireCode := strings.ToUpper(strings.TrimSpace(req.QuestionnaireCode))
		if questionnaireCode == "" {
			questionnaireCode = "EMS_PRIMARY_TRIAGE"
		}
		triage, triageErr := s.persistTriageOnCreate(ctx, created.ID, questionnaireCode, req.TriageResponses, req.TriageNotes, req.CreatedByUserID)
		if triageErr != nil {
			return CreateIncidentResponse{}, triageErr
		}
		updatedIncident, _ := s.repo.GetIncidentByID(ctx, created.ID)
		resp.Incident = updatedIncident
		resp.TriageSession = &triage
		resp.AutoDispatchEligible = triage.AutoDispatchEligible
		if triage.AutoDispatchEligible || triage.DerivedPriorityCode == "RED" || triage.DerivedPriorityCode == "ORANGE" {
			resp.DispatchRecommendationHint = "eligible for dispatch recommendations"
		} else {
			resp.DispatchRecommendationHint = "manual review"
		}
	}

	_ = s.bus.Publish(ctx, "incident.created", events.Event{
		ID:          uuid.NewString(),
		Topic:       "incident.created",
		AggregateID: created.ID,
		Type:        "incident.created",
		OccurredAt:  now,
		Payload: map[string]any{
			"incident_id":     created.ID,
			"incident_number": created.IncidentNumber,
			"source_channel":  created.SourceChannel,
			"priority_code":   resp.Incident.PriorityCode,
			"auto_dispatch":   resp.AutoDispatchEligible,
		},
	})

	return resp, nil
}

func (s *Service) persistTriageOnCreate(ctx context.Context, incidentID, questionnaireCode string, inputs []TriageResponseInput, notes string, actorUserID *string) (incidentdomain.PersistedTriageSession, error) {
	questionnaireID, err := s.repo.ResolveQuestionnaireIDByCode(ctx, questionnaireCode)
	if err != nil {
		return incidentdomain.PersistedTriageSession{}, err
	}
	defs, err := s.repo.GetQuestionDefinitions(ctx, questionnaireCode)
	if err != nil {
		return incidentdomain.PersistedTriageSession{}, err
	}

	session := incidentdomain.PersistedTriageSession{
		ID:                uuid.NewString(),
		IncidentID:        incidentID,
		QuestionnaireID:   questionnaireID,
		QuestionnaireCode: questionnaireCode,
		TriageMode:        "PRIMARY",
		Notes:             notes,
		TriagedByUserID:   actorUserID,
		TriagedAt:         time.Now().UTC(),
	}

	responses := make([]incidentdomain.PersistedTriageResponse, 0, len(inputs))
	booleanTrueCount := 0
	totalScore := 0

	for _, in := range inputs {
		code := strings.ToUpper(strings.TrimSpace(in.QuestionCode))
		raw := strings.TrimSpace(in.ResponseValue)
		def, ok := defs[code]
		if !ok {
			continue
		}
		resp := incidentdomain.PersistedTriageResponse{
			QuestionID:   def.QuestionID,
			QuestionCode: code,
			ResponseType: def.ResponseType,
		}

		switch def.ResponseType {
		case "BOOLEAN":
			v := strings.EqualFold(raw, "true") || strings.EqualFold(raw, "yes")
			resp.ResponseValueBool = &v
			txt := strings.ToLower(raw)
			resp.ResponseValueText = &txt
			if v {
				booleanTrueCount++
			}
			if def.TrueScore != nil && v {
				resp.ScoreAwarded = *def.TrueScore
			}
			if def.FalseScore != nil && !v {
				resp.ScoreAwarded = *def.FalseScore
			}
		case "INTEGER":
			n, err := strconv.Atoi(raw)
			if err == nil {
				resp.ResponseValueInt = &n
				txt := raw
				resp.ResponseValueText = &txt
				switch {
				case n >= 5:
					resp.ScoreAwarded = 90
				case n >= 3:
					resp.ScoreAwarded = 50
				case n >= 1:
					resp.ScoreAwarded = 10
				}
			}
		default:
			txt := raw
			resp.ResponseValueText = &txt
		}

		totalScore += resp.ScoreAwarded
		responses = append(responses, resp)
	}

	session.BooleanTrueCount = booleanTrueCount
	session.TotalScore = totalScore
	session.AutoDispatchEligible = booleanTrueCount >= 3

	priorityCode := "GREEN"
	switch {
	case session.AutoDispatchEligible:
		priorityCode = "RED"
	case totalScore >= 90:
		priorityCode = "RED"
	case totalScore >= 40:
		priorityCode = "ORANGE"
	default:
		priorityCode = "GREEN"
	}
	priorityID, err := s.repo.ResolvePriorityLevelIDByCode(ctx, priorityCode)
	if err != nil {
		return incidentdomain.PersistedTriageSession{}, err
	}
	session.DerivedPriorityLevelID = priorityID
	session.DerivedPriorityCode = priorityCode
	session.Responses = responses

	created, err := s.repo.CreatePersistedTriageSession(ctx, session)
	if err != nil {
		return incidentdomain.PersistedTriageSession{}, err
	}
	if err := s.repo.SetIncidentPriorityByCode(ctx, incidentID, priorityCode); err != nil {
		return incidentdomain.PersistedTriageSession{}, err
	}
	_ = s.repo.SetIncidentTriageSummary(ctx, incidentID, actorUserID)
	if created.AutoDispatchEligible || priorityCode == "RED" || priorityCode == "ORANGE" {
		_, _ = s.repo.UpdateIncidentStatus(ctx, incidentID, "AWAITING_ASSIGNMENT")
	}
	_ = s.repo.CreateIncidentUpdate(ctx, incidentID, "TRIAGE", "", priorityCode, fmt.Sprintf("triage persisted on incident creation: score=%d, boolean_true_count=%d", totalScore, booleanTrueCount), actorUserID)
	return created, nil
}

func (s *Service) GetIncidentByID(ctx context.Context, id string) (incidentdomain.Incident, error) {
	return s.repo.GetIncidentByID(ctx, id)
}

func (s *Service) ListIncidents(ctx context.Context, params ListIncidentsParams) (platformdb.PageResult[incidentdomain.Incident], error) {
	items, total, err := s.repo.ListIncidents(ctx, params)
	if err != nil {
		return platformdb.PageResult[incidentdomain.Incident]{}, err
	}
	return platformdb.PageResult[incidentdomain.Incident]{Items: items, Meta: platformdb.NewPageMeta(params.Pagination, total)}, nil
}

func (s *Service) UpdateIncidentStatus(ctx context.Context, id string, req UpdateIncidentStatusRequest, actorUserID *string) (incidentdomain.Incident, error) {
	updated, err := s.repo.UpdateIncidentStatus(ctx, id, strings.ToUpper(strings.TrimSpace(req.Status)))
	if err != nil {
		return incidentdomain.Incident{}, err
	}
	_ = s.repo.CreateIncidentUpdate(ctx, id, "STATUS_CHANGE", "", updated.Status, req.Notes, actorUserID)
	return updated, nil
}
