package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func date(d int, m time.Month, y int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestConference_Validate(t *testing.T) {
	type fields struct {
		DateFrom           time.Time
		DateUntil          time.Time
		EvaluationDeadline time.Time
	}
	tests := map[string]struct {
		fields  fields
		wantErr error
	}{
		"valid": {
			fields: fields{
				DateFrom:           date(24, time.June, 2024),
				DateUntil:          date(25, time.June, 2024),
				EvaluationDeadline: date(10, time.July, 2024),
			},
			wantErr: nil,
		},
		"date_until_before_date_from": {
			fields: fields{
				DateFrom:           date(24, time.June, 2024),
				DateUntil:          date(17, time.June, 2024),
				EvaluationDeadline: date(10, time.July, 2024),
			},
			wantErr: errInvalidDateUntil,
		},
		"date_until_equal_date_from": {
			fields: fields{
				DateFrom:           date(17, time.June, 2024),
				DateUntil:          date(17, time.June, 2024),
				EvaluationDeadline: date(10, time.July, 2024),
			},
			wantErr: errInvalidDateUntil,
		},
		"evaluation_deadline_before_date_until": {
			fields: fields{
				DateFrom:           date(17, time.June, 2024),
				DateUntil:          date(20, time.June, 2024),
				EvaluationDeadline: date(10, time.June, 2024),
			},
			wantErr: errInvalidEvaluationDeadline,
		},
		"evaluation_deadline_equal_date_until": {
			fields: fields{
				DateFrom:           date(17, time.June, 2024),
				DateUntil:          date(20, time.June, 2024),
				EvaluationDeadline: date(10, time.June, 2024),
			},
			wantErr: errInvalidEvaluationDeadline,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := Conference{
				DateFrom:           tt.fields.DateFrom,
				DateUntil:          tt.fields.DateUntil,
				EvaluationDeadline: tt.fields.EvaluationDeadline,
			}
			err := c.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
