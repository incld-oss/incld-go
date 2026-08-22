package client

import "time"

type ScheduleStatus string

const (
	StatusActive  ScheduleStatus = "active"
	StatusPaused  ScheduleStatus = "paused"
	StatusDeleted ScheduleStatus = "deleted"
)

type OverlapPolicy string

const (
	OverlapPolicyAllow OverlapPolicy = "allow"
	OverlapPolicySkip  OverlapPolicy = "skip"
)

type MisfirePolicy string

const (
	MisfirePolicySkip    MisfirePolicy = "skip"
	MisfirePolicyRunOnce MisfirePolicy = "run_once"
	MisfirePolicyCatchUp MisfirePolicy = "catch_up"
)

type Frequency string

const (
	FrequencyOnce    Frequency = "once"
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
)

type EndCondition struct {
	Type  string `json:"type"`
	Date  string `json:"date,omitempty"`
	Count int    `json:"count,omitempty"`
}

type Recurrence struct {
	Frequency   Frequency     `json:"frequency"`
	Interval    *int          `json:"interval,omitempty"`
	Weekdays    []string      `json:"weekdays,omitempty"`
	MonthlyMode *string       `json:"monthly_mode,omitempty"`
	DayOfMonth  *int          `json:"day_of_month,omitempty"`
	Date        *string       `json:"date,omitempty"`
	LocalTime   *string       `json:"local_time,omitempty"`
	Timezone    *string       `json:"timezone,omitempty"`
	StartsAt    *string       `json:"starts_at,omitempty"`
	Ends        *EndCondition `json:"ends,omitempty"`
}

type Schedule struct {
	ID                     string                 `json:"id"`
	Status                 ScheduleStatus         `json:"status"`
	ExternalUserID         string                 `json:"external_user_id,omitempty"`
	ExternalOrganizationID string                 `json:"external_organization_id,omitempty"`
	Action                 ActionSummary          `json:"action"`
	Payload                map[string]interface{} `json:"payload,omitempty"`
	Recurrence             Recurrence             `json:"recurrence"`
	Timezone               string                 `json:"timezone"`
	OverlapPolicy          OverlapPolicy          `json:"overlap_policy,omitempty"`
	MisfirePolicy          MisfirePolicy          `json:"misfire_policy,omitempty"`
	NextRunAt              *time.Time             `json:"next_run_at,omitempty"`
	LastRunAt              *time.Time             `json:"last_run_at,omitempty"`
	Revision               int                    `json:"revision,omitempty"`
	InsertedAt             *time.Time             `json:"inserted_at,omitempty"`
	UpdatedAt              *time.Time             `json:"updated_at,omitempty"`
}

type CreateScheduleParams struct {
	ExternalUserID         string                 `json:"external_user_id,omitempty"`
	ExternalOrganizationID string                 `json:"external_organization_id,omitempty"`
	Action                 string                 `json:"action"`
	Payload                map[string]interface{} `json:"payload,omitempty"`
	Recurrence             Recurrence             `json:"recurrence"`
	Timezone               string                 `json:"timezone"`
	OverlapPolicy          OverlapPolicy          `json:"overlap_policy,omitempty"`
	MisfirePolicy          MisfirePolicy          `json:"misfire_policy,omitempty"`
}

type UpdateScheduleParams struct {
	ExternalUserID         *string                 `json:"external_user_id,omitempty"`
	ExternalOrganizationID *string                 `json:"external_organization_id,omitempty"`
	Action                 *string                 `json:"action,omitempty"`
	Payload                *map[string]interface{} `json:"payload,omitempty"`
	Recurrence             *Recurrence             `json:"recurrence,omitempty"`
	Timezone               *string                 `json:"timezone,omitempty"`
	OverlapPolicy          *OverlapPolicy          `json:"overlap_policy,omitempty"`
	MisfirePolicy          *MisfirePolicy          `json:"misfire_policy,omitempty"`
}

type ListSchedulesParams struct {
	ExternalUserID string
	Action         string
	Status         string
	Limit          int
	Cursor         string
}

// ListParams is retained as a source-compatible alias.
type ListParams = ListSchedulesParams

type RunStatus string

const (
	RunStatusScheduled  RunStatus = "scheduled"
	RunStatusDelivering RunStatus = "delivering"
	RunStatusDelivered  RunStatus = "delivered"
	RunStatusRetrying   RunStatus = "retrying"
	RunStatusSucceeded  RunStatus = "succeeded"
	RunStatusFailed     RunStatus = "failed"
	RunStatusCancelled  RunStatus = "cancelled"
	RunStatusSkipped    RunStatus = "skipped"
)

type Run struct {
	ID               string                 `json:"id"`
	ExternalUserID   string                 `json:"external_user_id,omitempty"`
	ScheduleID       string                 `json:"schedule_id"`
	Action           ActionSummary          `json:"action"`
	Status           RunStatus              `json:"status"`
	NominalAt        *time.Time             `json:"nominal_at,omitempty"`
	ScheduledAt      *time.Time             `json:"scheduled_at,omitempty"`
	StartedAt        *time.Time             `json:"started_at,omitempty"`
	CompletedAt      *time.Time             `json:"completed_at,omitempty"`
	AttemptCount     int                    `json:"attempt_count,omitempty"`
	PayloadSnapshot  map[string]interface{} `json:"payload_snapshot,omitempty"`
	Error            *string                `json:"error,omitempty"`
	IdempotencyKey   string                 `json:"idempotency_key,omitempty"`
	ScheduleRevision int                    `json:"schedule_revision,omitempty"`
	InsertedAt       *time.Time             `json:"inserted_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
}

type ActionSummary struct {
	Identifier  string `json:"identifier"`
	DisplayName string `json:"display_name"`
}
