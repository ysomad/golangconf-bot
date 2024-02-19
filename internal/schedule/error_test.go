package schedule

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newCSVError(t *testing.T) {
	type args struct {
		cause error
		row   int
		col   int
	}
	tests := map[string]struct {
		args args
		want CSVError
	}{
		"ok": {
			args: args{
				cause: errors.New("test message"),
				row:   33,
				col:   55,
			},
			want: CSVError{
				cause: errors.New("test message"),
				row:   34,
				col:   56,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, newCSVError(tt.args.cause, tt.args.row, tt.args.col))
		})
	}
}

func TestCSVError_Error(t *testing.T) {
	type fields struct {
		cause error
		row   int
		col   int
	}
	tests := map[string]struct {
		fields fields
		want   string
	}{
		"ok": {
			fields: fields{
				cause: errors.New("CAUSE ERROR"),
				row:   123,
				col:   456,
			},
			want: "[ряд 123, столбец 456] CAUSE ERROR",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := CSVError{
				cause: tt.fields.cause,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}
