package client

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedulesService_Preview(t *testing.T) {
	localTime := "09:00"
	timezone := "Australia/Melbourne"

	client, server := setupTestClient(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/schedules/preview", r.URL.Path)

		var params SchedulePreviewParams
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&params))
		assert.Equal(t, FrequencyWeekly, params.Recurrence.Frequency)
		assert.Equal(t, 3, params.Count)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": SchedulePreview{
				Summary: "Every Monday at 09:00",
				Occurrences: []SchedulePreviewOccurrence{
					{Date: "2026-08-17", Time: "09:00", Timezone: timezone, Label: "Mon, 17 Aug 2026 at 09:00"},
				},
			},
		})
	})
	defer server.Close()

	preview, err := client.Schedules.Preview(SchedulePreviewParams{
		Recurrence: Recurrence{
			Frequency: FrequencyWeekly,
			Weekdays:  []string{"monday"},
			LocalTime: &localTime,
			Timezone:  &timezone,
		},
		Count: 3,
	})

	assert.NoError(t, err)
	assert.Equal(t, "Every Monday at 09:00", preview.Summary)
	assert.Len(t, preview.Occurrences, 1)
}
