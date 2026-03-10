-- +goose Up
CREATE TABLE triage_questionnaires (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE triage_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    questionnaire_id UUID NOT NULL REFERENCES triage_questionnaires(id) ON DELETE CASCADE,
    code TEXT NOT NULL UNIQUE,
    question_text TEXT NOT NULL,
    response_type TEXT NOT NULL CHECK (response_type IN ('BOOLEAN', 'INTEGER', 'TEXT', 'SINGLE_CHOICE')),
    display_order INT NOT NULL,
    is_required BOOLEAN NOT NULL DEFAULT TRUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE triage_question_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES triage_questions(id) ON DELETE CASCADE,
    option_code TEXT NOT NULL,
    option_label TEXT NOT NULL,
    option_value TEXT NOT NULL,
    score INT NOT NULL DEFAULT 0,
    display_order INT NOT NULL DEFAULT 1,
    UNIQUE(question_id, option_code)
);

CREATE TABLE triage_priority_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    questionnaire_id UUID NOT NULL REFERENCES triage_questionnaires(id) ON DELETE CASCADE,
    priority_level_id UUID NOT NULL REFERENCES ref_priority_levels(id) ON DELETE CASCADE,
    rule_name TEXT NOT NULL,
    rule_type TEXT NOT NULL CHECK (rule_type IN ('ALL', 'ANY', 'SCORE_RANGE')),
    min_score INT,
    max_score INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE triage_priority_rule_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_id UUID NOT NULL REFERENCES triage_priority_rules(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES triage_questions(id) ON DELETE CASCADE,
    operator TEXT NOT NULL CHECK (operator IN ('EQ', 'NEQ', 'GT', 'GTE', 'LT', 'LTE')),
    expected_value TEXT NOT NULL,
    score_delta INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO triage_questionnaires (code, name, description)
VALUES (
    'EMS_PRIMARY_TRIAGE',
    'EMS Primary Triage',
    'Primary triage questionnaire for dispatch priority grading.'
)
ON CONFLICT (code) DO NOTHING;

INSERT INTO triage_questions (questionnaire_id, code, question_text, response_type, display_order, is_required)
SELECT tq.id, v.code, v.question_text, v.response_type, v.display_order, TRUE
FROM triage_questionnaires tq
JOIN (
    VALUES
        ('IS_PATIENT_BREATHING', 'Is the patient breathing?', 'BOOLEAN', 1),
        ('IS_PATIENT_CONSCIOUS', 'Is the patient conscious?', 'BOOLEAN', 2),
        ('CAN_PATIENT_SPEAK', 'Can the patient speak?', 'BOOLEAN', 3),
        ('HAS_SEVERE_BLEEDING', 'Is there severe bleeding?', 'BOOLEAN', 4),
        ('HOW_MANY_PEOPLE_INJURED', 'How many people are injured?', 'INTEGER', 5)
) AS v(code, question_text, response_type, display_order)
ON tq.code = 'EMS_PRIMARY_TRIAGE'
ON CONFLICT (code) DO NOTHING;

INSERT INTO triage_question_options (question_id, option_code, option_label, option_value, score, display_order)
SELECT q.id, o.option_code, o.option_label, o.option_value, o.score, o.display_order
FROM triage_questions q
JOIN (
    VALUES
        ('IS_PATIENT_BREATHING', 'YES', 'Yes', 'true', 0, 1),
        ('IS_PATIENT_BREATHING', 'NO', 'No', 'false', 100, 2),
        ('IS_PATIENT_CONSCIOUS', 'YES', 'Yes', 'true', 0, 1),
        ('IS_PATIENT_CONSCIOUS', 'NO', 'No', 'false', 80, 2),
        ('CAN_PATIENT_SPEAK', 'YES', 'Yes', 'true', 0, 1),
        ('CAN_PATIENT_SPEAK', 'NO', 'No', 'false', 50, 2),
        ('HAS_SEVERE_BLEEDING', 'YES', 'Yes', 'true', 90, 1),
        ('HAS_SEVERE_BLEEDING', 'NO', 'No', 'false', 0, 2)
) AS o(question_code, option_code, option_label, option_value, score, display_order)
ON q.code = o.question_code
ON CONFLICT (question_id, option_code) DO NOTHING;

INSERT INTO triage_priority_rules (questionnaire_id, priority_level_id, rule_name, rule_type, min_score, max_score)
SELECT tq.id, rpl.id, x.rule_name, x.rule_type, x.min_score, x.max_score
FROM triage_questionnaires tq
JOIN (
    VALUES
        ('RED', 'Red score rule', 'SCORE_RANGE', 90, NULL),
        ('ORANGE', 'Orange score rule', 'SCORE_RANGE', 40, 89),
        ('GREEN', 'Green score rule', 'SCORE_RANGE', 0, 39)
) AS x(priority_code, rule_name, rule_type, min_score, max_score)
ON TRUE
JOIN ref_priority_levels rpl ON rpl.code = x.priority_code
WHERE tq.code = 'EMS_PRIMARY_TRIAGE';

INSERT INTO triage_priority_rule_conditions (rule_id, question_id, operator, expected_value, score_delta)
SELECT r.id, q.id, 'EQ', 'false', 0
FROM triage_priority_rules r
JOIN triage_questionnaires tq ON tq.id = r.questionnaire_id AND tq.code = 'EMS_PRIMARY_TRIAGE'
JOIN triage_questions q ON q.questionnaire_id = tq.id AND q.code = 'IS_PATIENT_BREATHING'
WHERE r.rule_name = 'Red score rule';

INSERT INTO triage_priority_rule_conditions (rule_id, question_id, operator, expected_value, score_delta)
SELECT r.id, q.id, 'EQ', 'true', 0
FROM triage_priority_rules r
JOIN triage_questionnaires tq ON tq.id = r.questionnaire_id AND tq.code = 'EMS_PRIMARY_TRIAGE'
JOIN triage_questions q ON q.questionnaire_id = tq.id AND q.code = 'HAS_SEVERE_BLEEDING'
WHERE r.rule_name = 'Red score rule';

-- +goose Down
DROP TABLE IF EXISTS triage_priority_rule_conditions;
DROP TABLE IF EXISTS triage_priority_rules;
DROP TABLE IF EXISTS triage_question_options;
DROP TABLE IF EXISTS triage_questions;
DROP TABLE IF EXISTS triage_questionnaires;