package client

import "time"

type SchedulePreviewParams struct {
	Recurrence Recurrence `json:"recurrence"`
	Count      int        `json:"count,omitempty"`
	From       *time.Time `json:"from,omitempty"`
}

type SchedulePreviewOccurrence struct {
	UTC                  time.Time `json:"utc"`
	Local                string    `json:"local"`
	Date                 string    `json:"date"`
	Time                 string    `json:"time"`
	Timezone             string    `json:"timezone"`
	TimezoneAbbreviation string    `json:"timezone_abbreviation"`
	Label                string    `json:"label"`
}

type SchedulePreview struct {
	Summary     string                      `json:"summary"`
	Occurrences []SchedulePreviewOccurrence `json:"occurrences"`
}

func (s *SchedulesService) Preview(params SchedulePreviewParams) (*SchedulePreview, error) {
	var response struct {
		Data SchedulePreview `json:"data"`
	}

	if err := s.client.request("POST", "/schedules/preview", params, nil, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
