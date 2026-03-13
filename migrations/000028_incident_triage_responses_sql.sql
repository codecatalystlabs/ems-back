-- +goose Up
-- +goose StatementBegin

CREATE TABLE incident_triage_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    questionnaire_id UUID NOT NULL REFERENCES triage_questionnaires(id) ON DELETE RESTRICT,
    triage_mode TEXT NOT NULL DEFAULT 'PRIMARY' CHECK (triage_mode IN ('PRIMARY', 'RETRIAGE', 'SECONDARY', 'MANUAL_OVERRIDE')),
    total_score INT NOT NULL DEFAULT 0,
    boolean_true_count INT NOT NULL DEFAULT 0,
    auto_dispatch_eligible BOOLEAN NOT NULL DEFAULT FALSE,
    derived_priority_level_id UUID REFERENCES ref_priority_levels(id),
    notes TEXT,
    triaged_by_user_id UUID REFERENCES users(id),
    triaged_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_incident_triage_sessions_incident_id
    ON incident_triage_sessions(incident_id, triaged_at DESC);

CREATE INDEX idx_incident_triage_sessions_questionnaire_id
    ON incident_triage_sessions(questionnaire_id);

CREATE INDEX idx_incident_triage_sessions_priority_level_id
    ON incident_triage_sessions(derived_priority_level_id);

CREATE TABLE incident_triage_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    triage_session_id UUID NOT NULL REFERENCES incident_triage_sessions(id) ON DELETE CASCADE,
    incident_id UUID NOT NULL REFERENCES incidents(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES triage_questions(id) ON DELETE RESTRICT,
    question_code TEXT NOT NULL,
    response_type TEXT NOT NULL CHECK (response_type IN ('BOOLEAN', 'INTEGER', 'TEXT', 'SINGLE_CHOICE')),
    response_value_text TEXT,
    response_value_bool BOOLEAN,
    response_value_int INT,
    selected_option_id UUID REFERENCES triage_question_options(id),
    selected_option_code TEXT,
    score_awarded INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(triage_session_id, question_id)
);

CREATE INDEX idx_incident_triage_responses_session_id
    ON incident_triage_responses(triage_session_id);

CREATE INDEX idx_incident_triage_responses_incident_id
    ON incident_triage_responses(incident_id);

CREATE INDEX idx_incident_triage_responses_question_code
    ON incident_triage_responses(question_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS incident_triage_responses;
DROP TABLE IF EXISTS incident_triage_sessions;
-- +goose StatementEnd
